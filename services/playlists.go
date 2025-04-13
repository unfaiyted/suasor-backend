package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
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
	if playlist.Type != mediatypes.TypePlaylist {
		playlist.Type = mediatypes.TypePlaylist
	}

	// Ensure playlist has a valid name
	if playlist.Data == nil || playlist.Data.Details.Title == "" {
		return nil, errors.New("playlist must have a title")
	}

	// Initialize items array if nil
	if playlist.Data.Items == nil {
		playlist.Data.Items = []mediatypes.PlaylistItem{}
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

	// Convert PlaylistItem to MediaItem
	var items []models.MediaItem[mediatypes.MediaData]
	
	// In a real implementation, we would query the database to get the actual media items
	// based on the playlist.Data.Items array. For now, this is just a placeholder.
	
	return items, nil
}

func (s *playlistService) AddItemToPlaylist(ctx context.Context, playlistID uint64, item models.MediaItem[mediatypes.MediaData]) error {
	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to add item to playlist: %w", err)
	}

	// Check if the item already exists in the playlist
	for _, existingItem := range playlist.Data.Items {
		if existingItem.ItemID == fmt.Sprintf("%d", item.ID) {
			return errors.New("item already exists in playlist")
		}
	}

	// Add the item to the playlist
	now := time.Now()
	newItem := mediatypes.PlaylistItem{
		ItemID:      fmt.Sprintf("%d", item.ID),
		Position:    len(playlist.Data.Items),
		LastChanged: now,
		ChangeHistory: []mediatypes.ChangeRecord{
			{
				ClientID:   0, // 0 indicates application level change
				ItemID:     fmt.Sprintf("%d", item.ID),
				ChangeType: "add",
				Timestamp:  now,
			},
		},
	}
	
	playlist.Data.Items = append(playlist.Data.Items, newItem)
	playlist.Data.LastModified = now

	// Update the playlist
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
	itemFound := false
	itemIDStr := fmt.Sprintf("%d", itemID)

	for _, item := range playlist.Data.Items {
		if item.ItemID != itemIDStr {
			newItems = append(newItems, item)
		} else {
			itemFound = true
		}
	}

	if !itemFound {
		return errors.New("item not found in playlist")
	}

	// Update the playlist with the new items
	playlist.Data.Items = newItems
	playlist.Data.LastModified = time.Now()
	
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
	}

	// Update the playlist
	playlist.Data.Items = newItems
	playlist.Data.LastModified = now
	
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
	}

	// Replace all items
	playlist.Data.Items = playlistItems
	playlist.Data.LastModified = now

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

	// In a real implementation, we would:
	// 1. For each target client, create or update a corresponding playlist
	// 2. Translate item IDs to client-specific IDs
	// 3. Sync the items to the target client
	// 4. Handle errors and conflicts

	// For now, we'll return a not implemented error
	return errors.New("playlist sync not implemented")
}