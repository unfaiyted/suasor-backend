// handlers/config.go
package handlers

import (
	"suasor/services"
	"suasor/types"
	"suasor/types/responses"

	"suasor/utils"

	"github.com/gin-gonic/gin"
)

// ConfigHandler handles configuration API endpoints
type ConfigHandler struct {
	configService services.ConfigService
}

// NewConfigHandler creates a new configuration handler
func NewConfigHandler(configService services.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// checkAdminAccess verifies if the request is from an admin user
func (h *ConfigHandler) checkAdminAccess(c *gin.Context) (uint64, bool) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

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

// GetConfig godoc
// @Summary Get current configuration
// @Description Returns the current system configuration
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} responses.APIResponse[types.Configuration] "Configuration retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized access"
// @Failure 403 {object} responses.ErrorResponse[responses.ErrorDetails] "Forbidden - admin access required"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /config [get]
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Check if user is authenticated and has admin role
	_, ok := h.checkAdminAccess(c)
	if !ok {
		return
	}

	log.Info().Msg("Retrieving system configuration")

	config := h.configService.GetConfig()

	log.Info().Msg("Configuration retrieved successfully")
	responses.RespondOK(c, config, "Configuration retrieved successfully")
}

// UpdateConfig godoc
// @Summary Update application configuration
// @Description Updates the system-wide application configuration (admin only)
// @Tags config
// @Accept json
// @Produce json
// @Param request body types.Configuration true "Configuration data"
// @Success 200 {object} responses.APIResponse[any] "Configuration updated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized access"
// @Failure 403 {object} responses.ErrorResponse[responses.ErrorDetails] "Forbidden - admin access required"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /config [put]
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Check if user is authenticated and has admin role
	userID, ok := h.checkAdminAccess(c)
	if !ok {
		return
	}

	var cfg types.Configuration
	if err := c.ShouldBindJSON(&cfg); err != nil {
		log.Error().Err(err).Msg("Invalid configuration format")
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Interface("userID", userID).
		Msg("Updating application configuration")

	if err := h.configService.SaveConfig(ctx, cfg); err != nil {
		log.Error().Err(err).Msg("Failed to update configuration")
		responses.RespondInternalError(c, err, "Failed to update configuration")
		return
	}

	log.Info().
		Interface("userID", userID).
		Msg("Application configuration updated successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Configuration updated successfully")
}

// GetFileConfig godoc
// @Summary Get file-based configuration
// @Description Returns the file-based system configuration (admin only)
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} responses.APIResponse[types.Configuration] "File configuration retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized access"
// @Failure 403 {object} responses.ErrorResponse[responses.ErrorDetails] "Forbidden - admin access required"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /config/file [get]
// GetFileConfig uses the admin access check helper
func (h *ConfigHandler) GetFileConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, ok := h.checkAdminAccess(c)
	if !ok {
		return
	}

	log.Info().
		Uint64("userID", userID).
		Msg("Retrieving file-based configuration")

	config := h.configService.GetFileConfig(ctx)
	if config == nil {
		log.Error().Msg("Failed to retrieve file configuration")
		responses.RespondInternalError(c, nil, "Failed to retrieve file configuration")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Msg("File configuration retrieved successfully")

	responses.RespondOK(c, config, "File configuration retrieved successfully")
}

// Similar refactoring should be applied to SaveFileConfig, UpdateConfig, and ResetConfig

// SaveFileConfig godoc
// @Summary Save configuration to file
// @Description Saves the configuration to file only (admin only)
// @Tags config
// @Accept json
// @Produce json
// @Param request body types.Configuration true "Configuration data"
// @Success 200 {object} responses.APIResponse[any] "Configuration saved to file successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized access"
// @Failure 403 {object} responses.ErrorResponse[responses.ErrorDetails] "Forbidden - admin access required"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /config/file [put]
func (h *ConfigHandler) SaveFileConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Check if user is authenticated and has admin role
	userID, ok := h.checkAdminAccess(c)
	if !ok {
		return
	}

	var cfg types.Configuration
	if err := c.ShouldBindJSON(&cfg); err != nil {
		log.Error().Err(err).Msg("Invalid configuration format")
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Interface("userID", userID).
		Msg("Saving configuration to file")

	if err := h.configService.SaveFileConfig(ctx, cfg); err != nil {
		log.Error().Err(err).Msg("Failed to save configuration to file")
		responses.RespondInternalError(c, err, "Failed to save configuration to file")
		return
	}

	log.Info().
		Interface("userID", userID).
		Msg("Configuration saved to file successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Configuration saved to file successfully")
}

// ResetConfig godoc
// @Summary Reset configuration to defaults
// @Description Resets the system configuration to default values (admin only)
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} responses.APIResponse[any] "Configuration reset successfully"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized access"
// @Failure 403 {object} responses.ErrorResponse[responses.ErrorDetails] "Forbidden - admin access required"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /config/reset [post]
func (h *ConfigHandler) ResetConfig(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, ok := h.checkAdminAccess(c)
	if !ok {
		return
	}

	log.Info().
		Interface("userID", userID).
		Msg("Resetting application configuration to defaults")

	if err := h.configService.ResetFileConfig(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to reset configuration")
		responses.RespondInternalError(c, err, "Failed to reset configuration")
		return
	}

	log.Info().
		Interface("userID", userID).
		Msg("Application configuration reset successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Configuration reset successfully")
}
