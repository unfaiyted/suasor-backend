package types

import "fmt"

type ClientConfig interface {
	isClientConfig()
	GetCategory() ClientCategory
	SetCategory(ClientCategory)
	GetType() ClientType
}

type BaseClientConfig struct {
	Type     ClientType     `json:"type"`
	Category ClientCategory `json:"category"`
	Name     string         `json:"name" mapstructure:"name" example:"My Client"`
	BaseURL  string         `json:"baseURL" mapstructure:"baseURL"`
	Enabled  bool           `json:"enabled" mapstructure:"enabled" example:"true"`
	ValidateConn bool       `json:"validateConn" mapstructure:"validateConn" example:"true"`
}

func (c *BaseClientConfig) GetType() ClientType {
	return c.Type
}

func (c *BaseClientConfig) GetCategory() ClientCategory {
	return c.Category
}
func (c *BaseClientConfig) SetCategory(category ClientCategory) {
	c.Category = category
}

func (BaseClientConfig) isClientConfig() {}

// Validate validates the client configuration
func (c *BaseClientConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("missing required field: name")
	}
	
	if c.BaseURL == "" {
		return fmt.Errorf("missing required field: baseURL")
	}
	
	return nil
}
