package interfaces

import (
	"context"
)

// Movie represents a movie item
type Movie struct {
	Details      MediaMetadata
	Cast         []Person `json:"cast,omitempty"`
	Crew         []Person `json:"crew,omitempty"`
	TrailerURL   string   `json:"trailerUrl,omitempty"`
	Resolution   string   `json:"resolution,omitempty"` // e.g., "4K", "1080p"
	VideoCodec   string   `json:"videoCodec,omitempty"`
	AudioCodec   string   `json:"audioCodec,omitempty"`
	SubtitleURLs []string `json:"subtitleUrls,omitempty"`
}

// Then in each media type
func (m Movie) GetDetails() MediaMetadata { return m.Details }
func (m Movie) GetMediaType() MediaType   { return MEDIATYPE_MOVIE }

// MovieProvider defines movie-related capabilities
type MovieProvider interface {
	SupportsMovies() bool
	GetMovies(ctx context.Context, options *QueryOptions) ([]MediaItem[Movie], error)
	GetMovieByID(ctx context.Context, id string) (MediaItem[Movie], error)
	GetMovieGenres(ctx context.Context) ([]string, error)
}
