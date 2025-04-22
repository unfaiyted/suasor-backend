package services

import (
	"context"
	"fmt"

	"suasor/clients/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// UserMediaItemService defines the interface for user-owned media item operations
// This service extends CoreMediaItemService with operations specific to media items
// that are directly owned by users, such as playlists and collections
type UserMediaItemService[T types.MediaData] interface {
	// Include all core service methods
	CoreMediaItemService[T]

	// User-specific operations
	GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error)
	GetUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error)

	// Specific to user-owned collections/playlists
	SearchUserContent(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error)
	GetRecentUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error)
}

// userMediaItemService implements UserMediaItemService
type userMediaItemService[T types.MediaData] struct {
	coreService CoreMediaItemService[T] // Embed the core service
	userRepo    repository.UserMediaItemRepository[T]
}

// NewUserMediaItemService creates a new user-owned media item service
func NewUserMediaItemService[T types.MediaData](
	coreService CoreMediaItemService[T],
	userRepo repository.UserMediaItemRepository[T],
) UserMediaItemService[T] {
	return &userMediaItemService[T]{
		coreService: coreService,
		userRepo:    userRepo,
	}
}

// Core service methods - delegate to embedded core service

func (s *userMediaItemService[T]) Create(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error) {
	return s.coreService.Create(ctx, item)
}

func (s *userMediaItemService[T]) Update(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error) {
	return s.coreService.Update(ctx, item)
}

func (s *userMediaItemService[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	return s.coreService.GetByID(ctx, id)
}

func (s *userMediaItemService[T]) GetMostPlayed(ctx context.Context, limit int) ([]*models.MediaItem[T], error) {
	return s.coreService.GetMostPlayed(ctx, limit)
}

func (s *userMediaItemService[T]) GetByClientItemID(ctx context.Context, clientItemID string, clientID uint64) (*models.MediaItem[T], error) {
	return s.coreService.GetByClientItemID(ctx, clientItemID, clientID)
}

func (s *userMediaItemService[T]) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error) {
	return s.coreService.GetAll(ctx, limit, offset)
}

func (s *userMediaItemService[T]) Delete(ctx context.Context, id uint64) error {
	return s.coreService.Delete(ctx, id)
}

func (s *userMediaItemService[T]) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error) {
	return s.coreService.GetByExternalID(ctx, source, externalID)
}

func (s *userMediaItemService[T]) GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error) {
	return s.coreService.GetByType(ctx, mediaType)
}

func (s *userMediaItemService[T]) Search(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error) {
	return s.coreService.Search(ctx, query)
}

func (s *userMediaItemService[T]) GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error) {
	return s.coreService.GetRecentItems(ctx, days, limit)
}

// User-specific methods

// GetByUserID retrieves all user-owned media items for a specific user
func (s *userMediaItemService[T]) GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting media items by user ID")

	// Delegate to user repository
	results, err := s.userRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get media items by user ID")
		return nil, fmt.Errorf("failed to get media items by user ID: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(results)).
		Msg("Media items retrieved by user ID")

	return results, nil
}

// GetUserContent retrieves all types of user-owned content in a single query
func (s *userMediaItemService[T]) GetUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting all user-owned content")

	// Delegate to user repository
	results, err := s.userRepo.GetUserContent(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get user-owned content")
		return nil, fmt.Errorf("failed to get user content: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(results)).
		Msg("User-owned content retrieved successfully")

	return results, nil
}

// SearchUserContent searches for user-owned content based on query parameters
func (s *userMediaItemService[T]) SearchUserContent(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Uint64("userID", query.OwnerID).
		Str("type", string(query.MediaType)).
		Msg("Searching user-owned content")

	// Delegate to user repository
	results, err := s.userRepo.Search(ctx, query)
	if err != nil {
		log.Error().Err(err).
			Str("query", query.Query).
			Uint64("userID", query.OwnerID).
			Str("type", string(query.MediaType)).
			Msg("Failed to search user-owned content")
		return nil, fmt.Errorf("failed to search user content: %w", err)
	}

	log.Info().
		Str("query", query.Query).
		Uint64("userID", query.OwnerID).
		Str("type", string(query.MediaType)).
		Int("count", len(results)).
		Msg("User-owned content found")

	return results, nil
}

// GetRecentUserContent retrieves recently created or updated user-owned content
func (s *userMediaItemService[T]) GetRecentUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)
	log.Debug().
		Uint64("userID", userID).
		Str("type", string(mediaType)).
		Int("limit", limit).
		Msg("Getting recent user-owned content")

	// Create query options for the repository
	options := types.QueryOptions{
		MediaType: mediaType,
		OwnerID:   userID,
		Limit:     limit,
		Sort:      "updated_at",
		SortOrder: "desc",
	}

	// Delegate to user repository
	results, err := s.userRepo.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Str("type", string(mediaType)).
			Msg("Failed to get recent user-owned content")
		return nil, fmt.Errorf("failed to get recent user content: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Str("type", string(mediaType)).
		Int("count", len(results)).
		Msg("Recent user-owned content retrieved")

	return results, nil
}
