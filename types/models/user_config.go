package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// UserConfig represents user-specific configuration preferences
// @Description User-specific configuration stored in the database
type UserConfig struct {
	BaseModel
	// UserID links this config to a specific user
	UserID      uint64 `json:"userId" gorm:"uniqueIndex;not null"`
	DisplayName string `json:"displayName" gorm:"default:''" binding:"omitempty" example:"John Doe"`
	// UI Preferences
	Theme            string `json:"theme" gorm:"default:'system'" example:"dark" binding:"omitempty,oneof=light dark system"`
	Language         string `json:"language" gorm:"default:'en-US'" example:"en-US" binding:"omitempty"`
	ItemsPerPage     int    `json:"itemsPerPage" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
	EnableAnimations bool   `json:"enableAnimations" gorm:"default:true" example:"true"`
	// What sidebar options to show based on the enabled content types.
	ContentTypes string `json:"contentTypes" gorm:"type:text;default:''" example:"movie,series,tv"`

	// Profile settings
	Bio string `json:"bio" gorm:"type:text;default:''" example:"I'm a developer"`

	SocialLinks     SocialLinks     `json:"socialLinks" gorm:"type:jsonb;serializer:json"`
	PrivacySettings PrivacySettings `json:"privacySettings" gorm:"type:jsonb;serializer:json"`

	// Recommendation Preferences
	// Automatically download and organized recommended media into a collection/playlist
	RecommendationSyncEnabled bool `json:"recommendationSyncEnabled" gorm:"default:false" example:"true"`
	// What type of list to create for the synced recommendations
	RecommendationSyncListType string `json:"recommendationSyncListType" gorm:"default:'collection'" example:"playlist" binding:"omitempty,oneof=playlist collection"`
	// How often to sync new recommendations
	RecommendationSyncFrequency string `json:"recommendationSyncFrequency" gorm:"default:'manual'" example:"daily" binding:"omitempty,oneof=manual daily weekly monthly"`
	// Prefux to add to teh beginning of the list name to identify that its part of the auto recommendations system
	RecommendationListPrefix string `json:"recommendationListPrefix" gorm:"default:'Recommendations'" example:"AI Picks" binding:"omitempty"`
	// What content types with be part of the auto sync recommendations
	RecommendationContentTypes   string  `json:"recommendationContentTypes" gorm:"type:text;default:''" example:"movie,series,tv,book"`
	RecommendationMinRating      float32 `json:"recommendationMinRating" gorm:"default:5.0" example:"6" binding:"omitempty"`
	RecommendationMaxAge         int     `json:"recommendationMaxAge" gorm:"default:0" example:"5" binding:"omitempty,min=0,max=100"` // In years, 0 = no limit
	RecommendationIncludeWatched bool    `json:"recommendationIncludeWatched" gorm:"default:false" example:"false"`
	RecommendationIncludeSimilar bool    `json:"recommendationIncludeSimilar" gorm:"default:false" example:"false"`
	// how many movie recommendations to generate
	MaxRecommendations *MaxRecommendations `json:"maxRecommendations" gorm:"type:jsonb;serializer:json"`

	ExcludedKeywords      string `json:"excludedKeywords" gorm:"type:text" example:"war,violence,politics"`
	IncludeUnratedContent bool   `json:"includeUnratedContent" gorm:"default:false" example:"false"`

	ShowAdultContent        bool   `json:"showAdultContent" gorm:"default:false" example:"false"`
	PreferredAudioLanguages string `json:"preferredAudioLanguages" gorm:"type:text;default:''" example:"en,ja"`
	PreferredContentLength  string `json:"preferredContentLength" gorm:"default:'medium'" example:"short" binding:"omitempty,oneof=short medium long"`

	PreferredGenres *Genres `json:"preferredGenres" gorm:"type:jsonb;serializer:json"`
	ExcludedGenres  *Genres `json:"excludedGenres" gorm:"type:jsonb;serializer:json"`

	MinContentRating string `json:"minContentRating" gorm:"default:'G'" example:"PG-13" binding:"omitempty"`
	MaxContentRating string `json:"maxContentRating" gorm:"default:'R'" example:"R" binding:"omitempty"`

	// AI Algorithm Settings
	AiChatPersonality      string  `json:"aiChatPersonality" gorm:"default:'friendly'" example:"serious" binding:"omitempty,oneof=friendly serious enthusiastic analytical custom"`
	RecommendationStrategy string  `json:"recommendationStrategy" gorm:"default:'balanced'" example:"popular" binding:"omitempty,oneof=similar recent popular balanced"`
	NewContentWeight       float32 `json:"newContentWeight" gorm:"default:0.5" example:"0.7" binding:"omitempty,min=0,max=1"`
	PopularityWeight       float32 `json:"popularityWeight" gorm:"default:0.5" example:"0.3" binding:"omitempty,min=0,max=1"`
	PersonalHistoryWeight  float32 `json:"personalHistoryWeight" gorm:"default:0.8" example:"0.8" binding:"omitempty,min=0,max=1"`
	DiscoveryModeEnabled   bool    `json:"discoveryModeEnabled" gorm:"default:false" example:"true"` // Emphasize new content discovery
	DiscoveryModeRatio     float32 `json:"discoveryModeRatio" gorm:"default:0.5" example:"0.5" binding:"omitempty,min=0,max=1"`

	DefaultClients *DefaultClients `json:"defaultClients" gorm:"type:jsonb;serializer:json"`

	// Smart Collections Settings
	SmartCollectionsEnabled bool `json:"smartCollectionsEnabled" gorm:"default:false" example:"true"`

	// Content Availability Settings
	ContentAvailabilityEnabled bool `json:"contentAvailabilityEnabled" gorm:"default:false" example:"true"`

	// New Release Notifications Settings
	NewReleaseNotificationsEnabled bool   `json:"newReleaseNotificationsEnabled" gorm:"default:false" example:"true"`
	NewReleaseMediaTypes           string `json:"newReleaseMediaTypes" gorm:"type:text;default:'movie,series'" example:"movie,series,music"`

	// Playlist Sync Settings
	PlaylistSyncEnabled   bool   `json:"playlistSyncEnabled" gorm:"default:false" example:"true"`
	PlaylistSyncDirection string `json:"playlistSyncDirection" gorm:"default:'bidirectional'" example:"bidirectional"`

	// Activity Analysis Settings
	ActivityAnalysisEnabled bool `json:"activityAnalysisEnabled" gorm:"default:false" example:"true"`

	// Notification Settings
	NotificationsEnabled       bool    `json:"notificationsEnabled" gorm:"default:true" example:"true"`
	NotifyRatingThreshold      float64 `json:"notifyRatingThreshold" gorm:"default:5.0" example:"5" binding:"omitempty,min=0,max=10"`
	NotifyUpcomingReleases     bool    `json:"notifyUpcomingReleases" gorm:"default:true" example:"true"`
	NotifyRecentReleases       bool    `json:"notifyRecentReleases" gorm:"default:true" example:"true"`
	NotifyEmail                bool    `json:"emailNotifications" gorm:"default:false" example:"true"`
	NotifyOnNewRecommendations bool    `json:"notifyOnNewRecommendations" gorm:"default:true" example:"true"`
	NotifyMediaTypes           string  `json:"notifyMediaTypes" gorm:"type:text;default:'movie,series,music'" example:"movie,series,music"`
	NotifyOnSync               bool    `json:"notifyOnSync" gorm:"default:false" example:"false"`
	MaxNotificationsPerDay     int     `json:"maxNotificationsPerDay" gorm:"default:10" example:"10" binding:"omitempty,min=1,max=100"`

	DigestFrequency string `json:"digestFrequency" gorm:"default:'never'" example:"weekly" binding:"omitempty,oneof=never daily weekly"`

	// Onboarding
	OnboardingCompleted bool `json:"onboardingCompleted" gorm:"default:false" example:"true"`
}
type Genres struct {
	Movies []string `json:"movies" gorm:"type:text" example:"action,comedy,drama"`
	Series []string `json:"series" gorm:"type:text" example:"action,comedy,drama"`
	Music  []string `json:"music" gorm:"type:text" example:"electronic,pop,rock"`
	Books  []string `json:"books" gorm:"type:text" example:"fantasy,horror,mystery"`
	Anime  []string `json:"anime" gorm:"type:text" example:"action,comedy,drama"`
	Games  []string `json:"games" gorm:"type:text" example:"action,comedy,drama"`
}

type SocialLinks struct {
	Letterboxd string `json:"letterboxd" gorm:"type:text;default:''" example:"https://letterboxd.com/faiyt"`
	Lastfm     string `json:"lastfm" gorm:"type:text;default:''" example:"https://last.fm/user/faiyt"`
	Trakt      string `json:"trakt" gorm:"type:text;default:''" example:"https://trakt.tv/users/faiyt"`
}

type PrivacySettings struct {
	ShowWatchHistory       bool `json:"showWatchHistory" gorm:"default:true" example:"true"`
	ShareRecommendations   bool `json:"shareRecommendations" gorm:"default:true" example:"true"`
	PublicProfile          bool `json:"publicProfile" gorm:"default:true" example:"true"`
	ShowRecommendationList bool `json:"showRecommendationList" gorm:"default:true" example:"true"`
}

// Scan implements the sql.Scanner interface
func (g *Genres) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &g)
}

// Value implements the driver.Valuer interface
func (g Genres) Value() (driver.Value, error) {
	return json.Marshal(g)
}

type MaxRecommendations struct {
	Movies int `json:"movies" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
	Series int `json:"series" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
	Music  int `json:"music" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
	Anime  int `json:"anime" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
	Games  int `json:"games" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
	Books  int `json:"books" gorm:"default:20" example:"20" binding:"omitempty,min=5,max=100"`
}

func (m *MaxRecommendations) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &m)
}
func (m MaxRecommendations) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type DefaultClients struct {
	VideoClientID uint64 `json:"videoClientId" gorm:"default:0" example:"1"`
	MusicClientID uint64 `json:"musicClientId" gorm:"default:0" example:"1"`
	// TODO: Add support for books. Should it be books vs audiobooks? How would I want to handle that?
	BookClientID       uint64 `json:"bookClientId" gorm:"default:0" example:"1"`
	AIClientID         uint64 `json:"aiClientId" gorm:"default:0" example:"1"`
	MovieAutomationID  uint64 `json:"movieAutomationId" gorm:"default:0" example:"1"`
	SeriesAutomationID uint64 `json:"seriesAutomationId" gorm:"default:0" example:"1"`
	MusicAutomationID  uint64 `json:"musicAutomationId" gorm:"default:0" example:"1"`
	BookAutomationID   uint64 `json:"bookAutomationId" gorm:"default:0" example:"1"`
}

func (d *DefaultClients) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &d)
}
func (d DefaultClients) Value() (driver.Value, error) {
	return json.Marshal(d)
}
