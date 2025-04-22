package services

import (
	"context"
	"fmt"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// PersonService handles operations on people
type PersonService struct {
	personRepo repository.PersonRepository
	creditRepo repository.CreditRepository
}

// NewPersonService creates a new person service
func NewPersonService(personRepo repository.PersonRepository, creditRepo repository.CreditRepository) *PersonService {
	return &PersonService{
		personRepo: personRepo,
		creditRepo: creditRepo,
	}
}

// GetPersonByID gets a person by ID
func (s *PersonService) GetPersonByID(ctx context.Context, id uint64) (*models.Person, error) {
	log := logger.LoggerFromContext(ctx)

	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get person")
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	return person, nil
}

// GetPersonWithCredits gets a person with their credits
func (s *PersonService) GetPersonWithCredits(ctx context.Context, id uint64) (*models.Person, []models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	// Get the person
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get person")
		return nil, nil, fmt.Errorf("failed to get person: %w", err)
	}

	if person == nil {
		return nil, nil, nil
	}

	// Get the person's credits
	credits, err := s.creditRepo.GetByPersonID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get credits for person")
		return person, nil, fmt.Errorf("failed to get credits: %w", err)
	}

	return person, credits, nil
}

// SearchPeople searches for people by name
func (s *PersonService) SearchPeople(ctx context.Context, query string, limit int) ([]models.Person, error) {
	log := logger.LoggerFromContext(ctx)

	people, err := s.personRepo.SearchByName(ctx, query, limit)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("Failed to search people")
		return nil, fmt.Errorf("failed to search people: %w", err)
	}

	return people, nil
}

// GetPopularPeople gets popular people
func (s *PersonService) GetPopularPeople(ctx context.Context, limit int) ([]models.Person, error) {
	log := logger.LoggerFromContext(ctx)

	people, err := s.personRepo.GetPopular(ctx, limit)
	if err != nil {
		log.Error().Err(err).Int("limit", limit).Msg("Failed to get popular people")
		return nil, fmt.Errorf("failed to get popular people: %w", err)
	}

	return people, nil
}

// GetPeopleByRole gets people by role
func (s *PersonService) GetPeopleByRole(ctx context.Context, role string) ([]models.Person, error) {
	log := logger.LoggerFromContext(ctx)

	people, err := s.personRepo.GetByRole(ctx, role)
	if err != nil {
		log.Error().Err(err).Str("role", role).Msg("Failed to get people by role")
		return nil, fmt.Errorf("failed to get people by role: %w", err)
	}

	return people, nil
}

// CreatePerson creates a new person
func (s *PersonService) CreatePerson(ctx context.Context, person *models.Person) (*models.Person, error) {
	log := logger.LoggerFromContext(ctx)

	// Create the person
	createdPerson, err := s.personRepo.Create(ctx, person)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create person")
		return nil, fmt.Errorf("failed to create person: %w", err)
	}

	return createdPerson, nil
}

// UpdatePerson updates an existing person
func (s *PersonService) UpdatePerson(ctx context.Context, person *models.Person) (*models.Person, error) {
	log := logger.LoggerFromContext(ctx)

	// Verify the person exists
	existingPerson, err := s.personRepo.GetByID(ctx, person.ID)
	if err != nil {
		log.Error().Err(err).Uint64("id", person.ID).Msg("Failed to get person")
		return nil, fmt.Errorf("failed to verify person: %w", err)
	}

	if existingPerson == nil {
		log.Error().Uint64("id", person.ID).Msg("Person not found")
		return nil, fmt.Errorf("person not found")
	}

	// Update the person
	updatedPerson, err := s.personRepo.Update(ctx, person)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update person")
		return nil, fmt.Errorf("failed to update person: %w", err)
	}

	return updatedPerson, nil
}

// DeletePerson deletes a person
func (s *PersonService) DeletePerson(ctx context.Context, id uint64) error {
	log := logger.LoggerFromContext(ctx)

	// Verify the person exists
	person, err := s.personRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get person")
		return fmt.Errorf("failed to verify person: %w", err)
	}

	if person == nil {
		log.Error().Uint64("id", id).Msg("Person not found")
		return fmt.Errorf("person not found")
	}

	// Delete associated credits
	credits, err := s.creditRepo.GetByPersonID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get credits for person")
		return fmt.Errorf("failed to get credits: %w", err)
	}

	for _, credit := range credits {
		if err := s.creditRepo.Delete(ctx, credit.ID); err != nil {
			log.Error().Err(err).Uint64("creditID", credit.ID).Msg("Failed to delete credit")
			// Continue with other credits even if one fails
		}
	}

	// Delete the person
	if err := s.personRepo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Msg("Failed to delete person")
		return fmt.Errorf("failed to delete person: %w", err)
	}

	return nil
}

// ImportPerson imports a person from an external source
func (s *PersonService) ImportPerson(ctx context.Context, source string, externalID string, personData *models.Person) (*models.Person, error) {
	log := logger.LoggerFromContext(ctx)

	// Check if the person already exists by external ID
	existingPerson, err := s.personRepo.GetByExternalID(ctx, source, externalID)
	if err != nil {
		log.Error().Err(err).Str("source", source).Str("externalID", externalID).Msg("Failed to check for existing person")
		return nil, fmt.Errorf("failed to check for existing person: %w", err)
	}

	// If the person already exists, update it
	if existingPerson != nil {
		log.Info().Uint64("id", existingPerson.ID).Msg("Person already exists, updating")

		// Update with new data but keep the ID
		personData.ID = existingPerson.ID

		// Add the external ID if not already present
		existingID := false
		for _, id := range existingPerson.ExternalIDs {
			if id.Source == source && id.ID == externalID {
				existingID = true
				break
			}
		}

		if !existingID {
			personData.ExternalIDs = append(personData.ExternalIDs, models.ExternalID{
				Source: source,
				ID:     externalID,
			})
		}

		// Update the person
		return s.UpdatePerson(ctx, personData)
	}

	// Add the external ID
	personData.ExternalIDs = append(personData.ExternalIDs, models.ExternalID{
		Source: source,
		ID:     externalID,
	})

	// Create a new person
	return s.CreatePerson(ctx, personData)
}

// GetPersonCreditsGrouped gets a person's credits grouped by type
func (s *PersonService) GetPersonCreditsGrouped(ctx context.Context, id uint64) (map[string][]models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	// Get the person's credits
	credits, err := s.creditRepo.GetByPersonID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get credits for person")
		return nil, fmt.Errorf("failed to get credits: %w", err)
	}

	// Group by type
	result := make(map[string][]models.Credit)

	// Initialize with known types
	result["cast"] = []models.Credit{}
	result["crew"] = []models.Credit{}
	result["directing"] = []models.Credit{}
	result["writing"] = []models.Credit{}
	result["producing"] = []models.Credit{}

	for _, credit := range credits {
		if credit.IsCast {
			result["cast"] = append(result["cast"], credit)
		}

		if credit.IsCrew {
			result["crew"] = append(result["crew"], credit)

			// Also add to specific departments
			switch credit.Department {
			case "Directing":
				result["directing"] = append(result["directing"], credit)
			case "Writing":
				result["writing"] = append(result["writing"], credit)
			case "Production":
				result["producing"] = append(result["producing"], credit)
			}
		}
	}

	return result, nil
}

// AddExternalIDToPerson adds an external ID to a person
func (s *PersonService) AddExternalIDToPerson(ctx context.Context, personID uint64, source string, id string) error {
	log := logger.LoggerFromContext(ctx)

	// Get the person
	person, err := s.personRepo.GetByID(ctx, personID)
	if err != nil {
		log.Error().Err(err).Uint64("id", personID).Msg("Failed to get person")
		return fmt.Errorf("failed to get person: %w", err)
	}

	if person == nil {
		log.Error().Uint64("id", personID).Msg("Person not found")
		return fmt.Errorf("person not found")
	}

	// Add the external ID
	person.AddExternalID(source, id)

	// Update the person
	_, err = s.personRepo.Update(ctx, person)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update person with external ID")
		return fmt.Errorf("failed to update person: %w", err)
	}

	return nil
}

