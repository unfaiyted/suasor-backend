package jobs

import (
	"context"
	"fmt"
	"log"
	"suasor/client/media"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
	"time"
)

// MediaClientInfo is defined in other job files

// WatchHistorySyncJob synchronizes watched media history from external clients
type WatchHistorySyncJob struct {
	jobRepo     repository.JobRepository
	userRepo    repository.UserRepository
	configRepo  repository.UserConfigRepository
	historyRepo repository.MediaPlayHistoryRepository
	movieRepo   repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo  repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo repository.MediaItemRepository[*mediatypes.Episode]
	musicRepo   repository.MediaItemRepository[*mediatypes.Track]
	clientRepo  repository.ClientRepository[*types.EmbyConfig] // Representative client repository
}

// NewWatchHistorySyncJob creates a new watch history sync job
func NewWatchHistorySyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	historyRepo repository.MediaPlayHistoryRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	episodeRepo repository.MediaItemRepository[*mediatypes.Episode],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	clientRepo repository.ClientRepository[*types.EmbyConfig], // Representative client repository
) *WatchHistorySyncJob {
	return &WatchHistorySyncJob{
		jobRepo:     jobRepo,
		userRepo:    userRepo,
		configRepo:  configRepo,
		historyRepo: historyRepo,
		movieRepo:   movieRepo,
		seriesRepo:  seriesRepo,
		episodeRepo: episodeRepo,
		musicRepo:   musicRepo,
		clientRepo:  clientRepo,
	}
}

// Name returns the unique name of the job
func (j *WatchHistorySyncJob) Name() string {
	return "system.history.sync"
}

// Schedule returns when the job should next run
func (j *WatchHistorySyncJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute runs the watch history sync job
func (j *WatchHistorySyncJob) Execute(ctx context.Context) error {
	log.Println("Starting watch history sync job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.processUserHistory(ctx, user); err != nil {
			log.Printf("Error processing history for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Watch history sync job completed")
	return nil
}

// processUserHistory syncs watch history for a single user
func (j *WatchHistorySyncJob) processUserHistory(ctx context.Context, user models.User) error {
	// Skip inactive users
	if !user.Active {
		log.Printf("Skipping inactive user: %s", user.Username)
		return nil
	}

	// Get user configuration
	_, err := j.configRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	// Create a job run record for this user
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSync,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &user.ID,
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s","type":"watchHistory"}`, user.ID, user.Username),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}
	
	// Get all media clients for the user
	clients, err := j.getUserMediaClients(ctx, user.ID)
	if err != nil {
		errorMsg := fmt.Sprintf("Error getting media clients: %v", err)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return fmt.Errorf(errorMsg)
	}
	
	if len(clients) == 0 {
		j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 100, "No media clients found")
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "No media clients found")
		return nil
	}
	
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 10, fmt.Sprintf("Found %d media clients", len(clients)))
	
	// Process each client
	totalClients := len(clients)
	processedClients := 0
	var lastError error
	
	for i, client := range clients {
		// Update progress
		progress := 10 + int(float64(i)/float64(totalClients)*80.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, progress, 
			fmt.Sprintf("Processing client %d/%d: %s", i+1, totalClients, client.Name))
		
		// Sync watch history for this client
		err := j.syncClientHistory(ctx, user.ID, client, jobRun.ID)
		if err != nil {
			log.Printf("Error syncing history for client %s: %v", client.Name, err)
			lastError = err
			continue
		}
		
		processedClients++
	}
	
	// Complete the job
	if lastError != nil {
		errorMsg := fmt.Sprintf("Completed with errors: %v", lastError)
		j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 100, 
			fmt.Sprintf("Processed %d/%d clients with errors", processedClients, totalClients))
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return lastError
	}
	
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 100, 
		fmt.Sprintf("Successfully processed %d/%d clients", processedClients, totalClients))
	j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "")
	return nil
}

// completeJobRun finalizes a job run with status and error info
func (j *WatchHistorySyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// SetupWatchHistorySyncSchedule creates or updates a watch history sync schedule for a user
func (j *WatchHistorySyncJob) SetupWatchHistorySyncSchedule(ctx context.Context, userID uint64, frequency string) error {
	jobName := fmt.Sprintf("%s.user.%d", j.Name(), userID)

	// Check if job already exists
	existing, err := j.jobRepo.GetJobSchedule(ctx, jobName)
	if err != nil {
		return fmt.Errorf("error checking for existing job: %w", err)
	}

	// If job exists, update it
	if existing != nil {
		existing.Frequency = frequency
		existing.Enabled = frequency != string(scheduler.FrequencyManual)
		return j.jobRepo.UpdateJobSchedule(ctx, existing)
	}

	// Create a new job schedule
	schedule := &models.JobSchedule{
		JobName:     jobName,
		JobType:     models.JobTypeSync,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		UserID:      &userID,
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// RunManualSync runs the watch history sync job manually for a specific user
func (j *WatchHistorySyncJob) RunManualSync(ctx context.Context, userID uint64) error {
	// Get the user
	user, err := j.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %d", userID)
	}

	// Run the sync job for this user
	return j.processUserHistory(ctx, *user)
}

// getUserMediaClients returns all media clients for a user
func (j *WatchHistorySyncJob) getUserMediaClients(ctx context.Context, userID uint64) ([]MediaClientInfo, error) {
	// Get clients from the database
	// In a real implementation, you would query the client repository
	// For now, return an empty slice
	return []MediaClientInfo{}, nil
}

// syncClientHistory syncs watch history for a specific client
func (j *WatchHistorySyncJob) syncClientHistory(ctx context.Context, userID uint64, client MediaClientInfo, jobRunID uint64) error {
	// Log the start of synchronization
	log.Printf("Syncing watch history for user %d from client %d", userID, client.ClientID)
	
	// Get the client using client factory
	_, err := j.getMediaClient(ctx, client.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get media client: %w", err)
	}
	
	// TODO: When implementing the full functionality,
	// uncomment the mediaClient variable and use it to retrieve history
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, "Fetching play history from client")
	
	// In a real implementation, we would:
	// 1. Check if client supports play history
	// 2. Get play history for movies, series, music
	// 3. Process each history item and match it to our database items using ClientIDs array
	
	// Example implementation stub for processing history items:
	/*
	historyItems, err := mediaClient.GetRecentlyPlayed(ctx)
	if err != nil {
		return fmt.Errorf("failed to get play history: %w", err)
	}
	
	for _, historyItem := range historyItems {
		// For each history item, find the corresponding media item in our database
		// using the new ClientIDs array-based structure
		
		// For movies:
		for _, movieID := range historyItem.MovieIDs {
			// Find the movie in our database using clientItemID
			existingMovie, err := j.movieRepo.GetByClientItemID(ctx, movieID, client.ClientID)
			if err != nil {
				// Skip if movie not found
				continue
			}
			
			// Record the watch history for this movie
			j.historyRepo.CreateMediaPlayHistory(ctx, &models.MediaPlayHistory{
				UserID:       userID,
				MediaItemID:  existingMovie.ID,
				MediaType:    models.MediaTypeMovie,
				PlayedAt:     historyItem.PlayedAt,
				ClientID:     client.ClientID,
				ClientItemID: movieID,
			})
		}
		
		// Similar processing for series, episodes, and music
	}
	*/
	
	// Simulate some work for now
	time.Sleep(100 * time.Millisecond)
	
	// Updated progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 90, "Finalizing watch history sync")
	
	return nil
}

// getMediaClient gets a media client from the client factory
func (j *WatchHistorySyncJob) getMediaClient(ctx context.Context, clientID uint64) (media.MediaClient, error) {
	// In a real implementation, we would get the client from the database
	// For now, just return an error
	return nil, fmt.Errorf("media client implementation not completed")
}