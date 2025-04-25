package requests

// RegisterRequest contains user registration information
//	@Description	Request payload for user registration
type RegisterRequest struct {
	// Email is the user's email address
	//	@Description	User's email address
	//	@Example		"user@example.com"
	Email string `json:"email" binding:"required,email" example:"user@example.com"`

	// Username is the user's desired username
	//	@Description	User's chosen username
	//	@Example		"johndoe"
	Username string `json:"username" binding:"required" example:"johndoe"`

	// Password is the user's chosen password
	//	@Description	User's password (plain text in request)
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest contains user  login information
//	@Description	Request payload for user login
type LoginRequest struct {
	// Email is the user's email address
	//	@Description	User's email address
	//	@Example		"user@example.com"
	Email string `json:"email" binding:"required,email" example:"user@example.com"`

	// Password is the user's chosen password
	//	@Description	User's password (plain text in request)
	Password string `json:"password" binding:"required,min=8"`
}

// LogoutRequest contains the refresh token for session termination
//	@Description	Request payload for user logout
type LogoutRequest struct {
	// RefreshToken identifies the session to terminate
	//	@Description	JWT refresh token to invalidate
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenRequest contains the refresh token for obtaining new access tokens
//	@Description	Request payload for refreshing access tokens
type RefreshTokenRequest struct {
	// RefreshToken is used to generate a new access token
	//	@Description	JWT refresh token to use for generating new access token
	RefreshToken string `json:"refreshToken" binding:"required"`
}
