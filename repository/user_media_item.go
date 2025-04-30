// Package repository provides data access layer implementations for the application.
package repository

// UserMediaItemRepository represents the user-specific repository for media items.
// This focuses on operations related to media items that are directly owned by users,
// such as playlists and collections. These items are created and managed within the app
// rather than synchronized from external clients.
//
// Relationships with other repositories:
// - MediaItemRepository: Core operations on media items without client or user associations
// - ClientMediaItemRepository: Operations for media items linked to specific clients
// - UserMediaItemRepository: Operations for media items owned by users (playlists, collections)
//
// This three-tier approach allows for clear separation of concerns while maintaining
// a single database table for all media items.

import (
	"context"
	"encoding/json"
	"fmt"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"

	"gorm.io/gorm"
)

// UserMediaItemRepository defines the interface for user-owned media item operations
// This specifically focuses on playlists, collections, and other user-owned media
// (as opposed to media from external clients)
type UserMediaItemRepository[T types.MediaData] interface {
	CoreMediaItemRepository[T]

	Create(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error)

	Delete(ctx context.Context, id uint64) error

	// User-specific operations
	GetUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error)

	// General retrieval operations
	GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error)
	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error)
	GetByExternalIDs(ctx context.Context, externalIDs types.ExternalIDs) (*models.MediaItem[T], error)
	Search(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[T], error)

	BatchCreate(ctx context.Context, items []*models.MediaItem[T]) ([]*models.MediaItem[T], error)
	BatchUpdate(ctx context.Context, items []*models.MediaItem[T]) ([]*models.MediaItem[T], error)
}

type userMediaItemRepository[T types.MediaData] struct {
	CoreMediaItemRepository[T]
	db *gorm.DB
}

// NewUserMediaItemRepository creates a new repository for user-owned media items
func NewUserMediaItemRepository[T types.MediaData](
	db *gorm.DB,
	itemRepo CoreMediaItemRepository[T],
) UserMediaItemRepository[T] {
	return &userMediaItemRepository[T]{
		CoreMediaItemRepository: itemRepo,
		db:                      db}
}

// Create adds a new user-owned media item to the database
func (r *userMediaItemRepository[T]) Create(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error) {

	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(item.Type)).
		Msg("Creating media item")

	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to create media item: %w", err)
	}
	return item, nil
}

// Update modifies an existing user-owned media item
func (r *userMediaItemRepository[T]) Update(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error) {
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

	// Convert SyncClients to JSON manually
	syncClientsJSON, err := json.Marshal(item.SyncClients)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sync clients: %w", err)
	}

	// Use updates map to ensure proper JSON formatting
	updates := map[string]interface{}{
		"updated_at":   time.Now(),
		"sync_clients": json.RawMessage(syncClientsJSON),
		"external_ids": item.ExternalIDs,
		"is_public":    item.IsPublic,
		"type":         item.Type,
		"title":        item.Title,
		"release_date": item.ReleaseDate,
		"release_year": item.ReleaseYear,
		"stream_url":   item.StreamURL,
		"download_url": item.DownloadURL,
		"data":         item.Data,
	}

	// Update the record using a map to ensure proper JSON handling
	if err := r.db.WithContext(ctx).Model(&item).Where("id = ?", item.ID).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update media item: %w", err)
	}

	return item, nil
}

// Delete removes a user-owned media item
func (r *userMediaItemRepository[T]) Delete(ctx context.Context, id uint64) error {
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

// BatchCreate adds multiple media items to the database
func (r *userMediaItemRepository[T]) BatchCreate(ctx context.Context, items []*models.MediaItem[T]) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("count", len(items)).
		Msg("Batch creating media items")

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Create a slice to hold the created items with pointers
	createdItems := make([]*models.MediaItem[T], 0, len(items))

	// Insert each item within the transaction
	for i := range items {
		if err := tx.Create(&items[i]).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create media item: %w", err)
		}
		createdItems = append(createdItems, items[i])
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdItems, nil
}

// BatchUpdate modifies multiple media items
func (r *userMediaItemRepository[T]) BatchUpdate(ctx context.Context, items []*models.MediaItem[T]) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("count", len(items)).
		Msg("Batch updating media items")

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Update each item within the transaction
	updatedItems := make([]*models.MediaItem[T], 0, len(items))
	for i := range items {
		// Get the existing item to preserve createdAt
		var existing models.MediaItem[T]
		if err := tx.Where("id = ?", items[i].ID).First(&existing).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("media item with ID %d not found", items[i].ID)
			}
			return nil, fmt.Errorf("failed to find media item: %w", err)
		}

		// Preserve createdAt
		items[i].CreatedAt = existing.CreatedAt

		// Update the item
		if err := tx.Save(&items[i]).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update media item: %w", err)
		}

		updatedItems = append(updatedItems, items[i])
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedItems, nil
}

// GetUserContent retrieves all types of user-owned content (playlists and collections)
func (r *userMediaItemRepository[T]) GetUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting all user-owned content")

	var items []*models.MediaItem[T]

	// Query for all user-owned content types
	query := r.db.WithContext(ctx).Where(
		"(type = ? OR type = ?) AND (data->'ItemList'->>'Owner')::integer = ?",
		types.MediaTypePlaylist,
		types.MediaTypeCollection,
		userID,
	)

	// Add limit if provided
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Order by most recently updated
	query = query.Order("updated_at DESC")

	if err := query.Find(&items).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get user-owned content")
		return nil, fmt.Errorf("failed to get user content: %w", err)
	}

	log.Info().
		Int("count", len(items)).
		Uint64("userID", userID).
		Msg("User-owned content retrieved successfully")

	return items, nil
}

// GetByExternalIDs retrieves a media item by any of the provided external IDs
func (r *userMediaItemRepository[T]) GetByExternalIDs(ctx context.Context, externalIDs types.ExternalIDs) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Msg("Retrieving media item by external IDs")

	var items []*models.MediaItem[T]

	if len(externalIDs) == 0 {
		return nil, fmt.Errorf("no external IDs provided")
	}

	// Start building the query
	db := r.db.WithContext(ctx)

	// For the first external ID, use Where; for subsequent IDs, use Or
	for i, externalID := range externalIDs {
		jsonPattern := fmt.Sprintf(`[{"source":"%s","id":"%s"}]`, externalID.Source, externalID.ID)

		if i == 0 {
			db = db.Where("external_ids @> ?", jsonPattern)
		} else {
			db = db.Or("external_ids @> ?", jsonPattern)
		}
	}

	// Execute the query
	if err := db.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by external IDs: %w", err)
	}

	log.Info().
		Int("count", len(items)).
		Msg("Media items retrieved successfully")

	if len(items) == 0 {
		return nil, fmt.Errorf("no media item found matching external IDs")
	}

	return items[0], nil
}
