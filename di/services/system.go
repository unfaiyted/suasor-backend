package services

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"suasor/di/container"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	"suasor/services"
	"suasor/services/jobs"
	"suasor/services/jobs/recommendation"
	"suasor/utils/logger"
	"time"
)

// Register system-level services
func registerSystemServices(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	// Health service
	log.Info().Msg("Registering health service")
	container.RegisterFactory[services.HealthService](c, func(c *container.Container) services.HealthService {
		db := container.MustGet[*gorm.DB](c)
		return services.NewHealthService(db)
	})

	// User service
	log.Info().Msg("Registering user service")
	container.RegisterFactory[services.UserService](c, func(c *container.Container) services.UserService {
		userRepo := container.MustGet[repository.UserRepository](c)
		return services.NewUserService(userRepo)
	})

	// Search service
	log.Info().Msg("Registering search service")
	registerSearchService(ctx, c)

	registerJobServices(ctx, c)

	// Auth service
	log.Info().Msg("Registering auth service")
	container.RegisterSingleton[services.AuthService](c, func(c *container.Container) services.AuthService {
		fmt.Println("Creating AuthService")
		fmt.Println("Getting UserRepository for AuthService")
		userRepo := container.MustGet[repository.UserRepository](c)
		fmt.Println("Got UserRepository for AuthService")

		fmt.Println("Getting SessionRepository for AuthService")
		sessionRepo := container.MustGet[repository.SessionRepository](c)
		fmt.Println("Got SessionRepository for AuthService")

		fmt.Println("Getting ConfigService for AuthService")
		configService := container.MustGet[services.ConfigService](c)
		fmt.Println("Got ConfigService for AuthService")

		// Get auth config from config service
		fmt.Println("Getting config for AuthService")
		appConfig := configService.GetConfig()
		fmt.Println("Got config for AuthService")

		// Verify auth config values
		fmt.Printf("AuthService config: JWTSecret=%s, TokenExpiration=%d, RefreshExpiryDays=%d, TokenIssuer=%s, TokenAudience=%s\n",
			appConfig.Auth.JWTSecret,
			appConfig.Auth.TokenExpiration,
			appConfig.Auth.RefreshExpiryDays,
			appConfig.Auth.TokenIssuer,
			appConfig.Auth.TokenAudience)

		// Set up auth service with config values
		fmt.Println("Creating new AuthService instance")
		authService := services.NewAuthService(
			userRepo,
			sessionRepo,
			appConfig.Auth.JWTSecret,
			time.Duration(appConfig.Auth.TokenExpiration)*time.Hour,
			time.Duration(appConfig.Auth.RefreshExpiryDays)*24*time.Hour,
			appConfig.Auth.TokenIssuer,
			appConfig.Auth.TokenAudience,
		)
		fmt.Println("AuthService created successfully")
		return authService
	})

	// Register job implementations
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
	log.Info().Msg("Registering job service")
	container.RegisterFactory[services.JobService](c, func(c *container.Container) services.JobService {
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
		watchHistorySyncJob := container.MustGet[*jobs.WatchHistorySyncJob](c)
		favoritesSyncJob := container.MustGet[*jobs.FavoritesSyncJob](c)
		mediaSyncJob := container.MustGet[*jobs.MediaSyncJob](c)
		recommendationJob := container.MustGet[*recommendation.RecommendationJob](c)

		// Job implementations
		return services.NewJobService(
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

	// UserConfig service
	log.Info().Msg("Registering user config service")
	container.RegisterFactory[services.UserConfigService](c, func(c *container.Container) services.UserConfigService {
		userConfigRepo := container.MustGet[repository.UserConfigRepository](c)
		jobService := container.MustGet[services.JobService](c)
		recommendationJob := container.MustGet[*recommendation.RecommendationJob](c)
		return services.NewUserConfigService(userConfigRepo, jobService, recommendationJob)
	})

	// Recommendation service
	container.RegisterFactory[services.RecommendationService](c, func(c *container.Container) services.RecommendationService {
		recommendationRepo := container.MustGet[repository.RecommendationRepository](c)
		return services.NewRecommendationService(recommendationRepo)
	})

}
