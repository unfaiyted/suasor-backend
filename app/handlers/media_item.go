package handlers

import (
	mediatypes "suasor/client/media/types"
	"suasor/handlers"
)

type MediaItemHandlers interface {
	CoreMediaItemHandlers
	UserMediaItemHandlers
	ClientMediaItemHandlers
}

// CoreMediaItemHandlers defines the core handlers for media items
type CoreMediaItemHandlers interface {
	MovieCoreHandler() *handlers.CoreMediaItemHandler[*mediatypes.Movie]
	SeriesCoreHandler() *handlers.CoreMediaItemHandler[*mediatypes.Series]
	EpisodeCoreHandler() *handlers.CoreMediaItemHandler[*mediatypes.Episode]
	TrackCoreHandler() *handlers.CoreMediaItemHandler[*mediatypes.Track]
	AlbumCoreHandler() *handlers.CoreMediaItemHandler[*mediatypes.Album]
	ArtistCoreHandler() *handlers.CoreMediaItemHandler[*mediatypes.Artist]
	CollectionCoreHandler() *handlers.CoreMediaItemHandler[*mediatypes.Collection]
	PlaylistCoreHandler() *handlers.CoreMediaItemHandler[*mediatypes.Playlist]
}

// User-layer handlers (extend core)
type UserMediaItemHandlers interface {
	MovieUserHandler() *handlers.UserMediaItemHandler[*mediatypes.Movie]
	SeriesUserHandler() *handlers.UserMediaItemHandler[*mediatypes.Series]
	EpisodeUserHandler() *handlers.UserMediaItemHandler[*mediatypes.Episode]
	TrackUserHandler() *handlers.UserMediaItemHandler[*mediatypes.Track]
	AlbumUserHandler() *handlers.UserMediaItemHandler[*mediatypes.Album]
	ArtistUserHandler() *handlers.UserMediaItemHandler[*mediatypes.Artist]
	CollectionUserHandler() *handlers.UserMediaItemHandler[*mediatypes.Collection]
	PlaylistUserHandler() *handlers.UserMediaItemHandler[*mediatypes.Playlist]
}

// Client-layer handlers (extend user)
type ClientMediaItemHandlers interface {
	MovieClientHandler() *handlers.ClientMediaItemHandler[*mediatypes.Movie]
	SeriesClientHandler() *handlers.ClientMediaItemHandler[*mediatypes.Series]
	EpisodeClientHandler() *handlers.ClientMediaItemHandler[*mediatypes.Episode]
	TrackClientHandler() *handlers.ClientMediaItemHandler[*mediatypes.Track]
	AlbumClientHandler() *handlers.ClientMediaItemHandler[*mediatypes.Album]
	ArtistClientHandler() *handlers.ClientMediaItemHandler[*mediatypes.Artist]
	CollectionClientHandler() *handlers.ClientMediaItemHandler[*mediatypes.Collection]
	PlaylistClientHandler() *handlers.ClientMediaItemHandler[*mediatypes.Playlist]
}

type CoreMediaTypeHandlers interface {
	MusicCoreHandler() *handlers.CoreMusicHandler
	MovieCoreHandler() *handlers.CoreMovieHandler
	SeriesCoreHandler() *handlers.CoreSeriesHandler
}

type UserMediaTypeHandlers interface {
	MusicUserHandler() *handlers.UserMusicHandler
	MovieUserHandler() *handlers.UserMovieHandler
	SeriesUserHandler() *handlers.UserSeriesHandler
}
}
