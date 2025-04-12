package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"suasor/types/models"
)

// SearchRepository defines operations for search-related data
type SearchRepository interface {
	// SaveSearchHistory saves a user's search query to the search history
	SaveSearchHistory(ctx context.Context, userID uint64, query string, resultCount int) (*models.SearchHistory, error)
	
	// GetRecentSearches retrieves recent searches for a user
	GetRecentSearches(ctx context.Context, userID uint64, limit int) ([]models.SearchHistory, error)
	
	// GetTrendingSearches retrieves popular searches across all users
	GetTrendingSearches(ctx context.Context, limit int) ([]models.SearchHistory, error)
	
	// GetSearchSuggestions retrieves search suggestions based on a partial query
	GetSearchSuggestions(ctx context.Context, partialQuery string, limit int) ([]string, error)
}

// searchRepository implements SearchRepository using GORM
type searchRepository struct {
	db *gorm.DB
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

// SaveSearchHistory saves a user's search query to the search history
func (r *searchRepository) SaveSearchHistory(ctx context.Context, userID uint64, query string, resultCount int) (*models.SearchHistory, error) {
	searchHistory := &models.SearchHistory{
		UserID:      userID,
		Query:       query,
		ResultCount: resultCount,
		SearchedAt:  time.Now(),
	}
	
	result := r.db.WithContext(ctx).Create(searchHistory)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to save search history: %w", result.Error)
	}
	
	return searchHistory, nil
}

// GetRecentSearches retrieves recent searches for a user
func (r *searchRepository) GetRecentSearches(ctx context.Context, userID uint64, limit int) ([]models.SearchHistory, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	
	var searches []models.SearchHistory
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("searched_at DESC").
		Limit(limit).
		Find(&searches)
	
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recent searches: %w", result.Error)
	}
	
	return searches, nil
}

// GetTrendingSearches retrieves popular searches across all users
func (r *searchRepository) GetTrendingSearches(ctx context.Context, limit int) ([]models.SearchHistory, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	
	var searches []models.SearchHistory
	
	// Get trending searches from the last 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	
	result := r.db.WithContext(ctx).
		Select("query, COUNT(*) as search_count").
		Where("searched_at > ?", sevenDaysAgo).
		Group("query").
		Order("search_count DESC").
		Limit(limit).
		Find(&searches)
	
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get trending searches: %w", result.Error)
	}
	
	return searches, nil
}

// GetSearchSuggestions retrieves search suggestions based on a partial query
func (r *searchRepository) GetSearchSuggestions(ctx context.Context, partialQuery string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 5 // Default limit
	}
	
	var suggestions []string
	
	// Get distinct queries that start with the partial query
	result := r.db.WithContext(ctx).
		Model(&models.SearchHistory{}).
		Select("DISTINCT query").
		Where("query LIKE ?", partialQuery+"%").
		Group("query").
		Order("COUNT(*) DESC").
		Limit(limit).
		Pluck("query", &suggestions)
	
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get search suggestions: %w", result.Error)
	}
	
	return suggestions, nil
}
