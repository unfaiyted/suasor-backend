package router

import (
	"suasor/di/container"
	"suasor/handlers"

	"github.com/gin-gonic/gin"
	mediatypes "suasor/clients/media/types"
)

// RegisterLocalMediaListRoutes sets up the routes for media lists
func RegisterLocalMediaListRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Register generic list routes for different media list types
	registerGenericListRoutes[*mediatypes.Playlist](c, rg)
	registerGenericListRoutes[*mediatypes.Collection](c, rg)
}

// registerGenericListRoutes sets up routes common to all media list types
func registerGenericListRoutes[T mediatypes.ListData](c *container.Container, rg *gin.RouterGroup) {
	// Get specialized handlers
	coreHandler := container.MustGet[handlers.CoreListHandler[T]](c)
	userHandler := container.MustGet[handlers.UserListHandler[T]](c)

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	// Create list group
	listGroup := rg.Group("/" + string(mediaType))

	// Basic CRUD operations
	listGroup.GET("", coreHandler.GetAll)
	listGroup.GET("/:listId", coreHandler.GetByID)
	listGroup.GET("/:listId/items", coreHandler.GetItemsByListID)
	listGroup.GET("/search", coreHandler.Search)

	// User-specific operations
	listGroup.POST("", userHandler.Create)
	listGroup.PUT("/:listId", userHandler.Update)
	listGroup.DELETE("/:listId", userHandler.Delete)
	listGroup.GET("/user", userHandler.GetUserLists)
	listGroup.POST("/:listId/items", userHandler.AddItem)
	listGroup.DELETE("/:listId/items/:itemId", userHandler.Delete)

	// Type-specific operations based on list type
	if mediaType == mediatypes.MediaTypePlaylist {
		// Playlist-specific routes
		listGroup.POST("/:listId/reorder", userHandler.ReorderItems)
	} else if mediaType == mediatypes.MediaTypeCollection {
		// Collection-specific routes (can be extended as needed)
	}
}
