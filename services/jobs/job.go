package jobs

import (
	"context"
	"fmt"
	"log"
	mediatypes "suasor/clients/media/types"
	"suasor/repository"
	"suasor/services/jobs/recommendation"
	"suasor/services/jobs/sync"
	"suasor/services/scheduler"
	"suasor/types/models"
	"time"
)

// JobService manages job scheduling and execution
type JobService interface {
	// StartScheduler starts the job scheduler
	StartScheduler() error
	// StopScheduler stops the job scheduler
	StopScheduler() error
	// RegisterJob adds a job to the scheduler
	RegisterJob(job scheduler.Job) error
	// SyncDatabaseJobs synchronizes database job schedules with the scheduler
	SyncDatabaseJobs(ctx context.Context) error

	// GetAllJobSchedules retrieves all job schedules
	GetAllJobSchedules(ctx context.Context) ([]models.JobSchedule, error)
	// GetJobScheduleByName retrieves a job schedule by name
	GetJobScheduleByName(ctx context.Context, name string) (*models.JobSchedule, error)
	// CreateJobSchedule creates a new job schedule
	CreateJobSchedule(ctx context.Context, schedule *models.JobSchedule) error
	// UpdateJobSchedule updates an existing job schedule
	UpdateJobSchedule(ctx context.Context, schedule *models.JobSchedule) error
	// DeleteJobSchedule deletes a job schedule
	DeleteJobSchedule(ctx context.Context, name string) error

	// GetRecentJobRuns retrieves recent job runs
	GetRecentJobRuns(ctx context.Context, limit int) ([]models.JobRun, error)
	// GetUserJobRuns retrieves job runs for a specific user
	GetUserJobRuns(ctx context.Context, userID uint64, limit int) ([]models.JobRun, error)
	// GetJobRunByID retrieves a specific job run by ID
	GetJobRunByID(ctx context.Context, jobRunID uint64) (*models.JobRun, error)
	// GetActiveJobRuns retrieves all currently active job runs
	GetActiveJobRuns(ctx context.Context) ([]models.JobRun, error)

	// RunJobManually triggers a job to run immediately
	RunJobManually(ctx context.Context, jobName string) error

	// Job progress tracking methods
	// UpdateJobProgress updates the progress of a job run
	UpdateJobProgress(ctx context.Context, jobRunID uint64, progress int, message string) error
	// SetJobTotalItems sets the total number of items to be processed in a job
	SetJobTotalItems(ctx context.Context, jobRunID uint64, totalItems int) error
	// IncrementJobProcessedItems increments the number of processed items in a job
	IncrementJobProcessedItems(ctx context.Context, jobRunID uint64, count int) error

	// GetUserRecommendations retrieves recommendations for a user
	GetUserRecommendations(ctx context.Context, userID uint64, active bool, limit int) ([]models.Recommendation, error)
	// DismissRecommendation marks a recommendation as dismissed
	DismissRecommendation(ctx context.Context, recommendationID uint64) error
	// UpdateRecommendationViewedStatus updates whether a recommendation has been viewed
	UpdateRecommendationViewedStatus(ctx context.Context, recommendationID uint64, viewed bool) error

	// SetupMediaSyncJob creates or updates a media sync job
	SetupMediaSyncJob(ctx context.Context, userID, clientID uint64, clientType, mediaType, frequency string) error
	// RunMediaSyncJob runs a media sync job manually
	RunMediaSyncJob(ctx context.Context, userID, clientID uint64, syncType models.SyncType) error
	// GetMediaSyncJobs retrieves all media sync jobs for a user
	GetMediaSyncJobs(ctx context.Context, userID uint64) ([]models.MediaSyncJob, error)
}

type jobService struct {
	jobRepo             repository.JobRepository
	userRepo            repository.UserRepository
	configRepo          repository.UserConfigRepository
	movieRepo           repository.CoreMediaItemRepository[*mediatypes.Movie]
	seriesRepo          repository.CoreMediaItemRepository[*mediatypes.Series]
	trackRepo           repository.CoreMediaItemRepository[*mediatypes.Track]
	userMovieDataRepo   repository.UserMediaItemDataRepository[*mediatypes.Movie]
	userSeriesDataRepo  repository.UserMediaItemDataRepository[*mediatypes.Series]
	userTrackDataRepo   repository.UserMediaItemDataRepository[*mediatypes.Track]
	scheduler           *scheduler.Scheduler
	jobs                map[string]scheduler.Job
	recommendationJob   *recommendation.RecommendationJob
	mediaSyncJob        *sync.MediaSyncJob
	watchHistorySyncJob *sync.WatchHistorySyncJob
	favoritesSyncJob    *sync.FavoritesSyncJob
}

// NewJobService creates a new job service
func NewJobService(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.CoreMediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.CoreMediaItemRepository[*mediatypes.Series],
	trackRepo repository.CoreMediaItemRepository[*mediatypes.Track],
	userMovieDataRepo repository.UserMediaItemDataRepository[*mediatypes.Movie],
	userSeriesDataRepo repository.UserMediaItemDataRepository[*mediatypes.Series],
	userTrackDataRepo repository.UserMediaItemDataRepository[*mediatypes.Track],
	recommendationJob *recommendation.RecommendationJob,
	mediaSyncJob *sync.MediaSyncJob,
	watchHistorySyncJob *sync.WatchHistorySyncJob,
	favoritesSyncJob *sync.FavoritesSyncJob,
) JobService {
	return &jobService{
		jobRepo:             jobRepo,
		userRepo:            userRepo,
		configRepo:          configRepo,
		movieRepo:           movieRepo,
		seriesRepo:          seriesRepo,
		trackRepo:           trackRepo,
		userMovieDataRepo:   userMovieDataRepo,
		userSeriesDataRepo:  userSeriesDataRepo,
		userTrackDataRepo:   userTrackDataRepo,
		scheduler:           scheduler.NewScheduler(),
		jobs:                make(map[string]scheduler.Job),
		recommendationJob:   recommendationJob,
		mediaSyncJob:        mediaSyncJob,
		watchHistorySyncJob: watchHistorySyncJob,
		favoritesSyncJob:    favoritesSyncJob,
	}
}

// StartScheduler starts the job scheduler
func (s *jobService) StartScheduler() error {
	// Check if we have any jobs before starting the scheduler
	if len(s.jobs) == 0 {
		log.Printf("Warning: Starting scheduler with no registered jobs")
		// Continue anyway - the scheduler can handle zero jobs
	}

	// Verify all job implementations are defined
	if s.recommendationJob == nil {
		log.Printf("Warning: recommendationJob is nil, some recommendation functionality will be unavailable")
	}
	if s.mediaSyncJob == nil {
		log.Printf("Warning: mediaSyncJob is nil, some media sync functionality will be unavailable")
	}
	if s.watchHistorySyncJob == nil {
		log.Printf("Warning: watchHistorySyncJob is nil, some watch history sync functionality will be unavailable")
	}
	if s.favoritesSyncJob == nil {
		log.Printf("Warning: favoritesSyncJob is nil, some favorites sync functionality will be unavailable")
	}

	// Start the scheduler
	s.scheduler.Start()
	return nil
}

// StopScheduler stops the job scheduler
func (s *jobService) StopScheduler() error {
	s.scheduler.Stop()
	return nil
}

// RegisterJob adds a job to the scheduler
func (s *jobService) RegisterJob(job scheduler.Job) error {
	s.jobs[job.Name()] = job
	s.scheduler.RegisterJob(job)
	return nil
}

// SyncDatabaseJobs synchronizes database job schedules with the scheduler
func (s *jobService) SyncDatabaseJobs(ctx context.Context) error {
	// First, register all system jobs to ensure they exist in the database
	systemJobs := make([]scheduler.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		systemJobs = append(systemJobs, job)
	}

	// Check if we have any jobs to register
	if len(systemJobs) == 0 {
		log.Printf("Warning: No job implementations registered. Job scheduler will not run any ")
		return nil
	}

	// Register system jobs in the database
	if err := RegisterSystemJobs(ctx, s.jobRepo, systemJobs); err != nil {
		return fmt.Errorf("error registering system jobs: %w", err)
	}

	// Now get all job schedules from the database
	schedules, err := s.jobRepo.GetAllJobSchedules(ctx)
	if err != nil {
		return fmt.Errorf("error getting job schedules: %w", err)
	}

	// Register each enabled job with the scheduler
	for _, schedule := range schedules {
		if !schedule.Enabled {
			continue
		}

		// Skip if we don't have a job implementation for this name
		job, ok := s.jobs[schedule.JobName]
		if !ok {
			log.Printf("No job implementation found for job: %s", schedule.JobName)
			continue
		}

		// Register the job
		s.scheduler.RegisterJob(job)
	}

	return nil
}

// GetAllJobSchedules retrieves all job schedules
func (s *jobService) GetAllJobSchedules(ctx context.Context) ([]models.JobSchedule, error) {
	return s.jobRepo.GetAllJobSchedules(ctx)
}

// GetJobScheduleByName retrieves a job schedule by name
func (s *jobService) GetJobScheduleByName(ctx context.Context, name string) (*models.JobSchedule, error) {
	return s.jobRepo.GetJobSchedule(ctx, name)
}

// CreateJobSchedule creates a new job schedule
func (s *jobService) CreateJobSchedule(ctx context.Context, schedule *models.JobSchedule) error {
	return s.jobRepo.CreateJobSchedule(ctx, schedule)
}

// UpdateJobSchedule updates an existing job schedule
func (s *jobService) UpdateJobSchedule(ctx context.Context, schedule *models.JobSchedule) error {
	return s.jobRepo.UpdateJobSchedule(ctx, schedule)
}

// DeleteJobSchedule deletes a job schedule
func (s *jobService) DeleteJobSchedule(ctx context.Context, name string) error {
	return s.jobRepo.DeleteJobSchedule(ctx, name)
}

// GetRecentJobRuns retrieves recent job runs
func (s *jobService) GetRecentJobRuns(ctx context.Context, limit int) ([]models.JobRun, error) {
	return s.jobRepo.GetRecentJobRuns(ctx, limit)
}

// GetUserJobRuns retrieves job runs for a specific user
func (s *jobService) GetUserJobRuns(ctx context.Context, userID uint64, limit int) ([]models.JobRun, error) {
	return s.jobRepo.GetJobRunsByUser(ctx, userID, limit)
}

// GetJobRunByID retrieves a specific job run by ID
func (s *jobService) GetJobRunByID(ctx context.Context, jobRunID uint64) (*models.JobRun, error) {
	return s.jobRepo.GetJobRunByID(ctx, jobRunID)
}

// GetActiveJobRuns retrieves all currently active job runs
func (s *jobService) GetActiveJobRuns(ctx context.Context) ([]models.JobRun, error) {
	return s.jobRepo.GetActiveJobRuns(ctx)
}

// UpdateJobProgress updates the progress of a job run
func (s *jobService) UpdateJobProgress(ctx context.Context, jobRunID uint64, progress int, message string) error {
	return s.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, message)
}

// SetJobTotalItems sets the total number of items to be processed in a job
func (s *jobService) SetJobTotalItems(ctx context.Context, jobRunID uint64, totalItems int) error {
	return s.jobRepo.SetJobTotalItems(ctx, jobRunID, totalItems)
}

// IncrementJobProcessedItems increments the number of processed items in a job
func (s *jobService) IncrementJobProcessedItems(ctx context.Context, jobRunID uint64, count int) error {
	return s.jobRepo.IncrementJobProcessedItems(ctx, jobRunID, count)
}

// RunJobManually triggers a job to run immediately
func (s *jobService) RunJobManually(ctx context.Context, jobName string) error {
	job, ok := s.jobs[jobName]
	if !ok {
		return fmt.Errorf("job not found: %s", jobName)
	}

	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   job.Name(),
		JobType:   getJobType(job),
		Status:    models.JobStatusRunning,
		StartTime: &now,
		Metadata:  "{}", // Initialize with an empty JSON object
	}

	err := s.jobRepo.CreateJobRun(ctx, jobRun)
	if err != nil {
		return fmt.Errorf("error creating job run record: %w", err)
	}

	// Execute the job
	go func() {
		// Create a new context for the job execution
		jobCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		// Execute the job
		err := job.Execute(jobCtx)

		// Update the job run record
		status := models.JobStatusCompleted
		errMsg := ""
		if err != nil {
			status = models.JobStatusFailed
			errMsg = err.Error()
			log.Printf("Error executing job %s: %v", job.Name(), err)
		}

		completeErr := s.jobRepo.CompleteJobRun(jobCtx, jobRun.ID, status, errMsg)
		if completeErr != nil {
			log.Printf("Error completing job run record: %v", completeErr)
		}

		// Update the job schedule's last run time
		updateErr := s.jobRepo.UpdateJobLastRunTime(jobCtx, job.Name(), time.Now())
		if updateErr != nil {
			log.Printf("Error updating job last run time: %v", updateErr)
		}
	}()

	return nil
}

// GetUserRecommendations retrieves recommendations for a user
func (s *jobService) GetUserRecommendations(ctx context.Context, userID uint64, active bool, limit int) ([]models.Recommendation, error) {
	return s.jobRepo.GetUserRecommendations(ctx, userID, active, limit)
}

// DismissRecommendation marks a recommendation as dismissed
func (s *jobService) DismissRecommendation(ctx context.Context, recommendationID uint64) error {
	return s.jobRepo.DismissRecommendation(ctx, recommendationID)
}

// UpdateRecommendationViewedStatus updates whether a recommendation has been viewed
func (s *jobService) UpdateRecommendationViewedStatus(ctx context.Context, recommendationID uint64, viewed bool) error {
	return s.jobRepo.UpdateRecommendationViewedStatus(ctx, recommendationID, viewed)
}

// SetupMediaSyncJob creates or updates a media sync job
func (s *jobService) SetupMediaSyncJob(ctx context.Context, userID, clientID uint64, clientType, mediaType, frequency string) error {
	// Validate inputs
	if userID == 0 || clientID == 0 || mediaType == "" || frequency == "" {
		return fmt.Errorf("invalid input parameters")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found: %d", userID)
	}

	// Create or update the sync job
	return s.recommendationJob.SetupMediaSyncJob(ctx, userID, clientID, clientType, mediaType, frequency)
}

// RunMediaSyncJob runs a media sync job manually
func (s *jobService) RunMediaSyncJob(ctx context.Context, userID, clientID uint64, syncType models.SyncType) error {
	// Validate inputs
	if userID == 0 || clientID == 0 || syncType == "" {
		return fmt.Errorf("invalid input parameters")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found: %d", userID)
	}

	// Run the sync job
	return s.mediaSyncJob.SyncUserMediaFromClient(ctx, userID, clientID, syncType)
}

// GetMediaSyncJobs retrieves all media sync jobs for a user
func (s *jobService) GetMediaSyncJobs(ctx context.Context, userID uint64) ([]models.MediaSyncJob, error) {
	return s.jobRepo.GetMediaSyncJobsByUser(ctx, userID)
}

// getJobType determines the job type from a job
func getJobType(job scheduler.Job) models.JobType {
	switch job.(type) {
	case *recommendation.RecommendationJob:
		return models.JobTypeRecommendation
	case *sync.MediaSyncJob, *sync.WatchHistorySyncJob, *sync.FavoritesSyncJob:
		return models.JobTypeSync
	default:
		return models.JobType("unknown")
	}
}
