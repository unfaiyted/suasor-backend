package models

import (
	"errors"
	"gorm.io/gorm"
	"suasor/client/media/types"
	"time"
)

// HistoryItem represents an item in watch history
type MediaPlayHistory[T types.MediaData] struct {
	ID               uint64        `json:"id" gorm:"primaryKey;autoIncrement"`
	MediaItemID      uint64        `json:"mediaItemId" gorm:"index"`     // Foreign key to MediaItem
	Item             *MediaItem[T] `json:"item" gorm:"-"`                // Not stored in DB, loaded via relationship
	Type             string        `json:"type" gorm:"type:varchar(50)"` // "movie", "episode", "show", "season"
	WatchedAt        time.Time     `json:"watchedAt" gorm:"index"`
	LastWatchedAt    time.Time     `json:"lastWatchedAt" gorm:"index"`
	IsFavorite       bool          `json:"isFavorite,omitempty"`
	PlayedPercentage float64       `json:"playedPercentage,omitempty"`
	PlayCount        int32         `json:"playCount,omitempty"`
	PositionSeconds  int           `json:"positionSeconds"`
	DurationSeconds  int           `json:"durationSeconds"`
	Completed        bool          `json:"completed"`
	CreatedAt        time.Time     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time     `json:"updatedAt" gorm:"autoUpdateTime"`
}

// Associate links this history record with a media item
func (h *MediaPlayHistory[T]) Associate(item *MediaItem[T]) {
	h.MediaItemID = item.ID
	h.Item = item
	h.Type = string(item.Type)
}

// LoadItem loads the associated MediaItem from the database using MediaItemID
// This would be called after retrieving the history record from the database
func (h *MediaPlayHistory[T]) LoadItem(db *gorm.DB) error {
	if h.MediaItemID == 0 {
		return errors.New("no media item ID associated with this history record")
	}

	item := &MediaItem[T]{}
	result := db.First(item, h.MediaItemID)
	if result.Error != nil {
		return result.Error
	}

	h.Item = item
	return nil
}

// BeforeSave ensures we have the proper MediaItemID before saving
func (h *MediaPlayHistory[T]) BeforeSave(tx *gorm.DB) error {
	if h.MediaItemID == 0 && h.Item != nil {
		h.MediaItemID = h.Item.ID
	}
	return nil
}
