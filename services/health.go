// In services/health.go
package services

import (
	"time"

	"gorm.io/gorm"
)

type HealthService interface {
	CheckApplicationStatus() bool
	CheckDatabaseConnection() bool
}

// HealthServiceImpl implements the HealthService interface
type HealthServiceImpl struct {
	db        *gorm.DB
	startTime time.Time
}

// NewHealthService creates a new instance of HealthService
func NewHealthService(db *gorm.DB) HealthService {
	return &HealthServiceImpl{
		db:        db,
		startTime: time.Now(),
	}
}

// CheckApplicationStatus verifies that the application is running correctly
func (h *HealthServiceImpl) CheckApplicationStatus() bool {
	// Simple check - application is considered healthy if it's running
	// We could add more sophisticated checks here (memory usage, goroutine count, etc.)
	return true
}

// CheckDatabaseConnection verifies database connectivity
func (h *HealthServiceImpl) CheckDatabaseConnection() bool {
	// Attempt to ping the database
	sqlDB, err := h.db.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
}

// GetUptime returns the duration the application has been running
func (h *HealthServiceImpl) GetUptime() time.Duration {
	return time.Since(h.startTime)
}
