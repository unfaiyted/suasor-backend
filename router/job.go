package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

// RegisterJobRoutes registers job-related routes
func RegisterJobRoutes(r *gin.RouterGroup, jobService services.JobService) {
	jobHandler := handlers.NewJobHandler(jobService)

	jobs := r.Group("/jobs")
	{
		// Job schedules
		jobs.GET("/schedules", jobHandler.GetAllJobSchedules)
		jobs.GET("/schedules/:name", jobHandler.GetJobScheduleByName)
		jobs.PUT("/schedules", jobHandler.UpdateJobSchedule)
		jobs.POST("/schedules", jobHandler.CreateJobSchedule)

		// Job runs
		jobs.GET("/runs", jobHandler.GetRecentJobRuns)
		jobs.GET("/runs/:id/progress", jobHandler.GetJobRunProgress)
		jobs.GET("/active", jobHandler.GetActiveJobRuns)
		// jobs.GET("/runs/user", jobHandler.GetUserJobRuns)
		jobs.POST("/:name/run", jobHandler.RunJobManually)

		// Media sync jobs
		jobs.GET("/media-sync", jobHandler.GetMediaSyncJobs)
		jobs.POST("/media-sync", jobHandler.SetupMediaSyncJob)
		jobs.POST("/media-sync/run", jobHandler.RunMediaSyncJob)

		// Recommendations
		jobs.GET("/recommendations", jobHandler.GetUserRecommendations)
		jobs.POST("/recommendations/:id/dismiss", jobHandler.DismissRecommendation)
		jobs.PUT("/recommendations/:id/viewed", jobHandler.UpdateRecommendationViewedStatus)
	}
}

