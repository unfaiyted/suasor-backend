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
	listGroup.GET("/:listID", coreHandler.GetByID)
	listGroup.GET("/:listID/items", coreHandler.GetItemsByListID)
	listGroup.GET("/search", coreHandler.Search)

	// User-specific operations
	listGroup.POST("", userHandler.Create)
	listGroup.PUT("/:listID", userHandler.Update)
	listGroup.DELETE("/:listID", userHandler.Delete)
	listGroup.GET("/user", userHandler.GetUserLists)
	listGroup.GET("/user/:userID", userHandler.GetUserLists)
	listGroup.POST("/:listID/item/:itemID", userHandler.AddItem)
	// Delets all instances of an item from a list
	listGroup.DELETE("/:listID/item/:itemID", userHandler.RemoveItem)

	// Type-specific operations based on list type
	if mediaType == mediatypes.MediaTypePlaylist {
		// Playlist-specific routes
		listGroup.POST("/:listID/reorder", userHandler.ReorderItems)
		// Delete item at specific position
		listGroup.DELETE("/:listID/item/:itemID/position/:position", userHandler.RemoveItemAtPosition)

	} else if mediaType == mediatypes.MediaTypeCollection {
		// listGroup.DELETE("/:listID/item/:itemID", userHandler.Delete)
		// Collection-specific routes (can be extended as needed)
	}
}
