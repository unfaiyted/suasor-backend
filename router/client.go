// router/client.go
package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

// SetupClientRoutes configures routes for client endpoints
func RegisterClientRoutes(r *gin.RouterGroup, downloadService services.DownloadClientService, mediaService services.MediaClientService) {
	downloadClientHandler := handlers.NewDownloadClientHandler(downloadService)

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

	mediaClientHandler := handlers.NewMediaClientHandler(mediaService)

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
