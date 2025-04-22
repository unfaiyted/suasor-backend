package subsonic

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (c *SubsonicClient) GetMusic(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music tracks from Subsonic server")

	var tracks []*models.MediaItem[*types.Track]
	var err error

	// If query or typed filters provided, use search3
	if options != nil && (options.Query != "" || hasAnyTypedFilter(options)) {
		queryString := buildQueryString(options)
		tracks, err = c.searchMusic(ctx, queryString, options.Limit)
	}

	if err != nil {
		return nil, err
	}

	log.Info().
		Int("trackCount", len(tracks)).
		Msg("Successfully retrieved music from Subsonic")

	return tracks, nil
}
func (c *SubsonicClient) GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Artist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

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
		return []*models.MediaItem[*types.Artist]{}, nil
	}

	var artists []*models.MediaItem[*types.Artist]

	// Flatten all artists from all indexes
	for _, index := range resp.Artists.Index {
		for _, artist := range index.Artist {
			// Apply pagination if needed
			if options != nil && options.Limit > 0 && len(artists) >= options.Limit {
				break
			}

			artistItem, err := GetArtistItem(ctx, c, artist)
			if err != nil {
				log.Error().
					Err(err).
					Str("artistID", artist.ID).
					Msg("Failed to convert artist")
				continue
			}

			artists = append(artists, artistItem)
		}
	}

	log.Info().
		Int("artistCount", len(artists)).
		Msg("Successfully retrieved music artists from Subsonic")

	return artists, nil
}

func (c *SubsonicClient) GetMusicTrackByID(ctx context.Context, id string) (*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

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
		return &models.MediaItem[*types.Track]{}, err
	}

	if resp.Song == nil {
		log.Error().
			Str("trackID", id).
			Msg("No track found with the specified ID")
		return &models.MediaItem[*types.Track]{}, fmt.Errorf("track with ID %s not found", id)
	}

	track, err := GetTrackItem(ctx, c, resp.Song)
	if err != nil {
		log.Error().
			Err(err).
			Str("trackID", id).
			Msg("Failed to convert track")
		return &models.MediaItem[*types.Track]{}, err
	}

	log.Info().
		Str("trackID", id).
		Str("title", track.Data.Details.Title).
		Str("artist", track.Data.ArtistName).
		Msg("Successfully retrieved music track from Subsonic")

	return track, nil
}

func (c *SubsonicClient) GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Album], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

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
		return []*models.MediaItem[*types.Album]{}, nil
	}

	albums := make([]*models.MediaItem[*types.Album], 0, len(resp.AlbumList2.Album))

	for _, album := range resp.AlbumList2.Album {
		albumItem, err := GetAlbumItem(ctx, c, album)
		if err != nil {
			log.Error().
				Err(err).
				Str("albumID", album.ID).
				Msg("Failed to convert album")
			continue
		}
		albums = append(albums, albumItem)
	}

	log.Info().
		Int("albumCount", len(albums)).
		Msg("Successfully retrieved music albums from Subsonic")

	return albums, nil
}

func (c *SubsonicClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

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

func (c *SubsonicClient) getRandomSongs(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

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
		return []*models.MediaItem[*types.Track]{}, nil
	}

	tracks := make([]*models.MediaItem[*types.Track], 0, len(resp.RandomSongs.Song))

	for _, song := range resp.RandomSongs.Song {
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

	return tracks, nil
}

func (c *SubsonicClient) searchMusic(ctx context.Context, query string, limit int) ([]*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

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
		return []*models.MediaItem[*types.Track]{}, nil
	}

	tracks := make([]*models.MediaItem[*types.Track], 0, len(resp.SearchResult3.Song))

	for _, song := range resp.SearchResult3.Song {
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

	return tracks, nil
}

// Helper function to check if any typed filters are set
func hasAnyTypedFilter(options *types.QueryOptions) bool {
	if options == nil {
		return false
	}

	return options.Favorites ||
		options.Genre != "" ||
		options.Year > 0 ||
		options.Actor != "" ||
		options.Director != "" ||
		options.Studio != "" ||
		options.Creator != "" ||
		options.MediaType != "" ||
		options.ContentRating != "" ||
		len(options.Tags) > 0 ||
		options.RecentlyAdded ||
		options.RecentlyPlayed ||
		options.Watched ||
		!options.DateAddedAfter.IsZero() ||
		!options.DateAddedBefore.IsZero() ||
		!options.ReleasedAfter.IsZero() ||
		!options.ReleasedBefore.IsZero() ||
		!options.PlayedAfter.IsZero() ||
		!options.PlayedBefore.IsZero() ||
		options.MinimumRating > 0 ||
		options.ExternalSourceID != ""
}

// Helper function to build a query string from typed filters
func buildQueryString(options *types.QueryOptions) string {
	if options == nil {
		return ""
	}

	var queryParts []string

	// First add the direct search query if provided
	if options.Query != "" {
		queryParts = append(queryParts, options.Query)
	}

	// Add typed filters
	if options.MediaType != "" {
		queryParts = append(queryParts, string(options.MediaType))
	}

	if options.Genre != "" {
		queryParts = append(queryParts, options.Genre)
	}

	if options.Year > 0 {
		queryParts = append(queryParts, strconv.Itoa(options.Year))
	}

	if options.Actor != "" {
		queryParts = append(queryParts, options.Actor)
	}

	if options.Director != "" {
		queryParts = append(queryParts, options.Director)
	}

	if options.Creator != "" {
		queryParts = append(queryParts, options.Creator)
	}

	if options.Studio != "" {
		queryParts = append(queryParts, options.Studio)
	}

	if options.ContentRating != "" {
		queryParts = append(queryParts, options.ContentRating)
	}

	if len(options.Tags) > 0 {
		for _, tag := range options.Tags {
			queryParts = append(queryParts, tag)
		}
	}

	// Join all parts with spaces
	return strings.Join(queryParts, " ")
}
