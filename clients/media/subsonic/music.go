package subsonic

import (
	"context"
	"fmt"
	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
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
		tracks, err = c.searchMusic(ctx, *options)
	} else {
		// If no query or filters, get a list of newest tracks
		// Use getAlbumList2 to get newest albums then fetch their tracks

		// params := map[string]string{
		// 	"size":   string(limit),
		// 	"offset": string(offset),
		// }

		// albums, err := c.GetAlbumList2("newest", params)
		// if err != nil {
		// 	log.Error().
		// 		Err(err).
		// 		Msg("Failed to retrieve albums for track listing from Subsonic")
		// 	return nil, fmt.Errorf("failed to retrieve albums for track listing: %w", err)
		// }

		// Collect tracks from each album
		// for _, album := range albums {
		// 	albumTracks, err := c.GetAlbumTracks(ctx, album.ID)
		// 	if err != nil {
		// 		log.Warn().
		// 			Err(err).
		// 			Str("albumID", album.ID).
		// 			Msg("Error retrieving tracks from album")
		// 		continue
		// 	}
		// 	tracks = append(tracks, albumTracks...)
		//
		// 	// Break if we have enough tracks
		// 	if limit > 0 && len(tracks) >= limit {
		// 		tracks = tracks[:limit]
		// 		break
		// 	}
		// }
	}

	if err != nil {
		return nil, err
	}

	log.Info().
		Int("trackCount", len(tracks)).
		Msg("Successfully retrieved music from Subsonic")

	return tracks, nil
}

// searchMusic searches for music tracks using the Subsonic search3 endpoint
func (c *SubsonicClient) searchMusic(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[*types.Track], error) {
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
		"query":       options.Query,
		"songOffset":  strconv.Itoa(offset),
		"songCount":   strconv.Itoa(limit),
		"artistCount": "0",
		"albumCount":  "0",
	}

	resp, err := c.client.Get("search3", params)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", options.Query).
			Msg("Failed to search for music tracks")
		return nil, fmt.Errorf("failed to search for music tracks: %w", err)
	}

	if resp.SearchResult3 == nil || resp.SearchResult3.Song == nil || len(resp.SearchResult3.Song) == 0 {
		log.Info().
			Str("query", options.Query).
			Msg("No music tracks found matching query")
		return []*models.MediaItem[*types.Track]{}, nil
	}

	// Convert the search results to MediaItems
	var tracks []*models.MediaItem[*types.Track]
	for _, track := range resp.SearchResult3.Song {
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

// GetAlbums retrieves albums from Subsonic
func (c *SubsonicClient) GetAlbums(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Album], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving albums from Subsonic server")

	limit := 50
	if options != nil && options.Limit > 0 {
		limit = options.Limit
	}

	// Get albums list using our helper method
	params := map[string]string{
		"size":   strconv.Itoa(limit),
		"offset": "0",
	}

	albumsResult, err := c.GetAlbumList2("newest", params)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to retrieve albums from Subsonic")
		return nil, fmt.Errorf("failed to retrieve albums: %w", err)
	}

	if albumsResult == nil || len(albumsResult) == 0 {
		log.Info().Msg("No albums found in Subsonic")
		return []*models.MediaItem[*types.Album]{}, nil
	}

	// Convert the albums to MediaItems
	var albums []*models.MediaItem[*types.Album]
	for _, album := range albumsResult {
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
		Int("albumCount", len(albums)).
		Msg("Successfully retrieved albums from Subsonic")

	return albums, nil
}

// GetAlbum returns an Album by ID.
func (c *SubsonicClient) GetAlbum(id string) (*gosonic.AlbumID3, error) {
	// Get album details from Subsonic
	params := map[string]string{
		"id": id,
	}
	resp, err := c.client.Get("getAlbum", params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve album from Subsonic: %w", err)
	}

	return resp.Album, nil
}

// GetArtistsWithContext retrieves artists from Subsonic with context
func (c *SubsonicClient) GetArtistsWithContext(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Artist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving artists from Subsonic server")

	// Get artists using the API method
	artistsResponse, err := c.GetArtistsAPI(nil)
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

// GetArtistsAPI returns all artists in the server.
func (c *SubsonicClient) GetArtistsAPI(parameters map[string]string) (*gosonic.ArtistsID3, error) {
	resp, err := c.client.Get("getArtists", parameters)
	if err != nil {
		return nil, err
	}
	return resp.Artists, nil
}

// GetAlbumList2 returns a list of albums using the getAlbumList2 endpoint.
//
// Parameters:
//
//	listType: Type of list, one of "newest", "frequent", "recent", "random", "alphabeticalByName", "alphabeticalByArtist", "starred", "byYear", "byGenre"
//	params: Optional parameters like size, offset, fromYear, toYear, genre, musicFolderId
func (c *SubsonicClient) GetAlbumList2(listType string, params map[string]string) ([]*gosonic.AlbumID3, error) {
	if params == nil {
		params = make(map[string]string)
	}

	// Set the type parameter
	params["type"] = listType

	resp, err := c.client.Get("getAlbumList2", params)
	if err != nil {
		return nil, err
	}

	if resp.AlbumList2 == nil {
		return []*gosonic.AlbumID3{}, nil
	}

	return resp.AlbumList2.Album, nil
}

// GetArtist returns an Artist by ID.
func (c *SubsonicClient) GetArtist(id string) (*gosonic.ArtistID3, error) {
	params := map[string]string{
		"id": id,
	}
	resp, err := c.client.Get("getArtist", params)
	if err != nil {
		return nil, err
	}
	return resp.Artist, nil
}

// GetAlbumTracks retrieves tracks from a specific album
func (c *SubsonicClient) GetAlbumTracks(ctx context.Context, albumID string) ([]*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("albumID", albumID).
		Msg("Retrieving tracks from album")

	// Get album details including tracks using the API method
	album, err := c.GetAlbum(albumID)
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

// GetArtistAlbums retrieves albums from a specific artist
func (c *SubsonicClient) GetArtistAlbums(ctx context.Context, artistID string) ([]*models.MediaItem[*types.Album], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("artistID", artistID).
		Msg("Retrieving albums from artist")

	// Get artist details including albums using the API method
	artist, err := c.GetArtist(artistID)
	if err != nil {
		log.Error().
			Err(err).
			Str("artistID", artistID).
			Msg("Failed to retrieve artist details from Subsonic")
		return nil, fmt.Errorf("failed to retrieve artist details: %w", err)
	}

	if artist.Album == nil || len(artist.Album) == 0 {
		log.Info().
			Str("artistID", artistID).
			Msg("No albums found for artist")
		return []*models.MediaItem[*types.Album]{}, nil
	}

	// Convert the albums to MediaItems
	var albums []*models.MediaItem[*types.Album]
	for _, album := range artist.Album {
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
		Str("artistID", artistID).
		Int("albumCount", len(albums)).
		Msg("Successfully retrieved albums from artist")

	return albums, nil
}

// Helper function to check if any typed filter is present in the options
func hasAnyTypedFilter(options *types.QueryOptions) bool {
	if options == nil {
		return false
	}

	return options.ClientAlbumID != "" ||
		options.ClientArtistID != "" ||
		options.Genre != "" ||
		options.Year != 0
}

// Helper function to build a search query string from options
func buildQueryString(options *types.QueryOptions) string {
	if options == nil {
		return ""
	}

	var parts []string

	// Add the main query
	if options.Query != "" {
		parts = append(parts, options.Query)
	}

	// Add artist filter
	if options.ClientArtistID != "" {
		parts = append(parts, "artist:"+options.ClientArtistID)
	}

	// Add album filter
	if options.ClientAlbumID != "" {
		parts = append(parts, "album:"+options.ClientAlbumID)
	}

	// Add genre filter
	if options.Genre != "" {
		parts = append(parts, "genre:"+options.Genre)
	}

	// Add year filter
	if options.Year != 0 {
		parts = append(parts, "year:"+strconv.Itoa(options.Year))
	}

	return strings.Join(parts, " ")
}
