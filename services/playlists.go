package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	// "strings"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
)

// PlaylistService manages application-specific playlist operations beyond the basic CRUD operations
type PlaylistService interface {
	// Base operations (wrapping MediaItemService)
	Create(ctx context.Context, playlist models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error)
	Update(ctx context.Context, playlist models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Playlist], error)
	GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error)
	Delete(ctx context.Context, id uint64) error

	// Playlist-specific operations
	GetPlaylistItems(ctx context.Context, playlistID uint64) ([]models.MediaItem[mediatypes.MediaData], error)
	AddItemToPlaylist(ctx context.Context, playlistID uint64, item models.MediaItem[mediatypes.MediaData]) error
	RemoveItemFromPlaylist(ctx context.Context, playlistID uint64, itemID uint64) error
	ReorderPlaylistItems(ctx context.Context, playlistID uint64, itemIDs []string) error
	UpdatePlaylistItems(ctx context.Context, playlistID uint64, items []models.MediaItem[mediatypes.MediaData]) error
	SearchPlaylists(ctx context.Context, query string, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error)
	GetRecentPlaylists(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Playlist], error)
	SyncPlaylist(ctx context.Context, playlistID uint64, targetClientIDs []uint64) error
}

type playlistService struct {
	repo         repository.MediaItemRepository[*mediatypes.Playlist]
	mediaItemSvc MediaItemService[*mediatypes.Playlist]
}

// NewPlaylistService creates a new playlist service
func NewPlaylistService(
	repo repository.MediaItemRepository[*mediatypes.Playlist],
) PlaylistService {
	// Create the base media item service
	mediaItemSvc := NewMediaItemService(repo)

	return &playlistService{
		repo:         repo,
		mediaItemSvc: mediaItemSvc,
	}
}

// Base operations (delegating to MediaItemService)

func (s *playlistService) Create(ctx context.Context, playlist models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error) {
	// Ensure playlist-specific validation
	if playlist.Type != mediatypes.MediaTypePlaylist {
		playlist.Type = mediatypes.MediaTypePlaylist
	}

	// Ensure playlist has a valid name
	if playlist.Data == nil || playlist.Data.Details.Title == "" {
		return nil, errors.New("playlist must have a title")
	}

	// Initialize items array if nil
	if playlist.Data.Items == nil {
		playlist.Data.Items = []mediatypes.PlaylistItem{}
	}

	// Initialize sync client states if nil
	if playlist.Data.SyncClientStates == nil {
		playlist.Data.SyncClientStates = mediatypes.SyncClientStates{}
	}

	// Set creation time for LastModified
	if playlist.Data.LastModified.IsZero() {
		playlist.Data.LastModified = time.Now()
	}

	return s.mediaItemSvc.Create(ctx, playlist)
}

func (s *playlistService) Update(ctx context.Context, playlist models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error) {
	// Ensure the playlist exists
	existing, err := s.GetByID(ctx, playlist.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update playlist: %w", err)
	}

	// Preserve items if not provided in the update
	if playlist.Data.Items == nil || len(playlist.Data.Items) == 0 {
		playlist.Data.Items = existing.Data.Items
	}

	// Preserve sync client states if not provided
	if playlist.Data.SyncClientStates == nil || len(playlist.Data.SyncClientStates) == 0 {
		playlist.Data.SyncClientStates = existing.Data.SyncClientStates
	}

	// Update last modified time
	playlist.Data.LastModified = time.Now()

	// Update the playlist
	return s.mediaItemSvc.Update(ctx, playlist)
}

func (s *playlistService) GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Playlist], error) {
	return s.mediaItemSvc.GetByID(ctx, id)
}

func (s *playlistService) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	return s.mediaItemSvc.GetByUserID(ctx, userID)
}

func (s *playlistService) Delete(ctx context.Context, id uint64) error {
	return s.mediaItemSvc.Delete(ctx, id)
}

// Playlist-specific operations

func (s *playlistService) GetPlaylistItems(ctx context.Context, playlistID uint64) ([]models.MediaItem[mediatypes.MediaData], error) {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist items: %w", err)
	}

	// Return empty array if the playlist has no items
	if len(playlist.Data.Items) == 0 && len(playlist.Data.ItemIDs) == 0 {
		return []models.MediaItem[mediatypes.MediaData]{}, nil
	}

	// Implementation note: In a real application, this would need access to:
	// 1. The mediaItemRepo to look up each item by ID
	// 2. Additional logic to handle the case where the item is stored in another client

	// As a placeholder implementation, we'll just create stub items based on the available IDs
	var items []models.MediaItem[mediatypes.MediaData]

	// First, try to get items from the Items array (the new format)
	if len(playlist.Data.Items) > 0 {
		for _, playlistItem := range playlist.Data.Items {
			// Try to parse the item ID
			itemID, err := uint64(playlistItem.ItemID)
			if err != nil {
				// If we can't parse it, it might be a client-specific ID format
				// In a complete implementation, we would search the mediaItemRepo
				// for items with this client-specific ID
				continue
			}

			// Create a stub item (in a real implementation, we would fetch from repository)
			stubItem := models.MediaItem[mediatypes.MediaData]{
				ID:          itemID,
				Type:        mediatypes.MediaTypeTrack, // Placeholder, would get actual type from DB
				ReleaseYear: 0,                         // Would be populated from DB
				CreatedAt:   playlistItem.LastChanged,  // Use last changed time as a proxy
				UpdatedAt:   playlistItem.LastChanged,
				// Would populate other fields from the actual item in the database
			}

			items = append(items, stubItem)
		}
	} else if len(playlist.Data.ItemIDs) > 0 {
		// Fall back to ItemIDs array (legacy format)
		for _, itemID := range playlist.Data.ItemIDs {
			// Create a stub item (in a real implementation, we would fetch from repository)
			stubItem := models.MediaItem[mediatypes.MediaData]{
				ID:          itemID,
				Type:        mediatypes.MediaTypeTrack, // Placeholder, would get actual type from DB
				ReleaseYear: 0,                         // Would be populated from DB
				CreatedAt:   time.Now(),                // Don't have a timestamp in the old format
				UpdatedAt:   time.Now(),
				// Would populate other fields from the actual item in the database
			}

			items = append(items, stubItem)
		}
	}

	// Sort items by position if available
	if len(playlist.Data.Items) > 0 {
		sortByPosition(items)
	}

	return items, nil
}

func (s *playlistService) AddItemToPlaylist(ctx context.Context, playlistID uint64, item models.MediaItem[mediatypes.MediaData]) error {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to add item to playlist: %w", err)
	}

	// Check if the item already exists in the playlist
	itemIDStr := fmt.Sprintf("%d", item.ID)
	for _, existingItem := range playlist.Data.Items {
		if existingItem.ItemID == itemIDStr {
			return errors.New("item already exists in playlist")
		}
	}

	// Add the item to the playlist
	now := time.Now()
	newItem := mediatypes.PlaylistItem{
		ItemID:      itemIDStr,
		Position:    len(playlist.Data.Items),
		LastChanged: now,
		ChangeHistory: []mediatypes.ChangeRecord{
			{
				ClientID:   0, // 0 indicates application level change
				ItemID:     itemIDStr,
				ChangeType: "add",
				Timestamp:  now,
			},
		},
	}

	// Update the playlist
	playlist.Data.Items = append(playlist.Data.Items, newItem)
	playlist.Data.LastModified = now
	playlist.Data.ModifiedBy = 0 // 0 indicates application level modification

	// Update ItemIDs for backward compatibility
	playlist.Data.ItemIDs = append(playlist.Data.ItemIDs, item.ID)

	// Update ItemCount
	playlist.Data.ItemCount = len(playlist.Data.Items)

	// Store the update
	_, err = s.Update(ctx, *playlist)
	return err
}

func (s *playlistService) RemoveItemFromPlaylist(ctx context.Context, playlistID uint64, itemID uint64) error {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to remove item from playlist: %w", err)
	}

	// Find and remove the item
	var newItems []mediatypes.PlaylistItem
	var newItemIDs []uint64
	itemFound := false
	itemIDStr := fmt.Sprintf("%d", itemID)

	for _, item := range playlist.Data.Items {
		if item.ItemID != itemIDStr {
			newItems = append(newItems, item)
		} else {
			itemFound = true
		}
	}

	// Update ItemIDs array for backward compatibility
	for _, id := range playlist.Data.ItemIDs {
		if id != itemID {
			newItemIDs = append(newItemIDs, id)
		}
	}

	if !itemFound {
		return errors.New("item not found in playlist")
	}

	// Update the playlist with the new items
	playlist.Data.Items = newItems
	playlist.Data.ItemIDs = newItemIDs
	playlist.Data.LastModified = time.Now()
	playlist.Data.ModifiedBy = 0 // 0 indicates application level modification
	playlist.Data.ItemCount = len(newItems)

	_, err = s.Update(ctx, *playlist)
	return err
}

func (s *playlistService) ReorderPlaylistItems(ctx context.Context, playlistID uint64, itemIDs []string) error {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to reorder playlist items: %w", err)
	}

	// Verify that the number of items matches
	if len(itemIDs) != len(playlist.Data.Items) {
		return errors.New("reorder operation must include all playlist items")
	}

	// Create a map of existing items
	itemMap := make(map[string]mediatypes.PlaylistItem)
	for _, item := range playlist.Data.Items {
		itemMap[item.ItemID] = item
	}

	// Verify that all items exist in the playlist
	for _, itemID := range itemIDs {
		if _, ok := itemMap[itemID]; !ok {
			return fmt.Errorf("item %s not found in playlist", itemID)
		}
	}

	// Reorder the items
	var newItems []mediatypes.PlaylistItem
	var newItemIDs []uint64
	now := time.Now()

	for i, itemID := range itemIDs {
		item := itemMap[itemID]
		item.Position = i
		item.LastChanged = now
		item.ChangeHistory = append(item.ChangeHistory, mediatypes.ChangeRecord{
			ClientID:   0, // 0 indicates application level change
			ItemID:     itemID,
			ChangeType: "reorder",
			Timestamp:  now,
		})
		newItems = append(newItems, item)

		// Parse the item ID for the ItemIDs array
		if id, err := parseUint64(itemID); err == nil {
			newItemIDs = append(newItemIDs, id)
		}
	}

	// Update the playlist
	playlist.Data.Items = newItems
	playlist.Data.ItemIDs = newItemIDs
	playlist.Data.LastModified = now
	playlist.Data.ModifiedBy = 0 // 0 indicates application level modification

	_, err = s.Update(ctx, *playlist)
	return err
}

func (s *playlistService) UpdatePlaylistItems(ctx context.Context, playlistID uint64, items []models.MediaItem[mediatypes.MediaData]) error {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to update playlist items: %w", err)
	}

	// Convert MediaItems to PlaylistItems
	var playlistItems []mediatypes.PlaylistItem
	var itemIDs []uint64
	now := time.Now()

	for i, item := range items {
		playlistItems = append(playlistItems, mediatypes.PlaylistItem{
			ItemID:      fmt.Sprintf("%d", item.ID),
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
		})
		itemIDs = append(itemIDs, item.ID)
	}

	// Replace all items
	playlist.Data.Items = playlistItems
	playlist.Data.ItemIDs = itemIDs
	playlist.Data.LastModified = now
	playlist.Data.ModifiedBy = 0 // 0 indicates application level modification
	playlist.Data.ItemCount = len(playlistItems)

	// Update the playlist
	_, err = s.Update(ctx, *playlist)
	return err
}

func (s *playlistService) SearchPlaylists(ctx context.Context, query string, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	// Get all playlists for the user
	playlists, err := s.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to search playlists: %w", err)
	}

	// Filter playlists by query
	var results []*models.MediaItem[*mediatypes.Playlist]
	for _, playlist := range playlists {
		if containsIgnoreCase(playlist.Data.Details.Title, query) ||
			containsIgnoreCase(playlist.Data.Details.Description, query) {
			results = append(results, playlist)
		}
	}

	return results, nil
}

func (s *playlistService) GetRecentPlaylists(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	// Get all playlists for the user
	playlists, err := s.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent playlists: %w", err)
	}

	// Sort by last modified date
	sort.Slice(playlists, func(i, j int) bool {
		return playlists[i].Data.LastModified.After(playlists[j].Data.LastModified)
	})

	// Limit results
	if len(playlists) > limit {
		playlists = playlists[:limit]
	}

	return playlists, nil
}

func (s *playlistService) SyncPlaylist(ctx context.Context, playlistID uint64, targetClientIDs []uint64) error {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to sync playlist: %w", err)
	}

	// In a real implementation, this would use the playlist sync job to:
	// 1. For each target client ID, create or update a playlist with the same name
	// 2. Map the item IDs in playlist.Data.Items to client-specific IDs using the mediaItemRepo
	// 3. Add these items to the client-specific playlist
	// 4. Store the client-specific playlist IDs and item IDs in the SyncClientStates
	// 5. Update the LastSynced timestamp

	// For now, just create a placeholder in the SyncClientStates to show intent
	now := time.Now()
	for _, clientID := range targetClientIDs {
		// Create a placeholder sync client state
		syncState := playlist.Data.SyncClientStates.GetSyncClientState(clientID)
		if syncState == nil {
			// Add a new state for this client
			playlist.Data.SyncClientStates.AddOrUpdateSyncClientState(
				clientID,
				[]string{}, // Empty for now, would contain client-specific item IDs
				"",         // Empty for now, would be the client-specific playlist ID
			)
		}
	}

	// Update the LastSynced timestamp
	playlist.Data.LastSynced = now

	// Save the updated playlist
	_, err = s.Update(ctx, *playlist)
	if err != nil {
		return fmt.Errorf("failed to update playlist sync state: %w", err)
	}

	// This would be implemented in a job that handles the actual sync
	return errors.New("playlist sync requires implementation in the playlist sync job")
}

// Helper functions

// func containsIgnoreCase(s, substr string) bool {
// 	s, substr = strings.ToLower(s), strings.ToLower(substr)
// 	return strings.Contains(s, substr)
// }
//
// func parseUint64(s string) (uint64, error) {
// 	return strings.ParseUint(s, 10, 64)
// }

// sortByPosition sorts a slice of media items by their position in the playlist
// This is a simplified version; in a real implementation we would need to
// match each media item to its position in the playlist item array
func sortByPosition(items []models.MediaItem[mediatypes.MediaData]) {
	// This is a placeholder sort - in a real implementation, we would:
	// 1. Create a map of item IDs to positions
	// 2. Use that map to determine the position of each item during sorting

	// For now, we'll just leave the items in the order they were added to the slice
	// which should match the order in the playlist.Data.Items array
	// if we processed them in order during collection
}

