// plex/client.go
package plex

import (
	"context"
	"fmt"
	"net/http"
	base "suasor/client"
	media "suasor/client/media"
	"suasor/client/media/types"
	"suasor/utils"

	"github.com/LukeHagar/plexgo"
)

// init is automatically called when package is imported
func init() {
	media.RegisterProvider(types.MediaClientTypePlex, NewPlexClient)
}

// PlexClient implements MediaContentProvider for Plex
type PlexClient struct {
	base.BaseMediaClient
	config     types.PlexConfig
	httpClient *http.Client
	baseURL    string
	plexAPI    *plexgo.PlexAPI
}

// NewPlexClient creates a new Plex client
func NewPlexClient(ctx context.Context, clientID uint64, config any) (media.MediaClient, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(types.MediaClientTypePlex)).
		Msg("Creating new Plex client")

	plexConfig, ok := config.(types.PlexConfig)
	if !ok {
		log.Error().
			Uint64("clientID", clientID).
			Str("clientType", string(types.MediaClientTypePlex)).
			Msg("Invalid Plex configuration")
		return nil, fmt.Errorf("invalid Plex configuration")
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("host", plexConfig.Host).
		Msg("Initializing Plex API client")

	// Initialize the Plex API client
	plexAPI := plexgo.New(
		plexgo.WithSecurity(plexConfig.Token),
		plexgo.WithServerURL(plexConfig.Host),
	)

	client := &PlexClient{
		BaseMediaClient: base.BaseMediaClient{
			ClientID:   clientID,
			ClientType: types.MediaClientTypePlex,
		},
		config:  plexConfig,
		plexAPI: plexAPI,
		baseURL: plexConfig.Host,
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(types.MediaClientTypePlex)).
		Str("host", plexConfig.Host).
		Msg("Successfully created Plex client")

	return client, nil
}

// Capability methods
func (c *PlexClient) SupportsMovies() bool      { return true }
func (c *PlexClient) SupportsTVShows() bool     { return true }
func (c *PlexClient) SupportsMusic() bool       { return true }
func (c *PlexClient) SupportsPlaylists() bool   { return true }
func (c *PlexClient) SupportsCollections() bool { return true }
