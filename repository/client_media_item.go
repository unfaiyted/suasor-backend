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
	"suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"gorm.io/gorm"
)

// ClientMediaItemRepository defines the interface for client-associated media item operations
// This focuses specifically on media items that are linked to external clients like Plex, Emby, etc.
type ClientMediaItemRepository[T types.MediaData] interface {
	MediaItemRepository[T]

	// Client-specific operations
	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error)
	GetByClientItemID(ctx context.Context, externalID string, clientID uint64) (*models.MediaItem[T], error)
	GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error)

	// Advanced operations
	GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[T], error)
	SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error

	DeleteClientItem(ctx context.Context, clientID uint64, itemID string) error
}

type clientMediaItemRepository[T types.MediaData] struct {
	MediaItemRepository[T]
	db *gorm.DB
}

// NewClientMediaItemRepository creates a new repository for client-associated media items
func NewClientMediaItemRepository[T types.MediaData](
	db *gorm.DB,
	mediaItemRepository MediaItemRepository[T],
) ClientMediaItemRepository[T] {
	return &clientMediaItemRepository[T]{
		MediaItemRepository: mediaItemRepository,
		db:                  db}
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

func (r *clientMediaItemRepository[T]) GetByClientUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]
	log := logger.LoggerFromContext(ctx)

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

// GetByMultipleClients retrieves all media items associated with any of the specified clients
func (r *clientMediaItemRepository[T]) GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
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
	log := logger.LoggerFromContext(ctx)
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
	targetSyncClient := r.getSyncClientByClientID(ctx, targetClientID, targetItemID)

	// If we get here, the target client isn't in the array yet, so add it
	item.SyncClients = append(item.SyncClients, *targetSyncClient)

	// Update the item
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return fmt.Errorf("failed to add new client mapping: %w", err)
	}

	return nil
}

func (s *clientMediaItemRepository[T]) getSyncClientByClientID(ctx context.Context, clientID uint64, targetItemID string) *models.SyncClient {
	var clientType string

	// Raw SQL approach
	row := s.db.WithContext(ctx).Raw("SELECT type FROM clients WHERE id = ?", clientID).Row()
	if err := row.Scan(&clientType); err != nil {
		return nil
	}

	return &models.SyncClient{
		ID:     clientID,
		ItemID: targetItemID,
		Type:   clienttypes.ClientType(clientType),
	}
}

func (s *clientMediaItemRepository[T]) DeleteClientItem(ctx context.Context, clientID uint64, itemID string) error {
	// TODO: Implement this, maybe should remove the records from SyncClients
	return nil

}
