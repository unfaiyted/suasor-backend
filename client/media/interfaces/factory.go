package interfaces

import (
	"context"
	"fmt"
)

// Provider factory type definition
type ProviderFactory func(ctx context.Context, clientID uint64, config any) (MediaClient, error)

// Registry to store provider factories
var providerFactories = make(map[MediaClientType]ProviderFactory)

// RegisterProvider adds a new provider factory to the registry
func RegisterProvider(clientType MediaClientType, factory ProviderFactory) {
	providerFactories[clientType] = factory
}

// NewMediaClient creates providers using the registry
func NewMediaClient(ctx context.Context, clientID uint64, clientType MediaClientType, config any) (MediaClient, error) {
	factory, exists := providerFactories[clientType]
	if !exists {
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}
	return factory(ctx, clientID, config)
}

// Helper functions to safely cast providers
func AsMovieProvider(client MediaClient) (MovieProvider, bool) {
	provider, ok := client.(MovieProvider)
	return provider, ok && provider.SupportsMovies()
}

func AsTVShowProvider(client MediaClient) (TVShowProvider, bool) {
	provider, ok := client.(TVShowProvider)
	return provider, ok && provider.SupportsTVShows()
}

func AsMusicProvider(client MediaClient) (MusicProvider, bool) {
	provider, ok := client.(MusicProvider)
	return provider, ok && provider.SupportsMusic()
}

func AsPlaylistProvider(client MediaClient) (PlaylistProvider, bool) {
	provider, ok := client.(PlaylistProvider)
	return provider, ok && provider.SupportsPlaylists()
}

func AsCollectionProvider(client MediaClient) (CollectionProvider, bool) {
	provider, ok := client.(CollectionProvider)
	return provider, ok && provider.SupportsCollections()
}

func AsWatchHistoryProvider(client MediaClient) (WatchHistoryProvider, bool) {
	provider, ok := client.(WatchHistoryProvider)
	return provider, ok
}
