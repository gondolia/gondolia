package telemetry

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
	// TraceIDHeader is the header name for trace ID
	TraceIDHeader = "X-Trace-ID"
)

// TracingMiddleware creates a Gin middleware for distributed tracing
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		// Extract trace context from incoming request
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Generate or extract request ID
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Start span
		spanName := c.Request.Method + " " + c.FullPath()
		if c.FullPath() == "" {
			spanName = c.Request.Method + " " + c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.target", c.Request.URL.Path),
				attribute.String("http.host", c.Request.Host),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("http.request_id", requestID),
				attribute.String("net.peer.ip", c.ClientIP()),
			),
		)
		defer span.End()

		// Store in context
		c.Request = c.Request.WithContext(ctx)
		c.Set("request_id", requestID)
		c.Set("trace_id", span.SpanContext().TraceID().String())

		// Add trace headers to response
		c.Header(RequestIDHeader, requestID)
		c.Header(TraceIDHeader, span.SpanContext().TraceID().String())

		// Record start time
		start := time.Now()

		// Process request
		c.Next()

		// Record response attributes
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int64("http.response_time_ms", duration.Milliseconds()),
			attribute.Int("http.response_size", c.Writer.Size()),
		)

		// Record errors
		if len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last().Err)
			span.SetAttributes(attribute.String("error.message", c.Errors.String()))
		}

		// Set span status based on HTTP status
		if statusCode >= 400 {
			span.SetAttributes(attribute.Bool("error", true))
		}
	}
}

// LoggingMiddleware creates a Gin middleware for structured logging
func LoggingMiddleware(logger interface {
	InfoContext(ctx interface{}, msg string, fields ...interface{})
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Build log message
		msg := "HTTP Request"
		if statusCode >= 500 {
			msg = "HTTP Request ERROR"
		} else if statusCode >= 400 {
			msg = "HTTP Request WARN"
		}

		// Log would be called here with context
		_ = map[string]interface{}{
			"method":       c.Request.Method,
			"path":         path,
			"query":        query,
			"status":       statusCode,
			"duration_ms":  duration.Milliseconds(),
			"client_ip":    c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"request_id":   c.GetString("request_id"),
			"trace_id":     c.GetString("trace_id"),
			"error":        c.Errors.String(),
			"body_size":    c.Writer.Size(),
		}

		_ = msg // Use structured logger here
	}
}

// ErrorMiddleware creates a Gin middleware for error handling
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Get trace context for error response
			traceID := c.GetString("trace_id")
			requestID := c.GetString("request_id")

			// Build error response
			response := gin.H{
				"error": gin.H{
					"message":    err.Error(),
					"trace_id":   traceID,
					"request_id": requestID,
				},
			}

			// Return appropriate status
			status := c.Writer.Status()
			if status == http.StatusOK {
				status = http.StatusInternalServerError
			}

			c.JSON(status, response)
		}
	}
}
