// app/di/media_factory.go
package di

// import (
// 	"gorm.io/gorm"
// 	"suasor/app/di/factories"
// 	apphandlers "suasor/app/handlers"
// 	"suasor/app/repository"
// 	"suasor/app/services"
// 	"suasor/client"
// 	mediatypes "suasor/client/media/types"
// 	clienttypes "suasor/client/types"
// 	"suasor/handlers"
// 	repo "suasor/repository"
// 	svc "suasor/services"
// )
//
// // mediaDataFactoryImpl is an implementation of the MediaDataFactory interface
// type mediaDataFactoryImpl struct {
// 	db            *gorm.DB
// 	clientFactory *client.ClientFactoryService
// }
//
// // CreateClientDataRepositories is the implementation of factories.MediaDataFactory interface
// func (f *mediaDataFactoryImpl) CreateClientDataRepositories() repository.ClientUserMediaDataRepositories {
// 	// Create repositories using the database
// 	// Get dependencies
// 	coreRepos := f.CreateCoreDataRepositories()
// 	userRepos := f.CreateUserDataRepositories()
//
// 	return &clientUserMediaDataRepositoriesImpl{
// 		movieDataRepo:      repo.NewClientUserMediaItemDataRepository[*mediatypes.Movie](f.db, coreRepos.MovieCoreService(), userRepos.MovieDataRepo()),
// 		seriesDataRepo:     repo.NewClientUserMediaItemDataRepository[*mediatypes.Series](f.db, coreRepos.SeriesCoreService(), userRepos.SeriesDataRepo()),
// 		episodeDataRepo:    repo.NewClientUserMediaItemDataRepository[*mediatypes.Episode](f.db, coreRepos.EpisodeCoreService(), userRepos.EpisodeDataRepo()),
// 		trackDataRepo:      repo.NewClientUserMediaItemDataRepository[*mediatypes.Track](f.db, coreRepos.TrackCoreService(), userRepos.TrackDataRepo()),
// 		albumDataRepo:      repo.NewClientUserMediaItemDataRepository[*mediatypes.Album](f.db, coreRepos.AlbumCoreService(), userRepos.AlbumDataRepo()),
// 		artistDataRepo:     repo.NewClientUserMediaItemDataRepository[*mediatypes.Artist](f.db, coreRepos.ArtistCoreService(), userRepos.ArtistDataRepo()),
// 		collectionDataRepo: repo.NewClientUserMediaItemDataRepository[*mediatypes.Collection](f.db, coreRepos.CollectionCoreService(), userRepos.CollectionDataRepo()),
// 		playlistDataRepo:   repo.NewClientUserMediaItemDataRepository[*mediatypes.Playlist](f.db, coreRepos.PlaylistCoreService(), userRepos.PlaylistDataRepo()),
// 	}
// }
//
// // createMediaDataFactory creates a new MediaDataFactory implementation
// func createMediaDataFactory(db *gorm.DB, clientFactory *client.ClientFactoryService) factories.MediaDataFactory {
// 	return &mediaDataFactoryImpl{
// 		db:            db,
// 		clientFactory: clientFactory,
// 	}
// }
//
// // --------------------------------------------------------
// // Core Repository Factory Methods
// // --------------------------------------------------------
//
// // CreateCoreRepositories initializes all core repositories
// func (f *mediaDataFactoryImpl) CreateCoreRepositories() repository.CoreMediaItemRepositories {
// 	return &coreRepositoriesImpl{
// 		movieRepo:      repo.NewMediaItemRepository[*mediatypes.Movie](f.db),
// 		seriesRepo:     repo.NewMediaItemRepository[*mediatypes.Series](f.db),
// 		seasonRepo:     repo.NewMediaItemRepository[*mediatypes.Season](f.db),
// 		episodeRepo:    repo.NewMediaItemRepository[*mediatypes.Episode](f.db),
// 		trackRepo:      repo.NewMediaItemRepository[*mediatypes.Track](f.db),
// 		albumRepo:      repo.NewMediaItemRepository[*mediatypes.Album](f.db),
// 		artistRepo:     repo.NewMediaItemRepository[*mediatypes.Artist](f.db),
// 		collectionRepo: repo.NewMediaItemRepository[*mediatypes.Collection](f.db),
// 		playlistRepo:   repo.NewMediaItemRepository[*mediatypes.Playlist](f.db),
// 	}
// }
//
// // CreateCoreDataRepositories initializes all core user data repositories
// func (f *mediaDataFactoryImpl) CreateCoreDataRepositories() repository.CoreUserMediaItemDataRepositories {
// 	return &coreCoreUserMediaItemDataRepositoriesImpl{
// 		movieCoreRepo:      repo.NewCoreUserMediaItemDataRepository[*mediatypes.Movie](f.db),
// 		seriesCoreRepo:     repo.NewCoreUserMediaItemDataRepository[*mediatypes.Series](f.db),
// 		episodeCoreRepo:    repo.NewCoreUserMediaItemDataRepository[*mediatypes.Episode](f.db),
// 		trackCoreRepo:      repo.NewCoreUserMediaItemDataRepository[*mediatypes.Track](f.db),
// 		albumCoreRepo:      repo.NewCoreUserMediaItemDataRepository[*mediatypes.Album](f.db),
// 		artistCoreRepo:     repo.NewCoreUserMediaItemDataRepository[*mediatypes.Artist](f.db),
// 		collectionCoreRepo: repo.NewCoreUserMediaItemDataRepository[*mediatypes.Collection](f.db),
// 		playlistCoreRepo:   repo.NewCoreUserMediaItemDataRepository[*mediatypes.Playlist](f.db),
// 	}
// }
//
// // --------------------------------------------------------
// // User Repository Factory Methods
// // --------------------------------------------------------
//
// // CreateUserRepositories initializes all user repositories
// func (f *mediaDataFactoryImpl) CreateUserRepositories() repository.UserMediaItemRepositories {
// 	return &userRepositoryFactoriesImpl{
// 		movieUserRepo:      repo.NewUserMediaItemRepository[*mediatypes.Movie](f.db),
// 		seriesUserRepo:     repo.NewUserMediaItemRepository[*mediatypes.Series](f.db),
// 		seasonUserRepo:     repo.NewUserMediaItemRepository[*mediatypes.Season](f.db),
// 		episodeUserRepo:    repo.NewUserMediaItemRepository[*mediatypes.Episode](f.db),
// 		trackUserRepo:      repo.NewUserMediaItemRepository[*mediatypes.Track](f.db),
// 		albumUserRepo:      repo.NewUserMediaItemRepository[*mediatypes.Album](f.db),
// 		artistUserRepo:     repo.NewUserMediaItemRepository[*mediatypes.Artist](f.db),
// 		collectionUserRepo: repo.NewUserMediaItemRepository[*mediatypes.Collection](f.db),
// 		playlistUserRepo:   repo.NewUserMediaItemRepository[*mediatypes.Playlist](f.db),
// 	}
// }
//
// // --------------------------------------------------------
// // Client Repository Factory Methods
// // --------------------------------------------------------
//
// // CreateClientRepositories initializes all client repositories
// func (f *mediaDataFactoryImpl) CreateClientRepositories() repository.ClientUserMediaDataRepositories {
// 	// Create core repositories if not already available
// 	coreRepos := f.CreateCoreDataRepositories()
// 	userRepos := f.CreateUserDataRepositories()
//
// 	return &clientUserMediaDataRepositoriesImpl{
// 		movieDataRepo:      repo.NewClientUserMediaItemDataRepository[*mediatypes.Movie](f.db, coreRepos.MovieCoreService(), userRepos.MovieDataRepo()),
// 		seriesDataRepo:     repo.NewClientUserMediaItemDataRepository[*mediatypes.Series](f.db, coreRepos.SeriesCoreService(), userRepos.SeriesDataRepo()),
// 		episodeDataRepo:    repo.NewClientUserMediaItemDataRepository[*mediatypes.Episode](f.db, coreRepos.EpisodeCoreService(), userRepos.EpisodeDataRepo()),
// 		trackDataRepo:      repo.NewClientUserMediaItemDataRepository[*mediatypes.Track](f.db, coreRepos.TrackCoreService(), userRepos.TrackDataRepo()),
// 		albumDataRepo:      repo.NewClientUserMediaItemDataRepository[*mediatypes.Album](f.db, coreRepos.AlbumCoreService(), userRepos.AlbumDataRepo()),
// 		artistDataRepo:     repo.NewClientUserMediaItemDataRepository[*mediatypes.Artist](f.db, coreRepos.ArtistCoreService(), userRepos.ArtistDataRepo()),
// 		collectionDataRepo: repo.NewClientUserMediaItemDataRepository[*mediatypes.Collection](f.db, coreRepos.CollectionCoreService(), userRepos.CollectionDataRepo()),
// 		playlistDataRepo:   repo.NewClientUserMediaItemDataRepository[*mediatypes.Playlist](f.db, coreRepos.PlaylistCoreService(), userRepos.PlaylistDataRepo()),
// 	}
// }
//
// // CreateClientMediaItemRepositories initializes all client media item repositories
// func (f *mediaDataFactoryImpl) CreateClientMediaItemRepositories() repository.ClientMediaItemRepositories {
// 	coreRepos := f.CreateCoreRepositories()
//
// 	return &clientMediaItemRepositoriesImpl{
// 		movieClientRepo:      repo.NewClientMediaItemRepository[*mediatypes.Movie](f.db, coreRepos.MovieRepo()),
// 		seriesClientRepo:     repo.NewClientMediaItemRepository[*mediatypes.Series](f.db, coreRepos.SeriesRepo()),
// 		seasonClientRepo:     repo.NewClientMediaItemRepository[*mediatypes.Season](f.db, coreRepos.SeasonRepo()),
// 		episodeClientRepo:    repo.NewClientMediaItemRepository[*mediatypes.Episode](f.db, coreRepos.EpisodeRepo()),
// 		trackClientRepo:      repo.NewClientMediaItemRepository[*mediatypes.Track](f.db, coreRepos.TrackRepo()),
// 		albumClientRepo:      repo.NewClientMediaItemRepository[*mediatypes.Album](f.db, coreRepos.AlbumRepo()),
// 		artistClientRepo:     repo.NewClientMediaItemRepository[*mediatypes.Artist](f.db, coreRepos.ArtistRepo()),
// 		collectionClientRepo: repo.NewClientMediaItemRepository[*mediatypes.Collection](f.db, coreRepos.CollectionRepo()),
// 		playlistClientRepo:   repo.NewClientMediaItemRepository[*mediatypes.Playlist](f.db, coreRepos.PlaylistRepo()),
// 	}
// }
//
// // CreateUserDataRepositories initializes all user data repositories
// func (f *mediaDataFactoryImpl) CreateUserDataRepositories() repository.UserMediaDataRepositories {
// 	coreDataRepos := f.CreateCoreDataRepositories()
//
// 	return &userDataRepositoriesImpl{
// 		movieDataRepo:      repo.NewUserMediaItemDataRepository[*mediatypes.Movie](f.db, coreDataRepos.MovieCoreService()),
// 		seriesDataRepo:     repo.NewUserMediaItemDataRepository[*mediatypes.Series](f.db, coreDataRepos.SeriesCoreService()),
// 		episodeDataRepo:    repo.NewUserMediaItemDataRepository[*mediatypes.Episode](f.db, coreDataRepos.EpisodeCoreService()),
// 		trackDataRepo:      repo.NewUserMediaItemDataRepository[*mediatypes.Track](f.db, coreDataRepos.TrackCoreService()),
// 		albumDataRepo:      repo.NewUserMediaItemDataRepository[*mediatypes.Album](f.db, coreDataRepos.AlbumCoreService()),
// 		artistDataRepo:     repo.NewUserMediaItemDataRepository[*mediatypes.Artist](f.db, coreDataRepos.ArtistCoreService()),
// 		collectionDataRepo: repo.NewUserMediaItemDataRepository[*mediatypes.Collection](f.db, coreDataRepos.CollectionCoreService()),
// 		playlistDataRepo:   repo.NewUserMediaItemDataRepository[*mediatypes.Playlist](f.db, coreDataRepos.PlaylistCoreService()),
// 	}
// }
//
// // --------------------------------------------------------
// // Core Service Factory Methods
// // --------------------------------------------------------
//
// // CreateCoreServices initializes all core services
// func (f *mediaDataFactoryImpl) CreateCoreServices(repos repository.CoreMediaItemRepositories) services.CoreMediaItemServices {
// 	return &coreMediaItemServicesImpl{
// 		movieCoreService:      svc.NewCoreMediaItemService[*mediatypes.Movie](repos.MovieRepo()),
// 		seriesCoreService:     svc.NewCoreMediaItemService[*mediatypes.Series](repos.SeriesRepo()),
// 		episodeCoreService:    svc.NewCoreMediaItemService[*mediatypes.Episode](repos.EpisodeRepo()),
// 		seasonCoreService:     svc.NewCoreMediaItemService[*mediatypes.Season](repos.SeasonRepo()), // Will need to implement Season repository
// 		trackCoreService:      svc.NewCoreMediaItemService[*mediatypes.Track](repos.TrackRepo()),
// 		albumCoreService:      svc.NewCoreMediaItemService[*mediatypes.Album](repos.AlbumRepo()),
// 		artistCoreService:     svc.NewCoreMediaItemService[*mediatypes.Artist](repos.ArtistRepo()),
// 		collectionCoreService: svc.NewCoreMediaItemService[*mediatypes.Collection](repos.CollectionRepo()),
// 		playlistCoreService:   svc.NewCoreMediaItemService[*mediatypes.Playlist](repos.PlaylistRepo()),
// 	}
// }
//
// // Updated implementation of CreateCoreDataRepositories
//
// // CreateCoreDataServices initializes all core data services
// func (f *mediaDataFactoryImpl) CreateCoreDataServices(repos repository.CoreMediaItemRepositories) services.CoreUserMediaItemDataServices {
// 	// Get core data repositories
// 	coreRepos := f.CreateCoreDataRepositories()
//
// 	return &coreUserMediaItemDataServicesImpl{
// 		movieCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Movie](
// 			svc.NewCoreMediaItemService[*mediatypes.Movie](repos.MovieRepo()),
// 			coreRepos.MovieCoreService()),
// 		seriesCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Series](
// 			svc.NewCoreMediaItemService[*mediatypes.Series](repos.SeriesRepo()),
// 			coreRepos.SeriesCoreService()),
// 		episodeCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Episode](
// 			svc.NewCoreMediaItemService[*mediatypes.Episode](repos.EpisodeRepo()),
// 			coreRepos.EpisodeCoreService()),
// 		trackCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Track](
// 			svc.NewCoreMediaItemService[*mediatypes.Track](repos.TrackRepo()),
// 			coreRepos.TrackCoreService()),
// 		albumCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Album](
// 			svc.NewCoreMediaItemService[*mediatypes.Album](repos.AlbumRepo()),
// 			coreRepos.AlbumCoreService()),
// 		artistCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Artist](
// 			svc.NewCoreMediaItemService[*mediatypes.Artist](repos.ArtistRepo()),
// 			coreRepos.ArtistCoreService()),
// 		collectionCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Collection](
// 			svc.NewCoreMediaItemService[*mediatypes.Collection](repos.CollectionRepo()),
// 			coreRepos.CollectionCoreService()),
// 		playlistCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Playlist](
// 			svc.NewCoreMediaItemService[*mediatypes.Playlist](repos.PlaylistRepo()),
// 			coreRepos.PlaylistCoreService()),
// 	}
// }
//
// // --------------------------------------------------------
// // List Service Factory Methods
// // --------------------------------------------------------
//
// // Create a wrapper type that properly combines UserMediaItemHandler with CoreMediaItemHandler
// type userMediaItemHandlerWithCore[T mediatypes.MediaData] struct {
// 	handlers.UserMediaItemHandler[T]
// 	handlers.CoreMediaItemHandler[T]
// }
//
// // Define core user media item data repositories implementation
// type coreUserMediaItemDataRepositoriesImpl struct {
// 	movieCoreService      repo.CoreUserMediaItemDataRepository[*mediatypes.Movie]
// 	seriesCoreService     repo.CoreUserMediaItemDataRepository[*mediatypes.Series]
// 	episodeCoreService    repo.CoreUserMediaItemDataRepository[*mediatypes.Episode]
// 	trackCoreService      repo.CoreUserMediaItemDataRepository[*mediatypes.Track]
// 	albumCoreService      repo.CoreUserMediaItemDataRepository[*mediatypes.Album]
// 	artistCoreService     repo.CoreUserMediaItemDataRepository[*mediatypes.Artist]
// 	collectionCoreService repo.CoreUserMediaItemDataRepository[*mediatypes.Collection]
// 	playlistCoreService   repo.CoreUserMediaItemDataRepository[*mediatypes.Playlist]
// }
//
// func (r *coreUserMediaItemDataRepositoriesImpl) MovieCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Movie] {
// 	return r.movieCoreService
// }
//
// func (r *coreUserMediaItemDataRepositoriesImpl) SeriesCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Series] {
// 	return r.seriesCoreService
// }
//
// func (r *coreUserMediaItemDataRepositoriesImpl) EpisodeCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Episode] {
// 	return r.episodeCoreService
// }
//
// func (r *coreUserMediaItemDataRepositoriesImpl) TrackCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Track] {
// 	return r.trackCoreService
// }
//
// func (r *coreUserMediaItemDataRepositoriesImpl) AlbumCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Album] {
// 	return r.albumCoreService
// }
//
// func (r *coreUserMediaItemDataRepositoriesImpl) ArtistCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Artist] {
// 	return r.artistCoreService
// }
//
// func (r *coreUserMediaItemDataRepositoriesImpl) CollectionCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Collection] {
// 	return r.collectionCoreService
// }
//
// func (r *coreUserMediaItemDataRepositoriesImpl) PlaylistCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Playlist] {
// 	return r.playlistCoreService
// }
//
// // Define core list services implementation
// type coreListServicesImpl struct {
// 	coreCollectionService svc.CoreListService[*mediatypes.Collection]
// 	corePlaylistService   svc.CoreListService[*mediatypes.Playlist]
// }
//
// func (s *coreListServicesImpl) CoreCollectionService() svc.CoreListService[*mediatypes.Collection] {
// 	return s.coreCollectionService
// }
//
// func (s *coreListServicesImpl) CorePlaylistService() svc.CoreListService[*mediatypes.Playlist] {
// 	return s.corePlaylistService
// }
//
// // CreateCoreListServices initializes core list services
// func (f *mediaDataFactoryImpl) CreateCoreListServices(coreServices services.CoreMediaItemServices) services.CoreListServices {
// 	return &coreListServicesImpl{
// 		// Temporary using nil until we have proper constructors
// 		coreCollectionService: nil,
// 		corePlaylistService:   nil,
// 	}
// }
//
// // CreateUserListServices initializes user list services
// func (f *mediaDataFactoryImpl) CreateUserListServices(
// 	userServices services.UserMediaItemServices,
// 	coreListServices services.CoreListServices) services.UserListServices {
//
// 	// Placeholder implementation - will be updated with proper service initialization
// 	return nil
// }
//
// // CreateClientListServices initializes client list services
// func (f *mediaDataFactoryImpl) CreateClientListServices(
// 	clientServices services.ClientMediaItemServices[clienttypes.ClientMediaConfig],
// 	coreListServices services.CoreListServices) services.ClientListServices {
//
// 	// Return empty implementation for now, will need to update properly
// 	return &clientListServicesImpl{
// 		// These will be properly initialized with per-client implementations
// 	}
// }
//
// // --------------------------------------------------------
// // MediaItem Handler Factory Methods
// // --------------------------------------------------------
//
// // CreateCoreMediaItemHandlers initializes all core media item handlers
// func (f *mediaDataFactoryImpl) CreateCoreMediaItemHandlers(
// 	coreServices services.CoreMediaItemServices) apphandlers.CoreMediaItemHandlers {
//
// 	return &coreMediaItemHandlersImpl{
// 		movieCoreHandler: handlers.NewCoreMediaItemHandler[*mediatypes.Movie](
// 			coreServices.MovieCoreService()),
// 		seriesCoreHandler: handlers.NewCoreMediaItemHandler[*mediatypes.Series](
// 			coreServices.SeriesCoreService()),
// 		episodeCoreHandler: handlers.NewCoreMediaItemHandler[*mediatypes.Episode](
// 			coreServices.EpisodeCoreService()),
// 		trackCoreHandler: handlers.NewCoreMediaItemHandler[*mediatypes.Track](
// 			coreServices.TrackCoreService()),
// 		albumCoreHandler: handlers.NewCoreMediaItemHandler[*mediatypes.Album](
// 			coreServices.AlbumCoreService()),
// 		artistCoreHandler: handlers.NewCoreMediaItemHandler[*mediatypes.Artist](
// 			coreServices.ArtistCoreService()),
// 		collectionCoreHandler: handlers.NewCoreMediaItemHandler[*mediatypes.Collection](
// 			coreServices.CollectionCoreService()),
// 		playlistCoreHandler: handlers.NewCoreMediaItemHandler[*mediatypes.Playlist](
// 			coreServices.PlaylistCoreService()),
// 	}
// }
//
// // CreateUserMediaItemHandlers initializes all user media item handlers
// func (f *mediaDataFactoryImpl) CreateUserMediaItemHandlers(
// 	userServices services.UserMediaItemServices,
// 	coreHandlers apphandlers.CoreMediaItemHandlers) apphandlers.UserMediaItemHandlers {
//
// 	// For now, we'll create handlers directly without trying to combine them with CoreHandlers
// 	// This is a simpler approach until we have a better implementation for embedding the core handlers
// 	return &userMediaItemHandlersImpl{
// 		movieUserHandler:      handlers.NewUserMediaItemHandler[*mediatypes.Movie](userServices.MovieUserService()),
// 		seriesUserHandler:     handlers.NewUserMediaItemHandler[*mediatypes.Series](userServices.SeriesUserService()),
// 		episodeUserHandler:    handlers.NewUserMediaItemHandler[*mediatypes.Episode](userServices.EpisodeUserService()),
// 		trackUserHandler:      handlers.NewUserMediaItemHandler[*mediatypes.Track](userServices.TrackUserService()),
// 		albumUserHandler:      handlers.NewUserMediaItemHandler[*mediatypes.Album](userServices.AlbumUserService()),
// 		artistUserHandler:     handlers.NewUserMediaItemHandler[*mediatypes.Artist](userServices.ArtistUserService()),
// 		collectionUserHandler: handlers.NewUserMediaItemHandler[*mediatypes.Collection](userServices.CollectionUserService()),
// 		playlistUserHandler:   handlers.NewUserMediaItemHandler[*mediatypes.Playlist](userServices.PlaylistUserService()),
// 	}
// }
//
// // CreateClientMediaItemHandlers initializes all client media item handlers
// func (f *mediaDataFactoryImpl) CreateClientMediaItemHandlers(
// 	clientServices services.ClientMediaItemServices[clienttypes.ClientMediaConfig],
// 	userServices services.UserMediaItemServices,
// 	userHandlers apphandlers.UserMediaItemHandlers,
// ) apphandlers.ClientMediaItemHandlers[clienttypes.ClientMediaConfig] {
//
// 	return &clientMediaItemHandlersImpl{
// 		movieClientHandler: handlers.NewClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Movie](
// 			userHandlers.MovieUserHandler(),
// 			clientServices.MovieClientService(),
// 		),
// 		seriesClientHandler: handlers.NewClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Series](
// 			userHandlers.SeriesUserHandler(),
// 			clientServices.SeriesClientService(),
// 		),
// 		episodeClientHandler: handlers.NewClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Episode](
// 			userHandlers.EpisodeUserHandler(),
// 			clientServices.EpisodeClientService(),
// 		),
// 		trackClientHandler: handlers.NewClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Track](
// 			userHandlers.TrackUserHandler(),
// 			clientServices.TrackClientService(),
// 		),
// 		albumClientHandler: handlers.NewClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Album](
// 			userHandlers.AlbumUserHandler(),
// 			clientServices.AlbumClientService(),
// 		),
// 		artistClientHandler: handlers.NewClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Artist](
// 			userHandlers.ArtistUserHandler(),
// 			clientServices.ArtistClientService(),
// 		),
// 		collectionClientHandler: handlers.NewClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Collection](
// 			userHandlers.CollectionUserHandler(),
// 			clientServices.CollectionClientService(),
// 		),
// 		playlistClientHandler: handlers.NewClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Playlist](
// 			userHandlers.PlaylistUserHandler(),
// 			clientServices.PlaylistClientService(),
// 		),
// 	}
// }
//
// // --------------------------------------------------------
// // User Service Factory Methods
// // --------------------------------------------------------
//
// // CreateUserServices initializes all user services
// func (f *mediaDataFactoryImpl) CreateUserServices(
// 	coreServices services.CoreMediaItemServices,
// 	userRepos repository.UserMediaItemRepositories) services.UserMediaItemServices {
//
// 	return &userMediaItemServicesImpl{
// 		movieUserService: svc.NewUserMediaItemService[*mediatypes.Movie](
// 			coreServices.MovieCoreService(), userRepos.MovieUserRepo()),
// 		seriesUserService: svc.NewUserMediaItemService[*mediatypes.Series](
// 			coreServices.SeriesCoreService(), userRepos.SeriesUserRepo()),
// 		episodeUserService: svc.NewUserMediaItemService[*mediatypes.Episode](
// 			coreServices.EpisodeCoreService(), userRepos.EpisodeUserRepo()),
// 		trackUserService: svc.NewUserMediaItemService[*mediatypes.Track](
// 			coreServices.TrackCoreService(), userRepos.TrackUserRepo()),
// 		albumUserService: svc.NewUserMediaItemService[*mediatypes.Album](
// 			coreServices.AlbumCoreService(), userRepos.AlbumUserRepo()),
// 		artistUserService: svc.NewUserMediaItemService[*mediatypes.Artist](
// 			coreServices.ArtistCoreService(), userRepos.ArtistUserRepo()),
// 		collectionUserService: svc.NewUserMediaItemService[*mediatypes.Collection](
// 			coreServices.CollectionCoreService(), userRepos.CollectionUserRepo()),
// 		playlistUserService: svc.NewUserMediaItemService[*mediatypes.Playlist](
// 			coreServices.PlaylistCoreService(), userRepos.PlaylistUserRepo()),
// 	}
// }
//
// // CreateUserDataServices initializes all user data services
// func (f *mediaDataFactoryImpl) CreateUserDataServices(
// 	coreDataServices services.CoreUserMediaItemDataServices,
// 	userRepos repository.UserMediaDataRepositories) services.UserMediaItemDataServices {
//
// 	return &userMediaItemDataServicesImpl{
// 		movieDataService: svc.NewUserMediaItemDataService[*mediatypes.Movie](
// 			coreDataServices.MovieCoreService(), userRepos.MovieDataRepo()),
// 		seriesDataService: svc.NewUserMediaItemDataService[*mediatypes.Series](
// 			coreDataServices.SeriesCoreService(), userRepos.SeriesDataRepo()),
// 		episodeDataService: svc.NewUserMediaItemDataService[*mediatypes.Episode](
// 			coreDataServices.EpisodeCoreService(), userRepos.EpisodeDataRepo()),
// 		trackDataService: svc.NewUserMediaItemDataService[*mediatypes.Track](
// 			coreDataServices.TrackCoreService(), userRepos.TrackDataRepo()),
// 		albumDataService: svc.NewUserMediaItemDataService[*mediatypes.Album](
// 			coreDataServices.AlbumCoreService(), userRepos.AlbumDataRepo()),
// 		artistDataService: svc.NewUserMediaItemDataService[*mediatypes.Artist](
// 			coreDataServices.ArtistCoreService(), userRepos.ArtistDataRepo()),
// 		collectionDataService: svc.NewUserMediaItemDataService[*mediatypes.Collection](
// 			coreDataServices.CollectionCoreService(), userRepos.CollectionDataRepo()),
// 		playlistDataService: svc.NewUserMediaItemDataService[*mediatypes.Playlist](
// 			coreDataServices.PlaylistCoreService(), userRepos.PlaylistDataRepo()),
// 	}
// }
//
// // --------------------------------------------------------
// // Client Service Factory Methods
// // --------------------------------------------------------
//
// // CreateClientServices initializes all client services
// func (f *mediaDataFactoryImpl) CreateClientServices(
// 	coreServices services.CoreMediaItemServices,
// 	clientRepos repo.ClientRepository[clienttypes.ClientMediaConfig],
// 	itemRepos repository.ClientMediaItemRepositories,
// ) services.ClientMediaItemServices[clienttypes.ClientMediaConfig] {
// 	return &clientMediaItemServicesImpl{
// 		movieClientService: svc.NewClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Movie](
// 			coreServices.MovieCoreService(), clientRepos, itemRepos.MovieClientRepo()),
// 		seriesClientService: svc.NewClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Series](
// 			coreServices.SeriesCoreService(), clientRepos, itemRepos.SeriesClientRepo()),
// 		episodeClientService: svc.NewClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Episode](
// 			coreServices.EpisodeCoreService(), clientRepos, itemRepos.EpisodeClientRepo()),
// 		trackClientService: svc.NewClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Track](
// 			coreServices.TrackCoreService(), clientRepos, itemRepos.TrackClientRepo()),
// 		albumClientService: svc.NewClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Album](
// 			coreServices.AlbumCoreService(), clientRepos, itemRepos.AlbumClientRepo()),
// 		artistClientService: svc.NewClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Artist](
// 			coreServices.ArtistCoreService(), clientRepos, itemRepos.ArtistClientRepo()),
// 		collectionClientService: svc.NewClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Collection](
// 			coreServices.CollectionCoreService(), clientRepos, itemRepos.CollectionClientRepo()),
// 		playlistClientService: svc.NewClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Playlist](
// 			coreServices.PlaylistCoreService(), clientRepos, itemRepos.PlaylistClientRepo()),
// 	}
// }
//
// // CreateClientDataServices initializes all client data services
// func (f *mediaDataFactoryImpl) CreateClientDataServices(
// 	userDataServices services.UserMediaItemDataServices,
// 	clientRepos repository.ClientUserMediaDataRepositories) services.ClientUserMediaItemDataServices {
//
// 	return &clientUserMediaItemDataServicesImpl{
// 		movieDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Movie](
// 			userDataServices.MovieDataService(), clientRepos.MovieDataRepo()),
// 		seriesDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Series](
// 			userDataServices.SeriesDataService(), clientRepos.SeriesDataRepo()),
// 		episodeDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Episode](
// 			userDataServices.EpisodeDataService(), clientRepos.EpisodeDataRepo()),
// 		trackDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Track](
// 			userDataServices.TrackDataService(), clientRepos.TrackDataRepo()),
// 		albumDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Album](
// 			userDataServices.AlbumDataService(), clientRepos.AlbumDataRepo()),
// 		artistDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Artist](
// 			userDataServices.ArtistDataService(), clientRepos.ArtistDataRepo()),
// 		collectionDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Collection](
// 			userDataServices.CollectionDataService(), clientRepos.CollectionDataRepo()),
// 		playlistDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Playlist](
// 			userDataServices.PlaylistDataService(), clientRepos.PlaylistDataRepo()),
// 	}
// }
//
// // --------------------------------------------------------
// // Specialized Collection Services
// // --------------------------------------------------------
//
// // CreateMediaCollectionServices creates collection and playlist services
// func (f *mediaDataFactoryImpl) CreateMediaCollectionServices(
// 	coreServices services.CoreMediaItemServices,
// 	userServices services.UserMediaItemServices,
// 	clientServices services.ClientMediaItemServices[clienttypes.ClientMediaConfig],
// 	coreCollectionService svc.CoreListService[*mediatypes.Collection],
//
// 	userCollectionService svc.UserListService[*mediatypes.Collection],
// 	clientCollectionService services.ClientListServices,
// 	playlistService svc.CoreListService[*mediatypes.Playlist]) interface{} {
//
// 	return &mediaCollectionServicesImpl{
// 		coreCollectionService:   coreCollectionService,
// 		userCollectionService:   userCollectionService,
// 		clientCollectionService: clientCollectionService,
//
// 		corePlaylistService:   coreServices.PlaylistCoreService(),
// 		userPlaylistService:   userServices.PlaylistUserService(),
// 		clientPlaylistService: clientServices.PlaylistClientService(),
//
// 		playlistService: playlistService,
// 	}
// }
//
// // --------------------------------------------------------
// // Core Handler Factory Methods
// // --------------------------------------------------------
//
// // CreateCoreDataHandlers initializes all core handlers
// func (f *mediaDataFactoryImpl) CreateCoreDataHandlers(
// 	coreServices services.CoreUserMediaItemDataServices) apphandlers.CoreMediaItemDataHandlers {
//
// 	return &coreMediaItemDataHandlersImpl{
// 		movieCoreDataHandler:      handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Movie](coreServices.MovieCoreService()),
// 		seriesCoreDataHandler:     handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Series](coreServices.SeriesCoreService()),
// 		episodeCoreDataHandler:    handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Episode](coreServices.EpisodeCoreService()),
// 		trackCoreDataHandler:      handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Track](coreServices.TrackCoreService()),
// 		albumCoreDataHandler:      handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Album](coreServices.AlbumCoreService()),
// 		artistCoreDataHandler:     handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Artist](coreServices.ArtistCoreService()),
// 		collectionCoreDataHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Collection](coreServices.CollectionCoreService()),
// 		playlistCoreDataHandler:   handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Playlist](coreServices.PlaylistCoreService()),
// 	}
// }
//
// // --------------------------------------------------------
// // User Handler Factory Methods
// // --------------------------------------------------------
//
// // CreateUserDataHandlers initializes all user handlers
// func (f *mediaDataFactoryImpl) CreateUserDataHandlers(
// 	userServices services.UserMediaItemDataServices,
// 	coreHandlers apphandlers.CoreMediaItemDataHandlers) apphandlers.UserMediaItemDataHandlers {
//
// 	return &userMediaItemDataHandlersImpl{
// 		movieUserDataHandler: handlers.NewUserMediaItemDataHandler[*mediatypes.Movie](
// 			coreHandlers.MovieCoreDataHandler(),
// 			userServices.MovieDataService()),
// 		seriesUserDataHandler: handlers.NewUserMediaItemDataHandler[*mediatypes.Series](
// 			coreHandlers.SeriesCoreDataHandler(),
// 			userServices.SeriesDataService()),
// 		episodeUserDataHandler: handlers.NewUserMediaItemDataHandler[*mediatypes.Episode](
// 			coreHandlers.EpisodeCoreDataHandler(),
// 			userServices.EpisodeDataService()),
// 		trackUserDataHandler: handlers.NewUserMediaItemDataHandler[*mediatypes.Track](
// 			coreHandlers.TrackCoreDataHandler(),
// 			userServices.TrackDataService()),
// 		albumUserDataHandler: handlers.NewUserMediaItemDataHandler[*mediatypes.Album](
// 			coreHandlers.AlbumCoreDataHandler(),
// 			userServices.AlbumDataService()),
// 		artistUserDataHandler: handlers.NewUserMediaItemDataHandler[*mediatypes.Artist](
// 			coreHandlers.ArtistCoreDataHandler(),
// 			userServices.ArtistDataService()),
// 		collectionUserDataHandler: handlers.NewUserMediaItemDataHandler[*mediatypes.Collection](
// 			coreHandlers.CollectionCoreDataHandler(),
// 			userServices.CollectionDataService()),
// 		playlistUserDataHandler: handlers.NewUserMediaItemDataHandler[*mediatypes.Playlist](
// 			coreHandlers.PlaylistCoreDataHandler(),
// 			userServices.PlaylistDataService()),
// 	}
// }
//
// // --------------------------------------------------------
// // Client Handler Factory Methods
// // --------------------------------------------------------
//
// // CreateClientDataHandlers initializes all client handlers
// func (f *mediaDataFactoryImpl) CreateClientDataHandlers(
// 	dataServices services.ClientUserMediaItemDataServices,
// 	userHandlers apphandlers.UserMediaItemDataHandlers) apphandlers.ClientMediaItemDataHandlers {
//
// 	return &clientMediaItemDataHandlersImpl{
// 		movieClientDataHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Movie](
// 			userHandlers.MovieUserDataHandler(),
// 			dataServices.MovieDataService()),
// 		seriesClientDataHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Series](
// 			userHandlers.SeriesUserDataHandler(),
// 			dataServices.SeriesDataService()),
// 		episodeClientDataHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Episode](
// 			userHandlers.EpisodeUserDataHandler(),
// 			dataServices.EpisodeDataService()),
// 		trackClientDataHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Track](
// 			userHandlers.TrackUserDataHandler(),
// 			dataServices.TrackDataService()),
// 		albumClientDataHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Album](
// 			userHandlers.AlbumUserDataHandler(),
// 			dataServices.AlbumDataService()),
// 		artistClientDataHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Artist](
// 			userHandlers.ArtistUserDataHandler(),
// 			dataServices.ArtistDataService()),
// 		collectionClientDataHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Collection](
// 			userHandlers.CollectionUserDataHandler(),
// 			dataServices.CollectionDataService()),
// 		playlistClientDataHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Playlist](
// 			userHandlers.PlaylistUserDataHandler(),
// 			dataServices.PlaylistDataService()),
// 	}
// }
//
// // --------------------------------------------------------
// // Specialized Media Handlers
// // --------------------------------------------------------
//
// // CreateSpecializedMediaHandlers creates specialized handlers for specific domains
// func (f *mediaDataFactoryImpl) CreateSpecializedMediaHandlers(
// 	coreServices services.CoreMediaItemServices,
// 	userServices services.UserMediaItemServices,
// 	clientServices services.ClientMediaItemServices[clienttypes.ClientMediaConfig]) *specializedMediaHandlersImpl {
//
// 	return &specializedMediaHandlersImpl{
// 		// These will be initialized elsewhere
// 	}
// }
//
// // Define specializedMediaHandlersImpl
// type specializedMediaHandlersImpl struct {
// 	// Fields will be defined based on specialized handler requirements
// }
//
// // Define mediaCollectionServicesImpl
// type mediaCollectionServicesImpl struct {
// 	coreCollectionService   svc.CoreListService[*mediatypes.Collection]
// 	userCollectionService   svc.UserListService[*mediatypes.Collection]
// 	clientCollectionService services.ClientListServices
//
// 	corePlaylistService   svc.CoreMediaItemService[*mediatypes.Playlist]
// 	userPlaylistService   svc.UserMediaItemService[*mediatypes.Playlist]
// 	clientPlaylistService svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Playlist]
//
// 	playlistService svc.CoreListService[*mediatypes.Playlist]
// }
//
// // --------------------------------------------------------
// // Implementation structs
// // --------------------------------------------------------
//
// type userListServicesImpl struct {
// 	userCollectionService svc.UserListService[*mediatypes.Collection]
// 	userPlaylistService   svc.UserListService[*mediatypes.Playlist]
// }
//
// func (s *userListServicesImpl) UserCollectionService() svc.UserListService[*mediatypes.Collection] {
// 	return s.userCollectionService
// }
//
// func (s *userListServicesImpl) UserPlaylistService() svc.UserListService[*mediatypes.Playlist] {
// 	return s.userPlaylistService
// }
//
// type clientListServicesImpl struct {
// 	embyClientCollectionService     svc.ClientListService[*clienttypes.EmbyConfig, *mediatypes.Collection]
// 	embyClientPlaylistService       svc.ClientListService[*clienttypes.EmbyConfig, *mediatypes.Playlist]
// 	jellyfinClientCollectionService svc.ClientListService[*clienttypes.JellyfinConfig, *mediatypes.Collection]
// 	jellyfinClientPlaylistService   svc.ClientListService[*clienttypes.JellyfinConfig, *mediatypes.Playlist]
// 	plexClientCollectionService     svc.ClientListService[*clienttypes.PlexConfig, *mediatypes.Collection]
// 	plexClientPlaylistService       svc.ClientListService[*clienttypes.PlexConfig, *mediatypes.Playlist]
// 	subsonicClientCollectionService svc.ClientListService[*clienttypes.SubsonicConfig, *mediatypes.Collection]
// 	subsonicClientPlaylistService   svc.ClientListService[*clienttypes.SubsonicConfig, *mediatypes.Playlist]
// }
//
// func (s *clientListServicesImpl) EmbyClientCollectionService() svc.ClientListService[*clienttypes.EmbyConfig, *mediatypes.Collection] {
// 	return s.embyClientCollectionService
// }
//
// func (s *clientListServicesImpl) EmbyClientPlaylistService() svc.ClientListService[*clienttypes.EmbyConfig, *mediatypes.Playlist] {
// 	return s.embyClientPlaylistService
// }
//
// func (s *clientListServicesImpl) JellyfinClientCollectionService() svc.ClientListService[*clienttypes.JellyfinConfig, *mediatypes.Collection] {
// 	return s.jellyfinClientCollectionService
// }
//
// func (s *clientListServicesImpl) JellyfinClientPlaylistService() svc.ClientListService[*clienttypes.JellyfinConfig, *mediatypes.Playlist] {
// 	return s.jellyfinClientPlaylistService
// }
//
// func (s *clientListServicesImpl) PlexClientCollectionService() svc.ClientListService[*clienttypes.PlexConfig, *mediatypes.Collection] {
// 	return s.plexClientCollectionService
// }
//
// func (s *clientListServicesImpl) PlexClientPlaylistService() svc.ClientListService[*clienttypes.PlexConfig, *mediatypes.Playlist] {
// 	return s.plexClientPlaylistService
// }
//
// func (s *clientListServicesImpl) SubsonicClientCollectionService() svc.ClientListService[*clienttypes.SubsonicConfig, *mediatypes.Collection] {
// 	return s.subsonicClientCollectionService
// }
//
// func (s *clientListServicesImpl) SubsonicClientPlaylistService() svc.ClientListService[*clienttypes.SubsonicConfig, *mediatypes.Playlist] {
// 	return s.subsonicClientPlaylistService
// }
//
// // ClientUserMediaDataRepositories implementation
// type clientUserMediaDataRepositoriesImpl struct {
// 	movieDataRepo      repo.ClientUserMediaItemDataRepository[*mediatypes.Movie]
// 	seriesDataRepo     repo.ClientUserMediaItemDataRepository[*mediatypes.Series]
// 	episodeDataRepo    repo.ClientUserMediaItemDataRepository[*mediatypes.Episode]
// 	trackDataRepo      repo.ClientUserMediaItemDataRepository[*mediatypes.Track]
// 	albumDataRepo      repo.ClientUserMediaItemDataRepository[*mediatypes.Album]
// 	artistDataRepo     repo.ClientUserMediaItemDataRepository[*mediatypes.Artist]
// 	collectionDataRepo repo.ClientUserMediaItemDataRepository[*mediatypes.Collection]
// 	playlistDataRepo   repo.ClientUserMediaItemDataRepository[*mediatypes.Playlist]
// }
//
// func (r *clientUserMediaDataRepositoriesImpl) MovieDataRepo() repo.ClientUserMediaItemDataRepository[*mediatypes.Movie] {
// 	return r.movieDataRepo
// }
//
// func (r *clientUserMediaDataRepositoriesImpl) SeriesDataRepo() repo.ClientUserMediaItemDataRepository[*mediatypes.Series] {
// 	return r.seriesDataRepo
// }
//
// func (r *clientUserMediaDataRepositoriesImpl) EpisodeDataRepo() repo.ClientUserMediaItemDataRepository[*mediatypes.Episode] {
// 	return r.episodeDataRepo
// }
//
// func (r *clientUserMediaDataRepositoriesImpl) TrackDataRepo() repo.ClientUserMediaItemDataRepository[*mediatypes.Track] {
// 	return r.trackDataRepo
// }
//
// func (r *clientUserMediaDataRepositoriesImpl) AlbumDataRepo() repo.ClientUserMediaItemDataRepository[*mediatypes.Album] {
// 	return r.albumDataRepo
// }
//
// func (r *clientUserMediaDataRepositoriesImpl) ArtistDataRepo() repo.ClientUserMediaItemDataRepository[*mediatypes.Artist] {
// 	return r.artistDataRepo
// }
//
// func (r *clientUserMediaDataRepositoriesImpl) CollectionDataRepo() repo.ClientUserMediaItemDataRepository[*mediatypes.Collection] {
// 	return r.collectionDataRepo
// }
//
// func (r *clientUserMediaDataRepositoriesImpl) PlaylistDataRepo() repo.ClientUserMediaItemDataRepository[*mediatypes.Playlist] {
// 	return r.playlistDataRepo
// }
//
// // Implementation for the ClientMediaItemRepositories
// type clientMediaItemRepositoriesImpl struct {
// 	movieClientRepo      repo.ClientMediaItemRepository[*mediatypes.Movie]
// 	seriesClientRepo     repo.ClientMediaItemRepository[*mediatypes.Series]
// 	seasonClientRepo     repo.ClientMediaItemRepository[*mediatypes.Season]
// 	episodeClientRepo    repo.ClientMediaItemRepository[*mediatypes.Episode]
// 	trackClientRepo      repo.ClientMediaItemRepository[*mediatypes.Track]
// 	albumClientRepo      repo.ClientMediaItemRepository[*mediatypes.Album]
// 	artistClientRepo     repo.ClientMediaItemRepository[*mediatypes.Artist]
// 	collectionClientRepo repo.ClientMediaItemRepository[*mediatypes.Collection]
// 	playlistClientRepo   repo.ClientMediaItemRepository[*mediatypes.Playlist]
// }
//
// func (r *clientMediaItemRepositoriesImpl) MovieClientRepo() repo.ClientMediaItemRepository[*mediatypes.Movie] {
// 	return r.movieClientRepo
// }
//
// func (r *clientMediaItemRepositoriesImpl) SeasonClientRepo() repo.ClientMediaItemRepository[*mediatypes.Season] {
// 	return r.seasonClientRepo
// }
//
// func (r *clientMediaItemRepositoriesImpl) SeriesClientRepo() repo.ClientMediaItemRepository[*mediatypes.Series] {
// 	return r.seriesClientRepo
// }
//
// func (r *clientMediaItemRepositoriesImpl) EpisodeClientRepo() repo.ClientMediaItemRepository[*mediatypes.Episode] {
// 	return r.episodeClientRepo
// }
//
// func (r *clientMediaItemRepositoriesImpl) TrackClientRepo() repo.ClientMediaItemRepository[*mediatypes.Track] {
// 	return r.trackClientRepo
// }
//
// func (r *clientMediaItemRepositoriesImpl) AlbumClientRepo() repo.ClientMediaItemRepository[*mediatypes.Album] {
// 	return r.albumClientRepo
// }
//
// func (r *clientMediaItemRepositoriesImpl) ArtistClientRepo() repo.ClientMediaItemRepository[*mediatypes.Artist] {
// 	return r.artistClientRepo
// }
//
// func (r *clientMediaItemRepositoriesImpl) CollectionClientRepo() repo.ClientMediaItemRepository[*mediatypes.Collection] {
// 	return r.collectionClientRepo
// }
//
// func (r *clientMediaItemRepositoriesImpl) PlaylistClientRepo() repo.ClientMediaItemRepository[*mediatypes.Playlist] {
// 	return r.playlistClientRepo
// }
//
// // Repository implementation structs
//
// type coreRepositoriesImpl struct {
// 	movieRepo      repo.MediaItemRepository[*mediatypes.Movie]
// 	seriesRepo     repo.MediaItemRepository[*mediatypes.Series]
// 	seasonRepo     repo.MediaItemRepository[*mediatypes.Season]
// 	episodeRepo    repo.MediaItemRepository[*mediatypes.Episode]
// 	trackRepo      repo.MediaItemRepository[*mediatypes.Track]
// 	albumRepo      repo.MediaItemRepository[*mediatypes.Album]
// 	artistRepo     repo.MediaItemRepository[*mediatypes.Artist]
// 	collectionRepo repo.MediaItemRepository[*mediatypes.Collection]
// 	playlistRepo   repo.MediaItemRepository[*mediatypes.Playlist]
// }
//
// func (r *coreRepositoriesImpl) MovieRepo() repo.MediaItemRepository[*mediatypes.Movie] {
// 	return r.movieRepo
// }
//
// func (r *coreRepositoriesImpl) SeriesRepo() repo.MediaItemRepository[*mediatypes.Series] {
// 	return r.seriesRepo
// }
//
// func (r *coreRepositoriesImpl) SeasonRepo() repo.MediaItemRepository[*mediatypes.Season] {
// 	return r.seasonRepo
// }
//
// func (r *coreRepositoriesImpl) EpisodeRepo() repo.MediaItemRepository[*mediatypes.Episode] {
// 	return r.episodeRepo
// }
//
// func (r *coreRepositoriesImpl) TrackRepo() repo.MediaItemRepository[*mediatypes.Track] {
// 	return r.trackRepo
// }
//
// func (r *coreRepositoriesImpl) AlbumRepo() repo.MediaItemRepository[*mediatypes.Album] {
// 	return r.albumRepo
// }
//
// func (r *coreRepositoriesImpl) ArtistRepo() repo.MediaItemRepository[*mediatypes.Artist] {
// 	return r.artistRepo
// }
//
// func (r *coreRepositoriesImpl) CollectionRepo() repo.MediaItemRepository[*mediatypes.Collection] {
// 	return r.collectionRepo
// }
//
// func (r *coreRepositoriesImpl) PlaylistRepo() repo.MediaItemRepository[*mediatypes.Playlist] {
// 	return r.playlistRepo
// }
//
// type userRepositoryFactoriesImpl struct {
// 	movieUserRepo      repo.UserMediaItemRepository[*mediatypes.Movie]
// 	seriesUserRepo     repo.UserMediaItemRepository[*mediatypes.Series]
// 	seasonUserRepo     repo.UserMediaItemRepository[*mediatypes.Season]
// 	episodeUserRepo    repo.UserMediaItemRepository[*mediatypes.Episode]
// 	trackUserRepo      repo.UserMediaItemRepository[*mediatypes.Track]
// 	albumUserRepo      repo.UserMediaItemRepository[*mediatypes.Album]
// 	artistUserRepo     repo.UserMediaItemRepository[*mediatypes.Artist]
// 	collectionUserRepo repo.UserMediaItemRepository[*mediatypes.Collection]
// 	playlistUserRepo   repo.UserMediaItemRepository[*mediatypes.Playlist]
// }
//
// func (r *userRepositoryFactoriesImpl) MovieUserRepo() repo.UserMediaItemRepository[*mediatypes.Movie] {
// 	return r.movieUserRepo
// }
//
// func (r *userRepositoryFactoriesImpl) SeasonUserRepo() repo.UserMediaItemRepository[*mediatypes.Season] {
// 	return r.seasonUserRepo
// }
//
// func (r *userRepositoryFactoriesImpl) SeriesUserRepo() repo.UserMediaItemRepository[*mediatypes.Series] {
// 	return r.seriesUserRepo
// }
//
// func (r *userRepositoryFactoriesImpl) EpisodeUserRepo() repo.UserMediaItemRepository[*mediatypes.Episode] {
// 	return r.episodeUserRepo
// }
//
// func (r *userRepositoryFactoriesImpl) TrackUserRepo() repo.UserMediaItemRepository[*mediatypes.Track] {
// 	return r.trackUserRepo
// }
//
// func (r *userRepositoryFactoriesImpl) AlbumUserRepo() repo.UserMediaItemRepository[*mediatypes.Album] {
// 	return r.albumUserRepo
// }
//
// func (r *userRepositoryFactoriesImpl) ArtistUserRepo() repo.UserMediaItemRepository[*mediatypes.Artist] {
// 	return r.artistUserRepo
// }
//
// func (r *userRepositoryFactoriesImpl) CollectionUserRepo() repo.UserMediaItemRepository[*mediatypes.Collection] {
// 	return r.collectionUserRepo
// }
//
// func (r *userRepositoryFactoriesImpl) PlaylistUserRepo() repo.UserMediaItemRepository[*mediatypes.Playlist] {
// 	return r.playlistUserRepo
// }
//
// type clientRepositoryFactoriesImpl struct {
// 	movieClientRepo      repo.ClientMediaItemRepository[*mediatypes.Movie]
// 	seriesClientRepo     repo.ClientMediaItemRepository[*mediatypes.Series]
// 	episodeClientRepo    repo.ClientMediaItemRepository[*mediatypes.Episode]
// 	trackClientRepo      repo.ClientMediaItemRepository[*mediatypes.Track]
// 	albumClientRepo      repo.ClientMediaItemRepository[*mediatypes.Album]
// 	artistClientRepo     repo.ClientMediaItemRepository[*mediatypes.Artist]
// 	collectionClientRepo repo.ClientMediaItemRepository[*mediatypes.Collection]
// 	playlistClientRepo   repo.ClientMediaItemRepository[*mediatypes.Playlist]
// }
//
// func (r *clientRepositoryFactoriesImpl) MovieClientRepo() repo.ClientMediaItemRepository[*mediatypes.Movie] {
// 	return r.movieClientRepo
// }
//
// func (r *clientRepositoryFactoriesImpl) SeriesClientRepo() repo.ClientMediaItemRepository[*mediatypes.Series] {
// 	return r.seriesClientRepo
// }
//
// func (r *clientRepositoryFactoriesImpl) EpisodeClientRepo() repo.ClientMediaItemRepository[*mediatypes.Episode] {
// 	return r.episodeClientRepo
// }
//
// func (r *clientRepositoryFactoriesImpl) TrackClientRepo() repo.ClientMediaItemRepository[*mediatypes.Track] {
// 	return r.trackClientRepo
// }
//
// func (r *clientRepositoryFactoriesImpl) AlbumClientRepo() repo.ClientMediaItemRepository[*mediatypes.Album] {
// 	return r.albumClientRepo
// }
//
// func (r *clientRepositoryFactoriesImpl) ArtistClientRepo() repo.ClientMediaItemRepository[*mediatypes.Artist] {
// 	return r.artistClientRepo
// }
//
// func (r *clientRepositoryFactoriesImpl) CollectionClientRepo() repo.ClientMediaItemRepository[*mediatypes.Collection] {
// 	return r.collectionClientRepo
// }
//
// func (r *clientRepositoryFactoriesImpl) PlaylistClientRepo() repo.ClientMediaItemRepository[*mediatypes.Playlist] {
// 	return r.playlistClientRepo
// }
//
// type userDataRepositoriesImpl struct {
// 	movieDataRepo      repo.UserMediaItemDataRepository[*mediatypes.Movie]
// 	seriesDataRepo     repo.UserMediaItemDataRepository[*mediatypes.Series]
// 	episodeDataRepo    repo.UserMediaItemDataRepository[*mediatypes.Episode]
// 	trackDataRepo      repo.UserMediaItemDataRepository[*mediatypes.Track]
// 	albumDataRepo      repo.UserMediaItemDataRepository[*mediatypes.Album]
// 	artistDataRepo     repo.UserMediaItemDataRepository[*mediatypes.Artist]
// 	collectionDataRepo repo.UserMediaItemDataRepository[*mediatypes.Collection]
// 	playlistDataRepo   repo.UserMediaItemDataRepository[*mediatypes.Playlist]
// }
//
// func (r *userDataRepositoriesImpl) MovieDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Movie] {
// 	return r.movieDataRepo
// }
//
// func (r *userDataRepositoriesImpl) SeriesDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Series] {
// 	return r.seriesDataRepo
// }
//
// func (r *userDataRepositoriesImpl) EpisodeDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Episode] {
// 	return r.episodeDataRepo
// }
//
// func (r *userDataRepositoriesImpl) TrackDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Track] {
// 	return r.trackDataRepo
// }
//
// func (r *userDataRepositoriesImpl) AlbumDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Album] {
// 	return r.albumDataRepo
// }
//
// func (r *userDataRepositoriesImpl) ArtistDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Artist] {
// 	return r.artistDataRepo
// }
//
// func (r *userDataRepositoriesImpl) CollectionDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Collection] {
// 	return r.collectionDataRepo
// }
//
// func (r *userDataRepositoriesImpl) PlaylistDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Playlist] {
// 	return r.playlistDataRepo
// }
//
// // Service implementation structs
//
// type coreMediaItemServicesImpl struct {
// 	movieCoreService      svc.CoreMediaItemService[*mediatypes.Movie]
// 	seriesCoreService     svc.CoreMediaItemService[*mediatypes.Series]
// 	episodeCoreService    svc.CoreMediaItemService[*mediatypes.Episode]
// 	seasonCoreService     svc.CoreMediaItemService[*mediatypes.Season]
// 	trackCoreService      svc.CoreMediaItemService[*mediatypes.Track]
// 	albumCoreService      svc.CoreMediaItemService[*mediatypes.Album]
// 	artistCoreService     svc.CoreMediaItemService[*mediatypes.Artist]
// 	collectionCoreService svc.CoreMediaItemService[*mediatypes.Collection]
// 	playlistCoreService   svc.CoreMediaItemService[*mediatypes.Playlist]
// }
//
// func (s *coreMediaItemServicesImpl) MovieCoreService() svc.CoreMediaItemService[*mediatypes.Movie] {
// 	return s.movieCoreService
// }
//
// func (s *coreMediaItemServicesImpl) SeriesCoreService() svc.CoreMediaItemService[*mediatypes.Series] {
// 	return s.seriesCoreService
// }
//
// func (s *coreMediaItemServicesImpl) EpisodeCoreService() svc.CoreMediaItemService[*mediatypes.Episode] {
// 	return s.episodeCoreService
// }
//
// func (s *coreMediaItemServicesImpl) SeasonCoreService() svc.CoreMediaItemService[*mediatypes.Season] {
// 	return s.seasonCoreService
// }
//
// func (s *coreMediaItemServicesImpl) TrackCoreService() svc.CoreMediaItemService[*mediatypes.Track] {
// 	return s.trackCoreService
// }
//
// func (s *coreMediaItemServicesImpl) AlbumCoreService() svc.CoreMediaItemService[*mediatypes.Album] {
// 	return s.albumCoreService
// }
//
// func (s *coreMediaItemServicesImpl) ArtistCoreService() svc.CoreMediaItemService[*mediatypes.Artist] {
// 	return s.artistCoreService
// }
//
// func (s *coreMediaItemServicesImpl) CollectionCoreService() svc.CoreMediaItemService[*mediatypes.Collection] {
// 	return s.collectionCoreService
// }
//
// func (s *coreMediaItemServicesImpl) PlaylistCoreService() svc.CoreMediaItemService[*mediatypes.Playlist] {
// 	return s.playlistCoreService
// }
//
// type coreUserMediaItemDataServicesImpl struct {
// 	movieCoreService      svc.CoreUserMediaItemDataService[*mediatypes.Movie]
// 	seriesCoreService     svc.CoreUserMediaItemDataService[*mediatypes.Series]
// 	episodeCoreService    svc.CoreUserMediaItemDataService[*mediatypes.Episode]
// 	trackCoreService      svc.CoreUserMediaItemDataService[*mediatypes.Track]
// 	albumCoreService      svc.CoreUserMediaItemDataService[*mediatypes.Album]
// 	artistCoreService     svc.CoreUserMediaItemDataService[*mediatypes.Artist]
// 	collectionCoreService svc.CoreUserMediaItemDataService[*mediatypes.Collection]
// 	playlistCoreService   svc.CoreUserMediaItemDataService[*mediatypes.Playlist]
// }
//
// func (s *coreUserMediaItemDataServicesImpl) MovieCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Movie] {
// 	return s.movieCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) SeriesCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Series] {
// 	return s.seriesCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) EpisodeCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Episode] {
// 	return s.episodeCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) TrackCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Track] {
// 	return s.trackCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) AlbumCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Album] {
// 	return s.albumCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) ArtistCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Artist] {
// 	return s.artistCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) CollectionCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Collection] {
// 	return s.collectionCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) PlaylistCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Playlist] {
// 	return s.playlistCoreService
// }
//
// type userMediaItemServicesImpl struct {
// 	movieUserService      svc.UserMediaItemService[*mediatypes.Movie]
// 	seriesUserService     svc.UserMediaItemService[*mediatypes.Series]
// 	episodeUserService    svc.UserMediaItemService[*mediatypes.Episode]
// 	seasonUserService     svc.UserMediaItemService[*mediatypes.Season]
// 	trackUserService      svc.UserMediaItemService[*mediatypes.Track]
// 	albumUserService      svc.UserMediaItemService[*mediatypes.Album]
// 	artistUserService     svc.UserMediaItemService[*mediatypes.Artist]
// 	collectionUserService svc.UserMediaItemService[*mediatypes.Collection]
// 	playlistUserService   svc.UserMediaItemService[*mediatypes.Playlist]
// }
//
// func (s *userMediaItemServicesImpl) MovieUserService() svc.UserMediaItemService[*mediatypes.Movie] {
// 	return s.movieUserService
// }
//
// func (s *userMediaItemServicesImpl) SeriesUserService() svc.UserMediaItemService[*mediatypes.Series] {
// 	return s.seriesUserService
// }
//
// func (s *userMediaItemServicesImpl) EpisodeUserService() svc.UserMediaItemService[*mediatypes.Episode] {
// 	return s.episodeUserService
// }
//
// func (s *userMediaItemServicesImpl) TrackUserService() svc.UserMediaItemService[*mediatypes.Track] {
// 	return s.trackUserService
// }
//
// func (s *userMediaItemServicesImpl) AlbumUserService() svc.UserMediaItemService[*mediatypes.Album] {
// 	return s.albumUserService
// }
//
// func (s *userMediaItemServicesImpl) ArtistUserService() svc.UserMediaItemService[*mediatypes.Artist] {
// 	return s.artistUserService
// }
//
// func (s *userMediaItemServicesImpl) CollectionUserService() svc.UserMediaItemService[*mediatypes.Collection] {
// 	return s.collectionUserService
// }
//
// func (s *userMediaItemServicesImpl) PlaylistUserService() svc.UserMediaItemService[*mediatypes.Playlist] {
// 	return s.playlistUserService
// }
//
// func (s *userMediaItemServicesImpl) SeasonUserService() svc.UserMediaItemService[*mediatypes.Season] {
// 	return s.seasonUserService
// }
//
// type userMediaItemDataServicesImpl struct {
// 	movieDataService      svc.UserMediaItemDataService[*mediatypes.Movie]
// 	seriesDataService     svc.UserMediaItemDataService[*mediatypes.Series]
// 	episodeDataService    svc.UserMediaItemDataService[*mediatypes.Episode]
// 	trackDataService      svc.UserMediaItemDataService[*mediatypes.Track]
// 	albumDataService      svc.UserMediaItemDataService[*mediatypes.Album]
// 	artistDataService     svc.UserMediaItemDataService[*mediatypes.Artist]
// 	collectionDataService svc.UserMediaItemDataService[*mediatypes.Collection]
// 	playlistDataService   svc.UserMediaItemDataService[*mediatypes.Playlist]
// }
//
// func (s *userMediaItemDataServicesImpl) MovieDataService() svc.UserMediaItemDataService[*mediatypes.Movie] {
// 	return s.movieDataService
// }
//
// func (s *userMediaItemDataServicesImpl) SeriesDataService() svc.UserMediaItemDataService[*mediatypes.Series] {
// 	return s.seriesDataService
// }
//
// func (s *userMediaItemDataServicesImpl) EpisodeDataService() svc.UserMediaItemDataService[*mediatypes.Episode] {
// 	return s.episodeDataService
// }
//
// func (s *userMediaItemDataServicesImpl) TrackDataService() svc.UserMediaItemDataService[*mediatypes.Track] {
// 	return s.trackDataService
// }
//
// func (s *userMediaItemDataServicesImpl) AlbumDataService() svc.UserMediaItemDataService[*mediatypes.Album] {
// 	return s.albumDataService
// }
//
// func (s *userMediaItemDataServicesImpl) ArtistDataService() svc.UserMediaItemDataService[*mediatypes.Artist] {
// 	return s.artistDataService
// }
//
// func (s *userMediaItemDataServicesImpl) CollectionDataService() svc.UserMediaItemDataService[*mediatypes.Collection] {
// 	return s.collectionDataService
// }
//
// func (s *userMediaItemDataServicesImpl) PlaylistDataService() svc.UserMediaItemDataService[*mediatypes.Playlist] {
// 	return s.playlistDataService
// }
//
// type clientMediaItemServicesImpl struct {
// 	movieClientService      svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Movie]
// 	seriesClientService     svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Series]
// 	episodeClientService    svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Episode]
// 	seasonClientService     svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Season]
// 	trackClientService      svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Track]
// 	albumClientService      svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Album]
// 	artistClientService     svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Artist]
// 	collectionClientService svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Collection]
// 	playlistClientService   svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Playlist]
// }
//
// func (s *clientMediaItemServicesImpl) MovieClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Movie] {
// 	return s.movieClientService
// }
//
// func (s *clientMediaItemServicesImpl) SeriesClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Series] {
// 	return s.seriesClientService
// }
//
// func (s *clientMediaItemServicesImpl) EpisodeClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Episode] {
// 	return s.episodeClientService
// }
//
// func (s *clientMediaItemServicesImpl) TrackClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Track] {
// 	return s.trackClientService
// }
//
// func (s *clientMediaItemServicesImpl) AlbumClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Album] {
// 	return s.albumClientService
// }
//
// func (s *clientMediaItemServicesImpl) ArtistClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Artist] {
// 	return s.artistClientService
// }
//
// func (s *clientMediaItemServicesImpl) CollectionClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Collection] {
// 	return s.collectionClientService
// }
//
// func (s *clientMediaItemServicesImpl) PlaylistClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Playlist] {
// 	return s.playlistClientService
// }
//
// func (s *clientMediaItemServicesImpl) SeasonClientService() svc.ClientMediaItemService[clienttypes.ClientMediaConfig, *mediatypes.Season] {
// 	return nil // TODO: Implement properly with a real service
// }
//
// type clientUserMediaItemDataServicesImpl struct {
// 	movieDataService      svc.ClientUserMediaItemDataService[*mediatypes.Movie]
// 	seriesDataService     svc.ClientUserMediaItemDataService[*mediatypes.Series]
// 	episodeDataService    svc.ClientUserMediaItemDataService[*mediatypes.Episode]
// 	trackDataService      svc.ClientUserMediaItemDataService[*mediatypes.Track]
// 	albumDataService      svc.ClientUserMediaItemDataService[*mediatypes.Album]
// 	artistDataService     svc.ClientUserMediaItemDataService[*mediatypes.Artist]
// 	collectionDataService svc.ClientUserMediaItemDataService[*mediatypes.Collection]
// 	playlistDataService   svc.ClientUserMediaItemDataService[*mediatypes.Playlist]
// }
//
// func (s *clientUserMediaItemDataServicesImpl) MovieDataService() svc.ClientUserMediaItemDataService[*mediatypes.Movie] {
// 	return s.movieDataService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) SeriesDataService() svc.ClientUserMediaItemDataService[*mediatypes.Series] {
// 	return s.seriesDataService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) EpisodeDataService() svc.ClientUserMediaItemDataService[*mediatypes.Episode] {
// 	return s.episodeDataService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) TrackDataService() svc.ClientUserMediaItemDataService[*mediatypes.Track] {
// 	return s.trackDataService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) AlbumDataService() svc.ClientUserMediaItemDataService[*mediatypes.Album] {
// 	return s.albumDataService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) ArtistDataService() svc.ClientUserMediaItemDataService[*mediatypes.Artist] {
// 	return s.artistDataService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) CollectionDataService() svc.ClientUserMediaItemDataService[*mediatypes.Collection] {
// 	return s.collectionDataService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) PlaylistDataService() svc.ClientUserMediaItemDataService[*mediatypes.Playlist] {
// 	return s.playlistDataService
// }
//
// // MediaItem Handlers implementation structs
// type coreMediaItemHandlersImpl struct {
// 	movieCoreHandler      handlers.CoreMediaItemHandler[*mediatypes.Movie]
// 	seriesCoreHandler     handlers.CoreMediaItemHandler[*mediatypes.Series]
// 	episodeCoreHandler    handlers.CoreMediaItemHandler[*mediatypes.Episode]
// 	trackCoreHandler      handlers.CoreMediaItemHandler[*mediatypes.Track]
// 	albumCoreHandler      handlers.CoreMediaItemHandler[*mediatypes.Album]
// 	artistCoreHandler     handlers.CoreMediaItemHandler[*mediatypes.Artist]
// 	collectionCoreHandler handlers.CoreMediaItemHandler[*mediatypes.Collection]
// 	playlistCoreHandler   handlers.CoreMediaItemHandler[*mediatypes.Playlist]
// }
//
// func (h *coreMediaItemHandlersImpl) MovieCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Movie] {
// 	return h.movieCoreHandler
// }
//
// func (h *coreMediaItemHandlersImpl) SeriesCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Series] {
// 	return h.seriesCoreHandler
// }
//
// func (h *coreMediaItemHandlersImpl) EpisodeCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Episode] {
// 	return h.episodeCoreHandler
// }
//
// func (h *coreMediaItemHandlersImpl) TrackCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Track] {
// 	return h.trackCoreHandler
// }
//
// func (h *coreMediaItemHandlersImpl) AlbumCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Album] {
// 	return h.albumCoreHandler
// }
//
// func (h *coreMediaItemHandlersImpl) ArtistCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Artist] {
// 	return h.artistCoreHandler
// }
//
// func (h *coreMediaItemHandlersImpl) CollectionCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Collection] {
// 	return h.collectionCoreHandler
// }
//
// func (h *coreMediaItemHandlersImpl) PlaylistCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Playlist] {
// 	return h.playlistCoreHandler
// }
//
// type userMediaItemHandlersImpl struct {
// 	movieUserHandler      handlers.UserMediaItemHandler[*mediatypes.Movie]
// 	seriesUserHandler     handlers.UserMediaItemHandler[*mediatypes.Series]
// 	episodeUserHandler    handlers.UserMediaItemHandler[*mediatypes.Episode]
// 	trackUserHandler      handlers.UserMediaItemHandler[*mediatypes.Track]
// 	albumUserHandler      handlers.UserMediaItemHandler[*mediatypes.Album]
// 	artistUserHandler     handlers.UserMediaItemHandler[*mediatypes.Artist]
// 	collectionUserHandler handlers.UserMediaItemHandler[*mediatypes.Collection]
// 	playlistUserHandler   handlers.UserMediaItemHandler[*mediatypes.Playlist]
// }
//
// func (h *userMediaItemHandlersImpl) MovieUserHandler() handlers.UserMediaItemHandler[*mediatypes.Movie] {
// 	return h.movieUserHandler
// }
//
// func (h *userMediaItemHandlersImpl) SeriesUserHandler() handlers.UserMediaItemHandler[*mediatypes.Series] {
// 	return h.seriesUserHandler
// }
//
// func (h *userMediaItemHandlersImpl) EpisodeUserHandler() handlers.UserMediaItemHandler[*mediatypes.Episode] {
// 	return h.episodeUserHandler
// }
//
// func (h *userMediaItemHandlersImpl) TrackUserHandler() handlers.UserMediaItemHandler[*mediatypes.Track] {
// 	return h.trackUserHandler
// }
//
// func (h *userMediaItemHandlersImpl) AlbumUserHandler() handlers.UserMediaItemHandler[*mediatypes.Album] {
// 	return h.albumUserHandler
// }
//
// func (h *userMediaItemHandlersImpl) ArtistUserHandler() handlers.UserMediaItemHandler[*mediatypes.Artist] {
// 	return h.artistUserHandler
// }
//
// func (h *userMediaItemHandlersImpl) CollectionUserHandler() handlers.UserMediaItemHandler[*mediatypes.Collection] {
// 	return h.collectionUserHandler
// }
//
// func (h *userMediaItemHandlersImpl) PlaylistUserHandler() handlers.UserMediaItemHandler[*mediatypes.Playlist] {
// 	return h.playlistUserHandler
// }
//
// type clientMediaItemHandlersImpl struct {
// 	movieClientHandler      handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Movie]
// 	seriesClientHandler     handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Series]
// 	episodeClientHandler    handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Episode]
// 	trackClientHandler      handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Track]
// 	albumClientHandler      handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Album]
// 	artistClientHandler     handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Artist]
// 	collectionClientHandler handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Collection]
// 	playlistClientHandler   handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Playlist]
// }
//
// func (h *clientMediaItemHandlersImpl) MovieClientHandler() handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Movie] {
// 	return h.movieClientHandler
// }
//
// func (h *clientMediaItemHandlersImpl) SeriesClientHandler() handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Series] {
// 	return h.seriesClientHandler
// }
//
// func (h *clientMediaItemHandlersImpl) EpisodeClientHandler() handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Episode] {
// 	return h.episodeClientHandler
// }
//
// func (h *clientMediaItemHandlersImpl) TrackClientHandler() handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Track] {
// 	return h.trackClientHandler
// }
//
// func (h *clientMediaItemHandlersImpl) AlbumClientHandler() handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Album] {
// 	return h.albumClientHandler
// }
//
// func (h *clientMediaItemHandlersImpl) ArtistClientHandler() handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Artist] {
// 	return h.artistClientHandler
// }
//
// func (h *clientMediaItemHandlersImpl) CollectionClientHandler() handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Collection] {
// 	return h.collectionClientHandler
// }
//
// func (h *clientMediaItemHandlersImpl) PlaylistClientHandler() handlers.ClientMediaItemHandler[clienttypes.ClientMediaConfig, *mediatypes.Playlist] {
// 	return h.playlistClientHandler
// }
//
// // MediaItemData Handlers implementation structs
// type coreMediaItemDataHandlersImpl struct {
// 	movieCoreDataHandler      handlers.CoreUserMediaItemDataHandler[*mediatypes.Movie]
// 	seriesCoreDataHandler     handlers.CoreUserMediaItemDataHandler[*mediatypes.Series]
// 	episodeCoreDataHandler    handlers.CoreUserMediaItemDataHandler[*mediatypes.Episode]
// 	trackCoreDataHandler      handlers.CoreUserMediaItemDataHandler[*mediatypes.Track]
// 	albumCoreDataHandler      handlers.CoreUserMediaItemDataHandler[*mediatypes.Album]
// 	artistCoreDataHandler     handlers.CoreUserMediaItemDataHandler[*mediatypes.Artist]
// 	collectionCoreDataHandler handlers.CoreUserMediaItemDataHandler[*mediatypes.Collection]
// 	playlistCoreDataHandler   handlers.CoreUserMediaItemDataHandler[*mediatypes.Playlist]
// }
//
// func (h *coreMediaItemDataHandlersImpl) MovieCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Movie] {
// 	return h.movieCoreDataHandler
// }
//
// func (h *coreMediaItemDataHandlersImpl) SeriesCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Series] {
// 	return h.seriesCoreDataHandler
// }
//
// func (h *coreMediaItemDataHandlersImpl) EpisodeCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Episode] {
// 	return h.episodeCoreDataHandler
// }
//
// func (h *coreMediaItemDataHandlersImpl) TrackCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Track] {
// 	return h.trackCoreDataHandler
// }
//
// func (h *coreMediaItemDataHandlersImpl) AlbumCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Album] {
// 	return h.albumCoreDataHandler
// }
//
// func (h *coreMediaItemDataHandlersImpl) ArtistCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Artist] {
// 	return h.artistCoreDataHandler
// }
//
// func (h *coreMediaItemDataHandlersImpl) CollectionCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Collection] {
// 	return h.collectionCoreDataHandler
// }
//
// func (h *coreMediaItemDataHandlersImpl) PlaylistCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Playlist] {
// 	return h.playlistCoreDataHandler
// }
//
// type userMediaItemDataHandlersImpl struct {
// 	movieUserDataHandler      handlers.UserMediaItemDataHandler[*mediatypes.Movie]
// 	seriesUserDataHandler     handlers.UserMediaItemDataHandler[*mediatypes.Series]
// 	episodeUserDataHandler    handlers.UserMediaItemDataHandler[*mediatypes.Episode]
// 	seasonUserDataHandler     handlers.UserMediaItemDataHandler[*mediatypes.Season]
// 	trackUserDataHandler      handlers.UserMediaItemDataHandler[*mediatypes.Track]
// 	albumUserDataHandler      handlers.UserMediaItemDataHandler[*mediatypes.Album]
// 	artistUserDataHandler     handlers.UserMediaItemDataHandler[*mediatypes.Artist]
// 	collectionUserDataHandler handlers.UserMediaItemDataHandler[*mediatypes.Collection]
// 	playlistUserDataHandler   handlers.UserMediaItemDataHandler[*mediatypes.Playlist]
// }
//
// func (h *userMediaItemDataHandlersImpl) MovieUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Movie] {
// 	return h.movieUserDataHandler
// }
//
// func (h *userMediaItemDataHandlersImpl) SeriesUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Series] {
// 	return h.seriesUserDataHandler
// }
//
// func (h *userMediaItemDataHandlersImpl) EpisodeUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Episode] {
// 	return h.episodeUserDataHandler
// }
//
// func (h *userMediaItemDataHandlersImpl) TrackUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Track] {
// 	return h.trackUserDataHandler
// }
//
// func (h *userMediaItemDataHandlersImpl) AlbumUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Album] {
// 	return h.albumUserDataHandler
// }
//
// func (h *userMediaItemDataHandlersImpl) ArtistUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Artist] {
// 	return h.artistUserDataHandler
// }
//
// func (h *userMediaItemDataHandlersImpl) CollectionUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Collection] {
// 	return h.collectionUserDataHandler
// }
//
// func (h *userMediaItemDataHandlersImpl) PlaylistUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Playlist] {
// 	return h.playlistUserDataHandler
// }
//
// func (h *userMediaItemDataHandlersImpl) SeasonUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Season] {
// 	return h.seasonUserDataHandler
// }
//
// type clientMediaItemDataHandlersImpl struct {
// 	movieClientDataHandler      handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie]
// 	seriesClientDataHandler     handlers.ClientUserMediaItemDataHandler[*mediatypes.Series]
// 	episodeClientDataHandler    handlers.ClientUserMediaItemDataHandler[*mediatypes.Episode]
// 	seasonClientDataHandler     handlers.ClientUserMediaItemDataHandler[*mediatypes.Season]
// 	trackClientDataHandler      handlers.ClientUserMediaItemDataHandler[*mediatypes.Track]
// 	albumClientDataHandler      handlers.ClientUserMediaItemDataHandler[*mediatypes.Album]
// 	artistClientDataHandler     handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist]
// 	collectionClientDataHandler handlers.ClientUserMediaItemDataHandler[*mediatypes.Collection]
// 	playlistClientDataHandler   handlers.ClientUserMediaItemDataHandler[*mediatypes.Playlist]
// }
//
// func (h *clientMediaItemDataHandlersImpl) MovieClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie] {
// 	return h.movieClientDataHandler
// }
//
// func (h *clientMediaItemDataHandlersImpl) SeriesClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Series] {
// 	return h.seriesClientDataHandler
// }
//
// func (h *clientMediaItemDataHandlersImpl) EpisodeClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Episode] {
// 	return h.episodeClientDataHandler
// }
//
// func (h *clientMediaItemDataHandlersImpl) TrackClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Track] {
// 	return h.trackClientDataHandler
// }
//
// func (h *clientMediaItemDataHandlersImpl) AlbumClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Album] {
// 	return h.albumClientDataHandler
// }
//
// func (h *clientMediaItemDataHandlersImpl) ArtistClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist] {
// 	return h.artistClientDataHandler
// }
//
// func (h *clientMediaItemDataHandlersImpl) CollectionClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Collection] {
// 	return h.collectionClientDataHandler
// }
//
// func (h *clientMediaItemDataHandlersImpl) PlaylistClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Playlist] {
// 	return h.playlistClientDataHandler
// }
//
// func (h *clientMediaItemDataHandlersImpl) SeasonClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Season] {
// 	return h.seasonClientDataHandler
// }
//
// // coreCoreUserMediaItemDataRepositoriesImpl
// type coreCoreUserMediaItemDataRepositoriesImpl struct {
// 	movieCoreRepo      repo.CoreUserMediaItemDataRepository[*mediatypes.Movie]
// 	seriesCoreRepo     repo.CoreUserMediaItemDataRepository[*mediatypes.Series]
// 	episodeCoreRepo    repo.CoreUserMediaItemDataRepository[*mediatypes.Episode]
// 	trackCoreRepo      repo.CoreUserMediaItemDataRepository[*mediatypes.Track]
// 	albumCoreRepo      repo.CoreUserMediaItemDataRepository[*mediatypes.Album]
// 	artistCoreRepo     repo.CoreUserMediaItemDataRepository[*mediatypes.Artist]
// 	collectionCoreRepo repo.CoreUserMediaItemDataRepository[*mediatypes.Collection]
// 	playlistCoreRepo   repo.CoreUserMediaItemDataRepository[*mediatypes.Playlist]
// }
//
// func (r *coreCoreUserMediaItemDataRepositoriesImpl) MovieCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Movie] {
// 	return r.movieCoreRepo
// }
//
// func (r *coreCoreUserMediaItemDataRepositoriesImpl) SeriesCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Series] {
// 	return r.seriesCoreRepo
// }
//
// func (r *coreCoreUserMediaItemDataRepositoriesImpl) EpisodeCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Episode] {
// 	return r.episodeCoreRepo
// }
//
// func (r *coreCoreUserMediaItemDataRepositoriesImpl) TrackCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Track] {
// 	return r.trackCoreRepo
// }
//
// func (r *coreCoreUserMediaItemDataRepositoriesImpl) AlbumCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Album] {
// 	return r.albumCoreRepo
// }
//
// func (r *coreCoreUserMediaItemDataRepositoriesImpl) ArtistCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Artist] {
// 	return r.artistCoreRepo
// }
//
// func (r *coreCoreUserMediaItemDataRepositoriesImpl) CollectionCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Collection] {
// 	return r.collectionCoreRepo
// }
// func (r *coreCoreUserMediaItemDataRepositoriesImpl) PlaylistCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Playlist] {
// 	return r.playlistCoreRepo
// }
