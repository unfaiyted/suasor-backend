package models

import (
	"time"
)

// SearchHistory represents a user's search history entry
type SearchHistory struct {
	BaseModel
	UserID      uint64    `json:"userId" gorm:"index"`
	Query       string    `json:"query" gorm:"index"`
	ResultCount int       `json:"resultCount"`
	SearchedAt  time.Time `json:"searchedAt" gorm:"index"`
}
