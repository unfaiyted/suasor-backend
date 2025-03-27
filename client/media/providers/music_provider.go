package providers

import (
	"context"
	"suasor/client/media/types"
	"suasor/types/models"
)

// MusicProvider defines music-related capabilities
type MusicProvider interface {
	SupportsMusic() bool
	GetMusic(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[types.Track], error)
	GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[types.Artist], error)
	GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[types.Album], error)
	GetMusicTrackByID(ctx context.Context, id string) (models.MediaItem[types.Track], error)
	GetMusicGenres(ctx context.Context) ([]string, error)
}
