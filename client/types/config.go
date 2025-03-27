package types

type ClientConfig interface {
	isClientConfig()
	GetSpecificType() string
}

type MediaClientConfig interface {
	isMediaClientConfig()
	GetClientType() MediaClientType
}

type AutomationClientConfig interface {
	isAutomationClientConfig()
	GetClientType() AutomationClientType
}

func (PlexConfig) GetSpecificType() string {
	return MediaClientTypePlex.String()
}

func (JellyfinConfig) GetSpecificType() string {
	return MediaClientTypeJellyfin.String()
}

func (EmbyConfig) GetSpecificType() string {
	return MediaClientTypeEmby.String()
}

func (SubsonicConfig) GetSpecificType() string {
	return MediaClientTypeSubsonic.String()
}

func (RadarrConfig) GetSpecificType() string {
	return AutomationClientTypeRadarr.String()
}

func (LidarrConfig) GetSpecificType() string {
	return AutomationClientTypeLidarr.String()
}
func (SonarrConfig) GetSpecificType() string {
	return AutomationClientTypeSonarr.String()
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

func (SonarrConfig) GetClientType() AutomationClientType {
	return AutomationClientTypeSonarr
}
func (RadarrConfig) GetClientType() AutomationClientType {
	return AutomationClientTypeRadarr
}
func (LidarrConfig) GetClientType() AutomationClientType {
	return AutomationClientTypeLidarr
}

func (EmbyConfig) GetClientType() MediaClientType {
	return MediaClientTypeEmby
}
func (JellyfinConfig) GetClientType() MediaClientType {
	return MediaClientTypeJellyfin
}
func (PlexConfig) GetClientType() MediaClientType {
	return MediaClientTypePlex
}
func (SubsonicConfig) GetClientType() MediaClientType {
	return MediaClientTypeSubsonic
}

// @Description Emby media server configuration
type SonarrConfig struct {
	Enabled bool                 `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string               `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string               `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool                 `json:"ssl" mapstructure:"ssl" example:"false"`
	Type    AutomationClientType `json:"type" mapstructure:"type" default:"sonarr"`
}

func NewSonarrConfig() SonarrConfig {
	return SonarrConfig{
		Type: AutomationClientTypeSonarr,
	}
}

// @Description Emby media server configuration
type RadarrConfig struct {
	Enabled bool                 `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string               `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string               `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool                 `json:"ssl" mapstructure:"ssl" example:"false"`
	Type    AutomationClientType `json:"type" mapstructure:"type" default:"radarr"`
}

func NewRadarrConfig() RadarrConfig {
	return RadarrConfig{
		Type: AutomationClientTypeRadarr,
	}
}

// @Description Jellyfin media server configuration
type LidarrConfig struct {
	Enabled bool                 `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string               `json:"baseURL" mapstructure:"host" example:"http://localhost:8096" binding:"required_if=Enabled true"`
	APIKey  string               `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool                 `json:"ssl" mapstructure:"ssl" example:"false"`
	Type    AutomationClientType `json:"type" mapstructure:"type" default:"lidarr"`
}

func NewLidarrConfig() LidarrConfig {
	return LidarrConfig{
		Type: AutomationClientTypeLidarr,
	}
}

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

// @Description Plex media server configuration
type PlexConfig struct {
	Enabled bool            `json:"enabled" mapstructure:"enabled" example:"false"`
	Host    string          `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port    int             `json:"port" mapstructure:"port" example:"32400" binding:"required_if=Enabled true"`
	Token   string          `json:"token" mapstructure:"token" example:"your-plex-token" binding:"required_if=Enabled true"`
	SSL     bool            `json:"ssl" mapstructure:"ssl" example:"false"`
	Type    MediaClientType `json:"type" mapstructure:"type" default:"plex"`
}

func NewPlexConfig() PlexConfig {
	return PlexConfig{
		Type: MediaClientTypePlex,
	}
}

// @Description Supersonic music server configuration
type SubsonicConfig struct {
	Enabled  bool            `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string          `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int             `json:"port" mapstructure:"port" example:"4533" binding:"required_if=Enabled true"`
	Username string          `json:"username" mapstructure:"username" example:"admin" binding:"required_if=Enabled true"`
	Password string          `json:"password" mapstructure:"password" example:"your-password" binding:"required_if=Enabled true"`
	SSL      bool            `json:"ssl" mapstructure:"ssl" example:"false"`
	Type     MediaClientType `json:"type" mapstructure:"type" default:"subsonic"`
}

func NewSubsonicConfig() SubsonicConfig {
	return SubsonicConfig{
		Type: MediaClientTypeSubsonic,
	}
}
