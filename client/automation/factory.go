package automation

import (
	"context"
	"fmt"

	p "suasor/client/automation/providers"
	client "suasor/client/types"
)

// Provider factory type definition
type ClientFactory func(ctx context.Context, clientID uint64, config client.ClientConfig) (AutomationClient, error)

// Registry to store provider factories
var clientFactories = make(map[client.AutomationClientType]ClientFactory)

// RegisterAutomationProvider adds a new provider factory to the registry
func RegisterAutomationClient(clientType client.AutomationClientType, factory ClientFactory) {
	clientFactories[clientType] = factory
}

// NewAutomationProvider creates providers using the registry
func NewAutomationClient(ctx context.Context, clientID uint64, clientType client.AutomationClientType, config client.ClientConfig) (AutomationClient, error) {
	factory, exists := clientFactories[clientType]
	if !exists {
		return nil, fmt.Errorf("unsupported automation tool type: %s", clientType)
	}
	return factory(ctx, clientID, config)
}

func AsProvider(client AutomationClient) (p.AutomationProvider, bool) {
	provider, ok := client.(p.AutomationProvider)
	return provider, ok
}
