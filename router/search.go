package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app/container"
	"suasor/handlers"
)

// RegisterSearchRoutes registers all search-related routes
func RegisterSearchRoutes(rg *gin.RouterGroup, c *container.Container) {
	handler := container.MustGet[handlers.SearchHandler](c)
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
