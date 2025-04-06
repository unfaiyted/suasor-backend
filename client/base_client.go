// interfaces/base_client.go
package client

import (
	"context"
	"errors"
	client "suasor/client/types"
)

// Common error definitions
var (
	ErrNotImplemented = errors.New("method not implemented")
)

type Client interface {
	GetClientID() uint64
	GetCategory() client.ClientCategory
	GetType() client.ClientType
	GetConfig() client.ClientConfig
	TestConnection(ctx context.Context) (bool, error)
}

// BaseClient provides common behavior for all media clients
type BaseClient struct {
	ClientID uint64
	Category client.ClientCategory
	Type     client.ClientType
	Config   client.ClientConfig
}

// Get client information
func (b *BaseClient) GetClientID() uint64 {
	return b.ClientID
}
func (b *BaseClient) GetCategory() client.ClientCategory {
	return b.Category
}
func (b *BaseClient) GetType() client.ClientType {
	return b.Type
}

func (b *BaseClient) GetConfig() client.ClientConfig {
	return b.Config
}

func (b *BaseClient) TestConnection(ctx context.Context) (bool, error) {
	return false, nil
}

// NewBaseClient creates a new base client
func NewBaseClient() *BaseClient {
	return &BaseClient{}
}

// SimpleClientFactory is a simpler factory type for direct client creation
type SimpleClientFactory func(config client.ClientConfig) (Client, error)

// RegisterClientType registers a client type with the existing factory system
func RegisterClientType(clientType client.ClientType, simpleFactory SimpleClientFactory) {
	// Adapt the simple factory to the expected factory signature
	factory := func(ctx context.Context, clientID uint64, config client.ClientConfig) (Client, error) {
		client, err := simpleFactory(config)
		if err != nil {
			return nil, err
		}
		
		// Set client ID if the client has a BaseClient
		if baseClient, ok := client.(*BaseClient); ok {
			baseClient.ClientID = clientID
		}
		
		return client, nil
	}
	
	RegisterClientFactory(clientType, factory)
}
