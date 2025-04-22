package clients

import (
	"context"
	"fmt"
	"suasor/clients/media/providers"
	"suasor/clients/types"
	"suasor/utils/logger"
	"sync"
)

type ClientKey struct {
	Type types.ClientType
	ID   uint64
}

type ClientProviderFactory func(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error)

// ClientProviderFactoryService provides client  provider creation functionality as a singleton
type ClientProviderFactoryService struct {
	factories map[types.ClientType]ClientProviderFactory
	instances map[ClientKey]Client
	mu        sync.RWMutex
}

// Singleton instance with thread-safe initialization
var (
	instance *ClientProviderFactoryService
	once     sync.Once
)

// GetClientProviderFactoryService returns the singleton instance
func GetClientProviderFactoryService() *ClientProviderFactoryService {
	once.Do(func() {
		instance = &ClientProviderFactoryService{
			factories: make(map[types.ClientType]ClientProviderFactory),
			instances: make(map[ClientKey]Client),
		}
	})
	return instance
}

// RegisterClientProviderFactory registers a factory function for a specific client provider type
func (s *ClientProviderFactoryService) RegisterClientProviderFactory(clientType types.ClientType, factory ClientProviderFactory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.factories[clientType] = factory
}

// UnregisterClient unregisters a client
func (s *ClientProviderFactoryService) UnregisterClient(ctx context.Context, clientID uint64, config types.ClientConfig) {
	log := logger.LoggerFromContext(ctx)
	clientType := config.GetType()
	key := ClientKey{Type: clientType, ID: clientID}
	log.Debug().
		Str("clientType", clientType.String()).
		Uint64("clientID", clientID).
		Msg("Factory service unregistering client")

	// Try to get existing client first (read lock)
	s.mu.RLock()
	_, exists := s.instances[key]
	s.mu.RUnlock()

	if exists {
		log.Info().
			Str("clientType", clientType.String()).
			Uint64("clientID", clientID).
			Msg("Factory unregistering existing client instance")
		delete(s.instances, key)
		return
	}

	// Need to create a new client (write lock)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring lock
	if _, exists := s.instances[key]; exists {
		delete(s.instances, key)
		return
	}

	log.Info().
		Str("clientType", clientType.String()).
		Uint64("clientID", clientID).
		Msg("Factory not found, no client to unregister")
}

// GetClient returns an existing client or creates a new one
func (s *ClientProviderFactoryService) GetClient(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error) {
	log := logger.LoggerFromContext(ctx)
	clientType := config.GetType()
	key := ClientKey{Type: clientType, ID: clientID}
	log.Debug().
		Str("clientType", clientType.String()).
		Uint64("clientID", clientID).
		Msg("Factory service retrieving client")

	// Try to get existing client first (read lock)
	s.mu.RLock()
	client, exists := s.instances[key]
	s.mu.RUnlock()

	if exists && clientID != 0 {
		log.Info().
			Str("clientType", clientType.String()).
			Uint64("clientID", clientID).
			Msg("Factory returning existing client instance")
		return client, nil
	}

	if clientID == 0 {
		// delete the exising client
		delete(s.instances, key)
	}

	// Need to create a new client (write lock)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring lock
	if client, exists := s.instances[key]; exists {
		log.Info().
			Str("clientType", clientType.String()).
			Uint64("clientID", clientID).
			Msg("Factory returning existing client instance")
		return client, nil
	}

	// Get factory for the client type
	factory, exists := s.factories[clientType]
	if !exists {
		log.Error().
			Str("clientType", clientType.String()).
			Uint64("clientID", clientID).
			Msg("No factory registered for client type")
		return nil, fmt.Errorf("no factory registered for client type: %s", clientType)
	}

	// log configuration to creat the new client instances
	log.Info().
		Str("clientType", clientType.String()).
		Uint64("clientID", clientID).
		Str("config.Category", config.GetCategory().String()).
		Str("config.Type", config.GetType().String()).
		Msg("Creating new client")

	// Create and cache new client
	client, err := factory(ctx, clientID, config)
	if err != nil {
		log.Error().
			Str("clientType", clientType.String()).
			Uint64("clientID", clientID).
			Msg("Error creating new client instance")
		return nil, err
	}

	s.instances[key] = client
	log.Info().
		Str("clientType", clientType.String()).
		Uint64("clientID", clientID).
		Msg("Created and cached new client instance")

	return client, nil
}

func (s *ClientProviderFactoryService) GetMovieProvider(ctx context.Context, clientID uint64, config types.ClientConfig) (providers.MovieProvider, error) {
	client, err := s.GetClient(ctx, clientID, config)
	if err != nil {
		return nil, err
	}
	return client.(providers.MovieProvider), nil
}

func (s *ClientProviderFactoryService) GetSeriesProvider(ctx context.Context, clientID uint64, config types.ClientConfig) (providers.SeriesProvider, error) {
	client, err := s.GetClient(ctx, clientID, config)
	if err != nil {
		return nil, err
	}
	return client.(providers.SeriesProvider), nil
}

func (s *ClientProviderFactoryService) GetMusicProvider(ctx context.Context, clientID uint64, config types.ClientConfig) (providers.MusicProvider, error) {
	client, err := s.GetClient(ctx, clientID, config)
	if err != nil {
		return nil, err
	}
	return client.(providers.MusicProvider), nil
}

// Convenience package-level functions for working with the singleton

// RegisterClientProviderFactory registers a factory at the package level
func RegisterClientProviderFactory(clientType types.ClientType, factory ClientProviderFactory) {
	GetClientProviderFactoryService().RegisterClientProviderFactory(clientType, factory)
}

// GetClient gets or creates a client at the package level
func GetClient(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error) {
	return GetClientProviderFactoryService().GetClient(ctx, clientID, config)
}

// UnregisterClient unregisters a client at the package level
func UnregisterClient(ctx context.Context, clientID uint64, config types.ClientConfig) {
	GetClientProviderFactoryService().UnregisterClient(ctx, clientID, config)
}

// GetClientFromModel creates a client from a client model
func (s *ClientProviderFactoryService) GetClientFromModel(ctx context.Context, model interface{}) (Client, error) {
	// Try to extract client ID and config from the model
	clientID, config, err := ExtractClientInfo(model)
	if err != nil {
		return nil, err
	}

	// Use existing method to get/create client
	return s.GetClient(ctx, clientID, config)
}

// ExtractClientInfo extracts client ID and config from a client model
func ExtractClientInfo(model interface{}) (uint64, types.ClientConfig, error) {
	// Generic extractor that tries to work with any client model type
	// This is a simplified implementation

	// Try to access ID and Config fields via reflection or type assertion
	// For now, we'll just return placeholders
	return 0, nil, fmt.Errorf("not implemented: extracting client info from model")
}
