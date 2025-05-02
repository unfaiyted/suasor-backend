package sync

import (
	"context"
	"fmt"
	"log"
	"suasor/clients"
	"suasor/clients/media"
	clienttypes "suasor/clients/types"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	servicebundles "suasor/services/bundles"
	"suasor/services/scheduler"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"
)

// MediaSyncJob handles syncing of media items from clients
type MediaSyncJob struct {
	jobRepo             repository.JobRepository
	userRepo            repository.UserRepository
	userConfigRepo      repository.UserConfigRepository
	clientRepos         repobundles.ClientRepositories
	dataRepos           repobundles.UserMediaDataRepositories
	clientItemRepos     repobundles.ClientMediaItemRepositories
	clientMusicServices servicebundles.ClientMusicServices
	itemRepos           repobundles.UserMediaItemRepositories
	clientFactories     *clients.ClientProviderFactoryService
}

// NewMediaSyncJob creates a new media sync job
func NewMediaSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	userConfigRepo repository.UserConfigRepository,
	clientRepos repobundles.ClientRepositories,
	dataRepos repobundles.UserMediaDataRepositories,
	clientItemRepos repobundles.ClientMediaItemRepositories,
	clientMusicServices servicebundles.ClientMusicServices,
	itemRepos repobundles.UserMediaItemRepositories,
	clientFactories *clients.ClientProviderFactoryService,
) *MediaSyncJob {
	return &MediaSyncJob{
		jobRepo:             jobRepo,
		userRepo:            userRepo,
		userConfigRepo:      userConfigRepo,
		clientRepos:         clientRepos,
		dataRepos:           dataRepos,
		clientItemRepos:     clientItemRepos,
		clientMusicServices: clientMusicServices,
		itemRepos:           itemRepos,
		clientFactories:     clientFactories,
	}
}

// Schedule returns how often the job should run (default)
func (j *MediaSyncJob) Schedule() time.Duration {
	// Default to checking hourly
	return 1 * time.Hour
}

// Name returns the unique name of the job
func (j *MediaSyncJob) Name() string {
	// Make sure we always return a valid name even if struct is empty
	if j == nil || j.jobRepo == nil {
		return "system.media.sync"
	}
	return "system.media.sync"
}

// Execute runs the job
func (j *MediaSyncJob) Execute(ctx context.Context) error {
	log.Println("Starting media sync job")

	// Check if job is properly initialized
	if j == nil || j.jobRepo == nil {
		log.Println("MediaSyncJob not properly initialized, using stub implementation")
		log.Println("Media sync job completed (no-op)")
		return nil
	}

	// Check if any sync jobs are scheduled and due
	syncJobs, err := j.jobRepo.GetMediaSyncJobsByUser(ctx, 0) // Get all sync jobs
	if err != nil {
		return fmt.Errorf("failed to get media sync jobs: %w", err)
	}

	// Process each sync job
	for _, syncJob := range syncJobs {
		// Check if this job is enabled and due to run
		if !syncJob.Enabled || !j.isDue(syncJob) {
			continue
		}

		// Run the sync job
		err := j.runSyncJob(ctx, syncJob)
		if err != nil {
			log.Printf("Error running sync job %d: %v", syncJob.ID, err)
			// Continue with other jobs even if one fails
			continue
		}

		// Update last run time
		now := time.Now()
		syncJob.LastSyncTime = &now
		err = j.jobRepo.UpdateMediaSyncJob(ctx, &syncJob)
		if err != nil {
			log.Printf("Error updating sync job last run time: %v", err)
		}
	}

	log.Println("Media sync job completed")
	return nil
}

// RunManualSync runs a manual sync for a specific user and client
func (j *MediaSyncJob) RunManualSync(ctx context.Context, userID uint64, clientID uint64, syncType models.SyncType) error {
	// First, we need to determine the client type
	// We'll need to get it from the client record in the database

	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("syncType", string(syncType)).
		Msg("Running manual sync job")

	// Validate input parameters
	if userID == 0 {
		return fmt.Errorf("invalid user ID: cannot be zero")
	}

	if clientID == 0 {
		return fmt.Errorf("invalid client ID: cannot be zero")
	}

	if !syncType.IsValid() {
		return fmt.Errorf("invalid sync type value: %s", syncType)
	}

	// Get a list of possible client types to check
	clientTypes := []clienttypes.ClientType{
		clienttypes.ClientTypeEmby,
		clienttypes.ClientTypeJellyfin,
		clienttypes.ClientTypePlex,
		clienttypes.ClientTypeSubsonic,
	}

	// Try to find which type this client is
	var clientType clienttypes.ClientType
	for _, cType := range clientTypes {
		// Try to get the config for this client type
		config, err := j.getClientConfig(ctx, clientID)
		if err == nil && config != nil {
			// Found the client type
			log.Info().
				Uint64("clientID", clientID).
				Str("clientType", string(cType)).
				Msg("Found client type")
			clientType = cType
			break
		}
	}

	// If we couldn't determine the client type, return an error
	if clientType == "" {
		return fmt.Errorf("couldn't determine client type for clientID=%d", clientID)
	}

	log.Info().
		Str("clientType", string(clientType)).
		Msg("Determined client type for manual sync")

	// Create a temporary sync job
	syncJob := models.MediaSyncJob{
		UserID:     userID,
		ClientID:   clientID,
		ClientType: clientType,
		SyncType:   syncType,
	}

	// Run the sync job
	return j.runSyncJob(ctx, syncJob)
}

// SyncUserMediaFromClient runs a sync job for a specific user and client
// This is an alias for RunManualSync for backward compatibility
func (j *MediaSyncJob) SyncUserMediaFromClient(ctx context.Context, userID uint64, clientID uint64, syncType models.SyncType) error {
	return j.RunManualSync(ctx, userID, clientID, syncType)
}

// completeJobRun marks a job run as completed with the given status and error message
func (j *MediaSyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

func (j *MediaSyncJob) getClientConfig(ctx context.Context, clientID uint64) (clienttypes.ClientConfig, error) {
	log := logger.LoggerFromContext(ctx)

	// Validate input parameters
	if clientID == 0 {
		return nil, fmt.Errorf("invalid client ID: cannot be zero")
	}

	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieving client config from database")

	// Get client config from database
	clientList, err := j.clientRepos.GetAllMediaClients(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get media clients: %w", err)
	}

	config := clientList.GetClientConfig(clientID)

	// Validate that config is not nil
	if config == nil {
		return nil, fmt.Errorf("retrieved nil config for clientID=%d", clientID)
	}

	return config, nil
}

// getClientMedia gets a media client from the database and initializes it
func (j *MediaSyncJob) getClientMedia(ctx context.Context, clientID uint64) (media.ClientMedia, error) {
	log := logger.LoggerFromContext(ctx)

	// Validate input parameters
	if clientID == 0 {
		return nil, fmt.Errorf("invalid client ID: cannot be zero")
	}

	log.Info().
		Uint64("clientID", clientID).
		Msg("Getting client media")

	// Use the type to get the client config by id
	clientConfig, err := j.getClientConfig(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client config: %w", err)
	}

	// Validate client config is not nil before proceeding
	if clientConfig == nil {
		return nil, fmt.Errorf("client config is nil for clientID=%d, clientType=%s", clientID)
	}

	// Cast media client from generic client
	client, err := j.clientFactories.GetClient(ctx, clientID, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Cast to media client
	clientMedia, ok := client.(media.ClientMedia)
	if !ok {
		return nil, fmt.Errorf("client is not a media client")
	}

	return clientMedia, nil
}

// isDue checks if a sync job is due to run
func (j *MediaSyncJob) isDue(job models.MediaSyncJob) bool {
	// If no last sync time, it's always due
	if job.LastSyncTime == nil || job.LastSyncTime.IsZero() {
		return true
	}

	// Parse the frequency
	var duration time.Duration
	freq := scheduler.Frequency(job.Frequency)

	switch freq {
	case scheduler.FrequencyDaily:
		duration = 24 * time.Hour
	case scheduler.FrequencyWeekly:
		duration = 7 * 24 * time.Hour
	case scheduler.FrequencyMonthly:
		// Approximate month as 30 days
		duration = 30 * 24 * time.Hour
	case scheduler.FrequencyManual:
		// Manual jobs are never automatically due
		return false
	default:
		// Default to daily if frequency is unknown
		duration = 24 * time.Hour
	}

	// Check if enough time has passed since the last run
	return time.Since(*job.LastSyncTime) >= duration
}

// runSyncJob executes a media sync job
func (j *MediaSyncJob) runSyncJob(ctx context.Context, syncJob models.MediaSyncJob) error {
	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   fmt.Sprintf("%s.%s", j.Name(), syncJob.SyncType),
		JobType:   models.JobTypeSync,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &syncJob.UserID,
		Metadata:  fmt.Sprintf(`{"clientID":%d,"mediaType":"%s"}`, syncJob.ClientID, syncJob.SyncType),
	}

	// Save the job run
	err := j.jobRepo.CreateJobRun(ctx, jobRun)
	if err != nil {
		return fmt.Errorf("failed to create job run: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 0, "Starting media sync")

	// Get the client from the database
	clientMedia, err := j.getClientMedia(ctx, syncJob.ClientID)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get media client: %v", err)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Process different media types
	var syncError error

	// Normalize media type to handle both singular and plural forms
	switch syncJob.SyncType {
	case models.SyncTypeMovies:
		syncError = j.syncMovies(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeSeries:
		syncError = j.syncSeries(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeMusic:
		syncError = j.syncMusic(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeHistory:
		syncError = j.syncHistory(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeFavorites:
		// syncError = j.syncFavorites(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeCollections:
		// syncError = j.syncCollections(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypePlaylists:
		// syncError = j.syncPlaylists(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	default:
		syncError = fmt.Errorf("unsupported media type: %s", syncJob.SyncType)
	}

	// Complete the job run
	status := models.JobStatusCompleted
	errorMessage := ""
	if syncError != nil {
		status = models.JobStatusFailed
		errorMessage = syncError.Error()
	}

	j.completeJobRun(ctx, jobRun.ID, status, errorMessage)
	return syncError
}
