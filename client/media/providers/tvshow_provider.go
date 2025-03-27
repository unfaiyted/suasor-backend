package providers

import (
	"context"
	"suasor/client/media/types"
	"suasor/types/models"
)

// TVShowProvider defines TV show-related capabilities
type TVShowProvider interface {
	SupportsTVShows() bool
	GetTVShows(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[types.TVShow], error)
	GetTVShowSeasons(ctx context.Context, showID string) ([]models.MediaItem[types.Season], error)
	GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]models.MediaItem[types.Episode], error)
	GetTVShowByID(ctx context.Context, id string) (models.MediaItem[types.TVShow], error)
	GetEpisodeByID(ctx context.Context, id string) (models.MediaItem[types.Episode], error)
}
