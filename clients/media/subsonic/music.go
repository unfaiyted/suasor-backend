package subsonic

import (
	"context"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (c *SubsonicClient) GetMusicTracks(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving music tracks from Subsonic server")

	var tracks []*models.MediaItem[*types.Track]
	var err error

	// If query or typed filters provided, use search3
	if options != nil && (options.Query != "" || hasAnyTypedFilter(options)) {
		queryString := buildQueryString(options)
		log.Info().
			Str("query", queryString).
			Msg("Searching for music tracks")
		tracks, err = c.searchTracks(ctx, *options)
	}

	if err != nil {
		return nil, err
	}

	log.Info().
		Int("trackCount", len(tracks)).
		Msg("Successfully retrieved music from Subsonic")

	return tracks, nil
}
