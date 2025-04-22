package providers

import (
	"github.com/google/wire"
	"suasor/client/media/types"
	"suasor/repository"
	"suasor/services"
)

// --- Core Media Item Services ---

func ProvideCoreMediaItemService[T types.MediaData](
	repo repository.MediaItemRepository[T],
) services.CoreMediaItemService[T] {
	return services.NewCoreMediaItemService[T](repo)
}

func ProvideCoreMovieService(
	repo repository.MediaItemRepository[*types.Movie],
) services.CoreMediaItemService[*types.Movie] {
	return services.NewCoreMediaItemService[*types.Movie](repo)
}

func ProvideCoreSeriesService(
	repo repository.MediaItemRepository[*types.Series],
) services.CoreMediaItemService[*types.Series] {
	return services.NewCoreMediaItemService[*types.Series](repo)
}

func ProvideCoreSeasonService(
	repo repository.MediaItemRepository[*types.Season],
) services.CoreMediaItemService[*types.Season] {
	return services.NewCoreMediaItemService[*types.Season](repo)
}

func ProvideCoreEpisodeService(
	repo repository.MediaItemRepository[*types.Episode],
) services.CoreMediaItemService[*types.Episode] {
	return services.NewCoreMediaItemService[*types.Episode](repo)
}

func ProvideCoreTrackService(
	repo repository.MediaItemRepository[*types.Track],
) services.CoreMediaItemService[*types.Track] {
	return services.NewCoreMediaItemService[*types.Track](repo)
}

func ProvideCoreAlbumService(
	repo repository.MediaItemRepository[*types.Album],
) services.CoreMediaItemService[*types.Album] {
	return services.NewCoreMediaItemService[*types.Album](repo)
}

func ProvideCoreArtistService(
	repo repository.MediaItemRepository[*types.Artist],
) services.CoreMediaItemService[*types.Artist] {
	return services.NewCoreMediaItemService[*types.Artist](repo)
}

func ProvideCorePlaylistService(
	repo repository.MediaItemRepository[*types.Playlist],
) services.CoreMediaItemService[*types.Playlist] {
	return services.NewCoreMediaItemService[*types.Playlist](repo)
}

func ProvideCoreCollectionService(
	repo repository.MediaItemRepository[*types.Collection],
) services.CoreMediaItemService[*types.Collection] {
	return services.NewCoreMediaItemService[*types.Collection](repo)
}

// --- Core List Services ---

func ProvideCoreListService[T types.ListData](
	repo repository.MediaItemRepository[T],
) services.CoreListService[T] {
	return services.NewCoreListService[T](repo)
}

func ProvideCorePlaylistListService(
	repo repository.MediaItemRepository[*types.Playlist],
) services.CoreListService[*types.Playlist] {
	return services.NewCoreListService[*types.Playlist](repo)
}

func ProvideCoreCollectionListService(
	repo repository.MediaItemRepository[*types.Collection],
) services.CoreListService[*types.Collection] {
	return services.NewCoreListService[*types.Collection](repo)
}

// --- User Media Item Services ---

func ProvideUserMediaItemService[T types.MediaData](
	coreService services.CoreMediaItemService[T],
	userRepo repository.UserMediaItemRepository[T],
) services.UserMediaItemService[T] {
	return services.NewUserMediaItemService[T](coreService, userRepo)
}

func ProvideUserMovieService(
	coreService services.CoreMediaItemService[*types.Movie],
	userRepo repository.UserMediaItemRepository[*types.Movie],
) services.UserMediaItemService[*types.Movie] {
	return services.NewUserMediaItemService[*types.Movie](coreService, userRepo)
}

func ProvideUserSeriesService(
	coreService services.CoreMediaItemService[*types.Series],
	userRepo repository.UserMediaItemRepository[*types.Series],
) services.UserMediaItemService[*types.Series] {
	return services.NewUserMediaItemService[*types.Series](coreService, userRepo)
}

func ProvideUserSeasonService(
	coreService services.CoreMediaItemService[*types.Season],
	userRepo repository.UserMediaItemRepository[*types.Season],
) services.UserMediaItemService[*types.Season] {
	return services.NewUserMediaItemService[*types.Season](coreService, userRepo)
}

func ProvideUserEpisodeService(
	coreService services.CoreMediaItemService[*types.Episode],
	userRepo repository.UserMediaItemRepository[*types.Episode],
) services.UserMediaItemService[*types.Episode] {
	return services.NewUserMediaItemService[*types.Episode](coreService, userRepo)
}

func ProvideUserTrackService(
	coreService services.CoreMediaItemService[*types.Track],
	userRepo repository.UserMediaItemRepository[*types.Track],
) services.UserMediaItemService[*types.Track] {
	return services.NewUserMediaItemService[*types.Track](coreService, userRepo)
}

func ProvideUserAlbumService(
	coreService services.CoreMediaItemService[*types.Album],
	userRepo repository.UserMediaItemRepository[*types.Album],
) services.UserMediaItemService[*types.Album] {
	return services.NewUserMediaItemService[*types.Album](coreService, userRepo)
}

func ProvideUserArtistService(
	coreService services.CoreMediaItemService[*types.Artist],
	userRepo repository.UserMediaItemRepository[*types.Artist],
) services.UserMediaItemService[*types.Artist] {
	return services.NewUserMediaItemService[*types.Artist](coreService, userRepo)
}

func ProvideUserPlaylistService(
	coreService services.CoreMediaItemService[*types.Playlist],
	userRepo repository.UserMediaItemRepository[*types.Playlist],
) services.UserMediaItemService[*types.Playlist] {
	return services.NewUserMediaItemService[*types.Playlist](coreService, userRepo)
}

func ProvideUserCollectionService(
	coreService services.CoreMediaItemService[*types.Collection],
	userRepo repository.UserMediaItemRepository[*types.Collection],
) services.UserMediaItemService[*types.Collection] {
	return services.NewUserMediaItemService[*types.Collection](coreService, userRepo)
}

// --- User List Services ---

func ProvideUserListService[T types.ListData](
	coreService services.CoreListService[T],
	userRepo repository.UserMediaItemRepository[T],
	userDataRepo repository.UserMediaItemDataRepository[T],
) services.UserListService[T] {
	return services.NewUserListService[T](coreService, userRepo, userDataRepo)
}

func ProvideUserPlaylistListService(
	coreService services.CoreListService[*types.Playlist],
	userRepo repository.UserMediaItemRepository[*types.Playlist],
	userDataRepo repository.UserMediaItemDataRepository[*types.Playlist],
) services.UserListService[*types.Playlist] {
	return services.NewUserListService[*types.Playlist](coreService, userRepo, userDataRepo)
}

func ProvideUserCollectionListService(
	coreService services.CoreListService[*types.Collection],
	userRepo repository.UserMediaItemRepository[*types.Collection],
	userDataRepo repository.UserMediaItemDataRepository[*types.Collection],
) services.UserListService[*types.Collection] {
	return services.NewUserListService[*types.Collection](coreService, userRepo, userDataRepo)
}

// --- Core Media Item Data Services ---

func ProvideCoreUserMediaItemDataService[T types.MediaData](
	coreService services.CoreMediaItemService[T],
	coreRepo repository.CoreUserMediaItemDataRepository[T],
) services.CoreUserMediaItemDataService[T] {
	return services.NewCoreUserMediaItemDataService[T](coreService, coreRepo)
}

func ProvideCoreMovieDataService(
	coreService services.CoreMediaItemService[*types.Movie],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Movie],
) services.CoreUserMediaItemDataService[*types.Movie] {
	return services.NewCoreUserMediaItemDataService[*types.Movie](coreService, coreRepo)
}

func ProvideCoreSeriesDataService(
	coreService services.CoreMediaItemService[*types.Series],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Series],
) services.CoreUserMediaItemDataService[*types.Series] {
	return services.NewCoreUserMediaItemDataService[*types.Series](coreService, coreRepo)
}

func ProvideCoreSeasonDataService(
	coreService services.CoreMediaItemService[*types.Season],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Season],
) services.CoreUserMediaItemDataService[*types.Season] {
	return services.NewCoreUserMediaItemDataService[*types.Season](coreService, coreRepo)
}

func ProvideCoreEpisodeDataService(
	coreService services.CoreMediaItemService[*types.Episode],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Episode],
) services.CoreUserMediaItemDataService[*types.Episode] {
	return services.NewCoreUserMediaItemDataService[*types.Episode](coreService, coreRepo)
}

func ProvideCoreTrackDataService(
	coreService services.CoreMediaItemService[*types.Track],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Track],
) services.CoreUserMediaItemDataService[*types.Track] {
	return services.NewCoreUserMediaItemDataService[*types.Track](coreService, coreRepo)
}

func ProvideCoreAlbumDataService(
	coreService services.CoreMediaItemService[*types.Album],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Album],
) services.CoreUserMediaItemDataService[*types.Album] {
	return services.NewCoreUserMediaItemDataService[*types.Album](coreService, coreRepo)
}

func ProvideCoreArtistDataService(
	coreService services.CoreMediaItemService[*types.Artist],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Artist],
) services.CoreUserMediaItemDataService[*types.Artist] {
	return services.NewCoreUserMediaItemDataService[*types.Artist](coreService, coreRepo)
}

func ProvideCorePlaylistDataService(
	coreService services.CoreMediaItemService[*types.Playlist],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Playlist],
) services.CoreUserMediaItemDataService[*types.Playlist] {
	return services.NewCoreUserMediaItemDataService[*types.Playlist](coreService, coreRepo)
}

func ProvideCoreCollectionDataService(
	coreService services.CoreMediaItemService[*types.Collection],
	coreRepo repository.CoreUserMediaItemDataRepository[*types.Collection],
) services.CoreUserMediaItemDataService[*types.Collection] {
	return services.NewCoreUserMediaItemDataService[*types.Collection](coreService, coreRepo)
}

// --- User Media Item Data Services ---

func ProvideUserMediaItemDataService[T types.MediaData](
	coreService services.CoreUserMediaItemDataService[T],
	userRepo repository.UserMediaItemDataRepository[T],
) services.UserMediaItemDataService[T] {
	return services.NewUserMediaItemDataService[T](coreService, userRepo)
}

func ProvideUserMovieDataService(
	coreService services.CoreUserMediaItemDataService[*types.Movie],
	userRepo repository.UserMediaItemDataRepository[*types.Movie],
) services.UserMediaItemDataService[*types.Movie] {
	return services.NewUserMediaItemDataService[*types.Movie](coreService, userRepo)
}

func ProvideUserSeriesDataService(
	coreService services.CoreUserMediaItemDataService[*types.Series],
	userRepo repository.UserMediaItemDataRepository[*types.Series],
) services.UserMediaItemDataService[*types.Series] {
	return services.NewUserMediaItemDataService[*types.Series](coreService, userRepo)
}

func ProvideUserSeasonDataService(
	coreService services.CoreUserMediaItemDataService[*types.Season],
	userRepo repository.UserMediaItemDataRepository[*types.Season],
) services.UserMediaItemDataService[*types.Season] {
	return services.NewUserMediaItemDataService[*types.Season](coreService, userRepo)
}

func ProvideUserEpisodeDataService(
	coreService services.CoreUserMediaItemDataService[*types.Episode],
	userRepo repository.UserMediaItemDataRepository[*types.Episode],
) services.UserMediaItemDataService[*types.Episode] {
	return services.NewUserMediaItemDataService[*types.Episode](coreService, userRepo)
}

func ProvideUserTrackDataService(
	coreService services.CoreUserMediaItemDataService[*types.Track],
	userRepo repository.UserMediaItemDataRepository[*types.Track],
) services.UserMediaItemDataService[*types.Track] {
	return services.NewUserMediaItemDataService[*types.Track](coreService, userRepo)
}

func ProvideUserAlbumDataService(
	coreService services.CoreUserMediaItemDataService[*types.Album],
	userRepo repository.UserMediaItemDataRepository[*types.Album],
) services.UserMediaItemDataService[*types.Album] {
	return services.NewUserMediaItemDataService[*types.Album](coreService, userRepo)
}

func ProvideUserArtistDataService(
	coreService services.CoreUserMediaItemDataService[*types.Artist],
	userRepo repository.UserMediaItemDataRepository[*types.Artist],
) services.UserMediaItemDataService[*types.Artist] {
	return services.NewUserMediaItemDataService[*types.Artist](coreService, userRepo)
}

func ProvideUserPlaylistDataService(
	coreService services.CoreUserMediaItemDataService[*types.Playlist],
	userRepo repository.UserMediaItemDataRepository[*types.Playlist],
) services.UserMediaItemDataService[*types.Playlist] {
	return services.NewUserMediaItemDataService[*types.Playlist](coreService, userRepo)
}

func ProvideUserCollectionDataService(
	coreService services.CoreUserMediaItemDataService[*types.Collection],
	userRepo repository.UserMediaItemDataRepository[*types.Collection],
) services.UserMediaItemDataService[*types.Collection] {
	return services.NewUserMediaItemDataService[*types.Collection](coreService, userRepo)
}

// --- Service Provider Sets ---

// CoreServiceSet is a provider set for core services
var CoreServiceSet = wire.NewSet(
	ProvideCoreMovieService,
	ProvideCoreSeriesService,
	ProvideCoreSeasonService,
	ProvideCoreEpisodeService,
	ProvideCoreTrackService,
	ProvideCoreAlbumService,
	ProvideCoreArtistService,
	ProvideCorePlaylistService,
	ProvideCoreCollectionService,
)

// CoreListServiceSet is a provider set for core list services
var CoreListServiceSet = wire.NewSet(
	ProvideCorePlaylistListService,
	ProvideCoreCollectionListService,
)

// UserServiceSet is a provider set for user services
var UserServiceSet = wire.NewSet(
	ProvideUserMovieService,
	ProvideUserSeriesService,
	ProvideUserSeasonService,
	ProvideUserEpisodeService,
	ProvideUserTrackService,
	ProvideUserAlbumService,
	ProvideUserArtistService,
	ProvideUserPlaylistService,
	ProvideUserCollectionService,
)

// UserListServiceSet is a provider set for user list services
var UserListServiceSet = wire.NewSet(
	ProvideUserPlaylistListService,
	ProvideUserCollectionListService,
)

// CoreDataServiceSet is a provider set for core data services
var CoreDataServiceSet = wire.NewSet(
	ProvideCoreMovieDataService,
	ProvideCoreSeriesDataService,
	ProvideCoreSeasonDataService,
	ProvideCoreEpisodeDataService,
	ProvideCoreTrackDataService,
	ProvideCoreAlbumDataService,
	ProvideCoreArtistDataService,
	ProvideCorePlaylistDataService,
	ProvideCoreCollectionDataService,
)

// UserDataServiceSet is a provider set for user data services
var UserDataServiceSet = wire.NewSet(
	ProvideUserMovieDataService,
	ProvideUserSeriesDataService,
	ProvideUserSeasonDataService,
	ProvideUserEpisodeDataService,
	ProvideUserTrackDataService,
	ProvideUserAlbumDataService,
	ProvideUserArtistDataService,
	ProvideUserPlaylistDataService,
	ProvideUserCollectionDataService,
)

// ServiceSet combines all service provider sets
var ServiceSet = wire.NewSet(
	CoreServiceSet,
	CoreListServiceSet,
	UserServiceSet,
	UserListServiceSet,
	CoreDataServiceSet,
	UserDataServiceSet,
)

