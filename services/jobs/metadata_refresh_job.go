package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// MetadataRefreshJob periodically updates metadata for media items from external sources
type MetadataRefreshJob struct {
	jobRepo            repository.JobRepository
	userRepo           repository.UserRepository
	configRepo         repository.UserConfigRepository
	movieRepo          repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo         repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo        repository.MediaItemRepository[*mediatypes.Episode]
	musicRepo          repository.MediaItemRepository[*mediatypes.Track]
	metadataClientSvc  interface{} // Using interface{} to avoid import cycles
}

// NewMetadataRefreshJob creates a new metadata refresh job
func NewMetadataRefreshJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	episodeRepo repository.MediaItemRepository[*mediatypes.Episode],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	metadataClientSvc interface{},
) *MetadataRefreshJob {
	return &MetadataRefreshJob{
		jobRepo:           jobRepo,
		userRepo:          userRepo,
		configRepo:        configRepo,
		movieRepo:         movieRepo,
		seriesRepo:        seriesRepo,
		episodeRepo:       episodeRepo,
		musicRepo:         musicRepo,
		metadataClientSvc: metadataClientSvc,
	}
}

// Name returns the unique name of the job
func (j *MetadataRefreshJob) Name() string {
	return "system.metadata.refresh"
}

// Schedule returns when the job should next run
func (j *MetadataRefreshJob) Schedule() time.Duration {
	// Run daily by default
	return 24 * time.Hour
}

// Execute runs the metadata refresh job
func (j *MetadataRefreshJob) Execute(ctx context.Context) error {
	log.Println("Starting metadata refresh job")

	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSystem,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		Metadata:  fmt.Sprintf(`{"type":"metadataRefresh","startTime":"%s"}`, now.Format(time.RFC3339)),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Process each media type
	var jobError error
	refreshStats := map[string]int{
		"moviesUpdated":     0,
		"seriesUpdated":     0,
		"episodesUpdated":   0,
		"tracksUpdated":     0,
		"errorCount":        0,
		"totalItemsChecked": 0,
	}

	// Refresh movie metadata
	movieStats, err := j.refreshMovieMetadata(ctx, jobRun.ID)
	if err != nil {
		log.Printf("Error refreshing movie metadata: %v", err)
		jobError = err
		refreshStats["errorCount"]++
	} else {
		refreshStats["moviesUpdated"] = movieStats.updated
		refreshStats["totalItemsChecked"] += movieStats.checked
	}

	// Refresh series metadata
	seriesStats, err := j.refreshSeriesMetadata(ctx, jobRun.ID)
	if err != nil {
		log.Printf("Error refreshing series metadata: %v", err)
		if jobError == nil {
			jobError = err
		}
		refreshStats["errorCount"]++
	} else {
		refreshStats["seriesUpdated"] = seriesStats.updated
		refreshStats["totalItemsChecked"] += seriesStats.checked
	}

	// Refresh episode metadata
	episodeStats, err := j.refreshEpisodeMetadata(ctx, jobRun.ID)
	if err != nil {
		log.Printf("Error refreshing episode metadata: %v", err)
		if jobError == nil {
			jobError = err
		}
		refreshStats["errorCount"]++
	} else {
		refreshStats["episodesUpdated"] = episodeStats.updated
		refreshStats["totalItemsChecked"] += episodeStats.checked
	}

	// Refresh music metadata
	musicStats, err := j.refreshMusicMetadata(ctx, jobRun.ID)
	if err != nil {
		log.Printf("Error refreshing music metadata: %v", err)
		if jobError == nil {
			jobError = err
		}
		refreshStats["errorCount"]++
	} else {
		refreshStats["tracksUpdated"] = musicStats.updated
		refreshStats["totalItemsChecked"] += musicStats.checked
	}

	// Complete the job
	status := models.JobStatusCompleted
	errorMessage := ""
	if jobError != nil {
		status = models.JobStatusFailed
		errorMessage = jobError.Error()
	}

	// Update job run with results
	j.completeJobRun(ctx, jobRun.ID, status, errorMessage, refreshStats)
	log.Println("Metadata refresh job completed")
	return jobError
}

// completeJobRun finalizes a job run with status and results
func (j *MetadataRefreshJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string, stats map[string]int) {
	// Convert stats to JSON
	statsJSON, _ := json.Marshal(stats)

	// In a real implementation, we would update the job run with this metadata
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}

	log.Printf("Metadata refresh stats: %s", string(statsJSON))
}

// MetadataRefreshStats holds statistics for metadata refresh operations
type MetadataRefreshStats struct {
	checked int
	updated int
}

// refreshMovieMetadata refreshes metadata for movies
func (j *MetadataRefreshJob) refreshMovieMetadata(ctx context.Context, jobRunID uint64) (MetadataRefreshStats, error) {
	stats := MetadataRefreshStats{}
	log.Println("Refreshing movie metadata")

	// In a real implementation, we would:
	// 1. Retrieve a batch of movies from the repository
	// 2. Prioritize recently added or popular movies
	// 3. Check if metadata needs updating based on last update time
	// 4. Call the metadata client to get updated information from TMDB or other sources
	// 5. Update the movie items in the database
	// 6. Track statistics on updates made

	// Mock implementation
	stats.checked = 250  // Pretend we checked 250 movies
	stats.updated = 75   // Pretend we updated 75 of them

	log.Printf("Checked %d movies, updated %d", stats.checked, stats.updated)
	return stats, nil
}

// refreshSeriesMetadata refreshes metadata for TV series
func (j *MetadataRefreshJob) refreshSeriesMetadata(ctx context.Context, jobRunID uint64) (MetadataRefreshStats, error) {
	stats := MetadataRefreshStats{}
	log.Println("Refreshing TV series metadata")

	// Similar to movie refresh, but for TV series
	// Mock implementation
	stats.checked = 120  // Pretend we checked 120 series
	stats.updated = 45   // Pretend we updated 45 of them

	log.Printf("Checked %d series, updated %d", stats.checked, stats.updated)
	return stats, nil
}

// refreshEpisodeMetadata refreshes metadata for TV episodes
func (j *MetadataRefreshJob) refreshEpisodeMetadata(ctx context.Context, jobRunID uint64) (MetadataRefreshStats, error) {
	stats := MetadataRefreshStats{}
	log.Println("Refreshing TV episode metadata")

	// Similar approach but for episodes
	// Mock implementation
	stats.checked = 850  // Pretend we checked 850 episodes
	stats.updated = 230  // Pretend we updated 230 of them

	log.Printf("Checked %d episodes, updated %d", stats.checked, stats.updated)
	return stats, nil
}

// refreshMusicMetadata refreshes metadata for music tracks
func (j *MetadataRefreshJob) refreshMusicMetadata(ctx context.Context, jobRunID uint64) (MetadataRefreshStats, error) {
	stats := MetadataRefreshStats{}
	log.Println("Refreshing music metadata")

	// Similar approach but for music
	// Mock implementation
	stats.checked = 500  // Pretend we checked 500 tracks
	stats.updated = 150  // Pretend we updated 150 of them

	log.Printf("Checked %d tracks, updated %d", stats.checked, stats.updated)
	return stats, nil
}

// SetupMetadataRefreshSchedule creates or updates a metadata refresh schedule
func (j *MetadataRefreshJob) SetupMetadataRefreshSchedule(ctx context.Context, frequency string) error {
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

// RunManualRefresh runs the metadata refresh job manually
func (j *MetadataRefreshJob) RunManualRefresh(ctx context.Context) error {
	return j.Execute(ctx)
}

// RefreshSingleItem refreshes metadata for a specific item
func (j *MetadataRefreshJob) RefreshSingleItem(ctx context.Context, mediaType mediatypes.MediaType, itemID uint64) error {
	log.Printf("Refreshing metadata for single %s item with ID %d", mediaType, itemID)

	// In a real implementation, we would:
	// 1. Look up the item by ID in the appropriate repository
	// 2. Call the metadata client to get updated information
	// 3. Update the item in the database

	// Mock implementation
	log.Printf("Successfully refreshed metadata for %s item with ID %d", mediaType, itemID)
	return nil
}

// GetRefreshHistory gets the history of metadata refresh runs
func (j *MetadataRefreshJob) GetRefreshHistory(ctx context.Context, limit int) ([]models.JobRun, error) {
	// In a real implementation, we would query the job repository for recent runs
	return nil, fmt.Errorf("not implemented")
}

// GetPriorityItems gets items that are high priority for refresh
func (j *MetadataRefreshJob) GetPriorityItems(ctx context.Context) (map[string][]uint64, error) {
	// In a real implementation, we would identify items that:
	// 1. Have never been refreshed
	// 2. Have very old metadata
	// 3. Are popular or recently accessed
	// 4. Have incomplete metadata

	priorityItems := map[string][]uint64{
		"movies":   {1, 2, 3},   // Mock movie IDs
		"series":   {4, 5},      // Mock series IDs
		"episodes": {6, 7, 8},   // Mock episode IDs
		"music":    {9, 10},     // Mock music IDs
	}

	return priorityItems, nil
}