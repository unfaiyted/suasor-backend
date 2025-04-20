package repository

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/types/models"
	"time"

	"gorm.io/gorm"
)

// ClientUserMediaItemDataRepository defines the interface for client-specific media item data operations
// This focuses on client-specific interactions and synchronization with external media systems
type ClientUserMediaItemDataRepository[T types.MediaData] interface {
	CoreUserMediaItemDataRepository[T]
	// SyncClientItemData synchronizes user media item data from an external client
	SyncClientItemData(ctx context.Context, userID uint64, clientID uint64, items []models.UserMediaItemData[T]) error

	// GetClientItemData retrieves user media item data for synchronization with a client
	GetClientItemData(ctx context.Context, userID uint64, clientID uint64, since time.Time) ([]*models.UserMediaItemData[T], error)

	// GetByClientID retrieves a user media item data entry by client ID
	GetByClientID(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error)

	// RecordClientPlay records a play event from a client
	RecordClientPlay(ctx context.Context, userID uint64, clientID uint64, clientItemID string, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)

	// MapClientMediaItemToInternal maps a client media item to an internal media item
	MapClientMediaItemToInternal(ctx context.Context, clientID uint64, clientItemID string) (uint64, error)

	// GetPlaybackState retrieves the current playback state for a client item
	GetPlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error)

	// UpdatePlaybackState updates the playback state for a client item
	UpdatePlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[T], error)
}

// clientUserMediaItemDataRepository implements ClientUserMediaItemDataRepository
type clientUserMediaItemDataRepository[T types.MediaData] struct {
	CoreUserMediaItemDataRepository[T]
	db       *gorm.DB
	userRepo UserMediaItemDataRepository[T]
	helper   *ClientMediaItemHelper
}

// NewClientUserMediaItemDataRepository creates a new client user media item data repository
func NewClientUserMediaItemDataRepository[T types.MediaData](
	db *gorm.DB,
	CoreUserMediaItemDataRepository CoreUserMediaItemDataRepository[T],
	userRepo UserMediaItemDataRepository[T],
) ClientUserMediaItemDataRepository[T] {
	return &clientUserMediaItemDataRepository[T]{
		CoreUserMediaItemDataRepository: CoreUserMediaItemDataRepository,
		db:                              db,
		userRepo:                        userRepo,
		helper:                          NewClientMediaItemHelper(db),
	}
}

// SyncClientItemData synchronizes user media item data from an external client
func (r *clientUserMediaItemDataRepository[T]) SyncClientItemData(ctx context.Context, userID uint64, clientID uint64, items []models.UserMediaItemData[T]) error {
	// Begin a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Defer rolling back the transaction in case something goes wrong
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Process each item
	for _, item := range items {
		// Map client's media item to our internal ID if needed
		if item.MediaItemID == 0 && item.Item != nil && item.Item.SyncClients.GetClientItemID(clientID) != "" {
			// Use our helper to map client item ID to internal ID
			clientItemID := item.Item.SyncClients.GetClientItemID(clientID)
			internalID, err := r.helper.GetMediaItemByClientID(ctx, clientID, clientItemID)

			if err != nil {
				// If not found and we have enough info, we could create a new media item
				// For now, let's skip items we can't map
				continue
			}

			item.MediaItemID = internalID
		}

		// Skip items with no internal ID mapping
		if item.MediaItemID == 0 {
			continue
		}

		// Check if an entry already exists
		var existingData models.UserMediaItemData[T]
		result := tx.Table("user_media_item_data").
			Where("user_id = ? AND media_item_id = ?", userID, item.MediaItemID).
			First(&existingData)

		if result.Error != nil {
			// If not found, create a new entry
			if result.Error == gorm.ErrRecordNotFound {
				// Set user ID and timestamps
				item.UserID = userID
				item.CreatedAt = time.Now()
				item.UpdatedAt = time.Now()

				// Create the item
				if err := tx.Table("user_media_item_data").Create(&item).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to create user media item data: %w", err)
				}

				continue
			}

			return fmt.Errorf("failed to check for existing user media item data: %w", result.Error)
		}

		// Update existing entry
		// Only update fields that are newer from the client
		if item.LastPlayedAt.After(existingData.LastPlayedAt) {
			updates := map[string]interface{}{
				"last_played_at":    item.LastPlayedAt,
				"played_percentage": item.PlayedPercentage,
				"position_seconds":  item.PositionSeconds,
				"duration_seconds":  item.DurationSeconds,
				"play_count":        item.PlayCount,
				"completed":         item.Completed,
				"updated_at":        time.Now(),
			}

			// Don't overwrite user preferences like ratings and favorites unless provided
			if item.UserRating > 0 {
				updates["user_rating"] = item.UserRating
			}

			// Only update favorite status if it's explicitly set
			if item.IsFavorite {
				updates["is_favorite"] = item.IsFavorite
			}

			if err := tx.Table("user_media_item_data").
				Where("id = ?", existingData.ID).
				Updates(updates).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update user media item data: %w", err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetClientItemData retrieves user media item data for synchronization with a client
func (r *clientUserMediaItemDataRepository[T]) GetClientItemData(ctx context.Context, userID uint64, clientID uint64, since time.Time) ([]*models.UserMediaItemData[T], error) {
	var items []*models.UserMediaItemData[T]

	query := r.db.WithContext(ctx).Table("user_media_item_data").
		Joins("JOIN media_items ON user_media_item_data.media_item_id = media_items.id").
		Where("user_media_item_data.user_id = ? AND media_items.provider = ? AND user_media_item_data.updated_at > ?",
			userID, clientID, since).
		Order("user_media_item_data.updated_at DESC")

	// Execute the query
	result := query.Find(&items)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get client item data: %w", result.Error)
	}

	// Load associated media items
	for i := range items {
		if err := items[i].LoadItem(r.db); err != nil {
			// Log the error but continue
			fmt.Printf("Error loading media item for data %d: %v\n", items[i].ID, err)
		}
	}

	return items, nil
}

// GetByClientID retrieves a user media item data entry by client ID
func (r *clientUserMediaItemDataRepository[T]) GetByClientID(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error) {
	// First, get the internal media item ID from the client item ID
	internalID, err := r.MapClientMediaItemToInternal(ctx, clientID, clientItemID)
	if err != nil {
		return nil, err
	}

	// Now get the user media item data using the internal ID
	return r.CoreUserMediaItemDataRepository.GetByUserIDAndMediaItemID(ctx, userID, internalID)
}

// RecordClientPlay records a play event from a client
func (r *clientUserMediaItemDataRepository[T]) RecordClientPlay(ctx context.Context, userID uint64, clientID uint64, clientItemID string, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	// First, get the internal media item ID from the client item ID
	internalID, err := r.MapClientMediaItemToInternal(ctx, clientID, clientItemID)
	if err != nil {
		return nil, err
	}

	// Set the internal media item ID
	data.MediaItemID = internalID
	data.UserID = userID

	// Use the user repository to record the play
	return r.userRepo.RecordPlay(ctx, data)
}

// MapClientMediaItemToInternal maps a client media item to an internal media item
func (r *clientUserMediaItemDataRepository[T]) MapClientMediaItemToInternal(ctx context.Context, clientID uint64, clientItemID string) (uint64, error) {
	// Use the helper to get media item by client ID
	return r.helper.GetMediaItemByClientID(ctx, clientID, clientItemID)
}

// GetPlaybackState retrieves the current playback state for a client item
func (r *clientUserMediaItemDataRepository[T]) GetPlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error) {
	// First, get the internal media item ID from the client item ID
	internalID, err := r.MapClientMediaItemToInternal(ctx, clientID, clientItemID)
	if err != nil {
		return nil, err
	}

	// Now get the user media item data using the internal ID
	return r.CoreUserMediaItemDataRepository.GetByUserIDAndMediaItemID(ctx, userID, internalID)
}

// UpdatePlaybackState updates the playback state for a client item
func (r *clientUserMediaItemDataRepository[T]) UpdatePlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[T], error) {
	// First, get the internal media item ID from the client item ID
	internalID, err := r.MapClientMediaItemToInternal(ctx, clientID, clientItemID)
	if err != nil {
		return nil, err
	}

	// Check if there's an existing record
	existingData, err := r.CoreUserMediaItemDataRepository.GetByUserIDAndMediaItemID(ctx, userID, internalID)
	if err != nil {
		// If it's not a "not found" error, return the error
		if err.Error() != "user media item data not found: record not found" {
			return nil, fmt.Errorf("error checking for existing data: %w", err)
		}

		// Get the media item type
		var mediaItem struct {
			Type types.MediaType
		}

		if err := r.db.WithContext(ctx).
			Table("media_items").
			Select("type").
			Where("id = ?", internalID).
			First(&mediaItem).Error; err != nil {
			return nil, fmt.Errorf("failed to get media item type: %w", err)
		}

		// No existing record, create a new one
		newData := models.UserMediaItemData[T]{
			UserID:           userID,
			MediaItemID:      internalID,
			Type:             mediaItem.Type,
			PositionSeconds:  position,
			DurationSeconds:  duration,
			PlayedPercentage: percentage,
			PlayedAt:         time.Now(),
			LastPlayedAt:     time.Now(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			PlayCount:        1,
			Completed:        percentage >= 0.9, // Consider complete if 90% played
		}

		result, err := r.CoreUserMediaItemDataRepository.Create(ctx, &newData)
		if err != nil {
			return nil, fmt.Errorf("failed to create playback state record: %w", err)
		}

		return result, nil
	}

	// Update existing record
	existingData.PositionSeconds = position
	existingData.DurationSeconds = duration
	existingData.PlayedPercentage = percentage
	existingData.LastPlayedAt = time.Now()
	existingData.UpdatedAt = time.Now()

	// If it wasn't completed before but is now, increment play count
	if !existingData.Completed && percentage >= 0.9 {
		existingData.PlayCount++
		existingData.Completed = true
	}

	// Update the record
	result, err := r.CoreUserMediaItemDataRepository.Update(ctx, existingData)
	if err != nil {
		return nil, fmt.Errorf("failed to update playback state: %w", err)
	}

	return result, nil
}
