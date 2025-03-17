package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

func RegisterUserConfigRoutes(rg *gin.RouterGroup, service services.UserConfigService) {
	configHandlers := handlers.NewUserConfigHandler(service)
	configs := rg.Group("/config/user")
	{

		configs.GET("", configHandlers.GetUserConfig)
		configs.PUT("", configHandlers.UpdateUserConfig)

	}
}
