package repository

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/types/models"
	"time"

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
	// GetByUserIDAndTypeMovie retrieves a user's movie play history with limit and offset
	GetByUserIDAndTypeMovie(ctx context.Context, userID uint64, limit int, offset int) ([]models.MediaPlayHistory[*types.Movie], error)
	// GetByUserIDAndTypeSeries retrieves a user's series play history with limit and offset
	GetByUserIDAndTypeSeries(ctx context.Context, userID uint64, limit int, offset int) ([]models.MediaPlayHistory[*types.Series], error)
	// GetByUserIDAndTypeMusic retrieves a user's music play history with limit and offset
	GetByUserIDAndTypeMusic(ctx context.Context, userID uint64, limit int, offset int) ([]models.MediaPlayHistory[*types.Track], error)
	
	// New methods for direct API endpoints
	
	// GetUserPlayHistory retrieves play history for a user with optional filtering
	GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) (interface{}, error)
	
	// GetContinueWatching retrieves items that a user has started but not completed
	GetContinueWatching(ctx context.Context, userID uint64, limit int) (interface{}, error)
	
	// GetByID retrieves a specific play history entry by ID
	GetByID(ctx context.Context, id uint64) (interface{}, error)
	
	// GetByMediaItemID retrieves play history for a specific media item
	GetByMediaItemID(ctx context.Context, mediaItemID, userID uint64) (interface{}, error)
	
	// RecordPlay records a new play event
	RecordPlay(ctx context.Context, history *models.MediaPlayHistoryGeneric) (interface{}, error)
	
	// ToggleFavorite marks or unmarks a media item as a favorite
	ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error
	
	// UpdateRating sets a user's rating for a media item
	UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error
	
	// GetFavorites retrieves favorite media items for a user
	GetFavorites(ctx context.Context, userID uint64, mediaType *types.MediaType, limit, offset int) (interface{}, error)
	
	// Delete removes a specific play history entry
	Delete(ctx context.Context, id uint64) error
	
	// ClearUserHistory removes all play history for a user
	ClearUserHistory(ctx context.Context, userID uint64, mediaType *types.MediaType) error
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
		Order("media_play_histories.last_played_at DESC").
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
		Order("media_play_histories.last_played_at DESC").
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
		Order("media_play_histories.last_played_at DESC").
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

// GetByUserIDAndTypeMovie retrieves a user's movie play history with limit and offset
func (r *mediaPlayHistoryRepository) GetByUserIDAndTypeMovie(ctx context.Context, userID uint64, limit int, offset int) ([]models.MediaPlayHistory[*types.Movie], error) {
	var history []models.MediaPlayHistory[*types.Movie]

	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Joins("JOIN media_items ON media_play_histories.media_item_id = media_items.id").
		Where("media_items.type = ? AND media_play_histories.user_id = ?", types.MediaTypeMovie, userID).
		Order("media_play_histories.last_played_at DESC").
		Limit(limit).
		Offset(offset).
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

// GetByUserIDAndTypeSeries retrieves a user's series play history with limit and offset
func (r *mediaPlayHistoryRepository) GetByUserIDAndTypeSeries(ctx context.Context, userID uint64, limit int, offset int) ([]models.MediaPlayHistory[*types.Series], error) {
	var history []models.MediaPlayHistory[*types.Series]

	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Joins("JOIN media_items ON media_play_histories.media_item_id = media_items.id").
		Where("media_items.type = ? AND media_play_histories.user_id = ?", types.MediaTypeSeries, userID).
		Order("media_play_histories.last_played_at DESC").
		Limit(limit).
		Offset(offset).
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

// GetByUserIDAndTypeMusic retrieves a user's music play history with limit and offset
func (r *mediaPlayHistoryRepository) GetByUserIDAndTypeMusic(ctx context.Context, userID uint64, limit int, offset int) ([]models.MediaPlayHistory[*types.Track], error) {
	var history []models.MediaPlayHistory[*types.Track]

	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Joins("JOIN media_items ON media_play_histories.media_item_id = media_items.id").
		Where("media_items.type = ? AND media_play_histories.user_id = ?", types.MediaTypeTrack, userID).
		Order("media_play_histories.last_played_at DESC").
		Limit(limit).
		Offset(offset).
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

// Implementation for the new methods

// GetUserPlayHistory retrieves play history for a user with optional filtering
func (r *mediaPlayHistoryRepository) GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) (interface{}, error) {
	// Simplified implementation that delegates to existing methods
	if mediaType != nil {
		switch *mediaType {
		case types.MediaTypeMovie:
			return r.GetByUserIDAndTypeMovie(ctx, userID, limit, offset)
		case types.MediaTypeSeries:
			return r.GetByUserIDAndTypeSeries(ctx, userID, limit, offset)
		case types.MediaTypeTrack:
			return r.GetByUserIDAndTypeMusic(ctx, userID, limit, offset)
		}
	}
	
	// If no specific type, just use the existing method for now
	return r.GetUserHistory(ctx, userID, limit)
}

// GetContinueWatching retrieves items that a user has started but not completed
func (r *mediaPlayHistoryRepository) GetContinueWatching(ctx context.Context, userID uint64, limit int) (interface{}, error) {
	// For now, just return recent movie history as a placeholder
	// This would need a proper implementation to filter for incomplete items
	return r.GetRecentUserMovieHistory(ctx, userID, limit)
}

// GetByID retrieves a specific play history entry by ID
func (r *mediaPlayHistoryRepository) GetByID(ctx context.Context, id uint64) (interface{}, error) {
	// Placeholder implementation
	var historyEntry struct {
		ID uint64
	}
	
	result := r.db.WithContext(ctx).Table("media_play_histories").First(&historyEntry, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get history entry: %w", result.Error)
	}
	
	// In a real implementation, you would determine the type and return the appropriate entry
	return historyEntry, nil
}

// GetByMediaItemID retrieves play history for a specific media item
func (r *mediaPlayHistoryRepository) GetByMediaItemID(ctx context.Context, mediaItemID, userID uint64) (interface{}, error) {
	// Placeholder implementation
	var mediaItems []struct {
		ID uint64
	}
	
	result := r.db.WithContext(ctx).
		Table("media_play_histories").
		Where("media_item_id = ? AND user_id = ?", mediaItemID, userID).
		Find(&mediaItems)
		
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get history entries: %w", result.Error)
	}
	
	return mediaItems, nil
}

// RecordPlay records a new play event
func (r *mediaPlayHistoryRepository) RecordPlay(ctx context.Context, history *models.MediaPlayHistoryGeneric) (interface{}, error) {
	// Simplified implementation that delegates to existing method
	result := r.db.WithContext(ctx).Create(history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create play history: %w", result.Error)
	}
	
	return history, nil
}

// ToggleFavorite marks or unmarks a media item as a favorite
func (r *mediaPlayHistoryRepository) ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error {
	// Check if there's an existing record
	var existingID uint64
	err := r.db.WithContext(ctx).
		Table("media_play_histories").
		Select("id").
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		Scan(&existingID).Error
		
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking for existing play history: %w", err)
	}
	
	if existingID > 0 {
		// Update existing record
		if err := r.db.WithContext(ctx).
			Table("media_play_histories").
			Where("id = ?", existingID).
			Update("is_favorite", favorite).Error; err != nil {
			return fmt.Errorf("failed to update favorite status: %w", err)
		}
	} else {
		// Get the media item type
		var mediaItem struct {
			Type types.MediaType
		}
		
		if err := r.db.WithContext(ctx).
			Table("media_items").
			Select("type").
			Where("id = ?", mediaItemID).
			First(&mediaItem).Error; err != nil {
			return fmt.Errorf("failed to get media item type: %w", err)
		}
		
		// Create new history record with favorite status
		newHistory := models.MediaPlayHistoryGeneric{
			UserID:      userID,
			MediaItemID: mediaItemID,
			Type:        mediaItem.Type,
			PlayedAt:    time.Now(),
			LastPlayedAt: time.Now(),
			IsFavorite:  favorite,
		}
		
		if err := r.db.WithContext(ctx).
			Table("media_play_histories").
			Create(&newHistory).Error; err != nil {
			return fmt.Errorf("failed to create play history with favorite status: %w", err)
		}
	}
	
	return nil
}

// UpdateRating sets a user's rating for a media item
func (r *mediaPlayHistoryRepository) UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error {
	// Check if there's an existing record
	var existingID uint64
	err := r.db.WithContext(ctx).
		Table("media_play_histories").
		Select("id").
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		Scan(&existingID).Error
		
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking for existing play history: %w", err)
	}
	
	if existingID > 0 {
		// Update existing record
		if err := r.db.WithContext(ctx).
			Table("media_play_histories").
			Where("id = ?", existingID).
			Update("user_rating", rating).Error; err != nil {
			return fmt.Errorf("failed to update rating: %w", err)
		}
	} else {
		// Get the media item type
		var mediaItem struct {
			Type types.MediaType
		}
		
		if err := r.db.WithContext(ctx).
			Table("media_items").
			Select("type").
			Where("id = ?", mediaItemID).
			First(&mediaItem).Error; err != nil {
			return fmt.Errorf("failed to get media item type: %w", err)
		}
		
		// Create new history record with rating
		newHistory := models.MediaPlayHistoryGeneric{
			UserID:      userID,
			MediaItemID: mediaItemID,
			Type:        mediaItem.Type,
			PlayedAt:    time.Now(),
			LastPlayedAt: time.Now(),
			UserRating:  rating,
		}
		
		if err := r.db.WithContext(ctx).
			Table("media_play_histories").
			Create(&newHistory).Error; err != nil {
			return fmt.Errorf("failed to create play history with rating: %w", err)
		}
	}
	
	return nil
}

// GetFavorites retrieves favorite media items for a user
func (r *mediaPlayHistoryRepository) GetFavorites(ctx context.Context, userID uint64, mediaType *types.MediaType, limit, offset int) (interface{}, error) {
	// Query to get favorite items
	query := r.db.WithContext(ctx).
		Table("media_play_histories").
		Joins("JOIN media_items ON media_play_histories.media_item_id = media_items.id").
		Where("media_play_histories.user_id = ? AND media_play_histories.is_favorite = ?", userID, true)
	
	// Add media type filter if provided
	if mediaType != nil {
		query = query.Where("media_items.type = ?", *mediaType)
	}
	
	// Add pagination
	query = query.
		Order("media_play_histories.last_played_at DESC").
		Limit(limit).
		Offset(offset)
	
	// If a specific media type is provided, return the appropriate type
	if mediaType != nil {
		switch *mediaType {
		case types.MediaTypeMovie:
			var history []models.MediaPlayHistory[*types.Movie]
			if err := query.Find(&history).Error; err != nil {
				return nil, fmt.Errorf("failed to get favorite movie items: %w", err)
			}
			
			// Load the associated MediaItems
			for i := range history {
				if err := history[i].LoadItem(r.db); err != nil {
					// Log the error but continue
					fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
				}
			}
			
			return history, nil
			
		case types.MediaTypeSeries:
			var history []models.MediaPlayHistory[*types.Series]
			if err := query.Find(&history).Error; err != nil {
				return nil, fmt.Errorf("failed to get favorite series items: %w", err)
			}
			
			// Load the associated MediaItems
			for i := range history {
				if err := history[i].LoadItem(r.db); err != nil {
					// Log the error but continue
					fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
				}
			}
			
			return history, nil
			
		case types.MediaTypeTrack:
			var history []models.MediaPlayHistory[*types.Track]
			if err := query.Find(&history).Error; err != nil {
				return nil, fmt.Errorf("failed to get favorite track items: %w", err)
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
	}
	
	// If no specific type or an unsupported type, return empty
	return []interface{}{}, nil
}

// Delete removes a specific play history entry
func (r *mediaPlayHistoryRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).
		Table("media_play_histories").
		Where("id = ?", id).
		Delete(nil).Error; err != nil {
		return fmt.Errorf("failed to delete play history entry: %w", err)
	}
	
	return nil
}

// ClearUserHistory removes all play history for a user
func (r *mediaPlayHistoryRepository) ClearUserHistory(ctx context.Context, userID uint64, mediaType *types.MediaType) error {
	query := r.db.WithContext(ctx).
		Table("media_play_histories").
		Where("user_id = ?", userID)
	
	// Add media type filter if provided
	if mediaType != nil {
		query = query.Where("type = ?", *mediaType)
	}
	
	if err := query.Delete(nil).Error; err != nil {
		return fmt.Errorf("failed to clear user play history: %w", err)
	}
	
	return nil
}