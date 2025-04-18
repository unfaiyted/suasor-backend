package providers

import (
	"context"
	"suasor/client/media/types"
	"suasor/types/models"
)

// SeriesProvider defines TV show-related capabilities
type SeriesProvider interface {
	SupportsSeries() bool
	GetSeries(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Series], error)
	GetSeriesSeasons(ctx context.Context, showID string) ([]*models.MediaItem[*types.Season], error)
	GetSeriesEpisodes(ctx context.Context, showID string, seasonNumber int) ([]*models.MediaItem[*types.Episode], error)
	GetSeriesByID(ctx context.Context, id string) (*models.MediaItem[*types.Series], error)
	GetEpisodeByID(ctx context.Context, id string) (*models.MediaItem[*types.Episode], error)
}
