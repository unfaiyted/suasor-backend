package providers

import (
	"context"
	"suasor/clients/media/types"
	"suasor/types/models"
)

// PlaylistProvider defines playlist capabilities
type ListProvider[T types.ListData] interface {
	SupportsLists() bool
	Search(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[T], error)

	// Full playlist management capabilities
	GetItems(ctx context.Context, playlistID string, options *types.QueryOptions) ([]*models.MediaItem[T], error)
	Create(ctx context.Context, name string, description string) (*models.MediaItem[T], error)
	Update(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[T], error)
	Delete(ctx context.Context, playlistID string) error
	AddItem(ctx context.Context, playlistID string, itemID string) error
	RemoveItem(ctx context.Context, playlistID string, itemID string) error
	ReorderItems(ctx context.Context, playlistID string, itemIDs []string) error
}

// CollectionProvider defines collection capabilities
type CollectionProvider interface {
	ListProvider[*types.Collection]
	SupportsCollections() bool

	collectionFactory(ctx context.Context, item *any) (*types.Collection, error)
}

type PlaylistProvider interface {
	ListProvider[*types.Playlist]
	SupportsPlaylists() bool

	playlistFactory(ctx context.Context, item *any) (*types.Playlist, error)
}
