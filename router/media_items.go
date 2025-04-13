package router

import (
	"fmt"
	"suasor/app"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
)

type MediaItemHandlerInterface interface {
	GetMediaItem(c *gin.Context)
	GetMediaItemsByPerson(c *gin.Context)
	GetMediaItemsByYear(c *gin.Context)
	GetLatestMediaItemsByAdded(c *gin.Context)
	GetAllMediaItems(c *gin.Context)
	GetMediaItemsByClient(c *gin.Context)
	GetMediaItemsByGenre(c *gin.Context)
	GetMediaItemByExternalSourceID(c *gin.Context)
	GetPopularMediaItems(c *gin.Context)
	GetTopRatedMediaItems(c *gin.Context)
	SearchMediaItems(c *gin.Context)
}

// RegisterDirectMediaItemRoutes configures routes for direct media item access
// These routes access the local database media items rather than client-specific items
func RegisterDirectMediaItemRoutes(rg *gin.RouterGroup, deps *app.AppDependencies) {
	// Get handlers
	mediaHandlers := deps.MediaItemHandlers

	handlerMap := map[string]MediaItemHandlerInterface{
		"movies": mediaHandlers.MovieHandler(),
		"series": mediaHandlers.SeriesHandler(),

		"tracks":  mediaHandlers.TrackHandler(),
		"albums":  mediaHandlers.AlbumHandler(),
		"artists": mediaHandlers.ArtistHandler(),

		"collections": mediaHandlers.CollectionHandler(),
		"playlists":   mediaHandlers.PlaylistHandler(),
	}

	getHandler := func(c *gin.Context) MediaItemHandlerInterface {
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
				handler.GetMediaItem(c)
			}
		})

		// Get all media with optional filtering
		media.GET("", func(c *gin.Context) {
			handler := getHandler(c)
			// Check for search query
			if q := c.Query("q"); q != "" {
				handler.SearchMediaItems(c)
			} else {
				handler.GetAllMediaItems(c)
			}
		})

		// Get media by genre
		media.GET("/genre/:genre", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMediaItemsByGenre(c)
			}
		})

		// Get media by year
		media.GET("/year/:year", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMediaItemsByYear(c)
			}
		})

		// Get popular media
		media.GET("/popular", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetPopularMediaItems(c)
			}
		})

		// Get latest media
		media.GET("/latest", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetLatestMediaItemsByAdded(c)
			}
		})

		// Get top rated media
		media.GET("/top-rated", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetTopRatedMediaItems(c)
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
				handler.SearchMediaItems(c)
			}
		})

		// Get by external ID
		media.GET("/external/:source/:externalId", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMediaItemByExternalSourceID(c)
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
				handler.GetMediaItemsByPerson(c)
			}
		})

	}

	// Series routes
	series := rg.Group("/series")
	{
		// Get specialized series handler
		seriesHandler := mediaHandlers.SeriesHandler()
		seriesSpecificHandler := mediaHandlers.SeriesSpecificHandler()

		// Get series by ID - use base handler
		series.GET("/:id", seriesHandler.GetMediaItem)

		// Get all series - use base handler
		series.GET("", func(c *gin.Context) {
			if q := c.Query("q"); q != "" {
				seriesHandler.SearchMediaItems(c)
			} else {
				seriesHandler.GetAllMediaItems(c)
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
		series.GET("/genre/:genre", seriesHandler.GetMediaItemsByGenre)
		series.GET("/year/:year", seriesHandler.GetMediaItemsByYear)
		series.GET("/popular", seriesHandler.GetPopularMediaItems)
		series.GET("/latest", seriesHandler.GetRecentMediaItems)
		series.GET("/top-rated", seriesHandler.GetTopRatedMediaItems)
		series.GET("/external/:source/:externalId", seriesHandler.GetMediaItemByExternalSourceID)
	}

	// Music routes
	music := rg.Group("/music")
	{
		// Get specialized music handler - this should always be available
		musicHandler := mediaHandlers.MusicHandler()

		// Get album and artist handlers
		albumHandler := mediaHandlers.AlbumHandler()
		artistHandler := mediaHandlers.ArtistHandler()
		trackHandler := mediaHandlers.TrackHandler()

		// Tracks routes
		tracks := music.Group("/tracks")
		{
			// Get track by ID
			tracks.GET("/:id", trackHandler.GetMediaItem)

			// Get all tracks with optional filtering
			tracks.GET("", func(c *gin.Context) {
				// Check for search query
				if q := c.Query("q"); q != "" {
					trackHandler.SearchMediaItems(c)
				} else {
					trackHandler.GetAllMediaItems(c)
				}
			})

			// Most played tracks
			tracks.GET("/most-played", musicHandler.GetMostPlayedTracks)

			// Get tracks by genre
			tracks.GET("/genre/:genre", trackHandler.GetMediaItemsByGenre)

			// Get latest tracks
			tracks.GET("/latest", trackHandler.GetRecentMediaItems)

			// Get by external ID
			tracks.GET("/external/:source/:externalId", trackHandler.GetMediaItemByExternalSourceID)
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
			albums.GET("/genre/:genre", albumHandler.GetMediaItemsByGenre)

			// Get albums by year
			albums.GET("/year/:year", albumHandler.GetMediaItemsByYear)

			// Get latest albums
			albums.GET("/latest", albumHandler.GetRecentMediaItems)

			// Get popular albums
			albums.GET("/popular", albumHandler.GetPopularMediaItems)

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
					artistHandler.SearchMediaItems(c)
				} else {
					artistHandler.GetAllMediaItems(c)
				}
			})

			// Get artist albums - use specialized handler
			artists.GET("/:id/albums", musicHandler.GetAlbumsByArtistID)

			// Get similar artists - use specialized handler
			artists.GET("/:id/similar", musicHandler.GetSimilarArtists)

			// Get artists by genre
			artists.GET("/genre/:genre", artistHandler.GetMediaItemsByGenre)

			// Get popular artists
			artists.GET("/popular", artistHandler.GetPopularMediaItems)

			// Get by external ID
			artists.GET("/external/:source/:externalId", artistHandler.GetMediaItemByExternalSourceID)
		}

		// General music routes
		music.GET("/recent", musicHandler.GetRecentlyAddedMusic)

		// Genre-based recommendations - use specialized handler
		music.GET("/recommendations/genre", musicHandler.GetGenreRecommendations)
	}

	// Playlist routes
	playlists := rg.Group("/playlists")
	{
		// Get specialized playlist handler
		playlistHandler := mediaHandlers.PlaylistSpecificHandler()
		
		// Basic CRUD operations
		playlists.GET("", playlistHandler.GetPlaylists)
		playlists.GET("/:id", playlistHandler.GetPlaylistByID)
		playlists.POST("", playlistHandler.CreatePlaylist)
		playlists.PUT("/:id", playlistHandler.UpdatePlaylist)
		playlists.DELETE("/:id", playlistHandler.DeletePlaylist)
		
		// Playlist items management
		playlists.GET("/:id/items", playlistHandler.GetPlaylistItems)
		playlists.POST("/:id/items", playlistHandler.AddItemToPlaylist)
		playlists.DELETE("/:id/items/:itemId", playlistHandler.RemoveItemFromPlaylist)
		
		// Playlist reordering
		playlists.POST("/:id/reorder", playlistHandler.ReorderPlaylistItems)
		
		// Search playlists
		playlists.GET("/search", playlistHandler.SearchPlaylists)
		
		// Sync playlist across clients
		playlists.POST("/:id/sync", playlistHandler.SyncPlaylist)
	}

	// Collections routes
	collections := rg.Group("/collections")
	{
		// Get specialized collection handler
		collectionHandler := mediaHandlers.CollectionSpecificHandler()
		
		// Basic CRUD operations
		collections.GET("", collectionHandler.GetCollections)
		collections.GET("/:id", collectionHandler.GetCollectionByID)
		
		// Collection items management
		collections.GET("/:id/items", collectionHandler.GetCollectionItems)
		
		// Special collection types
		collections.GET("/smart", collectionHandler.GetSmartCollections)
		collections.GET("/featured", collectionHandler.GetFeaturedCollections)
	}
}
