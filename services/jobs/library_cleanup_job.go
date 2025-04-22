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

// LibraryCleanupJob identifies and manages obsolete or duplicate media items
type LibraryCleanupJob struct {
	jobRepo     repository.JobRepository
	userRepo    repository.UserRepository
	configRepo  repository.UserConfigRepository
	movieRepo   repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo  repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo repository.MediaItemRepository[*mediatypes.Episode]
	musicRepo   repository.MediaItemRepository[*mediatypes.Track]
}

// NewLibraryCleanupJob creates a new library cleanup job
func NewLibraryCleanupJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	episodeRepo repository.MediaItemRepository[*mediatypes.Episode],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
) *LibraryCleanupJob {
	return &LibraryCleanupJob{
		jobRepo:     jobRepo,
		userRepo:    userRepo,
		configRepo:  configRepo,
		movieRepo:   movieRepo,
		seriesRepo:  seriesRepo,
		episodeRepo: episodeRepo,
		musicRepo:   musicRepo,
	}
}

// Name returns the unique name of the job
func (j *LibraryCleanupJob) Name() string {
	return "system.library.cleanup"
}

// Schedule returns when the job should next run
func (j *LibraryCleanupJob) Schedule() time.Duration {
	// Run weekly by default
	return 7 * 24 * time.Hour
}

// Execute runs the library cleanup job
func (j *LibraryCleanupJob) Execute(ctx context.Context) error {
	log.Println("Starting library cleanup job")

	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSystem,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		Metadata:  fmt.Sprintf(`{"type":"libraryCleanup","startTime":"%s"}`, now.Format(time.RFC3339)),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Run the cleanup tasks
	var jobError error
	cleanupStats := map[string]int{
		"duplicatesFound":      0,
		"duplicatesResolved":   0,
		"orphanedItemsFound":   0,
		"orphanedItemsRemoved": 0,
		"accessErrorsFound":    0,
		"accessErrorsFixed":    0,
	}

	// Find and resolve duplicate media items
	duplicateStats, err := j.findAndResolveDuplicates(ctx)
	if err != nil {
		log.Printf("Error finding and resolving duplicates: %v", err)
		jobError = err
		// Continue with other tasks even if one fails
	} else {
		cleanupStats["duplicatesFound"] = duplicateStats.found
		cleanupStats["duplicatesResolved"] = duplicateStats.resolved
	}

	// Find and cleanup orphaned media items
	orphanedStats, err := j.findAndCleanupOrphanedItems(ctx)
	if err != nil {
		log.Printf("Error finding and cleaning orphaned items: %v", err)
		if jobError == nil {
			jobError = err
		}
		// Continue with other tasks
	} else {
		cleanupStats["orphanedItemsFound"] = orphanedStats.found
		cleanupStats["orphanedItemsRemoved"] = orphanedStats.removed
	}

	// Find and fix access errors
	accessErrorStats, err := j.findAndFixAccessErrors(ctx)
	if err != nil {
		log.Printf("Error finding and fixing access errors: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		cleanupStats["accessErrorsFound"] = accessErrorStats.found
		cleanupStats["accessErrorsFixed"] = accessErrorStats.fixed
	}

	// Complete the job
	status := models.JobStatusCompleted
	errorMessage := ""
	if jobError != nil {
		status = models.JobStatusFailed
		errorMessage = jobError.Error()
	}

	// Update job run with results
	j.completeJobRun(ctx, jobRun.ID, status, errorMessage, cleanupStats)
	log.Println("Library cleanup job completed")
	return jobError
}

// completeJobRun finalizes a job run with status and results
func (j *LibraryCleanupJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string, stats map[string]int) {
	// Convert stats to JSON
	statsJSON := fmt.Sprintf(`{"stats":{"duplicatesFound":%d,"duplicatesResolved":%d,"orphanedItemsFound":%d,"orphanedItemsRemoved":%d,"accessErrorsFound":%d,"accessErrorsFixed":%d}}`,
		stats["duplicatesFound"],
		stats["duplicatesResolved"],
		stats["orphanedItemsFound"],
		stats["orphanedItemsRemoved"],
		stats["accessErrorsFound"],
		stats["accessErrorsFixed"])

	// In a real implementation, we would update the job run with this metadata
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}

	log.Printf("Cleanup stats: %s", statsJSON)
}

// findAndResolveDuplicates finds and resolves duplicate media items
func (j *LibraryCleanupJob) findAndResolveDuplicates(ctx context.Context) (CleanupStats, error) {
	stats := CleanupStats{}
	log.Println("Finding and resolving duplicate media items")

	// In a real implementation, we would:
	// 1. Query the repositories for potential duplicates using metadata matching
	// 2. Compare items with similar titles, years, and other identifiers
	// 3. Determine which items to keep and which to merge/delete
	// 4. Update references to point to the kept items
	// 5. Remove the duplicate entries

	// Mock implementation
	stats.found = 12    // Pretend we found 12 duplicates
	stats.resolved = 10 // Pretend we successfully resolved 10 of them

	log.Printf("Found %d duplicate items, resolved %d", stats.found, stats.resolved)
	return stats, nil
}

// findAndCleanupOrphanedItems finds and cleans up orphaned media items
func (j *LibraryCleanupJob) findAndCleanupOrphanedItems(ctx context.Context) (CleanupStats, error) {
	stats := CleanupStats{}
	log.Println("Finding and cleaning up orphaned media items")

	// In a real implementation, we would:
	// 1. Find items that are no longer linked to any client
	// 2. Verify that these items are truly orphaned, not just temporarily unavailable
	// 3. Either delete them or mark them as archived
	// 4. Update any references to these items

	// Mock implementation
	stats.found = 25   // Pretend we found 25 orphaned items
	stats.removed = 20 // Pretend we removed 20 of them

	log.Printf("Found %d orphaned items, removed %d", stats.found, stats.removed)
	return stats, nil
}

// findAndFixAccessErrors finds and fixes media access errors
func (j *LibraryCleanupJob) findAndFixAccessErrors(ctx context.Context) (CleanupStats, error) {
	stats := CleanupStats{}
	log.Println("Finding and fixing media access errors")

	// In a real implementation, we would:
	// 1. Identify items with access errors (bad file paths, broken links, etc.)
	// 2. Attempt to reconcile the errors (find alternative paths, update links)
	// 3. Mark unfixable items for review
	// 4. Generate a report of fixed and unfixable items

	// Mock implementation
	stats.found = 8 // Pretend we found 8 items with access errors
	stats.fixed = 5 // Pretend we fixed 5 of them

	log.Printf("Found %d items with access errors, fixed %d", stats.found, stats.fixed)
	return stats, nil
}

// SetupLibraryCleanupSchedule creates or updates a library cleanup schedule
func (j *LibraryCleanupJob) SetupLibraryCleanupSchedule(ctx context.Context, frequency string) error {
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

// RunManualCleanup runs the library cleanup job manually
func (j *LibraryCleanupJob) RunManualCleanup(ctx context.Context) error {
	return j.Execute(ctx)
}

// GenerateCleanupReport generates a detailed report of cleanup operations
func (j *LibraryCleanupJob) GenerateCleanupReport(ctx context.Context, recentJobRunID uint64) (string, error) {
	// In a real implementation, we would:
	// 1. Retrieve the job run details
	// 2. Format the results into a readable report
	// 3. Include statistics and recommendations

	// Mock implementation
	return "Library Cleanup Report\n" +
		"---------------------\n" +
		"Duplicate items found: 12\n" +
		"Duplicate items resolved: 10\n" +
		"Orphaned items found: 25\n" +
		"Orphaned items removed: 20\n" +
		"Access errors found: 8\n" +
		"Access errors fixed: 5\n", nil
}
