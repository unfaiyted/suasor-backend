// app/services/service_interfaces.go
package services

import (
	mediatypes "suasor/client/media/types"
	"suasor/services"
)

// CoreMediaItemServices defines the core services for media items
type CoreMediaItemServices interface {
	MovieCoreService() services.CoreMediaItemService[*mediatypes.Movie]
	SeriesCoreService() services.CoreMediaItemService[*mediatypes.Series]
	EpisodeCoreService() services.CoreMediaItemService[*mediatypes.Episode]
	TrackCoreService() services.CoreMediaItemService[*mediatypes.Track]
	AlbumCoreService() services.CoreMediaItemService[*mediatypes.Album]
	ArtistCoreService() services.CoreMediaItemService[*mediatypes.Artist]
	CollectionCoreService() services.CoreMediaItemService[*mediatypes.Collection]
	PlaylistCoreService() services.CoreMediaItemService[*mediatypes.Playlist]
}

// UserMediaItemServices defines the user-specific services for media items
type UserMediaItemServices interface {
	MovieUserService() services.UserMediaItemService[*mediatypes.Movie]
	SeriesUserService() services.UserMediaItemService[*mediatypes.Series]
	EpisodeUserService() services.UserMediaItemService[*mediatypes.Episode]
	TrackUserService() services.UserMediaItemService[*mediatypes.Track]
	AlbumUserService() services.UserMediaItemService[*mediatypes.Album]
	ArtistUserService() services.UserMediaItemService[*mediatypes.Artist]
	CollectionUserService() services.UserMediaItemService[*mediatypes.Collection]
	PlaylistUserService() services.UserMediaItemService[*mediatypes.Playlist]
}

// ClientMediaItemServices defines the client-specific services for media items
type ClientMediaItemServices interface {
	MovieClientService() services.ClientMediaItemService[*mediatypes.Movie]
	SeriesClientService() services.ClientMediaItemService[*mediatypes.Series]
	EpisodeClientService() services.ClientMediaItemService[*mediatypes.Episode]
	TrackClientService() services.ClientMediaItemService[*mediatypes.Track]
	AlbumClientService() services.ClientMediaItemService[*mediatypes.Album]
	ArtistClientService() services.ClientMediaItemService[*mediatypes.Artist]
	CollectionClientService() services.ClientMediaItemService[*mediatypes.Collection]
	PlaylistClientService() services.ClientMediaItemService[*mediatypes.Playlist]
}

type CoreListServices interface {
	CoreCollectionService() services.CoreCollectionService
	CorePlaylistService() services.CorePlaylistService
}

type UserListServices interface {
	UserCollectionService() services.UserCollectionService
	UserPlaylistService() services.UserPlaylistService
}

type ClientListServices interface {
	ClientCollectionService() services.ClientMediaCollectionService
	ClientPlaylistService() services.ClientPlaylistService
}

// MediaServices interface for media-related services
type MediaServices interface {
	PersonService() *services.PersonService
	CreditService() *services.CreditService
}
