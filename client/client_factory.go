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
type ClientFactory func(ctx context.Context, clientID uint64, clientType types.ClientType) (Client, error)

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

// GetClient returns an existing client or creates a new one
func (s *ClientFactoryService) GetClient(ctx context.Context, clientID uint64, clientType types.ClientType) (Client, error) {
	log := utils.LoggerFromContext(ctx)
	key := ClientKey{Type: clientType, ID: clientID}

	// Try to get existing client first (read lock)
	s.mu.RLock()
	client, exists := s.instances[key]
	s.mu.RUnlock()

	if exists {
		log.Info().
			Str("clientType", clientType.String()).
			Uint64("clientID", clientID).
			Msg("Returning existing client instance")
		return client, nil
	}

	// Need to create a new client (write lock)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring lock
	if client, exists := s.instances[key]; exists {
		return client, nil
	}

	// Get factory for the client type
	factory, exists := s.factories[clientType]
	if !exists {
		return nil, fmt.Errorf("no factory registered for client type: %s", clientType)
	}

	// Create and cache new client
	client, err := factory(ctx, clientID, clientType)
	if err != nil {
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
	return GetClientFactoryService().GetClient(ctx, clientID, config.GetType())
}
