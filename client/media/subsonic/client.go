package subsonic

import (
	"context"
	"fmt"
	"net/http"
	base "suasor/client"
	media "suasor/client/media"
	types "suasor/client/types"
	"suasor/utils"
	"time"

	gosonic "github.com/supersonic-app/go-subsonic/subsonic"

	c "suasor/client"
)

// Add this init function to register the subsonic client factory
func init() {
	c.GetClientFactoryService().RegisterClientFactory(types.ClientTypeSubsonic,
		func(ctx context.Context, clientID uint64, cfg types.ClientConfig) (base.Client, error) {
			// Type assert to subsonicConfig
			subsonicConfig, ok := cfg.(*types.SubsonicConfig)
			if !ok {
				return nil, fmt.Errorf("invalid config type for subsonic client, expected *EmbyConfig, got %T", cfg)
			}

			// Use your existing constructor
			return NewSubsonicClient(ctx, clientID, *subsonicConfig)
		})
}

// SubsonicClient implements MediaContentProvider for Subsonic
type SubsonicClient struct {
	media.BaseMediaClient
	httpClient *http.Client
	client     *gosonic.Client
}

// NewSubsonicClient creates a new Subsonic client
func NewSubsonicClient(ctx context.Context, clientID uint64, config types.SubsonicConfig) (media.MediaClient, error) {

	log := utils.LoggerFromContext(context.Background())

	log.Info().
		Uint64("clientID", clientID).
		Str("baseURL", config.BaseURL).
		Bool("ssl", config.SSL).
		Msg("Creating new Subsonic client")

	// protocol := "http"
	// if cfg.SSL {
	// 	protocol = "https"
	// }

	httpClient := &http.Client{Timeout: 30 * time.Second}

	// Create the go-subsonic client
	client := &gosonic.Client{
		Client:       httpClient,
		BaseUrl:      config.BaseURL,
		User:         config.Username,
		ClientName:   "suasor",
		UserAgent:    "Suasor/1.0",
		PasswordAuth: true, // Using plain password auth for simplicity
	}

	// Authenticate with the Subsonic server
	err := client.Authenticate(config.Password)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to authenticate with Subsonic server")
	} else {
		log.Info().Msg("Successfully authenticated with Subsonic server")
	}

	return &SubsonicClient{
		BaseMediaClient: media.BaseMediaClient{
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: types.MediaClientTypeSubsonic.AsCategory(),
				Config:   &config,
			},
			ClientType: types.MediaClientTypeSubsonic,
		},
		httpClient: httpClient,
		client:     client,
	}, nil
}

// Register the provider factory
// func init() {
// 	media.RegisterClient(types.MediaClientTypeSubsonic, NewSubsonicClient)
// }

// Capability methods - Subsonic only supports music
func (c *SubsonicClient) SupportsMusic() bool       { return true }
func (c *SubsonicClient) SupportsPlaylists() bool   { return true }
func (c *SubsonicClient) SupportsMovies() bool      { return false }
func (c *SubsonicClient) SupportsTVShows() bool     { return false }
func (c *SubsonicClient) SupportsBooks() bool       { return false }
func (c *SubsonicClient) SupportsCollections() bool { return false }

func (c *SubsonicClient) TestConnection(ctx context.Context) (bool, error) {
	isUp := c.client.Ping()
	return isUp, nil
}
