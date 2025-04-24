package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"suasor/services"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"
)

// AuthHandler handles authentication-related endpoints
type AuthHandler struct {
	service services.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with email, username, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.RegisterRequest true "User registration data"
// @Example request
//
//	{
//	  "email": "user@example.com",
//	  "username": "johndoe",
//	  "password": "securePassword123"
//	}
//
// @Success 201 {object} responses.APIResponse[responses.AuthDataResponse] "Successfully registered user"
// @Example response
//
//	{
//	  "success": true,
//	  "data": {
//	    "user": {
//	      "id": 1,
//	      "email": "user@example.com",
//	      "username": "johndoe",
//	      "role": "user"
//	    },
//	    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "expires_at": 1625097600
//	  },
//	  "message": "User registered successfully"
//	}
//
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 409 {object} responses.ErrorResponse[responses.ErrorDetails] "Email or username already in use"
// @Example response
//
//	{
//	  "error": "conflict",
//	  "message": "Email already registered",
//	  "details": {
//	    "error": "email already exists"
//	  },
//	  "timestamp": "2025-03-16T10:30:45Z",
//	  "requestID": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.RegisterRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().
		Str("email", req.Email).
		Str("username", req.Username).
		Msg("Registering new user")

	result, err := h.service.Register(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Warn().Err(err).Str("email", req.Email).Str("username", req.Username).Msg("Registration conflict")
			responses.RespondConflict(c, err, err.Error())
			return
		}

		log.Error().Err(err).Str("email", req.Email).Msg("Failed to register user")
		responses.RespondInternalError(c, err, "Failed to register user")
		return
	}

	log.Info().Uint64("userID", result.User.ID).Str("email", result.User.Email).Msg("Successfully registered user")

	responses.RespondCreated(c, result, "User registered successfully")
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.LoginRequest true "User login data"
// @Example request
//
//	{
//	  "email": "user@example.com",
//	  "password": "securePassword123"
//	}
//
// @Success 200 {object} responses.APIResponse[responses.AuthDataResponse] "Successfully authenticated user"
// @Example response
//
//	{
//	  "success": true,
//	  "data": {
//	    "user": {
//	      "id": 1,
//	      "email": "user@example.com",
//	      "username": "johndoe",
//	      "role": "user"
//	    },
//	    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "expires_at": 1625097600
//	  },
//	  "message": "Login successful"
//	}
//
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid credentials or inactive account"
// @Example response
//
//	{
//	  "error": "unauthorized",
//	  "message": "Invalid email or password",
//	  "details": {},
//	  "timestamp": "2025-03-16T10:30:45Z",
//	  "requestID": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.LoginRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().Str("email", req.Email).Msg("User attempting to login")

	result, err := h.service.Login(ctx, req)
	if err != nil {
		if checkErrorType(c, err, "Login error") {
			return
		}

		log.Error().Err(err).Str("email", req.Email).Msg("Login failed - server error")
		responses.RespondInternalError(c, err, "Failed to authenticate user")
		return
	}

	log.Info().Uint64("userID", result.User.ID).Str("email", result.User.Email).Msg("User successfully logged in")

	responses.RespondOK(c, result, "Login successful")
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.RefreshTokenRequest true "Refresh token data"
// @Example request
//
//	{
//	  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	}
//
// @Success 200 {object} responses.APIResponse[responses.AuthDataResponse] "Successfully refreshed token"
// @Example response
//
//	{
//	  "success": true,
//	  "data": {
//	    "user": {
//	      "id": 1,
//	      "email": "user@example.com",
//	      "username": "johndoe",
//	      "role": "user"
//	    },
//	    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//	    "expires_at": 1625097600
//	  },
//	  "message": "Token refreshed successfully"
//	}
//
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid or expired refresh token"
// @Example response
//
//	{
//	  "error": "unauthorized",
//	  "message": "Invalid refresh token",
//	  "details": {},
//	  "timestamp": "2025-03-16T10:30:45Z",
//	  "requestID": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.RefreshTokenRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().Msg("Token refresh requested")

	result, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if checkErrorType(c, err, "Token refresh error") {
			return
		}

		log.Error().Err(err).Msg("Token refresh failed - server error")
		responses.RespondInternalError(c, err, "Failed to refresh token")
		return
	}

	log.Info().Uint64("userID", result.User.ID).Msg("Token refreshed successfully")

	responses.RespondOK(c, result, "Token refreshed successfully")
}

// Logout godoc
// @Summary Log out a user
// @Description Invalidates the refresh token, effectively logging the user out
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.LogoutRequest true "Logout data"
// @Example request
//
//	{
//	  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	}
//
// @Success 200 {object} responses.APIResponse[any] "Successfully logged out"
// @Example response
//
//	{
//	  "success": true,
//	  "message": "Successfully logged out"
//	}
//
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid refresh token"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.LogoutRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().Msg("User logout requested")

	err := h.service.Logout(ctx, req.RefreshToken)
	if err != nil {
		if checkErrorType(c, err, "Logout error") {
			return
		}

		log.Error().Err(err).Msg("Logout failed - server error")
		responses.RespondInternalError(c, err, "Failed to log out")
		return
	}

	log.Info().Msg("User successfully logged out")

	c.JSON(http.StatusOK, responses.APIResponse[any]{
		Success: true,
		Message: "Successfully logged out",
	})
}

// ValidateSession godoc
// @Summary Validate user session
// @Description Validates the user's session token and returns current user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.APIResponse[responses.UserResponse] "Valid session with user details"
// @Example response
//
//	{
//	  "success": true,
//	  "data": {
//	    "id": 1,
//	    "email": "user@example.com",
//	    "username": "johndoe",
//	    "role": "user"
//	  },
//	  "message": "Session is valid"
//	}
//
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid or expired session token"
// @Example response
//
//	{
//	  "error": "unauthorized",
//	  "message": "Invalid or expired session token",
//	  "details": {},
//	  "timestamp": "2025-03-16T10:30:45Z",
//	  "requestID": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/auth/validate [get]
func (h *AuthHandler) ValidateSession(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	token, ok := extractToken(c)
	if !ok {
		return
	}

	log.Info().Msg("Validating user session")

	// Validate the token and get the user data
	_, err := h.service.ValidateToken(ctx, token)
	if err != nil {
		log.Warn().Err(err).Msg("Session validation failed")
		responses.RespondUnauthorized(c, err, "Invalid or expired session token")
		return
	}

	user, err := h.service.GetAuthorizedUser(ctx, token)
	if err != nil {
		log.Warn().Err(err).Msg("Getting Authorize User failed")
		responses.RespondInternalError(c, err, "Unable to get authorized user information")
	}
	log.Info().Uint64("userID", user.ID).Str("email", user.Email).Msg("Session validated successfully")

	// Return the user profile
	responses.RespondOK(c, user, "Session is valid")
}
