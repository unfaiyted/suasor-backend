package automation

import (
	"context"
	"fmt"

	p "suasor/client/automation/providers"
	client "suasor/client/types"
)

type ClientKey struct {
	Type client.AutomationClientType
	ID   uint64
}

// Provider factory type definition
type ClientFactory func(ctx context.Context, clientID uint64, config client.AutomationClientConfig) (AutomationClient, error)

// Registry to store provider factories

var clientFactories = make(map[ClientKey]ClientFactory)

// RegisterAutomationProvider adds a new provider factory to the registry
func RegisterAutomationClient(clientType client.AutomationClientType, clientID uint64, factory ClientFactory) {
	key := ClientKey{Type: clientType, ID: clientID}
	clientFactories[key] = factory

}

// NewAutomationProvider creates providers using the registry
func NewAutomationClient(ctx context.Context, clientID uint64, clientType client.AutomationClientType, config client.AutomationClientConfig) (AutomationClient, error) {
	key := ClientKey{Type: clientType, ID: clientID}
	factory, exists := clientFactories[key]
	if !exists {
		return nil, fmt.Errorf("unsupported automation tool type: %s", clientType)
	}
	return factory(ctx, clientID, config)
}

func AsProvider(client AutomationClient) (p.AutomationProvider, bool) {
	provider, ok := client.(p.AutomationProvider)
	return provider, ok
}

// func to get a registered client
func (ClientFactory) GetAutomationClient(ctx context.Context, clientID uint64, config client.AutomationClientConfig) (AutomationClient, error) {
	clientType := config.GetClientType()
	key := ClientKey{Type: clientType, ID: clientID}
	factory, exists := clientFactories[key]
	if !exists {
		return NewAutomationClient(ctx, clientID, clientType, config)
	}
	return factory(ctx, clientID, config)
}
