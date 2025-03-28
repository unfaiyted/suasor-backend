// router/client.go
package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"suasor/client/types"
)

// SetupClientRoutes configures routes for client endpoints
func RegisterClientRoutes(r *gin.RouterGroup, db *gorm.DB) {

	downloadService := services.NewClientService[types.AutomationClientConfig](db)
	mediaService := services.NewClientService[types.MediaClientConfig](db)

	downloadClientHandler := handlers.NewClientHandler[types.AutomationClientConfig](downloadService)

	clients := r.Group("/clients")

	// Download client routes
	download := clients.Group("/download")
	{
		download.POST("", downloadClientHandler.CreateClient)
		download.GET("", downloadClientHandler.GetAllClients)
		download.GET("/:id", downloadClientHandler.GetClient)
		download.PUT("/:id", downloadClientHandler.UpdateClient)
		download.DELETE("/:id", downloadClientHandler.DeleteClient)
		download.POST("/test", downloadClientHandler.TestConnection)
	}

	mediaClientHandler := handlers.NewClientHandler(mediaService)

	media := clients.Group("/media")
	{

		media.POST("", mediaClientHandler.CreateClient)
		media.GET("", mediaClientHandler.GetAllClients)
		media.GET("/:id", mediaClientHandler.GetClient)
		media.DELETE("/:id", mediaClientHandler.DeleteClient)
		media.PUT("/:id", mediaClientHandler.UpdateClient)
		media.POST("/test", mediaClientHandler.TestConnection)

	}

}
