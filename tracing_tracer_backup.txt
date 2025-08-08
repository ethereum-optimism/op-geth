package tracing

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	tracer     oteltrace.Tracer
	tracerOnce sync.Once
	isEnabled  bool
	config     *Config
)

// Initialize sets up the global tracer with the given configuration
func Initialize(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid tracing config: %w", err)
	}

	config = cfg
	isEnabled = cfg.Enabled

	if !cfg.Enabled {
		// Set up a no-op tracer
		otel.SetTracerProvider(oteltrace.NewNoopTracerProvider())
		tracer = otel.Tracer("geth-noop")
		return nil
	}

	return initializeTracer(cfg)
}

func initializeTracer(cfg *Config) error {
	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("service.version", cfg.ServiceVersion),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(cfg.Endpoint),
		otlptracehttp.WithHeaders(cfg.Headers),
		otlptracehttp.WithTimeout(cfg.Timeout),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.TraceIDRatioBased(cfg.SampleRate)),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Get tracer
	tracer = otel.Tracer("geth")

	return nil
}

// IsEnabled returns true if tracing is enabled
func IsEnabled() bool {
	return isEnabled
}

// GetTracer returns the global tracer instance
func GetTracer() oteltrace.Tracer {
	tracerOnce.Do(func() {
		if tracer == nil {
			// Initialize with default config if not already initialized
			_ = Initialize(nil)
		}
	})
	return tracer
}

// StartSpan starts a new span with the given name and options
func StartSpan(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	if !IsEnabled() {
		return ctx, oteltrace.SpanFromContext(ctx)
	}
	return GetTracer().Start(ctx, spanName, opts...)
}

// GetConfig returns the current tracing configuration
func GetConfig() *Config {
	if config == nil {
		return DefaultConfig()
	}
	return config
}

// Shutdown gracefully shuts down the tracer
func Shutdown(ctx context.Context) error {
	if !IsEnabled() {
		return nil
	}

	if tp, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
		return tp.Shutdown(ctx)
	}
	return nil
}

// ExtractTraceContext extracts trace context from environment variables
func ExtractTraceContext() (*Config, error) {
	cfg := DefaultConfig()

	// Check environment variables
	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		cfg.Endpoint = endpoint
		cfg.Enabled = true
	}

	if serviceName := os.Getenv("OTEL_SERVICE_NAME"); serviceName != "" {
		cfg.ServiceName = serviceName
	}

	if serviceVersion := os.Getenv("OTEL_SERVICE_VERSION"); serviceVersion != "" {
		cfg.ServiceVersion = serviceVersion
	}

	// Parse headers from OTEL_EXPORTER_OTLP_HEADERS
	if headers := os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"); headers != "" {
		cfg.Headers = parseHeaders(headers)
	}

	return cfg, nil
}

// parseHeaders parses comma-separated key=value pairs
func parseHeaders(headerStr string) map[string]string {
	headers := make(map[string]string)
	// Simple parsing - in production, you might want more robust parsing
	// Format: "key1=value1,key2=value2"
	return headers
}