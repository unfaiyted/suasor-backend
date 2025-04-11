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

	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error)
	GetByClientItemID(ctx context.Context, externalID string, clientID uint64) (*models.MediaItem[T], error)
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

func (r *mediaItemRepository[T]) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	// Use JSON contains operator to find items where externalIDs contains an entry with the given source and ID
	query := r.db.WithContext(ctx).
		Where("external_ids @> ?", fmt.Sprintf(`[{"source":"%s","id":"%s"}]`, source, externalID)).
		Find(&items)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("failed to get media item: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("media item not found")
	}

	// Return the first match
	return items[0], nil
}

func (r *mediaItemRepository[T]) GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	// Use JSON contains operator to find items where clientIDs contains an entry with the given client ID
	query := r.db.WithContext(ctx).
		Where("client_ids @> ?", fmt.Sprintf(`[{"id":%d}]`, clientID)).
		Find(&items)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("failed to get media items: %w", err)
	}

	return items, nil
}

// GetByClientItemID retrieves a media item by client ID and client item ID
func (r *mediaItemRepository[T]) GetByClientItemID(ctx context.Context, itemID string, clientID uint64) (*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	// Use JSON contains operator to find items where clientIDs contains an entry with the given ID and itemID
	query := r.db.WithContext(ctx).
		Where("client_ids @> ?", fmt.Sprintf(`[{"id":%d,"itemId":"%s"}]`, clientID, itemID)).
		Find(&items)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("failed to get media item: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("not found")
	}

	// Return the first match
	return items[0], nil
}

func (r *mediaItemRepository[T]) GetByType(ctx context.Context, mediaType types.MediaType, clientID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	// Find all items of the given type that also have a reference to the client ID
	query := r.db.WithContext(ctx).
		Where("type = ? AND client_ids @> ?", mediaType, fmt.Sprintf(`[{"id":%d}]`, clientID)).
		Find(&items)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by type: %w", err)
	}

	return items, nil
}

func (r *mediaItemRepository[T]) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]
	var clientIDs []uint64

	// First get all client IDs belonging to this user
	if err := r.db.WithContext(ctx).
		Table("clients").
		Where("user_id = ?", userID).
		Pluck("id", &clientIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get client IDs for user: %w", err)
	}

	if len(clientIDs) == 0 {
		// No clients for this user
		return items, nil
	}

	// Build a condition to match any client ID in the user's clients
	// This is more complex with JSON fields, so we'll use a different approach
	var conditions []string
	var args []interface{}

	for _, clientID := range clientIDs {
		jsonPattern := fmt.Sprintf(`[{"id":%d}]`, clientID)
		conditions = append(conditions, "client_ids @> ?")
		args = append(args, jsonPattern)
	}

	// Combine conditions with OR
	query := r.db.WithContext(ctx).Where(conditions[0], args[0])
	for i := 1; i < len(conditions); i++ {
		query = query.Or(conditions[i], args[i])
	}

	if err := query.Find(&items).Error; err != nil {
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
