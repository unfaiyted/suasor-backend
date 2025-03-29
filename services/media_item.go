package services

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
)

// MediaItemService defines the interface for media item operations
type MediaItemService[T types.MediaData] interface {
	Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
	GetByExternalID(ctx context.Context, externalID string, clientID uint64) (*models.MediaItem[T], error)
	GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error)
	GetByType(ctx context.Context, mediaType types.MediaType, clientID uint64) ([]*models.MediaItem[T], error)
	GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error)
	Delete(ctx context.Context, id uint64) error

	// Additional service methods
	SearchByTitle(ctx context.Context, title string, userID uint64) ([]*models.MediaItem[T], error)
	GetRecentItems(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error)
}

type mediaItemService[T types.MediaData] struct {
	repo repository.MediaItemRepository[T]
}

// NewMediaItemService creates a new media item service
func NewMediaItemService[T types.MediaData](repo repository.MediaItemRepository[T]) MediaItemService[T] {
	return &mediaItemService[T]{repo: repo}
}

func (s *mediaItemService[T]) Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	// Validate the media item
	if err := validateMediaItem(item); err != nil {
		return nil, fmt.Errorf("invalid media item: %w", err)
	}

	return s.repo.Create(ctx, item)
}

func (s *mediaItemService[T]) Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	// Validate the media item
	if err := validateMediaItem(item); err != nil {
		return nil, fmt.Errorf("invalid media item: %w", err)
	}

	return s.repo.Update(ctx, item)
}

func (s *mediaItemService[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	return s.repo.GetByID(ctx, id)
}

func (s *mediaItemService[T]) GetByExternalID(ctx context.Context, externalID string, clientID uint64) (*models.MediaItem[T], error) {
	return s.repo.GetByExternalID(ctx, externalID, clientID)
}

func (s *mediaItemService[T]) GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error) {
	return s.repo.GetByClientID(ctx, clientID)
}

func (s *mediaItemService[T]) GetByType(ctx context.Context, mediaType types.MediaType, clientID uint64) ([]*models.MediaItem[T], error) {
	return s.repo.GetByType(ctx, mediaType, clientID)
}

func (s *mediaItemService[T]) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *mediaItemService[T]) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

func (s *mediaItemService[T]) SearchByTitle(ctx context.Context, title string, userID uint64) ([]*models.MediaItem[T], error) {
	// Implementation would require extending the repository with a search method
	// This is just a placeholder showing how we might implement this functionality
	items, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter items by title (this would be more efficient in the database query)
	var filtered []*models.MediaItem[T]
	for _, item := range items {
		details := item.Data.GetDetails()
		if containsIgnoreCase(details.Title, title) {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

func (s *mediaItemService[T]) GetRecentItems(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error) {
	// This would ideally be implemented in the repository
	// For now, we'll get all items and sort/limit them here
	items, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Sort items by created time (newest first)
	// Note: In a real implementation, this would be handled at the database level
	sortByNewest(items)

	// Limit the number of results
	if len(items) > limit {
		items = items[:limit]
	}

	return items, nil
}

// Helper functions

func validateMediaItem[T types.MediaData](item models.MediaItem[T]) error {
	// Basic validation
	if item.ClientID == 0 {
		return fmt.Errorf("client ID is required")
	}

	if item.Type == "" {
		return fmt.Errorf("media type is required")
	}

	// Validate type-specific data
	details := item.Data.GetDetails()
	if details.Title == "" {
		return fmt.Errorf("title is required")
	}

	return nil
}

func containsIgnoreCase(s, substr string) bool {
	// Simple case-insensitive contains implementation
	// In a real app, you might use a more sophisticated search
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}

func sortByNewest[T types.MediaData](items []*models.MediaItem[T]) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
}
