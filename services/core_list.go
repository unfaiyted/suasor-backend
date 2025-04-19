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

// ListService manages application-specific List operations beyond the basic CRUD operations
type CoreListService[T mediatypes.ListData] interface {
	// Base operations (leveraging UserMediaItemService)
	Create(ctx context.Context, list *models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, list *models.MediaItem[T]) (*models.MediaItem[T], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
	GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error)
	GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error)
	Delete(ctx context.Context, id uint64) error

	// list-specific operations
	GetItems(ctx context.Context, listID uint64) (*models.MediaItems, error)
	AddItem(ctx context.Context, listID uint64, itemID uint64) error
	RemoveItem(ctx context.Context, listID uint64, itemID uint64) error
	ReorderItems(ctx context.Context, listID uint64, itemIDs []uint64) error
	UpdateItems(ctx context.Context, listID uint64, items []*models.MediaItem[T]) error
	Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[T], error)
	GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error)
	Sync(ctx context.Context, listID uint64, targetClientIDs []uint64) error
}

type coreListService[T mediatypes.ListData] struct {
	itemRepo repository.MediaItemRepository[T] // For fetching list items
}

// NewlistService creates a new list service
func NewCoreListService[T mediatypes.ListData](
	itemRepo repository.MediaItemRepository[T],
) CoreListService[T] {
	return &coreListService[T]{
		itemRepo: itemRepo,
	}
}

// Base operations (delegating to UserMediaItemService where appropriate)
func (s coreListService[T]) Create(ctx context.Context, list *models.MediaItem[T]) (*models.MediaItem[T], error) {

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("title", list.Title).
		Msg("Creating list")
	itemList := list.GetData().GetItemList()

	// Ensure list has a valid name
	if list.Title == "" || itemList.Details.Title == "" {
		return nil, errors.New("list must have a title")
	}

	// Initialize items array if nil
	if itemList.Items == nil {
		itemList.Items = []mediatypes.ListItem{}
	}

	// Initialize sync client states if nil
	if itemList.SyncClientStates == nil {
		itemList.SyncClientStates = mediatypes.SyncClientStates{}
	}

	// Set creation time for LastModified
	if itemList.LastModified.IsZero() {
		itemList.LastModified = time.Now()
	}

	// Initialize ItemCount
	itemList.ItemCount = len(itemList.Items)

	// Set owner if not set
	if itemList.OwnerID == 0 && list.OwnerID != 0 {
		itemList.OwnerID = list.OwnerID
	}

	// Set title at MediaItem level to match the Data.ItemList.Details.Title
	if list.Title == "" && itemList.Details.Title != "" {
		list.Title = itemList.Details.Title
	}

	// Use the underlying repository directly for better control over validation
	result, err := s.itemRepo.Create(ctx, list)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create list")
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	log.Info().
		Uint64("id", result.ID).
		Str("title", result.Title).
		Msg("list created successfully")

	return result, nil
}
func (s coreListService[T]) Update(ctx context.Context, list *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", list.ID).
		Str("title", list.Title).
		Msg("Updating list")

	// Ensure the list exists
	existing, err := s.GetByID(ctx, list.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update list: %w", err)
	}
	itemList := list.GetData().GetItemList()
	existingItemList := existing.GetData().GetItemList()

	// Preserve items if not provided in the update
	if itemList.Items == nil || len(itemList.Items) == 0 {
		itemList.Items = existingItemList.Items
	}

	// Preserve sync client states if not provided
	if itemList.SyncClientStates == nil || len(itemList.SyncClientStates) == 0 {
		itemList.SyncClientStates = existingItemList.SyncClientStates
	}

	// Update last modified time
	itemList.LastModified = time.Now()

	// Update ItemCount
	itemList.ItemCount = len(itemList.Items)

	// Ensure positions are normalized
	itemList.NormalizePositions()

	// Set title at MediaItem level to match the Data.ItemList.Details.Title
	if list.Title != itemList.Details.Title {
		list.Title = itemList.Details.Title
	}

	// Run validation to check for issues
	issues := itemList.ValidateItems()
	if len(issues) > 0 {
		// Log the issues but continue with the update
		for _, issue := range issues {
			log.Warn().Str("issue", issue).Msg("list validation issue")
		}
	}

	// Update using the user service for consistent behavior
	result, err := s.itemRepo.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", list.ID).
			Msg("Failed to update list")
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	log.Info().
		Uint64("id", result.ID).
		Str("title", result.Title).
		Msg("list updated successfully")

	return result, nil
}
func (s coreListService[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Getting list by ID")

	// Use the user service
	result, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to get list")
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	// Verify this is actually a list
	if !result.IsList() {
		log.Error().
			Uint64("id", id).
			Str("actualType", string(result.Type)).
			Msg("Item is not a list")
		return nil, fmt.Errorf("item with ID %d is not a list", id)
	}

	return result, nil
}
func (s coreListService[T]) GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting lists by user ID")
	// Use the user service
	return s.itemRepo.GetByUserID(ctx, userID, limit, offset)
}

// list-specific operations
func (s coreListService[T]) GetItems(ctx context.Context, listID uint64) (*models.MediaItems, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Msg("Getting list items")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return &models.MediaItems{}, fmt.Errorf("failed to get list items: %w", err)
	}

	itemList := list.GetData().GetItemList()

	// Return empty array if the list has no items
	if len(itemList.Items) == 0 {
		return &models.MediaItems{}, nil
	}

	// Extract item IDs for batch retrieval
	itemIDs := make([]uint64, len(itemList.Items))
	for i, item := range itemList.Items {
		itemIDs[i] = item.ItemID
	}

	// Fetch the actual media items using the core media repository
	actualItems, err := s.itemRepo.GetMixedMediaItemsByIDs(ctx, itemIDs)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to fetch actual list items")
		return nil, fmt.Errorf("failed to get list items: %w", err)
	}

	return actualItems, nil

}
func (s coreListService[T]) AddItem(ctx context.Context, listID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Adding item to list")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to add item to list: %w", err)
	}

	itemList := list.GetData().GetItemList()

	// Create a ListItem from the media item
	newItem := mediatypes.ListItem{
		ItemID:        itemID,
		Position:      len(itemList.Items),
		LastChanged:   time.Now(),
		ChangeHistory: []mediatypes.ChangeRecord{},
	}

	// Add the item using the built-in AddItem method
	// 0 indicates application level modification
	itemList.AddItem(newItem)

	// Store the update
	_, err = s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to add item to list")
		return fmt.Errorf("failed to update list after adding item: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Item added to list successfully")

	return nil
}
func (s coreListService[T]) RemoveItem(ctx context.Context, listID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Removing item from list")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to remove item from list: %w", err)
	}
	itemList := list.GetData().GetItemList()

	// Use the RemoveItem method provided by ItemList
	// 0 indicates application level modification
	err = itemList.RemoveItem(itemID, 0)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to remove item from list")
		return err
	}

	// Store the update
	_, err = s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to update list after removing item")
		return fmt.Errorf("failed to update list after removing item: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Item removed from list successfully")

	return nil
}
func (s coreListService[T]) ReorderItems(ctx context.Context, listID uint64, itemIDs []uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Interface("itemIDs", itemIDs).
		Msg("Reordering list items")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to reorder list items: %w", err)
	}
	itemList := list.GetData().GetItemList()

	// Verify that the number of items matches
	if len(itemIDs) != len(itemList.Items) {
		log.Error().
			Int("providedCount", len(itemIDs)).
			Int("actualCount", len(itemList.Items)).
			Msg("Reorder operation must include all list items")
		return errors.New("reorder operation must include all list items")
	}

	// Create a new ordered list of items
	newOrder := make([]mediatypes.ListItem, len(itemIDs))
	tempItems := make(map[uint64]mediatypes.ListItem)

	// Create a map of existing items for quick lookup
	for _, item := range itemList.Items {
		tempItems[item.ItemID] = item
	}

	// First verify all items exist
	missingItems := []uint64{}
	for _, id := range itemIDs {
		// id, err := strconv.ParseUint(idStr, 10, 64)
		// if err != nil {
		// 	log.Error().
		// 		Str("itemID", idStr).
		// 		Msg("Invalid item ID format")
		// 	return fmt.Errorf("invalid item ID format: %s", idStr)
		// }

		if _, exists := tempItems[id]; !exists {
			missingItems = append(missingItems, id)
		}
	}

	if len(missingItems) > 0 {
		log.Error().
			Interface("missingItems", missingItems).
			Msg("Items not found in list")
		return fmt.Errorf("items not found in list: %v", missingItems)
	}

	// Now build the new order
	for i, id := range itemIDs {
		// id, _ := strconv.ParseUint(id, 10, 64)
		item := tempItems[id]

		// Update position
		item.Position = i

		// Add change record
		item.AddChangeRecord(0, "reorder") // 0 indicates application level change

		newOrder[i] = item
	}

	// Update the list with the new item order
	itemList.Items = newOrder
	itemList.LastModified = time.Now()
	itemList.ModifiedBy = 0 // 0 indicates application level modification

	// Normalize positions to ensure they're sequential
	itemList.NormalizePositions()

	_, err = s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to update list after reordering items")
		return fmt.Errorf("failed to update list after reordering: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Msg("list items reordered successfully")

	return nil
}
func (s coreListService[T]) UpdateItems(ctx context.Context, listID uint64, items []*models.MediaItem[T]) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Int("itemCount", len(items)).
		Msg("Updating list items")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to update list items: %w", err)
	}

	// Convert MediaItems to ListItems
	listItems := make([]mediatypes.ListItem, len(items))
	now := time.Now()

	for i, item := range items {
		listItems[i] = mediatypes.ListItem{
			ItemID:      item.ID,
			Position:    i,
			LastChanged: now,
			ChangeHistory: []mediatypes.ChangeRecord{
				{
					ClientID:   0, // 0 indicates application level change
					ItemID:     fmt.Sprintf("%d", item.ID),
					ChangeType: "update",
					Timestamp:  now,
				},
			},
		}
	}

	itemList := list.GetData().GetItemList()
	// Replace all items
	itemList.Items = listItems
	itemList.LastModified = now
	itemList.ModifiedBy = 0 // 0 indicates application level modification
	itemList.ItemCount = len(listItems)

	// Ensure positions are normalized
	itemList.NormalizePositions()

	// Update the list
	_, err = s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to update list items")
		return fmt.Errorf("failed to update list items: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Int("itemCount", len(items)).
		Msg("list items updated successfully")

	return nil
}
func (s coreListService[T]) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Uint64("userID", query.OwnerID).
		Msg("Searching lists")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	query.MediaType = mediaType

	// Delegate to the user service
	results, err := s.itemRepo.Search(ctx, query)
	if err != nil {
		log.Error().Err(err).
			Str("query", query.Query).
			Uint64("userID", query.OwnerID).
			Msg("Failed to search lists")
		return nil, fmt.Errorf("failed to search lists: %w", err)
	}

	log.Info().
		Str("query", query.Query).
		Uint64("userID", query.OwnerID).
		Int("count", len(results)).
		Msg("lists found")

	return results, nil
}
func (s coreListService[T]) GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Int("limit", limit).
		Int("days", days).
		Msg("Getting recent lists")

	// Delegate to the user service
	results, err := s.itemRepo.GetRecentItems(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Msg("Failed to get recent lists")
		return nil, fmt.Errorf("failed to get recent lists: %w", err)
	}

	log.Info().
		Int("count", len(results)).
		Msg("Recent lists retrieved")

	return results, nil
}
func (s coreListService[T]) Sync(ctx context.Context, listID uint64, targetClientIDs []uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Interface("targetClientIDs", targetClientIDs).
		Msg("Syncing list")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to sync list: %w", err)
	}

	// In a real implementation, this would use the list sync job to:
	// 1. For each target client ID, create or update a list with the same name
	// 2. Map the item IDs in itemList.Items to client-specific IDs using the mediaItemRepo
	// 3. Add these items to the client-specific list
	// 4. Store the client-specific list IDs and item IDs in the SyncClientStates
	// 5. Update the LastSynced timestamp
	itemList := list.GetData().GetItemList()
	// For now, just create a placeholder in the SyncClientStates to show intent
	now := time.Now()
	for _, clientID := range targetClientIDs {

		// Create a placeholder sync client state
		syncState := itemList.SyncClientStates.GetSyncClientState(clientID)
		if syncState == nil {
			// Create empty sync list items
			syncItems := mediatypes.SyncListItems{}

			// Add a new state for this client
			newState := mediatypes.SyncClientState{
				ClientID:     clientID,
				Items:        syncItems,
				ClientListID: "", // Empty for now, would be the client-specific list ID
				LastSynced:   now,
			}

			itemList.SyncClientStates = append(itemList.SyncClientStates, newState)
		}
	}

	// Update the LastSynced timestamp
	itemList.LastSynced = now

	// Save the updated list
	_, err = s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to update list sync state")
		return fmt.Errorf("failed to update list sync state: %w", err)
	}

	log.Warn().Msg("list sync partially implemented - requires job implementation")

	// This would be implemented in a job that handles the actual sync
	return errors.New("list sync requires implementation in the list sync job")
}
func (s coreListService[T]) createMediaItems(ctx context.Context, items []*models.MediaItem[mediatypes.MediaData]) (*models.MediaItems, error) {
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Int("count", len(items)).
		Msg("Creating media items")

	return nil, nil
}

// Delete List
func (s coreListService[T]) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting list")

	err := s.itemRepo.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to delete list")
		return fmt.Errorf("failed to delete list: %w", err)
	}

	log.Info().
		Uint64("id", id).
		Msg("List deleted successfully")

	return nil
}

func (s *coreListService[T]) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting all lists")

	// Delegate to repository
	results, err := s.itemRepo.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Int("limit", limit).
			Int("offset", offset).
			Msg("Failed to get all lists")
		return nil, fmt.Errorf("failed to get all lists: %w", err)
	}

	log.Info().
		Int("count", len(results)).
		Msg("All lists retrieved successfully")

	return results, nil
}
