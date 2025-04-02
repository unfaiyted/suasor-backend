// app/dependencies.go
package app

import (
	"gorm.io/gorm"
	"suasor/client"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"
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
		embyRepo:     repository.NewClientRepository[*types.EmbyConfig](db),
		jellyfinRepo: repository.NewClientRepository[*types.JellyfinConfig](db),
		plexRepo:     repository.NewClientRepository[*types.PlexConfig](db),
		subsonicRepo: repository.NewClientRepository[*types.SubsonicConfig](db),
		sonarrRepo:   repository.NewClientRepository[*types.SonarrConfig](db),
		radarrRepo:   repository.NewClientRepository[*types.RadarrConfig](db),
		lidarrRepo:   repository.NewClientRepository[*types.LidarrConfig](db),
		claudeRepo:   repository.NewClientRepository[*types.ClaudeConfig](db),
		openaiRepo:   repository.NewClientRepository[*types.OpenAIConfig](db),
		ollamaRepo:   repository.NewClientRepository[*types.OllamaConfig](db),
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
		embyService:     services.NewClientService[*types.EmbyConfig](deps.ClientFactoryService, deps.ClientRepositories.EmbyRepo()),
		jellyfinService: services.NewClientService[*types.JellyfinConfig](deps.ClientFactoryService, deps.ClientRepositories.JellyfinRepo()),
		plexService:     services.NewClientService[*types.PlexConfig](deps.ClientFactoryService, deps.ClientRepositories.PlexRepo()),
		subsonicService: services.NewClientService[*types.SubsonicConfig](deps.ClientFactoryService, deps.ClientRepositories.SubsonicRepo()),
		sonarrService:   services.NewClientService[*types.SonarrConfig](deps.ClientFactoryService, deps.ClientRepositories.SonarrRepo()),
		radarrService:   services.NewClientService[*types.RadarrConfig](deps.ClientFactoryService, deps.ClientRepositories.RadarrRepo()),
		lidarrService:   services.NewClientService[*types.LidarrConfig](deps.ClientFactoryService, deps.ClientRepositories.LidarrRepo()),
		claudeService:   services.NewClientService[*types.ClaudeConfig](deps.ClientFactoryService, deps.ClientRepositories.ClaudeRepo()),
		openaiService:   services.NewClientService[*types.OpenAIConfig](deps.ClientFactoryService, deps.ClientRepositories.OpenAIRepo()),
		ollamaService:   services.NewClientService[*types.OllamaConfig](deps.ClientFactoryService, deps.ClientRepositories.OllamaRepo()),
	}

	// Initialize media client services
	deps.ClientMediaServices = &clientMediaServicesImpl{
		movieServices: clientMovieServicesImpl{
			embyMovieService:     services.NewMediaClientMovieService[*types.EmbyConfig](deps.ClientRepositories.EmbyRepo(), deps.ClientFactoryService),
			jellyfinMovieService: services.NewMediaClientMovieService[*types.JellyfinConfig](deps.ClientRepositories.JellyfinRepo(), deps.ClientFactoryService),
			plexMovieService:     services.NewMediaClientMovieService[*types.PlexConfig](deps.ClientRepositories.PlexRepo(), deps.ClientFactoryService),
			subsonicMovieService: services.NewMediaClientMovieService[*types.SubsonicConfig](deps.ClientRepositories.SubsonicRepo(), deps.ClientFactoryService),
		},
		seriesServices: clientSeriesServicesImpl{
			embySeriesService:     services.NewMediaClientSeriesService[*types.EmbyConfig](deps.ClientRepositories.EmbyRepo(), deps.ClientFactoryService),
			jellyfinSeriesService: services.NewMediaClientSeriesService[*types.JellyfinConfig](deps.ClientRepositories.JellyfinRepo(), deps.ClientFactoryService),
			plexSeriesService:     services.NewMediaClientSeriesService[*types.PlexConfig](deps.ClientRepositories.PlexRepo(), deps.ClientFactoryService),
			subsonicSeriesService: services.NewMediaClientSeriesService[*types.SubsonicConfig](deps.ClientRepositories.SubsonicRepo(), deps.ClientFactoryService),
		},
		musicServices: clientMusicServicesImpl{
			embyMusicService:     services.NewMediaClientMusicService[*types.EmbyConfig](deps.ClientRepositories.EmbyRepo(), deps.ClientFactoryService),
			jellyfinMusicService: services.NewMediaClientMusicService[*types.JellyfinConfig](deps.ClientRepositories.JellyfinRepo(), deps.ClientFactoryService),
			plexMusicService:     services.NewMediaClientMusicService[*types.PlexConfig](deps.ClientRepositories.PlexRepo(), deps.ClientFactoryService),
			subsonicMusicService: services.NewMediaClientMusicService[*types.SubsonicConfig](deps.ClientRepositories.SubsonicRepo(), deps.ClientFactoryService),
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
		embyHandler:     handlers.NewClientHandler[*types.EmbyConfig](deps.ClientServices.EmbyService()),
		jellyfinHandler: handlers.NewClientHandler[*types.JellyfinConfig](deps.ClientServices.JellyfinService()),
		plexHandler:     handlers.NewClientHandler[*types.PlexConfig](deps.ClientServices.PlexService()),
		subsonicHandler: handlers.NewClientHandler[*types.SubsonicConfig](deps.ClientServices.SubsonicService()),
		radarrHandler:   handlers.NewClientHandler[*types.RadarrConfig](deps.ClientServices.RadarrService()),
		lidarrHandler:   handlers.NewClientHandler[*types.LidarrConfig](deps.ClientServices.LidarrService()),
		sonarrHandler:   handlers.NewClientHandler[*types.SonarrConfig](deps.ClientServices.SonarrService()),
		claudeHandler:   handlers.NewClientHandler[*types.ClaudeConfig](deps.ClientServices.ClaudeService()),
		openaiHandler:   handlers.NewClientHandler[*types.OpenAIConfig](deps.ClientServices.OpenAIService()),
		ollamaHandler:   handlers.NewClientHandler[*types.OllamaConfig](deps.ClientServices.OllamaService()),
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

	deps.UserServices = &userServicesImpl{
		userService:       services.NewUserService(deps.UserRepo()),
		userConfigService: services.NewUserConfigService(deps.UserConfigRepo()),
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
		userHandler:       handlers.NewUserHandler(deps.UserService()),
		userConfigHandler: handlers.NewUserConfigHandler(deps.UserConfigService()),
	}

	clientMovieHandlers := &clientMediaMovieHandlersImpl{
		embyMovieHandler:     handlers.NewMediaClientMovieHandler[*types.EmbyConfig](deps.ClientMediaServices.EmbyMovieService()),
		jellyfinMovieHandler: handlers.NewMediaClientMovieHandler[*types.JellyfinConfig](deps.ClientMediaServices.JellyfinMovieService()),
		plexMovieHandler:     handlers.NewMediaClientMovieHandler[*types.PlexConfig](deps.ClientMediaServices.PlexMovieService()),
	}

	clientSeriesHandlers := &clientMediaSeriesHandlersImpl{
		embySeriesHandler:     handlers.NewMediaClientSeriesHandler[*types.EmbyConfig](deps.ClientMediaServices.EmbySeriesService()),
		jellyfinSeriesHandler: handlers.NewMediaClientSeriesHandler[*types.JellyfinConfig](deps.ClientMediaServices.JellyfinSeriesService()),
		plexSeriesHandler:     handlers.NewMediaClientSeriesHandler[*types.PlexConfig](deps.ClientMediaServices.PlexSeriesService()),
	}

	clientMusicHandlers := &clientMediaMusicHandlersImpl{
		embyMusicHandler:     handlers.NewMediaClientMusicHandler[*types.EmbyConfig](deps.ClientMediaServices.EmbyMusicService()),
		jellyfinMusicHandler: handlers.NewMediaClientMusicHandler[*types.JellyfinConfig](deps.ClientMediaServices.JellyfinMusicService()),
		plexMusicHandler:     handlers.NewMediaClientMusicHandler[*types.PlexConfig](deps.ClientMediaServices.PlexMusicService()),
		subsonicMusicHandler: handlers.NewMediaClientMusicHandler[*types.SubsonicConfig](deps.ClientMediaServices.SubsonicMusicService()),
	}

	deps.ClientMediaHandlers = &clientMediaHandlersImpl{
		movieHandlers:  clientMovieHandlers,
		seriesHandlers: clientSeriesHandlers,
		musicHandlers:  clientMusicHandlers,
	}

	return deps
}
