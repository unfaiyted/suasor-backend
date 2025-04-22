package types

type ClientMediaConfig interface {
	ClientConfig
	isClientMediaConfig()
	GetClientType() ClientMediaType

	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
	SupportsPlaylists() bool
	SupportsCollections() bool
	SupportsHistory() bool
}

type BaseClientMediaConfig struct {
	BaseClientConfig
	ClientType ClientMediaType `json:"clientType"`
	BaseURL    string          `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey     string          `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL        bool            `json:"ssl" mapstructure:"ssl" example:"false"`
}

func (BaseClientMediaConfig) isClientMediaConfig() {}

// func (c *BaseClientMediaConfig) GetType() ClientType {
// 	return c.BaseClientConfig.Type
// }

func (c *BaseClientMediaConfig) GetBaseURL() string {
	return c.BaseURL
}

func (c *BaseClientMediaConfig) GetAPIKey() string {
	return c.APIKey
}

func (c *BaseClientMediaConfig) GetClientType() ClientMediaType {
	return c.ClientType
}

func (c *BaseClientMediaConfig) GetCategory() ClientCategory {
	c.Category = ClientCategoryMedia
	return c.Category
}

func (c *BaseClientMediaConfig) SupportsMovies() bool {
	return false
}
func (c *BaseClientMediaConfig) SupportsSeries() bool {
	return false
}
func (c *BaseClientMediaConfig) SupportsMusic() bool {
	return false
}
func (c *BaseClientMediaConfig) SupportsPlaylists() bool {
	return false
}
func (c *BaseClientMediaConfig) SupportsCollections() bool {
	return false
}
