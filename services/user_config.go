package services

import (
	"context"

	"fmt"
	"suasor/models"
	"suasor/repository"
	"suasor/utils"
)

// UserConfigService provides methods to interact with configuration
type UserConfigService interface {
	GetUserConfig(ctx context.Context, id uint64) (*models.UserConfig, error)
	SaveUserConfig(cxt context.Context, config models.UserConfig) error
}

type userConfigService struct {
	userConfigRepo repository.UserConfigRepository
}

// NewUserConfigService creates a new configuration service
func NewUserConfigService(userConfigRepo repository.UserConfigRepository) UserConfigService {
	return &userConfigService{
		userConfigRepo: userConfigRepo,
	}
}

// GetUserConfig retrieves the configuration for a specific user
func (s *userConfigService) GetUserConfig(ctx context.Context, id uint64) (*models.UserConfig, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Uint64("userId", id).Msg("Retrieving user configuration")

	// Convert uint64 to uint for the user ID
	userID := id

	// Fetch user config from repository
	config, err := s.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint64("userId", id).Msg("Error retrieving user configuration")
		return nil, fmt.Errorf("error retrieving user configuration: %w", err)
	}

	log.Debug().Uint64("userId", id).Interface("config", config).Msg("User configuration retrieved")
	return config, nil
}

// SaveUserConfig creates or updates a user's configuration
func (s *userConfigService) SaveUserConfig(ctx context.Context, config models.UserConfig) error {
	log := utils.LoggerFromContext(ctx)
	log.Info().Uint64("userId", config.UserID).Msg("Saving user configuration")

	// Validate required fields
	if config.UserID == 0 {
		log.Error().Msg("UserID is required for saving user configuration")
		return fmt.Errorf("userID is required for saving user configuration")
	}

	// Save user config to repository
	if err := s.userConfigRepo.SaveUserConfig(ctx, &config); err != nil {
		log.Error().Err(err).Uint64("userId", config.UserID).Msg("Error saving user configuration")
		return fmt.Errorf("error saving user configuration: %w", err)
	}

	log.Info().Uint64("userId", config.UserID).Msg("User configuration saved successfully")
	return nil
}
