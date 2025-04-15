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
	"fmt"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"

	"gorm.io/gorm"
)

// UserMediaItemRepository defines the interface for user-owned media item operations
// This specifically focuses on playlists, collections, and other user-owned media
// (as opposed to media from external clients)
type UserMediaItemRepository[T types.MediaData] interface {
	// Basic CRUD operations
	Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
	Delete(ctx context.Context, id uint64) error

	// User-specific operations
	GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error)
	GetUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error)
	GetRecentItems(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error)

	// General retrieval operations
	GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error)
	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error)
	Search(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[T], error)
}

type userMediaItemRepository[T types.MediaData] struct {
	db *gorm.DB
}

// NewUserMediaItemRepository creates a new repository for user-owned media items
func NewUserMediaItemRepository[T types.MediaData](db *gorm.DB) UserMediaItemRepository[T] {
	return &userMediaItemRepository[T]{db: db}
}

// Create adds a new user-owned media item to the database
func (r *userMediaItemRepository[T]) Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to create %s media item: %w", item.Type, err)
	}
	return &item, nil
}

// Update modifies an existing user-owned media item
func (r *userMediaItemRepository[T]) Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
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

// GetByID retrieves a user-owned media item by its ID
func (r *userMediaItemRepository[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
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

// GetByExternalID retrieves a user-owned media item by an external ID
func (r *userMediaItemRepository[T]) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error) {
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

// GetByUserID retrieves all user-owned media items for a specific user
// This is a key method specifically for playlists and collections
func (r *userMediaItemRepository[T]) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Uint64("userID", userID).
		Msg("Getting user-owned media items")

	// Query for items where the Owner field in the JSON data matches the user ID
	// This covers both playlists and collections
	query := r.db.WithContext(ctx).Where(
		"(type = ? OR type = ?) AND (data->'ItemList'->>'Owner')::integer = ?",
		types.MediaTypePlaylist,
		types.MediaTypeCollection,
		userID,
	)

	if err := query.Find(&items).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get user-owned media items")
		return nil, fmt.Errorf("failed to get user media items: %w", err)
	}

	log.Info().
		Int("count", len(items)).
		Msg("User-owned media items retrieved successfully")

	return items, nil
}

// GetByType retrieves all user-owned media items of a specific type
func (r *userMediaItemRepository[T]) GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	// Find all items of the given type that are user-owned
	query := r.db.WithContext(ctx).
		Where("type = ?", mediaType)

	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by type: %w", err)
	}

	return items, nil
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

// Search finds user-owned media items based on query options
func (r *userMediaItemRepository[T]) Search(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(options.MediaType)).
		Str("query", options.Query).
		Uint64("ownerID", options.OwnerID).
		Msg("Searching user-owned media items")

	var items []*models.MediaItem[T]

	// Build the query for user-owned content
	query := r.db.WithContext(ctx)

	// Add type filter if provided
	if options.MediaType != "" {
		query = query.Where("type = ?", options.MediaType)
	} else {
		// Default to searching for user-owned content types
		query = query.Where("type IN (?, ?)",
			types.MediaTypePlaylist,
			types.MediaTypeCollection)
	}

	// If we're searching for playlists or collections, add owner restriction
	if options.MediaType == types.MediaTypePlaylist ||
		options.MediaType == types.MediaTypeCollection ||
		options.MediaType == "" {
		// Only include the owner condition if user ID is provided
		if options.OwnerID > 0 {
			query = query.Where("(data->'ItemList'->>'Owner')::integer = ?", options.OwnerID)
		}
	}

	// Add search criteria
	if options.Query != "" {
		query = query.Where("title ILIKE ?", "%"+options.Query+"%")
	}

	// Add pagination
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	// Order by most recently modified for user content
	query = query.Order("updated_at DESC")

	// Execute the query
	if err := query.Find(&items).Error; err != nil {
		log.Error().Err(err).Msg("Failed to search user-owned media items")
		return nil, fmt.Errorf("failed to search media items: %w", err)
	}

	log.Info().
		Int("count", len(items)).
		Str("type", string(options.MediaType)).
		Msg("User-owned media items found")

	return items, nil
}

// GetUserContent retrieves all types of user-owned content (playlists and collections)
// in a single query and returns them sorted by most recently updated
func (r *userMediaItemRepository[T]) GetUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
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

func (r *userMediaItemRepository[T]) GetRecentItems(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting recent user-owned content")

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
		log.Error().Err(err).Msg("Failed to get recent user-owned content")
		return nil, fmt.Errorf("failed to get recent user content: %w", err)
	}

	log.Info().
		Int("count", len(items)).
		Uint64("userID", userID).
		Msg("Recent user-owned content retrieved successfully")

	return items, nil
}
