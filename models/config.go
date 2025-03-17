// models/config.go
package models

import (
	"gorm.io/gorm"
)

// Configuration represents the complete application configuration
// @Description Complete application configuration settings
type Configuration struct {
	// App contains core application settings
	App struct {
		Name        string `json:"name" mapstructure:"name" example:"suasor" binding:"required"`
		Environment string `json:"environment" mapstructure:"environment" example:"development" binding:"required,oneof=development staging production"`
		AppURL      string `json:"appURL" mapstructure:"appURL" example:"http://localhost:3000" binding:"required,url"`
		APIBaseURL  string `json:"apiBaseURL" mapstructure:"apiBaseURL" example:"http://localhost:8080" binding:"required,url"`
		LogLevel    string `json:"logLevel" mapstructure:"logLevel" example:"info" binding:"required,oneof=debug info warn error"`
		MaxPageSize int    `json:"maxPageSize" mapstructure:"maxPageSize" example:"100" binding:"required,min=1,max=1000"`
	} `json:"app"`

	// Database contains database connection settings
	Db struct {
		Host     string `json:"host" mapstructure:"url" example:"localhost" binding:"required"`
		Port     string `json:"port" mapstructure:"port" example:"5432" binding:"required"`
		Name     string `json:"name" mapstructure:"name" example:"suasor" binding:"required"`
		User     string `json:"user" mapstructure:"user" example:"postgres_user" binding:"required"`
		Password string `json:"password" mapstructure:"password" example:"yourpassword" binding:"required"`
		MaxConns int    `json:"maxConns" mapstructure:"maxConns" example:"20" binding:"required,min=1"`
		Timeout  int    `json:"timeout" mapstructure:"timeout" example:"30" binding:"required,min=1"`
	} `json:"db" mapstructure:"db"`

	// HTTP contains HTTP server configuration
	HTTP struct {
		Port             string `json:"port" mapstructure:"port" example:"8080" binding:"required"`
		ReadTimeout      int    `json:"readTimeout" mapstructure:"readTimeout" example:"30" binding:"required,min=1"`
		WriteTimeout     int    `json:"writeTimeout" mapstructure:"writeTimeout" example:"30" binding:"required,min=1"`
		IdleTimeout      int    `json:"idleTimeout" mapstructure:"idleTimeout" example:"60" binding:"required,min=1"`
		EnableSSL        bool   `json:"enableSSL" mapstructure:"enableSSL" example:"false"`
		SSLCert          string `json:"sslCert" mapstructure:"sslCert" example:"/path/to/cert.pem"`
		SSLKey           string `json:"sslKey" mapstructure:"sslKey" example:"/path/to/key.pem"`
		ProxyEnabled     bool   `json:"proxyEnabled" mapstructure:"proxyEnabled" example:"false"`
		ProxyURL         string `json:"proxyURL" mapstructure:"proxyURL" example:"http://proxy:8080"`
		RateLimitEnabled bool   `json:"rateLimitEnabled" mapstructure:"rateLimitEnabled" example:"true"`
		RequestsPerMin   int    `json:"requestsPerMin" mapstructure:"requestsPerMin" example:"100" binding:"min=0"`
	} `json:"http"`

	// Auth contains authentication settings
	Auth struct {
		EnableLocal     bool     `json:"enableLocal" mapstructure:"enableLocal" example:"true"`
		SessionTimeout  int      `json:"sessionTimeout" mapstructure:"sessionTimeout" example:"60" binding:"required,min=1"`
		Enable2FA       bool     `json:"enable2FA" mapstructure:"enable2FA" example:"false"`
		JWTSecret       string   `json:"jwtSecret" mapstructure:"jwtSecret" example:"your-secret-key" binding:"required"`
		TokenExpiration int      `json:"tokenExpiration" mapstructure:"tokenExpiration" example:"24" binding:"required,min=1"`
		AllowedOrigins  []string `json:"allowedOrigins" mapstructure:"allowedOrigins" example:"http://localhost:3000"`
		// New fields to add
		AccessExpiryMinutes int    `json:"accessExpiryMinutes" mapstructure:"accessExpiryMinutes" example:"15" binding:"required,min=1"`
		RefreshExpiryDays   int    `json:"refreshExpiryDays" mapstructure:"refreshExpiryDays" example:"7" binding:"required,min=1"`
		TokenIssuer         string `json:"tokenIssuer" mapstructure:"tokenIssuer" example:"suasor-api" binding:"required"`
		TokenAudience       string `json:"tokenAudience" mapstructure:"tokenAudience" example:"suasor-client" binding:"required"`
	} `json:"auth"`

	// Integrations contains all third-party service configurations
	Integrations struct {
		Emby      EmbyConfig      `json:"emby" mapstructure:"emby"`
		Jellyfin  JellyfinConfig  `json:"jellyfin" mapstructure:"jellyfin"`
		Plex      PlexConfig      `json:"plex" mapstructure:"plex"`
		Trakt     TraktConfig     `json:"trakt" mapstructure:"trakt"`
		Navidrome NavidromeConfig `json:"navidrome" mapstructure:"navidrome"`
		Spotify   SpotifyConfig   `json:"spotify" mapstructure:"spotify"`
	} `json:"integrations"`
}

// Integration config types
// @Description Emby media server configuration
type EmbyConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int    `json:"port" mapstructure:"port" example:"8096" binding:"required_if=Enabled true"`
	APIKey   string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
	SSL      bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Jellyfin media server configuration
type JellyfinConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int    `json:"port" mapstructure:"port" example:"8096" binding:"required_if=Enabled true"`
	APIKey   string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
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

// @Description Trakt.tv configuration
type TraktConfig struct {
	Enabled      bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	ClientID     string `json:"clientId" mapstructure:"clientId" example:"your-client-id" binding:"required_if=Enabled true"`
	ClientSecret string `json:"clientSecret" mapstructure:"clientSecret" example:"your-client-secret" binding:"required_if=Enabled true"`
	RedirectURI  string `json:"redirectUri" mapstructure:"redirectUri" example:"http://localhost:8080/callback" binding:"required_if=Enabled true"`
}

// @Description Supersonic music server configuration
type NavidromeConfig struct {
	Enabled  bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int    `json:"port" mapstructure:"port" example:"4533" binding:"required_if=Enabled true"`
	Username string `json:"username" mapstructure:"username" example:"admin" binding:"required_if=Enabled true"`
	Password string `json:"password" mapstructure:"password" example:"your-password" binding:"required_if=Enabled true"`
	SSL      bool   `json:"ssl" mapstructure:"ssl" example:"false"`
}

// @Description Spotify configuration
type SpotifyConfig struct {
	Enabled      bool   `json:"enabled" mapstructure:"enabled" example:"false"`
	ClientID     string `json:"clientId" mapstructure:"clientId" example:"your-client-id" binding:"required_if=Enabled true"`
	ClientSecret string `json:"clientSecret" mapstructure:"clientSecret" example:"your-client-secret" binding:"required_if=Enabled true"`
	RedirectURI  string `json:"redirectUri" mapstructure:"redirectUri" example:"http://localhost:8080/callback" binding:"required_if=Enabled true"`
	Scopes       string `json:"scopes" mapstructure:"scopes" example:"user-library-read playlist-read-private"`
}

// ConfigResponse represents the response structure for configuration endpoints
// @Description Configuration response wrapper
type ConfigResponse struct {
	Data  *Configuration `json:"data,omitempty"`
	Error string         `json:"error,omitempty"`
}

// UserConfig represents user-specific configuration preferences
// @Description User-specific configuration stored in the database
type UserConfig struct {
	gorm.Model
	// UserID links this config to a specific user
	UserID uint64 `json:"userId" gorm:"uniqueIndex;not null"`
	User   User   `json:"-" gorm:"foreignKey:UserID"`

	// UI Preferences
	Theme            string `json:"theme" gorm:"default:'system'" example:"dark" binding:"omitempty,oneof=light dark system"`
	Language         string `json:"language" gorm:"default:'en-US'" example:"en-US" binding:"required"`
	ItemsPerPage     int    `json:"itemsPerPage" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
	EnableAnimations bool   `json:"enableAnimations" gorm:"default:true" example:"true"`

	// Recommendation Preferences
	RecommendationFrequency string `json:"recommendationFrequency" gorm:"default:'weekly'" example:"daily" binding:"omitempty,oneof=daily weekly monthly"`
	MaxRecommendations      int    `json:"maxRecommendations" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
	IncludeUnratedContent   bool   `json:"includeUnratedContent" gorm:"default:false" example:"false"`
	PreferredGenres         string `json:"preferredGenres" gorm:"type:text" example:"action,comedy,drama"`
	ExcludedGenres          string `json:"excludedGenres" gorm:"type:text" example:"horror,war"`
	MinContentRating        string `json:"minContentRating" gorm:"default:'G'" example:"PG-13" binding:"omitempty"`
	MaxContentRating        string `json:"maxContentRating" gorm:"default:'R'" example:"R" binding:"omitempty"`
	ContentLanguages        string `json:"contentLanguages" gorm:"type:text;default:'en'" example:"en,es,fr"`

	// AI Algorithm Settings
	RecommendationStrategy string  `json:"recommendationStrategy" gorm:"default:'balanced'" example:"diverse" binding:"omitempty,oneof=similar diverse balanced"`
	NewContentWeight       float32 `json:"newContentWeight" gorm:"default:0.5" example:"0.7" binding:"omitempty,min=0,max=1"`
	PopularityWeight       float32 `json:"popularityWeight" gorm:"default:0.5" example:"0.3" binding:"omitempty,min=0,max=1"`
	PersonalHistoryWeight  float32 `json:"personalHistoryWeight" gorm:"default:0.8" example:"0.8" binding:"omitempty,min=0,max=1"`
	EnableExperimentalAI   bool    `json:"enableExperimentalAI" gorm:"default:false" example:"false"`

	// Sync Preferences
	DefaultMediaServer      string `json:"defaultMediaServer" gorm:"default:''" example:"plex" binding:"omitempty,oneof= emby jellyfin plex"`
	AutoSyncRecommendations bool   `json:"autoSyncRecommendations" gorm:"default:false" example:"true"`
	CreateServerPlaylists   bool   `json:"createServerPlaylists" gorm:"default:false" example:"true"`
	SyncFrequency           string `json:"syncFrequency" gorm:"default:'manual'" example:"daily" binding:"omitempty,oneof=manual daily weekly"`
	DefaultCollection       string `json:"defaultCollection" gorm:"default:'Recommendations'" example:"AI Picks" binding:"omitempty"`

	// Notification Settings
	NotificationsEnabled       bool   `json:"notificationsEnabled" gorm:"default:true" example:"true"`
	EmailNotifications         bool   `json:"emailNotifications" gorm:"default:false" example:"true"`
	NotifyOnNewRecommendations bool   `json:"notifyOnNewRecommendations" gorm:"default:true" example:"true"`
	NotifyOnSync               bool   `json:"notifyOnSync" gorm:"default:false" example:"false"`
	DigestFrequency            string `json:"digestFrequency" gorm:"default:'never'" example:"weekly" binding:"omitempty,oneof=never daily weekly"`
}

// UserConfigResponse represents the user config data returned in API responses
// @Description User configuration information returned in API responses
type UserConfigResponse struct {
	ID                     uint64 `json:"id"`
	Theme                  string `json:"theme"`
	Language               string `json:"language"`
	RecommendationStrategy string `json:"recommendationStrategy"`
	DefaultMediaServer     string `json:"defaultMediaServer"`
	NotificationsEnabled   bool   `json:"notificationsEnabled"`
	// Include other fields that should be returned in responses
}

// UpdateUserConfigRequest represents the data for updating user configuration
// @Description Request payload for updating user configuration
type UpdateUserConfigRequest struct {
	Theme                  *string `json:"theme" binding:"omitempty,oneof=light dark system"`
	Language               *string `json:"language" binding:"omitempty"`
	RecommendationStrategy *string `json:"recommendationStrategy" binding:"omitempty,oneof=similar diverse balanced"`
	PreferredGenres        *string `json:"preferredGenres" binding:"omitempty"`
	ExcludedGenres         *string `json:"excludedGenres" binding:"omitempty"`
	DefaultMediaServer     *string `json:"defaultMediaServer" binding:"omitempty,oneof= emby jellyfin plex"`
	NotificationsEnabled   *bool   `json:"notificationsEnabled" binding:"omitempty"`
	// more
}
