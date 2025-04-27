package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"
)

// SeriesRepository interface defines operations specific to TV series in the database
// This repository handles specialized TV series queries and operations that work with
// the relationships between episodes, seasons, and series
type SeriesRepository interface {
	// Episode-related operations
	GetEpisodesBySeasonID(ctx context.Context, seasonID uint64) ([]*models.MediaItem[*types.Episode], error)
	GetEpisodesBySeriesID(ctx context.Context, seriesID uint64) ([]*models.MediaItem[*types.Episode], error)
	GetRecentlyAddedEpisodes(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Episode], error)
	GetUnwatchedEpisodes(ctx context.Context, userID uint64, seriesID uint64) ([]*models.MediaItem[*types.Episode], error)
	GetNextEpisodeToWatch(ctx context.Context, userID uint64, seriesID uint64) (*models.MediaItem[*types.Episode], error)
	GetEpisodeByNumber(ctx context.Context, seriesID uint64, seasonNumber int, episodeNumber int) (*models.MediaItem[*types.Episode], error)

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

	// Genre and attribute-based operations
	GetSeriesByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Series], error)
	GetEpisodesByAttribute(ctx context.Context, attribute string, value interface{}, limit int) ([]*models.MediaItem[*types.Episode], error)

	// Advanced search operations
	SearchSeriesLibrary(ctx context.Context, query types.QueryOptions) (*models.MediaItemList, error)
	GetSimilarSeries(ctx context.Context, seriesID uint64, limit int) ([]*models.MediaItem[*types.Series], error)

	// Calendar operations
	GetUpcomingEpisodes(ctx context.Context, days int) ([]*models.MediaItem[*types.Episode], error)
	GetEpisodesAiredBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*models.MediaItem[*types.Episode], error)

	// Collection operations
	GetSeriesInCollection(ctx context.Context, collectionID uint64) ([]*models.MediaItem[*types.Series], error)
}

// seriesRepository implements the SeriesRepository interface
type seriesRepository struct {
	db          *gorm.DB
	episodeRepo CoreMediaItemRepository[*types.Episode]
	seasonRepo  CoreMediaItemRepository[*types.Season]
	seriesRepo  CoreMediaItemRepository[*types.Series]
}

// NewSeriesRepository creates a new series repository
func NewSeriesRepository(
	db *gorm.DB,
	episodeRepo CoreMediaItemRepository[*types.Episode],
	seasonRepo CoreMediaItemRepository[*types.Season],
	seriesRepo CoreMediaItemRepository[*types.Series],
) SeriesRepository {
	return &seriesRepository{
		db:          db,
		episodeRepo: episodeRepo,
		seasonRepo:  seasonRepo,
		seriesRepo:  seriesRepo,
	}
}

// GetEpisodesBySeasonID retrieves all episodes for a specific season
func (r *seriesRepository) GetEpisodesBySeasonID(ctx context.Context, seasonID uint64) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seasonID", seasonID).
		Msg("Getting episodes by season ID")

	var episodes []*models.MediaItem[*types.Episode]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeEpisode).
		Where("data->>'seasonID' = ?", fmt.Sprint(seasonID)).
		Order("(data->>'episodeNumber')::int ASC").
		Find(&episodes).Error; err != nil {
		return nil, fmt.Errorf("failed to get episodes by season ID: %w", err)
	}

	return episodes, nil
}

// GetEpisodesBySeriesID retrieves all episodes for a specific series
func (r *seriesRepository) GetEpisodesBySeriesID(ctx context.Context, seriesID uint64) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Msg("Getting episodes by series ID")

	var episodes []*models.MediaItem[*types.Episode]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeEpisode).
		Where("data->>'seriesID' = ?", fmt.Sprint(seriesID)).
		Order("(data->>'seasonNumber')::int ASC, (data->>'episodeNumber')::int ASC").
		Find(&episodes).Error; err != nil {
		return nil, fmt.Errorf("failed to get episodes by series ID: %w", err)
	}

	return episodes, nil
}

// GetRecentlyAddedEpisodes retrieves recently added episodes
func (r *seriesRepository) GetRecentlyAddedEpisodes(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Episode], error) {
	return r.episodeRepo.GetRecentItems(ctx, days, limit)
}

// GetUnwatchedEpisodes retrieves unwatched episodes for a user
func (r *seriesRepository) GetUnwatchedEpisodes(ctx context.Context, userID uint64, seriesID uint64) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("seriesID", seriesID).
		Msg("Getting unwatched episodes")

	// This would require a join with user watch status data
	// For now, we'll implement a simpler version
	var episodes []*models.MediaItem[*types.Episode]

	// Get all episodes for the series first
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeEpisode).
		Where("data->>'seriesID' = ?", fmt.Sprint(seriesID)).
		Order("(data->>'seasonNumber')::int ASC, (data->>'episodeNumber')::int ASC").
		Find(&episodes).Error; err != nil {
		return nil, fmt.Errorf("failed to get episodes: %w", err)
	}

	// In a full implementation, we would filter out episodes that the user has watched
	// This would require a join with the user_media_item_data table
	// For now, we'll return all episodes

	return episodes, nil
}

// GetNextEpisodeToWatch retrieves the next episode to watch for a user
func (r *seriesRepository) GetNextEpisodeToWatch(ctx context.Context, userID uint64, seriesID uint64) (*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("seriesID", seriesID).
		Msg("Getting next episode to watch")

	// Get unwatched episodes
	episodes, err := r.GetUnwatchedEpisodes(ctx, userID, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unwatched episodes: %w", err)
	}

	// Return the first unwatched episode (they're already sorted by season and episode number)
	if len(episodes) > 0 {
		return episodes[0], nil
	}

	return nil, fmt.Errorf("no unwatched episodes found")
}

// GetEpisodeByNumber retrieves an episode by its season and episode number
func (r *seriesRepository) GetEpisodeByNumber(ctx context.Context, seriesID uint64, seasonNumber int, episodeNumber int) (*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeNumber", episodeNumber).
		Msg("Getting episode by number")

	var episode models.MediaItem[*types.Episode]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeEpisode).
		Where("data->>'seriesID' = ?", fmt.Sprint(seriesID)).
		Where("data->>'seasonNumber' = ?", fmt.Sprint(seasonNumber)).
		Where("data->>'episodeNumber' = ?", fmt.Sprint(episodeNumber)).
		First(&episode).Error; err != nil {
		return nil, fmt.Errorf("failed to get episode: %w", err)
	}

	return &episode, nil
}

// GetSeasonsBySeriesID retrieves all seasons for a specific series
func (r *seriesRepository) GetSeasonsBySeriesID(ctx context.Context, seriesID uint64) ([]*models.MediaItem[*types.Season], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Msg("Getting seasons by series ID")

	var seasons []*models.MediaItem[*types.Season]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeSeason).
		Where("data->>'seriesID' = ?", fmt.Sprint(seriesID)).
		Order("(data->>'seasonNumber')::int ASC").
		Find(&seasons).Error; err != nil {
		return nil, fmt.Errorf("failed to get seasons by series ID: %w", err)
	}

	return seasons, nil
}

// GetSeasonWithEpisodes retrieves a season and all its episodes
func (r *seriesRepository) GetSeasonWithEpisodes(ctx context.Context, seasonID uint64) (*models.MediaItem[*types.Season], []*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seasonID", seasonID).
		Msg("Getting season with episodes")

	// Get the season
	season, err := r.seasonRepo.GetByID(ctx, seasonID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get season: %w", err)
	}

	// Get the episodes
	episodes, err := r.GetEpisodesBySeasonID(ctx, seasonID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get episodes for season: %w", err)
	}

	return season, episodes, nil
}

// GetSeasonByNumber retrieves a season by its number within a series
func (r *seriesRepository) GetSeasonByNumber(ctx context.Context, seriesID uint64, seasonNumber int) (*models.MediaItem[*types.Season], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Msg("Getting season by number")

	var season models.MediaItem[*types.Season]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeSeason).
		Where("data->>'seriesID' = ?", fmt.Sprint(seriesID)).
		Where("data->>'seasonNumber' = ?", fmt.Sprint(seasonNumber)).
		First(&season).Error; err != nil {
		return nil, fmt.Errorf("failed to get season: %w", err)
	}

	return &season, nil
}

// GetRecentlyAddedSeasons retrieves recently added seasons
func (r *seriesRepository) GetRecentlyAddedSeasons(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Season], error) {
	return r.seasonRepo.GetRecentItems(ctx, days, limit)
}

// GetSeriesWithSeasons retrieves a series and all its seasons
func (r *seriesRepository) GetSeriesWithSeasons(ctx context.Context, seriesID uint64) (*models.MediaItem[*types.Series], []*models.MediaItem[*types.Season], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Msg("Getting series with seasons")

	// Get the series
	series, err := r.seriesRepo.GetByID(ctx, seriesID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get series: %w", err)
	}

	// Get the seasons
	seasons, err := r.GetSeasonsBySeriesID(ctx, seriesID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get seasons for series: %w", err)
	}

	return series, seasons, nil
}

// GetRecentlyAiredSeries retrieves recently aired series
func (r *seriesRepository) GetRecentlyAiredSeries(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recently aired series")

	// Calculate the cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -days)

	var series []*models.MediaItem[*types.Series]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeSeries).
		Where("data->>'lastAirDate' >= ?", cutoffDate.Format(time.RFC3339)).
		Order("data->>'lastAirDate' DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&series).Error; err != nil {
		return nil, fmt.Errorf("failed to get recently aired series: %w", err)
	}

	return series, nil
}

// GetPopularSeries retrieves popular series
func (r *seriesRepository) GetPopularSeries(ctx context.Context, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("Getting popular series")

	var series []*models.MediaItem[*types.Series]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeSeries).
		Order("(data->>'rating')::float DESC NULLS LAST")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&series).Error; err != nil {
		return nil, fmt.Errorf("failed to get popular series: %w", err)
	}

	return series, nil
}

// GetInProgressSeries retrieves series that the user is currently watching
func (r *seriesRepository) GetInProgressSeries(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting in-progress series")

	// This would require a join with user watch status data
	// For now, we'll implement a simpler version
	// var series []*models.MediaItem[*types.Series]

	// In a full implementation, we would find series where:
	// 1. The user has watched at least one episode
	// 2. There are still unwatched episodes
	// For now, we'll return popular series as a placeholder

	return r.GetPopularSeries(ctx, limit)
}

// GetSeriesByGenre retrieves series by genre
func (r *seriesRepository) GetSeriesByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting series by genre")

	var series []*models.MediaItem[*types.Series]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeSeries).
		Where("data->'genres' ? ?", genre)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&series).Error; err != nil {
		return nil, fmt.Errorf("failed to get series by genre: %w", err)
	}

	return series, nil
}

// GetEpisodesByAttribute retrieves episodes by a specific attribute
func (r *seriesRepository) GetEpisodesByAttribute(ctx context.Context, attribute string, value interface{}, limit int) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("attribute", attribute).
		Interface("value", value).
		Int("limit", limit).
		Msg("Getting episodes by attribute")

	var episodes []*models.MediaItem[*types.Episode]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeEpisode)

	switch attribute {
	case "writer", "director", "actor":
		// These would be array fields in the credits object
		query = query.Where(fmt.Sprintf("data->'credits'->'%s' ? ?", attribute), fmt.Sprintf("%v", value))
	case "rating", "isSpecial", "isSeason":
		// These are direct fields
		query = query.Where(fmt.Sprintf("data->>'%s' = ?", attribute), fmt.Sprintf("%v", value))
	default:
		// Default to exact match on any attribute
		query = query.Where(fmt.Sprintf("data->>'%s' = ?", attribute), fmt.Sprintf("%v", value))
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&episodes).Error; err != nil {
		return nil, fmt.Errorf("failed to get episodes by attribute: %w", err)
	}

	return episodes, nil
}

// SearchSeriesLibrary performs a comprehensive search across all series items
func (r *seriesRepository) SearchSeriesLibrary(ctx context.Context, query types.QueryOptions) (*models.MediaItemList, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Msg("Searching series library")

	seriesTypes := []types.MediaType{
		types.MediaTypeSeries,
		types.MediaTypeSeason,
		types.MediaTypeEpisode,
	}

	// Build the query for searching across all series types
	dbQuery := r.db.WithContext(ctx).
		Where("type IN ?", seriesTypes)

	if query.Query != "" {
		dbQuery = dbQuery.Where(
			"title ILIKE ? OR data->>'overview' ILIKE ?",
			"%"+query.Query+"%", "%"+query.Query+"%",
		)
	}

	// Execute separate queries for each type to populate the MediaItems struct
	var mediaItems models.MediaItemList = models.MediaItemList{}

	// Find series
	var seriesList []*models.MediaItem[*types.Series]
	if err := dbQuery.Where("type = ?", types.MediaTypeSeries).Find(&seriesList).Error; err != nil {
		return nil, fmt.Errorf("failed to search series: %w", err)
	}
	mediaItems.AddSeriesList(seriesList)

	// Find seasons
	var seasons []*models.MediaItem[*types.Season]
	if err := dbQuery.Where("type = ?", types.MediaTypeSeason).Find(&seasons).Error; err != nil {
		return nil, fmt.Errorf("failed to search seasons: %w", err)
	}
	mediaItems.AddSeasonList(seasons)

	// Find episodes
	var episodes []*models.MediaItem[*types.Episode]
	if err := dbQuery.Where("type = ?", types.MediaTypeEpisode).Find(&episodes).Error; err != nil {
		return nil, fmt.Errorf("failed to search episodes: %w", err)
	}
	mediaItems.AddEpisodeList(episodes)

	return &mediaItems, nil
}

// GetSimilarSeries finds series similar to a given series based on attributes
func (r *seriesRepository) GetSimilarSeries(ctx context.Context, seriesID uint64, limit int) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("seriesID", seriesID).
		Int("limit", limit).
		Msg("Getting similar series")

	// First get the source series
	sourceSeries, err := r.seriesRepo.GetByID(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source series: %w", err)
	}

	// Get the genres from the source series
	// In a real implementation, we'd also use networks, creators, etc.
	genres := sourceSeries.Data.Genres

	// Get similar series by genre
	var series []*models.MediaItem[*types.Series]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeSeries).
		Where("id != ?", seriesID)

	// Add genre conditions if there are any
	if len(genres) > 0 {
		for _, genre := range genres {
			query = query.Or("data->'genres' ? ?", genre)
		}
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&series).Error; err != nil {
		return nil, fmt.Errorf("failed to get similar series: %w", err)
	}

	return series, nil
}

// GetUpcomingEpisodes retrieves episodes that will air in the next few days
func (r *seriesRepository) GetUpcomingEpisodes(ctx context.Context, days int) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Msg("Getting upcoming episodes")

	now := time.Now()
	endDate := now.AddDate(0, 0, days)

	var episodes []*models.MediaItem[*types.Episode]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeEpisode).
		Where("data->>'airDate' >= ?", now.Format(time.RFC3339)).
		Where("data->>'airDate' <= ?", endDate.Format(time.RFC3339)).
		Order("data->>'airDate' ASC").
		Find(&episodes).Error; err != nil {
		return nil, fmt.Errorf("failed to get upcoming episodes: %w", err)
	}

	return episodes, nil
}

// GetEpisodesAiredBetween retrieves episodes that aired between two dates
func (r *seriesRepository) GetEpisodesAiredBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Time("startDate", startDate).
		Time("endDate", endDate).
		Msg("Getting episodes aired between dates")

	var episodes []*models.MediaItem[*types.Episode]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeEpisode).
		Where("data->>'airDate' >= ?", startDate.Format(time.RFC3339)).
		Where("data->>'airDate' <= ?", endDate.Format(time.RFC3339)).
		Order("data->>'airDate' ASC").
		Find(&episodes).Error; err != nil {
		return nil, fmt.Errorf("failed to get episodes aired between dates: %w", err)
	}

	return episodes, nil
}

// GetSeriesInCollection retrieves all series that are part of a collection
func (r *seriesRepository) GetSeriesInCollection(ctx context.Context, collectionID uint64) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("collectionID", collectionID).
		Msg("Getting series in collection")

	// First get the collection
	var collection models.MediaItem[*types.Collection]
	if err := r.db.WithContext(ctx).
		Where("id = ? AND type = ?", collectionID, types.MediaTypeCollection).
		First(&collection).Error; err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	itemList := collection.GetData().GetItemList()
	itemIDs := itemList.GetItemIDs()

	// Get the series IDs from the collection data
	if len(itemIDs) == 0 {
		return []*models.MediaItem[*types.Series]{}, nil
	}

	// Extract the series IDs
	var seriesIDs []uint64
	for _, itemID := range itemIDs {
		// Check if the item is a series
		var item models.MediaItem[types.ListData]
		if err := r.db.Where("id = ?", itemID).First(&item).Error; err == nil {
			if item.Type == types.MediaTypeSeries {
				seriesIDs = append(seriesIDs, itemID)
			}
		}
	}

	if len(seriesIDs) == 0 {
		return []*models.MediaItem[*types.Series]{}, nil
	}

	// Get the series
	var series []*models.MediaItem[*types.Series]
	if err := r.db.WithContext(ctx).
		Where("id IN ? AND type = ?", seriesIDs, types.MediaTypeSeries).
		Find(&series).Error; err != nil {
		return nil, fmt.Errorf("failed to get series in collection: %w", err)
	}

	return series, nil
}
