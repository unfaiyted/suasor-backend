package requests

// UpdateProfileRequest represents the data for updating a user's profile information
// @Description Request payload for updating user profile
type ProfileRequest struct {
	// Email is the new email address
	// @Description Updated email address for the user
	// @Example "newemail@example.com"
	Email string `json:"email" binding:"omitempty,email" example:"newemail@example.com"`

	// Username is the new username
	// @Description Updated username for the user
	// @Example "newusername"
	Username string `json:"username" binding:"omitempty" example:"newusername"`

	// Avatar is the path to the user's avatar image
	// @Description Path to the user's avatar image
	// @Example "/uploads/avatars/user_1.jpg"
	Avatar string `json:"avatar" binding:"omitempty" example:"/uploads/avatars/user_1.jpg"`
}

// ChangePasswordRequest represents the data needed to change a user's password
// @Description Request payload for changing user password
type ChangePasswordRequest struct {
	// CurrentPassword is the user's existing password for verification
	// @Description User's current password for verification
	// @Example "oldpassword123"
	CurrentPassword string `json:"currentPassword" binding:"required" example:"oldpassword123"`

	// NewPassword is the password to change to
	// @Description New password to set for the user
	// @Example "newpassword456"
	NewPassword string `json:"newPassword" binding:"required" example:"newpassword456"`
}

// ChangeRoleRequest represents the data needed to change a user's role
// @Description Request payload for changing user role
type ChangeRoleRequest struct {
	// Role is the new role to assign to the user
	// @Description New role to assign to the user
	// @Enum "user" "admin"
	// @Example "admin"
	Role string `json:"role" binding:"required,oneof=user admin" example:"admin"`
}

// UserConfigRequest represents the data for updating user configuration
// @Description Request payload for updating user configuration
type UserConfigRequest struct {
	Theme                  *string `json:"theme" binding:"omitempty,oneof=light dark system"`
	Language               *string `json:"language" binding:"omitempty"`
	RecommendationStrategy *string `json:"recommendationStrategy" binding:"omitempty,oneof=similar diverse balanced"`
	PreferredGenres        *string `json:"preferredGenres" binding:"omitempty"`
	ExcludedGenres         *string `json:"excludedGenres" binding:"omitempty"`
	DefaultMediaServer     *string `json:"defaultMediaServer" binding:"omitempty,oneof= emby jellyfin plex"`
	NotificationsEnabled   *bool   `json:"notificationsEnabled" binding:"omitempty"`
	// more
}

// AvatarUploadResponse represents the response after a successful avatar upload
// @Description Response data after avatar upload
type AvatarUploadResponse struct {
	// FilePath is the path to the uploaded avatar file
	// @Description Path to the uploaded avatar file
	// @Example "/uploads/avatars/user_1.jpg"
	FilePath string `json:"filePath" example:"/uploads/avatars/user_1.jpg"`
}
