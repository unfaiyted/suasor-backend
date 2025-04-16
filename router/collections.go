// router/collections.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
	"suasor/handlers"
)

// SetupCollectionRoutes sets up the routes for collection operations
func SetupCollectionRoutes(rg *gin.RouterGroup, app *app.App) {
	// Initialize handlers
	coreCollectionHandler := handlers.NewCoreCollectionHandler(
		app.Services().MediaItemServices().CoreCollectionService(),
		app.Services().CollectionServices(),
	)

	userCollectionHandler := handlers.NewUserCollectionHandler(
		app.Services().MediaItemServices().UserCollectionService(),
		app.Services().CollectionServices(),
	)

	// Core collection routes (database-focused)
	collections := rg.Group("/collections")
	{
		// Basic operations
		collections.GET("", coreCollectionHandler.GetAll)
		collections.GET("/:id", coreCollectionHandler.GetByID)
		collections.GET("/:id/items", coreCollectionHandler.GetCollectionItems)

		// Search and filtering
		collections.GET("/search", coreCollectionHandler.Search)
		collections.GET("/genre/:genre", coreCollectionHandler.GetByGenre)
		collections.GET("/public", coreCollectionHandler.GetPublicCollections)
	}

	// User-specific collection routes
	userCollections := rg.Group("/user/collections")
	{
		// Get user's collections
		userCollections.GET("", userCollectionHandler.GetUserCollections)

		// Collection CRUD operations
		userCollections.POST("", userCollectionHandler.CreateCollection)
		userCollections.PUT("/:id", userCollectionHandler.UpdateCollection)
		userCollections.DELETE("/:id", userCollectionHandler.DeleteCollection)

		// Collection item management
		userCollections.POST("/:id/items", userCollectionHandler.AddItemToCollection)
		userCollections.DELETE("/:id/items/:itemId", userCollectionHandler.RemoveItemFromCollection)
	}

	// Client-specific routes are handled in router/media.go
	// These routes follow the pattern /clients/media/{clientID}/collections/...
}

// Helper functions for client collections
func getCollectionHandler[T interface{}](clientMedia *handlers.ClientMediaHandler[T]) *handlers.ClientMediaCollectionHandler[T] {
	return handlers.NewClientMediaCollectionHandler[T](clientMedia.CollectionService())
}

