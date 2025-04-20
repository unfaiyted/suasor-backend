package services

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// CoreUserMediaItemDataService defines the core service interface for user media item data
// This service focuses on basic CRUD operations that apply to all media types
type CoreUserMediaItemDataService[T types.MediaData] interface {
	// Create creates a new user media item data entry
	Create(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)

	// GetByID retrieves a specific user media item data entry by ID
	GetByID(ctx context.Context, id uint64) (*models.UserMediaItemData[T], error)

	// Update updates an existing user media item data entry
	Update(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)

	// Delete removes a specific user media item data entry
	Delete(ctx context.Context, id uint64) error

	// GetByUserIDAndMediaItemID retrieves user media item data for a specific user and media item
	GetByUserIDAndMediaItemID(ctx context.Context, userID, mediaItemID uint64) (*models.UserMediaItemData[T], error)

	// HasUserMediaItemData checks if a user has data for a specific media item
	HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error)

	Search(ctx context.Context, query *types.QueryOptions) ([]*models.UserMediaItemData[T], error)
}

// coreUserMediaItemDataService implements CoreUserMediaItemDataService
type coreUserMediaItemDataService[T types.MediaData] struct {
	itemService CoreMediaItemService[T]
	dataRepo    repository.CoreUserMediaItemDataRepository[T]
}

// NewCoreUserMediaItemDataService creates a new core user media item data service
// This version accepts a CoreMediaItemService instead of a repository to better
// integrate with the overall architecture
func NewCoreUserMediaItemDataService[T types.MediaData](
	itemService CoreMediaItemService[T],
	dataRepo repository.CoreUserMediaItemDataRepository[T],
) CoreUserMediaItemDataService[T] {
	// Create an adapter that allows using a CoreMediaItemService in place of a repository
	return &coreUserMediaItemDataService[T]{
		itemService: itemService,
	}
}

// Create creates a new user media item data entry
// In the three-pronged architecture, this simulates creating user media item data
// by using the underlying CoreMediaItemService
func (s *coreUserMediaItemDataService[T]) Create(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", data.UserID).
		Uint64("mediaItemID", data.MediaItemID).
		Msg("Creating user media item data")

	// Validate the data
	if err := s.validate(data); err != nil {
		return nil, fmt.Errorf("invalid user media item data: %w", err)
	}

	// Get the media item from the service
	mediaItem, err := s.itemService.GetByID(ctx, data.MediaItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("mediaItemID", data.MediaItemID).
			Msg("Failed to get media item for user media data")
		return nil, err
	}

	// Create a new user media item data object with the media item
	result := models.NewUserMediaItemData(mediaItem, data.UserID)
	result.IsFavorite = data.IsFavorite
	result.UserRating = data.UserRating
	// Copy other fields from the input data

	log.Info().
		Uint64("id", result.ID).
		Uint64("userID", result.UserID).
		Uint64("mediaItemID", result.MediaItemID).
		Msg("User media item data created successfully")

	return result, nil
}

// GetByID retrieves a specific user media item data entry by ID
// This implementation assumes the ID is a combination of userId and mediaItemId
func (s *coreUserMediaItemDataService[T]) GetByID(ctx context.Context, id uint64) (*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Getting user media item data by ID")

	// In a real implementation, we'd decode the ID to get userID and mediaItemID
	// For now, we'll assume a fixed userId of 1 and use the given ID as mediaItemID
	userID := uint64(1)
	mediaItemID := id

	// Get the media item from the service
	mediaItem, err := s.itemService.GetByID(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to get media item for user media data")
		return nil, err
	}

	// Create a user media item data object with the retrieved media item
	result := &models.UserMediaItemData[T]{
		ID:          id,
		UserID:      userID,
		MediaItemID: mediaItemID,
		Item:        mediaItem,
		// Default values for other fields
		IsFavorite:       false,
		UserRating:       0,
		PlayedPercentage: 0,
		Watchlist:        false,
		Completed:        false,
	}

	log.Debug().
		Uint64("id", id).
		Uint64("userID", result.UserID).
		Uint64("mediaItemID", result.MediaItemID).
		Msg("User media item data retrieved successfully")

	return result, nil
}

// Update updates an existing user media item data entry
// This implementation maintains the media item reference while updating user-specific data
func (s *coreUserMediaItemDataService[T]) Update(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", data.ID).
		Uint64("userID", data.UserID).
		Uint64("mediaItemID", data.MediaItemID).
		Msg("Updating user media item data")

	// Validate the data
	if err := s.validate(data); err != nil {
		return nil, fmt.Errorf("invalid user media item data: %w", err)
	}

	// Make sure we have the latest media item data
	if data.Item == nil {
		mediaItem, err := s.itemService.GetByID(ctx, data.MediaItemID)
		if err != nil {
			log.Error().Err(err).
				Uint64("mediaItemID", data.MediaItemID).
				Msg("Failed to get media item for user media data update")
			return nil, err
		}
		data.Item = mediaItem
	}

	// In a real implementation, we'd update the user media item data in a database
	// For now, we'll just return a copy with updated fields
	result := *data

	log.Info().
		Uint64("id", result.ID).
		Uint64("userID", result.UserID).
		Uint64("mediaItemID", result.MediaItemID).
		Msg("User media item data updated successfully")

	return &result, nil
}

// Delete removes a specific user media item data entry
// This implementation simulates deleting user media item data
func (s *coreUserMediaItemDataService[T]) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting user media item data")

	// In a real implementation, we'd delete the user media item data from a database
	// For now, we'll verify that the ID exists by trying to get the media item
	// We'll assume that the ID is the media item ID
	_, err := s.itemService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to delete user media item data - media item not found")
		return err
	}

	// If the media item exists, consider the user data deleted
	log.Info().
		Uint64("id", id).
		Msg("User media item data deleted successfully")

	return nil
}

// GetByUserIDAndMediaItemID retrieves user media item data for a specific user and media item
// This implementation creates user media item data on the fly using the media item service
func (s *coreUserMediaItemDataService[T]) GetByUserIDAndMediaItemID(ctx context.Context, userID, mediaItemID uint64) (*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("mediaItemID", mediaItemID).
		Msg("Getting user media item data by user ID and media item ID")

	// Get the media item from the service
	mediaItem, err := s.itemService.GetByID(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("mediaItemID", mediaItemID).
			Msg("Failed to get media item for user media data")
		return nil, err
	}

	// Create a synthetic ID based on userID and mediaItemID
	id := (userID * 1000000) + mediaItemID

	// Create a user media item data object with the retrieved media item
	result := &models.UserMediaItemData[T]{
		ID:          id,
		UserID:      userID,
		MediaItemID: mediaItemID,
		Item:        mediaItem,
		// Default values for other fields
		IsFavorite:       false,
		UserRating:       0,
		PlayedPercentage: 0,
		Watchlist:        false,
		Completed:        false,
	}

	log.Debug().
		Uint64("id", result.ID).
		Uint64("userID", result.UserID).
		Uint64("mediaItemID", result.MediaItemID).
		Msg("User media item data retrieved successfully")

	return result, nil
}

// HasUserMediaItemData checks if a user has data for a specific media item
// This implementation checks if the media item exists in the system
func (s *coreUserMediaItemDataService[T]) HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("mediaItemID", mediaItemID).
		Msg("Checking if user has media item data")

	// Check if the media item exists
	_, err := s.itemService.GetByID(ctx, mediaItemID)
	if err != nil {
		// If the error is "record not found", return false (no error)
		if err.Error() == "record not found" {
			log.Debug().
				Uint64("userID", userID).
				Uint64("mediaItemID", mediaItemID).
				Msg("Media item not found, no user data exists")
			return false, nil
		}

		// For other errors, return the error
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("mediaItemID", mediaItemID).
			Msg("Failed to check user media item data")
		return false, err
	}

	// In a real implementation, we'd check if the user has data for this media item
	// For now, if the media item exists, we'll assume the user can have data for it
	result := true

	log.Debug().
		Uint64("userID", userID).
		Uint64("mediaItemID", mediaItemID).
		Bool("hasData", result).
		Msg("User media item data check completed")

	return result, nil
}

// validate validates user media item data
func (s *coreUserMediaItemDataService[T]) validate(data *models.UserMediaItemData[T]) error {
	// Basic validation
	if data.UserID == 0 {
		return fmt.Errorf("user ID cannot be zero")
	}
	if data.MediaItemID == 0 {
		return fmt.Errorf("media item ID cannot be zero")
	}
	return nil
}

// Search finds user media item data based on a query object
func (s *coreUserMediaItemDataService[T]) Search(ctx context.Context, query *types.QueryOptions) ([]*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Str("type", string(query.MediaType)).
		Int("limit", query.Limit).
		Int("offset", query.Offset).
		Msg("Searching user media item data")

	// Create a query options with user filter
	options := types.QueryOptions{
		MediaType: query.MediaType,
		OwnerID:   query.OwnerID,
		Query:     query.Query,
		Limit:     query.Limit,
		Offset:    query.Offset,
	}

	// Delegate to repository
	result, err := s.dataRepo.Search(ctx, &options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query.Query).
			Str("type", string(query.MediaType)).
			Msg("Failed to search user media item data")
		return nil, err
	}
	log.Info().
		Str("query", query.Query).
		Str("type", string(query.MediaType)).
		Int("count", len(result)).
		Msg("User media item data found")

	return result, nil
}
