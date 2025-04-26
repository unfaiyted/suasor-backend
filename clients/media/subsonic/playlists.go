package subsonic

import (
	"context"
	t "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (c *SubsonicClient) GetPlaylists(ctx context.Context, options *t.QueryOptions) ([]*models.MediaItem[*t.Playlist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving playlists from Subsonic")

	resp, err := c.client.Get("getPlaylists", nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch playlists from Subsonic")
		return nil, err
	}

	if resp.Playlists == nil || len(resp.Playlists.Playlist) == 0 {
		log.Info().Msg("No playlists returned from Subsonic")
		return []*models.MediaItem[*t.Playlist]{}, nil
	}

	playlists := make([]*models.MediaItem[*t.Playlist], 0, len(resp.Playlists.Playlist))

	for _, pl := range resp.Playlists.Playlist {
		playlistItem, err := GetPlaylistItem(ctx, c, pl)
		if err != nil {
			log.Error().
				Err(err).
				Str("playlistID", pl.ID).
				Msg("Failed to convert playlist")
			continue
		}

		playlists = append(playlists, playlistItem)
	}

	log.Info().
		Int("playlistCount", len(playlists)).
		Msg("Successfully retrieved playlists from Subsonic")

	return playlists, nil
}

func (c *SubsonicClient) GetPlaylistItems(ctx context.Context, playlistID string) ([]*models.MediaItem[*t.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist items from Subsonic")

	params := make(map[string]string)
	params["id"] = playlistID

	resp, err := c.client.Get("getPlaylist", params)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to fetch playlist items from Subsonic")
		return nil, err
	}

	if resp.Playlist == nil || len(resp.Playlist.Entry) == 0 {
		log.Info().
			Str("playlistID", playlistID).
			Msg("No tracks found in playlist")
		return []*models.MediaItem[*t.Track]{}, nil
	}

	tracks := make([]*models.MediaItem[*t.Track], 0, len(resp.Playlist.Entry))

	for _, song := range resp.Playlist.Entry {
		track, err := GetTrackItem(ctx, c, song)
		if err != nil {
			log.Error().
				Err(err).
				Str("trackID", song.ID).
				Msg("Failed to convert track")
			continue
		}
		tracks = append(tracks, track)
	}

	log.Info().
		Int("trackCount", len(tracks)).
		Str("playlistID", playlistID).
		Msg("Successfully retrieved playlist items from Subsonic")

	return tracks, nil
}

