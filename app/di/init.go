// app/di/init.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/services"
	"suasor/utils"
)

// Initialize registers all dependencies in the container
func Initialize(ctx context.Context, db *gorm.DB, configService services.ConfigService) *container.Container {
	// Create a new container
	c := container.NewContainer()
	log := utils.LoggerFromContext(ctx)

	// Register core dependencies
	log.Info().Msg("Registering core dependencies")
	RegisterCore(ctx, c, db, configService)

	// Register repositories
	log.Info().Msg("Registering repositories")
	RegisterRepositories(ctx, c)

	// Register media data factory and repositories
	log.Info().Msg("Registering media data")
	RegisterMediaData(ctx, c)

	// Register services
	log.Info().Msg("Registering services")
	RegisterServices(ctx, c)

	// Register handlers
	log.Info().Msg("Registering handlers")
	RegisterHandlers(ctx, c)

	return c
}
