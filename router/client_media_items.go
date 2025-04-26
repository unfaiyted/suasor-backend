package router

import (
	"context"
	"fmt"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/router/middleware"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterMediaItemRoutes configures routes for user media item data
func RegisterClientMediaItemRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {
	// Base routes for all media types
	db := container.MustGet[*gorm.DB](c)

	clientGroup := rg.Group("/:clientID/media")
	clientGroup.Use(middleware.ClientTypeMiddleware(db))
	{
		registerClientItemRoutes[*mediatypes.Movie](clientGroup, c)
		registerClientItemRoutes[*mediatypes.Series](clientGroup, c)
		registerClientItemRoutes[*mediatypes.Track](clientGroup, c)
		registerClientItemRoutes[*mediatypes.Album](clientGroup, c)
		registerClientItemRoutes[*mediatypes.Artist](clientGroup, c)
	}
}

func registerClientItemRoutes[T mediatypes.MediaData](rg *gin.RouterGroup, c *container.Container) {

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	// Print debug info
	fmt.Printf("Registering routes for media type: %T -> %s\n", zero, mediaType)

	// Don't register routes if media type is unknown
	if mediaType == mediatypes.MediaTypeUnknown {
		fmt.Printf("WARNING: Not registering routes for unknown media type: %T\n", zero)
		return
	}

	itemGroup := rg.Group("/" + string(mediaType))

	// Core routes
	// rg.POST("/sync", clientDataHandlers.SyncClientItemData)
	itemGroup.GET("", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			handler.GetAllClientItems(g)
		}
	})
	itemGroup.GET("/:clientItemID", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			handler.GetClientItemByItemID(g)
		}
	})
	itemGroup.POST("/:clientItemID/play", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.RecordClientPlay(g)
		}
	})
	itemGroup.GET("/:clientItemID/state", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.GetPlaybackState(g)
		}
	})
	itemGroup.PUT("/:clientItemID/state", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.UpdatePlaybackState(g)
		}
	})
	itemGroup.DELETE("/:clientItemID", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.DeleteMediaItemData(g)
		}
	})
	itemGroup.GET("/sync", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.SyncClientItemData(g)
		}
	})
	itemGroup.GET("/search", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			handler.SearchClient(g)
		}
	})

}

func getItemHandlerMap[T mediatypes.MediaData](ctx context.Context, c *container.Container, clientType clienttypes.ClientType) (
	handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, T], bool) {

	handlerMap := map[clienttypes.ClientType]handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, T]{}

	if CheckClientSupportsMediaType[*clienttypes.EmbyConfig, T](ctx) {
		handlerMap[clienttypes.ClientTypeEmby] = container.MustGet[handlers.ClientMediaItemHandler[*clienttypes.EmbyConfig, T]](c)
	}
	if CheckClientSupportsMediaType[*clienttypes.JellyfinConfig, T](ctx) {
		handlerMap[clienttypes.ClientTypeJellyfin] = container.MustGet[handlers.ClientMediaItemHandler[*clienttypes.JellyfinConfig, T]](c)
	}
	if CheckClientSupportsMediaType[*clienttypes.PlexConfig, T](ctx) {
		handlerMap[clienttypes.ClientTypePlex] = container.MustGet[handlers.ClientMediaItemHandler[*clienttypes.PlexConfig, T]](c)
	}
	if CheckClientSupportsMediaType[*clienttypes.SubsonicConfig, T](ctx) {
		handlerMap[clienttypes.ClientTypeSubsonic] = container.MustGet[handlers.ClientMediaItemHandler[*clienttypes.SubsonicConfig, T]](c)
	}

	fmt.Printf("HandlerMap: %v\n", handlerMap)
	fmt.Printf("ClientType: %v\n", clientType)
	handler, exists := handlerMap[clientType]
	return handler, exists
}

func getItemHandler[T mediatypes.MediaData](g *gin.Context, c *container.Container) handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, mediatypes.MediaData] {
	log := logger.LoggerFromContext(g.Request.Context())
	clientTypeStr, exists := g.Get("clientType")
	if !exists {
		log.Warn().Msg("Client type not found in request context")
		responses.RespondBadRequest(g, nil, "Client type not found")
		return nil
	}
	clientType := clientTypeStr.(clienttypes.ClientType)
	log.Debug().Str("clientType", string(clientType)).Msg("Getting client media item handler")
	handler, exists := getItemHandlerMap[T](g.Request.Context(), c, clientType)
	if !exists {
		err := fmt.Errorf("unsupported client type: %s", clientType)
		responses.RespondBadRequest(g, err, "Unsupported client type")
		return nil
	}
	return handler
}

// CheckSupportsMediaType

func CheckClientSupportsMediaType[T clienttypes.ClientMediaConfig, U mediatypes.MediaData](ctx context.Context) bool {
	var zero U
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	log := logger.LoggerFromContext(ctx)
	log.Debug().Str("mediaType", string(mediaType)).Msg("Checking client support for media type")

	// Create an instance of the client config
	var instance T
	switch any(instance).(type) {
	case *clienttypes.EmbyConfig:
		config := clienttypes.NewEmbyConfig("", "", "", "", true, true)
		return config.SupportsMediaType(mediaType)
	case *clienttypes.JellyfinConfig:
		config := clienttypes.NewJellyfinConfig("", "", "", "", true, true)
		return config.SupportsMediaType(mediaType)
	case *clienttypes.PlexConfig:
		config := clienttypes.NewPlexConfig("", "", true, true)
		return config.SupportsMediaType(mediaType)
	case *clienttypes.SubsonicConfig:
		config := clienttypes.NewSubsonicConfig("", "", "", true, true)
		return config.SupportsMediaType(mediaType)
	default:
		log.Warn().Msgf("Unsupported client config type: %T", instance)
		return false
	}
}
