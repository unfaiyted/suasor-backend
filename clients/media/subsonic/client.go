package subsonic

import (
	"context"
	"net/http"
	"suasor/clients/media"
	clienttypes "suasor/clients/types"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"

	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
)

// SubsonicClient implements MediaContentProvider for Subsonic
type SubsonicClient struct {
	media.ClientMedia
	httpClient *http.Client
	client     *gosonic.Client
}

// NewSubsonicClient creates a new Subsonic client
func NewSubsonicClient(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config *clienttypes.SubsonicConfig) (media.ClientMedia, error) {
	log := logger.LoggerFromContext(context.Background())

	log.Info().
		Uint64("clientID", clientID).
		Str("baseURL", config.GetBaseURL()).
		Msg("Creating new Subsonic client")

	httpClient := &http.Client{Timeout: 30 * time.Second}

	// Create the go-subsonic client
	client := &gosonic.Client{
		Client:       httpClient,
		BaseUrl:      config.GetBaseURL(),
		User:         config.GetUsername(),
		ClientName:   "suasor",
		UserAgent:    "Suasor/1.0",
		PasswordAuth: true, // Using plain password auth for simplicity
	}

	// Authenticate with the Subsonic server
	err := client.Authenticate(config.GetPassword())
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to authenticate with Subsonic server")
	} else {
		log.Info().Msg("Successfully authenticated with Subsonic server")
	}

	// Create the client media interface
	clientMedia, err := media.NewClientMedia(ctx, clientID, clienttypes.ClientMediaTypeSubsonic, registry, config)
	if err != nil {
		return nil, err
	}

	// Create the Subsonic client
	subsonicClient := &SubsonicClient{
		ClientMedia: clientMedia,
		httpClient:  httpClient,
		client:      client,
	}

	return subsonicClient, nil
}

// Capability methods - Subsonic only supports music
func (c *SubsonicClient) SupportsMusic() bool       { return true }
func (c *SubsonicClient) SupportsPlaylists() bool   { return true }
func (c *SubsonicClient) SupportsMovies() bool      { return false }
func (c *SubsonicClient) SupportsSeries() bool      { return false }
func (c *SubsonicClient) SupportsCollections() bool { return false }
func (c *SubsonicClient) SupportsHistory() bool     { return true }

// GetArtists is an alias for backward compatibility
func (c *SubsonicClient) GetArtists(ctx context.Context, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Artist], error) {
	return c.GetArtistsWithContext(ctx, options)
}

func (c *SubsonicClient) subsonicConfig() *clienttypes.SubsonicConfig {
	cfg := c.GetConfig().(*clienttypes.SubsonicConfig)
	return cfg
}

func (c *SubsonicClient) TestConnection(ctx context.Context) (bool, error) {
	isUp := c.client.Ping()
	return isUp, nil
}