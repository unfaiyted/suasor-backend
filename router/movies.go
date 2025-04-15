// router/movies.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
	"suasor/handlers"
)

// SetupMovieRoutes sets up the routes for movie operations
func SetupMovieRoutes(rg *gin.RouterGroup, app *app.App) {
	// Initialize handlers
	coreMovieHandler := handlers.NewCoreMovieHandler(
		app.Services().MediaItemServices().CoreMovieService(),
	)

	userMovieHandler := handlers.NewUserMovieHandler(
		app.Services().MediaItemServices().UserMovieService(),
	)

	// Core movie routes (database-focused)
	movies := rg.Group("/movies")
	{
		// Basic operations
		movies.GET("", coreMovieHandler.GetAll)
		movies.GET("/:id", coreMovieHandler.GetByID)
		
		// Search and filtering
		movies.GET("/search", coreMovieHandler.Search)
		movies.GET("/genre/:genre", coreMovieHandler.GetByGenre)
		movies.GET("/year/:year", coreMovieHandler.GetByYear)
		movies.GET("/actor/:actor", coreMovieHandler.GetByActor)
		movies.GET("/director/:director", coreMovieHandler.GetByDirector)
		
		// Discover
		movies.GET("/top-rated", coreMovieHandler.GetTopRated)
		movies.GET("/recently-added", coreMovieHandler.GetRecentlyAdded)
	}

	// User-specific movie routes
	userMovies := rg.Group("/user/movies")
	{
		// Get user's movies by status
		userMovies.GET("/favorites", userMovieHandler.GetFavoriteMovies)
		userMovies.GET("/watched", userMovieHandler.GetWatchedMovies)
		userMovies.GET("/watchlist", userMovieHandler.GetWatchlistMovies)
		userMovies.GET("/recommended", userMovieHandler.GetRecommendedMovies)
		
		// Update user data for a movie
		userMovies.PATCH("/:id", userMovieHandler.UpdateMovieUserData)
	}

	// Client-specific routes are handled in router/media.go
	// These routes follow the pattern /clients/media/{clientID}/movies/...
}

// Helper functions for client movies
func getMovieHandler[T interface{}](clientMedia *handlers.ClientMediaHandler[T]) *handlers.ClientMediaMovieHandler[T] {
	return handlers.NewClientMediaMovieHandler[T](clientMedia.MovieService())
}