package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
)

// CollectionService manages application-specific collection operations beyond the basic CRUD operations
type CollectionService interface {
	// Base operations (wrapping MediaItemService)
	Create(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error)
	Update(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Collection], error)
	GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Collection], error)
	Delete(ctx context.Context, id uint64) error

	// Collection-specific operations
	GetFeaturedCollections(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Collection], error)
	GetSmartCollections(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Collection], error)
	GetCollectionItems(ctx context.Context, collectionID uint64) (*models.MediaItems, error)
	AddItemToCollection(ctx context.Context, collectionID uint64, item models.MediaItem[mediatypes.MediaData]) error
	RemoveItemFromCollection(ctx context.Context, collectionID uint64, itemID uint64) error
	UpdateCollectionItems(ctx context.Context, collectionID uint64, items []models.MediaItem[mediatypes.MediaData]) error
	SearchCollections(ctx context.Context, query string, userID uint64) ([]*models.MediaItem[*mediatypes.Collection], error)

	// Additional methods needed by handler
	GetCollections(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Collection], error)
	GetCollectionByID(ctx context.Context, userID uint64, clientID uint64, collectionID uint64) (*models.MediaItem[*mediatypes.Collection], error)
	GetItems(ctx context.Context, itemIDs []uint64) (*models.MediaItems, error)
}

type collectionService struct {
	repo         repository.MediaItemRepository[*mediatypes.Collection]
	mediaItemSvc MediaItemService[*mediatypes.Collection]
}

// NewCollectionService creates a new collection service
func NewCollectionService(
	repo repository.MediaItemRepository[*mediatypes.Collection],
) CollectionService {
	// Create the base media item service
	mediaItemSvc := NewMediaItemService(repo)

	return &collectionService{
		repo:         repo,
		mediaItemSvc: mediaItemSvc,
	}
}

// Base operations (delegating to MediaItemService)

func (s *collectionService) Create(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error) {
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

	return s.mediaItemSvc.Create(ctx, collection)
}

func (s *collectionService) Update(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error) {
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
	return s.mediaItemSvc.Update(ctx, collection)
}

func (s *collectionService) GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Collection], error) {
	return s.mediaItemSvc.GetByID(ctx, id)
}

func (s *collectionService) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Collection], error) {
	return s.mediaItemSvc.GetByUserID(ctx, userID)
}

func (s *collectionService) Delete(ctx context.Context, id uint64) error {
	return s.mediaItemSvc.Delete(ctx, id)
}

// Collection-specific operations

func (s *collectionService) GetFeaturedCollections(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Collection], error) {
	// Get all collections for the user
	collections, err := s.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured collections: %w", err)
	}

	// Sort collections by created date (newest first)
	sort.Slice(collections, func(i, j int) bool {
		return collections[i].CreatedAt.After(collections[j].CreatedAt)
	})

	// Apply feature criteria (e.g., has cover image, has description, etc.)
	var featured []*models.MediaItem[*mediatypes.Collection]
	for _, collection := range collections {
		// Simple criteria: has items, has cover image, or is recently created
		isFeatured := len(collection.Data.Items) > 0 ||
			collection.Data.Details.Artwork.Poster != "" ||
			time.Since(collection.CreatedAt) < time.Hour*24*7 // One week

		if isFeatured {
			featured = append(featured, collection)
		}

		// Stop once we have enough featured collections
		if len(featured) >= limit {
			break
		}
	}

	// If we don't have enough featured collections, just return the most recent ones
	if len(featured) < limit && len(collections) > 0 {
		remaining := limit - len(featured)
		for _, collection := range collections {
			// Skip collections that are already featured
			if containsCollection(featured, collection.ID) {
				continue
			}

			featured = append(featured, collection)
			remaining--

			if remaining <= 0 {
				break
			}
		}
	}

	return featured, nil
}

func (s *collectionService) GetSmartCollections(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Collection], error) {
	// Smart collections are dynamically generated collections based on criteria
	// For now, we'll return an empty array and mark as not implemented
	// This would be a good place to implement algorithmic collection generation
	return []*models.MediaItem[*mediatypes.Collection]{}, errors.New("smart collections not implemented")
}

func (s *collectionService) GetCollectionItems(ctx context.Context, collectionID uint64) (*models.MediaItems, error) {
	// Get the collection
	collection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection items: %w", err)
	}

	// Return the items
	if collection.Data.Items == nil || len(collection.Data.Items) == 0 {
		return &models.MediaItems{}, nil
	}

	// Extract item IDs from the collection
	itemIDs := collection.Data.GetItemIDs()

	// Get all items
	items, err := s.repo.GetAllMediaItemsByIDs(ctx, itemIDs)

	return items, nil
}

func (s *collectionService) AddItemToCollection(ctx context.Context, collectionID uint64, item models.MediaItem[mediatypes.MediaData]) error {
	// Get the collection
	collection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("failed to add item to collection: %w", err)
	}

	// Check if the item already exists in the collection (using the Collection's optimized method)
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
	return err
}

func (s *collectionService) RemoveItemFromCollection(ctx context.Context, collectionID uint64, itemID uint64) error {
	// Get the collection
	collection, err := s.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("failed to remove item from collection: %w", err)
	}

	// Use the new RemoveItem method
	err = collection.Data.RemoveItem(itemID, collection.Data.ModifiedBy)
	if err != nil {
		return fmt.Errorf("failed to remove item: %w", err)
	}
	_, err = s.Update(ctx, *collection)
	return err
}

func (s *collectionService) UpdateCollectionItems(ctx context.Context, collectionID uint64, items []models.MediaItem[mediatypes.MediaData]) error {
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
	return err
}

func (s *collectionService) SearchCollections(ctx context.Context, query string, userID uint64) ([]*models.MediaItem[*mediatypes.Collection], error) {
	// Get all collections for the user
	collections, err := s.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to search collections: %w", err)
	}

	// Filter collections by query
	var results []*models.MediaItem[*mediatypes.Collection]
	for _, collection := range collections {
		if containsIgnoreCase(collection.Data.Details.Title, query) ||
			containsIgnoreCase(collection.Data.Details.Description, query) {
			results = append(results, collection)
		}
	}

	return results, nil
}

// Additional methods needed by handler

func (s *collectionService) GetCollections(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Collection], error) {
	// Get all collections for the user
	collections, err := s.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	// Sort collections by created date (newest first)
	sort.Slice(collections, func(i, j int) bool {
		return collections[i].CreatedAt.After(collections[j].CreatedAt)
	})

	// Apply limit if needed
	if limit > 0 && len(collections) > limit {
		collections = collections[:limit]
	}

	return collections, nil
}

func (s *collectionService) GetCollectionByID(ctx context.Context, userID uint64, clientID uint64, collectionID uint64) (*models.MediaItem[*mediatypes.Collection], error) {
	// This is a facade that will integrate with the client-specific media service
	// For this implementation, we'll try to use the internal ID if possible

	return s.GetByID(ctx, collectionID)

	// If the collectionID isn't a valid uint64, it might be a client-specific ID
	// In a complete implementation, we would need to check client-specific repositories
	// For now, we'll just return an error
}

func (s *collectionService) GetItems(ctx context.Context, itemIDs []uint64) (*models.MediaItems, error) {

	items, err := s.repo.GetAllMediaItemsByIDs(ctx, itemIDs)
	if err != nil {
		return nil, err
	}
	// For now, we'll return a not implemented error
	return items, errors.New("get items not implemented")
}

// Helper functions

func containsCollection(collections []*models.MediaItem[*mediatypes.Collection], id uint64) bool {
	for _, collection := range collections {
		if collection.ID == id {
			return true
		}
	}
	return false
}

// func parseID(id string) (uint64, error) {
// 	// Try to parse the ID as uint64
// 	return parseInt(id, 10, 64)
// }

// func parseInt(s string, base int, bitSize int) (uint64, error) {
// 	return strings.ParseUint(s, base, bitSize)
// }
