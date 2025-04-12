// router/recommendation.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
)

// RegisterRecommendationRoutes registers all recommendation routes
func RegisterRecommendationRoutes(rg *gin.RouterGroup, deps *app.AppDependencies) {
	// Get the recommendation handler from dependencies
	recommendationHandler := deps.RecommendationHandler()
	
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
