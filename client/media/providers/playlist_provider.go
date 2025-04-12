package providers

import (
	"context"
	"suasor/client/media/types"
	"suasor/types/models"
)

// PlaylistProvider defines playlist capabilities
type PlaylistProvider interface {
	SupportsPlaylists() bool
	GetPlaylists(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[*types.Playlist], error)
	
	// Full playlist management capabilities
	GetPlaylistItems(ctx context.Context, playlistID string, options *types.QueryOptions) ([]models.MediaItem[types.MediaData], error)
	CreatePlaylist(ctx context.Context, name string, description string) (*models.MediaItem[*types.Playlist], error)
	UpdatePlaylist(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[*types.Playlist], error)
	DeletePlaylist(ctx context.Context, playlistID string) error
	AddItemToPlaylist(ctx context.Context, playlistID string, itemID string) error
	RemoveItemFromPlaylist(ctx context.Context, playlistID string, itemID string) error
	ReorderPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error
}

// CollectionProvider defines collection capabilities
type CollectionProvider interface {
	SupportsCollections() bool
	GetCollections(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[*types.Collection], error)
	
	// Full collection management capabilities
	GetCollectionItems(ctx context.Context, collectionID string, options *types.QueryOptions) ([]models.MediaItem[types.MediaData], error)
	CreateCollection(ctx context.Context, name string, description string, collectionType string) (*models.MediaItem[*types.Collection], error)
	UpdateCollection(ctx context.Context, collectionID string, name string, description string) (*models.MediaItem[*types.Collection], error)
	DeleteCollection(ctx context.Context, collectionID string) error
	AddItemToCollection(ctx context.Context, collectionID string, itemID string) error
	RemoveItemFromCollection(ctx context.Context, collectionID string, itemID string) error
}
