// app/dependencies.go
package app

import (
	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"
	// "suasor/services/jobs"
	// "suasor/services/jobs/recommendation"
	"time"

	"gorm.io/gorm"
)

func InitializeDependencies(db *gorm.DB, configService services.ConfigService) *AppDependencies {
	deps := &AppDependencies{
		db: db,
	}

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

	// Initialize people and credits repositories
	personRepo := repository.NewPersonRepository(db)
	creditRepo := repository.NewCreditRepository(db)

	// Initialize media services for people and credits
	deps.MediaServices = NewMediaServices(personRepo, creditRepo)

	// Initialize client repositories
	clientRepos := &clientRepositoriesImpl{
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

	deps.RepositoryCollections = &repositoryCollectionsImpl{
		clientRepos: repository.NewClientRepositoryCollection(
			clientRepos.EmbyRepo(),
			clientRepos.JellyfinRepo(),
			clientRepos.PlexRepo(),
			clientRepos.SubsonicRepo(),
			clientRepos.SonarrRepo(),
			clientRepos.RadarrRepo(),
			clientRepos.LidarrRepo(),
			clientRepos.ClaudeRepo(),
			clientRepos.OpenAIRepo(),
			clientRepos.OllamaRepo(),
		),
	}

	// Initialize the media data factory
	mediaDataFactory := NewMediaDataFactory(db, clientFactory)

	// Initialize repositories using the factory
	deps.CoreMediaItemRepositories = mediaDataFactory.CreateCoreRepositories()
	deps.CoreUserMediaItemDataRepositories = mediaDataFactory.CreateCoreDataRepositories()
	deps.UserRepositoryFactories = mediaDataFactory.CreateUserRepositories()
	deps.ClientRepositoryFactories = mediaDataFactory.CreateClientRepositories()
	deps.UserDataFactories = mediaDataFactory.CreateUserDataRepositories()

	// We're no longer using the legacy MediaItemRepositories
	// Instead, we're using CoreRepositories, UserRepositoryFactories, and ClientRepositoryFactories directly

	// Store the client factory service
	deps.ClientFactoryService = clientFactory

	// Initialize client services
	deps.ClientServices = &clientServicesImpl{
		embyService:     services.NewClientService[*clienttypes.EmbyConfig](deps.ClientFactoryService, clientRepos.EmbyRepo()),
		jellyfinService: services.NewClientService[*clienttypes.JellyfinConfig](deps.ClientFactoryService, clientRepos.JellyfinRepo()),
		plexService:     services.NewClientService[*clienttypes.PlexConfig](deps.ClientFactoryService, clientRepos.PlexRepo()),
		subsonicService: services.NewClientService[*clienttypes.SubsonicConfig](deps.ClientFactoryService, clientRepos.SubsonicRepo()),
		sonarrService:   services.NewClientService[*clienttypes.SonarrConfig](deps.ClientFactoryService, clientRepos.SonarrRepo()),
		radarrService:   services.NewClientService[*clienttypes.RadarrConfig](deps.ClientFactoryService, clientRepos.RadarrRepo()),
		lidarrService:   services.NewClientService[*clienttypes.LidarrConfig](deps.ClientFactoryService, clientRepos.LidarrRepo()),
		claudeService:   services.NewClientService[*clienttypes.ClaudeConfig](deps.ClientFactoryService, clientRepos.ClaudeRepo()),
		openaiService:   services.NewClientService[*clienttypes.OpenAIConfig](deps.ClientFactoryService, clientRepos.OpenAIRepo()),
		ollamaService:   services.NewClientService[*clienttypes.OllamaConfig](deps.ClientFactoryService, clientRepos.OllamaRepo()),
	}

	// We're now using the pure three-pronged architecture for client media services
	// The legacy ClientMediaServices is being replaced by ClientMediaItemServices

	deps.SystemServices = &systemServicesImpl{
		healthService: services.NewHealthService(db),
		configService: configService,
	}

	// Initialize three-pronged services using the factory
	deps.CoreMediaItemServices = mediaDataFactory.CreateCoreServices(deps.CoreMediaItemRepositories)
	deps.UserMediaItemServices = mediaDataFactory.CreateUserServices(deps.CoreMediaItemServices, deps.UserRepositoryFactories)
	deps.ClientMediaItemServices = mediaDataFactory.CreateClientServices(deps.CoreMediaItemServices, deps.ClientRepositoryFactories)

	// Initialize collection services with three-pronged architecture
	// Core collection service
	collectionCoreService := services.NewCoreCollectionService(
		deps.CoreMediaItemRepositories.CollectionRepo())

	// Create core repository for media data
	coreMediaDataRepo := repository.NewMediaItemRepository[mediatypes.MediaData](db)

	// User collection service (extends core)
	collectionUserService := services.NewUserCollectionService(
		collectionCoreService,
		deps.UserRepositoryFactories.CollectionUserRepo(),
		coreMediaDataRepo,
	)
	
	// Client collection service (extends user)
	collectionClientService := services.NewClientMediaCollectionService(
		collectionUserService,
		deps.ClientRepositoryFactories.CollectionClientRepo(),
		nil, // client repo - we'll use nil since it's not directly used in this context
		deps.ClientFactoryService,
	)

	// Initialize specialized playlist service
	playlistExtendedService := services.NewPlaylistService(
		deps.UserRepositoryFactories.PlaylistUserRepo(),
		deps.UserMediaItemServices.PlaylistUserService(),
		coreMediaDataRepo,
	)

	// Store the factory in the dependencies
	deps.MediaDataFactory = mediaDataFactory

	// Create specialized media collection services and store it for use
	mediaCollectionServices := mediaDataFactory.CreateMediaCollectionServices(
		deps.CoreMediaItemServices,
		deps.UserMediaItemServices,
		deps.ClientMediaItemServices,
		collectionCoreService,
		collectionUserService,
		collectionClientService,
		playlistExtendedService,
	)

	// For the playlistSpecificHandler and collectionSpecificHandler creation
	// We need to use the new services instead of the old MediaItemServices

	// For the seriesSpecificHandler, using a direct instantiation for now
	// This will be replaced by a proper constructor later
	seriesSpecificHandler := CreateClientMediaSeriesHandler(
		deps.ClientServices.JellyfinService(),
		deps.ClientMediaItemServices.SeriesClientService(),
	)

	// Create specialized handler implementations
	// For the musicHandler, use the pure three-pronged structure
	specializedHandlers := mediaDataFactory.CreateSpecializedMediaHandlers(
		deps.CoreMediaItemServices,
		deps.UserMediaItemServices,
		deps.ClientMediaItemServices,
		handlers.NewCoreMusicHandler(
			deps.ClientMediaItemServices.TrackClientService(),
			deps.ClientMediaItemServices.AlbumClientService(),
			deps.ClientMediaItemServices.ArtistClientService(),
		),
		seriesSpecificHandler,
	)
	
	// Store specialized handlers
	deps.SpecializedMediaHandlers = specializedHandlers

	// We're no longer using playlist and collection specific handlers
	// These would be created using the code below but we're not storing them
	// since we're using the three-pronged architecture now
	/*
	playlistSpecificHandler := handlers.NewCorePlaylistHandler(
		deps.CoreMediaItemServices.PlaylistCoreService(),
		mediaCollectionServices.PlaylistService(),
	)

	collectionSpecificHandler := handlers.NewCoreCollectionHandler(
		deps.CoreMediaItemServices.CollectionCoreService(),
		mediaCollectionServices.CoreCollectionService(),
	)
	*/

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

	deps.UserMediaItemDataServices = &userMediaItemDataServicesImpl{
		movieDataService: services.NewUserMediaItemDataService[*mediatypes.Movie](
			deps.CoreUserMediaItemDataServices.MovieCoreService(),
			deps.UserDataFactories.MovieDataRepo(),
		),
		seriesDataService: services.NewUserMediaItemDataService[*mediatypes.Series](
			deps.CoreUserMediaItemDataServices.SeriesCoreService(),
			deps.UserDataFactories.SeriesDataRepo(),
		),
		episodeDataService: services.NewUserMediaItemDataService[*mediatypes.Episode](
			deps.CoreUserMediaItemDataServices.EpisodeCoreService(),
			deps.UserDataFactories.EpisodeDataRepo(),
		),
		trackDataService: services.NewUserMediaItemDataService[*mediatypes.Track](
			deps.CoreUserMediaItemDataServices.TrackCoreService(),
			deps.UserDataFactories.TrackDataRepo(),
		),
		albumDataService: services.NewUserMediaItemDataService[*mediatypes.Album](
			deps.CoreUserMediaItemDataServices.AlbumCoreService(),
			deps.UserDataFactories.AlbumDataRepo(),
		),
		artistDataService: services.NewUserMediaItemDataService[*mediatypes.Artist](
			deps.CoreUserMediaItemDataServices.ArtistCoreService(),
			deps.UserDataFactories.ArtistDataRepo(),
		),
		collectionDataService: services.NewUserMediaItemDataService[*mediatypes.Collection](
			deps.CoreUserMediaItemDataServices.CollectionCoreService(),
			deps.UserDataFactories.CollectionDataRepo(),
		),
		playlistDataService: services.NewUserMediaItemDataService[*mediatypes.Playlist](
			deps.CoreUserMediaItemDataServices.PlaylistCoreService(),
			deps.UserDataFactories.PlaylistDataRepo(),
		),
	}

	deps.ClientUserMediaItemDataServices = &clientUserMediaItemDataServicesImpl{
		movieClientService: services.NewClientUserMediaItemDataService[*mediatypes.Movie](
			deps.UserMediaItemDataServices.MovieDataService(),
			deps.ClientUserDataRepositories.MovieDataRepo(),
		),
		seriesClientService: services.NewClientUserMediaItemDataService[*mediatypes.Series](
			deps.UserMediaItemDataServices.SeriesDataService(),
			deps.ClientUserDataRepositories.SeriesDataRepo(),
		),
		episodeClientService: services.NewClientUserMediaItemDataService[*mediatypes.Episode](
			deps.UserMediaItemDataServices.EpisodeDataService(),
			deps.ClientUserDataRepositories.EpisodeDataRepo(),
		),
		trackClientService: services.NewClientUserMediaItemDataService[*mediatypes.Track](
			deps.UserMediaItemDataServices.TrackDataService(),
			deps.ClientUserDataRepositories.TrackDataRepo(),
		),
		albumClientService: services.NewClientUserMediaItemDataService[*mediatypes.Album](
			deps.UserMediaItemDataServices.AlbumDataService(),
			deps.ClientUserDataRepositories.AlbumDataRepo(),
		),
		artistClientService: services.NewClientUserMediaItemDataService[*mediatypes.Artist](
			deps.UserMediaItemDataServices.ArtistDataService(),
			deps.ClientUserDataRepositories.ArtistDataRepo(),
		),
		collectionClientService: services.NewClientUserMediaItemDataService[*mediatypes.Collection](
			deps.UserMediaItemDataServices.CollectionDataService(),
			deps.ClientUserDataRepositories.CollectionDataRepo(),
		),
		playlistClientService: services.NewClientUserMediaItemDataService[*mediatypes.Playlist](
			deps.UserMediaItemDataServices.PlaylistDataService(),
			deps.ClientUserDataRepositories.PlaylistDataRepo(),
		),
	}

	// Initialize three-pronged handlers using the factory
	deps.CoreMediaItemHandlers = mediaDataFactory.CreateCoreHandlers(deps.CoreMediaItemServices)
	deps.UserMediaItemHandlers = mediaDataFactory.CreateUserHandlers(
		deps.UserMediaItemServices,
		deps.UserMediaItemDataServices,
		deps.CoreMediaItemHandlers)
	deps.ClientMediaItemHandlers = mediaDataFactory.CreateClientHandlers(
		deps.ClientMediaItemServices,
		deps.ClientUserMediaItemDataServices,
		deps.UserMediaItemHandlers)

	// Specialized handlers are already created above

	// Store the MediaCollectionServices in dependencies
	deps.MediaCollectionServices = mediaCollectionServices

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

	// The job system needs to be updated to use our improved three-pronged architecture
	// For now, we'll initialize basic job components without the complex job implementations

	// Initialize job repositories
	deps.JobRepositories = &jobRepositoriesImpl{
		jobRepo: repository.NewJobRepository(db),
	}

	// Create a simple job service without the complex job implementations
	// This will allow the system to start without errors
	jobService := services.NewJobService(
		deps.JobRepo(),
		deps.UserRepo(),
		deps.UserConfigRepo(),
		nil, // movie repo
		nil, // series repo
		nil, // track repo
		nil, // user movie data repo
		nil, // user series data repo
		nil, // user music data repo
		nil, // recommendation job
		nil, // media sync job
		nil, // watch history sync job
		nil) // favorites sync job

	deps.JobServices = &jobServicesImpl{
		jobService: jobService,
	}

	// Initialize job handlers with just the basic job handler
	deps.JobHandlers = &jobHandlersImpl{
		jobHandler: handlers.NewJobHandler(jobService),
	}

	// No job registration for now, as we're removing legacy code

	deps.UserServices = &userServicesImpl{
		userService: services.NewUserService(deps.UserRepo()),
		userConfigService: services.NewUserConfigService(
			deps.UserConfigRepo(),
			deps.JobServices.JobService(),
			nil, // No recommendation job for now
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

	// Create client-specific media handlers using our three-pronged architecture
	// and the client service for authentication/connection
	clientMovieHandlers := &clientMediaMovieHandlersImpl{
		embyMovieHandler: CreateClientMediaMovieHandler(
			deps.ClientServices.EmbyService(),
			deps.ClientMediaItemServices.MovieClientService()),
		jellyfinMovieHandler: CreateClientMediaMovieHandler(
			deps.ClientServices.JellyfinService(),
			deps.ClientMediaItemServices.MovieClientService()),
		plexMovieHandler: CreateClientMediaMovieHandler(
			deps.ClientServices.PlexService(),
			deps.ClientMediaItemServices.MovieClientService()),
	}

	clientSeriesHandlers := &clientMediaSeriesHandlersImpl{
		embySeriesHandler: CreateClientMediaSeriesHandler(
			deps.ClientServices.EmbyService(),
			deps.ClientMediaItemServices.SeriesClientService()),
		jellyfinSeriesHandler: CreateClientMediaSeriesHandler(
			deps.ClientServices.JellyfinService(),
			deps.ClientMediaItemServices.SeriesClientService()),
		plexSeriesHandler: CreateClientMediaSeriesHandler(
			deps.ClientServices.PlexService(),
			deps.ClientMediaItemServices.SeriesClientService()),
	}

	clientMusicHandlers := &clientMediaMusicHandlersImpl{
		embyMusicHandler: CreateClientMediaMusicHandler(
			deps.ClientServices.EmbyService(),
			deps.ClientMediaItemServices.TrackClientService(),
			deps.ClientMediaItemServices.AlbumClientService(),
			deps.ClientMediaItemServices.ArtistClientService()),
		jellyfinMusicHandler: CreateClientMediaMusicHandler(
			deps.ClientServices.JellyfinService(),
			deps.ClientMediaItemServices.TrackClientService(),
			deps.ClientMediaItemServices.AlbumClientService(),
			deps.ClientMediaItemServices.ArtistClientService()),
		plexMusicHandler: CreateClientMediaMusicHandler(
			deps.ClientServices.PlexService(),
			deps.ClientMediaItemServices.TrackClientService(),
			deps.ClientMediaItemServices.AlbumClientService(),
			deps.ClientMediaItemServices.ArtistClientService()),
		subsonicMusicHandler: CreateClientMediaMusicHandler(
			deps.ClientServices.SubsonicService(),
			deps.ClientMediaItemServices.TrackClientService(),
			deps.ClientMediaItemServices.AlbumClientService(),
			deps.ClientMediaItemServices.ArtistClientService()),
	}

	deps.ClientMediaHandlers = &clientMediaHandlersImpl{
		movieHandlers:  clientMovieHandlers,
		seriesHandlers: clientSeriesHandlers,
		musicHandlers:  clientMusicHandlers,
	}

	// Initialize a basic search handler that doesn't use legacy repositories
	// We'll need to update the search system to work with our improved three-pronged architecture
	searchRepo := repository.NewSearchRepository(db)
	searchService := services.NewSearchService(
		searchRepo,
		deps.CoreMediaItemRepositories.MovieRepo(), // Use core repositories directly
		deps.CoreMediaItemRepositories.SeriesRepo(),
		deps.CoreMediaItemRepositories.EpisodeRepo(),
		deps.CoreMediaItemRepositories.TrackRepo(),
		deps.CoreMediaItemRepositories.AlbumRepo(),
		deps.CoreMediaItemRepositories.ArtistRepo(),
		deps.CoreMediaItemRepositories.CollectionRepo(),
		deps.CoreMediaItemRepositories.PlaylistRepo(),
		repository.NewPersonRepository(db),
		deps.RepositoryCollections.ClientRepositories(),
		client.GetClientFactoryService(),
	)
	deps.SearchHandler = handlers.NewSearchHandler(searchService)

	return deps
}
