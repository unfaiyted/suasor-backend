package services

import (
	"context"
	"fmt"
	"suasor/clients/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// UserMediaItemDataService defines the user service interface for user media item data
// This service focuses on user-specific operations like favorites, ratings, and history
type UserMediaItemDataService[T types.MediaData] interface {
	// Embed core service methods
	CoreUserMediaItemDataService[T]

	// GetUserHistory retrieves a user's media history
	GetUserHistory(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error)

	// GetRecentHistory retrieves a user's recent media history
	GetRecentHistory(ctx context.Context, userID uint64, days int, limit int) ([]*models.UserMediaItemData[T], error)

	// GetUserPlayHistory retrieves play history for a user with optional filtering
	GetUserPlayHistory(ctx context.Context, userID uint64, query *types.QueryOptions) ([]*models.UserMediaItemData[T], error)

	// GetContinueWatching retrieves items that a user has started but not completed
	GetContinueWatching(ctx context.Context, userID uint64, limit int) ([]*models.UserMediaItemData[T], error)

	// RecordPlay records a new play event
	RecordPlay(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)

	// ToggleFavorite marks or unmarks a media item as a favorite
	ToggleFavorite(ctx context.Context, itemID uint64, userID uint64, favorite bool) error

	// UpdateRating sets a user's rating for a media item
	UpdateRating(ctx context.Context, itemID, userID uint64, rating float32) error

	// GetFavorites retrieves favorite media items for a user
	GetFavorites(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error)

	// ClearUserHistory removes all data for a user
	ClearUserHistory(ctx context.Context, userID uint64) error
}

// userMediaItemDataService implements UserMediaItemDataService
type userMediaItemDataService[T types.MediaData] struct {
	CoreUserMediaItemDataService[T]
	repo repository.UserMediaItemDataRepository[T]
}

// NewUserMediaItemDataService creates a new user media item data service
func NewUserMediaItemDataService[T types.MediaData](
	coreService CoreUserMediaItemDataService[T],
	repo repository.UserMediaItemDataRepository[T],
) UserMediaItemDataService[T] {
	return &userMediaItemDataService[T]{
		CoreUserMediaItemDataService: coreService,
		repo:                         repo,
	}
}

// GetUserHistory retrieves a user's media history
func (s *userMediaItemDataService[T]) GetUserHistory(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user media history")

	// Delegate to repository
	result, err := s.repo.GetUserHistory(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get user media history")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(result)).
		Msg("User media history retrieved successfully")

	return result, nil
}

// GetRecentHistory retrieves a user's recent media history
func (s *userMediaItemDataService[T]) GetRecentHistory(ctx context.Context, userID uint64, days int, limit int) ([]*models.UserMediaItemData[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting user recent media history")

	// Delegate to repository
	result, err := s.repo.GetRecentHistory(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get user recent media history")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(result)).
		Msg("User recent media history retrieved successfully")

	return result, nil
}

// GetUserPlayHistory retrieves play history for a user with optional filtering
func (s *userMediaItemDataService[T]) GetUserPlayHistory(ctx context.Context, userID uint64, query *types.QueryOptions) ([]*models.UserMediaItemData[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", query.Limit).
		Int("offset", query.Offset).
		Msg("Getting user play history")

	completed := query.Watched || true
	// Delegate to repository
	result, err := s.repo.GetUserPlayHistory(ctx, userID, query.Limit, query.Offset, &completed)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get user play history")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(result)).
		Msg("User play history retrieved successfully")

	return result, nil
}

// GetContinueWatching retrieves items that a user has started but not completed
func (s *userMediaItemDataService[T]) GetContinueWatching(ctx context.Context, userID uint64, limit int) ([]*models.UserMediaItemData[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting continue watching items")

	// Delegate to repository
	result, err := s.repo.GetContinueWatching(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get continue watching items")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(result)).
		Msg("Continue watching items retrieved successfully")

	return result, nil
}

// RecordPlay records a new play event
func (s *userMediaItemDataService[T]) RecordPlay(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", data.UserID).
		Uint64("mediaItemID", data.MediaItemID).
		Msg("Recording play event")

	// Delegate to repository
	result, err := s.repo.RecordPlay(ctx, data)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", data.UserID).
			Uint64("mediaItemID", data.MediaItemID).
			Msg("Failed to record play event")
		return nil, err
	}

	log.Info().
		Uint64("id", result.ID).
		Uint64("userID", result.UserID).
		Uint64("mediaItemID", result.MediaItemID).
		Msg("Play event recorded successfully")

	return result, nil
}

// ToggleFavorite marks or unmarks a media item as a favorite
func (s *userMediaItemDataService[T]) ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("mediaItemID", mediaItemID).
		Bool("favorite", favorite).
		Msg("Toggling favorite status")

	// Delegate to repository
	err := s.repo.ToggleFavorite(ctx, mediaItemID, userID, favorite)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("mediaItemID", mediaItemID).
			Msg("Failed to toggle favorite status")
		return err
	}

	log.Info().
		Bool("favorite", favorite).
		Msg("Favorite status toggled successfully")

	return nil
}

// UpdateRating sets a user's rating for a media item
func (s *userMediaItemDataService[T]) UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("mediaItemID", mediaItemID).
		Float32("rating", rating).
		Msg("Updating rating")

		// Delegate to repository
	err := s.repo.UpdateRating(ctx, mediaItemID, userID, rating)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("mediaItemID", mediaItemID).
			Msg("Failed to update rating")
		return err
	}

	log.Info().
		Float32("rating", rating).
		Msg("Rating updated successfully")

	return nil
}

// GetFavorites retrieves favorite media items for a user
func (s *userMediaItemDataService[T]) GetFavorites(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user favorite items")

	// Delegate to repository
	result, err := s.repo.GetFavorites(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get user favorite items")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(result)).
		Msg("User favorite items retrieved successfully")

	return result, nil
}

// ClearUserHistory removes all data for a user
func (s *userMediaItemDataService[T]) ClearUserHistory(ctx context.Context, userID uint64) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Clearing user history")

	// Delegate to repository
	err := s.repo.ClearUserHistory(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to clear user history")
		return fmt.Errorf("failed to clear user history: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Msg("User history cleared successfully")

	return nil
}
