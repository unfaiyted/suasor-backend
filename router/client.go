// router/client.go
package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

// SetupClientRoutes configures routes for client endpoints
func RegisterClientRoutes(r *gin.RouterGroup, service services.DownloadClientService) {
	downloadClientHandler := handlers.NewDownloadClientHandler(service)

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

}
