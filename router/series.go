// router/series.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
	"suasor/handlers"
)

// SetupSeriesRoutes sets up the routes for series-related operations
func SetupSeriesRoutes(rg *gin.RouterGroup, app *app.App) {
	// Initialize handlers
	coreSeriesHandler := handlers.NewCoreSeriesHandler(
		app.Services().MediaItemServices().CoreSeriesService(),
		app.Services().MediaItemServices().CoreSeasonService(),
		app.Services().MediaItemServices().CoreEpisodeService(),
	)

	userSeriesHandler := handlers.NewUserSeriesHandler(
		app.Services().MediaItemServices().UserSeriesService(),
	)

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
		series.GET("/genre/:genre", coreSeriesHandler.GetSeriesByGenre)
		series.GET("/year/:year", coreSeriesHandler.GetSeriesByYear)
		
		// Discovery endpoints
		series.GET("/recently-aired", coreSeriesHandler.GetRecentlyAiredEpisodes)
		series.GET("/popular", coreSeriesHandler.GetPopularSeries)
		series.GET("/top-rated", coreSeriesHandler.GetTopRatedSeries)
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
		userSeries.GET("/continue-watching", userSeriesHandler.GetContinueWatchingSeries)
		userSeries.GET("/next-up", userSeriesHandler.GetNextUpEpisodes)
	}

	// Client-specific routes are handled in router/media.go
	// These routes follow the pattern /clients/media/{clientID}/series/...
}

// Series-specific helper functions
func getSeriesHandler[T interface{}](clientMedia *handlers.ClientMediaHandler[T]) *handlers.ClientMediaSeriesHandler[T] {
	return handlers.NewClientMediaSeriesHandler[T](clientMedia.SeriesService())
}