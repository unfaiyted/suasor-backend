package services

import (
	"context"
	"fmt"
	"suasor/clients/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"
)

// CoreSeriesService defines the interface for TV series-related operations
type CoreSeriesService interface {
	// Episode-related operations
	GetEpisodesBySeasonID(ctx context.Context, seasonID uint64) ([]*models.MediaItem[*types.Episode], error)
	GetEpisodesBySeriesID(ctx context.Context, seriesID uint64) ([]*models.MediaItem[*types.Episode], error)
	GetRecentlyAddedEpisodes(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Episode], error)
	GetUnwatchedEpisodes(ctx context.Context, userID uint64, seriesID uint64) ([]*models.MediaItem[*types.Episode], error)
	GetNextEpisodeToWatch(ctx context.Context, userID uint64, seriesID uint64) (*models.MediaItem[*types.Episode], error)
	GetEpisodeByNumber(ctx context.Context, seriesID uint64, seasonNumber int, episodeNumber int) (*models.MediaItem[*types.Episode], error)
	GetEpisodesByAttribute(ctx context.Context, attribute string, value interface{}, limit int) ([]*models.MediaItem[*types.Episode], error)

	// Season-related operations
	GetSeasonsBySeriesID(ctx context.Context, seriesID uint64) ([]*models.MediaItem[*types.Season], error)
	GetSeasonWithEpisodes(ctx context.Context, seasonID uint64) (*models.MediaItem[*types.Season], []*models.MediaItem[*types.Episode], error)
	GetSeasonByNumber(ctx context.Context, seriesID uint64, seasonNumber int) (*models.MediaItem[*types.Season], error)
	GetRecentlyAddedSeasons(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Season], error)

	// Series-related operations
	GetSeriesWithSeasons(ctx context.Context, seriesID uint64) (*models.MediaItem[*types.Series], []*models.MediaItem[*types.Season], error)
	GetRecentlyAiredSeries(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Series], error)
	GetPopularSeries(ctx context.Context, limit int) ([]*models.MediaItem[*types.Series], error)
	GetInProgressSeries(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*types.Series], error)
	GetSeriesByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Series], error)
	GetSimilarSeries(ctx context.Context, seriesID uint64, limit int) ([]*models.MediaItem[*types.Series], error)
	GetSeriesInCollection(ctx context.Context, collectionID uint64) ([]*models.MediaItem[*types.Series], error)

	// Calendar operations
	GetUpcomingEpisodes(ctx context.Context, days int) ([]*models.MediaItem[*types.Episode], error)
	GetEpisodesAiredBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*models.MediaItem[*types.Episode], error)

	// Search operations
	SearchSeriesLibrary(ctx context.Context, query types.QueryOptions) (*models.MediaItemList, error)
}

// coreSeriesService implements the CoreSeriesService interface
type coreSeriesService struct {
	seriesRepo repository.SeriesRepository
}

// NewCoreSeriesService creates a new core series service
func NewCoreSeriesService(seriesRepo repository.SeriesRepository) CoreSeriesService {
	return &coreSeriesService{
		seriesRepo: seriesRepo,
	}
}

// GetEpisodesBySeasonID gets all episodes for a specific season
func (s *coreSeriesService) GetEpisodesBySeasonID(ctx context.Context, seasonID uint64) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seasonID", seasonID).
		Msg("Getting episodes by season ID")

	episodes, err := s.seriesRepo.GetEpisodesBySeasonID(ctx, seasonID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seasonID", seasonID).
			Msg("Failed to get episodes for season")
		return nil, fmt.Errorf("failed to get episodes: %w", err)
	}

	return episodes, nil
}

// GetEpisodesBySeriesID gets all episodes for a specific series
func (s *coreSeriesService) GetEpisodesBySeriesID(ctx context.Context, seriesID uint64) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Msg("Getting episodes by series ID")

	episodes, err := s.seriesRepo.GetEpisodesBySeriesID(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to get episodes for series")
		return nil, fmt.Errorf("failed to get episodes: %w", err)
	}

	return episodes, nil
}

// GetRecentlyAddedEpisodes gets recently added episodes
func (s *coreSeriesService) GetRecentlyAddedEpisodes(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recently added episodes")

	episodes, err := s.seriesRepo.GetRecentlyAddedEpisodes(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Int("limit", limit).
			Msg("Failed to get recently added episodes")
		return nil, fmt.Errorf("failed to get recently added episodes: %w", err)
	}

	return episodes, nil
}

// GetUnwatchedEpisodes gets unwatched episodes for a user
func (s *coreSeriesService) GetUnwatchedEpisodes(ctx context.Context, userID uint64, seriesID uint64) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("seriesID", seriesID).
		Msg("Getting unwatched episodes")

	episodes, err := s.seriesRepo.GetUnwatchedEpisodes(ctx, userID, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("seriesID", seriesID).
			Msg("Failed to get unwatched episodes")
		return nil, fmt.Errorf("failed to get unwatched episodes: %w", err)
	}

	return episodes, nil
}

// GetNextEpisodeToWatch gets the next episode to watch for a user
func (s *coreSeriesService) GetNextEpisodeToWatch(ctx context.Context, userID uint64, seriesID uint64) (*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("seriesID", seriesID).
		Msg("Getting next episode to watch")

	episode, err := s.seriesRepo.GetNextEpisodeToWatch(ctx, userID, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("seriesID", seriesID).
			Msg("Failed to get next episode to watch")
		return nil, fmt.Errorf("failed to get next episode to watch: %w", err)
	}

	return episode, nil
}

// GetEpisodeByNumber gets an episode by its season and episode number
func (s *coreSeriesService) GetEpisodeByNumber(ctx context.Context, seriesID uint64, seasonNumber int, episodeNumber int) (*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeNumber", episodeNumber).
		Msg("Getting episode by number")

	episode, err := s.seriesRepo.GetEpisodeByNumber(ctx, seriesID, seasonNumber, episodeNumber)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Int("seasonNumber", seasonNumber).
			Int("episodeNumber", episodeNumber).
			Msg("Failed to get episode by number")
		return nil, fmt.Errorf("failed to get episode by number: %w", err)
	}

	return episode, nil
}

// GetEpisodesByAttribute gets episodes by a specific attribute
func (s *coreSeriesService) GetEpisodesByAttribute(ctx context.Context, attribute string, value interface{}, limit int) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("attribute", attribute).
		Interface("value", value).
		Int("limit", limit).
		Msg("Getting episodes by attribute")

	episodes, err := s.seriesRepo.GetEpisodesByAttribute(ctx, attribute, value, limit)
	if err != nil {
		log.Error().Err(err).
			Str("attribute", attribute).
			Interface("value", value).
			Int("limit", limit).
			Msg("Failed to get episodes by attribute")
		return nil, fmt.Errorf("failed to get episodes by attribute: %w", err)
	}

	return episodes, nil
}

// GetSeasonsBySeriesID gets all seasons for a specific series
func (s *coreSeriesService) GetSeasonsBySeriesID(ctx context.Context, seriesID uint64) ([]*models.MediaItem[*types.Season], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Msg("Getting seasons by series ID")

	seasons, err := s.seriesRepo.GetSeasonsBySeriesID(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to get seasons for series")
		return nil, fmt.Errorf("failed to get seasons: %w", err)
	}

	return seasons, nil
}

// GetSeasonWithEpisodes gets a season and all its episodes
func (s *coreSeriesService) GetSeasonWithEpisodes(ctx context.Context, seasonID uint64) (*models.MediaItem[*types.Season], []*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seasonID", seasonID).
		Msg("Getting season with episodes")

	season, episodes, err := s.seriesRepo.GetSeasonWithEpisodes(ctx, seasonID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seasonID", seasonID).
			Msg("Failed to get season with episodes")
		return nil, nil, fmt.Errorf("failed to get season with episodes: %w", err)
	}

	return season, episodes, nil
}

// GetSeasonByNumber gets a season by its number within a series
func (s *coreSeriesService) GetSeasonByNumber(ctx context.Context, seriesID uint64, seasonNumber int) (*models.MediaItem[*types.Season], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Msg("Getting season by number")

	season, err := s.seriesRepo.GetSeasonByNumber(ctx, seriesID, seasonNumber)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Int("seasonNumber", seasonNumber).
			Msg("Failed to get season by number")
		return nil, fmt.Errorf("failed to get season by number: %w", err)
	}

	return season, nil
}

// GetRecentlyAddedSeasons gets recently added seasons
func (s *coreSeriesService) GetRecentlyAddedSeasons(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Season], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recently added seasons")

	seasons, err := s.seriesRepo.GetRecentlyAddedSeasons(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Int("limit", limit).
			Msg("Failed to get recently added seasons")
		return nil, fmt.Errorf("failed to get recently added seasons: %w", err)
	}

	return seasons, nil
}

// GetSeriesWithSeasons gets a series and all its seasons
func (s *coreSeriesService) GetSeriesWithSeasons(ctx context.Context, seriesID uint64) (*models.MediaItem[*types.Series], []*models.MediaItem[*types.Season], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Msg("Getting series with seasons")

	series, seasons, err := s.seriesRepo.GetSeriesWithSeasons(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to get series with seasons")
		return nil, nil, fmt.Errorf("failed to get series with seasons: %w", err)
	}

	return series, seasons, nil
}

// GetRecentlyAiredSeries gets recently aired series
func (s *coreSeriesService) GetRecentlyAiredSeries(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recently aired series")

	series, err := s.seriesRepo.GetRecentlyAiredSeries(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Int("limit", limit).
			Msg("Failed to get recently aired series")
		return nil, fmt.Errorf("failed to get recently aired series: %w", err)
	}

	return series, nil
}

// GetPopularSeries gets popular series
func (s *coreSeriesService) GetPopularSeries(ctx context.Context, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("Getting popular series")

	series, err := s.seriesRepo.GetPopularSeries(ctx, limit)
	if err != nil {
		log.Error().Err(err).
			Int("limit", limit).
			Msg("Failed to get popular series")
		return nil, fmt.Errorf("failed to get popular series: %w", err)
	}

	return series, nil
}

// GetInProgressSeries gets series that the user is currently watching
func (s *coreSeriesService) GetInProgressSeries(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting in-progress series")

	series, err := s.seriesRepo.GetInProgressSeries(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Int("limit", limit).
			Msg("Failed to get in-progress series")
		return nil, fmt.Errorf("failed to get in-progress series: %w", err)
	}

	return series, nil
}

// GetSeriesByGenre gets series by genre
func (s *coreSeriesService) GetSeriesByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting series by genre")

	series, err := s.seriesRepo.GetSeriesByGenre(ctx, genre, limit)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Int("limit", limit).
			Msg("Failed to get series by genre")
		return nil, fmt.Errorf("failed to get series by genre: %w", err)
	}

	return series, nil
}

// GetSimilarSeries gets series similar to a given series
func (s *coreSeriesService) GetSimilarSeries(ctx context.Context, seriesID uint64, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Int("limit", limit).
		Msg("Getting similar series")

	series, err := s.seriesRepo.GetSimilarSeries(ctx, seriesID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Int("limit", limit).
			Msg("Failed to get similar series")
		return nil, fmt.Errorf("failed to get similar series: %w", err)
	}

	return series, nil
}

// GetSeriesInCollection gets all series in a collection
func (s *coreSeriesService) GetSeriesInCollection(ctx context.Context, collectionID uint64) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("collectionID", collectionID).
		Msg("Getting series in collection")

	series, err := s.seriesRepo.GetSeriesInCollection(ctx, collectionID)
	if err != nil {
		log.Error().Err(err).
			Uint64("collectionID", collectionID).
			Msg("Failed to get series in collection")
		return nil, fmt.Errorf("failed to get series in collection: %w", err)
	}

	return series, nil
}

// GetUpcomingEpisodes gets episodes that will air in the next few days
func (s *coreSeriesService) GetUpcomingEpisodes(ctx context.Context, days int) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Msg("Getting upcoming episodes")

	episodes, err := s.seriesRepo.GetUpcomingEpisodes(ctx, days)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Msg("Failed to get upcoming episodes")
		return nil, fmt.Errorf("failed to get upcoming episodes: %w", err)
	}

	return episodes, nil
}

// GetEpisodesAiredBetween gets episodes that aired between two dates
func (s *coreSeriesService) GetEpisodesAiredBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Time("startDate", startDate).
		Time("endDate", endDate).
		Msg("Getting episodes aired between dates")

	episodes, err := s.seriesRepo.GetEpisodesAiredBetween(ctx, startDate, endDate)
	if err != nil {
		log.Error().Err(err).
			Time("startDate", startDate).
			Time("endDate", endDate).
			Msg("Failed to get episodes aired between dates")
		return nil, fmt.Errorf("failed to get episodes aired between dates: %w", err)
	}

	return episodes, nil
}

// SearchSeriesLibrary performs a comprehensive search across all series items
func (s *coreSeriesService) SearchSeriesLibrary(ctx context.Context, query types.QueryOptions) (*models.MediaItemList, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Msg("Searching series library")

	results, err := s.seriesRepo.SearchSeriesLibrary(ctx, query)
	if err != nil {
		log.Error().Err(err).
			Str("query", query.Query).
			Msg("Failed to search series library")
		return nil, fmt.Errorf("failed to search series library: %w", err)
	}

	return results, nil
}

