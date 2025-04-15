package router

import (
	"suasor/app"
	"suasor/handlers"
	"suasor/client/media/types"

	"github.com/gin-gonic/gin"
)

// UserMediaItemDataHandlerInterface defines the interface for user media item data handlers
type UserMediaItemDataHandlerInterface interface {
	// Core methods
	GetMediaItemDataByID(c *gin.Context)
	CheckUserMediaItemData(c *gin.Context)
	GetMediaItemDataByUserAndMedia(c *gin.Context)
	DeleteMediaItemData(c *gin.Context)
	
	// User methods
	GetMediaPlayHistory(c *gin.Context)
	GetContinueWatching(c *gin.Context)
	GetRecentHistory(c *gin.Context)
	RecordMediaPlay(c *gin.Context)
	ToggleFavorite(c *gin.Context)
	UpdateUserRating(c *gin.Context)
	GetFavorites(c *gin.Context)
	ClearUserHistory(c *gin.Context)
	
	// Client methods
	SyncClientItemData(c *gin.Context)
	GetClientItemData(c *gin.Context)
	GetMediaItemDataByClientID(c *gin.Context)
	RecordClientPlay(c *gin.Context)
	GetPlaybackState(c *gin.Context)
	UpdatePlaybackState(c *gin.Context)
}

// RegisterUserMediaItemDataRoutes configures routes for user media item data
func RegisterUserMediaItemDataRoutes(rg *gin.RouterGroup, deps *app.AppDependencies) {
	// Get handlers
	mediaHandlers := deps.UserMediaItemDataHandlers
	
	handlerMap := map[string]UserMediaItemDataHandlerInterface{
		"movies": mediaHandlers.MovieHandler(),
		"series": mediaHandlers.SeriesHandler(),
		"tracks": mediaHandlers.TrackHandler(),
		"albums": mediaHandlers.AlbumHandler(),
		"artists": mediaHandlers.ArtistHandler(),
		"episodes": mediaHandlers.EpisodeHandler(),
	}
	
	getHandler := func(c *gin.Context) UserMediaItemDataHandlerInterface {
		mediaType := c.Param("mediaType")
		handler, exists := handlerMap[mediaType]
		if !exists {
			// Default to movie handler if type not specified or invalid
			return mediaHandlers.MovieHandler()
		}
		return handler
	}
	
	// Base routes for all media types
	userMediaData := rg.Group("/user-media-data")
	
	// Core routes
	userMediaData.GET("/:id", func(c *gin.Context) {
		getHandler(c).GetMediaItemDataByID(c)
	})
	userMediaData.GET("/check", func(c *gin.Context) {
		getHandler(c).CheckUserMediaItemData(c)
	})
	userMediaData.GET("/user-media", func(c *gin.Context) {
		getHandler(c).GetMediaItemDataByUserAndMedia(c)
	})
	userMediaData.DELETE("/:id", func(c *gin.Context) {
		getHandler(c).DeleteMediaItemData(c)
	})
	
	// User-specific routes
	userMediaData.GET("/history", func(c *gin.Context) {
		getHandler(c).GetMediaPlayHistory(c)
	})
	userMediaData.GET("/continue-watching", func(c *gin.Context) {
		getHandler(c).GetContinueWatching(c)
	})
	userMediaData.GET("/recent", func(c *gin.Context) {
		getHandler(c).GetRecentHistory(c)
	})
	userMediaData.POST("/record", func(c *gin.Context) {
		getHandler(c).RecordMediaPlay(c)
	})
	userMediaData.PUT("/media/:mediaItemId/favorite", func(c *gin.Context) {
		getHandler(c).ToggleFavorite(c)
	})
	userMediaData.PUT("/media/:mediaItemId/rating", func(c *gin.Context) {
		getHandler(c).UpdateUserRating(c)
	})
	userMediaData.GET("/favorites", func(c *gin.Context) {
		getHandler(c).GetFavorites(c)
	})
	userMediaData.DELETE("/clear", func(c *gin.Context) {
		getHandler(c).ClearUserHistory(c)
	})
	
	// Client-specific routes
	clientData := userMediaData.Group("/client/:clientId")
	{
		clientData.POST("/sync", func(c *gin.Context) {
			getHandler(c).SyncClientItemData(c)
		})
		clientData.GET("", func(c *gin.Context) {
			getHandler(c).GetClientItemData(c)
		})
		clientData.GET("/item/:clientItemId", func(c *gin.Context) {
			getHandler(c).GetMediaItemDataByClientID(c)
		})
		clientData.POST("/item/:clientItemId/play", func(c *gin.Context) {
			getHandler(c).RecordClientPlay(c)
		})
		clientData.GET("/item/:clientItemId/state", func(c *gin.Context) {
			getHandler(c).GetPlaybackState(c)
		})
		clientData.PUT("/item/:clientItemId/state", func(c *gin.Context) {
			getHandler(c).UpdatePlaybackState(c)
		})
	}
	
	// Media-type specific routes
	registerMovieUserMediaItemDataRoutes(userMediaData, mediaHandlers.MovieHandler())
	registerSeriesUserMediaItemDataRoutes(userMediaData, mediaHandlers.SeriesHandler())
	registerMusicUserMediaItemDataRoutes(userMediaData, mediaHandlers.TrackHandler())
}

// registerMovieUserMediaItemDataRoutes configures movie-specific routes
func registerMovieUserMediaItemDataRoutes(rg *gin.RouterGroup, handler UserMediaItemDataHandlerInterface) {
	movies := rg.Group("/movies")
	{
		// Get movie history
		movies.GET("/history", handler.GetMediaPlayHistory)
		
		// Get continue watching movies
		movies.GET("/continue-watching", handler.GetContinueWatching)
		
		// Get favorite movies
		movies.GET("/favorites", handler.GetFavorites)
	}
}

// registerSeriesUserMediaItemDataRoutes configures series-specific routes
func registerSeriesUserMediaItemDataRoutes(rg *gin.RouterGroup, handler UserMediaItemDataHandlerInterface) {
	series := rg.Group("/series")
	{
		// Get series history
		series.GET("/history", handler.GetMediaPlayHistory)
		
		// Get continue watching series
		series.GET("/continue-watching", handler.GetContinueWatching)
		
		// Get favorite series
		series.GET("/favorites", handler.GetFavorites)
	}
}

// registerMusicUserMediaItemDataRoutes configures music-specific routes
func registerMusicUserMediaItemDataRoutes(rg *gin.RouterGroup, handler UserMediaItemDataHandlerInterface) {
	music := rg.Group("/music")
	{
		// Get music history
		music.GET("/history", handler.GetMediaPlayHistory)
		
		// Get favorite music
		music.GET("/favorites", handler.GetFavorites)
	}
}