package handlers

import (
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/handlers"
)

type MediaItemHandlers interface {
	CoreMediaItemHandlers
	UserMediaItemHandlers
	// ClientMediaItemHandlers[T clienttypes.ClientMediaConfig]
}

// CoreMediaItemHandlers defines the core handlers for media items
type CoreMediaItemHandlers interface {
	MovieCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Movie]
	SeriesCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Series]
	EpisodeCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Episode]
	TrackCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Track]
	AlbumCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Album]
	ArtistCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Artist]
	CollectionCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Collection]
	PlaylistCoreHandler() handlers.CoreMediaItemHandler[*mediatypes.Playlist]
}

// User-layer handlers (extend core)
type UserMediaItemHandlers interface {
	MovieUserHandler() handlers.UserMediaItemHandler[*mediatypes.Movie]
	SeriesUserHandler() handlers.UserMediaItemHandler[*mediatypes.Series]
	EpisodeUserHandler() handlers.UserMediaItemHandler[*mediatypes.Episode]
	TrackUserHandler() handlers.UserMediaItemHandler[*mediatypes.Track]
	AlbumUserHandler() handlers.UserMediaItemHandler[*mediatypes.Album]
	ArtistUserHandler() handlers.UserMediaItemHandler[*mediatypes.Artist]
	CollectionUserHandler() handlers.UserMediaItemHandler[*mediatypes.Collection]
	PlaylistUserHandler() handlers.UserMediaItemHandler[*mediatypes.Playlist]
}

// Client-layer handlers (extend user)
type ClientMediaItemHandlers[T clienttypes.ClientMediaConfig] interface {
	MovieClientHandler() handlers.ClientMediaItemHandler[T, *mediatypes.Movie]
	SeriesClientHandler() handlers.ClientMediaItemHandler[T, *mediatypes.Series]
	EpisodeClientHandler() handlers.ClientMediaItemHandler[T, *mediatypes.Episode]
	TrackClientHandler() handlers.ClientMediaItemHandler[T, *mediatypes.Track]
	AlbumClientHandler() handlers.ClientMediaItemHandler[T, *mediatypes.Album]
	ArtistClientHandler() handlers.ClientMediaItemHandler[T, *mediatypes.Artist]
	CollectionClientHandler() handlers.ClientMediaItemHandler[T, *mediatypes.Collection]
	PlaylistClientHandler() handlers.ClientMediaItemHandler[T, *mediatypes.Playlist]
}
