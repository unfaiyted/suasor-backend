package services

import (
	"context"
	"fmt"

	"suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// ClientMediaItemService defines the interface for client-associated media item operations
// This service extends CoreMediaItemService with operations specific to media items
// that are linked to external clients like Plex, Emby, etc.
type ClientMediaItemService[T clienttypes.ClientMediaConfig, U types.MediaData] interface {
	// Include all core service methods
	CoreMediaItemService[U]

	// Client-specific operations
	GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[U], error)
	GetByClientItemID(ctx context.Context, itemID string, clientID uint64) (*models.MediaItem[U], error)

	// Multi-client operations
	GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[U], error)
	SearchAcrossClients(ctx context.Context, query types.QueryOptions, clientIDs []uint64) (map[uint64][]*models.MediaItem[U], error)

	// Sync operations
	SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error
}

// clientMediaItemService implements ClientMediaItemService
type clientMediaItemService[T clienttypes.ClientMediaConfig, U types.MediaData] struct {
	CoreMediaItemService[U] // Embed the core service
	itemRepo                repository.ClientMediaItemRepository[U]
}

// NewClientMediaItemService creates a new client-associated media item service
func NewClientMediaItemService[T clienttypes.ClientMediaConfig, U types.MediaData](
	coreService CoreMediaItemService[U],
	clientRepo repository.ClientRepository[T],
	itemRepo repository.ClientMediaItemRepository[U],
) ClientMediaItemService[T, U] {
	return &clientMediaItemService[T, U]{
		CoreMediaItemService: coreService,
		itemRepo:             itemRepo,
	}
}

// GetByClientID retrieves all media items associated with a specific client
func (s *clientMediaItemService[T, U]) GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[U], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Msg("Getting media items by client ID")

	// Delegate to client repository
	results, err := s.itemRepo.GetByClientID(ctx, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to get media items by client ID")
		return nil, err
	}

	log.Info().
		Uint64("clientID", clientID).
		Int("count", len(results)).
		Msg("Media items retrieved by client ID")

	return results, nil
}

// GetByClientItemID retrieves a media item by its client-specific ID
func (s *clientMediaItemService[T, U]) GetByClientItemID(ctx context.Context, itemID string, clientID uint64) (*models.MediaItem[U], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("itemID", itemID).
		Uint64("clientID", clientID).
		Msg("Getting media item by client item ID")

	// Delegate to client repository
	result, err := s.itemRepo.GetByClientItemID(ctx, itemID, clientID)
	if err != nil {
		log.Error().Err(err).
			Str("itemID", itemID).
			Uint64("clientID", clientID).
			Msg("Failed to get media item by client item ID")
		return nil, err
	}

	log.Debug().
		Str("itemID", itemID).
		Uint64("clientID", clientID).
		Uint64("id", result.ID).
		Msg("Media item retrieved by client item ID")

	return result, nil
}

// GetByMultipleClients retrieves all media items associated with any of the specified clients
func (s *clientMediaItemService[T, U]) GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[U], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Interface("clientIDs", clientIDs).
		Msg("Getting media items by multiple clients")

	// Delegate to client repository
	results, err := s.itemRepo.GetByMultipleClients(ctx, clientIDs)
	if err != nil {
		log.Error().Err(err).
			Interface("clientIDs", clientIDs).
			Msg("Failed to get media items by multiple clients")
		return nil, err
	}

	log.Info().
		Interface("clientIDs", clientIDs).
		Int("count", len(results)).
		Msg("Media items retrieved by multiple clients")

	return results, nil
}

// SearchAcrossClients searches for media items across multiple clients
// Maps by [clientID] for each of the set of MeidaItems[T]
func (s *clientMediaItemService[T, U]) SearchAcrossClients(ctx context.Context, query types.QueryOptions, clientIDs []uint64) (map[uint64][]*models.MediaItem[U], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Interface("clientIDs", clientIDs).
		Str("type", string(query.MediaType)).
		Msg("Searching media items across clients")

	// Create query options for the repository
	options := types.QueryOptions{
		MediaType: query.MediaType,
		Query:     query.Query,
	}

	var results map[uint64][]*models.MediaItem[U]

	for _, clientID := range clientIDs {
		// Delegate to client repository
		clientResult, err := s.itemRepo.Search(ctx, options)
		if err != nil {
			log.Error().Err(err).
				Str("query", query.Query).
				Interface("clientIDs", clientIDs).
				Str("type", string(query.MediaType)).
				Msg("Failed to search media items across clients")
			return nil, err
		}
		results[clientID] = clientResult
	}

	log.Info().
		Str("query", query.Query).
		Interface("clientIDs", clientIDs).
		Str("type", string(query.MediaType)).
		Int("count", len(results)).
		Msg("Media items found across clients")

	return results, nil
}

// SyncItemBetweenClients creates or updates a mapping between a media item and a target client
func (s *clientMediaItemService[T, U]) SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("itemID", itemID).
		Uint64("sourceClientID", sourceClientID).
		Uint64("targetClientID", targetClientID).
		Str("targetItemID", targetItemID).
		Msg("Syncing item between clients")

	// Delegate to client repository
	err := s.itemRepo.SyncItemBetweenClients(ctx, itemID, sourceClientID, targetClientID, targetItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("itemID", itemID).
			Uint64("sourceClientID", sourceClientID).
			Uint64("targetClientID", targetClientID).
			Msg("Failed to sync item between clients")
		return fmt.Errorf("failed to sync item between clients: %w", err)
	}

	log.Info().
		Uint64("itemID", itemID).
		Uint64("sourceClientID", sourceClientID).
		Uint64("targetClientID", targetClientID).
		Msg("Item synced between clients successfully")

	return nil
}
