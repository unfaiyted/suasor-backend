// movies.go
package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/utils"
)

// GetMovies retrieves movies from the Emby server
func (e *EmbyClient) GetMovies(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Movie], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving movies from Emby server")

	// Create query parameters
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Movie"),
		Recursive:        optional.NewBool(true),
	}

	// Apply options
	applyQueryOptions(&queryParams, options)

	// Call the Emby API
	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Msg("Failed to fetch movies from Emby")
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved movies from Emby")

	// Convert results to expected format
	movies := make([]types.MediaItem[types.Movie], 0)
	for _, item := range items.Items {
		if item.Type_ == "Movie" {
			movie, err := e.convertToMovie(ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("movieID", item.Id).
					Str("movieName", item.Name).
					Msg("Error converting Emby item to movie format")
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

// GetMovieByID retrieves a specific movie by ID
func (e *EmbyClient) GetMovieByID(ctx context.Context, id string) (types.MediaItem[types.Movie], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("movieID", id).
		Msg("Retrieving specific movie from Emby server")

	// Create query parameters
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids:              optional.NewString(id),
		IncludeItemTypes: optional.NewString("Movie"),
		Fields:           optional.NewString("ProductionYear,PremiereDate,ChannelMappingInfo,DateCreated,Genres,IndexOptions,HomePageUrl,Overview,ParentId,Path,ProviderIds,Studios,SortName"),
	}

	// Call the Emby API
	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("movieID", id).
			Msg("Failed to fetch movie from Emby")
		return types.MediaItem[types.Movie]{}, fmt.Errorf("failed to fetch movie: %w", err)
	}

	// Check if any items were returned
	if len(items.Items) == 0 {
		log.Error().
			Str("movieID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No movie found with the specified ID")
		return types.MediaItem[types.Movie]{}, fmt.Errorf("movie with ID %s not found", id)
	}

	item := items.Items[0]

	// Double-check that the returned item is a movie
	if item.Type_ != "Movie" {
		log.Error().
			Str("movieID", id).
			Str("actualType", item.Type_).
			Msg("Item with specified ID is not a movie")
		return types.MediaItem[types.Movie]{}, fmt.Errorf("item with ID %s is not a movie", id)
	}

	movie, err := e.convertToMovie(ctx, &item)
	if err != nil {
		log.Error().
			Err(err).
			Str("movieID", id).
			Str("movieName", item.Name).
			Msg("Error converting Emby item to movie format")
		return types.MediaItem[types.Movie]{}, fmt.Errorf("error converting movie data: %w", err)
	}

	return movie, nil
}

// GetMovieGenres retrieves movie genres from the Emby server
func (e *EmbyClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	opts := embyclient.GenresServiceApiGetGenresOpts{IsMovie: optional.NewBool(true)}

	result, _, err := e.client.GenresServiceApi.GetGenres(ctx, &opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie genres: %w", err)
	}

	genres := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		genres = append(genres, item.Name)
	}

	return genres, nil
}
