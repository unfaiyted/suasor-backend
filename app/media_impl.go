// app/dependencies.go
package app

import (
	mediatypes "suasor/client/media/types"
	"suasor/handlers"
	"suasor/repository"
)

type mediaItemRepositoriesImpl struct {
	movieRepo      repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo     repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo    repository.MediaItemRepository[*mediatypes.Episode]
	trackRepo      repository.MediaItemRepository[*mediatypes.Track]
	albumRepo      repository.MediaItemRepository[*mediatypes.Album]
	artistRepo     repository.MediaItemRepository[*mediatypes.Artist]
	collectionRepo repository.MediaItemRepository[*mediatypes.Collection]
	playlistRepo   repository.MediaItemRepository[*mediatypes.Playlist]
}

func (r *mediaItemRepositoriesImpl) AlbumRepo() repository.MediaItemRepository[*mediatypes.Album] {
	return r.albumRepo
}

func (r *mediaItemRepositoriesImpl) ArtistRepo() repository.MediaItemRepository[*mediatypes.Artist] {
	return r.artistRepo
}

func (r *mediaItemRepositoriesImpl) CollectionRepo() repository.MediaItemRepository[*mediatypes.Collection] {
	return r.collectionRepo
}

func (r *mediaItemRepositoriesImpl) PlaylistRepo() repository.MediaItemRepository[*mediatypes.Playlist] {
	return r.playlistRepo
}
func (r *mediaItemRepositoriesImpl) MovieRepo() repository.MediaItemRepository[*mediatypes.Movie] {
	return r.movieRepo
}

func (r *mediaItemRepositoriesImpl) SeriesRepo() repository.MediaItemRepository[*mediatypes.Series] {
	return r.seriesRepo
}

func (r *mediaItemRepositoriesImpl) EpisodeRepo() repository.MediaItemRepository[*mediatypes.Episode] {
	return r.episodeRepo
}

func (r *mediaItemRepositoriesImpl) TrackRepo() repository.MediaItemRepository[*mediatypes.Track] {
	return r.trackRepo
}

type mediaItemHandlersImpl struct {
	movieHandler      *handlers.MediaItemHandler[*mediatypes.Movie]
	seriesHandler     *handlers.MediaItemHandler[*mediatypes.Series]
	episodeHandler    *handlers.MediaItemHandler[*mediatypes.Episode]
	trackHandler      *handlers.MediaItemHandler[*mediatypes.Track]
	albumHandler      *handlers.MediaItemHandler[*mediatypes.Album]
	artistHandler     *handlers.MediaItemHandler[*mediatypes.Artist]
	collectionHandler *handlers.MediaItemHandler[*mediatypes.Collection]
	playlistHandler   *handlers.MediaItemHandler[*mediatypes.Playlist]
}

func (h *mediaItemHandlersImpl) MovieHandler() *handlers.MediaItemHandler[*mediatypes.Movie] {
	return h.movieHandler
}

func (h *mediaItemHandlersImpl) SeriesHandler() *handlers.MediaItemHandler[*mediatypes.Series] {
	return h.seriesHandler
}

func (h *mediaItemHandlersImpl) EpisodeHandler() *handlers.MediaItemHandler[*mediatypes.Episode] {
	return h.episodeHandler
}

func (h *mediaItemHandlersImpl) TrackHandler() *handlers.MediaItemHandler[*mediatypes.Track] {
	return h.trackHandler
}

func (h *mediaItemHandlersImpl) AlbumHandler() *handlers.MediaItemHandler[*mediatypes.Album] {
	return h.albumHandler
}

func (h *mediaItemHandlersImpl) ArtistHandler() *handlers.MediaItemHandler[*mediatypes.Artist] {
	return h.artistHandler
}

func (h *mediaItemHandlersImpl) CollectionHandler() *handlers.MediaItemHandler[*mediatypes.Collection] {
	return h.collectionHandler
}

func (h *mediaItemHandlersImpl) PlaylistHandler() *handlers.MediaItemHandler[*mediatypes.Playlist] {
	return h.playlistHandler
}
