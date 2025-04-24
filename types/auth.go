package types

import (
	"github.com/golang-jwt/jwt/v4"
)

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
	UserID uint64 `json:"userID"`
	UUID   string `json:"uuid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
