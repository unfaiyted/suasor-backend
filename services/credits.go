package services

import (
	"context"
	"fmt"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// CreditService handles operations on credits
type CreditService struct {
	creditRepo repository.CreditRepository
	personRepo repository.PersonRepository
}

// NewCreditService creates a new credit service
func NewCreditService(creditRepo repository.CreditRepository, personRepo repository.PersonRepository) *CreditService {
	return &CreditService{
		creditRepo: creditRepo,
		personRepo: personRepo,
	}
}

// GetCreditsForMediaItem gets all credits for a media item
func (s *CreditService) GetCreditsForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	credits, err := s.creditRepo.GetByMediaItemID(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get credits for media item")
		return nil, fmt.Errorf("failed to get credits: %w", err)
	}

	return credits, nil
}

// GetCastForMediaItem gets cast credits for a media item
func (s *CreditService) GetCastForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	credits, err := s.creditRepo.GetCastForMediaItem(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get cast for media item")
		return nil, fmt.Errorf("failed to get cast: %w", err)
	}

	return credits, nil
}

// GetCrewForMediaItem gets crew credits for a media item
func (s *CreditService) GetCrewForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	credits, err := s.creditRepo.GetCrewForMediaItem(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get crew for media item")
		return nil, fmt.Errorf("failed to get crew: %w", err)
	}

	return credits, nil
}

// GetCrewByDepartment gets crew credits for a media item filtered by department
func (s *CreditService) GetCrewByDepartment(ctx context.Context, mediaItemID uint64, department string) ([]models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	credits, err := s.creditRepo.GetCrewForMediaItem(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get crew for media item")
		return nil, fmt.Errorf("failed to get crew: %w", err)
	}

	// Filter by department
	var filteredCredits []models.Credit
	for _, credit := range credits {
		if credit.Department == department {
			filteredCredits = append(filteredCredits, credit)
		}
	}

	return filteredCredits, nil
}

// GetDirectorsForMediaItem gets director credits for a media item
func (s *CreditService) GetDirectorsForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	credits, err := s.creditRepo.GetDirectorsForMediaItem(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get directors for media item")
		return nil, fmt.Errorf("failed to get directors: %w", err)
	}

	return credits, nil
}

// GetCreditsByPerson gets all credits for a person
func (s *CreditService) GetCreditsByPerson(ctx context.Context, personID uint64) ([]models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	credits, err := s.creditRepo.GetByPersonID(ctx, personID)
	if err != nil {
		log.Error().Err(err).Uint64("personID", personID).Msg("Failed to get credits for person")
		return nil, fmt.Errorf("failed to get credits: %w", err)
	}

	return credits, nil
}

// CreateCredit creates a new credit
func (s *CreditService) CreateCredit(ctx context.Context, credit *models.Credit) (*models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	// Verify the person exists
	person, err := s.personRepo.GetByID(ctx, credit.PersonID)
	if err != nil {
		log.Error().Err(err).Uint64("personID", credit.PersonID).Msg("Failed to get person for credit")
		return nil, fmt.Errorf("failed to verify person: %w", err)
	}

	if person == nil {
		log.Error().Uint64("personID", credit.PersonID).Msg("Person not found for credit")
		return nil, fmt.Errorf("person not found")
	}

	// Create the credit
	createdCredit, err := s.creditRepo.Create(ctx, credit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create credit")
		return nil, fmt.Errorf("failed to create credit: %w", err)
	}

	return createdCredit, nil
}

// UpdateCredit updates an existing credit
func (s *CreditService) UpdateCredit(ctx context.Context, credit *models.Credit) (*models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	// Verify the credit exists
	existingCredit, err := s.creditRepo.GetByID(ctx, credit.ID)
	if err != nil {
		log.Error().Err(err).Uint64("id", credit.ID).Msg("Failed to get credit")
		return nil, fmt.Errorf("failed to verify credit: %w", err)
	}

	if existingCredit == nil {
		log.Error().Uint64("id", credit.ID).Msg("Credit not found")
		return nil, fmt.Errorf("credit not found")
	}

	// Update the credit
	updatedCredit, err := s.creditRepo.Update(ctx, credit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update credit")
		return nil, fmt.Errorf("failed to update credit: %w", err)
	}

	return updatedCredit, nil
}

// DeleteCredit deletes a credit
func (s *CreditService) DeleteCredit(ctx context.Context, id uint64) error {
	log := logger.LoggerFromContext(ctx)

	// Verify the credit exists
	credit, err := s.creditRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to get credit")
		return fmt.Errorf("failed to verify credit: %w", err)
	}

	if credit == nil {
		log.Error().Uint64("id", id).Msg("Credit not found")
		return fmt.Errorf("credit not found")
	}

	// Delete the credit
	if err := s.creditRepo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Msg("Failed to delete credit")
		return fmt.Errorf("failed to delete credit: %w", err)
	}

	return nil
}

// CreateCreditsForMediaItem creates multiple credits for a media item
func (s *CreditService) CreateCreditsForMediaItem(ctx context.Context, mediaItemID uint64, credits []models.Credit) ([]models.Credit, error) {
	log := logger.LoggerFromContext(ctx)

	// Set media item ID for all credits
	for i := range credits {
		credits[i].MediaItemID = mediaItemID
	}

	// Create credits
	createdCredits, err := s.creditRepo.CreateMany(ctx, credits)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Int("count", len(credits)).Msg("Failed to create credits for media item")
		return nil, fmt.Errorf("failed to create credits: %w", err)
	}

	return createdCredits, nil
}

// TabularizeCredits organizes credits by department and role
func (s *CreditService) TabularizeCredits(ctx context.Context, credits []models.Credit) map[string]map[string][]models.Credit {
	result := make(map[string]map[string][]models.Credit)

	for _, credit := range credits {
		// Determine department
		dept := credit.Department
		if dept == "" {
			if credit.IsCast {
				dept = "Cast"
			} else if credit.IsCrew {
				dept = "Crew"
			} else {
				dept = "Other"
			}
		}

		// Determine role
		role := credit.Role
		if role == "" {
			if credit.IsCast {
				role = "Actor"
			} else {
				role = "Unknown"
			}
		}

		// Add to result
		if _, ok := result[dept]; !ok {
			result[dept] = make(map[string][]models.Credit)
		}

		result[dept][role] = append(result[dept][role], credit)
	}

	return result
}

