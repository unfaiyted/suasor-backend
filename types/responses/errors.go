package responses

import (
	"net/http"
	"time"

	"suasor/types/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ErrorResponse represents an error response
type ErrorResponse[T any] struct {
	Type       errors.ErrorType `json:"type" example:"FAILED_CHECK"`
	Message    string           `json:"message" example:"This is a pretty message"`
	StatusCode uint16           `json:"statusCode" example:"201"`
	Details    T                `json:"details,omitempty"`
	Timestamp  time.Time        `json:"timestamp"`
	RequestID  string           `json:"request_id,omitempty"`
}

// BasicErrorResponse is a concrete error response type used for Swagger docs
type BasicErrorResponse = ErrorResponse[ErrorDetails]

// RespondWithError creates a standardized error response using models.ErrorResponse
func RespondWithError(c *gin.Context, statusCode int, err error, customMessage ...string) {
	// Get error type based on status code or default to internal error
	errorType, exists := errors.StatusCodeToErrorType[statusCode]
	if !exists {
		errorType = errors.ErrorTypeInternalError
	}

	// Get default message for this error type or use a generic message
	message, exists := errors.DefaultErrorMessages[errorType]
	if !exists {
		message = "An unexpected error occurred"
	}

	// Use custom message if provided
	if len(customMessage) > 0 && customMessage[0] != "" {
		message = customMessage[0]
	}

	// Get request ID from context or generate one
	requestID := c.GetString("RequestID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// Create response with error details
	errorResponse := ErrorResponse[ErrorDetails]{
		Type:      errorType,
		Message:   message,
		Details:   CreateErrorDetails(err),
		Timestamp: time.Now(),
		RequestID: requestID,
	}

	c.JSON(statusCode, errorResponse)
}

// Convenience functions for common error responses
func RespondBadRequest(c *gin.Context, err error, customMessage ...string) {
	RespondWithError(c, http.StatusBadRequest, err, customMessage...)
}

func RespondUnauthorized(c *gin.Context, err error, customMessage ...string) {
	RespondWithError(c, http.StatusUnauthorized, err, customMessage...)
}

func RespondForbidden(c *gin.Context, err error, customMessage ...string) {
	RespondWithError(c, http.StatusForbidden, err, customMessage...)
}

func RespondNotFound(c *gin.Context, err error, customMessage ...string) {
	RespondWithError(c, http.StatusNotFound, err, customMessage...)
}

func RespondConflict(c *gin.Context, err error, customMessage ...string) {
	RespondWithError(c, http.StatusConflict, err, customMessage...)
}

func RespondValidationError(c *gin.Context, err error, customMessage ...string) {
	errorType := errors.ErrorTypeValidation
	message := errors.DefaultErrorMessages[errorType]

	if len(customMessage) > 0 && customMessage[0] != "" {
		message = customMessage[0]
	}

	RespondWithError(c, http.StatusBadRequest, err, message)
}

func RespondInternalError(c *gin.Context, err error, customMessage ...string) {
	RespondWithError(c, http.StatusInternalServerError, err, customMessage...)
}

func RespondServiceUnavailable(c *gin.Context, err error, customMessage ...string) {
	RespondWithError(c, http.StatusServiceUnavailable, err, customMessage...)
}

func RespondNotImplemented(c *gin.Context, err error, customMessage ...string) {
	RespondWithError(c, http.StatusNotImplemented, err, customMessage...)
}

// Common error detail types
type ValidationErrorDetails struct {
	FieldErrors map[string]string `json:"fieldErrors,omitempty"`
}

type NotFoundErrorDetails struct {
	Resource string `json:"resource,omitempty"`
	ID       string `json:"id,omitempty"`
}

// EmptyErrorDetails for errors without specific details
type EmptyErrorDetails struct{}

// Error response constructors
func NewValidationError(message string, fieldErrors map[string]string, requestID string) ErrorResponse[ValidationErrorDetails] {
	return ErrorResponse[ValidationErrorDetails]{
		Type:       errors.ErrorTypeValidation,
		Message:    message,
		StatusCode: 422,
		Details:    ValidationErrorDetails{FieldErrors: fieldErrors},
		Timestamp:  time.Now(),
		RequestID:  requestID,
	}
}

func NewNotFoundError(message string, resource, id, requestID string) ErrorResponse[NotFoundErrorDetails] {
	return ErrorResponse[NotFoundErrorDetails]{
		Type:       errors.ErrorTypeNotFound,
		Message:    message,
		StatusCode: 404,
		Details:    NotFoundErrorDetails{Resource: resource, ID: id},
		Timestamp:  time.Now(),
		RequestID:  requestID,
	}
}

func NewGenericError(errorType errors.ErrorType, message string, statusCode uint16, requestID string) ErrorResponse[EmptyErrorDetails] {
	return ErrorResponse[EmptyErrorDetails]{
		Type:       errorType,
		Message:    message,
		StatusCode: statusCode,
		Details:    EmptyErrorDetails{},
		Timestamp:  time.Now(),
		RequestID:  requestID,
	}
}

// ErrorDetails captures all information from Go error objects
type ErrorDetails struct {
	Error     string         `json:"error"`
	StackInfo string         `json:"stackInfo,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// CreateErrorDetails extracts all available information from an error
func CreateErrorDetails(err error) ErrorDetails {
	if err == nil {
		return ErrorDetails{}
	}

	details := ErrorDetails{
		Error:    err.Error(),
		Metadata: make(map[string]interface{}),
	}

	// Extract additional info based on error type
	switch typedErr := err.(type) {
	case interface{ StackTrace() string }:
		details.StackInfo = typedErr.StackTrace()
	case interface{ Fields() map[string]interface{} }:
		details.Metadata = typedErr.Fields()
	case interface{ Unwrap() error }:
		// Handle wrapped errors if needed
		if unwrapped := typedErr.Unwrap(); unwrapped != nil {
			details.Metadata["wrapped"] = unwrapped.Error()
		}
	}

	return details
}
