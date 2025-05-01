package services

import (
	"context"
	"suasor/clients"
	"suasor/di/container"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	svcbundles "suasor/services/bundles"
	"suasor/services/jobs"
	"suasor/services/jobs/recommendation"
	"suasor/services/jobs/sync"
	"suasor/utils/logger"
)

func registerJobServices(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Registering job implementations")

	// Define empty jobs for different job types
	recommendationJobImpl := &jobs.EmptyJob{JobName: "system.recommendation"}
	mediaSyncJobImpl := &jobs.EmptyJob{JobName: "system.media.sync"}
	watchHistorySyncJobImpl := &jobs.EmptyJob{JobName: "system.watch.history.sync"}
	favoritesSyncJobImpl := &jobs.EmptyJob{JobName: "system.favorites.sync"}

	// Register the empty jobs directly in the container
	c.Register(recommendationJobImpl)
	c.Register(mediaSyncJobImpl)
	c.Register(watchHistorySyncJobImpl)
	c.Register(favoritesSyncJobImpl)

	// Job service
	container.RegisterFactory[jobs.JobService](c, func(c *container.Container) jobs.JobService {
		jobRepo := container.MustGet[repository.JobRepository](c)
		userRepo := container.MustGet[repository.UserRepository](c)
		configRepo := container.MustGet[repository.UserConfigRepository](c)

		// Media repositories needed for job service
		coreRepos := container.MustGet[repobundles.CoreMediaItemRepositories](c)
		movieRepo := coreRepos.MovieRepo()
		seriesRepo := coreRepos.SeriesRepo()
		musicRepo := coreRepos.TrackRepo()

		// User data repositories needed for job service
		userDataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)
		userMovieDataRepo := userDataRepos.MovieDataRepo()
		userSeriesDataRepo := userDataRepos.SeriesDataRepo()
		userMusicDataRepo := userDataRepos.TrackDataRepo()

		// Get job implementations
		watchHistorySyncJob := container.MustGet[*sync.WatchHistorySyncJob](c)
		favoritesSyncJob := container.MustGet[*sync.FavoritesSyncJob](c)
		mediaSyncJob := container.MustGet[*sync.MediaSyncJob](c)
		recommendationJob := container.MustGet[*recommendation.RecommendationJob](c)

		// Job implementations
		return jobs.NewJobService(
			jobRepo,
			userRepo,
			configRepo,
			movieRepo,
			seriesRepo,
			musicRepo,
			userMovieDataRepo,
			userSeriesDataRepo,
			userMusicDataRepo,
			recommendationJob,
			mediaSyncJob,
			watchHistorySyncJob,
			favoritesSyncJob,
		)
	})

	log.Info().Msg("Registering watch history sync job service")
	container.RegisterFactory[*sync.WatchHistorySyncJob](c, func(c *container.Container) *sync.WatchHistorySyncJob {
		jobRepo := container.MustGet[repository.JobRepository](c)
		userRepo := container.MustGet[repository.UserRepository](c)
		userConfigRepo := container.MustGet[repository.UserConfigRepository](c)
		clientRepos := container.MustGet[repobundles.ClientRepositories](c)
		dataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)
		clientItemRepos := container.MustGet[repobundles.ClientMediaItemRepositories](c)
		itemRepos := container.MustGet[repobundles.CoreMediaItemRepositories](c)
		clientFactories := container.MustGet[*clients.ClientProviderFactoryService](c)
		return sync.NewWatchHistorySyncJob(jobRepo, userRepo, userConfigRepo, clientRepos, dataRepos, clientItemRepos, itemRepos, clientFactories)
	})
	// Favorites Sync Job
	log.Info().Msg("Registering favorites sync job service")
	container.RegisterFactory[*sync.FavoritesSyncJob](c, func(c *container.Container) *sync.FavoritesSyncJob {
		jobRepo := container.MustGet[repository.JobRepository](c)
		userRepo := container.MustGet[repository.UserRepository](c)
		userConfigRepo := container.MustGet[repository.UserConfigRepository](c)
		clientRepos := container.MustGet[repobundles.ClientRepositories](c)
		dataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)
		clientItemRepos := container.MustGet[repobundles.ClientMediaItemRepositories](c)
		itemRepos := container.MustGet[repobundles.UserMediaItemRepositories](c)
		clientFactories := container.MustGet[*clients.ClientProviderFactoryService](c)
		return sync.NewFavoritesSyncJob(jobRepo, userRepo, userConfigRepo, clientRepos, dataRepos, clientItemRepos, itemRepos, clientFactories)

	})

	// Media Sync Job
	log.Info().Msg("Registering media sync job service")
	container.RegisterFactory[*sync.MediaSyncJob](c, func(c *container.Container) *sync.MediaSyncJob {

		jobRepo := container.MustGet[repository.JobRepository](c)
		userRepo := container.MustGet[repository.UserRepository](c)
		userConfigRepo := container.MustGet[repository.UserConfigRepository](c)
		clientRepos := container.MustGet[repobundles.ClientRepositories](c)
		dataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)
		clientItemRepos := container.MustGet[repobundles.ClientMediaItemRepositories](c)
		clientMusicServices := container.MustGet[svcbundles.ClientMusicServices](c)
		itemRepos := container.MustGet[repobundles.UserMediaItemRepositories](c)
		clientFactories := container.MustGet[*clients.ClientProviderFactoryService](c)
		return sync.NewMediaSyncJob(jobRepo, userRepo, userConfigRepo, clientRepos, dataRepos, clientItemRepos, clientMusicServices, itemRepos, clientFactories)
	})

	// Recommendation Job
	log.Info().Msg("Registering recommendation job service")
	container.RegisterFactory[*recommendation.RecommendationJob](c, func(c *container.Container) *recommendation.RecommendationJob {
		jobRepo := container.MustGet[repository.JobRepository](c)
		userRepo := container.MustGet[repository.UserRepository](c)
		userConfigRepo := container.MustGet[repository.UserConfigRepository](c)
		recommendationRepo := container.MustGet[repository.RecommendationRepository](c)
		clientRepos := container.MustGet[repobundles.ClientRepositories](c)
		itemRepos := container.MustGet[repobundles.CoreMediaItemRepositories](c)
		clientItemRepos := container.MustGet[repobundles.ClientMediaItemRepositories](c)
		dataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)
		clientFactories := container.MustGet[*clients.ClientProviderFactoryService](c)
		creditRepo := container.MustGet[repository.CreditRepository](c)
		peopleRepo := container.MustGet[repository.PersonRepository](c)
		return recommendation.NewRecommendationJob(ctx, jobRepo, userRepo, userConfigRepo, recommendationRepo, clientRepos, itemRepos, clientItemRepos, dataRepos, clientFactories, creditRepo, peopleRepo)

	})

}
