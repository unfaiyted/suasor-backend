package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
)

// This file contains integration tests that demonstrate how to use the test database utilities

// TestIntegrationWithTestDBHelper demonstrates using the TestDBHelper in an integration test
func TestIntegrationWithTestDBHelper(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 1. Initialize the test database helper
	helper := NewTestDBHelper(t)

	// 2. Demonstrate using transactions for test isolation
	t.Run("Use transactions for test isolation", func(t *testing.T) {
		// Create a user outside the transaction (will persist)
		permanentUser := helper.CreateTestUser(t, "permanentUser", "permanent@test.com", "user")

		// Use a transaction that will be rolled back
		helper.WithTransaction(t, func(tx *gorm.DB) {
			// Create test data within the transaction
			user := models.User{
				Username: "transactionUser",
				Email:    "transaction@test.com",
				Role:     "user",
			}
			user.SetPassword("TestPassword123")
			tx.Create(&user)

			// Verify the user exists within the transaction
			var foundUser models.User
			result := tx.Where("email = ?", "transaction@test.com").First(&foundUser)
			assert.NoError(t, result.Error)
			assert.Equal(t, "transactionUser", foundUser.Username)
		})

		// Verify that the transaction user doesn't exist after rollback
		var user models.User
		result := helper.DB.Where("email = ?", "transaction@test.com").First(&user)
		assert.Error(t, result.Error)
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error)

		// But the permanent user should still exist
		var persistentUser models.User
		result = helper.DB.Where("email = ?", "permanent@test.com").First(&persistentUser)
		assert.NoError(t, result.Error)
		assert.Equal(t, permanentUser.ID, persistentUser.ID)
	})

	// 3. Demonstrate creating and working with media items
	t.Run("Create and work with media items", func(t *testing.T) {
		// Create test data
		user := helper.CreateTestUser(t, "mediaTestUser", "media@test.com", "user")
		movie := helper.CreateTestMovie(t, "Test Movie for User", 2024)
		series := helper.CreateTestSeries(t, "Test Series for User", 2023)

		// Create user media data for the movie and series
		movieData := helper.CreateUserMediaDataForMovie(t, user, movie, true, 9.5, 2)
		_ = helper.CreateUserMediaDataForSeries(t, user, series, false, 8.0, 1)

		// Verify the data was created correctly
		var retrievedMovieData models.UserMediaItemData[*mediatypes.Movie]
		result := helper.DB.Table("user_media_item_data").Where("user_id = ? AND media_item_id = ?", user.ID, movie.ID).First(&retrievedMovieData)
		require.NoError(t, result.Error)
		assert.Equal(t, float32(9.5), retrievedMovieData.UserRating)
		assert.Equal(t, true, retrievedMovieData.IsFavorite)
		assert.Equal(t, int32(2), retrievedMovieData.PlayCount)

		var retrievedSeriesData models.UserMediaItemData[*mediatypes.Series]
		result = helper.DB.Table("user_media_item_data").Where("user_id = ? AND media_item_id = ?", user.ID, series.ID).First(&retrievedSeriesData)
		require.NoError(t, result.Error)
		assert.Equal(t, float32(8.0), retrievedSeriesData.UserRating)
		assert.Equal(t, false, retrievedSeriesData.IsFavorite)
		assert.Equal(t, int32(1), retrievedSeriesData.PlayCount)

		// Find all favorites
		var favorites []models.UserMediaItemData[*mediatypes.Movie]
		result = helper.DB.Table("user_media_item_data").Where("user_id = ? AND is_favorite = ?", user.ID, true).Find(&favorites)
		require.NoError(t, result.Error)
		assert.Equal(t, 1, len(favorites))
		assert.Equal(t, movieData.ID, favorites[0].ID)
	})

	// 4. Test more complex query capabilities
	t.Run("Complex media queries", func(t *testing.T) {
		// Create multiple users and items
		users := helper.CreateTestUsers(t, 3)
		
		// Create multiple movies with different ratings
		movies := make([]*models.MediaItem[*mediatypes.Movie], 5)
		for i := 0; i < 5; i++ {
			movies[i] = helper.CreateTestMovie(t, "Movie "+string(rune(i+65)), 2020+i)
		}

		// Create user data with varied ratings and favorite status
		for i, user := range users {
			for j, movie := range movies {
				// Create patterns: alternate favorite status and varied ratings
				isFavorite := (i+j)%2 == 0
				rating := float32(5.0 + float32(i) + float32(j)/2.0)
				playCount := int32(i + j + 1)
				
				helper.CreateUserMediaDataForMovie(t, user, movie, isFavorite, rating, playCount)
			}
		}

		// Verify we can find items with complex queries
		// Find all items with rating > 8.0
		var highRatedItems []models.UserMediaItemData[*mediatypes.Movie]
		result := helper.DB.Table("user_media_item_data").Where("user_rating > ?", 8.0).Find(&highRatedItems)
		require.NoError(t, result.Error)
		assert.Greater(t, len(highRatedItems), 0)
		
		for _, item := range highRatedItems {
			assert.Greater(t, item.UserRating, float32(8.0))
		}

		// Find most played items
		var mostPlayed []models.UserMediaItemData[*mediatypes.Movie]
		result = helper.DB.Table("user_media_item_data").Order("play_count DESC").Limit(3).Find(&mostPlayed)
		require.NoError(t, result.Error)
		assert.Equal(t, 3, len(mostPlayed))
		
		// Verify they're ordered by play count
		for i := 0; i < len(mostPlayed)-1; i++ {
			assert.GreaterOrEqual(t, mostPlayed[i].PlayCount, mostPlayed[i+1].PlayCount)
		}
	})
}

// TestMockUserMediaItemFunctions demonstrates working with mock test data
func TestMockUserMediaItemFunctions(t *testing.T) {
	// Create mock users
	user1 := MockUser(1, "mockUser1", "mock1@example.com", "user")
	user2 := MockUser(2, "mockUser2", "mock2@example.com", "admin")
	
	// Create mock media items
	movie := MockMovie(101, "Mock Movie 1")
	series := MockSeries(102, "Mock Series 1")
	
	// Create mock user media data
	movieData := MockUserMediaItemData[*mediatypes.Movie](1001, user1.ID, movie.ID, mediatypes.MediaTypeMovie)
	seriesData := MockUserMediaItemData[*mediatypes.Series](1002, user2.ID, series.ID, mediatypes.MediaTypeSeries)
	
	// Test that the mock data has expected properties
	assert.Equal(t, "mockUser1", user1.Username)
	assert.Equal(t, "Mock Movie 1", movie.Title)
	assert.Equal(t, "Mock Series 1", series.Title)
	assert.Equal(t, uint64(1), user1.ID)
	assert.Equal(t, uint64(101), movie.ID)
	
	// Test that the user media data was created with correct relationships
	assert.Equal(t, user1.ID, movieData.UserID)
	assert.Equal(t, movie.ID, movieData.MediaItemID)
	assert.Equal(t, mediatypes.MediaTypeMovie, movieData.Type)
	
	assert.Equal(t, user2.ID, seriesData.UserID)
	assert.Equal(t, series.ID, seriesData.MediaItemID)
	assert.Equal(t, mediatypes.MediaTypeSeries, seriesData.Type)
	
	// Test that other properties are set to expected defaults
	assert.True(t, movieData.IsFavorite)
	assert.Equal(t, float32(8.5), movieData.UserRating)
	assert.Equal(t, int32(1), movieData.PlayCount)
}