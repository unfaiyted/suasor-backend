// app/di/media_factory.go
package di

import (
	"gorm.io/gorm"
	"suasor/app/factory"
	apphandlers "suasor/app/handlers"
	"suasor/app/repository"
	"suasor/app/services"
	"suasor/client"
	mediatypes "suasor/client/media/types"
	"suasor/handlers"
	repo "suasor/repository"
	svc "suasor/services"
)

// mediaDataFactoryImpl is an implementation of the MediaDataFactory interface
type mediaDataFactoryImpl struct {
	db            *gorm.DB
	clientFactory *client.ClientFactoryService
}

// createMediaDataFactory creates a new MediaDataFactory implementation
func createMediaDataFactory(db *gorm.DB, clientFactory *client.ClientFactoryService) factory.MediaDataFactory {
	return &mediaDataFactoryImpl{
		db:            db,
		clientFactory: clientFactory,
	}
}

// --------------------------------------------------------
// Core Repository Factory Methods
// --------------------------------------------------------

// CreateCoreRepositories initializes all core repositories
func (f *mediaDataFactoryImpl) CreateCoreMediaItemRepositories() repository.CoreMediaItemRepositories {
	return &coreRepositoriesImpl{
		movieRepo:      repo.NewMediaItemRepository[*mediatypes.Movie](f.db),
		seriesRepo:     repo.NewMediaItemRepository[*mediatypes.Series](f.db),
		episodeRepo:    repo.NewMediaItemRepository[*mediatypes.Episode](f.db),
		trackRepo:      repo.NewMediaItemRepository[*mediatypes.Track](f.db),
		albumRepo:      repo.NewMediaItemRepository[*mediatypes.Album](f.db),
		artistRepo:     repo.NewMediaItemRepository[*mediatypes.Artist](f.db),
		collectionRepo: repo.NewMediaItemRepository[*mediatypes.Collection](f.db),
		playlistRepo:   repo.NewMediaItemRepository[*mediatypes.Playlist](f.db),
	}
}

// CreateCoreDataRepositories initializes all core user data repositories
func (f *mediaDataFactoryImpl) CreateCoreDataRepositories() repository.CoreUserMediaItemDataRepositories {
	return &coreUserMediaItemDataRepositoriesImpl{
		movieCoreService:      repo.NewCoreUserMediaItemDataRepository[*mediatypes.Movie](f.db),
		seriesCoreService:     repo.NewCoreUserMediaItemDataRepository[*mediatypes.Series](f.db),
		episodeCoreService:    repo.NewCoreUserMediaItemDataRepository[*mediatypes.Episode](f.db),
		trackCoreService:      repo.NewCoreUserMediaItemDataRepository[*mediatypes.Track](f.db),
		albumCoreService:      repo.NewCoreUserMediaItemDataRepository[*mediatypes.Album](f.db),
		artistCoreService:     repo.NewCoreUserMediaItemDataRepository[*mediatypes.Artist](f.db),
		collectionCoreService: repo.NewCoreUserMediaItemDataRepository[*mediatypes.Collection](f.db),
		playlistCoreService:   repo.NewCoreUserMediaItemDataRepository[*mediatypes.Playlist](f.db),
	}
}

// --------------------------------------------------------
// User Repository Factory Methods
// --------------------------------------------------------

// CreateUserRepositories initializes all user repositories
func (f *mediaDataFactoryImpl) CreateUserMediaItemRepositories() repository.UserMediaItemRepositories {
	return &userRepositoryFactoriesImpl{
		movieUserRepo:      repo.NewUserMediaItemRepository[*mediatypes.Movie](f.db),
		seriesUserRepo:     repo.NewUserMediaItemRepository[*mediatypes.Series](f.db),
		episodeUserRepo:    repo.NewUserMediaItemRepository[*mediatypes.Episode](f.db),
		trackUserRepo:      repo.NewUserMediaItemRepository[*mediatypes.Track](f.db),
		albumUserRepo:      repo.NewUserMediaItemRepository[*mediatypes.Album](f.db),
		artistUserRepo:     repo.NewUserMediaItemRepository[*mediatypes.Artist](f.db),
		collectionUserRepo: repo.NewUserMediaItemRepository[*mediatypes.Collection](f.db),
		playlistUserRepo:   repo.NewUserMediaItemRepository[*mediatypes.Playlist](f.db),
	}
}

// --------------------------------------------------------
// Client Repository Factory Methods
// --------------------------------------------------------

// CreateClientRepositories initializes all client repositories
func (f *mediaDataFactoryImpl) CreateClientRepositories() repository.ClientMediaItemRepositories {
	return &clientRepositoryFactoriesImpl{
		movieClientRepo:      repo.NewClientMediaItemRepository[*mediatypes.Movie](f.db),
		seriesClientRepo:     repo.NewClientMediaItemRepository[*mediatypes.Series](f.db),
		episodeClientRepo:    repo.NewClientMediaItemRepository[*mediatypes.Episode](f.db),
		trackClientRepo:      repo.NewClientMediaItemRepository[*mediatypes.Track](f.db),
		albumClientRepo:      repo.NewClientMediaItemRepository[*mediatypes.Album](f.db),
		artistClientRepo:     repo.NewClientMediaItemRepository[*mediatypes.Artist](f.db),
		collectionClientRepo: repo.NewClientMediaItemRepository[*mediatypes.Collection](f.db),
		playlistClientRepo:   repo.NewClientMediaItemRepository[*mediatypes.Playlist](f.db),
	}
}

// CreateUserDataRepositories initializes all user data repositories
func (f *mediaDataFactoryImpl) CreateUserDataRepositories() repository.UserMediaDataRepositories {
	return &userDataRepositoriesImpl{
		movieDataRepo:      repo.NewUserMediaItemDataRepository[*mediatypes.Movie](f.db),
		seriesDataRepo:     repo.NewUserMediaItemDataRepository[*mediatypes.Series](f.db),
		episodeDataRepo:    repo.NewUserMediaItemDataRepository[*mediatypes.Episode](f.db),
		trackDataRepo:      repo.NewUserMediaItemDataRepository[*mediatypes.Track](f.db),
		albumDataRepo:      repo.NewUserMediaItemDataRepository[*mediatypes.Album](f.db),
		artistDataRepo:     repo.NewUserMediaItemDataRepository[*mediatypes.Artist](f.db),
		collectionDataRepo: repo.NewUserMediaItemDataRepository[*mediatypes.Collection](f.db),
		playlistDataRepo:   repo.NewUserMediaItemDataRepository[*mediatypes.Playlist](f.db),
	}
}

// --------------------------------------------------------
// Core Service Factory Methods
// --------------------------------------------------------

// CreateCoreServices initializes all core services
func (f *mediaDataFactoryImpl) CreateCoreServices(repos repository.CoreMediaItemRepositories) services.CoreMediaItemServices {
	return &coreMediaItemServicesImpl{
		movieCoreService:      svc.NewCoreMediaItemService[*mediatypes.Movie](repos.MovieRepo()),
		seriesCoreService:     svc.NewCoreMediaItemService[*mediatypes.Series](repos.SeriesRepo()),
		episodeCoreService:    svc.NewCoreMediaItemService[*mediatypes.Episode](repos.EpisodeRepo()),
		trackCoreService:      svc.NewCoreMediaItemService[*mediatypes.Track](repos.TrackRepo()),
		albumCoreService:      svc.NewCoreMediaItemService[*mediatypes.Album](repos.AlbumRepo()),
		artistCoreService:     svc.NewCoreMediaItemService[*mediatypes.Artist](repos.ArtistRepo()),
		collectionCoreService: svc.NewCoreMediaItemService[*mediatypes.Collection](repos.CollectionRepo()),
		playlistCoreService:   svc.NewCoreMediaItemService[*mediatypes.Playlist](repos.PlaylistRepo()),
	}
}

// CreateCoreDataServices initializes all core data services
func (f *mediaDataFactoryImpl) CreateCoreDataServices(repos repository.CoreMediaItemRepositories) services.CoreUserMediaItemDataServices {
	return &coreUserMediaItemDataServicesImpl{
		movieCoreService:      svc.NewCoreUserMediaItemDataService[*mediatypes.Movie](svc.NewCoreMediaItemService[*mediatypes.Movie](repos.MovieRepo())),
		seriesCoreService:     svc.NewCoreUserMediaItemDataService[*mediatypes.Series](svc.NewCoreMediaItemService[*mediatypes.Series](repos.SeriesRepo())),
		episodeCoreService:    svc.NewCoreUserMediaItemDataService[*mediatypes.Episode](svc.NewCoreMediaItemService[*mediatypes.Episode](repos.EpisodeRepo())),
		trackCoreService:      svc.NewCoreUserMediaItemDataService[*mediatypes.Track](svc.NewCoreMediaItemService[*mediatypes.Track](repos.TrackRepo())),
		albumCoreService:      svc.NewCoreUserMediaItemDataService[*mediatypes.Album](svc.NewCoreMediaItemService[*mediatypes.Album](repos.AlbumRepo())),
		artistCoreService:     svc.NewCoreUserMediaItemDataService[*mediatypes.Artist](svc.NewCoreMediaItemService[*mediatypes.Artist](repos.ArtistRepo())),
		collectionCoreService: svc.NewCoreUserMediaItemDataService[*mediatypes.Collection](svc.NewCoreMediaItemService[*mediatypes.Collection](repos.CollectionRepo())),
		playlistCoreService:   svc.NewCoreUserMediaItemDataService[*mediatypes.Playlist](svc.NewCoreMediaItemService[*mediatypes.Playlist](repos.PlaylistRepo())),
	}
}

// --------------------------------------------------------
// User Service Factory Methods
// --------------------------------------------------------

// CreateUserServices initializes all user services
func (f *mediaDataFactoryImpl) CreateUserServices(
	coreServices services.CoreMediaItemServices,
	userRepos repository.UserMediaItemRepositories) services.UserMediaItemServices {

	return &userMediaItemServicesImpl{
		movieUserService: svc.NewUserMediaItemService[*mediatypes.Movie](
			coreServices.MovieCoreService(), userRepos.MovieUserRepo()),
		seriesUserService: svc.NewUserMediaItemService[*mediatypes.Series](
			coreServices.SeriesCoreService(), userRepos.SeriesUserRepo()),
		episodeUserService: svc.NewUserMediaItemService[*mediatypes.Episode](
			coreServices.EpisodeCoreService(), userRepos.EpisodeUserRepo()),
		trackUserService: svc.NewUserMediaItemService[*mediatypes.Track](
			coreServices.TrackCoreService(), userRepos.TrackUserRepo()),
		albumUserService: svc.NewUserMediaItemService[*mediatypes.Album](
			coreServices.AlbumCoreService(), userRepos.AlbumUserRepo()),
		artistUserService: svc.NewUserMediaItemService[*mediatypes.Artist](
			coreServices.ArtistCoreService(), userRepos.ArtistUserRepo()),
		collectionUserService: svc.NewUserMediaItemService[*mediatypes.Collection](
			coreServices.CollectionCoreService(), userRepos.CollectionUserRepo()),
		playlistUserService: svc.NewUserMediaItemService[*mediatypes.Playlist](
			coreServices.PlaylistCoreService(), userRepos.PlaylistUserRepo()),
	}
}

// CreateUserDataServices initializes all user data services
func (f *mediaDataFactoryImpl) CreateUserDataServices(
	coreDataServices services.CoreUserMediaItemDataServices,
	userRepos repository.UserMediaDataRepositories) services.UserMediaItemDataServices {

	return &userMediaItemDataServicesImpl{
		movieDataService: svc.NewUserMediaItemDataService[*mediatypes.Movie](
			coreDataServices.MovieCoreService(), userRepos.MovieDataRepo()),
		seriesDataService: svc.NewUserMediaItemDataService[*mediatypes.Series](
			coreDataServices.SeriesCoreService(), userRepos.SeriesDataRepo()),
		episodeDataService: svc.NewUserMediaItemDataService[*mediatypes.Episode](
			coreDataServices.EpisodeCoreService(), userRepos.EpisodeDataRepo()),
		trackDataService: svc.NewUserMediaItemDataService[*mediatypes.Track](
			coreDataServices.TrackCoreService(), userRepos.TrackDataRepo()),
		albumDataService: svc.NewUserMediaItemDataService[*mediatypes.Album](
			coreDataServices.AlbumCoreService(), userRepos.AlbumDataRepo()),
		artistDataService: svc.NewUserMediaItemDataService[*mediatypes.Artist](
			coreDataServices.ArtistCoreService(), userRepos.ArtistDataRepo()),
		collectionDataService: svc.NewUserMediaItemDataService[*mediatypes.Collection](
			coreDataServices.CollectionCoreService(), userRepos.CollectionDataRepo()),
		playlistDataService: svc.NewUserMediaItemDataService[*mediatypes.Playlist](
			coreDataServices.PlaylistCoreService(), userRepos.PlaylistDataRepo()),
	}
}

// --------------------------------------------------------
// Client Service Factory Methods
// --------------------------------------------------------

// CreateClientServices initializes all client services
func (f *mediaDataFactoryImpl) CreateClientServices(
	coreServices services.CoreMediaItemServices,
	clientRepos repository.ClientMediaItemRepositories) services.ClientMediaItemServices {

	return &clientMediaItemServicesImpl{
		movieClientService: svc.NewClientMediaItemService[*mediatypes.Movie](
			coreServices.MovieCoreService(), clientRepos.MovieClientRepo()),
		seriesClientService: svc.NewClientMediaItemService[*mediatypes.Series](
			coreServices.SeriesCoreService(), clientRepos.SeriesClientRepo()),
		episodeClientService: svc.NewClientMediaItemService[*mediatypes.Episode](
			coreServices.EpisodeCoreService(), clientRepos.EpisodeClientRepo()),
		trackClientService: svc.NewClientMediaItemService[*mediatypes.Track](
			coreServices.TrackCoreService(), clientRepos.TrackClientRepo()),
		albumClientService: svc.NewClientMediaItemService[*mediatypes.Album](
			coreServices.AlbumCoreService(), clientRepos.AlbumClientRepo()),
		artistClientService: svc.NewClientMediaItemService[*mediatypes.Artist](
			coreServices.ArtistCoreService(), clientRepos.ArtistClientRepo()),
		collectionClientService: svc.NewClientMediaItemService[*mediatypes.Collection](
			coreServices.CollectionCoreService(), clientRepos.CollectionClientRepo()),
		playlistClientService: svc.NewClientMediaItemService[*mediatypes.Playlist](
			coreServices.PlaylistCoreService(), clientRepos.PlaylistClientRepo()),
	}
}

// CreateClientDataServices initializes all client data services
func (f *mediaDataFactoryImpl) CreateClientDataServices(
	userDataServices services.UserMediaItemDataServices,
	clientRepos repository.ClientUserMediaDataRepositories) services.ClientUserMediaItemDataServices {

	return &clientUserMediaItemDataServicesImpl{
		movieDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Movie](
			userDataServices.MovieDataService(), clientRepos.MovieDataRepo()),
		seriesDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Series](
			userDataServices.SeriesDataService(), clientRepos.SeriesDataRepo()),
		episodeDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Episode](
			userDataServices.EpisodeDataService(), clientRepos.EpisodeDataRepo()),
		trackDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Track](
			userDataServices.TrackDataService(), clientRepos.TrackDataRepo()),
		albumDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Album](
			userDataServices.AlbumDataService(), clientRepos.AlbumDataRepo()),
		artistDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Artist](
			userDataServices.ArtistDataService(), clientRepos.ArtistDataRepo()),
		collectionDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Collection](
			userDataServices.CollectionDataService(), clientRepos.CollectionDataRepo()),
		playlistDataService: svc.NewClientUserMediaItemDataService[*mediatypes.Playlist](
			userDataServices.PlaylistDataService(), clientRepos.PlaylistDataRepo()),
	}
}

// --------------------------------------------------------
// Specialized Collection Services
// --------------------------------------------------------

// CreateMediaCollectionServices creates collection and playlist services
func (f *mediaDataFactoryImpl) CreateMediaCollectionServices(
	coreServices services.CoreMediaItemServices,
	userServices services.UserMediaItemServices,
	clientServices services.ClientMediaItemServices,
	coreCollectionService services.CoreCollectionService,
	userCollectionService services.UserCollectionService,
	clientCollectionService services.ClientMediaCollectionService,
	playlistService services.PlaylistService) services.MediaCollectionServices {

	return &mediaCollectionServicesImpl{
		coreCollectionService:   coreCollectionService,
		userCollectionService:   userCollectionService,
		clientCollectionService: clientCollectionService,

		corePlaylistService:   coreServices.PlaylistCoreService(),
		userPlaylistService:   userServices.PlaylistUserService(),
		clientPlaylistService: clientServices.PlaylistClientService(),

		playlistService: playlistService,
	}
}

// --------------------------------------------------------
// Core Handler Factory Methods
// --------------------------------------------------------

// CreateCoreHandlers initializes all core handlers
func (f *mediaDataFactoryImpl) CreateCoreHandlers(
	coreServices services.CoreUserMediaItemDataServices) apphandlers.CoreMediaItemDataHandlers {

	return &coreMediaItemDataHandlersImpl{
		movieCoreDataHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Movie](


			)
		//
		// seriesCoreDataHandler: apphandlers.NewCoreUserMediaItemDataHandler[*mediatypes.Series](
		// 	coreServices.SeriesCoreService()),
		// episodeCoreDataHandler: apphandlers.NewCoreUserMediaItemDataHandler[*mediatypes.Episode](
		// 	coreServices.EpisodeCoreService()),
		// trackCoreDataHandler: apphandlers.NewCoreMediaItemDataHandler[*mediatypes.Track](
		// 	coreServices.TrackCoreService()),
		// albumCoreDataHandler: apphandlers.NewCoreMediaItemDataHandler[*mediatypes.Album](
		// 	coreServices.AlbumCoreService()),
		// artistCoreDataHandler: apphandlers.NewCoreMediaItemDataHandler[*mediatypes.Artist](
		// 	coreServices.ArtistCoreService()),
		// collectionCoreDataHandler: apphandlers.NewCoreMediaItemDataHandler[*mediatypes.Collection](
		// 	coreServices.CollectionCoreService()),
		// playlistCoreDataHandler: apphandlers.NewCoreMediaItemDataHandler[*mediatypes.Playlist](
		// 	coreServices.PlaylistCoreService()),
	}
}

// --------------------------------------------------------
// User Handler Factory Methods
// --------------------------------------------------------

// CreateUserHandlers initializes all user handlers
func (f *mediaDataFactoryImpl) CreateUserHandlers(
	userServices services.UserMediaItemServices,
	dataServices services.UserMediaItemDataServices,
	coreHandlers apphandlers.CoreMediaItemDataHandlers) handlers.UserMediaItemDataHandlers {

	return &userMediaItemDataHandlersImpl{
		movieUserDataHandler: apphandlers.NewUserMediaItemDataHandler[*mediatypes.Movie](
			dataServices.MovieDataService(),
			coreapphandlers.MovieCoreDataHandler()),
		seriesUserDataHandler: apphandlers.NewUserMediaItemDataHandler[*mediatypes.Series](
			dataServices.SeriesDataService(),
			coreapphandlers.SeriesCoreDataHandler()),
		episodeUserDataHandler: apphandlers.NewUserMediaItemDataHandler[*mediatypes.Episode](
			dataServices.EpisodeDataService(),
			coreapphandlers.EpisodeCoreDataHandler()),
		trackUserDataHandler: apphandlers.NewUserMediaItemDataHandler[*mediatypes.Track](
			dataServices.TrackDataService(),
			coreapphandlers.TrackCoreDataHandler()),
		albumUserDataHandler: apphandlers.NewUserMediaItemDataHandler[*mediatypes.Album](
			dataServices.AlbumDataService(),
			coreapphandlers.AlbumCoreDataHandler()),
		artistUserDataHandler: apphandlers.NewUserMediaItemDataHandler[*mediatypes.Artist](
			dataServices.ArtistDataService(),
			coreapphandlers.ArtistCoreDataHandler()),
		collectionUserDataHandler: apphandlers.NewUserMediaItemDataHandler[*mediatypes.Collection](
			dataServices.CollectionDataService(),
			coreapphandlers.CollectionCoreDataHandler()),
		playlistUserDataHandler: apphandlers.NewUserMediaItemDataHandler[*mediatypes.Playlist](
			dataServices.PlaylistDataService(),
			coreapphandlers.PlaylistCoreDataHandler()),
	}
}

// --------------------------------------------------------
// Client Handler Factory Methods
// --------------------------------------------------------

// CreateClientHandlers initializes all client handlers
func (f *mediaDataFactoryImpl) CreateClientHandlers(
	clientServices services.ClientMediaItemServices,
	dataServices services.ClientUserMediaItemDataServices,
	userHandlers apphandlers.UserMediaItemDataHandlers) handlers.ClientMediaItemDataHandlers {

	return &clientMediaItemDataHandlersImpl{
		movieClientDataHandler: apphandlers.NewClientMediaItemDataHandler[*mediatypes.Movie](
			dataServices.MovieDataService(),
			userapphandlers.MovieUserDataHandler()),
		seriesClientDataHandler: apphandlers.NewClientMediaItemDataHandler[*mediatypes.Series](
			dataServices.SeriesDataService(),
			userapphandlers.SeriesUserDataHandler()),
		episodeClientDataHandler: apphandlers.NewClientMediaItemDataHandler[*mediatypes.Episode](
			dataServices.EpisodeDataService(),
			userapphandlers.EpisodeUserDataHandler()),
		trackClientDataHandler: apphandlers.NewClientMediaItemDataHandler[*mediatypes.Track](
			dataServices.TrackDataService(),
			userapphandlers.TrackUserDataHandler()),
		albumClientDataHandler: apphandlers.NewClientMediaItemDataHandler[*mediatypes.Album](
			dataServices.AlbumDataService(),
			userapphandlers.AlbumUserDataHandler()),
		artistClientDataHandler: apphandlers.NewClientMediaItemDataHandler[*mediatypes.Artist](
			dataServices.ArtistDataService(),
			userapphandlers.ArtistUserDataHandler()),
		collectionClientDataHandler: apphandlers.NewClientMediaItemDataHandler[*mediatypes.Collection](
			dataServices.CollectionDataService(),
			userapphandlers.CollectionUserDataHandler()),
		playlistClientDataHandler: apphandlers.NewClientMediaItemDataHandler[*mediatypes.Playlist](
			dataServices.PlaylistDataService(),
			userapphandlers.PlaylistUserDataHandler()),
	}
}

// --------------------------------------------------------
// Specialized Media Handlers
// --------------------------------------------------------

// CreateSpecializedMediaHandlers creates specialized handlers for specific domains
func (f *mediaDataFactoryImpl) CreateSpecializedMediaHandlers(
	coreServices services.CoreMediaItemServices,
	userServices services.UserMediaItemServices,
	clientServices services.ClientMediaItemServices,
	musicHandler apphandlers.MusicHandler,
	seriesSpecificHandler *apphandlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig]) handlers.SpecializedMediaHandlers {

	return &specializedMediaHandlersImpl{
		musicHandler:          musicHandler,
		seriesSpecificHandler: seriesSpecificHandler,
	}
}

// --------------------------------------------------------
// Implementation structs
// --------------------------------------------------------

// Repository implementation structs

type coreRepositoriesImpl struct {
	movieRepo      repo.MediaItemRepository[*mediatypes.Movie]
	seriesRepo     repo.MediaItemRepository[*mediatypes.Series]
	episodeRepo    repo.MediaItemRepository[*mediatypes.Episode]
	trackRepo      repo.MediaItemRepository[*mediatypes.Track]
	albumRepo      repo.MediaItemRepository[*mediatypes.Album]
	artistRepo     repo.MediaItemRepository[*mediatypes.Artist]
	collectionRepo repo.MediaItemRepository[*mediatypes.Collection]
	playlistRepo   repo.MediaItemRepository[*mediatypes.Playlist]
}

func (r *coreRepositoriesImpl) MovieRepo() repo.MediaItemRepository[*mediatypes.Movie] {
	return r.movieRepo
}

func (r *coreRepositoriesImpl) SeriesRepo() repo.MediaItemRepository[*mediatypes.Series] {
	return r.seriesRepo
}

func (r *coreRepositoriesImpl) EpisodeRepo() repo.MediaItemRepository[*mediatypes.Episode] {
	return r.episodeRepo
}

func (r *coreRepositoriesImpl) TrackRepo() repo.MediaItemRepository[*mediatypes.Track] {
	return r.trackRepo
}

func (r *coreRepositoriesImpl) AlbumRepo() repo.MediaItemRepository[*mediatypes.Album] {
	return r.albumRepo
}

func (r *coreRepositoriesImpl) ArtistRepo() repo.MediaItemRepository[*mediatypes.Artist] {
	return r.artistRepo
}

func (r *coreRepositoriesImpl) CollectionRepo() repo.MediaItemRepository[*mediatypes.Collection] {
	return r.collectionRepo
}

func (r *coreRepositoriesImpl) PlaylistRepo() repo.MediaItemRepository[*mediatypes.Playlist] {
	return r.playlistRepo
}

type coreUserMediaItemDataRepositoriesImpl struct {
	movieCoreService      repo.CoreUserMediaItemDataRepository[*mediatypes.Movie]
	seriesCoreService     repo.CoreUserMediaItemDataRepository[*mediatypes.Series]
	episodeCoreService    repo.CoreUserMediaItemDataRepository[*mediatypes.Episode]
	trackCoreService      repo.CoreUserMediaItemDataRepository[*mediatypes.Track]
	albumCoreService      repo.CoreUserMediaItemDataRepository[*mediatypes.Album]
	artistCoreService     repo.CoreUserMediaItemDataRepository[*mediatypes.Artist]
	collectionCoreService repo.CoreUserMediaItemDataRepository[*mediatypes.Collection]
	playlistCoreService   repo.CoreUserMediaItemDataRepository[*mediatypes.Playlist]
}

func (r *coreUserMediaItemDataRepositoriesImpl) MovieCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Movie] {
	return r.movieCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) SeriesCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Series] {
	return r.seriesCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) EpisodeCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Episode] {
	return r.episodeCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) TrackCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Track] {
	return r.trackCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) AlbumCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Album] {
	return r.albumCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) ArtistCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Artist] {
	return r.artistCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) CollectionCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Collection] {
	return r.collectionCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) PlaylistCoreService() repo.CoreUserMediaItemDataRepository[*mediatypes.Playlist] {
	return r.playlistCoreService
}

type userRepositoryFactoriesImpl struct {
	movieUserRepo      repo.UserMediaItemRepository[*mediatypes.Movie]
	seriesUserRepo     repo.UserMediaItemRepository[*mediatypes.Series]
	episodeUserRepo    repo.UserMediaItemRepository[*mediatypes.Episode]
	trackUserRepo      repo.UserMediaItemRepository[*mediatypes.Track]
	albumUserRepo      repo.UserMediaItemRepository[*mediatypes.Album]
	artistUserRepo     repo.UserMediaItemRepository[*mediatypes.Artist]
	collectionUserRepo repo.UserMediaItemRepository[*mediatypes.Collection]
	playlistUserRepo   repo.UserMediaItemRepository[*mediatypes.Playlist]
}

func (r *userRepositoryFactoriesImpl) MovieUserRepo() repo.UserMediaItemRepository[*mediatypes.Movie] {
	return r.movieUserRepo
}

func (r *userRepositoryFactoriesImpl) SeriesUserRepo() repo.UserMediaItemRepository[*mediatypes.Series] {
	return r.seriesUserRepo
}

func (r *userRepositoryFactoriesImpl) EpisodeUserRepo() repo.UserMediaItemRepository[*mediatypes.Episode] {
	return r.episodeUserRepo
}

func (r *userRepositoryFactoriesImpl) TrackUserRepo() repo.UserMediaItemRepository[*mediatypes.Track] {
	return r.trackUserRepo
}

func (r *userRepositoryFactoriesImpl) AlbumUserRepo() repo.UserMediaItemRepository[*mediatypes.Album] {
	return r.albumUserRepo
}

func (r *userRepositoryFactoriesImpl) ArtistUserRepo() repo.UserMediaItemRepository[*mediatypes.Artist] {
	return r.artistUserRepo
}

func (r *userRepositoryFactoriesImpl) CollectionUserRepo() repo.UserMediaItemRepository[*mediatypes.Collection] {
	return r.collectionUserRepo
}

func (r *userRepositoryFactoriesImpl) PlaylistUserRepo() repo.UserMediaItemRepository[*mediatypes.Playlist] {
	return r.playlistUserRepo
}

type clientRepositoryFactoriesImpl struct {
	movieClientRepo      repo.ClientMediaItemRepository[*mediatypes.Movie]
	seriesClientRepo     repo.ClientMediaItemRepository[*mediatypes.Series]
	episodeClientRepo    repo.ClientMediaItemRepository[*mediatypes.Episode]
	trackClientRepo      repo.ClientMediaItemRepository[*mediatypes.Track]
	albumClientRepo      repo.ClientMediaItemRepository[*mediatypes.Album]
	artistClientRepo     repo.ClientMediaItemRepository[*mediatypes.Artist]
	collectionClientRepo repo.ClientMediaItemRepository[*mediatypes.Collection]
	playlistClientRepo   repo.ClientMediaItemRepository[*mediatypes.Playlist]
}

func (r *clientRepositoryFactoriesImpl) MovieClientRepo() repo.ClientMediaItemRepository[*mediatypes.Movie] {
	return r.movieClientRepo
}

func (r *clientRepositoryFactoriesImpl) SeriesClientRepo() repo.ClientMediaItemRepository[*mediatypes.Series] {
	return r.seriesClientRepo
}

func (r *clientRepositoryFactoriesImpl) EpisodeClientRepo() repo.ClientMediaItemRepository[*mediatypes.Episode] {
	return r.episodeClientRepo
}

func (r *clientRepositoryFactoriesImpl) TrackClientRepo() repo.ClientMediaItemRepository[*mediatypes.Track] {
	return r.trackClientRepo
}

func (r *clientRepositoryFactoriesImpl) AlbumClientRepo() repo.ClientMediaItemRepository[*mediatypes.Album] {
	return r.albumClientRepo
}

func (r *clientRepositoryFactoriesImpl) ArtistClientRepo() repo.ClientMediaItemRepository[*mediatypes.Artist] {
	return r.artistClientRepo
}

func (r *clientRepositoryFactoriesImpl) CollectionClientRepo() repo.ClientMediaItemRepository[*mediatypes.Collection] {
	return r.collectionClientRepo
}

func (r *clientRepositoryFactoriesImpl) PlaylistClientRepo() repo.ClientMediaItemRepository[*mediatypes.Playlist] {
	return r.playlistClientRepo
}

type userDataRepositoriesImpl struct {
	movieDataRepo      repo.UserMediaItemDataRepository[*mediatypes.Movie]
	seriesDataRepo     repo.UserMediaItemDataRepository[*mediatypes.Series]
	episodeDataRepo    repo.UserMediaItemDataRepository[*mediatypes.Episode]
	trackDataRepo      repo.UserMediaItemDataRepository[*mediatypes.Track]
	albumDataRepo      repo.UserMediaItemDataRepository[*mediatypes.Album]
	artistDataRepo     repo.UserMediaItemDataRepository[*mediatypes.Artist]
	collectionDataRepo repo.UserMediaItemDataRepository[*mediatypes.Collection]
	playlistDataRepo   repo.UserMediaItemDataRepository[*mediatypes.Playlist]
}

func (r *userDataRepositoriesImpl) MovieDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Movie] {
	return r.movieDataRepo
}

func (r *userDataRepositoriesImpl) SeriesDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Series] {
	return r.seriesDataRepo
}

func (r *userDataRepositoriesImpl) EpisodeDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Episode] {
	return r.episodeDataRepo
}

func (r *userDataRepositoriesImpl) TrackDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Track] {
	return r.trackDataRepo
}

func (r *userDataRepositoriesImpl) AlbumDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Album] {
	return r.albumDataRepo
}

func (r *userDataRepositoriesImpl) ArtistDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Artist] {
	return r.artistDataRepo
}

func (r *userDataRepositoriesImpl) CollectionDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Collection] {
	return r.collectionDataRepo
}

func (r *userDataRepositoriesImpl) PlaylistDataRepo() repo.UserMediaItemDataRepository[*mediatypes.Playlist] {
	return r.playlistDataRepo
}

// Service implementation structs

type coreMediaItemServicesImpl struct {
	movieCoreService      svc.CoreMediaItemService[*mediatypes.Movie]
	seriesCoreService     svc.CoreMediaItemService[*mediatypes.Series]
	episodeCoreService    svc.CoreMediaItemService[*mediatypes.Episode]
	trackCoreService      svc.CoreMediaItemService[*mediatypes.Track]
	albumCoreService      svc.CoreMediaItemService[*mediatypes.Album]
	artistCoreService     svc.CoreMediaItemService[*mediatypes.Artist]
	collectionCoreService svc.CoreMediaItemService[*mediatypes.Collection]
	playlistCoreService   svc.CoreMediaItemService[*mediatypes.Playlist]
}

func (s *coreMediaItemServicesImpl) MovieCoreService() svc.CoreMediaItemService[*mediatypes.Movie] {
	return s.movieCoreService
}

func (s *coreMediaItemServicesImpl) SeriesCoreService() svc.CoreMediaItemService[*mediatypes.Series] {
	return s.seriesCoreService
}

func (s *coreMediaItemServicesImpl) EpisodeCoreService() svc.CoreMediaItemService[*mediatypes.Episode] {
	return s.episodeCoreService
}

func (s *coreMediaItemServicesImpl) TrackCoreService() svc.CoreMediaItemService[*mediatypes.Track] {
	return s.trackCoreService
}

func (s *coreMediaItemServicesImpl) AlbumCoreService() svc.CoreMediaItemService[*mediatypes.Album] {
	return s.albumCoreService
}

func (s *coreMediaItemServicesImpl) ArtistCoreService() svc.CoreMediaItemService[*mediatypes.Artist] {
	return s.artistCoreService
}

func (s *coreMediaItemServicesImpl) CollectionCoreService() svc.CoreMediaItemService[*mediatypes.Collection] {
	return s.collectionCoreService
}

func (s *coreMediaItemServicesImpl) PlaylistCoreService() svc.CoreMediaItemService[*mediatypes.Playlist] {
	return s.playlistCoreService
}

type coreUserMediaItemDataServicesImpl struct {
	movieCoreService      svc.CoreUserMediaItemDataService[*mediatypes.Movie]
	seriesCoreService     svc.CoreUserMediaItemDataService[*mediatypes.Series]
	episodeCoreService    svc.CoreUserMediaItemDataService[*mediatypes.Episode]
	trackCoreService      svc.CoreUserMediaItemDataService[*mediatypes.Track]
	albumCoreService      svc.CoreUserMediaItemDataService[*mediatypes.Album]
	artistCoreService     svc.CoreUserMediaItemDataService[*mediatypes.Artist]
	collectionCoreService svc.CoreUserMediaItemDataService[*mediatypes.Collection]
	playlistCoreService   svc.CoreUserMediaItemDataService[*mediatypes.Playlist]
}

func (s *coreUserMediaItemDataServicesImpl) MovieCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Movie] {
	return s.movieCoreService
}

func (s *coreUserMediaItemDataServicesImpl) SeriesCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Series] {
	return s.seriesCoreService
}

func (s *coreUserMediaItemDataServicesImpl) EpisodeCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Episode] {
	return s.episodeCoreService
}

func (s *coreUserMediaItemDataServicesImpl) TrackCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Track] {
	return s.trackCoreService
}

func (s *coreUserMediaItemDataServicesImpl) AlbumCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Album] {
	return s.albumCoreService
}

func (s *coreUserMediaItemDataServicesImpl) ArtistCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Artist] {
	return s.artistCoreService
}

func (s *coreUserMediaItemDataServicesImpl) CollectionCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Collection] {
	return s.collectionCoreService
}

func (s *coreUserMediaItemDataServicesImpl) PlaylistCoreService() svc.CoreUserMediaItemDataService[*mediatypes.Playlist] {
	return s.playlistCoreService
}

type userMediaItemServicesImpl struct {
	movieUserService      svc.UserMediaItemService[*mediatypes.Movie]
	seriesUserService     svc.UserMediaItemService[*mediatypes.Series]
	episodeUserService    svc.UserMediaItemService[*mediatypes.Episode]
	trackUserService      svc.UserMediaItemService[*mediatypes.Track]
	albumUserService      svc.UserMediaItemService[*mediatypes.Album]
	artistUserService     svc.UserMediaItemService[*mediatypes.Artist]
	collectionUserService svc.UserMediaItemService[*mediatypes.Collection]
	playlistUserService   svc.UserMediaItemService[*mediatypes.Playlist]
}

func (s *userMediaItemServicesImpl) MovieUserService() svc.UserMediaItemService[*mediatypes.Movie] {
	return s.movieUserService
}

func (s *userMediaItemServicesImpl) SeriesUserService() svc.UserMediaItemService[*mediatypes.Series] {
	return s.seriesUserService
}

func (s *userMediaItemServicesImpl) EpisodeUserService() svc.UserMediaItemService[*mediatypes.Episode] {
	return s.episodeUserService
}

func (s *userMediaItemServicesImpl) TrackUserService() svc.UserMediaItemService[*mediatypes.Track] {
	return s.trackUserService
}

func (s *userMediaItemServicesImpl) AlbumUserService() svc.UserMediaItemService[*mediatypes.Album] {
	return s.albumUserService
}

func (s *userMediaItemServicesImpl) ArtistUserService() svc.UserMediaItemService[*mediatypes.Artist] {
	return s.artistUserService
}

func (s *userMediaItemServicesImpl) CollectionUserService() svc.UserMediaItemService[*mediatypes.Collection] {
	return s.collectionUserService
}

func (s *userMediaItemServicesImpl) PlaylistUserService() svc.UserMediaItemService[*mediatypes.Playlist] {
	return s.playlistUserService
}

type userMediaItemDataServicesImpl struct {
	movieDataService      svc.UserMediaItemDataService[*mediatypes.Movie]
	seriesDataService     svc.UserMediaItemDataService[*mediatypes.Series]
	episodeDataService    svc.UserMediaItemDataService[*mediatypes.Episode]
	trackDataService      svc.UserMediaItemDataService[*mediatypes.Track]
	albumDataService      svc.UserMediaItemDataService[*mediatypes.Album]
	artistDataService     svc.UserMediaItemDataService[*mediatypes.Artist]
	collectionDataService svc.UserMediaItemDataService[*mediatypes.Collection]
	playlistDataService   svc.UserMediaItemDataService[*mediatypes.Playlist]
}

func (s *userMediaItemDataServicesImpl) MovieDataService() svc.UserMediaItemDataService[*mediatypes.Movie] {
	return s.movieDataService
}

func (s *userMediaItemDataServicesImpl) SeriesDataService() svc.UserMediaItemDataService[*mediatypes.Series] {
	return s.seriesDataService
}

func (s *userMediaItemDataServicesImpl) EpisodeDataService() svc.UserMediaItemDataService[*mediatypes.Episode] {
	return s.episodeDataService
}

func (s *userMediaItemDataServicesImpl) TrackDataService() svc.UserMediaItemDataService[*mediatypes.Track] {
	return s.trackDataService
}

func (s *userMediaItemDataServicesImpl) AlbumDataService() svc.UserMediaItemDataService[*mediatypes.Album] {
	return s.albumDataService
}

func (s *userMediaItemDataServicesImpl) ArtistDataService() svc.UserMediaItemDataService[*mediatypes.Artist] {
	return s.artistDataService
}

func (s *userMediaItemDataServicesImpl) CollectionDataService() svc.UserMediaItemDataService[*mediatypes.Collection] {
	return s.collectionDataService
}

func (s *userMediaItemDataServicesImpl) PlaylistDataService() svc.UserMediaItemDataService[*mediatypes.Playlist] {
	return s.playlistDataService
}

type clientMediaItemServicesImpl struct {
	movieClientService      svc.ClientMediaItemService[*mediatypes.Movie]
	seriesClientService     svc.ClientMediaItemService[*mediatypes.Series]
	episodeClientService    svc.ClientMediaItemService[*mediatypes.Episode]
	trackClientService      svc.ClientMediaItemService[*mediatypes.Track]
	albumClientService      svc.ClientMediaItemService[*mediatypes.Album]
	artistClientService     svc.ClientMediaItemService[*mediatypes.Artist]
	collectionClientService svc.ClientMediaItemService[*mediatypes.Collection]
	playlistClientService   svc.ClientMediaItemService[*mediatypes.Playlist]
}

func (s *clientMediaItemServicesImpl) MovieClientService() svc.ClientMediaItemService[*mediatypes.Movie] {
	return s.movieClientService
}

func (s *clientMediaItemServicesImpl) SeriesClientService() svc.ClientMediaItemService[*mediatypes.Series] {
	return s.seriesClientService
}

func (s *clientMediaItemServicesImpl) EpisodeClientService() svc.ClientMediaItemService[*mediatypes.Episode] {
	return s.episodeClientService
}

func (s *clientMediaItemServicesImpl) TrackClientService() svc.ClientMediaItemService[*mediatypes.Track] {
	return s.trackClientService
}

func (s *clientMediaItemServicesImpl) AlbumClientService() svc.ClientMediaItemService[*mediatypes.Album] {
	return s.albumClientService
}

func (s *clientMediaItemServicesImpl) ArtistClientService() svc.ClientMediaItemService[*mediatypes.Artist] {
	return s.artistClientService
}

func (s *clientMediaItemServicesImpl) CollectionClientService() svc.ClientMediaItemService[*mediatypes.Collection] {
	return s.collectionClientService
}

func (s *clientMediaItemServicesImpl) PlaylistClientService() svc.ClientMediaItemService[*mediatypes.Playlist] {
	return s.playlistClientService
}

type clientUserMediaItemDataServicesImpl struct {
	movieDataService      svc.ClientUserMediaItemDataService[*mediatypes.Movie]
	seriesDataService     svc.ClientUserMediaItemDataService[*mediatypes.Series]
	episodeDataService    svc.ClientUserMediaItemDataService[*mediatypes.Episode]
	trackDataService      svc.ClientUserMediaItemDataService[*mediatypes.Track]
	albumDataService      svc.ClientUserMediaItemDataService[*mediatypes.Album]
	artistDataService     svc.ClientUserMediaItemDataService[*mediatypes.Artist]
	collectionDataService svc.ClientUserMediaItemDataService[*mediatypes.Collection]
	playlistDataService   svc.ClientUserMediaItemDataService[*mediatypes.Playlist]
}

func (s *clientUserMediaItemDataServicesImpl) MovieDataService() svc.ClientUserMediaItemDataService[*mediatypes.Movie] {
	return s.movieDataService
}

func (s *clientUserMediaItemDataServicesImpl) SeriesDataService() svc.ClientUserMediaItemDataService[*mediatypes.Series] {
	return s.seriesDataService
}

func (s *clientUserMediaItemDataServicesImpl) EpisodeDataService() svc.ClientUserMediaItemDataService[*mediatypes.Episode] {
	return s.episodeDataService
}

func (s *clientUserMediaItemDataServicesImpl) TrackDataService() svc.ClientUserMediaItemDataService[*mediatypes.Track] {
	return s.trackDataService
}

func (s *clientUserMediaItemDataServicesImpl) AlbumDataService() svc.ClientUserMediaItemDataService[*mediatypes.Album] {
	return s.albumDataService
}

func (s *clientUserMediaItemDataServicesImpl) ArtistDataService() svc.ClientUserMediaItemDataService[*mediatypes.Artist] {
	return s.artistDataService
}

func (s *clientUserMediaItemDataServicesImpl) CollectionDataService() svc.ClientUserMediaItemDataService[*mediatypes.Collection] {
	return s.collectionDataService
}

func (s *clientUserMediaItemDataServicesImpl) PlaylistDataService() svc.ClientUserMediaItemDataService[*mediatypes.Playlist] {
	return s.playlistDataService
}

type mediaCollectionServicesImpl struct {
	coreCollectionService   services.CoreCollectionService
	userCollectionService   services.UserCollectionService
	clientCollectionService services.ClientMediaCollectionService

	corePlaylistService   svc.CoreMediaItemService[*mediatypes.Playlist]
	userPlaylistService   svc.UserMediaItemService[*mediatypes.Playlist]
	clientPlaylistService svc.ClientMediaItemService[*mediatypes.Playlist]

	playlistService services.PlaylistService
}

func (s *mediaCollectionServicesImpl) CoreCollectionService() services.CoreCollectionService {
	return s.coreCollectionService
}

func (s *mediaCollectionServicesImpl) UserCollectionService() services.UserCollectionService {
	return s.userCollectionService
}

func (s *mediaCollectionServicesImpl) ClientCollectionService() services.ClientMediaCollectionService {
	return s.clientCollectionService
}

func (s *mediaCollectionServicesImpl) CorePlaylistService() svc.CoreMediaItemService[*mediatypes.Playlist] {
	return s.corePlaylistService
}

func (s *mediaCollectionServicesImpl) UserPlaylistService() svc.UserMediaItemService[*mediatypes.Playlist] {
	return s.userPlaylistService
}

func (s *mediaCollectionServicesImpl) ClientPlaylistService() svc.ClientMediaItemService[*mediatypes.Playlist] {
	return s.clientPlaylistService
}

func (s *mediaCollectionServicesImpl) PlaylistService() services.PlaylistService {
	return s.playlistService
}

// Handler implementation structs

type coreMediaItemDataHandlersImpl struct {
	movieCoreDataHandler      *apphandlers.CoreMediaItemDataHandler[*mediatypes.Movie]
	seriesCoreDataHandler     *apphandlers.CoreMediaItemDataHandler[*mediatypes.Series]
	episodeCoreDataHandler    *apphandlers.CoreMediaItemDataHandler[*mediatypes.Episode]
	trackCoreDataHandler      *apphandlers.CoreMediaItemDataHandler[*mediatypes.Track]
	albumCoreDataHandler      *apphandlers.CoreMediaItemDataHandler[*mediatypes.Album]
	artistCoreDataHandler     *apphandlers.CoreMediaItemDataHandler[*mediatypes.Artist]
	collectionCoreDataHandler *apphandlers.CoreMediaItemDataHandler[*mediatypes.Collection]
	playlistCoreDataHandler   *apphandlers.CoreMediaItemDataHandler[*mediatypes.Playlist]
}

func (h *coreMediaItemDataHandlersImpl) MovieCoreDataHandler() *apphandlers.CoreMediaItemDataHandler[*mediatypes.Movie] {
	return h.movieCoreDataHandler
}

func (h *coreMediaItemDataHandlersImpl) SeriesCoreDataHandler() *apphandlers.CoreMediaItemDataHandler[*mediatypes.Series] {
	return h.seriesCoreDataHandler
}

func (h *coreMediaItemDataHandlersImpl) EpisodeCoreDataHandler() *apphandlers.CoreMediaItemDataHandler[*mediatypes.Episode] {
	return h.episodeCoreDataHandler
}

func (h *coreMediaItemDataHandlersImpl) TrackCoreDataHandler() *apphandlers.CoreMediaItemDataHandler[*mediatypes.Track] {
	return h.trackCoreDataHandler
}

func (h *coreMediaItemDataHandlersImpl) AlbumCoreDataHandler() *apphandlers.CoreMediaItemDataHandler[*mediatypes.Album] {
	return h.albumCoreDataHandler
}

func (h *coreMediaItemDataHandlersImpl) ArtistCoreDataHandler() *apphandlers.CoreMediaItemDataHandler[*mediatypes.Artist] {
	return h.artistCoreDataHandler
}

func (h *coreMediaItemDataHandlersImpl) CollectionCoreDataHandler() *apphandlers.CoreMediaItemDataHandler[*mediatypes.Collection] {
	return h.collectionCoreDataHandler
}

func (h *coreMediaItemDataHandlersImpl) PlaylistCoreDataHandler() *apphandlers.CoreMediaItemDataHandler[*mediatypes.Playlist] {
	return h.playlistCoreDataHandler
}

type userMediaItemDataHandlersImpl struct {
	movieUserDataHandler      *apphandlers.UserMediaItemDataHandler[*mediatypes.Movie]
	seriesUserDataHandler     *apphandlers.UserMediaItemDataHandler[*mediatypes.Series]
	episodeUserDataHandler    *apphandlers.UserMediaItemDataHandler[*mediatypes.Episode]
	trackUserDataHandler      *apphandlers.UserMediaItemDataHandler[*mediatypes.Track]
	albumUserDataHandler      *apphandlers.UserMediaItemDataHandler[*mediatypes.Album]
	artistUserDataHandler     *apphandlers.UserMediaItemDataHandler[*mediatypes.Artist]
	collectionUserDataHandler *apphandlers.UserMediaItemDataHandler[*mediatypes.Collection]
	playlistUserDataHandler   *apphandlers.UserMediaItemDataHandler[*mediatypes.Playlist]
}

func (h *userMediaItemDataHandlersImpl) MovieUserDataHandler() *apphandlers.UserMediaItemDataHandler[*mediatypes.Movie] {
	return h.movieUserDataHandler
}

func (h *userMediaItemDataHandlersImpl) SeriesUserDataHandler() *apphandlers.UserMediaItemDataHandler[*mediatypes.Series] {
	return h.seriesUserDataHandler
}

func (h *userMediaItemDataHandlersImpl) EpisodeUserDataHandler() *apphandlers.UserMediaItemDataHandler[*mediatypes.Episode] {
	return h.episodeUserDataHandler
}

func (h *userMediaItemDataHandlersImpl) TrackUserDataHandler() *apphandlers.UserMediaItemDataHandler[*mediatypes.Track] {
	return h.trackUserDataHandler
}

func (h *userMediaItemDataHandlersImpl) AlbumUserDataHandler() *apphandlers.UserMediaItemDataHandler[*mediatypes.Album] {
	return h.albumUserDataHandler
}

func (h *userMediaItemDataHandlersImpl) ArtistUserDataHandler() *apphandlers.UserMediaItemDataHandler[*mediatypes.Artist] {
	return h.artistUserDataHandler
}

func (h *userMediaItemDataHandlersImpl) CollectionUserDataHandler() *apphandlers.UserMediaItemDataHandler[*mediatypes.Collection] {
	return h.collectionUserDataHandler
}

func (h *userMediaItemDataHandlersImpl) PlaylistUserDataHandler() *apphandlers.UserMediaItemDataHandler[*mediatypes.Playlist] {
	return h.playlistUserDataHandler
}

type clientMediaItemDataHandlersImpl struct {
	movieClientDataHandler      *apphandlers.ClientMediaItemDataHandler[*mediatypes.Movie]
	seriesClientDataHandler     *apphandlers.ClientMediaItemDataHandler[*mediatypes.Series]
	episodeClientDataHandler    *apphandlers.ClientMediaItemDataHandler[*mediatypes.Episode]
	trackClientDataHandler      *apphandlers.ClientMediaItemDataHandler[*mediatypes.Track]
	albumClientDataHandler      *apphandlers.ClientMediaItemDataHandler[*mediatypes.Album]
	artistClientDataHandler     *apphandlers.ClientMediaItemDataHandler[*mediatypes.Artist]
	collectionClientDataHandler *apphandlers.ClientMediaItemDataHandler[*mediatypes.Collection]
	playlistClientDataHandler   *apphandlers.ClientMediaItemDataHandler[*mediatypes.Playlist]
}

func (h *clientMediaItemDataHandlersImpl) MovieClientDataHandler() *apphandlers.ClientMediaItemDataHandler[*mediatypes.Movie] {
	return h.movieClientDataHandler
}

func (h *clientMediaItemDataHandlersImpl) SeriesClientDataHandler() *apphandlers.ClientMediaItemDataHandler[*mediatypes.Series] {
	return h.seriesClientDataHandler
}

func (h *clientMediaItemDataHandlersImpl) EpisodeClientDataHandler() *apphandlers.ClientMediaItemDataHandler[*mediatypes.Episode] {
	return h.episodeClientDataHandler
}

func (h *clientMediaItemDataHandlersImpl) TrackClientDataHandler() *apphandlers.ClientMediaItemDataHandler[*mediatypes.Track] {
	return h.trackClientDataHandler
}

func (h *clientMediaItemDataHandlersImpl) AlbumClientDataHandler() *apphandlers.ClientMediaItemDataHandler[*mediatypes.Album] {
	return h.albumClientDataHandler
}

func (h *clientMediaItemDataHandlersImpl) ArtistClientDataHandler() *apphandlers.ClientMediaItemDataHandler[*mediatypes.Artist] {
	return h.artistClientDataHandler
}

func (h *clientMediaItemDataHandlersImpl) CollectionClientDataHandler() *apphandlers.ClientMediaItemDataHandler[*mediatypes.Collection] {
	return h.collectionClientDataHandler
}

func (h *clientMediaItemDataHandlersImpl) PlaylistClientDataHandler() *apphandlers.ClientMediaItemDataHandler[*mediatypes.Playlist] {
	return h.playlistClientDataHandler
}

type specializedMediaHandlersImpl struct {
	musicHandler          apphandlers.MusicHandler
	seriesSpecificHandler *apphandlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig]
}

func (h *specializedMediaHandlersImpl) MusicHandler() apphandlers.MusicHandler {
	return h.musicHandler
}

func (h *specializedMediaHandlersImpl) SeriesSpecificHandler() *apphandlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig] {
	return h.seriesSpecificHandler
}

