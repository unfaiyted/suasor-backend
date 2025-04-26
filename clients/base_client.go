// interfaces/base_types.go
package clients

import (
	"context"
	"errors"
	"suasor/clients/types"
)

// Common error definitions
var (
	ErrNotImplemented = errors.New("method not implemented")
)

type Client interface {
	GetClientID() uint64
	SetClientID(clientID uint64)
	GetCategory() types.ClientCategory
	GetClientType() types.ClientType
	GetConfig() types.ClientConfig
	TestConnection(ctx context.Context) (bool, error)
}

// BaseClient provides common behavior for all media clients
type client struct {
	ClientID uint64
	Category types.ClientCategory
	Type     types.ClientType
	Config   types.ClientConfig
}

func NewClient(clientID uint64, category types.ClientCategory, config types.ClientConfig) Client {
	return &client{
		ClientID: clientID,
		Category: category,
		Type:     config.GetType(),
		Config:   config,
	}
}

// Get client information
func (b *client) GetClientID() uint64 {
	return b.ClientID
}
func (b *client) SetClientID(clientID uint64) {
	b.ClientID = clientID
}
func (b *client) GetCategory() types.ClientCategory {
	return b.Category
}
func (b *client) GetClientType() types.ClientType {
	return b.Type
}
func (b *client) GetConfig() types.ClientConfig {
	return b.Config
}

func (b *client) TestConnection(ctx context.Context) (bool, error) {
	return false, nil
}

// SimpleClientFactory is a simpler factory type for direct client creation
type SimpleClientFactory func(config types.ClientConfig) (Client, error)

// RegisterClientType registers a client type with the existing factory system
func RegisterClientType(clientType types.ClientType, simpleFactory SimpleClientFactory) {
	// Adapt the simple factory to the expected factory signature
	factory := func(ctx context.Context, clientID uint64, config types.ClientConfig) (Client, error) {
		client, err := simpleFactory(config)
		if err != nil {
			return nil, err
		}

		// Set client ID if possible
		if baseClient, ok := client.(Client); ok {
			baseClient.SetClientID(clientID)
		}

		return client, nil
	}

	RegisterClientProviderFactory(clientType, factory)
}
