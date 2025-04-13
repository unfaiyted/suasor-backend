package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
)

// PlaylistService manages application-specific playlist operations beyond the basic CRUD operations
type PlaylistService interface {
	// Base operations (wrapping MediaItemService)
	Create(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error)
	Update(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Playlist], error)
	GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error)
	Delete(ctx context.Context, id uint64) error

	// Playlist-specific operations
	GetPlaylistItems(ctx context.Context, playlistID uint64) ([]models.MediaItem[mediatypes.MediaData], error)
	AddItemToPlaylist(ctx context.Context, playlistID uint64, itemID uint64) error
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

func (s *playlistService) Create(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error) {
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
	if playlist.Data.ItemList.Owner == 0 && playlist.Data.Owner != 0 {
		playlist.Data.ItemList.Owner = playlist.Data.Owner
	}

	return s.mediaItemSvc.Create(ctx, *playlist)
}

func (s *playlistService) Update(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error) {
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

	// Run validation to check for issues
	issues := playlist.Data.ValidateItems()
	if len(issues) > 0 {
		// Log the issues but continue with the update
		for _, issue := range issues {
			fmt.Printf("Warning: Playlist validation issue: %s\n", issue)
		}
	}

	// Update the playlist
	return s.mediaItemSvc.Update(ctx, *playlist)
}

func (s *playlistService) GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Playlist], error) {
	return s.mediaItemSvc.GetByID(ctx, id)
}

func (s *playlistService) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	return s.mediaItemSvc.GetByUserID(ctx, userID)
}

func (s *playlistService) Delete(ctx context.Context, id uint64) error {
	// TODO: check if logged in user owns the playlist before deleting
	return s.mediaItemSvc.Delete(ctx, id)
}

// Playlist-specific operations

func (s *playlistService) GetPlaylistItems(ctx context.Context, playlistID uint64) ([]models.MediaItem[mediatypes.MediaData], error) {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return []models.MediaItem[mediatypes.MediaData]{}, fmt.Errorf("failed to get playlist items: %w", err)
	}

	// Return empty array if the playlist has no items
	if len(playlist.Data.ItemList.Items) == 0 {
		return []models.MediaItem[mediatypes.MediaData]{}, nil
	}

	// Implementation note: this would need access to:
	// 1. The mediaItemRepo to look up each item by ID
	// 2. Additional logic to handle the case where the item is stored in another client

	// As a placeholder implementation, we'll just create stub items based on the available IDs
	var items []models.MediaItem[mediatypes.MediaData]

	// Get items from the ItemList.Items array
	for _, listItem := range playlist.Data.ItemList.Items {
		// Create a stub item (in a real implementation, we would fetch from repository)
		stubItem := models.MediaItem[mediatypes.MediaData]{
			ID:          listItem.ItemID,
			Type:        mediatypes.MediaTypeTrack, // Placeholder, would get actual type from DB
			ReleaseYear: 0,                         // Would be populated from DB
			CreatedAt:   listItem.LastChanged,      // Use last changed time as a proxy
			UpdatedAt:   listItem.LastChanged,
			// Would populate other fields from the actual item in the database
		}

		items = append(items, stubItem)
	}

	// The items are already sorted by position in the ItemList
	return items, nil
}

func (s *playlistService) AddItemToPlaylist(ctx context.Context, playlistID uint64, itemID uint64) error {
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
	playlist.Data.AddItem(newItem, 0)

	// Store the update
	_, err = s.Update(ctx, playlist)
	return err
}

func (s *playlistService) RemoveItemFromPlaylist(ctx context.Context, playlistID uint64, itemID uint64) error {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to remove item from playlist: %w", err)
	}

	// Use the RemoveItem method provided by ItemList
	// 0 indicates application level modification
	err = playlist.Data.RemoveItem(itemID, 0)
	if err != nil {
		return err
	}

	// Store the update
	_, err = s.Update(ctx, playlist)
	return err
}

func (s *playlistService) ReorderPlaylistItems(ctx context.Context, playlistID uint64, itemIDs []string) error {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to reorder playlist items: %w", err)
	}

	// Verify that the number of items matches
	if len(itemIDs) != len(playlist.Data.ItemList.Items) {
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
			return fmt.Errorf("invalid item ID format: %s", idStr)
		}

		if _, exists := tempItems[id]; !exists {
			missingItems = append(missingItems, idStr)
		}
	}

	if len(missingItems) > 0 {
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
	return err
}

func (s *playlistService) UpdatePlaylistItems(ctx context.Context, playlistID uint64, items []models.MediaItem[mediatypes.MediaData]) error {
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
		if containsIgnoreCase(playlist.Data.ItemList.Details.Title, query) ||
			containsIgnoreCase(playlist.Data.ItemList.Details.Description, query) {
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
		return playlists[i].Data.ItemList.LastModified.After(playlists[j].Data.ItemList.LastModified)
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
		return fmt.Errorf("failed to update playlist sync state: %w", err)
	}

	// This would be implemented in a job that handles the actual sync
	return errors.New("playlist sync requires implementation in the playlist sync job")
}
