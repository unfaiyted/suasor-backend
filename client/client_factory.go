package client

import (
	"context"
	"fmt"
	"suasor/client/types"
	"suasor/utils"
	"sync"
)

type ClientKey struct {
	Type types.ClientType
	ID   uint64
}

// Factory function type
type ClientFactory func(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error)

// ClientFactoryService provides client creation functionality as a singleton
type ClientFactoryService struct {
	factories map[types.ClientType]ClientFactory
	instances map[ClientKey]Client
	mu        sync.RWMutex
}

// Singleton instance with thread-safe initialization
var (
	instance *ClientFactoryService
	once     sync.Once
)

// GetClientFactoryService returns the singleton instance
func GetClientFactoryService() *ClientFactoryService {
	once.Do(func() {
		instance = &ClientFactoryService{
			factories: make(map[types.ClientType]ClientFactory),
			instances: make(map[ClientKey]Client),
		}
	})
	return instance
}

// RegisterClientFactory registers a factory function for a specific client type
func (s *ClientFactoryService) RegisterClientFactory(clientType types.ClientType, factory ClientFactory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.factories[clientType] = factory
}

// UnregisterClient unregisters a client
func (s *ClientFactoryService) UnregisterClient(ctx context.Context, clientID uint64, config types.ClientConfig) {
	log := utils.LoggerFromContext(ctx)
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
func (s *ClientFactoryService) GetClient(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error) {
	log := utils.LoggerFromContext(ctx)
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

// Convenience package-level functions for working with the singleton

// RegisterClientFactory registers a factory at the package level
func RegisterClientFactory(clientType types.ClientType, factory ClientFactory) {
	GetClientFactoryService().RegisterClientFactory(clientType, factory)
}

// GetClient gets or creates a client at the package level
func GetClient(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error) {
	return GetClientFactoryService().GetClient(ctx, clientID, config)
}

// UnregisterClient unregisters a client at the package level
func UnregisterClient(ctx context.Context, clientID uint64, config types.ClientConfig) {
	GetClientFactoryService().UnregisterClient(ctx, clientID, config)
}

// GetClientFromModel creates a client from a client model
func (s *ClientFactoryService) GetClientFromModel(ctx context.Context, model interface{}) (Client, error) {
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
