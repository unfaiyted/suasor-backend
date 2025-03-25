package providers

import (
	"context"

	"suasor/client/media/types"
)

// MovieProvider defines movie-related capabilities
type MovieProvider interface {
	SupportsMovies() bool
	GetMovies(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Movie], error)
	GetMovieByID(ctx context.Context, id string) (types.MediaItem[types.Movie], error)
	GetMovieGenres(ctx context.Context) ([]string, error)
}
