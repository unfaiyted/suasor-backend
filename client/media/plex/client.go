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
)

// PlexClient implements MediaContentProvider for Plex
type PlexClient struct {
	media.BaseClientMedia
	httpClient *http.Client
	plexAPI    *plexgo.PlexAPI
}

// NewPlexClient creates a new Plex client
func NewPlexClient(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config client.PlexConfig) (media.ClientMedia, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(client.ClientMediaTypePlex)).
		Str("baseURL", config.BaseURL).
		Msg("Initializing Plex API client")

	// Initialize the Plex API client
	plexAPI := plexgo.New(
		plexgo.WithSecurity(config.Token),
		plexgo.WithServerURL(config.BaseURL),
	)

	pClient := &PlexClient{
		BaseClientMedia: media.BaseClientMedia{
			ItemRegistry: registry,
			ClientType:   client.ClientMediaTypePlex,
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: client.ClientMediaTypePlex.AsCategory(),
				Type:     client.ClientTypePlex,
				Config:   &config,
			},
		},
		plexAPI: plexAPI,
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(client.ClientMediaTypePlex)).
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
func (c *PlexClient) SupportsHistory() bool     { return true }

// GetRegistry returns the client's item registry
func (c *PlexClient) GetRegistry() *media.ClientItemRegistry {
	return c.ItemRegistry
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
