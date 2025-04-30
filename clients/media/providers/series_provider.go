package providers

import (
	"context"
	"suasor/clients/media/types"
	"suasor/types/models"
)

// SeriesProvider defines TV show-related capabilities
type SeriesProvider interface {
	SupportsSeries() bool
	GetSeries(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Series], error)
	GetSeriesSeasons(ctx context.Context, clientSeriesID string) ([]*models.MediaItem[*types.Season], error)
	GetSeriesEpisodesBySeasonNbr(ctx context.Context, clientSeriesID string, seasonNumber int) ([]*models.MediaItem[*types.Episode], error)

	GetSeriesByID(ctx context.Context, clientSeriesID string) (*models.MediaItem[*types.Series], error)
	GetEpisodeByID(ctx context.Context, clientEpisodeID string) (*models.MediaItem[*types.Episode], error)
}
