// app/di/init.go
package di

import (
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/services"
)

// Initialize registers all dependencies in the container
func Initialize(db *gorm.DB, configService services.ConfigService) *container.Container {
	// Create a new container
	c := container.NewContainer()

	// Register core dependencies
	RegisterCore(c, db, configService)

	// Register repositories
	RegisterRepositories(c)

	// Register media data factory and repositories
	RegisterMediaData(c)

	// Register services
	RegisterServices(c)

	// Register handlers
	RegisterHandlers(c)

	return c
}

