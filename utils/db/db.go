// database/db.go
package database

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fmt"
	"os"
	media "suasor/clients/media/types"
	client "suasor/clients/types"
	"suasor/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// CreateTestAdminUser checks if this is a test environment and creates a default admin user if needed
func CreateTestAdminUser(ctx context.Context, db *gorm.DB) error {
	// Check if we're in a test environment
	log := logger.LoggerFromContext(ctx)
	isDevEnv := os.Getenv("GO_ENV") == "dev"

	log.Info().
		Str("environment", os.Getenv("GO_ENV")).
		Msg("Checking if we're in a test environment")

	if !isDevEnv {
		log.Info().
			Str("environment", os.Getenv("GO_ENV")).
			Msg("Not in a test environment, skipping admin creation")
		return nil // Not a test environment, do nothing
	}

	var count int64

	// Check if admin user already exists
	result := db.Model(&models.User{}).Where("email = ?", "admin@dev.com").Count(&count)
	if result.Error != nil {
		return fmt.Errorf("failed to check for existing admin: %w", result.Error)
	}

	if count > 0 {
		fmt.Println("Admin user already exists in test environment")
		return nil
	}

	// Create default admin user
	adminUser := models.User{
		Username: "devAdmin",
		Email:    "admin@dev.com",
		Role:     "admin",
		// Note: You should use your app's password hashing mechanism
		// Add any other required fields for your User model
	}
	adminUser.SetPassword("TestPassword123")

	if err := db.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("failed to create test admin user: %w", err)
	}

	fmt.Println("Default admin user created for test environment")
	return nil
}

// Initialize sets up the database connection and migrations
func Initialize(ctx context.Context, dbConfig types.DatabaseConfig) (*gorm.DB, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Initializing database connection and migrations")
	
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Port)

	postgresDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	var count int64

	postgresDB.Raw("SELECT count(*) from pg_database WHERE datname = ?", dbConfig.Name).Scan(&count)
	if count == 0 {
		createDBSQL := fmt.Sprintf("CREATE DATABASE %s", dbConfig.Name)
		if err := postgresDB.Exec(createDBSQL).Error; err != nil {
			return nil, fmt.Errorf("failed to create the database %s: %w", dbConfig.Name, err)
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// db.Exec("DELETE FROM clients")
	// db.Exec("DROP TABLE clients")
	// Auto Migrate the schema
	//&models.User{},
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserConfig{},
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

		&models.UserMediaItemData[*media.Movie]{},
		&models.UserMediaItemData[*media.Series]{},
		&models.UserMediaItemData[*media.Episode]{},
		&models.UserMediaItemData[*media.Season]{},
		&models.UserMediaItemData[*media.Track]{},
		&models.UserMediaItemData[*media.Album]{},
		&models.UserMediaItemData[*media.Artist]{},
		&models.UserMediaItemData[*media.Collection]{},
		&models.UserMediaItemData[*media.Playlist]{},

		&models.ListCollaborator{},

		&models.Session{},
		&models.JobSchedule{},
		&models.JobRun{},
		&models.Recommendation{},
		&models.MediaSyncJob{},
		
		// AI Conversation models
		&models.AIConversation{},
		&models.AIMessage{},
		&models.AIRecommendation{},
		&models.AIConversationAnalytics{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	if err := CreateTestAdminUser(ctx, db); err != nil {
		fmt.Printf("Warning: %v\n", err)
		// We don't return the error to avoid breaking the app initialization
	}

	log.Info().Msg("Database initialization completed successfully")
	return db, nil
}
