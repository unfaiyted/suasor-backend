package repository

import (
	"context"
	"suasor/types/models"
	"suasor/utils"

	"gorm.io/gorm"
)

// CreditRepository defines the interface for credit database operations
type CreditRepository interface {
	// CRUD operations
	Create(ctx context.Context, credit *models.Credit) (*models.Credit, error)
	GetByID(ctx context.Context, id uint64) (*models.Credit, error)
	Update(ctx context.Context, credit *models.Credit) (*models.Credit, error)
	Delete(ctx context.Context, id uint64) error

	// Query operations
	GetByMediaItemID(ctx context.Context, mediaItemID uint64) ([]models.Credit, error)
	GetByPersonID(ctx context.Context, personID uint64) ([]models.Credit, error)
	GetByRole(ctx context.Context, role string) ([]models.Credit, error)
	GetByDepartment(ctx context.Context, department string) ([]models.Credit, error)

	// Cast and crew specific operations
	GetCastForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error)
	GetCrewForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error)

	// Advanced operations
	GetDirectorsForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error)
	GetCreatorsForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error)

	// Bulk operations
	CreateMany(ctx context.Context, credits []models.Credit) ([]models.Credit, error)
}

// creditRepository is a GORM implementation of CreditRepository
type creditRepository struct {
	db *gorm.DB
}

// NewCreditRepository creates a new credit repository
func NewCreditRepository(db *gorm.DB) CreditRepository {
	return &creditRepository{
		db: db,
	}
}

// Create creates a new credit
func (r *creditRepository) Create(ctx context.Context, credit *models.Credit) (*models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	if err := r.db.Create(credit).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create credit")
		return nil, err
	}

	return credit, nil
}

// GetByID gets a credit by ID
func (r *creditRepository) GetByID(ctx context.Context, id uint64) (*models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credit models.Credit
	if err := r.db.First(&credit, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Info().Uint64("id", id).Msg("Credit not found")
			return nil, nil
		}
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get credit by ID")
		return nil, err
	}

	return &credit, nil
}

// Update updates a credit
func (r *creditRepository) Update(ctx context.Context, credit *models.Credit) (*models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	if err := r.db.Save(credit).Error; err != nil {
		log.Error().Err(err).Msg("Failed to update credit")
		return nil, err
	}

	return credit, nil
}

// Delete deletes a credit
func (r *creditRepository) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)

	if err := r.db.Delete(&models.Credit{}, id).Error; err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to delete credit")
		return err
	}

	return nil
}

// GetByMediaItemID gets all credits for a media item
func (r *creditRepository) GetByMediaItemID(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credits []models.Credit
	if err := r.db.Where("media_item_id = ?", mediaItemID).Find(&credits).Error; err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get credits for media item")
		return nil, err
	}

	return credits, nil
}

// GetByPersonID gets all credits for a person
func (r *creditRepository) GetByPersonID(ctx context.Context, personID uint64) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credits []models.Credit
	if err := r.db.Where("person_id = ?", personID).Find(&credits).Error; err != nil {
		log.Error().Err(err).Uint64("personID", personID).Msg("Failed to get credits for person")
		return nil, err
	}

	return credits, nil
}

// GetByRole gets all credits with a specific role
func (r *creditRepository) GetByRole(ctx context.Context, role string) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credits []models.Credit
	if err := r.db.Where("role = ?", role).Find(&credits).Error; err != nil {
		log.Error().Err(err).Str("role", role).Msg("Failed to get credits by role")
		return nil, err
	}

	return credits, nil
}

// GetByDepartment gets all credits with a specific department
func (r *creditRepository) GetByDepartment(ctx context.Context, department string) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credits []models.Credit
	if err := r.db.Where("department = ?", department).Find(&credits).Error; err != nil {
		log.Error().Err(err).Str("department", department).Msg("Failed to get credits by department")
		return nil, err
	}

	return credits, nil
}

// GetCastForMediaItem gets cast credits for a media item
func (r *creditRepository) GetCastForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credits []models.Credit
	if err := r.db.Where("media_item_id = ? AND is_cast = ?", mediaItemID, true).Order("\"order\" ASC").Find(&credits).Error; err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get cast for media item")
		return nil, err
	}

	return credits, nil
}

// GetCrewForMediaItem gets crew credits for a media item
func (r *creditRepository) GetCrewForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credits []models.Credit
	if err := r.db.Where("media_item_id = ? AND is_crew = ?", mediaItemID, true).Order("department ASC, job ASC").Find(&credits).Error; err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get crew for media item")
		return nil, err
	}

	return credits, nil
}

// GetDirectorsForMediaItem gets director credits for a media item
func (r *creditRepository) GetDirectorsForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credits []models.Credit
	if err := r.db.Where("media_item_id = ? AND department = ? AND job = ?", mediaItemID, "Directing", "Director").Find(&credits).Error; err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get directors for media item")
		return nil, err
	}

	return credits, nil
}

// GetCreatorsForMediaItem gets creator credits for a media item
func (r *creditRepository) GetCreatorsForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	var credits []models.Credit
	if err := r.db.Where("media_item_id = ? AND is_creator = ?", mediaItemID, true).Find(&credits).Error; err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get creators for media item")
		return nil, err
	}

	return credits, nil
}

// CreateMany creates multiple credits in a transaction
func (r *creditRepository) CreateMany(ctx context.Context, credits []models.Credit) ([]models.Credit, error) {
	log := utils.LoggerFromContext(ctx)

	if len(credits) == 0 {
		return []models.Credit{}, nil
	}

	// Use a transaction for bulk operations
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil, err
	}

	// Create credits in the transaction
	if err := tx.Create(&credits).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to create credits in transaction")
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil, err
	}

	return credits, nil
}

