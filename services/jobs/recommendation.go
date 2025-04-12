// Package jobs contains job implementations for the system
// This file is a wrapper for the recommendation package to maintain backward compatibility
package jobs

// import (
// 	"suasor/client"
// 	mediatypes "suasor/client/media/types"
// 	"suasor/repository"
// 	"suasor/services/jobs/recommendation"
// )
//
// // RecommendationJob is a wrapper around the new recommendation.RecommendationJob
// type RecommendationJob struct {
// 	impl *recommendation.RecommendationJob
// }
//
// // NewRecommendationJob creates a new recommendation job
// func NewRecommendationJob(
// 	jobRepo repository.JobRepository,
// 	userRepo repository.UserRepository,
// 	configRepo repository.UserConfigRepository,
// 	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
// 	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
// 	musicRepo repository.MediaItemRepository[*mediatypes.Track],
// 	historyRepo repository.MediaPlayHistoryRepository,
// 	clientRepos repository.ClientRepositoryCollection,
// 	clientFactories *client.ClientFactoryService,
// 	// Optional repositories for credits and people - can be nil for now
// 	creditRepo repository.CreditRepository,
// 	peopleRepo repository.PersonRepository,
// 	// Optional recommendation repository - can be nil for now
// 	recommendationRepo repository.RecommendationRepository,
// ) *RecommendationJob {
// 	// Create the new implementation from the recommendation package
// 	impl := recommendation.NewRecommendationJob(
// 		jobRepo,
// 		userRepo,
// 		configRepo,
// 		movieRepo,
// 		seriesRepo,
// 		musicRepo,
// 		historyRepo,
// 		clientRepos,
// 		clientFactories,
// 		creditRepo,
// 		peopleRepo,
// 		recommendationRepo,
// 	)
//
// 	return &RecommendationJob{
// 		impl: impl,
// 	}
// }

//
// // Name returns the unique name of the job
// func (j *RecommendationJob) Name() string {
// 	return j.impl.Name()
// }
//
// // Schedule returns when the job should next run
// func (j *RecommendationJob) Schedule() time.Duration {
// 	return j.impl.Schedule()
// }
//
// // Execute implements the scheduler.Job interface
// func (j *RecommendationJob) Execute(ctx context.Context) error {
// 	return j.impl.Execute(ctx)
// }
//
// // ExecuteWithParams runs the recommendation job with parameters
// func (j *RecommendationJob) ExecuteWithParams(ctx context.Context, jobID uint64, jobRunID uint64, params map[string]interface{}) error {
// 	return j.impl.ExecuteWithParams(ctx, jobID, jobRunID, params)
// }
//
// // IsContentTypeEnabled checks if a content type is enabled in the content type filter
// func (j *RecommendationJob) IsContentTypeEnabled(contentTypesFilter string, contentType string) bool {
// 	return j.impl.IsContentTypeEnabled(contentTypesFilter, contentType)
// }
//
// // getAIClient returns an AI client for the given user
// // It tries to get the default AI client from the user config, or falls back to the first active AI client
// func (j *RecommendationJob) getAIClient(ctx context.Context, userID uint64) (ai.AIClient, error) {
// 	logger := log.Logger{} // would ideally use structured logging from context
//
// 	// Get user config to check for default AI client
// 	config, err := j.configRepo.GetUserConfig(ctx, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting user config: %w", err)
// 	}
//
// 	// First try to get the default AI client if set
// 	if config.DefaultClients != nil && config.DefaultClients.AIClientID > 0 {
// 		// Try Claude repository first
// 		claudeRepo := j.clientRepos.AllRepos().ClaudeRepo
// 		if claudeRepo != nil {
// 			claudeClient, err := claudeRepo.GetByID(ctx, config.DefaultClients.AIClientID)
// 			if err == nil && claudeClient != nil {
// 				// Found the default Claude client
// 				client, err := j.clientFactories.GetClient(ctx, claudeClient.ID, claudeClient.Config.Data)
// 				if err == nil && client != nil {
// 					logger.Printf("Using default Claude AI client ID %d for user %d", claudeClient.ID, userID)
// 					return client.(ai.AIClient), nil
// 				}
// 			}
// 		}
//
// 		// Try OpenAI repository next
// 		openAIRepo := j.clientRepos.AllRepos().OpenAIRepo
// 		if openAIRepo != nil {
// 			openAIClient, err := openAIRepo.GetByID(ctx, config.DefaultClients.AIClientID)
// 			if err == nil && openAIClient != nil {
// 				// Found the default OpenAI client
// 				client, err := j.clientFactories.GetClient(ctx, openAIClient.ID, openAIClient.Config.Data)
// 				if err == nil && client != nil {
// 					logger.Printf("Using default OpenAI client ID %d for user %d", openAIClient.ID, userID)
// 					return client.(ai.AIClient), nil
// 				}
// 			}
// 		}
//
// 		// Try Ollama repository next
// 		ollamaRepo := j.clientRepos.AllRepos().OllamaRepo
// 		if ollamaRepo != nil {
// 			ollamaClient, err := ollamaRepo.GetByID(ctx, config.DefaultClients.AIClientID)
// 			if err == nil && ollamaClient != nil {
// 				// Found the default Ollama client
// 				client, err := j.clientFactories.GetClient(ctx, ollamaClient.ID, ollamaClient.Config.Data)
// 				if err == nil && client != nil {
// 					logger.Printf("Using default Ollama client ID %d for user %d", ollamaClient.ID, userID)
// 					return client.(ai.AIClient), nil
// 				}
// 			}
// 		}
//
// 		// If we get here, the default client couldn't be found or created
// 		logger.Printf("Default AI client ID %d for user %d not found or could not be created",
// 			config.DefaultClients.AIClientID, userID)
// 	}
//
// 	// If default client not set or couldn't be loaded, try to get any AI client
//
// 	// Try Claude clients first
// 	claudeRepo := j.clientRepos.AllRepos().ClaudeRepo
// 	if claudeRepo != nil {
// 		claudeClients, err := claudeRepo.GetByUserID(ctx, userID)
// 		if err == nil && len(claudeClients) > 0 {
// 			// Use the first active Claude client
// 			for _, clientConfig := range claudeClients {
// 				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
// 				if err == nil && client != nil {
// 					logger.Printf("Using first available Claude client ID %d for user %d", clientConfig.ID, userID)
// 					return client.(ai.AIClient), nil
// 				}
// 			}
// 		}
// 	}
//
// 	// Try OpenAI clients next
// 	openAIRepo := j.clientRepos.AllRepos().OpenAIRepo
// 	if openAIRepo != nil {
// 		openAIClients, err := openAIRepo.GetByUserID(ctx, userID)
// 		if err == nil && len(openAIClients) > 0 {
// 			// Use the first active OpenAI client
// 			for _, clientConfig := range openAIClients {
// 				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
// 				if err == nil && client != nil {
// 					logger.Printf("Using first available OpenAI client ID %d for user %d", clientConfig.ID, userID)
// 					return client.(ai.AIClient), nil
// 				}
// 			}
// 		}
// 	}
//
// 	// Try Ollama clients next
// 	ollamaRepo := j.clientRepos.AllRepos().OllamaRepo
// 	if ollamaRepo != nil {
// 		ollamaClients, err := ollamaRepo.GetByUserID(ctx, userID)
// 		if err == nil && len(ollamaClients) > 0 {
// 			// Use the first active Ollama client
// 			for _, clientConfig := range ollamaClients {
// 				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
// 				if err == nil && client != nil {
// 					logger.Printf("Using first available Ollama client ID %d for user %d", clientConfig.ID, userID)
// 					return client.(ai.AIClient), nil
// 				}
// 			}
// 		}
// 	}
//
// 	// No AI client found
// 	logger.Printf("No AI clients found for user %d", userID)
// 	return nil, fmt.Errorf("no AI clients found for user %d", userID)
// }
//
// // Name returns the unique name of the job
// func (j *RecommendationJob) Name() string {
// 	return "system.recommendation"
// }
//
// // Schedule returns when the job should next run
// func (j *RecommendationJob) Schedule() time.Duration {
// 	// Default to checking daily
// 	return 24 * time.Hour
// }
//
// // Execute implements the scheduler.Job interface
// func (j *RecommendationJob) Execute(ctx context.Context) error {
// 	// Since this implementation needs to match the scheduler.Job interface,
// 	// we'll create a basic version without parameters
// 	log.Println("Executing recommendation job")
// 	return nil
// }
//
// // ExecuteWithParams runs the recommendation job with parameters
// func (j *RecommendationJob) ExecuteWithParams(ctx context.Context, jobID uint64, jobRunID uint64, params map[string]interface{}) error {
// 	ctx, jobLog := utils.WithJobID(ctx, jobID)
//
// 	jobLog.Info().
// 		Uint64("jobID", jobID).
// 		Uint64("jobRunID", jobRunID).
// 		Interface("params", params).
// 		Msg("Starting recommendation job")
//
// 	// Update job status to in-progress
// 	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 0, "Starting recommendation generation")
//
// 	// Get all active users (or a specific user if provided in params)
// 	var users []models.User
// 	var err error
//
// 	if userIDParam, ok := params["userID"]; ok {
// 		// If a specific user ID was provided
// 		userIDint, _ := strconv.ParseUint(fmt.Sprintf("%v", userIDParam), 10, 64)
// 		user, err := j.userRepo.FindByID(ctx, userIDint)
// 		if err != nil {
// 			jobLog.Error().Err(err).Msg("Failed to get user")
// 			return err
// 		}
// 		users = []models.User{*user}
// 	} else {
// 		// Get all active users
// 		users, err = j.userRepo.FindAllActive(ctx) // Active
// 		if err != nil {
// 			jobLog.Error().Err(err).Msg("Failed to get users")
// 			return err
// 		}
// 	}
//
// 	total := len(users)
// 	jobLog.Info().
// 		Int("userCount", total).
// 		Msg("Generating recommendations for users")
//
// 	// Process each user
// 	for idx, user := range users {
// 		progress := (idx * 100) / total
// 		statusMsg := fmt.Sprintf("Processing user %d/%d", idx+1, total)
// 		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, statusMsg)
//
// 		err := j.processUserRecommendations(ctx, jobRunID, user)
// 		if err != nil {
// 			jobLog.Error().
// 				Err(err).
// 				Uint64("userID", user.ID).
// 				Msg("Failed to generate recommendations for user")
// 			// Continue to the next user
// 			continue
// 		}
// 	}
//
// 	// Mark job as completed
// 	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, "Recommendation generation completed")
// 	jobLog.Info().Msg("Recommendation job completed")
//
// 	return nil
// }
//
//
