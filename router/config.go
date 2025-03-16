package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

func RegisterConfigRoutes(rg *gin.RouterGroup, service services.ConfigService) {
	configHandlers := handlers.NewConfigHandler(service)
	configs := rg.Group("/config")
	{

		configs.GET("", configHandlers.GetConfig)
		configs.PUT("", configHandlers.UpdateConfig)
		configs.POST("/reset", configHandlers.ResetConfig)

	}
}
