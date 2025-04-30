package database

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
)

// TestDBHelper provides utilities for testing with a database
type TestDBHelper struct {
	DB  *gorm.DB
	Ctx context.Context
}

// NewTestDBHelper creates a new test database helper
func NewTestDBHelper(t *testing.T) *TestDBHelper {
	// Set test environment
	t.Setenv("GO_ENV", "dev")

	// Create context with logger
	ctx := context.Background()
	logger := NewTestLogger()
	ctx = ContextWithTestLogger(ctx, logger)

	// Initialize the in-memory database
	db, err := InitializeInMemoryDB(ctx)
	require.NoError(t, err, "Failed to initialize test database")

	helper := &TestDBHelper{
		DB:  db,
		Ctx: ctx,
	}

	return helper
}

// CreateTestUser creates a test user in the database
func (h *TestDBHelper) CreateTestUser(t *testing.T, username, email, role string) models.User {
	user := models.User{
		Username: username,
		Email:    email,
		Role:     role,
	}
	user.SetPassword("TestPassword123")

	err := h.DB.Create(&user).Error
	require.NoError(t, err, "Failed to create test user")
	return user
}

// CreateTestUsers creates multiple test users in the database
func (h *TestDBHelper) CreateTestUsers(t *testing.T, count int) []models.User {
	users := make([]models.User, count)
	for i := 0; i < count; i++ {
		users[i] = h.CreateTestUser(
			t,
			"testUser"+string(rune(i+65)), // testUserA, testUserB, etc.
			"user"+string(rune(i+65))+"@test.com",
			"user",
		)
	}
	return users
}

// CreateTestMovie creates a test movie in the database
func (h *TestDBHelper) CreateTestMovie(t *testing.T, title string, releaseYear int) *models.MediaItem[*mediatypes.Movie] {
	movieData := &mediatypes.Movie{
		Details: &mediatypes.MediaDetails{
			Title:       title,
			Description: "A test movie created for testing",
			ReleaseYear: releaseYear,
		},
	}

	mediaItem := &models.MediaItem[*mediatypes.Movie]{
		Type:        mediatypes.MediaTypeMovie,
		Title:       title,
		ReleaseYear: releaseYear,
		ReleaseDate: time.Now().AddDate(-1, 0, 0),
		Data:        movieData,
	}

	// For SQLite testing, insert the record using SQL
	now := time.Now()
	id := uint64(0)
	err := h.DB.Raw("INSERT INTO media_items (type, title, release_year, release_date, data, created_at, updated_at) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id",
		mediaItem.Type, mediaItem.Title, mediaItem.ReleaseYear, mediaItem.ReleaseDate,
		`{"Details":{"Title":"`+title+`","Description":"A test movie created for testing","ReleaseYear":`+string(rune(releaseYear))+`}}`,
		now, now).Scan(&id).Error
	require.NoError(t, err, "Failed to create test movie")

	// Create a simple result with the ID
	insertedMovie := &models.MediaItem[*mediatypes.Movie]{
		Data: &mediatypes.Movie{
			Details: &mediatypes.MediaDetails{
				Title:       title,
				Description: "A test movie created for testing",
				ReleaseYear: releaseYear,
			},
		},
	}
	insertedMovie.ID = id
	insertedMovie.Type = mediatypes.MediaTypeMovie
	insertedMovie.Title = title
	insertedMovie.ReleaseYear = releaseYear
	insertedMovie.ReleaseDate = mediaItem.ReleaseDate
	insertedMovie.CreatedAt = now
	insertedMovie.UpdatedAt = now

	return insertedMovie
}

// CreateTestSeries creates a test series in the database
func (h *TestDBHelper) CreateTestSeries(t *testing.T, title string, releaseYear int) *models.MediaItem[*mediatypes.Series] {
	seriesData := &mediatypes.Series{
		Details: &mediatypes.MediaDetails{
			Title:       title,
			Description: "A test series created for testing",
			ReleaseYear: releaseYear,
		},
	}

	mediaItem := &models.MediaItem[*mediatypes.Series]{
		Type:        mediatypes.MediaTypeSeries,
		Title:       title,
		ReleaseYear: releaseYear,
		ReleaseDate: time.Now().AddDate(-1, 0, 0),
		Data:        seriesData,
	}

	// For SQLite testing, insert the record using SQL
	now := time.Now()
	id := uint64(0)
	err := h.DB.Raw("INSERT INTO media_items (type, title, release_year, release_date, data, created_at, updated_at) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id",
		mediaItem.Type, mediaItem.Title, mediaItem.ReleaseYear, mediaItem.ReleaseDate,
		`{"Details":{"Title":"`+title+`","Description":"A test series created for testing","ReleaseYear":`+string(rune(releaseYear))+`}}`,
		now, now).Scan(&id).Error
	require.NoError(t, err, "Failed to create test series")

	// Create a simple result with the ID
	insertedSeries := &models.MediaItem[*mediatypes.Series]{
		Data: &mediatypes.Series{
			Details: &mediatypes.MediaDetails{
				Title:       title,
				Description: "A test series created for testing",
				ReleaseYear: releaseYear,
			},
		},
	}
	insertedSeries.ID = id
	insertedSeries.Type = mediatypes.MediaTypeSeries
	insertedSeries.Title = title
	insertedSeries.ReleaseYear = releaseYear
	insertedSeries.ReleaseDate = mediaItem.ReleaseDate
	insertedSeries.CreatedAt = now
	insertedSeries.UpdatedAt = now

	return insertedSeries
}

// CreateUserMediaDataForMovie creates a user media data entry for a movie
func (h *TestDBHelper) CreateUserMediaDataForMovie(
	t *testing.T,
	user models.User,
	mediaItem *models.MediaItem[*mediatypes.Movie],
	isFavorite bool,
	userRating float32,
	playCount int32,
) *models.UserMediaItemData[*mediatypes.Movie] {
	userData := &models.UserMediaItemData[*mediatypes.Movie]{
		UserID:          user.ID,
		MediaItemID:     mediaItem.ID,
		Type:            mediaItem.Type,
		IsFavorite:      isFavorite,
		UserRating:      userRating,
		PlayCount:       playCount,
		PlayedAt:        time.Now().Add(-24 * time.Hour),
		LastPlayedAt:    time.Now(),
		PositionSeconds: 3600,
		DurationSeconds: 7200,
	}
	userData.UUID = "movie-data-" + uuid.New().String()

	err := h.DB.Table("user_media_item_data").Create(userData).Error
	require.NoError(t, err, "Failed to create user media data")
	return userData
}

// CreateUserMediaDataForSeries creates a user media data entry for a series
func (h *TestDBHelper) CreateUserMediaDataForSeries(
	t *testing.T,
	user models.User,
	mediaItem *models.MediaItem[*mediatypes.Series],
	isFavorite bool,
	userRating float32,
	playCount int32,
) *models.UserMediaItemData[*mediatypes.Series] {
	userData := &models.UserMediaItemData[*mediatypes.Series]{
		UserID:          user.ID,
		MediaItemID:     mediaItem.ID,
		Type:            mediaItem.Type,
		IsFavorite:      isFavorite,
		UserRating:      userRating,
		PlayCount:       playCount,
		PlayedAt:        time.Now().Add(-24 * time.Hour),
		LastPlayedAt:    time.Now(),
		PositionSeconds: 3600,
		DurationSeconds: 7200,
	}
	userData.UUID = "series-data-" + uuid.New().String()

	err := h.DB.Table("user_media_item_data").Create(userData).Error
	require.NoError(t, err, "Failed to create user media data")
	return userData
}

// WithTransaction runs a function within a transaction and automatically rolls back
// Used to keep tests isolated and maintain a clean database state
func (h *TestDBHelper) WithTransaction(t *testing.T, fn func(tx *gorm.DB)) {
	tx := h.DB.Begin()
	require.NoError(t, tx.Error, "Failed to begin transaction")

	defer func() {
		// Always rollback test transactions to keep the database clean
		err := tx.Rollback().Error
		require.NoError(t, err, "Failed to rollback transaction")
	}()

	fn(tx)
}

// SeedInitialData adds basic test data to the database
func (h *TestDBHelper) SeedInitialData(t *testing.T) {
	// Create a few test users
	h.CreateTestUsers(t, 3)

	// Add other initial data as needed for tests
	// Create a test movie
	h.CreateTestMovie(t, "Seeded Test Movie", 2023)
}

// NewContextWithTestLogger creates a new context with a test logger
func NewContextWithTestLogger(t *testing.T) context.Context {
	ctx := context.Background()
	testLogger := NewTestLogger()
	return ContextWithTestLogger(ctx, testLogger)
}

