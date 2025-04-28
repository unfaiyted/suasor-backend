// router/series.go
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

// RegisterSeriesRoutes sets up the routes for series-related operations
func registerClientSeriesRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Registering series routes")

	db := container.MustGet[*gorm.DB](c)

	// Client series routes
	clientGroup := rg.Group("")
	clientGroup.Use(middleware.ClientTypeMiddleware(db))
	{

		// Get seasons by series ID
		clientGroup.GET("/:clientItemID/seasons", func(g *gin.Context) {
			if handler := getSeriesHandler(g, c); handler != nil {
				handler.GetSeasonsBySeriesID(g)
			}
		})

		// Get episodes by series ID
		clientGroup.GET("/:clientItemID/episodes", func(g *gin.Context) {
			if handler := getSeriesHandler(g, c); handler != nil {
				handler.GetEpisodesBySeriesID(g)
			}
		})

		// Discovery endpoints
		clientGroup.GET("/popular/:count", func(g *gin.Context) {
			if handler := getSeriesHandler(g, c); handler != nil {
				handler.GetPopularSeries(g)
			}
		})

		clientGroup.GET("/top-rated/:count", func(g *gin.Context) {
			if handler := getSeriesHandler(g, c); handler != nil {
				handler.GetTopRatedSeries(g)
			}
		})

		clientGroup.GET("/latest/:count", func(g *gin.Context) {
			if handler := getSeriesHandler(g, c); handler != nil {
				handler.GetLatestSeriesByAdded(g)
			}
		})

		clientGroup.GET("/actor/:actor", func(g *gin.Context) {
			if handler := getSeriesHandler(g, c); handler != nil {
				handler.GetSeriesByActor(g)
			}
		})

		clientGroup.GET("/creator/:creator", func(g *gin.Context) {
			if handler := getSeriesHandler(g, c); handler != nil {
				handler.GetSeriesByCreator(g)
			}
		})

	}
}

func getSeriesHandlerMap(c *container.Container, clientType clienttypes.ClientType) (handlers.ClientSeriesHandler[clienttypes.ClientMediaConfig], bool) {
	handlerMap := map[clienttypes.ClientType]handlers.ClientSeriesHandler[clienttypes.ClientMediaConfig]{
		clienttypes.ClientTypeEmby:     container.MustGet[handlers.ClientSeriesHandler[*clienttypes.EmbyConfig]](c),
		clienttypes.ClientTypeJellyfin: container.MustGet[handlers.ClientSeriesHandler[*clienttypes.JellyfinConfig]](c),
		clienttypes.ClientTypePlex:     container.MustGet[handlers.ClientSeriesHandler[*clienttypes.PlexConfig]](c),
	}

	handler, exists := handlerMap[clientType]
	return handler, exists
}

func getSeriesHandler(g *gin.Context, c *container.Container) handlers.ClientSeriesHandler[clienttypes.ClientMediaConfig] {
	log := logger.LoggerFromContext(g.Request.Context())

	clientTypeVal, exists := g.Get("clientType")
	if !exists {
		log.Warn().Msg("Client type not found in request context")
		responses.RespondBadRequest(g, nil, "Client type not found")
		return nil
	}

	clientType := clientTypeVal.(clienttypes.ClientType)
	log.Debug().Str("clientType", string(clientType)).Msg("Getting client series handler")

	handler, exists := getSeriesHandlerMap(c, clientType)
	if !exists {
		err := fmt.Errorf("unsupported client type: %s", clientType)
		responses.RespondBadRequest(g, err, "Unsupported client type")
		return nil
	}

	return handler
}

