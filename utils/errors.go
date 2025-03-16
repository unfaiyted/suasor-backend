package utils

import (
	"net/http"
	"suasor/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StatusCodeToErrorType maps HTTP status codes to ErrorType
var StatusCodeToErrorType = map[int]models.ErrorType{
	http.StatusBadRequest:          models.ErrorTypeBadRequest,
	http.StatusUnauthorized:        models.ErrorTypeUnauthorized,
	http.StatusForbidden:           models.ErrorTypeForbidden,
	http.StatusNotFound:            models.ErrorTypeNotFound,
	http.StatusConflict:            models.ErrorTypeConflict,
	http.StatusUnprocessableEntity: models.ErrorTypeUnprocessableEntity,
	http.StatusTooManyRequests:     models.ErrorTypeRateLimited,
	http.StatusInternalServerError: models.ErrorTypeInternalError,
	http.StatusServiceUnavailable:  models.ErrorTypeServiceUnavailable,
	http.StatusGatewayTimeout:      models.ErrorTypeTimeout,
}

// DefaultErrorMessages maps ErrorType to default human-readable messages
var DefaultErrorMessages = map[models.ErrorType]string{
	models.ErrorTypeBadRequest:          "The request could not be processed due to invalid parameters",
	models.ErrorTypeUnauthorized:        "Authentication is required to access this resource",
	models.ErrorTypeForbidden:           "You don't have permission to access this resource",
	models.ErrorTypeNotFound:            "The requested resource was not found",
	models.ErrorTypeConflict:            "The request conflicts with the current state of the resource",
	models.ErrorTypeValidation:          "The request contains validation errors",
	models.ErrorTypeRateLimited:         "Too many requests, please try again later",
	models.ErrorTypeTimeout:             "The operation timed out",
	models.ErrorTypeInternalError:       "An internal server error occurred",
	models.ErrorTypeServiceUnavailable:  "The service is currently unavailable",
	models.ErrorTypeUnprocessableEntity: "The request was well-formed but cannot be processed",
}

// RespondWithError creates a standardized error response using models.ErrorResponse
func RespondWithError(c *gin.Context, statusCode int, err error, customMessage ...string) {
	// Get error type based on status code or default to internal error
	errorType, exists := StatusCodeToErrorType[statusCode]
	if !exists {
		errorType = models.ErrorTypeInternalError
	}

	// Get default message for this error type or use a generic message
	message, exists := DefaultErrorMessages[errorType]
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
	errorResponse := models.ErrorResponse[error]{
		Type:      errorType,
		Message:   message,
		Details:   err,
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
	errorType := models.ErrorTypeValidation
	message := DefaultErrorMessages[errorType]

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

// RespondSuccess creates a standardized success response
func RespondSuccess[T any](c *gin.Context, statusCode int, data T, message string) {
	response := models.APIResponse[T]{
		Success: true,
		Data:    data,
		Message: message,
	}

	c.JSON(statusCode, response)
}

// Convenience functions for success responses
func RespondOK[T any](c *gin.Context, data T, message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	RespondSuccess(c, http.StatusOK, data, msg)
}

func RespondCreated[T any](c *gin.Context, data T, message ...string) {
	msg := "Resource created successfully"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	RespondSuccess(c, http.StatusCreated, data, msg)
}
