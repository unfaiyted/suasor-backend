package radarr

import (
	"context"
	"fmt"

	radarr "github.com/devopsarr/radarr-go/radarr"

	base "suasor/client"
	auto "suasor/client/automation"
	"suasor/client/automation/types"
	config "suasor/client/types"
)

// RadarrClient implements the AutomationProvider interface
type RadarrClient struct {
	auto.BaseAutomationClient
	client *radarr.APIClient
	config config.RadarrConfig
}

// NewRadarrClient creates a new Radarr client instance
func NewRadarrClient(ctx context.Context, clientID uint64, c config.ClientConfig) (auto.AutomationClient, error) {
	// Extract config
	cfg, ok := c.(config.RadarrConfig)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Radarr client")
	}

	// Create API client configuration
	apiConfig := radarr.NewConfiguration()
	apiConfig.AddDefaultHeader("X-Api-Key", cfg.APIKey)
	apiConfig.Servers = radarr.ServerConfigurations{
		{
			URL: cfg.BaseURL,
		},
	}

	client := radarr.NewAPIClient(apiConfig)

	radarrClient := &RadarrClient{
		BaseAutomationClient: auto.BaseAutomationClient{
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: config.AutomationClientTypeRadarr.AsCategory(),
			},
			ClientType: config.AutomationClientTypeRadarr,
		},
		client: client,
		config: cfg,
	}

	return radarrClient, nil
}

// Register the provider factory
// func init() {
// 	auto.RegisterAutomationClient(config.AutomationClientTypeRadarr, NewRadarrClient)
// }

// Capability methods
func (r *RadarrClient) SupportsMovies() bool  { return true }
func (r *RadarrClient) SupportsTVShows() bool { return false }
func (r *RadarrClient) SupportsMusic() bool   { return false }

func (r *RadarrClient) GetMetadataProfiles(ctx context.Context) ([]types.MetadataProfile, error) {
	return nil, types.ErrAutomationFeatureNotSupported
}

func (l *RadarrClient) TestConnection(ctx context.Context) (bool, error) {
	req := l.client.SystemAPI.GetSystemStatus(ctx)
	_, resp, err := req.Execute()
	if err != nil {
		return false, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("Radarr returned status code %d", resp.StatusCode)
	}
	return true, nil
}
