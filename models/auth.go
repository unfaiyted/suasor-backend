package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Session represents an active user session
// @Description Active user login session
type Session struct {
	gorm.Model

	// UserID is the foreign key to the user
	// @Description ID of the user this session belongs to
	UserID uint64 `json:"userId" gorm:"not null"`

	// RefreshToken is used to generate new access tokens
	// @Description Refresh token value (hashed in database)
	RefreshToken string `json:"-" gorm:"not null"`

	// UserAgent records the client that created this session
	// @Description Browser/client user agent string
	// @Example "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	UserAgent string `json:"userAgent" gorm:"not null" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"`

	// IP records the IP address of the client
	// @Description IP address of the client
	// @Example "192.168.1.1"
	IP string `json:"ip" gorm:"not null" example:"192.168.1.1"`

	// ExpiresAt indicates when this session should be invalidated
	// @Description When this session expires
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`

	// LastUsed tracks the most recent usage of this session
	// @Description Time the session was last used for authentication
	LastUsed time.Time `json:"lastUsed" gorm:"not null"`
}

// TokenDetails contains the metadata of both JWT token types
// @Description JWT token details
type TokenDetails struct {
	// AccessToken is the JWT token for API access
	// @Description JWT access token
	AccessToken string `json:"accessToken"`

	// RefreshToken is used to get new access tokens
	// @Description JWT refresh token
	RefreshToken string `json:"refreshToken"`

	// AccessUUID is the unique identifier for this access token
	// @Description Unique identifier for the access token
	AccessUUID string `json:"-"`

	// RefreshUUID is the unique identifier for this refresh token
	// @Description Unique identifier for the refresh token
	RefreshUUID string `json:"-"`

	// AtExpires is when the access token expires
	// @Description Access token expiration timestamp
	AtExpires int64 `json:"expiresAt"`

	// RtExpires is when the refresh token expires
	// @Description Refresh token expiration timestamp
	RtExpires int64 `json:"-"`
}

// JWTClaim defines the structure of the JWT claim
// @Description JWT claim st eructure
type JWTClaim struct {
	UserID uint64 `json:"userId"`
	UUID   string `json:"uuid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RegisterRequest contains user registration information
// @Description Request payload for user registration
type RegisterRequest struct {
	// Email is the user's email address
	// @Description User's email address
	// @Example "user@example.com"
	Email string `json:"email" binding:"required,email" example:"user@example.com"`

	// Username is the user's desired username
	// @Description User's chosen username
	// @Example "johndoe"
	Username string `json:"username" binding:"required" example:"johndoe"`

	// Password is the user's chosen password
	// @Description User's password (plain text in request)
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest contains user  login information
// @Description Request payload for user login
type LoginRequest struct {
	// Email is the user's email address
	// @Description User's email address
	// @Example "user@example.com"
	Email string `json:"email" binding:"required,email" example:"user@example.com"`

	// Password is the user's chosen password
	// @Description User's password (plain text in request)
	Password string `json:"password" binding:"required,min=8"`
}

// LogoutRequest contains the refresh token for session termination
// @Description Request payload for user logout
type LogoutRequest struct {
	// RefreshToken identifies the session to terminate
	// @Description JWT refresh token to invalidate
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenRequest contains the refresh token for obtaining new access tokens
// @Description Request payload for refreshing access tokens
type RefreshTokenRequest struct {
	// RefreshToken is used to generate a new access token
	// @Description JWT refresh token to use for generating new access token
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// AuthData contains authentication data returned after successful login/registration
// @Description Authentication data returned to client after successful authentication
type AuthData struct {
	// User contains the user profile information
	// @Description User profile data
	User UserResponse `json:"user"`

	// AccessToken is the JWT token for API access
	// @Description JWT access token for authenticated requests
	// @Example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	AccessToken string `json:"accessToken"`

	// RefreshToken is used to get new access tokens
	// @Description JWT refresh token for obtaining new access tokens
	// @Example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	RefreshToken string `json:"refreshToken"`

	// ExpiresAt indicates when the access token expires
	// @Description UNIX timestamp when the access token expires
	// @Example 1674140400
	ExpiresAt int64 `json:"expiresAt"`
}

// SetPassword hashes and sets the user's password
// @Description Sets a bcrypt-hashed password for the user
// @Return error If password hashing fails
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the provided password against the stored hash
// @Description Checks if the provided password matches the stored hash
// @Return bool True if password matches
// @Return error If password checking fails
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
// @Description GORM hook that runs before creating a user record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Additional validation or preparation could be added here
	return nil
}
