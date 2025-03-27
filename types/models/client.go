// models/automation_client.go
package models

import (
	media "suasor/client/media/types"
	client "suasor/client/types"
	"time"
)

// AutomationClient represents a download client configuration
type Client[T client.ClientConfig] struct {
	// type AutomationClient struct {
	ID         uint64              `json:"id" gorm:"primaryKey"`
	UserID     uint64              `json:"userId" gorm:"not null"`
	ClientType client.ClientType   `json:"clientType" gorm:"not null"`
	Config     client.ClientConfig `json:"config" gorm:"not null"`
	Name       string              `json:"name" gorm:"not null"`
	URL        string              `json:"url" gorm:"not null"`
	// APIKey    string    `json:"apiKey" gorm:"not null"`
	IsEnabled bool      `json:"isEnabled" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// MediaClient represents a media client configuration
type MediaClient[T types.ClientConfig] struct {
	ID         uint64                `json:"id" gorm:"primaryKey"`
	UserID     uint64                `json:"userId" gorm:"not null"`
	Name       string                `json:"name" gorm:"not null"`
	ClientType types.MediaClientType `json:"clientType" gorm:"not null"`
	Client     T                     `json:"client" gorm:"serializer:json"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

// // ToResponse converts a generic MediaClient[T] to a MediaClientResponse
// func ToResponse[T types.ClientConfig](client *MediaClient[T]) MediaClientResponse {
// 	return MediaClientResponse{
// 		ID:         client.ID,
// 		UserID:     client.UserID,
// 		Name:       client.Name,
// 		ClientType: client.ClientType,
// 		Client:     client.Client,
// 		CreatedAt:  client.CreatedAt,
// 		UpdatedAt:  client.UpdatedAt,
// 	}
// }
