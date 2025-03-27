package media

import (
	"context"
	"fmt"

	p "suasor/client/media/providers"
	client "suasor/client/types"
)

// Provider factory type definition
type ClientFactory func(ctx context.Context, clientID uint64, config client.ClientConfig) (MediaClient, error)

// Registry to store provider factories
var clientFactories = make(map[client.MediaClientType]ClientFactory)

// RegisterProvider adds a new provider factory to the registry
func RegisterClient(clientType client.MediaClientType, factory ClientFactory) {
	clientFactories[clientType] = factory
}

// NewMediaClient creates providers using the registry
func NewMediaClient(ctx context.Context, clientID uint64, clientType client.MediaClientType, config client.ClientConfig) (MediaClient, error) {
	factory, exists := clientFactories[clientType]
	if !exists {
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}
	return factory(ctx, clientID, config)
}

// Helper functions to safely cast providers
func AsMovieProvider(client MediaClient) (p.MovieProvider, bool) {
	provider, ok := client.(p.MovieProvider)
	return provider, ok && provider.SupportsMovies()
}

func AsTVShowProvider(client MediaClient) (p.TVShowProvider, bool) {
	provider, ok := client.(p.TVShowProvider)
	return provider, ok && provider.SupportsTVShows()
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
