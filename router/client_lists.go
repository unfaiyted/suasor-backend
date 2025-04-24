package router

import (
	"context"

	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/router/middleware"
	"suasor/types/responses"
)

// {baseURL}/api/v1/client/{clientID}/playlist/{action}
// We are using the clientID to get the clientType with middleware.
// The actions are where the handler methods are defined below.

// RegisterClientListRoutes registers all client-related routes with the gin router
func RegisterClientListRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	// Client resource group
	clientGroup := rg.Group("/:clientID")
	// Gets the client type based on the clientID
	clientGroup.Use(middleware.ClientTypeMiddleware(db))
	{
		// Register routes for playlists
		registerMediaRoutes[*mediatypes.Playlist](clientGroup, c)
		registerMediaRoutes[*mediatypes.Collection](clientGroup, c)
	}
}

// registerMediaRoutes sets up routes for a specific media resource type
func registerMediaRoutes[T mediatypes.ListData](rg *gin.RouterGroup, c *container.Container) {

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	// /playlists or /collections
	listGroup := rg.Group("/" + string(mediaType))
	{
		// Get all lists
		listGroup.GET("", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetAll(g)
			}
		})
		listGroup.GET("/:listID", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetByID(g)
			}
		})
		listGroup.GET("/genre/:genre", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetByGenre(g)
			}
		})
		listGroup.GET("/year/:year", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetByYear(g)
			}
		})
		// listGroup.GET("/actor/:actor", func(g *gin.Context) {
		// 	if handler := getHandler[T](g, c); handler != nil {
		// 		handler.GetByActor(g)
		// 	}
		// })
		// listGroup.GET("/creator/:creator", func(g *gin.Context) {
		// 	if handler := getHandler[T](g, c); handler != nil {
		// 		handler.GetByCreator(g)
		// 	}
		// })
		listGroup.GET("/rating", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetByRating(g)
			}
		})
		listGroup.GET("/latest/:count", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetLatestListsByAdded(g)
			}
		})
		listGroup.GET("/popular/:count", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetPopularLists(g)
			}
		})
		listGroup.GET("/top-rated/:count", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetTopRatedLists(g)
			}
		})
		listGroup.GET("/search", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.SearchLists(g)
			}
		})

		// Get items in a list - use specialized handler
		listGroup.GET("/:listID/items", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetItemsByListID(g)
			}
		})

	}

}

func getHandlerMap[T mediatypes.ListData](c *container.Container, clientType string) (handlers.ClientListHandler[clienttypes.ClientMediaConfig, T], bool) {
	handlers := map[string]handlers.ClientListHandler[clienttypes.ClientMediaConfig, T]{
		"emby":     container.MustGet[handlers.ClientListHandler[*clienttypes.EmbyConfig, T]](c),
		"jellyfin": container.MustGet[handlers.ClientListHandler[*clienttypes.JellyfinConfig, T]](c),
		"plex":     container.MustGet[handlers.ClientListHandler[*clienttypes.PlexConfig, T]](c),
		"subsonic": container.MustGet[handlers.ClientListHandler[*clienttypes.SubsonicConfig, T]](c),
	}
	handler, exists := handlers[clientType]
	return handler, exists
}

func getHandler[T mediatypes.ListData](g *gin.Context, c *container.Container) handlers.ClientListHandler[clienttypes.ClientMediaConfig, mediatypes.ListData] {
	clientType := g.Param("clientType")
	handler, exists := getHandlerMap[T](c, clientType)
	if !exists {
		err := fmt.Errorf("unsupported client type: %s", clientType)
		responses.RespondBadRequest(g, err, "Unsupported client type")
		return nil
	}
	return handler
}
