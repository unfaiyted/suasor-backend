package providers

import (
	"context"

	"suasor/clients/media/types"
	"suasor/types/models"
)

// MovieProvider defines movie-related capabilities
type MovieProvider interface {
	SupportsMovies() bool
	GetMovies(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Movie], error)
	GetMovieByID(ctx context.Context, id string) (*models.MediaItem[*types.Movie], error)
	GetMovieGenres(ctx context.Context) ([]string, error)

	// Note: The factory method is removed from the interface as it's implementation-specific
	// and should be handled via the registration system instead
}
