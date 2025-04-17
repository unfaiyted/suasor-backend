package router

import (
	"fmt"
	"suasor/app/container"
	"suasor/app/handlers"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
)

// type MediaItemHandlerInterface interface {
// 	GetMediaItem(c *gin.Context)
// 	GetByPerson(c *gin.Context)
// 	GetByYear(c *gin.Context)
// 	GetLatestByAdded(c *gin.Context)
// 	GetAll(c *gin.Context)
// 	GetByClient(c *gin.Context)
// 	GetByGenre(c *gin.Context)
// 	GetMediaItemByExternalSourceID(c *gin.Context)
// 	GetPopular(c *gin.Context)
// 	GetTopRated(c *gin.Context)
// 	Search(c *gin.Context)
// }

type CoreMediaItemHandler interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	GetByClientItemID(c *gin.Context)
	GetByExternalID(c *gin.Context)
	Search(c *gin.Context)
	GetRecentlyAdded(c *gin.Context)
	GetByType(c *gin.Context)
	GetByPerson(c *gin.Context)
	GetByYear(c *gin.Context)
	GetLatestByAdded(c *gin.Context)
	GetByClient(c *gin.Context)
	GetByGenre(c *gin.Context)
	GetPopular(c *gin.Context)
	GetTopRated(c *gin.Context)
}

// RegisterLocalMediaItemRoutes configures routes for direct media item access
// These routes access the local database media items rather than client-specific items
func RegisterLocalMediaItemRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Get handlers
	mediaHandlers := container.MustGet[handlers.MediaItemHandlers](c)
	mediaTypeHandlers := container.MustGet[handlers.MediaTypeHandlers](c)

	handlerMap := map[string]CoreMediaItemHandler{
		"movies": mediaHandlers.MovieCoreHandler(),
		"series": mediaHandlers.SeriesCoreHandler(),

		"tracks":  mediaHandlers.TrackCoreHandler(),
		"albums":  mediaHandlers.AlbumCoreHandler(),
		"artists": mediaHandlers.ArtistCoreHandler(),
	}

	getHandler := func(c *gin.Context) CoreMediaItemHandler {
		mediaType := c.Param("mediaType")
		handler, exists := handlerMap[mediaType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", mediaType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}
		return handler
	}

	media := rg.Group("/:mediaType")
	{
		// Get media item by ID
		media.GET("/:id", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetByID(c)
			}
		})

		// Get all media with optional filtering
		media.GET("", func(c *gin.Context) {
			handler := getHandler(c)
			// Check for search query
			if q := c.Query("q"); q != "" {
				handler.Search(c)
			} else {
				handler.GetAll(c)
			}
		})

		// Get media by genre
		media.GET("/genre/:genre", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetByGenre(c)
			}
		})

		// Get media by year
		media.GET("/year/:year", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetByYear(c)
			}
		})

		// Get popular media
		media.GET("/popular", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetPopular(c)
			}
		})

		// Get latest media
		media.GET("/latest", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetLatestByAdded(c)
			}
		})

		// Get top rated media
		media.GET("/top-rated", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetTopRated(c)
			}
		})

		media.GET("/search", func(c *gin.Context) {
			query := c.Query("q")
			if query == "" {
				responses.RespondBadRequest(c, nil, "Search query is required")
				return
			}

			// For now, we just search tracks
			if handler := getHandler(c); handler != nil {
				handler.Search(c)
			}
		})

		// Get by external ID
		media.GET("/external/:source/:externalId", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetByExternalID(c)
			}
		})
		// Recommended items
		media.GET("/recommended", func(c *gin.Context) {
			// This would need a recommendation service implementation
			responses.RespondNotImplemented(c, nil, "Getting recommended items across media types not implemented yet")
		})

		// roles: actor, director, creator, producer, writer, composer, cinematographer, editor, presenter, host, guest
		media.GET("/role/:role/:personId", func(c *gin.Context) {
			// Use the person-based method with a role filter by role
			role := c.Param("role")
			c.Request.URL.Query().Set("role", role)
			handler := getHandler(c)
			if handler != nil {
				handler.GetByPerson(c)
			}
		})

	}

	// Series routes
	series := rg.Group("/series")
	{
		// Get specialized series handler
		seriesHandler := mediaHandlers.SeriesCoreHandler()
		seriesSpecificHandler := mediaTypeHandlers.SeriesCoreHandler()

		// Get series by ID - use base handler
		series.GET("/:id", seriesHandler.GetByID)

		// Get all series - use base handler
		series.GET("", func(c *gin.Context) {
			if q := c.Query("q"); q != "" {
				seriesHandler.Search(c)
			} else {
				seriesHandler.GetAll(c)
			}
		})

		// Get seasons for a series - use specialized handler
		series.GET("/:id/seasons", seriesSpecificHandler.GetSeasonsBySeriesID)

		// Get episodes for a specific season - use specialized handler
		series.GET("/:id/seasons/:seasonNumber/episodes", seriesSpecificHandler.GetEpisodesBySeriesIDAndSeasonNumber)

		// Get all episodes for a series - use specialized handler
		series.GET("/:id/episodes", seriesSpecificHandler.GetAllEpisodes)

		// Get continue watching series - use specialized handler
		series.GET("/continue-watching", seriesSpecificHandler.GetContinueWatchingSeries)

		// Get next up episodes - use specialized handler
		series.GET("/next-up", seriesSpecificHandler.GetNextUpEpisodes)

		// Get recently aired episodes - use specialized handler
		series.GET("/recently-aired", seriesSpecificHandler.GetRecentlyAiredEpisodes)

		// Get series by network - use specialized handler
		series.GET("/network/:network", seriesSpecificHandler.GetSeriesByNetwork)

		// Standard handlers from the base MediaItemHandler
		series.GET("/genre/:genre", seriesHandler.GetByGenre)
		series.GET("/year/:year", seriesHandler.GetByYear)
		series.GET("/popular", seriesHandler.GetPopular)
		series.GET("/latest", seriesHandler.GetLatestByAdded)
		series.GET("/top-rated", seriesHandler.GetTopRated)
		series.GET("/external/:source/:externalId", seriesHandler.GetByExternalID)
	}

	// Music routes
	music := rg.Group("/music")
	{
		// Get specialized music handler - this should always be available
		musicHandler := mediaTypeHandlers.MusicCoreHandler()

		// Get album and artist handlers
		albumHandler := mediaHandlers.AlbumCoreHandler()
		artistHandler := mediaHandlers.ArtistCoreHandler()
		trackHandler := mediaHandlers.TrackCoreHandler()

		// Tracks routes
		tracks := music.Group("/tracks")
		{
			// Get track by ID
			tracks.GET("/:id", trackHandler.GetByID)

			// Get all tracks with optional filtering
			tracks.GET("", func(c *gin.Context) {
				// Check for search query
				if q := c.Query("q"); q != "" {
					trackHandler.Search(c)
				} else {
					trackHandler.GetAll(c)
				}
			})

			// Most played tracks
			tracks.GET("/most-played", musicHandler.GetMostPlayed)

			// Get tracks by genre
			tracks.GET("/genre/:genre", trackHandler.GetByGenre)

			// Get latest tracks
			tracks.GET("/latest", trackHandler.GetLatestByAdded)

			// Get by external ID
			tracks.GET("/external/:source/:externalId", trackHandler.GetByExternalID)
		}

		// Albums routes
		albums := music.Group("/albums")
		{
			// Get album by ID
			albums.GET("/:id", albumHandler.GetMediaItem)

			// Get album tracks - use specialized handler
			albums.GET("/:id/tracks", musicHandler.GetTracksByAlbumID)

			// Get top rated albums - use specialized handler
			albums.GET("/top-rated", musicHandler.GetTopRatedAlbums)

			// Get albums by genre
			albums.GET("/genre/:genre", albumHandler.GetByGenre)

			// Get albums by year
			albums.GET("/year/:year", albumHandler.GetByYear)

			// Get latest albums
			albums.GET("/latest", albumHandler.GetRecent)

			// Get popular albums
			albums.GET("/popular", albumHandler.GetPopular)

			// Get by external ID
			albums.GET("/external/:source/:externalId", albumHandler.GetMediaItemByExternalSourceID)
		}

		// Artists routes
		artists := music.Group("/artists")
		{
			// Get artist by ID
			artists.GET("/:id", artistHandler.GetMediaItem)

			// Get all artists with optional filtering
			artists.GET("", func(c *gin.Context) {
				// Check for search query
				if q := c.Query("q"); q != "" {
					artistHandler.Search(c)
				} else {
					artistHandler.GetAll(c)
				}
			})

			// Get artist albums - use specialized handler
			artists.GET("/:id/albums", musicHandler.GetAlbumsByArtistID)

			// Get similar artists - use specialized handler
			artists.GET("/:id/similar", musicHandler.GetSimilarArtists)

			// Get artists by genre
			artists.GET("/genre/:genre", artistHandler.GetByGenre)

			// Get popular artists
			artists.GET("/popular", artistHandler.GetPopular)

			// Get by external ID
			artists.GET("/external/:source/:externalId", artistHandler.GetMediaItemByExternalSourceID)
		}

		// General music routes
		music.GET("/recent", musicHandler.GetRecentlyAddedMusic)

		// Genre-based recommendations - use specialized handler
		music.GET("/recommendations/genre", musicHandler.GetGenreRecommendations)
	}

}
