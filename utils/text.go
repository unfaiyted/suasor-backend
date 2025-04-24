package utils

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"suasor/utils/logger"
)

// Helper function to truncate long strings for logging
func Truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

func GenerateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	// Create a byte array with the code length
	shortCode := make([]byte, codeLength)

	// Use crypto/rand for secure random generation
	randomBytes := make([]byte, codeLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fallback to less secure but working solution if crypto/rand fails
		for i := range shortCode {
			shortCode[i] = charset[time.Now().UnixNano()%int64(len(charset))]
			time.Sleep(1 * time.Nanosecond) // Add a tiny delay to change the seed
		}
		return string(shortCode)
	}

	// Map random bytes to characters in the charset
	for i, b := range randomBytes {
		shortCode[i] = charset[int(b)%len(charset)]
	}

	return string(shortCode)
}

// containsIgnoreCase checks if a string contains a substring, ignoring case
func ContainsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// GetUserID extracts the user ID from query parameters or from the authenticated context
// If userID query parameter is provided, it uses that value; otherwise, uses the authenticated user's ID
// Returns an error if no valid user ID is found or if the user ID is not a valid uint64
func GetUserID(c *gin.Context) (uint64, error) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Check if userId is provided in query parameters
	userIDStr := c.Query("userID")
	if userIDStr != "" {
		// Parse the userID from query parameter
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			log.Warn().Err(err).Str("userID", userIDStr).Msg("Invalid user ID format in query parameter")
			return 0, err
		}
		return userID, nil
	}

	// If not in query parameters, check the context (set by auth middleware)
	contextUserID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("User ID not found in context")
		return 0, fmt.Errorf("User ID not found")
	}

	// Convert the userID from context to uint64
	switch v := contextUserID.(type) {
	case uint64:
		return v, nil
	case int:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	default:
		log.Warn().Interface("contextUserID", contextUserID).Msg("User ID in context has unexpected type")
		return 0, fmt.Errorf("User ID has invalid type")
	}
}

// GetLimit extracts and validates the limit parameter from query parameters
// If limit is not provided or invalid, returns the default value
// The default value is used when limit is not provided, invalid, or outside the allowed range
// If forceMax is true, the limit will be capped at defaultMax regardless of the input value
func GetLimit(c *gin.Context, defaultValue, defaultMax int, forceMax bool) int {
	limitStr := c.Query("limit")
	if limitStr == "" {
		return defaultValue
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		// Invalid limit value, use default
		return defaultValue
	}

	// If limit is larger than max and we're forcing a maximum, cap it
	if forceMax && limit > defaultMax {
		return defaultMax
	}

	return limit
}

// GetOffset extracts and validates the offset parameter from query parameters
// If offset is not provided or invalid, returns the default value (usually 0)
func GetOffset(c *gin.Context, defaultValue int) int {
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		return defaultValue
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		// Invalid offset value, use default
		return defaultValue
	}

	return offset
}

// GetPage is a helper that calculates page from offset and limit
// This is useful when the underlying API uses page-based pagination but our API exposes offset/limit
func GetPage(offset, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (offset / limit) + 1
}

// GetListID extracts and validates the list ID from path or query parameters
// Checks both "id", "listId", and "list_id" parameters for flexibility
// Returns an error if no valid list ID can be found or if it's not a valid uint64
func GetListID(c *gin.Context) (uint64, error) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// First try path parameter "id" (most common)
	idStr := c.Param("id")
	if idStr == "" {
		// Try path parameter "listId"
		idStr = c.Param("listId")
	}
	if idStr == "" {
		// Try path parameter "list_id"
		idStr = c.Param("list_id")
	}

	// If still not found in path, check query parameters
	if idStr == "" {
		idStr = c.Query("id")
	}
	if idStr == "" {
		idStr = c.Query("listId")
	}
	if idStr == "" {
		idStr = c.Query("list_id")
	}

	// If no ID found in any parameter
	if idStr == "" {
		log.Warn().Msg("List ID not found in request parameters")
		return 0, fmt.Errorf("List ID is required")
	}

	// Parse the ID as uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("listId", idStr).Msg("Invalid list ID format")
		return 0, err
	}

	return id, nil
}

// GetDays extracts and validates the days parameter from query parameters
// If days is not provided or invalid, returns the default value
// The default value is used when days is not provided, invalid, or outside the allowed range
func GetDays(c *gin.Context, defaultValue int) int {
	daysStr := c.Query("days")
	if daysStr == "" {
		return defaultValue
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 0 {
		// Invalid days value, use default
		return defaultValue
	}

	return days
}

func GetRequiredParam(c *gin.Context, paramName string) (string, error) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	param := c.Param(paramName)
	if param == "" {
		log.Warn().Str(paramName, param).Msg("Required parameter not found in request parameters")
		return "", fmt.Errorf("required parameter not found in request parameters")
	}

	return param, nil
}
