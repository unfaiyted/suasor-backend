// list_provider_helper.go
package providers

import (
	"context"
	"fmt"

	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
)

// ListProviderFactory provides helper methods to create list providers
type ListProviderFactory struct{}

// NewListProviderFactory creates a new factory
func NewListProviderFactory() *ListProviderFactory {
	return &ListProviderFactory{}
}

// CreatePlaylistListProvider creates a ListProvider for playlists
func (f *ListProviderFactory) CreatePlaylistListProvider(provider PlaylistProvider) ListProvider[*mediatypes.Playlist] {
	return NewPlaylistListAdapter(provider)
}

// CreateCollectionListProvider creates a ListProvider for collections
func (f *ListProviderFactory) CreateCollectionListProvider(provider CollectionProvider) ListProvider[*mediatypes.Collection] {
	return NewCollectionListAdapter(provider)
}

// ListSyncAdapter is a generic adapter for syncing lists between providers
type ListSyncAdapter[T mediatypes.ListData] struct {
	sourceProvider ListProvider[T]
	sourceClientID uint64
	targetProvider ListProvider[T]
	targetClientID uint64
}

// NewListSyncAdapter creates a new adapter for syncing lists
func NewListSyncAdapter[T mediatypes.ListData](source, target ListProvider[T]) *ListSyncAdapter[T] {
	return &ListSyncAdapter[T]{
		sourceProvider: source,
		targetProvider: target,
	}
}

// SyncLists syncs lists between providers
func (a *ListSyncAdapter[T]) SyncLists(ctx context.Context, options *mediatypes.QueryOptions) error {
	// Get all lists from source provider
	sourceLists, err := a.sourceProvider.SearchLists(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to get lists from source: %w", err)
	}

	// Get all lists from target provider for comparison
	targetLists, err := a.targetProvider.SearchLists(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to get lists from target: %w", err)
	}

	// Create a map for quick lookup
	targetListsByTitle := make(map[string]*models.MediaItem[T])
	for i, list := range targetLists {
		targetListsByTitle[list.Title] = targetLists[i]
	}

	// Process each source list
	for _, sourceList := range sourceLists {
		// Check if list exists in target
		targetList, exists := targetListsByTitle[sourceList.Title]

		var targetListID string
		if exists {
			// Get the target list ID
			for _, id := range targetList.SyncClients {
				targetListID = id.ItemID
				break
			}
		} else {
			// Create a new list on the target
			newList, err := a.targetProvider.CreateList(ctx, sourceList.Title, sourceList.Data.GetDetails().Description)
			if err != nil {
				return fmt.Errorf("failed to create list on target: %w", err)
			}
			targetList = newList

			// Get the new list ID
			for _, id := range targetList.SyncClients {
				targetListID = id.ItemID
				break
			}
		}

		// Get source list ID
		var sourceListID string
		for _, id := range sourceList.SyncClients {
			sourceListID = id.ItemID
			break
		}

		// Sync items from source to target list
		err = a.SyncListItems(ctx, sourceListID, targetListID)
		if err != nil {
			return fmt.Errorf("failed to sync list items: %w", err)
		}
	}

	return nil
}

// SyncListItems syncs items between two lists
func (a *ListSyncAdapter[T]) SyncListItems(
	ctx context.Context,
	sourceListID string,
	targetListID string,
) error {
	// Get source list items

	// listType := mediatypes.GetListType[T]()
	sourceItems, err := a.sourceProvider.GetListItems(ctx, sourceListID)
	if err != nil {
		return fmt.Errorf("failed to get source list items: %w", err)
	}
	// Get target list items for comparison
	targetItems, err := a.targetProvider.GetListItems(ctx, targetListID)
	if err != nil {
		// If target is empty, just continue (it might be a new list)
		targetItems = models.NewMediaItemList[T](targetItems.ListOriginID, targetItems.OwnerID)
	}

	// Delete all liste items in target
	err = a.targetProvider.RemoveAllListItems(ctx, targetListID)
	if err != nil {
		return fmt.Errorf("failed to delete target list items: %w", err)
	}

	// Add each source item to target
	sourceItems.ForEach(func(UUID string, mediaType mediatypes.MediaType, item any) bool {
		typedItem, ok := item.(*models.MediaItem[T])
		if !ok {
			return true
		}
		// Add to target

		targetItemID := typedItem.SyncClients.GetClientItemID(a.targetClientID)
		err = a.targetProvider.AddListItem(ctx, targetListID, targetItemID)
		if err != nil {
			return false
		}
		return true
	})

	return nil
}

// ListMappingCollector collects list mappings across providers
type ListMappingCollector[T mediatypes.ListData] struct {
	providers []ListProvider[T]
}

// NewListMappingCollector creates a new collector
func NewListMappingCollector[T mediatypes.ListData](providers []ListProvider[T]) *ListMappingCollector[T] {
	return &ListMappingCollector[T]{
		providers: providers,
	}
}

// GetMappings returns a map of list titles to their IDs across providers
func (c *ListMappingCollector[T]) GetMappings(ctx context.Context) (map[string]map[int]string, error) {
	// Map of list titles to provider index to list ID
	mappings := make(map[string]map[int]string)

	// Get lists from each provider
	for i, provider := range c.providers {
		lists, err := provider.SearchLists(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get lists from provider %d: %w", i, err)
		}

		// Process each list
		for _, list := range lists {
			title := list.Title

			// Get the list ID
			var listID string
			for _, id := range list.SyncClients {
				listID = id.ItemID
				break
			}

			// Initialize the inner map if needed
			if mappings[title] == nil {
				mappings[title] = make(map[int]string)
			}

			// Store the mapping
			mappings[title][i] = listID
		}
	}

	return mappings, nil
}
