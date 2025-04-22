// app/di/repositories.go
package repositories

import (
	"context"
	"gorm.io/gorm"
	"suasor/di/container"
	"suasor/repository"
	"suasor/services/jobs"
	"suasor/services/jobs/recommendation"
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

	// Watch History Sync Job
	log.Info().Msg("Registering watch history sync job repository")
	container.RegisterFactory[jobs.WatchHistorySyncJob](c, func(c *container.Container) jobs.WatchHistorySyncJob {
		return *jobs.NewWatchHistorySyncJob(ctx, c)
	})
	// Favorites Sync Job
	log.Info().Msg("Registering favorites sync job repository")
	container.RegisterFactory[jobs.FavoritesSyncJob](c, func(c *container.Container) jobs.FavoritesSyncJob {
		return *jobs.NewFavoritesSyncJob(ctx, c)
	})

	// Media Sync Job
	log.Info().Msg("Registering media sync job repository")
	container.RegisterFactory[jobs.MediaSyncJob](c, func(c *container.Container) jobs.MediaSyncJob {
		return *jobs.NewMediaSyncJob(ctx, c)
	})

	// Recommendation Job
	log.Info().Msg("Registering recommendation job repository")
	container.RegisterFactory[recommendation.RecommendationJob](c, func(c *container.Container) recommendation.RecommendationJob {
		return *recommendation.NewRecommendationJob(ctx, c)
	})

}
