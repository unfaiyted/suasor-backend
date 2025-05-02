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

var defaultExtraFields = []jellyfin.ItemFields{
	jellyfin.ITEMFIELDS_DATE_CREATED,
	jellyfin.ITEMFIELDS_GENRES,
	jellyfin.ITEMFIELDS_PROVIDER_IDS,
	jellyfin.ITEMFIELDS_ORIGINAL_TITLE,
	jellyfin.ITEMFIELDS_AIR_TIME,
	jellyfin.ITEMFIELDS_EXTERNAL_URLS,
	jellyfin.ITEMFIELDS_STUDIOS,
}

func (j *JellyfinClient) SupportsMovies() bool { return true }

func (j *JellyfinClient) GetMovies(ctx context.Context, options *t.QueryOptions) ([]*models.MediaItem[*t.Movie], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.jellyfinConfig().GetBaseURL()).
		Msg("Retrieving movies from Jellyfin server")

	// Set up query parameters

	// Include movie type in the query
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}
	mediaTypes := []jellyfin.MediaType{jellyfin.MEDIATYPE_VIDEO}
	fields := defaultExtraFields
	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		IsMovie(true).
		Recursive(true).
		Fields(fields).
		MediaTypes(mediaTypes)

	// Set user ID first if available to ensure it's never nil
	if j.getUserID() != "" {
		log.Debug().
			Str("userID", j.getUserID()).
			Msg("Setting user ID")
		itemsReq.UserId(j.getUserID())
	}

	// Then apply any additional options
	if queryOptions := NewJellyfinQueryOptions(options); queryOptions != nil {
		queryOptions.SetItemsRequest(&itemsReq)
	}

	log.Debug().
		Bool("Recursive", true).
		Msg("Api Request with options")

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.jellyfinConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Int("statusCode", resp.StatusCode).
			Msg("Failed to fetch movies from Jellyfin")
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved movies from Jellyfin")

	// Convert results to expected format
	movies := make([]*models.MediaItem[*t.Movie], 0)

	for _, item := range result.Items {
		log.Info().
			Str("itemType", string(*item.Type)).
			Msg("Processing item")
		if *item.Type == jellyfin.BASEITEMKIND_MOVIE {
			itemMovie, err := GetItem[*t.Movie](ctx, j, &item)
			movie, err := GetMediaItem[*t.Movie](ctx, j, itemMovie, *item.Id)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("movieID", *item.Id).
					Str("movieName", *item.Name.Get()).
					Msg("Error converting Jellyfin item to movie format")
				continue
			}
			movies = append(movies, movie)
		}
	}

	log.Info().
		Int("moviesReturned", len(movies)).
		Msg("Completed GetMovies request")

	return movies, nil
}
func (j *JellyfinClient) GetMovieByID(ctx context.Context, id string) (*models.MediaItem[*t.Movie], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("movieID", id).
		Str("baseURL", j.jellyfinConfig().GetBaseURL()).
		Msg("Retrieving specific movie from Jellyfin server")

		// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}

	ids := id
	// fields := "ProductionYear,PremiereDate,ChannelMappingInfo,DateCreated,Genres,IndexOptions,HomePageUrl,Overview,ParentId,Path,ProviderIds,Studios,SortName"

	// Call the Jellyfin API
	log.Debug().
		Str("movieID", id).
		Msg("Making API request to Jellyfin server")

	fields := []jellyfin.ItemFields{
		jellyfin.ITEMFIELDS_DATE_CREATED,
		jellyfin.ITEMFIELDS_GENRES,
		jellyfin.ITEMFIELDS_PROVIDER_IDS,
		jellyfin.ITEMFIELDS_ORIGINAL_TITLE,
		jellyfin.ITEMFIELDS_AIR_TIME,
		jellyfin.ITEMFIELDS_EXTERNAL_URLS,
		jellyfin.ITEMFIELDS_STUDIOS,
	}

	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		Ids(strings.Split(ids, ",")).
		IncludeItemTypes(includeItemTypes).
		Fields(fields)

	// Set user ID if available
	if j.getUserID() != "" {
		itemsReq.UserId(j.getUserID())
	}

	result, resp, err := itemsReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.jellyfinConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("movieID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch movie from Jellyfin")
		return &models.MediaItem[*t.Movie]{}, fmt.Errorf("failed to fetch movie: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("movieID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No movie found with the specified ID")
		return &models.MediaItem[*t.Movie]{}, fmt.Errorf("movie with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is a movie
	if *item.Type != jellyfin.BASEITEMKIND_MOVIE {
		log.Error().
			Str("movieID", id).
			Str("actualType", string(*item.Type.Ptr())).
			Msg("Item with specified ID is not a movie")
		return &models.MediaItem[*t.Movie]{}, fmt.Errorf("item with ID %s is not a movie", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("movieID", id).
		Str("movieName", *item.Name.Get()).
		Msg("Successfully retrieved movie from Jellyfin")

	itemMovie, err := GetItem[*t.Movie](ctx, j, &item)
	movie, err := GetMediaItem[*t.Movie](ctx, j, itemMovie, *item.Id)
	if err != nil {
		log.Error().
			Err(err).
			Str("movieID", id).
			Str("movieName", *item.Name.Get()).
			Msg("Error converting Jellyfin item to movie format")
		return &models.MediaItem[*t.Movie]{}, fmt.Errorf("error converting movie data: %w", err)
	}

	log.Debug().
		Str("movieID", id).
		Str("movieName", movie.Data.Details.Title).
		Int("year", movie.Data.Details.ReleaseYear).
		Msg("Successfully returned movie data")

	return movie, nil
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
