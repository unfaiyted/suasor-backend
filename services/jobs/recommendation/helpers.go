package recommendation

import (
	"context"
	"fmt"
	"sort"
	"suasor/types/models"
)

// GetCastFromCredits extracts cast members from Credits, limited to maxCount
func GetCastFromCredits(credits []*models.Credit, maxCount int) []*models.Credit {
	var cast []*models.Credit

	// Get all cast members
	for _, credit := range credits {
		if credit.IsCast {
			cast = append(cast, credit)
		}
	}

	// Sort by order if available
	sort.Slice(cast, func(i, j int) bool {
		return cast[i].Order < cast[j].Order
	})

	// Limit to maxCount
	if len(cast) > maxCount {
		cast = cast[:maxCount]
	}

	return cast
}

// GetCrewByRole extracts crew members with a specific role
func GetCrewByRole(credits []*models.Credit, role models.MediaRole) []*models.Credit {
	var result []*models.Credit

	for _, credit := range credits {
		if credit.IsCrew && credit.Role == role {
			result = append(result, credit)
		}
	}

	return result
}

// GetCrewByDepartment extracts crew members from a specific department
func GetCrewByDepartment(credits []*models.Credit, department models.MediaDepartment) []*models.Credit {
	var result []*models.Credit

	for _, credit := range credits {
		if credit.IsCrew && credit.Department == department {
			result = append(result, credit)
		}
	}

	return result
}

// GetCreatorsFromCredits extracts creators from Credits
func GetCreatorsFromCredits(credits []*models.Credit) []*models.Credit {
	var creators []*models.Credit

	for _, credit := range credits {
		if credit.IsCreator {
			creators = append(creators, credit)
		}
	}

	return creators
}

// ExtractNamesFromCredits extracts just the names from a list of credits
func ExtractNamesFromCredits(credits []*models.Credit) []string {
	var names []string

	for _, credit := range credits {
		names = append(names, credit.Name)
	}

	return names
}

// GetPeopleByRole retrieves people from the repository who have a specific role
func (j *RecommendationJob) GetPeopleByRole(ctx context.Context, role models.MediaRole) ([]*models.Person, error) {
	// If we don't have a people repository, return an error
	if j.peopleRepo == nil {
		return nil, fmt.Errorf("people repository not available")
	}

	// Get people by role
	people, err := j.peopleRepo.GetByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	return people, nil
}

// GetPersonByID retrieves a person by ID
func (j *RecommendationJob) GetPersonByID(ctx context.Context, personID uint64) (*models.Person, error) {
	// If we don't have a people repository, return an error
	if j.peopleRepo == nil {
		return nil, fmt.Errorf("people repository not available")
	}

	// Get person by ID
	person, err := j.peopleRepo.GetByID(ctx, personID)
	if err != nil {
		return nil, err
	}

	return person, nil
}

// getCreditsForMediaItem retrieves all credits for a given media item
func (j *RecommendationJob) getCreditsForMediaItem(ctx context.Context, mediaItemID uint64) ([]*models.Credit, error) {
	// If we don't have a credit repository, return an error
	if j.creditRepo == nil {
		return nil, fmt.Errorf("credit repository not available")
	}

	// Get credits from the repository
	credits, err := j.creditRepo.GetByMediaItemID(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}

	return credits, nil
}

