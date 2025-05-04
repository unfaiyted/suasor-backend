// list_sync_helper.go
package sync

import (
	"context"
	"fmt"
	"time"

	"suasor/clients/media"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// ListSyncHelper provides utility functions for syncing lists (playlists and collections)
// between different media clients using the generic ListProvider interface
type ListSyncHelper struct {
	mediaJob *MediaSyncJob
}

// NewListSyncHelper creates a new helper for list syncing
func NewListSyncHelper(mediaJob *MediaSyncJob) *ListSyncHelper {
	return &ListSyncHelper{
		mediaJob: mediaJob,
	}
}

// syncListItems syncs items from one client to another using the generic ListProvider interface
func (j *MediaSyncJob) syncListItems(
	ctx context.Context,
	sourceClient media.ClientMedia,
	targetClient media.ClientMedia,
	sourceListID string,
	targetListID string,
	isCopyOperation bool,
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Str("sourceListID", sourceListID).
		Str("targetListID", targetListID).
		Bool("isCopyOperation", isCopyOperation).
		Msg("Syncing list items between clients")

	// Check if clients support lists
	sourceProvider, ok := sourceClient.(providers.ListProvider[mediatypes.ListData])
	if !ok {
		return fmt.Errorf("source client doesn't support lists")
	}

	targetProvider, ok := targetClient.(providers.ListProvider[mediatypes.ListData])
	if !ok {
		return fmt.Errorf("target client doesn't support lists")
	}

	// First, get the list from source client to determine its type
	sourceLists, err := sourceProvider.SearchLists(ctx, &mediatypes.QueryOptions{
		ExternalSourceID: sourceListID,
	})
	if err != nil || len(sourceLists) == 0 {
		return fmt.Errorf("could not find source list: %w", err)
	}
	sourceList := sourceLists[0]

	// Determine if this is a playlist or collection
	isPlaylist := sourceList.Type == mediatypes.MediaTypePlaylist
	log.Debug().
		Bool("isPlaylist", isPlaylist).
		Str("mediaType", string(sourceList.Type)).
		Msg("Determined list type")

	// Get detailed items for the list from source client
	listItems, err := sourceProvider.GetListItems(ctx, sourceListID, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get source list items: %w", err)
	}

	// Get or create target list if needed
	var targetList *models.MediaItem[mediatypes.ListData]
	if targetListID == "" || isCopyOperation {
		// Create a new list on target
		targetList, err = targetProvider.CreateList(ctx, sourceList.Title, sourceList.Data.GetDetails().Description)
		if err != nil {
			return fmt.Errorf("failed to create list on target: %w", err)
		}
		targetListID = getClientItemID(targetList, targetClient.GetClientID())
		log.Info().
			Str("listTitle", sourceList.Title).
			Str("targetListID", targetListID).
			Msg("Created new list on target client")
	} else {
		// Get the existing list on target
		targetLists, err := targetProvider.SearchLists(ctx, &mediatypes.QueryOptions{
			ExternalSourceID: targetListID,
		})
		if err != nil || len(targetLists) == 0 {
			return fmt.Errorf("could not find target list: %w", err)
		}
		targetList = targetLists[0]
		log.Debug().
			Str("targetListID", targetListID).
			Str("listTitle", targetList.Title).
			Msg("Found existing list on target client")
	}

	// Map items from source to target
	for _, item := range listItems {
		// Get corresponding item ID in target client
		sourceItemID := getClientItemID(item, sourceClient.GetClientID())
		if sourceItemID == "" {
			log.Printf("No item ID found in source client for item %s", item.Title)
			continue
		}

		// Find matching item in target client
		targetItemID, err := j.findMatchingItemInTargetClient(ctx, sourceClient.GetClientID(), sourceItemID, targetClient.GetClientID())
		if err != nil {
			log.Printf("Could not find matching item in target client: %v", err)
			continue
		}

		// Add item to target list
		if err := targetProvider.AddItemList(ctx, targetListID, targetItemID); err != nil {
			log.Printf("Failed to add item to target list: %v", err)
			continue
		}
	}

	return nil
}

// getClientItemID gets the client-specific ID for a media item
func getClientItemID(item *models.MediaItem[mediatypes.ListData], clientID uint64) string {
	for _, clientInfo := range item.SyncClients {
		if clientInfo.ID == clientID {
			return clientInfo.ItemID
		}
	}
	return ""
}

// findMatchingItemInTargetClient finds the corresponding item ID in the target client
func (j *MediaSyncJob) findMatchingItemInTargetClient(ctx context.Context, sourceClientID uint64, sourceItemID string, targetClientID uint64) (string, error) {
	// This is a placeholder - in a real implementation, you would:
	// 1. Look up the item by source client ID and item ID in your database
	// 2. Find if there's a matching entry for the target client
	// 3. Return the target client's item ID

	// Here's a simplified implementation that assumes direct 1:1 mapping
	// In a real system, you'd use your repository to look this up

	// Example implementation:
	// sourceItem, err := j.clientItemRepos.GetByClientItemID(ctx, sourceClientID, sourceItemID)
	// if err != nil {
	//    return "", fmt.Errorf("source item not found: %w", err)
	// }
	//
	// for _, clientID := range sourceItem.SyncClients {
	//    if clientID.ID == targetClientID {
	//        return clientID.ItemID, nil
	//    }
	// }

	return sourceItemID, nil // Simplified assumption that IDs match across clients
}

// Function to help determine if a client supports a particular list type
func clientSupportsListType(client media.ClientMedia, listType mediatypes.MediaType) bool {
	if listType == mediatypes.MediaTypePlaylist {
		if provider, ok := client.(providers.PlaylistProvider); ok {
			return provider.SupportsPlaylists()
		}
	} else if listType == mediatypes.MediaTypeCollection {
		if provider, ok := client.(providers.CollectionProvider); ok {
			return provider.SupportsCollections()
		}
	}
	return false
}

// syncLists syncs playlists or collections between clients
func (j *MediaSyncJob) syncLists(
	ctx context.Context,
	userID uint64,
	sourceClientID uint64,
	targetClientID uint64,
	listType mediatypes.MediaType,
) error {
	// Get the source client
	sourceClient, err := j.getMediaClient(ctx, userID, sourceClientID)
	if err != nil {
		return fmt.Errorf("failed to get source client: %w", err)
	}

	// Get the target client
	targetClient, err := j.getMediaClient(ctx, userID, targetClientID)
	if err != nil {
		return fmt.Errorf("failed to get target client: %w", err)
	}

	// Check if clients support the specified list type
	if !clientSupportsListType(sourceClient, listType) {
		return fmt.Errorf("source client doesn't support %s", listType)
	}

	if !clientSupportsListType(targetClient, listType) {
		return fmt.Errorf("target client doesn't support %s", listType)
	}

	// Cast to ListProvider
	sourceProvider, _ := sourceClient.(providers.ListProvider[mediatypes.ListData])
	targetProvider, _ := targetClient.(providers.ListProvider[mediatypes.ListData])

	// Get all lists from source client
	sourceLists, err := sourceProvider.SearchLists(ctx, &mediatypes.QueryOptions{
		MediaType: listType,
	})
	if err != nil {
		return fmt.Errorf("failed to get lists from source client: %w", err)
	}

	// Get all lists from target client
	targetLists, err := targetProvider.SearchLists(ctx, &mediatypes.QueryOptions{
		MediaType: listType,
	})
	if err != nil {
		return fmt.Errorf("failed to get lists from target client: %w", err)
	}

	// Create a map of target lists by title for easy lookup
	targetListsByTitle := make(map[string]*models.MediaItem[mediatypes.ListData])
	for _, list := range targetLists {
		targetListsByTitle[list.Title] = list
	}

	// Sync each source list to target
	for _, sourceList := range sourceLists {
		// Check if list exists in target by title
		targetList, exists := targetListsByTitle[sourceList.Title]

		// Get source list ID
		sourceListID := getClientItemID(sourceList, sourceClientID)

		if exists {
			// List exists on target, update it
			targetListID := getClientItemID(targetList, targetClientID)

			// Sync list items
			if err := j.syncListItems(ctx, sourceClient, targetClient, sourceListID, targetListID, false); err != nil {
				log.Printf("Error syncing list items for %s: %v", sourceList.Title, err)
				continue
			}
		} else {
			// List doesn't exist on target, create it
			if err := j.syncListItems(ctx, sourceClient, targetClient, sourceListID, "", true); err != nil {
				log.Printf("Error creating and syncing list for %s: %v", sourceList.Title, err)
				continue
			}
		}
	}

	return nil
}

// getMediaClient is a helper function to get a media client by ID
func (j *MediaSyncJob) getMediaClient(ctx context.Context, userID uint64, clientID uint64) (media.ClientMedia, error) {
	// This is a placeholder - in a real implementation, you would:
	// 1. Get the client configuration from your repositories
	// 2. Create or retrieve a client instance
	// 3. Return the client

	// Example implementation using your existing methods:
	client, _, err := j.getClientMedia(ctx, clientID)
	return client, err
}

// syncListItems is a method on ListSyncHelper for syncing items between lists
func (h *ListSyncHelper) syncListItems(
	ctx context.Context,
	sourceClient media.ClientMedia,
	targetClient media.ClientMedia,
	sourceListID string,
	targetListID string,
	sourceList *models.MediaItem[mediatypes.ListData],
	targetList *models.MediaItem[mediatypes.ListData],
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Str("sourceListID", sourceListID).
		Str("targetListID", targetListID).
		Msg("Syncing list items")

	// Get source items
	sourceProvider, _ := sourceClient.(providers.ListProvider[mediatypes.ListData])
	targetProvider, _ := targetClient.(providers.ListProvider[mediatypes.ListData])

	sourceItems, err := sourceProvider.GetListItems(ctx, sourceListID, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get items from source list: %w", err)
	}

	// Get existing target items to avoid duplicates
	targetItems, err := targetProvider.GetListItems(ctx, targetListID, &mediatypes.QueryOptions{})
	if err != nil {
		log.Warn().
			Err(err).
			Str("targetListID", targetListID).
			Msg("Failed to get items from target list, continuing with sync")
		// Continue with empty slice rather than aborting
		targetItems = []*models.MediaItem[mediatypes.ListData]{}
	}

	// Create a map of existing target items by title for quick lookups
	existingTargetTitles := make(map[string]bool)
	for _, item := range targetItems {
		existingTargetTitles[item.Title] = true
	}

	// Process each source item
	addedItems := 0
	for _, sourceItem := range sourceItems {
		// Skip if this item already exists in the target by title
		if existingTargetTitles[sourceItem.Title] {
			log.Debug().
				Str("itemTitle", sourceItem.Title).
				Msg("Item already exists in target list, skipping")
			continue
		}

		// Find source client-specific ID
		var sourceItemID string
		for _, clientID := range sourceItem.SyncClients {
			if clientID.ID == sourceClient.GetClientID() {
				sourceItemID = clientID.ItemID
				break
			}
		}

		if sourceItemID == "" {
			log.Warn().
				Str("itemTitle", sourceItem.Title).
				Msg("Could not find source item ID, skipping")
			continue
		}

		// Find matching item ID in target client using UUID if available
		var targetItemID string

		// First try to get the target ID from existing mappings
		for _, clientID := range sourceItem.SyncClients {
			if clientID.ID == targetClient.GetClientID() {
				targetItemID = clientID.ItemID
				break
			}
		}

		// If no direct mapping is found, try to find by UUID using the media sync job's helper
		if targetItemID == "" && h.mediaJob != nil {
			mappedID, err := h.mediaJob.findMatchingItemInTargetClient(ctx,
				sourceClient.GetClientID(),
				sourceItemID,
				targetClient.GetClientID())
			if err == nil {
				targetItemID = mappedID
			}
		}

		if targetItemID == "" {
			log.Warn().
				Str("itemTitle", sourceItem.Title).
				Msg("Could not find matching item in target client, skipping")
			continue
		}

		// Add to the target list
		err := targetProvider.AddItemList(ctx, targetListID, targetItemID)
		if err != nil {
			log.Error().
				Err(err).
				Str("targetListID", targetListID).
				Str("targetItemID", targetItemID).
				Msg("Failed to add item to target list")
			continue
		}

		addedItems++
		log.Debug().
			Str("itemTitle", sourceItem.Title).
			Str("targetItemID", targetItemID).
			Msg("Added item to target list")
	}

	// Update list sync states for tracking
	updateListSyncState(sourceList, sourceClient.GetClientID(), sourceListID)
	updateListSyncState(targetList, targetClient.GetClientID(), targetListID)

	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Int("addedItems", addedItems).
		Msg("Finished syncing list items")

	return nil
}

// updateListSyncState updates the sync state for a list
func updateListSyncState(list *models.MediaItem[mediatypes.ListData], clientID uint64, listID string) {
	if list.Data.ItemList.SyncStates == nil {
		list.Data.ItemList.SyncStates = mediatypes.ListSyncStates{}
	}

	// Check if we already have a sync state for this client
	var syncState *mediatypes.ListSyncState
	for i := range list.Data.ItemList.SyncStates {
		if list.Data.ItemList.SyncStates[i].ClientID == clientID {
			syncState = &list.Data.ItemList.SyncStates[i]
			break
		}
	}

	now := time.Now()

	// If no sync state exists, create a new one
	if syncState == nil {
		list.Data.ItemList.SyncStates = append(list.Data.ItemList.SyncStates, mediatypes.ListSyncState{
			ClientID:     clientID,
			ClientListID: listID,
			LastSynced:   now,
			Items:        mediatypes.ClientListItems{},
		})
	} else {
		// Update existing sync state
		syncState.LastSynced = now
		syncState.ClientListID = listID
	}

	// Update the overall last synced time
	list.Data.ItemList.LastSynced = now
}

// SyncOptions defines options for list syncing
type SyncOptions struct {
	MediaTypes    []mediatypes.MediaType // Filter by media type
	IncludeTitles []string               // Only sync lists with these titles
	ExcludeTitles []string               // Skip lists with these titles
}

// shouldSyncList checks if a list should be synced based on sync options
func shouldSyncList(list *models.MediaItem[mediatypes.ListData], options *SyncOptions) bool {
	if options == nil {
		return true // Default to sync all if no options provided
	}

	// Check media type filter if specified
	if options.MediaTypes != nil && len(options.MediaTypes) > 0 {
		mediaType := list.Data.GetMediaType()

		// Check if the list's media type is in the allowed types
		found := false
		for _, allowedType := range options.MediaTypes {
			if mediaType == allowedType || allowedType == mediatypes.MediaTypeAll {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	// Check specific list titles to include if specified
	if options.IncludeTitles != nil && len(options.IncludeTitles) > 0 {
		found := false
		for _, title := range options.IncludeTitles {
			if title == list.Data.GetTitle() {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	// Check list titles to exclude
	if options.ExcludeTitles != nil && len(options.ExcludeTitles) > 0 {
		for _, title := range options.ExcludeTitles {
			if title == list.Data.GetTitle() {
				return false
			}
		}
	}

	return true
}

// SyncLists syncs all lists (both playlists and collections) between two clients
func (h *ListSyncHelper) SyncLists(
	ctx context.Context,
	sourceClient media.ClientMedia,
	targetClient media.ClientMedia,
	syncOptions *SyncOptions,
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Msg("Syncing lists between clients")

	// Check if clients support lists
	sourceProvider, ok := sourceClient.(providers.ListProvider[mediatypes.ListData])
	if !ok {
		return fmt.Errorf("source client doesn't support lists")
	}

	targetProvider, ok := targetClient.(providers.ListProvider[mediatypes.ListData])
	if !ok {
		return fmt.Errorf("target client doesn't support lists")
	}

	// Get all lists from source client
	sourceLists, err := sourceProvider.SearchLists(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("sourceClientID", sourceClient.GetClientID()).
			Msg("Failed to get lists from source client")
		return fmt.Errorf("failed to get lists from source client: %w", err)
	}

	if len(sourceLists) == 0 {
		log.Info().
			Uint64("sourceClientID", sourceClient.GetClientID()).
			Msg("No lists found in source client")
		return nil
	}

	// Get all lists from target client
	targetLists, err := targetProvider.SearchLists(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("targetClientID", targetClient.GetClientID()).
			Msg("Failed to get lists from target client")
		return fmt.Errorf("failed to get lists from target client: %w", err)
	}

	// Create a map of target lists by title for easier lookup
	targetListsByTitle := make(map[string]*models.MediaItem[mediatypes.ListData])
	for i, list := range targetLists {
		targetListsByTitle[list.Data.GetTitle()] = targetLists[i]
	}

	// Sync each source list to target
	syncedCount := 0
	for _, sourceList := range sourceLists {
		// Check if this list should be synced based on filters
		if !shouldSyncList(sourceList, syncOptions) {
			log.Info().
				Uint64("sourceClientID", sourceClient.GetClientID()).
				Str("listTitle", sourceList.Data.GetTitle()).
				Msg("Skipping list based on sync options")
			continue
		}

		// Find matching list in target or create it
		targetList, exists := targetListsByTitle[sourceList.Data.GetTitle()]
		if !exists {
			// Create a new list in the target client
			newList, err := targetProvider.CreateList(ctx,
				sourceList.Data.GetTitle(),
				sourceList.Data.ItemList.Details.Description)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("targetClientID", targetClient.GetClientID()).
					Str("listTitle", sourceList.Data.GetTitle()).
					Msg("Failed to create list in target client")
				continue
			}
			targetList = newList
			log.Info().
				Uint64("targetClientID", targetClient.GetClientID()).
				Str("listTitle", targetList.Data.GetTitle()).
				Msg("Created new list in target client")
		}

		// Get source list ID
		var sourceListID string
		for _, clientID := range sourceList.SyncClients {
			if clientID.ID == sourceClient.GetClientID() {
				sourceListID = clientID.ItemID
				break
			}
		}

		// Get target list ID
		var targetListID string
		for _, clientID := range targetList.SyncClients {
			if clientID.ID == targetClient.GetClientID() {
				targetListID = clientID.ItemID
				break
			}
		}

		// Sync the list items
		err = h.syncListItems(ctx,
			sourceClient,
			targetClient,
			sourceListID,
			targetListID,
			sourceList,
			targetList)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("sourceClientID", sourceClient.GetClientID()).
				Uint64("targetClientID", targetClient.GetClientID()).
				Str("sourceListID", sourceListID).
				Str("targetListID", targetListID).
				Msg("Failed to sync list items")
			continue
		}

		syncedCount++
		log.Info().
			Uint64("sourceClientID", sourceClient.GetClientID()).
			Uint64("targetClientID", targetClient.GetClientID()).
			Str("listTitle", sourceList.Data.GetTitle()).
			Msg("Successfully synced list")
	}

	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Int("syncedCount", syncedCount).
		Int("totalSourceLists", len(sourceLists)).
		Msg("Completed list synchronization")

	return nil
}

