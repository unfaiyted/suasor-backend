// app/dependencies.go
package app

import (
	"gorm.io/gorm"
	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"
	"suasor/services/jobs"
	"time"
)

func InitializeDependencies(db *gorm.DB, configService services.ConfigService) *AppDependencies {
	deps := &AppDependencies{}

	// NOTE: The Config Service represents the file configuration for the app itself. Not the user configuraiton.
	appConfig := configService.GetConfig()

	clientFactory := client.GetClientFactoryService()

	deps.SystemRepositories = &systemRepositoriesImpl{
		configRepo: configService.GetRepo(),
	}

	// Initialize repositories
	deps.UserRepositories = &userRepositoriesImpl{
		userRepo:       repository.NewUserRepository(db),
		userConfigRepo: repository.NewUserConfigRepository(db),
		sessionRepo:    repository.NewSessionRepository(db),
	}

	deps.ClientRepositories = &clientRepositoriesImpl{
		embyRepo:     repository.NewClientRepository[*clienttypes.EmbyConfig](db),
		jellyfinRepo: repository.NewClientRepository[*clienttypes.JellyfinConfig](db),
		plexRepo:     repository.NewClientRepository[*clienttypes.PlexConfig](db),
		subsonicRepo: repository.NewClientRepository[*clienttypes.SubsonicConfig](db),
		sonarrRepo:   repository.NewClientRepository[*clienttypes.SonarrConfig](db),
		radarrRepo:   repository.NewClientRepository[*clienttypes.RadarrConfig](db),
		lidarrRepo:   repository.NewClientRepository[*clienttypes.LidarrConfig](db),
		claudeRepo:   repository.NewClientRepository[*clienttypes.ClaudeConfig](db),
		openaiRepo:   repository.NewClientRepository[*clienttypes.OpenAIConfig](db),
		ollamaRepo:   repository.NewClientRepository[*clienttypes.OllamaConfig](db),
	}

	deps.MediaItemRepositories = &mediaItemRepositoriesImpl{
		movieRepo:      repository.NewMediaItemRepository[*mediatypes.Movie](db),
		seriesRepo:     repository.NewMediaItemRepository[*mediatypes.Series](db),
		episodeRepo:    repository.NewMediaItemRepository[*mediatypes.Episode](db),
		trackRepo:      repository.NewMediaItemRepository[*mediatypes.Track](db),
		albumRepo:      repository.NewMediaItemRepository[*mediatypes.Album](db),
		artistRepo:     repository.NewMediaItemRepository[*mediatypes.Artist](db),
		collectionRepo: repository.NewMediaItemRepository[*mediatypes.Collection](db),
		playlistRepo:   repository.NewMediaItemRepository[*mediatypes.Playlist](db),
	}

	// Store the client factory service
	deps.ClientFactoryService = clientFactory

	// Initialize client services
	deps.ClientServices = &clientServicesImpl{
		embyService:     services.NewClientService[*clienttypes.EmbyConfig](deps.ClientFactoryService, deps.ClientRepositories.EmbyRepo()),
		jellyfinService: services.NewClientService[*clienttypes.JellyfinConfig](deps.ClientFactoryService, deps.ClientRepositories.JellyfinRepo()),
		plexService:     services.NewClientService[*clienttypes.PlexConfig](deps.ClientFactoryService, deps.ClientRepositories.PlexRepo()),
		subsonicService: services.NewClientService[*clienttypes.SubsonicConfig](deps.ClientFactoryService, deps.ClientRepositories.SubsonicRepo()),
		sonarrService:   services.NewClientService[*clienttypes.SonarrConfig](deps.ClientFactoryService, deps.ClientRepositories.SonarrRepo()),
		radarrService:   services.NewClientService[*clienttypes.RadarrConfig](deps.ClientFactoryService, deps.ClientRepositories.RadarrRepo()),
		lidarrService:   services.NewClientService[*clienttypes.LidarrConfig](deps.ClientFactoryService, deps.ClientRepositories.LidarrRepo()),
		claudeService:   services.NewClientService[*clienttypes.ClaudeConfig](deps.ClientFactoryService, deps.ClientRepositories.ClaudeRepo()),
		openaiService:   services.NewClientService[*clienttypes.OpenAIConfig](deps.ClientFactoryService, deps.ClientRepositories.OpenAIRepo()),
		ollamaService:   services.NewClientService[*clienttypes.OllamaConfig](deps.ClientFactoryService, deps.ClientRepositories.OllamaRepo()),
	}

	// Initialize media client services
	deps.ClientMediaServices = &clientMediaServicesImpl{
		movieServices: clientMovieServicesImpl{
			embyMovieService:     services.NewMediaClientMovieService[*clienttypes.EmbyConfig](deps.ClientRepositories.EmbyRepo(), deps.ClientFactoryService),
			jellyfinMovieService: services.NewMediaClientMovieService[*clienttypes.JellyfinConfig](deps.ClientRepositories.JellyfinRepo(), deps.ClientFactoryService),
			plexMovieService:     services.NewMediaClientMovieService[*clienttypes.PlexConfig](deps.ClientRepositories.PlexRepo(), deps.ClientFactoryService),
			subsonicMovieService: services.NewMediaClientMovieService[*clienttypes.SubsonicConfig](deps.ClientRepositories.SubsonicRepo(), deps.ClientFactoryService),
		},
		seriesServices: clientSeriesServicesImpl{
			embySeriesService:     services.NewMediaClientSeriesService[*clienttypes.EmbyConfig](deps.ClientRepositories.EmbyRepo(), deps.ClientFactoryService),
			jellyfinSeriesService: services.NewMediaClientSeriesService[*clienttypes.JellyfinConfig](deps.ClientRepositories.JellyfinRepo(), deps.ClientFactoryService),
			plexSeriesService:     services.NewMediaClientSeriesService[*clienttypes.PlexConfig](deps.ClientRepositories.PlexRepo(), deps.ClientFactoryService),
			subsonicSeriesService: services.NewMediaClientSeriesService[*clienttypes.SubsonicConfig](deps.ClientRepositories.SubsonicRepo(), deps.ClientFactoryService),
		},
		musicServices: clientMusicServicesImpl{
			embyMusicService:     services.NewMediaClientMusicService[*clienttypes.EmbyConfig](deps.ClientRepositories.EmbyRepo(), deps.ClientFactoryService),
			jellyfinMusicService: services.NewMediaClientMusicService[*clienttypes.JellyfinConfig](deps.ClientRepositories.JellyfinRepo(), deps.ClientFactoryService),
			plexMusicService:     services.NewMediaClientMusicService[*clienttypes.PlexConfig](deps.ClientRepositories.PlexRepo(), deps.ClientFactoryService),
			subsonicMusicService: services.NewMediaClientMusicService[*clienttypes.SubsonicConfig](deps.ClientRepositories.SubsonicRepo(), deps.ClientFactoryService),
		},
		episodeServices:  clientEpisodeServicesImpl{},
		playlistServices: clientPlaylistServicesImpl{},
	}

	deps.SystemServices = &systemServicesImpl{
		healthService: services.NewHealthService(db),
		configService: configService,
	}

	// Initialize media item services
	deps.MediaItemServices = &mediaItemServicesImpl{
		movieService:      services.NewMediaItemService[*mediatypes.Movie](deps.MediaItemRepositories.MovieRepo()),
		seriesService:     services.NewMediaItemService[*mediatypes.Series](deps.MediaItemRepositories.SeriesRepo()),
		episodeService:    services.NewMediaItemService[*mediatypes.Episode](deps.MediaItemRepositories.EpisodeRepo()),
		trackService:      services.NewMediaItemService[*mediatypes.Track](deps.MediaItemRepositories.TrackRepo()),
		albumService:      services.NewMediaItemService[*mediatypes.Album](deps.MediaItemRepositories.AlbumRepo()),
		artistService:     services.NewMediaItemService[*mediatypes.Artist](deps.MediaItemRepositories.ArtistRepo()),
		collectionService: services.NewMediaItemService[*mediatypes.Collection](deps.MediaItemRepositories.CollectionRepo()),
		playlistService:   services.NewMediaItemService[*mediatypes.Playlist](deps.MediaItemRepositories.PlaylistRepo()),
	}

	// Initialize client handlers
	deps.ClientHandlers = &clientHandlersImpl{
		embyHandler:     handlers.NewClientHandler[*clienttypes.EmbyConfig](deps.ClientServices.EmbyService()),
		jellyfinHandler: handlers.NewClientHandler[*clienttypes.JellyfinConfig](deps.ClientServices.JellyfinService()),
		plexHandler:     handlers.NewClientHandler[*clienttypes.PlexConfig](deps.ClientServices.PlexService()),
		subsonicHandler: handlers.NewClientHandler[*clienttypes.SubsonicConfig](deps.ClientServices.SubsonicService()),
		radarrHandler:   handlers.NewClientHandler[*clienttypes.RadarrConfig](deps.ClientServices.RadarrService()),
		lidarrHandler:   handlers.NewClientHandler[*clienttypes.LidarrConfig](deps.ClientServices.LidarrService()),
		sonarrHandler:   handlers.NewClientHandler[*clienttypes.SonarrConfig](deps.ClientServices.SonarrService()),
		claudeHandler:   handlers.NewClientHandler[*clienttypes.ClaudeConfig](deps.ClientServices.ClaudeService()),
		openaiHandler:   handlers.NewClientHandler[*clienttypes.OpenAIConfig](deps.ClientServices.OpenAIService()),
		ollamaHandler:   handlers.NewClientHandler[*clienttypes.OllamaConfig](deps.ClientServices.OllamaService()),
	}

	deps.MediaItemHandlers = &mediaItemHandlersImpl{
		movieHandler:      handlers.NewMediaItemHandler[*mediatypes.Movie](deps.MediaItemServices.MovieService()),
		seriesHandler:     handlers.NewMediaItemHandler[*mediatypes.Series](deps.MediaItemServices.SeriesService()),
		episodeHandler:    handlers.NewMediaItemHandler[*mediatypes.Episode](deps.MediaItemServices.EpisodeService()),
		trackHandler:      handlers.NewMediaItemHandler[*mediatypes.Track](deps.MediaItemServices.TrackService()),
		albumHandler:      handlers.NewMediaItemHandler[*mediatypes.Album](deps.MediaItemServices.AlbumService()),
		artistHandler:     handlers.NewMediaItemHandler[*mediatypes.Artist](deps.MediaItemServices.ArtistService()),
		collectionHandler: handlers.NewMediaItemHandler[*mediatypes.Collection](deps.MediaItemServices.CollectionService()),
		playlistHandler:   handlers.NewMediaItemHandler[*mediatypes.Playlist](deps.MediaItemServices.PlaylistService()),
	}

	// System Handlers
	deps.SystemHandlers = &systemHandlersImpl{
		configHandler: handlers.NewConfigHandler(deps.SystemServices.ConfigService()),
		healthHandler: handlers.NewHealthHandler(deps.SystemServices.HealthService()),

		clientsHandler: handlers.NewClientsHandler(
			deps.ClientServices.EmbyService(),
			deps.ClientServices.JellyfinService(),
			deps.ClientServices.PlexService(),
			deps.ClientServices.SubsonicService(),
			deps.ClientServices.SonarrService(),
			deps.ClientServices.RadarrService(),
			deps.ClientServices.LidarrService(),
			deps.ClientServices.ClaudeService(),
			deps.ClientServices.OpenAIService(),
			deps.ClientServices.OllamaService(),
		),
	}

	deps.AIHandlers = &aiHandlersImpl{
		claudeAIHandler: *handlers.NewAIHandler(
			deps.ClientFactoryService,
			deps.ClientServices.ClaudeService(),
		),
		openaiAIHandler: *handlers.NewAIHandler(
			deps.ClientFactoryService,
			deps.ClientServices.OpenAIService(),
		),
		ollamaAIHandler: *handlers.NewAIHandler(
			deps.ClientFactoryService,
			deps.ClientServices.OllamaService(),
		),
	}

	// Initialize additional repositories
	historyRepo := repository.NewMediaPlayHistoryRepository(db)

	// Initialize job repositories
	deps.JobRepositories = &jobRepositoriesImpl{
		jobRepo: repository.NewJobRepository(db),
	}

	// Initialize job services

	// Get a Claude AI client for recommendations (if configured)
	// This could be any AI client type that implements the required interface
	// var aiClientService interface{}
	// claudeClients, _ := deps.ClientRepositories.ClaudeRepo.GetByUserId(ctx)
	// if len(claudeClients) > 0 {
	// 	// Use the first Claude client found
	// 	clientID := claudeClients[0].ID
	// 	clientConfig := claudeClients[0].Config.Data
	// 	aiClient, err := deps.ClientFactoryService.GetClient(ctx, clientID, clientConfig)
	// 	if err == nil {
	// 		aiClientService = aiClient
	// 		log.Info().
	// 			Uint64("clientID", clientID).
	// 			Msg("Initialized AI client for recommendation service")
	// 	}
	// }

	recommendationJob := jobs.NewRecommendationJob(
		deps.JobRepo(),
		deps.UserRepo(),
		deps.UserConfigRepo(),
		deps.MovieRepo(),
		deps.SeriesRepo(),
		deps.TrackRepo(),
		historyRepo,
		aiClientService,
	)

	mediaSyncJob := jobs.NewMediaSyncJob(
		deps.JobRepo(),
		deps.UserRepo(),
		deps.UserConfigRepo(),
		deps.MovieRepo(),
		deps.SeriesRepo(),
		deps.EpisodeRepo(),
		deps.TrackRepo(),
		deps.AlbumRepo(),
		deps.ArtistRepo(),
		deps.ClientRepositories,
		client.GetClientFactoryService(),
	)

	watchHistorySyncJob := jobs.NewWatchHistorySyncJob(
		deps.JobRepo(),
		deps.UserRepo(),
		deps.UserConfigRepo(),
		historyRepo,
		deps.MovieRepo(),
		deps.SeriesRepo(),
		deps.MediaItemRepositories.EpisodeRepo(),
		deps.TrackRepo(),
		deps.ClientRepositories,
		deps.ClientFactoryService,
	)

	favoritesSyncJob := jobs.NewFavoritesSyncJob(
		deps.JobRepo(),
		deps.UserRepo(),
		deps.UserConfigRepo(),
		historyRepo,
		deps.MovieRepo(),
		deps.SeriesRepo(),
		deps.MediaItemRepositories.EpisodeRepo(),
		deps.TrackRepo(),
		deps.ClientRepositories,
		deps.ClientFactoryService,
	)

	jobService := services.NewJobService(
		deps.JobRepo(),
		deps.UserRepo(),
		deps.UserConfigRepo(),
		deps.MovieRepo(),
		deps.SeriesRepo(),
		deps.TrackRepo(),
		historyRepo,
		recommendationJob,
		mediaSyncJob,
		watchHistorySyncJob,
		favoritesSyncJob,
	)

	deps.JobServices = &jobServicesImpl{
		jobService:          jobService,
		recommendationJob:   recommendationJob,
		mediaSyncJob:        mediaSyncJob,
		watchHistorySyncJob: watchHistorySyncJob,
		favoritesSyncJob:    favoritesSyncJob,
	}

	// Initialize job handlers
	deps.JobHandlers = &jobHandlersImpl{
		jobHandler: handlers.NewJobHandler(jobService),
	}

	// Register jobs with the job service
	jobService.RegisterJob(recommendationJob)
	jobService.RegisterJob(mediaSyncJob)
	jobService.RegisterJob(watchHistorySyncJob)
	jobService.RegisterJob(favoritesSyncJob)

	// Initialize and register additional system jobs where we have implementations
	// Only register jobs that we can create properly based on available implementations

	// Database maintenance job only needs JobRepo
	databaseMaintenanceJob := jobs.NewDatabaseMaintenanceJob(
		deps.JobRepo(),
	)
	jobService.RegisterJob(databaseMaintenanceJob)

	deps.UserServices = &userServicesImpl{
		userService: services.NewUserService(deps.UserRepo()),
		userConfigService: services.NewUserConfigService(
			deps.UserConfigRepo(),
			deps.JobServices.JobService(),
			deps.JobServices.RecommendationJob(),
		),
		authService: services.NewAuthService(deps.UserRepo(),
			deps.SessionRepo(),
			appConfig.Auth.JWTSecret,
			time.Duration(appConfig.Auth.AccessExpiryMinutes)*time.Minute,
			time.Duration(appConfig.Auth.RefreshExpiryDays)*24*time.Hour,
			appConfig.Auth.TokenIssuer,
			appConfig.Auth.TokenAudience),
	}

	// User Handlers
	deps.UserHandlers = &userHandlersImpl{
		authHandler:       handlers.NewAuthHandler(deps.AuthService()),
		userHandler:       handlers.NewUserHandler(deps.UserService(), deps.SystemServices.ConfigService()),
		userConfigHandler: handlers.NewUserConfigHandler(deps.UserConfigService()),
	}

	clientMovieHandlers := &clientMediaMovieHandlersImpl{
		embyMovieHandler:     handlers.NewMediaClientMovieHandler[*clienttypes.EmbyConfig](deps.ClientMediaServices.EmbyMovieService()),
		jellyfinMovieHandler: handlers.NewMediaClientMovieHandler[*clienttypes.JellyfinConfig](deps.ClientMediaServices.JellyfinMovieService()),
		plexMovieHandler:     handlers.NewMediaClientMovieHandler[*clienttypes.PlexConfig](deps.ClientMediaServices.PlexMovieService()),
	}

	clientSeriesHandlers := &clientMediaSeriesHandlersImpl{
		embySeriesHandler:     handlers.NewMediaClientSeriesHandler[*clienttypes.EmbyConfig](deps.ClientMediaServices.EmbySeriesService()),
		jellyfinSeriesHandler: handlers.NewMediaClientSeriesHandler[*clienttypes.JellyfinConfig](deps.ClientMediaServices.JellyfinSeriesService()),
		plexSeriesHandler:     handlers.NewMediaClientSeriesHandler[*clienttypes.PlexConfig](deps.ClientMediaServices.PlexSeriesService()),
	}

	clientMusicHandlers := &clientMediaMusicHandlersImpl{
		embyMusicHandler:     handlers.NewMediaClientMusicHandler[*clienttypes.EmbyConfig](deps.ClientMediaServices.EmbyMusicService()),
		jellyfinMusicHandler: handlers.NewMediaClientMusicHandler[*clienttypes.JellyfinConfig](deps.ClientMediaServices.JellyfinMusicService()),
		plexMusicHandler:     handlers.NewMediaClientMusicHandler[*clienttypes.PlexConfig](deps.ClientMediaServices.PlexMusicService()),
		subsonicMusicHandler: handlers.NewMediaClientMusicHandler[*clienttypes.SubsonicConfig](deps.ClientMediaServices.SubsonicMusicService()),
	}

	deps.ClientMediaHandlers = &clientMediaHandlersImpl{
		movieHandlers:  clientMovieHandlers,
		seriesHandlers: clientSeriesHandlers,
		musicHandlers:  clientMusicHandlers,
	}

	return deps
}
