package repository

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/types/models"

	"gorm.io/gorm"
)

// MediaItemRepository defines the interface for media item database operations
type MediaItemRepository[T types.MediaData] interface {
	Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)

	// Retrieval operations
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
	GetByExternalID(ctx context.Context, externalID string, clientID uint64) (*models.MediaItem[T], error)
	GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error)
	GetByType(ctx context.Context, mediaType types.MediaType, clientID uint64) ([]*models.MediaItem[T], error)
	GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error)

	// Delete operation
	Delete(ctx context.Context, id uint64) error
}

type mediaItemRepository[T types.MediaData] struct {
	db *gorm.DB
}

// NewMediaItemRepository creates a new media item repository
func NewMediaItemRepository[T types.MediaData](db *gorm.DB) MediaItemRepository[T] {
	return &mediaItemRepository[T]{db: db}
}

func (r *mediaItemRepository[T]) Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to create %s media item: %w", item.Type, err)
	}
	return &item, nil
}

func (r *mediaItemRepository[T]) Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	// Get existing record first to check if it exists and preserve createdAt
	var existing models.MediaItem[T]

	if err := r.db.WithContext(ctx).Where("id = ?", item.ID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media item not found")
		}
		return nil, fmt.Errorf("failed to find media item: %w", err)
	}

	// Preserve createdAt
	item.CreatedAt = existing.CreatedAt

	// Update the record
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to update media item: %w", err)
	}
	return &item, nil
}

func (r *mediaItemRepository[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	var item models.MediaItem[T]

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media item not found")
		}
		return nil, fmt.Errorf("failed to get media item: %w", err)
	}

	return &item, nil
}

func (r *mediaItemRepository[T]) GetByExternalID(ctx context.Context, externalID string, clientID uint64) (*models.MediaItem[T], error) {
	var item models.MediaItem[T]

	if err := r.db.WithContext(ctx).
		Where("external_id = ? AND client_id = ?", externalID, clientID).
		First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media item not found")
		}
		return nil, fmt.Errorf("failed to get media item: %w", err)
	}

	return &item, nil
}

func (r *mediaItemRepository[T]) GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	if err := r.db.WithContext(ctx).
		Where("client_id = ?", clientID).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items: %w", err)
	}

	return items, nil
}

func (r *mediaItemRepository[T]) GetByType(ctx context.Context, mediaType types.MediaType, clientID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	if err := r.db.WithContext(ctx).
		Where("type = ? AND client_id = ?", mediaType, clientID).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by type: %w", err)
	}

	return items, nil
}

func (r *mediaItemRepository[T]) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	if err := r.db.WithContext(ctx).
		Joins("JOIN clients ON media_items.client_id = clients.id").
		Where("clients.user_id = ?", userID).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items for user: %w", err)
	}

	return items, nil
}

func (r *mediaItemRepository[T]) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.MediaItem[T]{})

	if err := result.Error; err != nil {
		return fmt.Errorf("failed to delete media item: %w", err)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("media item not found")
	}

	return nil
}
