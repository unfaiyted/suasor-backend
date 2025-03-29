// router/client.go
package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"suasor/client/media/types"
	"suasor/repository"
)

// SetupClientRoutes configures routes for client endpoints
func RegisterMediaItemRoutes(r *gin.RouterGroup, db *gorm.DB) {

	mediaItemRepo := repository.NewMediaItemRepository[types.MediaData](db)
	mediaItemService := services.NewMediaItemService[types.MediaData](mediaItemRepo)
	mediaItemHandler := handlers.NewMediaItemHandler[types.MediaData](mediaItemService)

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
		mediaItem.GET("/:clientId", mediaItemHandler.GetMediaItemsByClient)
		//mediaItem.GET("/:clientId/search", mediaItemHandler.SearchMediaItemsByClient)
	}

}
