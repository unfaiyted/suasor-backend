package plex

import (
	"context"
	"fmt"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"github.com/LukeHagar/plexgo/models/operations"
)

// GetPlaylists retrieves playlists from Plex
func (c *PlexClient) GetPlaylists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

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
		return nil, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved playlists from Plex")

	playlists, err := GetMediaItemListFromPlaylist[*types.Playlist](ctx, c, res.Object.MediaContainer.Metadata)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("playlistsReturned", len(playlists)).
		Msg("Completed GetPlaylists request")

	return playlists, nil
}
