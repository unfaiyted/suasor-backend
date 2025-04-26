// app/di/handlers/job.go
package handlers

import (
	"context"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/services/jobs"
)

// RegisterJobHandlers registers job-related handlers
func RegisterJobHandlers(ctx context.Context, c *container.Container) {
	container.RegisterFactory[*handlers.JobHandler](c, func(c *container.Container) *handlers.JobHandler {
		jobService := container.MustGet[jobs.JobService](c)
		return handlers.NewJobHandler(jobService)
	})
}
