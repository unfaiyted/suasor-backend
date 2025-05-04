package providers

import (
	"context"
	"suasor/clients/media/types"
	"suasor/types/models"
)

// CollectionProvider defines collection capabilities
type CollectionProvider interface {
	GetCollectionItems(ctx context.Context, playlistID string, options *types.QueryOptions) ([]*models.MediaItem[*types.Collection], error)
	CreateCollection(ctx context.Context, name string, description string) (*models.MediaItem[*types.Collection], error)
	UpdateCollection(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[*types.Collection], error)
	DeleteCollection(ctx context.Context, playlistID string) error
	AddItemCollection(ctx context.Context, playlistID string, itemID string) error
	RemoveCollectionItem(ctx context.Context, playlistID string, itemID string) error
	ReorderCollectionItems(ctx context.Context, playlistID string, itemIDs []string) error

	SearchCollections(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Collection], error)

	SupportsCollections() bool
}

// PlaylistProvider defines playlist capabilities
type PlaylistProvider interface {
	GetPlaylistItems(ctx context.Context, playlistID string, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error)
	CreatePlaylist(ctx context.Context, name string, description string) (*models.MediaItem[*types.Playlist], error)
	UpdatePlaylist(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[*types.Playlist], error)
	DeletePlaylist(ctx context.Context, playlistID string) error
	AddItemPlaylist(ctx context.Context, playlistID string, itemID string) error
	RemovePlaylistItem(ctx context.Context, playlistID string, itemID string) error
	ReorderPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error

	SearchPlaylists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error)

	SupportsPlaylists() bool
}

type ListProvider[T types.ListData] interface {

	// Full playlist management capabilities
	GetListItems(ctx context.Context, playlistID string, options *types.QueryOptions) ([]*models.MediaItem[T], error)
	CreateList(ctx context.Context, name string, description string) (*models.MediaItem[T], error)
	UpdateList(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[T], error)
	DeleteList(ctx context.Context, playlistID string) error
	AddItemList(ctx context.Context, playlistID string, itemID string) error
	RemoveListItem(ctx context.Context, playlistID string, itemID string) error
	ReorderListItems(ctx context.Context, playlistID string, itemIDs []string) error

	SearchLists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[T], error)

	SupportsLists() bool
}
