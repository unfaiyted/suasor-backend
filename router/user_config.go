package router

import (
	"suasor/handlers"

	"github.com/gin-gonic/gin"
	"suasor/di/container"
)

func RegisterUserConfigRoutes(rg *gin.RouterGroup, c *container.Container) {
	configHandlers := container.MustGet[*handlers.UserConfigHandler](c)
	configs := rg.Group("/user-config")
	{

		configs.GET("", configHandlers.GetUserConfig)
		configs.PUT("", configHandlers.UpdateUserConfig)

	}
}
