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

	// User-specific operations (not implemented in CoreMediaItemService)
	Update(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error)
	Delete(ctx context.Context, id uint64) error
	Create(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error)

	// User-specific operations
	GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error)
	GetUserContent(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error)

	// Specific to user-owned collections/playlists
	SearchUserContent(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error)
	GetRecentUserContent(ctx context.Context, userID uint64, days int, limit int) ([]*models.MediaItem[T], error)
}

// userMediaItemService implements UserMediaItemService
type userMediaItemService[T types.MediaData] struct {
	CoreMediaItemService[T] // Embed the core service
	userRepo                repository.UserMediaItemRepository[T]
}

// NewUserMediaItemService creates a new user-owned media item service
func NewUserMediaItemService[T types.MediaData](
	coreService CoreMediaItemService[T],
	userRepo repository.UserMediaItemRepository[T],
) UserMediaItemService[T] {
	return &userMediaItemService[T]{
		CoreMediaItemService: coreService,
		userRepo:             userRepo,
	}
}

// Create adds a new media item
func (s *userMediaItemService[T]) Create(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(item.Type)).
		Msg("Creating media item")

	// Validate the media item
	if err := s.validateMediaItem(item); err != nil {
		return nil, fmt.Errorf("invalid media item: %w", err)
	}

	// Delegate to repository
	result, err := s.userRepo.Create(ctx, item)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create media item")
		return nil, err
	}

	log.Info().
		Uint64("id", result.ID).
		Str("type", string(result.Type)).
		Msg("Media item created successfully")

	return result, nil
}

// Update modifies an existing media item
func (s *userMediaItemService[T]) Update(ctx context.Context, item *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", item.ID).
		Str("type", string(item.Type)).
		Msg("Updating media item")

	// Validate the media item
	if err := s.validateMediaItem(item); err != nil {
		return nil, fmt.Errorf("invalid media item: %w", err)
	}

	// Delegate to repository
	result, err := s.userRepo.Update(ctx, item)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", item.ID).
			Msg("Failed to update media item")
		return nil, err
	}

	log.Info().
		Uint64("id", result.ID).
		Str("type", string(result.Type)).
		Msg("Media item updated successfully")

	return result, nil
}

// Delete removes a media item
func (s *userMediaItemService[T]) Delete(ctx context.Context, id uint64) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting media item")

	// Delegate to repository
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to delete media item")
		return err
	}

	log.Info().
		Uint64("id", id).
		Msg("Media item deleted successfully")

	return nil
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
func (s *userMediaItemService[T]) GetRecentUserContent(ctx context.Context, userID uint64, days int, limit int) ([]*models.MediaItem[T], error) {
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

func (s *userMediaItemService[T]) validateMediaItem(item *models.MediaItem[T]) error {
	// Validate the media item
	// if err := item.Validate(); err != nil {
	// 	return fmt.Errorf("invalid media item: %w", err)
	// }
	// TODO: Add validation for media item

	// Delegate to repository
	// return s.repo.ValidateMediaItem(item)
	return nil
}
