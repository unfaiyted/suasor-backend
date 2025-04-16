package router

import (
	"suasor/handlers"

	"github.com/gin-gonic/gin"
	"suasor/app/container"
)

func RegisterUserConfigRoutes(rg *gin.RouterGroup, c *container.Container) {
	configHandlers := container.MustGet[handlers.UserConfigHandler](c)
	configs := rg.Group("/config/user")
	{

		configs.GET("", configHandlers.GetUserConfig)
		configs.PUT("", configHandlers.UpdateUserConfig)

	}
}
