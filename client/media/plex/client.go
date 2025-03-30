// plex/client.go
package plex

import (
	"context"
	"fmt"
	"net/http"
	base "suasor/client"
	media "suasor/client/media"
	client "suasor/client/types"
	"suasor/utils"

	"github.com/LukeHagar/plexgo"
	c "suasor/client"
)

// Add this init function to register the plex client factory
func init() {
	c.GetClientFactoryService().RegisterClientFactory(client.ClientTypePlex,
		func(ctx context.Context, clientID uint64, cfg client.ClientConfig) (base.Client, error) {
			// Type assert to plexConfig
			plexConfig, ok := cfg.(*client.PlexConfig)
			if !ok {
				return nil, fmt.Errorf("invalid config type for plex client, expected *EmbyConfig, got %T", cfg)
			}

			// Use your existing constructor
			return NewPlexClient(ctx, clientID, *plexConfig)
		})
}

// PlexClient implements MediaContentProvider for Plex
type PlexClient struct {
	media.BaseMediaClient
	httpClient *http.Client
	plexAPI    *plexgo.PlexAPI
}

// NewPlexClient creates a new Plex client
func NewPlexClient(ctx context.Context, clientID uint64, config client.PlexConfig) (media.MediaClient, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(client.MediaClientTypePlex)).
		Str("baseURL", config.BaseURL).
		Msg("Initializing Plex API client")

	// Initialize the Plex API client
	plexAPI := plexgo.New(
		plexgo.WithSecurity(config.Token),
		plexgo.WithServerURL(config.BaseURL),
	)

	pClient := &PlexClient{
		BaseMediaClient: media.BaseMediaClient{
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: client.MediaClientTypePlex.AsCategory(),
				Config:   &config,
			},
		},
		plexAPI: plexAPI,
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(client.MediaClientTypePlex)).
		Str("baseUrl", config.BaseURL).
		Msg("Successfully created Plex client")

	return pClient, nil
}

// Capability methods
func (c *PlexClient) SupportsMovies() bool      { return true }
func (c *PlexClient) SupportsTVShows() bool     { return true }
func (c *PlexClient) SupportsMusic() bool       { return true }
func (c *PlexClient) SupportsPlaylists() bool   { return true }
func (c *PlexClient) SupportsCollections() bool { return true }

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
