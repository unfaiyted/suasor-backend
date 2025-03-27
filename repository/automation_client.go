// repository/download_client.go
package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"suasor/types/models"
)

// AutomationClientRepository defines the interface for download client operations
type AutomationClientRepository interface {
	Create(ctx context.Context, client models.AutomationClient) (models.AutomationClient, error)
	GetByID(ctx context.Context, id, userID uint64) (models.AutomationClient, error)
	GetByUserID(ctx context.Context, userID uint64) ([]models.AutomationClient, error)
	Update(ctx context.Context, client models.AutomationClient) error
	Delete(ctx context.Context, id, userID uint64) error
}

type downloadClientRepo struct {
	db *gorm.DB
}

func NewAutomationClientRepository(db *gorm.DB) AutomationClientRepository {
	return &downloadClientRepo{db: db}
}

func (r *downloadClientRepo) Create(ctx context.Context, client models.AutomationClient) (models.AutomationClient, error) {
	if err := r.db.WithContext(ctx).Create(&client).Error; err != nil {
		return models.AutomationClient{}, err
	}
	return client, nil
}

func (r *downloadClientRepo) GetByID(ctx context.Context, id, userID uint64) (models.AutomationClient, error) {
	var client models.AutomationClient
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
		return models.AutomationClient{}, err
	}
	return client, nil
}

func (r *downloadClientRepo) GetByUserID(ctx context.Context, userID uint64) ([]models.AutomationClient, error) {
	var clients []models.AutomationClient
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *downloadClientRepo) Update(ctx context.Context, client models.AutomationClient) error {
	updateMap := map[string]interface{}{
		"name":        client.Name,
		"client_type": client.ClientType,
		"url":         client.URL,
		"api_key":     client.APIKey,
		"is_enabled":  client.IsEnabled,
		"updated_at":  client.UpdatedAt,
	}

	// Debug: print the SQL query
	result := r.db.WithContext(ctx).Model(&models.AutomationClient{}).
		Where("id = ? AND user_id = ?", client.ID, client.UserID).
		Updates(updateMap)

	// Check if rows were affected
	if result.RowsAffected == 0 {
		// No rows affected can mean the WHERE clause didn't match any records
		return fmt.Errorf("no records updated: id=%d, user_id=%d may not exist",
			client.ID, client.UserID)
	}

	return result.Error
}

func (r *downloadClientRepo) Delete(ctx context.Context, id, userID uint64) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.AutomationClient{}).Error
}
