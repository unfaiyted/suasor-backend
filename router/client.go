// router/client.go
package router

import (
	"context"
	"suasor/clients/types"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/router/middleware"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterClientRoutes configures routes for client endpoints
// These endpoints are specific to a single client instance.
func RegisterClientRoutes(ctx context.Context, r *gin.RouterGroup, c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	registerClientRoutes[*types.EmbyConfig](ctx, r, c)
	registerClientRoutes[*types.JellyfinConfig](ctx, r, c)
	registerClientRoutes[*types.PlexConfig](ctx, r, c)
	registerClientRoutes[*types.SubsonicConfig](ctx, r, c)
	registerClientRoutes[*types.SonarrConfig](ctx, r, c)
	registerClientRoutes[*types.RadarrConfig](ctx, r, c)
	registerClientRoutes[*types.LidarrConfig](ctx, r, c)
	registerClientRoutes[*types.ClaudeConfig](ctx, r, c)
	registerClientRoutes[*types.OpenAIConfig](ctx, r, c)
	registerClientRoutes[*types.OllamaConfig](ctx, r, c)

	existingClientGroup := r.Group("/client/:clientID")
	existingClientGroup.Use(middleware.ClientTypeMiddleware(db))
	{
		existingClientGroup.GET("", func(g *gin.Context) {
			clientType := g.Param("clientType")
			getClientHandler(g, c, clientType).GetClient(g)
		})
		existingClientGroup.PUT("", func(g *gin.Context) {
			clientType := g.Param("clientType")
			getClientHandler(g, c, clientType).UpdateClient(g)
		})
		existingClientGroup.DELETE("", func(g *gin.Context) {
			clientType := g.Param("clientType")
			getClientHandler(g, c, clientType).DeleteClient(g)
		})

		existingClientGroup.GET("/test", func(g *gin.Context) {
			clientType := g.Param("clientType")
			getClientHandler(g, c, clientType).TestConnection(g)
		})
	}

	automationHandler := container.MustGet[*handlers.ClientAutomationHandler](c)
	automationClientGroup := r.Group("client/:clientID/automation")
	{
		automationClientGroup.GET("/status", automationHandler.GetSystemStatus)
		automationClientGroup.POST("/library", automationHandler.GetLibraryItems)
		automationClientGroup.GET("/item/:itemID", automationHandler.GetMediaByID)
		automationClientGroup.POST("/item", automationHandler.AddMedia)
		automationClientGroup.PUT("/item/:itemID", automationHandler.UpdateMedia)
		automationClientGroup.DELETE("/item/:itemID", automationHandler.DeleteMedia)
		automationClientGroup.GET("/command", automationHandler.ExecuteCommand)
		automationClientGroup.GET("/calendar", automationHandler.GetCalendar)
		automationClientGroup.GET("/search", automationHandler.SearchMedia)

	}
}

func getClientHandler(g *gin.Context, c *container.Container, clientType string) handlers.ClientHandler[clienttypes.ClientConfig] {
	log := logger.LoggerFromContext(g.Request.Context())
	log.Info().
		Str("clientType", clientType).
		Msg("Retrieving client handler")
	handlers := map[string]handlers.ClientHandler[clienttypes.ClientConfig]{
		"emby":     container.MustGet[handlers.ClientHandler[*clienttypes.EmbyConfig]](c),
		"jellyfin": container.MustGet[handlers.ClientHandler[*clienttypes.JellyfinConfig]](c),
		"plex":     container.MustGet[handlers.ClientHandler[*clienttypes.PlexConfig]](c),
		"subsonic": container.MustGet[handlers.ClientHandler[*clienttypes.SubsonicConfig]](c),
		"sonarr":   container.MustGet[handlers.ClientHandler[*clienttypes.SonarrConfig]](c),
		"radarr":   container.MustGet[handlers.ClientHandler[*clienttypes.RadarrConfig]](c),
		"lidarr":   container.MustGet[handlers.ClientHandler[*clienttypes.LidarrConfig]](c),
		"claude":   container.MustGet[handlers.ClientHandler[*clienttypes.ClaudeConfig]](c),
		"openai":   container.MustGet[handlers.ClientHandler[*clienttypes.OpenAIConfig]](c),
		"ollama":   container.MustGet[handlers.ClientHandler[*clienttypes.OllamaConfig]](c),
	}
	handler, exists := handlers[clientType]
	if !exists {
		responses.RespondInternalError(g, nil, "Client handler not found")
		return nil
	}
	return handler
}

func registerClientRoutes[T types.ClientConfig](ctx context.Context, r *gin.RouterGroup, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	clientType := types.GetClientType[T]()

	log.Info().
		Str("clientType", string(clientType)).
		Msg("Registering client routes")

	clientGroup := r.Group("/client/" + string(clientType))
	{
		clientGroup.GET("", func(g *gin.Context) {
			clientType := string(clientType)
			log.Info().
				Str("clientType", string(clientType)).
				Msg("Retrieving clients")
			getClientHandler(g, c, clientType).GetAllOfType(g)
		})
		clientGroup.POST("", func(g *gin.Context) {
			clientType := string(clientType)
			getClientHandler(g, c, clientType).CreateClient(g)
		})
	}
}
