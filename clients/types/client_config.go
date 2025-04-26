package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ClientConfig interface {
	isClientConfig()
	GetCategory() ClientCategory
	SetCategory(ClientCategory)
	GetType() ClientType
	GetBaseURL() string
	SetBaseURL(baseURL string)
}

type clientConfig struct {
	Type         ClientType     `json:"type"`
	Category     ClientCategory `json:"category"`
	Name         string         `json:"name" mapstructure:"name" example:"My Client"`
	BaseURL      string         `json:"baseURL" mapstructure:"baseURL"`
	Enabled      bool           `json:"enabled" mapstructure:"enabled" example:"true"`
	ValidateConn bool           `json:"validateConn" mapstructure:"validateConn" example:"true"`
}

func NewClientConfig(clientType ClientType, category ClientCategory, name string, baseURL string, enabled bool, validateConn bool) ClientConfig {
	return &clientConfig{
		Type:         clientType,
		Category:     category,
		Name:         name,
		BaseURL:      baseURL,
		Enabled:      enabled,
		ValidateConn: validateConn,
	}
}

func (c *clientConfig) GetType() ClientType {
	return c.Type
}

func (c *clientConfig) GetCategory() ClientCategory {
	return c.Category
}
func (c *clientConfig) SetCategory(category ClientCategory) {
	c.Category = category
}

func (clientConfig) isClientConfig() {}

// Validate validates the client configuration
func (c *clientConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("missing required field: name")
	}

	if c.BaseURL == "" {
		return fmt.Errorf("missing required field: baseURL")
	}

	return nil
}

func (c *clientConfig) GetBaseURL() string {
	return c.BaseURL
}

// setBaseURL sets the base URL for the client
func (c *clientConfig) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
}

// Value implements driver.Valuer for database storage
func (c *clientConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *clientConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// First try to unmarshal into the wrapper structure
	err := json.Unmarshal(bytes, &m)
	if err == nil {
		return nil
	}

	return nil
}
