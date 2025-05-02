package subsonic

import (
	"context"
	"fmt"
	"strings"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
)

// GetArtist returns an Artist by ID.
func (c *SubsonicClient) GetMusicArtistByID(ctx context.Context, artistID string) (*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("artistID", artistID).
		Msg("Retrieving specific artist from Subsonic server")

	resp, err := c.client.GetArtist(artistID)
	if err != nil {
		log.Error().Err(err).Str("artistID", artistID).Msg("Failed to fetch artist from Subsonic")
		return nil, err
	}

	// Convert to MediaItem
	log.Debug().
		Str("artistID", artistID).
		Msg("Converting Subsonic artist to MediaItem")
	return GetArtistItem(ctx, c, resp)

}

// GetArtists retrieves all artists in the server.
func (c *SubsonicClient) GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Artist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving artists from Subsonic server")

	// Get artists using the API method
	artistsResponse, err := c.client.GetArtists(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to retrieve artists from Subsonic")
		return nil, fmt.Errorf("failed to retrieve artists: %w", err)
	}

	if artistsResponse.Index == nil || len(artistsResponse.Index) == 0 {
		log.Info().Msg("No artists found in Subsonic")
		return []*models.MediaItem[*types.Artist]{}, nil
	}

	// Collect all artists from all indexes
	var allArtists []*gosonic.ArtistID3
	for _, index := range artistsResponse.Index {
		if index.Artist == nil {
			continue
		}
		allArtists = append(allArtists, index.Artist...)
	}

	// Apply limit if specified
	limit := len(allArtists)
	if options != nil && options.Limit > 0 && options.Limit < limit {
		limit = options.Limit
	}

	// Filter by query if specified
	filteredArtists := allArtists
	if options != nil && options.Query != "" {
		query := strings.ToLower(options.Query)
		filteredArtists = nil
		for _, artist := range allArtists {
			if strings.Contains(strings.ToLower(artist.Name), query) {
				filteredArtists = append(filteredArtists, artist)
			}
		}
	}

	// Apply limit after filtering
	if len(filteredArtists) > limit {
		filteredArtists = filteredArtists[:limit]
	}

	// Convert the artists to MediaItems
	var artists []*models.MediaItem[*types.Artist]
	for _, artist := range filteredArtists {
		artistItem, err := GetArtistItem(ctx, c, artist)
		if err != nil {
			log.Warn().
				Err(err).
				Str("artistID", artist.ID).
				Str("artistName", artist.Name).
				Msg("Error converting artist to MediaItem")
			continue
		}
		artists = append(artists, artistItem)
	}

	log.Info().
		Int("artistCount", len(artists)).
		Msg("Successfully retrieved artists from Subsonic")

	return artists, nil
}
