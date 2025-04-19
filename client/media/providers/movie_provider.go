package providers

import (
	"context"

	"suasor/client/media/types"
	"suasor/types/models"
)

// MovieProvider defines movie-related capabilities
type MovieProvider interface {
	SupportsMovies() bool
	GetMovies(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Movie], error)
	GetMovieByID(ctx context.Context, id string) (*models.MediaItem[*types.Movie], error)
	GetMovieGenres(ctx context.Context) ([]string, error)

	movieFactory(ctx context.Context, item *any) (*types.Movie, error)
}
