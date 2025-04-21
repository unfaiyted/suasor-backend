// router/movies.go
package router

//
// import (
// 	"fmt"
// 	"suasor/app/container"
// 	apphandlers "suasor/app/handlers"
// 	"suasor/handlers"
// 	"suasor/types/responses"
//
// 	"github.com/gin-gonic/gin"
// )
//
// type ClientMovieHandlerInterface interface {
// 	GetMovieByID(c *gin.Context)
// 	GetMoviesByGenre(c *gin.Context)
// 	GetMoviesByYear(c *gin.Context)
// 	GetMoviesByActor(c *gin.Context)
// 	GetMoviesByDirector(c *gin.Context)
// 	GetMoviesByRating(c *gin.Context)
// 	GetLatestMoviesByAdded(c *gin.Context)
// 	GetPopularMovies(c *gin.Context)
// }
//
// // SetupMovieRoutes sets up the routes for movie operations
// func RegisterMovieRoutes(rg *gin.RouterGroup, c *container.Container) {
// 	// Initialize handlers
// 	coreMovieHandler := container.MustGet[handlers.CoreMovieHandler](c)
// 	userMovieHandler := container.MustGet[handlers.UserMovieHandler](c)
// 	clientMovieHandlers := container.MustGet[apphandlers.ClientMovieHandlers](c)
//
// 	// Core movie routes (database-focused)
// 	movies := rg.Group("/movies")
// 	{
// 		// Basic operations
// 		movies.GET("", coreMovieHandler.GetAll)
// 		movies.GET("/:id", coreMovieHandler.GetByID)
//
// 		// Search and filtering
// 		movies.GET("/search", coreMovieHandler.Search)
// 		movies.GET("/genre/:genre", coreMovieHandler.GetByGenre)
// 		movies.GET("/year/:year", coreMovieHandler.GetByYear)
// 		movies.GET("/actor/:actor", coreMovieHandler.GetByActor)
// 		movies.GET("/director/:director", coreMovieHandler.GetByDirector)
//
// 		// Discover
// 		movies.GET("/top-rated", coreMovieHandler.GetTopRated)
// 		movies.GET("/recently-added", coreMovieHandler.GetRecentlyAdded)
// 	}
//
// 	// User-specific movie routes
// 	userMovies := rg.Group("/user/movies")
// 	{
// 		// Get user's movies by status
// 		userMovies.GET("/favorites", userMovieHandler.GetFavoriteMovies)
// 		userMovies.GET("/watched", userMovieHandler.GetWatchedMovies)
// 		userMovies.GET("/watchlist", userMovieHandler.GetWatchlistMovies)
// 		userMovies.GET("/recommended", userMovieHandler.GetRecommendedMovies)
//
// 		// Update user data for a movie
// 		userMovies.PATCH("/:id", userMovieHandler.UpdateMovieUserData)
// 	}
//
// 	clientHandlerMap := map[string]ClientMovieHandlerInterface{
// 		"emby":     clientMovieHandlers.EmbyMovieHandler(),
// 		"jellyfin": clientMovieHandlers.JellyfinMovieHandler(),
// 		"plex":     clientMovieHandlers.PlexMovieHandler(),
// 	}
//
// 	getClientHandler := func(c *gin.Context) ClientMovieHandlerInterface {
// 		clientType := c.Param("clientType")
// 		handler, exists := clientHandlerMap[clientType]
// 		if !exists {
// 			err := fmt.Errorf("unsupported client type: %s", clientType)
// 			responses.RespondBadRequest(c, err, "Unsupported client type")
// 			return nil
// 		}
// 		return handler
// 	}
//
// 	// These routes follow the pattern /clients/media/{clientID}/movies/...
// 	clientMovies := rg.Group("/clients/:clientType/:clientID/movies")
// 	{
// 		// Get movie by ID
// 		clientMovies.GET("/:id", func(c *gin.Context) {
// 			getClientHandler(c).GetMovieByID(c)
// 		})
//
// 		// Get movies by genre
// 		clientMovies.GET("/genre/:genre", func(c *gin.Context) {
// 			getClientHandler(c).GetMoviesByGenre(c)
// 		})
//
// 		// Get movies by year
// 		clientMovies.GET("/year/:year", func(c *gin.Context) {
// 			getClientHandler(c).GetMoviesByYear(c)
// 		})
//
// 		// Get movies by actor
// 		clientMovies.GET("/actor/:actor", func(c *gin.Context) {
// 			getClientHandler(c).GetMoviesByActor(c)
// 		})
//
// 		// Get movies by director
// 		clientMovies.GET("/director/:director", func(c *gin.Context) {
// 			getClientHandler(c).GetMoviesByDirector(c)
// 		})
//
// 		// Get movies by rating range
// 		clientMovies.GET("/rating", func(c *gin.Context) {
// 			getClientHandler(c).GetMoviesByRating(c)
// 		})
//
// 		// Get latest movies
// 		clientMovies.GET("/latest/:count", func(c *gin.Context) {
// 			getClientHandler(c).GetLatestMoviesByAdded(c)
// 		})
//
// 		// Get popular movies
// 		clientMovies.GET("/popular/:count", func(c *gin.Context) {
// 			getClientHandler(c).GetPopularMovies(c)
// 		})
//
// 	}
//
// }
