package pid

import (
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// FastPIDController manages execution layer gas utilization with real-time PID control
type FastPIDController struct {
	mu sync.RWMutex

	// PID Parameters (configurable via RPC)
	Kp, Ki, Kd float64
	setpoint   float64 // Target utilization (1.0 = 100% of gas target)
	maxChange  float64 // Maximum fee change per block

	// Gas configuration
	gasTarget uint64
	gasLimit  uint64

	// PID state variables
	previousError float64
	integralSum   float64
	lastTime      time.Time

	// Safety limits
	minBaseFee *big.Int
	maxBaseFee *big.Int

	// External pressure signals (from batcher)
	daPressure  float64 // DA pressure from batcher PID
	l1Influence float64 // L1 economics influence

	// Metrics and logging
	logger log.Logger

	// Enable/disable PID control
	enabled bool
}

// NewFastPIDController creates a new fast PID controller for sequencer
func NewFastPIDController(gasTarget, gasLimit uint64, logger log.Logger) *FastPIDController {
	if logger == nil {
		logger = log.New("module", "fast-pid")
	}

	return &FastPIDController{
		// Conservative initial parameters
		Kp:        1.8,  // Proportional gain - responsive but stable
		Ki:        0.12, // Integral gain - eliminate steady-state errors
		Kd:        0.35, // Derivative gain - dampen oscillations
		setpoint:  1.0,  // 100% of gas target
		maxChange: 0.12, // ±12% max change per block (vs ±2% in current EIP-1559)

		gasTarget: gasTarget,
		gasLimit:  gasLimit,

		lastTime:   time.Now(),
		minBaseFee: big.NewInt(1000000),    // 0.001 gwei minimum
		maxBaseFee: big.NewInt(1000000000), // 1 gwei maximum (reasonable for L2)

		enabled: true,

		logger: logger,
	}
}

// CalculateBaseFee computes new base fee using PID control
func (f *FastPIDController) CalculateBaseFee(gasUsed uint64, currentBaseFee *big.Int, parentHeader *types.Header) *big.Int {
	if !f.enabled {
		// Fall back to EIP-1559 when disabled
		return f.calculateEIP1559BaseFee(parentHeader)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	deltaTime := now.Sub(f.lastTime).Seconds()
	if deltaTime <= 0 || deltaTime > 10 { // Sanity check for time
		deltaTime = 2.0 // Default 2s block time for Base
	}

	// Calculate utilization error
	utilization := float64(gasUsed) / float64(f.gasTarget)
	error := f.setpoint - utilization

	// PID calculation
	proportional := f.Kp * error

	// Integral with windup protection
	f.integralSum += error * deltaTime
	maxIntegral := f.maxChange / math.Max(f.Ki, 0.001) // Prevent divide by zero
	if f.integralSum > maxIntegral {
		f.integralSum = maxIntegral
	} else if f.integralSum < -maxIntegral {
		f.integralSum = -maxIntegral
	}
	integral := f.Ki * f.integralSum

	// Derivative with basic filtering
	derivativeError := error - f.previousError
	if math.Abs(derivativeError) > 2.0 { // Filter noise spikes
		// 2.0 is an extremely unlikely value, we would want to make it configurable in the future
		derivativeError = f.previousError
	}
	derivative := f.Kd * derivativeError / deltaTime

	// Combined PID output
	pidOutput := proportional + integral + derivative

	// Apply external pressure signals (will be zero initially)
	totalOutput := pidOutput + f.daPressure + f.l1Influence

	// Apply safety limits
	if totalOutput > f.maxChange {
		totalOutput = f.maxChange
	} else if totalOutput < -f.maxChange {
		totalOutput = -f.maxChange
	}

	// Calculate new base fee
	newBaseFee := f.applyPIDAdjustment(currentBaseFee, totalOutput)

	// Apply absolute limits
	if newBaseFee.Cmp(f.minBaseFee) < 0 {
		newBaseFee.Set(f.minBaseFee)
	}
	// TODO: I think minBaseFee makes sense, but maxBaseFee is a bit arbitrary
	// It's more of a safety limit, so that the fees don't get too high
	if newBaseFee.Cmp(f.maxBaseFee) > 0 {
		newBaseFee.Set(f.maxBaseFee)
	}

	// Update state
	f.previousError = error
	f.lastTime = now

	f.logger.Debug("PID calculation",
		"gasUsed", gasUsed,
		"gasTarget", f.gasTarget,
		"utilization", utilization,
		"error", error,
		"P", proportional,
		"I", integral,
		"D", derivative,
		"pidOutput", pidOutput,
		"daPressure", f.daPressure,
		"l1Influence", f.l1Influence,
		"totalOutput", totalOutput,
		"oldBaseFee", currentBaseFee,
		"newBaseFee", newBaseFee,
	)

	return newBaseFee
}

// applyPIDAdjustment applies PID output to current base fee
func (f *FastPIDController) applyPIDAdjustment(currentBaseFee *big.Int, adjustment float64) *big.Int {
	currentFeeFloat := new(big.Float).SetInt(currentBaseFee)
	adjustmentFactor := big.NewFloat(1.0 + adjustment)
	newFeeFloat := new(big.Float).Mul(currentFeeFloat, adjustmentFactor)
	newBaseFee, _ := newFeeFloat.Int(nil)
	return newBaseFee
}

// calculateEIP1559BaseFee computes EIP-1559 base fee for hybrid/fallback mode
func (f *FastPIDController) calculateEIP1559BaseFee(parentHeader *types.Header) *big.Int {
	// This would call the existing EIP-1559 calculation
	// For now, return parent base fee (will be replaced with actual EIP-1559 calc)
	if parentHeader.BaseFee != nil {
		return new(big.Int).Set(parentHeader.BaseFee)
	}
	return f.minBaseFee
}

// UpdateExternalPressure receives pressure signals from batcher (via RPC)
func (f *FastPIDController) UpdateExternalPressure(daPressure, l1Influence float64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.daPressure = daPressure
	f.l1Influence = l1Influence

	f.logger.Debug("Updated external pressure",
		"daPressure", daPressure,
		"l1Influence", l1Influence,
	)
}

// UpdateParameters allows runtime parameter adjustment
func (f *FastPIDController) UpdateParameters(kp, ki, kd, maxChange float64) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Validate parameters
	if kp < 0 || kp > 10 || ki < 0 || ki > 1 || kd < 0 || kd > 2 || maxChange <= 0 || maxChange > 0.5 {
		f.logger.Error("Invalid PID parameters", "kp", kp, "ki", ki, "kd", kd, "maxChange", maxChange)
		return
	}

	f.Kp = kp
	f.Ki = ki
	f.Kd = kd
	f.maxChange = maxChange

	f.logger.Info("Updated PID parameters",
		"Kp", kp, "Ki", ki, "Kd", kd, "maxChange", maxChange)
}

// SetEnabled enables or disables PID control
func (f *FastPIDController) SetEnabled(enabled bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.enabled = enabled
	f.logger.Info("PID controller enabled", "enabled", enabled)
}

func (f *FastPIDController) IsEnabled() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.enabled
}

// ValidateConfiguration validates all PID controller parameters and configuration
func (f *FastPIDController) ValidateConfiguration() error {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Validate PID parameters
	if f.Kp < 0 || f.Kp > 10 {
		return fmt.Errorf("invalid Kp parameter: %f, must be between 0 and 10", f.Kp)
	}
	if f.Ki < 0 || f.Ki > 1 {
		return fmt.Errorf("invalid Ki parameter: %f, must be between 0 and 1", f.Ki)
	}
	if f.Kd < 0 || f.Kd > 2 {
		return fmt.Errorf("invalid Kd parameter: %f, must be between 0 and 2", f.Kd)
	}

	// Validate setpoint
	if f.setpoint <= 0 || f.setpoint > 2 {
		return fmt.Errorf("invalid setpoint: %f, must be between 0 and 2", f.setpoint)
	}

	// Validate maxChange
	if f.maxChange <= 0 || f.maxChange > 0.5 {
		return fmt.Errorf("invalid maxChange: %f, must be between 0 and 0.5", f.maxChange)
	}

	// Validate gas configuration
	if f.gasTarget == 0 {
		return fmt.Errorf("gasTarget cannot be zero")
	}
	if f.gasLimit == 0 {
		return fmt.Errorf("gasLimit cannot be zero")
	}
	if f.gasTarget > f.gasLimit {
		return fmt.Errorf("gasTarget (%d) cannot exceed gasLimit (%d)", f.gasTarget, f.gasLimit)
	}

	// Validate base fee limits
	if f.minBaseFee == nil || f.minBaseFee.Sign() <= 0 {
		return fmt.Errorf("minBaseFee must be positive")
	}
	if f.maxBaseFee == nil || f.maxBaseFee.Sign() <= 0 {
		return fmt.Errorf("maxBaseFee must be positive")
	}
	if f.minBaseFee.Cmp(f.maxBaseFee) >= 0 {
		return fmt.Errorf("minBaseFee (%s) must be less than maxBaseFee (%s)", f.minBaseFee.String(), f.maxBaseFee.String())
	}

	// Validate external pressure signals (should be reasonable bounds)
	if math.Abs(f.daPressure) > 1.0 {
		return fmt.Errorf("daPressure out of bounds: %f, should be between -1.0 and 1.0", f.daPressure)
	}
	if math.Abs(f.l1Influence) > 1.0 {
		return fmt.Errorf("l1Influence out of bounds: %f, should be between -1.0 and 1.0", f.l1Influence)
	}

	return nil
}

// GetStatus returns current PID controller status
func (f *FastPIDController) GetStatus() map[string]interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return map[string]interface{}{
		"enabled":       f.enabled,
		"kp":            f.Kp,
		"ki":            f.Ki,
		"kd":            f.Kd,
		"maxChange":     f.maxChange,
		"gasTarget":     f.gasTarget,
		"gasLimit":      f.gasLimit,
		"integralSum":   f.integralSum,
		"previousError": f.previousError,
		"daPressure":    f.daPressure,
		"l1Influence":   f.l1Influence,
	}
}
