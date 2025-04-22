// app/di/handlers.go
package di

import (
	"context"
	"suasor/container"
	"suasor/di/handlers"
)

// RegisterHandlers registers all handler dependencies
func RegisterHandlers(ctx context.Context, c *container.Container) {
	// Register core handlers (system, user, client)
	handlers.RegisterSystemHandlers(ctx, c)
	handlers.RegisterUserHandlers(ctx, c)
	handlers.RegisterClientHandlers(ctx, c)

	// Register media handlers
	handlers.RegisterMediaHandlers(ctx, c)

	// Register job handlers
	handlers.RegisterJobHandlers(ctx, c)

	// Register recommendation handlers
	handlers.RegisterRecommendationHandlers(ctx, c)

	// Register media list handlers
	handlers.RegisterMediaListHandlers(ctx, c)
}
