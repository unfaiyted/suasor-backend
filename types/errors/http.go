package errors

import (
	"net/http"
)

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

// StatusCodeToErrorType maps HTTP status codes to ErrorType
var StatusCodeToErrorType = map[int]ErrorType{
	http.StatusBadRequest:          ErrorTypeBadRequest,
	http.StatusUnauthorized:        ErrorTypeUnauthorized,
	http.StatusForbidden:           ErrorTypeForbidden,
	http.StatusNotFound:            ErrorTypeNotFound,
	http.StatusConflict:            ErrorTypeConflict,
	http.StatusUnprocessableEntity: ErrorTypeUnprocessableEntity,
	http.StatusTooManyRequests:     ErrorTypeRateLimited,
	http.StatusInternalServerError: ErrorTypeInternalError,
	http.StatusServiceUnavailable:  ErrorTypeServiceUnavailable,
	http.StatusGatewayTimeout:      ErrorTypeTimeout,
}

// DefaultErrorMessages maps ErrorType to default human-readable messages
var DefaultErrorMessages = map[ErrorType]string{
	ErrorTypeBadRequest:          "The request could not be processed due to invalid parameters",
	ErrorTypeUnauthorized:        "Authentication is required to access this resource",
	ErrorTypeForbidden:           "You don't have permission to access this resource",
	ErrorTypeNotFound:            "The requested resource was not found",
	ErrorTypeConflict:            "The request conflicts with the current state of the resource",
	ErrorTypeValidation:          "The request contains validation errors",
	ErrorTypeRateLimited:         "Too many requests, please try again later",
	ErrorTypeTimeout:             "The operation timed out",
	ErrorTypeInternalError:       "An internal server error occurred",
	ErrorTypeServiceUnavailable:  "The service is currently unavailable",
	ErrorTypeUnprocessableEntity: "The request was well-formed but cannot be processed",
}
