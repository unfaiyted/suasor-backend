package automation

import (
	"context"
	"errors"
	client "suasor/clients"
	types "suasor/clients/types"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this automation client")

type AutomationClient interface {
	client.Client
	SupportsMovies() bool
	SupportsTVShows() bool
	SupportsMusic() bool
}

type BaseAutomationClient struct {
	client.BaseClient
	ClientType types.AutomationClientType
	Config     types.AutomationClientConfig
}

func NewAutomationClient(ctx context.Context, clientID uint64, clientType types.AutomationClientType, config types.AutomationClientConfig) (AutomationClient, error) {
	return &BaseAutomationClient{
		BaseClient: client.BaseClient{
			ClientID: clientID,
			Category: clientType.AsCategory(),
			Config:   config,
		},
		ClientType: clientType,
	}, nil
}

// Default caity implementations (all false by default)
func (m *BaseAutomationClient) SupportsMovies() bool  { return false }
func (m *BaseAutomationClient) SupportsTVShows() bool { return false }
func (m *BaseAutomationClient) SupportsMusic() bool   { return false }
