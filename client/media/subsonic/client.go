package subsonic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"suasor/client/media/interfaces"
	"suasor/models"
	"time"
)

// SubsonicClient implements MediaContentProvider for Subsonic/Subsonic
type SubsonicClient struct {
	interfaces.BaseMediaClient
	config     models.SubsonicConfig
	httpClient *http.Client
	baseURL    string
}

// NewSubsonicClient creates a new Subsonic client
func NewSubsonicClient(clientID uint64, config models.SubsonicConfig) *SubsonicClient {
	protocol := "http"
	if config.SSL {
		protocol = "https"
	}
	baseURL := fmt.Sprintf("%s://%s:%d", protocol, config.Host, config.Port)

	return &SubsonicClient{
		BaseMediaClient: interfaces.BaseMediaClient{
			ClientID:   clientID,
			ClientType: models.MediaClientTypeSubsonic,
		},
		config:     config,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
	}
}

// Capability methods - Subsonic only supports music
func (c *SubsonicClient) SupportsMusic() bool     { return true }
func (c *SubsonicClient) SupportsPlaylists() bool { return true }

// Get basic auth parameters for Subsonic API
func (c *SubsonicClient) getAuthParams() url.Values {
	params := url.Values{}
	params.Add("u", c.config.Username)
	params.Add("p", c.config.Password)
	params.Add("v", "1.16.1")
	params.Add("c", "suasor")
	params.Add("f", "json")
	return params
}

// GetMusic retrieves music tracks
func (c *SubsonicClient) GetMusic(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicTrack, error) {
	params := c.getAuthParams()

	if options != nil && options.Limit > 0 {
		params.Add("size", fmt.Sprintf("%d", options.Limit))
	}

	reqURL := fmt.Sprintf("%s/rest/getRandomSongs.view?%s", c.baseURL, params.Encode())

	// Make API request and parse response...
	// This would be your actual implementation calling the Subsonic API

	// Example result (in real code, this would come from API response)
	tracks := []interfaces.MusicTrack{
		{
			MediaItem: interfaces.MediaItem{
				ID:   "song1",
				Type: "music",
				Metadata: interfaces.MediaMetadata{
					Title:  "Example Song",
					Genres: []string{"Rock"},
					ExternalIDs: interfaces.ExternalIDs{
						{Source: "musicbrainz", ID: "mb-123456"},
					},
				},
			},
			ArtistName: "Example Artist",
			AlbumTitle: "Example Album",
		},
	}

	// Add client info to each track
	for i := range tracks {
		c.AddClientInfo(&tracks[i].MediaItem)
	}

	return tracks, nil
}

// GetMusicGenres retrieves available music genres
func (c *SubsonicClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	params := c.getAuthParams()
	reqURL := fmt.Sprintf("%s/rest/getGenres.view?%s", c.baseURL, params.Encode())

	// Make API request and parse response...

	// Example result
	return []string{"Rock", "Pop", "Jazz", "Classical"}, nil
}

// Unsupported methods just inherit from BaseMediaClient
// All movie/TV methods will return ErrFeatureNotSupported
