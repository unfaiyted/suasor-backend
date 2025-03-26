package subsonic

import (
	"context"
	"fmt"
	"net/url"
	"suasor/utils"
)

// GetStreamURL returns the URL to stream a music track
func (c *SubsonicClient) GetStreamURL(ctx context.Context, trackID string) (string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", trackID).
		Msg("Generating stream URL for track")

	// We can't access the unexported setupRequest method, so build the URL manually
	protocol := "http"
	if c.config.SSL {
		protocol = "https"
	}

	// Create query parameters
	params := url.Values{}
	params.Add("id", trackID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", c.config.Username)
	params.Add("p", c.config.Password)

	streamURL := fmt.Sprintf("%s://%s:%d/rest/stream.view?%s",
		protocol, c.config.Host, c.config.Port, params.Encode())

	log.Debug().
		Str("trackID", trackID).
		Str("streamURL", streamURL).
		Msg("Generated stream URL for track")

	return streamURL, nil
}

// GetCoverArtURL returns the URL to download cover art
func (c *SubsonicClient) GetCoverArtURL(coverArtID string) string {
	if coverArtID == "" {
		return ""
	}

	// We can't access the unexported setupRequest method, so build the URL manually
	protocol := "http"
	if c.config.SSL {
		protocol = "https"
	}

	// Create query parameters
	params := url.Values{}
	params.Add("id", coverArtID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", c.config.Username)
	params.Add("p", c.config.Password)

	return fmt.Sprintf("%s://%s:%d/rest/getCoverArt.view?%s",
		protocol, c.config.Host, c.config.Port, params.Encode())
}
