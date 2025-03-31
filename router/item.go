// router/client.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
)

// SetupClientRoutes configures routes for client endpoints
func RegisterMediaItemRoutes(r *gin.RouterGroup, deps *app.AppDependencies) {

	mediaItemHandler := deps.MediaHandlers.MovieHandler()

	clients := r.Group("/item")

	// mediaItem client routes
	mediaItem := clients.Group("/media")
	{
		mediaItem.POST("", mediaItemHandler.CreateMediaItem)
		// mediaItem.GET("", mediaItemHandler.GetMediaItem)
		mediaItem.GET("/:id", mediaItemHandler.GetMediaItem)
		mediaItem.PUT("/:id", mediaItemHandler.UpdateMediaItem)
		mediaItem.DELETE("/:id", mediaItemHandler.DeleteMediaItem)
		mediaItem.GET("/search", mediaItemHandler.SearchMediaItems)
		mediaItem.GET("/recent", mediaItemHandler.GetRecentMediaItems)
		mediaItem.GET("/client/:clientId", mediaItemHandler.GetMediaItemsByClient)
		//mediaItem.GET("/:clientId/search", mediaItemHandler.SearchMediaItemsByClient)
	}

}
