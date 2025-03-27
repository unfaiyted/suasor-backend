package types

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
