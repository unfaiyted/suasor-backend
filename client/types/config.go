package types

type ClientConfig interface {
	isClientConfig()
}

type MediaClientConfig interface {
	isMediaClientConfig()
}

type AutomationClientConfig interface {
	isAutomationClientConfig()
}

func (EmbyConfig) isClientConfig()     {}
func (JellyfinConfig) isClientConfig() {}
func (PlexConfig) isClientConfig()     {}
func (SubsonicConfig) isClientConfig() {}

func (RadarrConfig) isClientConfig() {}
func (LidarrConfig) isClientConfig() {}
func (SonarrConfig) isClientConfig() {}

func (RadarrConfig) isAutomationClientConfig() {}
func (LidarrConfig) isAutomationClientConfig() {}
func (SonarrConfig) isAutomationClientConfig() {}

func (EmbyConfig) isMediaClientConfig()     {}
func (JellyfinConfig) isMediaClientConfig() {}
func (PlexConfig) isMediaClientConfig()     {}
func (SubsonicConfig) isMediaClientConfig() {}

// @Description Emby media server configuration
type SonarrConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Emby media server configuration
type RadarrConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Jellyfin media server configuration
type LidarrConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string `json:"baseURL" mapstructure:"host" example:"http://localhost:8096" binding:"required_if=Enabled true"`
	APIKey  string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Emby media server configuration
type EmbyConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL  string `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey   string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
	UserID   string `json:"userID,omitempty" mapstructure:"userID" example:"your-internal-user-id"`
	SSL      bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Jellyfin media server configuration
type JellyfinConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL  string `json:"baseURL" mapstructure:"host" example:"http://localhost:8096" binding:"required_if=Enabled true"`
	APIKey   string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
	UserID   string `json:"userID,omitempty" mapstructure:"userID" example:"your-internal-user-id"`
	SSL      bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Plex media server configuration
type PlexConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host    string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port    int    `json:"port" mapstructure:"port" example:"32400" binding:"required_if=Enabled true"`
	Token   string `json:"token" mapstructure:"token" example:"your-plex-token" binding:"required_if=Enabled true"`
	SSL     bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Supersonic music server configuration
type SubsonicConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int    `json:"port" mapstructure:"port" example:"4533" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin" binding:"required_if=Enabled true"`
	Password string `json:"password" mapstructure:"password" example:"your-password" binding:"required_if=Enabled true"`
	SSL      bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}
