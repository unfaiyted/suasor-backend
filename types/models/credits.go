package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"suasor/clients/media/types"
)

type MediaDepartment string

const (
	DepartmentCast       MediaDepartment = "Cast"
	DepartmentCrew       MediaDepartment = "Crew"
	DepartmentDirecting  MediaDepartment = "Directing"
	DepartmentWriting    MediaDepartment = "Writing"
	DepartmentProduction MediaDepartment = "Production"
	DepartmentCamera     MediaDepartment = "Camera"
	DepartmentEditing    MediaDepartment = "Editing"
	DepartmentSound      MediaDepartment = "Sound"
	DepartmentArt        MediaDepartment = "Art"
	DepartmentOther      MediaDepartment = "Other"
)

type MediaRole string

const (
	// Cast roles
	RoleActor MediaRole = "Actor"
	RoleVoice MediaRole = "Voice"

	// Directing roles
	RoleDirector MediaRole = "Director"

	// Writing roles
	RoleWriter     MediaRole = "Writer"
	RoleScreenplay MediaRole = "Screenplay"
	RoleStory      MediaRole = "Story"

	// Production roles
	RoleProducer          MediaRole = "Producer"
	RoleExecutiveProducer MediaRole = "Executive Producer"

	// Add other roles as needed
	RoleOther MediaRole = "Other"
)

// RoleToDepartment maps specific roles to their departments
var RoleToDepartment = map[MediaRole]MediaDepartment{
	RoleActor:             DepartmentCast,
	RoleVoice:             DepartmentCast,
	RoleDirector:          DepartmentDirecting,
	RoleWriter:            DepartmentWriting,
	RoleScreenplay:        DepartmentWriting,
	RoleStory:             DepartmentWriting,
	RoleProducer:          DepartmentProduction,
	RoleExecutiveProducer: DepartmentProduction,
	// Add other mappings
}

type PersonCreditsByRole struct {
	Person *Person
	// Credits is a map of credits grouped by role
	Credits map[MediaRole][]*Credit
}

// Credit represents a person's involvement with a particular media item
type Credit struct {
	BaseModel                               // Include base fields (ID, timestamps)
	PersonID    uint64                      `json:"personID" gorm:"index;not null"`
	Person      Person                      `json:"person,omitempty" gorm:"foreignKey:PersonID"`
	MediaItemID uint64                      `json:"mediaItemID" gorm:"index;not null"`
	MediaItem   *MediaItem[types.MediaData] `json:"-" gorm:"foreignKey:MediaItemID"` // Use pointer to avoid recursion issues

	Name         string          `json:"name" gorm:"type:varchar(255)"`                 // Name as it appears in the credits
	Role         MediaRole       `json:"role,omitempty" gorm:"type:varchar(100);index"` // e.g., "Director", "Actor"
	Character    string          `json:"character,omitempty" gorm:"type:varchar(255)"`  // For actors
	Department   MediaDepartment `json:"department,omitempty" gorm:"type:varchar(100)"` // e.g., "Directing", "Writing", "Sound"
	Job          string          `json:"job,omitempty" gorm:"type:varchar(100)"`        // Specific job title
	Order        int             `json:"order,omitempty"`                               // Display order in credits
	SeasonNumber int             `json:"seasonNumber,omitempty"`                        // For TV series credits
	EpisodeCount int             `json:"episodeCount,omitempty"`                        // Number of episodes for TV series

	// Credit type flags
	IsCast    bool `json:"isCast,omitempty" gorm:"index"`
	IsCrew    bool `json:"isCrew,omitempty" gorm:"index"`
	IsGuest   bool `json:"isGuest,omitempty"`
	IsCreator bool `json:"isCreator,omitempty"`
	IsArtist  bool `json:"isArtist,omitempty"`

	// Credit metadata (awards, notes, etc.)
	Metadata CreditMetadata `json:"metadata,omitempty" gorm:"type:jsonb"`
}

// CreditMetadata contains additional information about a credit
type CreditMetadata struct {
	Notes              string         `json:"notes,omitempty"`
	Awards             []CreditAward  `json:"awards,omitempty"`
	Uncredited         bool           `json:"uncredited,omitempty"`
	VoiceOnly          bool           `json:"voiceOnly,omitempty"`
	SpecialPerformance bool           `json:"specialPerformance,omitempty"`
	AdditionalInfo     map[string]any `json:"additionalInfo,omitempty"`
}

// Value implements the driver.Valuer interface for database serialization
func (cm CreditMetadata) Value() (driver.Value, error) {
	return json.Marshal(cm)
}

// Scan implements the sql.Scanner interface for database deserialization
func (cm *CreditMetadata) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &cm)
}

// CreditAward represents an award for a specific credit
type CreditAward struct {
	Name      string `json:"name"`
	Year      int    `json:"year"`
	Category  string `json:"category,omitempty"`
	IsWinner  bool   `json:"isWinner"`
	IsNominee bool   `json:"isNominee"`
}

// Credits is a collection of Credit objects
type Credits []Credit

// GetCastCredits returns all cast credits
func (c Credits) GetCast() Credits {
	var cast Credits
	for _, credit := range c {
		if credit.IsCast {
			cast = append(cast, credit)
		}
	}
	return cast
}

// GetCrewCredits returns all crew credits
func (c Credits) GetCrew() Credits {
	var crew Credits
	for _, credit := range c {
		if credit.IsCrew {
			crew = append(crew, credit)
		}
	}
	return crew
}

// GetGuestCredits returns all guest credits
func (c Credits) GetGuests() Credits {
	var guests Credits
	for _, credit := range c {
		if credit.IsGuest {
			guests = append(guests, credit)
		}
	}
	return guests
}

// GetCreatorCredits returns all creator credits
func (c Credits) GetCreators() Credits {
	var creators Credits
	for _, credit := range c {
		if credit.IsCreator {
			creators = append(creators, credit)
		}
	}
	return creators
}

// GetByDepartment returns credits filtered by department
func (c Credits) GetByDepartment(department MediaDepartment) Credits {
	var results Credits
	for _, credit := range c {
		if credit.Department == department {
			results = append(results, credit)
		}
	}
	return results
}

// GetByRole returns credits filtered by role
func (c Credits) GetByRole(role MediaRole) Credits {
	var results Credits
	for _, credit := range c {
		if credit.Role == role {
			results = append(results, credit)
		}
	}
	return results
}

// GetTVEpisodeCredits returns credits for specific TV episode
func (c Credits) GetDirectors() Credits {
	return c.GetByRole("Director")
}

// GetWriters returns all writing credits
func (c Credits) GetWriters() Credits {
	var writers Credits
	for _, credit := range c {
		if credit.Department == "Writing" || credit.Role == "Writer" || credit.Role == "Screenplay" {
			writers = append(writers, credit)
		}
	}
	return writers
}

// GetCreditPublicView returns a view of the credit for API responses
func (c *Credit) GetCreditPublicView() map[string]any {
	result := map[string]any{
		"id":       c.ID,
		"name":     c.Name,
		"role":     c.Role,
		"personID": c.PersonID,
	}

	if c.IsCast {
		result["character"] = c.Character
	}

	if c.Department != "" {
		result["department"] = c.Department
	}

	if c.Job != "" {
		result["job"] = c.Job
	}

	if c.Order > 0 {
		result["order"] = c.Order
	}

	return result
}

// TabularizedCredits returns credits organized by department and role
func TabularizedCredits(credits Credits) map[MediaDepartment]map[MediaRole][]Credit {
	result := make(map[MediaDepartment]map[MediaRole][]Credit)

	for _, credit := range credits {
		dept := credit.Department
		if dept == "" {
			if credit.IsCast {
				dept = DepartmentCast
			} else if credit.IsCrew {
				dept = DepartmentCrew
			} else {
				dept = DepartmentOther
			}
		}

		role := credit.Role
		if role == "" {
			role = RoleOther
		}

		if _, exists := result[dept]; !exists {
			result[dept] = make(map[MediaRole][]Credit)
		}

		result[dept][role] = append(result[dept][role], credit)
	}

	return result
}

// NewCredit creates a new credit with the given person and media item
func NewCredit(personID, mediaItemID uint64, name string, role MediaRole, isCast bool) *Credit {
	credit := &Credit{
		PersonID:    personID,
		MediaItemID: mediaItemID,
		Name:        name,
		Role:        role,
		IsCast:      isCast,
		IsCrew:      !isCast, // Default to crew if not cast
	}

	// Set department based on role
	switch role {
	case "Director":
		credit.Department = "Directing"
	case "Writer", "Screenplay", "Story":
		credit.Department = "Writing"
	case "Producer", "Executive Producer":
		credit.Department = "Production"
	case "Cinematographer", "Director of Photography":
		credit.Department = "Camera"
	case "Editor":
		credit.Department = "Editing"
	case "Composer", "Original Music Composer":
		credit.Department = "Sound"
	case "Actor", "Actress":
		credit.Department = "Cast"
		credit.IsCast = true
		credit.IsCrew = false
	}

	return credit
}

// NewCastCredit creates a new credit for a cast member
func NewCastCredit(personID, mediaItemID uint64, name string, character string, order int) *Credit {
	credit := &Credit{
		PersonID:    personID,
		MediaItemID: mediaItemID,
		Name:        name,
		Role:        "Actor",
		Character:   character,
		Order:       order,
		Department:  "Cast",
		IsCast:      true,
		IsCrew:      false,
	}
	return credit
}

// NewCrewCredit creates a new credit for a crew member
func NewCrewCredit(personID, mediaItemID uint64, name string, department MediaDepartment, job string) *Credit {
	credit := &Credit{
		PersonID:    personID,
		MediaItemID: mediaItemID,
		Name:        name,
		Department:  department,
		Job:         job,
		IsCast:      false,
		IsCrew:      true,
	}

	// Set appropriate role based on job
	switch job {
	case "Director":
		credit.Role = "Director"
	case "Screenplay", "Writer":
		credit.Role = RoleWriter
	case "Producer", "Executive Producer":
		credit.Role = RoleProducer
	case "Director of Photography", "Cinematographer":
		credit.Role = "Cinematographer"
	default:
		credit.Role = RoleWriter
	}

	return credit
}

// GetDepartmentForRole returns the appropriate department for a given role
func GetDepartmentForRole(role MediaRole) MediaDepartment {
	if dept, exists := RoleToDepartment[role]; exists {
		return dept
	}
	return DepartmentCrew // Default department
}

// ToTableFormatted returns credits formatted for display in a table
func (c Credits) ToTableFormatted() []map[string]any {
	result := make([]map[string]any, 0, len(c))

	for _, credit := range c {
		item := map[string]any{
			"id":       credit.ID,
			"name":     credit.Name,
			"role":     credit.Role,
			"personID": credit.PersonID,
		}

		if credit.Department != "" {
			item["department"] = credit.Department
		}

		if credit.IsCast && credit.Character != "" {
			item["character"] = credit.Character
		}

		if credit.Job != "" {
			item["job"] = credit.Job
		}

		result = append(result, item)
	}

	return result
}

// TableName specifies the database table name for Credit
func (Credit) TableName() string {
	return "credits"
}

// GetCreditWithoutPerson returns a copy of the credit with the person field set to nil
// This helps avoid circular JSON serialization issues
func (c *Credit) GetCreditWithoutPerson() Credit {
	creditCopy := *c
	creditCopy.Person = Person{}
	creditCopy.MediaItem = nil
	return creditCopy
}

// GetCreditWithDetails returns a credit with minimal person details
// Used for API responses to avoid circular references
func (c *Credit) GetCreditWithDetails() map[string]any {
	result := map[string]any{
		"id":          c.ID,
		"personID":    c.PersonID,
		"mediaItemID": c.MediaItemID,
		"name":        c.Name,
		"role":        c.Role,
		"department":  c.Department,
	}

	if c.IsCast {
		result["character"] = c.Character
		result["order"] = c.Order
	}

	if c.Job != "" {
		result["job"] = c.Job
	}

	// Add basic person details
	if c.Person.ID > 0 {
		result["person"] = map[string]any{
			"id":    c.Person.ID,
			"name":  c.Person.Name,
			"photo": c.Person.Photo,
		}
	}

	return result
}
