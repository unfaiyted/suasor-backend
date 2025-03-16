// database/db.go
package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"suasor/models"
)

// Config holds database configuration
type Config struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

// Initialize sets up the database connection and migrations
func Initialize(dbConfig Config) (*gorm.DB, error) {
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
	if err := db.AutoMigrate(&models.User{}, &models.Session{}, &models.Shorten{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return db, nil
}
