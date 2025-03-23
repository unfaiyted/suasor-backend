package interfaces

import (
	"context"
	"time"
)

// Season represents a TV season
type Season struct {
	Details      MediaMetadata
	Number       int       `json:"seasonNumber"`
	Title        string    `json:"title,omitempty"`
	Overview     string    `json:"overview,omitempty"`
	EpisodeCount int       `json:"episodeCount"`
	Artwork      Artwork   `json:"artwork,omitempty"`
	ReleaseDate  time.Time `json:"releaseDate,omitempty"`
	ParentID     string    `json:"parentID,omitempty"`
}

// Episode represents a TV episode
type Episode struct {
	Details      MediaMetadata
	Number       int64  `json:"number"`
	ShowID       string `json:"showID"`
	SeasonID     string `json:"seasonID"`
	SeasonNumber int    `json:"seasonNumber"`
	ShowTitle    string `json:"showTitle,omitempty"`
}

// TVShow represents a TV series
type TVShow struct {
	Details       MediaMetadata
	Seasons       []Season `json:"seasons,omitempty"`
	EpisodeCount  int      `json:"episodeCount"`
	SeasonCount   int      `json:"seasonCount"`
	ReleaseYear   int      `json:"releaseYear"`
	ContentRating string   `json:"contentRating"`
	Rating        float64  `json:"rating"`
	Network       string   `json:"network,omitempty"`
	Status        string   `json:"status,omitempty"` // e.g., "Ended", "Continuing"
	Genres        []string `json:"genres,omitempty"`
}

func (t TVShow) GetDetails() MediaMetadata { return t.Details }
func (t TVShow) GetMediaType() MediaType   { return MEDIATYPE_SHOW }

func (s Season) GetDetails() MediaMetadata { return s.Details }
func (s Season) GetMediaType() MediaType   { return MEDIATYPE_SEASON }

func (e Episode) GetDetails() MediaMetadata { return e.Details }
func (e Episode) GetMediaType() MediaType   { return MEDIATYPE_EPISODE }

// TVShowProvider defines TV show-related capabilities
type TVShowProvider interface {
	SupportsTVShows() bool
	GetTVShows(ctx context.Context, options *QueryOptions) ([]MediaItem[TVShow], error)
	GetTVShowSeasons(ctx context.Context, showID string) ([]MediaItem[Season], error)
	GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]MediaItem[Episode], error)
	GetTVShowByID(ctx context.Context, id string) (MediaItem[TVShow], error)
	GetEpisodeByID(ctx context.Context, id string) (MediaItem[Episode], error)
}
