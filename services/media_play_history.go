package services

import (
	"context"
	"suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
)

// MediaPlayHistoryService defines operations for media play history
type MediaPlayHistoryService interface {
	// GetUserPlayHistory retrieves play history for a user with optional filtering
	GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) (interface{}, error)
	
	// GetContinueWatching retrieves items that a user has started but not completed
	GetContinueWatching(ctx context.Context, userID uint64, limit int) (interface{}, error)
	
	// GetByID retrieves a specific play history entry by ID
	GetByID(ctx context.Context, id uint64) (interface{}, error)
	
	// GetByMediaItemID retrieves play history for a specific media item
	GetByMediaItemID(ctx context.Context, mediaItemID, userID uint64) (interface{}, error)
	
	// RecordPlay records a new play event
	RecordPlay(ctx context.Context, history *models.MediaPlayHistoryGeneric) (interface{}, error)
	
	// ToggleFavorite marks or unmarks a media item as a favorite
	ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error
	
	// UpdateRating sets a user's rating for a media item
	UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error
	
	// GetFavorites retrieves favorite media items for a user
	GetFavorites(ctx context.Context, userID uint64, mediaType *types.MediaType, limit, offset int) (interface{}, error)
	
	// Delete removes a specific play history entry
	Delete(ctx context.Context, id uint64) error
	
	// ClearUserHistory removes all play history for a user
	ClearUserHistory(ctx context.Context, userID uint64, mediaType *types.MediaType) error
}

// mediaPlayHistoryService implements MediaPlayHistoryService
type mediaPlayHistoryService struct {
	repo repository.MediaPlayHistoryRepository
}

// NewMediaPlayHistoryService creates a new media play history service
func NewMediaPlayHistoryService(repo repository.MediaPlayHistoryRepository) MediaPlayHistoryService {
	return &mediaPlayHistoryService{
		repo: repo,
	}
}

// GetUserPlayHistory retrieves play history for a user with optional filtering
func (s *mediaPlayHistoryService) GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) (interface{}, error) {
	// Placeholder - this would be implemented with repository calls
	return s.repo.GetUserPlayHistory(ctx, userID, limit, offset, mediaType, completed)
}

// GetContinueWatching retrieves items that a user has started but not completed
func (s *mediaPlayHistoryService) GetContinueWatching(ctx context.Context, userID uint64, limit int) (interface{}, error) {
	// Items that are not completed and have been played recently
	return s.repo.GetContinueWatching(ctx, userID, limit)
}

// GetByID retrieves a specific play history entry by ID
func (s *mediaPlayHistoryService) GetByID(ctx context.Context, id uint64) (interface{}, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByMediaItemID retrieves play history for a specific media item
func (s *mediaPlayHistoryService) GetByMediaItemID(ctx context.Context, mediaItemID, userID uint64) (interface{}, error) {
	return s.repo.GetByMediaItemID(ctx, mediaItemID, userID)
}

// RecordPlay records a new play event
func (s *mediaPlayHistoryService) RecordPlay(ctx context.Context, history *models.MediaPlayHistoryGeneric) (interface{}, error) {
	return s.repo.RecordPlay(ctx, history)
}

// ToggleFavorite marks or unmarks a media item as a favorite
func (s *mediaPlayHistoryService) ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error {
	return s.repo.ToggleFavorite(ctx, mediaItemID, userID, favorite)
}

// UpdateRating sets a user's rating for a media item
func (s *mediaPlayHistoryService) UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error {
	return s.repo.UpdateRating(ctx, mediaItemID, userID, rating)
}

// GetFavorites retrieves favorite media items for a user
func (s *mediaPlayHistoryService) GetFavorites(ctx context.Context, userID uint64, mediaType *types.MediaType, limit, offset int) (interface{}, error) {
	return s.repo.GetFavorites(ctx, userID, mediaType, limit, offset)
}

// Delete removes a specific play history entry
func (s *mediaPlayHistoryService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

// ClearUserHistory removes all play history for a user
func (s *mediaPlayHistoryService) ClearUserHistory(ctx context.Context, userID uint64, mediaType *types.MediaType) error {
	return s.repo.ClearUserHistory(ctx, userID, mediaType)
}