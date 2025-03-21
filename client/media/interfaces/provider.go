package interfaces

import (
	"context"
	"errors"
	"fmt"
	"suasor/models"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// MediaContentProvider defines a common interface for all media client implementations
type MediaContentProvider interface {
	// Capability methods
	SupportsMovies() bool
	SupportsTVShows() bool
	SupportsMusic() bool
	SupportsPlaylists() bool
	SupportsCollections() bool

	// Media retrieval methods
	GetPlaylists(ctx context.Context, options *QueryOptions) ([]Playlist, error)
	GetCollections(ctx context.Context, options *QueryOptions) ([]Collection, error)
	GetMovies(ctx context.Context, options *QueryOptions) ([]Movie, error)
	GetTVShows(ctx context.Context, options *QueryOptions) ([]TVShow, error)
	GetTVShowSeasons(ctx context.Context, showID string) ([]Season, error)
	GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]Episode, error)
	GetMusic(ctx context.Context, options *QueryOptions) ([]MusicTrack, error)
	GetMusicArtists(ctx context.Context, options *QueryOptions) ([]MusicArtist, error)
	GetMusicAlbums(ctx context.Context, options *QueryOptions) ([]MusicAlbum, error)
	GetWatchHistory(ctx context.Context, options *QueryOptions) ([]WatchHistoryItem, error)

	// Item retrieval methods
	GetMovieByID(ctx context.Context, id string) (Movie, error)
	GetTVShowByID(ctx context.Context, id string) (TVShow, error)
	GetEpisodeByID(ctx context.Context, id string) (Episode, error)
	GetMusicTrackByID(ctx context.Context, id string) (MusicTrack, error)

	// Genre methods
	GetMusicGenres(ctx context.Context) ([]string, error)
	GetMovieGenres(ctx context.Context) ([]string, error)
}

// BaseMediaClient provides common behavior for all media clients
type BaseMediaClient struct {
	ClientID   uint64
	ClientType models.MediaClientType
}

// Default "not supported" implementations
func (b *BaseMediaClient) SupportsMovies() bool      { return false }
func (b *BaseMediaClient) SupportsTVShows() bool     { return false }
func (b *BaseMediaClient) SupportsMusic() bool       { return false }
func (b *BaseMediaClient) SupportsPlaylists() bool   { return false }
func (b *BaseMediaClient) SupportsCollections() bool { return false }

// Default implementation for unsupported features
func (b *BaseMediaClient) GetPlaylists(ctx context.Context, options *QueryOptions) ([]Playlist, error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetMovies(ctx context.Context, options *QueryOptions) ([]Movie, error) {
	return nil, ErrFeatureNotSupported
}

// TODO: Add other method implementations with ErrFeatureNotSupported

// Helper to add client information to items
func (b *BaseMediaClient) AddClientInfo(item *MediaItem) {
	item.ClientID = b.ClientID
	item.ClientType = string(b.ClientType)
}

// Provider factory type definition
type ProviderFactory func(ctx context.Context, clientID uint64, config interface{}) (MediaContentProvider, error)

// Registry to store provider factories
var providerFactories = make(map[models.MediaClientType]ProviderFactory)

// RegisterProvider adds a new provider factory to the registry
func RegisterProvider(clientType models.MediaClientType, factory ProviderFactory) {
	providerFactories[clientType] = factory
}

// NewMediaContentProvider creates providers using the registry
func NewMediaContentProvider(ctx context.Context, clientID uint64, clientType models.MediaClientType, config interface{}) (MediaContentProvider, error) {
	factory, exists := providerFactories[clientType]
	if !exists {
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}
	return factory(ctx, clientID, config)
}
