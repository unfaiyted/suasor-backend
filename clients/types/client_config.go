package types

import "fmt"

type ClientConfig interface {
	isClientConfig()
	GetCategory() ClientCategory
	SetCategory(ClientCategory)
	GetType() ClientType
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
