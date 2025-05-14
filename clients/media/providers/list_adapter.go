// list_adapter.go
package providers

import (
	"context"
	"fmt"

	"suasor/types/models"

	mediatypes "suasor/clients/media/types"
)

// PlaylistListAdapter adapts a PlaylistProvider to a ListProvider[*types.Playlist]
type PlaylistListAdapter struct {
	provider PlaylistProvider
}

// NewPlaylistListAdapter creates a new adapter
func NewPlaylistListAdapter(provider PlaylistProvider) ListProvider[*mediatypes.Playlist] {
	return &PlaylistListAdapter{provider: provider}
}

// Implementation of ListProvider[*types.Playlist] interface methods

func (a *PlaylistListAdapter) GetListItems(ctx context.Context, listID string) (*models.MediaItemList[*mediatypes.Playlist], error) {
	return a.provider.GetPlaylistItems(ctx, listID)
}

func (a *PlaylistListAdapter) CreateListWithItems(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*mediatypes.Playlist], error) {
	return a.provider.CreatePlaylistWithItems(ctx, name, description, itemIDs)
}

func (a *PlaylistListAdapter) GetList(ctx context.Context, listID string) (*models.MediaItem[*mediatypes.Playlist], error) {
	return a.provider.GetPlaylist(ctx, listID)
}

func (a *PlaylistListAdapter) CreateList(ctx context.Context, name string, description string) (*models.MediaItem[*mediatypes.Playlist], error) {
	return a.provider.CreatePlaylist(ctx, name, description)
}

func (a *PlaylistListAdapter) UpdateList(ctx context.Context, listID string, name string, description string) (*models.MediaItem[*mediatypes.Playlist], error) {
	return a.provider.UpdatePlaylist(ctx, listID, name, description)
}

func (a *PlaylistListAdapter) DeleteList(ctx context.Context, listID string) error {
	return a.provider.DeletePlaylist(ctx, listID)
}

func (a *PlaylistListAdapter) AddListItem(ctx context.Context, listID string, itemID string) error {
	return a.provider.AddPlaylistItem(ctx, listID, itemID)
}

func (a *PlaylistListAdapter) AddListItems(ctx context.Context, listID string, itemIDs []string) error {
	return a.provider.AddPlaylistItems(ctx, listID, itemIDs)
}

func (a *PlaylistListAdapter) RemoveListItem(ctx context.Context, listID string, itemID string) error {
	return a.provider.RemovePlaylistItem(ctx, listID, itemID)
}

func (a *PlaylistListAdapter) RemoveListItems(ctx context.Context, listID string, itemIDs []string) error {
	return a.provider.RemovePlaylistItems(ctx, listID, itemIDs)
}

func (a *PlaylistListAdapter) RemoveAllListItems(ctx context.Context, listID string) error {
	return a.provider.RemoveAllPlaylistItems(ctx, listID)
}

func (a *PlaylistListAdapter) ReorderListItems(ctx context.Context, listID string, itemIDs []string) error {
	return a.provider.ReorderPlaylistItems(ctx, listID, itemIDs)
}

func (a *PlaylistListAdapter) SearchLists(ctx context.Context, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	return a.provider.SearchPlaylists(ctx, options)
}
func (a *PlaylistListAdapter) SearchListItems(ctx context.Context, listID string, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	return a.provider.SearchPlaylistItems(ctx, listID, options)
}

func (a *PlaylistListAdapter) SupportsLists() bool {
	return a.provider.SupportsPlaylists()
}

// CollectionListAdapter adapts a CollectionProvider to a ListProvider[*types.Collection]
type CollectionListAdapter struct {
	provider CollectionProvider
}

// NewCollectionListAdapter creates a new adapter
func NewCollectionListAdapter(provider CollectionProvider) ListProvider[*mediatypes.Collection] {
	return &CollectionListAdapter{provider: provider}
}

// Implementation of ListProvider[*types.Collection] interface methods

func (a *CollectionListAdapter) GetListItems(ctx context.Context, listID string) (*models.MediaItemList[*mediatypes.Collection], error) {
	return a.provider.GetCollectionItems(ctx, listID)
}

func (a *CollectionListAdapter) CreateList(ctx context.Context, name string, description string) (*models.MediaItem[*mediatypes.Collection], error) {
	return a.provider.CreateCollection(ctx, name, description)
}

func (a *CollectionListAdapter) GetList(ctx context.Context, listID string) (*models.MediaItem[*mediatypes.Collection], error) {
	return a.provider.GetCollection(ctx, listID)
}

func (a *CollectionListAdapter) CreateListWithItems(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*mediatypes.Collection], error) {
	return a.provider.CreateCollectionWithItems(ctx, name, description, itemIDs)
}

func (a *CollectionListAdapter) UpdateList(ctx context.Context, listID string, name string, description string) (*models.MediaItem[*mediatypes.Collection], error) {
	return a.provider.UpdateCollection(ctx, listID, name, description)
}

func (a *CollectionListAdapter) DeleteList(ctx context.Context, listID string) error {
	return a.provider.DeleteCollection(ctx, listID)
}

func (a *CollectionListAdapter) AddListItem(ctx context.Context, listID string, itemID string) error {
	return a.provider.AddCollectionItem(ctx, listID, itemID)
}

func (a *CollectionListAdapter) AddListItems(ctx context.Context, listID string, itemIDs []string) error {
	return a.provider.AddCollectionItems(ctx, listID, itemIDs)
}

func (a *CollectionListAdapter) RemoveListItem(ctx context.Context, listID string, itemID string) error {
	return a.provider.RemoveCollectionItem(ctx, listID, itemID)
}

func (a *CollectionListAdapter) RemoveListItems(ctx context.Context, listID string, itemIDs []string) error {
	return a.provider.RemoveCollectionItems(ctx, listID, itemIDs)
}

func (a *CollectionListAdapter) RemoveAllListItems(ctx context.Context, listID string) error {
	return a.provider.RemoveAllCollectionItems(ctx, listID)
}

func (a *CollectionListAdapter) ReorderListItems(ctx context.Context, listID string, itemIDs []string) error {
	return a.provider.ReorderCollectionItems(ctx, listID, itemIDs)
}

func (a *CollectionListAdapter) SearchLists(ctx context.Context, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
	return a.provider.SearchCollections(ctx, options)
}

func (a *CollectionListAdapter) SearchListItems(ctx context.Context, listID string, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
	return a.provider.SearchCollectionItems(ctx, listID, options)
}

func (a *CollectionListAdapter) SupportsLists() bool {
	return a.provider.SupportsCollections()
}

// DynamicListAdapter is a type-switching adapter that implements ListProvider[mediatypes.ListData]
// This adapter determines the type of list at runtime and dispatches to the appropriate concrete provider
type DynamicListAdapter struct {
	playlistAdapter   ListProvider[*mediatypes.Playlist]
	collectionAdapter ListProvider[*mediatypes.Collection]
}

// NewDynamicListAdapter creates a new adapter that can dynamically handle both playlists and collections
func NewDynamicListAdapter(
	playlistAdapter ListProvider[*mediatypes.Playlist],
	collectionAdapter ListProvider[*mediatypes.Collection],
) *DynamicListAdapter {
	return &DynamicListAdapter{
		playlistAdapter:   playlistAdapter,
		collectionAdapter: collectionAdapter,
	}
}

// determineListType tries to identify if a list is a playlist or collection
func (a *DynamicListAdapter) determineListType(ctx context.Context, listID string) (mediatypes.MediaType, error) {
	// Try playlist first
	if a.playlistAdapter != nil {
		playlists, err := a.playlistAdapter.SearchLists(ctx, &mediatypes.QueryOptions{
			ExternalSourceID: listID,
		})
		if err == nil && len(playlists) > 0 {
			return mediatypes.MediaTypePlaylist, nil
		}
	}

	// Then try collection
	if a.collectionAdapter != nil {
		collections, err := a.collectionAdapter.SearchLists(ctx, &mediatypes.QueryOptions{
			ExternalSourceID: listID,
		})
		if err == nil && len(collections) > 0 {
			return mediatypes.MediaTypeCollection, nil
		}
	}

	// If not found, but we have only one type of adapter, assume it's that type
	if a.playlistAdapter != nil && a.collectionAdapter == nil {
		return mediatypes.MediaTypePlaylist, nil
	} else if a.collectionAdapter != nil && a.playlistAdapter == nil {
		return mediatypes.MediaTypeCollection, nil
	}

	return "", fmt.Errorf("unable to determine list type for ID: %s", listID)
}

// GetListItems retrieves items from a list based on the list type
func (a *DynamicListAdapter) GetListItems(
	ctx context.Context,
	listID string,
	options *mediatypes.QueryOptions,
) ([]*models.MediaItem[mediatypes.ListData], error) {
	listType, err := a.determineListType(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to determine list type: %w", err)
	}

	// Implementation would need type conversion from concrete types to ListData interface
	// This is challenging in Go with generics since we can't directly convert
	// []MediaItem[*Playlist] to []MediaItem[ListData]

	// For a working implementation, we might need:
	// 1. Ensure all items implement a GetAsListData() method
	// 2. Use type assertions in a conversion helper

	// For now, let's return a placeholder error
	return nil, fmt.Errorf("dynamic conversion not implemented for type: %s", listType)
}
