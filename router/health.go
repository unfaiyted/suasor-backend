package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(rg *gin.RouterGroup, service services.HealthService) {
	healthHandlers := handlers.NewHealthHandler(service)

	// Create a health endpoint
	rg.GET("/health", healthHandlers.CheckHealth)
}
