package sonarr

import (
	"context"
	"fmt"

	sonarr "github.com/devopsarr/sonarr-go/sonarr"

	base "suasor/clients"
	auto "suasor/clients/automation"
	"suasor/clients/automation/types"
	config "suasor/clients/types"

	c "suasor/clients"
)

// Add this init function to register the plex client factory
func init() {
	c.GetClientFactoryService().RegisterClientFactory(config.ClientTypeSonarr,
		func(ctx context.Context, clientID uint64, cfg config.ClientConfig) (base.Client, error) {
			// Type assert to plexConfig
			plexConfig, ok := cfg.(*config.SonarrConfig)
			if !ok {
				return nil, fmt.Errorf("invalid config type for plex client, expected *EmbyConfig, got %T", cfg)
			}

			// Use your existing constructor
			return NewSonarrClient(ctx, clientID, *plexConfig)
		})
}

type SonarrClient struct {
	auto.BaseAutomationClient
	client *sonarr.APIClient
	config config.SonarrConfig
}

// NewSonarrClient creates a new Sonarr client instance
func NewSonarrClient(ctx context.Context, clientID uint64, c config.SonarrConfig) (auto.AutomationClient, error) {

	// Create API client configuration
	apiConfig := sonarr.NewConfiguration()
	apiConfig.AddDefaultHeader("X-Api-Key", c.APIKey)
	apiConfig.Servers = sonarr.ServerConfigurations{
		{
			URL: c.BaseURL,
		},
	}

	client := sonarr.NewAPIClient(apiConfig)

	sonarrClient := &SonarrClient{
		BaseAutomationClient: auto.BaseAutomationClient{
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: config.AutomationClientTypeSonarr.AsCategory(),
			},
			ClientType: config.AutomationClientTypeSonarr,
		},
		client: client,
		config: c,
	}

	return sonarrClient, nil
}

// Register the provider factory
// func init() {
// 	auto.RegisterAutomationClient(config.AutomationClientTypeSonarr, NewSonarrClient)
// }

// Capability methods
func (s *SonarrClient) SupportsMovies() bool  { return false }
func (s *SonarrClient) SupportsTVShows() bool { return true }
func (s *SonarrClient) SupportsMusic() bool   { return false }

// GetSystemStatus retrieves system information from Sonarr

// GetLibraryItems retrieves all series from Sonarr

// Helper function to convert Sonarr series to generic MediaItem

// GetMediaByID retrieves a specific series by ID

// AddMedia adds a new series to Sonarr

// UpdateMedia updates an existing series in Sonarr

// DeleteMedia removes a series from Sonarr

// SearchMedia searches for series in Sonarr

// GetQualityProfiles retrieves available quality profiles from Sonarr

// GetTags retrieves all tags from Sonarr

// CreateTag creates a new tag in Sonarr

// GetCalendar retrieves upcoming releases from Sonarr

func (r *SonarrClient) GetMetadataProfiles(ctx context.Context) ([]types.MetadataProfile, error) {
	return nil, types.ErrAutomationFeatureNotSupported
}

func (l *SonarrClient) TestConnection(ctx context.Context) (bool, error) {
	req := l.client.SystemAPI.GetSystemStatus(ctx)
	_, resp, err := req.Execute()
	if err != nil {
		return false, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("Sonarr returned status code %d", resp.StatusCode)
	}
	return true, nil
}
