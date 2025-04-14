package services

import (
	"context"
	"fmt"
	
	"suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// ClientMediaItemService defines the interface for client-associated media item operations
// This service extends CoreMediaItemService with operations specific to media items
// that are linked to external clients like Plex, Emby, etc.
type ClientMediaItemService[T types.MediaData] interface {
	// Include all core service methods
	CoreMediaItemService[T]
	
	// Client-specific operations
	GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error)
	GetByClientItemID(ctx context.Context, itemID string, clientID uint64) (*models.MediaItem[T], error)
	GetByClientAndType(ctx context.Context, clientID uint64, mediaType types.MediaType) ([]*models.MediaItem[T], error)
	
	// Multi-client operations
	GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[T], error)
	SearchAcrossClients(ctx context.Context, query string, clientIDs []uint64, mediaType types.MediaType) ([]*models.MediaItem[T], error)
	
	// Sync operations
	SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error
}

// clientMediaItemService implements ClientMediaItemService
type clientMediaItemService[T types.MediaData] struct {
	coreService CoreMediaItemService[T] // Embed the core service
	clientRepo  repository.ClientMediaItemRepository[T]
}

// NewClientMediaItemService creates a new client-associated media item service
func NewClientMediaItemService[T types.MediaData](
	coreService CoreMediaItemService[T],
	clientRepo repository.ClientMediaItemRepository[T],
) ClientMediaItemService[T] {
	return &clientMediaItemService[T]{
		coreService: coreService,
		clientRepo:  clientRepo,
	}
}

// Core service methods - delegate to embedded core service

func (s *clientMediaItemService[T]) Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	return s.coreService.Create(ctx, item)
}

func (s *clientMediaItemService[T]) Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error) {
	return s.coreService.Update(ctx, item)
}

func (s *clientMediaItemService[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	return s.coreService.GetByID(ctx, id)
}

func (s *clientMediaItemService[T]) Delete(ctx context.Context, id uint64) error {
	return s.coreService.Delete(ctx, id)
}

func (s *clientMediaItemService[T]) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error) {
	return s.coreService.GetByExternalID(ctx, source, externalID)
}

func (s *clientMediaItemService[T]) GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error) {
	return s.coreService.GetByType(ctx, mediaType)
}

func (s *clientMediaItemService[T]) Search(ctx context.Context, query string, mediaType types.MediaType, limit int, offset int) ([]*models.MediaItem[T], error) {
	return s.coreService.Search(ctx, query, mediaType, limit, offset)
}

func (s *clientMediaItemService[T]) GetRecentItems(ctx context.Context, mediaType types.MediaType, days int, limit int) ([]*models.MediaItem[T], error) {
	return s.coreService.GetRecentItems(ctx, mediaType, days, limit)
}

// Client-specific methods

// GetByClientID retrieves all media items associated with a specific client
func (s *clientMediaItemService[T]) GetByClientID(ctx context.Context, clientID uint64) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Msg("Getting media items by client ID")
		
	// Delegate to client repository
	results, err := s.clientRepo.GetByClientID(ctx, clientID)
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
func (s *clientMediaItemService[T]) GetByClientItemID(ctx context.Context, itemID string, clientID uint64) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("itemID", itemID).
		Uint64("clientID", clientID).
		Msg("Getting media item by client item ID")
		
	// Delegate to client repository
	result, err := s.clientRepo.GetByClientItemID(ctx, itemID, clientID)
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

// GetByClientAndType retrieves all media items of a specific type from a client
func (s *clientMediaItemService[T]) GetByClientAndType(ctx context.Context, clientID uint64, mediaType types.MediaType) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Str("type", string(mediaType)).
		Msg("Getting media items by client and type")
		
	// Delegate to client repository
	results, err := s.clientRepo.GetByType(ctx, mediaType, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("type", string(mediaType)).
			Msg("Failed to get media items by client and type")
		return nil, err
	}
	
	log.Info().
		Uint64("clientID", clientID).
		Str("type", string(mediaType)).
		Int("count", len(results)).
		Msg("Media items retrieved by client and type")
		
	return results, nil
}

// GetByMultipleClients retrieves all media items associated with any of the specified clients
func (s *clientMediaItemService[T]) GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Interface("clientIDs", clientIDs).
		Msg("Getting media items by multiple clients")
		
	// Delegate to client repository
	results, err := s.clientRepo.GetByMultipleClients(ctx, clientIDs)
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
func (s *clientMediaItemService[T]) SearchAcrossClients(ctx context.Context, query string, clientIDs []uint64, mediaType types.MediaType) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query).
		Interface("clientIDs", clientIDs).
		Str("type", string(mediaType)).
		Msg("Searching media items across clients")
		
	// Create query options for the repository
	options := types.QueryOptions{
		MediaType: mediaType,
		Query:     query,
		ClientIDs: clientIDs,
	}
	
	// Delegate to client repository
	results, err := s.clientRepo.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Interface("clientIDs", clientIDs).
			Str("type", string(mediaType)).
			Msg("Failed to search media items across clients")
		return nil, err
	}
	
	log.Info().
		Str("query", query).
		Interface("clientIDs", clientIDs).
		Str("type", string(mediaType)).
		Int("count", len(results)).
		Msg("Media items found across clients")
		
	return results, nil
}

// SyncItemBetweenClients creates or updates a mapping between a media item and a target client
func (s *clientMediaItemService[T]) SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("itemID", itemID).
		Uint64("sourceClientID", sourceClientID).
		Uint64("targetClientID", targetClientID).
		Str("targetItemID", targetItemID).
		Msg("Syncing item between clients")
		
	// Delegate to client repository
	err := s.clientRepo.SyncItemBetweenClients(ctx, itemID, sourceClientID, targetClientID, targetItemID)
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