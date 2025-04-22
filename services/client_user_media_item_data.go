package services

import (
	"context"
	"fmt"
	"suasor/clients"
	"suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"
)

// ClientUserMediaItemDataService defines the client service interface for user media item data
// This service focuses on client-specific operations including synchronization with external media systems
type ClientUserMediaItemDataService[T clienttypes.ClientMediaConfig, U types.MediaData] interface {
	// Embed the user service methods
	UserMediaItemDataService[U]

	// SyncClientItemData synchronizes user media item data from an external client
	SyncClientItemData(ctx context.Context, userID uint64, clientID uint64, items []models.UserMediaItemData[U]) error

	// GetClientItemData retrieves user media item data for synchronization with a client
	GetClientItemData(ctx context.Context, userID uint64, clientID uint64, since *string) ([]*models.UserMediaItemData[U], error)

	// GetByClientID retrieves a user media item data entry by client ID
	GetByClientID(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[U], error)

	// RecordClientPlay records a play event from a client
	RecordClientPlay(ctx context.Context, userID uint64, clientID uint64, clientItemID string, data *models.UserMediaItemData[U]) (*models.UserMediaItemData[U], error)

	// GetPlaybackState retrieves the current playback state for a client item
	GetPlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[U], error)

	// UpdatePlaybackState updates the playback state for a client item
	UpdatePlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[U], error)
}

// clientUserMediaItemDataService implements ClientUserMediaItemDataService
type clientUserMediaItemDataService[T clienttypes.ClientMediaConfig, U types.MediaData] struct {
	UserMediaItemDataService[U]
	dataRepo      repository.ClientUserMediaItemDataRepository[U]
	clientRepo    repository.ClientRepository[T]
	clientFactory *clients.ClientProviderFactoryService
}

// NewClientUserMediaItemDataService creates a new client user media item data service
func NewClientUserMediaItemDataService[T clienttypes.ClientMediaConfig, U types.MediaData](
	userService UserMediaItemDataService[U],
	dataRepo repository.ClientUserMediaItemDataRepository[U],
	clientRepo repository.ClientRepository[T],
	clientFactory *clients.ClientProviderFactoryService,
) ClientUserMediaItemDataService[T, U] {
	return &clientUserMediaItemDataService[T, U]{
		UserMediaItemDataService: userService,
		dataRepo:                 dataRepo,
		clientRepo:               clientRepo,
		clientFactory:            clientFactory,
	}
}

// SyncClientItemData synchronizes user media item data from an external client
func (s *clientUserMediaItemDataService[T, U]) SyncClientItemData(ctx context.Context, userID uint64, clientID uint64, items []models.UserMediaItemData[U]) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Int("itemCount", len(items)).
		Msg("Synchronizing client media item data")

	// Delegate to repository
	err := s.dataRepo.SyncClientItemData(ctx, userID, clientID, items)
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
func (s *clientUserMediaItemDataService[T, U]) GetClientItemData(ctx context.Context, userID uint64, clientID uint64, sinceDateStr *string) ([]*models.UserMediaItemData[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("since", s.getSinceString(sinceDateStr)).
		Msg("Getting client media item data")

	// conver str to date
	sinceDate, err := time.Parse(time.RFC3339, *sinceDateStr)

	// Delegate to repository
	result, err := s.dataRepo.GetClientItemData(ctx, userID, clientID, sinceDate)
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
func (s *clientUserMediaItemDataService[T, U]) GetByClientID(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Getting user media item data by client ID")

	// Delegate to repository
	result, err := s.dataRepo.GetByClientID(ctx, userID, clientID, clientItemID)
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
func (s *clientUserMediaItemDataService[T, U]) RecordClientPlay(ctx context.Context, userID uint64, clientID uint64, clientItemID string, data *models.UserMediaItemData[U]) (*models.UserMediaItemData[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Recording client play event")

	// Delegate to repository
	result, err := s.dataRepo.RecordClientPlay(ctx, userID, clientID, clientItemID, data)
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
func (s *clientUserMediaItemDataService[T, U]) GetPlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Getting playback state")

	// Delegate to repository
	result, err := s.dataRepo.GetPlaybackState(ctx, userID, clientID, clientItemID)
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
func (s *clientUserMediaItemDataService[T, U]) UpdatePlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Int("position", position).
		Int("duration", duration).
		Float64("percentage", percentage).
		Msg("Updating playback state")

	// Delegate to repository
	result, err := s.dataRepo.UpdatePlaybackState(ctx, userID, clientID, clientItemID, position, duration, percentage)
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

func (s *clientUserMediaItemDataService[T, U]) getSinceString(sinceDateStr *string) string {
	if sinceDateStr == nil {
		return "24 hours ago"
	}
	return *sinceDateStr
}
