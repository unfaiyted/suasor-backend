// models/responses.go
package models

import "time"

type ErrorType string

const (
	ErrorTypeFailedCheck         ErrorType = "FAILED_CHECK"
	ErrorTypeUnauthorized        ErrorType = "UNAUTHORIZED"
	ErrorTypeNotFound            ErrorType = "NOT_FOUND"
	ErrorTypeBadRequest          ErrorType = "BAD_REQUEST"
	ErrorTypeInternalError       ErrorType = "INTERNAL_ERROR"
	ErrorTypeForbidden           ErrorType = "FORBIDDEN"
	ErrorTypeConflict            ErrorType = "CONFLICT"
	ErrorTypeValidation          ErrorType = "VALIDATION_ERROR"
	ErrorTypeRateLimited         ErrorType = "RATE_LIMITED"
	ErrorTypeTimeout             ErrorType = "TIMEOUT"
	ErrorTypeServiceUnavailable  ErrorType = "SERVICE_UNAVAILABLE"
	ErrorTypeUnprocessableEntity ErrorType = "UNPROCESSABLE_ENTITY"
)

// ErrorResponse represents an error response
type ErrorResponse[T any] struct {
	Type       ErrorType `json:"type" example:"FAILED_CHECK"`
	Message    string    `json:"message" example:"This is a pretty message"`
	StatusCode uint16    `json:"statusCode" example:"201"`
	Details    T         `json:"details,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	RequestID  string    `json:"request_id,omitempty"`
}

// APIResponse represents a generic API response
type APIResponse[T any] struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"Operation successful"`
	Data    T      `json:"data,omitempty"`
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
		Type:       ErrorTypeValidation,
		Message:    message,
		StatusCode: 422,
		Details:    ValidationErrorDetails{FieldErrors: fieldErrors},
		Timestamp:  time.Now(),
		RequestID:  requestID,
	}
}

func NewNotFoundError(message string, resource, id, requestID string) ErrorResponse[NotFoundErrorDetails] {
	return ErrorResponse[NotFoundErrorDetails]{
		Type:       ErrorTypeNotFound,
		Message:    message,
		StatusCode: 404,
		Details:    NotFoundErrorDetails{Resource: resource, ID: id},
		Timestamp:  time.Now(),
		RequestID:  requestID,
	}
}

func NewGenericError(errorType ErrorType, message string, statusCode uint16, requestID string) ErrorResponse[EmptyErrorDetails] {
	return ErrorResponse[EmptyErrorDetails]{
		Type:       errorType,
		Message:    message,
		StatusCode: statusCode,
		Details:    EmptyErrorDetails{},
		Timestamp:  time.Now(),
		RequestID:  requestID,
	}
}

// Type-specific response creators
func NewShortenResponse(shorten *Shorten, message string) APIResponse[ShortenData] {
	return APIResponse[ShortenData]{
		Success: true,
		Message: message,
		Data:    ShortenData{Shorten: shorten},
	}
}
