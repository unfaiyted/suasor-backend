// interfaces/base_client.go
package client

import (
	"context"
	client "suasor/client/types"
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
