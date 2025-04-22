// app/di/handlers/users.go
package handlers

import (
	"context"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/services"
)

// RegisterUserHandlers registers the user-related handlers
func RegisterUserHandlers(ctx context.Context, c *container.Container) {
	// User handler
	container.RegisterFactory[*handlers.UserHandler](c, func(c *container.Container) *handlers.UserHandler {
		userService := container.MustGet[services.UserService](c)
		configService := container.MustGet[services.ConfigService](c)
		return handlers.NewUserHandler(userService, configService)
	})

	// Auth handler
	container.RegisterFactory[*handlers.AuthHandler](c, func(c *container.Container) *handlers.AuthHandler {
		authService := container.MustGet[services.AuthService](c)
		return handlers.NewAuthHandler(authService)
	})

	// User config handler
	container.RegisterFactory[*handlers.UserConfigHandler](c, func(c *container.Container) *handlers.UserConfigHandler {
		userConfigService := container.MustGet[services.UserConfigService](c)
		return handlers.NewUserConfigHandler(userConfigService)
	})
}
