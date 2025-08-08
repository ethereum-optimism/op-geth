package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	// Span names
	SpanNameRPCRequest        = "rpc.request"
	SpanNameSendRawTx         = "eth.sendRawTransaction"
	SpanNameSubmitTransaction = "eth.submitTransaction"
	SpanNameTxPoolAdd         = "txpool.add"
	SpanNameEngineAPI         = "engine.api"
	SpanNameGetPayload        = "engine.getPayload"
	SpanNameNewPayload        = "engine.newPayload"
	SpanNameForkchoiceUpdate  = "engine.forkchoiceUpdate"
	SpanNameTxForward         = "tx.forward"
	SpanNameBlockBuilding     = "block.building"
)

// Attribute keys
const (
	AttrMethod          = "rpc.method"
	AttrTxHash          = "tx.hash"
	AttrTxFrom          = "tx.from"
	AttrTxTo            = "tx.to"
	AttrTxGas           = "tx.gas"
	AttrTxGasPrice      = "tx.gasPrice"
	AttrTxValue         = "tx.value"
	AttrBlockNumber     = "block.number"
	AttrBlockHash       = "block.hash"
	AttrPayloadID       = "payload.id"
	AttrEngineMethod    = "engine.method"
	AttrBackendURL      = "backend.url"
	AttrErrorCode       = "error.code"
	AttrErrorMessage    = "error.message"
	AttrRequestID       = "request.id"
	AttrUserAgent       = "http.user_agent"
	AttrRemoteAddr      = "http.remote_addr"
)

// StartRPCSpan starts a span for an RPC request
func StartRPCSpan(ctx context.Context, method string) (context.Context, oteltrace.Span) {
	ctx, span := StartSpan(ctx, SpanNameRPCRequest,
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(
			attribute.String(AttrMethod, method),
		),
	)
	return ctx, span
}

// StartSendRawTransactionSpan starts a span for sendRawTransaction
func StartSendRawTransactionSpan(ctx context.Context, txHash string) (context.Context, oteltrace.Span) {
	ctx, span := StartSpan(ctx, SpanNameSendRawTx,
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(
			attribute.String(AttrTxHash, txHash),
		),
	)
	return ctx, span
}

// StartTxForwardSpan starts a span for transaction forwarding
func StartTxForwardSpan(ctx context.Context, txHash string, backendURL string) (context.Context, oteltrace.Span) {
	ctx, span := StartSpan(ctx, SpanNameTxForward,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(
			attribute.String(AttrTxHash, txHash),
			attribute.String(AttrBackendURL, backendURL),
		),
	)
	return ctx, span
}

// StartEngineAPISpan starts a span for Engine API calls
func StartEngineAPISpan(ctx context.Context, method string) (context.Context, oteltrace.Span) {
	ctx, span := StartSpan(ctx, SpanNameEngineAPI,
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(
			attribute.String(AttrEngineMethod, method),
		),
	)
	return ctx, span
}

// StartTxPoolSpan starts a span for transaction pool operations
func StartTxPoolSpan(ctx context.Context, txHash string) (context.Context, oteltrace.Span) {
	ctx, span := StartSpan(ctx, SpanNameTxPoolAdd,
		oteltrace.WithAttributes(
			attribute.String(AttrTxHash, txHash),
		),
	)
	return ctx, span
}

// AddTransactionAttributes adds transaction-related attributes to a span
func AddTransactionAttributes(span oteltrace.Span, txHash, from, to string, gas uint64, gasPrice, value string) {
	if !IsEnabled() {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String(AttrTxHash, txHash),
	}
	
	if from != "" {
		attrs = append(attrs, attribute.String(AttrTxFrom, from))
	}
	if to != "" {
		attrs = append(attrs, attribute.String(AttrTxTo, to))
	}
	if gas > 0 {
		attrs = append(attrs, attribute.Int64(AttrTxGas, int64(gas)))
	}
	if gasPrice != "" {
		attrs = append(attrs, attribute.String(AttrTxGasPrice, gasPrice))
	}
	if value != "" {
		attrs = append(attrs, attribute.String(AttrTxValue, value))
	}
	
	span.SetAttributes(attrs...)
}

// AddBlockAttributes adds block-related attributes to a span
func AddBlockAttributes(span oteltrace.Span, blockNumber uint64, blockHash string) {
	if !IsEnabled() {
		return
	}
	
	attrs := []attribute.KeyValue{}
	if blockNumber > 0 {
		attrs = append(attrs, attribute.Int64(AttrBlockNumber, int64(blockNumber)))
	}
	if blockHash != "" {
		attrs = append(attrs, attribute.String(AttrBlockHash, blockHash))
	}
	
	span.SetAttributes(attrs...)
}

// AddHTTPAttributes adds HTTP-related attributes to a span
func AddHTTPAttributes(span oteltrace.Span, req *http.Request) {
	if !IsEnabled() || req == nil {
		return
	}
	
	attrs := []attribute.KeyValue{}
	if req.UserAgent() != "" {
		attrs = append(attrs, attribute.String(AttrUserAgent, req.UserAgent()))
	}
	if req.RemoteAddr != "" {
		attrs = append(attrs, attribute.String(AttrRemoteAddr, req.RemoteAddr))
	}
	
	span.SetAttributes(attrs...)
}

// SetSpanError sets error information on a span
func SetSpanError(span oteltrace.Span, err error) {
	if !IsEnabled() || err == nil {
		return
	}
	
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(
		attribute.String(AttrErrorMessage, err.Error()),
	)
}

// SetSpanErrorWithCode sets error information with a code on a span
func SetSpanErrorWithCode(span oteltrace.Span, code int, message string) {
	if !IsEnabled() {
		return
	}
	
	span.SetStatus(codes.Error, message)
	span.SetAttributes(
		attribute.Int(AttrErrorCode, code),
		attribute.String(AttrErrorMessage, message),
	)
}

// InjectTraceContext injects trace context into HTTP headers (simplified)
func InjectTraceContext(ctx context.Context, req *http.Request) {
	// Simplified implementation - would normally inject W3C trace context
	// For now, just a no-op to allow building
}

// ExtractTraceContextFromHeaders extracts trace context from HTTP headers (simplified)
func ExtractTraceContextFromHeaders(ctx context.Context, headers http.Header) context.Context {
	// Simplified implementation - would normally extract W3C trace context
	// For now, just return the original context
	return ctx
}

// RecordEvent records an event on a span
func RecordEvent(span oteltrace.Span, name string, attrs ...attribute.KeyValue) {
	if !IsEnabled() {
		return
	}
	
	span.AddEvent(name, oteltrace.WithAttributes(attrs...))
}

// SetSpanAttributes sets multiple attributes on a span
func SetSpanAttributes(span oteltrace.Span, attrs ...attribute.KeyValue) {
	if !IsEnabled() {
		return
	}
	
	span.SetAttributes(attrs...)
}