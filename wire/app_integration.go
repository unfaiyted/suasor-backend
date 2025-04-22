package wire

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
)

// This is an example of how to integrate Wire with your application

// SetupRoutes demonstrates how to use all the Wire-generated handlers
func SetupRoutes(router *gin.Engine) {
	// Create context
	ctx := context.Background()

	// Initialize handlers
	handlers, err := InitializeAllHandlers(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	// Set up routes for all media types
	
	// Movies
	movieRoutes := router.Group("/movies")
	{
		movieRoutes.GET("/:id", handlers.Media.MovieHandler.GetByID)
		movieRoutes.POST("", handlers.Media.MovieHandler.Create)
		movieRoutes.PUT("/:id", handlers.Media.MovieHandler.Update)
		movieRoutes.DELETE("/:id", handlers.Media.MovieHandler.Delete)
	}
	
	// Series
	seriesRoutes := router.Group("/series")
	{
		seriesRoutes.GET("/:id", handlers.Media.SeriesHandler.GetByID)
		seriesRoutes.POST("", handlers.Media.SeriesHandler.Create)
		seriesRoutes.PUT("/:id", handlers.Media.SeriesHandler.Update)
		seriesRoutes.DELETE("/:id", handlers.Media.SeriesHandler.Delete)
	}
	
	// Episodes
	episodeRoutes := router.Group("/episodes")
	{
		episodeRoutes.GET("/:id", handlers.Media.EpisodeHandler.GetByID)
		episodeRoutes.POST("", handlers.Media.EpisodeHandler.Create)
		episodeRoutes.PUT("/:id", handlers.Media.EpisodeHandler.Update)
		episodeRoutes.DELETE("/:id", handlers.Media.EpisodeHandler.Delete)
	}
	
	// Music
	musicRoutes := router.Group("/music")
	{
		// Tracks
		musicRoutes.GET("/tracks/:id", handlers.Media.TrackHandler.GetByID)
		musicRoutes.POST("/tracks", handlers.Media.TrackHandler.Create)
		musicRoutes.PUT("/tracks/:id", handlers.Media.TrackHandler.Update)
		musicRoutes.DELETE("/tracks/:id", handlers.Media.TrackHandler.Delete)
		
		// Albums
		musicRoutes.GET("/albums/:id", handlers.Media.AlbumHandler.GetByID)
		musicRoutes.POST("/albums", handlers.Media.AlbumHandler.Create)
		musicRoutes.PUT("/albums/:id", handlers.Media.AlbumHandler.Update)
		musicRoutes.DELETE("/albums/:id", handlers.Media.AlbumHandler.Delete)
		
		// Artists
		musicRoutes.GET("/artists/:id", handlers.Media.ArtistHandler.GetByID)
		musicRoutes.POST("/artists", handlers.Media.ArtistHandler.Create)
		musicRoutes.PUT("/artists/:id", handlers.Media.ArtistHandler.Update)
		musicRoutes.DELETE("/artists/:id", handlers.Media.ArtistHandler.Delete)
	}
	
	// Playlists
	playlistRoutes := router.Group("/playlists")
	{
		playlistRoutes.GET("/:id", handlers.Media.PlaylistHandler.GetByID)
		playlistRoutes.POST("", handlers.Media.PlaylistHandler.Create)
		playlistRoutes.PUT("/:id", handlers.Media.PlaylistHandler.Update)
		playlistRoutes.DELETE("/:id", handlers.Media.PlaylistHandler.Delete)
		
		// Playlist items management
		playlistRoutes.GET("/:id/items", handlers.Media.PlaylistHandler.GetItemsByListID)
		playlistRoutes.POST("/:id/items", handlers.Media.PlaylistHandler.AddItem)
		playlistRoutes.DELETE("/:id/items/:itemId", handlers.Media.PlaylistHandler.RemoveItem)
		
		// Playlist reordering
		playlistRoutes.POST("/:id/reorder", handlers.Media.PlaylistHandler.ReorderItems)
	}
	
	// Collections
	collectionRoutes := router.Group("/collections")
	{
		collectionRoutes.GET("/:id", handlers.Media.CollectionHandler.GetByID)
		collectionRoutes.POST("", handlers.Media.CollectionHandler.Create)
		collectionRoutes.PUT("/:id", handlers.Media.CollectionHandler.Update)
		collectionRoutes.DELETE("/:id", handlers.Media.CollectionHandler.Delete)
		
		// Collection items management
		collectionRoutes.GET("/:id/items", handlers.Media.CollectionHandler.GetItemsByListID)
	}
	
	// User Media Data routes
	userDataRoutes := router.Group("/user/data")
	{
		// Movie data
		userDataRoutes.GET("/movies/:id", handlers.MediaData.MovieDataHandler.GetMediaItemDataByID)
		
		// Series data
		userDataRoutes.GET("/series/:id", handlers.MediaData.SeriesDataHandler.GetMediaItemDataByID)
		
		// Episode data
		userDataRoutes.GET("/episodes/:id", handlers.MediaData.EpisodeDataHandler.GetMediaItemDataByID)
		
		// Track data
		userDataRoutes.GET("/tracks/:id", handlers.MediaData.TrackDataHandler.GetMediaItemDataByID)
		
		// Album data
		userDataRoutes.GET("/albums/:id", handlers.MediaData.AlbumDataHandler.GetMediaItemDataByID)
		
		// Artist data
		userDataRoutes.GET("/artists/:id", handlers.MediaData.ArtistDataHandler.GetMediaItemDataByID)
		
		// Playlist data
		userDataRoutes.GET("/playlists/:id", handlers.MediaData.PlaylistDataHandler.GetMediaItemDataByID)
		
		// Collection data
		userDataRoutes.GET("/collections/:id", handlers.MediaData.CollectionDataHandler.GetMediaItemDataByID)
		
		// Media play history
		userDataRoutes.GET("/history", handlers.MediaData.MovieDataHandler.GetMediaPlayHistory)
		
		// Continue watching
		userDataRoutes.GET("/continue", handlers.MediaData.MovieDataHandler.GetContinuePlaying)
		
		// Recent history
		userDataRoutes.GET("/recent", handlers.MediaData.MovieDataHandler.GetRecentHistory)
		
		// Favorites
		userDataRoutes.GET("/favorites", handlers.MediaData.MovieDataHandler.GetFavorites)
		
		// Record play
		userDataRoutes.POST("/record", handlers.MediaData.MovieDataHandler.RecordMediaPlay)
		
		// Clear history
		userDataRoutes.DELETE("/clear", handlers.MediaData.MovieDataHandler.ClearUserHistory)
	}
}

// InitializeApplication is a central point for initializing the entire application
func InitializeApplication() (*gin.Engine, error) {
	// Create a new Gin router
	router := gin.Default()
	
	// Set up routes using Wire-generated handlers
	SetupRoutes(router)
	
	return router, nil
}