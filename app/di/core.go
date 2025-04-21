// app/di/core.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/app/di/factories"
	"suasor/repository"
	"suasor/services"
	"suasor/utils"
)

// Register core dependencies that are used throughout the application
func RegisterCore(ctx context.Context, c *container.Container, db *gorm.DB, configService services.ConfigService) {
	// Register core components
	log := utils.LoggerFromContext(ctx)
	c.Register(db)
	c.Register(configService)

	// Register config service as a factory for reuse
	log.Info().Msg("Registering config service as a factory")
	container.RegisterFactory[services.ConfigService](c,
		func(c *container.Container) services.ConfigService {
			return configService
		})

	// Register client factory service
	// This is responsible for creating clients based on the client type and client ID
	// This is a singleton service that is shared across the application
	// It ensures that our external clients are created only once and reused throughout the application

	// First register the ClientItemRegistry before it's needed by client factories
	log.Info().Msg("Registering client media item factories")
	factories.RegisterClientMediaItemFactories(ctx, c)
	// Then register client factories which depend on the registry
	log.Info().Msg("Registering client factories")
	factories.RegisterClientFactories(ctx, c)

	// Register config repository
	log.Info().Msg("Registering config repository")
	container.RegisterFactory[repository.ConfigRepository](c, func(c *container.Container) repository.ConfigRepository {
		return container.MustGet[services.ConfigService](c).GetRepo()
	})
}
