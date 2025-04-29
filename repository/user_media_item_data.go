package repository

import (
	"context"
	"fmt"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"

	"gorm.io/gorm"
)

// UserUserMediaItemDataRepository defines the interface for user-specific media item data operations
// This focuses on user-specific actions like favorites, ratings, watchlist, etc.
type UserMediaItemDataRepository[T types.MediaData] interface {
	CoreUserMediaItemDataRepository[T]

	// GetUserHistory retrieves a user's media item history
	GetUserHistory(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error)

	// GetRecentHistory retrieves a user's recent media history
	GetRecentHistory(ctx context.Context, userID uint64, limit int) ([]*models.UserMediaItemData[T], error)

	// GetUserPlayHistory retrieves play history for a user with optional filtering
	GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, completed *bool) ([]*models.UserMediaItemData[T], error)

	// GetContinueWatching retrieves items that a user has started but not completed
	GetContinueWatching(ctx context.Context, userID uint64, limit int) ([]*models.UserMediaItemData[T], error)

	// RecordPlay records a new play event
	RecordPlay(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)

	// ToggleFavorite marks or unmarks a media item as a favorite
	ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error

	// UpdateRating sets a user's rating for a media item
	UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error

	// GetFavorites retrieves favorite media items for a user
	GetFavorites(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error)

	HasUserViewedMedia(ctx context.Context, userID, mediaItemID uint64) (bool, error)

	GetItemPlayCount(ctx context.Context, userID, mediaItemID uint64) (int, error)

	// ClearUserHistory removes all data for a user
	ClearUserHistory(ctx context.Context, userID uint64) error
}

// userMediaItemDataRepository implements UserUserMediaItemDataRepository
type userMediaItemDataRepository[T types.MediaData] struct {
	CoreUserMediaItemDataRepository[T]
	db *gorm.DB
}

// NewUserUserMediaItemDataRepository creates a new user media item data repository
func NewUserMediaItemDataRepository[T types.MediaData](db *gorm.DB, CoreUserMediaItemDataRepository CoreUserMediaItemDataRepository[T]) UserMediaItemDataRepository[T] {
	return &userMediaItemDataRepository[T]{
		CoreUserMediaItemDataRepository: CoreUserMediaItemDataRepository,
		db:                              db,
	}
}

// GetUserHistory retrieves a user's media item history
func (r *userMediaItemDataRepository[T]) GetUserHistory(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error) {
	var history []*models.UserMediaItemData[T]

	query := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ?", userID).
		Order("last_played_at DESC")

	// Apply media type filter if provided
	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)
	query = query.Where("type = ?", mediaType)

	// Apply pagination
	query = query.Limit(limit).Offset(offset)

	// Execute the query
	result := query.Find(&history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user history: %w", result.Error)
	}

	// Load associated media items
	for i := range history {
		if err := history[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
		}
	}

	return history, nil
}

// GetRecentHistory retrieves a user's recent media history
func (r *userMediaItemDataRepository[T]) GetRecentHistory(ctx context.Context, userID uint64, limit int) ([]*models.UserMediaItemData[T], error) {
	var history []*models.UserMediaItemData[T]

	query := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ?", userID).
		Order("last_played_at DESC").
		Limit(limit)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)
	query = query.Where("type = ?", mediaType)

	// Execute the query
	result := query.Find(&history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recent user history: %w", result.Error)
	}

	// Load associated media items
	for i := range history {
		if err := history[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
		}
	}

	return history, nil
}

// GetUserPlayHistory retrieves play history for a user with optional filtering
func (r *userMediaItemDataRepository[T]) GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, completed *bool) ([]*models.UserMediaItemData[T], error) {
	var history []*models.UserMediaItemData[T]

	query := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ?", userID).
		Order("last_played_at DESC")

	// Apply media type filter if provided
	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)
	query = query.Where("type = ?", mediaType)

	// Filter by completion status if provided
	if completed != nil {
		query = query.Where("completed = ?", *completed)
	}

	// Apply pagination
	query = query.Limit(limit).Offset(offset)

	// Execute the query
	result := query.Find(&history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user play history: %w", result.Error)
	}

	// Load associated media items
	for i := range history {
		if err := history[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
		}
	}

	return history, nil
}

// GetContinueWatching retrieves items that a user has started but not completed
func (r *userMediaItemDataRepository[T]) GetContinueWatching(ctx context.Context, userID uint64, limit int) ([]*models.UserMediaItemData[T], error) {
	var history []*models.UserMediaItemData[T]

	query := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ? AND completed = ? AND played_percentage > ? AND played_percentage < ?",
			userID, false, 0.0, 0.95).
		Order("last_played_at DESC").
		Limit(limit)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	query = query.Where("type = ?", mediaType)

	// Execute the query
	result := query.Find(&history)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get continue watching items: %w", result.Error)
	}

	// Load associated media items
	for i := range history {
		if err := history[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for history %d: %v\n", history[i].ID, err)
		}
	}

	return history, nil
}

// RecordPlay records a new play event
func (r *userMediaItemDataRepository[T]) RecordPlay(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	// Check if there's an existing record
	existingData, err := r.CoreUserMediaItemDataRepository.GetByUserIDAndMediaItemID(ctx, data.UserID, data.MediaItemID)
	if err != nil {
		// If it's not a "not found" error, return the error
		if err.Error() != "user media item data not found: record not found" {
			return nil, fmt.Errorf("error checking for existing data: %w", err)
		}

		// No existing record, create a new one
		data.CreatedAt = time.Now()
		data.UpdatedAt = time.Now()
		data.PlayedAt = time.Now()
		data.LastPlayedAt = time.Now()

		// Increment play count
		data.PlayCount = 1

		result, err := r.CoreUserMediaItemDataRepository.Create(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("failed to create play record: %w", err)
		}

		return result, nil
	}

	// Update existing record
	existingData.LastPlayedAt = time.Now()
	existingData.PlayCount++
	existingData.PlayedPercentage = data.PlayedPercentage
	existingData.PositionSeconds = data.PositionSeconds
	existingData.DurationSeconds = data.DurationSeconds
	existingData.Completed = data.Completed

	// Update the record
	result, err := r.CoreUserMediaItemDataRepository.Update(ctx, existingData)
	if err != nil {
		return nil, fmt.Errorf("failed to update play record: %w", err)
	}

	return result, nil
}

// ToggleFavorite marks or unmarks a media item as a favorite
func (r *userMediaItemDataRepository[T]) ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error {
	// Check if there's an existing record
	log := logger.LoggerFromContext(ctx)
	existingData, err := r.CoreUserMediaItemDataRepository.
		GetByUserIDAndMediaItemID(ctx, userID, mediaItemID)
	if err != nil {
		// If it's not a "not found" error, return the error
		if err.Error() != "user media item data not found: record not found" {
			return fmt.Errorf("error checking for existing data: %w", err)
		}

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

		// No existing record, create a new one
		newData := models.UserMediaItemData[T]{
			UserID:       userID,
			MediaItemID:  mediaItemID,
			Type:         mediaItem.Type,
			IsFavorite:   favorite,
			PlayedAt:     time.Now(),
			LastPlayedAt: time.Now(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		result, err := r.CoreUserMediaItemDataRepository.Create(ctx, &newData)
		if err != nil {
			return fmt.Errorf("failed to create favorite record: %w", err)
		}

		log.Info().
			Uint64("mediaItemID", result.MediaItemID).
			Msg("Created favorite record")

		return nil
	}

	// Update existing record
	existingData.IsFavorite = favorite
	existingData.UpdatedAt = time.Now()

	// Update the record
	result, err := r.CoreUserMediaItemDataRepository.Update(ctx, existingData)
	if err != nil {
		return fmt.Errorf("failed to update favorite status: %w", err)
	}
	log.Info().
		Uint64("mediaItemID", result.MediaItemID).
		Msg("Updated favorite status")

	return nil
}

// UpdateRating sets a user's rating for a media item
func (r *userMediaItemDataRepository[T]) UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error {
	// Check if there's an existing record
	log := logger.LoggerFromContext(ctx)
	existingData, err := r.CoreUserMediaItemDataRepository.GetByUserIDAndMediaItemID(ctx, userID, mediaItemID)
	if err != nil {
		// If it's not a "not found" error, return the error
		if err.Error() != "user media item data not found: record not found" {
			return fmt.Errorf("error checking for existing data: %w", err)
		}

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

		// No existing record, create a new one
		newData := models.UserMediaItemData[T]{
			UserID:       userID,
			MediaItemID:  mediaItemID,
			Type:         mediaItem.Type,
			UserRating:   rating,
			PlayedAt:     time.Now(),
			LastPlayedAt: time.Now(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		result, err := r.CoreUserMediaItemDataRepository.Create(ctx, &newData)
		if err != nil {
			return fmt.Errorf("failed to create rating record: %w", err)
		}
		log.Info().
			Uint64("mediaItemID", result.MediaItemID).
			Msg("Created rating record")

		return nil
	}

	// Update existing record
	existingData.UserRating = rating
	existingData.UpdatedAt = time.Now()

	// Update the record
	result, err := r.CoreUserMediaItemDataRepository.Update(ctx, existingData)
	if err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}
	log.Info().
		Uint64("mediaItemID", result.MediaItemID).
		Msg("Updated rating")

	return nil
}

// GetFavorites retrieves favorite media items for a user
func (r *userMediaItemDataRepository[T]) GetFavorites(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error) {
	var favorites []*models.UserMediaItemData[T]

	query := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ? AND is_favorite = ?", userID, true).
		Order("last_played_at DESC").
		Limit(limit).Offset(offset)

	// Execute the query
	result := query.Find(&favorites)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user favorites: %w", result.Error)
	}

	// Load associated media items
	for i := range favorites {
		if err := favorites[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for favorite %d: %v\n", favorites[i].ID, err)
		}
	}

	return favorites, nil
}

// ClearUserHistory removes all data for a user
func (r *userMediaItemDataRepository[T]) ClearUserHistory(ctx context.Context, userID uint64) error {
	result := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ?", userID).
		Delete(&models.UserMediaItemData[T]{})

	if result.Error != nil {
		return fmt.Errorf("failed to clear user history: %w", result.Error)
	}

	return nil
}

func (r *userMediaItemDataRepository[T]) HasUserViewedMedia(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check user viewing history: %w", result.Error)
	}

	return count > 0, nil
}

func (r *userMediaItemDataRepository[T]) GetItemPlayCount(ctx context.Context, userID, mediaItemID uint64) (int, error) {
	var history struct {
		PlayCount int32
	}

	result := r.db.WithContext(ctx).Table("user_media_item_data").
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

func (r *userMediaItemDataRepository[T]) HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check user viewing history: %w", result.Error)
	}

	return count > 0, nil
}
