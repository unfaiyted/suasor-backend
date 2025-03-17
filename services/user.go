package services

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"suasor/models"
	"suasor/repository"
)

// Service errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidRole        = errors.New("invalid role")
)

// UserService defines the user service interface
type UserService interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*models.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*models.UserResponse, error)
	GetByUsername(ctx context.Context, username string) (*models.UserResponse, error)
	UpdatePassword(ctx context.Context, id uint64, currentPassword, newPassword string) error
	UpdateProfile(ctx context.Context, id uint64, updateData map[string]interface{}) error
	ChangeRole(ctx context.Context, id uint64, newRole string) error
	ActivateUser(ctx context.Context, id uint64) error
	DeactivateUser(ctx context.Context, id uint64) error
	RecordLogin(ctx context.Context, id uint64) error
}

// userService implements the UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Create creates a new user with a hashed password
func (s *userService) Create(ctx context.Context, user *models.User) error {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return ErrEmailExists
	}

	// Check if username already exists
	existingUser, err = s.userRepo.FindByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return ErrUsernameExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Set hashed password
	user.Password = string(hashedPassword)

	// Create user
	return s.userRepo.Create(ctx, user)
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

// Delete deletes a user by ID
func (s *userService) Delete(ctx context.Context, id uint64) error {
	return s.userRepo.Delete(ctx, id)
}

// GetByID retrieves a user by ID and returns a UserResponse
func (s *userService) GetByID(ctx context.Context, id uint64) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &models.UserResponse{
		ID:       uint64(user.ID),
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

// GetByEmail retrieves a user by email and returns a UserResponse
func (s *userService) GetByEmail(ctx context.Context, email string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &models.UserResponse{
		ID:       uint64(user.ID),
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

// GetByUsername retrieves a user by username and returns a UserResponse
func (s *userService) GetByUsername(ctx context.Context, username string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &models.UserResponse{
		ID:       uint64(user.ID),
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

// UpdatePassword updates a user's password after verifying the current password
func (s *userService) UpdatePassword(ctx context.Context, id uint64, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// UpdateProfile updates user profile information
func (s *userService) UpdateProfile(ctx context.Context, id uint64, updateData map[string]interface{}) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Check if email is being updated and is unique
	if email, ok := updateData["email"].(string); ok && email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(ctx, email)
		if err == nil && existingUser != nil {
			return ErrEmailExists
		}
		user.Email = email
	}

	// Check if username is being updated and is unique
	if username, ok := updateData["username"].(string); ok && username != user.Username {
		existingUser, err := s.userRepo.FindByUsername(ctx, username)
		if err == nil && existingUser != nil {
			return ErrUsernameExists
		}
		user.Username = username
	}

	return s.userRepo.Update(ctx, user)
}

// ChangeRole changes a user's role
func (s *userService) ChangeRole(ctx context.Context, id uint64, newRole string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Only allow valid roles
	if newRole != "user" && newRole != "admin" {
		return ErrInvalidRole
	}

	user.Role = newRole
	return s.userRepo.Update(ctx, user)
}

// ActivateUser activates a user account
func (s *userService) ActivateUser(ctx context.Context, id uint64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	user.Active = true
	return s.userRepo.Update(ctx, user)
}

// DeactivateUser deactivates a user account
func (s *userService) DeactivateUser(ctx context.Context, id uint64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	user.Active = false
	return s.userRepo.Update(ctx, user)
}

// RecordLogin updates the LastLogin timestamp for a user
func (s *userService) RecordLogin(ctx context.Context, id uint64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	now := time.Now()
	user.LastLogin = &now
	return s.userRepo.Update(ctx, user)
}
