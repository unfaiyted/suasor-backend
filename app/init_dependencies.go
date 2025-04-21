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

type AppDependencies struct {
	db        *gorm.DB
	container *container.Container
}

func (deps *AppDependencies) GetContainer() *container.Container {
	container := deps.container
	if container == nil {
		panic("Container not initialized")
	}
	return container
}

func InitializeDependencies(ctx context.Context, db *gorm.DB, configService services.ConfigService) *AppDependencies {
	log := utils.LoggerFromContext(ctx)
	// Create and initialize the container using the new DI structure
	c := di.Initialize(ctx, db, configService)

	// Create the application dependencies with the initialized container
	deps := &AppDependencies{
		container: c,
		db:        db,
	}

	log.Info().Msg("Initializing core services")
	return deps
}
