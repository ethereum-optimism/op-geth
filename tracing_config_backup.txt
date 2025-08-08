package tracing

import (
	"fmt"
	"time"
)

// Config represents the tracing configuration
type Config struct {
	// Enabled controls whether tracing is active
	Enabled bool

	// Endpoint is the OTLP endpoint URL for trace export
	Endpoint string

	// ServiceName is the service name to use in traces
	ServiceName string

	// ServiceVersion is the service version to include in traces
	ServiceVersion string

	// Headers are additional headers to send with trace exports
	Headers map[string]string

	// Timeout is the timeout for trace exports
	Timeout time.Duration

	// SampleRate is the sampling rate (0.0 to 1.0)
	SampleRate float64

	// EnableRPCTracing controls whether RPC calls are traced
	EnableRPCTracing bool

	// EnableEngineAPITracing controls whether Engine API calls are traced
	EnableEngineAPITracing bool

	// EnableTransactionTracing controls whether individual transactions are traced
	EnableTransactionTracing bool
}

// DefaultConfig returns a default tracing configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled:                  false,
		Endpoint:                 "",
		ServiceName:             "geth",
		ServiceVersion:          "",
		Headers:                 make(map[string]string),
		Timeout:                 10 * time.Second,
		SampleRate:              1.0,
		EnableRPCTracing:        true,
		EnableEngineAPITracing:  true,
		EnableTransactionTracing: false,
	}
}

// Validate validates the tracing configuration
func (c *Config) Validate() error {
	if c.Enabled {
		if c.Endpoint == "" {
			return fmt.Errorf("tracing endpoint is required when tracing is enabled")
		}
		if c.ServiceName == "" {
			return fmt.Errorf("service name is required when tracing is enabled")
		}
		if c.SampleRate < 0.0 || c.SampleRate > 1.0 {
			return fmt.Errorf("sample rate must be between 0.0 and 1.0, got %f", c.SampleRate)
		}
		if c.Timeout <= 0 {
			return fmt.Errorf("timeout must be positive, got %v", c.Timeout)
		}
	}
	return nil
}

// IsRPCTracingEnabled returns true if RPC tracing should be active
func (c *Config) IsRPCTracingEnabled() bool {
	return c.Enabled && c.EnableRPCTracing
}

// IsEngineAPITracingEnabled returns true if Engine API tracing should be active
func (c *Config) IsEngineAPITracingEnabled() bool {
	return c.Enabled && c.EnableEngineAPITracing
}

// IsTransactionTracingEnabled returns true if transaction tracing should be active
func (c *Config) IsTransactionTracingEnabled() bool {
	return c.Enabled && c.EnableTransactionTracing
}