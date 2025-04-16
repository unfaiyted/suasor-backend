package router

import (
	"fmt"
	"suasor/app/container"
	"suasor/app/handlers"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
)

// movieHandlerInterface defines common operations for all media handlers
type movieHandlerInterface interface {
	GetMovieByID(c *gin.Context)
	GetMoviesByGenre(c *gin.Context)
	GetMoviesByYear(c *gin.Context)
	GetMoviesByActor(c *gin.Context)
	GetMoviesByDirector(c *gin.Context)
	GetMoviesByRating(c *gin.Context)
	GetLatestMoviesByAdded(c *gin.Context)
	GetPopularMovies(c *gin.Context)
	GetTopRatedMovies(c *gin.Context)
	SearchMovies(c *gin.Context)
}

func RegisterClientMediaRoutes(rg *gin.RouterGroup, c *container.Container) {

	// Initialize handlers
	mediaHandler := container.MustGet[handlers.ClientMediaHandlers](c)

	// Create a map of movie types to handlers
	movieHandlerMap := map[string]movieHandlerInterface{
		"jellyfin": mediaHandler.JellyfinMovieHandler(),
		"emby":     mediaHandler.EmbyMovieHandler(),
		"plex":     mediaHandler.PlexMovieHandler(),
	}

	// Define series handler interface
	type seriesHandlerInterface interface {
		GetSeriesByID(c *gin.Context)
		GetSeriesByGenre(c *gin.Context)
		GetSeriesByYear(c *gin.Context)
		GetSeriesByActor(c *gin.Context)
		GetSeriesByCreator(c *gin.Context)
		GetSeriesByRating(c *gin.Context)
		GetLatestSeriesByAdded(c *gin.Context)
		GetPopularSeries(c *gin.Context)
		GetTopRatedSeries(c *gin.Context)
		SearchSeries(c *gin.Context)
		GetSeasonsBySeriesID(c *gin.Context)
	}

	seriesHandlerMap := map[string]seriesHandlerInterface{
		"jellyfin": mediaHandler.JellyfinSeriesHandler(),
		"emby":     mediaHandler.EmbySeriesHandler(),
		"plex":     mediaHandler.PlexSeriesHandler(),
	}

	// Define music handler interface
	type musicHandlerInterface interface {
		GetTrackByID(c *gin.Context)
		GetAlbumByID(c *gin.Context)
		GetArtistByID(c *gin.Context)
		GetTracksByAlbum(c *gin.Context)
		GetAlbumsByArtist(c *gin.Context)
		GetArtistsByGenre(c *gin.Context)
		GetAlbumsByGenre(c *gin.Context)
		GetTracksByGenre(c *gin.Context)
		GetAlbumsByYear(c *gin.Context)
		GetLatestAlbumsByAdded(c *gin.Context)
		GetPopularAlbums(c *gin.Context)
		GetPopularArtists(c *gin.Context)
		SearchMusic(c *gin.Context)
	}

	musicHandlerMap := map[string]musicHandlerInterface{
		"jellyfin": mediaHandler.JellyfinMusicHandler(),
		"emby":     mediaHandler.EmbyMusicHandler(),
		"plex":     mediaHandler.PlexMusicHandler(),
		"subsonic": mediaHandler.SubsonicMusicHandler(),
	}

	// Define playlist handler interface
	type playlistHandlerInterface interface {
		GetPlaylistByID(c *gin.Context)
		GetPlaylists(c *gin.Context)
		CreatePlaylist(c *gin.Context)
		UpdatePlaylist(c *gin.Context)
		DeletePlaylist(c *gin.Context)
		AddItemToPlaylist(c *gin.Context)
		RemoveItemFromPlaylist(c *gin.Context)
		SearchPlaylists(c *gin.Context)
	}

	// For now these will be placeholders until we implement the interface methods
	playlistHandlerMap := map[string]playlistHandlerInterface{
		"jellyfin": mediaHandler.JellyfinPlaylistHandler(),
		"emby":     mediaHandler.EmbyPlaylistHandler(),
		"plex":     mediaHandler.PlexPlaylistHandler(),
		"subsonic": mediaHandler.SubsonicPlaylistHandler(),
	}

	// Helper function to get the appropriate movie handler
	getMovieHandler := func(c *gin.Context) movieHandlerInterface {
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
	getSeriesHandler := func(c *gin.Context) seriesHandlerInterface {
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
	getMusicHandler := func(c *gin.Context) musicHandlerInterface {
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
	getPlaylistHandler := func(c *gin.Context) playlistHandlerInterface {
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

	// client/emby/1/movie/10
	clientGroup := rg.Group("/client/:clientType")
	client := clientGroup.Group("/:clientID")

	// Add movieType parameter to enable handler selection
	movie := client.Group("/movie")
	{
		movie.GET("/:id", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetMovieByID(c)
			}
		})

		movie.GET("/genre/:genre", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetMoviesByGenre(c)
			}
		})

		movie.GET("/year/:year", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetMoviesByYear(c)
			}
		})

		movie.GET("/actor/:actor", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetMoviesByActor(c)
			}
		})

		movie.GET("/director/:director", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetMoviesByDirector(c)
			}
		})

		movie.GET("/rating", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetMoviesByRating(c)
			}
		})

		movie.GET("/latest/:count", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetLatestMoviesByAdded(c)
			}
		})

		movie.GET("/popular/:count", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetPopularMovies(c)
			}
		})

		movie.GET("/top-rated/:count", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.GetTopRatedMovies(c)
			}
		})

		movie.GET("/search", func(c *gin.Context) {
			if handler := getMovieHandler(c); handler != nil {
				handler.SearchMovies(c)
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
				handler.GetSeriesByCreator(c)
			}
		})

		series.GET("/rating", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetSeriesByRating(c)
			}
		})

		series.GET("/latest/:count", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetLatestSeriesByAdded(c)
			}
		})

		series.GET("/popular/:count", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetPopularSeries(c)
			}
		})

		series.GET("/top-rated/:count", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.GetTopRatedSeries(c)
			}
		})

		series.GET("/search", func(c *gin.Context) {
			if handler := getSeriesHandler(c); handler != nil {
				handler.SearchSeries(c)
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
				handler.GetPlaylistByID(c)
			}
		})

		playlists.POST("", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.CreatePlaylist(c)
			}
		})

		playlists.PUT("/:id", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.UpdatePlaylist(c)
			}
		})

		playlists.DELETE("/:id", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.DeletePlaylist(c)
			}
		})

		playlists.POST("/:id/items", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.AddItemToPlaylist(c)
			}
		})

		playlists.DELETE("/:id/items/:itemID", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.RemoveItemFromPlaylist(c)
			}
		})

		// Search within playlists
		playlists.GET("/search", func(c *gin.Context) {
			if handler := getPlaylistHandler(c); handler != nil {
				handler.SearchPlaylists(c)
			}
		})
	}
}

