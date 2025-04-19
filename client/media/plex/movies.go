package plex

import (
	"context"
	"fmt"
	"strconv"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"

	"github.com/LukeHagar/plexgo/models/operations"
)

// GetMovies retrieves movies from Plex
func (c *PlexClient) GetMovies(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Movie], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving movies from Plex server")

	// First, find the movie library section
	log.Debug().Msg("Finding movie library section")
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find movie library section")
		return nil, err
	}

	if movieSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No movie library section found in Plex")
		return nil, nil
	}

	// Get movies from the movie section
	sectionKey, _ := strconv.Atoi(movieSectionKey)
	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for movies")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		Tag:         "all",
		Type:        operations.GetLibraryItemsQueryParamTypeMovie,
		SectionKey:  sectionKey,
		IncludeMeta: operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
	})

	// TODO: this interface here is wrong, need differernt type
	// log.Debug().Interface("response", res).Msg("Response from Plex")

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get movies from Plex")
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No movies found in Plex")
		return nil, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved movies from Plex")

	movies, err := GetMediaItemList[*types.Movie](ctx, c, res.Object.MediaContainer.Metadata)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to get movies from Plex")
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("moviesReturned", len(movies)).
		Msg("Completed GetMovies request")

	return movies, nil
}

// GetMovieByID retrieves a specific movie by ID
func (c *PlexClient) GetMovieByID(ctx context.Context, id string) (*models.MediaItem[*types.Movie], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("movieID", id).
		Msg("Retrieving specific movie from Plex server")

	ratingKey, _ := strconv.Atoi(id)
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("movieID", id).
			Msg("Failed to get movie from Plex")
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("movieID", id).
			Msg("Movie not found in Plex")
		return nil, fmt.Errorf("movie not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "movie" {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("movieID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a movie")
		return nil, fmt.Errorf("item is not a movie")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("movieID", id).
		Str("movieTitle", item.Title).
		Msg("Successfully retrieved movie from Plex")

	itemMovie, err := GetItemFromMetadata[*types.Movie](ctx, c, &item)
	movie, err := GetMediaItem[*types.Movie](ctx, c, itemMovie, item.RatingKey)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("movieID", id).
		Str("movieTitle", movie.Data.Details.Title).
		Msg("Successfully converted movie data")

	return movie, nil
}

// GetMovieGenres retrieves movie genres from Plex
func (c *PlexClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving movie genres from Plex server")

	// Find the movie library section
	log.Debug().Msg("Finding movie library section")
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find movie library section")
		return nil, err
	}

	if movieSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
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
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("genresFound", len(genres)).
		Msg("Successfully retrieved movie genres from Plex")

	return genres, nil
}
