package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"suasor/types/models"
)

// SessionRepository defines the session repository interface
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*models.Session, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	FindByUserID(ctx context.Context, userID uint64) ([]*models.Session, error)
	DeleteExpired(ctx context.Context) error
}

// sessionRepository implements the SessionRepository interface
type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new SessionRepository instance
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

// Create creates a new session
func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	result := r.db.Create(session)
	return result.Error
}

// Update updates an existing session
func (r *sessionRepository) Update(ctx context.Context, session *models.Session) error {
	result := r.db.Save(session)
	return result.Error
}

// Delete deletes a session by ID
func (r *sessionRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.Delete(&models.Session{}, id)
	return result.Error
}

// FindByID finds a session by ID
func (r *sessionRepository) FindByID(ctx context.Context, id uint64) (*models.Session, error) {
	var session models.Session
	result := r.db.First(&session, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}

	return &session, nil
}

// FindByRefreshToken finds a session by refresh token
func (r *sessionRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var session models.Session
	result := r.db.Where("refresh_token = ?", refreshToken).First(&session)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}

	return &session, nil
}

// FindByUserID finds all sessions for a specific user
func (r *sessionRepository) FindByUserID(ctx context.Context, userID uint64) ([]*models.Session, error) {
	var sessions []*models.Session
	result := r.db.Where("user_id = ?", userID).Find(&sessions)

	if result.Error != nil {
		return nil, result.Error
	}

	return sessions, nil
}

// DeleteExpired deletes all expired sessions
func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	result := r.db.Where("expires_at < ?", now).Delete(&models.Session{})
	return result.Error
}
