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

func (e *EmbyClient) SupportsMusic() bool { return true }

// GetMusic retrieves music tracks from the Emby server
func (e *EmbyClient) GetMusicTracks(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Msg("Retrieving music tracks from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Audio"),
		Fields:           optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSourceCount,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines,ProviderIds"),
		Recursive:        optional.NewBool(true),
	}

	ApplyClientQueryOptions(ctx, &queryParams, options)

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
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

func (e *EmbyClient) GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Msg("Retrieving music artists from Emby server")

	opts := embyclient.ArtistsServiceApiGetArtistsOpts{
		Recursive: optional.NewBool(true),
		Fields:    optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSourceCount,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines,ProviderIds"),
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
			Str("baseURL", e.embyConfig().GetBaseURL()).
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

func (e *EmbyClient) GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Msg("Retrieving music albums from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Album"),
		Fields:           optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSourceCount,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines,ProviderIds"),
		Recursive:        optional.NewBool(true),
	}

	ApplyClientQueryOptions(ctx, &queryParams, options)

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
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
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("trackID", id).
		Msg("Retrieving specific music track from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids:              optional.NewString(id),
		Fields:           optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSourceCount,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines,ProviderIds"),
		IncludeItemTypes: optional.NewString("Audio"),
	}

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
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

func (e *EmbyClient) GetMusicArtistByID(ctx context.Context, id string) (*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("artistID", id).
		Msg("Retrieving specific music artist from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids:              optional.NewString(id),
		Fields:           optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSourceCount,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines,ProviderIds"),
		IncludeItemTypes: optional.NewString("MusicArtist"),
	}

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("artistID", id).
			Msg("Failed to fetch music artist from Emby")
		return &models.MediaItem[*types.Artist]{}, fmt.Errorf("failed to fetch music artist: %w", err)
	}

	if len(items.Items) == 0 {
		log.Error().
			Str("artistID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No music artist found with the specified ID")
		return &models.MediaItem[*types.Artist]{}, fmt.Errorf("music artist with ID %s not found", id)
	}

	item := items.Items[0]
	if item.Type_ != "MusicArtist" {
		log.Error().
			Str("artistID", id).
			Str("actualType", item.Type_).
			Msg("Item with specified ID is not a music artist")
		return &models.MediaItem[*types.Artist]{}, fmt.Errorf("item with ID %s is not a music artist", id)
	}

	itemArtist, err := GetItem[*types.Artist](ctx, e, &item)
	mediaItemArtist, err := GetMediaItem[*types.Artist](ctx, e, itemArtist, item.Id)

	return mediaItemArtist, nil
}

func (e *EmbyClient) GetMusicAlbumByID(ctx context.Context, id string) (*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("albumID", id).
		Msg("Retrieving specific music album from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids:              optional.NewString(id),
		Fields:           optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSourceCount,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines,ProviderIds"),
		IncludeItemTypes: optional.NewString("MusicAlbum"),
	}

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("albumID", id).
			Msg("Failed to fetch music album from Emby")
		return &models.MediaItem[*types.Album]{}, fmt.Errorf("failed to fetch music album: %w", err)
	}

	if len(items.Items) == 0 {
		log.Error().
			Str("albumID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No music album found with the specified ID")
		return &models.MediaItem[*types.Album]{}, fmt.Errorf("music album with ID %s not found", id)
	}

	item := items.Items[0]
	if item.Type_ != "MusicAlbum" {
		log.Error().
			Str("albumID", id).
			Str("actualType", item.Type_).
			Msg("Item with specified ID is not a music album")
		return &models.MediaItem[*types.Album]{}, fmt.Errorf("item with ID %s is not a music album", id)
	}

	itemAlbum, err := GetItem[*types.Album](ctx, e, &item)
	mediaItemAlbum, err := GetMediaItem[*types.Album](ctx, e, itemAlbum, item.Id)

	return mediaItemAlbum, nil
}

// GetMusicGenres retrieves music genres from the Emby server
func (e *EmbyClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Msg("Retrieving music genres from Emby server")

	opts := embyclient.MusicGenresServiceApiGetMusicgenresOpts{}

	result, resp, err := e.client.MusicGenresServiceApi.GetMusicgenres(ctx, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
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
