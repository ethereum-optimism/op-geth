package utils

import (
	"time"

	"github.com/ethereum/go-ethereum/tracing"
	"github.com/urfave/cli/v2"
)

var (
	// Tracing flags
	TracingEnabledFlag = &cli.BoolFlag{
		Name:  "tracing.enabled",
		Usage: "Enable OpenTelemetry tracing",
	}
	TracingEndpointFlag = &cli.StringFlag{
		Name:  "tracing.endpoint",
		Usage: "OpenTelemetry OTLP endpoint URL",
		Value: "",
	}
	TracingServiceNameFlag = &cli.StringFlag{
		Name:  "tracing.service-name",
		Usage: "Service name for tracing",
		Value: "geth",
	}
	TracingServiceVersionFlag = &cli.StringFlag{
		Name:  "tracing.service-version",
		Usage: "Service version for tracing",
		Value: "",
	}
	TracingSampleRateFlag = &cli.Float64Flag{
		Name:  "tracing.sample-rate",
		Usage: "Tracing sample rate (0.0 to 1.0)",
		Value: 1.0,
	}
	TracingTimeoutFlag = &cli.DurationFlag{
		Name:  "tracing.timeout",
		Usage: "Tracing export timeout",
		Value: 10 * time.Second,
	}
	TracingRPCFlag = &cli.BoolFlag{
		Name:  "tracing.rpc",
		Usage: "Enable RPC tracing",
		Value: true,
	}
	TracingEngineAPIFlag = &cli.BoolFlag{
		Name:  "tracing.engine-api",
		Usage: "Enable Engine API tracing",
		Value: true,
	}
	TracingTransactionFlag = &cli.BoolFlag{
		Name:  "tracing.transactions",
		Usage: "Enable transaction-level tracing",
		Value: false,
	}
)

// TracingFlags contains all tracing-related flags
var TracingFlags = []cli.Flag{
	TracingEnabledFlag,
	TracingEndpointFlag,
	TracingServiceNameFlag,
	TracingServiceVersionFlag,
	TracingSampleRateFlag,
	TracingTimeoutFlag,
	TracingRPCFlag,
	TracingEngineAPIFlag,
	TracingTransactionFlag,
}

// SetTracing configures tracing from CLI context
func SetTracing(ctx *cli.Context) *tracing.Config {
	cfg := tracing.DefaultConfig()

	if ctx.IsSet(TracingEnabledFlag.Name) {
		cfg.Enabled = ctx.Bool(TracingEnabledFlag.Name)
	}

	if ctx.IsSet(TracingEndpointFlag.Name) {
		cfg.Endpoint = ctx.String(TracingEndpointFlag.Name)
		if cfg.Endpoint != "" {
			cfg.Enabled = true // Auto-enable if endpoint is provided
		}
	}

	if ctx.IsSet(TracingServiceNameFlag.Name) {
		cfg.ServiceName = ctx.String(TracingServiceNameFlag.Name)
	}

	if ctx.IsSet(TracingServiceVersionFlag.Name) {
		cfg.ServiceVersion = ctx.String(TracingServiceVersionFlag.Name)
	}

	if ctx.IsSet(TracingSampleRateFlag.Name) {
		cfg.SampleRate = ctx.Float64(TracingSampleRateFlag.Name)
	}

	if ctx.IsSet(TracingTimeoutFlag.Name) {
		cfg.Timeout = ctx.Duration(TracingTimeoutFlag.Name)
	}

	if ctx.IsSet(TracingRPCFlag.Name) {
		cfg.EnableRPCTracing = ctx.Bool(TracingRPCFlag.Name)
	}

	if ctx.IsSet(TracingEngineAPIFlag.Name) {
		cfg.EnableEngineAPITracing = ctx.Bool(TracingEngineAPIFlag.Name)
	}

	if ctx.IsSet(TracingTransactionFlag.Name) {
		cfg.EnableTransactionTracing = ctx.Bool(TracingTransactionFlag.Name)
	}

	return cfg
}