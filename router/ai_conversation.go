package router

import (
	"suasor/di/container"
	"suasor/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterAIConversationRoutes registers routes for AI conversation history
func RegisterAIConversationRoutes(r *gin.RouterGroup, c *container.Container) {
	// Get the handler from the DI container
	handler := container.MustGet[handlers.AIConversationHandler](c)

	// Two sets of routes:
	// 1. Client-specific routes under /client/:clientID/ai/...
	// 2. User-centric routes under /user/ai/...

	// Client-specific routes for AI conversation
	clientRoute := r.Group("/client/:clientID/ai")
	{
		// History endpoints for specific client
		clientRoute.GET("/conversations", handler.GetUserConversations)
		clientRoute.GET("/conversations/:conversationId/messages", handler.GetConversationHistory)
		// clientRoute.GET("/recommendations", handler.GetUserRecommendations)
		clientRoute.POST("/conversations/:conversationId/continue", handler.ContinueConversation)
		clientRoute.PUT("/conversations/:conversationId/archive", handler.ArchiveConversation)
		clientRoute.DELETE("/conversations/:conversationId", handler.DeleteConversation)
	}

	// User-centric routes (across all AI clients)
	userGroup := r.Group("/user/ai")
	{
		// Same endpoints but for all client conversations for the user
		userGroup.GET("/conversations", handler.GetUserConversations)
		userGroup.GET("/conversations/:conversationId/messages", handler.GetConversationHistory)
		userGroup.GET("/recommendations", handler.GetUserRecommendations)
		userGroup.POST("/conversations/:conversationId/continue", handler.ContinueConversation)
		userGroup.PUT("/conversations/:conversationId/archive", handler.ArchiveConversation)
		userGroup.DELETE("/conversations/:conversationId", handler.DeleteConversation)
	}
}

