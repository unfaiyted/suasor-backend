package router

import (
	"suasor/app"
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

// RegisterMediaPlayHistoryRoutes configures routes for media play history
func RegisterMediaPlayHistoryRoutes(rg *gin.RouterGroup, deps *app.AppDependencies) {
	// Initialize the repository, service, and handler
	// Note: Eventually this should be moved to the dependency injection system
	repo := repository.NewMediaPlayHistoryRepository(deps.GetDB())
	service := services.NewMediaPlayHistoryService(repo)
	handler := handlers.NewMediaPlayHistoryHandler(service)

	// History routes
	history := rg.Group("/history")
	{
		// Get play history
		history.GET("", handler.GetMediaPlayHistory)
		
		// Record a play
		history.POST("", handler.RecordMediaPlay)
		
		// Get continue watching
		history.GET("/continue-watching", handler.GetContinueWatching)
		
		// Get specific history entry
		history.GET("/:id", handler.GetMediaPlayHistoryByID)
		
		// Delete history entry
		history.DELETE("/:id", handler.DeleteHistory)
		
		// Clear user history
		history.DELETE("/clear", handler.ClearUserHistory)
		
		// Get play history for a specific media item
		history.GET("/media/:mediaItemId", handler.GetMediaPlayHistoryByMediaItem)
		
		// Toggle favorite status for a media item
		history.PUT("/media/:mediaItemId/favorite", handler.ToggleFavorite)
		
		// Update rating for a media item
		history.PUT("/media/:mediaItemId/rating", handler.UpdateUserRating)
	}
	
	// Favorites routes
	favorites := rg.Group("/favorites")
	{
		// Get all favorites
		favorites.GET("", handler.GetFavorites)
	}
}