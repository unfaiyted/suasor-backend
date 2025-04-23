package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	media "suasor/clients/media/types"
	client "suasor/clients/types"
	"suasor/types/models"
)

// InitializeInMemoryDB creates an in-memory SQLite database for testing
// It performs the same migrations as the production database
func InitializeInMemoryDB(ctx context.Context) (*gorm.DB, error) {
	// Use in-memory SQLite database with a persistent connection
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		// Silence logs in tests unless explicitly configured otherwise
		Logger: logger.Default.LogMode(logger.Silent),
		// This is needed for generic types and better SQLite compatibility
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open in-memory database: %w", err)
	}

	// SQLite pragmas for better performance in tests
	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")
	db.Exec("PRAGMA synchronous = NORMAL")

	// Auto migrate the schema for specific models
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserConfig{},
		&models.Session{},
		&models.Client[*client.EmbyConfig]{},
		&models.Client[*client.JellyfinConfig]{},
		&models.Client[*client.PlexConfig]{},
		&models.Client[*client.SubsonicConfig]{},
		&models.Client[*client.LidarrConfig]{},
		&models.Client[*client.RadarrConfig]{},
		&models.Client[*client.SonarrConfig]{},
		&models.MediaItem[*media.Movie]{},
		&models.MediaItem[*media.Series]{},
		&models.MediaItem[*media.Episode]{},
		&models.MediaItem[*media.Season]{},
		&models.MediaItem[*media.Track]{},
		&models.MediaItem[*media.Album]{},
		&models.MediaItem[*media.Artist]{},
		&models.MediaItem[*media.Collection]{},
		&models.MediaItem[*media.Playlist]{},
		// We explicitly create the user_media_item_data table with SQL instead of migrations
		&models.JobSchedule{},
		&models.JobRun{},
		&models.Recommendation{},
		&models.MediaSyncJob{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	// Create user_media_item_data table explicitly for SQLite compatibility
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS user_media_item_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE,
		user_id INTEGER NOT NULL,
		media_item_id INTEGER NOT NULL,
		type TEXT NOT NULL,
		is_favorite BOOLEAN DEFAULT false,
		is_disliked BOOLEAN DEFAULT false,
		watchlist BOOLEAN DEFAULT false,
		user_rating REAL DEFAULT 0,
		play_count INTEGER DEFAULT 0,
		position_seconds INTEGER DEFAULT 0,
		duration_seconds INTEGER DEFAULT 0,
		played_percentage REAL DEFAULT 0,
		completed BOOLEAN DEFAULT false,
		played_at DATETIME,
		last_played_at DATETIME,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME
	);
	`
	if err := db.Exec(createTableSQL).Error; err != nil {
		return nil, fmt.Errorf("failed to create user_media_item_data table: %w", err)
	}

	// Create admin user for testing if running in test environment
	if err := CreateTestAdminUser(ctx, db); err != nil {
		// Log the error but don't fail the test setup
		log := TestLoggerFromContext(ctx)
		log.Warn().Err(err).Msg("Failed to create test admin user")
	}

	return db, nil
}

// InitializeTestDB creates an in-memory database populated with test data
func InitializeTestDB(ctx context.Context) (*gorm.DB, error) {
	db, err := InitializeInMemoryDB(ctx)
	if err != nil {
		return nil, err
	}

	// Seed with basic test data
	if err := seedTestData(db); err != nil {
		return nil, fmt.Errorf("failed to seed test data: %w", err)
	}

	return db, nil
}

// seedTestData adds standard test data to the database
func seedTestData(db *gorm.DB) error {
	// Create test users
	testUsers := []models.User{
		{
			Username: "testUser",
			Email:    "test@example.com",
			Role:     "user",
		},
		{
			Username: "testAdmin",
			Email:    "admin@example.com",
			Role:     "admin",
		},
	}

	for i := range testUsers {
		if err := testUsers[i].SetPassword("TestPassword123"); err != nil {
			return fmt.Errorf("failed to set password for test user: %w", err)
		}
		if err := db.Create(&testUsers[i]).Error; err != nil {
			return fmt.Errorf("failed to create test user: %w", err)
		}
	}

	// Create a test movie item using a simplified approach for SQLite
	releaseDate := time.Now().AddDate(-1, 0, 0)
	now := time.Now()
	nowMinus := time.Now().Add(-30 * 24 * time.Hour)

	// For SQLite in testing, we need to pre-marshal complex fields with proper JSON structure
	movieData := `{"Data":{"Details":{"Title":"Test Movie","Description":"This is a test movie for testing","ReleaseYear":2023,"ImdbID":"tt9999999"}}}`
	
	if err := db.Exec("INSERT INTO media_items (type, title, release_year, release_date, created_at, updated_at, data) VALUES (?, ?, ?, ?, ?, ?, ?)",
		media.MediaTypeMovie, "Test Movie", 2023, releaseDate, nowMinus, now, movieData).Error; err != nil {
		return fmt.Errorf("failed to create test movie: %w", err)
	}

	// Create a test series with proper JSON structure
	seriesJsonData := `{"Data":{"Details":{"Title":"Test Series","Description":"This is a test series for testing","ReleaseYear":2022,"ImdbID":"tt8888888"}}}`
	
	if err := db.Exec("INSERT INTO media_items (type, title, release_year, release_date, created_at, updated_at, data) VALUES (?, ?, ?, ?, ?, ?, ?)",
		media.MediaTypeSeries, "Test Series", 2022, releaseDate.AddDate(-1, 0, 0), nowMinus, now, seriesJsonData).Error; err != nil {
		return fmt.Errorf("failed to create test series: %w", err)
	}

	// Get the ID of the inserted movie
	var movieID uint64
	if err := db.Raw("SELECT id FROM media_items WHERE title = ? AND type = ?", "Test Movie", media.MediaTypeMovie).Scan(&movieID).Error; err != nil {
		return fmt.Errorf("failed to retrieve test movie ID: %w", err)
	}

	// Get the ID of the inserted series
	var seriesID uint64
	if err := db.Raw("SELECT id FROM media_items WHERE title = ? AND type = ?", "Test Series", media.MediaTypeSeries).Scan(&seriesID).Error; err != nil {
		return fmt.Errorf("failed to retrieve test series ID: %w", err)
	}

	// Create user media data for movie
	userData := models.UserMediaItemData[*media.Movie]{
		UserID:           testUsers[0].ID,
		MediaItemID:      movieID,
		Type:             media.MediaTypeMovie,
		IsFavorite:       true,
		UserRating:       4.5,
		PlayCount:        2,
		PlayedAt:         time.Now().Add(-24 * time.Hour),
		LastPlayedAt:     time.Now(),
		PositionSeconds:  3600,
		DurationSeconds:  7200,
		PlayedPercentage: 50.0,
	}
	userData.CreatedAt = nowMinus
	userData.UpdatedAt = now
	userData.UUID = "movie-" + uuid.New().String()

	if err := db.Table("user_media_item_data").Create(&userData).Error; err != nil {
		return fmt.Errorf("failed to create user media data for movie: %w", err)
	}

	// Create user media data for series
	seriesData := models.UserMediaItemData[*media.Series]{
		UserID:           testUsers[0].ID,
		MediaItemID:      seriesID,
		Type:             media.MediaTypeSeries,
		IsFavorite:       false,
		UserRating:       3.5,
		PlayCount:        1,
		PlayedAt:         time.Now().Add(-48 * time.Hour),
		LastPlayedAt:     time.Now().Add(-48 * time.Hour),
		PositionSeconds:  1800,
		DurationSeconds:  2700,
		PlayedPercentage: 66.7,
	}
	seriesData.CreatedAt = nowMinus
	seriesData.UpdatedAt = now
	seriesData.UUID = "series-" + uuid.New().String()

	if err := db.Table("user_media_item_data").Create(&seriesData).Error; err != nil {
		return fmt.Errorf("failed to create user media data for series: %w", err)
	}

	// Create a test collection with the actual IDs and proper JSON structure
	collectionData := fmt.Sprintf(`{"Data":{"Details":{"Title":"Test Collection","Description":"This is a test collection for testing"},"Items":[{"ID":%d,"Title":"Test Movie","Type":"movie"},{"ID":%d,"Title":"Test Series","Type":"series"}]}}`, 
		movieID, seriesID)
		
	if err := db.Exec("INSERT INTO media_items (type, title, release_year, release_date, created_at, updated_at, data) VALUES (?, ?, ?, ?, ?, ?, ?)",
		media.MediaTypeCollection, "Test Collection", 2023, releaseDate, nowMinus, now, collectionData).Error; err != nil {
		return fmt.Errorf("failed to create test collection: %w", err)
	}

	// Create a client - adding user_id and category since they're required in the schema
	if err := db.Exec("INSERT INTO clients (name, type, created_at, updated_at, config, user_id, category) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"Test Emby", "emby", nowMinus, now,
		`{"BaseURL":"http://localhost:8096","APIKey":"testkey123","Username":"testuser","Password":"testpassword"}`,
		testUsers[0].ID, "media").Error; err != nil {
		return fmt.Errorf("failed to create test client: %w", err)
	}

	return nil
}

// CleanupInMemoryDB removes all data from tables but keeps the schema
// Useful for tests that need a fresh database but don't want to recreate schema
func CleanupInMemoryDB(db *gorm.DB) error {
	// Get all table names
	var tables []string
	rows, err := db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Rows()
	if err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		// Skip sqlite_sequence which is used for autoincrement
		if name != "sqlite_sequence" && name != "sqlite_master" {
			tables = append(tables, name)
		}
	}

	// Disable foreign key constraints temporarily
	if err := db.Exec("PRAGMA foreign_keys = OFF").Error; err != nil {
		return fmt.Errorf("failed to disable foreign keys: %w", err)
	}

	// Begin transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Delete all data from tables
	for _, table := range tables {
		if err := tx.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Re-enable foreign key constraints
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return fmt.Errorf("failed to re-enable foreign keys: %w", err)
	}

	return nil
}