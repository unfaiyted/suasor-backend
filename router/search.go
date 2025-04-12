package router

import (
	"github.com/gin-gonic/gin"
	"suasor/handlers"
)

// RegisterSearchRoutes registers all search-related routes
func RegisterSearchRoutes(rg *gin.RouterGroup, handler *handlers.SearchHandler) {
	searchGroup := rg.Group("/search")
	{
		// Main search endpoint
		searchGroup.GET("", handler.Search)

		// Search history and suggestions
		searchGroup.GET("/recent", handler.GetRecentSearches)
		searchGroup.GET("/trending", handler.GetTrendingSearches)
		searchGroup.GET("/suggestions", handler.GetSearchSuggestions)
	}
}

