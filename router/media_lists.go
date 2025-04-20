package router

import (
	"suasor/app/container"
	"suasor/app/handlers"

	"github.com/gin-gonic/gin"
)

// SetupMediaListRoutes sets up the routes for media lists
func RegisterLocalMediaListRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Initialize handlers
	listHandlers := container.MustGet[handlers.MediaListHandlers](c)

	// Playlist routes
	playlists := rg.Group("/playlists")
	{
		// Get specialized playlist handler
		playlistHandler := listHandlers.CorePlaylistsHandler()
		userPlaylistHandler := listHandlers.UserPlaylistsHandler()

		// Basic CRUD operations
		playlists.GET("", playlistHandler.GetAll)
		playlists.GET("/:id", playlistHandler.GetByID)
		playlists.POST("", userPlaylistHandler.Create)
		playlists.PUT("/:id", userPlaylistHandler.Update)
		playlists.DELETE("/:id", userPlaylistHandler.Delete)

		// Playlist items management
		playlists.GET("/:id/items", playlistHandler.GetItemsByListID)
		playlists.POST("/:id/items", userPlaylistHandler.AddItem)
		playlists.DELETE("/:id/items/:itemId", userPlaylistHandler.Delete)

		// Playlist reordering
		playlists.POST("/:id/reorder", userPlaylistHandler.ReorderItems)

		// Search playlists
		playlists.GET("/search", playlistHandler.Search)

		// Sync playlist across clients

		// playlists.POST("/:id/sync", playlistHandler.Sync)
	}

	// Collections routes
	collections := rg.Group("/collections")
	{
		// Get specialized collection handler
		collectionHandler := listHandlers.CoreCollectionsHandler()
		userCollectionHandler := listHandlers.UserCollectionsHandler()

		// Basic CRUD operations
		collections.GET("", collectionHandler.GetAll)
		collections.GET("/:id", collectionHandler.GetByID)
		collections.POST("", userCollectionHandler.Create)
		collections.PUT("/:id", userCollectionHandler.Update)
		collections.DELETE("/:id", userCollectionHandler.Delete)

		// Collection items management
		collections.GET("/:id/items", collectionHandler.GetItemsByListID)

		// Special collection types
		// collections.GET("/smart", collectionHandler.GetSmart)
		// collections.GET("/featured", collectionHandler.GetFeatured)
	}
}
