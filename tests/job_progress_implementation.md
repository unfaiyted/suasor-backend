# Job Progress Tracking Implementation

This document outlines a proposed enhancement to add job progress tracking to the Suasor backend.

## Overview

The current job system doesn't provide detailed progress information during job execution. We need to enhance it to:

1. Track progress percentage
2. Report items processed/remaining
3. Provide status messages during execution
4. Allow polling for current status

## Database Changes

Extend the `JobRun` model to include progress information:

```go
// JobRun represents a single execution of a scheduled job
type JobRun struct {
    BaseModel
    // Existing fields...
    
    // Progress percentage (0-100)
    Progress int `json:"progress" gorm:"not null;default:0"`
    // Total items to process
    TotalItems int `json:"totalItems" gorm:"default:0"`
    // Items processed so far
    ProcessedItems int `json:"processedItems" gorm:"default:0"`
    // Current status message
    StatusMessage string `json:"statusMessage"`
    // Detailed progress data (stored as JSON)
    ProgressData string `json:"progressData" gorm:"type:jsonb"`
}
```

## API Endpoints

Add these new endpoints:

1. `GET /api/v1/jobs/runs/:id/progress` - Get progress for a specific job run
2. `GET /api/v1/jobs/active` - Get all currently active job runs with progress
3. `GET /api/v1/jobs/media-sync/:id/progress` - Get progress for a media sync job

## Service Layer

Enhance the JobService interface:

```go
type JobService interface {
    // Existing methods...
    
    // Get progress for a specific job run
    GetJobRunProgress(ctx context.Context, jobRunID uint64) (*JobRunProgress, error)
    // Get all active job runs
    GetActiveJobRuns(ctx context.Context) ([]JobRunWithProgress, error)
    // Update job progress
    UpdateJobProgress(ctx context.Context, jobRunID uint64, progress int, message string) error
    // Set total items to process
    SetJobTotalItems(ctx context.Context, jobRunID uint64, totalItems int) error
    // Increment processed items
    IncrementJobProcessedItems(ctx context.Context, jobRunID uint64, count int) error
}
```

## Implementation for Media Sync Jobs

Enhance the MediaSyncJob to report progress:

```go
// SyncUserMediaFromClient synchronizes media from a client to the local database
func (j *MediaSyncJob) SyncUserMediaFromClient(ctx context.Context, userID, clientID uint64, mediaType string) error {
    // Create job run record
    jobRun, err := j.createJobRun(ctx, userID)
    if err != nil {
        return err
    }
    
    // Update progress to 0%
    j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 0, "Starting media sync")
    
    // Fetch media items from client
    j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 10, "Fetching media items from client")
    items, err := j.fetchMediaItems(ctx, clientID, mediaType)
    if err != nil {
        j.jobRepo.CompleteJobRun(ctx, jobRun.ID, models.JobStatusFailed, err.Error())
        return err
    }
    
    // Set total items
    j.jobRepo.SetJobTotalItems(ctx, jobRun.ID, len(items))
    j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 20, fmt.Sprintf("Found %d items to process", len(items)))
    
    // Process each item
    for i, item := range items {
        // Process the item...
        
        // Update progress (20-90%)
        progress := 20 + int(70.0*float64(i+1)/float64(len(items)))
        j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, progress, fmt.Sprintf("Processing item %d of %d", i+1, len(items)))
        j.jobRepo.IncrementJobProcessedItems(ctx, jobRun.ID, 1)
    }
    
    // Update progress to 100%
    j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 100, "Media sync completed successfully")
    
    // Mark job as completed
    j.jobRepo.CompleteJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "")
    
    return nil
}
```

## Frontend Integration

The frontend can poll these endpoints to show job progress:

```javascript
// Poll job progress every 2 seconds
function pollJobProgress(jobRunId) {
  const interval = setInterval(async () => {
    const response = await fetch(`/api/v1/jobs/runs/${jobRunId}/progress`);
    const data = await response.json();
    
    updateProgressUI(data);
    
    if (data.status !== 'running') {
      clearInterval(interval);
    }
  }, 2000);
}
```

## Implementation Steps

1. Add the new database fields to JobRun model
2. Create migration script to update existing tables
3. Implement repository methods for progress tracking
4. Enhance JobService with progress-related methods
5. Update job handlers to expose new endpoints
6. Register new routes in router
7. Modify job implementations to report progress