// models/automation_client.go
package models

import (
	client "suasor/client/types"

	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// AutomationClient represents a download client configuration
type Client[T client.ClientConfig] struct {
	BaseModel
	UserID    uint64                 `json:"userId" gorm:"not null"`
	Type      client.ClientType      `json:"type" gorm:"not null"`
	Config    ClientConfigWrapper[T] `json:"config" gorm:"jsonb"`
	Name      string                 `json:"name" gorm:"not null"`
	IsEnabled bool                   `json:"isEnabled" gorm:"default:true"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

func (Client[T]) TableName() string {
	return "clients"
}

// ClientConfigWrapper wraps generic client configuration with database serialization
type ClientConfigWrapper[T client.ClientConfig] struct {
	Data T
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

	return json.Unmarshal(bytes, &m)
}

// NewAutomationClient creates a new AutomationClient instance
func (*Client[T]) NewAutomationClient(userID uint64, clientType client.ClientType, config client.ClientConfig, name string, url string, isEnabled bool) *Client[T] {
	return &Client[T]{
		UserID:    userID,
		Type:      clientType,
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

func (c *Client[T]) GetClientType() client.ClientType {
	return c.Type
}

func (c *Client[T]) GetConfig() client.ClientConfig {
	return c.Config.Data
}

func (c *Client[T]) GetName() string {
	return c.Name
}
