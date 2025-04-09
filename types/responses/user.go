package responses

// UserResponse represents the user data returned in API responses
// @Description User information returned in API responses
type UserResponse struct {
	// ID is the unique identifier for the user
	// @Description User's unique identifier
	// @Example 1
	ID uint64 `json:"id"`

	// Email is the unique email address of the user
	// @Description User's email address
	// @Example "user@example.com"
	Email string `json:"email"`

	// Username is the display name chosen by the user
	// @Description User's chosen username
	// @Example "johndoe"
	Username string `json:"username"`

	// Avatar is the path to the user's avatar image
	// @Description Path to the user's avatar image
	// @Example "/uploads/avatars/user_1.jpg"
	Avatar string `json:"avatar,omitempty" example:"/uploads/avatars/user_1.jpg"`

	// Role defines the user's permission level
	// @Description User's role in the system
	// @Enum "user" "admin"
	// @Example "user"
	Role string `json:"role"`
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
