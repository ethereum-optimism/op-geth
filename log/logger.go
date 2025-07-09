package log

import (
	"context"
	"log/slog"
	"math"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	legacyLevelCrit = iota
	legacyLevelError
	legacyLevelWarn
	legacyLevelInfo
	legacyLevelDebug
	legacyLevelTrace
)

const (
	levelMaxVerbosity slog.Level = math.MinInt
	LevelTrace        slog.Level = -8
	LevelDebug                   = slog.LevelDebug
	LevelInfo                    = slog.LevelInfo
	LevelWarn                    = slog.LevelWarn
	LevelError                   = slog.LevelError
	LevelCrit         slog.Level = 12

	// for backward-compatibility
	LvlTrace = LevelTrace
	LvlInfo  = LevelInfo
	LvlDebug = LevelDebug
)

// FromLegacyLevel converts from old Geth verbosity level constants
// to levels defined by slog
func FromLegacyLevel(lvl int) slog.Level {
	switch lvl {
	case legacyLevelCrit:
		return LevelCrit
	case legacyLevelError:
		return slog.LevelError
	case legacyLevelWarn:
		return slog.LevelWarn
	case legacyLevelInfo:
		return slog.LevelInfo
	case legacyLevelDebug:
		return slog.LevelDebug
	case legacyLevelTrace:
		return LevelTrace
	default:
		break
	}

	// TODO: should we allow use of custom levels or force them to match existing max/min if they fall outside the range as I am doing here?
	if lvl > legacyLevelTrace {
		return LevelTrace
	}
	return LevelCrit
}

// LevelAlignedString returns a 5-character string containing the name of a Lvl.
func LevelAlignedString(l slog.Level) string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO "
	case slog.LevelWarn:
		return "WARN "
	case slog.LevelError:
		return "ERROR"
	case LevelCrit:
		return "CRIT "
	default:
		return "unknown level"
	}
}

// LevelString returns a string containing the name of a Lvl.
func LevelString(l slog.Level) string {
	switch l {
	case LevelTrace:
		return "trace"
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	case LevelCrit:
		return "crit"
	default:
		return "unknown"
	}
}

// A Logger writes key/value pairs to a Handler.
// Each key/value pair can be two consecutive references, or a slog.Attr, and both may occur in the same log call.
type Logger interface {
	// With returns a new Logger that has this logger's attributes plus the given attributes
	With(args ...any) Logger

	// New returns a new Logger that has this logger's attributes plus the given attributes. Identical to 'With'.
	New(args ...any) Logger

	// Log logs a message at the specified level with context key/value pairs.
	Log(level slog.Level, msg string, args ...any)

	// Trace log a message at the trace level with context key/value pairs.
	Trace(msg string, args ...any)

	// Debug logs a message at the debug level with context key/value pairs
	Debug(msg string, args ...any)

	// Info logs a message at the info level with context key/value pairs
	Info(msg string, args ...any)

	// Warn logs a message at the warn level with context key/value pairs
	Warn(msg string, args ...any)

	// Error logs a message at the error level with context key/value pairs
	Error(msg string, args ...any)

	// Crit logs a message at the crit level with context key/value pairs, and exits.
	// Warning: for legacy compatibility this runs os.Exit(1).
	Crit(msg string, args ...any)

	// Write logs a message at the specified level
	Write(level slog.Level, msg string, attrs ...any)

	// Enabled reports whether l emits log records at the given context and level.
	Enabled(ctx context.Context, level slog.Level) bool

	// Handler returns the underlying handler of the inner logger.
	Handler() slog.Handler

	// SetContext sets the default context: this is used for every log call without specified context.
	// Sub-loggers inherit this context.
	// Contexts may be used to filter log records before attributes processing: see slog.Handler.Enabled.
	SetContext(ctx context.Context)

	// WriteCtx logs a message at the specified level with attribute key/value pairs and/or slog.Attr pairs.
	WriteCtx(ctx context.Context, level slog.Level, msg string, attrs ...any)

	// LogAttrs is a more efficient version of [Logger.Log] that accepts only Attrs.
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)

	// TraceContext logs at [LevelTrace] with the given context.
	TraceContext(ctx context.Context, msg string, args ...any)

	// DebugContext logs at [LevelDebug] with the given context.
	DebugContext(ctx context.Context, msg string, args ...any)

	// InfoContext logs at [LevelInfo] with the given context.
	InfoContext(ctx context.Context, msg string, args ...any)

	// WarnContext logs at [LevelWarn] with the given context.
	WarnContext(ctx context.Context, msg string, args ...any)

	// ErrorContext logs at [LevelError] with the given context.
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type logger struct {
	inner      *slog.Logger
	defaultCtx atomic.Pointer[context.Context] // used when no context is specified
}

// NewLogger returns a logger with the specified handler set
func NewLogger(h slog.Handler) Logger {
	out := &logger{inner: slog.New(h)}
	ctx := context.Background()
	out.defaultCtx.Store(&ctx)
	return out
}

// SetContext sets the default context: this is used for every log call without specified context.
// Sub-loggers inherit this context.
func (l *logger) SetContext(ctx context.Context) {
	l.defaultCtx.Store(&ctx)
}

func (l *logger) Handler() slog.Handler {
	return l.inner.Handler()
}

// Write logs a message at the specified level.
func (l *logger) Write(level slog.Level, msg string, attrs ...any) {
	l.writeCtx(*l.defaultCtx.Load(), level, msg, attrs...)
}

// WriteCtx logs a message at the specified level, with context.
func (l *logger) WriteCtx(ctx context.Context, level slog.Level, msg string, attrs ...any) {
	l.writeCtx(ctx, level, msg, attrs...)
}

// writeCtx basically does what the inner slog.Logger.log function would do,
// but adjusts to PC, so call-site calculation is still accurate (see TestLoggingWithVmodule).
func (l *logger) writeCtx(ctx context.Context, level slog.Level, msg string, attrs ...any) {
	if !l.inner.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller (Write/WriteCtx), and the caller's caller]
	runtime.Callers(4, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(attrs...)
	_ = l.inner.Handler().Handle(ctx, r)
}

func (l *logger) writeCtxAttr(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if !l.inner.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = l.inner.Handler().Handle(ctx, r)
}

func (l *logger) Log(level slog.Level, msg string, attrs ...any) {
	l.Write(level, msg, attrs...)
}

func (l *logger) With(args ...interface{}) Logger {
	out := &logger{inner: l.inner.With(args...)}
	out.defaultCtx.Store(l.defaultCtx.Load())
	return out
}

func (l *logger) New(args ...interface{}) Logger {
	return l.With(args...)
}

// Enabled reports whether l emits log records at the given context and level.
func (l *logger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.inner.Enabled(ctx, level)
}

func (l *logger) Trace(msg string, ctx ...interface{}) {
	l.Write(LevelTrace, msg, ctx...)
}

func (l *logger) Debug(msg string, ctx ...interface{}) {
	l.Write(slog.LevelDebug, msg, ctx...)
}

func (l *logger) Info(msg string, ctx ...interface{}) {
	l.Write(slog.LevelInfo, msg, ctx...)
}

func (l *logger) Warn(msg string, ctx ...any) {
	l.Write(slog.LevelWarn, msg, ctx...)
}

func (l *logger) Error(msg string, ctx ...interface{}) {
	l.Write(slog.LevelError, msg, ctx...)
}

func (l *logger) Crit(msg string, ctx ...interface{}) {
	l.Write(LevelCrit, msg, ctx...)
	os.Exit(1)
}

func (l *logger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	l.writeCtxAttr(ctx, level, msg, attrs...)
}

func (l *logger) TraceContext(ctx context.Context, msg string, args ...any) {
	l.WriteCtx(ctx, LevelTrace, msg, args...)
}

func (l *logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.WriteCtx(ctx, LevelDebug, msg, args...)
}

func (l *logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.WriteCtx(ctx, LevelInfo, msg, args...)
}

func (l *logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.WriteCtx(ctx, LevelWarn, msg, args...)
}

func (l *logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.WriteCtx(ctx, LevelError, msg, args...)
}
