// app/repository/repository/media_data.go
package repository

import (
	mediatypes "suasor/client/media/types"
	"suasor/repository"
)

// CoreMediaItemRepositories defines the core repositories for media items
type CoreMediaItemRepositories interface {
	MovieRepo() repository.MediaItemRepository[*mediatypes.Movie]
	SeriesRepo() repository.MediaItemRepository[*mediatypes.Series]
	EpisodeRepo() repository.MediaItemRepository[*mediatypes.Episode]
	TrackRepo() repository.MediaItemRepository[*mediatypes.Track]
	AlbumRepo() repository.MediaItemRepository[*mediatypes.Album]
	ArtistRepo() repository.MediaItemRepository[*mediatypes.Artist]
	CollectionRepo() repository.MediaItemRepository[*mediatypes.Collection]
	PlaylistRepo() repository.MediaItemRepository[*mediatypes.Playlist]
}

// UserMediaItemRepositories defines the user-specific repository
type UserMediaItemRepositories interface {
	MovieUserRepo() repository.UserMediaItemRepository[*mediatypes.Movie]
	SeriesUserRepo() repository.UserMediaItemRepository[*mediatypes.Series]
	EpisodeUserRepo() repository.UserMediaItemRepository[*mediatypes.Episode]
	TrackUserRepo() repository.UserMediaItemRepository[*mediatypes.Track]
	AlbumUserRepo() repository.UserMediaItemRepository[*mediatypes.Album]
	ArtistUserRepo() repository.UserMediaItemRepository[*mediatypes.Artist]
	CollectionUserRepo() repository.UserMediaItemRepository[*mediatypes.Collection]
	PlaylistUserRepo() repository.UserMediaItemRepository[*mediatypes.Playlist]
}

// Client MediaItems Repositories defines the client-specific repository
type ClientMediaItemRepositories interface {
	MovieClientRepo() repository.ClientMediaItemRepository[*mediatypes.Movie]
	SeriesClientRepo() repository.ClientMediaItemRepository[*mediatypes.Series]
	EpisodeClientRepo() repository.ClientMediaItemRepository[*mediatypes.Episode]
	TrackClientRepo() repository.ClientMediaItemRepository[*mediatypes.Track]
	AlbumClientRepo() repository.ClientMediaItemRepository[*mediatypes.Album]
	ArtistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Artist]
	CollectionClientRepo() repository.ClientMediaItemRepository[*mediatypes.Collection]
	PlaylistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Playlist]
}
