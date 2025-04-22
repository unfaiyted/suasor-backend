// app/di/services.go
package services

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	mediatypes "suasor/clients/media/types"
	"suasor/clients/types"
	clienttypes "suasor/clients/types"
	"suasor/container"
	"suasor/di/factories"
	"suasor/repository"
	apprepository "suasor/repository/bundles"
	"suasor/services"
	appservices "suasor/services/bundles"
	"suasor/services/jobs"
	"suasor/services/jobs/recommendation"
	"suasor/utils/logger"
	"time"
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

	// Register three-pronged architecture services
	log.Info().Msg("Registering three-pronged architecture services")
	registerThreeProngedServices(ctx, c)
}

// Register system-level services
func registerSystemServices(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	// Health service
	log.Info().Msg("Registering health service")
	container.RegisterFactory[services.HealthService](c, func(c *container.Container) services.HealthService {
		db := container.MustGet[*gorm.DB](c)
		return services.NewHealthService(db)
	})

	// Search service
	log.Info().Msg("Registering search service")
	diservices.RegisterSearchService(ctx, c)

	// User service
	log.Info().Msg("Registering user service")
	container.RegisterFactory[services.UserService](c, func(c *container.Container) services.UserService {
		userRepo := container.MustGet[repository.UserRepository](c)
		return services.NewUserService(userRepo)
	})

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

	// Media Sync Job - using a simple empty implementation
	log.Info().Msg("Registering media sync job")
	container.RegisterFactory[*jobs.MediaSyncJob](c, func(c *container.Container) *jobs.MediaSyncJob {
		job := jobs.NewMediaSyncJob(ctx, c)
		return job
	})

	// Watch History Sync Job - using existing definition but with fallback
	log.Info().Msg("Registering watch history sync job")
	container.RegisterFactory[*jobs.WatchHistorySyncJob](c, func(c *container.Container) *jobs.WatchHistorySyncJob {
		job := jobs.NewWatchHistorySyncJob(ctx, c)

		return job
	})

	// Favorites Sync Job - using existing definition but with fallback
	log.Info().Msg("Registering favorites sync job")
	container.RegisterFactory[*jobs.FavoritesSyncJob](c, func(c *container.Container) *jobs.FavoritesSyncJob {
		job := jobs.NewFavoritesSyncJob(ctx, c)

		return job
	})

	// Recommendation Job - using existing definition but with fallback
	log.Info().Msg("Registering recommendation job")
	container.RegisterFactory[*recommendation.RecommendationJob](c, func(c *container.Container) *recommendation.RecommendationJob {
		job := recommendation.NewRecommendationJob(ctx, c)

		return job
	})

	// Job service
	log.Info().Msg("Registering job service")
	container.RegisterFactory[services.JobService](c, func(c *container.Container) services.JobService {
		jobRepo := container.MustGet[repository.JobRepository](c)
		userRepo := container.MustGet[repository.UserRepository](c)
		configRepo := container.MustGet[repository.UserConfigRepository](c)

		// Media repositories needed for job service
		coreRepos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		movieRepo := coreRepos.MovieRepo()
		seriesRepo := coreRepos.SeriesRepo()
		musicRepo := coreRepos.TrackRepo()

		// User data repositories needed for job service
		userDataRepos := container.MustGet[apprepository.UserMediaDataRepositories](c)
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

	// People and credits services
	container.RegisterFactory[*services.PersonService](c, func(c *container.Container) *services.PersonService {
		personRepo := container.MustGet[repository.PersonRepository](c)
		creditRepo := container.MustGet[repository.CreditRepository](c)
		return services.NewPersonService(personRepo, creditRepo)
	})

	container.RegisterFactory[*services.CreditService](c, func(c *container.Container) *services.CreditService {
		creditRepo := container.MustGet[repository.CreditRepository](c)
		personRepo := container.MustGet[repository.PersonRepository](c)
		return services.NewCreditService(creditRepo, personRepo)
	})

	// Recommendation service
	container.RegisterFactory[services.RecommendationService](c, func(c *container.Container) services.RecommendationService {
		recommendationRepo := container.MustGet[repository.RecommendationRepository](c)
		return services.NewRecommendationService(recommendationRepo)
	})

	// Media services
	container.RegisterFactory[appservices.PeopleServices](c, func(c *container.Container) appservices.PeopleServices[*types.Person] {
		personService := container.MustGet[services.PersonService](c)
		creditService := container.MustGet[services.CreditService](c)
		return appservices.NewPeopleServices(&personService, &creditService)
	})
}

// Register client-specific services
func registerClientServices(ctx context.Context, c *container.Container) {
	// Media clients
	registerClientService[*types.JellyfinConfig](c)
	registerClientService[*types.EmbyConfig](c)
	registerClientService[*types.PlexConfig](c)
	registerClientService[*types.SubsonicConfig](c)
	registerClientService[*types.RadarrConfig](c)
	registerClientService[*types.LidarrConfig](c)
	registerClientService[*types.SonarrConfig](c)
	registerClientService[*types.ClaudeConfig](c)
	registerClientService[*types.OpenAIConfig](c)
	registerClientService[*types.OllamaConfig](c)

	registerClientListService[*types.JellyfinConfig, *mediatypes.Collection](c)
	registerClientListService[*types.EmbyConfig, *mediatypes.Collection](c)
	registerClientListService[*types.PlexConfig, *mediatypes.Collection](c)
	registerClientListService[*types.SubsonicConfig, *mediatypes.Collection](c)

	registerClientListService[*types.EmbyConfig, *mediatypes.Playlist](c)
	registerClientListService[*types.JellyfinConfig, *mediatypes.Playlist](c)
	registerClientListService[*types.PlexConfig, *mediatypes.Playlist](c)
	registerClientListService[*types.SubsonicConfig, *mediatypes.Playlist](c)

}

func registerClientService[T types.ClientConfig](c *container.Container) {
	container.RegisterFactory[services.ClientService[T]](c, func(c *container.Container) services.ClientService[T] {
		clientFactory := container.MustGet[*client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[T]](c)
		return services.NewClientService[T](clientFactory, repo)
	})
}

func registerClientListService[T types.ClientMediaConfig, U mediatypes.ListData](c *container.Container) {
	container.RegisterFactory[services.ClientListService[T, U]](c, func(c *container.Container) services.ClientListService[T, U] {
		coreListService := container.MustGet[services.CoreListService[U]](c)
		clientRepo := container.MustGet[repository.ClientRepository[T]](c)
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		return services.NewClientListService[T, U](coreListService, clientRepo, &clientFactory)
	})
}

// Register services for the three-pronged architecture
func registerThreeProngedServices(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	// Core media item services
	log.Info().Msg("Registering core media item services")
	container.RegisterFactory[appservices.CoreMediaItemServices](c, func(c *container.Container) appservices.CoreMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		repos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		return factory.CreateCoreServices(repos)
	})

	// User media item services
	log.Info().Msg("Registering user media item services")
	container.RegisterFactory[appservices.UserMediaItemServices](c, func(c *container.Container) appservices.UserMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[appservices.CoreMediaItemServices](c)
		userRepos := container.MustGet[apprepository.UserMediaItemRepositories](c)
		return factory.CreateUserServices(coreServices, userRepos)
	})

	// Client media item services
	log.Info().Msg("Registering client media item services")
	container.RegisterFactory[appservices.ClientMediaItemServices[clienttypes.ClientMediaConfig]](c, func(c *container.Container) appservices.ClientMediaItemServices[clienttypes.ClientMediaConfig] {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[appservices.CoreMediaItemServices](c)
		clientRepo := container.MustGet[repository.ClientRepository[clienttypes.ClientMediaConfig]](c)
		clientRepos := container.MustGet[apprepository.ClientMediaItemRepositories](c)
		return factory.CreateClientServices(coreServices, clientRepo, clientRepos)
	})

	// Register Media Lists Services - Playlists and Collections
	log.Info().Msg("Registering media list services")
	diservices.RegisterMediaListServices(ctx, c)

	// Collection services
	log.Info().Msg("Registering collection services")
	container.RegisterFactory[services.CoreListService[*mediatypes.Collection]](c, func(c *container.Container) services.CoreListService[*mediatypes.Collection] {
		repos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		return services.NewCoreListService(repos.CollectionRepo())
	})

	log.Info().Msg("Registering user collection services")
	container.RegisterFactory[services.UserListService[*mediatypes.Collection]](c, func(c *container.Container) services.UserListService[*mediatypes.Collection] {
		coreService := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		userItemRepos := container.MustGet[apprepository.UserMediaItemRepositories](c)
		userDataRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Collection]](c)

		return services.NewUserListService(coreService, userItemRepos.CollectionUserRepo(), userDataRepo)
	})
}
