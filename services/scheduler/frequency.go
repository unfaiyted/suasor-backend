package scheduler

import (
	"time"
)

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

// ToDuration converts a frequency string to a time.Duration
func (f Frequency) ToDuration() time.Duration {
	switch f {
	case FrequencyDaily:
		return 24 * time.Hour
	case FrequencyWeekly:
		return 7 * 24 * time.Hour
	case FrequencyMonthly:
		return 30 * 24 * time.Hour // Approximation
	default:
		// For manual frequency or any unrecognized value, use a high value
		// that effectively means "never auto-run"
		return 365 * 24 * time.Hour
	}
}

// NextRunTime calculates when a job with this frequency should next run
// based on the last run time
func (f Frequency) NextRunTime(lastRun time.Time) time.Time {
	switch f {
	case FrequencyDaily:
		// Run at the same time tomorrow
		return lastRun.AddDate(0, 0, 1)
	case FrequencyWeekly:
		// Run on the same day next week
		return lastRun.AddDate(0, 0, 7)
	case FrequencyMonthly:
		// Run on the same day next month
		return lastRun.AddDate(0, 1, 0)
	default:
		// For manual or unrecognized frequency, set a far future time
		return lastRun.AddDate(1, 0, 0)
	}
}

// ShouldRunNow determines if a job should run now based on its last run time and frequency
func (f Frequency) ShouldRunNow(lastRun time.Time) bool {
	// Manual jobs should never auto-run
	if f == FrequencyManual {
		return false
	}

	// If job has never run, run it now
	if lastRun.IsZero() {
		return true
	}

	// Check if enough time has passed since the last run
	return time.Now().After(f.NextRunTime(lastRun))
}