package jellyfin

import (
	"context"
	"fmt"
	"strings"

	jellyfin "github.com/sj14/jellyfin-go/api"
	mediatype "suasor/clients/media/types"
	t "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// Array of music types
func (j *JellyfinClient) GetMusicTracks(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[*t.Track], error) {

	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving music tracks from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_AUDIO}

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music tracks")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true)

	// Set user ID if available
	if j.getUserID() != "" {
		itemsReq.UserId(j.getUserID())
	}

	NewJellyfinQueryOptions(options).
		SetItemsRequest(&itemsReq)

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch music tracks from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music tracks: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved music tracks from Jellyfin")

	// Convert results to expected format
	tracks := make([]models.MediaItem[*t.Track], 0)
	for _, item := range result.Items {
		if *item.Type == "Audio" {
			track := models.MediaItem[*t.Track]{
				Data: &t.Track{
					Details: &t.MediaDetails{
						Title:       *item.Name.Get(),
						Description: *item.Overview.Get(),
						Duration:    getDurationFromTicks(item.RunTimeTicks.Get()),
						Artwork:     *j.getArtworkURLs(&item),
					},
					Number: int(*item.IndexNumber.Get()),
				},
				Type: "track",
			}

			track.SetClientInfo(j.GetClientID(), j.GetClientType(), *item.Id)
			track.Data.AlbumName = *item.Album.Get()
			// Set album info if available
			if item.AlbumId.IsSet() && item.ArtistItems != nil && len(item.ArtistItems) > 0 {
				// TODO: check if we need to do something
			}

			// Add artist information if available
			if item.ArtistItems != nil && len(item.ArtistItems) > 0 {
				// track.Data.ArtistID = *item.ArtistItems[0].Id
				track.Data.ArtistName = *item.ArtistItems[0].Name.Get()
			}

			embedProviderIDs(ctx, &item.ProviderIds, &track.Data.Details.ExternalIDs)

			tracks = append(tracks, track)
		}
	}

	log.Info().
		Int("tracksReturned", len(tracks)).
		Msg("Completed GetMusic request")

	return tracks, nil
}

func (j *JellyfinClient) GetMusicArtists(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[*t.Artist], error) {

	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving music artists from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music artists")
	artistReq := j.client.ArtistsAPI.GetArtists(ctx)

	// Set user ID if available
	if j.getUserID() != "" {
		artistReq.UserId(j.getUserID())
	}

	NewJellyfinQueryOptions(options).
		SetArtistsRequest(&artistReq)

	result, resp, err := artistReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Artists").
			Int("statusCode", 0).
			Msg("Failed to fetch music artists from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music artists: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved music artists from Jellyfin")

	// Convert results to expected format
	artists := make([]models.MediaItem[*t.Artist], 0)

	for _, item := range result.Items {
		artist := models.MediaItem[*t.Artist]{
			Data: &t.Artist{
				Details: &t.MediaDetails{
					Title:       *item.Name.Get(),
					Description: *item.Overview.Get(),
					Artwork:     *j.getArtworkURLs(&item),
					Genres:      item.Genres,
				},
			},

			Type: "artist",
		}
		artist.SetClientInfo(j.GetClientID(), j.GetClientType(), *item.Id)

		embedProviderIDs(ctx, &item.ProviderIds, &artist.Data.Details.ExternalIDs)

		artists = append(artists, artist)
	}

	log.Info().
		Int("artistsReturned", len(artists)).
		Msg("Completed GetMusicArtists request")

	return artists, nil
}

func (j *JellyfinClient) GetMusicAlbums(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[*t.Album], error) {

	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving music albums from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MUSIC_ALBUM}

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music albums")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true)

	// Set user ID if available
	if j.getUserID() != "" {
		itemsReq.UserId(j.getUserID())
	}

	NewJellyfinQueryOptions(options).
		SetItemsRequest(&itemsReq)

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch music albums from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music albums: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved music albums from Jellyfin")

	// Convert results to expected format
	albums := make([]models.MediaItem[*t.Album], 0)
	for _, item := range result.Items {
		album := models.MediaItem[*t.Album]{
			Data: &t.Album{
				Details: &t.MediaDetails{
					Title:       *item.Name.Get(),
					Description: *item.Overview.Get(),
					ReleaseYear: int(*item.ProductionYear.Get()),
					Genres:      item.Genres,
					Artwork:     *j.getArtworkURLs(&item),
				},
				TrackCount: int(*item.ChildCount.Get()),
			},
			Type: "album",
		}

		album.SetClientInfo(j.GetClientID(), j.GetClientType(), *item.Id)

		// Set album artist if available
		if item.AlbumArtist.IsSet() {
			album.Data.ArtistName = *item.AlbumArtist.Get()
		}

		embedProviderIDs(ctx, &item.ProviderIds, &album.Data.Details.ExternalIDs)

		albums = append(albums, album)
	}

	log.Info().
		Int("albumsReturned", len(albums)).
		Msg("Completed GetMusicAlbums request")

	return albums, nil
}

// Single music type
func (j *JellyfinClient) GetMusicTrackByID(ctx context.Context, trackID string) (models.MediaItem[*t.Track], error) {

	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("trackID", trackID).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving specific music track from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().
		Str("trackID", trackID).
		Msg("Making API request to Jellyfin server")

	resultItems, err := j.getItemByIDs(ctx, trackID)
	if err != nil {
		return models.MediaItem[*t.Track]{}, err
	}

	// Check if any items were returned
	if len(resultItems) == 0 {
		log.Error().
			Str("trackID", trackID).
			Msg("No music track found with the specified ID")
		return models.MediaItem[*t.Track]{}, fmt.Errorf("music track with ID %s not found", trackID)
	}

	resultItem := resultItems[0]

	// Double-check that the returned item is a audio track
	if *resultItem.Type != jellyfin.BASEITEMKIND_AUDIO {
		log.Error().
			Str("trackID", trackID).
			Str("actualType", string(*resultItem.Type)).
			Msg("Item with specified ID is not a music track")
		return models.MediaItem[*t.Track]{}, fmt.Errorf("item with ID %s is not a music track", resultItem.Id)
	}

	log.Info().
		Str("trackID", trackID).
		Str("trackName", *resultItem.Name.Get()).
		Msg("Successfully retrieved music track from Jellyfin")

	track := models.MediaItem[*t.Track]{
		Data: &t.Track{
			Details: &t.MediaDetails{
				Title:       *resultItem.Name.Get(),
				Description: *resultItem.Overview.Get(),
				Duration:    getDurationFromTicks(resultItem.RunTimeTicks.Get()),
				Artwork:     *j.getArtworkURLs(&resultItem),
			},
			Number: int(*resultItem.IndexNumber.Get()),
		},
		Type: mediatype.MediaTypeTrack,
	}

	// Set album info if available
	if resultItem.AlbumId.IsSet() && resultItem.ArtistItems != nil {
		// TODO: check if we need to do something
	}

	if resultItem.Album.IsSet() {
		track.Data.AlbumName = *resultItem.Album.Get()
	}

	// Add artist information if available
	if resultItem.ArtistItems != nil && len(resultItem.ArtistItems) > 0 {
		track.Data.ArtistName = *resultItem.ArtistItems[0].Name.Get()
	}

	// Extract provider IDs
	embedProviderIDs(ctx, &resultItem.ProviderIds, &track.Data.Details.ExternalIDs)

	log.Debug().
		Str("trackID", trackID).
		Str("trackName", track.Data.Details.Title).
		Int("trackNumber", track.Data.Number).
		Msg("Successfully returned music track data")

	return track, nil
}

func (j *JellyfinClient) GetMusicArtistByID(ctx context.Context, artistID string) (*models.MediaItem[*t.Artist], error) {

	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("artistID", artistID).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving specific music track from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().
		Str("artistID", artistID).
		Msg("Making API request to Jellyfin server")

	resultItems, err := j.getItemByIDs(ctx, artistID)
	if err != nil {
		return &models.MediaItem[*t.Artist]{}, err
	}

	// Check if any items were returned
	if len(resultItems) == 0 {
		log.Error().
			Str("artistID", artistID).
			Msg("No music track found with the specified ID")
		return &models.MediaItem[*t.Artist]{}, fmt.Errorf("music track with ID %s not found", artistID)
	}

	resultItem := resultItems[0]

	// Double-check that the returned item is a audio track
	if *resultItem.Type != jellyfin.BASEITEMKIND_MUSIC_ARTIST {
		log.Error().
			Str("artistID", artistID).
			Str("actualType", string(*resultItem.Type)).
			Msg("Item with specified ID is not a music track")
		return &models.MediaItem[*t.Artist]{}, fmt.Errorf("item with ID %s is not a music track", resultItem.Id)
	}

	log.Info().
		Str("artistID", artistID).
		Str("trackName", *resultItem.Name.Get()).
		Msg("Successfully retrieved music track from Jellyfin")

	artist := t.Artist{
		Details: &t.MediaDetails{
			Title:       *resultItem.Name.Get(),
			Description: *resultItem.Overview.Get(),
			Duration:    getDurationFromTicks(resultItem.RunTimeTicks.Get()),
			Artwork:     *j.getArtworkURLs(&resultItem),
		},
	}

	mediaItem := models.NewMediaItem[*t.Artist](mediatype.MediaTypeArtist, &artist)

	// Set album info if available
	if resultItem.AlbumId.IsSet() && resultItem.ArtistItems != nil {
		// TODO: check if we need to do something
	}

	embedProviderIDs(ctx, &resultItem.ProviderIds, &artist.Details.ExternalIDs)

	log.Debug().
		Str("artistID", artistID).
		Str("trackName", artist.Details.Title).
		Msg("Successfully returned music track data")

	return mediaItem, nil
}

func (j *JellyfinClient) GetMusicAlbumByID(ctx context.Context, albumID string) (*models.MediaItem[*t.Album], error) {

	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("albumID", albumID).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving specific music album from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().
		Str("albumID", albumID).
		Msg("Making API request to Jellyfin server")

	resultItems, err := j.getItemByIDs(ctx, albumID)
	if err != nil {
		return &models.MediaItem[*t.Album]{}, err
	}

	// Check if any items were returned
	if len(resultItems) == 0 {
		log.Error().
			Str("albumID", albumID).
			Msg("No music album found with the specified ID")
		return &models.MediaItem[*t.Album]{}, fmt.Errorf("music album with ID %s not found", albumID)
	}

	resultItem := resultItems[0]

	// Double-check that the returned item is a audio track
	if *resultItem.Type != jellyfin.BASEITEMKIND_MUSIC_ALBUM {
		log.Error().
			Str("albumID", albumID).
			Str("actualType", string(*resultItem.Type)).
			Msg("Item with specified ID is not a music album")
		return &models.MediaItem[*t.Album]{}, fmt.Errorf("item with ID %s is not a music album", resultItem.Id)
	}
	log.Info().
		Str("albumID", albumID).
		Str("trackName", *resultItem.Name.Get()).
		Msg("Successfully retrieved music album from Jellyfin")

	album := t.Album{
		Details: &t.MediaDetails{
			Title:       *resultItem.Name.Get(),
			Description: *resultItem.Overview.Get(),
			ReleaseYear: int(*resultItem.ProductionYear.Get()),
			Genres:      resultItem.Genres,
			Artwork:     *j.getArtworkURLs(&resultItem),
		},
		TrackCount: int(*resultItem.ChildCount.Get()),
	}

	mediaItem := models.NewMediaItem[*t.Album](mediatype.MediaTypeAlbum, &album)
	mediaItem.SetClientInfo(j.GetClientID(), j.GetClientType(), *resultItem.Id)

	embedProviderIDs(ctx, &resultItem.ProviderIds, &album.Details.ExternalIDs)

	log.Debug().
		Str("albumID", albumID).
		Str("trackName", album.Details.Title).
		Msg("Successfully returned music album data")

	return mediaItem, nil
}

// GetMusicGenres retrieves music genres from the Jellyfin server
func (j *JellyfinClient) GetMusicGenres(ctx context.Context) ([]string, error) {

	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving music genres from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music genres")
	musicGenresReq := j.client.MusicGenresAPI.GetMusicGenres(ctx)
	result, resp, err := musicGenresReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/MusicGenres").
			Int("statusCode", 0).
			Msg("Failed to fetch music genres from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music genres: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved music genres from Jellyfin")

	// Convert results to expected format
	genres := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		if item.Name.Get() != nil {
			genres = append(genres, *item.Name.Get())
		}
	}

	log.Info().
		Int("genresReturned", len(genres)).
		Msg("Completed GetMusicGenres request")

	return genres, nil
}

// Helper function to get item by IDs
func (j *JellyfinClient) getItemByIDs(ctx context.Context, IDs string) ([]jellyfin.BaseItemDto, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving specific music track from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().
		Str("ID", IDs).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx)
	itemsReq.Ids(strings.Split(IDs, ","))

	// Set user ID if available
	if j.getUserID() != "" {
		itemsReq.UserId(j.getUserID())
	}

	result, resp, err := itemsReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("IDs", IDs).
			Int("statusCode", 0).
			Msg("Failed to fetch music track from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music track: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("IDs", IDs).
			Int("statusCode", resp.StatusCode).
			Msg("No music track found with the specified ID")
		return nil, fmt.Errorf("music track with ID %s not found", IDs)
	}

	return result.Items, nil
}
