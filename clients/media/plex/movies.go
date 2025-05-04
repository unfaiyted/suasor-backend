package plex

import (
	"context"
	"fmt"
	"strconv"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"github.com/LukeHagar/plexgo/models/operations"
)

// GetMovies retrieves movies from Plex
func (c *PlexClient) GetMovies(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Movie], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving movies from Plex server")

	// First, find the movie library section
	log.Debug().Msg("Finding movie library section")
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to find movie library section")
		return nil, err
	}

	if movieSectionKey == "" {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No movie library section found in Plex")
		return nil, nil
	}

	// Get movies from the movie section
	sectionKey, _ := strconv.Atoi(movieSectionKey)
	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for movies")

	// Handle pagination when fetching all movies (limit=0, offset=0)
	if options.Limit == 0 && options.Offset == 0 {
		log.Debug().Msg("Fetching all movies, NO LIMITS!")
		return c.getAllMovies(ctx, sectionKey)
	}

	// Regular case with specified limit and offset
	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		Tag:                 "all",
		Type:                operations.GetLibraryItemsQueryParamTypeMovie,
		SectionKey:          sectionKey,
		XPlexContainerStart: &options.Offset,
		XPlexContainerSize:  &options.Limit,
		IncludeGuids:        operations.IncludeGuidsEnable.ToPointer(),
		IncludeMeta:         operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
	})

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Int("sectionKey", sectionKey).
			Msg("Failed to get movies from Plex")
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No movies found in Plex")
		return nil, nil
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved movies from Plex")

	movies, err := GetMediaItemList[*types.Movie](ctx, c, res.Object.MediaContainer.Metadata)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to get movies from Plex")
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("moviesReturned", len(movies)).
		Msg("Completed GetMovies request")

	return movies, nil
}

// GetMovieByID retrieves a specific movie by ID
func (c *PlexClient) GetMovieByID(ctx context.Context, id string) (*models.MediaItem[*types.Movie], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("movieID", id).
		Msg("Retrieving specific movie from Plex server")

	ratingKey, err := strconv.Atoi(id)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("movieID", id).
			Msg("Invalid movie ID format")
		return nil, fmt.Errorf("invalid movie ID format: %w", err)
	}
	int64RatingKey := int64(ratingKey)

	log.Debug().
		Str("movieID", id).
		Int64("ratingKey", int64RatingKey).
		Msg("Making API request to Plex server for movie")

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("movieID", id).
			Msg("Failed to get movie from Plex")
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("movieID", id).
			Msg("Movie not found in Plex")
		return nil, fmt.Errorf("movie not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "movie" {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("movieID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a movie")
		return nil, fmt.Errorf("item is not a movie")
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("movieID", id).
		Str("movieTitle", item.Title).
		Msg("Successfully retrieved movie from Plex")

	// Convert Plex metadata to Movie object
	itemMovie, err := GetItemFromMetadata[*types.Movie](ctx, c, &item)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("movieID", id).
			Msg("Failed to convert movie metadata")
		return nil, fmt.Errorf("failed to convert movie metadata: %w", err)
	}

	// Verify that the movie has a Details field
	if itemMovie.Details == nil {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("movieID", id).
			Msg("Movie metadata is missing Details field")
			
		// Create a minimal Details field if it's nil
		itemMovie.Details = &types.MediaDetails{
			Title: item.Title,
		}
	}

	// Create MediaItem from Movie object
	movie, err := GetMediaItem[*types.Movie](ctx, c, itemMovie, item.RatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("movieID", id).
			Msg("Failed to create media item")
		return nil, fmt.Errorf("failed to create media item: %w", err)
	}

	// Verify the movie's data before logging
	if movie.Data == nil || movie.Data.Details == nil {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("movieID", id).
			Msg("Movie data or details is nil")
		return nil, fmt.Errorf("movie data conversion failed: nil data or details")
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("movieID", id).
		Str("movieTitle", movie.Data.Details.Title).
		Msg("Successfully converted movie data")

	return movie, nil
}

// GetMovieGenres retrieves movie genres from Plex
func (c *PlexClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving movie genres from Plex server")

	// Find the movie library section
	log.Debug().Msg("Finding movie library section")
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to find movie library section")
		return nil, err
	}

	if movieSectionKey == "" {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No movie library section found in Plex")
		return []string{}, nil
	}

	// Get genres directly from the genre endpoint
	sectionKey, _ := strconv.Atoi(movieSectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for movie genres")

	res, err := c.plexAPI.Library.GetGenresLibrary(ctx, sectionKey, operations.GetGenresLibraryQueryParamTypeMovie)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Int("sectionKey", sectionKey).
			Msg("Failed to get movie genres from Plex")
		return nil, fmt.Errorf("failed to get movie genres: %w", err)
	}

	genreMap := make(map[string]bool)
	if res.Object.MediaContainer != nil {
		directories := res.Object.MediaContainer.GetDirectory()
		log.Debug().
			Int("genreCount", len(directories)).
			Msg("Extracting genres from directories")

		for _, item := range directories {
			genreMap[item.Title] = true
			log.Debug().
				Str("genre", item.Title).
				Msg("Found movie genre")
		}
	}

	genres := make([]string, 0, len(genreMap))
	for genre := range genreMap {
		genres = append(genres, genre)
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("genresFound", len(genres)).
		Msg("Successfully retrieved movie genres from Plex")

	return genres, nil
}

// getAllMovies fetches all movies from Plex using pagination with batch size of 50
func (c *PlexClient) getAllMovies(ctx context.Context, sectionKey int) ([]*models.MediaItem[*types.Movie], error) {
	log := logger.LoggerFromContext(ctx)
	var allMovies []*models.MediaItem[*types.Movie]

	batchSize := 50
	offset := 0

	for {
		log.Debug().
			Int("sectionKey", sectionKey).
			Int("offset", offset).
			Int("batchSize", batchSize).
			Msg("Fetching batch of movies from Plex")

		res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
			Tag:                 "all",
			Type:                operations.GetLibraryItemsQueryParamTypeMovie,
			SectionKey:          sectionKey,
			XPlexContainerStart: &offset,
			XPlexContainerSize:  &batchSize,
			IncludeGuids:        operations.IncludeGuidsEnable.ToPointer(),
			IncludeMeta:         operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
		})

		if err != nil {
			log.Error().
				Err(err).
				Uint64("clientID", c.GetClientID()).
				Str("clientType", string(c.GetClientType())).
				Int("sectionKey", sectionKey).
				Int("offset", offset).
				Msg("Failed to get batch of movies from Plex")
			return nil, fmt.Errorf("failed to get movies batch: %w", err)
		}

		// Check if we have results
		if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
			// No more results, break the loop
			break
		}

		// Process the current batch with error recovery for individual items
		var batchMovies []*models.MediaItem[*types.Movie]
		for i, item := range res.Object.MediaContainer.Metadata {
			// Skip potential problematic items
			if item.RatingKey == "" {
				log.Warn().
					Int("batchIndex", i).
					Int("offset", offset).
					Msg("Skipping movie with missing RatingKey")
				continue
			}

			// Try to process each movie individually
			movieID := item.RatingKey
			movieTitle := item.Title

			try := func() (movie *models.MediaItem[*types.Movie], err error) {
				// Recover from panics during processing
				defer func() {
					if r := recover(); r != nil {
						log.Error().
							Interface("panic", r).
							Str("movieID", movieID).
							Str("movieTitle", movieTitle).
							Int("offset", offset+i).
							Msg("Panic while processing movie, skipping")
						err = fmt.Errorf("panic while processing movie: %v", r)
					}
				}()

				itemT, err := GetItemFromLibraryMetadata[*types.Movie](ctx, c, &item)
				if err != nil {
					return nil, err
				}

				movie, err = GetMediaItem[*types.Movie](ctx, c, itemT, item.RatingKey)
				return movie, err
			}

			movie, err := try()
			if err != nil {
				log.Error().
					Err(err).
					Str("movieID", movieID).
					Str("movieTitle", movieTitle).
					Int("offset", offset+i).
					Msg("Error processing movie, skipping")
				continue
			}

			if movie != nil {
				batchMovies = append(batchMovies, movie)
			}
		}

		// Add to our accumulated results
		allMovies = append(allMovies, batchMovies...)

		batchCount := len(res.Object.MediaContainer.Metadata)
		log.Debug().
			Int("batchCount", batchCount).
			Int("processedCount", len(batchMovies)).
			Int("totalSoFar", len(allMovies)).
			Msg("Retrieved batch of movies")

		// If we got fewer items than requested, we've reached the end
		if batchCount < batchSize {
			break
		}

		// Move to the next batch
		offset += batchSize
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("totalMovies", len(allMovies)).
		Msg("Successfully retrieved all movies from Plex")

	return allMovies, nil
}
