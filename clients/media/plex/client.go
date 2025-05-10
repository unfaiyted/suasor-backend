// plex/client.go
package plex

import (
	"context"
	"fmt"
	"net/http"
	"suasor/clients/media"
	clienttypes "suasor/clients/types"
	"suasor/utils/logger"

	"github.com/unfaiyted/plexgo"
)

// PlexClient implements MediaContentProvider for Plex
type PlexClient struct {
	media.ClientMedia
	httpClient *http.Client
	plexAPI    *plexgo.PlexAPI
	config     *clienttypes.PlexConfig
}

// NewPlexClient creates a new Plex client
func NewPlexClient(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config *clienttypes.PlexConfig) (media.ClientMedia, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(clienttypes.ClientMediaTypePlex)).
		Str("baseURL", config.GetBaseURL()).
		Msg("Initializing Plex API client")

	// Initialize the Plex API client
	plexAPI := plexgo.New(
		plexgo.WithSecurity(config.GetToken()),
		plexgo.WithServerURL(config.GetBaseURL()),
	)

	clientMedia, err := media.NewClientMedia(ctx, clientID, clienttypes.ClientMediaTypePlex, registry, config)
	if err != nil {
		return nil, err
	}

	pClient := &PlexClient{
		ClientMedia: clientMedia,
		plexAPI:     plexAPI,
		config:      config,
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(clienttypes.ClientMediaTypePlex)).
		Str("baseUrl", config.GetBaseURL()).
		Msg("Successfully created Plex client")

	return pClient, nil
}

// Capability methods
func (c *PlexClient) SupportsMovies() bool  { return true }
func (c *PlexClient) SupportsSeries() bool  { return true }
func (c *PlexClient) SupportsMusic() bool   { return true }
func (c *PlexClient) SupportsHistory() bool { return true }

func (c *PlexClient) plexConfig() *clienttypes.PlexConfig {
	// First check if c.config is already set
	if c.config != nil {
		return c.config
	}

	// If not, try to get from the client interface
	cfg, ok := c.GetConfig().(*clienttypes.PlexConfig)
	if !ok {
		return nil
	}
	return cfg
}

func (c *PlexClient) TestConnection(ctx context.Context) (bool, error) {
	sysInfo, err := c.plexAPI.Server.GetServerCapabilities(ctx)
	if err != nil {
		return false, err
	}
	if sysInfo.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to retrieve Plex server version, Response code: %i", sysInfo.StatusCode)
	}
	return true, nil
}
