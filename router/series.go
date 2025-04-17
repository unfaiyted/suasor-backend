// router/series.go
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	"suasor/handlers"
	"suasor/types/responses"
)

type ClientSeriesHandlerInterface interface {
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

// RegisterSeriesRoutes sets up the routes for series-related operations
func RegisterSeriesRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Initialize handlers
	coreSeriesHandler := container.MustGet[handlers.CoreSeriesHandler](c)
	userSeriesHandler := container.MustGet[handlers.UserSeriesHandler](c)
	clientSeriesHandler := container.MustGet[apphandlers.ClientSeriesHandlers](c)

	// Core series routes (database-focused)
	series := rg.Group("/series")
	{
		// Basic CRUD operations
		series.GET("/:id", coreSeriesHandler.GetByID)
		series.GET("", coreSeriesHandler.GetAll)

		// Core metadata operations
		series.GET("/:id/seasons", coreSeriesHandler.GetSeasonsBySeriesID)
		series.GET("/:id/seasons/:seasonNumber/episodes", coreSeriesHandler.GetEpisodesBySeriesIDAndSeasonNumber)
		series.GET("/:id/episodes", coreSeriesHandler.GetAllEpisodes)

		// Specialized metadata filters
		series.GET("/network/:network", coreSeriesHandler.GetSeriesByNetwork)
		series.GET("/genre/:genre", coreSeriesHandler.GetByGenre)
		series.GET("/year/:year", coreSeriesHandler.GetByYear)

		// Discovery endpoints
		series.GET("/recently-aired", coreSeriesHandler.GetRecentlyAiredEpisodes)
		series.GET("/popular", coreSeriesHandler.GetPopular)
		series.GET("/top-rated", coreSeriesHandler.GetTopRated)

		// Search
		series.GET("/search", func(c *gin.Context) {
			// This is a universal search route that might redirect to client-specific search
			// Later we might want to add a centralized search service that combines results
			clientSeriesHandler.PlexSeriesHandler().SearchSeries(c)
		})
	}

	// User-specific series routes
	userSeries := rg.Group("/user/series")
	{
		// User's series collections
		userSeries.GET("/favorites", userSeriesHandler.GetFavoriteSeries)
		userSeries.GET("/watched", userSeriesHandler.GetWatchedSeries)
		userSeries.GET("/watchlist", userSeriesHandler.GetWatchlistSeries)

		// User interactions with series
		userSeries.PATCH("/:id", userSeriesHandler.UpdateSeriesUserData)

		// Personalized recommendations
		// userSeries.GET("/continue-watching", userSeriesHandler.GetContinueWatching)
		// userSeries.GET("/next-up", userSeriesHandler.GetNextUpEpisodes)
	}

	// Set up client-specific handler map for dynamic routing
	clientHandlerMap := map[string]ClientSeriesHandlerInterface{
		"emby":     clientSeriesHandler.EmbySeriesHandler(),
		"jellyfin": clientSeriesHandler.JellyfinSeriesHandler(),
		"plex":     clientSeriesHandler.PlexSeriesHandler(),
	}

	getClientHandler := func(c *gin.Context) ClientSeriesHandlerInterface {
		clientType := c.Param("clientType")
		handler, exists := clientHandlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}
		return handler
	}

	// Client-specific series routes
	// These routes follow the pattern /clients/:clientType/:clientID/series/...
	clientSeries := rg.Group("/clients/:clientType/:clientID/series")
	{
		// Basic operations
		clientSeries.GET("/:seriesID", func(c *gin.Context) {
			getClientHandler(c).GetSeriesByID(c)
		})

		// Seasons and episodes
		clientSeries.GET("/:seriesID/seasons", func(c *gin.Context) {
			getClientHandler(c).GetSeasonsBySeriesID(c)
		})

		// Discovery endpoints
		clientSeries.GET("/popular/:count", func(c *gin.Context) {
			getClientHandler(c).GetPopularSeries(c)
		})

		clientSeries.GET("/top-rated/:count", func(c *gin.Context) {
			getClientHandler(c).GetTopRatedSeries(c)
		})

		clientSeries.GET("/latest/:count", func(c *gin.Context) {
			getClientHandler(c).GetLatestSeriesByAdded(c)
		})

		// Filters
		clientSeries.GET("/genre/:genre", func(c *gin.Context) {
			getClientHandler(c).GetSeriesByGenre(c)
		})

		clientSeries.GET("/year/:year", func(c *gin.Context) {
			getClientHandler(c).GetSeriesByYear(c)
		})

		clientSeries.GET("/actor/:actor", func(c *gin.Context) {
			getClientHandler(c).GetSeriesByActor(c)
		})

		clientSeries.GET("/creator/:creator", func(c *gin.Context) {
			getClientHandler(c).GetSeriesByCreator(c)
		})

		// Search endpoint
		clientSeries.GET("/search", func(c *gin.Context) {
			getClientHandler(c).SearchSeries(c)
		})

		// Rating-based filtering
		clientSeries.GET("/rating", func(c *gin.Context) {
			getClientHandler(c).GetSeriesByRating(c)
		})
	}
}
