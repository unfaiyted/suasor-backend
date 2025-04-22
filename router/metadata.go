package router

import (
	"github.com/gin-gonic/gin"
	"suasor/di/container"
)

// RegisterMetadataRoutes registers all metadata client routes
func RegisterMetadataRoutes(api *gin.RouterGroup, c *container.Container) {
	metadata := api.Group("/metadata")

	// TMDB Routes
	tmdb := metadata.Group("/tmdb")

	// Note: We'll add actual route handlers when we have the TMDBHandler implemented
	// For now, just define the route structure

	// Movie routes
	tmdbMovies := tmdb.Group("/movies")
	tmdbMovies.GET("", func(c *gin.Context) { c.JSON(200, gin.H{"message": "TMDB movie endpoint - Not implemented yet"}) })
	tmdbMovies.GET("/search", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB movie search endpoint - Not implemented yet"})
	})
	tmdbMovies.GET("/recommendations", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB movie recommendations endpoint - Not implemented yet"})
	})
	tmdbMovies.GET("/popular", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB popular movies endpoint - Not implemented yet"})
	})
	tmdbMovies.GET("/trending", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB trending movies endpoint - Not implemented yet"})
	})

	// TV Show routes
	tmdbTV := tmdb.Group("/tv")
	tmdbTV.GET("", func(c *gin.Context) { c.JSON(200, gin.H{"message": "TMDB TV show endpoint - Not implemented yet"}) })
	tmdbTV.GET("/search", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB TV show search endpoint - Not implemented yet"})
	})
	tmdbTV.GET("/recommendations", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB TV show recommendations endpoint - Not implemented yet"})
	})
	tmdbTV.GET("/popular", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB popular TV shows endpoint - Not implemented yet"})
	})
	tmdbTV.GET("/trending", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB trending TV shows endpoint - Not implemented yet"})
	})
	tmdbTV.GET("/season", func(c *gin.Context) { c.JSON(200, gin.H{"message": "TMDB TV season endpoint - Not implemented yet"}) })
	tmdbTV.GET("/episode", func(c *gin.Context) { c.JSON(200, gin.H{"message": "TMDB TV episode endpoint - Not implemented yet"}) })

	// Person routes
	tmdbPeople := tmdb.Group("/people")
	tmdbPeople.GET("", func(c *gin.Context) { c.JSON(200, gin.H{"message": "TMDB person endpoint - Not implemented yet"}) })
	tmdbPeople.GET("/search", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB person search endpoint - Not implemented yet"})
	})

	// Collection routes
	tmdbCollections := tmdb.Group("/collections")
	tmdbCollections.GET("", func(c *gin.Context) { c.JSON(200, gin.H{"message": "TMDB collection endpoint - Not implemented yet"}) })
	tmdbCollections.GET("/search", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "TMDB collection search endpoint - Not implemented yet"})
	})
}

// The following handlers will be implemented later as part of app dependencies:
// TMDBHandler will return the TMDB metadata client handler
// func TMDBHandler(deps *app.AppDependencies) *handlers.MetadataClientHandler[*types.TMDBConfig] {
//    return deps.TMDBHandler()
// }
