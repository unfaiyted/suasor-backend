// Package repository provides data access layer implementations for the application.
package repository

// ClientMediaItemRepository represents the client-specific repository for media items.
// This focuses on operations related to media items that are linked to external clients
// such as Plex, Emby, Jellyfin, etc.
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
	"strings"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"

	"gorm.io/gorm"
)

// ClientMediaItemRepository defines the interface for client-associated media item operations
// This focuses specifically on media items that are linked to external clients like Plex, Emby, etc.
type ClientMediaItemRepository[T types.MediaData] interface {
	// Core CRUD operations
	Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
	Delete(ctx context.Context, id uint64) error

	// Client-specific operations
	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error)
	GetByClientItemID(ctx context.Context, externalID string, clientID uint64) (*models.MediaItem[T], error)
	GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error)
	GetByType(ctx context.Context, mediaType types.MediaType, clientID uint64) ([]*models.MediaItem[T], error)
	
	// User and client combined operations
	GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error)
	
	// Search operation
	Search(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[T], error)
	
	// Advanced operations
	GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[T], error)
	SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error
}

type clientMediaItemRepository[T types.MediaData] struct {
	db *gorm.DB
}

// NewClientMediaItemRepository creates a new repository for client-associated media items
func NewClientMediaItemRepository[T types.MediaData](db *gorm.DB) ClientMediaItemRepository[T] {
	return &clientMediaItemRepository[T]{db: db}
}

func (r *clientMediaItemRepository[T]) Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to create %s media item: %w", item.Type, err)
	}
	return &item, nil
}

func (r *clientMediaItemRepository[T]) Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
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

func (r *clientMediaItemRepository[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
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

func (r *clientMediaItemRepository[T]) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error) {
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

func (r *clientMediaItemRepository[T]) GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	// Use JSON contains operator to find items where clientIDs contains an entry with the given client ID
	query := r.db.WithContext(ctx).
		Where("sync_clients @> ?", fmt.Sprintf(`[{"id":%d}]`, clientID)).
		Find(&items)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("failed to get media items: %w", err)
	}

	return items, nil
}

func (r *clientMediaItemRepository[T]) GetByClientItemID(ctx context.Context, itemID string, clientID uint64) (*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	// Use JSON contains operator to find items where clientIDs contains an entry with the given ID and itemID
	query := r.db.WithContext(ctx).
		Where("sync_clients @> ?", fmt.Sprintf(`[{"id":%d,"itemId":"%s"}]`, clientID, itemID)).
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

func (r *clientMediaItemRepository[T]) GetByType(ctx context.Context, mediaType types.MediaType, clientID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]

	// Find all items of the given type that also have a reference to the client ID
	query := r.db.WithContext(ctx).
		Where("type = ? AND sync_clients @> ?", mediaType, fmt.Sprintf(`[{"id":%d}]`, clientID)).
		Find(&items)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by type: %w", err)
	}

	return items, nil
}

func (r *clientMediaItemRepository[T]) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Uint64("userID", userID).
		Msg("Getting media items by user ID")

	// Get all client IDs belonging to this user
	var clientIDs []uint64
	if err := r.db.WithContext(ctx).
		Table("clients").
		Where("user_id = ?", userID).
		Pluck("id", &clientIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get client IDs for user: %w", err)
	}

	// Build client-based conditions
	clientConditions := ""
	if len(clientIDs) > 0 {
		// Build JSON pattern conditions for each client ID
		clientPatterns := make([]string, len(clientIDs))
		for i, clientID := range clientIDs {
			clientPatterns[i] = fmt.Sprintf("sync_clients @> '[{\"id\":%d}]'", clientID)
		}
		clientConditions = "(" + strings.Join(clientPatterns, " OR ") + ")"
	}

	// Build the comprehensive query
	query := r.db.WithContext(ctx)

	// We always want to include user-owned playlists and collections
	ownershipCondition := fmt.Sprintf("("+
		"(type = '%s' AND (data->'ItemList'->>'Owner')::integer = %d) OR "+
		"(type = '%s' AND (data->'ItemList'->>'Owner')::integer = %d)"+
		")",
		types.MediaTypePlaylist, userID,
		types.MediaTypeCollection, userID)

	// Combine owner-based and client-based conditions
	if clientConditions != "" {
		query = query.Where(ownershipCondition + " OR " + clientConditions)
	} else {
		query = query.Where(ownershipCondition)
	}

	if err := query.Find(&items).Error; err != nil {
		log.Error().Err(err).Msg("Failed to get media items")
		return nil, fmt.Errorf("failed to get media items for user: %w", err)
	}

	log.Info().
		Int("count", len(items)).
		Msg("Media items retrieved successfully")

	return items, nil
}

func (r *clientMediaItemRepository[T]) Delete(ctx context.Context, id uint64) error {
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

func (r *clientMediaItemRepository[T]) Search(ctx context.Context, options types.QueryOptions) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(options.MediaType)).
		Str("query", options.Query).
		Uint64("clientID", options.ClientID).
		Msg("Searching client media items")

	var items []*models.MediaItem[T]

	// Build the query
	query := r.db.WithContext(ctx)
	
	// Add type filter if provided
	if options.MediaType != "" {
		query = query.Where("type = ?", options.MediaType)
	}
	
	// Add client filter if specified
	if options.ClientID > 0 {
		query = query.Where("sync_clients @> ?", fmt.Sprintf(`[{"id":%d}]`, options.ClientID))
	}
	
	// Add text search on title
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
	
	// Order by most recently created
	query = query.Order("created_at DESC")

	// Execute the query
	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to search media items: %w", err)
	}

	return items, nil
}

// GetByMultipleClients retrieves all media items associated with any of the specified clients
func (r *clientMediaItemRepository[T]) GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Msg(fmt.Sprintf("Getting media items by multiple clients: %v", clientIDs))

	if len(clientIDs) == 0 {
		return []*models.MediaItem[T]{}, nil
	}

	var items []*models.MediaItem[T]
	
	// Build a query with OR conditions for each client ID
	query := r.db.WithContext(ctx)
	
	// Create array of client pattern conditions
	clientPatterns := make([]string, len(clientIDs))
	for i, clientID := range clientIDs {
		clientPatterns[i] = fmt.Sprintf("sync_clients @> '[{\"id\":%d}]'", clientID)
	}
	combinedCondition := "(" + strings.Join(clientPatterns, " OR ") + ")"
	
	// Apply the combined condition
	query = query.Where(combinedCondition)
	
	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by multiple clients: %w", err)
	}

	return items, nil
}

// SyncItemBetweenClients creates or updates a mapping between a media item and a target client
func (r *clientMediaItemRepository[T]) SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("itemID", itemID).
		Uint64("sourceClientID", sourceClientID).
		Uint64("targetClientID", targetClientID).
		Str("targetItemID", targetItemID).
		Msg("Syncing item between clients")

	// Get the item to update
	var item models.MediaItem[T]
	if err := r.db.WithContext(ctx).Where("id = ?", itemID).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("media item not found")
		}
		return fmt.Errorf("failed to find media item for sync: %w", err)
	}
	
	// Check if this item already has the source client
	hasSourceClient := false
	for _, client := range item.SyncClients {
		if client.ID == sourceClientID {
			hasSourceClient = true
			break
		}
	}
	
	if !hasSourceClient {
		return fmt.Errorf("item is not associated with source client")
	}
	
	// Check if the target client is already in the SyncClients array
	for i, client := range item.SyncClients {
		if client.ID == targetClientID {
			// Update the existing client mapping
			item.SyncClients[i].ItemID = targetItemID
			// Update the item
			if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
				return fmt.Errorf("failed to update client mapping: %w", err)
			}
			return nil
		}
	}
	
	// If we get here, the target client isn't in the array yet, so add it
	item.SyncClients = append(item.SyncClients, models.ClientMapping{
		ID:     targetClientID,
		ItemID: targetItemID,
	})
	
	// Update the item
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return fmt.Errorf("failed to add new client mapping: %w", err)
	}
	
	return nil
}
