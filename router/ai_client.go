package router

import (
	"context"
	"net/http"
	"suasor/clients/types"
	"suasor/di/container"
	apphandlers "suasor/handlers/bundles"
	"suasor/router/middleware"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAIClientRoutes registers routes for AI client operations
func RegisterAIClientRoutes(ctx context.Context, r *gin.RouterGroup, c *container.Container) {
	// Get dependencies from container
	db := container.MustGet[*gorm.DB](c)
	handler := container.MustGet[apphandlers.AIClientHandlers](c)
	log := logger.LoggerFromContext(ctx)

	// Get specific handlers
	claudeHandler := handler.ClaudeAIHandler()
	openaiHandler := handler.OpenAIHandler()
	ollamaHandler := handler.OllamaHandler()

	// Map of client type to handler
	handlerMap := map[types.ClientType]interface{}{
		types.ClientTypeClaude: claudeHandler,
		types.ClientTypeOpenAI: openaiHandler,
		types.ClientTypeOllama: ollamaHandler,
	}

	// AI routes for specific clients
	clientRoute := r.Group("/client/:clientID/ai")
	clientRoute.Use(middleware.ClientTypeMiddleware(db))
	{
		// Core AI endpoints that work with specific clients
		clientRoute.POST("/recommendations", func(c *gin.Context) {
			clientType, ok := c.Get("clientType")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Client type not found"})
				return
			}

			// Get appropriate handler based on client type
			if h, exists := handlerMap[clientType.(types.ClientType)]; exists {
				if handler, ok := h.(interface{ RequestRecommendation(c *gin.Context) }); ok {
					handler.RequestRecommendation(c)
					return
				}
			}
			
			log.Error().
				Str("clientType", clientType.(types.ClientType).String()).
				Str("endpoint", "/recommendations").
				Msg("No compatible handler found for client type")
			c.JSON(http.StatusBadRequest, gin.H{"error": "This AI client does not support recommendations"})
		})

		clientRoute.POST("/analyze", func(c *gin.Context) {
			clientType, ok := c.Get("clientType")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Client type not found"})
				return
			}

			// Get appropriate handler based on client type
			if h, exists := handlerMap[clientType.(types.ClientType)]; exists {
				if handler, ok := h.(interface{ AnalyzeContent(c *gin.Context) }); ok {
					handler.AnalyzeContent(c)
					return
				}
			}
			
			log.Error().
				Str("clientType", clientType.(types.ClientType).String()).
				Str("endpoint", "/analyze").
				Msg("No compatible handler found for client type")
			c.JSON(http.StatusBadRequest, gin.H{"error": "This AI client does not support content analysis"})
		})

		// Conversation endpoints
		conversation := clientRoute.Group("/conversation")
		{
			conversation.POST("/start", func(c *gin.Context) {
				clientType, ok := c.Get("clientType")
				if !ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Client type not found"})
					return
				}

				// Get appropriate handler based on client type
				if h, exists := handlerMap[clientType.(types.ClientType)]; exists {
					if handler, ok := h.(interface{ StartConversation(c *gin.Context) }); ok {
						handler.StartConversation(c)
						return
					}
				}
				
				log.Error().
					Str("clientType", clientType.(types.ClientType).String()).
					Str("endpoint", "/conversation/start").
					Msg("No compatible handler found for client type")
				c.JSON(http.StatusBadRequest, gin.H{"error": "This AI client does not support conversations"})
			})

			conversation.POST("/message", func(c *gin.Context) {
				clientType, ok := c.Get("clientType")
				if !ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Client type not found"})
					return
				}

				// Get appropriate handler based on client type
				if h, exists := handlerMap[clientType.(types.ClientType)]; exists {
					if handler, ok := h.(interface{ SendConversationMessage(c *gin.Context) }); ok {
						handler.SendConversationMessage(c)
						return
					}
				}
				
				log.Error().
					Str("clientType", clientType.(types.ClientType).String()).
					Str("endpoint", "/conversation/message").
					Msg("No compatible handler found for client type")
				c.JSON(http.StatusBadRequest, gin.H{"error": "This AI client does not support conversations"})
			})
			
			// Add the history/get list of conversations endpoint
			conversation.GET("", func(c *gin.Context) {
				clientType, ok := c.Get("clientType")
				if !ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Client type not found"})
					return
				}

				// Get appropriate handler based on client type
				if h, exists := handlerMap[clientType.(types.ClientType)]; exists {
					if handler, ok := h.(interface{ GetConversations(c *gin.Context) }); ok {
						handler.GetConversations(c)
						return
					}
				}
				
				log.Error().
					Str("clientType", clientType.(types.ClientType).String()).
					Str("endpoint", "/conversation").
					Msg("No compatible handler found for client type")
				c.JSON(http.StatusBadRequest, gin.H{"error": "This AI client does not support conversation history"})
			})
			
			// Get specific conversation history
			conversation.GET("/:conversationID", func(c *gin.Context) {
				clientType, ok := c.Get("clientType")
				if !ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Client type not found"})
					return
				}

				// Get appropriate handler based on client type
				if h, exists := handlerMap[clientType.(types.ClientType)]; exists {
					if handler, ok := h.(interface{ GetConversationHistory(c *gin.Context) }); ok {
						handler.GetConversationHistory(c)
						return
					}
				}
				
				log.Error().
					Str("clientType", clientType.(types.ClientType).String()).
					Str("endpoint", "/conversation/:conversationID").
					Msg("No compatible handler found for client type")
				c.JSON(http.StatusBadRequest, gin.H{"error": "This AI client does not support conversation history"})
			})
		}
	}
}