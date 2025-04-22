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

	clientGroup := rg.Group("/client/:clientId")
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

	// Core routes
	// rg.POST("/sync", clientDataHandlers.SyncClientItemData)
	rg.GET("", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			handler.GetAllClientItems(g)
		}
	})
	rg.GET("/item/:clientItemId", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			handler.GetClientItemByItemID(g)
		}
	})
	rg.POST("/item/:clientItemId/play", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.RecordClientPlay(g)
		}
	})
	rg.GET("/item/:clientItemId/state", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.GetPlaybackState(g)
		}
	})
	rg.PUT("/item/:clientItemId/state", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.UpdatePlaybackState(g)
		}
	})
	rg.DELETE("/item/:clientItemId", func(g *gin.Context) {
		if handler := getItemHandler[T](g, c); handler != nil {
			// handler.DeleteMediaItemData(g)
		}
	})
	rg.GET("/sync", func(g *gin.Context) {
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
