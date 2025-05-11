package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"suasor/repository"
	"suasor/types"
	"suasor/types/models"
	req "suasor/types/requests"
	res "suasor/types/responses"
)

// AuthService defines the authentication service interface
type AuthService interface {
	Register(ctx context.Context, request req.RegisterRequest) (*res.AuthDataResponse, error)
	Login(ctx context.Context, request req.LoginRequest) (*res.AuthDataResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*res.AuthDataResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateToken(ctx context.Context, token string) (*types.JWTClaim, error)
	GetAuthorizedUser(ctx context.Context, token string) (*res.UserResponse, error)
}

// authService implements the AuthService interface
type authService struct {
	userRepo      repository.UserRepository
	sessionRepo   repository.SessionRepository
	jwtSecret     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	tokenIssuer   string
	tokenAudience string
}

// NewAuthService creates a new AuthService instance
func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwtSecret string,
	accessExpiry time.Duration,
	refreshExpiry time.Duration,
	tokenIssuer string,
	tokenAudience string,
) AuthService {
	return &authService{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		jwtSecret:     []byte(jwtSecret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		tokenIssuer:   tokenIssuer,
		tokenAudience: tokenAudience,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, request req.RegisterRequest) (*res.AuthDataResponse, error) {
	// Check if email already exists
	_, err := s.userRepo.FindByEmail(ctx, request.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("error checking email: %w", err)
	}

	// Check if username already exists
	_, err = s.userRepo.FindByUsername(ctx, request.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("error checking username: %w", err)
	}

	// Create new user
	user := &models.User{
		Email:    request.Email,
		Username: request.Username,
		Role:     "user", // Default role
		Active:   true,
	}

	if err := user.SetPassword(request.Password); err != nil {
		return nil, fmt.Errorf("error setting password: %w", err)
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	// Create session
	now := time.Now()
	session := &models.Session{
		UserID:       uint64(user.ID),
		RefreshToken: tokens.RefreshUUID,
		UserAgent:    "User registration", // This would normally come from the request
		IP:           "0.0.0.0",           // This would normally come from the request
		ExpiresAt:    time.Unix(tokens.RtExpires, 0),
		LastUsed:     now,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	// Update last login time
	user.LastLogin = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &res.AuthDataResponse{
		User: res.UserResponse{
			ID:       uint64(user.ID),
			Email:    user.Email,
			Avatar:   user.Avatar,
			Username: user.Username,
			Role:     user.Role,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.AtExpires,
	}, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, request req.LoginRequest) (*res.AuthDataResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	// Check if account is active
	if !user.Active {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	match, err := user.CheckPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("error checking password: %w", err)
	}
	if !match {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	// Create session
	now := time.Now()
	session := &models.Session{
		UserID:       uint64(user.ID),
		RefreshToken: tokens.RefreshUUID,
		UserAgent:    "User login", // This would normally come from the request
		IP:           "0.0.0.0",    // This would normally come from the request
		ExpiresAt:    time.Unix(tokens.RtExpires, 0),
		LastUsed:     now,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	// Update last login time
	user.LastLogin = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &res.AuthDataResponse{
		User: res.UserResponse{
			ID:       uint64(user.ID),
			Email:    user.Email,
			Avatar:   user.Avatar,
			Username: user.Username,
			Role:     user.Role,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.AtExpires,
	}, nil
}

// RefreshToken refreshes the access token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*res.AuthDataResponse, error) {
	// Parse the refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &types.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*types.JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Find the session
	session, err := s.sessionRepo.FindByRefreshToken(ctx, claims.UUID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("invalid session")
		}
		return nil, fmt.Errorf("error finding session: %w", err)
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		if err := s.sessionRepo.Delete(ctx, uint64(session.ID)); err != nil {
			return nil, fmt.Errorf("error deleting expired session: %w", err)
		}
		return nil, errors.New("session expired")
	}

	// Find the user
	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	// Generate new tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	// Update session
	now := time.Now()
	session.RefreshToken = tokens.RefreshUUID
	session.ExpiresAt = time.Unix(tokens.RtExpires, 0)
	session.LastUsed = now

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("error updating session: %w", err)
	}

	return &res.AuthDataResponse{
		User: res.UserResponse{
			ID:       uint64(user.ID),
			Email:    user.Email,
			Avatar:   user.Avatar,
			Username: user.Username,
			Role:     user.Role,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.AtExpires,
	}, nil
}

// Logout invalidates the refresh token
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	// Parse the refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &types.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*types.JWTClaim)
	if !ok || !token.Valid {
		return errors.New("invalid token claims")
	}

	// Find and delete the session
	session, err := s.sessionRepo.FindByRefreshToken(ctx, claims.UUID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil // Session already doesn't exist, consider it a success
		}
		return fmt.Errorf("error finding session: %w", err)
	}

	if err := s.sessionRepo.Delete(ctx, uint64(session.ID)); err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}

// ValidateToken validates the JWT token
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*types.JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*types.JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// generateTokens generates access and refresh tokens
func (s *authService) generateTokens(user *models.User) (*types.TokenDetails, error) {
	td := &types.TokenDetails{
		AccessUUID:  uuid.New().String(),
		RefreshUUID: uuid.New().String(),
		AtExpires:   time.Now().Add(s.accessExpiry).Unix(),
		RtExpires:   time.Now().Add(s.refreshExpiry).Unix(),
	}

	// Create access token
	atClaims := types.JWTClaim{
		UserID: uint64(user.ID),
		UUID:   td.AccessUUID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.tokenIssuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			Audience:  []string{s.tokenAudience},
			ID:        td.AccessUUID,
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}
	td.AccessToken = accessToken

	// Create refresh token
	rtClaims := types.JWTClaim{
		UserID: uint64(user.ID),
		UUID:   td.RefreshUUID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.tokenIssuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			Audience:  []string{s.tokenAudience},
			ID:        td.RefreshUUID,
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}
	td.RefreshToken = refreshToken

	return td, nil
}

func (s *authService) GetAuthorizedUser(ctx context.Context, tokenString string) (*res.UserResponse, error) {
	var userID uint64

	// Try to get and validate token from context
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}
	userID = claims.UserID

	// Fetch user from repository
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	// Check if account is active
	if !user.Active {
		return nil, errors.New("account is inactive")
	}

	// Return user response
	return &res.UserResponse{
		ID:       uint64(user.ID),
		Email:    user.Email,
		Avatar:   user.Avatar,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}
