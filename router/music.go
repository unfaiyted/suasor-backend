// router/music.go
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	"suasor/handlers"
	"suasor/types/responses"
)

type ClientMusicHandlerInterface interface {
	GetTopTracks(c *gin.Context)
	GetRecentlyAddedTracks(c *gin.Context)
	GetTopAlbums(c *gin.Context)
	GetTopArtists(c *gin.Context)
	GetFavoriteArtists(c *gin.Context)
}

// SetupMusicRoutes sets up the routes for music-related operations
func RegisterMusicRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Initialize handlers
	coreMusicHandler := container.MustGet[handlers.CoreMusicHandler](c)
	userMusicHandler := container.MustGet[handlers.UserMusicHandler](c)
	clientMusicHandler := container.MustGet[apphandlers.ClientMusicHandlers](c)

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

	clientHandlerMap := map[string]ClientMusicHandlerInterface{
		"emby":     clientMusicHandler.EmbyMusicHandler(),
		"jellyfin": clientMusicHandler.JellyfinMusicHandler(),
		"plex":     clientMusicHandler.PlexMusicHandler(),
		"subsonic": clientMusicHandler.SubsonicMusicHandler(),
	}

	getClientHandler := func(c *gin.Context) ClientMusicHandlerInterface {
		clientType := c.Param("clientType")
		handler, exists := clientHandlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}
		return handler
	}

	// These routes follow the pattern /clients/media/{clientID}/music/...
	clientMusic := rg.Group("/clients/:clientType/:clientID/music")
	{
		// Tracks
		clientMusic.GET("/tracks/top", func(c *gin.Context) {
			getClientHandler(c).GetTopTracks(c)
		})

		clientMusic.GET("/tracks/recently-added", func(c *gin.Context) {
			getClientHandler(c).GetRecentlyAddedTracks(c)
		})

		// Albums
		clientMusic.GET("/albums/top", func(c *gin.Context) {
			getClientHandler(c).GetTopAlbums(c)
		})

		// Artists
		clientMusic.GET("/artists/top", func(c *gin.Context) {
			getClientHandler(c).GetTopArtists(c)
		})
		clientMusic.GET("/artists/favorites", func(c *gin.Context) {
			getClientHandler(c).GetFavoriteArtists(c)
		})
	}

}

