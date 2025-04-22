// app/di/repositories.go
package repositories

import (
	"context"
	"gorm.io/gorm"
	"suasor/di/container"
	"suasor/repository"
	"suasor/utils/logger"
)

func registerJobRepositories(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Registering job repository")
	container.RegisterFactory[repository.JobRepository](c, func(c *container.Container) repository.JobRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewJobRepository(db)
	})

	// Recommendation Repo
	log.Info().Msg("Registering recommendation repository")
	container.RegisterFactory[repository.RecommendationRepository](c, func(c *container.Container) repository.RecommendationRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewRecommendationRepository(db)
	})

}
