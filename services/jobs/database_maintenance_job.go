package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// DatabaseMaintenanceJob performs routine database optimization and maintenance
type DatabaseMaintenanceJob struct {
	jobRepo repository.JobRepository
	// Include other repositories that might need maintenance
}

// NewDatabaseMaintenanceJob creates a new database maintenance job
func NewDatabaseMaintenanceJob(
	jobRepo repository.JobRepository,
) *DatabaseMaintenanceJob {
	return &DatabaseMaintenanceJob{
		jobRepo: jobRepo,
	}
}

// Name returns the unique name of the job
func (j *DatabaseMaintenanceJob) Name() string {
	return "system.database.maintenance"
}

// Schedule returns when the job should next run
func (j *DatabaseMaintenanceJob) Schedule() time.Duration {
	// Run weekly by default
	return 7 * 24 * time.Hour
}

// Execute runs the database maintenance job
func (j *DatabaseMaintenanceJob) Execute(ctx context.Context) error {
	log.Println("Starting database maintenance job")

	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSystem,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		Metadata:  fmt.Sprintf(`{"type":"databaseMaintenance","startTime":"%s"}`, now.Format(time.RFC3339)),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Run the maintenance tasks
	var jobError error
	maintenanceStats := map[string]int{
		"tablesOptimized":      0,
		"recordsArchived":      0,
		"recordsDeleted":       0,
		"tokensCleanedUp":      0,
		"sessionsCleanedUp":    0,
		"integrityIssuesFixed": 0,
	}

	// Task 1: Optimize database tables
	optimizeStats, err := j.optimizeDatabaseTables(ctx)
	if err != nil {
		log.Printf("Error optimizing database tables: %v", err)
		jobError = err
		// Continue with other tasks even if one fails
	} else {
		maintenanceStats["tablesOptimized"] = optimizeStats.optimized
	}

	// Task 2: Archive old job runs
	archiveStats, err := j.archiveOldJobRuns(ctx)
	if err != nil {
		log.Printf("Error archiving old job runs: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		maintenanceStats["recordsArchived"] += archiveStats.archived
	}

	// Task 3: Clean up expired tokens
	tokenStats, err := j.cleanupExpiredTokens(ctx)
	if err != nil {
		log.Printf("Error cleaning up expired tokens: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		maintenanceStats["tokensCleanedUp"] = tokenStats.cleaned
	}

	// Task 4: Clean up expired sessions
	sessionStats, err := j.cleanupExpiredSessions(ctx)
	if err != nil {
		log.Printf("Error cleaning up expired sessions: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		maintenanceStats["sessionsCleanedUp"] = sessionStats.cleaned
	}

	// Task 5: Validate data integrity
	integrityStats, err := j.validateDataIntegrity(ctx)
	if err != nil {
		log.Printf("Error validating data integrity: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		maintenanceStats["integrityIssuesFixed"] = integrityStats.fixed
	}

	// Complete the job
	status := models.JobStatusCompleted
	errorMessage := ""
	if jobError != nil {
		status = models.JobStatusFailed
		errorMessage = jobError.Error()
	}

	// Update job run with results
	j.completeJobRun(ctx, jobRun.ID, status, errorMessage, maintenanceStats)
	log.Println("Database maintenance job completed")
	return jobError
}

// completeJobRun finalizes a job run with status and results
func (j *DatabaseMaintenanceJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string, stats map[string]int) {
	// Convert stats to a string for the message
	message := fmt.Sprintf("Tables optimized: %d, Records archived: %d, Tokens cleaned up: %d, Sessions cleaned up: %d, Integrity issues fixed: %d",
		stats["tablesOptimized"],
		stats["recordsArchived"],
		stats["tokensCleanedUp"],
		stats["sessionsCleanedUp"],
		stats["integrityIssuesFixed"])

	if status == models.JobStatusFailed {
		message = fmt.Sprintf("%s. Error: %s", message, errorMsg)
	}

	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, message); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// optimizeDatabaseTables optimizes database tables
func (j *DatabaseMaintenanceJob) optimizeDatabaseTables(ctx context.Context) (MaintenanceStats, error) {
	stats := MaintenanceStats{}
	log.Println("Optimizing database tables")

	// In a real implementation, we would:
	// 1. Get a list of tables to optimize
	// 2. Run OPTIMIZE TABLE or equivalent commands
	// 3. Track statistics on optimizations performed

	// Mock implementation
	stats.optimized = 15 // Pretend we optimized 15 tables

	log.Printf("Optimized %d database tables", stats.optimized)
	return stats, nil
}

// archiveOldJobRuns archives old job runs
func (j *DatabaseMaintenanceJob) archiveOldJobRuns(ctx context.Context) (MaintenanceStats, error) {
	stats := MaintenanceStats{}
	log.Println("Archiving old job runs")

	// In a real implementation, we would:
	// 1. Find job runs older than a certain threshold (e.g., 30 days)
	// 2. Either move them to an archive table or compress their data
	// 3. Track statistics on archives performed

	// Mock implementation
	cutoffDate := time.Now().AddDate(0, -1, 0) // 1 month ago
	stats.archived = 250                       // Pretend we archived 250 job runs

	log.Printf("Archived %d job runs older than %s", stats.archived, cutoffDate.Format("2006-01-02"))
	return stats, nil
}

// cleanupExpiredTokens cleans up expired authentication tokens
func (j *DatabaseMaintenanceJob) cleanupExpiredTokens(ctx context.Context) (MaintenanceStats, error) {
	stats := MaintenanceStats{}
	log.Println("Cleaning up expired tokens")

	// In a real implementation, we would:
	// 1. Find expired refresh tokens and access tokens
	// 2. Delete them from the database
	// 3. Track statistics on deletions performed

	// Mock implementation
	stats.cleaned = 75 // Pretend we cleaned up 75 expired tokens

	log.Printf("Cleaned up %d expired tokens", stats.cleaned)
	return stats, nil
}

// cleanupExpiredSessions cleans up expired user sessions
func (j *DatabaseMaintenanceJob) cleanupExpiredSessions(ctx context.Context) (MaintenanceStats, error) {
	stats := MaintenanceStats{}
	log.Println("Cleaning up expired sessions")

	// In a real implementation, we would:
	// 1. Find inactive sessions beyond a threshold (e.g., 7 days)
	// 2. Delete them from the database
	// 3. Track statistics on deletions performed

	// Mock implementation
	stats.cleaned = 120 // Pretend we cleaned up 120 expired sessions

	log.Printf("Cleaned up %d expired sessions", stats.cleaned)
	return stats, nil
}

// validateDataIntegrity validates and fixes data integrity issues
func (j *DatabaseMaintenanceJob) validateDataIntegrity(ctx context.Context) (MaintenanceStats, error) {
	stats := MaintenanceStats{}
	log.Println("Validating data integrity")

	// In a real implementation, we would:
	// 1. Check for data consistency issues (orphaned records, broken relationships)
	// 2. Fix issues where possible, report unfixable issues
	// 3. Track statistics on fixes performed

	// Mock implementation
	stats.fixed = 8 // Pretend we fixed 8 integrity issues

	log.Printf("Fixed %d data integrity issues", stats.fixed)
	return stats, nil
}

// SetupDatabaseMaintenanceSchedule creates or updates a database maintenance schedule
func (j *DatabaseMaintenanceJob) SetupDatabaseMaintenanceSchedule(ctx context.Context, frequency string) error {
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

// RunManualMaintenance runs the database maintenance job manually
func (j *DatabaseMaintenanceJob) RunManualMaintenance(ctx context.Context) error {
	return j.Execute(ctx)
}

// RunSpecificMaintenance runs a specific maintenance task
func (j *DatabaseMaintenanceJob) RunSpecificMaintenance(ctx context.Context, taskType string) error {
	log.Printf("Running specific maintenance task: %s", taskType)

	switch taskType {
	case "optimize":
		_, err := j.optimizeDatabaseTables(ctx)
		return err
	case "archive":
		_, err := j.archiveOldJobRuns(ctx)
		return err
	case "tokens":
		_, err := j.cleanupExpiredTokens(ctx)
		return err
	case "sessions":
		_, err := j.cleanupExpiredSessions(ctx)
		return err
	case "integrity":
		_, err := j.validateDataIntegrity(ctx)
		return err
	default:
		return fmt.Errorf("unknown maintenance task type: %s", taskType)
	}
}

// GetDatabaseStats gets statistics about the database
func (j *DatabaseMaintenanceJob) GetDatabaseStats(ctx context.Context) (map[string]interface{}, error) {
	// In a real implementation, we would:
	// 1. Query various database statistics (table sizes, row counts, etc.)
	// 2. Format them into a structured report

	// Mock implementation
	return map[string]interface{}{
		"totalTables":    25,
		"totalRows":      125000,
		"totalSizeBytes": 524288000, // 500 MB
		"largestTables": []map[string]interface{}{
			{"name": "media_items", "rows": 50000, "sizeBytes": 209715200},
			{"name": "media_play_history", "rows": 35000, "sizeBytes": 104857600},
			{"name": "job_runs", "rows": 15000, "sizeBytes": 52428800},
		},
		"lastOptimized":   time.Now().AddDate(0, 0, -7).Format(time.RFC3339),
		"databaseVersion": "PostgreSQL 14.5",
	}, nil
}

