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
		registerListRoutes[*mediatypes.Playlist](clientGroup, c)
		registerListRoutes[*mediatypes.Collection](clientGroup, c)

	}
}

// registerListRoutes sets up routes for a specific media resource type
func registerListRoutes[T mediatypes.ListData](rg *gin.RouterGroup, c *container.Container) {

	mediaType := mediatypes.GetMediaType[T]()

	// /playlists or /collections
	listGroup := rg.Group("/" + string(mediaType))
	{
		// Get all lists
		listGroup.GET("", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetClientLists(g)
			}
		})
		listGroup.GET("/:listID", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetClientListByID(g)
			}
		})
		listGroup.GET("/genre/:genre", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetClientListsByGenre(g)
			}
		})
		listGroup.GET("/year/:year", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetClientListsByYear(g)
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
				handler.GetClientListsByRating(g)
			}
		})
		listGroup.GET("/latest/:count", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetClientLatestListsByAdded(g)
			}
		})
		listGroup.GET("/popular/:count", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetClientPopularLists(g)
			}
		})
		listGroup.GET("/top-rated/:count", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetClientTopRatedLists(g)
			}
		})
		listGroup.GET("/search", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.SearchClientLists(g)
			}
		})

		// Get items in a list
		listGroup.GET("/:listID/items", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.GetClientListItems(g)
			}
		})

	}

	syncGroup := rg.Group("/sync/" + string(mediaType))
	{
		syncGroup.POST("/:listID", func(g *gin.Context) {
			if handler := getHandler[T](g, c); handler != nil {
				handler.SyncLocalListToClient(g)
			}
		})
	}
}

func getHandlerMap[T mediatypes.ListData](c *container.Container, clientType clienttypes.ClientType) (handlers.ClientListHandler[clienttypes.ClientMediaConfig, T], bool) {
	handlers := map[clienttypes.ClientType]handlers.ClientListHandler[clienttypes.ClientMediaConfig, T]{
		clienttypes.ClientTypeEmby:     container.MustGet[handlers.ClientListHandler[*clienttypes.EmbyConfig, T]](c),
		clienttypes.ClientTypeJellyfin: container.MustGet[handlers.ClientListHandler[*clienttypes.JellyfinConfig, T]](c),
		clienttypes.ClientTypePlex:     container.MustGet[handlers.ClientListHandler[*clienttypes.PlexConfig, T]](c),
		clienttypes.ClientTypeSubsonic: container.MustGet[handlers.ClientListHandler[*clienttypes.SubsonicConfig, T]](c),
	}
	handler, exists := handlers[clientType]
	return handler, exists
}

func getHandler[T mediatypes.ListData](g *gin.Context, c *container.Container) handlers.ClientListHandler[clienttypes.ClientMediaConfig, mediatypes.ListData] {
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

	handler, exists := getHandlerMap[T](c, clientType)
	if !exists {
		err := fmt.Errorf("unsupported client type: %s", clientType)
		responses.RespondBadRequest(g, err, "Unsupported client type")
		return nil
	}
	return handler
}
