package router

import (
	"fmt"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
)

// RegisterMediaItemDataRoutes configures routes for user media item data
func RegisterMediaItemDataRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Base routes for all media types
	userMediaData := rg.Group("/user-data")
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

	clientMediaData := rg.Group("/client/:clientID/user-data")
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
	dataGroup.GET("/data/:dataId", userDataHandlers.GetMediaItemDataByID)

	// Primariy working with mediaItemID (id) to keep things simple
	dataGroup.GET("/:itemID/check", userDataHandlers.CheckUserMediaItemData)
	dataGroup.GET("/:itemID", userDataHandlers.GetUserMediaItemDataByItemID)
	dataGroup.DELETE("/:itemID", userDataHandlers.DeleteMediaItemData)

	dataGroup.GET("/history", userDataHandlers.GetMediaPlayHistory)
	dataGroup.GET("/continue-watching", userDataHandlers.GetContinuePlaying)
	dataGroup.GET("/recent", userDataHandlers.GetRecentHistory)
	dataGroup.POST(":itemID/record", userDataHandlers.RecordMediaPlay)
	dataGroup.PUT("/:itemID/favorite", userDataHandlers.ToggleFavorite)
	dataGroup.PUT("/:itemID/rating", userDataHandlers.UpdateUserRating)
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

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	clientDataGroup := rg.Group("/" + string(mediaType))

	// Core routes
	clientDataGroup.POST("/sync", func(g *gin.Context) {
		if handler := getClientDataHandler[T](g, c); handler != nil {
			handler.SyncClientItemData(g)
		}
	})
	clientDataGroup.GET("/", func(g *gin.Context) {
		if handler := getClientDataHandler[T](g, c); handler != nil {
			handler.GetClientItemData(g)
		}
	})
	clientDataGroup.GET("/:clientItemID", func(g *gin.Context) {
		if handler := getClientDataHandler[T](g, c); handler != nil {
			handler.GetMediaItemDataByClientID(g)
		}
	})
	clientDataGroup.POST("/:clientItemID/play", func(g *gin.Context) {
		if handler := getClientDataHandler[T](g, c); handler != nil {
			handler.RecordClientPlay(g)
		}
	})
	clientDataGroup.GET("/item/:clientItemID/state", func(g *gin.Context) {
		if handler := getClientDataHandler[T](g, c); handler != nil {
			handler.GetPlaybackState(g)
		}
	})
	clientDataGroup.PUT("/item/:clientItemID/state", func(g *gin.Context) {
		if handler := getClientDataHandler[T](g, c); handler != nil {
			handler.UpdatePlaybackState(g)
		}
	})

}

func getClientDataHandlerMap[T mediatypes.MediaData](c *container.Container, clientType clienttypes.ClientType) (handlers.ClientUserMediaItemDataHandler[clienttypes.ClientMediaConfig, T], bool) {
	handlers := map[clienttypes.ClientType]handlers.ClientUserMediaItemDataHandler[clienttypes.ClientMediaConfig, T]{
		clienttypes.ClientTypeEmby:     container.MustGet[handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, T]](c),
		clienttypes.ClientTypeJellyfin: container.MustGet[handlers.ClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, T]](c),
		clienttypes.ClientTypePlex:     container.MustGet[handlers.ClientUserMediaItemDataHandler[*clienttypes.PlexConfig, T]](c),
		clienttypes.ClientTypeSubsonic: container.MustGet[handlers.ClientUserMediaItemDataHandler[*clienttypes.SubsonicConfig, T]](c),
	}
	handler, exists := handlers[clientType]
	return handler, exists
}

func getClientDataHandler[T mediatypes.MediaData](g *gin.Context, c *container.Container) handlers.ClientUserMediaItemDataHandler[clienttypes.ClientMediaConfig, T] {
	// Get the client type from the context which was set by the ClientTypeMiddleware
	clientTypeValue, exists := g.Get("clientType")
	if !exists {
		err := fmt.Errorf("client type not found in context")
		responses.RespondBadRequest(g, err, "Client type not found")
		return nil
	}
	
	clientType, ok := clientTypeValue.(clienttypes.ClientType)
	if !ok {
		err := fmt.Errorf("invalid client type format")
		responses.RespondBadRequest(g, err, "Invalid client type format")
		return nil
	}
	
	handler, exists := getClientDataHandlerMap[T](c, clientType)
	if !exists {
		err := fmt.Errorf("unsupported client type: %s", clientType)
		responses.RespondBadRequest(g, err, "Unsupported client type")
		return nil
	}
	return handler
}
