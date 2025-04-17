package services

import (
	"context"
	"fmt"

	"suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// CoreMediaItemService defines the core interface for media item operations
// This service focuses on basic operations that are common to all media items
// regardless of whether they are client-associated or user-owned
type CoreMediaItemService[T types.MediaData] interface {
	// Basic CRUD operations
	Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
	Delete(ctx context.Context, id uint64) error
	GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error)
	GetByClientItemID(ctx context.Context, clientItemID string, clientID uint64) (*models.MediaItem[T], error)

	// Basic query operations
	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error)
	GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error)

	// Search operations
	Search(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error)
	GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error)
}

// coreMediaItemService implements CoreMediaItemService
type coreMediaItemService[T types.MediaData] struct {
	repo repository.MediaItemRepository[T]
}

// NewCoreMediaItemService creates a new core media item service
func NewCoreMediaItemService[T types.MediaData](repo repository.MediaItemRepository[T]) CoreMediaItemService[T] {
	return &coreMediaItemService[T]{repo: repo}
}

// Create adds a new media item
func (s *coreMediaItemService[T]) Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(item.Type)).
		Msg("Creating media item")

	// Validate the media item
	if err := s.validateMediaItem(&item); err != nil {
		return nil, fmt.Errorf("invalid media item: %w", err)
	}

	// Delegate to repository
	result, err := s.repo.Create(ctx, item)
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
func (s *coreMediaItemService[T]) Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", item.ID).
		Str("type", string(item.Type)).
		Msg("Updating media item")

	// Validate the media item
	if err := s.validateMediaItem(&item); err != nil {
		return nil, fmt.Errorf("invalid media item: %w", err)
	}

	// Delegate to repository
	result, err := s.repo.Update(ctx, item)
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

// GetByID retrieves a media item by its ID
func (s *coreMediaItemService[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Getting media item by ID")

	// Delegate to repository
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to get media item")
		return nil, err
	}

	log.Debug().
		Uint64("id", id).
		Str("type", string(result.Type)).
		Msg("Media item retrieved successfully")

	return result, nil
}

// Delete removes a media item
func (s *coreMediaItemService[T]) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting media item")

	// Delegate to repository
	err := s.repo.Delete(ctx, id)
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

// GetByExternalID retrieves a media item by its external ID
func (s *coreMediaItemService[T]) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("source", source).
		Str("externalID", externalID).
		Msg("Getting media item by external ID")

	// Delegate to repository
	result, err := s.repo.GetByExternalID(ctx, source, externalID)
	if err != nil {
		log.Error().Err(err).
			Str("source", source).
			Str("externalID", externalID).
			Msg("Failed to get media item by external ID")
		return nil, err
	}

	log.Debug().
		Str("source", source).
		Str("externalID", externalID).
		Uint64("id", result.ID).
		Msg("Media item retrieved by external ID")

	return result, nil
}

// GetByType retrieves all media items of a specific type
func (s *coreMediaItemService[T]) GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(mediaType)).
		Msg("Getting media items by type")

	// Delegate to repository
	results, err := s.repo.GetByType(ctx, mediaType)
	if err != nil {
		log.Error().Err(err).
			Str("type", string(mediaType)).
			Msg("Failed to get media items by type")
		return nil, err
	}

	log.Info().
		Str("type", string(mediaType)).
		Int("count", len(results)).
		Msg("Media items retrieved by type")

	return results, nil
}

// Search finds media items based on a query string
func (s *coreMediaItemService[T]) Search(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Str("type", string(query.MediaType)).
		Int("limit", query.Limit).
		Int("offset", query.Offset).
		Msg("Searching media items")

	// Delegate to repository
	results, err := s.repo.Search(ctx, query)
	if err != nil {
		log.Error().Err(err).
			Str("query", query.Query).
			Str("type", string(query.MediaType)).
			Msg("Failed to search media items")
		return nil, err
	}

	log.Info().
		Str("query", query.Query).
		Str("type", string(query.MediaType)).
		Int("count", len(results)).
		Msg("Media items found")

	return results, nil
}

// GetRecentItems retrieves recently added items of a specific type
func (s *coreMediaItemService[T]) GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Str("type", string(mediaType)).
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recent media items")

	// Delegate to repository
	results, err := s.repo.GetRecentItems(ctx, mediaType, days, limit)
	if err != nil {
		log.Error().Err(err).
			Str("type", string(mediaType)).
			Msg("Failed to get recent media items")
		return nil, err
	}

	log.Info().
		Str("type", string(mediaType)).
		Int("count", len(results)).
		Msg("Recent media items retrieved")

	return results, nil
}

func (s *coreMediaItemService[T]) validateMediaItem(item *models.MediaItem[T]) error {
	// Validate the media item
	// if err := item.Validate(); err != nil {
	// 	return fmt.Errorf("invalid media item: %w", err)
	// }
	// TODO: Add validation for media item

	// Delegate to repository
	// return s.repo.ValidateMediaItem(item)
	return nil
}

func (s *coreMediaItemService[T]) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting all media items")

	// Delegate to repository
	results, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Int("limit", limit).
			Int("offset", offset).
			Msg("Failed to get all media items")
		return nil, fmt.Errorf("failed to get all media items: %w", err)
	}

	log.Info().
		Int("count", len(results)).
		Msg("All media items retrieved successfully")

	return results, nil
}

func (s *coreMediaItemService[T]) GetByClientItemID(ctx context.Context, clientItemID string, clientID uint64) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("clientItemID", clientItemID).
		Uint64("clientID", clientID).
		Msg("Getting media item by client item ID")

	// Delegate to repository
	result, err := s.repo.GetByClientItemID(ctx, clientItemID, clientID)
	if err != nil {
		log.Error().Err(err).
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Failed to get media item by client item ID")
		return nil, err
	}

	log.Debug().
		Str("clientItemID", clientItemID).
		Uint64("clientID", clientID).
		Str("type", string(result.Type)).
		Msg("Media item retrieved by client item ID")

	return result, nil
}
