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
		registerClientItemRoutes[mediatypes.Movie](clientGroup, c)
		registerClientItemRoutes[mediatypes.Series](clientGroup, c)
		registerClientItemRoutes[mediatypes.Track](clientGroup, c)
		registerClientItemRoutes[mediatypes.Album](clientGroup, c)
		registerClientItemRoutes[mediatypes.Artist](clientGroup, c)
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

}

func getItemHandlerMap[T mediatypes.MediaData](c *container.Container, clientType string) (
	handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, T], bool) {

	handlers := map[string]handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, T]{
		"emby":     container.MustGet[handlers.ClientMediaItemHandler[*clienttypes.EmbyConfig, T]](c),
		"jellyfin": container.MustGet[handlers.ClientMediaItemHandler[*clienttypes.JellyfinConfig, T]](c),
		"plex":     container.MustGet[handlers.ClientMediaItemHandler[*clienttypes.PlexConfig, T]](c),
		"subsonic": container.MustGet[handlers.ClientMediaItemHandler[*clienttypes.SubsonicConfig, T]](c),
	}
	handler, exists := handlers[clientType]
	return handler, exists
}

func getItemHandler[T mediatypes.MediaData](g *gin.Context, c *container.Container) handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, mediatypes.MediaData] {
	clientType := g.Param("clientType")
	handler, exists := getItemHandlerMap[T](c, clientType)
	if !exists {
		err := fmt.Errorf("unsupported client type: %s", clientType)
		responses.RespondBadRequest(g, err, "Unsupported client type")
		return nil
	}
	return handler
}
