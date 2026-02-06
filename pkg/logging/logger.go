// Package logging provides structured logging with trace correlation.
// All logs include trace_id, span_id, and request_id for easy debugging.
package logging

import (
	"context"
	"os"

	"github.com/gondolia/gondolia/pkg/telemetry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ContextKey is the key type for context values
type ContextKey string

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// TenantIDKey is the context key for tenant ID
	TenantIDKey ContextKey = "tenant_id"
)

// Logger wraps zap.Logger with context-aware logging
type Logger struct {
	*zap.Logger
	serviceName string
}

// Config holds the logger configuration
type Config struct {
	ServiceName string
	Environment string // "development" or "production"
	Level       string // "debug", "info", "warn", "error"
}

// NewLogger creates a new logger instance
func NewLogger(cfg Config) (*Logger, error) {
	var zapConfig zap.Config

	if cfg.Environment == "development" {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Set log level
	switch cfg.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(
			zap.String("service", cfg.ServiceName),
			zap.String("hostname", getHostname()),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger:      logger,
		serviceName: cfg.ServiceName,
	}, nil
}

// WithContext returns a logger with trace context fields
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	fields := []zap.Field{}

	// Add trace context
	if traceID := telemetry.TraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if spanID := telemetry.SpanID(ctx); spanID != "" {
		fields = append(fields, zap.String("span_id", spanID))
	}

	// Add request context
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		fields = append(fields, zap.String("request_id", requestID))
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		fields = append(fields, zap.String("user_id", userID))
	}
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		fields = append(fields, zap.String("tenant_id", tenantID))
	}

	return l.Logger.With(fields...)
}

// Debug logs a debug message with context
func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Debug(msg, fields...)
}

// Info logs an info message with context
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Info(msg, fields...)
}

// Warn logs a warning message with context
func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Warn(msg, fields...)
}

// Error logs an error message with context
func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Error(msg, fields...)
}

// Fatal logs a fatal message with context and exits
func (l *Logger) FatalContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Fatal(msg, fields...)
}

// WithError adds an error field
func WithError(err error) zap.Field {
	return zap.Error(err)
}

// WithField adds a field
func WithField(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}