package handlers

import (
	"net/http"
	"strconv"
	"suasor/models"
	"suasor/services"
	"suasor/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "User registration information"
// @Example request
//
//	{
//	  "email": "user@example.com",
//	  "username": "johndoe",
//	  "password": "securepassword123"
//	}
//
// @Success 201 {object} models.APIResponse[models.UserResponse] "Successfully registered user"
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
//	  "message": "User registered successfully"
//	}
//
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format or email/username already exists"
// @Example response
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
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for user registration")
		utils.RespondValidationError(c, err)
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
			utils.RespondBadRequest(c, err, "Email already exists")
			return
		}
		if err.Error() == "username already exists" {
			log.Warn().Err(err).Str("username", req.Username).Msg("Username already exists")
			utils.RespondBadRequest(c, err, "Username already exists")
			return
		}

		log.Error().Err(err).Str("email", req.Email).Msg("Failed to register user")
		utils.RespondInternalError(c, err, "Failed to register user")
		return
	}

	// Get the user response to return
	userResponse, err := h.service.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to retrieve registered user")
		utils.RespondInternalError(c, err, "Failed to retrieve registered user")
		return
	}

	log.Info().Uint64("id", userResponse.ID).Str("email", userResponse.Email).Msg("Successfully registered user")
	utils.RespondCreated(c, userResponse, "User registered successfully")
}

// GetProfile godoc
// @Summary Get the current user's profile
// @Description Retrieves the profile information for the currently authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse[models.UserResponse] "Successfully retrieved user profile"
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
//	  "message": "Profile retrieved successfully"
//	}
//
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized - Not logged in"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("User ID not found in context")
		utils.RespondUnauthorized(c, nil, "Not authenticated")
		return
	}

	id := userID.(uint64)
	log.Info().Uint64("id", id).Msg("Retrieving user profile")

	userResponse, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", id).Msg("User not found")
			utils.RespondNotFound(c, err, "User not found")
			return
		}
		log.Error().Err(err).Uint64("id", id).Msg("Failed to retrieve user profile")
		utils.RespondInternalError(c, err, "Failed to retrieve user profile")
		return
	}

	log.Info().Uint64("id", id).Msg("Successfully retrieved user profile")
	utils.RespondOK(c, userResponse, "Profile retrieved successfully")
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Updates the profile information for the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateProfileRequest true "Updated profile information"
// @Example request
//
//	{
//	  "email": "newemail@example.com",
//	  "username": "newusername"
//	}
//
// @Success 200 {object} models.APIResponse[models.UserResponse] "Successfully updated user profile"
// @Example response
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
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format or email/username already exists"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized - Not logged in"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("User ID not found in context")
		utils.RespondUnauthorized(c, nil, "Not authenticated")
		return
	}

	id := userID.(uint64)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for profile update")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().Uint64("id", id).Str("email", req.Email).Str("username", req.Username).Msg("Updating user profile")

	// Create map of fields to update
	updateData := make(map[string]interface{})
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Username != "" {
		updateData["username"] = req.Username
	}

	if err := h.service.UpdateProfile(ctx, id, updateData); err != nil {
		if err.Error() == "email already exists" {
			log.Warn().Err(err).Str("email", req.Email).Msg("Email already exists")
			utils.RespondBadRequest(c, err, "Email already exists")
			return
		}
		if err.Error() == "username already exists" {
			log.Warn().Err(err).Str("username", req.Username).Msg("Username already exists")
			utils.RespondBadRequest(c, err, "Username already exists")
			return
		}
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", id).Msg("User not found")
			utils.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("id", id).Msg("Failed to update user profile")
		utils.RespondInternalError(c, err, "Failed to update user profile")
		return
	}

	// Get updated user response
	userResponse, err := h.service.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to retrieve updated user profile")
		utils.RespondInternalError(c, err, "Failed to retrieve updated user profile")
		return
	}

	log.Info().Uint64("id", id).Msg("Successfully updated user profile")
	utils.RespondOK(c, userResponse, "Profile updated successfully")
}

// ChangePassword godoc
// @Summary Change user password
// @Description Changes the password for the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Password change information"
// @Example request
//
//	{
//	  "currentPassword": "oldpassword123",
//	  "newPassword": "newpassword456"
//	}
//
// @Success 200 {object} models.APIResponse[string] "Successfully changed password"
// @Example response
//
//	{
//	  "success": true,
//	  "data": null,
//	  "message": "Password changed successfully"
//	}
//
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format or incorrect current password"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized - Not logged in"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("User ID not found in context")
		utils.RespondUnauthorized(c, nil, "Not authenticated")
		return
	}

	id := userID.(uint64)

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for password change")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().Uint64("id", id).Msg("Changing user password")

	if err := h.service.UpdatePassword(ctx, id, req.CurrentPassword, req.NewPassword); err != nil {
		if err.Error() == "invalid credentials" {
			log.Warn().Uint64("id", id).Msg("Invalid current password")
			utils.RespondBadRequest(c, err, "Current password is incorrect")
			return
		}
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", id).Msg("User not found")
			utils.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("id", id).Msg("Failed to change password")
		utils.RespondInternalError(c, err, "Failed to change password")
		return
	}

	log.Info().Uint64("id", id).Msg("Successfully changed password")
	utils.RespondOK(c, http.StatusOK, "Password changed successfully")
}

// GetByID godoc
// @Summary Get user by ID
// @Description Retrieves a user by their ID (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID" example:"1"
// @Success 200 {object} models.APIResponse[models.UserResponse] "Successfully retrieved user"
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
//	  "message": "User retrieved successfully"
//	}
//
// @Failure 400 {object} models.ErrorResponse[error] "Invalid user ID format"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized - Not logged in"
// @Failure 403 {object} models.ErrorResponse[error] "Forbidden - Not an admin"
// @Failure 404 {object} models.ErrorResponse[error] "User not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Verify admin role (should be done by middleware, but double-checking)
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		log.Warn().Msg("Non-admin attempted to access user by ID")
		utils.RespondForbidden(c, nil, "Admin access required")
		return
	}

	// Parse user ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Warn().Str("id", idStr).Msg("Invalid user ID format")
		utils.RespondBadRequest(c, err, "Invalid user ID format")
		return
	}

	log.Info().Uint64("id", id).Msg("Admin retrieving user by ID")

	userResponse, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", id).Msg("User not found")
			utils.RespondNotFound(c, err, "User not found")
			return
		}
		log.Error().Err(err).Uint64("id", id).Msg("Failed to retrieve user")
		utils.RespondInternalError(c, err, "Failed to retrieve user")
		return
	}

	log.Info().Uint64("id", id).Msg("Successfully retrieved user")
	utils.RespondOK(c, userResponse, "User retrieved successfully")
}

// ChangeRole godoc
// @Summary Change user role
// @Description Changes a user's role (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID" example:"1"
// @Param request body models.ChangeRoleRequest true "New role information"
// @Example request
//
//	{
//	  "role": "admin"
//	}
//
// @Success 200 {object} models.APIResponse[models.UserResponse] "Successfully changed user role"
// @Example response
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
// @Failure 400 {object} models.ErrorResponse[error] "Invalid user ID format or invalid role"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized - Not logged in"
// @Failure 403 {object} models.ErrorResponse[error] "Forbidden - Not an admin"
// @Failure 404 {object} models.ErrorResponse[error] "User not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/{id}/role [put]
func (h *UserHandler) ChangeRole(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Verify admin role (should be done by middleware, but double-checking)
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		log.Warn().Msg("Non-admin attempted to change user role")
		utils.RespondForbidden(c, nil, "Admin access required")
		return
	}

	// Parse user ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Warn().Str("id", idStr).Msg("Invalid user ID format")
		utils.RespondBadRequest(c, err, "Invalid user ID format")
		return
	}

	var req models.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for role change")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().Uint64("id", id).Str("newRole", req.Role).Msg("Changing user role")

	if err := h.service.ChangeRole(ctx, id, req.Role); err != nil {
		if err.Error() == "invalid role" {
			log.Warn().Uint64("id", id).Str("role", req.Role).Msg("Invalid role specified")
			utils.RespondBadRequest(c, err, "Invalid role specified")
			return
		}
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", id).Msg("User not found")
			utils.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("id", id).Msg("Failed to change user role")
		utils.RespondInternalError(c, err, "Failed to change user role")
		return
	}

	// Get updated user response
	userResponse, err := h.service.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to retrieve updated user")
		utils.RespondInternalError(c, err, "Failed to retrieve updated user")
		return
	}

	log.Info().Uint64("id", id).Str("newRole", req.Role).Msg("Successfully changed user role")
	utils.RespondOK(c, userResponse, "User role changed successfully")
}

// ActivateUser godoc
// @Summary Activate a user account
// @Description Activates a user account (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID" example:"1"
// @Success 200 {object} models.APIResponse[models.UserResponse] "Successfully activated user account"
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
//	  "message": "User account activated successfully"
//	}
//
// @Failure 400 {object} models.ErrorResponse[error] "Invalid user ID format"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized - Not logged in"
// @Failure 403 {object} models.ErrorResponse[error] "Forbidden - Not an admin"
// @Failure 404 {object} models.ErrorResponse[error] "User not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/{id}/activate [post]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Verify admin role (should be done by middleware, but double-checking)
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		log.Warn().Msg("Non-admin attempted to activate user account")
		utils.RespondForbidden(c, nil, "Admin access required")
		return
	}

	// Parse user ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Warn().Str("id", idStr).Msg("Invalid user ID format")
		utils.RespondBadRequest(c, err, "Invalid user ID format")
		return
	}

	log.Info().Uint64("id", id).Msg("Activating user account")

	if err := h.service.ActivateUser(ctx, id); err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", id).Msg("User not found")
			utils.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("id", id).Msg("Failed to activate user account")
		utils.RespondInternalError(c, err, "Failed to activate user account")
		return
	}

	// Get updated user response
	userResponse, err := h.service.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to retrieve updated user")
		utils.RespondInternalError(c, err, "Failed to retrieve updated user")
		return
	}

	log.Info().Uint64("id", id).Msg("Successfully activated user account")
	utils.RespondOK(c, userResponse, "User account activated successfully")
}

// DeactivateUser godoc
// @Summary Deactivate a user account
// @Description Deactivates a user account (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID" example:"1"
// @Success 200 {object} models.APIResponse[models.UserResponse] "Successfully deactivated user account"
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
//	  "message": "User account deactivated successfully"
//	}
//
// @Failure 400 {object} models.ErrorResponse[error] "Invalid user ID format"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized - Not logged in"
// @Failure 403 {object} models.ErrorResponse[error] "Forbidden - Not an admin"
// @Failure 404 {object} models.ErrorResponse[error] "User not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/{id}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Verify admin role (should be done by middleware, but double-checking)
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		log.Warn().Msg("Non-admin attempted to deactivate user account")
		utils.RespondForbidden(c, nil, "Admin access required")
		return
	}

	// Parse user ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Warn().Str("id", idStr).Msg("Invalid user ID format")
		utils.RespondBadRequest(c, err, "Invalid user ID format")
		return
	}

	log.Info().Uint64("id", id).Msg("Deactivating user account")

	if err := h.service.DeactivateUser(ctx, id); err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", id).Msg("User not found")
			utils.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("id", id).Msg("Failed to deactivate user account")
		utils.RespondInternalError(c, err, "Failed to deactivate user account")
		return
	}

	// Get updated user response
	userResponse, err := h.service.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to retrieve updated user")
		utils.RespondInternalError(c, err, "Failed to retrieve updated user")
		return
	}

	log.Info().Uint64("id", id).Msg("Successfully deactivated user account")
	utils.RespondOK(c, userResponse, "User account deactivated successfully")
}

// Delete godoc
// @Summary Delete a user account
// @Description Deletes a user account (admin only)
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID" example:"1"
// @Success 204 "No Content - User successfully deleted"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid user ID format"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized - Not logged in"
// @Failure 403 {object} models.ErrorResponse[error] "Forbidden - Not an admin"
// @Failure 404 {object} models.ErrorResponse[error] "User not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Verify admin role (should be done by middleware, but double-checking)
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		log.Warn().Msg("Non-admin attempted to delete user account")
		utils.RespondForbidden(c, nil, "Admin access required")
		return
	}

	// Parse user ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Warn().Str("id", idStr).Msg("Invalid user ID format")
		utils.RespondBadRequest(c, err, "Invalid user ID format")
		return
	}

	log.Info().Uint64("id", id).Msg("Deleting user account")

	if err := h.service.Delete(ctx, id); err != nil {
		if err.Error() == "user not found" {
			log.Warn().Uint64("id", id).Msg("User not found")
			utils.RespondNotFound(c, err, "User not found")
			return
		}

		log.Error().Err(err).Uint64("id", id).Msg("Failed to delete user account")
		utils.RespondInternalError(c, err, "Failed to delete user account")
		return
	}

	log.Info().Uint64("id", id).Msg("Successfully deleted user account")
	c.Status(http.StatusNoContent)
}
