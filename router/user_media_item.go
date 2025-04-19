package router

import (
	"suasor/app/container"
	// "suasor/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserMediaItemRoutes(rg *gin.RouterGroup, c *container.Container) {
	// handlers := container.MustGet[handlers.UserMediaItemHandlers](c)
	// userMediaItems := rg.Group("/user/item")
	// {
	// 	// Get user's media items by status
	// 	userMediaItems.GET("/favorites", handlers.GetFavoriteMovies)
	// 	userMediaItems.GET("/watched", handlers.GetWatchedMovies)
	// 	userMediaItems.GET("/watchlist", handlers.GetWatchlistMovies)
	// 	userMediaItems.GET("/recommended", handlers.GetRecommendedMovies)
	//
	// 	// Update user data for a movie
	// 	userMediaItems.PATCH("/:id", handlers.UpdateMovieUserData)
	// }
	//
	// // User-specific routes for media items
	// userMediaItems := rg.Group("/user-media-item")
	// {
	// 	// Get user's media items by status
	// 	userMediaItems.GET("/favorites", handlers.GetFavoriteMovies)
	// 	userMediaItems.GET("/watched", handlers.GetWatchedMovies)
	// 	userMediaItems.GET("/watch", handlers.GetWatchlistMovies)
	// 	userMediaItems.GET("/recommended", handlers.GetRecommendedMovies)
	//
	// 	// Update user data for a movie
	// 	userMediaItems.PATCH("/:id", handlers.UpdateMovieUserData)
	// }
}
