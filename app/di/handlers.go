// app/di/handlers.go
package di

import (
	"suasor/app/container"
	"suasor/app/di/handlers"
)

// RegisterHandlers registers all handler dependencies
func RegisterHandlers(c *container.Container) {
	// Register core handlers (system, user, client)
	handlers.RegisterSystemHandlers(c)
	handlers.RegisterUserHandlers(c)
	handlers.RegisterClientHandlers(c)

	// Register media handlers
	handlers.RegisterMediaHandlers(c)
}

