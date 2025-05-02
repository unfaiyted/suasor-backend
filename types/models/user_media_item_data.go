package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"suasor/clients/media/types"
	"time"
)

// Represents the user's personal data for a specific media item
type UserMediaItemData[T types.MediaData] struct {
	ID               uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID             string          `json:"uuid" gorm:"type:uuid;uniqueIndex"` // Stable UUID for syncing
	UserID           uint64          `json:"userId" gorm:"index"`               // Foreign key to User
	MediaItemID      uint64          `json:"mediaItemId" gorm:"index"`          // Foreign key to MediaItem
	Item             *MediaItem[T]   `json:"item" gorm:"-"`                     // Not stored in DB, loaded via relationship
	Type             types.MediaType `json:"type" gorm:"type:varchar(50)"`      // "movie", "episode", "show", "season"
	PlayedAt         time.Time       `json:"playedAt" gorm:"index"`
	LastPlayedAt     time.Time       `json:"lastPlayedAt" gorm:"index"`
	IsFavorite       bool            `json:"isFavorite,omitempty"`
	IsDisliked       bool            `json:"isDisliked,omitempty"`
	UserRating       float32         `json:"userRating,omitempty"`
	Watchlist        bool            `json:"watchlist,omitempty"`
	PlayedPercentage float64         `json:"playedPercentage,omitempty"`
	PlayCount        int32           `json:"playCount,omitempty"`
	PositionSeconds  int             `json:"positionSeconds"`
	DurationSeconds  int             `json:"durationSeconds"`
	Completed        bool            `json:"completed"`
	CreatedAt        time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (h *UserMediaItemData[T]) TableName() string {
	return "user_media_item_data"
}

// Associate links this history record with a media item
func (h *UserMediaItemData[T]) Associate(item *MediaItem[T]) {
	h.MediaItemID = item.ID
	h.Item = item
	h.Type = item.Type
}

// LoadItem loads the associated MediaItem from the database using MediaItemID
// This would be called after retrieving the history record from the database
func (h *UserMediaItemData[T]) LoadItem(db *gorm.DB) error {
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
func (h *UserMediaItemData[T]) BeforeSave(tx *gorm.DB) error {
	if h.MediaItemID == 0 && h.Item != nil {
		h.MediaItemID = h.Item.ID
	}
	return nil
}

func NewUserMediaItemData[T types.MediaData](item *MediaItem[T], userID uint64) *UserMediaItemData[T] {
	// Create a new user media item data object with the media item
	result := &UserMediaItemData[T]{
		ID:          0, // Placeholder ID
		UUID:        uuid.New().String(),
		UserID:      userID,
		MediaItemID: item.ID,
		Item:        item,
		Type:        item.Type,
		// Default values for other fields
		IsFavorite:       false,
		UserRating:       0,
		PlayedPercentage: 0,
		Watchlist:        false,
		Completed:        false,
	}

	return result
}
