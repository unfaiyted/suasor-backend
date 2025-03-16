package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user account in the system
// @Description User account information
type User struct {
	gorm.Model
	// Email is the unique identifier for the user
	// @Description User's email address
	// @Example "user@example.com"
	Email string `json:"email" gorm:"uniqueIndex;not null" binding:"required,email" example:"user@example.com"`

	// Username is an optional display name
	// @Description User's chosen username
	// @Example "johndoe"
	Username string `json:"username" gorm:"uniqueIndex;not null" binding:"required" example:"johndoe"`

	// Password is stored as a bcrypt hash
	// @Description Hashed password (never returned in responses)
	Password string `json:"-" gorm:"not null" binding:"required"`

	// Role defines the user's permission level
	// @Description User's role in the system
	// @Enum "user" "admin"
	// @Example "user"
	Role string `json:"role" gorm:"default:'user'" example:"user"`

	// Active indicates if the account is currently active
	// @Description Whether the user account is active
	// @Example true
	Active bool `json:"active" gorm:"default:true" example:"true"`

	// LastLogin tracks the most recent login time
	// @Description Time of the user's most recent login
	LastLogin *time.Time `json:"lastLogin,omitempty"`

	// Sessions holds all active sessions for this user
	Sessions []Session `json:"-" gorm:"foreignKey:UserID"`
}

// UserResponse represents the user data returned in API responses
// @Description User information returned in API responses
type UserResponse struct {
	// ID is the unique identifier for the user
	// @Description User's unique identifier
	// @Example 1
	ID uint `json:"id"`

	// Email is the unique email address of the user
	// @Description User's email address
	// @Example "user@example.com"
	Email string `json:"email"`

	// Username is the display name chosen by the user
	// @Description User's chosen username
	// @Example "johndoe"
	Username string `json:"username"`

	// Role defines the user's permission level
	// @Description User's role in the system
	// @Enum "user" "admin"
	// @Example "user"
	Role string `json:"role"`
}

// UpdateProfileRequest represents the data for updating a user's profile information
// @Description Request payload for updating user profile
type UpdateProfileRequest struct {
	// Email is the new email address
	// @Description Updated email address for the user
	// @Example "newemail@example.com"
	Email string `json:"email" binding:"omitempty,email" example:"newemail@example.com"`

	// Username is the new username
	// @Description Updated username for the user
	// @Example "newusername"
	Username string `json:"username" binding:"omitempty" example:"newusername"`
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
