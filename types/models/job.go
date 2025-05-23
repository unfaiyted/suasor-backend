package models

import (
	clienttypes "suasor/clients/types"
	"time"
)

type SyncType string

const (
	// Media types
	SyncTypeMovies SyncType = "movies"
	SyncTypeSeries SyncType = "series"
	SyncTypeMusic  SyncType = "music"

	// Other types
	SyncTypeHistory     SyncType = "history"
	SyncTypeFavorites   SyncType = "favorites"
	SyncTypeCollections SyncType = "collections"
	SyncTypePlaylists   SyncType = "playlists"
	
	// Special types
	SyncTypeFull        SyncType = "full"    // Syncs everything from all clients
)

func (s SyncType) String() string {
	return string(s)
}

func (s SyncType) IsValid() bool {
	switch s {
	case SyncTypeMovies, SyncTypeSeries, SyncTypeMusic,
		SyncTypeHistory, SyncTypeFavorites, SyncTypeCollections, SyncTypePlaylists,
		SyncTypeFull:
		return true
	default:
		return false
	}
}

// JobType defines the type of scheduled job
type JobType string

const (
	// JobTypeRecommendation represents a recommendation generation job
	JobTypeRecommendation JobType = "recommendation"
	// JobTypeSync represents a media synchronization job
	JobTypeSync JobType = "sync"
	// JobTypeSystem represents a system maintenance job
	JobTypeSystem JobType = "system"
	// JobTypeNotification represents a notification job
	JobTypeNotification JobType = "notification"
	// JobTypeAnalysis represents an analysis job
	JobTypeAnalysis JobType = "analysis"
)

// JobStatus defines the status of a job run
type JobStatus string

const (
	// JobStatusPending job is scheduled but hasn't run yet
	JobStatusPending JobStatus = "pending"
	// JobStatusRunning job is currently running
	JobStatusRunning JobStatus = "running"
	// JobStatusCompleted job completed successfully
	JobStatusCompleted JobStatus = "completed"
	// JobStatusFailed job failed to complete
	JobStatusFailed JobStatus = "failed"
)

// JobRun represents a single execution of a scheduled job
type JobRun struct {
	BaseModel
	// The name of the job
	JobName string `json:"jobName" gorm:"index;not null"`
	// Type of job (recommendation, sync, etc.)
	JobType JobType `json:"jobType" gorm:"index;not null"`
	// Status of the job run
	Status JobStatus `json:"status" gorm:"not null"`
	// When the job started running
	StartTime *time.Time `json:"startTime"`
	// When the job finished running
	EndTime *time.Time `json:"endTime"`
	// User ID associated with the job, if applicable
	UserID *uint64 `json:"userID" gorm:"index"`
	// Any error message from the job run
	ErrorMessage string `json:"errorMessage"`
	// Progress percentage (0-100)
	Progress int `json:"progress" gorm:"not null;default:0"`
	// Total items to process
	TotalItems int `json:"totalItems" gorm:"default:0"`
	// Items processed so far
	ProcessedItems int `json:"processedItems" gorm:"default:0"`
	// Current status message
	StatusMessage string `json:"statusMessage"`
	// Metadata related to the job (stored as JSON)
	Metadata string `json:"metadata" gorm:"type:jsonb"`
}

// JobSchedule represents a scheduled job
type JobSchedule struct {
	BaseModel
	// Unique name of the job
	JobName string `json:"jobName" gorm:"uniqueIndex;not null"`
	// Type of job (recommendation, sync, etc.)
	JobType JobType `json:"jobType" gorm:"index;not null"`
	// How often the job should run
	Frequency string `json:"frequency" gorm:"not null"`
	// When the job last ran
	LastRunTime *time.Time `json:"lastRunTime"`
	// Whether the job is enabled
	Enabled bool `json:"enabled" gorm:"not null;default:true"`
	// User ID associated with the job, if applicable (for user-specific jobs)
	UserID *uint64 `json:"userID" gorm:"index"`
	// Any configuration for the job (stored as JSON)
	Config string `json:"config" gorm:"type:jsonb"`
}

// MediaSyncJob represents a job to sync media from external clients
type MediaSyncJob struct {
	BaseModel
	// ID of the user
	UserID uint64 `json:"userID" gorm:"index;not null"`
	// ID of the client to sync from
	ClientID uint64 `json:"clientID" gorm:"index;not null"`
	// Type of the client
	ClientType clienttypes.ClientType `json:"clientType" gorm:"index;not null"`
	// Type of media to sync (movies, series, music, etc.)
	SyncType SyncType `json:"syncType" gorm:"index;not null"`
	// Last sync time
	LastSyncTime *time.Time `json:"lastSyncTime"`
	// Sync frequency
	Frequency string `json:"frequency" gorm:"not null;default:'daily'"`
	// Whether sync is enabled
	Enabled bool `json:"enabled" gorm:"not null;default:true"`
	// Sync filter criteria (stored as JSON)
	Filters string `json:"filters" gorm:"type:jsonb;default:'{}'"`
}
