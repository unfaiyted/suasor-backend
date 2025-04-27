package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// CoreListRepository defines the interface for list operations
// This focuses on generic operations that are not specifically tied to clients or users.
// It provides the core functionality for working with lists directly.
//
// Relationships with other repositories:
// - MediaItemRepository: Core operations on media items without client or user associations
// - UserMediaItemRepository: Operations for media items owned by users (playlists, collections)
//
// This three-tier approach allows for clear separation of concerns while maintaining
// a single database table for all media items and lists.

type CoreListRepository[T types.ListData] interface {
	CoreMediaItemRepository[T]

	GetCollaborator(ctx context.Context, listID uint64, collaboratorID uint64) (*models.ListCollaborator, error)
	GetCollaborators(ctx context.Context, listID uint64) ([]*models.ListCollaborator, error)
	RemoveCollaborator(ctx context.Context, listID uint64, collaboratorID uint64) error
	AddCollaborator(ctx context.Context, listID uint64, collaboratorID uint64, permission models.CollaboratorPermission) error
}

// coreListRepository implements the CoreListRepository interface
type coreListRepository[T types.ListData] struct {
	CoreMediaItemRepository[T]
	db *gorm.DB
}

// NewCoreListRepository creates a new list repository
func NewCoreListRepository[T types.ListData](

	db *gorm.DB,
	mediaItemRepo CoreMediaItemRepository[T]) CoreListRepository[T] {
	return &coreListRepository[T]{
		db:                      db,
		CoreMediaItemRepository: mediaItemRepo,
	}
}
func (r *coreListRepository[T]) GetCollaborator(ctx context.Context, listID uint64, collaboratorID uint64) (*models.ListCollaborator, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("collaboratorID", collaboratorID).
		Msg("Getting list collaborator")

	var collaborator models.ListCollaborator

	// Query collaborators table by list ID and collaborator ID
	if err := r.db.WithContext(ctx).
		Where("list_id = ? AND user_id = ?", listID, collaboratorID).
		First(&collaborator).Error; err != nil {
		return nil, fmt.Errorf("failed to get list collaborator: %w", err)
	}

	return &collaborator, nil
}

func (r *coreListRepository[T]) GetCollaborators(ctx context.Context, listID uint64) ([]*models.ListCollaborator, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Msg("Getting list collaborators")

	var collaborators []*models.ListCollaborator

	// Query collaborators table by list ID
	if err := r.db.WithContext(ctx).
		Where("list_id = ?", listID).
		Find(&collaborators).Error; err != nil {
		return nil, fmt.Errorf("failed to get list collaborators: %w", err)
	}

	return collaborators, nil

}
func (r *coreListRepository[T]) AddCollaborator(ctx context.Context, listID uint64, collaboratorID uint64, permission models.CollaboratorPermission) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("collaboratorID", collaboratorID).
		Str("permission", string(permission)).
		Msg("Adding list collaborator")

	// Create a new collaborator
	collaborator := models.ListCollaborator{
		ListID:     listID,
		UserID:     collaboratorID,
		Permission: permission,
	}

	// Insert the collaborator
	if err := r.db.WithContext(ctx).Create(&collaborator).Error; err != nil {
		return fmt.Errorf("failed to add list collaborator: %w", err)
	}

	return nil
}

func (r *coreListRepository[T]) RemoveCollaborator(ctx context.Context, listID uint64, collaboratorID uint64) error {
	log := logger.LoggerFromContext(ctx)

	log.Debug().
		Uint64("listID", listID).
		Uint64("collaboratorID", collaboratorID).
		Msg("Removing list collaborator")

	// Delete the collaborator
	if err := r.db.WithContext(ctx).
		Where("list_id = ? AND user_id = ?", listID, collaboratorID).
		Delete(&models.ListCollaborator{}).Error; err != nil {
		return fmt.Errorf("failed to remove list collaborator: %w", err)
	}

	return nil
}
