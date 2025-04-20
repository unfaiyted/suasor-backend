package repository

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/types/models"
	"time"

	"gorm.io/gorm"
)

// CoreUserMediaItemDataRepository defines the interface for basic user media item data operations
// This focuses on core CRUD operations that apply to all media types
type CoreUserMediaItemDataRepository[T types.MediaData] interface {
	// Create creates a new user media item data entry
	Create(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)

	// GetByID retrieves a specific user media item data entry by ID
	GetByID(ctx context.Context, id uint64) (*models.UserMediaItemData[T], error)

	// Update updates an existing user media item data entry
	Update(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)

	// Delete removes a specific user media item data entry
	Delete(ctx context.Context, id uint64) error

	// GetByUserIDAndMediaItemID retrieves user media item data for a specific user and media item
	GetByUserIDAndMediaItemID(ctx context.Context, userID, mediaItemID uint64) (*models.UserMediaItemData[T], error)

	// HasUserMediaItemData checks if a user has data for a specific media item
	HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error)

	// Search finds user media item data based on a query object
	Search(ctx context.Context, query *types.QueryOptions) ([]*models.UserMediaItemData[T], error)
}

// coreUserMediaItemDataRepository implements CoreUserMediaItemDataRepository
type coreUserMediaItemDataRepository[T types.MediaData] struct {
	db *gorm.DB
}

// NewCoreUserMediaItemDataRepository creates a new core user media item data repository
func NewCoreUserMediaItemDataRepository[T types.MediaData](db *gorm.DB) CoreUserMediaItemDataRepository[T] {
	return &coreUserMediaItemDataRepository[T]{db: db}
}

// Create creates a new user media item data entry
func (r *coreUserMediaItemDataRepository[T]) Create(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	// Set timestamps if not already set
	if data.CreatedAt.IsZero() {
		data.CreatedAt = time.Now()
	}
	if data.UpdatedAt.IsZero() {
		data.UpdatedAt = time.Now()
	}
	if data.PlayedAt.IsZero() {
		data.PlayedAt = time.Now()
	}
	if data.LastPlayedAt.IsZero() {
		data.LastPlayedAt = time.Now()
	}

	result := r.db.WithContext(ctx).Table("user_media_item_data").Create(&data)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user media item data: %w", result.Error)
	}

	return data, nil
}

// GetByID retrieves a specific user media item data entry by ID
func (r *coreUserMediaItemDataRepository[T]) GetByID(ctx context.Context, id uint64) (*models.UserMediaItemData[T], error) {
	var data models.UserMediaItemData[T]

	result := r.db.WithContext(ctx).Table("user_media_item_data").Where("id = ?", id).First(&data)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user media item data not found: %w", result.Error)
		}
		return nil, fmt.Errorf("failed to get user media item data: %w", result.Error)
	}

	// Load associated media item
	if err := data.LoadItem(r.db); err != nil {
		// Log error but continue
		fmt.Printf("Error loading media item for data %d: %v\n", data.ID, err)
	}

	return &data, nil
}

// Update updates an existing user media item data entry
func (r *coreUserMediaItemDataRepository[T]) Update(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	// Set updated timestamp
	data.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Table("user_media_item_data").Where("id = ?", data.ID).Updates(&data)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update user media item data: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user media item data not found")
	}

	// Get the updated record
	return r.GetByID(ctx, data.ID)
}

// Delete removes a specific user media item data entry
func (r *coreUserMediaItemDataRepository[T]) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Table("user_media_item_data").Where("id = ?", id).Delete(&models.UserMediaItemData[T]{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user media item data: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user media item data not found")
	}

	return nil
}

// GetByUserIDAndMediaItemID retrieves user media item data for a specific user and media item
func (r *coreUserMediaItemDataRepository[T]) GetByUserIDAndMediaItemID(ctx context.Context, userID, mediaItemID uint64) (*models.UserMediaItemData[T], error) {
	var data models.UserMediaItemData[T]

	result := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		First(&data)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user media item data not found: %w", result.Error)
		}
		return nil, fmt.Errorf("failed to get user media item data: %w", result.Error)
	}

	// Load associated media item
	if err := data.LoadItem(r.db); err != nil {
		// Log error but continue
		fmt.Printf("Error loading media item for data %d: %v\n", data.ID, err)
	}

	return &data, nil
}

// HasUserMediaItemData checks if a user has data for a specific media item
func (r *coreUserMediaItemDataRepository[T]) HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).Table("user_media_item_data").
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check user media item data: %w", result.Error)
	}

	return count > 0, nil
}

// Search finds user media item data based on a query object
func (r *coreUserMediaItemDataRepository[T]) Search(ctx context.Context, query *types.QueryOptions) ([]*models.UserMediaItemData[T], error) {
	// var data models.UserMediaItemData[T]
	//
	// // Create a query options with user filter
	// options := types.QueryOptions{
	// 	MediaType: query.MediaType,
	// 	OwnerID:   query.OwnerID,
	// 	Query:     query.Query,
	// 	Limit:     query.Limit,
	// 	Offset:    query.Offset,
	// }
	//

	return nil, nil
}
