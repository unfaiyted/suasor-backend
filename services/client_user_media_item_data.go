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
	userService UserMediaItemDataService[T]
	repo        repository.ClientUserMediaItemDataRepository[T]
}

// NewClientUserMediaItemDataService creates a new client user media item data service
func NewClientUserMediaItemDataService[T types.MediaData](
	userService UserMediaItemDataService[T],
	repo repository.ClientUserMediaItemDataRepository[T],
) ClientUserMediaItemDataService[T] {
	return &clientUserMediaItemDataService[T]{
		userService: userService,
		repo:        repo,
	}
}

// User service methods - delegate to embedded user service

func (s *clientUserMediaItemDataService[T]) Create(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	return s.userService.Create(ctx, data)
}

func (s *clientUserMediaItemDataService[T]) GetByID(ctx context.Context, id uint64) (*models.UserMediaItemData[T], error) {
	return s.userService.GetByID(ctx, id)
}

func (s *clientUserMediaItemDataService[T]) Update(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	return s.userService.Update(ctx, data)
}

func (s *clientUserMediaItemDataService[T]) Delete(ctx context.Context, id uint64) error {
	return s.userService.Delete(ctx, id)
}

func (s *clientUserMediaItemDataService[T]) GetByUserIDAndMediaItemID(ctx context.Context, userID, mediaItemID uint64) (*models.UserMediaItemData[T], error) {
	return s.userService.GetByUserIDAndMediaItemID(ctx, userID, mediaItemID)
}

func (s *clientUserMediaItemDataService[T]) HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
	return s.userService.HasUserMediaItemData(ctx, userID, mediaItemID)
}

func (s *clientUserMediaItemDataService[T]) GetUserHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error) {
	return s.userService.GetUserHistory(ctx, userID, limit, offset, mediaType)
}

func (s *clientUserMediaItemDataService[T]) GetRecentHistory(ctx context.Context, userID uint64, limit int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error) {
	return s.userService.GetRecentHistory(ctx, userID, limit, mediaType)
}

func (s *clientUserMediaItemDataService[T]) GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) ([]*models.UserMediaItemData[T], error) {
	return s.userService.GetUserPlayHistory(ctx, userID, limit, offset, mediaType, completed)
}

func (s *clientUserMediaItemDataService[T]) GetContinueWatching(ctx context.Context, userID uint64, limit int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error) {
	return s.userService.GetContinueWatching(ctx, userID, limit, mediaType)
}

func (s *clientUserMediaItemDataService[T]) RecordPlay(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
	return s.userService.RecordPlay(ctx, data)
}

func (s *clientUserMediaItemDataService[T]) ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error {
	return s.userService.ToggleFavorite(ctx, mediaItemID, userID, favorite)
}

func (s *clientUserMediaItemDataService[T]) UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error {
	return s.userService.UpdateRating(ctx, mediaItemID, userID, rating)
}

func (s *clientUserMediaItemDataService[T]) GetFavorites(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error) {
	return s.userService.GetFavorites(ctx, userID, limit, offset)
}

func (s *clientUserMediaItemDataService[T]) ClearUserHistory(ctx context.Context, userID uint64) error {
	return s.userService.ClearUserHistory(ctx, userID)
}

// Client-specific methods

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
