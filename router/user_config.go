package router

import (
	"suasor/handlers"

	"github.com/gin-gonic/gin"
	"suasor/app"
)

func RegisterUserConfigRoutes(rg *gin.RouterGroup, deps *app.AppDependencies) {
	configHandlers := handlers.NewUserConfigHandler(deps.UserConfigService())
	configs := rg.Group("/config/user")
	{

		configs.GET("", configHandlers.GetUserConfig)
		configs.PUT("", configHandlers.UpdateUserConfig)

	}
}
