package subsonic

import (
	"context"
	"fmt"
	"strconv"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
)

type searchType string

const (
	searchTypeMusic  searchType = "music"
	searchTypeAlbum  searchType = "album"
	searchTypeArtist searchType = "artist"
)

func (c *SubsonicClient) search(ctx context.Context, searchType searchType, options types.QueryOptions) (*gosonic.SearchResult3, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	limit := options.Limit
	offset := options.Offset

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("query", options.Query).
		Int("limit", options.Limit).
		Msg("Searching for music tracks")

	// Default to 50 items if not specified
	if limit <= 0 {
		limit = 50
	}

	// Call the Subsonic search3 API with proper parameters
	params := map[string]string{
		"query":        options.Query,
		"artistCount":  "0",
		"albumCount":   "0",
		"songCount":    "0",
		"songOffset":   "0",
		"albumOffset":  "0",
		"artistOffset": "0",
	}

	if searchType == searchTypeMusic {
		params["songOffset"] = strconv.Itoa(offset)
		params["songCount"] = strconv.Itoa(limit)
	} else if searchType == searchTypeAlbum {
		params["albumOffset"] = strconv.Itoa(offset)
		params["albumCount"] = strconv.Itoa(limit)
	} else if searchType == searchTypeArtist {
		params["artistOffset"] = strconv.Itoa(offset)
		params["artistCount"] = strconv.Itoa(limit)
	}

	return c.client.Search3(options.Query, params)
}

func (c *SubsonicClient) searchTracks(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	searchType := searchTypeMusic

	results, err := c.search(ctx, searchType, options)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", options.Query).
			Msg("Failed to search for music tracks")
		return nil, fmt.Errorf("failed to search for music tracks: %w", err)
	}

	if results == nil || results.Song == nil || len(results.Song) == 0 {
		log.Info().
			Str("query", options.Query).
			Msg("No music tracks found matching query")
		return []*models.MediaItem[*types.Track]{}, nil
	}

	// Convert the search results to MediaItems
	var tracks []*models.MediaItem[*types.Track]
	for _, track := range results.Song {
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
		Str("query", options.Query).
		Int("trackCount", len(tracks)).
		Msg("Successfully searched for music tracks")

	return tracks, nil

}

func (c *SubsonicClient) searchAlbums(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	searchType := searchTypeAlbum

	results, err := c.search(ctx, searchType, options)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", options.Query).
			Msg("Failed to search for albums")
		return nil, fmt.Errorf("failed to search for albums: %w", err)
	}

	if results == nil || results.Album == nil || len(results.Album) == 0 {
		log.Info().
			Str("query", options.Query).
			Msg("No albums found matching query")
		return []*models.MediaItem[*types.Album]{}, nil
	}

	// Convert the search results to MediaItems
	var albums []*models.MediaItem[*types.Album]
	for _, album := range results.Album {
		albumItem, err := GetAlbumItem(ctx, c, album)
		if err != nil {
			log.Warn().
				Err(err).
				Str("albumID", album.ID).
				Str("albumName", album.Name).
				Msg("Error converting album to MediaItem")
			continue
		}
		albums = append(albums, albumItem)
	}

	log.Info().
		Str("query", options.Query).
		Int("albumCount", len(albums)).
		Msg("Successfully searched for albums")

	return albums, nil
}

func (c *SubsonicClient) searchArtists(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)
	searchType := searchTypeArtist

	results, err := c.search(ctx, searchType, options)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", options.Query).
			Msg("Failed to search for artists")
		return nil, fmt.Errorf("failed to search for artists: %w", err)
	}

	if results == nil || results.Artist == nil || len(results.Artist) == 0 {
		log.Info().
			Str("query", options.Query).
			Msg("No artists found matching query")
		return []*models.MediaItem[*types.Artist]{}, nil
	}

	// Convert the search results to MediaItems
	var artists []*models.MediaItem[*types.Artist]
	for _, artist := range results.Artist {
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
		Str("query", options.Query).
		Int("artistCount", len(artists)).
		Msg("Successfully searched for artists")

	return artists, nil
}
