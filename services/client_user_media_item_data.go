package services

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
	"time"
)

// ClientUserMediaItemDataService defines the client service interface for user media item data
// This service focuses on client-specific operations including synchronization with external media systems
type ClientUserMediaItemDataService[T types.MediaData] interface {
	// Embed the user service methods
	UserMediaItemDataService[T]

	// SyncClientItemData synchronizes user media item data from an external client
	SyncClientItemData(ctx context.Context, userID uint64, clientID uint64, items []models.UserMediaItemData[T]) error

	// GetClientItemData retrieves user media item data for synchronization with a client
	GetClientItemData(ctx context.Context, userID uint64, clientID uint64, since *string) ([]*models.UserMediaItemData[T], error)

	// GetByClientID retrieves a user media item data entry by client ID
	GetByClientID(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error)

	// RecordClientPlay records a play event from a client
	RecordClientPlay(ctx context.Context, userID uint64, clientID uint64, clientItemID string, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)

	// GetPlaybackState retrieves the current playback state for a client item
	GetPlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error)

	// UpdatePlaybackState updates the playback state for a client item
	UpdatePlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[T], error)
}

// clientUserMediaItemDataService implements ClientUserMediaItemDataService
type clientUserMediaItemDataService[T types.MediaData] struct {
	UserMediaItemDataService[T]
	repo repository.ClientUserMediaItemDataRepository[T]
}

// NewClientUserMediaItemDataService creates a new client user media item data service
func NewClientUserMediaItemDataService[T types.MediaData](
	userService UserMediaItemDataService[T],
	repo repository.ClientUserMediaItemDataRepository[T],
) ClientUserMediaItemDataService[T] {
	return &clientUserMediaItemDataService[T]{
		UserMediaItemDataService: userService,
		repo:                     repo,
	}
}

// SyncClientItemData synchronizes user media item data from an external client
func (s *clientUserMediaItemDataService[T]) SyncClientItemData(ctx context.Context, userID uint64, clientID uint64, items []models.UserMediaItemData[T]) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Int("itemCount", len(items)).
		Msg("Synchronizing client media item data")

	// Delegate to repository
	err := s.repo.SyncClientItemData(ctx, userID, clientID, items)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to synchronize client media item data")
		return fmt.Errorf("failed to synchronize client media item data: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Int("itemCount", len(items)).
		Msg("Client media item data synchronized successfully")

	return nil
}

// GetClientItemData retrieves user media item data for synchronization with a client
func (s *clientUserMediaItemDataService[T]) GetClientItemData(ctx context.Context, userID uint64, clientID uint64, sinceDateStr *string) ([]*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("since", s.getSinceString(sinceDateStr)).
		Msg("Getting client media item data")

	// conver str to date
	sinceDate, err := time.Parse(time.RFC3339, *sinceDateStr)

	// Delegate to repository
	result, err := s.repo.GetClientItemData(ctx, userID, clientID, sinceDate)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to get client media item data")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Int("count", len(result)).
		Msg("Client media item data retrieved successfully")

	return result, nil
}

// GetByClientID retrieves a user media item data entry by client ID
func (s *clientUserMediaItemDataService[T]) GetByClientID(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Getting user media item data by client ID")

	// Delegate to repository
	result, err := s.repo.GetByClientID(ctx, userID, clientID, clientItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to get user media item data by client ID")
		return nil, err
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Uint64("mediaItemID", result.MediaItemID).
		Msg("User media item data retrieved by client ID")

	return result, nil
}

// RecordClientPlay records a play event from a client
func (s *clientUserMediaItemDataService[T]) RecordClientPlay(ctx context.Context, userID uint64, clientID uint64, clientItemID string, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Recording client play event")

	// Delegate to repository
	result, err := s.repo.RecordClientPlay(ctx, userID, clientID, clientItemID, data)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to record client play event")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Uint64("id", result.ID).
		Msg("Client play event recorded successfully")

	return result, nil
}

// GetPlaybackState retrieves the current playback state for a client item
func (s *clientUserMediaItemDataService[T]) GetPlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Getting playback state")

	// Delegate to repository
	result, err := s.repo.GetPlaybackState(ctx, userID, clientID, clientItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to get playback state")
		return nil, err
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Float64("percentage", result.PlayedPercentage).
		Msg("Playback state retrieved successfully")

	return result, nil
}

// UpdatePlaybackState updates the playback state for a client item
func (s *clientUserMediaItemDataService[T]) UpdatePlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Int("position", position).
		Int("duration", duration).
		Float64("percentage", percentage).
		Msg("Updating playback state")

	// Delegate to repository
	result, err := s.repo.UpdatePlaybackState(ctx, userID, clientID, clientItemID, position, duration, percentage)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to update playback state")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Float64("percentage", percentage).
		Msg("Playback state updated successfully")

	return result, nil
}

// Helper methods

func (s *clientUserMediaItemDataService[T]) getSinceString(sinceDateStr *string) string {
	if sinceDateStr == nil {
		return "24 hours ago"
	}
	return *sinceDateStr
}
