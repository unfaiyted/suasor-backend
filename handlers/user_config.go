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
// @Summary Get user configuration
// @Description Returns the configuration for the current user
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} responses.APIResponse[models.UserConfig] "User configuration retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized access"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /config/user [get]
func (h *UserConfigHandler) GetUserConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Unauthorized attempt to get user configuration")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	log.Info().Uint64("userID", uid).Msg("Retrieving user configuration")

	config, err := h.userConfigService.GetUserConfig(ctx, uid)
	if err != nil {
		log.Error().Err(err).Uint64("userID", uid).Msg("Failed to retrieve user configuration")
		responses.RespondInternalError(c, err, "Failed to retrieve user configuration")
		return
	}

	log.Info().Uint64("userID", uid).Msg("User configuration retrieved successfully")
	responses.RespondOK(c, config, "User configuration retrieved successfully")
}

// UpdateUserConfig godoc
// @Summary Update user configuration
// @Description Updates the configuration for the current user
// @Tags config
// @Accept json
// @Produce json
// @Param request body models.UserConfig true "User configuration data"
// @Success 200 {object} responses.APIResponse[any] "User configuration updated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized access"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /config/user [put]
func (h *UserConfigHandler) UpdateUserConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Unauthorized attempt to update user configuration")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	var cfg models.UserConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		log.Error().Err(err).Msg("Invalid user configuration format")
		responses.RespondValidationError(c, err)
		return
	}

	// Ensure the user can only modify their own config
	cfg.UserID = uid

	log.Info().Uint64("userID", uid).Msg("Updating user configuration")

	if err := h.userConfigService.SaveUserConfig(ctx, cfg); err != nil {
		log.Error().Err(err).Uint64("userID", uid).Msg("Failed to update user configuration")
		responses.RespondInternalError(c, err, "Failed to update user configuration")
		return
	}

	log.Info().Uint64("userID", uid).Msg("User configuration updated successfully")
	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "User configuration updated successfully")
}
