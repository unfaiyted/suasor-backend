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
		registerDataRoutes[*mediatypes.Movie](userMediaData, c)
		registerDataRoutes[*mediatypes.Series](userMediaData, c)
		registerDataRoutes[*mediatypes.Track](userMediaData, c)
		registerDataRoutes[*mediatypes.Album](userMediaData, c)
		registerDataRoutes[*mediatypes.Artist](userMediaData, c)

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
		registerClientDataRoutes[*mediatypes.Movie](clientMediaData, c)
		registerClientDataRoutes[*mediatypes.Series](clientMediaData, c)
		registerClientDataRoutes[*mediatypes.Track](clientMediaData, c)
		registerClientDataRoutes[*mediatypes.Album](clientMediaData, c)
		registerClientDataRoutes[*mediatypes.Artist](clientMediaData, c)

		// Specialized routes
		// TODO: consider adding these extra routes.
		// registerMusic
	}
}

func registerDataRoutes[T mediatypes.MediaData](rg *gin.RouterGroup, c *container.Container) {
	// Get handlers
	userDataHandlers := container.MustGet[handlers.UserMediaItemDataHandler[T]](c)

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	dataGroup := rg.Group("/" + string(mediaType))

	// Core routes
	dataGroup.GET("/:id", userDataHandlers.GetMediaItemDataByID)
	dataGroup.GET("/check", userDataHandlers.CheckUserMediaItemData)
	dataGroup.GET("/user-media", userDataHandlers.GetMediaItemDataByUserAndMedia)
	dataGroup.DELETE("/:id", userDataHandlers.DeleteMediaItemData)
	dataGroup.GET("/history", userDataHandlers.GetMediaPlayHistory)
	dataGroup.GET("/continue-watching", userDataHandlers.GetContinuePlaying)
	dataGroup.GET("/recent", userDataHandlers.GetRecentHistory)
	dataGroup.POST("/record", userDataHandlers.RecordMediaPlay)
	dataGroup.PUT("/media/:mediaItemId/favorite", userDataHandlers.ToggleFavorite)
	dataGroup.PUT("/media/:mediaItemId/rating", userDataHandlers.UpdateUserRating)
	dataGroup.GET("/favorites", userDataHandlers.GetFavorites)
	dataGroup.DELETE("/clear", userDataHandlers.ClearUserHistory)

}

// registerMovieUserMediaItemDataRoutes configures movie-specific routes
func registerMovieUserMediaItemDataRoutes(rg *gin.RouterGroup, handler handlers.UserMediaItemDataHandler[*mediatypes.Movie]) {
	// movies := rg.Group("/movies")
	// {
	// // Get movie history
	// movies.GET("/history", handler.GetMediaPlayHistory)
	//
	// // Get favorite movies
	// movies.GET("/favorites", handler.GetFavorites)
	// }
}

// registerSeriesUserMediaItemDataRoutes configures series-specific routes
func registerSeriesUserMediaItemDataRoutes(rg *gin.RouterGroup, handler handlers.UserMediaItemDataHandler[*mediatypes.Series]) {
	// series := rg.Group("/series")
	// {
	// 	// Get series history
	// 	// series.GET("/history", handler.GetMediaPlayHistory)
	//
	// 	// Get favorite series
	// 	series.GET("/favorites", handler.GetFavorites)
	// }
}

// registerMusicUserMediaItemDataRoutes configures music-specific routes
func registerMusicUserMediaItemDataRoutes(rg *gin.RouterGroup, handler handlers.UserMediaItemDataHandler[*mediatypes.Track]) {
	// music := rg.Group("/music")
	// {
	// 	// Get music history
	// 	music.GET("/history", handler.GetMediaPlayHistory)
	//
	// 	// Get favorite music
	// 	music.GET("/favorites", handler.GetFavorites)
	// }
}

func registerClientDataRoutes[T mediatypes.MediaData](rg *gin.RouterGroup, c *container.Container) {
	// Get handlers
	clientDataHandlers := container.MustGet[handlers.ClientUserMediaItemDataHandler[T]](c)

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	clientDataGroup := rg.Group("/" + string(mediaType))
	// Core routes
	clientDataGroup.POST("/sync", clientDataHandlers.SyncClientItemData)
	clientDataGroup.GET("", clientDataHandlers.GetClientItemData)
	clientDataGroup.GET("/item/:clientItemId", clientDataHandlers.GetMediaItemDataByClientID)
	clientDataGroup.POST("/item/:clientItemId/play", clientDataHandlers.RecordClientPlay)
	clientDataGroup.GET("/item/:clientItemId/state", clientDataHandlers.GetPlaybackState)
	clientDataGroup.PUT("/item/:clientItemId/state", clientDataHandlers.UpdatePlaybackState)
}
