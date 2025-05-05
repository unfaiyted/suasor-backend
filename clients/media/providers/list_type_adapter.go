// list_type_adapter.go
package providers

import (
	"context"
	"fmt"
	"reflect"

	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
)

// TypedListProvider is a type-erased adapter that allows working with the ListData interface
// This is useful when you need to handle both playlist and collection types uniformly
type TypedListProvider struct {
	// The client that will be queried
	client any

	// Functions to access the specialized providers
	getPlaylistProvider   func() PlaylistProvider
	getCollectionProvider func() CollectionProvider
}

// NewTypedListProvider creates a new adapter
func NewTypedListProvider(
	client any,
	getPlaylistProvider func() PlaylistProvider,
	getCollectionProvider func() CollectionProvider,
) *TypedListProvider {
	return &TypedListProvider{
		client:                client,
		getPlaylistProvider:   getPlaylistProvider,
		getCollectionProvider: getCollectionProvider,
	}
}

// ListResult represents a generic list result with its type information
type ListResult struct {
	// The type of list (playlist, collection)
	Type mediatypes.MediaType

	// The actual data, to be type-asserted to the correct type
	Data interface{}
}

// TypedSearchResult is a result from a typed search operation
type TypedSearchResult struct {
	// Results from playlist provider, if any
	Playlists []*models.MediaItem[*mediatypes.Playlist]

	// Results from collection provider, if any
	Collections []*models.MediaItem[*mediatypes.Collection]
}

// SearchAllLists searches for all lists across all supported providers
func (p *TypedListProvider) SearchAllLists(
	ctx context.Context,
	options *mediatypes.QueryOptions,
) (*TypedSearchResult, error) {
	result := &TypedSearchResult{}

	// Search playlists if supported
	playlistProvider := p.getPlaylistProvider()
	if playlistProvider != nil && playlistProvider.SupportsPlaylists() {
		playlists, err := playlistProvider.SearchPlaylists(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to search playlists: %w", err)
		}
		result.Playlists = playlists
	}

	// Search collections if supported
	collectionProvider := p.getCollectionProvider()
	if collectionProvider != nil && collectionProvider.SupportsCollections() {
		collections, err := collectionProvider.SearchCollections(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to search collections: %w", err)
		}
		result.Collections = collections
	}

	return result, nil
}

// GetList retrieves a specific list by ID
func (p *TypedListProvider) GetList(
	ctx context.Context,
	listID string,
) (*ListResult, error) {
	// Try to determine the list type
	listType, err := p.determineListType(ctx, listID)
	if err != nil {
		return nil, err
	}

	// Get the list based on its type
	switch listType {
	case mediatypes.MediaTypePlaylist:
		playlistProvider := p.getPlaylistProvider()
		if playlistProvider == nil {
			return nil, fmt.Errorf("playlist provider not available")
		}

		playlists, err := playlistProvider.SearchPlaylists(ctx, &mediatypes.QueryOptions{
			ExternalSourceID: listID,
		})
		if err != nil || len(playlists) == 0 {
			return nil, fmt.Errorf("playlist not found: %w", err)
		}

		return &ListResult{
			Type: mediatypes.MediaTypePlaylist,
			Data: playlists[0],
		}, nil

	case mediatypes.MediaTypeCollection:
		collectionProvider := p.getCollectionProvider()
		if collectionProvider == nil {
			return nil, fmt.Errorf("collection provider not available")
		}

		collections, err := collectionProvider.SearchCollections(ctx, &mediatypes.QueryOptions{
			ExternalSourceID: listID,
		})
		if err != nil || len(collections) == 0 {
			return nil, fmt.Errorf("collection not found: %w", err)
		}

		return &ListResult{
			Type: mediatypes.MediaTypeCollection,
			Data: collections[0],
		}, nil

	default:
		return nil, fmt.Errorf("unsupported list type: %s", listType)
	}
}

// GetListItems gets items from a list
func (p *TypedListProvider) GetListItems(
	ctx context.Context,
	listID string,
	options *mediatypes.QueryOptions,
) (*ListResult, error) {
	// Determine the list type
	listType, err := p.determineListType(ctx, listID)
	if err != nil {
		return nil, err
	}

	// Get items based on list type
	switch listType {
	case mediatypes.MediaTypePlaylist:
		playlistProvider := p.getPlaylistProvider()
		if playlistProvider == nil {
			return nil, fmt.Errorf("playlist provider not available")
		}

		items, err := playlistProvider.GetPlaylistItems(ctx, listID)
		if err != nil {
			return nil, fmt.Errorf("failed to get playlist items: %w", err)
		}

		return &ListResult{
			Type: mediatypes.MediaTypePlaylist,
			Data: items,
		}, nil

	case mediatypes.MediaTypeCollection:
		collectionProvider := p.getCollectionProvider()
		if collectionProvider == nil {
			return nil, fmt.Errorf("collection provider not available")
		}

		items, err := collectionProvider.GetCollectionItems(ctx, listID)
		if err != nil {
			return nil, fmt.Errorf("failed to get collection items: %w", err)
		}

		return &ListResult{
			Type: mediatypes.MediaTypeCollection,
			Data: items,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported list type: %s", listType)
	}
}

// CreateList creates a new list
func (p *TypedListProvider) CreateList(
	ctx context.Context,
	listType mediatypes.MediaType,
	name string,
	description string,
) (*ListResult, error) {
	switch listType {
	case mediatypes.MediaTypePlaylist:
		playlistProvider := p.getPlaylistProvider()
		if playlistProvider == nil {
			return nil, fmt.Errorf("playlist provider not available")
		}

		playlist, err := playlistProvider.CreatePlaylist(ctx, name, description)
		if err != nil {
			return nil, fmt.Errorf("failed to create playlist: %w", err)
		}

		return &ListResult{
			Type: mediatypes.MediaTypePlaylist,
			Data: playlist,
		}, nil

	case mediatypes.MediaTypeCollection:
		collectionProvider := p.getCollectionProvider()
		if collectionProvider == nil {
			return nil, fmt.Errorf("collection provider not available")
		}

		collection, err := collectionProvider.CreateCollection(ctx, name, description)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}

		return &ListResult{
			Type: mediatypes.MediaTypeCollection,
			Data: collection,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported list type: %s", listType)
	}
}

// AddItemToList adds an item to a list
func (p *TypedListProvider) AddItemToList(
	ctx context.Context,
	listID string,
	itemID string,
) error {
	// Determine the list type
	listType, err := p.determineListType(ctx, listID)
	if err != nil {
		return err
	}

	// Add item based on list type
	switch listType {
	case mediatypes.MediaTypePlaylist:
		playlistProvider := p.getPlaylistProvider()
		if playlistProvider == nil {
			return fmt.Errorf("playlist provider not available")
		}

		return playlistProvider.AddPlaylistItem(ctx, listID, itemID)

	case mediatypes.MediaTypeCollection:
		collectionProvider := p.getCollectionProvider()
		if collectionProvider == nil {
			return fmt.Errorf("collection provider not available")
		}

		return collectionProvider.AddCollectionItem(ctx, listID, itemID)

	default:
		return fmt.Errorf("unsupported list type: %s", listType)
	}
}

// determineListType identifies the type of a list by ID
func (p *TypedListProvider) determineListType(ctx context.Context, listID string) (mediatypes.MediaType, error) {
	// Try as playlist first
	playlistProvider := p.getPlaylistProvider()
	if playlistProvider != nil && playlistProvider.SupportsPlaylists() {
		playlists, err := playlistProvider.SearchPlaylists(ctx, &mediatypes.QueryOptions{
			ExternalSourceID: listID,
		})
		if err == nil && len(playlists) > 0 {
			return mediatypes.MediaTypePlaylist, nil
		}
	}

	// Then try as collection
	collectionProvider := p.getCollectionProvider()
	if collectionProvider != nil && collectionProvider.SupportsCollections() {
		collections, err := collectionProvider.SearchCollections(ctx, &mediatypes.QueryOptions{
			ExternalSourceID: listID,
		})
		if err == nil && len(collections) > 0 {
			return mediatypes.MediaTypeCollection, nil
		}
	}

	// Default to playlist if only that provider is available
	if playlistProvider != nil && collectionProvider == nil {
		return mediatypes.MediaTypePlaylist, nil
	}

	// Default to collection if only that provider is available
	if collectionProvider != nil && playlistProvider == nil {
		return mediatypes.MediaTypeCollection, nil
	}

	return "", fmt.Errorf("unable to determine list type for ID: %s", listID)
}

// TypeHelper is a utility for working with list types
type TypeHelper struct{}

// NewTypeHelper creates a new type helper
func NewTypeHelper() *TypeHelper {
	return &TypeHelper{}
}

// IsPlaylist checks if a media item is a playlist
func (h *TypeHelper) IsPlaylist(item interface{}) bool {
	_, ok := item.(*models.MediaItem[*mediatypes.Playlist])
	return ok
}

// IsCollection checks if a media item is a collection
func (h *TypeHelper) IsCollection(item interface{}) bool {
	_, ok := item.(*models.MediaItem[*mediatypes.Collection])
	return ok
}

// GetListTitle safely extracts the title from a list item
func (h *TypeHelper) GetListTitle(item interface{}) (string, error) {
	if playlist, ok := item.(*models.MediaItem[*mediatypes.Playlist]); ok {
		return playlist.Title, nil
	}
	if collection, ok := item.(*models.MediaItem[*mediatypes.Collection]); ok {
		return collection.Title, nil
	}

	return "", fmt.Errorf("unknown item type: %s", reflect.TypeOf(item))
}

// GetListDescription safely extracts the description from a list item
func (h *TypeHelper) GetListDescription(item any) (string, error) {
	if playlist, ok := item.(*models.MediaItem[*mediatypes.Playlist]); ok {
		return playlist.Data.ItemList.Details.Description, nil
	}
	if collection, ok := item.(*models.MediaItem[*mediatypes.Collection]); ok {
		return collection.Data.ItemList.Details.Description, nil
	}

	return "", fmt.Errorf("unknown item type: %s", reflect.TypeOf(item))
}
