package lidarr

import (
	"context"
	"fmt"

	lidarr "github.com/devopsarr/lidarr-go/lidarr"

	base "suasor/client"
	c "suasor/client"
	auto "suasor/client/automation"
	config "suasor/client/types"
)

// Add this init function to register the liadarr client factory
func init() {
	c.GetClientFactoryService().RegisterClientFactory(config.ClientTypeLidarr,
		func(ctx context.Context, clientID uint64, cfg config.ClientConfig) (base.Client, error) {
			// Type assert to liadarrConfig
			liadarrConfig, ok := cfg.(*config.LidarrConfig)
			if !ok {
				return nil, fmt.Errorf("invalid config type for liadarr client, expected *EmbyConfig, got %T", cfg)
			}

			// Use your existing constructor
			return NewLidarrClient(ctx, clientID, *liadarrConfig)
		})
}

// Capability methods
func (l *LidarrClient) SupportsMusic() bool { return true }

// LidarrClient implements the AutomationProvider interface
type LidarrClient struct {
	auto.BaseAutomationClient
	client *lidarr.APIClient
}

// NewLidarrClient creates a new Lidarr client instance
func NewLidarrClient(ctx context.Context, clientID uint64, cfg config.LidarrConfig) (auto.AutomationClient, error) {

	// Create API client configuration
	apiConfig := lidarr.NewConfiguration()
	apiConfig.AddDefaultHeader("X-Api-Key", cfg.APIKey)
	apiConfig.Servers = lidarr.ServerConfigurations{
		{
			URL: cfg.BaseURL,
		},
	}

	client := lidarr.NewAPIClient(apiConfig)

	lidarrClient := &LidarrClient{
		BaseAutomationClient: auto.BaseAutomationClient{
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: config.AutomationClientTypeLidarr.AsCategory(),
				Config:   &cfg,
			},
			ClientType: config.AutomationClientTypeLidarr,
		},
		client: client,
	}

	return lidarrClient, nil
}

// Register the provider factory
// func init() {
// 	auto.RegisterAutomationClient(config.AutomationClientTypeLidarr, NewLidarrClient)
// }

func (l *LidarrClient) TestConnection(ctx context.Context) (bool, error) {
	req := l.client.SystemAPI.GetSystemStatus(ctx)
	_, resp, err := req.Execute()
	if err != nil {
		return false, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("Lidarr returned status code %d", resp.StatusCode)
	}
	return true, nil
}
