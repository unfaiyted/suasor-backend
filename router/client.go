// router/client.go
package router

import (
	"context"
	"fmt"
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
	log := logger.LoggerFromContext(ctx)

	log.Info().Msg("Registering client routes")

	// Media Clients
	registerClientRoutes[*types.EmbyConfig](ctx, r, c)
	registerClientRoutes[*types.JellyfinConfig](ctx, r, c)
	registerClientRoutes[*types.PlexConfig](ctx, r, c)
	registerClientRoutes[*types.SubsonicConfig](ctx, r, c)

	// Automation Clients
	registerClientRoutes[*types.SonarrConfig](ctx, r, c)
	registerClientRoutes[*types.RadarrConfig](ctx, r, c)
	registerClientRoutes[*types.LidarrConfig](ctx, r, c)

	// Ai Clients
	registerClientRoutes[*types.ClaudeConfig](ctx, r, c)
	registerClientRoutes[*types.OpenAIConfig](ctx, r, c)
	registerClientRoutes[*types.OllamaConfig](ctx, r, c)

	existingClientGroup := r.Group("/client/:clientID")
	existingClientGroup.Use(middleware.ClientTypeMiddleware(db))
	{
		existingClientGroup.GET("", func(g *gin.Context) {
			getClientHandler(g, c).GetClient(g)
		})
		existingClientGroup.PUT("", func(g *gin.Context) {
			getClientHandler(g, c).UpdateClient(g)
		})
		existingClientGroup.DELETE("", func(g *gin.Context) {
			getClientHandler(g, c).DeleteClient(g)
		})

		existingClientGroup.GET("/test", func(g *gin.Context) {
			getClientHandler(g, c).TestConnection(g)
		})
	}

	automationHandler := container.MustGet[*handlers.ClientAutomationHandler](c)
	automationClientGroup := r.Group("client/:clientID/automation")
	automationClientGroup.Use(middleware.ClientTypeMiddleware(db))

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
	aiGroup := r.Group("/client/:clientID/ai")
	aiGroup.Use(middleware.ClientTypeMiddleware(db))
	{
		aiGroup.GET("/recommendations", func(g *gin.Context) {
			getAIClientHandler(g, c).RequestRecommendation(g)
		})
		aiGroup.POST("/recommendations", func(g *gin.Context) {
			getAIClientHandler(g, c).AnalyzeContent(g)
		})
		aiGroup.POST("/conversation/start", func(g *gin.Context) {
			getAIClientHandler(g, c).StartConversation(g)
		})
		aiGroup.POST("/conversation/message", func(g *gin.Context) {
			getAIClientHandler(g, c).SendConversationMessage(g)
		})
	}
}

func getClientHandler(g *gin.Context, c *container.Container) handlers.ClientHandler[clienttypes.ClientConfig] {
	log := logger.LoggerFromContext(g.Request.Context())
	clientType := getClientType(g)

	log.Info().
		Str("clientType", string(clientType)).
		Msg("Retrieving client handler")
	handlers := map[clienttypes.ClientType]handlers.ClientHandler[clienttypes.ClientConfig]{
		clienttypes.ClientTypeEmby:     container.MustGet[handlers.ClientHandler[*clienttypes.EmbyConfig]](c),
		clienttypes.ClientTypeJellyfin: container.MustGet[handlers.ClientHandler[*clienttypes.JellyfinConfig]](c),
		clienttypes.ClientTypePlex:     container.MustGet[handlers.ClientHandler[*clienttypes.PlexConfig]](c),
		clienttypes.ClientTypeSubsonic: container.MustGet[handlers.ClientHandler[*clienttypes.SubsonicConfig]](c),
		clienttypes.ClientTypeSonarr:   container.MustGet[handlers.ClientHandler[*clienttypes.SonarrConfig]](c),
		clienttypes.ClientTypeRadarr:   container.MustGet[handlers.ClientHandler[*clienttypes.RadarrConfig]](c),
		clienttypes.ClientTypeLidarr:   container.MustGet[handlers.ClientHandler[*clienttypes.LidarrConfig]](c),
		clienttypes.ClientTypeClaude:   container.MustGet[handlers.ClientHandler[*clienttypes.ClaudeConfig]](c),
		clienttypes.ClientTypeOpenAI:   container.MustGet[handlers.ClientHandler[*clienttypes.OpenAIConfig]](c),
		clienttypes.ClientTypeOllama:   container.MustGet[handlers.ClientHandler[*clienttypes.OllamaConfig]](c),
	}
	handler, exists := handlers[clientType]
	if !exists {
		responses.RespondInternalError(g, nil, "Client handler not found")
		return nil
	}
	return handler
}

func getAIClientHandler(g *gin.Context, c *container.Container) handlers.AIHandler[clienttypes.AIClientConfig] {
	log := logger.LoggerFromContext(g.Request.Context())
	clientType := getClientType(g)

	log.Info().
		Str("clientType", string(clientType)).
		Msg("Retrieving AI client handler")
	handlers := map[clienttypes.ClientType]handlers.AIHandler[clienttypes.AIClientConfig]{
		clienttypes.ClientTypeClaude: container.MustGet[handlers.AIHandler[*clienttypes.ClaudeConfig]](c),
		clienttypes.ClientTypeOpenAI: container.MustGet[handlers.AIHandler[*clienttypes.OpenAIConfig]](c),
		clienttypes.ClientTypeOllama: container.MustGet[handlers.AIHandler[*clienttypes.OllamaConfig]](c),
	}
	handler, exists := handlers[clientType]
	if !exists {
		responses.RespondInternalError(g, nil, "AI client handler not found")
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
			log.Info().
				Str("clientType", string(clientType)).
				Msg("Retrieving clients")
			g.Set("clientType", clientType)
			getClientHandler(g, c).GetAllOfType(g)
		})
		clientGroup.POST("", func(g *gin.Context) {
			g.Set("clientType", clientType)
			getClientHandler(g, c).CreateClient(g)
		})
	}
}

func getClientType(g *gin.Context) clienttypes.ClientType {
	log := logger.LoggerFromContext(g.Request.Context())
	clientTypeVal, exists := g.Get("clientType")
	if !exists {
		responses.RespondBadRequest(g, nil, "Client type not found in context")
		return ""
	}

	// Cast to ClientType
	clientType, ok := clientTypeVal.(clienttypes.ClientType)
	if !ok {
		log.Error().Str("actual_type", fmt.Sprintf("%T", clientTypeVal)).Msg("Invalid client type in context")
		responses.RespondInternalError(g, nil, "Invalid client type in context")
		return ""
	}

	log.Debug().Str("clientType", string(clientType)).Msg("Got client type")
	return clientType
}
