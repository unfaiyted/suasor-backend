package types

// @Description Jellyfin media server configuration
type JellyfinConfig struct {
	Enabled  bool            `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL  string          `json:"baseURL" mapstructure:"host" example:"http://localhost:8096" binding:"required_if=Enabled true"`
	APIKey   string          `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Username string          `json:"username" mapstructure:"username" example:"admin"`
	UserID   string          `json:"userID,omitempty" mapstructure:"userID" example:"your-internal-user-id"`
	SSL      bool            `json:"ssl" mapstructure:"ssl" example:"false"`
	Type     MediaClientType `json:"type" mapstructure:"type" default:"jellyfin"`
}

func NewJellyfinConfig() JellyfinConfig {
	return JellyfinConfig{
		Type: MediaClientTypeJellyfin,
	}
}

func (JellyfinConfig) isMediaClientConfig() {}
func (JellyfinConfig) isClientConfig()      {}

func (JellyfinConfig) GetClientType() MediaClientType {
	return MediaClientTypeJellyfin
}
func (JellyfinConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (JellyfinConfig) SupportsMovies() bool {
	return true
}
func (JellyfinConfig) SupportsSeries() bool {
	return true
}
func (JellyfinConfig) SupportsMusic() bool {
	return true
}

func (JellyfinConfig) SupportsPlaylists() bool {
	return true
}
func (JellyfinConfig) SupportsCollections() bool {
	return true
}
func (JellyfinConfig) SupportsHistory() bool {
	return true
}
