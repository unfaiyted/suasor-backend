package repository

import (
	"context"
	"fmt"
	"time"

	"suasor/types/models"

	"gorm.io/gorm"
)

// RecommendationRepository defines the interface for recommendation operations
type RecommendationRepository interface {
	// Create creates a new recommendation
	Create(ctx context.Context, recommendation *models.Recommendation) error
	// CreateMany creates multiple recommendations in a batch
	CreateMany(ctx context.Context, recommendations []*models.Recommendation) error
	// GetByID retrieves a recommendation by ID
	GetByID(ctx context.Context, id uint64) (*models.Recommendation, error)
	// GetByUserID retrieves recommendations for a specific user
	GetByUserID(ctx context.Context, userID uint64, limit, offset int) ([]models.Recommendation, error)
	// GetByMediaType retrieves recommendations for a specific user and media type
	GetByMediaType(ctx context.Context, userID uint64, mediaType string, limit, offset int) ([]models.Recommendation, error)
	// GetRecentByUserID retrieves recent recommendations for a user
	GetRecentByUserID(ctx context.Context, userID uint64, since time.Time, limit int) ([]models.Recommendation, error)
	// GetTopByUserID retrieves top-scored recommendations for a user
	GetTopByUserID(ctx context.Context, userID uint64, minScore float32, limit int) ([]models.Recommendation, error)
	// MarkAsViewed marks a recommendation as viewed
	MarkAsViewed(ctx context.Context, id uint64) error
	// RateRecommendation sets a user rating for a recommendation
	RateRecommendation(ctx context.Context, id uint64, rating float32) error
	// DeleteByJobRunID deletes recommendations by job run ID
	DeleteByJobRunID(ctx context.Context, jobRunID uint64) error
	// DeleteByUserID deletes all recommendations for a user
	DeleteByUserID(ctx context.Context, userID uint64) error
	// DeleteExpired deletes all expired recommendations
	DeleteExpired(ctx context.Context) error
	// GetCount returns the count of recommendations for a user
	GetCount(ctx context.Context, userID uint64, mediaType string) (int, error)
}

// recommendationRepository implements the RecommendationRepository interface
type recommendationRepository struct {
	db *gorm.DB
}

// NewRecommendationRepository creates a new recommendation repository
func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &recommendationRepository{db: db}
}

// Create creates a new recommendation
func (r *recommendationRepository) Create(ctx context.Context, recommendation *models.Recommendation) error {
	result := r.db.WithContext(ctx).Create(recommendation)
	if result.Error != nil {
		return fmt.Errorf("failed to create recommendation: %w", result.Error)
	}
	return nil
}

// CreateMany creates multiple recommendations in a batch
func (r *recommendationRepository) CreateMany(ctx context.Context, recommendations []*models.Recommendation) error {
	if len(recommendations) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Create(recommendations)
	if result.Error != nil {
		return fmt.Errorf("failed to create recommendations in batch: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a recommendation by ID
func (r *recommendationRepository) GetByID(ctx context.Context, id uint64) (*models.Recommendation, error) {
	var recommendation models.Recommendation
	result := r.db.WithContext(ctx).First(&recommendation, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get recommendation by ID: %w", result.Error)
	}
	return &recommendation, nil
}

// GetByUserID retrieves recommendations for a specific user
func (r *recommendationRepository) GetByUserID(ctx context.Context, userID uint64, limit, offset int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation
	
	if limit <= 0 {
		limit = 20 // Default limit
	}

	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&recommendations)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recommendations by user ID: %w", result.Error)
	}
	
	return recommendations, nil
}

// GetByMediaType retrieves recommendations for a specific user and media type
func (r *recommendationRepository) GetByMediaType(ctx context.Context, userID uint64, mediaType string, limit, offset int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation
	
	if limit <= 0 {
		limit = 20 // Default limit
	}

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND media_type = ?", userID, mediaType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&recommendations)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recommendations by media type: %w", result.Error)
	}
	
	return recommendations, nil
}

// GetRecentByUserID retrieves recent recommendations for a user
func (r *recommendationRepository) GetRecentByUserID(ctx context.Context, userID uint64, since time.Time, limit int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation
	
	if limit <= 0 {
		limit = 20 // Default limit
	}

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Order("created_at DESC").
		Limit(limit).
		Find(&recommendations)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recent recommendations: %w", result.Error)
	}
	
	return recommendations, nil
}

// GetTopByUserID retrieves top-scored recommendations for a user
func (r *recommendationRepository) GetTopByUserID(ctx context.Context, userID uint64, minScore float32, limit int) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation
	
	if limit <= 0 {
		limit = 20 // Default limit
	}

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND score >= ?", userID, minScore).
		Order("score DESC").
		Limit(limit).
		Find(&recommendations)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get top recommendations: %w", result.Error)
	}
	
	return recommendations, nil
}

// MarkAsViewed marks a recommendation as viewed
func (r *recommendationRepository) MarkAsViewed(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).
		Model(&models.Recommendation{}).
		Where("id = ?", id).
		Update("is_viewed", true)

	if result.Error != nil {
		return fmt.Errorf("failed to mark recommendation as viewed: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// RateRecommendation sets a user rating for a recommendation
func (r *recommendationRepository) RateRecommendation(ctx context.Context, id uint64, rating float32) error {
	result := r.db.WithContext(ctx).
		Model(&models.Recommendation{}).
		Where("id = ?", id).
		Update("user_rating", rating)

	if result.Error != nil {
		return fmt.Errorf("failed to rate recommendation: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// DeleteByJobRunID deletes recommendations by job run ID
func (r *recommendationRepository) DeleteByJobRunID(ctx context.Context, jobRunID uint64) error {
	result := r.db.WithContext(ctx).
		Where("job_run_id = ?", jobRunID).
		Delete(&models.Recommendation{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete recommendations by job run ID: %w", result.Error)
	}
	
	return nil
}

// DeleteByUserID deletes all recommendations for a user
func (r *recommendationRepository) DeleteByUserID(ctx context.Context, userID uint64) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.Recommendation{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete recommendations by user ID: %w", result.Error)
	}
	
	return nil
}

// DeleteExpired deletes all expired recommendations
func (r *recommendationRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	
	result := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Delete(&models.Recommendation{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete expired recommendations: %w", result.Error)
	}
	
	return nil
}

// GetCount returns the count of recommendations for a user
func (r *recommendationRepository) GetCount(ctx context.Context, userID uint64, mediaType string) (int, error) {
	var count int64
	
	query := r.db.WithContext(ctx).Model(&models.Recommendation{}).Where("user_id = ?", userID)
	
	if mediaType != "" {
		query = query.Where("media_type = ?", mediaType)
	}
	
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get recommendation count: %w", err)
	}
	
	return int(count), nil
}