package router

import (
	"suasor/app/container"
	mediatypes "suasor/client/media/types"
	"suasor/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterMediaItemDataRoutes configures routes for user media item data
func RegisterMediaItemDataRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Base routes for all media types
	userMediaData := rg.Group("/media-data")
	{
		registerDataRoutes[mediatypes.Movie](userMediaData, c)
		registerDataRoutes[mediatypes.Series](userMediaData, c)
		registerDataRoutes[mediatypes.Track](userMediaData, c)
		registerDataRoutes[mediatypes.Album](userMediaData, c)
		registerDataRoutes[mediatypes.Artist](userMediaData, c)

		// Media-type specific routes
		registerMovieUserMediaItemDataRoutes(userMediaData,
			container.MustGet[handlers.UserMediaItemDataHandler[*mediatypes.Movie]](c))
		registerSeriesUserMediaItemDataRoutes(userMediaData,
			container.MustGet[handlers.UserMediaItemDataHandler[*mediatypes.Series]](c))
		registerMusicUserMediaItemDataRoutes(userMediaData,
			container.MustGet[handlers.UserMediaItemDataHandler[*mediatypes.Track]](c))
	}

	clientMediaData := rg.Group("/client/:clientId/media-data")
	{
		registerClientDataRoutes[mediatypes.Movie](clientMediaData, c)
		registerClientDataRoutes[mediatypes.Series](clientMediaData, c)
		registerClientDataRoutes[mediatypes.Track](clientMediaData, c)
		registerClientDataRoutes[mediatypes.Album](clientMediaData, c)
		registerClientDataRoutes[mediatypes.Artist](clientMediaData, c)

		// Specialized routes
		// TODO: consider adding these extra routes.
		// registerMusic
	}
}

func registerDataRoutes[T mediatypes.MediaData](rg *gin.RouterGroup, c *container.Container) {
	// Get handlers
	userDataHandlers := container.MustGet[handlers.UserMediaItemDataHandler[T]](c)

	// Core routes
	rg.GET("/:id", userDataHandlers.GetMediaItemDataByID)
	rg.GET("/check", userDataHandlers.CheckUserMediaItemData)
	rg.GET("/user-media", userDataHandlers.GetMediaItemDataByUserAndMedia)
	rg.DELETE("/:id", userDataHandlers.DeleteMediaItemData)
	rg.GET("/history", userDataHandlers.GetMediaPlayHistory)
	rg.GET("/continue-watching", userDataHandlers.GetContinuePlaying)
	rg.GET("/recent", userDataHandlers.GetRecentHistory)
	rg.POST("/record", userDataHandlers.RecordMediaPlay)
	rg.PUT("/media/:mediaItemId/favorite", userDataHandlers.ToggleFavorite)
	rg.PUT("/media/:mediaItemId/rating", userDataHandlers.UpdateUserRating)
	rg.GET("/favorites", userDataHandlers.GetFavorites)
	rg.DELETE("/clear", userDataHandlers.ClearUserHistory)

}

// registerMovieUserMediaItemDataRoutes configures movie-specific routes
func registerMovieUserMediaItemDataRoutes(rg *gin.RouterGroup, handler handlers.UserMediaItemDataHandler[*mediatypes.Movie]) {
	movies := rg.Group("/movies")
	{
		// Get movie history
		movies.GET("/history", handler.GetMediaPlayHistory)

		// Get favorite movies
		movies.GET("/favorites", handler.GetFavorites)
	}
}

// registerSeriesUserMediaItemDataRoutes configures series-specific routes
func registerSeriesUserMediaItemDataRoutes(rg *gin.RouterGroup, handler handlers.UserMediaItemDataHandler[*mediatypes.Series]) {
	series := rg.Group("/series")
	{
		// Get series history
		series.GET("/history", handler.GetMediaPlayHistory)

		// Get favorite series
		series.GET("/favorites", handler.GetFavorites)
	}
}

// registerMusicUserMediaItemDataRoutes configures music-specific routes
func registerMusicUserMediaItemDataRoutes(rg *gin.RouterGroup, handler handlers.UserMediaItemDataHandler[*mediatypes.Track]) {
	music := rg.Group("/music")
	{
		// Get music history
		music.GET("/history", handler.GetMediaPlayHistory)

		// Get favorite music
		music.GET("/favorites", handler.GetFavorites)
	}
}

func registerClientDataRoutes[T mediatypes.MediaData](rg *gin.RouterGroup, c *container.Container) {
	// Get handlers
	clientDataHandlers := container.MustGet[handlers.ClientUserMediaItemDataHandler[T]](c)

	// Core routes
	rg.POST("/sync", clientDataHandlers.SyncClientItemData)
	rg.GET("", clientDataHandlers.GetClientItemData)
	rg.GET("/item/:clientItemId", clientDataHandlers.GetMediaItemDataByClientID)
	rg.POST("/item/:clientItemId/play", clientDataHandlers.RecordClientPlay)
	rg.GET("/item/:clientItemId/state", clientDataHandlers.GetPlaybackState)
	rg.PUT("/item/:clientItemId/state", clientDataHandlers.UpdatePlaybackState)
}
