// app/repository/repository/media_items.go
package repository

import (
	mediatypes "suasor/client/media/types"
	"suasor/repository"
)

// CoreUserMediaItemDataRepositories defines the core data repositories
type CoreUserMediaItemDataRepositories interface {
	MovieCoreService() repository.UserMediaItemDataRepository[*mediatypes.Movie]
	SeriesCoreService() repository.UserMediaItemDataRepository[*mediatypes.Series]
	EpisodeCoreService() repository.UserMediaItemDataRepository[*mediatypes.Episode]
	TrackCoreService() repository.UserMediaItemDataRepository[*mediatypes.Track]
	AlbumCoreService() repository.UserMediaItemDataRepository[*mediatypes.Album]
	ArtistCoreService() repository.UserMediaItemDataRepository[*mediatypes.Artist]
	CollectionCoreService() repository.UserMediaItemDataRepository[*mediatypes.Collection]
	PlaylistCoreService() repository.UserMediaItemDataRepository[*mediatypes.Playlist]
}

// UserDataFactories defines the user data repositories
type UserMediaDataRepositories interface {
	MovieDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Movie]
	SeriesDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Series]
	EpisodeDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Episode]
	TrackDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Track]
	AlbumDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Album]
	ArtistDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Artist]
	CollectionDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Collection]
	PlaylistDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Playlist]
}

// ClientUserDataRepositories defines the client-specific user data repositories
type ClientUserMediaDataRepositories interface {
	MovieDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Movie]
	SeriesDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Series]
	EpisodeDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Episode]
	TrackDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Track]
	AlbumDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Album]
	ArtistDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Artist]
	CollectionDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Collection]
	PlaylistDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Playlist]
}
