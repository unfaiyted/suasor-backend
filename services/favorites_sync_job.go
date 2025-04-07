package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// MediaClientInfo holds information about a media client for syncing
type MediaClientInfo struct {
	ClientID   uint64
	ClientType clienttypes.MediaClientType
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
	clientRepos   map[clienttypes.MediaClientType]interface{}
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
	embyRepo interface{},
	jellyfinRepo interface{},
	plexRepo interface{},
	subsonicRepo interface{},
	clientFactory *client.ClientFactoryService,
) *FavoritesSyncJob {
	clientRepos := map[clienttypes.MediaClientType]interface{}{
		clienttypes.MediaClientTypeEmby:     embyRepo,
		clienttypes.MediaClientTypeJellyfin: jellyfinRepo,
		clienttypes.MediaClientTypePlex:     plexRepo,
		clienttypes.MediaClientTypeSubsonic: subsonicRepo,
	}

	return &FavoritesSyncJob{
		jobRepo:       jobRepo,
		userRepo:      userRepo,
		configRepo:    configRepo,
		movieRepo:     movieRepo,
		seriesRepo:    seriesRepo,
		episodeRepo:   episodeRepo,
		musicRepo:     musicRepo,
		clientRepos:   clientRepos,
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

	// For now, log what we would sync
	log.Printf("Would sync favorites for user %s (ID: %d)", user.Username, user.ID)

	// Complete the job
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

