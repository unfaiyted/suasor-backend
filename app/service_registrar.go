// app/service_registrar.go
package app

//
// import (
// 	"gorm.io/gorm"
// 	"suasor/client"
// 	mediatypes "suasor/client/media/types"
// 	clienttypes "suasor/client/types"
// 	"suasor/handlers"
// 	"suasor/repository"
// 	"suasor/services"
// )
//
// // ServiceRegistrar is responsible for registering all services in the dependency injection system
// // This provides a single location to view and modify the service registration logic
// type ServiceRegistrar struct {
// 	db            *gorm.DB
// 	clientFactory *client.ClientFactoryService
// 	dependencies  *AppDependencies
// }
//
// // NewServiceRegistrar creates a new service registrar
// func NewServiceRegistrar(db *gorm.DB, clientFactory *client.ClientFactoryService, dependencies *AppDependencies) *ServiceRegistrar {
// 	return &ServiceRegistrar{
// 		db:            db,
// 		clientFactory: clientFactory,
// 		dependencies:  dependencies,
// 	}
// }
//
// // RegisterCoreServices registers all core services
// func (r *ServiceRegistrar) RegisterCoreServices(configService services.ConfigService) {
// 	// Initialize the media data factory if not already initialized
// 	if r.dependencies.MediaDataFactory == nil {
// 		r.dependencies.MediaDataFactory = NewMediaDataFactory(r.db, r.clientFactory)
// 	}
//
// 	// Register repositories
// 	r.RegisterRepositories()
//
// 	// Register system services
// 	r.dependencies.SystemServices = &systemServicesImpl{
// 		healthService: services.NewHealthService(r.db),
// 		configService: configService,
// 	}
// }
//
// // RegisterRepositories registers all repositories
// func (r *ServiceRegistrar) RegisterRepositories() {
// 	// Core repositories
// 	r.dependencies.CoreMediaItemRepositories = r.dependencies.MediaDataFactory.CreateCoreRepositories()
//
// 	// User repositories
// 	r.dependencies.UserRepositoryFactories = r.dependencies.MediaDataFactory.CreateUserRepositories()
//
// 	// Client repositories
// 	r.dependencies.ClientRepositoryFactories = r.dependencies.MediaDataFactory.CreateClientRepositories()
//
// 	// System repositories
// 	r.dependencies.SystemRepositories = &systemRepositoriesImpl{
// 		configRepo: r.dependencies.SystemServices.ConfigService().GetRepo(),
// 	}
//
// 	// User repositories
// 	r.dependencies.UserRepositories = &userRepositoriesImpl{
// 		userRepo:       repository.NewUserRepository(r.db),
// 		userConfigRepo: repository.NewUserConfigRepository(r.db),
// 		sessionRepo:    repository.NewSessionRepository(r.db),
// 	}
//
// 	// Define client repositories directly for the repository collection
// 	// This avoids the need for legacy ClientRepositories
// 	embyRepo := repository.NewClientRepository[*clienttypes.EmbyConfig](r.db)
// 	jellyfinRepo := repository.NewClientRepository[*clienttypes.JellyfinConfig](r.db)
// 	plexRepo := repository.NewClientRepository[*clienttypes.PlexConfig](r.db)
// 	subsonicRepo := repository.NewClientRepository[*clienttypes.SubsonicConfig](r.db)
// 	sonarrRepo := repository.NewClientRepository[*clienttypes.SonarrConfig](r.db)
// 	radarrRepo := repository.NewClientRepository[*clienttypes.RadarrConfig](r.db)
// 	lidarrRepo := repository.NewClientRepository[*clienttypes.LidarrConfig](r.db)
// 	claudeRepo := repository.NewClientRepository[*clienttypes.ClaudeConfig](r.db)
// 	openaiRepo := repository.NewClientRepository[*clienttypes.OpenAIConfig](r.db)
// 	ollamaRepo := repository.NewClientRepository[*clienttypes.OllamaConfig](r.db)
//
// 	// Repository collections
// 	r.dependencies.RepositoryCollections = &repositoryCollectionsImpl{
// 		clientRepos: repository.NewClientRepositoryCollection(
// 			embyRepo,
// 			jellyfinRepo,
// 			plexRepo,
// 			subsonicRepo,
// 			sonarrRepo,
// 			radarrRepo,
// 			lidarrRepo,
// 			claudeRepo,
// 			openaiRepo,
// 			ollamaRepo,
// 		),
// 	}
// }
//
// // RegisterMediaDataServices registers all media data services using the three-pronged approach
// func (r *ServiceRegistrar) RegisterMediaDataServices() {
// 	// Initialize three-pronged services using the factory
// 	r.dependencies.CoreMediaItemServices = r.dependencies.MediaDataFactory.CreateCoreServices(r.dependencies.CoreMediaItemRepositories)
// 	r.dependencies.UserMediaItemServices = r.dependencies.MediaDataFactory.CreateUserServices(r.dependencies.CoreMediaItemServices, r.dependencies.UserRepositoryFactories)
// 	r.dependencies.ClientMediaItemServices = r.dependencies.MediaDataFactory.CreateClientServices(r.dependencies.CoreMediaItemServices, r.dependencies.ClientRepositoryFactories)
//
// 	// Initialize collection services
// 	r.RegisterCollectionServices()
// }
//
// // RegisterCollectionServices registers specialized collection services
// func (r *ServiceRegistrar) RegisterCollectionServices() {
// 	// Create core collection service
// 	collectionCoreService := services.NewCoreCollectionService(
// 		r.dependencies.CoreMediaItemRepositories.CollectionRepo())
//
// 	// Create client collection service (extends core)
// 	collectionClientService := services.NewClientMediaCollectionService(
// 		collectionCoreService,
// 		r.dependencies.ClientRepositoryFactories.CollectionClientRepo(),
// 		nil, // client repo - we'll use nil since it's not directly used in this context
// 		r.dependencies.ClientFactoryService,
// 	)
//
// 	// Create core repository for media data
// 	coreMediaDataRepo := repository.NewMediaItemRepository[mediatypes.MediaData](r.db)
//
// 	// Create user collection service (extends core)
// 	collectionUserService := services.NewUserCollectionService(
// 		collectionCoreService,
// 		r.dependencies.UserRepositoryFactories.CollectionUserRepo(),
// 		coreMediaDataRepo,
// 	)
//
// 	// Create playlist service
// 	playlistExtendedService := services.NewPlaylistService(
// 		r.dependencies.UserRepositoryFactories.PlaylistUserRepo(),
// 		r.dependencies.UserMediaItemServices.PlaylistUserService(),
// 		coreMediaDataRepo,
// 	)
//
// 	// Register specialized collection services
// 	r.dependencies.MediaCollectionServices = r.dependencies.MediaDataFactory.CreateMediaCollectionServices(
// 		r.dependencies.CoreMediaItemServices,
// 		r.dependencies.UserMediaItemServices,
// 		r.dependencies.ClientMediaItemServices,
// 		collectionCoreService,
// 		collectionUserService,
// 		collectionClientService,
// 		playlistExtendedService,
// 	)
// }
//
// // RegisterMediaDataHandlers registers all handlers using the three-pronged approach
// func (r *ServiceRegistrar) RegisterMediaDataHandlers() {
// 	// Initialize core handlers - create directly using the factory
// 	r.dependencies.CoreMediaItemHandlers = r.dependencies.MediaDataFactory.CreateCoreHandlers(
// 		r.dependencies.CoreMediaItemServices)
//
// 	// Initialize CoreUserMediaItemDataServices
// 	r.dependencies.CoreUserMediaItemDataServices = &coreUserMediaItemDataServicesImpl{
// 		movieCoreService: services.NewCoreUserMediaItemDataService[*mediatypes.Movie](
// 			r.dependencies.CoreMediaItemServices.MovieCoreService()),
// 		seriesCoreService: services.NewCoreUserMediaItemDataService[*mediatypes.Series](
// 			r.dependencies.CoreMediaItemServices.SeriesCoreService()),
// 		episodeCoreService: services.NewCoreUserMediaItemDataService[*mediatypes.Episode](
// 			r.dependencies.CoreMediaItemServices.EpisodeCoreService()),
// 		trackCoreService: services.NewCoreUserMediaItemDataService[*mediatypes.Track](
// 			r.dependencies.CoreMediaItemServices.TrackCoreService()),
// 		albumCoreService: services.NewCoreUserMediaItemDataService[*mediatypes.Album](
// 			r.dependencies.CoreMediaItemServices.AlbumCoreService()),
// 		artistCoreService: services.NewCoreUserMediaItemDataService[*mediatypes.Artist](
// 			r.dependencies.CoreMediaItemServices.ArtistCoreService()),
// 		collectionCoreService: services.NewCoreUserMediaItemDataService[*mediatypes.Collection](
// 			r.dependencies.CoreMediaItemServices.CollectionCoreService()),
// 		playlistCoreService: services.NewCoreUserMediaItemDataService[*mediatypes.Playlist](
// 			r.dependencies.CoreMediaItemServices.PlaylistCoreService()),
// 	}
//
// 	// Initialize user data services - these were already initialized in RegisterMediaDataServices
// 	r.dependencies.UserMediaItemDataServices = &userMediaItemDataServicesImpl{
// 		movieDataService: services.NewUserMediaItemDataService[*mediatypes.Movie](
// 			r.dependencies.CoreUserMediaItemDataServices.MovieCoreService(),
// 			r.dependencies.UserDataFactories.MovieDataRepo(),
// 		),
// 		seriesDataService: services.NewUserMediaItemDataService[*mediatypes.Series](
// 			r.dependencies.CoreUserMediaItemDataServices.SeriesCoreService(),
// 			r.dependencies.UserDataFactories.SeriesDataRepo(),
// 		),
// 		episodeDataService: services.NewUserMediaItemDataService[*mediatypes.Episode](
// 			r.dependencies.CoreUserMediaItemDataServices.EpisodeCoreService(),
// 			r.dependencies.UserDataFactories.EpisodeDataRepo(),
// 		),
// 		trackDataService: services.NewUserMediaItemDataService[*mediatypes.Track](
// 			r.dependencies.CoreUserMediaItemDataServices.TrackCoreService(),
// 			r.dependencies.UserDataFactories.TrackDataRepo(),
// 		),
// 		albumDataService: services.NewUserMediaItemDataService[*mediatypes.Album](
// 			r.dependencies.CoreUserMediaItemDataServices.AlbumCoreService(),
// 			r.dependencies.UserDataFactories.AlbumDataRepo(),
// 		),
// 		artistDataService: services.NewUserMediaItemDataService[*mediatypes.Artist](
// 			r.dependencies.CoreUserMediaItemDataServices.ArtistCoreService(),
// 			r.dependencies.UserDataFactories.ArtistDataRepo(),
// 		),
// 		collectionDataService: services.NewUserMediaItemDataService[*mediatypes.Collection](
// 			r.dependencies.CoreUserMediaItemDataServices.CollectionCoreService(),
// 			r.dependencies.UserDataFactories.CollectionDataRepo(),
// 		),
// 		playlistDataService: services.NewUserMediaItemDataService[*mediatypes.Playlist](
// 			r.dependencies.CoreUserMediaItemDataServices.PlaylistCoreService(),
// 			r.dependencies.UserDataFactories.PlaylistDataRepo(),
// 		),
// 	}
//
// 	// Initialize user handlers
// 	r.dependencies.UserMediaItemHandlers = r.dependencies.MediaDataFactory.CreateUserHandlers(
// 		r.dependencies.UserMediaItemServices,
// 		r.dependencies.UserMediaItemDataServices,
// 		r.dependencies.CoreMediaItemHandlers,
// 	)
//
// 	// Initialize client handlers
// 	// Initialize client-user data services
// 	// emptyClientUserDataRepos := &clientUserDataRepositoriesImpl{}
//
// 	r.dependencies.ClientUserMediaItemDataServices = &clientUserMediaItemDataServicesImpl{
// 		movieClientService: services.NewClientUserMediaItemDataService[*mediatypes.Movie](
// 			r.dependencies.UserMediaItemDataServices.MovieDataService(),
// 			nil, // No client data repositories needed
// 		),
// 		seriesClientService: services.NewClientUserMediaItemDataService[*mediatypes.Series](
// 			r.dependencies.UserMediaItemDataServices.SeriesDataService(),
// 			nil,
// 		),
// 		episodeClientService: services.NewClientUserMediaItemDataService[*mediatypes.Episode](
// 			r.dependencies.UserMediaItemDataServices.EpisodeDataService(),
// 			nil,
// 		),
// 		trackClientService: services.NewClientUserMediaItemDataService[*mediatypes.Track](
// 			r.dependencies.UserMediaItemDataServices.TrackDataService(),
// 			nil,
// 		),
// 		albumClientService: services.NewClientUserMediaItemDataService[*mediatypes.Album](
// 			r.dependencies.UserMediaItemDataServices.AlbumDataService(),
// 			nil,
// 		),
// 		artistClientService: services.NewClientUserMediaItemDataService[*mediatypes.Artist](
// 			r.dependencies.UserMediaItemDataServices.ArtistDataService(),
// 			nil,
// 		),
// 		collectionClientService: services.NewClientUserMediaItemDataService[*mediatypes.Collection](
// 			r.dependencies.UserMediaItemDataServices.CollectionDataService(),
// 			nil,
// 		),
// 		playlistClientService: services.NewClientUserMediaItemDataService[*mediatypes.Playlist](
// 			r.dependencies.UserMediaItemDataServices.PlaylistDataService(),
// 			nil,
// 		),
// 	}
//
// 	r.dependencies.ClientMediaItemHandlers = r.dependencies.MediaDataFactory.CreateClientHandlers(
// 		r.dependencies.ClientMediaItemServices,
// 		r.dependencies.ClientUserMediaItemDataServices,
// 		r.dependencies.UserMediaItemHandlers,
// 	)
//
// 	// Initialize specialized handlers
// 	r.RegisterSpecializedHandlers()
// }
//
// // RegisterSpecializedHandlers registers specialized domain-specific handlers
// func (r *ServiceRegistrar) RegisterSpecializedHandlers() {
// 	// Create specialized handlers
//
// 	// Music handler
// 	musicHandler := handlers.NewCoreMusicHandler(
// 		r.dependencies.ClientMediaItemServices.TrackClientService(),
// 		r.dependencies.ClientMediaItemServices.AlbumClientService(),
// 		r.dependencies.ClientMediaItemServices.ArtistClientService(),
// 	)
//
// 	// Series handler - using JellyfinConfig as example
// 	// This is a specialized handler for series that adds specific functionality for Jellyfin
// 	seriesSpecificHandler := &handlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig]{}
//
// 	// Register specialized handlers
// 	r.dependencies.SpecializedMediaHandlers = r.dependencies.MediaDataFactory.CreateSpecializedMediaHandlers(
// 		r.dependencies.CoreMediaItemServices,
// 		r.dependencies.UserMediaItemServices,
// 		r.dependencies.ClientMediaItemServices,
// 		musicHandler,
// 		seriesSpecificHandler,
// 	)
// }
//
// // RegisterStandardHandlers registers standard application handlers
// func (r *ServiceRegistrar) RegisterStandardHandlers() {
// 	// Initialize client repositories directly for client handlers
// 	embyRepo := repository.NewClientRepository[*clienttypes.EmbyConfig](r.db)
// 	jellyfinRepo := repository.NewClientRepository[*clienttypes.JellyfinConfig](r.db)
// 	plexRepo := repository.NewClientRepository[*clienttypes.PlexConfig](r.db)
// 	subsonicRepo := repository.NewClientRepository[*clienttypes.SubsonicConfig](r.db)
// 	sonarrRepo := repository.NewClientRepository[*clienttypes.SonarrConfig](r.db)
// 	radarrRepo := repository.NewClientRepository[*clienttypes.RadarrConfig](r.db)
// 	lidarrRepo := repository.NewClientRepository[*clienttypes.LidarrConfig](r.db)
// 	claudeRepo := repository.NewClientRepository[*clienttypes.ClaudeConfig](r.db)
// 	openaiRepo := repository.NewClientRepository[*clienttypes.OpenAIConfig](r.db)
// 	ollamaRepo := repository.NewClientRepository[*clienttypes.OllamaConfig](r.db)
//
// 	// Initialize client services
// 	embyService := services.NewClientService[*clienttypes.EmbyConfig](r.dependencies.ClientFactoryService, embyRepo)
// 	jellyfinService := services.NewClientService[*clienttypes.JellyfinConfig](r.dependencies.ClientFactoryService, jellyfinRepo)
// 	plexService := services.NewClientService[*clienttypes.PlexConfig](r.dependencies.ClientFactoryService, plexRepo)
// 	subsonicService := services.NewClientService[*clienttypes.SubsonicConfig](r.dependencies.ClientFactoryService, subsonicRepo)
// 	sonarrService := services.NewClientService[*clienttypes.SonarrConfig](r.dependencies.ClientFactoryService, sonarrRepo)
// 	radarrService := services.NewClientService[*clienttypes.RadarrConfig](r.dependencies.ClientFactoryService, radarrRepo)
// 	lidarrService := services.NewClientService[*clienttypes.LidarrConfig](r.dependencies.ClientFactoryService, lidarrRepo)
// 	claudeService := services.NewClientService[*clienttypes.ClaudeConfig](r.dependencies.ClientFactoryService, claudeRepo)
// 	openaiService := services.NewClientService[*clienttypes.OpenAIConfig](r.dependencies.ClientFactoryService, openaiRepo)
// 	ollamaService := services.NewClientService[*clienttypes.OllamaConfig](r.dependencies.ClientFactoryService, ollamaRepo)
//
// 	// Store client services
// 	r.dependencies.ClientServices = &clientServicesImpl{
// 		embyService:     embyService,
// 		jellyfinService: jellyfinService,
// 		plexService:     plexService,
// 		subsonicService: subsonicService,
// 		sonarrService:   sonarrService,
// 		radarrService:   radarrService,
// 		lidarrService:   lidarrService,
// 		claudeService:   claudeService,
// 		openaiService:   openaiService,
// 		ollamaService:   ollamaService,
// 	}
//
// 	// Client handlers
// 	r.dependencies.ClientHandlers = &clientHandlersImpl{
// 		embyHandler:     handlers.NewClientHandler[*clienttypes.EmbyConfig](embyService),
// 		jellyfinHandler: handlers.NewClientHandler[*clienttypes.JellyfinConfig](jellyfinService),
// 		plexHandler:     handlers.NewClientHandler[*clienttypes.PlexConfig](plexService),
// 		subsonicHandler: handlers.NewClientHandler[*clienttypes.SubsonicConfig](subsonicService),
// 		radarrHandler:   handlers.NewClientHandler[*clienttypes.RadarrConfig](radarrService),
// 		lidarrHandler:   handlers.NewClientHandler[*clienttypes.LidarrConfig](lidarrService),
// 		sonarrHandler:   handlers.NewClientHandler[*clienttypes.SonarrConfig](sonarrService),
// 		claudeHandler:   handlers.NewClientHandler[*clienttypes.ClaudeConfig](claudeService),
// 		openaiHandler:   handlers.NewClientHandler[*clienttypes.OpenAIConfig](openaiService),
// 		ollamaHandler:   handlers.NewClientHandler[*clienttypes.OllamaConfig](ollamaService),
// 	}
//
// 	// System handlers
// 	r.dependencies.SystemHandlers = &systemHandlersImpl{
// 		configHandler: handlers.NewConfigHandler(r.dependencies.SystemServices.ConfigService()),
// 		healthHandler: handlers.NewHealthHandler(r.dependencies.SystemServices.HealthService()),
//
// 		clientsHandler: handlers.NewClientsHandler(
// 			embyService,
// 			jellyfinService,
// 			plexService,
// 			subsonicService,
// 			sonarrService,
// 			radarrService,
// 			lidarrService,
// 			claudeService,
// 			openaiService,
// 			ollamaService,
// 		),
// 	}
//
// 	// AI handlers
// 	r.dependencies.AIHandlers = &aiHandlersImpl{
// 		claudeAIHandler: *handlers.NewAIHandler(
// 			r.dependencies.ClientFactoryService,
// 			claudeService,
// 		),
// 		openaiAIHandler: *handlers.NewAIHandler(
// 			r.dependencies.ClientFactoryService,
// 			openaiService,
// 		),
// 		ollamaAIHandler: *handlers.NewAIHandler(
// 			r.dependencies.ClientFactoryService,
// 			ollamaService,
// 		),
// 	}
// }
//
// // RegisterAllServices registers all services in the correct order
// func (r *ServiceRegistrar) RegisterAllServices(configService services.ConfigService) {
// 	// Core services
// 	r.RegisterCoreServices(configService)
//
// 	// Media data services with three-pronged approach
// 	r.RegisterMediaDataServices()
//
// 	// Media data handlers with three-pronged approach
// 	r.RegisterMediaDataHandlers()
//
// 	// Standard handlers
// 	r.RegisterStandardHandlers()
// }

