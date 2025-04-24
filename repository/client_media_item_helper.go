package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/types/models"

	"gorm.io/gorm"
)

// ClientMediaItemHelper provides helper methods for working with client media items
type ClientMediaItemHelper struct {
	db *gorm.DB
}

// NewClientMediaItemHelper creates a new client media item helper
func NewClientMediaItemHelper(db *gorm.DB) *ClientMediaItemHelper {
	return &ClientMediaItemHelper{db: db}
}

// GetMediaItemByClientID retrieves a media item by client ID
// This handles the complex case of mapping client-specific IDs stored in the SyncClients JSONB field
func (h *ClientMediaItemHelper) GetMediaItemByClientID(ctx context.Context, clientID uint64, clientItemID string) (uint64, error) {
	var mediaItemID uint64

	// Query to find a media item with the given client ID in its SyncClients
	// This uses a JSONB query to search within the SyncClients array
	result := h.db.WithContext(ctx).
		Table("media_items").
		Select("id").
		Where("sync_clients @> ?::jsonb",
			fmt.Sprintf(`[{"id": %d, "itemID": "%s"}]`, clientID, clientItemID)).
		First(&mediaItemID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no media item found for client ID %d and item ID %s", clientID, clientItemID)
		}
		return 0, fmt.Errorf("failed to find media item by client ID: %w", result.Error)
	}

	return mediaItemID, nil
}

// GetMediaItemByClientIDAndType retrieves a media item by client ID and media type
func (h *ClientMediaItemHelper) GetMediaItemByClientIDAndType(ctx context.Context, clientID uint64, clientItemID string, mediaType types.MediaType) (uint64, error) {
	var mediaItemID uint64

	// Query to find a media item with the given client ID in its SyncClients and the specified type
	result := h.db.WithContext(ctx).
		Table("media_items").
		Select("id").
		Where("sync_clients @> ?::jsonb AND type = ?",
			fmt.Sprintf(`[{"id": %d, "itemID": "%s"}]`, clientID, clientItemID),
			mediaType).
		First(&mediaItemID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no %s found for client ID %d and item ID %s", mediaType, clientID, clientItemID)
		}
		return 0, fmt.Errorf("failed to find %s by client ID: %w", mediaType, result.Error)
	}

	return mediaItemID, nil
}

// GetMediaItemsByClientIDs retrieves multiple media items by client IDs
// Useful for batch operations
func (h *ClientMediaItemHelper) GetMediaItemsByClientIDs(ctx context.Context, clientID uint64, clientItemIDs []string) (map[string]uint64, error) {
	results := make(map[string]uint64)

	// For each client item ID, we need to search the SyncClients array
	// This is not ideal for performance but necessary for batch operations
	for _, itemID := range clientItemIDs {
		id, err := h.GetMediaItemByClientID(ctx, clientID, itemID)
		if err == nil {
			results[itemID] = id
		}
		// If there's an error, we just skip this item
	}

	return results, nil
}

// GetOrCreateMediaItemMapping ensures a mapping exists between a client item and an internal media item
// If no mapping exists, it will create a new media item
func (h *ClientMediaItemHelper) GetOrCreateMediaItemMapping(
	ctx context.Context,
	clientID uint64,
	clientType clienttypes.ClientType,
	clientItemID string,
	mediaType types.MediaType,
	title string,
	data types.MediaData) (uint64, error) {

	// Try to find existing mapping
	mediaItemID, err := h.GetMediaItemByClientID(ctx, clientID, clientItemID)
	if err == nil {
		// Found an existing mapping
		return mediaItemID, nil
	}

	// No mapping found, create a new media item
	syncClient := models.SyncClient{
		ID:     clientID,
		Type:   clientType,
		ItemID: clientItemID,
	}

	// Create a new media item with this sync client
	mediaItem := models.MediaItem[types.MediaData]{
		Type:        mediaType,
		Title:       title,
		SyncClients: []models.SyncClient{syncClient},
		ExternalIDs: []models.ExternalID{},
		Data:        data,
	}

	// Insert the new media item
	result := h.db.WithContext(ctx).
		Table("media_items").
		Create(&mediaItem)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to create new media item: %w", result.Error)
	}

	return mediaItem.ID, nil
}

// SyncClientsList retrieves all media items with client IDs for a specific client
// Useful for synchronization tasks
func (h *ClientMediaItemHelper) SyncClientsList(ctx context.Context, clientID uint64) (map[string]uint64, error) {
	type Result struct {
		ID          uint64
		SyncClients json.RawMessage
	}

	var results []Result

	// Query all media items that have the specified client ID in their SyncClients
	err := h.db.WithContext(ctx).
		Table("media_items").
		Select("id, sync_clients").
		Where("sync_clients @> ?::jsonb",
			fmt.Sprintf(`[{"id": %d}]`, clientID)).
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query media items with client IDs: %w", err)
	}

	// Map client item IDs to internal IDs
	mapping := make(map[string]uint64)

	for _, result := range results {
		var syncClients []models.SyncClient
		if err := json.Unmarshal(result.SyncClients, &syncClients); err != nil {
			continue
		}

		// Find the client ID for this specific client
		for _, client := range syncClients {
			if client.ID == clientID {
				mapping[client.ItemID] = result.ID
				break
			}
		}
	}

	return mapping, nil
}
