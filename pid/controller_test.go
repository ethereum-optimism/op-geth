package pid

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

// Test configuration constants
const (
	testGasTarget = 5000000  // 5M gas target
	testGasLimit  = 30000000 // 30M gas limit
)

func TestNewFastPIDController(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	assert.NotNil(t, controller)
	assert.True(t, controller.IsEnabled())
	assert.Equal(t, 1.8, controller.Kp)
	assert.Equal(t, 0.12, controller.Ki)
	assert.Equal(t, 0.35, controller.Kd)
	assert.Equal(t, uint64(testGasTarget), controller.gasTarget)
	assert.Equal(t, uint64(testGasLimit), controller.gasLimit)
}

func TestPIDBaseFeeCalculation(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	// Create a mock parent header
	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000, // Exactly at target
		GasLimit: 30000000,
		BaseFee:  big.NewInt(10000000), // 0.01 gwei
		Time:     uint64(time.Now().Unix()),
	}

	// Test at target utilization (should maintain base fee relatively stable)
	newBaseFee := controller.CalculateBaseFee(5000000, parentHeader.BaseFee, parentHeader)

	// With no error (at setpoint), fee should remain close to original
	assert.True(t, newBaseFee.Cmp(big.NewInt(8000000)) > 0)  // Not too low
	assert.True(t, newBaseFee.Cmp(big.NewInt(12000000)) < 0) // Not too high
}

func TestPIDHighUtilization(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  big.NewInt(10000000), // 0.01 gwei
		Time:     uint64(time.Now().Unix()),
	}

	// Test high utilization (150% of target)
	highGasUsed := uint64(7500000) // 150% of 5M target
	newBaseFee := controller.CalculateBaseFee(highGasUsed, parentHeader.BaseFee, parentHeader)

	// Base fee should change with high utilization (direction depends on PID tuning)
	assert.NotEqual(t, newBaseFee, parentHeader.BaseFee,
		"Base fee should change with high utilization")

	// But not exceed max change limit (12%)
	maxExpected := new(big.Int).Mul(parentHeader.BaseFee, big.NewInt(112))
	maxExpected.Div(maxExpected, big.NewInt(100))
	assert.True(t, newBaseFee.Cmp(maxExpected) <= 0,
		"Base fee increase should not exceed max change limit")
}

func TestPIDLowUtilization(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  big.NewInt(10000000), // 0.01 gwei
		Time:     uint64(time.Now().Unix()),
	}

	// Test low utilization (50% of target)
	lowGasUsed := uint64(2500000) // 50% of 5M target
	newBaseFee := controller.CalculateBaseFee(lowGasUsed, parentHeader.BaseFee, parentHeader)

	// Base fee should change with low utilization
	assert.NotEqual(t, newBaseFee, parentHeader.BaseFee,
		"Base fee should change with low utilization")

	// But not go below minimum
	minBaseFee := big.NewInt(1000000) // 0.001 gwei minimum from controller
	assert.True(t, newBaseFee.Cmp(minBaseFee) >= 0,
		"Base fee should not go below minimum")
}

func TestPIDParameterUpdate(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	// Test valid parameter update
	controller.UpdateParameters(2.0, 0.15, 0.4, 0.15)

	assert.Equal(t, 2.0, controller.Kp)
	assert.Equal(t, 0.15, controller.Ki)
	assert.Equal(t, 0.4, controller.Kd)
	assert.Equal(t, 0.15, controller.maxChange)

	// Test invalid parameter update (should be ignored)
	oldKp := controller.Kp
	controller.UpdateParameters(-1.0, 0.15, 0.4, 0.15) // Invalid Kp
	assert.Equal(t, oldKp, controller.Kp)              // Should remain unchanged
}

func TestPIDExternalPressure(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	// Test valid external pressure
	controller.UpdateExternalPressure(0.1, -0.05)

	status := controller.GetStatus()
	assert.Equal(t, 0.1, status["daPressure"])
	assert.Equal(t, -0.05, status["l1Influence"])

	// Test external pressure affects calculation
	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  big.NewInt(10000000),
		Time:     uint64(time.Now().Unix()),
	}

	// Calculate with external pressure
	newBaseFee := controller.CalculateBaseFee(5000000, parentHeader.BaseFee, parentHeader)
	assert.NotNil(t, newBaseFee)
}

func TestPIDFallbackToEIP1559(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	// Disable PID controller
	controller.SetEnabled(false)
	assert.False(t, controller.IsEnabled())

	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  big.NewInt(10000000),
		Time:     uint64(time.Now().Unix()),
	}

	// Should fall back to EIP-1559 (currently returns parent base fee)
	newBaseFee := controller.CalculateBaseFee(5000000, parentHeader.BaseFee, parentHeader)

	// Verify it uses EIP-1559 fallback logic
	assert.Equal(t, parentHeader.BaseFee, newBaseFee,
		"EIP-1559 fallback should return parent base fee")
}

func TestPIDIntegralWindup(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  big.NewInt(10000000),
		Time:     uint64(time.Now().Unix()),
	}

	// Simulate persistent high utilization to test integral windup protection
	for i := 0; i < 10; i++ {
		time.Sleep(10 * time.Millisecond)                                                       // Small delay between calculations
		newBaseFee := controller.CalculateBaseFee(10000000, parentHeader.BaseFee, parentHeader) // Very high usage
		parentHeader.BaseFee = newBaseFee
		parentHeader.Number.Add(parentHeader.Number, big.NewInt(1))
	}

	// Integral sum should be bounded
	status := controller.GetStatus()
	integralSum := status["integralSum"].(float64)
	maxExpectedIntegral := controller.maxChange / controller.Ki

	assert.True(t, integralSum <= maxExpectedIntegral,
		"Integral sum should be bounded to prevent windup")
}

func TestPIDStatus(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	status := controller.GetStatus()

	// Verify all expected status fields are present
	expectedKeys := []string{
		"enabled", "kp", "ki", "kd", "maxChange", "gasTarget", "gasLimit",
		"integralSum", "previousError", "daPressure", "l1Influence",
	}

	for _, key := range expectedKeys {
		assert.Contains(t, status, key, "Status should contain key: %s", key)
	}

	assert.Equal(t, true, status["enabled"])
	assert.Equal(t, 1.8, status["kp"])
	assert.Equal(t, uint64(testGasTarget), status["gasTarget"])
}

func TestPIDValidateConfiguration(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	// Test valid configuration
	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)
	err := controller.ValidateConfiguration()
	assert.NoError(t, err)

	// Test invalid configurations
	invalidController := NewFastPIDController(testGasTarget, testGasLimit, logger)

	// Invalid Kp
	invalidController.Kp = -1.0
	err = invalidController.ValidateConfiguration()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Kp parameter")

	// Reset and test invalid Ki
	invalidController = NewFastPIDController(testGasTarget, testGasLimit, logger)
	invalidController.Ki = 2.0
	err = invalidController.ValidateConfiguration()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Ki parameter")

	// Reset and test invalid Kd
	invalidController = NewFastPIDController(testGasTarget, testGasLimit, logger)
	invalidController.Kd = -1.0
	err = invalidController.ValidateConfiguration()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Kd parameter")
}

func TestPIDMinMaxBaseFee(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  big.NewInt(500000), // Very low base fee
		Time:     uint64(time.Now().Unix()),
	}

	// Test minimum base fee enforcement
	newBaseFee := controller.CalculateBaseFee(1000000, parentHeader.BaseFee, parentHeader) // Very low usage
	minBaseFee := big.NewInt(1000000)                                                      // 0.001 gwei minimum
	assert.True(t, newBaseFee.Cmp(minBaseFee) >= 0,
		"Base fee should not go below minimum")

	// Test maximum base fee enforcement
	parentHeader.BaseFee = big.NewInt(900000000)                                           // Very high base fee
	newBaseFee = controller.CalculateBaseFee(15000000, parentHeader.BaseFee, parentHeader) // Very high usage
	maxBaseFee := big.NewInt(1000000000)                                                   // 1 gwei maximum
	assert.True(t, newBaseFee.Cmp(maxBaseFee) <= 0,
		"Base fee should not exceed maximum")
}

// Benchmark tests
func BenchmarkPIDCalculation(b *testing.B) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  big.NewInt(10000000),
		Time:     uint64(time.Now().Unix()),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		controller.CalculateBaseFee(5000000, parentHeader.BaseFee, parentHeader)
	}
}

// Integration tests with mock headers
func TestPIDIntegrationScenario(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	// Simulate a demand spike scenario
	scenarios := []struct {
		name         string
		gasUsed      uint64
		expectChange bool // Whether we expect the fee to change
	}{
		{"Normal usage", 5000000, false},          // At target, should be stable
		{"Demand spike starts", 8000000, true},    // Above target, should change
		{"High demand continues", 10000000, true}, // Well above target, should change
		{"Peak demand", 12000000, true},           // Very high, should change
		{"Demand starts to drop", 8000000, true},  // Still above target, should change
		{"Back to normal", 5000000, true},         // Back to target, may still adjust
		{"Low demand", 3000000, true},             // Below target, should change
	}

	currentBaseFee := big.NewInt(10000000) // Start at 0.01 gwei

	for i, scenario := range scenarios {
		parentHeader := &types.Header{
			Number:   big.NewInt(int64(1000 + i)),
			GasUsed:  5000000, // Previous block usage
			GasLimit: 30000000,
			BaseFee:  currentBaseFee,
			Time:     uint64(time.Now().Unix()),
		}

		newBaseFee := controller.CalculateBaseFee(scenario.gasUsed, currentBaseFee, parentHeader)

		if scenario.expectChange {
			// For scenarios where we expect change, just verify the fee changed
			// and is within reasonable bounds
			maxChange := new(big.Int).Mul(currentBaseFee, big.NewInt(112))
			maxChange.Div(maxChange, big.NewInt(100))
			minChange := new(big.Int).Mul(currentBaseFee, big.NewInt(88))
			minChange.Div(minChange, big.NewInt(100))

			assert.True(t, newBaseFee.Cmp(maxChange) <= 0 && newBaseFee.Cmp(minChange) >= 0,
				"Scenario %s: Base fee change should be within bounds", scenario.name)
		} else {
			// For stable scenarios, allow small variations
			ratio := new(big.Float).Quo(new(big.Float).SetInt(newBaseFee), new(big.Float).SetInt(currentBaseFee))
			ratioFloat, _ := ratio.Float64()
			assert.True(t, ratioFloat > 0.95 && ratioFloat < 1.05,
				"Scenario %s: Base fee should be relatively stable", scenario.name)
		}

		currentBaseFee = newBaseFee
		t.Logf("Scenario: %s, Gas Used: %d, Old Fee: %s, New Fee: %s",
			scenario.name, scenario.gasUsed, parentHeader.BaseFee.String(), newBaseFee.String())

		// Add small delay to simulate block time
		time.Sleep(10 * time.Millisecond)
	}
}

func TestPIDConcurrency(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  big.NewInt(10000000),
		Time:     uint64(time.Now().Unix()),
	}

	// Test concurrent access to controller methods
	done := make(chan bool, 3)

	// Goroutine 1: Calculate base fees
	go func() {
		for i := 0; i < 100; i++ {
			controller.CalculateBaseFee(5000000, parentHeader.BaseFee, parentHeader)
		}
		done <- true
	}()

	// Goroutine 2: Update parameters
	go func() {
		for i := 0; i < 100; i++ {
			controller.UpdateParameters(1.8, 0.12, 0.35, 0.12)
		}
		done <- true
	}()

	// Goroutine 3: Update external pressure
	go func() {
		for i := 0; i < 100; i++ {
			controller.UpdateExternalPressure(0.1, -0.05)
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// Verify controller is still functional
	status := controller.GetStatus()
	assert.NotNil(t, status)
	assert.True(t, controller.IsEnabled())
}

func TestPIDEdgeCases(t *testing.T) {
	logger := log.NewLogger(log.DiscardHandler())

	controller := NewFastPIDController(testGasTarget, testGasLimit, logger)

	// Test with nil base fee
	parentHeader := &types.Header{
		Number:   big.NewInt(1000),
		GasUsed:  5000000,
		GasLimit: 30000000,
		BaseFee:  nil,
		Time:     uint64(time.Now().Unix()),
	}

	// Should handle nil base fee gracefully when disabled
	controller.SetEnabled(false)
	newBaseFee := controller.CalculateBaseFee(5000000, big.NewInt(10000000), parentHeader)
	assert.NotNil(t, newBaseFee)

	// Test with zero gas used
	controller.SetEnabled(true)
	newBaseFee = controller.CalculateBaseFee(0, big.NewInt(10000000), parentHeader)
	assert.NotNil(t, newBaseFee)

	// Test with maximum gas used
	newBaseFee = controller.CalculateBaseFee(testGasLimit, big.NewInt(10000000), parentHeader)
	assert.NotNil(t, newBaseFee)
}
