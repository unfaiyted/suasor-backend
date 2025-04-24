package requests

// UpdateJobScheduleRequest represents a request to update a job schedule
type UpdateJobScheduleRequest struct {
	JobName   string `json:"jobName" binding:"required"`
	Frequency string `json:"frequency" binding:"required"`
	Enabled   bool   `json:"enabled"`
}

// SetupMediaSyncJobRequest represents a request to setup a media sync job
type SetupMediaSyncJobRequest struct {
	ClientID   uint64 `json:"clientID" binding:"required"`
	ClientType string `json:"clientType" binding:"required"`
	MediaType  string `json:"mediaType" binding:"required"`
	Frequency  string `json:"frequency" binding:"required"`
}

// RunMediaSyncJobRequest represents a request to run a media sync job
type RunMediaSyncJobRequest struct {
	ClientID  uint64 `json:"clientID" binding:"required"`
	MediaType string `json:"mediaType" binding:"required"`
}

// UpdateRecommendationViewedRequest represents a request to update recommendation viewed status
type UpdateRecommendationViewedRequest struct {
	Viewed bool `json:"viewed"`
}

