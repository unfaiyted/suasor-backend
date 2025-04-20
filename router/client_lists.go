package router

import (
	"context"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	clienttypes "suasor/client/types"
	"suasor/router/middleware"
	"suasor/types/responses"
)

// {baseURL}/api/v1/client/{clientID}/playlists/{action}
// We are using the clientID to get the clientType with middleware.
// The actions are where the handler methods are defined below.

// HandlerMethod represents a method in a client handler
type HandlerMethod string

// Define constants for all handler methods to avoid string literals
const (
	// Core methods
	GetAll           HandlerMethod = "GetAll"
	GetByID          HandlerMethod = "GetByID"
	GetItemsByListID HandlerMethod = "GetItemsByListID"
	GetByGenre       HandlerMethod = "GetByGenre"
	Search           HandlerMethod = "Search"

	// Client list specific methods
	GetListByID           HandlerMethod = "GetListByID"
	GetListsByGenre       HandlerMethod = "GetListsByGenre"
	GetListsByYear        HandlerMethod = "GetListsByYear"
	GetListsByActor       HandlerMethod = "GetListsByActor"
	GetListsByCreator     HandlerMethod = "GetListsByCreator"
	GetListsByRating      HandlerMethod = "GetListsByRating"
	GetLatestListsByAdded HandlerMethod = "GetLatestListsByAdded"
	GetPopularLists       HandlerMethod = "GetPopularLists"
	GetTopRatedLists      HandlerMethod = "GetTopRatedLists"
	SearchLists           HandlerMethod = "SearchLists"
)

// MediaResourceType identifies the type of media resource
type MediaResourceType string

const (
	PlaylistResource   MediaResourceType = "playlists"
	CollectionResource MediaResourceType = "collections"
)

// RegisterClientListRoutes registers all client-related routes with the gin router
func RegisterClientListRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {
	mediaHandler := container.MustGet[apphandlers.ClientMediaHandlers](c)
	db := container.MustGet[*gorm.DB](c)

	// Client resource group
	clientGroup := rg.Group("/client/:clientID")
	clientGroup.Use(middleware.ClientTypeMiddleware(db))

	// Register routes for playlists
	// playlistsGroup := clientGroup.Group("/playlists")
	registerMediaRoutes(clientGroup, mediaHandler, PlaylistResource)

	// Register routes for collections
	// collectionsGroup := clientGroup.Group("/collections")
	registerMediaRoutes(clientGroup, mediaHandler, CollectionResource)

	// Register other media types as needed
}

// registerMediaRoutes sets up routes for a specific media resource type
func registerMediaRoutes(rg *gin.RouterGroup, mediaHandler apphandlers.ClientMediaHandlers, resourceType MediaResourceType) {
	// Create resource group
	resourceGroup := rg.Group(string(resourceType))

	// Map routes to handler methods
	routeMappings := []struct {
		path    string
		method  string
		handler HandlerMethod
	}{
		// Core routes
		{"", "GET", GetAll},
		{"/:id", "GET", GetByID},
		{"/:id/items", "GET", GetItemsByListID},
		{"/genre/:genre", "GET", GetByGenre},
		{"/search", "GET", Search},

		// Resource-specific routes
		{"/id/:id", "GET", GetListByID},
		{"/genre", "GET", GetListsByGenre},
		{"/year", "GET", GetListsByYear},
		{"/actor", "GET", GetListsByActor},
		{"/creator", "GET", GetListsByCreator},
		{"/rating", "GET", GetListsByRating},
		{"/latest", "GET", GetLatestListsByAdded},
		{"/popular", "GET", GetPopularLists},
		{"/top-rated", "GET", GetTopRatedLists},
		{"/search-lists", "GET", SearchLists},
	}

	// Register all routes
	for _, route := range routeMappings {
		switch route.method {
		case "GET":
			resourceGroup.GET(route.path, handleClientResourceAction(mediaHandler, route.handler, resourceType))
		case "POST":
			resourceGroup.POST(route.path, handleClientResourceAction(mediaHandler, route.handler, resourceType))
		case "PUT":
			resourceGroup.PUT(route.path, handleClientResourceAction(mediaHandler, route.handler, resourceType))
		case "DELETE":
			resourceGroup.DELETE(route.path, handleClientResourceAction(mediaHandler, route.handler, resourceType))
		}
	}
}

// Handler factory for client resource actions
func handleClientResourceAction(mediaHandler apphandlers.ClientMediaHandlers, methodName HandlerMethod, resourceType MediaResourceType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client type from context
		clientType, exists := c.Get("clientType")
		if !exists {
			responses.RespondBadRequest(c, nil, "Client type not found in context")
			return
		}

		ct, ok := clientType.(clienttypes.ClientType)
		if !ok {
			responses.RespondBadRequest(c, nil, "Invalid client type in context")
			return
		}

		// Get the appropriate handler based on client type and resource type
		var handler interface{}

		switch resourceType {
		case PlaylistResource:
			switch ct {
			case clienttypes.ClientTypeEmby:
				handler = mediaHandler.EmbyPlaylistHandler()
			case clienttypes.ClientTypeJellyfin:
				handler = mediaHandler.JellyfinPlaylistHandler()
			case clienttypes.ClientTypePlex:
				handler = mediaHandler.PlexPlaylistHandler()
			case clienttypes.ClientTypeSubsonic:
				handler = mediaHandler.SubsonicPlaylistHandler()
			default:
				responses.RespondBadRequest(c, nil, "Unsupported client type for playlists")
				return
			}
		case CollectionResource:
			switch ct {
			case clienttypes.ClientTypeEmby:
				handler = mediaHandler.EmbyCollectionHandler()
			case clienttypes.ClientTypeJellyfin:
				handler = mediaHandler.JellyfinCollectionHandler()
			case clienttypes.ClientTypePlex:
				handler = mediaHandler.PlexCollectionHandler()
			case clienttypes.ClientTypeSubsonic:
				handler = mediaHandler.SubsonicCollectionHandler()
			default:
				responses.RespondBadRequest(c, nil, "Unsupported client type for collections")
				return
			}
		default:
			responses.RespondBadRequest(c, nil, "Unsupported resource type")
			return
		}

		// Use reflection to call the method by name
		handlerValue := reflect.ValueOf(handler)
		method := handlerValue.MethodByName(string(methodName))

		if !method.IsValid() {
			responses.RespondBadRequest(c, nil, "Method not found: "+string(methodName))
			return
		}

		// Call the method with the context
		method.Call([]reflect.Value{reflect.ValueOf(c)})
	}
}

