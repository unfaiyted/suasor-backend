package requests

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
