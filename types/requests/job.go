package requests

import "suasor/types/models"

// UpdateJobScheduleRequest represents a request to update a job schedule
type UpdateJobScheduleRequest struct {
	JobName   string `json:"jobName" binding:"required"`
	Frequency string `json:"frequency" binding:"required"`
	Enabled   bool   `json:"enabled"`
}

// SetupMediaSyncJobRequest represents a request to setup a media sync job
type SetupMediaSyncJobRequest struct {
	ClientID   uint64          `json:"clientID" binding:"required"`
	ClientType string          `json:"clientType" binding:"required"`
	SyncType   models.SyncType `json:"syncType" binding:"required"`
	Frequency  string          `json:"frequency" binding:"required"`
}

// RunMediaSyncJobRequest represents a request to run a media sync job
type RunMediaSyncJobRequest struct {
	ClientID uint64          `json:"clientID" binding:"required"`
	SyncType models.SyncType `json:"syncType" binding:"required"`
}

// RunFullSyncRequest represents a request to run a full sync for all clients
type RunFullSyncRequest struct {
	// No parameters needed - we'll use the user ID from the context
}

// UpdateRecommendationViewedRequest represents a request to update recommendation viewed status
type UpdateRecommendationViewedRequest struct {
	Viewed bool `json:"viewed"`
}
