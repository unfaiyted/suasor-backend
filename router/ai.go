package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
)

// RegisterAIRoutes registers routes for AI operations
func RegisterAIRoutes(r *gin.RouterGroup, deps *app.AppDependencies) {
	handler := deps.AIHandlers.ClaudeAIHandler()

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

