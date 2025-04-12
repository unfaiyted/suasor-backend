package requests

// GetRecommendationsRequest represents a request to retrieve recommendations
type GetRecommendationsRequest struct {
	MediaType string `form:"mediaType" binding:"omitempty,oneof=movie series music" example:"movie"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100" example:"20"`
	Offset    int    `form:"offset" binding:"omitempty,min=0" example:"0"`
}

// GetRecentRecommendationsRequest represents a request to retrieve recent recommendations
type GetRecentRecommendationsRequest struct {
	Days      int    `form:"days" binding:"omitempty,min=1,max=90" example:"7"`
	MediaType string `form:"mediaType" binding:"omitempty,oneof=movie series music" example:"movie"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100" example:"20"`
}

// GetTopRecommendationsRequest represents a request to retrieve top recommendations
type GetTopRecommendationsRequest struct {
	MinScore  float32 `form:"minScore" binding:"omitempty,min=0,max=1" example:"0.7"`
	MediaType string  `form:"mediaType" binding:"omitempty,oneof=movie series music" example:"movie"`
	Limit     int     `form:"limit" binding:"omitempty,min=1,max=100" example:"20"`
}

// MarkRecommendationAsViewedRequest represents a request to mark a recommendation as viewed
type MarkRecommendationAsViewedRequest struct {
	RecommendationID uint64 `json:"recommendationId" binding:"required" example:"123"`
}

// RateRecommendationRequest represents a request to rate a recommendation
type RateRecommendationRequest struct {
	RecommendationID uint64  `json:"recommendationId" binding:"required" example:"123"`
	Rating           float32 `json:"rating" binding:"required,min=0,max=5" example:"4.5"`
}