package jellyfin

import (
	"context"
	"fmt"
	"strings"

	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (j *JellyfinClient) GetMusic(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[*t.Track], error) {
	// Get logger from context
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
		Recursive(true).
		UserId(j.config.UserID)

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
					Details: t.MediaDetails{
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
				track.Data.AddSyncClient(j.GetClientID(), *item.AlbumId.Get(), *item.ArtistItems[0].Id)
			}

			// Add artist information if available
			if item.ArtistItems != nil && len(item.ArtistItems) > 0 {
				// track.Data.ArtistID = *item.ArtistItems[0].Id
				track.Data.ArtistName = *item.ArtistItems[0].Name.Get()
			}

			extractProviderIDs(&item.ProviderIds, &track.Data.Details.ExternalIDs)

			tracks = append(tracks, track)
		}
	}

	log.Info().
		Int("tracksReturned", len(tracks)).
		Msg("Completed GetMusic request")

	return tracks, nil
}

func (j *JellyfinClient) GetMusicArtists(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[*t.Artist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving music artists from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music artists")
	artistReq := j.client.ArtistsAPI.GetArtists(ctx).
		UserId(j.config.UserID)

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
				Details: t.MediaDetails{
					Title:       *item.Name.Get(),
					Description: *item.Overview.Get(),
					Artwork:     *j.getArtworkURLs(&item),
					Genres:      item.Genres,
				},
			},

			Type: "artist",
		}
		artist.SetClientInfo(j.GetClientID(), j.GetClientType(), *item.Id)

		extractProviderIDs(&item.ProviderIds, &artist.Data.Details.ExternalIDs)

		artists = append(artists, artist)
	}

	log.Info().
		Int("artistsReturned", len(artists)).
		Msg("Completed GetMusicArtists request")

	return artists, nil
}

func (j *JellyfinClient) GetMusicAlbums(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[*t.Album], error) {
	// Get logger from context
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
		Recursive(true).
		UserId(j.config.UserID)

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
				Details: t.MediaDetails{
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

		extractProviderIDs(&item.ProviderIds, &album.Data.Details.ExternalIDs)

		albums = append(albums, album)
	}

	log.Info().
		Int("albumsReturned", len(albums)).
		Msg("Completed GetMusicAlbums request")

	return albums, nil
}

func (j *JellyfinClient) GetMusicTrackByID(ctx context.Context, id string) (models.MediaItem[*t.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("trackID", id).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving specific music track from Jellyfin server")

	// Set up query parameters
	ids := id

	// Call the Jellyfin API
	log.Debug().
		Str("trackID", id).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx)

	itemsReq.Ids(strings.Split(ids, ","))

	result, resp, err := itemsReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("trackID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch music track from Jellyfin")
		return models.MediaItem[*t.Track]{}, fmt.Errorf("failed to fetch music track: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("trackID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No music track found with the specified ID")
		return models.MediaItem[*t.Track]{}, fmt.Errorf("music track with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is a music track
	if *item.Type != "Audio" {
		log.Error().
			Str("trackID", id).
			Str("actualType", string(*item.Type)).
			Msg("Item with specified ID is not a music track")
		return models.MediaItem[*t.Track]{}, fmt.Errorf("item with ID %s is not a music track", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("trackID", id).
		Str("trackName", *item.Name.Get()).
		Msg("Successfully retrieved music track from Jellyfin")

	track := models.MediaItem[*t.Track]{
		Data: &t.Track{
			Details: t.MediaDetails{
				Title:       *item.Name.Get(),
				Description: *item.Overview.Get(),
				Duration:    getDurationFromTicks(item.RunTimeTicks.Get()),
				Artwork:     *j.getArtworkURLs(&item),
			},
			Number: int(*item.IndexNumber.Get()),
		},
		Type: "track",
	}

	// Set album info if available
	if item.AlbumId.IsSet() && item.ArtistItems != nil && len(item.ArtistItems) > 0 {
		track.Data.AddSyncClient(j.GetClientID(), *item.AlbumId.Get(), *item.ArtistItems[0].Id)
	}

	if item.Album.IsSet() {
		track.Data.AlbumName = *item.Album.Get()
	}

	// Add artist information if available
	if item.ArtistItems != nil && len(item.ArtistItems) > 0 {
		track.Data.ArtistName = *item.ArtistItems[0].Name.Get()
	}

	// Extract provider IDs
	extractProviderIDs(&item.ProviderIds, &track.Data.Details.ExternalIDs)

	log.Debug().
		Str("trackID", id).
		Str("trackName", track.Data.Details.Title).
		Int("trackNumber", track.Data.Number).
		Msg("Successfully returned music track data")

	return track, nil
}

func (j *JellyfinClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
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

func (j *JellyfinClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving movie genres from Jellyfin server")

	// Set up query parameters to get only movie genres
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}
	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for movie genres")
	genresReq := j.client.GenresAPI.GetGenres(ctx)

	genresReq.IncludeItemTypes(includeItemTypes)
	result, resp, err := genresReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Genres").
			Int("statusCode", 0).
			Msg("Failed to fetch movie genres from Jellyfin")
		return nil, fmt.Errorf("failed to fetch movie genres: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved movie genres from Jellyfin")

	// Convert results to expected format
	genres := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		if item.Name.Get() != nil {
			genres = append(genres, *item.Name.Get())
		}
	}

	log.Info().
		Int("genresReturned", len(genres)).
		Msg("Completed GetMovieGenres request")

	return genres, nil
}
