package router

import (
	"suasor/di/container"

	"github.com/gin-gonic/gin"
	mediatypes "suasor/clients/media/types"
	"suasor/handlers"
)

// RegisterMediaPlayHistoryRoutes configures routes for media play history
func RegisterMediaPlayHistoryRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Initialize the repository, service, and handler
	// We need to initialize these here since they're not fully integrated into the dependency system yet
	handler := container.MustGet[handlers.UserMediaItemDataHandler[*mediatypes.Movie]](c)

	// History routes
	history := rg.Group("/history")
	{
		// Get play history
		history.GET("", handler.GetMediaPlayHistory)

		// Record a play
		history.POST("", handler.RecordMediaPlay)

		// Get continue watching
		// history.GET("/continue-watching", handler.GetContinueWatching)

		// Get specific history entry
		// history.GET("/:id", handler.GetMediaPlayHistoryByID)

		// Delete history entry
		// history.DELETE("/:id", handler.DeleteHstory)

		// Clear user history
		history.DELETE("/clear", handler.ClearUserHistory)

		// Get play history for a specific media item
		// history.GET("/media/:mediaItemId", handler.GetMediaPlayHistoryByMediaItem)

		// Toggle favorite status for a media item
		history.PUT("/media/:itemId/favorite", handler.ToggleFavorite)

		// Update rating for a media item
		history.PUT("/media/:itemId/rating", handler.UpdateUserRating)
	}

	// Favorites routes
	favorites := rg.Group("/favorites")
	{
		// Get all favorites
		favorites.GET("", handler.GetFavorites)
	}
}
