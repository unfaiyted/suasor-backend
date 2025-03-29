package lidarr

import (
	"context"
	"fmt"

	lidarr "github.com/devopsarr/lidarr-go/lidarr"

	base "suasor/client"
	auto "suasor/client/automation"
	config "suasor/client/types"
)

// Capability methods
func (l *LidarrClient) SupportsMusic() bool { return true }

// LidarrClient implements the AutomationProvider interface
type LidarrClient struct {
	auto.BaseAutomationClient
	client *lidarr.APIClient
	config config.LidarrConfig
}

// NewLidarrClient creates a new Lidarr client instance
func NewLidarrClient(ctx context.Context, clientID uint64, cfg config.ClientConfig) (auto.AutomationClient, error) {
	// Extract config
	lidarrCfg, ok := cfg.(config.LidarrConfig)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Lidarr client")
	}

	// Create API client configuration
	apiConfig := lidarr.NewConfiguration()
	apiConfig.AddDefaultHeader("X-Api-Key", lidarrCfg.APIKey)
	apiConfig.Servers = lidarr.ServerConfigurations{
		{
			URL: lidarrCfg.BaseURL,
		},
	}

	client := lidarr.NewAPIClient(apiConfig)

	lidarrClient := &LidarrClient{
		BaseAutomationClient: auto.BaseAutomationClient{
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: config.AutomationClientTypeLidarr.AsCategory(),
			},
			ClientType: config.AutomationClientTypeLidarr,
		},
		client: client,
		config: lidarrCfg,
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
