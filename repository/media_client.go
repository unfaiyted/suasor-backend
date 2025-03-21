// repository/media_client.go
package repository

import (
	"context"
	"fmt"
	"suasor/models"

	"gorm.io/gorm"
)

// MediaClientRepository defines the interface for media client database operations
type MediaClientRepository interface {
	// Typed operations
	CreatePlex(ctx context.Context, client models.MediaClient[models.PlexConfig]) (*models.MediaClient[models.PlexConfig], error)
	UpdatePlex(ctx context.Context, client models.MediaClient[models.PlexConfig]) (*models.MediaClient[models.PlexConfig], error)
	CreateJellyfin(ctx context.Context, client models.MediaClient[models.JellyfinConfig]) (*models.MediaClient[models.JellyfinConfig], error)
	UpdateJellyfin(ctx context.Context, client models.MediaClient[models.JellyfinConfig]) (*models.MediaClient[models.JellyfinConfig], error)
	CreateEmby(ctx context.Context, client models.MediaClient[models.EmbyConfig]) (*models.MediaClient[models.EmbyConfig], error)
	UpdateEmby(ctx context.Context, client models.MediaClient[models.EmbyConfig]) (*models.MediaClient[models.EmbyConfig], error)
	CreateNavidrome(ctx context.Context, client models.MediaClient[models.SubsonicConfig]) (*models.MediaClient[models.SubsonicConfig], error)
	UpdateNavidrome(ctx context.Context, client models.MediaClient[models.SubsonicConfig]) (*models.MediaClient[models.SubsonicConfig], error)

	// Common operations
	GetByID(ctx context.Context, id, userID uint64) (models.MediaClientResponse, error)
	GetByUserID(ctx context.Context, userID uint64) ([]models.MediaClientResponse, error)
	Delete(ctx context.Context, id, userID uint64) error
}

type mediaClientRepository struct {
	db *gorm.DB
}

// NewMediaClientRepository creates a new media client repository
func NewMediaClientRepository(db *gorm.DB) MediaClientRepository {
	return &mediaClientRepository{db: db}
}

// Generic helper functions to reduce code duplication (as package-level functions)

// createClient is a generic helper for creating any media client type
func createClient[T models.ClientConfig](db *gorm.DB, ctx context.Context, client models.MediaClient[T], clientType string) (*models.MediaClient[T], error) {
	if err := db.WithContext(ctx).Create(&client).Error; err != nil {
		return nil, fmt.Errorf("failed to create %s client: %w", clientType, err)
	}
	return &client, nil
}

// updateClient is a generic helper for updating any media client type
func updateClient[T models.ClientConfig](db *gorm.DB, ctx context.Context, client models.MediaClient[T], clientType string) (*models.MediaClient[T], error) {
	// Get existing record first to check if it exists and preserve createdAt
	var existing models.MediaClient[T]
	if err := db.WithContext(ctx).Where("id = ? AND user_id = ?", client.ID, client.UserID).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media client not found")
		}
		return nil, fmt.Errorf("failed to find %s client: %w", clientType, err)
	}

	// Preserve createdAt
	client.CreatedAt = existing.CreatedAt

	// Update the record
	if err := db.WithContext(ctx).Save(&client).Error; err != nil {
		return nil, fmt.Errorf("failed to update %s client: %w", clientType, err)
	}
	return &client, nil
}

// Type-specific methods using the generic helpers

// CreatePlex creates a new Plex media client
func (r *mediaClientRepository) CreatePlex(ctx context.Context, client models.MediaClient[models.PlexConfig]) (*models.MediaClient[models.PlexConfig], error) {
	return createClient(r.db, ctx, client, "plex")
}

// UpdatePlex updates an existing Plex media client
func (r *mediaClientRepository) UpdatePlex(ctx context.Context, client models.MediaClient[models.PlexConfig]) (*models.MediaClient[models.PlexConfig], error) {
	return updateClient(r.db, ctx, client, "plex")
}

// CreateJellyfin creates a new Jellyfin media client
func (r *mediaClientRepository) CreateJellyfin(ctx context.Context, client models.MediaClient[models.JellyfinConfig]) (*models.MediaClient[models.JellyfinConfig], error) {
	return createClient(r.db, ctx, client, "jellyfin")
}

// UpdateJellyfin updates an existing Jellyfin media client
func (r *mediaClientRepository) UpdateJellyfin(ctx context.Context, client models.MediaClient[models.JellyfinConfig]) (*models.MediaClient[models.JellyfinConfig], error) {
	return updateClient(r.db, ctx, client, "jellyfin")
}

// CreateEmby creates a new Emby media client
func (r *mediaClientRepository) CreateEmby(ctx context.Context, client models.MediaClient[models.EmbyConfig]) (*models.MediaClient[models.EmbyConfig], error) {
	return createClient(r.db, ctx, client, "emby")
}

// UpdateEmby updates an existing Emby media client
func (r *mediaClientRepository) UpdateEmby(ctx context.Context, client models.MediaClient[models.EmbyConfig]) (*models.MediaClient[models.EmbyConfig], error) {
	return updateClient(r.db, ctx, client, "emby")
}

// CreateNavidrome creates a new Navidrome/Subsonic media client
func (r *mediaClientRepository) CreateNavidrome(ctx context.Context, client models.MediaClient[models.SubsonicConfig]) (*models.MediaClient[models.SubsonicConfig], error) {
	return createClient(r.db, ctx, client, "subsonic")
}

// UpdateNavidrome updates an existing Navidrome/Subsonic media client
func (r *mediaClientRepository) UpdateNavidrome(ctx context.Context, client models.MediaClient[models.SubsonicConfig]) (*models.MediaClient[models.SubsonicConfig], error) {
	return updateClient(r.db, ctx, client, "subsonic")
}

// GetByID retrieves a media client by ID
func (r *mediaClientRepository) GetByID(ctx context.Context, id, userID uint64) (models.MediaClientResponse, error) {
	// First, try to get the basic client info to determine the type
	var baseClient struct {
		ID         uint64                 `json:"id"`
		UserID     uint64                 `json:"userId"`
		ClientType models.MediaClientType `json:"clientType"`
	}

	if err := r.db.WithContext(ctx).Table("media_clients").
		Select("id, user_id, client_type").
		Where("id = ? AND user_id = ?", id, userID).
		First(&baseClient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.MediaClientResponse{}, fmt.Errorf("media client not found")
		}
		return models.MediaClientResponse{}, fmt.Errorf("failed to get media client: %w", err)
	}

	// Based on the client type, get the full client with the proper type
	switch baseClient.ClientType {
	case models.MediaClientTypePlex:
		var client models.MediaClient[models.PlexConfig]
		if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
			return models.MediaClientResponse{}, fmt.Errorf("failed to get plex client: %w", err)
		}
		return models.ToResponse(&client), nil

	case models.MediaClientTypeJellyfin:
		var client models.MediaClient[models.JellyfinConfig]
		if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
			return models.MediaClientResponse{}, fmt.Errorf("failed to get jellyfin client: %w", err)
		}
		return models.ToResponse(&client), nil

	case models.MediaClientTypeEmby:
		var client models.MediaClient[models.EmbyConfig]
		if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
			return models.MediaClientResponse{}, fmt.Errorf("failed to get emby client: %w", err)
		}
		return models.ToResponse(&client), nil

	case models.MediaClientTypeSubsonic:
		var client models.MediaClient[models.SubsonicConfig]
		if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
			return models.MediaClientResponse{}, fmt.Errorf("failed to get subsonic client: %w", err)
		}
		return models.ToResponse(&client), nil

	default:
		return models.MediaClientResponse{}, fmt.Errorf("unsupported client type: %s", baseClient.ClientType)
	}
}

// GetByUserID retrieves all media clients for a user
func (r *mediaClientRepository) GetByUserID(ctx context.Context, userID uint64) ([]models.MediaClientResponse, error) {
	// Get all client IDs and types for this user
	var baseClients []struct {
		ID         uint64                 `json:"id"`
		ClientType models.MediaClientType `json:"clientType"`
	}

	if err := r.db.WithContext(ctx).Table("media_clients").
		Select("id, client_type").
		Where("user_id = ?", userID).
		Find(&baseClients).Error; err != nil {
		return nil, fmt.Errorf("failed to get media clients: %w", err)
	}

	// No clients found
	if len(baseClients) == 0 {
		return []models.MediaClientResponse{}, nil
	}

	// For each client, get the full client with the proper type
	results := make([]models.MediaClientResponse, 0, len(baseClients))

	for _, base := range baseClients {
		client, err := r.GetByID(ctx, base.ID, userID)
		if err != nil {
			// Log the error but continue with other clients
			continue
		}
		results = append(results, client)
	}

	return results, nil
}

// Delete deletes a media client
func (r *mediaClientRepository) Delete(ctx context.Context, id, userID uint64) error {
	// Using the Unscoped raw delete approach to avoid needing a concrete type
	result := r.db.WithContext(ctx).Exec(
		"DELETE FROM media_clients WHERE id = ? AND user_id = ?",
		id, userID,
	)

	if err := result.Error; err != nil {
		return fmt.Errorf("failed to delete media client: %w", err)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("media client not found")
	}

	return nil
}
