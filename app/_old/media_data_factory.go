// app/factory.go
package app

import (
	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"

	"gorm.io/gorm"
)

// MediaDataFactory is a factory for creating components in the three-pronged architecture
// This factory makes it easy to create properly configured repositories, services, and handlers
// for media data operations that work together in the three-pronged architecture pattern.
type MediaDataFactory struct {
	db               *gorm.DB
	clientFactory    *client.ClientFactoryService
	coreRepositories CoreMediaItemRepositories
}

// NewMediaDataFactory creates a new factory for media data components
func NewMediaDataFactory(db *gorm.DB, clientFactory *client.ClientFactoryService) *MediaDataFactory {
	return &MediaDataFactory{
		db:            db,
		clientFactory: clientFactory,
	}
}

// --------------------------------------------------------
// Core Repository Factory Methods
// --------------------------------------------------------

// CreateCoreRepositories initializes all core repositories
func (f *MediaDataFactory) CreateCoreRepositories() CoreMediaItemRepositories {
	return &coreRepositoriesImpl{
		movieRepo:      repository.NewMediaItemRepository[*mediatypes.Movie](f.db),
		seriesRepo:     repository.NewMediaItemRepository[*mediatypes.Series](f.db),
		episodeRepo:    repository.NewMediaItemRepository[*mediatypes.Episode](f.db),
		trackRepo:      repository.NewMediaItemRepository[*mediatypes.Track](f.db),
		albumRepo:      repository.NewMediaItemRepository[*mediatypes.Album](f.db),
		artistRepo:     repository.NewMediaItemRepository[*mediatypes.Artist](f.db),
		collectionRepo: repository.NewMediaItemRepository[*mediatypes.Collection](f.db),
		playlistRepo:   repository.NewMediaItemRepository[*mediatypes.Playlist](f.db),
	}
}

// CreateCoreDataRepositories initializes all core user data repositories
func (f *MediaDataFactory) CreateCoreDataRepositories() CoreUserMediaItemDataRepositories {
	return &coreUserMediaItemDataRepositoriesImpl{
		movieDataRepo:      repository.NewCoreUserMediaItemDataRepository[*mediatypes.Movie](f.db),
		seriesDataRepo:     repository.NewCoreUserMediaItemDataRepository[*mediatypes.Series](f.db),
		episodeDataRepo:    repository.NewCoreUserMediaItemDataRepository[*mediatypes.Episode](f.db),
		trackDataRepo:      repository.NewCoreUserMediaItemDataRepository[*mediatypes.Track](f.db),
		albumDataRepo:      repository.NewCoreUserMediaItemDataRepository[*mediatypes.Album](f.db),
		artistDataRepo:     repository.NewCoreUserMediaItemDataRepository[*mediatypes.Artist](f.db),
		collectionDataRepo: repository.NewCoreUserMediaItemDataRepository[*mediatypes.Collection](f.db),
		playlistDataRepo:   repository.NewCoreUserMediaItemDataRepository[*mediatypes.Playlist](f.db),
	}
}

// --------------------------------------------------------
// User Repository Factory Methods
// --------------------------------------------------------

// CreateUserRepositories initializes all user repositories
func (f *MediaDataFactory) CreateUserRepositories() UserRepositoryFactories {
	return &userRepositoryFactoriesImpl{
		movieUserRepo:      repository.NewUserMediaItemRepository[*mediatypes.Movie](f.db),
		seriesUserRepo:     repository.NewUserMediaItemRepository[*mediatypes.Series](f.db),
		episodeUserRepo:    repository.NewUserMediaItemRepository[*mediatypes.Episode](f.db),
		trackUserRepo:      repository.NewUserMediaItemRepository[*mediatypes.Track](f.db),
		albumUserRepo:      repository.NewUserMediaItemRepository[*mediatypes.Album](f.db),
		artistUserRepo:     repository.NewUserMediaItemRepository[*mediatypes.Artist](f.db),
		collectionUserRepo: repository.NewUserMediaItemRepository[*mediatypes.Collection](f.db),
		playlistUserRepo:   repository.NewUserMediaItemRepository[*mediatypes.Playlist](f.db),
	}
}

// --------------------------------------------------------
// Client Repository Factory Methods
// --------------------------------------------------------

// CreateClientRepositories initializes all client repositories
func (f *MediaDataFactory) CreateClientRepositories() ClientRepositoryFactories {
	return &clientRepositoryFactoriesImpl{
		movieClientRepo:      repository.NewClientMediaItemRepository[*mediatypes.Movie](f.db),
		seriesClientRepo:     repository.NewClientMediaItemRepository[*mediatypes.Series](f.db),
		episodeClientRepo:    repository.NewClientMediaItemRepository[*mediatypes.Episode](f.db),
		trackClientRepo:      repository.NewClientMediaItemRepository[*mediatypes.Track](f.db),
		albumClientRepo:      repository.NewClientMediaItemRepository[*mediatypes.Album](f.db),
		artistClientRepo:     repository.NewClientMediaItemRepository[*mediatypes.Artist](f.db),
		collectionClientRepo: repository.NewClientMediaItemRepository[*mediatypes.Collection](f.db),
		playlistClientRepo:   repository.NewClientMediaItemRepository[*mediatypes.Playlist](f.db),
	}
}

// CreateUserDataRepositories initializes all user data repositories
func (f *MediaDataFactory) CreateUserDataRepositories() UserDataFactories {
	return &userDataRepositoriesImpl{
		movieDataRepo: repository.NewUserMediaItemDataRepository[*mediatypes.Movie](f.db),
		seriesDataRepo: repository.NewUserMediaItemDataRepository[*mediatypes.Series](f.db),
		episodeDataRepo: repository.NewUserMediaItemDataRepository[*mediatypes.Episode](f.db),
		trackDataRepo: repository.NewUserMediaItemDataRepository[*mediatypes.Track](f.db),
		albumDataRepo: repository.NewUserMediaItemDataRepository[*mediatypes.Album](f.db),
		artistDataRepo: repository.NewUserMediaItemDataRepository[*mediatypes.Artist](f.db),
		collectionDataRepo: repository.NewUserMediaItemDataRepository[*mediatypes.Collection](f.db),
		playlistDataRepo: repository.NewUserMediaItemDataRepository[*mediatypes.Playlist](f.db),
	}
}

// --------------------------------------------------------
// Core Service Factory Methods
// --------------------------------------------------------

// CreateCoreServices initializes all core services
func (f *MediaDataFactory) CreateCoreServices(repos CoreMediaItemRepositories) CoreMediaItemServices {
	return &coreMediaItemServicesImpl{
		movieCoreService:      services.NewCoreMediaItemService[*mediatypes.Movie](repos.MovieRepo()),
		seriesCoreService:     services.NewCoreMediaItemService[*mediatypes.Series](repos.SeriesRepo()),
		episodeCoreService:    services.NewCoreMediaItemService[*mediatypes.Episode](repos.EpisodeRepo()),
		trackCoreService:      services.NewCoreMediaItemService[*mediatypes.Track](repos.TrackRepo()),
		albumCoreService:      services.NewCoreMediaItemService[*mediatypes.Album](repos.AlbumRepo()),
		artistCoreService:     services.NewCoreMediaItemService[*mediatypes.Artist](repos.ArtistRepo()),
		collectionCoreService: services.NewCoreMediaItemService[*mediatypes.Collection](repos.CollectionRepo()),
		playlistCoreService:   services.NewCoreMediaItemService[*mediatypes.Playlist](repos.PlaylistRepo()),
	}
}

// --------------------------------------------------------
// User Service Factory Methods
// --------------------------------------------------------

// CreateUserServices initializes all user services
func (f *MediaDataFactory) CreateUserServices(
	coreServices CoreMediaItemServices,
	userRepos UserRepositoryFactories) UserMediaItemServices {

	return &userMediaItemServicesImpl{
		movieUserService: services.NewUserMediaItemService[*mediatypes.Movie](
			coreServices.MovieCoreService(), userRepos.MovieUserRepo()),
		seriesUserService: services.NewUserMediaItemService[*mediatypes.Series](
			coreServices.SeriesCoreService(), userRepos.SeriesUserRepo()),
		episodeUserService: services.NewUserMediaItemService[*mediatypes.Episode](
			coreServices.EpisodeCoreService(), userRepos.EpisodeUserRepo()),
		trackUserService: services.NewUserMediaItemService[*mediatypes.Track](
			coreServices.TrackCoreService(), userRepos.TrackUserRepo()),
		albumUserService: services.NewUserMediaItemService[*mediatypes.Album](
			coreServices.AlbumCoreService(), userRepos.AlbumUserRepo()),
		artistUserService: services.NewUserMediaItemService[*mediatypes.Artist](
			coreServices.ArtistCoreService(), userRepos.ArtistUserRepo()),
		collectionUserService: services.NewUserMediaItemService[*mediatypes.Collection](
			coreServices.CollectionCoreService(), userRepos.CollectionUserRepo()),
		playlistUserService: services.NewUserMediaItemService[*mediatypes.Playlist](
			coreServices.PlaylistCoreService(), userRepos.PlaylistUserRepo()),
	}
}

// --------------------------------------------------------
// Client Service Factory Methods
// --------------------------------------------------------

// CreateClientServices initializes all client services
func (f *MediaDataFactory) CreateClientServices(
	coreServices CoreMediaItemServices,
	clientRepos ClientRepositoryFactories) ClientMediaItemServices {

	return &clientMediaItemServicesImpl{
		movieClientService: services.NewClientMediaItemService[*mediatypes.Movie](
			coreServices.MovieCoreService(), clientRepos.MovieClientRepo()),
		seriesClientService: services.NewClientMediaItemService[*mediatypes.Series](
			coreServices.SeriesCoreService(), clientRepos.SeriesClientRepo()),
		episodeClientService: services.NewClientMediaItemService[*mediatypes.Episode](
			coreServices.EpisodeCoreService(), clientRepos.EpisodeClientRepo()),
		trackClientService: services.NewClientMediaItemService[*mediatypes.Track](
			coreServices.TrackCoreService(), clientRepos.TrackClientRepo()),
		albumClientService: services.NewClientMediaItemService[*mediatypes.Album](
			coreServices.AlbumCoreService(), clientRepos.AlbumClientRepo()),
		artistClientService: services.NewClientMediaItemService[*mediatypes.Artist](
			coreServices.ArtistCoreService(), clientRepos.ArtistClientRepo()),
		collectionClientService: services.NewClientMediaItemService[*mediatypes.Collection](
			coreServices.CollectionCoreService(), clientRepos.CollectionClientRepo()),
		playlistClientService: services.NewClientMediaItemService[*mediatypes.Playlist](
			coreServices.PlaylistCoreService(), clientRepos.PlaylistClientRepo()),
	}
}

// --------------------------------------------------------
// Specialized Collection Services
// --------------------------------------------------------

// CreateMediaCollectionServices creates collection and playlist services
func (f *MediaDataFactory) CreateMediaCollectionServices(
	coreServices CoreMediaItemServices,
	userServices UserMediaItemServices,
	clientServices ClientMediaItemServices,
	coreCollectionService services.CoreCollectionService,
	userCollectionService services.UserCollectionService,
	clientCollectionService services.ClientMediaCollectionService,
	playlistService services.PlaylistService) MediaCollectionServices {

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

// Define mediaCollectionServicesImpl
type mediaCollectionServicesImpl struct {
	coreCollectionService   services.CoreCollectionService
	userCollectionService   services.UserCollectionService
	clientCollectionService services.ClientMediaCollectionService

	corePlaylistService   services.CoreMediaItemService[*mediatypes.Playlist]
	userPlaylistService   services.UserMediaItemService[*mediatypes.Playlist]
	clientPlaylistService services.ClientMediaItemService[*mediatypes.Playlist]

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

func (s *mediaCollectionServicesImpl) CorePlaylistService() services.CoreMediaItemService[*mediatypes.Playlist] {
	return s.corePlaylistService
}

func (s *mediaCollectionServicesImpl) UserPlaylistService() services.UserMediaItemService[*mediatypes.Playlist] {
	return s.userPlaylistService
}

func (s *mediaCollectionServicesImpl) ClientPlaylistService() services.ClientMediaItemService[*mediatypes.Playlist] {
	return s.clientPlaylistService
}

func (s *mediaCollectionServicesImpl) PlaylistService() services.PlaylistService {
	return s.playlistService
}

// --------------------------------------------------------
// Core Handler Factory Methods
// --------------------------------------------------------

// CreateCoreHandlers initializes all core handlers
func (f *MediaDataFactory) CreateCoreHandlers(coreServices CoreMediaItemServices) CoreMediaItemHandlers {
	return &coreMediaItemHandlersImpl{
		movieCoreHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Movie](
			services.NewCoreUserMediaItemDataService[*mediatypes.Movie](coreServices.MovieCoreService())),
		seriesCoreHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Series](
			services.NewCoreUserMediaItemDataService[*mediatypes.Series](coreServices.SeriesCoreService())),
		episodeCoreHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Episode](
			services.NewCoreUserMediaItemDataService[*mediatypes.Episode](coreServices.EpisodeCoreService())),
		trackCoreHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Track](
			services.NewCoreUserMediaItemDataService[*mediatypes.Track](coreServices.TrackCoreService())),
		albumCoreHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Album](
			services.NewCoreUserMediaItemDataService[*mediatypes.Album](coreServices.AlbumCoreService())),
		artistCoreHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Artist](
			services.NewCoreUserMediaItemDataService[*mediatypes.Artist](coreServices.ArtistCoreService())),
		collectionCoreHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Collection](
			services.NewCoreUserMediaItemDataService[*mediatypes.Collection](coreServices.CollectionCoreService())),
		playlistCoreHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Playlist](
			services.NewCoreUserMediaItemDataService[*mediatypes.Playlist](coreServices.PlaylistCoreService())),
	}
}

// --------------------------------------------------------
// User Handler Factory Methods
// --------------------------------------------------------

// CreateUserHandlers initializes all user handlers
func (f *MediaDataFactory) CreateUserHandlers(
	userServices UserMediaItemServices,
	dataServices UserMediaItemDataServices,
	coreHandlers CoreMediaItemHandlers) UserMediaItemHandlers {

	// For now, we'll use the userServices directly and create simple mock handlers
	// This is a temporary solution until we properly implement the repositories

	return &userMediaItemHandlersImpl{
		movieUserHandler: handlers.NewUserUserMediaItemDataHandler[*mediatypes.Movie](
			dataServices.MovieDataService(), // Use service directly instead of creating new one
			coreHandlers.MovieCoreHandler()),
		seriesUserHandler: handlers.NewUserUserMediaItemDataHandler[*mediatypes.Series](
			dataServices.SeriesDataService(), // Use service directly instead of creating new one
			coreHandlers.SeriesCoreHandler()),
		episodeUserHandler: handlers.NewUserUserMediaItemDataHandler[*mediatypes.Episode](
			dataServices.EpisodeDataService(), // Use service directly instead of creating new one
			coreHandlers.EpisodeCoreHandler()),
		trackUserHandler: handlers.NewUserUserMediaItemDataHandler[*mediatypes.Track](
			dataServices.TrackDataService(), // Use service directly instead of creating new one
			coreHandlers.TrackCoreHandler()),
		albumUserHandler: handlers.NewUserUserMediaItemDataHandler[*mediatypes.Album](
			dataServices.AlbumDataService(), // Use service directly instead of creating new one
			coreHandlers.AlbumCoreHandler()),
		artistUserHandler: handlers.NewUserUserMediaItemDataHandler[*mediatypes.Artist](
			dataServices.ArtistDataService(), // Use service directly instead of creating new one
			coreHandlers.ArtistCoreHandler()),
		collectionUserHandler: handlers.NewUserUserMediaItemDataHandler[*mediatypes.Collection](
			dataServices.CollectionDataService(), // Use service directly instead of creating new one
			coreHandlers.CollectionCoreHandler()),
		playlistUserHandler: handlers.NewUserUserMediaItemDataHandler[*mediatypes.Playlist](
			dataServices.PlaylistDataService(), // Use service directly instead of creating new one
			coreHandlers.PlaylistCoreHandler()),
	}
}

// --------------------------------------------------------
// Client Handler Factory Methods
// --------------------------------------------------------

// CreateClientHandlers initializes all client handlers
func (f *MediaDataFactory) CreateClientHandlers(
	clientServices ClientMediaItemServices,
	dataServices ClientUserMediaItemDataServices,
	userHandlers UserMediaItemHandlers) ClientMediaItemHandlers {

	// For now, we'll use a simpler approach with just the client service and user handler
	// This is a temporary solution until we properly implement the repositories

	return &clientMediaItemHandlersImpl{
		movieClientHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Movie](
			dataServices.MovieDataService(), // Use service directly instead of creating new one
			userHandlers.MovieUserHandler()),
		seriesClientHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Series](
			dataServices.SeriesDataService(), // Use service directly instead of creating new one
			userHandlers.SeriesUserHandler()),
		episodeClientHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Episode](
			dataServices.EpisodeDataService(), // Use service directly instead of creating new one
			userHandlers.EpisodeUserHandler()),
		trackClientHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Track](
			dataServices.TrackDataService(), // Use service directly instead of creating new one
			userHandlers.TrackUserHandler()),
		albumClientHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Album](
			dataServices.AlbumDataService(), // Use service directly instead of creating new one
			userHandlers.AlbumUserHandler()),
		artistClientHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Artist](
			dataServices.ArtistDataService(), // Use service directly instead of creating new one
			userHandlers.ArtistUserHandler()),
		collectionClientHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Collection](
			dataServices.CollectionDataService(), // Use service directly instead of creating new one
			userHandlers.CollectionUserHandler()),
		playlistClientHandler: handlers.NewClientUserMediaItemDataHandler[*mediatypes.Playlist](
			dataServices.PlaylistDataService(), // Use service directly instead of creating new one
			userHandlers.PlaylistUserHandler()),
	}
}

// --------------------------------------------------------
// Specialized Media Handlers
// --------------------------------------------------------

// CreateSpecializedMediaHandlers creates specialized handlers for specific domains
// This version supports the legacy signature with the SeriesSpecificHandler type for backward compatibility
func (f *MediaDataFactory) CreateSpecializedMediaHandlers(
	coreServices CoreMediaItemServices,
	userServices UserMediaItemServices,
	clientServices ClientMediaItemServices,
	musicHandler *handlers.CoreMusicHandler,
	seriesSpecificHandler *handlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig]) SpecializedMediaHandlers {

	return &specializedMediaHandlersImpl{
		musicHandler:          musicHandler,
		seriesSpecificHandler: seriesSpecificHandler,
		seasonHandler: handlers.NewCoreUserMediaItemDataHandler[*mediatypes.Season](
			services.NewCoreUserMediaItemDataService[*mediatypes.Season](
				services.NewCoreMediaItemService[*mediatypes.Season](nil))), // Placeholder implementation
	}
}

// --------------------------------------------------------
// Repository Implementation
// --------------------------------------------------------

// Core Repository implementations
type coreRepositoriesImpl struct {
	movieRepo      repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo     repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo    repository.MediaItemRepository[*mediatypes.Episode]
	trackRepo      repository.MediaItemRepository[*mediatypes.Track]
	albumRepo      repository.MediaItemRepository[*mediatypes.Album]
	artistRepo     repository.MediaItemRepository[*mediatypes.Artist]
	collectionRepo repository.MediaItemRepository[*mediatypes.Collection]
	playlistRepo   repository.MediaItemRepository[*mediatypes.Playlist]
}

func (r *coreRepositoriesImpl) MovieRepo() repository.MediaItemRepository[*mediatypes.Movie] {
	return r.movieRepo
}

func (r *coreRepositoriesImpl) SeriesRepo() repository.MediaItemRepository[*mediatypes.Series] {
	return r.seriesRepo
}

func (r *coreRepositoriesImpl) EpisodeRepo() repository.MediaItemRepository[*mediatypes.Episode] {
	return r.episodeRepo
}

func (r *coreRepositoriesImpl) TrackRepo() repository.MediaItemRepository[*mediatypes.Track] {
	return r.trackRepo
}

func (r *coreRepositoriesImpl) AlbumRepo() repository.MediaItemRepository[*mediatypes.Album] {
	return r.albumRepo
}

func (r *coreRepositoriesImpl) ArtistRepo() repository.MediaItemRepository[*mediatypes.Artist] {
	return r.artistRepo
}

func (r *coreRepositoriesImpl) CollectionRepo() repository.MediaItemRepository[*mediatypes.Collection] {
	return r.collectionRepo
}

func (r *coreRepositoriesImpl) PlaylistRepo() repository.MediaItemRepository[*mediatypes.Playlist] {
	return r.playlistRepo
}

// User Repository implementations
type userRepositoryFactoriesImpl struct {
	movieUserRepo      repository.UserMediaItemRepository[*mediatypes.Movie]
	seriesUserRepo     repository.UserMediaItemRepository[*mediatypes.Series]
	episodeUserRepo    repository.UserMediaItemRepository[*mediatypes.Episode]
	trackUserRepo      repository.UserMediaItemRepository[*mediatypes.Track]
	albumUserRepo      repository.UserMediaItemRepository[*mediatypes.Album]
	artistUserRepo     repository.UserMediaItemRepository[*mediatypes.Artist]
	collectionUserRepo repository.UserMediaItemRepository[*mediatypes.Collection]
	playlistUserRepo   repository.UserMediaItemRepository[*mediatypes.Playlist]
}

func (r *userRepositoryFactoriesImpl) MovieUserRepo() repository.UserMediaItemRepository[*mediatypes.Movie] {
	return r.movieUserRepo
}

func (r *userRepositoryFactoriesImpl) SeriesUserRepo() repository.UserMediaItemRepository[*mediatypes.Series] {
	return r.seriesUserRepo
}

func (r *userRepositoryFactoriesImpl) EpisodeUserRepo() repository.UserMediaItemRepository[*mediatypes.Episode] {
	return r.episodeUserRepo
}

func (r *userRepositoryFactoriesImpl) TrackUserRepo() repository.UserMediaItemRepository[*mediatypes.Track] {
	return r.trackUserRepo
}

func (r *userRepositoryFactoriesImpl) AlbumUserRepo() repository.UserMediaItemRepository[*mediatypes.Album] {
	return r.albumUserRepo
}

func (r *userRepositoryFactoriesImpl) ArtistUserRepo() repository.UserMediaItemRepository[*mediatypes.Artist] {
	return r.artistUserRepo
}

func (r *userRepositoryFactoriesImpl) CollectionUserRepo() repository.UserMediaItemRepository[*mediatypes.Collection] {
	return r.collectionUserRepo
}

func (r *userRepositoryFactoriesImpl) PlaylistUserRepo() repository.UserMediaItemRepository[*mediatypes.Playlist] {
	return r.playlistUserRepo
}

// Client Repository implementations
type clientRepositoryFactoriesImpl struct {
	movieClientRepo      repository.ClientMediaItemRepository[*mediatypes.Movie]
	seriesClientRepo     repository.ClientMediaItemRepository[*mediatypes.Series]
	episodeClientRepo    repository.ClientMediaItemRepository[*mediatypes.Episode]
	trackClientRepo      repository.ClientMediaItemRepository[*mediatypes.Track]
	albumClientRepo      repository.ClientMediaItemRepository[*mediatypes.Album]
	artistClientRepo     repository.ClientMediaItemRepository[*mediatypes.Artist]
	collectionClientRepo repository.ClientMediaItemRepository[*mediatypes.Collection]
	playlistClientRepo   repository.ClientMediaItemRepository[*mediatypes.Playlist]
}

func (r *clientRepositoryFactoriesImpl) MovieClientRepo() repository.ClientMediaItemRepository[*mediatypes.Movie] {
	return r.movieClientRepo
}

func (r *clientRepositoryFactoriesImpl) SeriesClientRepo() repository.ClientMediaItemRepository[*mediatypes.Series] {
	return r.seriesClientRepo
}

func (r *clientRepositoryFactoriesImpl) EpisodeClientRepo() repository.ClientMediaItemRepository[*mediatypes.Episode] {
	return r.episodeClientRepo
}

func (r *clientRepositoryFactoriesImpl) TrackClientRepo() repository.ClientMediaItemRepository[*mediatypes.Track] {
	return r.trackClientRepo
}

func (r *clientRepositoryFactoriesImpl) AlbumClientRepo() repository.ClientMediaItemRepository[*mediatypes.Album] {
	return r.albumClientRepo
}

func (r *clientRepositoryFactoriesImpl) ArtistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Artist] {
	return r.artistClientRepo
}

func (r *clientRepositoryFactoriesImpl) CollectionClientRepo() repository.ClientMediaItemRepository[*mediatypes.Collection] {
	return r.collectionClientRepo
}

func (r *clientRepositoryFactoriesImpl) PlaylistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Playlist] {
	return r.playlistClientRepo
}

// --------------------------------------------------------
// Service Implementation
// --------------------------------------------------------

// Core Service implementation
type coreMediaItemServicesImpl struct {
	movieCoreService      services.CoreMediaItemService[*mediatypes.Movie]
	seriesCoreService     services.CoreMediaItemService[*mediatypes.Series]
	episodeCoreService    services.CoreMediaItemService[*mediatypes.Episode]
	trackCoreService      services.CoreMediaItemService[*mediatypes.Track]
	albumCoreService      services.CoreMediaItemService[*mediatypes.Album]
	artistCoreService     services.CoreMediaItemService[*mediatypes.Artist]
	collectionCoreService services.CoreMediaItemService[*mediatypes.Collection]
	playlistCoreService   services.CoreMediaItemService[*mediatypes.Playlist]
}

func (s *coreMediaItemServicesImpl) MovieCoreService() services.CoreMediaItemService[*mediatypes.Movie] {
	return s.movieCoreService
}

func (s *coreMediaItemServicesImpl) SeriesCoreService() services.CoreMediaItemService[*mediatypes.Series] {
	return s.seriesCoreService
}

func (s *coreMediaItemServicesImpl) EpisodeCoreService() services.CoreMediaItemService[*mediatypes.Episode] {
	return s.episodeCoreService
}

func (s *coreMediaItemServicesImpl) TrackCoreService() services.CoreMediaItemService[*mediatypes.Track] {
	return s.trackCoreService
}

func (s *coreMediaItemServicesImpl) AlbumCoreService() services.CoreMediaItemService[*mediatypes.Album] {
	return s.albumCoreService
}

func (s *coreMediaItemServicesImpl) ArtistCoreService() services.CoreMediaItemService[*mediatypes.Artist] {
	return s.artistCoreService
}

func (s *coreMediaItemServicesImpl) CollectionCoreService() services.CoreMediaItemService[*mediatypes.Collection] {
	return s.collectionCoreService
}

func (s *coreMediaItemServicesImpl) PlaylistCoreService() services.CoreMediaItemService[*mediatypes.Playlist] {
	return s.playlistCoreService
}

// User Service implementation
type userMediaItemServicesImpl struct {
	movieUserService      services.UserMediaItemService[*mediatypes.Movie]
	seriesUserService     services.UserMediaItemService[*mediatypes.Series]
	episodeUserService    services.UserMediaItemService[*mediatypes.Episode]
	trackUserService      services.UserMediaItemService[*mediatypes.Track]
	albumUserService      services.UserMediaItemService[*mediatypes.Album]
	artistUserService     services.UserMediaItemService[*mediatypes.Artist]
	collectionUserService services.UserMediaItemService[*mediatypes.Collection]
	playlistUserService   services.UserMediaItemService[*mediatypes.Playlist]
}

func (s *userMediaItemServicesImpl) MovieUserService() services.UserMediaItemService[*mediatypes.Movie] {
	return s.movieUserService
}

func (s *userMediaItemServicesImpl) SeriesUserService() services.UserMediaItemService[*mediatypes.Series] {
	return s.seriesUserService
}

func (s *userMediaItemServicesImpl) EpisodeUserService() services.UserMediaItemService[*mediatypes.Episode] {
	return s.episodeUserService
}

func (s *userMediaItemServicesImpl) TrackUserService() services.UserMediaItemService[*mediatypes.Track] {
	return s.trackUserService
}

func (s *userMediaItemServicesImpl) AlbumUserService() services.UserMediaItemService[*mediatypes.Album] {
	return s.albumUserService
}

func (s *userMediaItemServicesImpl) ArtistUserService() services.UserMediaItemService[*mediatypes.Artist] {
	return s.artistUserService
}

func (s *userMediaItemServicesImpl) CollectionUserService() services.UserMediaItemService[*mediatypes.Collection] {
	return s.collectionUserService
}

func (s *userMediaItemServicesImpl) PlaylistUserService() services.UserMediaItemService[*mediatypes.Playlist] {
	return s.playlistUserService
}

// Client Service implementation
type clientMediaItemServicesImpl struct {
	movieClientService      services.ClientMediaItemService[*mediatypes.Movie]
	seriesClientService     services.ClientMediaItemService[*mediatypes.Series]
	episodeClientService    services.ClientMediaItemService[*mediatypes.Episode]
	trackClientService      services.ClientMediaItemService[*mediatypes.Track]
	albumClientService      services.ClientMediaItemService[*mediatypes.Album]
	artistClientService     services.ClientMediaItemService[*mediatypes.Artist]
	collectionClientService services.ClientMediaItemService[*mediatypes.Collection]
	playlistClientService   services.ClientMediaItemService[*mediatypes.Playlist]
}

func (s *clientMediaItemServicesImpl) MovieClientService() services.ClientMediaItemService[*mediatypes.Movie] {
	return s.movieClientService
}

func (s *clientMediaItemServicesImpl) SeriesClientService() services.ClientMediaItemService[*mediatypes.Series] {
	return s.seriesClientService
}

func (s *clientMediaItemServicesImpl) EpisodeClientService() services.ClientMediaItemService[*mediatypes.Episode] {
	return s.episodeClientService
}

func (s *clientMediaItemServicesImpl) TrackClientService() services.ClientMediaItemService[*mediatypes.Track] {
	return s.trackClientService
}

func (s *clientMediaItemServicesImpl) AlbumClientService() services.ClientMediaItemService[*mediatypes.Album] {
	return s.albumClientService
}

func (s *clientMediaItemServicesImpl) ArtistClientService() services.ClientMediaItemService[*mediatypes.Artist] {
	return s.artistClientService
}

func (s *clientMediaItemServicesImpl) CollectionClientService() services.ClientMediaItemService[*mediatypes.Collection] {
	return s.collectionClientService
}

func (s *clientMediaItemServicesImpl) PlaylistClientService() services.ClientMediaItemService[*mediatypes.Playlist] {
	return s.playlistClientService
}

// --------------------------------------------------------
// Handler Implementation
// --------------------------------------------------------

// Core handler implementation
type coreMediaItemHandlersImpl struct {
	movieCoreHandler      *handlers.CoreUserMediaItemDataHandler[*mediatypes.Movie]
	seriesCoreHandler     *handlers.CoreUserMediaItemDataHandler[*mediatypes.Series]
	episodeCoreHandler    *handlers.CoreUserMediaItemDataHandler[*mediatypes.Episode]
	trackCoreHandler      *handlers.CoreUserMediaItemDataHandler[*mediatypes.Track]
	albumCoreHandler      *handlers.CoreUserMediaItemDataHandler[*mediatypes.Album]
	artistCoreHandler     *handlers.CoreUserMediaItemDataHandler[*mediatypes.Artist]
	collectionCoreHandler *handlers.CoreUserMediaItemDataHandler[*mediatypes.Collection]
	playlistCoreHandler   *handlers.CoreUserMediaItemDataHandler[*mediatypes.Playlist]
}

func (h *coreMediaItemHandlersImpl) MovieCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Movie] {
	return h.movieCoreHandler
}

func (h *coreMediaItemHandlersImpl) SeriesCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Series] {
	return h.seriesCoreHandler
}

func (h *coreMediaItemHandlersImpl) EpisodeCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Episode] {
	return h.episodeCoreHandler
}

func (h *coreMediaItemHandlersImpl) TrackCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Track] {
	return h.trackCoreHandler
}

func (h *coreMediaItemHandlersImpl) AlbumCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Album] {
	return h.albumCoreHandler
}

func (h *coreMediaItemHandlersImpl) ArtistCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Artist] {
	return h.artistCoreHandler
}

func (h *coreMediaItemHandlersImpl) CollectionCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Collection] {
	return h.collectionCoreHandler
}

func (h *coreMediaItemHandlersImpl) PlaylistCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Playlist] {
	return h.playlistCoreHandler
}

// User handler implementation
type userMediaItemHandlersImpl struct {
	movieUserHandler      *handlers.UserUserMediaItemDataHandler[*mediatypes.Movie]
	seriesUserHandler     *handlers.UserUserMediaItemDataHandler[*mediatypes.Series]
	episodeUserHandler    *handlers.UserUserMediaItemDataHandler[*mediatypes.Episode]
	trackUserHandler      *handlers.UserUserMediaItemDataHandler[*mediatypes.Track]
	albumUserHandler      *handlers.UserUserMediaItemDataHandler[*mediatypes.Album]
	artistUserHandler     *handlers.UserUserMediaItemDataHandler[*mediatypes.Artist]
	collectionUserHandler *handlers.UserUserMediaItemDataHandler[*mediatypes.Collection]
	playlistUserHandler   *handlers.UserUserMediaItemDataHandler[*mediatypes.Playlist]
}

func (h *userMediaItemHandlersImpl) MovieUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Movie] {
	return h.movieUserHandler
}

func (h *userMediaItemHandlersImpl) SeriesUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Series] {
	return h.seriesUserHandler
}

func (h *userMediaItemHandlersImpl) EpisodeUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Episode] {
	return h.episodeUserHandler
}

func (h *userMediaItemHandlersImpl) TrackUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Track] {
	return h.trackUserHandler
}

func (h *userMediaItemHandlersImpl) AlbumUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Album] {
	return h.albumUserHandler
}

func (h *userMediaItemHandlersImpl) ArtistUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Artist] {
	return h.artistUserHandler
}

func (h *userMediaItemHandlersImpl) CollectionUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Collection] {
	return h.collectionUserHandler
}

func (h *userMediaItemHandlersImpl) PlaylistUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Playlist] {
	return h.playlistUserHandler
}

// Client handler implementation
type clientMediaItemHandlersImpl struct {
	movieClientHandler      *handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie]
	seriesClientHandler     *handlers.ClientUserMediaItemDataHandler[*mediatypes.Series]
	episodeClientHandler    *handlers.ClientUserMediaItemDataHandler[*mediatypes.Episode]
	trackClientHandler      *handlers.ClientUserMediaItemDataHandler[*mediatypes.Track]
	albumClientHandler      *handlers.ClientUserMediaItemDataHandler[*mediatypes.Album]
	artistClientHandler     *handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist]
	collectionClientHandler *handlers.ClientUserMediaItemDataHandler[*mediatypes.Collection]
	playlistClientHandler   *handlers.ClientUserMediaItemDataHandler[*mediatypes.Playlist]
}

func (h *clientMediaItemHandlersImpl) MovieClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie] {
	return h.movieClientHandler
}

func (h *clientMediaItemHandlersImpl) SeriesClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Series] {
	return h.seriesClientHandler
}

func (h *clientMediaItemHandlersImpl) EpisodeClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Episode] {
	return h.episodeClientHandler
}

func (h *clientMediaItemHandlersImpl) TrackClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Track] {
	return h.trackClientHandler
}

func (h *clientMediaItemHandlersImpl) AlbumClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Album] {
	return h.albumClientHandler
}

func (h *clientMediaItemHandlersImpl) ArtistClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist] {
	return h.artistClientHandler
}

func (h *clientMediaItemHandlersImpl) CollectionClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Collection] {
	return h.collectionClientHandler
}

func (h *clientMediaItemHandlersImpl) PlaylistClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Playlist] {
	return h.playlistClientHandler
}

// Specialized media handlers implementation
type specializedMediaHandlersImpl struct {
	// Domain-specific handlers
	musicHandler          *handlers.CoreMusicHandler
	seriesSpecificHandler *handlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig]

	// Special case handlers
	seasonHandler *handlers.CoreUserMediaItemDataHandler[*mediatypes.Season]
}

// Domain-specific handler access methods
func (h *specializedMediaHandlersImpl) MusicHandler() *handlers.CoreMusicHandler {
	return h.musicHandler
}

func (h *specializedMediaHandlersImpl) SeriesSpecificHandler() *handlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig] {
	return h.seriesSpecificHandler
}

// Special case handler access methods
func (h *specializedMediaHandlersImpl) SeasonHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Season] {
	return h.seasonHandler
}

type coreUserMediaItemDataRepositoriesImpl struct {
	movieDataRepo      repository.CoreUserMediaItemDataRepository[*mediatypes.Movie]
	seriesDataRepo     repository.CoreUserMediaItemDataRepository[*mediatypes.Series]
	episodeDataRepo    repository.CoreUserMediaItemDataRepository[*mediatypes.Episode]
	trackDataRepo      repository.CoreUserMediaItemDataRepository[*mediatypes.Track]
	albumDataRepo      repository.CoreUserMediaItemDataRepository[*mediatypes.Album]
	artistDataRepo     repository.CoreUserMediaItemDataRepository[*mediatypes.Artist]
	collectionDataRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Collection]
	playlistDataRepo   repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist]
}

func (r *coreUserMediaItemDataRepositoriesImpl) MovieDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Movie] {
	return r.movieDataRepo
}

func (r *coreUserMediaItemDataRepositoriesImpl) SeriesDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Series] {
	return r.seriesDataRepo
}

func (r *coreUserMediaItemDataRepositoriesImpl) EpisodeDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Episode] {
	return r.episodeDataRepo
}

func (r *coreUserMediaItemDataRepositoriesImpl) TrackDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Track] {
	return r.trackDataRepo
}

func (r *coreUserMediaItemDataRepositoriesImpl) AlbumDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Album] {
	return r.albumDataRepo
}

func (r *coreUserMediaItemDataRepositoriesImpl) ArtistDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Artist] {
	return r.artistDataRepo
}

func (r *coreUserMediaItemDataRepositoriesImpl) CollectionDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Collection] {
	return r.collectionDataRepo
}

func (r *coreUserMediaItemDataRepositoriesImpl) PlaylistDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist] {
	return r.playlistDataRepo
}

type userDataRepositoriesImpl struct {
	movieDataRepo      repository.UserMediaItemDataRepository[*mediatypes.Movie]
	seriesDataRepo     repository.UserMediaItemDataRepository[*mediatypes.Series]
	episodeDataRepo    repository.UserMediaItemDataRepository[*mediatypes.Episode]
	trackDataRepo      repository.UserMediaItemDataRepository[*mediatypes.Track]
	albumDataRepo      repository.UserMediaItemDataRepository[*mediatypes.Album]
	artistDataRepo     repository.UserMediaItemDataRepository[*mediatypes.Artist]
	collectionDataRepo repository.UserMediaItemDataRepository[*mediatypes.Collection]
	playlistDataRepo   repository.UserMediaItemDataRepository[*mediatypes.Playlist]
}

func (r *userDataRepositoriesImpl) MovieDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Movie] {
	return r.movieDataRepo
}

func (r *userDataRepositoriesImpl) SeriesDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Series] {
	return r.seriesDataRepo
}

func (r *userDataRepositoriesImpl) EpisodeDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Episode] {
	return r.episodeDataRepo
}

func (r *userDataRepositoriesImpl) TrackDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Track] {
	return r.trackDataRepo
}

func (r *userDataRepositoriesImpl) AlbumDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Album] {
	return r.albumDataRepo
}

func (r *userDataRepositoriesImpl) ArtistDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Artist] {
	return r.artistDataRepo
}

func (r *userDataRepositoriesImpl) CollectionDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Collection] {
	return r.collectionDataRepo
}

func (r *userDataRepositoriesImpl) PlaylistDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Playlist] {
	return r.playlistDataRepo
}
