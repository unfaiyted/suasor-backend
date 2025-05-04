package repository

import (
	"context"
	"testing"
	
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/db"
	"suasor/utils/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// This test demonstrates the issue with the ownerID field in playlists/collections
func TestPlaylistOwnerIDType(t *testing.T) {
	// Set up a test database
	testDB, err := db.NewInMemoryDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	
	// Create a context with logger
	ctx := logger.NewLoggerContext(context.Background())
	
	// Create a repository for playlists
	repo := NewMediaItemRepository[*types.Playlist](testDB)
	
	// Create a test playlist with a user ID
	const userID uint64 = 1
	
	// Create a playlist
	itemList := types.ItemList{
		Details: &types.MediaDetails{
			Title:       "Test Playlist",
			Description: "A test playlist for user ID testing",
		},
		Items:     []types.ListItem{},
		ItemCount: 0,
		OwnerID:   userID,
		IsPublic:  true,
	}
	
	playlist := &types.Playlist{
		ItemList: itemList,
	}
	
	// Create the media item
	mediaItem := models.NewMediaItem[*types.Playlist](types.MediaTypePlaylist, playlist)
	mediaItem.OwnerID = userID
	mediaItem.UUID = uuid.New().String()
	
	// Save the playlist
	err = repo.Create(ctx, mediaItem)
	if err != nil {
		t.Fatalf("Failed to create playlist: %v", err)
	}
	
	// Test direct query using the current approach
	var items []*models.MediaItem[*types.Playlist]
	
	// This is the query that's failing - we'll test both approaches
	// 1. Current approach with ->>' operator
	queryFailing := testDB.WithContext(ctx).
		Where("type IN (?) AND data->'itemList'->>'ownerID' = ?", types.MediaTypePlaylist, userID)
	
	// 2. Fixed approach with -> operator and casting
	queryFixed := testDB.WithContext(ctx).
		Where("type IN (?) AND CAST(data->'itemList'->'ownerID' AS INTEGER) = ?", types.MediaTypePlaylist, userID)
		
	// Try the failing query and log the error
	err = queryFailing.Find(&items).Error
	t.Logf("Original query error: %v", err)
	
	// Try the fixed query
	items = nil // Reset items
	err = queryFixed.Find(&items).Error
	
	// This should work
	assert.NoError(t, err, "Fixed query should work")
	assert.Len(t, items, 1, "Should find one playlist")
	
	if len(items) > 0 {
		assert.Equal(t, userID, items[0].OwnerID, "Owner ID should match")
		assert.Equal(t, "Test Playlist", items[0].Title, "Title should match")
	}
	
	// Clean up
	testDB.Unscoped().Delete(&models.MediaItem[*types.Playlist]{}, "type = ?", types.MediaTypePlaylist)
}