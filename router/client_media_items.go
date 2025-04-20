package router

import (
	"fmt"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/handlers"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
)

func RegisterClientMediaItemRoutes(rg *gin.RouterGroup, c *container.Container) {

	// Initialize handlers
	mediaHandler := container.MustGet[apphandlers.ClientMediaHandlers](c)

	// Create a map of movie types to handlers
	movieHandlerMap := map[string]handlers.CoreMovieHandler{
		"jellyfin": mediaHandler.JellyfinMovieHandler(),
		"emby":     mediaHandler.EmbyMovieHandler(),
		"plex":     mediaHandler.PlexMovieHandler(),
	}

	seriesHandlerMap := map[string]handlers.ClientSeriesHandler[*types.EmbyConfig]{
		"jellyfin": mediaHandler.JellyfinSeriesHandler(),
		"emby":     mediaHandler.EmbySeriesHandler(),
		"plex":     mediaHandler.PlexSeriesHandler(),
	}

	// musicHandlerMap := map[string]handlers.ClientMusicHandler{
	// 	"jellyfin": mediaHandler.JellyfinMusicHandler(),
	// 	"emby":     mediaHandler.EmbyMusicHandler(),
	// 	"plex":     mediaHandler.PlexMusicHandler(),
	// 	"subsonic": mediaHandler.SubsonicMusicHandler(),
	// }
	//
	// // For now these will be placeholders until we implement the interface methods
	// playlistHandlerMap := map[string]handlers.ClientListHandler[*mediatypes.Playlist]{
	// 	"jellyfin": mediaHandler.JellyfinPlaylistHandler(),
	// 	"emby":     mediaHandler.EmbyPlaylistHandler(),
	// 	"plex":     mediaHandler.PlexPlaylistHandler(),
	// 	"subsonic": mediaHandler.SubsonicPlaylistHandler(),
	// }

	// Helper function to get the appropriate movie handler
	getMovieHandler := func(c *gin.Context) handlers.ClientMovieHandler {
		clientType := c.Param("clientType")
		handler, exists := movieHandlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}
		return handler
	}

	// Helper function to get the appropriate series handler
	getSeriesHandler := func(c *gin.Context) handlers.ClientSeriesHandler[*types.EmbyConfig] {
		clientType := c.Param("clientType")
		handler, exists := seriesHandlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}
		return handler
	}

	// Helper function to get the appropriate music handler
	getMusicHandler := func(c *gin.Context) handlers.ClientMusicHandler {
		clientType := c.Param("clientType")
		handler, exists := musicHandlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}
		return handler
	}

	// Helper function to get the appropriate playlist handler
	getPlaylistHandler := func(c *gin.Context) handlers.ClientListHandler[*types.EmbyConfig, *mediatypes.Playlist] {
		clientType := c.Param("clientType")
		handler, exists := playlistHandlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}

		// For now, since handlers are not implemented
		if handler == nil {
			err := fmt.Errorf("playlist support not implemented for client type: %s", clientType)
			responses.RespondInternalError(c, err, "Feature not implemented")
			return nil
		}

		return handler
	}

	// client/11/movie/10
	clientGroup := rg.Group("/client")
	client := clientGroup.Group("/:clientID")

	// Add movieType parameter to enable handler selection
	movie := client.Group("/movie")
	{
		movie.GET("/:id", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetByID(c)
			}
		})

		movie.GET("/genre/:genre", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetByGenre(c)
			}
		})

		movie.GET("/year/:year", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetByYear(c)
			}
		})

		movie.GET("/actor/:actor", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetByActor(c)
			}
		})

		movie.GET("/director/:director", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetByDirector(c)
			}
		})

		movie.GET("/rating", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetByRating(c)
			}
		})

		movie.GET("/latest/:count", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetLatestByAdded(c)
			}
		})

		movie.GET("/popular/:count", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetPopular(c)
			}
		})

		movie.GET("/top-rated/:count", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetTopRated(c)
			}
		})

		movie.GET("/search", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.Search(c)
			}
		})
	}

	series := client.Group("/series")
	{
		series.GET("/:seriesID", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetSeriesByID(c)
			}
		})

		series.GET("/:seriesID/seasons", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetSeasonsBySeriesID(c)
			}
		})

		series.GET("/genre/:genre", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetSeriesByGenre(c)
			}
		})

		series.GET("/year/:year", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetSeriesByYear(c)
			}
		})

		series.GET("/actor/:actor", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetSeriesByActor(c)
			}
		})

		series.GET("/creator/:creator", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetByCreator(c)
			}
		})

		series.GET("/rating", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetSeriesByRating(c)
			}
		})

		series.GET("/latest/:count", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetLatestByAdded(c)
			}
		})

		series.GET("/popular/:count", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetPopular(c)
			}
		})

		series.GET("/top-rated/:count", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetTopRated(c)
			}
		})

		series.GET("/search", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.Search(c)
			}
		})
	}

	// Music routes
	music := client.Group("/music")
	{
		// Track routes
		tracks := music.Group("/tracks")
		{
			tracks.GET("/:trackID", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetTrackByID(c)
				}
			})

			tracks.GET("/genre/:genre", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetTracksByGenre(c)
				}
			})
		}

		// Album routes
		albums := music.Group("/albums")
		{
			albums.GET("/:albumID", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetAlbumByID(c)
				}
			})

			albums.GET("/:albumID/tracks", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetTracksByAlbum(c)
				}
			})

			albums.GET("/genre/:genre", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetAlbumsByGenre(c)
				}
			})

			albums.GET("/year/:year", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetAlbumsByYear(c)
				}
			})

			albums.GET("/latest/:count", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetLatestAlbumsByAdded(c)
				}
			})

			albums.GET("/popular/:count", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetPopularAlbums(c)
				}
			})
		}

		// Artist routes
		artists := music.Group("/artists")
		{
			artists.GET("/:artistID", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetArtistByID(c)
				}
			})

			artists.GET("/:artistID/albums", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetAlbumsByArtist(c)
				}
			})

			artists.GET("/genre/:genre", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetArtistsByGenre(c)
				}
			})

			artists.GET("/popular/:count", func(c *gin.Context) {
				if handler := getMusicHandler(c); handler != nil {
					handler.GetPopularArtists(c)
				}
			})
		}

		// General music search
		music.GET("/search", func(c *gin.Context) {
			if handler := getMusicHandler(c); handler != nil {
				handler.SearchMusic(c)
			}
		})
	}

	// Playlist routes
	playlists := client.Group("/playlists")
	{
		playlists.GET("", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.GetPlaylists(c)
			}
		})

		playlists.GET("/:id", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.GetByID(c)
			}
		})

		playlists.POST("", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.Create(c)
			}
		})

		playlists.PUT("/:id", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.Update(c)
			}
		})

		playlists.DELETE("/:id", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.Delete(c)
			}
		})

		playlists.POST("/:id/items", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.AddItem(c)
			}
		})

		playlists.DELETE("/:id/items/:itemID", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.RemoveItem(c)
			}
		})

		// Search within playlists
		playlists.GET("/search", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.Search(c)
			}
		})
	}
}
