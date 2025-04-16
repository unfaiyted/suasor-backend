package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	"suasor/handlers"
)

// RegisterAIRoutes registers routes for AI operations
func RegisterAIRoutes(r *gin.RouterGroup, c *container.Container) {
	handlers := container.MustGet[apphandlers.AIClientHandlers](c)

	claude := handlers.ClaudeAIHandler()
	openai := handlers.OpenAIHandler()
	ollama := handlers.OllamaHandler()

	handlerMap := map[string]handlers.AIHandler[*handlers.ClientMediaItemDataHandlers]{
		"claude": claude,
		"openai": openai,
		"ollama": ollama,
	}

	getHandler := func(c *gin.Context) handlers.AIHandler[*handlers.ClientMediaItemDataHandlers] {
		clientType := c.Param("clientType")
		handler, exists := handlerMap[clientType]
		if !exists {
			// Default to Claude if type not specified or invalid
			return claude
		}
		return handler
	}

	ai := r.Group("/ai")
	client := ai.Group(":clientType")
	{
		// Recommendations and analysis endpoints
		client.POST("/recommendations", handler.RequestRecommendation)
		client.POST("/analyze", handler.AnalyzeContent)

		// Conversational recommendation endpoints
		conversation := client.Group("/conversation")
		{
			conversation.POST("/start", handler.StartConversation)
			conversation.POST("/message", handler.SendConversationMessage)
		}
	}
}
