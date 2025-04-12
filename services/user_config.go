package services

import (
	"context"

	"fmt"
	"suasor/repository"
	"suasor/services/jobs/recommendation"
	"suasor/types/models"
	"suasor/utils"
)

// UserConfigService provides methods to interact with configuration
type UserConfigService interface {
	GetUserConfig(ctx context.Context, id uint64) (*models.UserConfig, error)
	SaveUserConfig(ctx context.Context, config models.UserConfig) error
}

type userConfigService struct {
	userConfigRepo    repository.UserConfigRepository
	jobService        JobService
	recommendationJob *recommendation.RecommendationJob
}

// NewUserConfigService creates a new configuration service
func NewUserConfigService(
	userConfigRepo repository.UserConfigRepository,
	jobService JobService,
	recommendationJob *recommendation.RecommendationJob,
) UserConfigService {
	return &userConfigService{
		userConfigRepo:    userConfigRepo,
		jobService:        jobService,
		recommendationJob: recommendationJob,
	}
}

// GetUserConfig retrieves the configuration for a specific user
func (s *userConfigService) GetUserConfig(ctx context.Context, id uint64) (*models.UserConfig, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Uint64("userId", id).Msg("Retrieving user configuration")

	// Fetch user config from repository
	config, err := s.userConfigRepo.GetUserConfig(ctx, id)
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

	// Get the existing config to check for changes
	existingConfig, err := s.userConfigRepo.GetUserConfig(ctx, config.UserID)
	if err != nil {
		log.Error().Err(err).Uint64("userId", config.UserID).Msg("Error retrieving existing user configuration")
		// Continue anyway, as we might be creating a new config
	}

	// Save user config to repository
	if err := s.userConfigRepo.SaveUserConfig(ctx, &config); err != nil {
		log.Error().Err(err).Uint64("userId", config.UserID).Msg("Error saving user configuration")
		return fmt.Errorf("error saving user configuration: %w", err)
	}

	// Check if recommendation settings have changed
	if existingConfig == nil ||
		existingConfig.RecommendationSyncEnabled != config.RecommendationSyncEnabled ||
		existingConfig.RecommendationSyncFrequency != config.RecommendationSyncFrequency {

		log.Info().
			Uint64("userId", config.UserID).
			Bool("syncEnabled", config.RecommendationSyncEnabled).
			Str("frequency", config.RecommendationSyncFrequency).
			Msg("Recommendation settings changed, updating job schedule")

		// Update the recommendation job schedule for this user
		if s.recommendationJob != nil {
			if err := s.recommendationJob.UpdateUserRecommendationSchedule(ctx, config.UserID); err != nil {
				log.Error().Err(err).Uint64("userId", config.UserID).Msg("Error updating recommendation schedule")
				// Don't fail the overall operation if this fails
			}
		}

		// If recommendations are now enabled and set to run immediately, trigger a manual run
		if config.RecommendationSyncEnabled && config.RecommendationSyncFrequency != "manual" {
			jobName := fmt.Sprintf("%s.user.%d", s.recommendationJob.Name(), config.UserID)
			if err := s.jobService.RunJobManually(ctx, jobName); err != nil {
				log.Error().Err(err).Uint64("userId", config.UserID).Msg("Error triggering initial recommendation job")
				// Don't fail the overall operation if this fails
			}
		}
	}

	log.Info().Uint64("userId", config.UserID).Msg("User configuration saved successfully")
	return nil
}
