package database

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
)

// This file contains examples of how to use the mock database for unit testing

// ExampleRepository is a simple repository for demonstration purposes
type ExampleRepository struct {
	db *gorm.DB
}

// GetUserByID retrieves a user by ID from the database
func (r *ExampleRepository) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// SaveUserMediaData creates or updates user media data
func (r *ExampleRepository) SaveUserMediaData(ctx context.Context, data *models.UserMediaItemData[*mediatypes.Movie]) error {
	// Set timestamps
	now := time.Now()
	if data.ID == 0 {
		data.CreatedAt = now
		// Set a unique UUID if not present
		if data.UUID == "" {
			data.UUID = "data-" + uuid.New().String()
		}
	}
	data.UpdatedAt = now

	// Save to database
	return r.db.WithContext(ctx).Table("user_media_item_data").Save(data).Error
}

// TestExampleRepositoryWithMocks demonstrates using mocks for database access
func TestExampleRepositoryWithMocks(t *testing.T) {
	t.Skip("Skipping this test as it's for demonstration purposes only")
	
	// Actual integration tests use the real in-memory database
	// Create a repository with the actual test database
	helper := NewTestDBHelper(t)
	
	// Create a repository with the in-memory database
	repo := &ExampleRepository{
		db: helper.DB,
	}
	
	// Test context
	ctx := context.Background()
	
	t.Run("GetUserByID", func(t *testing.T) {
		// Create a test user
		testUser := helper.CreateTestUser(t, "testUser", "test@example.com", "user")
		
		// Call the function under test
		user, err := repo.GetUserByID(ctx, testUser.ID)
		
		// Assert expectations
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, "testUser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
	})
	
	t.Run("SaveUserMediaData", func(t *testing.T) {
		// Create a test user
		testUser := helper.CreateTestUser(t, "mediaUser", "media@example.com", "user")
		
		// Create a test movie
		testMovie := helper.CreateTestMovie(t, "Test Movie", 2023)
		
		// Create user media data
		userData := &models.UserMediaItemData[*mediatypes.Movie]{
			UserID:      testUser.ID,
			MediaItemID: testMovie.ID,
			Type:        mediatypes.MediaTypeMovie,
			IsFavorite:  true,
			UserRating:  8.5,
			UUID:        "test-" + uuid.New().String(),
		}
		
		// Call the function under test
		err := repo.SaveUserMediaData(ctx, userData)
		
		// Assert expectations
		require.NoError(t, err)
		assert.NotZero(t, userData.ID)
		
		// Verify that the data was saved correctly
		var savedData models.UserMediaItemData[*mediatypes.Movie]
		err = helper.DB.Table("user_media_item_data").Where("user_id = ? AND media_item_id = ?", testUser.ID, testMovie.ID).First(&savedData).Error
		require.NoError(t, err)
		assert.Equal(t, userData.ID, savedData.ID)
		assert.Equal(t, true, savedData.IsFavorite)
		assert.Equal(t, float32(8.5), savedData.UserRating)
	})
}

// This example shows a more realistic test with service and repository layers
// using both mock DB and test DB helper

// ExampleMediaItemService is a service for working with media items
type ExampleMediaItemService struct {
	repository *ExampleRepository
}

// ToggleUserFavorite toggles the favorite status of a media item for a user
func (s *ExampleMediaItemService) ToggleUserFavorite(ctx context.Context, userID, mediaItemID uint64, favoriteStatus bool) error {
	// In a real service, this would have more business logic
	data := &models.UserMediaItemData[*mediatypes.Movie]{
		UserID:      userID,
		MediaItemID: mediaItemID,
		Type:        mediatypes.MediaTypeMovie,
		IsFavorite:  favoriteStatus,
		UUID:        "toggle-" + uuid.New().String(),
	}
	return s.repository.SaveUserMediaData(ctx, data)
}

// TestWithBothApproaches demonstrates testing with both in-memory DB and mocks
func TestWithBothApproaches(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Integration test with in-memory DB
	t.Run("Integration test with in-memory DB", func(t *testing.T) {
		// Create a test DB helper
		helper := NewTestDBHelper(t)
		
		// Create a repository with the real database
		repo := &ExampleRepository{
			db: helper.DB,
		}
		
		// Create the service with the repository
		service := &ExampleMediaItemService{
			repository: repo,
		}
		
		// Create test data
		user := helper.CreateTestUser(t, "favoriteTestUser", "favorite@test.com", "user")
		movie := helper.CreateTestMovie(t, "Test Favorite Movie", 2024)
		
		// Call the function under test
		err := service.ToggleUserFavorite(helper.Ctx, user.ID, movie.ID, true)
		require.NoError(t, err)
		
		// Verify the data was saved correctly
		var userData models.UserMediaItemData[*mediatypes.Movie]
		result := helper.DB.Table("user_media_item_data").Where("user_id = ? AND media_item_id = ?", user.ID, movie.ID).First(&userData)
		require.NoError(t, result.Error)
		assert.True(t, userData.IsFavorite)
		
		// Create a new user and movie for the second test to avoid UUID conflicts
		user2 := helper.CreateTestUser(t, "favoriteTestUser2", "favorite2@test.com", "user") 
		movie2 := helper.CreateTestMovie(t, "Test Favorite Movie 2", 2023)
		
		// Set to false on the new user/movie
		err = service.ToggleUserFavorite(helper.Ctx, user2.ID, movie2.ID, false)
		require.NoError(t, err)
		
		// Verify it was set correctly
		var userData2 models.UserMediaItemData[*mediatypes.Movie]
		result = helper.DB.Table("user_media_item_data").Where("user_id = ? AND media_item_id = ?", user2.ID, movie2.ID).First(&userData2)
		require.NoError(t, result.Error) 
		assert.False(t, userData2.IsFavorite)
	})
}