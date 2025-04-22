// app/di/services.go
package services

import (
	"context"
	"suasor/di/container"
	"suasor/utils/logger"
)

// RegisterServices registers all service dependencies
func RegisterServices(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	// Register system services
	log.Info().Msg("Registering system services")
	registerSystemServices(ctx, c)

	// Register client services
	log.Info().Msg("Registering client services")
	registerClientServices(ctx, c)

	// Register media item services
	log.Info().Msg("Registering media item services")
	registerMediaItemServices(ctx, c)

	// Register media data services
	log.Info().Msg("Registering media data services")
	registerMediaDataServices(ctx, c)

	// Register list services
	log.Info().Msg("Registering list services")
	registerListServices(ctx, c)

	// Register jobs
	log.Info().Msg("Registering jobs")
	registerJobServices(ctx, c)

	// Search service
	log.Info().Msg("Registering search service")
	registerSearchService(ctx, c)

	// Recommendation service
	log.Info().Msg("Registering recommendation service")
	registerRecommendationService(ctx, c)

}
