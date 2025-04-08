package jobs

import (
	"context"
	"fmt"
	"log"
	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"time"
)

// MediaSyncJob synchronizes media from external clients to the local database
type MediaSyncJob struct {
	jobRepo     repository.JobRepository
	userRepo    repository.UserRepository
	configRepo  repository.UserConfigRepository
	movieRepo   repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo  repository.MediaItemRepository[*mediatypes.Series]
	musicRepo   repository.MediaItemRepository[*mediatypes.Track]
	historyRepo repository.MediaPlayHistoryRepository
}

// NewMediaSyncJob creates a new media sync job
func NewMediaSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	historyRepo repository.MediaPlayHistoryRepository,
	_ interface{},  // embyService
	_ interface{},  // jellyfinService
	_ interface{},  // plexService
	_ interface{},  // subsonicService
) *MediaSyncJob {
	return &MediaSyncJob{
		jobRepo:         jobRepo,
		userRepo:        userRepo,
		configRepo:      configRepo,
		movieRepo:       movieRepo,
		seriesRepo:      seriesRepo,
		musicRepo:       musicRepo,
		historyRepo:     historyRepo,
	}
}

// Name returns the unique name of the job
func (j *MediaSyncJob) Name() string {
	return "system.media.sync"
}

// Schedule returns when the job should next run
func (j *MediaSyncJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute runs the media sync job
func (j *MediaSyncJob) Execute(ctx context.Context) error {
	log.Println("Starting media sync job")

	// This is a placeholder implementation
	log.Println("Media sync job completed (placeholder)")
	return nil
}

// SyncUserMediaFromClient synchronizes media from a client to the local database
func (j *MediaSyncJob) SyncUserMediaFromClient(ctx context.Context, userID, clientID uint64, mediaType string) error {
	log.Printf("Syncing %s media for user %d from client %d", mediaType, userID, clientID)
	
	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:       j.Name(),
		JobType:       models.JobTypeSync,
		Status:        models.JobStatusRunning,
		StartTime:     &now,
		UserID:        &userID,
		Progress:      0,
		StatusMessage: fmt.Sprintf("Starting %s sync for client %d", mediaType, clientID),
	}
	
	err := j.jobRepo.CreateJobRun(ctx, jobRun)
	if err != nil {
		return fmt.Errorf("error creating job run record: %w", err)
	}
	
	// Run the sync process
	go func() {
		// Set up a new context with timeout
		syncCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		
		// Simulate fetching items (10% progress)
		j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, 10, fmt.Sprintf("Fetching %s items from client", mediaType))
		time.Sleep(2 * time.Second) // Simulate work
		
		// Set total items (for demonstration, using a fixed number)
		totalItems := 50
		j.jobRepo.SetJobTotalItems(syncCtx, jobRun.ID, totalItems)
		j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, 20, fmt.Sprintf("Found %d items to process", totalItems))
		time.Sleep(1 * time.Second) // Simulate work
		
		// Process items in batches, updating progress along the way
		batchSize := 5
		for i := 0; i < totalItems; i += batchSize {
			// Check if context is cancelled
			if syncCtx.Err() != nil {
				j.jobRepo.CompleteJobRun(syncCtx, jobRun.ID, models.JobStatusFailed, "Job was cancelled")
				return
			}
			
			// Calculate actual batch size (might be smaller at the end)
			currentBatch := batchSize
			if i+currentBatch > totalItems {
				currentBatch = totalItems - i
			}
			
			// Update processed items and progress message
			j.jobRepo.IncrementJobProcessedItems(syncCtx, jobRun.ID, currentBatch)
			progressMsg := fmt.Sprintf("Processed %d/%d items", i+currentBatch, totalItems)
			progress := 20 + int(float64(i+currentBatch)/float64(totalItems)*70.0)
			j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, progress, progressMsg)
			
			// Simulate processing time
			time.Sleep(1 * time.Second)
		}
		
		// Finalize the sync (last 10% of progress)
		j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, 90, "Finalizing sync...")
		time.Sleep(1 * time.Second)
		
		// Complete the job
		j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, 100, fmt.Sprintf("%s sync completed successfully", mediaType))
		j.jobRepo.CompleteJobRun(syncCtx, jobRun.ID, models.JobStatusCompleted, "")
		
		// Update the last sync time for this media sync job
		// We'd typically get the sync job ID from the database
		// This is a simplification
		log.Printf("%s sync for user %d from client %d completed", mediaType, userID, clientID)
	}()
	
	return nil
}