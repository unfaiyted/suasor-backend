package subsonic

import (
	"context"
	"fmt"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	mediatypes "suasor/clients/media/types"
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

func (c *SubsonicClient) GetMusicTracksByAlbumID(ctx context.Context, albumID string) ([]*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("albumID", albumID).
		Msg("Retrieving tracks from album")

	// Get album details including tracks using the API method
	album, err := c.GetMusicAlbumByID(albumID)
	if err != nil {
		log.Error().
			Err(err).
			Str("albumID", albumID).
			Msg("Failed to retrieve album details from Subsonic")
		return nil, fmt.Errorf("failed to retrieve album details: %w", err)
	}

	if album.Song == nil || len(album.Song) == 0 {
		log.Info().
			Str("albumID", albumID).
			Msg("No tracks found in album")
		return []*models.MediaItem[*types.Track]{}, nil
	}

	// Convert the tracks to MediaItems
	var tracks []*models.MediaItem[*types.Track]
	for _, track := range album.Song {
		trackItem, err := GetTrackItem(ctx, c, track)
		if err != nil {
			log.Warn().
				Err(err).
				Str("trackID", track.ID).
				Str("trackTitle", track.Title).
				Msg("Error converting track to MediaItem")
			continue
		}
		tracks = append(tracks, trackItem)
	}

	log.Info().
		Str("albumID", albumID).
		Int("trackCount", len(tracks)).
		Msg("Successfully retrieved tracks from album")

	return tracks, nil
}

func (c *SubsonicClient) GetMusicTrackByID(ctx context.Context, trackID string) (*models.MediaItem[*mediatypes.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("trackID", trackID).
		Msg("Retrieving specific music track from Subsonic server")
	// Call Subsonic getSong endpoint
	params := map[string]string{"id": trackID}
	resp, err := c.client.Get("getSong", params)
	if err != nil {
		log.Error().Err(err).Str("trackID", trackID).Msg("Failed to fetch music track from Subsonic")
		return nil, fmt.Errorf("failed to fetch music track: %w", err)
	}
	// Ensure a track was returned
	if resp.Song == nil {
		return nil, fmt.Errorf("music track with ID %s not found", trackID)
	}
	// Convert to MediaItem
	return GetTrackItem(ctx, c, resp.Song)
}
