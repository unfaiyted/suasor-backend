package responses

// @Description	Authentication data returned to client after successful authentication
type AuthDataResponse struct {
	// User contains the user profile information
	//	@Description	User profile data
	User UserResponse `json:"user"`

	// AccessToken is the JWT token for API access
	//	@Description	JWT access token for authenticated requests
	//	@Example		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	AccessToken string `json:"accessToken"`

	// RefreshToken is used to get new access tokens
	//	@Description	JWT refresh token for obtaining new access tokens
	//	@Example		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	RefreshToken string `json:"refreshToken"`

	// ExpiresAt indicates when the access token expires
	//	@Description	UNIX timestamp when the access token expires
	//	@Example		1674140400
	ExpiresAt int64 `json:"expiresAt"`
}

type SystemStatusResponse struct {
	Status string `json:"status"`
}
