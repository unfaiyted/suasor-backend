// app/di/handlers.go
package handlers

import (
	"context"
	"suasor/di/container"
)

// RegisterHandlers registers all handler dependencies
func RegisterHandlers(ctx context.Context, c *container.Container) {
	// Register core handlers (system, user, client)
	RegisterSystemHandlers(ctx, c)
	RegisterUserHandlers(ctx, c)
	RegisterClientHandlers(ctx, c)

	// Register media handlers
	RegisterMediaItemHandlers(ctx, c)

	// Register meida data handlers
	RegisterMediaDataHandlers(ctx, c)

	// Register job handlers
	RegisterJobHandlers(ctx, c)

	// Register recommendation handlers
	RegisterRecommendationHandlers(ctx, c)

	// Register media list handlers
	RegisterMediaListHandlers(ctx, c)
}
