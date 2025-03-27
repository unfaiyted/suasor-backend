// repository/config_repository.go
package repository

import (
	"context"
	"fmt"
	"suasor/types/models"

	"gorm.io/gorm"
)

// ConfigRepository handles configuration storage operations
type UserConfigRepository interface {
	GetUserConfig(ctx context.Context, userID uint64) (*models.UserConfig, error)
	SaveUserConfig(ctx context.Context, config *models.UserConfig) error
}

type userConfigRepository struct {
	configPath string
	db         *gorm.DB
}

// NewConfigRepository creates a new configuration repository
func NewUserConfigRepository(db *gorm.DB) UserConfigRepository {
	return &userConfigRepository{
		db: db,
	}
}

// GetUserConfig retrieves a user's configuration from the database
func (r *userConfigRepository) GetUserConfig(ctx context.Context, userID uint64) (*models.UserConfig, error) {
	var config models.UserConfig

	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&config)
	if result.Error != nil {
		// If the user config doesn't exist yet, return a new default config
		if result.Error == gorm.ErrRecordNotFound {
			return &models.UserConfig{
				UserID:                 userID,
				Theme:                  "system",
				Language:               "en-US",
				ItemsPerPage:           20,
				EnableAnimations:       true,
				PreferredGenres:        "",
				ExcludedGenres:         "",
				ContentLanguages:       "en",
				RecommendationStrategy: "balanced",
				NotificationsEnabled:   true,
			}, nil
		}
		return nil, fmt.Errorf("error retrieving user config: %w", result.Error)
	}

	return &config, nil
}

// SaveUserConfig creates or updates a user's configuration in the database
func (r *userConfigRepository) SaveUserConfig(ctx context.Context, config *models.UserConfig) error {
	// Check if the config already exists
	var existingConfig models.UserConfig
	result := r.db.WithContext(ctx).Where("user_id = ?", config.UserID).First(&existingConfig)

	// If config exists, update it
	if result.Error == nil {
		config.ID = existingConfig.ID
		config.CreatedAt = existingConfig.CreatedAt
		result = r.db.WithContext(ctx).Save(config)
		if result.Error != nil {
			return fmt.Errorf("error updating user config: %w", result.Error)
		}
		return nil
	}

	// If config doesn't exist or there was an error other than record not found
	if result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("error checking for existing user config: %w", result.Error)
	}

	// Create new config
	result = r.db.WithContext(ctx).Create(config)
	if result.Error != nil {
		return fmt.Errorf("error creating user config: %w", result.Error)
	}

	return nil
}
