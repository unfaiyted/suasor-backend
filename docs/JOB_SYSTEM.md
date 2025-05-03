# Job System

---
Created: 2025-05-03
Last Updated: 2025-05-03
Update Frequency: As needed when job system changes
Owner: Backend Team
---

## Overview

Suasor includes a robust job system for handling background processing tasks such as media synchronization, metadata refreshing, recommendation generation, and system maintenance. This document explains the architecture, components, and usage patterns of the job system.

## Architecture

The job system consists of the following components:

1. **Job Interface**: Defines the contract for all job implementations
2. **Scheduler**: Manages job scheduling and execution
3. **Job Repository**: Persists job schedules and execution history
4. **Job Service**: Provides application-level job management
5. **Job Implementations**: Concrete job types for various tasks

```
┌─────────────┐      ┌───────────┐     ┌─────────────┐
│   Job       │──────▶ Scheduler │────▶│   Job       │
│  Service    │      │           │     │ Repository  │
└─────────────┘      └───────────┘     └─────────────┘
       │                 │                  │
       │                 │                  │
       ▼                 ▼                  ▼
┌─────────────┐      ┌───────────┐     ┌─────────────┐
│    Job      │      │   Job     │     │  Database   │
│Implementations│    │ Execution │     │   Storage   │
└─────────────┘      └───────────┘     └─────────────┘
```

## Core Components

### Job Interface

All jobs implement the `scheduler.Job` interface:

```go
// Job represents a scheduled job that can be executed
type Job interface {
    // Execute runs the job with the given context
    Execute(ctx context.Context) error
    // Name returns the unique name of the job
    Name() string
    // Schedule returns when the job should next run
    Schedule() time.Duration
}
```

### Scheduler

The scheduler manages the execution of jobs according to their schedule:

```go
// Scheduler manages the execution of scheduled jobs
type Scheduler struct {
    jobs       map[string]Job
    jobTimers  map[string]*time.Timer
    mutex      sync.Mutex
    cancelFunc context.CancelFunc
    ctx        context.Context
    wg         sync.WaitGroup
}
```

Key scheduler methods:

- `RegisterJob(job Job)`: Adds a job to the scheduler
- `Start()`: Begins the scheduler, executing all registered jobs
- `Stop()`: Cancels all scheduled jobs
- `scheduleJob(name string, job Job, delay time.Duration)`: Creates a timer for job execution
- `executeJob(name string, job Job)`: Runs a job and reschedules it

### Job Repository

The job repository persists job schedules and execution history:

```go
// JobRepository handles job scheduling and tracking operations
type JobRepository interface {
    // Job schedule methods
    GetAllJobSchedules(ctx context.Context) ([]models.JobSchedule, error)
    GetJobSchedulesByType(ctx context.Context, jobType models.JobType) ([]models.JobSchedule, error)
    GetUserJobSchedules(ctx context.Context, userID uint64) ([]models.JobSchedule, error)
    GetJobSchedule(ctx context.Context, jobName string) (*models.JobSchedule, error)
    CreateJobSchedule(ctx context.Context, jobSchedule *models.JobSchedule) error
    UpdateJobSchedule(ctx context.Context, jobSchedule *models.JobSchedule) error
    UpdateJobLastRunTime(ctx context.Context, jobName string, lastRunTime time.Time) error
    DeleteJobSchedule(ctx context.Context, jobName string) error
    
    // Job run methods
    CreateJobRun(ctx context.Context, jobRun *models.JobRun) error
    UpdateJobRunStatus(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) error
    CompleteJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) error
    GetRecentJobRuns(ctx context.Context, limit int) ([]models.JobRun, error)
    GetJobRunsByUser(ctx context.Context, userID uint64, limit int) ([]models.JobRun, error)
    GetJobRunByID(ctx context.Context, jobRunID uint64) (*models.JobRun, error)
    GetActiveJobRuns(ctx context.Context) ([]models.JobRun, error)
    
    // Job progress methods
    UpdateJobProgress(ctx context.Context, jobRunID uint64, progress int, message string) error
    SetJobTotalItems(ctx context.Context, jobRunID uint64, totalItems int) error
    IncrementJobProcessedItems(ctx context.Context, jobRunID uint64, count int) error
    
    // Additional methods for specific job types...
}
```

### Job Service

The job service provides application-level job management:

```go
// JobService manages job scheduling and execution
type JobService interface {
    // Scheduler management
    StartScheduler() error
    StopScheduler() error
    RegisterJob(job scheduler.Job) error
    SyncDatabaseJobs(ctx context.Context) error
    
    // Job schedule management
    GetAllJobSchedules(ctx context.Context) ([]models.JobSchedule, error)
    GetJobScheduleByName(ctx context.Context, name string) (*models.JobSchedule, error)
    CreateJobSchedule(ctx context.Context, schedule *models.JobSchedule) error
    UpdateJobSchedule(ctx context.Context, schedule *models.JobSchedule) error
    DeleteJobSchedule(ctx context.Context, name string) error
    
    // Job run management
    GetRecentJobRuns(ctx context.Context, limit int) ([]models.JobRun, error)
    GetUserJobRuns(ctx context.Context, userID uint64, limit int) ([]models.JobRun, error)
    GetJobRunByID(ctx context.Context, jobRunID uint64) (*models.JobRun, error)
    GetActiveJobRuns(ctx context.Context) ([]models.JobRun, error)
    RunJobManually(ctx context.Context, jobName string) error
    
    // Job progress tracking
    UpdateJobProgress(ctx context.Context, jobRunID uint64, progress int, message string) error
    SetJobTotalItems(ctx context.Context, jobRunID uint64, totalItems int) error
    IncrementJobProcessedItems(ctx context.Context, jobRunID uint64, count int) error
    
    // Additional methods for specific job types...
}
```

## Job Types

Suasor implements several types of background jobs:

### System Jobs

1. **Database Maintenance Job**:
   - Optimizes database performance
   - Cleans up old records
   - Performs database integrity checks

2. **Metadata Refresh Job**:
   - Updates metadata for media items
   - Refreshes information from external sources (TMDB, etc.)
   - Prioritizes recently added or popular items

3. **Library Cleanup Job**:
   - Identifies orphaned items
   - Removes invalid references
   - Updates stale metadata

### User-Specific Jobs

1. **Recommendation Job**:
   - Analyzes user preferences
   - Generates personalized recommendations
   - Updates recommendation database

2. **New Release Notification Job**:
   - Checks for new releases matching user preferences
   - Creates notifications for users
   - Filters based on user watch history

3. **User Activity Analysis Job**:
   - Analyzes user viewing patterns
   - Updates user profile data
   - Improves recommendation quality

### Media Synchronization Jobs

1. **Media Sync Job**:
   - Synchronizes media between clients and the database
   - Updates metadata from client sources
   - Handles incremental updates

2. **Favorites Sync Job**:
   - Synchronizes user favorites across clients
   - Maintains consistent favorite status
   - Resolves conflicts between clients

3. **Smart Collection Job**:
   - Updates dynamic collections based on rules
   - Refreshes collection contents
   - Applies smart filters

## Job Implementation Example

Here's an example of a job implementation:

```go
// MetadataRefreshJob periodically updates metadata for media items from external sources
type MetadataRefreshJob struct {
    jobRepo           repository.JobRepository
    userRepo          repository.UserRepository
    configRepo        repository.UserConfigRepository
    movieRepo         repository.CoreMediaItemRepository[*mediatypes.Movie]
    seriesRepo        repository.CoreMediaItemRepository[*mediatypes.Series]
    episodeRepo       repository.CoreMediaItemRepository[*mediatypes.Episode]
    musicRepo         repository.CoreMediaItemRepository[*mediatypes.Track]
    metadataClientSvc interface{}
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

    // [Job implementation details]

    // Complete the job
    status := models.JobStatusCompleted
    errorMessage := ""
    if jobError != nil {
        status = models.JobStatusFailed
        errorMessage = jobError.Error()
    }

    // Update job run with results
    j.completeJobRun(ctx, jobRun.ID, status, errorMessage, refreshStats)
    return jobError
}
```

## Job Scheduling

Jobs can be scheduled with different frequencies:

```go
// Frequency represents how often a job should run
type Frequency string

const (
    // FrequencyManual job only runs manually
    FrequencyManual Frequency = "manual"
    // FrequencyDaily job runs daily
    FrequencyDaily Frequency = "daily"
    // FrequencyWeekly job runs weekly
    FrequencyWeekly Frequency = "weekly"
    // FrequencyMonthly job runs monthly
    FrequencyMonthly Frequency = "monthly"
)
```

The `Frequency` type provides methods to calculate next run time:

```go
// NextRunTime calculates when a job with this frequency should next run
func (f Frequency) NextRunTime(lastRun time.Time) time.Time {
    switch f {
    case FrequencyDaily:
        return lastRun.AddDate(0, 0, 1)
    case FrequencyWeekly:
        return lastRun.AddDate(0, 0, 7)
    case FrequencyMonthly:
        return lastRun.AddDate(0, 1, 0)
    default:
        return lastRun.AddDate(1, 0, 0)
    }
}
```

## Job Progress Tracking

Jobs can report progress, which is stored in the database:

```go
// Update progress
func (s *jobService) UpdateJobProgress(ctx context.Context, jobRunID uint64, progress int, message string) error {
    return s.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, message)
}

// Set total items to process
func (s *jobService) SetJobTotalItems(ctx context.Context, jobRunID uint64, totalItems int) error {
    return s.jobRepo.SetJobTotalItems(ctx, jobRunID, totalItems)
}

// Increment processed items count
func (s *jobService) IncrementJobProcessedItems(ctx context.Context, jobRunID uint64, count int) error {
    return s.jobRepo.IncrementJobProcessedItems(ctx, jobRunID, count)
}
```

## Job Registration and Initialization

The job system is initialized during application startup:

```go
// Initialize and start the job system
func InitializeJobSystem(ctx context.Context, deps *AppDependencies) error {
    // Register system jobs
    deps.JobService.RegisterJob(deps.MetadataRefreshJob)
    deps.JobService.RegisterJob(deps.DatabaseMaintenanceJob)
    deps.JobService.RegisterJob(deps.LibraryCleanupJob)
    deps.JobService.RegisterJob(deps.NewReleaseNotificationJob)
    deps.JobService.RegisterJob(deps.UserActivityAnalysisJob)
    deps.JobService.RegisterJob(deps.RecommendationJob)
    
    // Sync job schedules from database
    if err := deps.JobService.SyncDatabaseJobs(ctx); err != nil {
        return fmt.Errorf("error syncing job schedules: %w", err)
    }
    
    // Start the scheduler
    if err := deps.JobService.StartScheduler(); err != nil {
        return fmt.Errorf("error starting job scheduler: %w", err)
    }
    
    return nil
}
```

## Handler Integration

The job system exposes HTTP endpoints for job management:

```go
// JobHandler provides HTTP endpoints for job management
type JobHandler struct {
    jobService services.JobService
}

// GetJobSchedules returns all job schedules
func (h *JobHandler) GetJobSchedules(c *gin.Context) {
    ctx := c.Request.Context()
    schedules, err := h.jobService.GetAllJobSchedules(ctx)
    if err != nil {
        responses.RespondInternalError(c, err, "Failed to retrieve job schedules")
        return
    }
    responses.RespondOK(c, schedules, "Job schedules retrieved successfully")
}

// RunJob executes a job manually
func (h *JobHandler) RunJob(c *gin.Context) {
    ctx := c.Request.Context()
    jobName := c.Param("jobName")
    
    if err := h.jobService.RunJobManually(ctx, jobName); err != nil {
        responses.RespondInternalError(c, err, "Failed to run job")
        return
    }
    
    responses.RespondOK(c, nil, "Job started successfully")
}

// GetJobStatus returns the status of a job run
func (h *JobHandler) GetJobStatus(c *gin.Context) {
    ctx := c.Request.Context()
    jobRunID, err := strconv.ParseUint(c.Param("jobRunID"), 10, 64)
    if err != nil {
        responses.RespondBadRequest(c, err, "Invalid job run ID")
        return
    }
    
    jobRun, err := h.jobService.GetJobRunByID(ctx, jobRunID)
    if err != nil {
        responses.RespondInternalError(c, err, "Failed to retrieve job status")
        return
    }
    
    if jobRun == nil {
        responses.RespondNotFound(c, nil, "Job run not found")
        return
    }
    
    responses.RespondOK(c, jobRun, "Job status retrieved successfully")
}
```

## Job Models

The job system uses several data models:

### JobSchedule

```go
// JobSchedule represents a scheduled job
type JobSchedule struct {
    ID          uint64     `json:"id" gorm:"primaryKey"`
    JobName     string     `json:"jobName" gorm:"uniqueIndex"`
    JobType     JobType    `json:"jobType"`
    UserID      uint64     `json:"userId,omitempty"`
    Frequency   string     `json:"frequency"`
    Enabled     bool       `json:"enabled"`
    LastRunTime *time.Time `json:"lastRunTime,omitempty"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}
```

### JobRun

```go
// JobRun represents a single execution of a job
type JobRun struct {
    ID             uint64     `json:"id" gorm:"primaryKey"`
    JobName        string     `json:"jobName"`
    JobType        JobType    `json:"jobType"`
    UserID         uint64     `json:"userId,omitempty"`
    Status         JobStatus  `json:"status"`
    StartTime      *time.Time `json:"startTime,omitempty"`
    EndTime        *time.Time `json:"endTime,omitempty"`
    Progress       int        `json:"progress"`         // 0-100 percentage
    TotalItems     int        `json:"totalItems"`       // Total items to process
    ProcessedItems int        `json:"processedItems"`   // Number of items processed
    StatusMessage  string     `json:"statusMessage"`    // Current status message
    ErrorMessage   string     `json:"errorMessage"`     // Error message if failed
    Metadata       string     `json:"metadata"`         // Additional JSON metadata
    CreatedAt      time.Time  `json:"createdAt"`
    UpdatedAt      time.Time  `json:"updatedAt"`
}
```

### JobStatus

```go
// JobStatus represents the current status of a job run
type JobStatus string

const (
    JobStatusPending   JobStatus = "pending"
    JobStatusRunning   JobStatus = "running"
    JobStatusCompleted JobStatus = "completed"
    JobStatusFailed    JobStatus = "failed"
    JobStatusCancelled JobStatus = "cancelled"
)
```

### JobType

```go
// JobType categorizes different kinds of jobs
type JobType string

const (
    JobTypeSystem        JobType = "system"
    JobTypeUser          JobType = "user"
    JobTypeRecommendation JobType = "recommendation"
    JobTypeSync          JobType = "sync"
    JobTypeNotification  JobType = "notification"
)
```

## Best Practices

### Implementing Jobs

1. **Context Awareness**: Always respect the context cancellation signal
2. **Progress Reporting**: Update progress regularly for long-running jobs
3. **Error Handling**: Handle and report errors properly
4. **Transaction Safety**: Use database transactions where appropriate
5. **Resource Management**: Clean up resources in all cases

### Job Implementation Pattern

```go
func (j *SomeJob) Execute(ctx context.Context) error {
    // Create job run record
    jobRun := createJobRun(ctx, j)
    
    // Set up progress tracking
    j.jobRepo.SetJobTotalItems(ctx, jobRun.ID, totalItems)
    
    // Check for cancellation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Continue processing
    }
    
    // Process items in batches
    for i, batch := range batches {
        // Process batch
        
        // Update progress
        progress := int((i + 1) * 100 / len(batches))
        j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, progress, "Processing batch")
        
        // Check for cancellation between batches
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue processing
        }
    }
    
    // Complete the job
    completeJobRun(ctx, jobRun.ID, status, errorMessage)
    return nil
}
```

## Job System Initialization

The job system depends on several components and is initialized as part of the dependency injection system:

```go
// Register job components
container.RegisterFactory[services.JobService](c, func(c *container.Container) services.JobService {
    jobRepo := container.MustGet[repository.JobRepository](c)
    userRepo := container.MustGet[repository.UserRepository](c)
    configRepo := container.MustGet[repository.UserConfigRepository](c)
    movieRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Movie]](c)
    seriesRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Series]](c)
    trackRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Track]](c)
    userMovieDataRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Movie]](c)
    userSeriesDataRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Series]](c)
    userTrackDataRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Track]](c)
    recommendationJob := container.MustGet[*recommendation.RecommendationJob](c)
    mediaSyncJob := container.MustGet[*sync.MediaSyncJob](c)
    favoritesSyncJob := container.MustGet[*sync.FavoritesSyncJob](c)
    
    return services.NewJobService(
        jobRepo, userRepo, configRepo,
        movieRepo, seriesRepo, trackRepo,
        userMovieDataRepo, userSeriesDataRepo, userTrackDataRepo,
        recommendationJob, mediaSyncJob, favoritesSyncJob,
    )
})
```

## Conclusion

Suasor's job system provides a robust framework for background processing, with features for scheduling, progress tracking, and error handling. The system is designed to be flexible and extensible, allowing new job types to be added as needed.

Key benefits of this architecture include:

1. **Decoupling**: Background tasks are separated from user-facing operations
2. **Scalability**: Jobs can be scheduled and executed independently
3. **Reliability**: Progress tracking and error handling ensure robustness
4. **Visibility**: Job status and history are stored for monitoring
5. **Flexibility**: New job types can be added without changing the core system