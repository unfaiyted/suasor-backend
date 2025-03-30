package subsonic

import (
	"context"
	"fmt"
	"net/url"
	types "suasor/client/types"
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

	subsonicConfig := c.Config.(*types.SubsonicConfig)

	// Create query parameters
	params := url.Values{}
	params.Add("id", trackID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", subsonicConfig.Username)
	params.Add("p", subsonicConfig.Password)

	streamURL := fmt.Sprintf("%s/rest/stream.view?%s",
		subsonicConfig.BaseURL, params.Encode())

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

	subsonicConfig := c.Config.(*types.SubsonicConfig)

	// Create query parameters
	params := url.Values{}
	params.Add("id", coverArtID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", subsonicConfig.Username)
	params.Add("p", subsonicConfig.Password)

	return fmt.Sprintf("%s/rest/getCoverArt.view?%s",
		subsonicConfig.BaseURL, params.Encode())
}
