package router

import (
	"context"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/router/middleware"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAIClientRoutes(ctx context.Context, r *gin.RouterGroup, c *container.Container) {
	db := container.MustGet[*gorm.DB](c)
	log := logger.LoggerFromContext(ctx)

	log.Info().Msg("Registering AI client routes")

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
