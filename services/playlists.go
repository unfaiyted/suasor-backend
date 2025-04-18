package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// PlaylistService manages application-specific playlist operations beyond the basic CRUD operations
type PlaylistService interface {
	// Base operations (leveraging UserMediaItemService)
	Create(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error)
	Update(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Playlist], error)
	GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[*mediatypes.Playlist], error)
	Delete(ctx context.Context, id uint64) error

	// Playlist-specific operations
	GetItems(ctx context.Context, playlistID uint64) (*models.MediaItems, error)
	AddItem(ctx context.Context, playlistID uint64, itemID uint64) error
	RemoveItem(ctx context.Context, playlistID uint64, itemID uint64) error
	ReorderItems(ctx context.Context, playlistID uint64, itemIDs []string) error
	UpdateItems(ctx context.Context, playlistID uint64, items []*models.MediaItem[mediatypes.MediaData]) error
	Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Playlist], error)
	GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Playlist], error)
	Sync(ctx context.Context, playlistID uint64, targetClientIDs []uint64) error
}

type playlistService struct {
	repo repository.MediaItemRepository[*mediatypes.Playlist] // For fetching playlist items
}

// NewPlaylistService creates a new playlist service
func NewPlaylistService(
	coreMediaRepo repository.MediaItemRepository[*mediatypes.Playlist],
) PlaylistService {
	return &playlistService{
		repo: coreMediaRepo,
	}
}

// Base operations (delegating to UserMediaItemService where appropriate)

func (s *playlistService) Create(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("title", playlist.Title).
		Msg("Creating playlist")

	// Ensure playlist-specific validation
	if playlist.Type != mediatypes.MediaTypePlaylist {
		playlist.Type = mediatypes.MediaTypePlaylist
	}

	// Ensure playlist has a valid name
	if playlist.Data == nil || playlist.Data.ItemList.Details.Title == "" {
		return nil, errors.New("playlist must have a title")
	}

	// Initialize items array if nil
	if playlist.Data.ItemList.Items == nil {
		playlist.Data.ItemList.Items = []mediatypes.ListItem{}
	}

	// Initialize sync client states if nil
	if playlist.Data.ItemList.SyncClientStates == nil {
		playlist.Data.ItemList.SyncClientStates = mediatypes.SyncClientStates{}
	}

	// Set creation time for LastModified
	if playlist.Data.ItemList.LastModified.IsZero() {
		playlist.Data.ItemList.LastModified = time.Now()
	}

	// Initialize ItemCount
	playlist.Data.ItemList.ItemCount = len(playlist.Data.ItemList.Items)

	// Set owner if not set
	if playlist.Data.ItemList.OwnerID == 0 && playlist.Data.OwnerID != 0 {
		playlist.Data.ItemList.OwnerID = playlist.Data.OwnerID
	}

	// Set title at MediaItem level to match the Data.ItemList.Details.Title
	if playlist.Title == "" && playlist.Data.ItemList.Details.Title != "" {
		playlist.Title = playlist.Data.ItemList.Details.Title
	}

	// Use the underlying repository directly for better control over validation
	result, err := s.repo.Create(ctx, playlist)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create playlist")
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	log.Info().
		Uint64("id", result.ID).
		Str("title", result.Title).
		Msg("Playlist created successfully")

	return result, nil
}

func (s *playlistService) Update(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", playlist.ID).
		Str("title", playlist.Title).
		Msg("Updating playlist")

	// Ensure the playlist exists
	existing, err := s.GetByID(ctx, playlist.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update playlist: %w", err)
	}

	// Preserve items if not provided in the update
	if playlist.Data.ItemList.Items == nil || len(playlist.Data.ItemList.Items) == 0 {
		playlist.Data.ItemList.Items = existing.Data.ItemList.Items
	}

	// Preserve sync client states if not provided
	if playlist.Data.ItemList.SyncClientStates == nil || len(playlist.Data.ItemList.SyncClientStates) == 0 {
		playlist.Data.ItemList.SyncClientStates = existing.Data.ItemList.SyncClientStates
	}

	// Update last modified time
	playlist.Data.ItemList.LastModified = time.Now()

	// Update ItemCount
	playlist.Data.ItemList.ItemCount = len(playlist.Data.ItemList.Items)

	// Ensure positions are normalized
	playlist.Data.NormalizePositions()

	// Set title at MediaItem level to match the Data.ItemList.Details.Title
	if playlist.Title != playlist.Data.ItemList.Details.Title {
		playlist.Title = playlist.Data.ItemList.Details.Title
	}

	// Run validation to check for issues
	issues := playlist.Data.ValidateItems()
	if len(issues) > 0 {
		// Log the issues but continue with the update
		for _, issue := range issues {
			log.Warn().Str("issue", issue).Msg("Playlist validation issue")
		}
	}

	// Update using the user service for consistent behavior
	result, err := s.repo.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", playlist.ID).
			Msg("Failed to update playlist")
		return nil, fmt.Errorf("failed to update playlist: %w", err)
	}

	log.Info().
		Uint64("id", result.ID).
		Str("title", result.Title).
		Msg("Playlist updated successfully")

	return result, nil
}

func (s *playlistService) GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Getting playlist by ID")

	// Use the user service
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to get playlist")
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	// Verify this is actually a playlist
	if result.Type != mediatypes.MediaTypePlaylist {
		log.Error().
			Uint64("id", id).
			Str("actualType", string(result.Type)).
			Msg("Item is not a playlist")
		return nil, fmt.Errorf("item with ID %d is not a playlist", id)
	}

	return result, nil
}

func (s *playlistService) GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting playlists by user ID")
	// Use the user service
	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

// Playlist-specific operations
func (s *playlistService) GetItems(ctx context.Context, playlistID uint64) (*models.MediaItems, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Msg("Getting playlist items")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return &models.MediaItems{}, fmt.Errorf("failed to get playlist items: %w", err)
	}

	// Return empty array if the playlist has no items
	if len(playlist.Data.ItemList.Items) == 0 {
		return &models.MediaItems{}, nil
	}

	// Extract item IDs for batch retrieval
	itemIDs := make([]uint64, len(playlist.Data.ItemList.Items))
	for i, item := range playlist.Data.ItemList.Items {
		itemIDs[i] = item.ItemID
	}

	// Fetch the actual media items using the core media repository
	actualItems, err := s.repo.GetMixedMediaItemsByIDs(ctx, itemIDs)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Msg("Failed to fetch actual playlist items")
		return nil, fmt.Errorf("failed to get playlist items: %w", err)
	}

	return actualItems, nil

}

func (s *playlistService) AddItem(ctx context.Context, playlistID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Uint64("itemID", itemID).
		Msg("Adding item to playlist")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to add item to playlist: %w", err)
	}

	// Create a ListItem from the media item
	newItem := mediatypes.ListItem{
		ItemID:        itemID,
		Position:      len(playlist.Data.ItemList.Items),
		LastChanged:   time.Now(),
		ChangeHistory: []mediatypes.ChangeRecord{},
	}

	// Add the item using the built-in AddItem method
	// 0 indicates application level modification
	playlist.Data.AddItem(newItem)

	// Store the update
	_, err = s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Uint64("itemID", itemID).
			Msg("Failed to add item to playlist")
		return fmt.Errorf("failed to update playlist after adding item: %w", err)
	}

	log.Info().
		Uint64("playlistID", playlistID).
		Uint64("itemID", itemID).
		Msg("Item added to playlist successfully")

	return nil
}

func (s *playlistService) RemoveItem(ctx context.Context, playlistID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Uint64("itemID", itemID).
		Msg("Removing item from playlist")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to remove item from playlist: %w", err)
	}

	// Use the RemoveItem method provided by ItemList
	// 0 indicates application level modification
	err = playlist.Data.RemoveItem(itemID, 0)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Uint64("itemID", itemID).
			Msg("Failed to remove item from playlist")
		return err
	}

	// Store the update
	_, err = s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Uint64("itemID", itemID).
			Msg("Failed to update playlist after removing item")
		return fmt.Errorf("failed to update playlist after removing item: %w", err)
	}

	log.Info().
		Uint64("playlistID", playlistID).
		Uint64("itemID", itemID).
		Msg("Item removed from playlist successfully")

	return nil
}

func (s *playlistService) ReorderItems(ctx context.Context, playlistID uint64, itemIDs []string) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Interface("itemIDs", itemIDs).
		Msg("Reordering playlist items")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to reorder playlist items: %w", err)
	}

	// Verify that the number of items matches
	if len(itemIDs) != len(playlist.Data.ItemList.Items) {
		log.Error().
			Int("providedCount", len(itemIDs)).
			Int("actualCount", len(playlist.Data.ItemList.Items)).
			Msg("Reorder operation must include all playlist items")
		return errors.New("reorder operation must include all playlist items")
	}

	// Create a new ordered list of items
	newOrder := make([]mediatypes.ListItem, len(itemIDs))
	tempItems := make(map[uint64]mediatypes.ListItem)

	// Create a map of existing items for quick lookup
	for _, item := range playlist.Data.ItemList.Items {
		tempItems[item.ItemID] = item
	}

	// First verify all items exist
	missingItems := []string{}
	for _, idStr := range itemIDs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Error().
				Str("itemID", idStr).
				Msg("Invalid item ID format")
			return fmt.Errorf("invalid item ID format: %s", idStr)
		}

		if _, exists := tempItems[id]; !exists {
			missingItems = append(missingItems, idStr)
		}
	}

	if len(missingItems) > 0 {
		log.Error().
			Interface("missingItems", missingItems).
			Msg("Items not found in playlist")
		return fmt.Errorf("items not found in playlist: %v", missingItems)
	}

	// Now build the new order
	for i, idStr := range itemIDs {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		item := tempItems[id]

		// Update position
		item.Position = i

		// Add change record
		item.AddChangeRecord(0, "reorder") // 0 indicates application level change

		newOrder[i] = item
	}

	// Update the playlist with the new item order
	playlist.Data.ItemList.Items = newOrder
	playlist.Data.ItemList.LastModified = time.Now()
	playlist.Data.ItemList.ModifiedBy = 0 // 0 indicates application level modification

	// Normalize positions to ensure they're sequential
	playlist.Data.NormalizePositions()

	_, err = s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Msg("Failed to update playlist after reordering items")
		return fmt.Errorf("failed to update playlist after reordering: %w", err)
	}

	log.Info().
		Uint64("playlistID", playlistID).
		Msg("Playlist items reordered successfully")

	return nil
}

func (s *playlistService) UpdateItems(ctx context.Context, playlistID uint64, items []*models.MediaItem[mediatypes.MediaData]) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Int("itemCount", len(items)).
		Msg("Updating playlist items")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to update playlist items: %w", err)
	}

	// Convert MediaItems to ListItems
	listItems := make([]mediatypes.ListItem, len(items))
	now := time.Now()

	for i, item := range items {
		listItems[i] = mediatypes.ListItem{
			ItemID:      item.ID,
			Position:    i,
			LastChanged: now,
			ChangeHistory: []mediatypes.ChangeRecord{
				{
					ClientID:   0, // 0 indicates application level change
					ItemID:     fmt.Sprintf("%d", item.ID),
					ChangeType: "update",
					Timestamp:  now,
				},
			},
		}
	}

	// Replace all items
	playlist.Data.ItemList.Items = listItems
	playlist.Data.ItemList.LastModified = now
	playlist.Data.ItemList.ModifiedBy = 0 // 0 indicates application level modification
	playlist.Data.ItemList.ItemCount = len(listItems)

	// Ensure positions are normalized
	playlist.Data.NormalizePositions()

	// Update the playlist
	_, err = s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Msg("Failed to update playlist items")
		return fmt.Errorf("failed to update playlist items: %w", err)
	}

	log.Info().
		Uint64("playlistID", playlistID).
		Int("itemCount", len(items)).
		Msg("Playlist items updated successfully")

	return nil
}

func (s *playlistService) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Uint64("userID", query.OwnerID).
		Msg("Searching playlists")

	query.MediaType = mediatypes.MediaTypePlaylist

	// Delegate to the user service
	results, err := s.repo.Search(ctx, query)
	if err != nil {
		log.Error().Err(err).
			Str("query", query.Query).
			Uint64("userID", query.OwnerID).
			Msg("Failed to search playlists")
		return nil, fmt.Errorf("failed to search playlists: %w", err)
	}

	log.Info().
		Str("query", query.Query).
		Uint64("userID", query.OwnerID).
		Int("count", len(results)).
		Msg("Playlists found")

	return results, nil
}

func (s *playlistService) GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Int("limit", limit).
		Int("days", days).
		Msg("Getting recent playlists")

	// Delegate to the user service
	results, err := s.repo.GetRecentItems(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Msg("Failed to get recent playlists")
		return nil, fmt.Errorf("failed to get recent playlists: %w", err)
	}

	log.Info().
		Int("count", len(results)).
		Msg("Recent playlists retrieved")

	return results, nil
}

func (s *playlistService) Sync(ctx context.Context, playlistID uint64, targetClientIDs []uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Interface("targetClientIDs", targetClientIDs).
		Msg("Syncing playlist")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to sync playlist: %w", err)
	}

	// In a real implementation, this would use the playlist sync job to:
	// 1. For each target client ID, create or update a playlist with the same name
	// 2. Map the item IDs in playlist.Data.ItemList.Items to client-specific IDs using the mediaItemRepo
	// 3. Add these items to the client-specific playlist
	// 4. Store the client-specific playlist IDs and item IDs in the SyncClientStates
	// 5. Update the LastSynced timestamp

	// For now, just create a placeholder in the SyncClientStates to show intent
	now := time.Now()
	for _, clientID := range targetClientIDs {
		// Create a placeholder sync client state
		syncState := playlist.Data.ItemList.SyncClientStates.GetSyncClientState(clientID)
		if syncState == nil {
			// Create empty sync list items
			syncItems := mediatypes.SyncListItems{}

			// Add a new state for this client
			newState := mediatypes.SyncClientState{
				ClientID:     clientID,
				Items:        syncItems,
				ClientListID: "", // Empty for now, would be the client-specific playlist ID
				LastSynced:   now,
			}

			playlist.Data.ItemList.SyncClientStates = append(playlist.Data.ItemList.SyncClientStates, newState)
		}
	}

	// Update the LastSynced timestamp
	playlist.Data.ItemList.LastSynced = now

	// Save the updated playlist
	_, err = s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Msg("Failed to update playlist sync state")
		return fmt.Errorf("failed to update playlist sync state: %w", err)
	}

	log.Warn().Msg("Playlist sync partially implemented - requires job implementation")

	// This would be implemented in a job that handles the actual sync
	return errors.New("playlist sync requires implementation in the playlist sync job")
}

//
// // Helper functions
//
// // createStubItems creates stub media items for all items in a playlist
// func createStubItems(listItems []mediatypes.ListItem) []models.MediaItem[mediatypes.MediaData] {
// 	var items []models.MediaItem[mediatypes.MediaData]
// 	for _, listItem := range listItems {
// 		items = append(items, createStubItem(listItem))
// 	}
// 	return items
// }
//
// // createStubItem creates a stub media item for a list item
// func createStubItem(listItem mediatypes.ListItem) models.MediaItem[mediatypes.MediaData] {
// 	return models.MediaItem[mediatypes.MediaData]{
// 		ID:          listItem.ItemID,
// 		Type:        mediatypes.MediaTypeTrack, // Placeholder, would get actual type from DB
// 		ReleaseYear: 0,                         // Would be populated from DB
// 		CreatedAt:   listItem.LastChanged,      // Use last changed time as a proxy
// 		UpdatedAt:   listItem.LastChanged,
// 		// Would populate other fields from the actual item in the database
// 	}
// }

func (s *playlistService) createMediaItems(ctx context.Context, items []*models.MediaItem[mediatypes.MediaData]) (*models.MediaItems, error) {
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Int("count", len(items)).
		Msg("Creating media items")

	return nil, nil
}
