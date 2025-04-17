package router

import (
	"fmt"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	mediatypes "suasor/client/media/types"
	"suasor/handlers"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
)

// // UserMediaItemDataHandlerInterface defines the interface for user media item data handlers
// type UserMediaItemDataHandlerInterface interface {
// 	// Core methods
// 	GetMediaItemDataByID(c *gin.Context)
// 	CheckUserMediaItemData(c *gin.Context)
// 	GetMediaItemDataByUserAndMedia(c *gin.Context)
// 	DeleteMediaItemData(c *gin.Context)
//
// 	// User methods
// 	GetMediaPlayHistory(c *gin.Context)
// 	GetContinueWatching(c *gin.Context)
// 	GetRecentHistory(c *gin.Context)
// 	RecordMediaPlay(c *gin.Context)
// 	ToggleFavorite(c *gin.Context)
// 	UpdateUserRating(c *gin.Context)
// 	GetFavorites(c *gin.Context)
// 	ClearUserHistory(c *gin.Context)
//
// 	// Client methods
// 	SyncClientItemData(c *gin.Context)
// 	GetClientItemData(c *gin.Context)
// 	GetMediaItemDataByClientID(c *gin.Context)
// 	RecordClientPlay(c *gin.Context)
// 	GetPlaybackState(c *gin.Context)
// 	UpdatePlaybackState(c *gin.Context)
// }

type UserDataHandlerInterface interface {
	GetMediaItemDataByID(c *gin.Context)
	CheckUserMediaItemData(c *gin.Context)
	GetMediaItemDataByUserAndMedia(c *gin.Context)
	DeleteMediaItemData(c *gin.Context)

	GetMediaPlayHistory(c *gin.Context)
	GetContinuePlaying(c *gin.Context)
	GetRecentHistory(c *gin.Context)
	RecordMediaPlay(c *gin.Context)
	ToggleFavorite(c *gin.Context)
	UpdateUserRating(c *gin.Context)
	GetFavorites(c *gin.Context)
	ClearUserHistory(c *gin.Context)
}

type ClientDataHandlerInterface interface {
	SyncClientItemData(c *gin.Context)
	GetClientItemData(c *gin.Context)
	GetMediaItemDataByClientID(c *gin.Context)
	RecordClientPlay(c *gin.Context)
	GetPlaybackState(c *gin.Context)
	UpdatePlaybackState(c *gin.Context)
}

// RegisterMediaItemDataRoutes configures routes for user media item data
func RegisterMediaItemDataRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Get handlers
	userDataHandlers := container.MustGet[apphandlers.UserMediaItemDataHandlers](c)
	clientDataHandlers := container.MustGet[apphandlers.ClientMediaItemDataHandlers](c)

	// Base routes for all media types
	userMediaData := rg.Group("/media-data")

	handlerMap := map[string]UserDataHandlerInterface{
		"movies":   userDataHandlers.MovieUserDataHandler(),
		"series":   userDataHandlers.SeriesUserDataHandler(),
		"seasons":  userDataHandlers.SeasonUserDataHandler(),
		"episodes": userDataHandlers.EpisodeUserDataHandler(),
		"tracks":   userDataHandlers.TrackUserDataHandler(),
		"albums":   userDataHandlers.AlbumUserDataHandler(),
		"artists":  userDataHandlers.ArtistUserDataHandler(),
	}

	clientHandlerMap := map[string]ClientDataHandlerInterface{
		"movies":   clientDataHandlers.MovieClientDataHandler(),
		"series":   clientDataHandlers.SeriesClientDataHandler(),
		"seasons":  clientDataHandlers.SeasonClientDataHandler(),
		"episodes": clientDataHandlers.EpisodeClientDataHandler(),
		"tracks":   clientDataHandlers.TrackClientDataHandler(),
		"albums":   clientDataHandlers.AlbumClientDataHandler(),
		"artists":  clientDataHandlers.ArtistClientDataHandler(),
	}

	getHandler := func(c *gin.Context) UserDataHandlerInterface {

		clientType := c.Param("clientType")
		handler, exists := handlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil

		}
		return handler
	}

	getClientHandler := func(c *gin.Context) ClientDataHandlerInterface {
		clientType := c.Param("clientType")
		handler, exists := clientHandlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}
		return handler
	}

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
	userMediaData.GET("/continue-playing", func(c *gin.Context) {
		getHandler(c).GetContinuePlaying(c)
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
			getClientHandler(c).SyncClientItemData(c)
		})
		clientData.GET("", func(c *gin.Context) {
			getClientHandler(c).GetClientItemData(c)
		})
		clientData.GET("/item/:clientItemId", func(c *gin.Context) {
			getClientHandler(c).GetMediaItemDataByClientID(c)
		})
		clientData.POST("/item/:clientItemId/play", func(c *gin.Context) {
			getClientHandler(c).RecordClientPlay(c)
		})
		clientData.GET("/item/:clientItemId/state", func(c *gin.Context) {
			getClientHandler(c).GetPlaybackState(c)
		})
		clientData.PUT("/item/:clientItemId/state", func(c *gin.Context) {
			getClientHandler(c).UpdatePlaybackState(c)
		})
	}

	// Media-type specific routes
	registerMovieUserMediaItemDataRoutes(userMediaData,
		container.MustGet[handlers.UserMediaItemDataHandler[*mediatypes.Movie]](c))
	registerSeriesUserMediaItemDataRoutes(userMediaData,
		container.MustGet[handlers.UserMediaItemDataHandler[*mediatypes.Series]](c))
	registerMusicUserMediaItemDataRoutes(userMediaData,
		container.MustGet[handlers.UserMediaItemDataHandler[*mediatypes.Track]](c))
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
