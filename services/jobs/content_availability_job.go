package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	mediatypes "suasor/clients/media/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// ContentAvailabilityJob monitors the availability of content across streaming services
type ContentAvailabilityJob struct {
	jobRepo        repository.JobRepository
	userRepo       repository.UserRepository
	configRepo     repository.UserConfigRepository
	movieRepo      repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo     repository.MediaItemRepository[*mediatypes.Series]
	metadataClient interface{} // Using interface{} to avoid import cycles
}

// NewContentAvailabilityJob creates a new content availability monitoring job
func NewContentAvailabilityJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	metadataClient interface{},
) *ContentAvailabilityJob {
	return &ContentAvailabilityJob{
		jobRepo:        jobRepo,
		userRepo:       userRepo,
		configRepo:     configRepo,
		movieRepo:      movieRepo,
		seriesRepo:     seriesRepo,
		metadataClient: metadataClient,
	}
}

// Name returns the unique name of the job
func (j *ContentAvailabilityJob) Name() string {
	return "system.content.availability"
}

// Schedule returns when the job should next run
func (j *ContentAvailabilityJob) Schedule() time.Duration {
	// Run daily by default
	return 24 * time.Hour
}

// Execute runs the content availability monitoring job
func (j *ContentAvailabilityJob) Execute(ctx context.Context) error {
	log.Println("Starting content availability monitoring job")

	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSystem,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		Metadata:  fmt.Sprintf(`{"type":"contentAvailability","startTime":"%s"}`, now.Format(time.RFC3339)),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		msg := fmt.Sprintf("Error getting users: %v", err)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, msg)
		return fmt.Errorf(msg)
	}

	var jobError error
	var totalChanges int

	// Process each user
	for _, user := range users {
		if !user.Active {
			continue
		}

		// Check if the user has enabled content availability monitoring
		config, err := j.configRepo.GetUserConfig(ctx, user.ID)
		if err != nil {
			log.Printf("Error getting config for user %d: %v", user.ID, err)
			continue
		}

		if !config.ContentAvailabilityEnabled {
			continue
		}

		changes, err := j.checkContentForUser(ctx, user, jobRun.ID)
		if err != nil {
			log.Printf("Error checking content for user %d: %v", user.ID, err)
			if jobError == nil {
				jobError = err
			}
		}

		totalChanges += changes
	}

	// Complete the job run
	status := models.JobStatusCompleted
	message := fmt.Sprintf("Processed content availability. Found %d changes.", totalChanges)
	if jobError != nil {
		status = models.JobStatusFailed
		message = jobError.Error()
	}

	j.completeJobRun(ctx, jobRun.ID, status, message)
	log.Println("Content availability monitoring job completed")
	return jobError
}

// completeJobRun finalizes a job run with status and message
func (j *ContentAvailabilityJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, message string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, message); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// checkContentForUser checks content availability for a specific user
func (j *ContentAvailabilityJob) checkContentForUser(ctx context.Context, user models.User, jobRunID uint64) (int, error) {
	log.Printf("Checking content availability for user %s", user.Username)

	// In a real implementation, we would:
	// 1. Get the user's watchlist or saved content
	// 2. Check for updates in availability across streaming services
	// 3. Send notifications about changes
	// 4. Update the availability data in our database

	// Mock implementation
	watchlist, err := j.getUserWatchlist(ctx, user.ID)
	if err != nil {
		return 0, fmt.Errorf("error getting user watchlist: %w", err)
	}

	totalChanges := 0

	// Check each item in the watchlist
	for _, item := range watchlist {
		changes, err := j.checkContentAvailability(ctx, item, user.ID)
		if err != nil {
			log.Printf("Error checking availability for %s '%s': %v",
				item.Type, item.Title, err)
			continue
		}

		if changes.New > 0 || changes.Removed > 0 {
			j.notifyUserOfChanges(ctx, user.ID, item, changes)
			totalChanges++
		}
	}

	log.Printf("Found %d availability changes for user %s", totalChanges, user.Username)
	return totalChanges, nil
}

// getUserWatchlist gets a user's watchlist items
func (j *ContentAvailabilityJob) getUserWatchlist(ctx context.Context, userID uint64) ([]WatchlistItem, error) {
	// In a real implementation, we would:
	// 1. Query the database for the user's saved/watchlisted content
	// 2. Format it into a list of items to check

	// Mock implementation
	return []WatchlistItem{
		{ID: 1, Type: "movie", Title: "The Matrix", Year: 1999, TMDB_ID: 603},
		{ID: 2, Type: "movie", Title: "Inception", Year: 2010, TMDB_ID: 27205},
		{ID: 3, Type: "series", Title: "Stranger Things", Year: 2016, TMDB_ID: 66732},
		{ID: 4, Type: "series", Title: "Breaking Bad", Year: 2008, TMDB_ID: 1396},
	}, nil
}

// checkContentAvailability checks availability of a specific content item
func (j *ContentAvailabilityJob) checkContentAvailability(ctx context.Context, item WatchlistItem, userID uint64) (AvailabilityChanges, error) {
	changes := AvailabilityChanges{}

	// In a real implementation, we would:
	// 1. Use the metadata client to check current availability
	// 2. Compare with previous availability data
	// 3. Identify new and removed streaming services

	// Mock implementation - randomly simulate changes for testing
	itemID := item.ID % 4
	if itemID == 0 {
		// New service added
		changes.New = 1
		changes.NewServices = []string{"HBO Max"}
	} else if itemID == 1 {
		// Service removed
		changes.Removed = 1
		changes.RemovedServices = []string{"Netflix"}
	} else if itemID == 2 {
		// Both added and removed
		changes.New = 1
		changes.Removed = 1
		changes.NewServices = []string{"Disney+"}
		changes.RemovedServices = []string{"Hulu"}
	}

	// Update the stored availability information
	j.updateStoredAvailability(ctx, item, userID)

	return changes, nil
}

// updateStoredAvailability updates the stored availability information
func (j *ContentAvailabilityJob) updateStoredAvailability(ctx context.Context, item WatchlistItem, userID uint64) {
	// In a real implementation, we would:
	// 1. Update the database with the current availability information
	// 2. Track historical availability for trend analysis

	log.Printf("Updated availability information for %s '%s'", item.Type, item.Title)
}

// notifyUserOfChanges sends a notification about availability changes
func (j *ContentAvailabilityJob) notifyUserOfChanges(ctx context.Context, userID uint64, item WatchlistItem, changes AvailabilityChanges) {
	// In a real implementation, we would:
	// 1. Create a notification in the database
	// 2. Send an email or push notification if configured
	// 3. Include details about what changed

	var message string
	if changes.New > 0 && changes.Removed > 0 {
		message = fmt.Sprintf("%s '%s' is now available on %s but no longer on %s",
			item.Type, item.Title,
			changes.NewServices[0], changes.RemovedServices[0])
	} else if changes.New > 0 {
		message = fmt.Sprintf("%s '%s' is now available on %s",
			item.Type, item.Title, changes.NewServices[0])
	} else if changes.Removed > 0 {
		message = fmt.Sprintf("%s '%s' is no longer available on %s",
			item.Type, item.Title, changes.RemovedServices[0])
	}

	log.Printf("Would notify user %d: %s", userID, message)
}

// SetupContentAvailabilitySchedule creates or updates a content availability monitoring schedule
func (j *ContentAvailabilityJob) SetupContentAvailabilitySchedule(ctx context.Context, frequency string) error {
	// Check if job already exists
	existing, err := j.jobRepo.GetJobSchedule(ctx, j.Name())
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
		JobName:     j.Name(),
		JobType:     models.JobTypeSystem,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// RunManualAvailabilityCheck runs the content availability job manually
func (j *ContentAvailabilityJob) RunManualAvailabilityCheck(ctx context.Context) error {
	return j.Execute(ctx)
}

// CheckSingleItem checks availability for a single content item
func (j *ContentAvailabilityJob) CheckSingleItem(ctx context.Context, mediaType string, tmdbID int, userID uint64) (map[string]interface{}, error) {
	log.Printf("Checking availability for %s with TMDB ID %d", mediaType, tmdbID)

	// In a real implementation, we would:
	// 1. Use the metadata client to get current availability
	// 2. Format the results

	// Mock implementation
	services := []map[string]interface{}{
		{
			"name":   "Netflix",
			"region": "US",
			"type":   "subscription",
			"link":   "https://www.netflix.com/title/12345",
		},
		{
			"name":   "Amazon Prime",
			"region": "US",
			"type":   "rent",
			"price":  "$3.99",
			"link":   "https://www.amazon.com/gp/video/detail/12345",
		},
	}

	title := "Example Series"
	if mediaType == "movie" {
		title = "Example Movie"
	}

	result := map[string]interface{}{
		"mediaType":   mediaType,
		"tmdbID":      tmdbID,
		"title":       title,
		"availableOn": services,
		"lastUpdated": time.Now().Format(time.RFC3339),
	}

	return result, nil
}

// GetUserAvailabilityReport gets an availability report for a user's content
func (j *ContentAvailabilityJob) GetUserAvailabilityReport(ctx context.Context, userID uint64) (map[string]interface{}, error) {
	// In a real implementation, we would:
	// 1. Get the user's watchlist
	// 2. For each item, get the current availability
	// 3. Compile into a comprehensive report

	// Mock implementation
	return map[string]interface{}{
		"userId":    userID,
		"itemCount": 25,
		"availabilityByService": map[string]int{
			"Netflix":      12,
			"Hulu":         8,
			"Amazon Prime": 15,
			"Disney+":      5,
			"HBO Max":      10,
		},
		"recentChanges": []map[string]interface{}{
			{
				"title":   "The Matrix",
				"type":    "movie",
				"added":   []string{"HBO Max"},
				"removed": []string{},
				"date":    time.Now().AddDate(0, 0, -2).Format(time.RFC3339),
			},
			{
				"title":   "Stranger Things",
				"type":    "series",
				"added":   []string{},
				"removed": []string{"Hulu"},
				"date":    time.Now().AddDate(0, 0, -5).Format(time.RFC3339),
			},
		},
	}, nil
}
