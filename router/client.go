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
	embyService := services.NewClientService[*types.EmbyConfig](db)
	jellyfinService := services.NewClientService[*types.JellyfinConfig](db)
	lidarrService := services.NewClientService[*types.LidarrConfig](db)

	downloadClientHandler := handlers.NewClientHandler[types.AutomationClientConfig](downloadService)
	aiClientHandler := handlers.NewClientHandler[types.AIClientConfig](aiService)
	mediaClientHandler := handlers.NewClientHandler[types.MediaClientConfig](mediaService)

	embyClientHandler := handlers.NewClientHandler[*types.EmbyConfig](embyService)
	jellyfinClientHandler := handlers.NewClientHandler[*types.JellyfinConfig](jellyfinService)
	lidarrClientHandler := handlers.NewClientHandler[*types.LidarrConfig](lidarrService)

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
	emby := clients.Group("/emby")
	{
		emby.POST("", embyClientHandler.CreateClient)
		emby.GET("", embyClientHandler.GetAllClients)
		emby.GET("/:id", embyClientHandler.GetClient)
		emby.DELETE("/:id", embyClientHandler.DeleteClient)
		emby.PUT("/:id", embyClientHandler.UpdateClient)
		emby.POST("/test", embyClientHandler.TestConnection)
	}
	jellyfin := clients.Group("/jellyfin")
	{
		jellyfin.POST("", jellyfinClientHandler.CreateClient)
		jellyfin.GET("", jellyfinClientHandler.GetAllClients)
		jellyfin.GET("/:id", jellyfinClientHandler.GetClient)
		jellyfin.DELETE("/:id", jellyfinClientHandler.DeleteClient)
		jellyfin.PUT("/:id", jellyfinClientHandler.UpdateClient)
		jellyfin.POST("/test", jellyfinClientHandler.TestConnection)
	}
	lidarr := clients.Group("/lidarr")
	{
		lidarr.POST("", lidarrClientHandler.CreateClient)
		lidarr.GET("", lidarrClientHandler.GetAllClients)
		lidarr.GET("/:id", lidarrClientHandler.GetClient)
		lidarr.DELETE("/:id", lidarrClientHandler.DeleteClient)
		lidarr.PUT("/:id", lidarrClientHandler.UpdateClient)
		lidarr.POST("/test", lidarrClientHandler.TestConnection)
	}

}
