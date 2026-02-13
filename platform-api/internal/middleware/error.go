package middleware

import (
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse format
type ErrorResponse struct {
	Error      string                 `json:"error"`   // Error code (INVALID_INPUT, NOT_FOUND, etc)
	Message    string                 `json:"message"` // Human-readable message
	StatusCode int                    `json:"statusCode"`
	TraceID    string                 `json:"traceId,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"` // Additional context
}

// Custom error types
type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Details    map[string]interface{}
	Err        error // underlying error for logging
}

func (e *AppError) Error() string {
	return e.Message
}

// Factory functions for common errors
func NewValidationError(field, reason string) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    "Invalid request: " + reason,
		StatusCode: http.StatusBadRequest,
		Details: map[string]interface{}{
			"field": field,
		},
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    resource + " not found",
		StatusCode: http.StatusNotFound,
	}
}

func NewConflictError(reason string) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    reason,
		StatusCode: http.StatusConflict,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "An internal error occurred",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewUnauthorizedError(reason string) *AppError {
	return &AppError{
		Code:       "UNAUTHORIZED",
		Message:    reason,
		StatusCode: http.StatusUnauthorized,
	}
}

// ErrorHandlerMiddleware catches panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger := RequestLogger(c, "error_handler")
				logger.Error("panic recovered",
					slog.Any("panic", r),
				)

				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:      "INTERNAL_ERROR",
					Message:    "An unexpected error occurred",
					StatusCode: http.StatusInternalServerError,
					TraceID:    getTraceID(c),
				})
			}
		}()

		c.Next()

		// Check if error was set by handler
		if len(c.Errors) > 0 {
			handleError(c)
		}
	}
}

// handleError formats and returns error
func handleError(c *gin.Context) {
	logger := RequestLogger(c, "error_handler")

	lastErr := c.Errors.Last()

	var appErr *AppError
	if errors.As(lastErr.Err, &appErr) {
		logger.WarnContext(c.Request.Context(),
			"application error",
			slog.String("code", appErr.Code),
			slog.String("message", appErr.Message),
			slog.Int("statusCode", appErr.StatusCode),
		)

		c.JSON(appErr.StatusCode, ErrorResponse{
			Error:      appErr.Code,
			Message:    appErr.Message,
			StatusCode: appErr.StatusCode,
			TraceID:    getTraceID(c),
			Details:    appErr.Details,
		})
	} else {
		// Unknown error type
		logger.ErrorContext(c.Request.Context(),
			"unknown error",
			slog.String("error", lastErr.Error()),
		)

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:      "INTERNAL_ERROR",
			Message:    "An unexpected error occurred",
			StatusCode: http.StatusInternalServerError,
			TraceID:    getTraceID(c),
		})
	}
}

// Helper to get trace ID
func getTraceID(c *gin.Context) string {
	if traceID, ok := c.Get("traceId"); ok {
		return traceID.(string)
	}
	return ""
}

// RequestLogger creates a logger instance for request handling
func RequestLogger(c *gin.Context, name string) *slog.Logger {
	return slog.Default().WithGroup(name)
}

// CheckNetworkError converts network errors to app errors
func CheckNetworkError(err error) *AppError {
	if err == nil {
		return nil
	}

	// Handle specific network errors
	if errors.Is(err, net.ErrClosed) {
		return NewInternalError(err)
	}

	// Could also check for DNS errors, timeouts, etc.

	return NewInternalError(err)
}
