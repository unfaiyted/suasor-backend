package repository

import (
	"context"
	"fmt"
	"suasor/types/models"
	"time"

	"gorm.io/gorm"
)

// JobRepository handles job scheduling and tracking operations
type JobRepository interface {
	// GetAllJobSchedules retrieves all job schedules
	GetAllJobSchedules(ctx context.Context) ([]models.JobSchedule, error)
	// GetJobSchedulesByType retrieves all job schedules of a specific type
	GetJobSchedulesByType(ctx context.Context, jobType models.JobType) ([]models.JobSchedule, error)
	// GetUserJobSchedules retrieves all job schedules for a specific user
	GetUserJobSchedules(ctx context.Context, userID uint64) ([]models.JobSchedule, error)
	// GetJobSchedule retrieves a specific job schedule by name
	GetJobSchedule(ctx context.Context, jobName string) (*models.JobSchedule, error)
	// CreateJobSchedule creates a new job schedule
	CreateJobSchedule(ctx context.Context, jobSchedule *models.JobSchedule) error
	// UpdateJobSchedule updates an existing job schedule
	UpdateJobSchedule(ctx context.Context, jobSchedule *models.JobSchedule) error
	// UpdateJobLastRunTime updates the last run time for a job
	UpdateJobLastRunTime(ctx context.Context, jobName string, lastRunTime time.Time) error
	// DeleteJobSchedule deletes a job schedule
	DeleteJobSchedule(ctx context.Context, jobName string) error

	// CreateJobRun creates a new job run record
	CreateJobRun(ctx context.Context, jobRun *models.JobRun) error
	// UpdateJobRunStatus updates the status of a job run
	UpdateJobRunStatus(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) error
	// CompleteJobRun marks a job run as completed
	CompleteJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) error
	// GetRecentJobRuns retrieves recent job runs
	GetRecentJobRuns(ctx context.Context, limit int) ([]models.JobRun, error)
	// GetJobRunsByUser retrieves job runs for a specific user
	GetJobRunsByUser(ctx context.Context, userID uint64, limit int) ([]models.JobRun, error)
	
	// Recommendation methods
	// CreateRecommendation creates a new recommendation
	CreateRecommendation(ctx context.Context, recommendation *models.Recommendation) error
	// GetUserRecommendations retrieves recommendations for a specific user
	GetUserRecommendations(ctx context.Context, userID uint64, active bool, limit int) ([]models.Recommendation, error)
	// GetRecommendationsByJobRun retrieves recommendations created by a specific job run
	GetRecommendationsByJobRun(ctx context.Context, jobRunID uint64) ([]models.Recommendation, error)
	// DismissRecommendation marks a recommendation as dismissed
	DismissRecommendation(ctx context.Context, recommendationID uint64) error
	// UpdateRecommendationViewedStatus updates whether the recommendation has been viewed
	UpdateRecommendationViewedStatus(ctx context.Context, recommendationID uint64, viewed bool) error
	// BatchCreateRecommendations creates multiple recommendations at once
	BatchCreateRecommendations(ctx context.Context, recommendations []*models.Recommendation) error
	// GetRecommendationByMediaItem gets a recommendation for a specific user and media item
	GetRecommendationByMediaItem(ctx context.Context, userID uint64, mediaItemID uint64) (*models.Recommendation, error)
	
	// Media sync job methods
	// CreateMediaSyncJob creates a new media sync job
	CreateMediaSyncJob(ctx context.Context, syncJob *models.MediaSyncJob) error
	// GetMediaSyncJobsByUser retrieves media sync jobs for a specific user
	GetMediaSyncJobsByUser(ctx context.Context, userID uint64) ([]models.MediaSyncJob, error)
	// UpdateMediaSyncJob updates an existing media sync job
	UpdateMediaSyncJob(ctx context.Context, syncJob *models.MediaSyncJob) error
	// UpdateMediaSyncLastRunTime updates the last run time for a media sync job
	UpdateMediaSyncLastRunTime(ctx context.Context, syncJobID uint64, lastRunTime time.Time) error
	// DeleteMediaSyncJob deletes a media sync job
	DeleteMediaSyncJob(ctx context.Context, syncJobID uint64) error
}

type jobRepository struct {
	db *gorm.DB
}

// NewJobRepository creates a new job repository
func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{
		db: db,
	}
}

// GetAllJobSchedules retrieves all job schedules
func (r *jobRepository) GetAllJobSchedules(ctx context.Context) ([]models.JobSchedule, error) {
	var schedules []models.JobSchedule
	result := r.db.WithContext(ctx).Find(&schedules)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving job schedules: %w", result.Error)
	}
	return schedules, nil
}

// GetJobSchedulesByType retrieves all job schedules of a specific type
func (r *jobRepository) GetJobSchedulesByType(ctx context.Context, jobType models.JobType) ([]models.JobSchedule, error) {
	var schedules []models.JobSchedule
	result := r.db.WithContext(ctx).Where("job_type = ?", jobType).Find(&schedules)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving job schedules of type %s: %w", jobType, result.Error)
	}
	return schedules, nil
}

// GetUserJobSchedules retrieves all job schedules for a specific user
func (r *jobRepository) GetUserJobSchedules(ctx context.Context, userID uint64) ([]models.JobSchedule, error) {
	var schedules []models.JobSchedule
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&schedules)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving job schedules for user %d: %w", userID, result.Error)
	}
	return schedules, nil
}

// GetJobSchedule retrieves a specific job schedule by name
func (r *jobRepository) GetJobSchedule(ctx context.Context, jobName string) (*models.JobSchedule, error) {
	var schedule models.JobSchedule
	result := r.db.WithContext(ctx).Where("job_name = ?", jobName).First(&schedule)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving job schedule %s: %w", jobName, result.Error)
	}
	return &schedule, nil
}

// CreateJobSchedule creates a new job schedule
func (r *jobRepository) CreateJobSchedule(ctx context.Context, jobSchedule *models.JobSchedule) error {
	result := r.db.WithContext(ctx).Create(jobSchedule)
	if result.Error != nil {
		return fmt.Errorf("error creating job schedule: %w", result.Error)
	}
	return nil
}

// UpdateJobSchedule updates an existing job schedule
func (r *jobRepository) UpdateJobSchedule(ctx context.Context, jobSchedule *models.JobSchedule) error {
	result := r.db.WithContext(ctx).Save(jobSchedule)
	if result.Error != nil {
		return fmt.Errorf("error updating job schedule: %w", result.Error)
	}
	return nil
}

// UpdateJobLastRunTime updates the last run time for a job
func (r *jobRepository) UpdateJobLastRunTime(ctx context.Context, jobName string, lastRunTime time.Time) error {
	result := r.db.WithContext(ctx).Model(&models.JobSchedule{}).
		Where("job_name = ?", jobName).
		Update("last_run_time", lastRunTime)
	if result.Error != nil {
		return fmt.Errorf("error updating job last run time: %w", result.Error)
	}
	return nil
}

// DeleteJobSchedule deletes a job schedule
func (r *jobRepository) DeleteJobSchedule(ctx context.Context, jobName string) error {
	result := r.db.WithContext(ctx).Where("job_name = ?", jobName).Delete(&models.JobSchedule{})
	if result.Error != nil {
		return fmt.Errorf("error deleting job schedule: %w", result.Error)
	}
	return nil
}

// CreateJobRun creates a new job run record
func (r *jobRepository) CreateJobRun(ctx context.Context, jobRun *models.JobRun) error {
	result := r.db.WithContext(ctx).Create(jobRun)
	if result.Error != nil {
		return fmt.Errorf("error creating job run: %w", result.Error)
	}
	return nil
}

// UpdateJobRunStatus updates the status of a job run
func (r *jobRepository) UpdateJobRunStatus(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	result := r.db.WithContext(ctx).Model(&models.JobRun{}).
		Where("id = ?", jobRunID).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("error updating job run status: %w", result.Error)
	}
	return nil
}

// CompleteJobRun marks a job run as completed
func (r *jobRepository) CompleteJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":   status,
		"end_time": now,
	}
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	result := r.db.WithContext(ctx).Model(&models.JobRun{}).
		Where("id = ?", jobRunID).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("error completing job run: %w", result.Error)
	}
	return nil
}

// GetRecentJobRuns retrieves recent job runs
func (r *jobRepository) GetRecentJobRuns(ctx context.Context, limit int) ([]models.JobRun, error) {
	var runs []models.JobRun
	result := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&runs)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving recent job runs: %w", result.Error)
	}
	return runs, nil
}

// GetJobRunsByUser retrieves job runs for a specific user
func (r *jobRepository) GetJobRunsByUser(ctx context.Context, userID uint64, limit int) ([]models.JobRun, error) {
	var runs []models.JobRun
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&runs)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving job runs for user %d: %w", userID, result.Error)
	}
	return runs, nil
}

// CreateRecommendation creates a new recommendation
func (r *jobRepository) CreateRecommendation(ctx context.Context, recommendation *models.Recommendation) error {
	result := r.db.WithContext(ctx).Create(recommendation)
	if result.Error != nil {
		return fmt.Errorf("error creating recommendation: %w", result.Error)
	}
	return nil
}

// GetUserRecommendations retrieves recommendations for a specific user
func (r *jobRepository) GetUserRecommendations(ctx context.Context, userID uint64, active bool, limit int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	if active {
		query = query.Where("active = ? AND dismissed = ?", true, false)
	}
	
	result := query.Order("created_at DESC").
		Limit(limit).
		Find(&recommendations)
	
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving recommendations for user %d: %w", userID, result.Error)
	}
	return recommendations, nil
}

// GetRecommendationsByJobRun retrieves recommendations created by a specific job run
func (r *jobRepository) GetRecommendationsByJobRun(ctx context.Context, jobRunID uint64) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation
	result := r.db.WithContext(ctx).
		Where("job_run_id = ?", jobRunID).
		Find(&recommendations)
	
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving recommendations for job run %d: %w", jobRunID, result.Error)
	}
	return recommendations, nil
}

// DismissRecommendation marks a recommendation as dismissed
func (r *jobRepository) DismissRecommendation(ctx context.Context, recommendationID uint64) error {
	result := r.db.WithContext(ctx).
		Model(&models.Recommendation{}).
		Where("id = ?", recommendationID).
		Updates(map[string]interface{}{
			"dismissed": true,
			"active":    false,
		})
	
	if result.Error != nil {
		return fmt.Errorf("error dismissing recommendation %d: %w", recommendationID, result.Error)
	}
	return nil
}

// UpdateRecommendationViewedStatus updates whether the recommendation has been viewed
func (r *jobRepository) UpdateRecommendationViewedStatus(ctx context.Context, recommendationID uint64, viewed bool) error {
	result := r.db.WithContext(ctx).
		Model(&models.Recommendation{}).
		Where("id = ?", recommendationID).
		Update("viewed", viewed)
	
	if result.Error != nil {
		return fmt.Errorf("error updating recommendation viewed status %d: %w", recommendationID, result.Error)
	}
	return nil
}

// BatchCreateRecommendations creates multiple recommendations at once
func (r *jobRepository) BatchCreateRecommendations(ctx context.Context, recommendations []*models.Recommendation) error {
	result := r.db.WithContext(ctx).Create(recommendations)
	if result.Error != nil {
		return fmt.Errorf("error batch creating recommendations: %w", result.Error)
	}
	return nil
}

// GetRecommendationByMediaItem gets a recommendation for a specific user and media item
func (r *jobRepository) GetRecommendationByMediaItem(ctx context.Context, userID uint64, mediaItemID uint64) (*models.Recommendation, error) {
	var recommendation models.Recommendation
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND media_item_id = ?", userID, mediaItemID).
		First(&recommendation)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting recommendation: %w", result.Error)
	}
	return &recommendation, nil
}

// CreateMediaSyncJob creates a new media sync job
func (r *jobRepository) CreateMediaSyncJob(ctx context.Context, syncJob *models.MediaSyncJob) error {
	result := r.db.WithContext(ctx).Create(syncJob)
	if result.Error != nil {
		return fmt.Errorf("error creating media sync job: %w", result.Error)
	}
	return nil
}

// GetMediaSyncJobsByUser retrieves media sync jobs for a specific user
func (r *jobRepository) GetMediaSyncJobsByUser(ctx context.Context, userID uint64) ([]models.MediaSyncJob, error) {
	var syncJobs []models.MediaSyncJob
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&syncJobs)
	
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving media sync jobs for user %d: %w", userID, result.Error)
	}
	return syncJobs, nil
}

// UpdateMediaSyncJob updates an existing media sync job
func (r *jobRepository) UpdateMediaSyncJob(ctx context.Context, syncJob *models.MediaSyncJob) error {
	result := r.db.WithContext(ctx).Save(syncJob)
	if result.Error != nil {
		return fmt.Errorf("error updating media sync job: %w", result.Error)
	}
	return nil
}

// UpdateMediaSyncLastRunTime updates the last run time for a media sync job
func (r *jobRepository) UpdateMediaSyncLastRunTime(ctx context.Context, syncJobID uint64, lastRunTime time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&models.MediaSyncJob{}).
		Where("id = ?", syncJobID).
		Update("last_sync_time", lastRunTime)
	
	if result.Error != nil {
		return fmt.Errorf("error updating media sync job last run time: %w", result.Error)
	}
	return nil
}

// DeleteMediaSyncJob deletes a media sync job
func (r *jobRepository) DeleteMediaSyncJob(ctx context.Context, syncJobID uint64) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", syncJobID).
		Delete(&models.MediaSyncJob{})
	
	if result.Error != nil {
		return fmt.Errorf("error deleting media sync job: %w", result.Error)
	}
	return nil
}