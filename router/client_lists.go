package router // (M) Comment

import ( // (M) Comment
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	gorm "gorm.io/gorm"
	"strconv"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	"suasor/app/services"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/responses"
)

func RegisterClientListRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {
	mediaHandler := container.MustGet[apphandlers.ClientMediaHandlers](c)
	clientServices := container.MustGet[services.ClientServices](c)

	db := container.MustGet[*gorm.DB](c)

	// Create generic handler function for playlists
	handlePlaylists := func(c *gin.Context) {
		clientIDStr := c.Param("clientID")

		clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
		if err != nil {
			responses.RespondBadRequest(c, err, "Invalid client ID")
			return
		}

		// Get client type from ID
		clientType, err := repository.GetClientTypeFromID(ctx, db, clientID)
		if err != nil {
			responses.RespondNotFound(c, err, "Client not found")
			return
		}

		// Use type switching to call the correct handler with proper types
		switch clientType {
		case "jellyfin":
			config, err := clientServices.JellyfinService().GetClientConfig(clientID)
			if err != nil {
				responses.RespondInternalError(c, err, "Failed to get client config")
				return
			}
			handler := mediaHandler.JellyfinPlaylistHandler()

			// playlists, err := handler.GetAll(c, config)
			// if err != nil {
			// 	responses.RespondWithError(c, err)
			// 	return
			// }
			// responses.RespondOK(c, playlists)

		case "emby":
			config, err := clientService.GetClientConfig[types.EmbyConfig](clientID)
			if err != nil {
				responses.RespondInternalError(c, err, "Failed to get client config")
				return
			}
			handler := mediaHandler.GetPlaylistHandler[types.EmbyConfig, mediatypes.Playlist]()
			playlists, err := handler.GetPlaylists(c, config)
			if err != nil {
				responses.RespondWithError(c, err)
				return
			}
			responses.RespondOK(c, playlists)

		case "plex":
			config, err := clientService.GetClientConfig[types.PlexConfig](clientID)
			if err != nil {
				responses.RespondInternalError(c, err, "Failed to get client config")
				return
			}
			handler := mediaHandler.GetPlaylistHandler[types.PlexConfig, mediatypes.Playlist]()
			playlists, err := handler.GetPlaylists(c, config)
			if err != nil {
				responses.RespondWithError(c, err)
				return
			}
			responses.RespondOK(c, playlists)

		case "subsonic":
			config, err := clientService.GetClientConfig[types.SubsonicConfig](clientID)
			if err != nil {
				responses.RespondInternalError(c, err, "Failed to get client config")
				return
			}
			handler := mediaHandler.GetPlaylistHandler[types.SubsonicConfig, mediatypes.Playlist]()
			playlists, err := handler.GetPlaylists(c, config)
			if err != nil {
				responses.RespondWithError(c, err)
				return
			}
			responses.RespondOK(c, playlists)

		default:
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
		}
	}

	// Similar function for collections
	handleCollections := func(c *gin.Context) {
		// Similar implementation as handlePlaylists but for collections
		// ...
	}

	rg.Group("/client/") {
		rg.GET("/:clientID/playlists", handlePlaylists)
		rg.GET("/media/:clientID/collections", handleCollections)
	}

	// Register routes
	rg.GET("/client/:clientID/playlists", handlePlaylists)
	rg.GET("/clients/media/:clientID/collections", handleCollections)
}

