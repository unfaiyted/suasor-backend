package jobs

import (
	"context"
	"log"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// SystemJobConfig holds the configuration for a system job
type SystemJobConfig struct {
	Name        string
	Type        models.JobType
	Frequency   string
	Description string
	Enabled     bool
}

// RegisterSystemJobs registers all system jobs in the database during application startup
func RegisterSystemJobs(ctx context.Context, jobRepo repository.JobRepository, systemJobs []scheduler.Job) error {
	log.Println("Registering system jobs...")
	
	// Create a map of job configurations with defaults
	defaultJobs := map[string]SystemJobConfig{
		"system.database.maintenance": {
			Name:        "system.database.maintenance",
			Type:        models.JobTypeSystem,
			Frequency:   string(scheduler.FrequencyWeekly),
			Description: "Performs routine database maintenance, cleanup, and optimization",
			Enabled:     true,
		},
		"system.content.availability": {
			Name:        "system.content.availability",
			Type:        models.JobTypeSystem,
			Frequency:   string(scheduler.FrequencyDaily),
			Description: "Monitors changes in content availability across streaming services",
			Enabled:     true,
		},
		"system.metadata.refresh": {
			Name:        "system.metadata.refresh",
			Type:        models.JobTypeSystem,
			Frequency:   string(scheduler.FrequencyDaily),
			Description: "Refreshes metadata for media items in the library",
			Enabled:     true,
		},
		"system.library.cleanup": {
			Name:        "system.library.cleanup",
			Type:        models.JobTypeSystem,
			Frequency:   string(scheduler.FrequencyWeekly),
			Description: "Cleans up the media library by removing stale references",
			Enabled:     true,
		},
		"system.user.activity.analysis": {
			Name:        "system.user.activity.analysis",
			Type:        models.JobTypeAnalysis,
			Frequency:   string(scheduler.FrequencyWeekly),
			Description: "Analyzes user activity patterns to improve recommendations",
			Enabled:     true,
		},
		"system.new.release.notification": {
			Name:        "system.new.release.notification",
			Type:        models.JobTypeNotification,
			Frequency:   string(scheduler.FrequencyDaily),
			Description: "Notifies users of new content releases relevant to their interests",
			Enabled:     true,
		},
		"system.playlist.sync": {
			Name:        "system.playlist.sync",
			Type:        models.JobTypeSync,
			Frequency:   string(scheduler.FrequencyDaily),
			Description: "Synchronizes playlists across multiple media servers",
			Enabled:     true,
		},
		"system.smart.collection": {
			Name:        "system.smart.collection",
			Type:        models.JobTypeSystem,
			Frequency:   string(scheduler.FrequencyDaily),
			Description: "Updates smart collections based on configured rules",
			Enabled:     true,
		},
		"system.recommendation": {
			Name:        "system.recommendation",
			Type:        models.JobTypeRecommendation,
			Frequency:   string(scheduler.FrequencyDaily),
			Description: "Generates personalized content recommendations for users",
			Enabled:     true,
		},
		"system.watch.history.sync": {
			Name:        "system.watch.history.sync",
			Type:        models.JobTypeSync,
			Frequency:   string(scheduler.FrequencyDaily),
			Description: "Synchronizes watch history across multiple media servers",
			Enabled:     true,
		},
		"system.favorites.sync": {
			Name:        "system.favorites.sync",
			Type:        models.JobTypeSync,
			Frequency:   string(scheduler.FrequencyDaily),
			Description: "Synchronizes favorite items across multiple media servers",
			Enabled:     true,
		},
	}

	// Get existing job schedules from the database
	existingSchedules, err := jobRepo.GetAllJobSchedules(ctx)
	if err != nil {
		return err
	}

	// Create a map of existing job schedules for easy lookup
	existingMap := make(map[string]*models.JobSchedule)
	for i, schedule := range existingSchedules {
		existingMap[schedule.JobName] = &existingSchedules[i]
	}

	// Register all system jobs
	for _, job := range systemJobs {
		jobName := job.Name()
		
		// Check if we have a default configuration for this job
		config, hasConfig := defaultJobs[jobName]
		if !hasConfig {
			// If no default config exists, use basic defaults
			config = SystemJobConfig{
				Name:        jobName,
				Type:        models.JobTypeSystem,
				Frequency:   string(scheduler.FrequencyDaily),
				Description: "System job",
				Enabled:     true,
			}
		}
		
		// Check if the job already exists in the database
		existingJob, exists := existingMap[jobName]
		if exists {
			// Job already exists, update if needed but preserve user settings
			log.Printf("System job already exists: %s", jobName)
			
			// Only update if necessary (to preserve user modifications)
			needsUpdate := false
			
			// Don't override user-configured frequency
			// Don't override enabled/disabled state
			
			if needsUpdate {
				err := jobRepo.UpdateJobSchedule(ctx, existingJob)
				if err != nil {
					log.Printf("Error updating job schedule for %s: %v", jobName, err)
				}
			}
		} else {
			// Job doesn't exist, create it
			log.Printf("Creating system job: %s", jobName)
			
			newSchedule := &models.JobSchedule{
				JobName:     config.Name,
				JobType:     config.Type,
				Frequency:   config.Frequency,
				Enabled:     config.Enabled,
				LastRunTime: nil, // Never run yet
				Config:      "{\"description\":\"" + config.Description + "\"}",
			}
			
			err := jobRepo.CreateJobSchedule(ctx, newSchedule)
			if err != nil {
				log.Printf("Error creating job schedule for %s: %v", jobName, err)
			}
		}
	}
	
	log.Println("System jobs registration completed")
	return nil
}