package models

import "time"

// BaseModel defines common fields for all models.
type BaseModel struct {
	ID        uint64     `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}
