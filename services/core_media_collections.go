package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// CoreCollectionService defines the base interface for working with collections
// without client or user-specific concerns
type CoreCollectionService interface {
	// Core collection operations
	Create(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error)
	Update(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Collection], error)
	GetItemsByID(ctx context.Context, collectionID uint64) (*models.MediaItems, error)
	Delete(ctx context.Context, id uint64) error

	// Query operations
	GetByType(ctx context.Context, mediaType mediatypes.MediaType) ([]*models.MediaItem[*mediatypes.Collection], error)
	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[*mediatypes.Collection], error)
	Search(ctx context.Context, query string, limit int, offset int) ([]*models.MediaItem[*mediatypes.Collection], error)
	GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Collection], error)

	// Collection-specific operations
	GetCollectionItems(ctx context.Context, collectionID uint64) (*models.MediaItems, error)
	AddItemToCollection(ctx context.Context, collectionID uint64, item models.MediaItem[mediatypes.MediaData]) error
	RemoveItemFromCollection(ctx context.Context, collectionID uint64, itemID uint64) error
	UpdateCollectionItems(ctx context.Context, collectionID uint64, items []models.MediaItem[mediatypes.MediaData]) error
}

type coreCollectionService struct {
	repo repository.MediaItemRepository[*mediatypes.Collection]
}

// NewCoreCollectionService creates a new core collection service
func NewCoreCollectionService(
	repo repository.MediaItemRepository[*mediatypes.Collection],
) CoreCollectionService {
	return &coreCollectionService{
		repo: repo,
	}
}

// Create adds a new collection
func (s *coreCollectionService) Create(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("title", collection.Title).
		Msg("Creating collection")

	// Ensure collection-specific validation
	if collection.Type != mediatypes.MediaTypeCollection {
		collection.Type = mediatypes.MediaTypeCollection
	}

	// Ensure collection has a valid name
	if collection.Data == nil || collection.Data.Details.Title == "" {
		return nil, errors.New("collection must have a title")
	}

	// Initialize Items array if nil
	if collection.Data.Items == nil {
		collection.Data.Items = []mediatypes.ListItem{}
	}

	// Create the collection
	result, err := s.repo.Create(ctx, collection)
	if err != nil {
		log.Error().Err(err).
			Str("title", collection.Title).
			Msg("Failed to create collection")
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	log.Info().
		Uint64("id", result.ID).
		Str("title", result.Title).
		Msg("Collection created successfully")

	return result, nil
}

// Update modifies an existing collection
func (s *coreCollectionService) Update(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", collection.ID).
		Str("title", collection.Title).
		Msg("Updating collection")

	// Ensure the collection exists
	existing, err := s.GetByID(ctx, collection.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection: %w", err)
	}

	// Preserve items if not provided in the update
	if collection.Data.Items == nil || len(collection.Data.Items) == 0 {
		collection.Data.Items = existing.Data.Items
		collection.Data.ItemCount = existing.Data.ItemCount
	}

	// Update the collection
	result, err := s.repo.Update(ctx, collection)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", collection.ID).
			Msg("Failed to update collection")
		return nil, fmt.Errorf("failed to update collection: %w", err)
	}

	log.Info().
		Uint64("id", result.ID).
		Str("title", result.Title).
		Msg("Collection updated successfully")

	return result, nil
}

// GetByID retrieves a collection by its ID
func (s *coreCollectionService) GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Getting collection by ID")

	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to get collection by ID")
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	// Verify this is actually a collection
	if result.Type != mediatypes.MediaTypeCollection {
		log.Error().
			Uint64("id", id).
			Str("actualType", string(result.Type)).
			Msg("Item is not a collection")
		return nil, fmt.Errorf("item with ID %d is not a collection", id)
	}

	return result, nil
}

// Delete removes a collection
func (s *coreCollectionService) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting collection")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to delete collection")
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	log.Info().
		Uint64("id", id).
		Msg("Collection deleted successfully")

	return nil
}

// GetByType retrieves all collections of a specific type
func (s *coreCollectionService) GetByType(ctx context.Context, mediaType mediatypes.MediaType) ([]*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(mediaType)).
		Msg("Getting collections by type")

	// This method is typically used with mediaType = MediaTypeCollection
	results, err := s.repo.GetByType(ctx, mediaType)
	if err != nil {
		log.Error().Err(err).
			Str("type", string(mediaType)).
			Msg("Failed to get collections by type")
		return nil, fmt.Errorf("failed to get collections by type: %w", err)
	}

	log.Info().
		Str("type", string(mediaType)).
		Int("count", len(results)).
		Msg("Retrieved collections by type")

	return results, nil
}

// GetByExternalID retrieves a collection by its external ID
func (s *coreCollectionService) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("source", source).
		Str("externalID", externalID).
		Msg("Getting collection by external ID")

	result, err := s.repo.GetByExternalID(ctx, source, externalID)
	if err != nil {
		log.Error().Err(err).
			Str("source", source).
			Str("externalID", externalID).
			Msg("Failed to get collection by external ID")
		return nil, fmt.Errorf("failed to get collection by external ID: %w", err)
	}

	return result, nil
}

// Search finds collections based on a query string
func (s *coreCollectionService) Search(ctx context.Context, query string, limit int, offset int) ([]*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Searching collections")

	// Use repository search with collection type filter
	results, err := s.repo.Search(ctx, query, mediatypes.MediaTypeCollection, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Msg("Failed to search collections")
		return nil, fmt.Errorf("failed to search collections: %w", err)
	}

	log.Info().
		Str("query", query).
		Int("count", len(results)).
		Msg("Collections found")

	return results, nil
}

// GetRecentItems retrieves recently added collections
func (s *coreCollectionService) GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Collection], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recent collections")

	results, err := s.repo.GetRecentItems(ctx, mediatypes.MediaTypeCollection, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Int("limit", limit).
			Msg("Failed to get recent collections")
		return nil, fmt.Errorf("failed to get recent collections: %w", err)
	}

	log.Info().
		Int("count", len(results)).
		Msg("Recent collections retrieved")

	return results, nil
}

// GetCollectionItems retrieves all items in a collection
func (s *coreCollectionService) GetCollectionItems(ctx context.Context, collectionID uint64) (*models.MediaItems, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("collectionID", collectionID).
		Msg("Getting collection items")

	// Get the collection
	collection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection items: %w", err)
	}

	// Return empty result if no items
	if collection.Data.Items == nil || len(collection.Data.Items) == 0 {
		return &models.MediaItems{}, nil
	}

	// Extract item IDs
	itemIDs := collection.Data.GetItemIDs()

	// Fetch all items using the repository's batch method
	items, err := s.repo.GetMixedMediaItemsByIDs(ctx, itemIDs)
	if err != nil {
		log.Error().Err(err).
			Uint64("collectionID", collectionID).
			Msg("Failed to fetch collection items")
		return &models.MediaItems{}, fmt.Errorf("failed to fetch collection items: %w", err)
	}

	log.Info().
		Uint64("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Collection items retrieved")

	return items, nil
}

// AddItemToCollection adds an item to a collection
func (s *coreCollectionService) AddItemToCollection(ctx context.Context, collectionID uint64, item models.MediaItem[mediatypes.MediaData]) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("collectionID", collectionID).
		Uint64("itemID", item.ID).
		Msg("Adding item to collection")

	// Get the collection
	collection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("failed to add item to collection: %w", err)
	}

	// Check if the item already exists in the collection
	_, _, found := collection.Data.FindItemByID(item.ID)
	if found {
		return errors.New("item already exists in collection")
	}

	// Create a new ListItem
	newItem := mediatypes.ListItem{
		ItemID:      item.ID,
		Position:    len(collection.Data.Items),
		LastChanged: time.Now(),
	}

	// Add the item to the collection
	collection.Data.AddItem(newItem, collection.Data.ModifiedBy)

	// Update the collection
	_, err = s.Update(ctx, *collection)
	if err != nil {
		log.Error().Err(err).
			Uint64("collectionID", collectionID).
			Uint64("itemID", item.ID).
			Msg("Failed to update collection after adding item")
		return fmt.Errorf("failed to update collection after adding item: %w", err)
	}

	log.Info().
		Uint64("collectionID", collectionID).
		Uint64("itemID", item.ID).
		Msg("Item added to collection successfully")

	return nil
}

// RemoveItemFromCollection removes an item from a collection
func (s *coreCollectionService) RemoveItemFromCollection(ctx context.Context, collectionID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("collectionID", collectionID).
		Uint64("itemID", itemID).
		Msg("Removing item from collection")

	// Get the collection
	collection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("failed to remove item from collection: %w", err)
	}

	// Remove the item using the collection's method
	err = collection.Data.RemoveItem(itemID, collection.Data.ModifiedBy)
	if err != nil {
		log.Error().Err(err).
			Uint64("collectionID", collectionID).
			Uint64("itemID", itemID).
			Msg("Failed to remove item from collection")
		return fmt.Errorf("failed to remove item from collection: %w", err)
	}

	// Update the collection
	_, err = s.Update(ctx, *collection)
	if err != nil {
		log.Error().Err(err).
			Uint64("collectionID", collectionID).
			Uint64("itemID", itemID).
			Msg("Failed to update collection after removing item")
		return fmt.Errorf("failed to update collection after removing item: %w", err)
	}

	log.Info().
		Uint64("collectionID", collectionID).
		Uint64("itemID", itemID).
		Msg("Item removed from collection successfully")

	return nil
}

// UpdateCollectionItems replaces all items in a collection
func (s *coreCollectionService) UpdateCollectionItems(ctx context.Context, collectionID uint64, items []models.MediaItem[mediatypes.MediaData]) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("collectionID", collectionID).
		Int("itemCount", len(items)).
		Msg("Updating collection items")

	// Get the collection
	collection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("failed to update collection items: %w", err)
	}

	// Clear existing items
	collection.Data.Items = []mediatypes.ListItem{}

	// Add each item with position
	for i, item := range items {
		newItem := mediatypes.ListItem{
			ItemID:      item.ID,
			Position:    i,
			LastChanged: time.Now(),
		}
		collection.Data.AddItem(newItem, collection.Data.ModifiedBy)
	}

	// Ensure positions are normalized
	collection.Data.NormalizePositions()

	// Update the collection
	_, err = s.Update(ctx, *collection)
	if err != nil {
		log.Error().Err(err).
			Uint64("collectionID", collectionID).
			Msg("Failed to update collection items")
		return fmt.Errorf("failed to update collection items: %w", err)
	}

	log.Info().
		Uint64("collectionID", collectionID).
		Int("itemCount", len(items)).
		Msg("Collection items updated successfully")

	return nil
}

func (s *coreCollectionService) GetItemsByID(ctx context.Context, collectionID uint64) (*models.MediaItems, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("collectionID", collectionID).
		Msg("Getting collection items")

	// Get the collection
	collection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	// Search for
	items := collection.Data.Items
	if len(items) == 0 {
		return &models.MediaItems{}, nil
	}
	// get item ids out of the collection
	itemIDs := collection.Data.GetItemIDs()

	mediaItems, err := s.repo.GetMixedMediaItemsByIDs(ctx, itemIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get mixed media items: %w", err)
	}

	return mediaItems, nil
}
