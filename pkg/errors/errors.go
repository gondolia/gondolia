// Package errors provides error types with context for debugging.
// All errors include stack traces, error codes, and can be traced.
package errors

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gondolia/gondolia/pkg/telemetry"
)

// ErrorCode represents a unique error code
type ErrorCode string

// Error codes
const (
	// Client errors (4xx)
	ErrCodeBadRequest       ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeConflict         ErrorCode = "CONFLICT"
	ErrCodeValidation       ErrorCode = "VALIDATION_ERROR"
	ErrCodeRateLimited      ErrorCode = "RATE_LIMITED"

	// Server errors (5xx)
	ErrCodeInternal         ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase         ErrorCode = "DATABASE_ERROR"
	ErrCodeExternalService  ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeTimeout          ErrorCode = "TIMEOUT"
	ErrCodeUnavailable      ErrorCode = "SERVICE_UNAVAILABLE"
)

// AppError represents an application error with context
type AppError struct {
	Code       ErrorCode         `json:"code"`
	Message    string            `json:"message"`
	Details    string            `json:"details,omitempty"`
	HTTPStatus int               `json:"-"`
	TraceID    string            `json:"trace_id,omitempty"`
	SpanID     string            `json:"span_id,omitempty"`
	Stack      string            `json:"-"`
	Cause      error             `json:"-"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithContext adds trace context to the error
func (e *AppError) WithContext(ctx context.Context) *AppError {
	e.TraceID = telemetry.TraceID(ctx)
	e.SpanID = telemetry.SpanID(ctx)

	// Record error in trace
	telemetry.SetError(ctx, e)

	return e
}

// WithMetadata adds metadata to the error
func (e *AppError) WithMetadata(key, value string) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: codeToHTTPStatus(code),
		Stack:      captureStack(),
	}
}

// Wrap wraps an existing error with context
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: codeToHTTPStatus(code),
		Cause:      err,
		Stack:      captureStack(),
	}
}

// Convenience constructors

// BadRequest creates a 400 error
func BadRequest(message string) *AppError {
	return New(ErrCodeBadRequest, message)
}

// Unauthorized creates a 401 error
func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message)
}

// Forbidden creates a 403 error
func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message)
}

// NotFound creates a 404 error
func NotFound(message string) *AppError {
	return New(ErrCodeNotFound, message)
}

// Conflict creates a 409 error
func Conflict(message string) *AppError {
	return New(ErrCodeConflict, message)
}

// Validation creates a validation error
func Validation(message string) *AppError {
	return New(ErrCodeValidation, message)
}

// Internal creates a 500 error
func Internal(message string) *AppError {
	return New(ErrCodeInternal, message)
}

// Database creates a database error
func Database(err error, message string) *AppError {
	return Wrap(err, ErrCodeDatabase, message)
}

// ExternalService creates an external service error
func ExternalService(err error, service string) *AppError {
	return Wrap(err, ErrCodeExternalService, fmt.Sprintf("external service error: %s", service))
}

// Timeout creates a timeout error
func Timeout(message string) *AppError {
	return New(ErrCodeTimeout, message)
}

// codeToHTTPStatus maps error codes to HTTP status codes
func codeToHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest, ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeRateLimited:
		return http.StatusTooManyRequests
	case ErrCodeTimeout:
		return http.StatusGatewayTimeout
	case ErrCodeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// captureStack captures the stack trace
func captureStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCodeNotFound
	}
	return false
}

// IsValidation checks if an error is a validation error
func IsValidation(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCodeValidation
	}
	return false
}

// GetHTTPStatus returns the HTTP status code for an error
func GetHTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}