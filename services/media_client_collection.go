// services/media_client_collection.go
package services

import (
	"context"
	"errors"
	"sort"

	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// MediaClientCollectionService defines operations for interacting with collection clients
type MediaClientCollectionService[T types.ClientConfig] interface {
	GetCollectionByID(ctx context.Context, userID uint64, clientID uint64, collectionID string) (*models.MediaItem[*mediatypes.Collection], error)
	GetCollections(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Collection], error)
}

type mediaCollectionService[T types.MediaClientConfig] struct {
	repo    repository.ClientRepository[T]
	factory *client.ClientFactoryService
}

// NewMediaClientCollectionService creates a new media collection service
func NewMediaClientCollectionService[T types.MediaClientConfig](
	repo repository.ClientRepository[T],
	factory *client.ClientFactoryService,
) MediaClientCollectionService[T] {
	return &mediaCollectionService[T]{
		repo:    repo,
		factory: factory,
	}
}

// getCollectionClients gets all collection clients for a user
func (s *mediaCollectionService[T]) getCollectionClients(ctx context.Context, userID uint64) ([]media.MediaClient, error) {
	repo := s.repo
	// Get all media clients for the user
	clients, err := repo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	var collectionClients []media.MediaClient

	// Filter and instantiate clients that support collections
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsCollections() {
			clientId := clientConfig.GetID()
			client, err := s.factory.GetClient(ctx, clientId, clientConfig.Config.Data)
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			collectionClients = append(collectionClients, client.(media.MediaClient))
		}
	}

	return collectionClients, nil
}

// getSpecificCollectionClient gets a specific collection client
func (s *mediaCollectionService[T]) getSpecificCollectionClient(ctx context.Context, userID, clientID uint64) (media.MediaClient, error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := (s.repo).GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client config")

	if !clientConfig.Config.Data.SupportsCollections() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Client does not support collections")
		return nil, ErrUnsupportedFeature
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Client supports collections")

	client, err := s.factory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client")
	return client.(media.MediaClient), nil
}

func (s *mediaCollectionService[T]) GetCollectionByID(ctx context.Context, userID uint64, clientID uint64, collectionID string) (*models.MediaItem[*mediatypes.Collection], error) {
	client, err := s.getSpecificCollectionClient(ctx, userID, clientID)
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("collectionID", collectionID).
		Msg("Retrieving collection")

	collectionProvider, ok := client.(providers.CollectionProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("collectionID", collectionID).
			Msg("Client does not support collections")
		return nil, ErrUnsupportedFeature
	}

	// Check if the client supports getting collection by ID
	if !collectionProvider.SupportsCollections() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("collectionID", collectionID).
			Msg("Client does not support collections")
		return nil, ErrUnsupportedFeature
	}

	// Get all collections and find by ID
	options := &mediatypes.QueryOptions{
		ExternalSourceID: collectionID,
	}

	collections, err := collectionProvider.GetCollections(ctx, options)
	if err != nil {
		return nil, err
	}

	// Check if we found any collections
	if len(collections) == 0 {
		return nil, errors.New("collection not found")
	}

	// Return the first matching collection
	return &collections[0], nil
}

func (s *mediaCollectionService[T]) GetCollections(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Collection], error) {
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

	return allCollections, nil
}
