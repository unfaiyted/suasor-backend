package services

import (
	"context"
	"fmt"
	"time"

	"suasor/repository"
	"suasor/types/models"
)

// RecommendationService defines the interface for recommendation operations
type RecommendationService interface {
	// GetRecommendations retrieves recommendations for a user with optional filtering
	GetRecommendations(ctx context.Context, userID uint64, mediaType string, limit, offset int) ([]models.Recommendation, error)
	// GetRecommendationByID retrieves a specific recommendation by ID
	GetRecommendationByID(ctx context.Context, id uint64) (*models.Recommendation, error)
	// MarkRecommendationAsViewed marks a recommendation as viewed
	MarkRecommendationAsViewed(ctx context.Context, id uint64, userID uint64) error
	// RateRecommendation sets a user rating for a recommendation
	RateRecommendation(ctx context.Context, id uint64, userID uint64, rating float32) error
	// StoreRecommendations stores recommendations for a user
	StoreRecommendations(ctx context.Context, recommendations []*models.Recommendation) error
	// GetRecentRecommendations retrieves recently created recommendations for a user
	GetRecentRecommendations(ctx context.Context, userID uint64, days int, limit int) ([]models.Recommendation, error)
	// GetTopRecommendations retrieves top-scored recommendations for a user
	GetTopRecommendations(ctx context.Context, userID uint64, minScore float32, limit int) ([]models.Recommendation, error)
	// GetRecommendationCount returns the count of recommendations for a user
	GetRecommendationCount(ctx context.Context, userID uint64, mediaType string) (int, error)
	// DeleteExpiredRecommendations deletes all expired recommendations
	DeleteExpiredRecommendations(ctx context.Context) error
}

// recommendationService implements the RecommendationService interface
type recommendationService struct {
	recommendationRepo repository.RecommendationRepository
}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService(recommendationRepo repository.RecommendationRepository) RecommendationService {
	return &recommendationService{
		recommendationRepo: recommendationRepo,
	}
}

// GetRecommendations retrieves recommendations for a user with optional filtering
func (s *recommendationService) GetRecommendations(ctx context.Context, userID uint64, mediaType string, limit, offset int) ([]models.Recommendation, error) {
	if mediaType != "" {
		return s.recommendationRepo.GetByMediaType(ctx, userID, mediaType, limit, offset)
	}
	return s.recommendationRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetRecommendationByID retrieves a specific recommendation by ID
func (s *recommendationService) GetRecommendationByID(ctx context.Context, id uint64) (*models.Recommendation, error) {
	return s.recommendationRepo.GetByID(ctx, id)
}

// MarkRecommendationAsViewed marks a recommendation as viewed
func (s *recommendationService) MarkRecommendationAsViewed(ctx context.Context, id uint64, userID uint64) error {
	// First check if this recommendation belongs to the user
	rec, err := s.recommendationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if rec.UserID != userID {
		return fmt.Errorf("recommendation does not belong to the user")
	}
	
	return s.recommendationRepo.MarkAsViewed(ctx, id)
}

// RateRecommendation sets a user rating for a recommendation
func (s *recommendationService) RateRecommendation(ctx context.Context, id uint64, userID uint64, rating float32) error {
	// Validate rating
	if rating < 0 || rating > 5 {
		return fmt.Errorf("rating must be between 0 and 5")
	}
	
	// First check if this recommendation belongs to the user
	rec, err := s.recommendationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if rec.UserID != userID {
		return fmt.Errorf("recommendation does not belong to the user")
	}
	
	return s.recommendationRepo.RateRecommendation(ctx, id, rating)
}

// StoreRecommendations stores recommendations for a user
func (s *recommendationService) StoreRecommendations(ctx context.Context, recommendations []*models.Recommendation) error {
	if len(recommendations) == 0 {
		return nil
	}
	
	// Set creation timestamp if not already set
	now := time.Now()
	for _, rec := range recommendations {
		if rec.CreatedAt.IsZero() {
			rec.CreatedAt = now
		}
	}
	
	return s.recommendationRepo.CreateMany(ctx, recommendations)
}

// GetRecentRecommendations retrieves recently created recommendations for a user
func (s *recommendationService) GetRecentRecommendations(ctx context.Context, userID uint64, days int, limit int) ([]models.Recommendation, error) {
	if days <= 0 {
		days = 7 // Default to 7 days
	}
	
	since := time.Now().AddDate(0, 0, -days)
	return s.recommendationRepo.GetRecentByUserID(ctx, userID, since, limit)
}

// GetTopRecommendations retrieves top-scored recommendations for a user
func (s *recommendationService) GetTopRecommendations(ctx context.Context, userID uint64, minScore float32, limit int) ([]models.Recommendation, error) {
	if minScore <= 0 {
		minScore = 0.7 // Default minimum score
	}
	
	return s.recommendationRepo.GetTopByUserID(ctx, userID, minScore, limit)
}

// GetRecommendationCount returns the count of recommendations for a user
func (s *recommendationService) GetRecommendationCount(ctx context.Context, userID uint64, mediaType string) (int, error) {
	return s.recommendationRepo.GetCount(ctx, userID, mediaType)
}

// DeleteExpiredRecommendations deletes all expired recommendations
func (s *recommendationService) DeleteExpiredRecommendations(ctx context.Context) error {
	return s.recommendationRepo.DeleteExpired(ctx)
}