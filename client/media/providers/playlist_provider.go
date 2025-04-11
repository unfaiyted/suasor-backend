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
}

// CollectionProvider defines collection capabilities
type CollectionProvider interface {
	SupportsCollections() bool
	GetCollections(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[*types.Collection], error)
}
