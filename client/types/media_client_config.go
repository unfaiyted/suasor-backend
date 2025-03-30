package types

type MediaClientConfig interface {
	ClientConfig
	isMediaClientConfig()
	GetClientType() MediaClientType

	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
	SupportsPlaylists() bool
	SupportsCollections() bool
	SupportsHistory() bool
}

type BaseMediaClientConfig struct {
	BaseClientConfig
	Type    MediaClientType
	BaseURL string `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

func (BaseMediaClientConfig) isMediaClientConfig() {}

func (c *BaseMediaClientConfig) GetName() string {
	return c.Name
}

func (c *BaseMediaClientConfig) GetBaseURL() string {
	return c.BaseURL
}

func (c *BaseMediaClientConfig) GetAPIKey() string {
	return c.APIKey
}

func (c *BaseMediaClientConfig) GetClientType() MediaClientType {
	return c.Type
}

func (c *BaseMediaClientConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (c *BaseMediaClientConfig) SupportsMovies() bool {
	return false
}
func (c *BaseMediaClientConfig) SupportsSeries() bool {
	return false
}
func (c *BaseMediaClientConfig) SupportsMusic() bool {
	return false
}
func (c *BaseMediaClientConfig) SupportsPlaylists() bool {
	return false
}
func (c *BaseMediaClientConfig) SupportsCollections() bool {
	return false
}
