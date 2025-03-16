// In handlers/health.go
package handlers

import (
	"net/http"
	"suasor/models"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service services.HealthService
}

func NewHealthHandler(service services.HealthService) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

// CheckHealth godoc
// @Summary checks app and database health
// @Description returns JSON object with health statuses.
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 500 {object} models.ErrorResponse[error]
// @Router /health [get]
func (h *HealthHandler) CheckHealth(c *gin.Context) {
	appStatus := h.service.CheckApplicationStatus()
	dbStatus := h.service.CheckDatabaseConnection()

	// Set overall status based on individual component statuses
	overallStatus := "up"
	httpStatus := http.StatusOK

	// Check if any component is down
	if !appStatus || !dbStatus {
		overallStatus = "down"
		httpStatus = http.StatusInternalServerError

		// Create error response
		errorResponse := models.ErrorResponse[models.HealthResponse]{
			Type: models.ErrorTypeFailedCheck,
			Details: models.HealthResponse{
				Status:      "down",
				Application: appStatus,
				Database:    dbStatus,
			},
		}

		c.JSON(httpStatus, errorResponse)
		return
	}

	// All components are healthy
	response := models.HealthResponse{
		Status:      overallStatus,
		Application: appStatus,
		Database:    dbStatus,
	}

	c.JSON(httpStatus, response)
}
