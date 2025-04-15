// router/playlists.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
	"suasor/handlers"
)

// SetupPlaylistRoutes sets up the routes for playlist operations
func SetupPlaylistRoutes(rg *gin.RouterGroup, app *app.App) {
	// Initialize handlers
	corePlaylistHandler := handlers.NewCorePlaylistHandler(
		app.Services().MediaItemServices().CorePlaylistService(),
		app.Services().PlaylistServices(),
	)

	userPlaylistHandler := handlers.NewUserPlaylistHandler(
		app.Services().MediaItemServices().UserPlaylistService(),
		app.Services().PlaylistServices(),
	)

	// Core playlist routes (database-focused)
	playlists := rg.Group("/playlists")
	{
		// Basic operations
		playlists.GET("", corePlaylistHandler.GetAll)
		playlists.GET("/:id", corePlaylistHandler.GetByID)
		playlists.GET("/:id/tracks", corePlaylistHandler.GetPlaylistTracks)
		
		// Search and filtering
		playlists.GET("/search", corePlaylistHandler.Search)
		playlists.GET("/genre/:genre", corePlaylistHandler.GetByGenre)
	}

	// User-specific playlist routes
	userPlaylists := rg.Group("/user/playlists")
	{
		// Get user's playlists
		userPlaylists.GET("", userPlaylistHandler.GetUserPlaylists)
		
		// Playlist CRUD operations
		userPlaylists.POST("", userPlaylistHandler.CreatePlaylist)
		userPlaylists.PUT("/:id", userPlaylistHandler.UpdatePlaylist)
		userPlaylists.DELETE("/:id", userPlaylistHandler.DeletePlaylist)
		
		// Playlist track management
		userPlaylists.POST("/:id/tracks", userPlaylistHandler.AddTrackToPlaylist)
		userPlaylists.DELETE("/:id/tracks/:trackId", userPlaylistHandler.RemoveTrackFromPlaylist)
	}

	// Client-specific routes are handled in router/media.go
	// These routes follow the pattern /clients/media/{clientID}/playlists/...
}

// Helper functions for client playlists
func getPlaylistHandler[T interface{}](clientMedia *handlers.ClientMediaHandler[T]) *handlers.ClientMediaPlaylistHandler[T] {
	return handlers.NewClientMediaPlaylistHandler[T](clientMedia.PlaylistService())
}