// Package repository provides data access layer implementations for the application.
package repository

// MediaItemRepository represents the base repository for media items.
// This focuses on generic operations that are not specifically tied to clients or users.
// It provides the core functionality for working with media items directly.
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
	"time"

	"gorm.io/gorm"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
)

// MediaItemRepository defines the interface for generic media item operations
// This focuses solely on the media items themselves without user or client associations
type MediaItemRepository[T types.MediaData] interface {
	// Core CRUD operations
	Create(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
	GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error)
	GetByClientItemID(ctx context.Context, clientItemID string, clientID uint64) (*models.MediaItem[T], error)
	GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error)
	Delete(ctx context.Context, id uint64) error

	// Batch operations
	GetMediaItemsByIDs(ctx context.Context, ids []uint64) ([]*models.MediaItem[T], error)
	GetMixedMediaItemsByIDs(ctx context.Context, ids []uint64) (*models.MediaItems, error)
	BatchCreate(ctx context.Context, items []*models.MediaItem[T]) ([]*models.MediaItem[T], error)
	BatchUpdate(ctx context.Context, items []*models.MediaItem[T]) ([]*models.MediaItem[T], error)

	// Query operations
	GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error)
	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error)
	Search(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error)

	// Specialized queries
	GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error)
	GetPopularItems(ctx context.Context, limit int) ([]*models.MediaItem[T], error)
	GetItemsByAttributes(ctx context.Context, attributes map[string]interface{}, limit int) ([]*models.MediaItem[T], error)
}

type mediaItemRepository[T types.MediaData] struct {
	db *gorm.DB
}

// NewMediaItemRepository creates a new media item repository
func NewMediaItemRepository[T types.MediaData](db *gorm.DB) MediaItemRepository[T] {
	return &mediaItemRepository[T]{db: db}
}

// Create adds a new media item to the database
func (r *mediaItemRepository[T]) Create(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(item.Type)).
		Msg("Creating media item")

	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to create media item: %w", err)
	}
	return item, nil
}

// Update modifies an existing media item
func (r *mediaItemRepository[T]) Update(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", item.ID).
		Str("type", string(item.Type)).
		Msg("Updating media item")

	// Get existing record first to check if it exists and preserve createdAt
	var existing models.MediaItem[T]
	if err := r.db.WithContext(ctx).Where("id = ?", item.ID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media item not found")
		}
		return nil, fmt.Errorf("failed to find media item: %w", err)
	}

	// Preserve createdAt timestamp
	item.CreatedAt = existing.CreatedAt

	// Update the record
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to update media item: %w", err)
	}
	return item, nil
}

// GetByID retrieves a media item by its ID
func (r *mediaItemRepository[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Getting media item by ID")

	var item models.MediaItem[T]
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media item not found")
		}
		return nil, fmt.Errorf("failed to get media item: %w", err)
	}
	return &item, nil
}

// Delete removes a media item
func (r *mediaItemRepository[T]) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting media item")

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

// GetMediaItemsByIDs retrieves multiple media items by their IDs
func (r *mediaItemRepository[T]) GetMediaItemsByIDs(ctx context.Context, ids []uint64) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("count", len(ids)).
		Msg("Getting media items by IDs")

	var items []*models.MediaItem[T]
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by IDs: %w", err)
	}

	return items, nil
}

// BatchCreate adds multiple media items to the database
func (r *mediaItemRepository[T]) BatchCreate(ctx context.Context, items []*models.MediaItem[T]) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
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
func (r *mediaItemRepository[T]) BatchUpdate(ctx context.Context, items []*models.MediaItem[T]) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
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

// GetByType retrieves all media items of a specific type
func (r *mediaItemRepository[T]) GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(mediaType)).
		Msg("Getting media items by type")

	var items []*models.MediaItem[T]
	if err := r.db.WithContext(ctx).
		Where("type = ?", mediaType).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by type: %w", err)
	}

	return items, nil
}

// GetByExternalID retrieves a media item by an external ID
func (r *mediaItemRepository[T]) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("source", source).
		Str("externalID", externalID).
		Msg("Getting media item by external ID")

	var items []*models.MediaItem[T]

	// Use JSON contains operator to find items where externalIDs contains an entry with the given source and ID
	query := r.db.WithContext(ctx).
		Where("external_ids @> ?", fmt.Sprintf(`[{"source":"%s","id":"%s"}]`, source, externalID)).
		Find(&items)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("failed to get media item by external ID: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("media item not found")
	}

	// Return the first match
	return items[0], nil
}

// Search finds media items based on a query string
func (r *mediaItemRepository[T]) Search(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Str("type", string(query.MediaType)).
		Int("limit", query.Limit).
		Int("offset", query.Offset).
		Msg("Searching media items")

	dbQuery := r.db.WithContext(ctx)

	// Add type filter if provided
	if query.MediaType != "" {
		dbQuery = dbQuery.Where("type = ?", query.MediaType)
	}

	// Add search condition
	if query.Query != "" {
		// Use ILIKE for case-insensitive search in PostgreSQL
		// TODOL: user paramater string
		dbQuery = dbQuery.Where("title ILIKE ?", "%"+query.Query+"%")
	}

	// Add pagination
	if query.Limit > 0 {
		dbQuery = dbQuery.Limit(query.Limit)
	}

	if query.Offset > 0 {
		dbQuery = dbQuery.Offset(query.Offset)
	}

	// Order by most recently created
	dbQuery = dbQuery.Order("created_at DESC")

	var items []*models.MediaItem[T]
	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to search media items: %w", err)
	}

	return items, nil
}

// GetRecentItems retrieves recently added items of a specific type
func (r *mediaItemRepository[T]) GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Str("type", string(mediaType)).
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recent media items")

	var items []*models.MediaItem[T]

	// Calculate the cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -days)

	dbQuery := r.db.WithContext(ctx).
		Where("type = ?", mediaType).
		Where("created_at >= ?", cutoffDate)

	// Add limit if provided
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	// Order by most recently created
	dbQuery = dbQuery.Order("created_at DESC")

	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent media items: %w", err)
	}

	return items, nil
}

// GetPopularItems retrieves popular items of a specific type
// Note: This implementation assumes a "play_count" or similar field in the data JSON
// You may need to adapt this based on your actual schema
func (r *mediaItemRepository[T]) GetPopularItems(ctx context.Context, limit int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Str("type", string(mediaType)).
		Int("limit", limit).
		Msg("Getting popular media items")

	var items []*models.MediaItem[T]

	dbQuery := r.db.WithContext(ctx).
		Where("type = ?", mediaType)

	// Add an order by play_count or a similar metric from the JSON data
	// This is PostgreSQL-specific JSON path syntax
	dbQuery = dbQuery.Order("(data->>'PlayCount')::int DESC NULLS LAST")

	// Add limit if provided
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get popular media items: %w", err)
	}

	return items, nil
}

// GetItemsByAttributes retrieves items matching specific attributes
func (r *mediaItemRepository[T]) GetItemsByAttributes(ctx context.Context, attributes map[string]interface{}, limit int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("attributeCount", len(attributes)).
		Int("limit", limit).
		Msg("Getting media items by attributes")

	dbQuery := r.db.WithContext(ctx)

	// Add filters for each attribute
	for key, value := range attributes {
		// For JSON attributes, use the PostgreSQL JSON operators
		if key == "genre" || key == "tags" || key == "categories" {
			// Use the @> operator for array containment
			dbQuery = dbQuery.Where(fmt.Sprintf("data->'%s' @> ?", key), fmt.Sprintf(`["%v"]`, value))
		} else if key == "year" || key == "runtime" || key == "rating" {
			// These are likely numeric fields
			dbQuery = dbQuery.Where(fmt.Sprintf("data->>'%s' = ?", key), fmt.Sprintf("%v", value))
		} else {
			// For other fields, use direct column matching if it's a column, or JSON path if it's in the data
			if key == "type" || key == "title" {
				dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
			} else {
				dbQuery = dbQuery.Where(fmt.Sprintf("data->>'%s' = ?", key), fmt.Sprintf("%v", value))
			}
		}
	}

	// Add limit if provided
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	var items []*models.MediaItem[T]
	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by attributes: %w", err)
	}

	return items, nil
}

func (r *mediaItemRepository[T]) GetMixedMediaItemsByIDs(ctx context.Context, ids []uint64) (*models.MediaItems, error) {
	// Fetch movies
	movies, err := fetchMediaItemsByType[*types.Movie](ctx, r.db, ids, types.MediaTypeMovie)
	if err != nil {
		return nil, err
	}
	series, err := fetchMediaItemsByType[*types.Series](ctx, r.db, ids, types.MediaTypeSeries)
	if err != nil {
		return nil, err
	}
	episodes, err := fetchMediaItemsByType[*types.Episode](ctx, r.db, ids, types.MediaTypeEpisode)
	if err != nil {
		return nil, err
	}
	seasons, err := fetchMediaItemsByType[*types.Season](ctx, r.db, ids, types.MediaTypeSeason)
	if err != nil {
		return nil, err
	}
	tracks, err := fetchMediaItemsByType[*types.Track](ctx, r.db, ids, types.MediaTypeTrack)
	if err != nil {
		return nil, err
	}
	albums, err := fetchMediaItemsByType[*types.Album](ctx, r.db, ids, types.MediaTypeAlbum)
	if err != nil {
		return nil, err
	}
	artists, err := fetchMediaItemsByType[*types.Artist](ctx, r.db, ids, types.MediaTypeArtist)
	if err != nil {
		return nil, err
	}
	playlists, err := fetchMediaItemsByType[*types.Playlist](ctx, r.db, ids, types.MediaTypePlaylist)
	if err != nil {
		return nil, err
	}
	collections, err := fetchMediaItemsByType[*types.Collection](ctx, r.db, ids, types.MediaTypeCollection)
	if err != nil {
		return nil, err
	}

	return &models.MediaItems{
		Movies:      movies,
		Series:      series,
		Seasons:     seasons,
		Episodes:    episodes,
		Albums:      albums,
		Artists:     artists,
		Tracks:      tracks,
		Playlists:   playlists,
		Collections: collections,
	}, nil
}

func (r *mediaItemRepository[T]) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting all media items")

	var items []*models.MediaItem[T]

	dbQuery := r.db.WithContext(ctx)

	// Add limit if provided
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	// Add offset if provided
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}

	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get all media items: %w", err)
	}

	return items, nil
}

func (r *mediaItemRepository[T]) GetByClientItemID(ctx context.Context, clientItemID string, clientID uint64) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("clientItemID", clientItemID).
		Uint64("clientID", clientID).
		Msg("Getting media item by client item ID")

	var item *models.MediaItem[T]
	if err := r.db.WithContext(ctx).
		Where("client_id = ?", clientID).
		Where("client_item_id = ?", clientItemID).
		First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media item not found")
		}
		return nil, fmt.Errorf("failed to get media item by client item ID: %w", err)
	}

	return item, nil
}

func (r *mediaItemRepository[T]) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting media items by user ID")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	if mediaType == types.MediaTypePlaylist || mediaType == types.MediaTypeCollection {

		// Should for now be limited to user-owned playlists and collections
		query := r.db.WithContext(ctx).
			Where("type IN (?) AND data->'ItemList'->>'Owner' = ?", mediaType, userID)

		if err := query.Find(&items).Error; err != nil {
			log.Error().Err(err).Msg("Failed to get media items")
			return nil, fmt.Errorf("failed to get media items for user: %w", err)
		}

		log.Info().
			Int("count", len(items)).
			Msg("Media items retrieved successfully")

		return items, nil
	}
	return nil, fmt.Errorf("media type not supported")

}
