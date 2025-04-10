package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// NewReleaseNotificationJob identifies and notifies users about new releases
type NewReleaseNotificationJob struct {
	jobRepo        repository.JobRepository
	userRepo       repository.UserRepository
	configRepo     repository.UserConfigRepository
	movieRepo      repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo     repository.MediaItemRepository[*mediatypes.Series]
	musicRepo      repository.MediaItemRepository[*mediatypes.Track]
	historyRepo    repository.MediaPlayHistoryRepository
	metadataClient interface{} // Using interface{} to avoid import cycles
}

// NewNewReleaseNotificationJob creates a new release notification job
func NewNewReleaseNotificationJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	historyRepo repository.MediaPlayHistoryRepository,
	metadataClient interface{},
) *NewReleaseNotificationJob {
	return &NewReleaseNotificationJob{
		jobRepo:        jobRepo,
		userRepo:       userRepo,
		configRepo:     configRepo,
		movieRepo:      movieRepo,
		seriesRepo:     seriesRepo,
		musicRepo:      musicRepo,
		historyRepo:    historyRepo,
		metadataClient: metadataClient,
	}
}

// Name returns the unique name of the job
func (j *NewReleaseNotificationJob) Name() string {
	return "system.newrelease.notification"
}

// Schedule returns when the job should next run
func (j *NewReleaseNotificationJob) Schedule() time.Duration {
	// Run daily by default
	return 24 * time.Hour
}

// Execute runs the new release notification job
func (j *NewReleaseNotificationJob) Execute(ctx context.Context) error {
	log.Println("Starting new release notification job")

	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeNotification,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		Metadata:  fmt.Sprintf(`{"type":"newReleaseNotification","startTime":"%s"}`, now.Format(time.RFC3339)),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// In a full implementation, we would:
	// 1. Fetch new releases from metadata providers
	// 2. Process each user's preferences to determine relevant notifications
	// 3. Create notifications for matching new releases
	
	// For now, simulate some work with a progress update
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 25, "Searching for new releases")
	time.Sleep(100 * time.Millisecond)
	
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 50, "Processing user preferences")
	time.Sleep(100 * time.Millisecond)
	
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 75, "Creating notifications")
	time.Sleep(100 * time.Millisecond)

	// Complete the job
	notificationStats := map[string]int{
		"usersNotified":      0,
		"movieNotifications": 0,
		"seriesNotifications": 0,
		"musicNotifications": 0,
		"totalNotifications": 0,
	}
	
	statsJSON, _ := json.Marshal(notificationStats)
	j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, string(statsJSON))

	log.Println("New release notification job completed")
	return nil
}

// NewReleases holds new release information by media type
type NewReleases struct {
	movies []NewRelease
	series []NewRelease
	music  []NewRelease
}

// NewRelease holds information about a single new release
type NewRelease struct {
	ID          string
	Title       string
	ReleaseDate time.Time
	MediaType   string
	Genres      []string
	Creators    []string // Directors for movies, showrunners for series, artists for music
}

// NotificationStats holds statistics about notifications sent
type NotificationStats struct {
	movieNotifications  int
	seriesNotifications int
	musicNotifications  int
	totalNotifications  int
}

// completeJobRun finalizes a job run with status and results
func (j *NewReleaseNotificationJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, message string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, message); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// SetupNewReleaseNotificationSchedule creates or updates a notification schedule for a user
func (j *NewReleaseNotificationJob) SetupNewReleaseNotificationSchedule(ctx context.Context, userID uint64, frequency string) error {
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
		JobType:     models.JobTypeNotification,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		UserID:      &userID,
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// processUserNotifications processes notifications for a single user
// This is a stub that would be replaced with actual notification logic
func (j *NewReleaseNotificationJob) processUserNotifications(ctx context.Context, user models.User, config *models.UserConfig, newReleases NewReleases) (NotificationStats, error) {
	// This would be implemented in a full version
	return NotificationStats{}, nil
}

// fetchNewReleases fetches new releases from metadata providers
// This is a stub that would be replaced with actual implementation
func (j *NewReleaseNotificationJob) fetchNewReleases(ctx context.Context) (NewReleases, error) {
	// This would be implemented in a full version
	return NewReleases{}, nil
}

// getMetadataClient gets a metadata client for a specific provider
func (j *NewReleaseNotificationJob) getMetadataClient(ctx context.Context, providerType types.ClientType) (interface{}, error) {
	// This would be implemented in a full version
	return nil, fmt.Errorf("metadata client not implemented")
}