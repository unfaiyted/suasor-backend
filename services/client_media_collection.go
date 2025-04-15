// services/client_media_collection.go
package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// ClientMediaCollectionService defines the interface for client-associated collection operations
// This service extends CoreCollectionService with operations specific to media collections
// that are linked to external clients like Plex, Emby, etc.
type ClientMediaCollectionService interface {
	// Include all core service methods
	CoreCollectionService

	// Client-specific operations
	GetCollectionByClientID(ctx context.Context, clientID uint64, collectionID string) (*models.MediaItem[*mediatypes.Collection], error)
	GetCollectionsByClient(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Collection], error)
	GetCollectionsByMultipleClients(ctx context.Context, clientIDs []uint64, limit int) (map[uint64][]*models.MediaItem[*mediatypes.Collection], error)

	// Sync operations
	SyncCollectionBetweenClients(ctx context.Context, collectionID uint64, sourceClientID uint64, targetClientID uint64) error

	// Legacy operations (to maintain compatibility)
	GetCollectionByID(ctx context.Context, userID uint64, clientID uint64, collectionID string) (*models.MediaItem[*mediatypes.Collection], error)
	GetCollections(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Collection], error)
}

type clientCollectionService struct {
	coreService    CoreCollectionService
	collectionRepo repository.ClientMediaItemRepository[*mediatypes.Collection]
	clientRepo     repository.ClientRepository[types.ClientMediaConfig]
	factory        *client.ClientFactoryService
}

// NewClientMediaCollectionService creates a new client collection service
func NewClientMediaCollectionService(
	coreService CoreCollectionService,
	collectionRepo repository.ClientMediaItemRepository[*mediatypes.Collection],
	clientRepo repository.ClientRepository[types.ClientMediaConfig],
	factory *client.ClientFactoryService,
) ClientMediaCollectionService {
	return &clientCollectionService{
		coreService:    coreService,
		collectionRepo: collectionRepo,
		clientRepo:     clientRepo,
		factory:        factory,
	}
}

// Implement all methods from CoreCollectionService through delegation

// Create adds a new collection
func (s *clientCollectionService) Create(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error) {
	return s.coreService.Create(ctx, collection)
}

// Update modifies an existing collection
func (s *clientCollectionService) Update(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error) {
	return s.coreService.Update(ctx, collection)
}

// GetByID retrieves a collection by its ID
func (s *clientCollectionService) GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Collection], error) {
	return s.coreService.GetByID(ctx, id)
}

func (s *clientCollectionService) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[*mediatypes.Collection], error) {
	return s.coreService.GetAll(ctx, limit, offset)
}

// Delete removes a collection
func (s *clientCollectionService) Delete(ctx context.Context, id uint64) error {
	return s.coreService.Delete(ctx, id)
}

// GetByType retrieves all collections of a specific type
func (s *clientCollectionService) GetByType(ctx context.Context, mediaType mediatypes.MediaType) ([]*models.MediaItem[*mediatypes.Collection], error) {
	return s.coreService.GetByType(ctx, mediaType)
}

// GetByExternalID retrieves a collection by its external ID
func (s *clientCollectionService) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[*mediatypes.Collection], error) {
	return s.coreService.GetByExternalID(ctx, source, externalID)
}

// Search finds collections based on a query string
func (s *clientCollectionService) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
	return s.coreService.Search(ctx, query)
}

// GetRecentItems retrieves recently added collections
func (s *clientCollectionService) GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Collection], error) {
	return s.coreService.GetRecentItems(ctx, days, limit)
}

// GetCollectionItems retrieves all items in a collection
func (s *clientCollectionService) GetCollectionItems(ctx context.Context, collectionID uint64) (*models.MediaItems, error) {
	return s.coreService.GetCollectionItems(ctx, collectionID)
}

// AddItemToCollection adds an item to a collection
func (s *clientCollectionService) AddItemToCollection(ctx context.Context, collectionID uint64, itemID uint64) error {
	return s.coreService.AddItemToCollection(ctx, collectionID, itemID)
}

// RemoveItemFromCollection removes an item from a collection
func (s *clientCollectionService) RemoveItemFromCollection(ctx context.Context, collectionID uint64, itemID uint64) error {
	return s.coreService.RemoveItemFromCollection(ctx, collectionID, itemID)
}

// UpdateCollectionItems replaces all items in a collection
func (s *clientCollectionService) UpdateCollectionItems(ctx context.Context, collectionID uint64, items []models.MediaItem[mediatypes.MediaData]) error {
	return s.coreService.UpdateCollectionItems(ctx, collectionID, items)
}

// getCollectionClients gets all collection clients for a user
func (s *clientCollectionService) getCollectionClients(ctx context.Context, userID uint64) ([]media.ClientMedia, error) {
	log := utils.LoggerFromContext(ctx)
	// Get all media clients for the user
	clients, err := s.clientRepo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get media clients for user")
		return nil, err
	}

	var collectionClients []media.ClientMedia

	// Filter and instantiate clients that support collections
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsCollections() {
			clientId := clientConfig.GetID()
			client, err := s.factory.GetClient(ctx, clientId, clientConfig.Config.Data)
			if err != nil {
				// Log error but continue with other clients
				log.Warn().Err(err).
					Uint64("clientID", clientId).
					Msg("Failed to initialize client, skipping")
				continue
			}
			clientMedia, ok := client.(media.ClientMedia)
			if !ok {
				log.Warn().
					Uint64("clientID", clientId).
					Msg("Client is not a media client, skipping")
				continue
			}
			collectionClients = append(collectionClients, clientMedia)
		}
	}

	return collectionClients, nil
}

// getSpecificCollectionClient gets a specific collection client
func (s *clientCollectionService) getSpecificCollectionClient(ctx context.Context, clientID uint64) (media.ClientMedia, error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to get client config")
		return nil, err
	}
	log.Debug().
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client config")

	if !clientConfig.Config.Data.SupportsCollections() {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Client does not support collections")
		return nil, ErrUnsupportedFeature
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Client supports collections")

	client, err := s.factory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to initialize client")
		return nil, err
	}

	clientMedia, ok := client.(media.ClientMedia)
	if !ok {
		log.Error().
			Uint64("clientID", clientID).
			Msg("Client is not a media client")
		return nil, errors.New("client is not a media client")
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client")

	return clientMedia, nil
}

// GetCollectionByClientID retrieves a collection by its client-specific ID
func (s *clientCollectionService) GetCollectionByClientID(ctx context.Context, clientID uint64, collectionID string) (*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Str("collectionID", collectionID).
		Msg("Getting collection by client ID")

	// I already have the collectionID for the client and the client ID we just need to search MediaItem table to see if it exists
	workingCollection, err := s.collectionRepo.GetByClientItemID(ctx, collectionID, clientID)
	if err == nil {
		log.Info().
			Uint64("clientID", clientID).
			Str("collectionID", collectionID).
			Uint64("internalID", workingCollection.ID).
			Msg("Found existing collection in database")
		return workingCollection, nil
	}

	// If not found in database, fetch directly from client
	client, err := s.getSpecificCollectionClient(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection client: %w", err)
	}

	collectionProvider, ok := client.(providers.CollectionProvider)
	if !ok || !collectionProvider.SupportsCollections() {
		log.Warn().
			Uint64("clientID", clientID).
			Str("collectionID", collectionID).
			Msg("Client does not support collections")
		return nil, ErrUnsupportedFeature
	}

	// Get the collection directly from the client
	options := &mediatypes.QueryOptions{
		ItemIDs: collectionID,
		Limit:   1,
	}

	collections, err := collectionProvider.GetCollections(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("collectionID", collectionID).
			Msg("Failed to get collection from client")
		return nil, err
	}

	// Check if we found any collections
	if len(collections) == 0 {
		log.Warn().
			Uint64("clientID", clientID).
			Str("collectionID", collectionID).
			Msg("Collection not found on client")
		return nil, errors.New("collection not found")
	}

	// Get the first matching collection
	newCollection := collections[0]

	// Persist or update the collection in our database
	savedCollection, err := s.Create(ctx, newCollection)
	if err != nil {
		// Log the error but continue with the client result
		log.Warn().Err(err).
			Uint64("clientID", clientID).
			Str("collectionID", collectionID).
			Msg("Failed to save collection to database, continuing with client result")
		return &newCollection, nil
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("collectionID", collectionID).
		Uint64("id", savedCollection.ID).
		Msg("Collection saved successfully")
	return savedCollection, nil
}

// GetCollectionsByClient retrieves all collections from a specific client
func (s *clientCollectionService) GetCollectionsByClient(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Getting collections by client")

	client, err := s.getSpecificCollectionClient(ctx, clientID) // userID 0 since we're looking by clientID directly
	if err != nil {
		return nil, fmt.Errorf("failed to get collection client: %w", err)
	}

	collectionProvider, ok := client.(providers.CollectionProvider)
	if !ok || !collectionProvider.SupportsCollections() {
		log.Warn().
			Uint64("clientID", clientID).
			Msg("Client does not support collections")
		return nil, ErrUnsupportedFeature
	}

	options := &mediatypes.QueryOptions{
		Limit: limit,
	}

	clientCollections, err := collectionProvider.GetCollections(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to get collections from client")
		return nil, err
	}

	// Convert to pointer slice
	collections := make([]*models.MediaItem[*mediatypes.Collection], len(clientCollections))
	for i := range clientCollections {
		collections[i] = &clientCollections[i]
	}

	log.Info().
		Uint64("clientID", clientID).
		Int("count", len(collections)).
		Msg("Retrieved collections from client")

	return collections, nil
}

// GetCollectionsByMultipleClients retrieves collections from multiple clients
func (s *clientCollectionService) GetCollectionsByMultipleClients(ctx context.Context, clientIDs []uint64, limit int) (map[uint64][]*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Interface("clientIDs", clientIDs).
		Int("limit", limit).
		Msg("Getting collections from multiple clients")

	result := make(map[uint64][]*models.MediaItem[*mediatypes.Collection])

	// Get collections from each client
	for _, clientID := range clientIDs {
		clientCollections, err := s.GetCollectionsByClient(ctx, clientID, limit)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to get collections from client, continuing with others")
			continue
		}

		// Store collections mapped to this client ID
		result[clientID] = clientCollections
	}

	log.Info().
		Int("clientCount", len(clientIDs)).
		Int("resultClientCount", len(result)).
		Msg("Retrieved collections from multiple clients")

	return result, nil
}

// SyncCollectionBetweenClients syncs a collection between two clients
func (s *clientCollectionService) SyncCollectionBetweenClients(ctx context.Context, collectionID uint64, sourceClientID uint64, targetClientID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("collectionID", collectionID).
		Uint64("sourceClientID", sourceClientID).
		Uint64("targetClientID", targetClientID).
		Msg("Syncing collection between clients")

	// Get the collection
	workingCollection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	// Check syncClients to see if they already have ClientIDs for this collection
	sourceSyncClient, exists := workingCollection.SyncClients.GetByClientID(sourceClientID)
	if !exists {
		return fmt.Errorf("source client does not have a ClientID for this collection")
	}
	sourceClientID = sourceSyncClient.ID
	targetSyncClient, targetExists := workingCollection.SyncClients.GetByClientID(targetClientID)
	targetClient, err := s.getSpecificCollectionClient(ctx, targetClientID)

	targetProvider, ok := targetClient.(providers.CollectionProvider)
	if !ok || !targetProvider.SupportsCollections() {
		return errors.New("target client does not support collections")
	}

	targetCollection, err := targetProvider.CreateCollection(ctx, workingCollection.Data.Details.Title, "", "")
	if err != nil {
		return fmt.Errorf("failed to get create target collection: %w", err)
	}
	if !targetExists {
		itemID := targetCollection.SyncClients.GetClientItemID(targetClientID)
		workingCollection.SyncClients.AddClient(targetClientID, targetClient.GetConfig().GetType(), itemID)
	}
	targetClientID = targetSyncClient.ID

	// Get source and target clients
	sourceClient, err := s.getSpecificCollectionClient(ctx, sourceClientID)
	if err != nil {
		return fmt.Errorf("failed to get source client: %w", err)
	}

	// Ensure both clients support collections
	sourceProvider, ok := sourceClient.(providers.CollectionProvider)
	if !ok || !sourceProvider.SupportsCollections() {
		return errors.New("source client does not support collections")
	}

	// TODO: Needs to pull the collection itesm from the source client. Then get the items details.
	// Try to match those up with the target client. We have a media sync method so we should be able to search
	// for the items. as par of the MediaItems struct we keep the ID's of the items in the SyncClientStates
	// when we get the sourceClient we should update the SyncClientStates with the ID's of the items in the client and do the same for the// target upon update
	// Get collection items
	// items, err := s.GetCollectionItems(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("failed to get collection items: %w", err)
	}

	// targetCollection.SyncClients.AddClient(targetClientID, targetClient.GetConfig().GetType(), targetCollection.ID)
	targetCollection.Data.LastSynced = time.Now()

	// Check if the target client already has this collection or if we have the collection ID stored in the media items

	// Create the collection in target client
	// name string, description string, collectionType string
	newCollection, err := targetProvider.CreateCollection(ctx, targetCollection.Data.Details.Title, "", "")
	if err != nil {
		log.Error().Err(err).
			Uint64("sourceClientID", sourceClientID).
			Uint64("targetClientID", targetClientID).
			Msg("Failed to create collection in target client")
		return fmt.Errorf("failed to create collection in target client: %w", err)
	}

	// Add each item to the new collection
	// for _, item := range items.Items {
	// 	// Find item in target provider
	// 	targetItem, err := targetProvider.GetCollectionItems(ctx, item.ID, options)
	// 	if err != nil {
	// 		// Log error but continue with other items
	// 		log.Warn().Err(err).
	// 			Uint64("collectionID", newCollection.ID).
	// 			Uint64("itemID", item.ID).
	// 			Msg("Failed to get item in target client, continuing with others")
	// 		continue
	// 	}
	//
	// 	err = targetProvider.AddItemToCollection(ctx, newCollection.ID, item.ID)
	// 	if err != nil {
	// 		// Log error but continue with other items
	// 		log.Warn().Err(err).
	// 			Uint64("collectionID", newCollection.ID).
	// 			Uint64("itemID", item.ID).
	// 			Msg("Failed to add item to collection in target client, continuing with others")
	// 	}
	// }

	log.Info().
		Uint64("sourceCollectionID", collectionID).
		Uint64("targetCollectionID", newCollection.ID).
		Uint64("sourceClientID", sourceClientID).
		Uint64("targetClientID", targetClientID).
		// Int("itemCount", len(items.Items)).
		Msg("Collection synced successfully between clients")

	return nil
}

// Legacy methods to maintain compatibility

// GetCollectionByID retrieves a collection by its ID (legacy method)
func (s *clientCollectionService) GetCollectionByID(ctx context.Context, userID uint64, clientID uint64, collectionID string) (*models.MediaItem[*mediatypes.Collection], error) {
	return s.GetCollectionByClientID(ctx, clientID, collectionID)
}

// GetCollections retrieves collections for a user (legacy method)
func (s *clientCollectionService) GetCollections(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)

	clients, err := s.getCollectionClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allCollections []models.MediaItem[*mediatypes.Collection]

	for _, client := range clients {
		collectionProvider, ok := client.(providers.CollectionProvider)
		if !ok || !collectionProvider.SupportsCollections() {
			continue
		}

		options := &mediatypes.QueryOptions{
			Limit: count,
		}

		collections, err := collectionProvider.GetCollections(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().Err(err).
				Str("clientType", client.GetConfig().GetType().String()).
				Msg("Failed to get collections from client, continuing with others")
			continue
		}

		allCollections = append(allCollections, collections...)
	}

	// Sort by added date
	sort.Slice(allCollections, func(i, j int) bool {
		return allCollections[i].Data.GetDetails().AddedAt.After(allCollections[j].Data.GetDetails().AddedAt)
	})

	// Limit to requested count if specified
	if count > 0 && len(allCollections) > count {
		allCollections = allCollections[:count]
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(allCollections)).
		Msg("Retrieved collections for user")

	return allCollections, nil
}
