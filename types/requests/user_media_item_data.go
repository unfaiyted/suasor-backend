package requests

import (
	"suasor/clients/media/types"
)

// UserMediaItemDataRequest represents the data for recording a new play history entry
// @Description Request payload for recording a new play history entry
type UserMediaItemDataRequest struct {
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

// UserMediaItemDataSyncRequest represents the data for synchronizing media item data
// @Description Request payload for synchronizing media item data
type UserMediaItemDataSyncRequest struct {
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

type UserMediaItemDataUpdateRequest struct {
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
