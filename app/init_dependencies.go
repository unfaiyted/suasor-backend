// app/init_dependencies.go
package app

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/app/di"
	"suasor/services"
	"suasor/utils"
)

// AppDependencies contains all the application's dependencies
type AppDependencies struct {
	db        *gorm.DB
	container *container.Container
	handlers  *di.ApplicationHandlers
}

// GetContainer returns the application's dependency container
func (deps *AppDependencies) GetContainer() *container.Container {
	container := deps.container
	if container == nil {
		panic("Container not initialized")
	}
	return container
}

// GetHandlers returns the application's handlers
func (deps *AppDependencies) GetHandlers() *di.ApplicationHandlers {
	handlers := deps.handlers
	if handlers == nil {
		panic("Handlers not initialized")
	}
	return handlers
}

// InitializeDependencies initializes all application dependencies
func InitializeDependencies(ctx context.Context, db *gorm.DB, configService services.ConfigService) *AppDependencies {
	log := utils.LoggerFromContext(ctx)
	
	// Create and initialize the container using the new DI structure
	log.Info().Msg("Initializing dependency container")
	c := di.Initialize(ctx, db, configService)

	// Get all handlers from the container
	log.Info().Msg("Getting organized handlers from container")
	handlers, err := di.GetAllHandlers(ctx, c)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get handlers from container")
	}

	// Create the application dependencies with the initialized container and handlers
	deps := &AppDependencies{
		container: c,
		db:        db,
		handlers:  handlers,
	}

	log.Info().Msg("Application dependencies initialized successfully")
	return deps
}