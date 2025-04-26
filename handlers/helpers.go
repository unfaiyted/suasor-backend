package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/types/responses"
	"suasor/utils/logger"
)

func checkAdminAccess(c *gin.Context) (uint64, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Check if user is authenticated
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Unauthorized access attempt")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return 0, false
	}

	// Check if user has admin role
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		log.Warn().
			Interface("userID", userID).
			Msg("Forbidden access attempt - admin required")
		responses.RespondForbidden(c, nil, "Admin privileges required")
		return 0, false
	}

	return userID.(uint64), true
}

func checkUserAccess(c *gin.Context) (uint64, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Check if user is authenticated
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Unauthorized access attempt")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return 0, false
	}

	log.Info().
		Uint64("userID", userID.(uint64)).
		Msg("User access granted")
	return userID.(uint64), true
}

func checkUserRole(c *gin.Context) (string, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Check if user is authenticated
	userRole, exists := c.Get("userRole")
	if !exists {
		log.Warn().Msg("Unauthorized access attempt")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return "", false
	}

	return userRole.(string), true
}

func checkClientID(c *gin.Context) (uint64, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Check if clientID is a valid uint64
	clientIDStr, exists := c.Get("clientID")
	if !exists {
		log.Warn().Msg("Client ID not found in context")
		responses.RespondBadRequest(c, nil, "Client ID not found in context")
		return 0, false
	}
	clientID, err := strconv.ParseUint(clientIDStr.(string), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientID", clientIDStr.(string)).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID format")
		return 0, false
	}

	return clientID, true
}

func checkClientType(c *gin.Context) (clienttypes.ClientType, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Check if clientType exists in context
	clientTypeVal, exists := c.Get("clientType")
	if !exists {
		log.Warn().Msg("Client type not found in context")
		responses.RespondBadRequest(c, nil, "Client type not found in context")
		return "", false
	}

	// Handle different types that might be in the context
	var clientType clienttypes.ClientType

	switch ct := clientTypeVal.(type) {
	case clienttypes.ClientType:
		// Direct ClientType object
		clientType = ct
		log.Debug().Str("clientType", string(clientType)).Msg("Got ClientType directly")
	case string:
		// String that needs conversion
		clientType = clienttypes.ClientType(ct)
		log.Debug().Str("clientType", ct).Msg("Got string client type")
	default:
		// Unknown type
		log.Warn().
			Str("type", fmt.Sprintf("%T", clientTypeVal)).
			Msg("Invalid client type format in context")
		responses.RespondBadRequest(c, nil, "Invalid client type format")
		return "", false
	}

	// Validate the client type
	switch clientType {
	case clienttypes.ClientTypeEmby, clienttypes.ClientTypeJellyfin, clienttypes.ClientTypePlex,
		clienttypes.ClientTypeSubsonic, clienttypes.ClientTypeSonarr, clienttypes.ClientTypeLidarr,
		clienttypes.ClientTypeRadarr, clienttypes.ClientTypeClaude, clienttypes.ClientTypeOpenAI,
		clienttypes.ClientTypeOllama:
		log.Debug().Str("clientType", string(clientType)).Msg("Client type valid")
		return clientType, true
	default:
		log.Warn().Str("clientType", string(clientType)).Msg("Invalid client type")
		responses.RespondBadRequest(c, nil, "Invalid client type")
		return "", false
	}

}

func checkItemID(c *gin.Context, paramName string) (uint64, error) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	itemIDStr := c.Param(paramName)
	if itemIDStr == "" {
		log.Warn().Msg("Item ID not found in request parameters")
		responses.RespondBadRequest(c, nil, "Item ID not found in request parameters")
		return 0, fmt.Errorf("item ID not found in request parameters")
	}

	itemID, err := strconv.ParseUint(itemIDStr, 10, 64)
	if err != nil {
		log.Warn().Err(err).Str(paramName, itemIDStr).Msg("Invalid item ID format")
		responses.RespondBadRequest(c, err, "Invalid item ID format")
		return 0, fmt.Errorf("invalid item ID format: %s", itemIDStr)
	}

	return itemID, nil
}

func checkClientItemID(c *gin.Context, paramName string) (string, error) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	itemIDStr := c.Param(paramName)
	if itemIDStr == "" {
		log.Warn().Msg("Client item ID not found in request parameters")
		responses.RespondBadRequest(c, nil, "Client item ID not found in request parameters")
		return "", fmt.Errorf("client item ID not found in request parameters")
	}

	return itemIDStr, nil
}

// checkJSONBinding binds JSON request and handles validation errors consistently
func checkJSONBinding(c *gin.Context, req any) bool {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	if err := c.ShouldBindJSON(req); err != nil {
		log.Error().Err(err).Msg("Invalid request format")
		responses.RespondValidationError(c, err)
		return false
	}

	return true
}

// extractToken extracts and validates bearer token from Authorization header
func extractToken(c *gin.Context) (string, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Extract the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Warn().Msg("Missing Authorization header")
		responses.RespondUnauthorized(c, nil, "Missing Authorization header")
		return "", false
	}

	// Check if the Authorization header has the correct format
	bearerPrefix := "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		log.Warn().Msg("Invalid Authorization header format")
		responses.RespondUnauthorized(c, nil, "Invalid Authorization header format")
		return "", false
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		log.Warn().Msg("Empty token provided")
		responses.RespondUnauthorized(c, nil, "Empty token provided")
		return "", false
	}

	return token, true
}

// checkErrorType helps categorize and handle common API error types
func checkErrorType(c *gin.Context, err error, logMsg string) bool {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	if err == nil {
		return false
	}

	errStr := err.Error()

	// Handle common error types
	if strings.Contains(errStr, "not found") {
		log.Warn().Err(err).Msg(logMsg + " - not found")
		responses.RespondNotFound(c, err, errStr)
		return true
	}

	if strings.Contains(errStr, "already exists") {
		log.Warn().Err(err).Msg(logMsg + " - conflict")
		responses.RespondConflict(c, err, errStr)
		return true
	}

	if strings.Contains(errStr, "invalid") ||
		strings.Contains(errStr, "unauthorized") ||
		strings.Contains(errStr, "expired") {
		log.Warn().Err(err).Msg(logMsg + " - unauthorized")
		responses.RespondUnauthorized(c, err, errStr)
		return true
	}

	return false
}

// handleServiceError is a simple helper to check for common error patterns in service calls and respond appropriately
func handleServiceError(c *gin.Context, err error, logMsg string, notFoundMsg string, internalErrMsg string) bool {
	if err == nil {
		return false
	}

	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	if checkErrorType(c, err, logMsg) {
		return true
	}

	// For not found checks (sometimes error messages don't use "not found" pattern)
	if notFoundMsg != "" && err.Error() == notFoundMsg {
		log.Warn().Msg(logMsg + " - not found")
		responses.RespondNotFound(c, err, notFoundMsg)
		return true
	}

	// If no special handling was applicable, log as internal error
	log.Error().Err(err).Msg(logMsg + " - server error")
	responses.RespondInternalError(c, err, internalErrMsg)
	return true
}

// checkYear parses a year parameter from the context
// If the year is invalid, it responds with an error and returns 0, false
func checkYear(c *gin.Context, paramName string) (int, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	yearStr := c.Param(paramName)
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Warn().Err(err).Str(paramName, yearStr).Msg("Invalid year format")
		responses.RespondBadRequest(c, err, "Invalid year format")
		return 0, false
	}

	return year, true
}

// checkDaysParam parses the 'days' query parameter with a default value
func checkDaysParam(c *gin.Context, defaultDays int) int {
	daysParam := c.DefaultQuery("days", strconv.Itoa(defaultDays))
	if daysVal, err := strconv.Atoi(daysParam); err == nil {
		return daysVal
	}
	return defaultDays
}

// checkRating parses a rating parameter from the context
// If the rating is invalid, it responds with an error and returns 0, false
func checkRating(c *gin.Context, paramName string) (float64, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	ratingStr := c.Param(paramName)
	rating, err := strconv.ParseFloat(ratingStr, 32)
	if err != nil {
		log.Warn().Err(err).Str(paramName, ratingStr).Msg("Invalid rating value")
		responses.RespondBadRequest(c, err, "Invalid rating value")
		return 0, false
	}

	return rating, true
}

// checkRequiredStringParam validates that a string parameter is not empty
// If the parameter is empty, it responds with an error and returns false
func checkRequiredStringParam(c *gin.Context, paramName string, errorMsg string) (string, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	paramValue := c.Param(paramName)
	if paramValue == "" {
		log.Warn().Str(paramName, paramValue).Msg(errorMsg)
		responses.RespondBadRequest(c, nil, errorMsg)
		return "", false
	}

	return paramValue, true
}

// checkRequiredQueryParam validates that a query parameter is not empty
// If the parameter is empty, it responds with an error and returns false
func checkRequiredQueryParam(c *gin.Context, paramName string, errorMsg string) (string, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	paramValue := c.Query(paramName)
	if paramValue == "" {
		log.Warn().Str(paramName, paramValue).Msg(errorMsg)
		responses.RespondBadRequest(c, nil, errorMsg)
		return "", false
	}

	return paramValue, true
}

// checkSeasonNumber parses a season number parameter from the context
// If the season number is invalid, it responds with an error and returns 0, false
func checkSeasonNumber(c *gin.Context, paramName string) (int, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	seasonNumStr := c.Param(paramName)
	seasonNum, err := strconv.Atoi(seasonNumStr)
	if err != nil {
		log.Warn().Err(err).Str(paramName, seasonNumStr).Msg("Invalid season number")
		responses.RespondBadRequest(c, err, "Invalid season number")
		return 0, false
	}

	return seasonNum, true
}

func checkClientCategory(c *gin.Context) (clienttypes.ClientCategory, bool) {

	clientCategory, valid := checkOptionalClientCategory(c)
	if !valid {
		err := fmt.Errorf("invalid client category")
		responses.RespondBadRequest(c, err, "Invalid client category")
		return "", false
	}

	return clientCategory, true
}

func checkOptionalClientCategory(c *gin.Context) (clienttypes.ClientCategory, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientGroupStr := c.Query("clientCategory")
	clientGroup := clienttypes.ClientCategory(clientGroupStr)

	if clientGroup == "" {
		log.Debug().Msg("Client group not provided")
		return "", false
	}
	switch clientGroup {
	case clienttypes.ClientCategoryMedia, clienttypes.ClientCategoryAutomation, clienttypes.ClientCategoryAI, clienttypes.ClientCategoryMetadata:
		return clientGroup, true
	default:
		log.Warn().Msg("Invalid client group")
		responses.RespondBadRequest(c, nil, "Invalid client group")
		return "", false
	}

}

func checkMediaType(c *gin.Context) (types.MediaType, bool) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	mediaTypeStr := c.Param("mediaType")
	mediaType := types.MediaType(mediaTypeStr)

	if mediaType == "" {
		log.Debug().Msg("Media type not provided")
		return "", false
	}
	switch mediaType {
	case types.MediaTypeMovie, types.MediaTypeSeries, types.MediaTypeSeason, types.MediaTypeEpisode,
		types.MediaTypeArtist, types.MediaTypeAlbum, types.MediaTypeTrack, types.MediaTypePlaylist,
		types.MediaTypeCollection:
		return mediaType, true
	default:
		log.Warn().Msg("Invalid media type")
		responses.RespondBadRequest(c, nil, "Invalid media type")
		return "", false
	}

}
