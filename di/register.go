// app/di/init.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/di/container"
	dihandlers "suasor/di/handlers"
	direpos "suasor/di/repositories"
	diservices "suasor/di/services"

	"suasor/services"
	"suasor/utils/logger"
)

// Initialize registers all dependencies in the container
func RegisterAppContainers(ctx context.Context, db *gorm.DB, configService services.ConfigService) *container.Container {
	// Create a new container
	c := container.NewContainer()
	log := logger.LoggerFromContext(ctx)

	// Register core dependencies
	log.Info().Msg("Registering core dependencies")
	RegisterCore(ctx, c, db, configService)

	// Register repositories
	log.Info().Msg("Registering repositories")
	direpos.RegisterRepositories(ctx, c)

	// Register services
	log.Info().Msg("Registering services")
	diservices.RegisterServices(ctx, c)

	// Register handlers
	log.Info().Msg("Registering handlers")
	dihandlers.RegisterHandlers(ctx, c)

	// Register Bundles
	log.Info().Msg("Registering bundles")
	// RegisterBundles(ctx, c)

	return c
}
