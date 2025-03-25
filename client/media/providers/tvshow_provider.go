package providers

import (
	"context"
	"suasor/client/media/types"
)

// TVShowProvider defines TV show-related capabilities
type TVShowProvider interface {
	SupportsTVShows() bool
	GetTVShows(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.TVShow], error)
	GetTVShowSeasons(ctx context.Context, showID string) ([]types.MediaItem[types.Season], error)
	GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]types.MediaItem[types.Episode], error)
	GetTVShowByID(ctx context.Context, id string) (types.MediaItem[types.TVShow], error)
	GetEpisodeByID(ctx context.Context, id string) (types.MediaItem[types.Episode], error)
}
