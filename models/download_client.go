// models/download_client.go
package models

import "time"

// ClientType represents different types of download clients
type ClientType string

const (
	ClientTypeRadarr ClientType = "radarr"
	ClientTypeSonarr ClientType = "sonarr"
	ClientTypeLidarr ClientType = "lidarr"
)

// AutomationClient represents a download client configuration
type AutomationClient struct {
	ID         uint64     `json:"id" gorm:"primaryKey"`
	UserID     uint64     `json:"userId" gorm:"not null"`
	ClientType ClientType `json:"clientType" gorm:"not null"`
	Name       string     `json:"name" gorm:"not null"`
	URL        string     `json:"url" gorm:"not null"`
	APIKey     string     `json:"apiKey" gorm:"not null"`
	IsEnabled  bool       `json:"isEnabled" gorm:"default:true"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// AutomationClientRequest is used for creating/updating a download client
type AutomationClientRequest struct {
	Name       string     `json:"name" binding:"required"`
	ClientType ClientType `json:"clientType" binding:"required,oneof=radarr sonarr lidarr"`
	URL        string     `json:"url" binding:"required,url"`
	APIKey     string     `json:"apiKey" binding:"required"`
	IsEnabled  bool       `json:"isEnabled"`
}

// ClientTestRequest is used for testing a client connection
type ClientTestRequest struct {
	URL        string     `json:"url" binding:"required,url"`
	APIKey     string     `json:"apiKey" binding:"required"`
	ClientType ClientType `json:"clientType" binding:"required,oneof=radarr sonarr lidarr"`
}

// ClientTestResponse is the response for a client connection test
type ClientTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}
