package providers

import (
	"context"

	"suasor/clients/media/types"
	"suasor/types/models"
)

// CollectionProvider defines collection capabilities
type CollectionProvider interface {
	GetCollection(ctx context.Context, collectionID string) (*models.MediaItem[*types.Collection], error)
	GetCollectionItems(ctx context.Context, collectionID string) (*models.MediaItemList[*types.Collection], error)
	CreateCollection(ctx context.Context, name string, description string) (*models.MediaItem[*types.Collection], error)
	CreateCollectionWithItems(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*types.Collection], error)
	UpdateCollection(ctx context.Context, collectionID string, name string, description string) (*models.MediaItem[*types.Collection], error)
	DeleteCollection(ctx context.Context, collectionID string) error
	AddCollectionItem(ctx context.Context, collectionID string, itemID string) error
	AddCollectionItems(ctx context.Context, collectionID string, itemIDs []string) error
	RemoveCollectionItem(ctx context.Context, collectionID string, itemID string) error
	RemoveCollectionItems(ctx context.Context, collectionID string, itemIDs []string) error
	RemoveAllCollectionItems(ctx context.Context, collectionID string) error
	ReorderCollectionItems(ctx context.Context, collectionID string, itemIDs []string) error

	SearchCollections(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Collection], error)
	SearchCollectionItems(ctx context.Context, collectionID string, options *types.QueryOptions) ([]*models.MediaItem[*types.Collection], error)

	SupportsCollections() bool
}

// PlaylistProvider defines playlist capabilities
type PlaylistProvider interface {
	GetPlaylist(ctx context.Context, playlistID string) (*models.MediaItem[*types.Playlist], error)
	GetPlaylistItems(ctx context.Context, playlistID string) (*models.MediaItemList[*types.Playlist], error)
	CreatePlaylist(ctx context.Context, name string, description string) (*models.MediaItem[*types.Playlist], error)
	CreatePlaylistWithItems(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*types.Playlist], error)
	UpdatePlaylist(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[*types.Playlist], error)
	DeletePlaylist(ctx context.Context, playlistID string) error
	AddPlaylistItem(ctx context.Context, playlistID string, itemID string) error
	AddPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error
	RemovePlaylistItem(ctx context.Context, playlistID string, itemID string) error
	RemovePlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error
	RemoveAllPlaylistItems(ctx context.Context, playlistID string) error
	ReorderPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error

	SearchPlaylists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error)
	SearchPlaylistItems(ctx context.Context, playlistID string, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error)

	SupportsPlaylists() bool
}

type ListProvider[T types.ListData] interface {
	// Full collection management capabilities
	GetList(ctx context.Context, listID string) (*models.MediaItem[T], error)
	GetListItems(ctx context.Context, collectionID string) (*models.MediaItemList[T], error)
	CreateList(ctx context.Context, name string, description string) (*models.MediaItem[T], error)
	CreateListWithItems(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[T], error)
	UpdateList(ctx context.Context, collectionID string, name string, description string) (*models.MediaItem[T], error)
	DeleteList(ctx context.Context, collectionID string) error
	AddListItem(ctx context.Context, collectionID string, itemID string) error
	AddListItems(ctx context.Context, collectionID string, itemIDs []string) error
	RemoveListItem(ctx context.Context, collectionID string, itemID string) error
	RemoveListItems(ctx context.Context, collectionID string, itemIDs []string) error
	RemoveAllListItems(ctx context.Context, collectionID string) error
	ReorderListItems(ctx context.Context, collectionID string, itemIDs []string) error

	SearchLists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[T], error)
	SearchListItems(ctx context.Context, collectionID string, options *types.QueryOptions) ([]*models.MediaItem[T], error)

	SupportsLists() bool
}
