package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"suasor/client"
	"suasor/client/media"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// MediaClientInfo holds information about a media client for syncing
type MediaClientInfo struct {
	ClientID   uint64
	ClientType types.MediaClientType
	Name       string
	UserID     uint64
}

// FavoritesSyncJob synchronizes favorite/liked media from external clients
type FavoritesSyncJob struct {
	jobRepo       repository.JobRepository
	userRepo      repository.UserRepository
	configRepo    repository.UserConfigRepository
	movieRepo     repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo    repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo   repository.MediaItemRepository[*mediatypes.Episode]
	musicRepo     repository.MediaItemRepository[*mediatypes.Track]
	clientRepo    interface{} // Generic client repository
	clientFactory *client.ClientFactoryService
}

// NewFavoritesSyncJob creates a new favorites sync job
func NewFavoritesSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	episodeRepo repository.MediaItemRepository[*mediatypes.Episode],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	clientRepo interface{}, // Generic client repository
	clientFactory *client.ClientFactoryService,
) *FavoritesSyncJob {
	return &FavoritesSyncJob{
		jobRepo:       jobRepo,
		userRepo:      userRepo,
		configRepo:    configRepo,
		movieRepo:     movieRepo,
		seriesRepo:    seriesRepo,
		episodeRepo:   episodeRepo,
		musicRepo:     musicRepo,
		clientRepo:    clientRepo,
		clientFactory: clientFactory,
	}
}

// Name returns the unique name of the job
func (j *FavoritesSyncJob) Name() string {
	return "system.favorites.sync"
}

// Schedule returns when the job should next run
func (j *FavoritesSyncJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute runs the favorites sync job
func (j *FavoritesSyncJob) Execute(ctx context.Context) error {
	log.Println("Starting favorites sync job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.processUserFavorites(ctx, user); err != nil {
			log.Printf("Error processing favorites for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Favorites sync job completed")
	return nil
}

// processUserFavorites syncs favorites for a single user
func (j *FavoritesSyncJob) processUserFavorites(ctx context.Context, user models.User) error {
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

	// Check if favorites sync is enabled for the user
	// First check if sync notifications are enabled as a proxy for sync being enabled
	if !config.NotifyOnSync {
		log.Printf("Favorites sync not enabled for user: %s", user.Username)
		return nil
	}

	// Create a job run record for this user
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSync,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &user.ID,
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s","type":"favorites"}`, user.ID, user.Username),
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
		
		// Sync favorites for this client
		err := j.syncClientFavorites(ctx, user.ID, client, jobRun.ID)
		if err != nil {
			log.Printf("Error syncing favorites for client %s: %v", client.Name, err)
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
func (j *FavoritesSyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// SetupFavoritesSyncSchedule creates or updates a favorites sync schedule for a user
func (j *FavoritesSyncJob) SetupFavoritesSyncSchedule(ctx context.Context, userID uint64, frequency string) error {
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

// RunManualSync runs the favorites sync job manually for a specific user
func (j *FavoritesSyncJob) RunManualSync(ctx context.Context, userID uint64) error {
	// Get the user
	user, err := j.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %d", userID)
	}

	// Run the sync job for this user
	return j.processUserFavorites(ctx, *user)
}

// getUserMediaClients returns all media clients for a user
func (j *FavoritesSyncJob) getUserMediaClients(ctx context.Context, userID uint64) ([]MediaClientInfo, error) {
	// Get clients from the database
	// In a real implementation, you would query the client repository
	// For now we'll return an empty slice
	return []MediaClientInfo{}, nil
}

// syncClientFavorites syncs favorites for a specific client
func (j *FavoritesSyncJob) syncClientFavorites(ctx context.Context, userID uint64, client MediaClientInfo, jobRunID uint64) error {
	// Log the start of synchronization
	log.Printf("Syncing favorites for user %d from client %d", userID, client.ClientID)
	
	// Get the client using client factory
	mediaClient, err := j.getMediaClient(ctx, client.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get media client: %w", err)
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, "Fetching favorites from client")
	
	// This will depend on the client implementation
	// Here we'll show an example for movies
	favoriteMovies, err := j.getFavoriteMovies(ctx, mediaClient)
	if err != nil {
		return fmt.Errorf("failed to get favorite movies: %w", err)
	}
	
	// Update movies with the favorite flag
	for _, movie := range favoriteMovies {
		existingMovie, err := j.movieRepo.GetByExternalID(ctx, movie.ExternalID, client.ClientID)
		if err != nil {
			// Skip if movie not found
			log.Printf("Movie not found in database: %s", movie.Data.Details.Title)
			continue
		}
		
		// Update favorite flag
		existingMovie.Data.Details.IsFavorite = true
		_, err = j.movieRepo.Update(ctx, *existingMovie)
		if err != nil {
			log.Printf("Error updating favorite flag for movie %s: %v", existingMovie.Data.Details.Title, err)
			// Continue with other movies
			continue
		}
		
		log.Printf("Updated favorite status for movie: %s", existingMovie.Data.Details.Title)
	}
	
	// Similarly, implement for series, episodes, and music tracks
	
	// Updated progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 90, "Finalizing favorites sync")
	
	return nil
}

// getMediaClient gets a media client from the client factory
func (j *FavoritesSyncJob) getMediaClient(ctx context.Context, clientID uint64) (media.MediaClient, error) {
	// In a real implementation, we would get the client from the database
	// For now, we'll simulate the client
	
	// Currently config doesn't implement ClientConfig, so we'll use a placeholder
	// until proper implementation is available
	// For now, return an error to indicate this is not fully implemented
	return nil, fmt.Errorf("media client implementation not completed")
}

// getFavoriteMovies gets favorite/liked movies from a client
func (j *FavoritesSyncJob) getFavoriteMovies(ctx context.Context, mediaClient media.MediaClient) ([]models.MediaItem[*mediatypes.Movie], error) {
	// In a full implementation, we would cast client to MovieProvider and call GetMovies
	
	// For now, return an empty slice as a placeholder
	return []models.MediaItem[*mediatypes.Movie]{}, nil
}

// Helper to serialize metadata to JSON for favorites
func makeFavoriteMetadataJson(data map[string]interface{}) string {
	if data == nil {
		return "{}"
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}

	return string(jsonBytes)
}