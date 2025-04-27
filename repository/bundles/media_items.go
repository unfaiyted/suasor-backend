package bundles

import (
	mediatypes "suasor/clients/media/types"
	"suasor/repository"
)

// CoreMediaItemRepositories provides access to all core media item repositories
type CoreMediaItemRepositories interface {
	MovieRepo() repository.CoreMediaItemRepository[*mediatypes.Movie]
	SeriesRepo() repository.CoreMediaItemRepository[*mediatypes.Series]
	SeasonRepo() repository.CoreMediaItemRepository[*mediatypes.Season]
	EpisodeRepo() repository.CoreMediaItemRepository[*mediatypes.Episode]
	TrackRepo() repository.CoreMediaItemRepository[*mediatypes.Track]
	AlbumRepo() repository.CoreMediaItemRepository[*mediatypes.Album]
	ArtistRepo() repository.CoreMediaItemRepository[*mediatypes.Artist]
	CollectionRepo() repository.CoreMediaItemRepository[*mediatypes.Collection]
	PlaylistRepo() repository.CoreMediaItemRepository[*mediatypes.Playlist]
}

// UserMediaItemRepositories provides access to all user media item repositories
type UserMediaItemRepositories interface {
	MovieUserRepo() repository.UserMediaItemRepository[*mediatypes.Movie]
	SeriesUserRepo() repository.UserMediaItemRepository[*mediatypes.Series]
	SeasonUserRepo() repository.UserMediaItemRepository[*mediatypes.Season]
	EpisodeUserRepo() repository.UserMediaItemRepository[*mediatypes.Episode]
	TrackUserRepo() repository.UserMediaItemRepository[*mediatypes.Track]
	AlbumUserRepo() repository.UserMediaItemRepository[*mediatypes.Album]
	ArtistUserRepo() repository.UserMediaItemRepository[*mediatypes.Artist]
	CollectionUserRepo() repository.UserMediaItemRepository[*mediatypes.Collection]
	PlaylistUserRepo() repository.UserMediaItemRepository[*mediatypes.Playlist]
}

// ClientMediaItemRepositories provides access to all client media item repositories
type ClientMediaItemRepositories interface {
	MovieClientRepo() repository.ClientMediaItemRepository[*mediatypes.Movie]
	SeriesClientRepo() repository.ClientMediaItemRepository[*mediatypes.Series]
	SeasonClientRepo() repository.ClientMediaItemRepository[*mediatypes.Season]
	EpisodeClientRepo() repository.ClientMediaItemRepository[*mediatypes.Episode]
	TrackClientRepo() repository.ClientMediaItemRepository[*mediatypes.Track]
	AlbumClientRepo() repository.ClientMediaItemRepository[*mediatypes.Album]
	ArtistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Artist]
	CollectionClientRepo() repository.ClientMediaItemRepository[*mediatypes.Collection]
	PlaylistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Playlist]
}

// CoreMediaItemRepositories implementation
type coreMediaItemRepositoriesImpl struct {
	movieRepo      repository.CoreMediaItemRepository[*mediatypes.Movie]
	seriesRepo     repository.CoreMediaItemRepository[*mediatypes.Series]
	seasonRepo     repository.CoreMediaItemRepository[*mediatypes.Season]
	episodeRepo    repository.CoreMediaItemRepository[*mediatypes.Episode]
	trackRepo      repository.CoreMediaItemRepository[*mediatypes.Track]
	albumRepo      repository.CoreMediaItemRepository[*mediatypes.Album]
	artistRepo     repository.CoreMediaItemRepository[*mediatypes.Artist]
	collectionRepo repository.CoreMediaItemRepository[*mediatypes.Collection]
	playlistRepo   repository.CoreMediaItemRepository[*mediatypes.Playlist]
}

func NewCoreMediaItemRepositories(
	movieRepo repository.CoreMediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.CoreMediaItemRepository[*mediatypes.Series],
	seasonRepo repository.CoreMediaItemRepository[*mediatypes.Season],
	episodeRepo repository.CoreMediaItemRepository[*mediatypes.Episode],
	trackRepo repository.CoreMediaItemRepository[*mediatypes.Track],
	albumRepo repository.CoreMediaItemRepository[*mediatypes.Album],
	artistRepo repository.CoreMediaItemRepository[*mediatypes.Artist],
	collectionRepo repository.CoreMediaItemRepository[*mediatypes.Collection],
	playlistRepo repository.CoreMediaItemRepository[*mediatypes.Playlist],
) CoreMediaItemRepositories {
	return &coreMediaItemRepositoriesImpl{
		movieRepo:      movieRepo,
		seriesRepo:     seriesRepo,
		seasonRepo:     seasonRepo,
		episodeRepo:    episodeRepo,
		trackRepo:      trackRepo,
		albumRepo:      albumRepo,
		artistRepo:     artistRepo,
		collectionRepo: collectionRepo,
		playlistRepo:   playlistRepo,
	}
}

func (r *coreMediaItemRepositoriesImpl) MovieRepo() repository.CoreMediaItemRepository[*mediatypes.Movie] {
	return r.movieRepo
}

func (r *coreMediaItemRepositoriesImpl) SeriesRepo() repository.CoreMediaItemRepository[*mediatypes.Series] {
	return r.seriesRepo
}

func (r *coreMediaItemRepositoriesImpl) SeasonRepo() repository.CoreMediaItemRepository[*mediatypes.Season] {
	return r.seasonRepo
}

func (r *coreMediaItemRepositoriesImpl) EpisodeRepo() repository.CoreMediaItemRepository[*mediatypes.Episode] {
	return r.episodeRepo
}

func (r *coreMediaItemRepositoriesImpl) TrackRepo() repository.CoreMediaItemRepository[*mediatypes.Track] {
	return r.trackRepo
}

func (r *coreMediaItemRepositoriesImpl) AlbumRepo() repository.CoreMediaItemRepository[*mediatypes.Album] {
	return r.albumRepo
}

func (r *coreMediaItemRepositoriesImpl) ArtistRepo() repository.CoreMediaItemRepository[*mediatypes.Artist] {
	return r.artistRepo
}

func (r *coreMediaItemRepositoriesImpl) CollectionRepo() repository.CoreMediaItemRepository[*mediatypes.Collection] {
	return r.collectionRepo
}

func (r *coreMediaItemRepositoriesImpl) PlaylistRepo() repository.CoreMediaItemRepository[*mediatypes.Playlist] {
	return r.playlistRepo
}

// UserMediaItemRepositories implementation
type userMediaItemRepositoriesImpl struct {
	movieUserRepo      repository.UserMediaItemRepository[*mediatypes.Movie]
	seriesUserRepo     repository.UserMediaItemRepository[*mediatypes.Series]
	seasonUserRepo     repository.UserMediaItemRepository[*mediatypes.Season]
	episodeUserRepo    repository.UserMediaItemRepository[*mediatypes.Episode]
	trackUserRepo      repository.UserMediaItemRepository[*mediatypes.Track]
	albumUserRepo      repository.UserMediaItemRepository[*mediatypes.Album]
	artistUserRepo     repository.UserMediaItemRepository[*mediatypes.Artist]
	collectionUserRepo repository.UserMediaItemRepository[*mediatypes.Collection]
	playlistUserRepo   repository.UserMediaItemRepository[*mediatypes.Playlist]
}

func (r *userMediaItemRepositoriesImpl) MovieUserRepo() repository.UserMediaItemRepository[*mediatypes.Movie] {
	return r.movieUserRepo
}

func (r *userMediaItemRepositoriesImpl) SeriesUserRepo() repository.UserMediaItemRepository[*mediatypes.Series] {
	return r.seriesUserRepo
}

func (r *userMediaItemRepositoriesImpl) SeasonUserRepo() repository.UserMediaItemRepository[*mediatypes.Season] {
	return r.seasonUserRepo
}

func (r *userMediaItemRepositoriesImpl) EpisodeUserRepo() repository.UserMediaItemRepository[*mediatypes.Episode] {
	return r.episodeUserRepo
}

func (r *userMediaItemRepositoriesImpl) TrackUserRepo() repository.UserMediaItemRepository[*mediatypes.Track] {
	return r.trackUserRepo
}

func (r *userMediaItemRepositoriesImpl) AlbumUserRepo() repository.UserMediaItemRepository[*mediatypes.Album] {
	return r.albumUserRepo
}

func (r *userMediaItemRepositoriesImpl) ArtistUserRepo() repository.UserMediaItemRepository[*mediatypes.Artist] {
	return r.artistUserRepo
}

func (r *userMediaItemRepositoriesImpl) CollectionUserRepo() repository.UserMediaItemRepository[*mediatypes.Collection] {
	return r.collectionUserRepo
}

func (r *userMediaItemRepositoriesImpl) PlaylistUserRepo() repository.UserMediaItemRepository[*mediatypes.Playlist] {
	return r.playlistUserRepo
}

// ClientMediaItemRepositories implementation
type clientMediaItemRepositoriesImpl struct {
	movieClientRepo      repository.ClientMediaItemRepository[*mediatypes.Movie]
	seriesClientRepo     repository.ClientMediaItemRepository[*mediatypes.Series]
	seasonClientRepo     repository.ClientMediaItemRepository[*mediatypes.Season]
	episodeClientRepo    repository.ClientMediaItemRepository[*mediatypes.Episode]
	trackClientRepo      repository.ClientMediaItemRepository[*mediatypes.Track]
	albumClientRepo      repository.ClientMediaItemRepository[*mediatypes.Album]
	artistClientRepo     repository.ClientMediaItemRepository[*mediatypes.Artist]
	collectionClientRepo repository.ClientMediaItemRepository[*mediatypes.Collection]
	playlistClientRepo   repository.ClientMediaItemRepository[*mediatypes.Playlist]
}

func (r *clientMediaItemRepositoriesImpl) MovieClientRepo() repository.ClientMediaItemRepository[*mediatypes.Movie] {
	return r.movieClientRepo
}

func (r *clientMediaItemRepositoriesImpl) SeriesClientRepo() repository.ClientMediaItemRepository[*mediatypes.Series] {
	return r.seriesClientRepo
}

func (r *clientMediaItemRepositoriesImpl) SeasonClientRepo() repository.ClientMediaItemRepository[*mediatypes.Season] {
	return r.seasonClientRepo
}

func (r *clientMediaItemRepositoriesImpl) EpisodeClientRepo() repository.ClientMediaItemRepository[*mediatypes.Episode] {
	return r.episodeClientRepo
}

func (r *clientMediaItemRepositoriesImpl) TrackClientRepo() repository.ClientMediaItemRepository[*mediatypes.Track] {
	return r.trackClientRepo
}

func (r *clientMediaItemRepositoriesImpl) AlbumClientRepo() repository.ClientMediaItemRepository[*mediatypes.Album] {
	return r.albumClientRepo
}

func (r *clientMediaItemRepositoriesImpl) ArtistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Artist] {
	return r.artistClientRepo
}

func (r *clientMediaItemRepositoriesImpl) CollectionClientRepo() repository.ClientMediaItemRepository[*mediatypes.Collection] {
	return r.collectionClientRepo
}

func (r *clientMediaItemRepositoriesImpl) PlaylistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Playlist] {
	return r.playlistClientRepo

}

func NewUserMediaItemRepositories(
	movieRepo repository.UserMediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.UserMediaItemRepository[*mediatypes.Series],
	seasonRepo repository.UserMediaItemRepository[*mediatypes.Season],
	episodeRepo repository.UserMediaItemRepository[*mediatypes.Episode],
	trackRepo repository.UserMediaItemRepository[*mediatypes.Track],
	albumRepo repository.UserMediaItemRepository[*mediatypes.Album],
	artistRepo repository.UserMediaItemRepository[*mediatypes.Artist],
	collectionRepo repository.UserMediaItemRepository[*mediatypes.Collection],
	playlistRepo repository.UserMediaItemRepository[*mediatypes.Playlist],
) UserMediaItemRepositories {
	return &userMediaItemRepositoriesImpl{
		movieUserRepo:      movieRepo,
		seriesUserRepo:     seriesRepo,
		seasonUserRepo:     seasonRepo,
		episodeUserRepo:    episodeRepo,
		trackUserRepo:      trackRepo,
		albumUserRepo:      albumRepo,
		artistUserRepo:     artistRepo,
		collectionUserRepo: collectionRepo,
		playlistUserRepo:   playlistRepo,
	}
}

func NewClientMediaItemRepositories(
	movieRepo repository.ClientMediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.ClientMediaItemRepository[*mediatypes.Series],
	seasonRepo repository.ClientMediaItemRepository[*mediatypes.Season],
	episodeRepo repository.ClientMediaItemRepository[*mediatypes.Episode],
	trackRepo repository.ClientMediaItemRepository[*mediatypes.Track],
	albumRepo repository.ClientMediaItemRepository[*mediatypes.Album],
	artistRepo repository.ClientMediaItemRepository[*mediatypes.Artist],
	collectionRepo repository.ClientMediaItemRepository[*mediatypes.Collection],
	playlistRepo repository.ClientMediaItemRepository[*mediatypes.Playlist],
) ClientMediaItemRepositories {
	return &clientMediaItemRepositoriesImpl{
		movieClientRepo:      movieRepo,
		seriesClientRepo:     seriesRepo,
		seasonClientRepo:     seasonRepo,
		episodeClientRepo:    episodeRepo,
		trackClientRepo:      trackRepo,
		albumClientRepo:      albumRepo,
		artistClientRepo:     artistRepo,
		collectionClientRepo: collectionRepo,
		playlistClientRepo:   playlistRepo,
	}
}

// NewCoreUserMediaItemDataRepositories creates a new instance of CoreUserMediaItemDataRepositories
func NewCoreUserMediaItemDataRepositories(
	movieRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Movie],
	seriesRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Series],
	seasonRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Season],
	episodeRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Episode],
	trackRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Track],
	albumRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Album],
	artistRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Artist],
	collectionRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Collection],
	playlistRepo repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist],
) CoreUserMediaItemDataRepositories {
	return &coreUserMediaItemDataRepositoriesImpl{
		movieCoreService:      movieRepo,
		seriesCoreService:     seriesRepo,
		seasonCoreService:     seasonRepo,
		episodeCoreService:    episodeRepo,
		trackCoreService:      trackRepo,
		albumCoreService:      albumRepo,
		artistCoreService:     artistRepo,
		collectionCoreService: collectionRepo,
		playlistCoreService:   playlistRepo,
	}
}

// NewUserMediaDataRepositories creates a new instance of UserMediaDataRepositories
func NewUserMediaDataRepositories(
	movieRepo repository.UserMediaItemDataRepository[*mediatypes.Movie],
	seriesRepo repository.UserMediaItemDataRepository[*mediatypes.Series],
	seasonRepo repository.UserMediaItemDataRepository[*mediatypes.Season],
	episodeRepo repository.UserMediaItemDataRepository[*mediatypes.Episode],
	trackRepo repository.UserMediaItemDataRepository[*mediatypes.Track],
	albumRepo repository.UserMediaItemDataRepository[*mediatypes.Album],
	artistRepo repository.UserMediaItemDataRepository[*mediatypes.Artist],
	collectionRepo repository.UserMediaItemDataRepository[*mediatypes.Collection],
	playlistRepo repository.UserMediaItemDataRepository[*mediatypes.Playlist],
) UserMediaDataRepositories {
	return &userDataRepositoriesImpl{
		movieDataRepo:      movieRepo,
		seriesDataRepo:     seriesRepo,
		seasonDataRepo:     seasonRepo,
		episodeDataRepo:    episodeRepo,
		trackDataRepo:      trackRepo,
		albumDataRepo:      albumRepo,
		artistDataRepo:     artistRepo,
		collectionDataRepo: collectionRepo,
		playlistDataRepo:   playlistRepo,
	}
}
