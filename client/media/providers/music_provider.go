package providers

import (
	"context"
	"suasor/client/media/types"
)

// MusicProvider defines music-related capabilities
type MusicProvider interface {
	SupportsMusic() bool
	GetMusic(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Track], error)
	GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Artist], error)
	GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Album], error)
	GetMusicTrackByID(ctx context.Context, id string) (types.MediaItem[types.Track], error)
	GetMusicGenres(ctx context.Context) ([]string, error)
}
