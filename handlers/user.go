package handlers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"

	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service    services.UserService
	configSvc  services.ConfigService
	avatarPath string
	maxSize    int
}

func NewUserHandler(service services.UserService, configSvc services.ConfigService) *UserHandler {
	config := configSvc.GetConfig()
	return &UserHandler{
		service:    service,
		configSvc:  configSvc,
		avatarPath: config.App.AvatarPath,
		maxSize:    config.App.MaxAvatarSize,
	}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account with the provided information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body	requests.RegisterRequest	true	"User registration information"
//	@Example		request
//
//	{
//	  "email": "user@example.com",
//	  "username": "johndoe",
//	  "password": "securepassword123"
//	}
//
//	@Success		201	{object}	responses.APIResponse[responses.UserResponse]	"Successfully registered user"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": {
//	    "id": 1,
//	    "email": "user@example.com",
//	    "username": "johndoe",
//	    "role": "user"
//	  },
//	  "message": "User registered successfully"
//	}
//
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid request format or email/username already exists"
//	@Example		response
//
//	{
//	  "error": "bad_request",
//	  "message": "Email already exists",
//	  "details": {
//	    "error": "email already exists"
//	  },
//	  "timestamp": "2023-06-08T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.RegisterRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().Str("email", req.Email).Str("username", req.Username).Msg("Registering new user")

	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Role:     "user", // Default role
		Active:   true,   // Default to active
	}

	if err := h.service.Create(ctx, user); err != nil {
		if err.Error() == "email already exists" {
			log.Warn().Err(err).Str("email", req.Email).Msg("Email already exists")
			responses.RespondBadRequest(c, err, "Email already exists")
			return
		}
		if err.Error() == "username already exists" {
			log.Warn().Err(err).Str("username", req.Username).Msg("Username already exists")
			responses.RespondBadRequest(c, err, "Username already exists")
			return
		}

		log.Error().Err(err).Str("email", req.Email).Msg("Failed to register user")
		responses.RespondInternalError(c, err, "Failed to register user")
		return
	}

	// Get the user response to return
	userResponse, err := h.service.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to retrieve registered user")
		responses.RespondInternalError(c, err, "Failed to retrieve registered user")
		return
	}

	log.Info().Uint64("id", userResponse.ID).Str("email", userResponse.Email).Msg("Successfully registered user")
	responses.RespondCreated(c, userResponse, "User registered successfully")
}

// GetProfile godoc
//
//	@Summary		Get the current user's profile
//	@Description	Retrieves the profile information for the currently authenticated user
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	responses.APIResponse[responses.UserResponse]	"Successfully retrieved user profile"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": {
//	    "id": 1,
//	    "email": "user@example.com",
//	    "username": "johndoe",
//	    "role": "user"
//	  },
//	  "message": "Profile retrieved successfully"
//	}
//
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get user userID from context (set by auth middleware)
	userID, _ := checkUserAccess(c)

	log.Info().Uint64("id", userID).Msg("Retrieving user profile")

	userResponse, err := h.service.GetByID(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", userID).Msg("User not found")
			responses.RespondNotFound(c, err, "User not found")
			return
		}
		log.Error().Err(err).Uint64("id", userID).Msg("Failed to retrieve user profile")
		responses.RespondInternalError(c, err, "Failed to retrieve user profile")
		return
	}

	log.Info().Uint64("id", userID).Msg("Successfully retrieved user profile")
	responses.RespondOK(c, userResponse, "Profile retrieved successfully")
}

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Updates the profile information for the currently authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	requests.UpdateUserRequest	true	"Updated profile information"
//	@Example		request
//
//	{
//	  "email": "newemail@example.com",
//	  "username": "newusername"
//	}
//
//	@Success		200	{object}	responses.APIResponse[responses.UserResponse]	"Successfully updated user profile"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": {
//	    "id": 1,
//	    "email": "newemail@example.com",
//	    "username": "newusername",
//	    "role": "user"
//	  },
//	  "message": "Profile updated successfully"
//	}
//
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid request format or email/username already exists"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get user userID from context (set by auth middleware)
	userID, _ := checkUserAccess(c)

	var req requests.UpdateUserRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().Uint64("id", userID).Str("email", req.Email).Str("username", req.Username).Msg("Updating user profile")

	// Create map of fields to update
	updateData := make(map[string]interface{})
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Username != "" {
		updateData["username"] = req.Username
	}
	if req.Avatar != "" {
		updateData["avatar"] = req.Avatar
	}

	if err := h.service.UpdateProfile(ctx, userID, updateData); err != nil {
		if err.Error() == "email already exists" {
			log.Warn().Err(err).Str("email", req.Email).Msg("Email already exists")
			responses.RespondBadRequest(c, err, "Email already exists")
			return
		}
		if err.Error() == "username already exists" {
			log.Warn().Err(err).Str("username", req.Username).Msg("Username already exists")
			responses.RespondBadRequest(c, err, "Username already exists")
			return
		}
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", userID).Msg("User not found")
			responses.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("id", userID).Msg("Failed to update user profile")
		responses.RespondInternalError(c, err, "Failed to update user profile")
		return
	}

	// Get updated user response
	userResponse, err := h.service.GetByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint64("id", userID).Msg("Failed to retrieve updated user profile")
		responses.RespondInternalError(c, err, "Failed to retrieve updated user profile")
		return
	}

	log.Info().Uint64("id", userID).Msg("Successfully updated user profile")
	responses.RespondOK(c, userResponse, "Profile updated successfully")
}

// ChangePassword godoc
//
//	@Summary		Change user password
//	@Description	Changes the password for the currently authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	requests.ChangePasswordRequest	true	"Password change information"
//	@Example		request
//
//	{
//	  "currentPassword": "oldpassword123",
//	  "newPassword": "newpassword456"
//	}
//
//	@Success		200	{object}	responses.SuccessResponse	"Password changed successfully"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": null,
//	  "message": "Password changed successfully"
//	}
//
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid request format or incorrect current password"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get user userID from context (set by auth middleware)
	userID, _ := checkUserAccess(c)

	var req requests.ChangePasswordRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().Uint64("id", userID).Msg("Changing user password")

	if err := h.service.UpdatePassword(ctx, userID, req.CurrentPassword, req.NewPassword); err != nil {
		if err.Error() == "invalid credentials" {
			log.Warn().Uint64("id", userID).Msg("Invalid current password")
			responses.RespondBadRequest(c, err, "Current password is incorrect")
			return
		}
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", userID).Msg("User not found")
			responses.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("id", userID).Msg("Failed to change password")
		responses.RespondInternalError(c, err, "Failed to change password")
		return
	}

	log.Info().Uint64("id", userID).Msg("Successfully changed password")
	responses.RespondOK(c, http.StatusOK, "Password changed successfully")
}

// GetByID godoc
//
//	@Summary		Get user by ID
//	@Description	Retrieves a user by their userID (admin only)
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userID	path		int												true	"User ID"	example:"1"
//	@Success		200		{object}	responses.APIResponse[responses.UserResponse]	"Successfully retrieved user"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": {
//	    "id": 1,
//	    "email": "user@example.com",
//	    "username": "johndoe",
//	    "role": "user"
//	  },
//	  "message": "User retrieved successfully"
//	}
//
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid user userID format"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		403	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Forbidden - Not an admin"
//	@Failure		404	{object}	responses.ErrorResponse[responses.ErrorDetails]	"User not found"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/{userID} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Verify admin role
	_, ok := checkAdminAccess(c)
	if !ok {
		return
	}

	targetUserID, _ := checkItemID(c, "userID")

	user, err := h.service.GetByID(ctx, targetUserID)
	if err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", targetUserID).Msg("User not found")
			responses.RespondNotFound(c, err, "User not found")
			return
		}
		log.Error().Err(err).Uint64("id", targetUserID).Msg("Failed to retrieve user")
		responses.RespondInternalError(c, err, "Failed to retrieve user")
		return
	}

	userResponse := responses.NewUserResponse(user)
	log.Info().Uint64("id", targetUserID).Msg("Successfully retrieved user")
	responses.RespondOK(c, userResponse, "User retrieved successfully")
}

// ChangeRole godoc
//
//	@Summary		Change user role
//	@Description	Changes a user's role (admin only)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userID	path	int							true	"User ID"	example:"1"
//	@Param			request	body	requests.ChangeRoleRequest	true	"New role information"
//	@Example		request
//
//	{
//	  "role": "admin"
//	}
//
//	@Success		200	{object}	responses.APIResponse[responses.UserResponse]	"Successfully changed user role"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": {
//	    "id": 1,
//	    "email": "user@example.com",
//	    "username": "johndoe",
//	    "role": "admin"
//	  },
//	  "message": "User role changed successfully"
//	}
//
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid user userID format or invalid role"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		403	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Forbidden - Not an admin"
//	@Failure		404	{object}	responses.ErrorResponse[responses.ErrorDetails]	"User not found"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/{userID}/role [put]
func (h *UserHandler) ChangeRole(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Verify admin role
	userID, ok := checkAdminAccess(c)
	if !ok {
		return
	}

	var req requests.ChangeRoleRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().Uint64("userID", userID).Str("newRole", req.Role).Msg("Changing user role")

	if err := h.service.ChangeRole(ctx, userID, req.Role); err != nil {
		if err.Error() == "invalid role" {
			log.Warn().Uint64("userID", userID).Str("role", req.Role).Msg("Invalid role specified")
			responses.RespondBadRequest(c, err, "Invalid role specified")
			return
		}
		if err.Error() == "user not found" {
			log.Warn().Uint64("userID", userID).Msg("User not found")
			responses.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to change user role")
		responses.RespondInternalError(c, err, "Failed to change user role")
		return
	}

	// Get updated user response
	userResponse, err := h.service.GetByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to retrieve updated user")
		responses.RespondInternalError(c, err, "Failed to retrieve updated user")
		return
	}

	log.Info().Uint64("userID", userID).Str("newRole", req.Role).Msg("Successfully changed user role")
	responses.RespondOK(c, userResponse, "User role changed successfully")
}

// ActivateUser godoc
//
//	@Summary		Activate a user account
//	@Description	Activates a user account (admin only)
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userID	path		int												true	"User ID"	example:"1"
//	@Success		200		{object}	responses.APIResponse[responses.UserResponse]	"Successfully activated user account"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": {
//	    "id": 1,
//	    "email": "user@example.com",
//	    "username": "johndoe",
//	    "role": "user"
//	  },
//	  "message": "User account activated successfully"
//	}
//
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid user userID format"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		403	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Forbidden - Not an admin"
//	@Failure		404	{object}	responses.ErrorResponse[responses.ErrorDetails]	"User not found"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/{userID}/activate [post]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Verify admin role
	userID, ok := checkAdminAccess(c)
	if !ok {
		return
	}

	log.Info().Uint64("userID", userID).Msg("Activating user account")

	if err := h.service.ActivateUser(ctx, userID); err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", userID).Msg("User not found")
			responses.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to activate user account")
		responses.RespondInternalError(c, err, "Failed to activate user account")
		return
	}

	// Get updated user response
	userResponse, err := h.service.GetByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to retrieve updated user")
		responses.RespondInternalError(c, err, "Failed to retrieve updated user")
		return
	}

	log.Info().Uint64("userID", userID).Msg("Successfully activated user account")
	responses.RespondOK(c, userResponse, "User account activated successfully")
}

// DeactivateUser godoc
//
//	@Summary		Deactivate a user account
//	@Description	Deactivates a user account (admin only)
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userID	path		int												true	"User ID"	example:"1"
//	@Success		200		{object}	responses.APIResponse[responses.UserResponse]	"Successfully deactivated user account"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": {
//	    "id": 1,
//	    "email": "user@example.com",
//	    "username": "johndoe",
//	    "role": "user"
//	  },
//	  "message": "User account deactivated successfully"
//	}
//
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid user userID format"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		403	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Forbidden - Not an admin"
//	@Failure		404	{object}	responses.ErrorResponse[responses.ErrorDetails]	"User not found"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/{userID}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Verify admin role
	userID, ok := checkAdminAccess(c)
	if !ok {
		return
	}

	log.Info().Uint64("userID", userID).Msg("Deactivating user account")

	if err := h.service.DeactivateUser(ctx, userID); err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", userID).Msg("User not found")
			responses.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to deactivate user account")
		responses.RespondInternalError(c, err, "Failed to deactivate user account")
		return
	}

	// Get updated user response
	userResponse, err := h.service.GetByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to retrieve updated user")
		responses.RespondInternalError(c, err, "Failed to retrieve updated user")
		return
	}

	log.Info().Uint64("userID", userID).Msg("Successfully deactivated user account")
	responses.RespondOK(c, userResponse, "User account deactivated successfully")
}

// Delete godoc
//
//	@Summary		Delete a user account
//	@Description	Deletes a user account (admin only)
//	@Tags			users
//	@Security		BearerAuth
//	@Param			userID	path	int	true	"User ID"	example:"1"
//	@Success		204		"No Content - User successfully deleted"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid user userID format"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		403		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Forbidden - Not an admin"
//	@Failure		404		{object}	responses.ErrorResponse[responses.ErrorDetails]	"User not found"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/{userID} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, _ := checkAdminAccess(c)

	log.Info().Uint64("id", userID).Msg("Deleting user account")

	if err := h.service.Delete(ctx, userID); err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", userID).Msg("User not found")
			responses.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to delete user account")
		responses.RespondInternalError(c, err, "Failed to delete user account")
		return
	}

	log.Info().Uint64("userID", userID).Msg("Successfully deleted user account")
	c.Status(http.StatusNoContent)
}

// UploadAvatar godoc
//
//	@Summary		Upload user avatar
//	@Description	Uploads a new avatar image for the currently authenticated user
//	@Tags			users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			avatar	formData	file													true	"Avatar image file (jpeg, png, gif only)"
//	@Success		200		{object}	responses.APIResponse[requests.AvatarUploadResponse]	"Successfully uploaded avatar"
//	@Example		response
//
//	{
//	  "success": true,
//	  "data": {
//	    "filePath": "/uploads/avatars/user_1.jpg"
//	  },
//	  "message": "Avatar uploaded successfully"
//	}
//
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid file format or size"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get user userID from context (set by auth middleware)
	id, ok := checkUserAccess(c)
	if !ok {
		return
	}

	// Log request details for debugging
	log.Debug().
		Str("content-type", c.GetHeader("Content-Type")).
		Str("content-length", c.GetHeader("Content-Length")).
		Msg("Avatar upload request received")

	// Get file from form
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get avatar file from request")
		responses.RespondBadRequest(c, err, "Failed to get avatar file: "+err.Error())
		return
	}

	// Log file details
	log.Debug().
		Str("filename", header.Filename).
		Int64("size", header.Size).
		Str("content-type", header.Header.Get("Content-Type")).
		Msg("Avatar file received")

	defer file.Close()

	// Validate file
	if err := h.validateAvatarFile(header); err != nil {
		log.Error().Err(err).Str("filename", header.Filename).Msg("Invalid avatar file")
		responses.RespondBadRequest(c, err, err.Error())
		return
	}

	// Create avatar directory if it doesn't exist
	if err := os.MkdirAll(h.avatarPath, os.ModePerm); err != nil {
		log.Error().Err(err).Str("path", h.avatarPath).Msg("Failed to create avatar directory")
		responses.RespondInternalError(c, err, "Failed to create avatar directory")
		return
	}

	// Generate unique filename based on user userID and original extension
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("user_%d%s", id, ext)
	filePath := filepath.Join(h.avatarPath, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("Failed to create avatar file")
		responses.RespondInternalError(c, err, "Failed to save avatar file")
		return
	}
	defer dst.Close()

	// Read the file into memory first to diagnose EOF issues
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("Failed to read avatar file contents")
		responses.RespondInternalError(c, err, "Failed to read avatar file: "+err.Error())
		return
	}

	log.Debug().Int("bytesRead", len(fileBytes)).Msg("Read file bytes into memory")

	// Write the file to disk
	if _, err = dst.Write(fileBytes); err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("Failed to write avatar file")
		responses.RespondInternalError(c, err, "Failed to save avatar file: "+err.Error())
		return
	}

	// Get base API URL from config for full avatar URL
	config := h.configSvc.GetConfig()
	baseURL := config.HTTP.BaseURL
	if baseURL == "" {
		// Default to host header if baseURL not set in config
		baseURL = fmt.Sprintf("http://%s", c.Request.Host)
	}

	// Create relative path for storage in DB
	relativePath := fmt.Sprintf("/uploads/avatars/%s", filename)

	// Generate full web-accessible URL for avatar
	fullURL := fmt.Sprintf("%s%s", baseURL, relativePath)

	log.Debug().
		Str("baseURL", baseURL).
		Str("relativePath", relativePath).
		Str("fullURL", fullURL).
		Msg("Generated avatar URLs")

	// Update user's avatar field in profile with the full URL
	updateData := map[string]interface{}{
		"avatar": fullURL,
	}

	if err := h.service.UpdateProfile(ctx, id, updateData); err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to update user avatar")
		responses.RespondInternalError(c, err, "Failed to update user avatar")
		return
	}

	log.Info().Uint64("id", id).Str("path", fullURL).Msg("Successfully uploaded avatar")
	responses.RespondOK(c, &requests.AvatarUploadResponse{FilePath: fullURL}, "Avatar uploaded successfully")
}

// validateAvatarFile checks if the file is valid
func (h *UserHandler) validateAvatarFile(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > int64(h.maxSize) {
		return fmt.Errorf("file size exceeds the limit of %d bytes", h.maxSize)
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return errors.New("only JPEG, PNG and GIF images are allowed")
	}

	return nil
}

// ForgotPassword godoc
//
//	@Summary		Forgot password
//	@Description	Request a password reset email
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	requests.ForgotPasswordRequest	true	"Forgot password request"
//	@Success		200		{object}	responses.SuccessResponse	"Password reset email sent"
//	@Example		response
//
//	{
//	  "success": true,
//	  "message": "Password reset email sent"
//	}
//
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid request format"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/forgot-password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.ForgotPasswordRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().
		Str("email", req.Email).
		Msg("Sending password reset email")

	if err := h.service.ForgotPassword(ctx, req.Email); err != nil {
		log.Error().Err(err).
			Str("email", req.Email).
			Msg("Failed to send password reset email")
		responses.RespondInternalError(c, err, "Failed to send password reset email")
		return
	}

	log.Info().
		Str("email", req.Email).
		Msg("Password reset email sent")
	responses.RespondOK(c, http.StatusOK, "Password reset email sent")
}

// ResetPassword godoc
//
//	@Summary		Reset password
//	@Description	Reset the user's password using a password reset token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	requests.ResetPasswordRequest	true	"Reset password request"
//	@Success		200		{object}	responses.SuccessResponse	"Password reset successfully"
//	@Example		response
//
//	{
//	  "success": true,
//	  "message": "Password reset successfully"
//	}
//
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid request format"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized - Not logged in"
//	@Failure		404		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Password reset token not found"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/user/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.ResetPasswordRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	log.Info().
		Str("email", req.Email).
		Str("token", req.Token).
		Msg("Resetting password")

	if err := h.service.ResetPassword(ctx, req.Email, req.Token, req.NewPassword); err != nil {
		log.Error().Err(err).
			Str("email", req.Email).
			Str("token", req.Token).
			Msg("Failed to reset password")
		responses.RespondInternalError(c, err, "Failed to reset password")
		return
	}

	log.Info().
		Str("email", req.Email).
		Str("token", req.Token).
		Msg("Password reset successfully")
	responses.RespondOK(c, http.StatusOK, "Password reset successfully")
}
