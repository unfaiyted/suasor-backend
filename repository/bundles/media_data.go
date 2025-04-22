package bundles

import (
	mediatypes "suasor/client/media/types"
	"suasor/repository"
)

// CoreUserMediaItemDataRepositories defines the core data repositories
type CoreUserMediaItemDataRepositories interface {
	MovieCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Movie]
	SeriesCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Series]
	EpisodeCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Episode]
	TrackCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Track]
	AlbumCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Album]
	ArtistCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Artist]
	CollectionCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Collection]
	PlaylistCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist]
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

// CoreUserMediaItemDataRepositories implementation
type coreUserMediaItemDataRepositoriesImpl struct {
	movieCoreService      repository.CoreUserMediaItemDataRepository[*mediatypes.Movie]
	seriesCoreService     repository.CoreUserMediaItemDataRepository[*mediatypes.Series]
	episodeCoreService    repository.CoreUserMediaItemDataRepository[*mediatypes.Episode]
	trackCoreService      repository.CoreUserMediaItemDataRepository[*mediatypes.Track]
	albumCoreService      repository.CoreUserMediaItemDataRepository[*mediatypes.Album]
	artistCoreService     repository.CoreUserMediaItemDataRepository[*mediatypes.Artist]
	collectionCoreService repository.CoreUserMediaItemDataRepository[*mediatypes.Collection]
	playlistCoreService   repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist]
}

func (r *coreUserMediaItemDataRepositoriesImpl) MovieCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Movie] {
	return r.movieCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) SeriesCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Series] {
	return r.seriesCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) EpisodeCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Episode] {
	return r.episodeCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) TrackCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Track] {
	return r.trackCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) AlbumCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Album] {
	return r.albumCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) ArtistCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Artist] {
	return r.artistCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) CollectionCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Collection] {
	return r.collectionCoreService
}

func (r *coreUserMediaItemDataRepositoriesImpl) PlaylistCoreService() repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist] {
	return r.playlistCoreService
}

// UserMediaDataRepositories implementation
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
