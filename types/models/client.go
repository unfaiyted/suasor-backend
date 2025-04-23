// models/automation_client.go
package models

import (
	client "suasor/clients/types"

	"database/sql/driver"
	"encoding/json"
	"errors"
)

type AIClients struct {
	Claude []*Client[*client.ClaudeConfig]
	OpenAI []*Client[*client.OpenAIConfig]
	Ollama []*Client[*client.OllamaConfig]

	Total int
}

type AutomationClients struct {
	Sonarr []*Client[*client.SonarrConfig]
	Radarr []*Client[*client.RadarrConfig]
	Lidarr []*Client[*client.LidarrConfig]

	Total int
}

type MetadataClients struct {
	// Tmdb []*Client[*client.TmdbConfig]
}

// Client represents a download client configuration
type Client[T client.ClientConfig] struct {
	BaseModel
	UserID    uint64                 `json:"userId" gorm:"not null"`
	Category  client.ClientCategory  `json:"category" gorm:"not null"`
	Type      client.ClientType      `json:"type" gorm:"not null"`
	Config    ClientConfigWrapper[T] `json:"config" gorm:"type:jsonb"`
	Name      string                 `json:"name" gorm:"not null"`
	IsEnabled bool                   `json:"isEnabled" gorm:"default:true"`
}

func (c Client[T]) GetConfig() T {
	return c.Config.Data
}

func (c Client[T]) SupportsMovies() bool {
	if mediaConfig, ok := any(c.Config.Data).(client.ClientMediaConfig); ok {
		return mediaConfig.SupportsMovies()
	}
	return false
}

func (c Client[T]) SupportsSeries() bool {
	if mediaConfig, ok := any(c.Config.Data).(client.ClientMediaConfig); ok {
		return mediaConfig.SupportsSeries()
	}
	return false
}

func (c Client[T]) SupportsMusic() bool {
	if mediaConfig, ok := any(c.Config.Data).(client.ClientMediaConfig); ok {
		return mediaConfig.SupportsMusic()
	}
	return false
}
func (c Client[T]) SupportsPlaylists() bool {
	if mediaConfig, ok := any(c.Config.Data).(client.ClientMediaConfig); ok {
		return mediaConfig.SupportsPlaylists()
	}
	return false
}
func (c Client[T]) SupportsCollections() bool {
	if mediaConfig, ok := any(c.Config.Data).(client.ClientMediaConfig); ok {
		return mediaConfig.SupportsCollections()
	}
	return false
}
func (c Client[T]) SupportsHistory() bool {
	if mediaConfig, ok := any(c.Config.Data).(client.ClientMediaConfig); ok {
		return mediaConfig.SupportsHistory()
	}
	return false
}

func (c Client[T]) GetClientType() client.ClientType {
	if automationConfig, ok := any(c.Config.Data).(client.AutomationClientConfig); ok {
		return automationConfig.GetClientType().AsGenericClient()
	} else if mediaConfig, ok := any(c.Config.Data).(client.ClientMediaConfig); ok {
		return mediaConfig.GetClientType().AsGenericClient()
	}
	return client.ClientTypeUnknown
}

func (Client[T]) TableName() string {
	return "clients"
}

// ClientConfigWrapper wraps generic client configuration with database serialization
type ClientConfigWrapper[T client.ClientConfig] struct {
	Data T `json:"data" gorm:"type:jsonb"`
}

// Value implements driver.Valuer for database storage
func (m ClientConfigWrapper[T]) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(m)
}

// Scan implements sql.Scanner for database retrieval
func (m *ClientConfigWrapper[T]) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// First try to unmarshal into the wrapper structure
	err := json.Unmarshal(bytes, &m)
	if err == nil {
		return nil
	}

	// If that fails, try to unmarshal directly into the Data field
	var data T
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	m.Data = data
	return nil
}

// NewAutomationClient creates a new AutomationClient instance
func (*Client[T]) NewAutomationClient(userID uint64, clientType client.AutomationClientType, config client.ClientConfig, name string, url string, isEnabled bool) *Client[T] {
	return &Client[T]{
		UserID:    userID,
		Category:  clientType.AsCategory(),
		Config:    ClientConfigWrapper[T]{config.(T)},
		Name:      name,
		IsEnabled: isEnabled,
	}
}

func (c *Client[T]) GetID() uint64 {
	return c.ID
}

func (c *Client[T]) GetUserID() uint64 {
	return c.UserID
}

func (c *Client[T]) GetCategory() client.ClientCategory {
	return c.Category
}

func (c *Client[T]) GetName() string {
	return c.Name
}

func (c *Client[T]) GetType() client.ClientType {
	return c.Type
}
