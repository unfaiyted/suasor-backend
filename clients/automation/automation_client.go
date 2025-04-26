package automation

import (
	"context"
	"errors"
	client "suasor/clients"
	types "suasor/clients/types"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this automation client")

type ClientAutomation interface {
	client.Client
	SupportsMovies() bool
	SupportsTVShows() bool
	SupportsMusic() bool
}

type BaseAutomationClient struct {
	client.Client
	ClientType types.AutomationClientType
	Config     types.ClientAutomationConfig
}

func NewAutomationClient(ctx context.Context, clientID uint64, clientType types.AutomationClientType, config types.ClientAutomationConfig) (ClientAutomation, error) {
	return &BaseAutomationClient{
		Client:     client.NewClient(clientID, clientType.AsCategory(), config),
		ClientType: clientType,
	}, nil
}

// Default caity implementations (all false by default)
func (m *BaseAutomationClient) SupportsMovies() bool  { return false }
func (m *BaseAutomationClient) SupportsTVShows() bool { return false }
func (m *BaseAutomationClient) SupportsMusic() bool   { return false }
