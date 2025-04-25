// handlers/config.go
package handlers

import (
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

// UserConfigHandler handles configuration API endpoints
type UserConfigHandler struct {
	userConfigService services.UserConfigService
}

// NewConfigHandler creates a new configuration handler
func NewUserConfigHandler(userConfigService services.UserConfigService) *UserConfigHandler {
	return &UserConfigHandler{
		userConfigService: userConfigService,
	}
}

// GetUserConfig godoc
//
//	@Summary		Get user configuration
//	@Description	Returns the configuration for the current user
//	@Tags			config
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.APIResponse[models.UserConfig]		"User configuration retrieved successfully"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized access"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user-config [get]
func (h *UserConfigHandler) GetUserConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)

	config, err := h.userConfigService.GetUserConfig(ctx, userID)
	if handleServiceError(c, err, "Failed to retrieve user configuration", "", "Failed to retrieve user configuration") {
		return
	}

	log.Info().Uint64("userID", userID).Msg("User configuration retrieved successfully")
	responses.RespondOK(c, config, "User configuration retrieved successfully")
}

// UpdateUserConfig godoc
//
//	@Summary		Update user configuration
//	@Description	Updates the configuration for the current user
//	@Tags			config
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.UserConfig								true	"User configuration data"
//	@Success		200		{object}	responses.APIResponse[any]						"User configuration updated successfully"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid request format"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized access"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user-config [put]
func (h *UserConfigHandler) UpdateUserConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)

	var cfg models.UserConfig
	if !checkJSONBinding(c, &cfg) {
		return
	}

	// Ensure the user can only modify their own config
	cfg.UserID = userID

	log.Info().Uint64("userID", userID).Msg("Updating user configuration")

	err := h.userConfigService.SaveUserConfig(ctx, cfg)
	if handleServiceError(c, err, "Failed to update user configuration", "", "Failed to update user configuration") {
		return
	}

	log.Info().Uint64("userID", userID).Msg("User configuration updated successfully")
	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "User configuration updated successfully")
}
