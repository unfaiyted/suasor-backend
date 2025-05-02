package subsonic

import (
	"context"
	"fmt"
	"net/http"
	"suasor/clients/media"
	clienttypes "suasor/clients/types"
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

func (c *SubsonicClient) subsonicConfig() *clienttypes.SubsonicConfig {
	cfg := c.GetConfig().(*clienttypes.SubsonicConfig)
	return cfg
}

func (c *SubsonicClient) TestConnection(ctx context.Context) (bool, error) {
	isUp := c.client.Ping()
	return isUp, nil
}

// GetMusicGenres retrieves music genres from Subsonic
func (c *SubsonicClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving music genres from Subsonic server")
	// Call Subsonic getGenres endpoint
	resp, err := c.client.Get("getGenres", nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch music genres from Subsonic")
		return nil, fmt.Errorf("failed to fetch music genres: %w", err)
	}
	// Return list (may be empty)
	if resp.Genres == nil || resp.Genres.Genre == nil {
		return []string{}, nil
	}

	// loop over genres and remove any duplicates
	genreStrArr := make([]string, 0, len(resp.Genres.Genre))
	for _, genre := range resp.Genres.Genre {
		genreStrArr = append(genreStrArr, genre.Name)
	}

	return genreStrArr, nil
}
