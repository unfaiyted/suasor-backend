package jobs

import (
	"context"
	"fmt"
	"log"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
	"time"
)

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
	clientRepos map[clienttypes.MediaClientType]interface{}
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
	// Client repositories (simplified for now)
	embyRepo interface{},
	jellyfinRepo interface{},
	plexRepo interface{},
	subsonicRepo interface{},
) *WatchHistorySyncJob {
	clientRepos := map[clienttypes.MediaClientType]interface{}{
		clienttypes.MediaClientTypeEmby:     embyRepo,
		clienttypes.MediaClientTypeJellyfin: jellyfinRepo,
		clienttypes.MediaClientTypePlex:     plexRepo,
		clienttypes.MediaClientTypeSubsonic: subsonicRepo,
	}

	return &WatchHistorySyncJob{
		jobRepo:     jobRepo,
		userRepo:    userRepo,
		configRepo:  configRepo,
		historyRepo: historyRepo,
		movieRepo:   movieRepo,
		seriesRepo:  seriesRepo,
		episodeRepo: episodeRepo,
		musicRepo:   musicRepo,
		clientRepos: clientRepos,
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
		if err := j.processUserWatchHistory(ctx, user); err != nil {
			log.Printf("Error processing watch history for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Watch history sync job completed")
	return nil
}

// processUserWatchHistory syncs watch history for a single user
func (j *WatchHistorySyncJob) processUserWatchHistory(ctx context.Context, user models.User) error {
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

	// Check if syncing is enabled for the user
	if !config.NotifyOnSync {
		log.Printf("History sync not enabled for user: %s", user.Username)
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
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s","type":"history"}`, user.ID, user.Username),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Get all media clients for this user
	// For now, we'll just use some placeholder logic
	clients, err := j.getUserMediaClients(ctx, user.ID)
	if err != nil {
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Error getting media clients: %v", err))
		return err
	}

	var jobError error
	// Process each client
	for _, client := range clients {
		clientError := j.syncClientHistory(ctx, user.ID, client, jobRun.ID)
		if clientError != nil {
			log.Printf("Error syncing history from client %d: %v", client.ClientID, clientError)
			// Record the error but continue with other clients
			if jobError == nil {
				jobError = clientError
			}
		}
	}

	// Complete the job
	status := models.JobStatusCompleted
	errorMessage := ""
	if jobError != nil {
		status = models.JobStatusFailed
		errorMessage = jobError.Error()
	}

	j.completeJobRun(ctx, jobRun.ID, status, errorMessage)
	return jobError
}

// completeJobRun finalizes a job run with status and error info
func (j *WatchHistorySyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// WatchHistoryClientInfo holds basic client information for watch history sync operations
type WatchHistoryClientInfo struct {
	ClientID   uint64
	ClientType clienttypes.MediaClientType
	Name       string
	Config     interface{}
}

// getUserMediaClients returns all media clients for a user
// This is a placeholder implementation
func (j *WatchHistorySyncJob) getUserMediaClients(ctx context.Context, userID uint64) ([]WatchHistoryClientInfo, error) {
	// For a real implementation, you would:
	// 1. Query each client repository for clients belonging to this user
	// 2. Filter to only include clients that support watch history
	// 3. Return the compiled list

	// For now, return a placeholder client
	return []WatchHistoryClientInfo{
		{
			ClientID:   1,
			ClientType: clienttypes.MediaClientTypeEmby,
			Name:       "Test Emby Server",
			Config:     nil,
		},
	}, nil
}

// syncClientHistory syncs watch history for a specific client
func (j *WatchHistorySyncJob) syncClientHistory(ctx context.Context, userID uint64, client WatchHistoryClientInfo, jobRunID uint64) error {
	log.Printf("Syncing history for user %d from client %d (%s)", userID, client.ClientID, client.Name)

	// Since this is a simplified implementation, we'll just log what would happen
	log.Printf("Would fetch watch history from %s client", client.ClientType)
	log.Printf("Would process each history item and store in our database")
	log.Printf("Would update existing history entries if they already exist")

	// Real implementation would:
	// 1. Create or get a client connection
	// 2. Fetch watch history from the client API
	// 3. Process each history item:
	//    - Find or create the corresponding media item in our database
	//    - Create or update the history record with watch date, completion status, etc.
	// 4. Update sync status

	switch client.ClientType {
	case clienttypes.MediaClientTypeEmby:
		return j.syncEmbyHistory(ctx, userID, client, jobRunID)
	case clienttypes.MediaClientTypeJellyfin:
		return j.syncJellyfinHistory(ctx, userID, client, jobRunID)
	case clienttypes.MediaClientTypePlex:
		return j.syncPlexHistory(ctx, userID, client, jobRunID)
	case clienttypes.MediaClientTypeSubsonic:
		return j.syncSubsonicHistory(ctx, userID, client, jobRunID)
	default:
		return fmt.Errorf("unsupported client type: %s", client.ClientType)
	}
}

// syncEmbyHistory syncs watch history from an Emby server
func (j *WatchHistorySyncJob) syncEmbyHistory(ctx context.Context, userID uint64, client WatchHistoryClientInfo, jobRunID uint64) error {
	// Placeholder implementation
	log.Printf("Syncing Emby history for user %d from client %d", userID, client.ClientID)

	// In a real implementation:
	// 1. Get an authenticated Emby client
	// 2. Call the Emby API to get playback history
	// 3. Process each history item

	return nil
}

// syncJellyfinHistory syncs watch history from a Jellyfin server
func (j *WatchHistorySyncJob) syncJellyfinHistory(ctx context.Context, userID uint64, client WatchHistoryClientInfo, jobRunID uint64) error {
	// Placeholder implementation
	log.Printf("Syncing Jellyfin history for user %d from client %d", userID, client.ClientID)

	// In a real implementation:
	// 1. Get an authenticated Jellyfin client
	// 2. Call the Jellyfin API to get playback history
	// 3. Process each history item

	return nil
}

// syncPlexHistory syncs watch history from a Plex server
func (j *WatchHistorySyncJob) syncPlexHistory(ctx context.Context, userID uint64, client WatchHistoryClientInfo, jobRunID uint64) error {
	// Placeholder implementation
	log.Printf("Syncing Plex history for user %d from client %d", userID, client.ClientID)

	// In a real implementation:
	// 1. Get an authenticated Plex client
	// 2. Call the Plex API to get watch history
	// 3. Process each history item

	return nil
}

// syncSubsonicHistory syncs played history from a Subsonic server
func (j *WatchHistorySyncJob) syncSubsonicHistory(ctx context.Context, userID uint64, client WatchHistoryClientInfo, jobRunID uint64) error {
	// Placeholder implementation
	log.Printf("Syncing Subsonic play history for user %d from client %d", userID, client.ClientID)

	// In a real implementation:
	// 1. Get an authenticated Subsonic client
	// 2. Call the Subsonic API to get play history
	// 3. Process each history item

	return nil
}

// ProcessHistoryItem processes a single history item from a client
// This is a generic approach that would be implemented for each client
func (j *WatchHistorySyncJob) ProcessHistoryItem(ctx context.Context, userID uint64, clientID uint64, clientType clienttypes.MediaClientType, historyItem interface{}) error {
	// This would be a common processing pipeline:
	// 1. Extract data from the client-specific history item
	// 2. Find or create the media item in our database
	// 3. Create a MediaPlayHistory record linking the user to the media item
	// 4. Store completion status, play date, position, etc.

	// For now, just log what would happen
	log.Printf("Would process history item for user %d from client %d (%s)", userID, clientID, clientType)

	return nil
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

// RunManualSync runs the history sync job manually for a specific user
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
	return j.processUserWatchHistory(ctx, *user)
}
