package subsonic

import (
	"context"
	t "suasor/client/media/types"
	"suasor/models"
	"suasor/utils"
	"time"
)

func (c *SubsonicClient) GetPlaylists(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[t.Playlist], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
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
		return []models.MediaItem[t.Playlist]{}, nil
	}

	playlists := make([]models.MediaItem[t.Playlist], 0, len(resp.Playlists.Playlist))

	for _, pl := range resp.Playlists.Playlist {
		playlist := models.MediaItem[t.Playlist]{
			Type: "playlist",
			Data: t.Playlist{
				Details: t.MediaMetadata{
					Title:       pl.Name,
					Description: pl.Comment,
					Duration:    time.Duration(pl.Duration) * time.Second,
				},
				ItemCount: pl.SongCount,
				Owner:     pl.Owner,
				IsPublic:  pl.Public,
			},
		}
		playlist.SetClientInfo(c.ClientID, c.ClientType, pl.ID)

		// Add cover art if available
		if pl.CoverArt != "" {
			coverURL := c.GetCoverArtURL(pl.CoverArt)
			playlist.Data.Details.Artwork.Poster = coverURL
		}

		playlists = append(playlists, playlist)
	}

	log.Info().
		Int("playlistCount", len(playlists)).
		Msg("Successfully retrieved playlists from Subsonic")

	return playlists, nil
}

func (c *SubsonicClient) GetPlaylistItems(ctx context.Context, playlistID string) ([]models.MediaItem[t.Track], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
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
		return []models.MediaItem[t.Track]{}, nil
	}

	tracks := make([]models.MediaItem[t.Track], 0, len(resp.Playlist.Entry))

	for _, song := range resp.Playlist.Entry {
		track := c.convertChildToTrack(*song)
		track.SetClientInfo(c.ClientID, c.ClientType, song.ID)
		tracks = append(tracks, track)
	}

	log.Info().
		Int("trackCount", len(tracks)).
		Str("playlistID", playlistID).
		Msg("Successfully retrieved playlist items from Subsonic")

	return tracks, nil
}
