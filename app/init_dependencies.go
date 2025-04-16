// app/init_dependencies.go
package app

import (
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/app/di"
	"suasor/services"
)

type AppDependencies struct {
	db        *gorm.DB
	container *container.Container
}

func InitializeDependencies(db *gorm.DB, configService services.ConfigService) *AppDependencies {
	// Create and initialize the container using the new DI structure
	c := di.Initialize(db, configService)

	// Create the application dependencies with the initialized container
	deps := &AppDependencies{
		container: c,
		db:        db,
	}

	return deps
}

