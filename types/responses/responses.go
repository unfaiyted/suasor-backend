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

// The following type aliases are for Swagger documentation
// They help swagger resolve generic type issues

// StringResponse is the APIResponse for a string
type StringResponse = APIResponse[string]

type SuccessResponse = APIResponse[EmptyResponse]

// HealthCheckResponse is for health check responses
type HealthCheckResponse = APIResponse[HealthResponse]

// UserProfileResponse is the API response for user profile
type UserProfileResponse = APIResponse[UserResponse]

// ClientsResponse is the API response for clients
type ClientsResponse = APIResponse[[]ClientResponse]

// EmptyAPIResponse is the APIResponse for empty data
type EmptyAPIResponse = APIResponse[any]

// Type-specific response creators
type EmptyResponse struct {
	Success bool `json:"success"`
}

// TestConnectionResponse contains details about a connection test
type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
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

func RespondListOK[T any](c *gin.Context, data T, count int, message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	type temp struct {
		Items T   `json:"items"`
		Total int `json:"total"`
	}

	wrappedData := temp{
		Items: data,
		Total: count,
	}

	RespondSuccess(c, http.StatusOK, wrappedData, msg)
}

func RespondCreated[T any](c *gin.Context, data T, message ...string) {
	msg := "Resource created successfully"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	RespondSuccess(c, http.StatusCreated, data, msg)
}
