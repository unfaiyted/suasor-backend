package interfaces

import (
	"context"
	"time"
)

// WatchHistoryItem represents an item in watch history
type WatchHistoryItem[T MediaData] struct {
	Item             MediaItem[T]
	ItemType         string    `json:"itemType"` // "movie", "episode" , "show","season"
	WatchedAt        time.Time `json:"watchedAt"`
	LastWatchedAt    time.Time `json:"lastWatchedAt"`
	IsFavorite       bool      `json:"isFavorite,omitempty"`
	PlayedPercentage float64   `json:"playedPercentage,omitempty"`
	PlayCount        int32     `json:"playCount,omitempty"`
	PositionSeconds  int       `json:"positionSeconds"`
	DurationSeconds  int       `json:"durationSeconds"`
	Completed        bool      `json:"completed"`
	SeriesName       string    `json:"seriesName,omiempty"`
	SeasonNumber     int       `json:"seasonNumber,omitempty"`
	EpisodeNumber    int       `json:"episodeNumber,omitempty"`
}

// WatchHistoryProvider defines watch history capabilities
type WatchHistoryProvider interface {
	GetWatchHistory(ctx context.Context, options *QueryOptions) ([]WatchHistoryItem[MediaData], error)
}
