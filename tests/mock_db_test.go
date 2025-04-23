package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	database "suasor/utils/db"
)

// ExampleRepositoryFunction is an example repository function that would use the database
func ExampleRepositoryFunction(ctx context.Context, db *gorm.DB, userID uint64) ([]models.UserMediaItemData[*mediatypes.Movie], error) {
	var results []models.UserMediaItemData[*mediatypes.Movie]
	if err := db.Where("user_id = ?", userID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// TestMockDBUsage demonstrates how to use the mock database
func TestMockDBUsage(t *testing.T) {
	// Create a mock DB
	mockDB := database.NewMockDB()

	// Create a mock GormDB from the mock
	mockGormDB := &database.MockGormDB{
		mock: mockDB,
	}

	// Set up expectations
	mockDB.On("Called", mock.Anything, mock.Anything).Return(nil)

	// Create test data
	mockUser := database.MockUser(1, "testUser", "test@example.com", "user")
	mockMovieData := &mediatypes.Movie{
		Details: mediatypes.MediaDetails{
			Title:       "Mock Movie",
			ReleaseYear: 2023,
		},
	}
	mockMovie := database.MockMediaItem(2, "Mock Movie", mediatypes.MediaTypeMovie, mockMovieData)
	mockUserData := database.MockUserMediaItemData[*mediatypes.Movie](3, mockUser.ID, mockMovie.ID, mediatypes.MediaTypeMovie)

	// Mock user media item data results
	mockResults := []models.UserMediaItemData[*mediatypes.Movie]{mockUserData}

	// Set up expectations for the Where and Find methods
	mockDB.On("Where", "user_id = ?", uint64(1)).Return(mockGormDB)
	mockDB.On("Find", mock.AnythingOfType("*[]models.UserMediaItemData[*types.Movie]"), []interface{}(nil)).Run(func(args mock.Arguments) {
		// When Find is called, we set the results in the provided slice
		dest := args.Get(0).(*[]models.UserMediaItemData[*mediatypes.Movie])
		*dest = mockResults
	}).Return(&gorm.DB{Error: nil})

	// Call the function that uses the database
	ctx := context.Background()
	results, err := ExampleRepositoryFunction(ctx, mockGormDB, 1)
	require.NoError(t, err)

	// Assert the results
	assert.Equal(t, 1, len(results))
	assert.Equal(t, mockUser.ID, results[0].UserID)
	assert.Equal(t, mockMovie.ID, results[0].MediaItemID)
	assert.Equal(t, mediatypes.MediaTypeMovie, results[0].Type)

	// Verify expectations
	mockDB.AssertExpectations(t)
}