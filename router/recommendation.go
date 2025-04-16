// router/recommendation.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app/container"
	"suasor/handlers"
)

// RegisterRecommendationRoutes registers all recommendation routes
func RegisterRecommendationRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Get the recommendation handler from dependencies
	recommendationHandler := container.MustGet[handlers.RecommendationHandler](c)

	recommendations := rg.Group("/recommendations")
	{
		// Base routes
		recommendations.GET("", recommendationHandler.GetRecommendations)
		recommendations.GET("/:id", recommendationHandler.GetRecommendationByID)

		// Special routes
		recommendations.GET("/recent", recommendationHandler.GetRecentRecommendations)
		recommendations.GET("/top", recommendationHandler.GetTopRecommendations)

		// Action routes
		recommendations.POST("/view", recommendationHandler.MarkRecommendationAsViewed)
		recommendations.POST("/rate", recommendationHandler.RateRecommendation)
	}
}
