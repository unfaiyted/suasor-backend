// app/init_dependencies.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/di/container"
	"suasor/services"
	"suasor/utils/logger"
)

// AppDependencies contains all the application's dependencies
type AppDependencies struct {
	db        *gorm.DB
	container *container.Container
}

// GetContainer returns the application's dependency container
func (deps *AppDependencies) GetContainer() *container.Container {
	container := deps.container
	if container == nil {
		panic("Container not initialized")
	}
	return container
}

// InitializeDependencies initializes all application dependencies
func InitializeDependencies(ctx context.Context, db *gorm.DB, configService services.ConfigService) *AppDependencies {
	log := logger.LoggerFromContext(ctx)

	// Create and initialize the container using the new DI structure
	log.Info().Msg("Initializing dependency container")
	c := RegisterAppContainers(ctx, db, configService)

	// Create the application dependencies with the initialized container and handlers
	deps := &AppDependencies{
		container: c,
		db:        db,
	}

	log.Info().Msg("Application dependencies initialized successfully")
	return deps
}
