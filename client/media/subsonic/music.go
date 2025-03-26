package subsonic

import (
	"context"
	"fmt"
	"strconv"
	t "suasor/client/media/types"
	"suasor/utils"
	"time"
)

func (c *SubsonicClient) GetMusic(ctx context.Context, options *t.QueryOptions) ([]t.MediaItem[t.Track], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music tracks from Subsonic server")

	var tracks []t.MediaItem[t.Track]
	var err error

	// If filter/query provided, use search3
	if options != nil && len(options.Filters) > 0 {

		queryString := ""
		for _, value := range options.Filters {
			if queryString != "" {
				queryString += " "
			}
			queryString += value
		}
		tracks, err = c.searchMusic(ctx, queryString, options.Limit)
	} else {
		// Otherwise get random songs
		tracks, err = c.getRandomSongs(ctx, options)
	}

	if err != nil {
		return nil, err
	}

	// // Add client info to each track
	// for i := range tracks {
	// 	tracks[i].SetClientInfo(c.ClientID, c.ClientType, *tracks[i].ID)
	// }

	log.Info().
		Int("trackCount", len(tracks)).
		Msg("Successfully retrieved music from Subsonic")

	return tracks, nil
}
func (c *SubsonicClient) GetMusicArtists(ctx context.Context, options *t.QueryOptions) ([]t.MediaItem[t.Artist], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music artists from Subsonic")

	resp, err := c.client.Get("getArtists", nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch artists from Subsonic")
		return nil, err
	}

	if resp.Artists == nil || len(resp.Artists.Index) == 0 {
		log.Info().Msg("No artists returned from Subsonic")
		return []t.MediaItem[t.Artist]{}, nil
	}

	var artists []t.MediaItem[t.Artist]

	// Flatten all artists from all indexes
	for _, index := range resp.Artists.Index {
		for _, artist := range index.Artist {
			// Apply pagination if needed
			if options != nil && options.Limit > 0 && len(artists) >= options.Limit {
				break
			}

			musicArtist := t.MediaItem[t.Artist]{
				Type: "artist",
				Data: t.Artist{
					Details: t.MediaMetadata{
						Title: artist.Name,
					},
				},
			}
			musicArtist.SetClientInfo(c.ClientID, c.ClientType, artist.ID)

			// Add cover art if available
			if artist.CoverArt != "" {
				musicArtist.Data.Details.Artwork.Poster = c.GetCoverArtURL(artist.CoverArt)
			}

			artists = append(artists, musicArtist)
		}
	}

	log.Info().
		Int("artistCount", len(artists)).
		Msg("Successfully retrieved music artists from Subsonic")

	return artists, nil
}

func (c *SubsonicClient) GetMusicTrackByID(ctx context.Context, id string) (t.MediaItem[t.Track], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", id).
		Msg("Retrieving specific music track from Subsonic")

	params := make(map[string]string)
	params["id"] = id

	resp, err := c.client.Get("getSong", params)
	if err != nil {
		log.Error().
			Err(err).
			Str("trackID", id).
			Msg("Failed to fetch track from Subsonic")
		return t.MediaItem[t.Track]{}, err
	}

	if resp.Song == nil {
		log.Error().
			Str("trackID", id).
			Msg("No track found with the specified ID")
		return t.MediaItem[t.Track]{}, fmt.Errorf("track with ID %s not found", id)
	}

	track := c.convertChildToTrack(*resp.Song)

	log.Info().
		Str("trackID", id).
		Str("title", track.Data.Details.Title).
		Str("artist", track.Data.ArtistName).
		Msg("Successfully retrieved music track from Subsonic")

	return track, nil
}

func (c *SubsonicClient) GetMusicAlbums(ctx context.Context, options *t.QueryOptions) ([]t.MediaItem[t.Album], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music albums from Subsonic")

	params := make(map[string]string)
	params["type"] = "alphabeticalByName"

	if options != nil {
		if options.Limit > 0 {
			params["size"] = strconv.Itoa(options.Limit)
		}
		if options.Offset > 0 {
			params["offset"] = strconv.Itoa(options.Offset)
		}
	}

	// Use getAlbumList2 which is tag-based instead of folder-based
	resp, err := c.client.Get("getAlbumList2", params)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch albums from Subsonic")
		return nil, err
	}

	if resp.AlbumList2 == nil || len(resp.AlbumList2.Album) == 0 {
		log.Info().Msg("No albums returned from Subsonic")
		return []t.MediaItem[t.Album]{}, nil
	}

	albums := make([]t.MediaItem[t.Album], 0, len(resp.AlbumList2.Album))

	for _, album := range resp.AlbumList2.Album {
		musicAlbum := t.MediaItem[t.Album]{
			Type: "album",
			Data: t.Album{
				Details: t.MediaMetadata{
					Title:       album.Name,
					ReleaseYear: album.Year,
					Duration:    time.Duration(album.Duration) * time.Second,
					Genres:      []string{album.Genre},
					Artwork: t.Artwork{
						Poster: c.GetCoverArtURL(album.CoverArt),
					},
				},
				ArtistName: album.Artist,
				TrackCount: album.SongCount,
			},
		}
		musicAlbum.SetClientInfo(c.ClientID, c.ClientType, album.ID)
		albums = append(albums, musicAlbum)
	}

	log.Info().
		Int("albumCount", len(albums)).
		Msg("Successfully retrieved music albums from Subsonic")

	return albums, nil
}

func (c *SubsonicClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music genres from Subsonic")

	resp, err := c.client.Get("getGenres", nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch genres from Subsonic")
		return nil, err
	}

	if resp.Genres == nil {
		log.Info().Msg("No genres returned from Subsonic")
		return []string{}, nil
	}

	genres := make([]string, 0, len(resp.Genres.Genre))
	for _, genre := range resp.Genres.Genre {
		genres = append(genres, genre.Name)
	}

	log.Info().
		Int("genreCount", len(genres)).
		Msg("Successfully retrieved music genres from Subsonic")

	return genres, nil
}

func (c *SubsonicClient) getRandomSongs(ctx context.Context, options *t.QueryOptions) ([]t.MediaItem[t.Track], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Debug().Msg("Fetching random songs from Subsonic")

	params := make(map[string]string)

	if options != nil && options.Limit > 0 {
		params["size"] = strconv.Itoa(options.Limit)
	}

	resp, err := c.client.Get("getRandomSongs", params)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch random songs from Subsonic")
		return nil, err
	}

	if resp.RandomSongs == nil || len(resp.RandomSongs.Song) == 0 {
		log.Info().Msg("No songs returned from Subsonic")
		return []t.MediaItem[t.Track]{}, nil
	}

	tracks := make([]t.MediaItem[t.Track], 0, len(resp.RandomSongs.Song))

	for _, song := range resp.RandomSongs.Song {
		track := c.convertChildToTrack(*song)
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (c *SubsonicClient) searchMusic(ctx context.Context, query string, limit int) ([]t.MediaItem[t.Track], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Str("query", query).
		Int("limit", limit).
		Msg("Searching for music in Subsonic")

	params := make(map[string]string)
	params["query"] = query

	if limit > 0 {
		params["songCount"] = strconv.Itoa(limit)
	}

	resp, err := c.client.Get("search3", params)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", query).
			Msg("Failed to search for music in Subsonic")
		return nil, err
	}

	if resp.SearchResult3 == nil || len(resp.SearchResult3.Song) == 0 {
		log.Info().
			Str("query", query).
			Msg("No songs found matching query")
		return []t.MediaItem[t.Track]{}, nil
	}

	tracks := make([]t.MediaItem[t.Track], 0, len(resp.SearchResult3.Song))

	for _, song := range resp.SearchResult3.Song {
		track := c.convertChildToTrack(*song)
		tracks = append(tracks, track)
	}

	return tracks, nil
}
