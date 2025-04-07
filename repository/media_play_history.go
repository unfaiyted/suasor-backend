package repository

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/types/models"

	"gorm.io/gorm"
)

// MediaPlayHistoryRepository defines the interface for media play history operations
type MediaPlayHistoryRepository interface {
	// CreateHistory creates a new play history entry
	CreateHistory(ctx context.Context, history interface{}) error
	// GetUserHistory retrieves a user's play history
	GetUserHistory(ctx context.Context, userID uint64, limit int) ([]interface{}, error)
	// GetRecentUserMovieHistory retrieves a user's recent movie history
	GetRecentUserMovieHistory(ctx context.Context, userID uint64, limit int) ([]models.MediaPlayHistory[*types.Movie], error)
	// GetRecentUserSeriesHistory retrieves a user's recent series history
	GetRecentUserSeriesHistory(ctx context.Context, userID uint64, limit int) ([]models.MediaPlayHistory[*types.Series], error)
	// GetRecentUserMusicHistory retrieves a user's recent music history
	GetRecentUserMusicHistory(ctx context.Context, userID uint64, limit int) ([]models.MediaPlayHistory[*types.Track], error)
	// HasUserViewedMedia checks if a user has viewed a specific media item
	HasUserViewedMedia(ctx context.Context, userID, mediaItemID uint64) (bool, error)
	// GetItemPlayCount gets the number of times a user has played a media item
	GetItemPlayCount(ctx context.Context, userID, mediaItemID uint64) (int, error)
}

type mediaPlayHistoryRepository struct {
	db *gorm.DB
}

// NewMediaPlayHistoryRepository creates a new media play history repository
func NewMediaPlayHistoryRepository(db *gorm.DB) MediaPlayHistoryRepository {
	return &mediaPlayHistoryRepository{db: db}
}

// CreateHistory creates a new play history entry
func (r *mediaPlayHistoryRepository) CreateHistory(ctx context.Context, history interface{}) error {
	result := r.db.WithContext(ctx).Create(history)
	if result.Error != nil {
		return fmt.Errorf("failed to create play history: %w", result.Error)
	}
	return nil
}

// GetUserHistory retrieves a user's play history
func (r *mediaPlayHistoryRepository) GetUserHistory(ctx context.Context, userID uint64, limit int) ([]interface{}, error) {
	// This is a simplified implementation that returns an empty slice
	// In a full implementation, we'd need to query each type of media history and combine the results
	return []interface{}{}, nil
}

// GetRecentUserMovieHistory retrieves a user's recent movie history
func (r *mediaPlayHistoryRepository) GetRecentUserMovieHistory(ctx context.Context, userID uint64, limit int) ([]models.MediaPlayHistory[*types.Movie], error) {
	var history []models.MediaPlayHistory[*types.Movie]

	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Joins("JOIN media_items ON media_play_histories.media_item_id = media_items.id").
		Where("media_items.type = ? AND media_play_histories.user_id = ?", types.MediaTypeMovie, userID).
		Order("media_play_histories.last_watched_at DESC").
		Limit(limit).
		Find(&history)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user movie history: %w", result.Error)
	}

	// Load the associated MediaItems
	for i := range history {
		if err := history[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
		}
	}

	return history, nil
}

// GetRecentUserSeriesHistory retrieves a user's recent series history
func (r *mediaPlayHistoryRepository) GetRecentUserSeriesHistory(ctx context.Context, userID uint64, limit int) ([]models.MediaPlayHistory[*types.Series], error) {
	var history []models.MediaPlayHistory[*types.Series]

	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Joins("JOIN media_items ON media_play_histories.media_item_id = media_items.id").
		Where("media_items.type = ? AND media_play_histories.user_id = ?", types.MediaTypeSeries, userID).
		Order("media_play_histories.last_watched_at DESC").
		Limit(limit).
		Find(&history)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user series history: %w", result.Error)
	}

	// Load the associated MediaItems
	for i := range history {
		if err := history[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
		}
	}

	return history, nil
}

// GetRecentUserMusicHistory retrieves a user's recent music history
func (r *mediaPlayHistoryRepository) GetRecentUserMusicHistory(ctx context.Context, userID uint64, limit int) ([]models.MediaPlayHistory[*types.Track], error) {
	var history []models.MediaPlayHistory[*types.Track]

	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Joins("JOIN media_items ON media_play_histories.media_item_id = media_items.id").
		Where("media_items.type = ? AND media_play_histories.user_id = ?", types.MediaTypeTrack, userID).
		Order("media_play_histories.last_watched_at DESC").
		Limit(limit).
		Find(&history)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user music history: %w", result.Error)
	}

	// Load the associated MediaItems
	for i := range history {
		if err := history[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
		}
	}

	return history, nil
}

// HasUserViewedMedia checks if a user has viewed a specific media item
func (r *mediaPlayHistoryRepository) HasUserViewedMedia(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check user viewing history: %w", result.Error)
	}

	return count > 0, nil
}

// GetItemPlayCount gets the number of times a user has played a media item
func (r *mediaPlayHistoryRepository) GetItemPlayCount(ctx context.Context, userID, mediaItemID uint64) (int, error) {
	var history struct {
		PlayCount int32
	}

	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Select("play_count").
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		First(&history)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get play count: %w", result.Error)
	}

	return int(history.PlayCount), nil
}

