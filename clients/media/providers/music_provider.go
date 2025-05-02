package providers

import (
	"context"
	"suasor/clients/media/types"
	"suasor/types/models"
)

// MusicProvider defines music-related capabilities
type MusicProvider interface {
	SupportsMusic() bool

	GetMusicTracks(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Track], error)
	GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Artist], error)
	GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Album], error)

	GetMusicTrackByID(ctx context.Context, trackID string) (*models.MediaItem[*types.Track], error)
	GetMusicArtistByID(ctx context.Context, artistID string) (*models.MediaItem[*types.Artist], error)
	GetMusicAlbumByID(ctx context.Context, albumID string) (*models.MediaItem[*types.Album], error)

	GetMusicGenres(ctx context.Context) ([]string, error)
}
