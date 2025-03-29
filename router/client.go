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
	aiService := services.NewClientService[types.AIClientConfig](db)

	downloadClientHandler := handlers.NewClientHandler[types.AutomationClientConfig](downloadService)
	aiClientHandler := handlers.NewClientHandler[types.AIClientConfig](aiService)
	mediaClientHandler := handlers.NewClientHandler[types.MediaClientConfig](mediaService)

	clients := r.Group("/client")

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

	media := clients.Group("/media")
	{

		media.POST("", mediaClientHandler.CreateClient)
		media.GET("", mediaClientHandler.GetAllClients)
		media.GET("/:id", mediaClientHandler.GetClient)
		media.DELETE("/:id", mediaClientHandler.DeleteClient)
		media.PUT("/:id", mediaClientHandler.UpdateClient)
		media.POST("/test", mediaClientHandler.TestConnection)

	}
	ai := clients.Group("/ai")
	{
		ai.POST("", aiClientHandler.CreateClient)
		ai.GET("", aiClientHandler.GetAllClients)
		ai.GET("/:id", aiClientHandler.GetClient)
		ai.DELETE("/:id", aiClientHandler.DeleteClient)
		ai.PUT("/:id", aiClientHandler.UpdateClient)
		ai.POST("/test", aiClientHandler.TestConnection)
	}

}
