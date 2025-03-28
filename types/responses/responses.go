package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// APIResponse represents a generic API response
type APIResponse[T any] struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"Operation successful"`
	Data    T      `json:"data,omitempty"`
}

// Type-specific response creators
type EmptyResponse struct {
	Success bool `json:"success"`
}

// RespondSuccess creates a standardized success response
func RespondSuccess[T any](c *gin.Context, statusCode int, data T, message string) {
	response := APIResponse[T]{
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
