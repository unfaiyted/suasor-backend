package router

import (
	"suasor/app/container"
	"suasor/app/handlers"

	"github.com/gin-gonic/gin"
)

// SetupMediaListRoutes sets up the routes for media lists
func SetupMediaListRoutes(rg *gin.RouterGroup, c *container.Container) {
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
		playlists.POST("", userPlaylistHandler.CreatePlaylist)
		playlists.PUT("/:id", playlistHandler.Update)
		playlists.DELETE("/:id", playlistHandler.Delete)

		// Playlist items management
		playlists.GET("/:id/items", playlistHandler.GetItems)
		playlists.POST("/:id/items", playlistHandler.AddItem)
		playlists.DELETE("/:id/items/:itemId", playlistHandler.RemoveItem)

		// Playlist reordering
		playlists.POST("/:id/reorder", playlistHandler.ReorderItems)

		// Search playlists
		playlists.GET("/search", playlistHandler.Search)

		// Sync playlist across clients
		playlists.POST("/:id/sync", playlistHandler.Sync)
	}

	// Collections routes
	collections := rg.Group("/collections")
	{
		// Get specialized collection handler
		collectionHandler := listHandlers.CoreCollectionsHandler()

		// Basic CRUD operations
		collections.GET("", collectionHandler.GetAll)
		collections.GET("/:id", collectionHandler.GetByID)

		// Collection items management
		collections.GET("/:id/items", collectionHandler.GetItems)

		// Special collection types
		collections.GET("/smart", collectionHandler.GetSmart)
		collections.GET("/featured", collectionHandler.GetFeatured)
	}
}
