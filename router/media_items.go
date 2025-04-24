package router

import (
	mediatypes "suasor/clients/media/types"
	"suasor/di/container"
	"suasor/handlers"

	"context"
	"github.com/gin-gonic/gin"
)

// RegisterLocalMediaItemRoutes configures routes for direct media item access
func RegisterLocalMediaItemRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {

	// Register generic media routes that work across all types
	registerGenericMediaRoutes[*mediatypes.Movie](c, rg)
	registerGenericMediaRoutes[*mediatypes.Series](c, rg)
	registerGenericMediaRoutes[*mediatypes.Season](c, rg)
	registerGenericMediaRoutes[*mediatypes.Episode](c, rg)
	registerGenericMediaRoutes[*mediatypes.Track](c, rg)
	registerGenericMediaRoutes[*mediatypes.Album](c, rg)
	registerGenericMediaRoutes[*mediatypes.Artist](c, rg)

	// Register specialized routes for series
	registerSeriesRoutes(c, rg)

	// Register specialized routes for music (tracks, albums, artists)
	registerMusicRoutes(c, rg)
}

// registerGenericMediaRoutes sets up routes common to all media types
func registerGenericMediaRoutes[T mediatypes.MediaData](c *container.Container, rg *gin.RouterGroup) {

	coreHandler := container.MustGet[handlers.CoreMediaItemHandler[T]](c)
	userHandler := container.MustGet[handlers.UserMediaItemHandler[T]](c)

	// Create media group i.e (/movies, /series, etc.)
	// full route is now /api/v1/media/{mediaType}
	mediaGroup := rg.Group("/" + coreHandler.GetType())

	// Register the routes
	mediaGroup.GET("", coreHandler.GetAll)
	mediaGroup.GET("/:itemID", coreHandler.GetByID)
	mediaGroup.GET("/genre/:genre", coreHandler.GetByGenre)
	mediaGroup.GET("/year/:year", coreHandler.GetByYear)
	mediaGroup.GET("/popular", coreHandler.GetPopular)
	mediaGroup.GET("/latest", coreHandler.GetLatestByAdded)
	mediaGroup.GET("/top-rated", coreHandler.GetTopRated)
	mediaGroup.GET("/search", coreHandler.Search)
	mediaGroup.GET("/external/:source/:externalID", coreHandler.GetByExternalID)
	mediaGroup.GET("/role/:role/:personID", coreHandler.GetByPerson)

	// User routes
	mediaGroup.GET("/user/:userID", userHandler.GetByUserID)
	mediaGroup.PUT("/:itemID", userHandler.Update)
	mediaGroup.DELETE("/:itemID", userHandler.Delete)
	mediaGroup.POST("", userHandler.Create)

}

// registerSeriesRoutes sets up routes specific to series
func registerSeriesRoutes(c *container.Container, rg *gin.RouterGroup) {
	coreSeriesHandler := container.MustGet[handlers.CoreSeriesHandler](c)
	coreUserHandler := container.MustGet[handlers.UserSeriesHandler](c)

	// Series-specific routes
	seriesGroup := rg.Group("/series")

	seriesGroup.GET("/:itemID/seasons", coreSeriesHandler.GetSeasonsBySeriesID)
	seriesGroup.GET("/:itemID/seasons/:seasonNumber/episodes", coreSeriesHandler.GetEpisodesBySeriesIDAndSeasonNumber)
	seriesGroup.GET("/:itemID/episodes", coreSeriesHandler.GetAllEpisodes)
	seriesGroup.GET("/continue-watching", coreUserHandler.GetNextUpEpisodes)
	seriesGroup.GET("/next-up", coreUserHandler.GetNextUpEpisodes)
	seriesGroup.GET("/recently-aired", coreSeriesHandler.GetRecentlyAiredEpisodes)
	seriesGroup.GET("/network/:network", coreSeriesHandler.GetSeriesByNetwork)
}

// registerMusicRoutes sets up routes specific to music (tracks, albums, artists)
func registerMusicRoutes(c *container.Container, rg *gin.RouterGroup) {

	// Core Music specific routes
	coreHandler := container.MustGet[handlers.CoreMusicHandler](c)

	// Media Item specific routes
	albumHandler := container.MustGet[handlers.CoreMediaItemHandler[*mediatypes.Album]](c)
	artistHandler := container.MustGet[handlers.CoreMediaItemHandler[*mediatypes.Artist]](c)
	trackHandler := container.MustGet[handlers.CoreMediaItemHandler[*mediatypes.Track]](c)

	// Main music group
	musicGroup := rg.Group("/music")
	{
		artistsGroup := musicGroup.Group("/artists")
		{
			artistsGroup.GET("/:itemID", artistHandler.GetByID)
			artistsGroup.GET("/:itemID/albums", coreHandler.GetAlbumsByArtistID)
			artistsGroup.GET("/:itemID/similar", coreHandler.GetSimilarArtists)
			artistsGroup.GET("/genre/:genre", artistHandler.GetByGenre)
			artistsGroup.GET("/popular", artistHandler.GetPopular)
			// artistsGroup.GET("/external/:source/:externalID", artistHandler.GetMediaItemByExternalSourceID)
		}
		albumsGroup := musicGroup.Group("/albums")
		{
			albumsGroup.GET("/:itemID", albumHandler.GetByID)
			albumsGroup.GET("/:itemID/tracks", coreHandler.GetTracksByAlbum)
			albumsGroup.GET("/genre/:genre", albumHandler.GetByGenre)
			albumsGroup.GET("/year/:year", albumHandler.GetByYear)
			albumsGroup.GET("/latest", albumHandler.GetLatestByAdded)
			albumsGroup.GET("/popular", albumHandler.GetPopular)
			// albumsGroup.GET("/external/:source/:externalID", albumHandler.GetByExternalSourceID)
		}
		tracksGroup := musicGroup.Group("/tracks")
		{
			tracksGroup.GET("/:itemID", trackHandler.GetByID)
			tracksGroup.GET("/most-played", trackHandler.GetMostPlayed)
			tracksGroup.GET("/genre/:genre", trackHandler.GetByGenre)
			tracksGroup.GET("/latest", trackHandler.GetLatestByAdded)
			tracksGroup.GET("/external/:source/:externalID", trackHandler.GetByExternalID)
		}

		// General music routes
		musicGroup.GET("/recent", coreHandler.GetRecentlyAddedMusic)
		musicGroup.GET("/recommendations/genre", coreHandler.GetGenreRecommendations)
	}
}
