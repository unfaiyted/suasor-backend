package router

import (
	"suasor/di/container"
	"suasor/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterJobRoutes registers job-related routes
func RegisterJobRoutes(r *gin.RouterGroup, c *container.Container) {
	jobHandler := container.MustGet[*handlers.JobHandler](c)

	jobs := r.Group("/jobs")
	{
		// Job schedules
		jobs.GET("/schedules", jobHandler.GetAllJobSchedules)
		jobs.GET("/schedules/:name", jobHandler.GetJobScheduleByName)
		jobs.PUT("/schedules", jobHandler.UpdateJobSchedule)
		jobs.POST("/schedules", jobHandler.CreateJobSchedule)

		// Job runs
		jobs.GET("/runs", jobHandler.GetRecentJobRuns)
		jobs.GET("/runs/:jobID/progress", jobHandler.GetJobRunProgress)
		jobs.GET("/active", jobHandler.GetActiveJobRuns)
		// jobs.GET("/runs/user", jobHandler.GetUserJobRuns)
		jobs.POST("/:name/run", jobHandler.RunJobManually)

		// Media sync jobs
		jobs.GET("/media-sync", jobHandler.GetMediaSyncJobs)
		jobs.POST("/media-sync", jobHandler.SetupMediaSyncJob)
		jobs.POST("/media-sync/run", jobHandler.RunMediaSyncJob)
		jobs.POST("/full-sync", jobHandler.RunFullSync)  // New endpoint for full sync

		// Recommendations
		jobs.GET("/recommendations", jobHandler.GetUserRecommendations)
		jobs.POST("/recommendations/:jobID/dismiss", jobHandler.DismissRecommendation)
		jobs.PUT("/recommendations/:jobID/viewed", jobHandler.UpdateRecommendationViewedStatus)
	}
}
