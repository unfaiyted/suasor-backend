package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	database "suasor/utils/db"
)

// TestDatabaseIntegration tests PostgreSQL database integration
// This test demonstrates how to use the database test helpers
func TestDatabaseIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database using our helper
	dbHelper := database.NewTestDBHelper(t)

	// Create a test user
	user := dbHelper.CreateTestUser(t, "integrationUser", "integration@test.com", "user")
	assert.NotZero(t, user.ID, "User should have a non-zero ID")

	// Run test within a transaction for automatic cleanup
	dbHelper.WithTransaction(t, func(tx *gorm.DB) {
		// Create a test movie
		movieData := &mediatypes.Movie{
			Details: mediatypes.MediaDetails{
				Title:       "Integration Test Movie",
				ReleaseYear: 2023,
			},
		}

		mediaItem := &models.MediaItem[*mediatypes.Movie]{
			Type:        mediatypes.MediaTypeMovie,
			Title:       "Integration Test Movie",
			ReleaseDate: time.Now().AddDate(-1, 0, 0),
			Data:        movieData,
		}

		// Save to database
		result := tx.Create(mediaItem)
		require.NoError(t, result.Error)
		assert.NotZero(t, mediaItem.ID, "Media item should have a non-zero ID")

		// Create user media item data
		userData := &models.UserMediaItemData[*mediatypes.Movie]{
			UserID:        user.ID,
			MediaItemID:   mediaItem.ID,
			Type:          mediatypes.MediaTypeMovie,
			PlayedAt:      time.Now().Add(-24 * time.Hour),
			LastPlayedAt:  time.Now(),
			PlayCount:     1,
			IsFavorite:    true,
			UserRating:    8.5,
		}

		// Save to database
		result = tx.Create(userData)
		require.NoError(t, result.Error)
		assert.NotZero(t, userData.ID, "User media data should have a non-zero ID")

		// Retrieve user media item data
		var retrievedData models.UserMediaItemData[*mediatypes.Movie]
		result = tx.First(&retrievedData, userData.ID)
		require.NoError(t, result.Error)

		// Verify data integrity
		assert.Equal(t, user.ID, retrievedData.UserID)
		assert.Equal(t, mediaItem.ID, retrievedData.MediaItemID)
		assert.Equal(t, mediatypes.MediaTypeMovie, retrievedData.Type)
		assert.True(t, retrievedData.IsFavorite)
		assert.Equal(t, float32(8.5), retrievedData.UserRating)
		assert.Equal(t, int32(1), retrievedData.PlayCount)

		// Load the associated media item
		err := retrievedData.LoadItem(tx)
		require.NoError(t, err)

		// Verify the media item
		assert.NotNil(t, retrievedData.Item)
		assert.Equal(t, mediaItem.ID, retrievedData.Item.ID)
		assert.Equal(t, "Integration Test Movie", retrievedData.Item.Title)
	})
}

// TestMultipleUsersIntegration tests working with multiple users in the database
func TestMultipleUsersIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database using our helper
	dbHelper := database.NewTestDBHelper(t)

	// Create multiple test users
	users := dbHelper.CreateTestUsers(t, 3)
	assert.Equal(t, 3, len(users), "Should create 3 users")

	// Verify created users
	for i, user := range users {
		// Check that the username follows expected pattern (testUserA, testUserB, etc.)
		expectedChar := rune(i + 65) // A, B, C, etc.
		expectedUsername := "testUser" + string(expectedChar)
		expectedEmail := "user" + string(expectedChar) + "@test.com"
		
		assert.Equal(t, expectedUsername, user.Username)
		assert.Equal(t, expectedEmail, user.Email)
		assert.Equal(t, "user", user.Role)
		assert.NotZero(t, user.ID, "User should have a non-zero ID")
		
		// Verify the user exists in the database
		var retrievedUser models.User
		result := dbHelper.DB.First(&retrievedUser, user.ID)
		require.NoError(t, result.Error)
		assert.Equal(t, user.Username, retrievedUser.Username)
	}
}

// TestTransactionRollback verifies that transactions are properly rolled back
func TestTransactionRollback(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	dbHelper := database.NewTestDBHelper(t)

	// Create a test user
	user := dbHelper.CreateTestUser(t, "rollbackUser", "rollback@test.com", "user")
	
	// Get initial user count
	var initialCount int64
	dbHelper.DB.Model(&models.User{}).Count(&initialCount)
	
	// Run in a transaction that will be rolled back
	dbHelper.WithTransaction(t, func(tx *gorm.DB) {
		// Create more users in the transaction
		for i := 0; i < 5; i++ {
			newUser := models.User{
				Username:  "txUser" + string(rune(i+65)),
				Email:     "txuser" + string(rune(i+65)) + "@test.com",
				Role:      "user",
			}
			newUser.SetPassword("TestPassword123")
			
			err := tx.Create(&newUser).Error
			require.NoError(t, err, "Failed to create user in transaction")
		}
		
		// Verify users exist within the transaction
		var txCount int64
		tx.Model(&models.User{}).Count(&txCount)
		assert.Equal(t, initialCount+5, txCount, "Transaction should contain new users")
	})
	
	// After transaction rollback, the count should be the same as initial
	var finalCount int64
	dbHelper.DB.Model(&models.User{}).Count(&finalCount)
	assert.Equal(t, initialCount, finalCount, "Count should be same after transaction rollback")
	
	// Verify no "txUser" exists in the database
	var count int64
	dbHelper.DB.Model(&models.User{}).Where("username LIKE 'txUser%'").Count(&count)
	assert.Equal(t, int64(0), count, "No transactional users should exist after rollback")
}