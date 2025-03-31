package router

import (
	"fmt"
	"suasor/app"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
)

// MediaHandlerInterface defines common operations for all media handlers
type MediaHandlerInterface interface {
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
	movieHandler := deps.ClientMediaHandlers

	// Create a map of media types to handlers
	handlerMap := map[string]MediaHandlerInterface{
		"jellyfin": movieHandler.JellyfinMovieHandler(),
		"emby":     movieHandler.EmbyMovieHandler(),
		"plex":     movieHandler.PlexMovieHandler(),
		// Add more handlers for different media types if needed
		// "tv": tvShowHandler,
		// "music": musicHandler,
	}

	// Helper function to get the appropriate handler
	getHandler := func(c *gin.Context) MediaHandlerInterface {
		mediaType := c.Param("mediaType")
		handler, exists := handlerMap[mediaType]
		if !exists {
			err := fmt.Errorf("unsupported media type: %s", mediaType)
			responses.RespondBadRequest(c, err, "Unsupported media type")
			return nil
		}
		return handler
	}

	clientGroup := rg.Group("/client/:clientType")
	client := clientGroup.Group("/:clientID")

	// Add mediaType parameter to enable handler selection
	media := client.Group("/:mediaType")
	{
		media.GET("/:movieID", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMovieByID(c)
			}
		})

		media.GET("/genre/:genre", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMoviesByGenre(c)
			}
		})

		media.GET("/year/:year", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMoviesByYear(c)
			}
		})

		media.GET("/actor/:actor", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMoviesByActor(c)
			}
		})

		media.GET("/director/:director", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMoviesByDirector(c)
			}
		})

		media.GET("/rating", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetMoviesByRating(c)
			}
		})

		media.GET("/latest/:count", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetLatestMoviesByAdded(c)
			}
		})

		media.GET("/popular/:count", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetPopularMovies(c)
			}
		})

		media.GET("/top-rated/:count", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetTopRatedMovies(c)
			}
		})

		media.GET("/search", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.SearchMovies(c)
			}
		})
	}
}
