// handlers/config.go
package handlers

import (
	"net/http"
	"suasor/models"
	"suasor/services"

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

// GetConfig returns the current configuration
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, h.configService.GetConfig())
}

// UpdateConfig handles configuration updates
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var cfg models.Configuration
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.configService.SaveConfig(c.Request.Context(), cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// ResetConfig resets the configuration to default values
func (h *ConfigHandler) ResetConfig(c *gin.Context) {
	if err := h.configService.ResetFileConfig(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
