// database/db.go
package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	media "suasor/client/media/types"
	client "suasor/client/types"
	"suasor/types"
	"suasor/types/models"
)

// Initialize sets up the database connection and migrations
func Initialize(dbConfig types.DatabaseConfig) (*gorm.DB, error) {
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

		&models.MediaItem[media.Movie]{},
		&models.MediaItem[media.Series]{},
		&models.MediaItem[media.Episode]{},
		&models.MediaItem[media.Season]{},
		&models.MediaItem[media.Track]{},
		&models.MediaItem[media.Album]{},
		&models.MediaItem[media.Artist]{},
		&models.MediaItem[media.Collection]{},
		&models.MediaItem[media.Playlist]{},

		&models.Session{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return db, nil
}
