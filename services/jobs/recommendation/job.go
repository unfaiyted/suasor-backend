// Package recommendation provides implementations for recommendation jobs
package recommendation

import (
	"context"
	"fmt"
	"log"
	"strings"
	"suasor/app/container"
	apprepos "suasor/app/repository"
	"suasor/client"
	"suasor/client/ai"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
	"time"
)

// RecommendationJob creates recommendations for users based on their preferences
type RecommendationJob struct {
	ctx                context.Context
	c                  *container.Container
	jobRepo            repository.JobRepository
	userRepo           repository.UserRepository
	userConfigRepo     repository.UserConfigRepository
	recommendationRepo repository.RecommendationRepository
	clientRepos        apprepos.ClientRepositories
	itemRepos          apprepos.CoreMediaItemRepositories
	clientItemRepos    apprepos.ClientMediaItemRepositories
	dataRepos          apprepos.UserMediaDataRepositories

	// New repositories for credits and people
	clientFactories *client.ClientFactoryService
	creditRepo      repository.CreditRepository // Will be implemented in the future
	peopleRepo      repository.PersonRepository // Will be implemented in the future
}

// NewRecommendationJob creates a new recommendation job
func NewRecommendationJob(
	ctx context.Context,
	c *container.Container,
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	userConfigRepo repository.UserConfigRepository,
	recommendationRepo repository.RecommendationRepository,
	clientFactories *client.ClientFactoryService,
	creditRepo repository.CreditRepository,
	peopleRepo repository.PersonRepository,

) *RecommendationJob {
	return &RecommendationJob{
		ctx:                ctx,
		c:                  c,
		jobRepo:            container.MustGet[repository.JobRepository](c),
		userRepo:           container.MustGet[repository.UserRepository](c),
		userConfigRepo:     container.MustGet[repository.UserConfigRepository](c),
		recommendationRepo: container.MustGet[repository.RecommendationRepository](c),
		clientFactories:    container.MustGet[*client.ClientFactoryService](c),
		clientRepos:        container.MustGet[apprepos.ClientRepositories](c),
		creditRepo:         container.MustGet[repository.CreditRepository](c),
		peopleRepo:         container.MustGet[repository.PersonRepository](c),
	}
}

// getAIClient returns an AI client for the given user
// It tries to get the default AI client from the user config, or falls back to the first active AI client
func (j *RecommendationJob) getAIClient(ctx context.Context, userID uint64) (ai.AIClient, error) {
	logger := log.Logger{} // would ideally use structured logging from context

	// Get user config to check for default AI client
	config, err := j.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user config: %w", err)
	}

	// First try to get the default AI client if set
	if config.DefaultClients != nil && config.DefaultClients.AIClientID > 0 {
		// Try Claude repository first
		claudeRepo := j.clientRepos.ClaudeRepo()
		if claudeRepo != nil {
			claudeClient, err := claudeRepo.GetByID(ctx, config.DefaultClients.AIClientID)
			if err == nil && claudeClient != nil {
				// Found the default Claude client
				client, err := j.clientFactories.GetClient(ctx, claudeClient.ID, claudeClient.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using default Claude AI client ID %d for user %d", claudeClient.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}

		// Try OpenAI repository next
		openAIRepo := j.clientRepos.OpenAIRepo()
		if openAIRepo != nil {
			openAIClient, err := openAIRepo.GetByID(ctx, config.DefaultClients.AIClientID)
			if err == nil && openAIClient != nil {
				// Found the default OpenAI client
				client, err := j.clientFactories.GetClient(ctx, openAIClient.ID, openAIClient.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using default OpenAI client ID %d for user %d", openAIClient.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}

		// Try Ollama repository next
		ollamaRepo := j.clientRepos.OllamaRepo()
		if ollamaRepo != nil {
			ollamaClient, err := ollamaRepo.GetByID(ctx, config.DefaultClients.AIClientID)
			if err == nil && ollamaClient != nil {
				// Found the default Ollama client
				client, err := j.clientFactories.GetClient(ctx, ollamaClient.ID, ollamaClient.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using default Ollama client ID %d for user %d", ollamaClient.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}

		// If we get here, the default client couldn't be found or created
		logger.Printf("Default AI client ID %d for user %d not found or could not be created",
			config.DefaultClients.AIClientID, userID)
	}

	// Try Claude clients first
	claudeRepo := j.clientRepos.ClaudeRepo()
	if claudeRepo != nil {
		claudeClients, err := claudeRepo.GetByUserID(ctx, userID)
		if err == nil && len(claudeClients) > 0 {
			// Use the first active Claude client
			for _, clientConfig := range claudeClients {
				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using first available Claude client ID %d for user %d", clientConfig.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}
	}

	// Try OpenAI clients next
	openAIRepo := j.clientRepos.OpenAIRepo()
	if openAIRepo != nil {
		openAIClients, err := openAIRepo.GetByUserID(ctx, userID)
		if err == nil && len(openAIClients) > 0 {
			// Use the first active OpenAI client
			for _, clientConfig := range openAIClients {
				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using first available OpenAI client ID %d for user %d", clientConfig.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}
	}

	// Try Ollama clients next
	ollamaRepo := j.clientRepos.OllamaRepo()
	if ollamaRepo != nil {
		ollamaClients, err := ollamaRepo.GetByUserID(ctx, userID)
		if err == nil && len(ollamaClients) > 0 {
			// Use the first active Ollama client
			for _, clientConfig := range ollamaClients {
				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using first available Ollama client ID %d for user %d", clientConfig.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}
	}

	// No AI client found
	return nil, fmt.Errorf("no AI client available for user %d", userID)
}

// Name returns the job name
func (j *RecommendationJob) Name() string {
	return "system.recommendation"
}

// Schedule returns how often the job should run
func (j *RecommendationJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute implements the standard job interface
func (j *RecommendationJob) Execute(ctx context.Context) error {
	// Since this implementation needs to match the scheduler.Job interface,
	// we'll delegate to the full implementation
	return j.ExecuteWithParams(ctx, 0, 0, nil)
}

// ExecuteWithParams runs the recommendation job with the specified parameters
func (j *RecommendationJob) ExecuteWithParams(ctx context.Context, jobID uint64, jobRunID uint64, params map[string]interface{}) error {
	// Create a logger using job ID
	logger := log.Logger{}
	logger.Printf("Starting recommendation job (ID: %d, Run: %d)", jobID, jobRunID)

	// Find active users
	users, err := j.userRepo.FindAllActive(ctx)
	if err != nil {
		logger.Printf("Failed to get active users: %s", err)
		return err
	}

	logger.Printf("Processing recommendations for %d users", len(users))

	// Process each user
	for _, user := range users {
		// Add user ID to context for better logging
		userCtx := ctx

		err := j.processUserRecommendations(userCtx, jobRunID, user)
		if err != nil {
			logger.Printf("Failed to process recommendations for user %d: %s", user.ID, err)
			// Continue with next user
			continue
		}
	}

	logger.Printf("Recommendation job completed")
	return nil
}

// IsContentTypeEnabled checks if a content type is enabled in the content type filter
func (j *RecommendationJob) IsContentTypeEnabled(contentTypesFilter string, contentType string) bool {
	// If no filter is specified, all content types are enabled
	if contentTypesFilter == "" {
		return true
	}

	// Split the filter by comma and check if the content type is in the list
	contentTypes := strings.Split(contentTypesFilter, ",")
	for _, ct := range contentTypes {
		// Trim spaces and compare
		if strings.TrimSpace(ct) == contentType {
			return true
		}
	}

	return false
}

// processUserRecommendations handles generating recommendations for a single user
func (j *RecommendationJob) processUserRecommendations(ctx context.Context, jobRunID uint64, user models.User) error {
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Processing recommendations for user")

	// Get user configuration
	config, err := j.userConfigRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user configuration")
		return err
	}

	// Build user preference profile from history
	profile, err := j.buildUserPreferenceProfile(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to build user preference profile")
		return err
	}

	// Calculate some advanced metrics based on the profile
	j.calculateAdvancedMetrics(profile)

	// Generate movie recommendations
	err = j.generateMovieRecommendations(ctx, jobRunID, user, profile, config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate movie recommendations")
		// Continue with other types
	}

	// TODO: Add functions for generating other media type recommendations
	// j.generateSeriesRecommendations(ctx, jobRunID, user, profile, config)
	// j.generateMusicRecommendations(ctx, jobRunID, user, profile, config)

	log.Info().Msg("Finished processing recommendations for user")
	return nil
}

// SetupMediaSyncJob creates or updates a media sync job for a user
func (j *RecommendationJob) SetupMediaSyncJob(ctx context.Context, userID, clientID uint64, clientType, mediaType, frequency string) error {
	// Implementation would set up a media sync job
	// This is just a stub to satisfy the interface
	log.Printf("Setting up media sync job for user %d, client %d, type %s", userID, clientID, mediaType)
	return nil
}
