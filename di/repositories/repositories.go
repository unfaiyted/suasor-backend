// app/di/repositories.go
package repositories

import (
	"context"
	"suasor/di/container"
	"suasor/utils/logger"
)

// RegisterRepositories registers all repository dependencies
func RegisterRepositories(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)

	// System repositories
	log.Info().Msg("Registering system repositories")
	registerSystemRepositories(c)

	// Register client repositories
	log.Info().Msg("Registering client repositories")
	registerClientRepositories(ctx, c)

	// Register three-pronged architecture repositories
	log.Info().Msg("Registering media repositories")
	RegisterMediaRepositories(ctx, c)

	// Job repository
	log.Info().Msg("Registering job repository")
	registerJobRepositories(ctx, c)
}
