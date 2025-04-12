package router

import (
	"suasor/app"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
)

// RegisterDirectMediaItemRoutes configures routes for direct media item access
// These routes access the local database media items rather than client-specific items
func RegisterDirectMediaItemRoutes(rg *gin.RouterGroup, deps *app.AppDependencies) {
	// Get handlers
	mediaHandlers := deps.MediaItemHandlers

	// Movie routes
	movies := rg.Group("/movies")
	{
		// Get movie by ID
		movies.GET("/:id", func(c *gin.Context) {
			mediaHandlers.MovieHandler().GetMediaItem(c)
		})

		// Get all movies with optional filtering
		movies.GET("", func(c *gin.Context) {
			// Query params handled in handler
			if handler := mediaHandlers.MovieHandler(); handler != nil {
				// Check for search query
				if q := c.Query("q"); q != "" {
					handler.SearchMediaItems(c)
				} else {
					handler.GetRecentMediaItems(c)
				}
			}
		})

		// Get movies by genre
		movies.GET("/genre/:genre", func(c *gin.Context) {
			// This is a placeholder - we'll need to implement this method
			responses.RespondNotImplemented(c, nil, "Getting movies by genre not implemented yet")
		})

		// Get movies by year
		movies.GET("/year/:year", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting movies by year not implemented yet")
		})

		// Get movies by actor
		movies.GET("/actor/:actor", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting movies by actor not implemented yet")
		})

		// Get movies by director
		movies.GET("/director/:director", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting movies by director not implemented yet")
		})

		// Get popular movies
		movies.GET("/popular", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting popular movies not implemented yet")
		})

		// Get latest movies
		movies.GET("/latest", func(c *gin.Context) {
			if handler := mediaHandlers.MovieHandler(); handler != nil {
				handler.GetRecentMediaItems(c)
			}
		})

		// Get top rated movies
		movies.GET("/top-rated", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting top rated movies not implemented yet")
		})
	}

	// Series routes
	series := rg.Group("/series")
	{
		// Get series by ID
		series.GET("/:id", func(c *gin.Context) {
			mediaHandlers.SeriesHandler().GetMediaItem(c)
		})

		// Get all series with optional filtering
		series.GET("", func(c *gin.Context) {
			if handler := mediaHandlers.SeriesHandler(); handler != nil {
				// Check for search query
				if q := c.Query("q"); q != "" {
					handler.SearchMediaItems(c)
				} else {
					handler.GetRecentMediaItems(c)
				}
			}
		})

		// Get series by genre
		series.GET("/genre/:genre", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting series by genre not implemented yet")
		})

		// Get series by year
		series.GET("/year/:year", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting series by year not implemented yet")
		})

		// Get series by actor
		series.GET("/actor/:actor", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting series by actor not implemented yet")
		})

		// Get series by creator
		series.GET("/creator/:creator", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting series by creator not implemented yet")
		})

		// Get popular series
		series.GET("/popular", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting popular series not implemented yet")
		})

		// Get latest series
		series.GET("/latest", func(c *gin.Context) {
			if handler := mediaHandlers.SeriesHandler(); handler != nil {
				handler.GetRecentMediaItems(c)
			}
		})

		// Get top rated series
		series.GET("/top-rated", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting top rated series not implemented yet")
		})

		// Get seasons for a series (this will need special implementation)
		series.GET("/:id/seasons", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting seasons for a series not implemented yet")
		})
	}

	// Music routes
	music := rg.Group("/music")
	{
		// Tracks routes
		tracks := music.Group("/tracks")
		{
			// Get track by ID
			tracks.GET("/:id", func(c *gin.Context) {
				mediaHandlers.TrackHandler().GetMediaItem(c)
			})

			// Get all tracks with optional filtering
			tracks.GET("", func(c *gin.Context) {
				if handler := mediaHandlers.TrackHandler(); handler != nil {
					// Check for search query
					if q := c.Query("q"); q != "" {
						handler.SearchMediaItems(c)
					} else {
						handler.GetRecentMediaItems(c)
					}
				}
			})

			// Get tracks by genre
			tracks.GET("/genre/:genre", func(c *gin.Context) {
				responses.RespondNotImplemented(c, nil, "Getting tracks by genre not implemented yet")
			})
		}

		// Albums routes
		albums := music.Group("/albums")
		{
			// Get album by ID
			albums.GET("/:id", func(c *gin.Context) {
				mediaHandlers.AlbumHandler().GetMediaItem(c)
			})

			// Get all albums with optional filtering
			albums.GET("", func(c *gin.Context) {
				if handler := mediaHandlers.AlbumHandler(); handler != nil {
					// Check for search query
					if q := c.Query("q"); q != "" {
						handler.SearchMediaItems(c)
					} else {
						handler.GetRecentMediaItems(c)
					}
				}
			})

			// Get album tracks
			albums.GET("/:id/tracks", func(c *gin.Context) {
				responses.RespondNotImplemented(c, nil, "Getting tracks for an album not implemented yet")
			})

			// Get albums by genre
			albums.GET("/genre/:genre", func(c *gin.Context) {
				responses.RespondNotImplemented(c, nil, "Getting albums by genre not implemented yet")
			})

			// Get albums by year
			albums.GET("/year/:year", func(c *gin.Context) {
				responses.RespondNotImplemented(c, nil, "Getting albums by year not implemented yet")
			})

			// Get latest albums
			albums.GET("/latest", func(c *gin.Context) {
				if handler := mediaHandlers.AlbumHandler(); handler != nil {
					handler.GetRecentMediaItems(c)
				}
			})

			// Get popular albums
			albums.GET("/popular", func(c *gin.Context) {
				responses.RespondNotImplemented(c, nil, "Getting popular albums not implemented yet")
			})
		}

		// Artists routes
		artists := music.Group("/artists")
		{
			// Get artist by ID
			artists.GET("/:id", func(c *gin.Context) {
				mediaHandlers.ArtistHandler().GetMediaItem(c)
			})

			// Get all artists with optional filtering
			artists.GET("", func(c *gin.Context) {
				if handler := mediaHandlers.ArtistHandler(); handler != nil {
					// Check for search query
					if q := c.Query("q"); q != "" {
						handler.SearchMediaItems(c)
					} else {
						handler.GetRecentMediaItems(c)
					}
				}
			})

			// Get artist albums
			artists.GET("/:id/albums", func(c *gin.Context) {
				responses.RespondNotImplemented(c, nil, "Getting albums for an artist not implemented yet")
			})

			// Get artists by genre
			artists.GET("/genre/:genre", func(c *gin.Context) {
				responses.RespondNotImplemented(c, nil, "Getting artists by genre not implemented yet")
			})

			// Get popular artists
			artists.GET("/popular", func(c *gin.Context) {
				responses.RespondNotImplemented(c, nil, "Getting popular artists not implemented yet")
			})
		}

		// General music search
		music.GET("/search", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Searching across all music types not implemented yet")
		})
	}

	// Playlist routes
	playlists := rg.Group("/playlists")
	{
		// Get playlist by ID
		playlists.GET("/:id", func(c *gin.Context) {
			mediaHandlers.PlaylistHandler().GetMediaItem(c)
		})

		// Get all playlists
		playlists.GET("", func(c *gin.Context) {
			if handler := mediaHandlers.PlaylistHandler(); handler != nil {
				// Check for search query
				if q := c.Query("q"); q != "" {
					handler.SearchMediaItems(c)
				} else {
					handler.GetRecentMediaItems(c)
				}
			}
		})

		// Create a new playlist
		playlists.POST("", func(c *gin.Context) {
			mediaHandlers.PlaylistHandler().CreateMediaItem(c)
		})

		// Update a playlist
		playlists.PUT("/:id", func(c *gin.Context) {
			mediaHandlers.PlaylistHandler().UpdateMediaItem(c)
		})

		// Delete a playlist
		playlists.DELETE("/:id", func(c *gin.Context) {
			mediaHandlers.PlaylistHandler().DeleteMediaItem(c)
		})

		// Add item to playlist - needs special implementation
		playlists.POST("/:id/items", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Adding items to a playlist not implemented yet")
		})

		// Remove item from playlist - needs special implementation
		playlists.DELETE("/:id/items/:itemID", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Removing items from a playlist not implemented yet")
		})

		// Search playlists
		playlists.GET("/search", func(c *gin.Context) {
			if handler := mediaHandlers.PlaylistHandler(); handler != nil {
				handler.SearchMediaItems(c)
			}
		})
	}

	// Collections routes
	collections := rg.Group("/collections")
	{
		// Get collection by ID
		collections.GET("/:id", func(c *gin.Context) {
			mediaHandlers.CollectionHandler().GetMediaItem(c)
		})

		// Get all collections
		collections.GET("", func(c *gin.Context) {
			if handler := mediaHandlers.CollectionHandler(); handler != nil {
				// Check for search query
				if q := c.Query("q"); q != "" {
					handler.SearchMediaItems(c)
				} else {
					handler.GetRecentMediaItems(c)
				}
			}
		})

		// Create a new collection
		collections.POST("", func(c *gin.Context) {
			mediaHandlers.CollectionHandler().CreateMediaItem(c)
		})

		// Update a collection
		collections.PUT("/:id", func(c *gin.Context) {
			mediaHandlers.CollectionHandler().UpdateMediaItem(c)
		})

		// Delete a collection
		collections.DELETE("/:id", func(c *gin.Context) {
			mediaHandlers.CollectionHandler().DeleteMediaItem(c)
		})

		// Get items in a collection - needs special implementation
		collections.GET("/:id/items", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting items in a collection not implemented yet")
		})

		// Add item to collection - needs special implementation
		collections.POST("/:id/items", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Adding items to a collection not implemented yet")
		})

		// Remove item from collection - needs special implementation
		collections.DELETE("/:id/items/:itemID", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Removing items from a collection not implemented yet")
		})
	}

	// Generic media items endpoint
	items := rg.Group("/media-items")
	{
		// Cross-media type search
		items.GET("/search", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Searching across all media types not implemented yet")
		})

		// Cross-media latest items
		items.GET("/latest", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting latest items across media types not implemented yet")
		})

		// Cross-media popular items
		items.GET("/popular", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting popular items across media types not implemented yet")
		})

		// Cross-media recommended items
		items.GET("/recommended", func(c *gin.Context) {
			responses.RespondNotImplemented(c, nil, "Getting recommended items across media types not implemented yet")
		})
	}

}

