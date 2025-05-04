// router/episode.go
package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/router/middleware"
	"suasor/types/responses"
	"suasor/utils/logger"
)

// registerClientEpisodeRoutes sets up the routes for episode-related operations
func registerClientEpisodeRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Registering episode routes")

	db := container.MustGet[*gorm.DB](c)

	// Client episode routes
	clientGroup := rg.Group("")
	clientGroup.Use(middleware.ClientTypeMiddleware(db))
	{
		// // Individual episode access
		// clientGroup.GET("/:clientItemID", func(g *gin.Context) {
		// 	if handler := getEpisodeHandler(g, c); handler != nil {
		// 		handler.GetEpisodeByID(g)
		// 	}
		// })

		// // Episode playback
		// clientGroup.POST("/:clientItemID/play", func(g *gin.Context) {
		// 	if handler := getEpisodeHandler(g, c); handler != nil {
		// 		handler.RecordEpisodePlay(g)
		// 	}
		// })

		// Episode state
		// clientGroup.GET("/:clientItemID/state", func(g *gin.Context) {
		// 	if handler := getEpisodeHandler(g, c); handler != nil {
		// 		handler.GetEpisodePlaybackState(g)
		// 	}
		// })

		// clientGroup.PUT("/:clientItemID/state", func(g *gin.Context) {
		// 	if handler := getEpisodeHandler(g, c); handler != nil {
		// 		handler.UpdateEpisodePlaybackState(g)
		// 	}
		// })
	}
}

func getEpisodeHandlerMap(c *container.Container, clientType clienttypes.ClientType) (handlers.ClientEpisodeHandler[clienttypes.ClientMediaConfig], bool) {
	handlerMap := map[clienttypes.ClientType]handlers.ClientEpisodeHandler[clienttypes.ClientMediaConfig]{
		clienttypes.ClientTypeEmby:     container.MustGet[handlers.ClientEpisodeHandler[*clienttypes.EmbyConfig]](c),
		clienttypes.ClientTypeJellyfin: container.MustGet[handlers.ClientEpisodeHandler[*clienttypes.JellyfinConfig]](c),
		clienttypes.ClientTypePlex:     container.MustGet[handlers.ClientEpisodeHandler[*clienttypes.PlexConfig]](c),
	}

	handler, exists := handlerMap[clientType]
	return handler, exists
}

func getEpisodeHandler(g *gin.Context, c *container.Container) handlers.ClientEpisodeHandler[clienttypes.ClientMediaConfig] {
	log := logger.LoggerFromContext(g.Request.Context())

	clientTypeVal, exists := g.Get("clientType")
	if !exists {
		log.Warn().Msg("Client type not found in request context")
		responses.RespondBadRequest(g, nil, "Client type not found")
		return nil
	}

	clientType := clientTypeVal.(clienttypes.ClientType)
	log.Debug().Str("clientType", string(clientType)).Msg("Getting client episode handler")

	handler, exists := getEpisodeHandlerMap(c, clientType)
	if !exists {
		err := fmt.Errorf("unsupported client type: %s", clientType)
		responses.RespondBadRequest(g, err, "Unsupported client type")
		return nil
	}

	return handler
}

