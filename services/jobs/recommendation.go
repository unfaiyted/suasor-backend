package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	mediatypes "suasor/client/media/types"
	clientTypes "suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
	"time"
)

// RecommendationJob creates recommendations for users based on their preferences
type RecommendationJob struct {
	jobRepo         repository.JobRepository
	userRepo        repository.UserRepository
	configRepo      repository.UserConfigRepository
	movieRepo       repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo      repository.MediaItemRepository[*mediatypes.Series]
	musicRepo       repository.MediaItemRepository[*mediatypes.Track]
	historyRepo     repository.MediaPlayHistoryRepository
	aiClientService interface{} // Using interface{} to avoid import cycles
}

// NewRecommendationJob creates a new recommendation job
func NewRecommendationJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	historyRepo repository.MediaPlayHistoryRepository,
	aiClientService interface{},
) *RecommendationJob {
	return &RecommendationJob{
		jobRepo:         jobRepo,
		userRepo:        userRepo,
		configRepo:      configRepo,
		movieRepo:       movieRepo,
		seriesRepo:      seriesRepo,
		musicRepo:       musicRepo,
		historyRepo:     historyRepo,
		aiClientService: aiClientService,
	}
}

// Name returns the unique name of the job
func (j *RecommendationJob) Name() string {
	return "system.recommendation"
}

// Schedule returns when the job should next run
func (j *RecommendationJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute runs the recommendation job
func (j *RecommendationJob) Execute(ctx context.Context) error {
	log.Println("Starting recommendation job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.processUserRecommendations(ctx, user); err != nil {
			log.Printf("Error processing recommendations for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Recommendation job completed")
	return nil
}

// processUserRecommendations generates recommendations for a single user
func (j *RecommendationJob) processUserRecommendations(ctx context.Context, user models.User) error {
	// Skip inactive users
	if !user.Active {
		log.Printf("Skipping inactive user: %s", user.Username)
		return nil
	}

	// Get user configuration
	config, err := j.configRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	// Check if user has automated recommendations enabled
	if !config.RecommendationSyncEnabled {
		log.Printf("Recommendation sync not enabled for user: %s", user.Username)
		return nil
	}

	// Check if it's time to generate recommendations based on frequency
	shouldRun := j.shouldRunForUser(ctx, user.ID, config.RecommendationSyncFrequency)
	if !shouldRun {
		log.Printf("Not time to run recommendations for user: %s", user.Username)
		return nil
	}

	log.Printf("Generating recommendations for user: %s", user.Username)

	// Create a job run record for this user
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeRecommendation,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &user.ID,
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s"}`, user.ID, user.Username),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Process recommendations for different content types
	var jobError error

	// Generate movie recommendations if enabled
	if j.isContentTypeEnabled(config.RecommendationContentTypes, "movie") {
		if err := j.generateMovieRecommendations(ctx, user, config, jobRun.ID); err != nil {
			log.Printf("Error generating movie recommendations: %v", err)
			jobError = err
		}
	}

	// Generate TV show recommendations if enabled
	if j.isContentTypeEnabled(config.RecommendationContentTypes, "series") {
		if err := j.generateSeriesRecommendations(ctx, user, config, jobRun.ID); err != nil {
			log.Printf("Error generating series recommendations: %v", err)
			if jobError == nil {
				jobError = err
			}
		}
	}

	// Generate music recommendations if enabled
	if j.isContentTypeEnabled(config.RecommendationContentTypes, "music") {
		if err := j.generateMusicRecommendations(ctx, user, config, jobRun.ID); err != nil {
			log.Printf("Error generating music recommendations: %v", err)
			if jobError == nil {
				jobError = err
			}
		}
	}

	// Set job status based on outcome
	status := models.JobStatusCompleted
	errorMessage := ""
	if jobError != nil {
		status = models.JobStatusFailed
		errorMessage = jobError.Error()
	}

	if err := j.jobRepo.CompleteJobRun(ctx, jobRun.ID, status, errorMessage); err != nil {
		log.Printf("Error completing job run: %v", err)
	}

	// Update the job schedule's last run time
	jobName := fmt.Sprintf("%s.user.%d", j.Name(), user.ID)
	if err := j.jobRepo.UpdateJobLastRunTime(ctx, jobName, now); err != nil {
		log.Printf("Error updating job last run time: %v", err)
	}

	return jobError
}

// shouldRunForUser determines if recommendations should be generated for a user
func (j *RecommendationJob) shouldRunForUser(ctx context.Context, userID uint64, frequency string) bool {
	// Convert to scheduler.Frequency
	freq := scheduler.Frequency(frequency)

	// Manual frequency means never auto-run
	if freq == scheduler.FrequencyManual {
		return false
	}

	// Get the last run time for this user
	jobName := fmt.Sprintf("%s.user.%d", j.Name(), userID)
	schedule, err := j.jobRepo.GetJobSchedule(ctx, jobName)
	if err != nil {
		log.Printf("Error getting job schedule for user %d: %v", userID, err)
		// If we can't get the schedule, assume it should run
		return true
	}

	// If no schedule exists or it has never run, it should run
	if schedule == nil || schedule.LastRunTime == nil {
		return true
	}

	// Check if enough time has passed since the last run
	return freq.ShouldRunNow(*schedule.LastRunTime)
}

// isContentTypeEnabled checks if a content type is enabled in the comma-separated list
func (j *RecommendationJob) isContentTypeEnabled(contentTypes string, contentType string) bool {
	// TODO: Implement this properly with string parsing
	// For now, always return true for testing
	return true
}

// generateMovieRecommendations creates movie recommendations for a user
func (j *RecommendationJob) generateMovieRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating movie recommendations for user %s", user.Username)

	// For now, just create a few placeholder recommendations
	// In a real implementation, we would use AI to generate personalized recommendations
	movieTitles := []string{"The Matrix", "Inception", "Interstellar"}
	movieYears := []int{1999, 2010, 2014}

	for i, title := range movieTitles {
		// Create placeholder movie item
		mediaItem, err := j.findOrCreateMovieItem(ctx, title, movieYears[i], user.ID)
		if err != nil {
			log.Printf("Error creating movie item: %v", err)
			continue
		}

		// Create recommendation
		recommendation := &models.Recommendation{
			UserID:      user.ID,
			MediaItemID: mediaItem.ID,
			MediaType:   "movie",
			Source:      models.RecommendationSourceSystem,
			Reason:      "This is a placeholder recommendation",
			Confidence:  0.9,
			InLibrary:   false,
			Viewed:      false,
			Active:      true,
			JobRunID:    &jobRunID,
			Metadata:    "{}",
		}

		// Save recommendation
		if err := j.jobRepo.CreateRecommendation(ctx, recommendation); err != nil {
			log.Printf("Error creating recommendation: %v", err)
			continue
		}
	}

	return nil
}

// generateSystemMovieRecommendations creates movie recommendations using system algorithms
func (j *RecommendationJob) generateSystemMovieRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating system movie recommendations for user %s", user.Username)
	return j.generateMovieRecommendations(ctx, user, config, jobRunID)
}

// generateSeriesRecommendations creates TV series recommendations for a user
func (j *RecommendationJob) generateSeriesRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating series recommendations for user %s", user.Username)
	// TODO: Implement series recommendations similar to movies
	return nil
}

// generateMusicRecommendations creates music recommendations for a user
func (j *RecommendationJob) generateMusicRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating music recommendations for user %s", user.Username)
	// TODO: Implement music recommendations similar to movies
	return nil
}

// getRecentUserMovies gets a list of recent movies the user has watched
func (j *RecommendationJob) getRecentUserMovies(ctx context.Context, userID uint64, limit int) ([]map[string]interface{}, error) {
	// TODO: Implement this by querying the media history repository
	// For now, return an empty list
	return []map[string]interface{}{}, nil
}

// findOrCreateMovieItem finds a movie in the database or creates it if it doesn't exist
func (j *RecommendationJob) findOrCreateMovieItem(ctx context.Context, title string, year int, userID uint64) (*models.MediaItem[*mediatypes.Movie], error) {
	// Create a placeholder movie with the correct structure
	details := mediatypes.MediaDetails{
		Title:       title,
		ReleaseYear: year,
		Description: "Placeholder description",
	}

	movie := &mediatypes.Movie{
		Details: details,
	}

	mediaItem := models.MediaItem[*mediatypes.Movie]{
		ExternalID: fmt.Sprintf("placeholder-%s-%d", title, year),
		ClientID:   0,                               // Placeholder - would come from the actual media client
		ClientType: clientTypes.MediaClientTypePlex, // Placeholder
		Type:       mediatypes.MediaTypeMovie,
		Data:       movie,
	}

	// In a real implementation, we would:
	// 1. Look for matching movies in the database
	// 2. If not found, create the movie
	// 3. Associate it with the appropriate client

	// For now, just return the mock item with ID 1
	mediaItem.ID = 1
	return &mediaItem, nil
}

// hasUserViewed checks if a user has already viewed/played a media item
func (j *RecommendationJob) hasUserViewed(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
	// TODO: Implement by checking media play history
	// For testing purposes, return false
	return false, nil
}

// isInUserLibrary checks if a media item is in the user's library
func (j *RecommendationJob) isInUserLibrary(ctx context.Context, userID, mediaItemID uint64, mediaType string) (bool, error) {
	// TODO: Implement by checking if the media item is in any of the user's media clients
	// For testing purposes, return false
	return false, nil
}

// Helper to build filters for AI movie recommendations
func buildMovieFilters(config *models.UserConfig, recentMovies []map[string]interface{}) map[string]interface{} {
	filters := map[string]interface{}{}

	// Add preferred genres if available
	if config.PreferredGenres != nil && len(config.PreferredGenres.Movies) > 0 {
		filters["preferredGenres"] = config.PreferredGenres.Movies
	}

	// Add excluded genres if available
	if config.ExcludedGenres != nil && len(config.ExcludedGenres.Movies) > 0 {
		filters["excludedGenres"] = config.ExcludedGenres.Movies
	}

	// Add content rating preferences
	if config.MinContentRating != "" {
		filters["minContentRating"] = config.MinContentRating
	}
	if config.MaxContentRating != "" {
		filters["maxContentRating"] = config.MaxContentRating
	}

	// Add age preference
	if config.RecommendationMaxAge > 0 {
		filters["maxYearsOld"] = config.RecommendationMaxAge
	}

	// Add excluded keywords
	if config.ExcludedKeywords != "" {
		filters["excludedKeywords"] = config.ExcludedKeywords
	}

	// Add recommendation strategy
	if config.RecommendationStrategy != "" {
		filters["strategy"] = config.RecommendationStrategy
	}

	// Add recently watched movies
	if len(recentMovies) > 0 {
		filters["recentlyWatched"] = recentMovies
	}

	return filters
}

// Helper to serialize metadata to JSON
func makeMetadataJson(data map[string]interface{}) string {
	if data == nil {
		return "{}"
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}

	return string(jsonBytes)
}

// UpdateUserRecommendationSchedule updates a user's recommendation schedule based on their config
func (j *RecommendationJob) UpdateUserRecommendationSchedule(ctx context.Context, userID uint64) error {
	// Get user configuration
	config, err := j.configRepo.GetUserConfig(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	frequency := config.RecommendationSyncFrequency

	// Create or update the user's job schedule
	return j.createUserRecommendationJob(ctx, userID, frequency)
}

// createUserRecommendationJob creates or updates a job schedule for a user
func (j *RecommendationJob) createUserRecommendationJob(ctx context.Context, userID uint64, frequency string) error {
	jobName := fmt.Sprintf("%s.user.%d", j.Name(), userID)

	// Check if the job already exists
	existing, err := j.jobRepo.GetJobSchedule(ctx, jobName)
	if err != nil {
		return fmt.Errorf("error checking for existing job: %w", err)
	}

	// If the job exists, update it
	if existing != nil {
		existing.Frequency = frequency
		existing.Enabled = frequency != string(scheduler.FrequencyManual)
		return j.jobRepo.UpdateJobSchedule(ctx, existing)
	}

	// Create a new job schedule
	schedule := &models.JobSchedule{
		JobName:     jobName,
		JobType:     models.JobTypeRecommendation,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		UserID:      &userID,
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// SetupMediaSyncJob creates or updates a media sync job for a user and client
func (j *RecommendationJob) SetupMediaSyncJob(ctx context.Context, userID, clientID uint64, clientType string, mediaType string, frequency string) error {
	// Check if sync job already exists
	syncJobs, err := j.jobRepo.GetMediaSyncJobsByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error checking for existing sync jobs: %w", err)
	}

	// Look for matching job
	var existingJob *models.MediaSyncJob
	for i := range syncJobs {
		if syncJobs[i].ClientID == clientID && syncJobs[i].MediaType == mediaType {
			existingJob = &syncJobs[i]
			break
		}
	}

	// If job exists, update it
	if existingJob != nil {
		existingJob.Frequency = frequency
		existingJob.Enabled = frequency != string(scheduler.FrequencyManual)
		return j.jobRepo.UpdateMediaSyncJob(ctx, existingJob)
	}

	// Create new sync job
	syncJob := &models.MediaSyncJob{
		UserID:     userID,
		ClientID:   clientID,
		ClientType: clientType,
		MediaType:  mediaType,
		Frequency:  frequency,
		Enabled:    frequency != string(scheduler.FrequencyManual),
	}

	return j.jobRepo.CreateMediaSyncJob(ctx, syncJob)
}
