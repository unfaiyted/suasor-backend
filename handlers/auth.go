package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"suasor/models"
	"suasor/services"
	"suasor/utils"
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
// @Param request body models.RegisterRequest true "User registration data"
// @Example request
//
//	{
//	  "email": "user@example.com",
//	  "username": "johndoe",
//	  "password": "securePassword123"
//	}
//
// @Success 201 {object} models.APIResponse[models.AuthData] "Successfully registered user"
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
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format"
// @Failure 409 {object} models.ErrorResponse[error] "Email or username already in use"
// @Example response
//
//	{
//	  "error": "conflict",
//	  "message": "Email already registered",
//	  "details": {
//	    "error": "email already exists"
//	  },
//	  "timestamp": "2025-03-16T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for user registration")
		utils.RespondValidationError(c, err)
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
			utils.RespondConflict(c, err, err.Error())
			return
		}

		log.Error().Err(err).Str("email", req.Email).Msg("Failed to register user")
		utils.RespondInternalError(c, err, "Failed to register user")
		return
	}

	log.Info().Uint("userId", result.User.ID).Str("email", result.User.Email).Msg("Successfully registered user")

	utils.RespondCreated(c, result, "User registered successfully")
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "User login data"
// @Example request
//
//	{
//	  "email": "user@example.com",
//	  "password": "securePassword123"
//	}
//
// @Success 200 {object} models.APIResponse[models.AuthData] "Successfully authenticated user"
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
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format"
// @Failure 401 {object} models.ErrorResponse[error] "Invalid credentials or inactive account"
// @Example response
//
//	{
//	  "error": "unauthorized",
//	  "message": "Invalid email or password",
//	  "details": {},
//	  "timestamp": "2025-03-16T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for user login")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().Str("email", req.Email).Msg("User attempting to login")

	result, err := h.service.Login(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "inactive") {
			log.Warn().Err(err).Str("email", req.Email).Msg("Login failed - authentication error")
			utils.RespondUnauthorized(c, err, err.Error())
			return
		}

		log.Error().Err(err).Str("email", req.Email).Msg("Login failed - server error")
		utils.RespondInternalError(c, err, "Failed to authenticate user")
		return
	}

	log.Info().Uint("userId", result.User.ID).Str("email", result.User.Email).Msg("User successfully logged in")

	utils.RespondOK(c, result, "Login successful")
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token data"
// @Example request
//
//	{
//	  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	}
//
// @Success 200 {object} models.APIResponse[models.AuthData] "Successfully refreshed token"
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
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format"
// @Failure 401 {object} models.ErrorResponse[error] "Invalid or expired refresh token"
// @Example response
//
//	{
//	  "error": "unauthorized",
//	  "message": "Invalid refresh token",
//	  "details": {},
//	  "timestamp": "2025-03-16T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for token refresh")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().Msg("Token refresh requested")

	result, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "expired") {
			log.Warn().Err(err).Msg("Token refresh failed - invalid token")
			utils.RespondUnauthorized(c, err, err.Error())
			return
		}

		log.Error().Err(err).Msg("Token refresh failed - server error")
		utils.RespondInternalError(c, err, "Failed to refresh token")
		return
	}

	log.Info().Uint("userId", result.User.ID).Msg("Token refreshed successfully")

	utils.RespondOK(c, result, "Token refreshed successfully")
}

// Logout godoc
// @Summary Log out a user
// @Description Invalidates the refresh token, effectively logging the user out
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LogoutRequest true "Logout data"
// @Example request
//
//	{
//	  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//	}
//
// @Success 200 {object} models.APIResponse[any] "Successfully logged out"
// @Example response
//
//	{
//	  "success": true,
//	  "message": "Successfully logged out"
//	}
//
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format"
// @Failure 401 {object} models.ErrorResponse[error] "Invalid refresh token"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req models.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for logout")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().Msg("User logout requested")

	err := h.service.Logout(ctx, req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			log.Warn().Err(err).Msg("Logout failed - invalid token")
			utils.RespondUnauthorized(c, err, err.Error())
			return
		}

		log.Error().Err(err).Msg("Logout failed - server error")
		utils.RespondInternalError(c, err, "Failed to log out")
		return
	}

	log.Info().Msg("User successfully logged out")

	c.JSON(http.StatusOK, models.APIResponse[any]{
		Success: true,
		Message: "Successfully logged out",
	})
}
