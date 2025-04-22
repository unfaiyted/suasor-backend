// music.go
package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/clients/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"
	"suasor/utils/logger"
)

// GetMusic retrieves music tracks from the Emby server
func (e *EmbyClient) GetMusic(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving music tracks from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Audio"),
		Recursive:        optional.NewBool(true),
	}

	ApplyClientQueryOptions(&queryParams, options)

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().BaseURL).
			Str("apiEndpoint", "/Items").
			Msg("Failed to fetch music tracks from Emby")
		return nil, fmt.Errorf("failed to fetch music tracks: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved music tracks from Emby")

	tracks := make([]*models.MediaItem[*types.Track], 0)
	for _, item := range items.Items {

		itemTrack, err := GetItem[*types.Track](ctx, e, &item)
		mediaItemTrack, err := GetMediaItem[*types.Track](ctx, e, itemTrack, item.Id)

		if err != nil {
			log.Warn().
				Err(err).
				Str("trackID", item.Id).
				Str("trackName", item.Name).
				Msg("Error converting Emby item to music track format")
			continue
		}
		tracks = append(tracks, mediaItemTrack)
	}

	return tracks, nil
}

// GetMusicArtists retrieves music artists from the Emby server
func (e *EmbyClient) GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving music artists from Emby server")

	opts := embyclient.ArtistsServiceApiGetArtistsOpts{
		Recursive: optional.NewBool(true),
	}

	// Apply pagination and sorting
	if options != nil {
		if options.Limit > 0 {
			opts.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			opts.StartIndex = optional.NewInt32(int32(options.Offset))
		}
		if options.Sort != "" {
			// TODO: work on translating types to external sortBy,
			// they dont have any type definitions on this so we might need to look into it a bit
			opts.SortBy = optional.NewString(string(options.Sort))
			if options.SortOrder == "desc" {
				opts.SortOrder = optional.NewString("Descending")
			} else {
				opts.SortOrder = optional.NewString("Ascending")
			}
		}
	}

	result, resp, err := e.client.ArtistsServiceApi.GetArtists(ctx, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().BaseURL).
			Str("apiEndpoint", "/Artists").
			Msg("Failed to fetch music artists from Emby")
		return nil, fmt.Errorf("failed to fetch music artists: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(result.TotalRecordCount)).
		Msg("Successfully retrieved music artists from Emby")

	artists := make([]*models.MediaItem[*types.Artist], 0)
	for _, item := range result.Items {
		itemArtist, err := GetItem[*types.Artist](ctx, e, &item)
		mediaItemArtist, err := GetMediaItem[*types.Artist](ctx, e, itemArtist, item.Id)
		if err != nil {
			log.Warn().
				Err(err).
				Str("artistID", item.Id).
				Str("artistName", item.Name).
				Msg("Error converting Emby item to music artist format")
			continue
		}
		artists = append(artists, mediaItemArtist)
	}

	return artists, nil
}

// GetAlbums retrieves music albums from the Emby server
func (e *EmbyClient) GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving music albums from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Album"),
		Recursive:        optional.NewBool(true),
	}

	ApplyClientQueryOptions(&queryParams, options)

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().BaseURL).
			Str("apiEndpoint", "/Items").
			Msg("Failed to fetch music albums from Emby")
		return nil, fmt.Errorf("failed to fetch music albums: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved music albums from Emby")

	albums := make([]*models.MediaItem[*types.Album], 0)
	for _, item := range items.Items {
		itemAlbum, err := GetItem[*types.Album](ctx, e, &item)
		mediaItemAlbum, err := GetMediaItem[*types.Album](ctx, e, itemAlbum, item.Id)
		if err != nil {
			log.Warn().
				Err(err).
				Str("albumID", item.Id).
				Str("albumName", item.Name).
				Msg("Error converting Emby item to music album format")
			continue
		}
		albums = append(albums, mediaItemAlbum)
	}

	return albums, nil
}

// GetMusicTrackByID retrieves a specific music track by ID
func (e *EmbyClient) GetMusicTrackByID(ctx context.Context, id string) (*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("trackID", id).
		Msg("Retrieving specific music track from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids:              optional.NewString(id),
		IncludeItemTypes: optional.NewString("Audio"),
	}

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().BaseURL).
			Str("apiEndpoint", "/Items").
			Str("trackID", id).
			Msg("Failed to fetch music track from Emby")
		return &models.MediaItem[*types.Track]{}, fmt.Errorf("failed to fetch music track: %w", err)
	}

	if len(items.Items) == 0 {
		log.Error().
			Str("trackID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No music track found with the specified ID")
		return &models.MediaItem[*types.Track]{}, fmt.Errorf("music track with ID %s not found", id)
	}

	item := items.Items[0]
	if item.Type_ != "Audio" {
		log.Error().
			Str("trackID", id).
			Str("actualType", item.Type_).
			Msg("Item with specified ID is not a music track")
		return &models.MediaItem[*types.Track]{}, fmt.Errorf("item with ID %s is not a music track", id)
	}

	itemTrack, err := GetItem[*types.Track](ctx, e, &item)
	mediaItemTrack, err := GetMediaItem[*types.Track](ctx, e, itemTrack, item.Id)

	return mediaItemTrack, nil
}

// GetMusicGenres retrieves music genres from the Emby server
func (e *EmbyClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving music genres from Emby server")

	opts := embyclient.MusicGenresServiceApiGetMusicgenresOpts{}

	result, resp, err := e.client.MusicGenresServiceApi.GetMusicgenres(ctx, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().BaseURL).
			Str("apiEndpoint", "/MusicGenres").
			Msg("Failed to fetch music genres from Emby")
		return nil, fmt.Errorf("failed to fetch music genres: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("genreCount", len(result.Items)).
		Msg("Successfully retrieved music genres from Emby")

	genres := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		genres = append(genres, item.Name)
	}

	return genres, nil
}
