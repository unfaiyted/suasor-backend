// interfaces/media_client.go
package ai

import (
	"context"
	"errors"
	"suasor/client"
	types "suasor/client/types"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// MediaClient defines basic client information that all providers must implement

type AIClient interface {
	client.Client
}

type BaseAIClient struct {
	client.BaseClient
	ClientType types.AIClientType
	config     *types.AIClientConfig
}

func NewAIClient(ctx context.Context, clientID uint64, clientType types.AIClientType, config types.AIClientConfig) (AIClient, error) {
	return &BaseAIClient{
		BaseClient: client.BaseClient{
			ClientID: clientID,
			Category: clientType.AsCategory(),
		},
		config:     &config,
		ClientType: clientType,
	}, nil
}

func (b *BaseAIClient) TestConnection(ctx context.Context) (bool, error) {
	return false, nil
}
