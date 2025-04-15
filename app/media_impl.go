// app/media_impl.go
package app

import (
	mediatypes "suasor/client/media/types"
	"suasor/repository"
)

// mediaItemRepositoriesImpl implements the legacy MediaItemRepositories interface for backward compatibility
type mediaItemRepositoriesImpl struct {
	movieRepo      repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo     repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo    repository.MediaItemRepository[*mediatypes.Episode]
	trackRepo      repository.MediaItemRepository[*mediatypes.Track]
	albumRepo      repository.MediaItemRepository[*mediatypes.Album]
	artistRepo     repository.MediaItemRepository[*mediatypes.Artist]
	collectionRepo repository.MediaItemRepository[*mediatypes.Collection]
	playlistRepo   repository.MediaItemRepository[*mediatypes.Playlist]
	
	// User-owned media repositories
	userMediaPlaylistRepo repository.UserMediaItemRepository[*mediatypes.Playlist]
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

// UserMediaPlaylistRepo returns the user media playlist repository
func (r *mediaItemRepositoriesImpl) UserMediaPlaylistRepo() repository.UserMediaItemRepository[*mediatypes.Playlist] {
	return r.userMediaPlaylistRepo
}

// Note: Core repositories and service implementations are now defined in media_data_factory.go
// to avoid duplication and ensure consistency across the codebase.

// NOTE: Legacy media item handlers implementation is commented out
// since it's no longer needed with the new three-pronged architecture
// and also has some type compatibility issues.
/*
// Specialized media handlers
type legacyMediaItemHandlersImpl struct {
	// Three-pronged architecture handlers
	coreHandlers   CoreMediaItemHandlers
	userHandlers   UserMediaItemHandlers
	clientHandlers ClientMediaItemHandlers
	
	// Specialized handlers
	musicHandler *handlers.CoreMusicHandler
	seriesHandler *handlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig]
	seasonHandler *handlers.CoreUserMediaItemDataHandler[*mediatypes.Season]
}

// Implementations that simulate the old handler interface by delegating to the new architecture
func (h *legacyMediaItemHandlersImpl) MovieHandler() *handlers.CoreMovieHandler {
	// Properly delegate to the core movie handler using the core service
	return handlers.NewCoreMovieHandler(
		h.coreHandlers.MovieCoreHandler().Service(),
	)
}

func (h *legacyMediaItemHandlersImpl) SeriesHandler() *handlers.CoreMovieHandler {
	// This is a mismatch in the legacy interface - the actual return type should be adjusted
	// For now, we're returning a movie handler to maintain compatibility
	return handlers.NewCoreMovieHandler(
		h.coreHandlers.MovieCoreHandler().Service(),
	)
}

func (h *legacyMediaItemHandlersImpl) EpisodeHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Episode] {
	return h.coreHandlers.EpisodeCoreHandler()
}

func (h *legacyMediaItemHandlersImpl) SeasonHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Season] {
	return h.seasonHandler
}

func (h *legacyMediaItemHandlersImpl) TrackHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Track] {
	return h.coreHandlers.TrackCoreHandler()
}

func (h *legacyMediaItemHandlersImpl) AlbumHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Album] {
	return h.coreHandlers.AlbumCoreHandler()
}

func (h *legacyMediaItemHandlersImpl) ArtistHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Artist] {
	return h.coreHandlers.ArtistCoreHandler()
}

func (h *legacyMediaItemHandlersImpl) CollectionHandler() *handlers.CoreCollectionHandler {
	// Properly initialize the collection handler with proper services
	return handlers.NewCoreCollectionHandler(
		h.coreHandlers.CollectionCoreHandler().Service(),
		services.NewCoreCollectionService(nil), // This is a placeholder, pass nil for simplicity
	)
}

func (h *legacyMediaItemHandlersImpl) PlaylistHandler() *handlers.CorePlaylistHandler {
	// Properly initialize the playlist handler with proper services
	return handlers.NewCorePlaylistHandler(
		h.coreHandlers.PlaylistCoreHandler().Service(),
		services.NewPlaylistService(nil, nil, nil), // This is a placeholder, pass nil for simplicity
	)
}

func (h *legacyMediaItemHandlersImpl) GetCoreHandlers() CoreMediaItemHandlers {
	return h.coreHandlers
}

func (h *legacyMediaItemHandlersImpl) GetUserHandlers() UserMediaItemHandlers {
	return h.userHandlers
}

func (h *legacyMediaItemHandlersImpl) GetClientHandlers() ClientMediaItemHandlers {
	return h.clientHandlers
}

func (h *legacyMediaItemHandlersImpl) MusicHandler() *handlers.CoreMusicHandler {
	return h.musicHandler
}

func (h *legacyMediaItemHandlersImpl) SeriesSpecificHandler() *handlers.ClientMediaSeriesHandler[*clienttypes.JellyfinConfig] {
	return h.seriesHandler
}

// These handlers are specialized versions of the generic handlers
// They provide domain-specific functionality beyond the basic CRUD operations
func (h *legacyMediaItemHandlersImpl) PlaylistSpecificHandler() *handlers.CorePlaylistHandler {
	// Reuse the same implementation as the standard playlist handler
	// but we could extend this with additional specialized functionality if needed
	return handlers.NewCorePlaylistHandler(
		h.coreHandlers.PlaylistCoreHandler().Service(),
		services.NewPlaylistService(nil, nil, nil), // This is a placeholder, pass nil for simplicity
	)
}

func (h *legacyMediaItemHandlersImpl) CollectionSpecificHandler() *handlers.CoreCollectionHandler {
	// Reuse the same implementation as the standard collection handler
	// but we could extend this with additional specialized functionality if needed
	return handlers.NewCoreCollectionHandler(
		h.coreHandlers.CollectionCoreHandler().Service(),
		services.NewCoreCollectionService(nil), // This is a placeholder, pass nil for simplicity
	)
}
*/