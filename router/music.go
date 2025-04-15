// router/music.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
	"suasor/handlers"
)

// SetupMusicRoutes sets up the routes for music-related operations
func SetupMusicRoutes(rg *gin.RouterGroup, app *app.App) {
	// Initialize handlers
	coreMusicHandler := handlers.NewCoreMusicHandler(
		app.Services().MediaItemServices().CoreTrackService(),
		app.Services().MediaItemServices().CoreAlbumService(),
		app.Services().MediaItemServices().CoreArtistService(),
	)

	userMusicHandler := handlers.NewUserMusicHandler(
		app.Services().MediaItemServices().UserTrackService(),
		app.Services().MediaItemServices().UserAlbumService(),
		app.Services().MediaItemServices().UserArtistService(),
	)

	// Core music routes (database-focused)
	// Tracks
	tracks := rg.Group("/music/tracks")
	{
		tracks.GET("/top", coreMusicHandler.GetTopTracks)
		tracks.GET("/recently-added", coreMusicHandler.GetRecentlyAddedTracks)
	}

	// Albums
	albums := rg.Group("/music/albums")
	{
		albums.GET("/:id/tracks", coreMusicHandler.GetAlbumTracks)
		albums.GET("/top", coreMusicHandler.GetTopAlbums)
	}

	// Artists
	artists := rg.Group("/music/artists")
	{
		artists.GET("/:id/albums", coreMusicHandler.GetArtistAlbums)
	}

	// User-specific music routes
	// User tracks
	userTracks := rg.Group("/user/music/tracks")
	{
		userTracks.GET("/favorites", userMusicHandler.GetFavoriteTracks)
		userTracks.GET("/recently-played", userMusicHandler.GetRecentlyPlayedTracks)
		userTracks.PATCH("/:id", userMusicHandler.UpdateTrackUserData)
	}

	// User albums
	userAlbums := rg.Group("/user/music/albums")
	{
		userAlbums.GET("/favorites", userMusicHandler.GetFavoriteAlbums)
	}

	// User artists
	userArtists := rg.Group("/user/music/artists")
	{
		userArtists.GET("/favorites", userMusicHandler.GetFavoriteArtists)
	}

	// Client-specific routes are handled in router/media.go
	// These routes follow the pattern /clients/media/{clientID}/music/...
}

// Music-specific helper functions
func getMusicHandler[T interface{}](clientMedia *handlers.ClientMediaHandler[T]) *handlers.ClientMediaMusicHandler[T] {
	return handlers.NewClientMediaMusicHandler[T](clientMedia.MusicService())
}