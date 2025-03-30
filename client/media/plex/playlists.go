package plex

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
	"time"

	"github.com/LukeHagar/plexgo/models/operations"
)

// GetPlaylists retrieves playlists from Plex
func (c *PlexClient) GetPlaylists(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[types.Playlist], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving playlists from Plex server")

	log.Debug().Msg("Making API request to Plex server for playlists")
	res, err := c.plexAPI.Playlists.GetPlaylists(ctx, operations.PlaylistTypeAudio.ToPointer(), operations.QueryParamSmartOne.ToPointer())
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to get playlists from Plex")
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No playlists found in Plex")
		return []models.MediaItem[types.Playlist]{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved playlists from Plex")

	playlists := make([]models.MediaItem[types.Playlist], 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		playlist := models.MediaItem[types.Playlist]{
			ExternalID: *item.RatingKey,
			Data: types.Playlist{
				Details: types.MediaDetails{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork:     types.Artwork{
						// Thumbnail: c.makeFullURL(*item.Thumb),
					},
					ExternalIDs: types.ExternalIDs{types.ExternalID{
						Source: "plex",
						ID:     *item.RatingKey,
					}},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
		}
		playlist.SetClientInfo(c.ClientID, c.ClientType, *item.RatingKey)

		playlists = append(playlists, playlist)
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("playlistsReturned", len(playlists)).
		Msg("Completed GetPlaylists request")

	return playlists, nil
}
