package router

import (
	"suasor/app/container"
	mediatypes "suasor/client/media/types"
	"suasor/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterLocalMediaItemRoutes configures routes for direct media item access
func RegisterLocalMediaItemRoutes(rg *gin.RouterGroup, c *container.Container) {

	// Register generic media routes that work across all types
	registerGenericMediaRoutes[mediatypes.Movie](c, rg)
	registerGenericMediaRoutes[mediatypes.Series](c, rg)
	registerGenericMediaRoutes[mediatypes.Season](c, rg)
	registerGenericMediaRoutes[mediatypes.Episode](c, rg)
	registerGenericMediaRoutes[mediatypes.Track](c, rg)
	registerGenericMediaRoutes[mediatypes.Album](c, rg)
	registerGenericMediaRoutes[mediatypes.Artist](c, rg)

	// Register specialized routes for series
	registerSeriesRoutes(c, rg)

	// Register specialized routes for music (tracks, albums, artists)
	registerMusicRoutes(c, rg)
}

// registerGenericMediaRoutes sets up routes common to all media types
func registerGenericMediaRoutes[T mediatypes.MediaData](c *container.Container, rg *gin.RouterGroup) {

	handler := container.MustGet[handlers.CoreMediaItemHandler[T]](c)

	// Create media group i.e (/movies, /series, etc.)
	mediaGroup := rg.Group("/" + handler.GetType())

	// Register the routes
	mediaGroup.GET("", handler.GetAll)
	mediaGroup.GET("/:id", handler.GetByID)
	mediaGroup.GET("/genre/:genre", handler.GetByGenre)
	mediaGroup.GET("/year/:year", handler.GetByYear)
	mediaGroup.GET("/popular", handler.GetPopular)
	mediaGroup.GET("/latest", handler.GetLatestByAdded)
	mediaGroup.GET("/top-rated", handler.GetTopRated)
	mediaGroup.GET("/search", handler.Search)
	mediaGroup.GET("/external/:source/:externalId", handler.GetByExternalID)
	mediaGroup.GET("/role/:role/:personId", handler.GetByPerson)

}

// registerSeriesRoutes sets up routes specific to series
func registerSeriesRoutes(c *container.Container, rg *gin.RouterGroup) {
	coreHandler := container.MustGet[handlers.CoreSeriesHandler](c)
	userHandler := container.MustGet[handlers.UserSeriesHandler](c)

	// Series-specific routes
	seriesGroup := rg.Group("/series")

	seriesGroup.GET("/:id/seasons", coreHandler.GetSeasonsBySeriesID)
	seriesGroup.GET("/:id/seasons/:seasonNumber/episodes", coreHandler.GetEpisodesBySeriesIDAndSeasonNumber)
	seriesGroup.GET("/:id/episodes", coreHandler.GetAllEpisodes)
	seriesGroup.GET("/continue-watching", userHandler.GetNextUpEpisodes)
	seriesGroup.GET("/next-up", userHandler.GetNextUpEpisodes)
	seriesGroup.GET("/recently-aired", coreHandler.GetRecentlyAiredEpisodes)
	seriesGroup.GET("/network/:network", coreHandler.GetSeriesByNetwork)
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
			artistsGroup.GET("/:id", artistHandler.GetByID)
			artistsGroup.GET("/:id/albums", coreHandler.GetAlbumsByArtistID)
			artistsGroup.GET("/:id/similar", coreHandler.GetSimilarArtists)
			artistsGroup.GET("/genre/:genre", artistHandler.GetByGenre)
			artistsGroup.GET("/popular", artistHandler.GetPopular)
			// artistsGroup.GET("/external/:source/:externalId", artistHandler.GetMediaItemByExternalSourceID)
		}
		albumsGroup := musicGroup.Group("/albums")
		{
			albumsGroup.GET("/:id", albumHandler.GetByID)
			albumsGroup.GET("/:id/tracks", coreHandler.GetTracksByAlbum)
			albumsGroup.GET("/genre/:genre", albumHandler.GetByGenre)
			albumsGroup.GET("/year/:year", albumHandler.GetByYear)
			albumsGroup.GET("/latest", albumHandler.GetLatestByAdded)
			albumsGroup.GET("/popular", albumHandler.GetPopular)
			// albumsGroup.GET("/external/:source/:externalId", albumHandler.GetByExternalSourceID)
		}
		tracksGroup := musicGroup.Group("/tracks")
		{
			tracksGroup.GET("/:id", trackHandler.GetByID)
			tracksGroup.GET("/most-played", trackHandler.GetMostPlayed)
			tracksGroup.GET("/genre/:genre", trackHandler.GetByGenre)
			tracksGroup.GET("/latest", trackHandler.GetLatestByAdded)
			tracksGroup.GET("/external/:source/:externalId", trackHandler.GetByExternalID)
		}

		// General music routes
		musicGroup.GET("/recent", coreHandler.GetRecentlyAddedMusic)
		musicGroup.GET("/recommendations/genre", coreHandler.GetGenreRecommendations)
	}
}

