package responses

import (
	"suasor/client/media/types"
	"suasor/types/models"
	"time"
)

// RecommendationResponse represents a recommendation in the API response
type RecommendationResponse struct {
	ID               uint64            `json:"id"`
	Title            string            `json:"title"`
	Year             int               `json:"year,omitempty"`
	MediaType        types.MediaType   `json:"mediaType"`
	Genres           []string          `json:"genres,omitempty"`
	Score            float32           `json:"score"`
	Reasoning        string            `json:"reasoning"`
	SimilarItems     []string          `json:"similarItems,omitempty"`
	MatchesActors    []string          `json:"matchesActors,omitempty"`
	MatchesDirectors []string          `json:"matchesDirectors,omitempty"`
	MatchesGenres    []string          `json:"matchesGenres,omitempty"`
	RecommendedBy    string            `json:"recommendedBy"`
	AIModel          string            `json:"aiModel,omitempty"`
	CreatedAt        time.Time         `json:"createdAt"`
	ExternalIDs      map[string]string `json:"externalIds,omitempty"`
	IsViewed         bool              `json:"isViewed"`
	UserRating       float32           `json:"userRating,omitempty"`
}

// RecommendationsListResponse represents a paginated list of recommendations
type RecommendationsListResponse struct {
	Recommendations []RecommendationResponse `json:"recommendations"`
	Total           int                      `json:"total"`
	Limit           int                      `json:"limit"`
	Offset          int                      `json:"offset"`
	MediaType       string                   `json:"mediaType,omitempty"`
}

// ConvertToRecommendationResponse converts a model to a response
func ConvertToRecommendationResponse(recommendation models.Recommendation) RecommendationResponse {
	response := RecommendationResponse{
		ID:               recommendation.ID,
		Title:            recommendation.Title,
		Year:             recommendation.Year,
		MediaType:        recommendation.MediaType,
		Genres:           []string(recommendation.Genres),
		Score:            recommendation.Confidence,
		Reasoning:        recommendation.Reasoning,
		SimilarItems:     []string(recommendation.SimilarItems),
		MatchesActors:    []string(recommendation.MatchesActors),
		MatchesDirectors: []string(recommendation.MatchesDirectors),
		MatchesGenres:    []string(recommendation.MatchesGenres),
		RecommendedBy:    recommendation.RecommendedBy,
		AIModel:          recommendation.AIModel,
		CreatedAt:        recommendation.CreatedAt,
		IsViewed:         recommendation.IsViewed,
		UserRating:       recommendation.UserRating,
	}

	// Convert ExternalIDs if present
	if recommendation.ExternalIDs != nil {
		response.ExternalIDs = map[string]string(*recommendation.ExternalIDs)
	}

	return response
}

// ConvertToRecommendationsListResponse converts a slice of models to a response
func ConvertToRecommendationsListResponse(
	recommendations []models.Recommendation,
	total int,
	limit int,
	offset int,
	mediaType string,
) RecommendationsListResponse {
	response := RecommendationsListResponse{
		Recommendations: make([]RecommendationResponse, len(recommendations)),
		Total:           total,
		Limit:           limit,
		Offset:          offset,
		MediaType:       mediaType,
	}

	for i, rec := range recommendations {
		response.Recommendations[i] = ConvertToRecommendationResponse(rec)
	}

	return response
}

