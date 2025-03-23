package interfaces

import (
	"context"
)

// Collection represents a collection of media items
type Collection struct {
	Details        MediaMetadata
	ItemIDs        []string `json:"itemIDs"`
	ItemCount      int      `json:"itemCount"`
	CollectionType string   `json:"collectionType"` // e.g., "movie", "tvshow"
}

// Playlist represents a user-created playlist of media items
type Playlist struct {
	Details   MediaMetadata
	ItemIDs   []string `json:"itemIDs"`
	ItemCount int      `json:"itemCount"`
	Owner     string   `json:"owner,omitempty"`
	IsPublic  bool     `json:"isPublic"`
}

func (c Collection) GetDetails() MediaMetadata { return c.Details }
func (c Collection) GetMediaType() MediaType   { return MEDIATYPE_COLLECTION }

func (p Playlist) GetDetails() MediaMetadata { return p.Details }
func (p Playlist) GetMediaType() MediaType   { return MEDIATYPE_PLAYLIST }

// PlaylistProvider defines playlist capabilities
type PlaylistProvider interface {
	SupportsPlaylists() bool
	GetPlaylists(ctx context.Context, options *QueryOptions) ([]MediaItem[Playlist], error)
}

// CollectionProvider defines collection capabilities
type CollectionProvider interface {
	SupportsCollections() bool
	GetCollections(ctx context.Context, options *QueryOptions) ([]MediaItem[Collection], error)
}
