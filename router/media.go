package router

import (
	"fmt"
	"suasor/app"
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

func RegisterMediaClientRoutes(rg *gin.RouterGroup, deps *app.AppDependencies) {

	// Initialize handlers
	mediaHandler := deps.ClientMediaHandlers

	// Create a map of movie types to handlers
	movieHandlerMap := map[string]movieHandlerInterface{
		"jellyfin": mediaHandler.JellyfinMovieHandler(),
		"emby":     mediaHandler.EmbyMovieHandler(),
		"plex":     mediaHandler.PlexMovieHandler(),
	}

	// seriesHandlerMap := map[string]seriesHandlerInterface{
	// "jellyfin": mediaHandler.JellyfinSeriesHandler(),
	// "emby":     mediaHandler.EmbySeriesHandler(),
	// "plex":     mediaHandler.PlexSeriesHandler(),
	// }

	// Helper function to get the appropriate handler
	getMovieHandler := func(c *gin.Context) movieHandlerInterface {
		clientType := c.Param("clientType")
		handler, exists := movieHandlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported movie type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported movie type")
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

	// series := client.Group("/series")
	// {
	// series.GET("/:id", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetSeriesByID(c)
	// 	}
	// })
	//
	// series.GET("/genre/:genre", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetSeriesByGenre(c)
	// 	}
	// })
	//
	// series.GET("/year/:year", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetSeriesByYear(c)
	// 	}
	// })
	//
	// series.GET("/actor/:actor", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetSeriesByActor(c)
	// 	}
	// })
	//
	// series.GET("/director/:director", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetSeriesByDirector(c)
	// 	}
	// })
	//
	// series.GET("/rating", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetSeriesByRating(c)
	// 	}
	// })
	// series.GET("/latest/:count", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetLatestSeriesByAdded(c)
	// 	}
	// })
	//
	// series.GET("/popular/:count", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetPopularSeries(c)
	// 	}
	// })
	//
	// series.GET("/top-rated/:count", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.GetTopRatedSeries(c)
	// 	}
	// })
	//
	// series.GET("/search", func(c *gin.Context) {
	// 	if handler := getSeriesHandler(c); handler != nil {
	// 		handler.SearchSeries(c)
	// 	}
	// })
	// }
}
