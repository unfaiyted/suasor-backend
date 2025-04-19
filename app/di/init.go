// app/di/init.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/services"
)

// Initialize registers all dependencies in the container
func Initialize(ctx context.Context, db *gorm.DB, configService services.ConfigService) *container.Container {
	// Create a new container
	c := container.NewContainer()

	// Register core dependencies
	RegisterCore(ctx, c, db, configService)

	// Register repositories
	RegisterRepositories(ctx, c)

	// Register media data factory and repositories
	RegisterMediaData(ctx, c)

	// Register services
	RegisterServices(ctx, c)

	// Register handlers
	RegisterHandlers(ctx, c)

	return c
}
