package types

// @Description Emby media server configuration
type EmbyConfig struct {
	Enabled  bool            `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL  string          `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey   string          `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Username string          `json:"username" mapstructure:"username" example:"admin"`
	UserID   string          `json:"userID,omitempty" mapstructure:"userID" example:"your-internal-user-id"`
	SSL      bool            `json:"ssl" mapstructure:"ssl" example:"false"`
	Type     MediaClientType `json:"type" mapstructure:"type" default:"emby"`
}

func NewEmbyConfig() EmbyConfig {
	return EmbyConfig{
		Type: MediaClientTypeEmby,
	}
}

func (EmbyConfig) isClientConfig()      {}
func (EmbyConfig) isMediaClientConfig() {}

func (EmbyConfig) GetClientType() MediaClientType {
	return MediaClientTypeEmby
}
func (EmbyConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (EmbyConfig) SupportsMovies() bool {
	return true
}
func (EmbyConfig) SupportsSeries() bool {
	return true
}
func (EmbyConfig) SupportsMusic() bool {
	return true
}
func (EmbyConfig) SupportsPlaylists() bool {
	return true
}
func (EmbyConfig) SupportsCollections() bool {
	return true
}
func (EmbyConfig) SupportsHistory() bool {
	return true
}
