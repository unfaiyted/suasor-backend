package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the user account in the system
//
//	@Description	User account information
type User struct {
	BaseModel
	// Email is the unique identifier for the user
	//	@Description	User's email address
	//	@Example		"user@example.com"
	Email string `json:"email" gorm:"uniqueIndex;not null" binding:"required,email" example:"user@example.com"`

	// Username is an optional display name
	//	@Description	User's chosen username
	//	@Example		"johndoe"
	Username string `json:"username" gorm:"uniqueIndex;not null" binding:"required" example:"johndoe"`

	// Password is stored as a bcrypt hash
	//	@Description	Hashed password (never returned in responses)
	Password string `json:"-" gorm:"not null" binding:"required"`

	// Avatar is the file path to the user's avatar image
	//	@Description	Path to the user's avatar image
	//	@Example		"/uploads/avatars/user_1.jpg"
	Avatar string `json:"avatar" gorm:"default:''" example:"/uploads/avatars/user_1.jpg"`

	// Role defines the user's permission level
	//	@Description	User's role in the system
	//	@Enum			"user" "admin"
	//	@Example		"user"
	Role string `json:"role" gorm:"default:'user'" example:"user"`

	// Active indicates if the account is currently active
	//	@Description	Whether the user account is active
	//	@Example		true
	Active bool `json:"active" gorm:"default:true" example:"true"`

	// LastLogin tracks the most recent login time
	//	@Description	Time of the user's most recent login
	LastLogin *time.Time `json:"lastLogin,omitempty"`

	// Sessions holds all active sessions for this user
	Sessions []Session `json:"-" gorm:"foreignKey:UserID"`

	// PasswordResetToken is used to reset a user's password
	//	@Description	Password reset token (optional)
	//	@Example		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	PasswordResetToken string `json:"passwordResetToken,omitempty" gorm:"size:1024"`
}

// SetPassword hashes and sets the user's password
//
//	@Description	Sets a bcrypt-hashed password for the user
//	@Return			error If password hashing fails
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the provided password against the stored hash
//
//	@Description	Checks if the provided password matches the stored hash
//	@Return			bool True if password matches
//	@Return			error If password checking fails
func (u *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// BeforeCreate is a GORM hook that runs before creating a new user
//
//	@Description	GORM hook that runs before creating a user record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Additional validation or preparation could be added here
	return nil
}
