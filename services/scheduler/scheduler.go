package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

// Job represents a scheduled job that can be executed
type Job interface {
	// Execute runs the job with the given context
	Execute(ctx context.Context) error
	// Name returns the unique name of the job
	Name() string
	// Schedule returns when the job should next run
	Schedule() time.Duration
}

// Scheduler manages the execution of scheduled jobs
type Scheduler struct {
	jobs       map[string]Job
	jobTimers  map[string]*time.Timer
	mutex      sync.Mutex
	cancelFunc context.CancelFunc
	ctx        context.Context
	wg         sync.WaitGroup
}

// NewScheduler creates a new job scheduler
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		jobs:       make(map[string]Job),
		jobTimers:  make(map[string]*time.Timer),
		cancelFunc: cancel,
		ctx:        ctx,
	}
}

// RegisterJob adds a job to the scheduler
func (s *Scheduler) RegisterJob(job Job) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.jobs[job.Name()] = job
}

// Start begins the scheduler, executing all registered jobs according to their schedule
func (s *Scheduler) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.jobs) == 0 {
		log.Println("Starting scheduler with no registered jobs")
		return
	}

	for name, job := range s.jobs {
		s.scheduleJob(name, job, 0) // Schedule immediately for the first run
	}
}

// Stop cancels all scheduled jobs and waits for any executing jobs to complete
func (s *Scheduler) Stop() {
	s.cancelFunc()

	// Cancel all timers
	s.mutex.Lock()
	for _, timer := range s.jobTimers {
		if timer != nil {
			timer.Stop()
		}
	}
	s.mutex.Unlock()

	// Wait for all jobs to complete
	s.wg.Wait()
}

// scheduleJob creates a timer for the next execution of a job
func (s *Scheduler) scheduleJob(name string, job Job, delay time.Duration) {
	// Stop existing timer if any
	if timer, exists := s.jobTimers[name]; exists && timer != nil {
		timer.Stop()
	}

	// Use the job's schedule if no delay is specified
	if delay == 0 {
		delay = job.Schedule()
	}

	s.jobTimers[name] = time.AfterFunc(delay, func() {
		s.executeJob(name, job)
	})
}

// executeJob runs a job and reschedules it
func (s *Scheduler) executeJob(name string, job Job) {
	// Mark that we're executing a job
	s.wg.Add(1)
	defer s.wg.Done()

	// Create a context that will be canceled if the scheduler is stopped
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Minute)
	defer cancel()

	// Execute the job
	if err := job.Execute(ctx); err != nil {
		log.Printf("Error executing job %s: %v", name, err)
	}

	// Reschedule the job if scheduler hasn't been stopped
	s.mutex.Lock()
	defer s.mutex.Unlock()

	select {
	case <-s.ctx.Done():
		// Scheduler has been stopped, don't reschedule
		return
	default:
		// Reschedule the job
		s.scheduleJob(name, job, job.Schedule())
	}
}