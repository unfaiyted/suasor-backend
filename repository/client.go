// repository/client.go
package repository

import (
	"context"
	"fmt"
	"suasor/client/types"

	"suasor/types/models"

	"gorm.io/gorm"
)

// ClientRepository defines the interface for media client database operations
type ClientRepository[T types.ClientConfig] interface {
	Create(ctx context.Context, client models.Client[T]) (*models.Client[T], error)
	Update(ctx context.Context, client models.Client[T]) (*models.Client[T], error)

	// Common operations
	GetByID(ctx context.Context, id, userID uint64) (*models.Client[T], error)
	GetByUserID(ctx context.Context, userID uint64) ([]*models.Client[T], error)
	GetByCategory(ctx context.Context, clientType types.ClientCategory, userID uint64) ([]*models.Client[T], error)
	GetByType(ctx context.Context, clientType types.ClientType, userID uint64) ([]*models.Client[T], error)
	Delete(ctx context.Context, id, userID uint64) error
}

type clientRepository[T types.ClientConfig] struct {
	db *gorm.DB
}

// NewClientRepository creates a new media client repository
func NewClientRepository[T types.ClientConfig](db *gorm.DB) ClientRepository[T] {
	return &clientRepository[T]{db: db}
}

func (r *clientRepository[T]) Create(ctx context.Context, client models.Client[T]) (*models.Client[T], error) {

	if err := r.db.WithContext(ctx).Create(&client).Error; err != nil {
		clientType := client.GetConfig()
		return nil, fmt.Errorf("failed to create %s client: %w", clientType, err)
	}

	updatedClient, err := r.GetByID(ctx, client.ID, client.UserID)
	if err != nil {
		return nil, err
	}
	return updatedClient, nil

}

func (r *clientRepository[T]) Update(ctx context.Context, client models.Client[T]) (*models.Client[T], error) {
	// Get existing record first to check if it exists and preserve createdAt
	var existing models.Client[T]

	category := client.GetCategory()
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", client.ID, client.UserID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media client not found")
		}
		return nil, fmt.Errorf("failed to find %s client: %w", category, err)
	}

	// Preserve createdAt
	client.CreatedAt = existing.CreatedAt

	// Update the record
	if err := r.db.WithContext(ctx).Save(&client).Error; err != nil {
		return nil, fmt.Errorf("failed to update %s client: %w", category, err)
	}
	return &client, nil

}

// GetByID retrieves a media client by ID
// GetByID retrieves a media client by ID
func (r *clientRepository[T]) GetByID(ctx context.Context, id, userID uint64) (*models.Client[T], error) {
	var client models.Client[T]

	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return &client, nil
}

// GetByUserID retrieves all media clients for a user
// GetByUserID retrieves all media clients for a user
func (r *clientRepository[T]) GetByUserID(ctx context.Context, userID uint64) ([]*models.Client[T], error) {
	var clients []*models.Client[T]

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&clients).Error; err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}

	return clients, nil
}

// Delete deletes a media client
// Delete deletes a media client
func (r *clientRepository[T]) Delete(ctx context.Context, id, userID uint64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Client[T]{})

	if err := result.Error; err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("client not found")
	}

	return nil
}

func (r *clientRepository[T]) GetByCategory(ctx context.Context, category types.ClientCategory, userID uint64) ([]*models.Client[T], error) {
	var clients []*models.Client[T]

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND category = ?", userID, category).
		Find(&clients).Error; err != nil {
		return nil, fmt.Errorf("failed to get clients by type: %w", err)
	}

	return clients, nil
}

func (r *clientRepository[T]) GetByType(ctx context.Context, clientType types.ClientType, userID uint64) ([]*models.Client[T], error) {
	var clients []*models.Client[T]

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND  config ->> 'type' = ?", userID, clientType).
		Find(&clients).Error; err != nil {
		return nil, fmt.Errorf("failed to get clients by type: %w", err)
	}

	return clients, nil
}
