package providers

import (
	"github.com/google/wire"
	"suasor/client/media/types"
	"suasor/handlers"
	"suasor/services"
)

// --- Core Media Item Handlers ---

func ProvideCoreMediaItemHandler[T types.MediaData](
	service services.CoreMediaItemService[T],
) handlers.CoreMediaItemHandler[T] {
	return handlers.NewCoreMediaItemHandler[T](service)
}

func ProvideCoreMovieHandler(
	service services.CoreMediaItemService[*types.Movie],
) handlers.CoreMediaItemHandler[*types.Movie] {
	return handlers.NewCoreMediaItemHandler[*types.Movie](service)
}

func ProvideCoreSeriesHandler(
	service services.CoreMediaItemService[*types.Series],
) handlers.CoreMediaItemHandler[*types.Series] {
	return handlers.NewCoreMediaItemHandler[*types.Series](service)
}

func ProvideCoreEpisodeHandler(
	service services.CoreMediaItemService[*types.Episode],
) handlers.CoreMediaItemHandler[*types.Episode] {
	return handlers.NewCoreMediaItemHandler[*types.Episode](service)
}

func ProvideCoreSeasonHandler(
	service services.CoreMediaItemService[*types.Season],
) handlers.CoreMediaItemHandler[*types.Season] {
	return handlers.NewCoreMediaItemHandler[*types.Season](service)
}

func ProvideCoreTrackHandler(
	service services.CoreMediaItemService[*types.Track],
) handlers.CoreMediaItemHandler[*types.Track] {
	return handlers.NewCoreMediaItemHandler[*types.Track](service)
}

func ProvideCoreAlbumHandler(
	service services.CoreMediaItemService[*types.Album],
) handlers.CoreMediaItemHandler[*types.Album] {
	return handlers.NewCoreMediaItemHandler[*types.Album](service)
}

func ProvideCoreArtistHandler(
	service services.CoreMediaItemService[*types.Artist],
) handlers.CoreMediaItemHandler[*types.Artist] {
	return handlers.NewCoreMediaItemHandler[*types.Artist](service)
}

func ProvideCorePlaylistHandler(
	service services.CoreMediaItemService[*types.Playlist],
) handlers.CoreMediaItemHandler[*types.Playlist] {
	return handlers.NewCoreMediaItemHandler[*types.Playlist](service)
}

func ProvideCoreCollectionHandler(
	service services.CoreMediaItemService[*types.Collection],
) handlers.CoreMediaItemHandler[*types.Collection] {
	return handlers.NewCoreMediaItemHandler[*types.Collection](service)
}

// --- Core List Handlers ---

func ProvideCoreListHandler[T types.ListData](
	coreHandler handlers.CoreMediaItemHandler[T],
	listService services.CoreListService[T],
) handlers.CoreListHandler[T] {
	return handlers.NewCoreListHandler[T](coreHandler, listService)
}

func ProvideCorePlaylistListHandler(
	coreHandler handlers.CoreMediaItemHandler[*types.Playlist],
	listService services.CoreListService[*types.Playlist],
) handlers.CoreListHandler[*types.Playlist] {
	return handlers.NewCoreListHandler[*types.Playlist](coreHandler, listService)
}

func ProvideCoreCollectionListHandler(
	coreHandler handlers.CoreMediaItemHandler[*types.Collection],
	listService services.CoreListService[*types.Collection],
) handlers.CoreListHandler[*types.Collection] {
	return handlers.NewCoreListHandler[*types.Collection](coreHandler, listService)
}

// --- User Media Item Handlers ---

func ProvideUserMediaItemHandler[T types.MediaData](
	service services.UserMediaItemService[T],
) handlers.UserMediaItemHandler[T] {
	return handlers.NewUserMediaItemHandler[T](service)
}

func ProvideUserMovieHandler(
	service services.UserMediaItemService[*types.Movie],
) handlers.UserMediaItemHandler[*types.Movie] {
	return handlers.NewUserMediaItemHandler[*types.Movie](service)
}

func ProvideUserSeriesHandler(
	service services.UserMediaItemService[*types.Series],
) handlers.UserMediaItemHandler[*types.Series] {
	return handlers.NewUserMediaItemHandler[*types.Series](service)
}

func ProvideUserEpisodeHandler(
	service services.UserMediaItemService[*types.Episode],
) handlers.UserMediaItemHandler[*types.Episode] {
	return handlers.NewUserMediaItemHandler[*types.Episode](service)
}

func ProvideUserSeasonHandler(
	service services.UserMediaItemService[*types.Season],
) handlers.UserMediaItemHandler[*types.Season] {
	return handlers.NewUserMediaItemHandler[*types.Season](service)
}

func ProvideUserTrackHandler(
	service services.UserMediaItemService[*types.Track],
) handlers.UserMediaItemHandler[*types.Track] {
	return handlers.NewUserMediaItemHandler[*types.Track](service)
}

func ProvideUserAlbumHandler(
	service services.UserMediaItemService[*types.Album],
) handlers.UserMediaItemHandler[*types.Album] {
	return handlers.NewUserMediaItemHandler[*types.Album](service)
}

func ProvideUserArtistHandler(
	service services.UserMediaItemService[*types.Artist],
) handlers.UserMediaItemHandler[*types.Artist] {
	return handlers.NewUserMediaItemHandler[*types.Artist](service)
}

func ProvideUserPlaylistHandler(
	service services.UserMediaItemService[*types.Playlist],
) handlers.UserMediaItemHandler[*types.Playlist] {
	return handlers.NewUserMediaItemHandler[*types.Playlist](service)
}

func ProvideUserCollectionHandler(
	service services.UserMediaItemService[*types.Collection],
) handlers.UserMediaItemHandler[*types.Collection] {
	return handlers.NewUserMediaItemHandler[*types.Collection](service)
}

// --- User List Handlers ---

func ProvideUserListHandler[T types.ListData](
	coreHandler handlers.CoreListHandler[T],
	itemService services.UserMediaItemService[T],
	listService services.UserListService[T],
) handlers.UserListHandler[T] {
	return handlers.NewUserListHandler[T](coreHandler, itemService, listService)
}

func ProvideUserPlaylistListHandler(
	coreHandler handlers.CoreListHandler[*types.Playlist],
	itemService services.UserMediaItemService[*types.Playlist],
	listService services.UserListService[*types.Playlist],
) handlers.UserListHandler[*types.Playlist] {
	return handlers.NewUserListHandler[*types.Playlist](coreHandler, itemService, listService)
}

func ProvideUserCollectionListHandler(
	coreHandler handlers.CoreListHandler[*types.Collection],
	itemService services.UserMediaItemService[*types.Collection],
	listService services.UserListService[*types.Collection],
) handlers.UserListHandler[*types.Collection] {
	return handlers.NewUserListHandler[*types.Collection](coreHandler, itemService, listService)
}

// --- Core User Media Item Data Handlers ---

func ProvideCoreUserMediaItemDataHandler[T types.MediaData](
	service services.CoreUserMediaItemDataService[T],
) handlers.CoreUserMediaItemDataHandler[T] {
	return handlers.NewCoreUserMediaItemDataHandler[T](service)
}

func ProvideCoreMovieDataHandler(
	service services.CoreUserMediaItemDataService[*types.Movie],
) handlers.CoreUserMediaItemDataHandler[*types.Movie] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Movie](service)
}

func ProvideCoreSeriesDataHandler(
	service services.CoreUserMediaItemDataService[*types.Series],
) handlers.CoreUserMediaItemDataHandler[*types.Series] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Series](service)
}

func ProvideCoreSeasonDataHandler(
	service services.CoreUserMediaItemDataService[*types.Season],
) handlers.CoreUserMediaItemDataHandler[*types.Season] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Season](service)
}

func ProvideCoreEpisodeDataHandler(
	service services.CoreUserMediaItemDataService[*types.Episode],
) handlers.CoreUserMediaItemDataHandler[*types.Episode] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Episode](service)
}

func ProvideCoreTrackDataHandler(
	service services.CoreUserMediaItemDataService[*types.Track],
) handlers.CoreUserMediaItemDataHandler[*types.Track] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Track](service)
}

func ProvideCoreAlbumDataHandler(
	service services.CoreUserMediaItemDataService[*types.Album],
) handlers.CoreUserMediaItemDataHandler[*types.Album] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Album](service)
}

func ProvideCoreArtistDataHandler(
	service services.CoreUserMediaItemDataService[*types.Artist],
) handlers.CoreUserMediaItemDataHandler[*types.Artist] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Artist](service)
}

func ProvideCorePlaylistDataHandler(
	service services.CoreUserMediaItemDataService[*types.Playlist],
) handlers.CoreUserMediaItemDataHandler[*types.Playlist] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Playlist](service)
}

func ProvideCoreCollectionDataHandler(
	service services.CoreUserMediaItemDataService[*types.Collection],
) handlers.CoreUserMediaItemDataHandler[*types.Collection] {
	return handlers.NewCoreUserMediaItemDataHandler[*types.Collection](service)
}

// --- User Media Item Data Handlers ---

func ProvideUserMediaItemDataHandler[T types.MediaData](
	coreHandler handlers.CoreUserMediaItemDataHandler[T],
	service services.UserMediaItemDataService[T],
) handlers.UserMediaItemDataHandler[T] {
	return handlers.NewUserMediaItemDataHandler[T](coreHandler, service)
}

func ProvideUserMovieDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Movie],
	service services.UserMediaItemDataService[*types.Movie],
) handlers.UserMediaItemDataHandler[*types.Movie] {
	return handlers.NewUserMediaItemDataHandler[*types.Movie](coreHandler, service)
}

func ProvideUserSeriesDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Series],
	service services.UserMediaItemDataService[*types.Series],
) handlers.UserMediaItemDataHandler[*types.Series] {
	return handlers.NewUserMediaItemDataHandler[*types.Series](coreHandler, service)
}

func ProvideUserSeasonDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Season],
	service services.UserMediaItemDataService[*types.Season],
) handlers.UserMediaItemDataHandler[*types.Season] {
	return handlers.NewUserMediaItemDataHandler[*types.Season](coreHandler, service)
}

func ProvideUserEpisodeDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Episode],
	service services.UserMediaItemDataService[*types.Episode],
) handlers.UserMediaItemDataHandler[*types.Episode] {
	return handlers.NewUserMediaItemDataHandler[*types.Episode](coreHandler, service)
}

func ProvideUserTrackDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Track],
	service services.UserMediaItemDataService[*types.Track],
) handlers.UserMediaItemDataHandler[*types.Track] {
	return handlers.NewUserMediaItemDataHandler[*types.Track](coreHandler, service)
}

func ProvideUserAlbumDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Album],
	service services.UserMediaItemDataService[*types.Album],
) handlers.UserMediaItemDataHandler[*types.Album] {
	return handlers.NewUserMediaItemDataHandler[*types.Album](coreHandler, service)
}

func ProvideUserArtistDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Artist],
	service services.UserMediaItemDataService[*types.Artist],
) handlers.UserMediaItemDataHandler[*types.Artist] {
	return handlers.NewUserMediaItemDataHandler[*types.Artist](coreHandler, service)
}

func ProvideUserPlaylistDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Playlist],
	service services.UserMediaItemDataService[*types.Playlist],
) handlers.UserMediaItemDataHandler[*types.Playlist] {
	return handlers.NewUserMediaItemDataHandler[*types.Playlist](coreHandler, service)
}

func ProvideUserCollectionDataHandler(
	coreHandler handlers.CoreUserMediaItemDataHandler[*types.Collection],
	service services.UserMediaItemDataService[*types.Collection],
) handlers.UserMediaItemDataHandler[*types.Collection] {
	return handlers.NewUserMediaItemDataHandler[*types.Collection](coreHandler, service)
}

// --- Handler Provider Sets ---

// CoreHandlerSet is a provider set for core handlers
var CoreHandlerSet = wire.NewSet(
	ProvideCoreMovieHandler,
	ProvideCoreSeriesHandler,
	ProvideCoreEpisodeHandler,
	ProvideCoreSeasonHandler,
	ProvideCoreTrackHandler,
	ProvideCoreAlbumHandler,
	ProvideCoreArtistHandler,
	ProvideCorePlaylistHandler,
	ProvideCoreCollectionHandler,
)

// CoreListHandlerSet is a provider set for core list handlers
var CoreListHandlerSet = wire.NewSet(
	ProvideCorePlaylistListHandler,
	ProvideCoreCollectionListHandler,
)

// UserHandlerSet is a provider set for user handlers
var UserHandlerSet = wire.NewSet(
	ProvideUserMovieHandler,
	ProvideUserSeriesHandler,
	ProvideUserEpisodeHandler,
	ProvideUserSeasonHandler,
	ProvideUserTrackHandler,
	ProvideUserAlbumHandler,
	ProvideUserArtistHandler,
	ProvideUserPlaylistHandler,
	ProvideUserCollectionHandler,
)

// UserListHandlerSet is a provider set for user list handlers
var UserListHandlerSet = wire.NewSet(
	ProvideUserPlaylistListHandler,
	ProvideUserCollectionListHandler,
)

// CoreDataHandlerSet is a provider set for core data handlers
var CoreDataHandlerSet = wire.NewSet(
	ProvideCoreMovieDataHandler,
	ProvideCoreSeriesDataHandler,
	ProvideCoreSeasonDataHandler,
	ProvideCoreEpisodeDataHandler,
	ProvideCoreTrackDataHandler,
	ProvideCoreAlbumDataHandler,
	ProvideCoreArtistDataHandler,
	ProvideCorePlaylistDataHandler,
	ProvideCoreCollectionDataHandler,
)

// UserDataHandlerSet is a provider set for user data handlers
var UserDataHandlerSet = wire.NewSet(
	ProvideUserMovieDataHandler,
	ProvideUserSeriesDataHandler,
	ProvideUserEpisodeDataHandler,
	ProvideUserSeasonDataHandler,
	ProvideUserTrackDataHandler,
	ProvideUserAlbumDataHandler,
	ProvideUserArtistDataHandler,
	ProvideUserPlaylistDataHandler,
	ProvideUserCollectionDataHandler,
)

// HandlerSet combines all handler provider sets
var HandlerSet = wire.NewSet(
	CoreHandlerSet,
	CoreListHandlerSet,
	UserHandlerSet,
	UserListHandlerSet,
	CoreDataHandlerSet,
	UserDataHandlerSet,
)

