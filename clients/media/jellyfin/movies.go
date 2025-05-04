package jellyfin

import (
	"context"
	"fmt"

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

	// Prepare all parameters first before building the request
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}
	mediaTypes := []jellyfin.MediaType{jellyfin.MEDIATYPE_VIDEO}
	fields := defaultExtraFields

	// Default values
	limit := int32(100)
	startIndex := int32(0)
	recursive := true
	isMovie := true

	// If there are any query options, prepare everything from them first
	if options != nil {
		log.Debug().
			Interface("options", options).
			Msg("Processing query options")

		// Handle limit
		if options.Limit > 0 {
			limit = int32(options.Limit)
			log.Debug().Int32("limit", limit).Msg("Using custom limit")
		}

		// Handle offset/startIndex
		if options.Offset > 0 {
			startIndex = int32(options.Offset)
			log.Debug().Int32("startIndex", startIndex).Msg("Using custom startIndex")
		}
	}

	// Log all parameters before creating the request
	log.Debug().
		Interface("includeItemTypes", includeItemTypes).
		Interface("mediaTypes", mediaTypes).
		Int32("limit", limit).
		Int32("startIndex", startIndex).
		Bool("recursive", recursive).
		Bool("isMovie", isMovie).
		Msg("Building Jellyfin API request with parameters")

	// Create the request with ALL parameters in a single builder chain
	requestBuilder := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		IsMovie(isMovie).
		Recursive(recursive).
		Fields(fields).
		MediaTypes(mediaTypes).
		Limit(limit).
		StartIndex(startIndex)

	// Add userId if available
	if j.getUserID() != "" {
		log.Debug().Str("userID", j.getUserID()).Msg("Adding user ID to request")
		requestBuilder = requestBuilder.UserId(j.getUserID())
	}

	// Apply additional query parameters if provided
	if options != nil {
		if options.Sort != "" {
			sortBy := []jellyfin.ItemSortBy{jellyfin.ItemSortBy(options.Sort)}
			log.Debug().Interface("sortBy", sortBy).Msg("Adding sort option")
			requestBuilder = requestBuilder.SortBy(sortBy)

			// Add sort order
			if options.SortOrder == "desc" {
				log.Debug().Msg("Setting descending sort order")
				requestBuilder = requestBuilder.SortOrder([]jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING})
			} else {
				log.Debug().Msg("Setting ascending sort order")
				requestBuilder = requestBuilder.SortOrder([]jellyfin.SortOrder{jellyfin.SORTORDER_ASCENDING})
			}
		}

		// Search term
		if options.Query != "" {
			log.Debug().Str("query", options.Query).Msg("Adding search term")
			requestBuilder = requestBuilder.SearchTerm(options.Query)
		}

		// Genre filter
		if options.Genre != "" {
			log.Debug().Str("genre", options.Genre).Msg("Adding genre filter")
			requestBuilder = requestBuilder.Genres([]string{options.Genre})
		}

		// Favorite filter
		if options.Favorites {
			log.Debug().Msg("Adding favorites filter")
			requestBuilder = requestBuilder.IsFavorite(true)
		}

		// Year filter
		if options.Year > 0 {
			log.Debug().Int("year", options.Year).Msg("Adding year filter")
			requestBuilder = requestBuilder.Years([]int32{int32(options.Year)})
		}

		// Person filters
		if options.Actor != "" {
			log.Debug().Str("actor", options.Actor).Msg("Adding actor filter")
			requestBuilder = requestBuilder.Person(options.Actor)
		} else if options.Director != "" {
			log.Debug().Str("director", options.Director).Msg("Adding director filter")
			requestBuilder = requestBuilder.Person(options.Director)
		} else if options.Creator != "" {
			log.Debug().Str("creator", options.Creator).Msg("Adding creator filter")
			requestBuilder = requestBuilder.Person(options.Creator)
		}

		// Content rating filter
		if options.ContentRating != "" {
			log.Debug().Str("contentRating", options.ContentRating).Msg("Adding content rating filter")
			requestBuilder = requestBuilder.OfficialRatings([]string{options.ContentRating})
		}

		// Tags filter
		if len(options.Tags) > 0 {
			log.Debug().Strs("tags", options.Tags).Msg("Adding tags filter")
			requestBuilder = requestBuilder.Tags(options.Tags)
		}

		// Watched filter
		if options.Watched {
			log.Debug().Msg("Adding watched filter")
			requestBuilder = requestBuilder.IsPlayed(true)
		}

		// Rating filter
		if options.MinimumRating > 0 {
			minRating := float64(options.MinimumRating)
			log.Debug().Float64("minRating", minRating).Msg("Adding minimum rating filter")
			requestBuilder = requestBuilder.MinCommunityRating(minRating)
		}
	}

	// Always enable these options
	requestBuilder = requestBuilder.
		EnableUserData(true).
		EnableImages(true).
		EnableTotalRecordCount(true)

	// Execute the request
	log.Debug().Msg("Executing Jellyfin API request")
	result, resp, err := requestBuilder.Execute()

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

// getMovieByIDWithFilter is a helper function that gets a movie by ID using the GetItems endpoint with a filter
func (j *JellyfinClient) getMovieByIDWithFilter(ctx context.Context, id string) (*models.MediaItem[*t.Movie], error) {
	log := logger.LoggerFromContext(ctx)

	log.Debug().
		Str("movieID", id).
		Msg("Using GetItems with ID filter")

	// Prepare parameters for GetItems
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}
	ids := []string{id}
	fields := defaultExtraFields

	// Build the GetItems request - important to set limit to 1 and specific filters
	requestBuilder := j.client.ItemsAPI.GetItems(ctx).
		Ids(ids).
		IncludeItemTypes(includeItemTypes).
		Fields(fields).
		EnableUserData(true).
		EnableImages(true).
		// Explicitly set to get only ONE item and force an exact match
		Limit(1).
		EnableTotalRecordCount(true)

	// Add user ID if available
	if j.getUserID() != "" {
		log.Debug().Str("userID", j.getUserID()).Msg("Adding user ID to request")
		requestBuilder = requestBuilder.UserId(j.getUserID())
	}

	// Execute the request
	log.Debug().Msg("Executing Jellyfin API GetItems request with ID filter")
	result, resp, err := requestBuilder.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.jellyfinConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("movieID", id).
			Int("statusCode", resp.StatusCode).
			Msg("Failed to fetch movie from Jellyfin with filter method")
		return nil, fmt.Errorf("failed to fetch movie: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("movieID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No movie found with the specified ID")
		return nil, fmt.Errorf("movie with ID %s not found", id)
	}

	// If we got multiple results, that's an error for GetMovieByID
	if len(result.Items) > 1 {
		log.Error().
			Str("movieID", id).
			Int("resultCount", len(result.Items)).
			Msg("Multiple movies found with the specified ID")
		return nil, fmt.Errorf("multiple movies found")
	}

	item := result.Items[0]

	// Double-check that the returned item is a movie
	if *item.Type != jellyfin.BASEITEMKIND_MOVIE {
		log.Error().
			Str("movieID", id).
			Str("actualType", string(*item.Type)).
			Msg("Item with specified ID is not a movie")
		return nil, fmt.Errorf("item with ID %s is not a movie", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("movieID", id).
		Str("movieName", *item.Name.Get()).
		Msg("Successfully retrieved movie from Jellyfin using filter method")

	itemMovie, err := GetItem[*t.Movie](ctx, j, &item)
	if err != nil {
		log.Error().
			Err(err).
			Str("movieID", id).
			Msg("Error converting Jellyfin item to movie format")
		return nil, fmt.Errorf("error converting movie data: %w", err)
	}

	movie, err := GetMediaItem[*t.Movie](ctx, j, itemMovie, *item.Id)
	if err != nil {
		log.Error().
			Err(err).
			Str("movieID", id).
			Str("movieName", *item.Name.Get()).
			Msg("Error creating media item")
		return nil, fmt.Errorf("error creating media item: %w", err)
	}

	return movie, nil
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

	// Prepare all parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}
	// Use a single ID rather than splitting which might mess up the query
	ids := []string{id}
	fields := defaultExtraFields

	// Log parameters
	log.Debug().
		Str("movieID", id).
		Interface("includeItemTypes", includeItemTypes).
		Msg("Building Jellyfin API request for movie by ID")

	// Build the entire request in a single chain
	requestBuilder := j.client.ItemsAPI.GetItems(ctx).
		Ids(ids).
		IncludeItemTypes(includeItemTypes).
		Fields(fields).
		EnableUserData(true).
		EnableImages(true).
		Limit(1).
		IsMovie(true).
		EnableTotalRecordCount(true)

	// Add user ID if available
	if j.getUserID() != "" {
		log.Debug().Str("userID", j.getUserID()).Msg("Adding user ID to request")
		requestBuilder = requestBuilder.UserId(j.getUserID())
	}

	// Execute the request
	log.Debug().Msg("Executing Jellyfin API request for movie by ID")
	result, resp, err := requestBuilder.Execute()
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

	// Check if we have the expected number of items
	log.Debug().
		Int("resultCount", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Jellyfin API returned items")

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("movieID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No movie found with the specified ID")
		return &models.MediaItem[*t.Movie]{}, fmt.Errorf("movie with ID %s not found", id)
	}

	// Explicitly check for multiple items matching the ID
	if *result.TotalRecordCount > 1 {
		log.Error().
			Str("movieID", id).
			Int("totalItems", int(*result.TotalRecordCount)).
			Msg("Multiple movies found with the specified ID")
		return &models.MediaItem[*t.Movie]{}, fmt.Errorf("multiple movies found")
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

// SearchMovies searches for movies using the provided query options
func (j *JellyfinClient) SearchMovies(ctx context.Context, options *t.QueryOptions) ([]*models.MediaItem[*t.Movie], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	// Ensure we have a search query
	if options == nil || options.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.jellyfinConfig().GetBaseURL()).
		Str("query", options.Query).
		Msg("Searching movies in Jellyfin server")

	// Set up search parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}
	mediaTypes := []jellyfin.MediaType{jellyfin.MEDIATYPE_VIDEO}
	fields := defaultExtraFields

	// Default values
	limit := int32(50) // Higher limit for search results
	startIndex := int32(0)
	recursive := true
	isMovie := true

	// Apply options from query parameters
	if options.Limit > 0 {
		limit = int32(options.Limit)
	}

	if options.Offset > 0 {
		startIndex = int32(options.Offset)
	}

	// Log search parameters
	log.Debug().
		Str("query", options.Query).
		Int32("limit", limit).
		Int32("startIndex", startIndex).
		Msg("Building Jellyfin search request")

	// Create the search request
	requestBuilder := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		IsMovie(isMovie).
		Recursive(recursive).
		SearchTerm(options.Query). // This is the key parameter for search
		Fields(fields).
		MediaTypes(mediaTypes).
		Limit(limit).
		StartIndex(startIndex)

	// Add userId if available
	if j.getUserID() != "" {
		requestBuilder = requestBuilder.UserId(j.getUserID())
	}

	// Additional filtering
	if options.Genre != "" {
		requestBuilder = requestBuilder.Genres([]string{options.Genre})
	}

	if options.Year > 0 {
		requestBuilder = requestBuilder.Years([]int32{int32(options.Year)})
	}

	if options.MinimumRating > 0 {
		minRating := float64(options.MinimumRating)
		requestBuilder = requestBuilder.MinCommunityRating(minRating)
	}

	// Execute the search
	log.Debug().Msg("Executing Jellyfin movie search request")
	result, resp, err := requestBuilder.
		EnableUserData(true).
		EnableImages(true).
		EnableTotalRecordCount(true).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.jellyfinConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("query", options.Query).
			Int("statusCode", resp.StatusCode).
			Msg("Failed to search movies from Jellyfin")
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully searched movies from Jellyfin")

	// Convert results to expected format
	movies := make([]*models.MediaItem[*t.Movie], 0)

	for _, item := range result.Items {
		if item.Type == nil {
			continue
		}

		log.Debug().
			Str("itemType", string(*item.Type)).
			Str("itemName", getItemName(&item)).
			Msg("Processing search result")

		if *item.Type == jellyfin.BASEITEMKIND_MOVIE {
			itemMovie, err := GetItem[*t.Movie](ctx, j, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("movieID", *item.Id).
					Str("movieName", getItemName(&item)).
					Msg("Error converting Jellyfin item to movie format")
				continue
			}

			movie, err := GetMediaItem[*t.Movie](ctx, j, itemMovie, *item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("movieID", *item.Id).
					Str("movieName", getItemName(&item)).
					Msg("Error creating media item")
				continue
			}

			movies = append(movies, movie)
		}
	}

	log.Info().
		Int("moviesReturned", len(movies)).
		Msg("Completed SearchMovies request")

	return movies, nil
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
