package jobs

import (
	"context"
	"log"
	"time"
)

// EmptyJob is a placeholder job implementation that satisfies scheduler.Job interface
// It doesn't do anything when executed but allows the scheduler to run without errors
type EmptyJob struct {
	JobName string
}

// Execute implements the Job interface
func (e *EmptyJob) Execute(ctx context.Context) error {
	log.Printf("Empty job %s executed (no-op)", e.JobName)
	return nil
}

// Name returns the job name
func (e *EmptyJob) Name() string {
	return e.JobName
}

// Schedule returns how often the job should run
func (e *EmptyJob) Schedule() time.Duration {
	// Default to a daily schedule
	return 24 * time.Hour
}
