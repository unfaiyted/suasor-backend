package media

import (
	"context"
	"fmt"

	p "suasor/client/media/providers"
	client "suasor/client/types"
)

type ClientKey struct {
	Type client.MediaClientType
	ID   uint64
}

// Provider factory type definition
type ClientFactory func(ctx context.Context, clientID uint64, config client.MediaClientConfig) (MediaClient, error)

// Registry to store provider factories
var clientFactories = make(map[ClientKey]ClientFactory)

// RegisterProvider adds a new provider factory to the registry
func RegisterClient(clientType client.MediaClientType, clientID uint64, factory ClientFactory) {
	key := ClientKey{Type: clientType, ID: clientID}
	clientFactories[key] = factory
}

// NewMediaClient creates providers using the registry
func NewMediaClient(ctx context.Context, clientID uint64, clientType client.MediaClientType, config client.MediaClientConfig) (MediaClient, error) {
	key := ClientKey{Type: clientType, ID: clientID}
	factory, exists := clientFactories[key]
	if !exists {
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}
	return factory(ctx, clientID, config)
}

func (ClientFactory) GetMediaClient(ctx context.Context, clientID uint64, config client.MediaClientConfig) (MediaClient, error) {
	key := ClientKey{Type: config.GetClientType(), ID: clientID}
	factory, exists := clientFactories[key]
	if !exists {
		return NewMediaClient(ctx, clientID, config.GetClientType(), config)
	}
	return factory(ctx, clientID, config)
}

// Helper functions to safely cast providers
func AsMovieProvider(client MediaClient) (p.MovieProvider, bool) {
	provider, ok := client.(p.MovieProvider)
	return provider, ok && provider.SupportsMovies()
}

func AsSeriesProvider(client MediaClient) (p.SeriesProvider, bool) {
	provider, ok := client.(p.SeriesProvider)
	return provider, ok && provider.SupportsSeries()
}

func AsMusicProvider(client MediaClient) (p.MusicProvider, bool) {
	provider, ok := client.(p.MusicProvider)
	return provider, ok && provider.SupportsMusic()
}

func AsPlaylistProvider(client MediaClient) (p.PlaylistProvider, bool) {
	provider, ok := client.(p.PlaylistProvider)
	return provider, ok && provider.SupportsPlaylists()
}

func AsCollectionProvider(client MediaClient) (p.CollectionProvider, bool) {
	provider, ok := client.(p.CollectionProvider)
	return provider, ok && provider.SupportsCollections()
}

func AsHistoryProvider(client MediaClient) (p.HistoryProvider, bool) {
	provider, ok := client.(p.HistoryProvider)
	return provider, ok && provider.SupportsHistory()
}
