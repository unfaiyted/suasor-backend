package models

import (
	"errors"
	"gorm.io/gorm"
	"suasor/client/media/types"
	"time"
)

// HistoryItem represents an item in watch history
type MediaPlayHistory[T types.MediaData] struct {
	ID               uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           uint64          `json:"userId" gorm:"index"`          // Foreign key to User
	MediaItemID      uint64          `json:"mediaItemId" gorm:"index"`     // Foreign key to MediaItem
	Item             *MediaItem[T]   `json:"item" gorm:"-"`                // Not stored in DB, loaded via relationship
	Type             types.MediaType `json:"type" gorm:"type:varchar(50)"` // "movie", "episode", "show", "season"
	PlayedAt         time.Time       `json:"playedAt" gorm:"index"`
	LastPlayedAt     time.Time       `json:"lastPlayedAt" gorm:"index"`
	IsFavorite       bool            `json:"isFavorite,omitempty"`
	IsDisliked       bool            `json:"isDisliked,omitempty"`
	UserRating       float32         `json:"userRating,omitempty"`
	PlayedPercentage float64         `json:"playedPercentage,omitempty"`
	PlayCount        int32           `json:"playCount,omitempty"`
	PositionSeconds  int             `json:"positionSeconds"`
	DurationSeconds  int             `json:"durationSeconds"`
	Completed        bool            `json:"completed"`
	CreatedAt        time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}

// Associate links this history record with a media item
func (h *MediaPlayHistory[T]) Associate(item *MediaItem[T]) {
	h.MediaItemID = item.ID
	h.Item = item
	h.Type = item.Type
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

// MediaPlayHistoryGeneric is a non-generic version of MediaPlayHistory to avoid type issues
type MediaPlayHistoryGeneric struct {
	ID               uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           uint64          `json:"userId" gorm:"index"`          
	MediaItemID      uint64          `json:"mediaItemId" gorm:"index"`     
	Type             types.MediaType `json:"type" gorm:"type:varchar(50)"` 
	PlayedAt         time.Time       `json:"playedAt" gorm:"index"`
	LastPlayedAt     time.Time       `json:"lastPlayedAt" gorm:"index"`
	IsFavorite       bool            `json:"isFavorite,omitempty"`
	IsDisliked       bool            `json:"isDisliked,omitempty"`
	UserRating       float32         `json:"userRating,omitempty"`
	PlayedPercentage float64         `json:"playedPercentage,omitempty"`
	PlayCount        int32           `json:"playCount,omitempty"`
	PositionSeconds  int             `json:"positionSeconds"`
	DurationSeconds  int             `json:"durationSeconds"`
	Completed        bool            `json:"completed"`
	CreatedAt        time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}

// MediaPlayHistoryRequest is used to record a new play history entry
type MediaPlayHistoryRequest struct {
	UserID           uint64          `json:"userId" binding:"required"`
	MediaItemID      uint64          `json:"mediaItemId" binding:"required"`
	Type             types.MediaType `json:"type" binding:"required"`
	IsFavorite       bool            `json:"isFavorite,omitempty"`
	UserRating       float32         `json:"userRating,omitempty"`
	PlayedPercentage float64         `json:"playedPercentage,omitempty"`
	PositionSeconds  int             `json:"positionSeconds"`
	DurationSeconds  int             `json:"durationSeconds"`
	Completed        bool            `json:"completed"`
	Continued        bool            `json:"continued"` // If this is a continuation of a previous play
}
