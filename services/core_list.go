package services

import (
	"context"
	"fmt"

	mediatypes "suasor/clients/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// ListService manages application-specific List operations beyond the basic CRUD operations
type CoreListService[T mediatypes.ListData] interface {
	// Base operations (leveraging UserMediaItemService)
	GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error)
	GetByID(ctx context.Context, listID uint64) (*models.MediaItem[T], error)
	GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error)

	// list-specific operations
	GetItems(ctx context.Context, listID uint64) (*models.MediaItemList, error)
	GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error)

	Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[T], error)
}

type coreListService[T mediatypes.ListData] struct {
	itemRepo repository.CoreMediaItemRepository[T] // For fetching list items
}

// NewlistService creates a new list service
func NewCoreListService[T mediatypes.ListData](
	itemRepo repository.CoreMediaItemRepository[T],
) CoreListService[T] {
	return &coreListService[T]{
		itemRepo: itemRepo,
	}
}

// Base operations (delegating to UserMediaItemService where appropriate)
func (s coreListService[T]) GetByID(ctx context.Context, listID uint64) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", listID).
		Msg("Getting list by ID")

	// Use the user service
	result, err := s.itemRepo.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", listID).
			Msg("Failed to get list")
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	// Verify this is actually a list
	if !result.IsList() {
		log.Error().
			Uint64("id", listID).
			Str("actualType", string(result.Type)).
			Msg("Item is not a list")
		return nil, fmt.Errorf("item with ID %d is not a list", listID)
	}

	return result, nil
}
func (s coreListService[T]) GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting lists by user ID")
	// Use the user service
	return s.itemRepo.GetByUserID(ctx, userID, limit, offset)
}

// list-specific operations
func (s coreListService[T]) GetItems(ctx context.Context, listID uint64) (*models.MediaItemList, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Msg("Getting list items")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return &models.MediaItemList{}, fmt.Errorf("failed to get list items: %w", err)
	}

	itemList := list.GetData().GetItemList()

	// Return empty array if the list has no items
	if len(itemList.Items) == 0 {
		return &models.MediaItemList{}, nil
	}

	// Extract item IDs for batch retrieval
	itemIDs := make([]uint64, len(itemList.Items))
	for i, item := range itemList.Items {
		itemIDs[i] = item.ItemID
	}

	// Fetch the actual media items using the core media repository
	actualItems, err := s.itemRepo.GetMixedMediaItemsByIDs(ctx, itemIDs)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to fetch actual list items")
		return nil, fmt.Errorf("failed to get list items: %w", err)
	}

	return actualItems, nil

}
func (s coreListService[T]) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Uint64("userID", query.OwnerID).
		Msg("Searching lists")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	query.MediaType = mediaType

	// Delegate to the user service
	results, err := s.itemRepo.Search(ctx, query)
	if err != nil {
		log.Error().Err(err).
			Str("query", query.Query).
			Uint64("userID", query.OwnerID).
			Msg("Failed to search lists")
		return nil, fmt.Errorf("failed to search lists: %w", err)
	}

	log.Info().
		Str("query", query.Query).
		Uint64("userID", query.OwnerID).
		Int("count", len(results)).
		Msg("lists found")

	return results, nil
}
func (s coreListService[T]) GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)

	log.Debug().
		Int("limit", limit).
		Int("days", days).
		Msg("Getting recent lists")

	// Delegate to the user service
	results, err := s.itemRepo.GetRecentItems(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Msg("Failed to get recent lists")
		return nil, fmt.Errorf("failed to get recent lists: %w", err)
	}

	log.Info().
		Int("count", len(results)).
		Msg("Recent lists retrieved")

	return results, nil
}

func (s *coreListService[T]) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting all lists")

	publicOnly := true

	// Delegate to repository
	results, err := s.itemRepo.GetAll(ctx, limit, offset, publicOnly)
	if err != nil {
		log.Error().Err(err).
			Int("limit", limit).
			Int("offset", offset).
			Msg("Failed to get all lists")
		return nil, fmt.Errorf("failed to get all lists: %w", err)
	}

	log.Info().
		Int("count", len(results)).
		Msg("All lists retrieved successfully")

	return results, nil
}
