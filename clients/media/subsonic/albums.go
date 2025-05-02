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

// GetAlbums retrieves albums from Subsonic
func (c *SubsonicClient) GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Album], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving albums from Subsonic server")

	limit := 50
	offset := 0
	if options != nil && options.Limit > 0 {
		limit = options.Limit
	}
	if options != nil && options.Offset > 0 {
		offset = options.Offset
	}

	// Get albums list using our helper method
	params := map[string]string{
		"size":   strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
	}

	// List types:
	// "random"               ,
	// "newest"               ,
	// "highest"              ,
	// "frequent"             ,
	// "recent"              ,
	// "alphabeticalByName"   ,
	// "alphabeticalByArtist",
	// "starred"              ,
	// "byYear"               ,
	// "byGenre"              ,

	// Optional Parameters:
	//
	//	size:           The number of albums to return. Max 500, default 10.
	//	offset:         The list offset. Useful if you for example want to page through the list of newest albums.
	//	fromYear:       The first year in the range. If fromYear > toYear a reverse chronological list is returned.
	//	toYear:         The last year in the range.
	//	genre:          The name of the genre, e.g., "Rock".
	//	musicFolderId:  (Since 1.11.0) Only return albums in the music folder with the given ID. See getMusicFolders.
	//
	// toYear and fromYear are required parameters when type == "byYear". genre is a required parameter when type == "byGenre".

	albumsResult, err := c.client.GetAlbumList2("alphabeticalByArtist", params)
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

func (c *SubsonicClient) GetMusicAlbumsByArtistID(ctx context.Context, artistID string) ([]*models.MediaItem[*types.Album], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("artistID", artistID).
		Msg("Retrieving albums from artist")

	// Get artist details including albums using the API method
	artist, err := c.client.GetArtist(artistID)
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

// GetAlbum returns an Album by ID.
func (c *SubsonicClient) GetMusicAlbumByID(albumID string) (*gosonic.AlbumID3, error) {
	// Get album details from Subsonic
	resp, err := c.client.GetAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve album from Subsonic: %w", err)
	}

	return resp, nil
}
