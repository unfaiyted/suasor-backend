package sonarr

import (
	"context"
	"fmt"

	sonarr "github.com/devopsarr/sonarr-go/sonarr"

	base "suasor/client"
	auto "suasor/client/automation"
	"suasor/client/automation/types"
	config "suasor/client/types"
)

type SonarrClient struct {
	auto.BaseAutomationClient
	client *sonarr.APIClient
	config config.SonarrConfig
}

// NewSonarrClient creates a new Sonarr client instance
func NewSonarrClient(ctx context.Context, clientID uint64, c config.ClientConfig) (auto.AutomationClient, error) {
	// Extract config
	cfg, ok := c.(config.SonarrConfig)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Sonarr client")
	}

	// Create API client configuration
	apiConfig := sonarr.NewConfiguration()
	apiConfig.AddDefaultHeader("X-Api-Key", cfg.APIKey)
	apiConfig.Servers = sonarr.ServerConfigurations{
		{
			URL: cfg.BaseURL,
		},
	}

	client := sonarr.NewAPIClient(apiConfig)

	sonarrClient := &SonarrClient{
		BaseAutomationClient: auto.BaseAutomationClient{
			BaseClient: base.BaseClient{
				ClientID:   clientID,
				ClientType: config.AutomationClientTypeSonarr.AsClientType(),
			},
			ClientType: config.AutomationClientTypeSonarr,
		},
		client: client,
		config: cfg,
	}

	return sonarrClient, nil
}

// Register the provider factory
func init() {
	auto.RegisterAutomationClient(config.AutomationClientTypeSonarr, NewSonarrClient)
}

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
