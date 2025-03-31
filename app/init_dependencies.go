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

func InitializeDependencies(db *gorm.DB, configService services.ConfigService, factory *client.ClientFactoryService) *AppDependencies {
	deps := &AppDependencies{}

	appConfig := configService.GetConfig()

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
	deps.ClientFactoryService = factory

	// Initialize client services
	deps.ClientServices = &clientServicesImpl{
		embyService:     services.NewClientService[*types.EmbyConfig](factory, deps.ClientRepositories.EmbyRepo()),
		jellyfinService: services.NewClientService[*types.JellyfinConfig](factory, deps.ClientRepositories.JellyfinRepo()),
		plexService:     services.NewClientService[*types.PlexConfig](factory, deps.ClientRepositories.PlexRepo()),
		subsonicService: services.NewClientService[*types.SubsonicConfig](factory, deps.ClientRepositories.SubsonicRepo()),
		sonarrService:   services.NewClientService[*types.SonarrConfig](factory, deps.ClientRepositories.SonarrRepo()),
		radarrService:   services.NewClientService[*types.RadarrConfig](factory, deps.ClientRepositories.RadarrRepo()),
		lidarrService:   services.NewClientService[*types.LidarrConfig](factory, deps.ClientRepositories.LidarrRepo()),
	}

	// Initialize media client services
	deps.ClientMediaServices = &clientMediaServicesImpl{
		movieServices: clientMovieServicesImpl{
			embyMovieService:     services.NewMediaClientMovieService[*types.EmbyConfig](deps.ClientRepositories.EmbyRepo(), factory),
			jellyfinMovieService: services.NewMediaClientMovieService[*types.JellyfinConfig](deps.ClientRepositories.JellyfinRepo(), factory),
			plexMovieService:     services.NewMediaClientMovieService[*types.PlexConfig](deps.ClientRepositories.PlexRepo(), factory),
		},
		seriesServices:   clientSeriesServicesImpl{},
		episodeServices:  clientEpisodeServicesImpl{},
		playlistServices: clientPlaylistServicesImpl{},
	}

	deps.SystemServices = &systemServicesImpl{
		healthService: services.NewHealthService(db),
		configService: services.NewConfigService(deps.ConfigRepo()),
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
	}

	deps.MediaHandlers = &mediaHandlersImpl{
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
	// Client Handlers
	deps.ClientHandlers = &clientHandlersImpl{
		embyHandler:     handlers.NewClientHandler[*types.EmbyConfig](deps.ClientServices.EmbyService()),
		jellyfinHandler: handlers.NewClientHandler[*types.JellyfinConfig](deps.ClientServices.JellyfinService()),
		plexHandler:     handlers.NewClientHandler[*types.PlexConfig](deps.ClientServices.PlexService()),
		subsonicHandler: handlers.NewClientHandler[*types.SubsonicConfig](deps.ClientServices.SubsonicService()),
		radarrHandler:   handlers.NewClientHandler[*types.RadarrConfig](deps.ClientServices.RadarrService()),
		lidarrHandler:   handlers.NewClientHandler[*types.LidarrConfig](deps.ClientServices.LidarrService()),
		sonarrHandler:   handlers.NewClientHandler[*types.SonarrConfig](deps.ClientServices.SonarrService()),
	}

	return deps
}
