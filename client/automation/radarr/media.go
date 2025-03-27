package radarr

import (
	"context"
	"fmt"

	radarr "github.com/devopsarr/radarr-go/radarr"
	"strconv"
	"suasor/client/automation/types"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/utils"
)

func (r *RadarrClient) GetLibraryItems(ctx context.Context, options *types.LibraryQueryOptions) ([]models.AutomationMediaItem[types.AutomationData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("baseURL", r.config.BaseURL).
		Msg("Retrieving library items from Radarr server")

	// Call the Radarr API
	log.Debug().Msg("Making API request to Radarr server for movie library")

	moviesResult, resp, err := r.client.MovieAPI.ListMovie(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", r.config.BaseURL).
			Str("apiEndpoint", "/movie").
			Int("statusCode", 0).
			Msg("Failed to fetch movies from Radarr")
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("movieCount", len(moviesResult)).
		Msg("Successfully retrieved movies from Radarr")

	// Apply paging if options provided
	var start, end int
	if options != nil {
		if options.Offset > 0 {
			start = options.Offset
		}
		if options.Limit > 0 {
			end = start + options.Limit
			if end > len(moviesResult) {
				end = len(moviesResult)
			}
		} else {
			end = len(moviesResult)
		}
	} else {
		end = len(moviesResult)
	}

	// Ensure valid slice bounds
	if start >= len(moviesResult) {
		start = 0
		end = 0
	}

	// Apply paging
	var pagedMovies []radarr.MovieResource
	if start < end {
		pagedMovies = moviesResult[start:end]
	} else {
		pagedMovies = []radarr.MovieResource{}
	}

	// Convert to our internal type
	mediaItems := make([]models.AutomationMediaItem[types.AutomationData], 0, len(pagedMovies))
	for _, movie := range pagedMovies {
		mediaItem := r.convertMovieToMediaItem(&movie)
		mediaItems = append(mediaItems, mediaItem)
	}

	log.Info().
		Int("itemsReturned", len(mediaItems)).
		Msg("Completed GetLibraryItems request")

	return mediaItems, nil
}

func (r *RadarrClient) GetMediaByID(ctx context.Context, id int64) (models.AutomationMediaItem[types.AutomationData], error) {

	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Int64("movieID", id).
		Str("baseURL", r.config.BaseURL).
		Msg("Retrieving specific movie from Radarr server")

	// Call the Radarr API
	log.Debug().
		Int64("movieID", id).
		Msg("Making API request to Radarr server")

	movie, resp, err := r.client.MovieAPI.GetMovieById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", r.config.BaseURL).
			Str("apiEndpoint", fmt.Sprintf("/movie/%d", id)).
			Int("statusCode", 0).
			Msg("Failed to fetch movie from Radarr")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to fetch movie: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int64("movieID", id).
		Str("movieTitle", movie.GetTitle()).
		Msg("Successfully retrieved movie from Radarr")

	// Convert to our internal type
	mediaItem := r.convertMovieToMediaItem(movie)

	log.Debug().
		Int64("movieID", id).
		Str("movieTitle", mediaItem.Title).
		Int32("year", mediaItem.Year).
		Msg("Successfully returned movie data")

	return mediaItem, nil
}

func (r *RadarrClient) AddMedia(ctx context.Context, req requests.AutomationMediaAddRequest) (models.AutomationMediaItem[types.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("title", req.Title).
		Msg("Adding movie to Radarr")

	// Create new movie resource
	newMovie := radarr.NewMovieResource()
	newMovie.SetTitle(req.Title)
	newMovie.SetQualityProfileId(int32(req.QualityProfileID))
	newMovie.SetTmdbId(int32(req.TMDBID))
	newMovie.SetYear(int32(req.Year))
	newMovie.SetMonitored(req.Monitored)
	newMovie.SetRootFolderPath(req.Path)
	newMovie.SetTags(req.Tags)

	// Set minimum availability if provided
	// if req.MinimumAvailability != "" {
	// 	newMovie.SetMinimumAvailability(req.MinimumAvailability)
	// }

	// Make API request
	result, resp, err := r.client.MovieAPI.CreateMovie(ctx).MovieResource(*newMovie).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", r.config.BaseURL).
			Str("title", req.Title).
			Msg("Failed to add movie to Radarr")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to add movie: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("movieID", result.GetId()).
		Str("title", result.GetTitle()).
		Msg("Successfully added movie to Radarr")

	return r.convertMovieToMediaItem(result), nil
}

func (r *RadarrClient) UpdateMedia(ctx context.Context, id int64, item requests.AutomationMediaUpdateRequest) (models.AutomationMediaItem[types.AutomationData], error) {

	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Int64("movieID", id).
		Msg("Updating movie in Radarr")

	// First get the existing movie
	existingMovie, resp, err := r.client.MovieAPI.GetMovieById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("movieID", id).
			Msg("Failed to fetch movie for update")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to fetch movie for update: %w", err)
	}

	// Update fields as needed
	existingMovie.SetMonitored(item.Monitored)

	if item.QualityProfileID > 0 {
		existingMovie.SetQualityProfileId(int32(item.QualityProfileID))
	}

	if item.Path != "" {
		existingMovie.SetPath(item.Path)
	}

	if item.Tags != nil {
		existingMovie.SetTags(convertInt64SliceToInt32(item.Tags))
	}

	stringId := strconv.FormatInt(id, 10)

	// Send update request
	updatedMovie, resp, err := r.client.MovieAPI.UpdateMovie(ctx, stringId).MovieResource(*existingMovie).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("movieID", id).
			Msg("Failed to update movie in Radarr")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to update movie: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("movieID", updatedMovie.GetId()).
		Msg("Successfully updated movie in Radarr")

	return r.convertMovieToMediaItem(updatedMovie), nil
}

func (r *RadarrClient) DeleteMedia(ctx context.Context, id int64) error {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Int64("movieID", id).
		Msg("Deleting movie from Radarr")

	// Optional deletion flags
	deleteFiles := false
	addExclusion := false

	resp, err := r.client.MovieAPI.DeleteMovie(ctx, int32(id)).
		DeleteFiles(deleteFiles).
		AddImportExclusion(addExclusion).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Int64("movieID", id).
			Msg("Failed to delete movie from Radarr")
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int64("movieID", id).
		Msg("Successfully deleted movie from Radarr")

	return nil
}

func (r *RadarrClient) SearchMedia(ctx context.Context, query string, options *types.SearchOptions) ([]models.AutomationMediaItem[types.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("query", query).
		Msg("Searching for movies in Radarr")

	searchResult, resp, err := r.client.MovieLookupAPI.ListMovieLookup(ctx).Term(query).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("query", query).
			Msg("Failed to search for movies in Radarr")
		return nil, fmt.Errorf("failed to search for movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("resultCount", len(searchResult)).
		Msg("Successfully searched for movies in Radarr")

	// Convert results to MediaItems
	mediaItems := make([]models.AutomationMediaItem[types.AutomationData], 0, len(searchResult))
	for _, movie := range searchResult {
		mediaItem := r.convertMovieToMediaItem(&movie)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}
