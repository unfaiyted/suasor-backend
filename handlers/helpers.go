package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
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

	// Check if clientType is a valid string
	clientTypeStr, exists := c.Get("clientType")
	log.Debug().
		Str("clientType", clientTypeStr.(string)).
		Bool("exists", exists).
		Msg("Checking client type")
	if !exists {
		log.Warn().Msg("Client type not found in context")
		responses.RespondBadRequest(c, nil, "Client type not found in context")
		return "", false
	}
	clientType := clienttypes.ClientType(clientTypeStr.(string))
	switch clientType {
	case clienttypes.ClientTypeEmby, clienttypes.ClientTypeJellyfin, clienttypes.ClientTypePlex,
		clienttypes.ClientTypeSubsonic, clienttypes.ClientTypeSonarr, clienttypes.ClientTypeLidarr,
		clienttypes.ClientTypeRadarr, clienttypes.ClientTypeClaude, clienttypes.ClientTypeOpenAI,
		clienttypes.ClientTypeOllama:
		log.Debug().
			Str("clientType", clientTypeStr.(string)).
			Msg("Client type valid")
		return clientType, true
	default:
		log.Warn().Str("clientType", clientTypeStr.(string)).Msg("Invalid client type")
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
