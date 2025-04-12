package repository

import (
	"context"
	mediatypes "suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"

	"fmt"
	"gorm.io/gorm"
)

// PersonRepository defines the interface for person database operations
type PersonRepository interface {
	// CRUD operations
	Create(ctx context.Context, person *models.Person) (*models.Person, error)
	GetByID(ctx context.Context, id uint64) (*models.Person, error)
	Update(ctx context.Context, person *models.Person) (*models.Person, error)
	Delete(ctx context.Context, id uint64) error

	// Query operations
	GetAll(ctx context.Context, limit, offset int) ([]models.Person, error)
	GetByName(ctx context.Context, name string) ([]models.Person, error)
	GetByRole(ctx context.Context, role string) ([]models.Person, error)
	GetByExternalID(ctx context.Context, source, id string) (*models.Person, error)
	SearchByName(ctx context.Context, name string, limit int) ([]models.Person, error)

	// Advanced operations
	GetPopular(ctx context.Context, limit int) ([]models.Person, error)
	GetWithMostCredits(ctx context.Context, limit int) ([]models.Person, error)

	Search(ctx context.Context, options mediatypes.QueryOptions) ([]*models.Person, error)
}

// PersonRepository is a GORM implementation of PersonRepository
type personRepository struct {
	db *gorm.DB
}

// NewPersonRepository creates a new person repository
func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepository{
		db: db,
	}
}

// Create creates a new person
func (r *personRepository) Create(ctx context.Context, person *models.Person) (*models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	if err := r.db.Create(person).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create person")
		return nil, err
	}

	return person, nil
}

// GetByID gets a person by ID
func (r *personRepository) GetByID(ctx context.Context, id uint64) (*models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	var person models.Person
	if err := r.db.First(&person, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Info().Uint64("id", id).Msg("Person not found")
			return nil, nil
		}
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get person by ID")
		return nil, err
	}

	return &person, nil
}

// Update updates a person
func (r *personRepository) Update(ctx context.Context, person *models.Person) (*models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	if err := r.db.Save(person).Error; err != nil {
		log.Error().Err(err).Msg("Failed to update person")
		return nil, err
	}

	return person, nil
}

// Delete deletes a person
func (r *personRepository) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)

	if err := r.db.Delete(&models.Person{}, id).Error; err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to delete person")
		return err
	}

	return nil
}

// GetAll gets all people with pagination
func (r *personRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	var people []models.Person
	if err := r.db.Limit(limit).Offset(offset).Find(&people).Error; err != nil {
		log.Error().Err(err).Int("limit", limit).Int("offset", offset).Msg("Failed to get all people")
		return nil, err
	}

	return people, nil
}

// GetByName gets all people with a specific name
func (r *personRepository) GetByName(ctx context.Context, name string) ([]models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	var people []models.Person
	if err := r.db.Where("name = ?", name).Find(&people).Error; err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to get people by name")
		return nil, err
	}

	return people, nil
}

// GetByRole gets all people with a specific role
func (r *personRepository) GetByRole(ctx context.Context, role string) ([]models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	var people []models.Person
	if err := r.db.Where("known_for = ?", role).Find(&people).Error; err != nil {
		log.Error().Err(err).Str("role", role).Msg("Failed to get people by role")
		return nil, err
	}

	return people, nil
}

// GetByExternalID gets a person by external ID
func (r *personRepository) GetByExternalID(ctx context.Context, source, id string) (*models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	var people []models.Person
	if err := r.db.Find(&people).Error; err != nil {
		log.Error().Err(err).Str("source", source).Str("id", id).Msg("Failed to get people for external ID search")
		return nil, err
	}

	// Search through external IDs
	for _, person := range people {
		for _, extID := range person.ExternalIDs {
			if extID.Source == source && extID.ID == id {
				return &person, nil
			}
		}
	}

	return nil, nil
}

// SearchByName searches for people by name with a LIKE query
func (r *personRepository) SearchByName(ctx context.Context, name string, limit int) ([]models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	var people []models.Person
	if err := r.db.Where("name ILIKE ?", "%"+name+"%").Limit(limit).Find(&people).Error; err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to search people by name")
		return nil, err
	}

	return people, nil
}

// GetPopular gets the most popular people
func (r *personRepository) GetPopular(ctx context.Context, limit int) ([]models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	var people []models.Person
	if err := r.db.Order("popularity DESC").Limit(limit).Find(&people).Error; err != nil {
		log.Error().Err(err).Int("limit", limit).Msg("Failed to get popular people")
		return nil, err
	}

	return people, nil
}

// GetWithMostCredits gets people with the most credits
func (r *personRepository) GetWithMostCredits(ctx context.Context, limit int) ([]models.Person, error) {
	log := utils.LoggerFromContext(ctx)

	var people []models.Person
	if err := r.db.Joins("LEFT JOIN credits ON people.id = credits.person_id").
		Group("people.id").
		Order("COUNT(credits.id) DESC").
		Limit(limit).
		Find(&people).Error; err != nil {
		log.Error().Err(err).Int("limit", limit).Msg("Failed to get people with most credits")
		return nil, err
	}

	return people, nil
}

func (r *personRepository) Search(ctx context.Context, options mediatypes.QueryOptions) ([]*models.Person, error) {
	var people []*models.Person
	log := utils.LoggerFromContext(ctx)

	log.Info().Str("query", options.Query).Int("limit", options.Limit).Int("offset", options.Offset).Msg("Searching people")

	query := r.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+options.Query+"%")

	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	if err := query.Find(&people).Error; err != nil {
		return nil, fmt.Errorf("failed to search people: %w", err)
	}

	return people, nil

}
