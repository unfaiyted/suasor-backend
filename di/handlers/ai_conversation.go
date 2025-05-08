package handlers

import (
	"context"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/services"
)

// RegisterAIConversationHandlers registers AI conversation handlers
func RegisterAIConversationHandlers(ctx context.Context, c *container.Container) {
	// Register the AIConversationHandler
	container.RegisterFactory[handlers.AIConversationHandler](c, func(c *container.Container) handlers.AIConversationHandler {
		conversationService := container.MustGet[services.AIConversationService](c)
		return handlers.NewAIConversationHandler(conversationService)
	})
}