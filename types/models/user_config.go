package models

// UserConfig represents user-specific configuration preferences
// @Description User-specific configuration stored in the database
type UserConfig struct {
	BaseModel
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
