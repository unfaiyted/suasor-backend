package client

import (
	"context"
	"fmt"
	"suasor/client/types"
)

type ClientKey struct {
	Type types.ClientType
	ID   uint64
}

// Provider factory type definition
type ClientFactory func(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error)

// Registry to store client factories
var clientFactories = make(map[ClientKey]ClientFactory)

// RegisterProvider adds a new provider factory to the registry
func RegisterClient(clientType types.ClientType, clientID uint64, factory ClientFactory) {
	key := ClientKey{Type: clientType, ID: clientID}
	clientFactories[key] = factory
}

func NewClient(ctx context.Context, clientID uint64, clientType types.ClientType, config types.ClientConfig) (Client, error) {
	key := ClientKey{Type: clientType, ID: clientID}
	factory, exists := clientFactories[key]
	if !exists {
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}
	return factory(ctx, clientID, config)
}

func (ClientFactory) GetClient(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error) {
	key := ClientKey{Type: config.GetType(), ID: clientID}
	factory, exists := clientFactories[key]
	if !exists {
		return NewClient(ctx, clientID, config.GetType(), config)
	}
	return factory(ctx, clientID, config)

}

// // Helper functions to safely cast providers
// func AsMovieProvider(client MediaClient) (p.MovieProvider, bool) {
// 	provider, ok := client.(p.MovieProvider)
// 	return provider, ok && provider.SupportsMovies()
// }
//
// func AsSeriesProvider(client MediaClient) (p.SeriesProvider, bool) {
// 	provider, ok := client.(p.SeriesProvider)
// 	return provider, ok && provider.SupportsSeries()
// }
//
// func AsMusicProvider(client MediaClient) (p.MusicProvider, bool) {
// 	provider, ok := client.(p.MusicProvider)
// 	return provider, ok && provider.SupportsMusic()
// }
//
// func AsPlaylistProvider(client MediaClient) (p.PlaylistProvider, bool) {
// 	provider, ok := client.(p.PlaylistProvider)
// 	return provider, ok && provider.SupportsPlaylists()
// }
//
// func AsCollectionProvider(client MediaClient) (p.CollectionProvider, bool) {
// 	provider, ok := client.(p.CollectionProvider)
// 	return provider, ok && provider.SupportsCollections()
// }
//
// func AsHistoryProvider(client MediaClient) (p.HistoryProvider, bool) {
// 	provider, ok := client.(p.HistoryProvider)
// 	return provider, ok && provider.SupportsHistory()
// }

// ClientFactoryService provides client creation functionality
type ClientFactoryService struct{}

// NewClientFactoryService creates a new factory service
func NewClientFactoryService() *ClientFactoryService {
	return &ClientFactoryService{}
}

// GetClient retrieves or creates a client based on ID and config
func (s *ClientFactoryService) GetClient(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error) {
	key := ClientKey{Type: config.GetType(), ID: clientID}
	factory, exists := clientFactories[key]
	if !exists {
		return NewClient(ctx, clientID, config.GetType(), config)
	}
	return factory(ctx, clientID, config)
}
