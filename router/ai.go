package router

import (
	"suasor/app/container"
	// "suasor/app/handlers"
	apphandlers "suasor/app/handlers"
	// "suasor/handlers"

	"github.com/gin-gonic/gin"
)

type AIHandlerInterface interface {
	RequestRecommendation(c *gin.Context)
	AnalyzeContent(c *gin.Context)
	StartConversation(c *gin.Context)
	SendConversationMessage(c *gin.Context)
}

// RegisterAIRoutes registers routes for AI operations
func RegisterAIRoutes(r *gin.RouterGroup, c *container.Container) {
	handlers := container.MustGet[apphandlers.AIClientHandlers](c)

	claude := handlers.ClaudeAIHandler()
	openai := handlers.OpenAIHandler()
	ollama := handlers.OllamaHandler()

	handlerMap := map[string]AIHandlerInterface{
		"claude": claude,
		"openai": openai,
		"ollama": ollama,
	}

	getHandler := func(c *gin.Context) AIHandlerInterface {
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
		client.POST("/recommendations", func(c *gin.Context) {
			getHandler(c).RequestRecommendation(c)
		})
		client.POST("/analyze", func(c *gin.Context) {
			getHandler(c).AnalyzeContent(c)
		})

		// Conversational recommendation endpoints
		conversation := client.Group("/conversation")
		{
			conversation.POST("/start", func(c *gin.Context) {
				getHandler(c).StartConversation(c)
			})
			conversation.POST("/message", func(c *gin.Context) {
				getHandler(c).SendConversationMessage(c)
			})
		}
	}
}
