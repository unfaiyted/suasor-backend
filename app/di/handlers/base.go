// app/di/handlers/base.go
package handlers

import (
	"context"
	"suasor/app/container"
	"suasor/handlers"
	"suasor/services"
)

// RegisterSystemHandlers registers the system-level handlers
func RegisterSystemHandlers(ctx context.Context, c *container.Container) {
	// Health handler
	container.RegisterFactory[*handlers.HealthHandler](c, func(c *container.Container) *handlers.HealthHandler {
		healthService := container.MustGet[services.HealthService](c)
		return handlers.NewHealthHandler(healthService)
	})

	// Config handler
	container.RegisterFactory[*handlers.ConfigHandler](c, func(c *container.Container) *handlers.ConfigHandler {
		configService := container.MustGet[services.ConfigService](c)
		return handlers.NewConfigHandler(configService)
	})

	// Search handler
	container.RegisterFactory[*handlers.SearchHandler](c, func(c *container.Container) *handlers.SearchHandler {
		searchService := container.MustGet[services.SearchService](c)
		return handlers.NewSearchHandler(searchService)
	})
}

