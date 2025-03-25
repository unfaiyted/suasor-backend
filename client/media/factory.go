package media

import (
	"context"
	"fmt"

	p "suasor/client/media/providers"
	t "suasor/client/media/types"
)

// Provider factory type definition
type ProviderFactory func(ctx context.Context, clientID uint64, config t.ClientConfig) (MediaClient, error)

// Registry to store provider factories
var providerFactories = make(map[t.MediaClientType]ProviderFactory)

// RegisterProvider adds a new provider factory to the registry
func RegisterProvider(clientType t.MediaClientType, factory ProviderFactory) {
	providerFactories[clientType] = factory
}

// NewMediaClient creates providers using the registry
func NewMediaClient(ctx context.Context, clientID uint64, clientType t.MediaClientType, config t.ClientConfig) (MediaClient, error) {
	factory, exists := providerFactories[clientType]
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

func AsWatchHistoryProvider(client MediaClient) (p.WatchHistoryProvider, bool) {
	provider, ok := client.(p.WatchHistoryProvider)
	return provider, ok && provider.SupportsWatchHistory()
}
