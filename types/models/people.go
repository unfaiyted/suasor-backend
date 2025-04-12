package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	client "suasor/client/types"
	"time"
)

// Person represents someone involved with media (actors, directors, etc.)
type Person struct {
	BaseModel               // Include base fields (ID, timestamps)
	Name        string      `json:"name" gorm:"type:varchar(255);index"`
	ClientIDs   ClientIDs   `json:"clientIds" gorm:"type:jsonb"`
	ExternalIDs ExternalIDs `json:"externalIds" gorm:"type:jsonb"`

	// Biographical information
	Photo       string     `json:"photo,omitempty" gorm:"type:text"`
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`
	DateOfDeath *time.Time `json:"dateOfDeath,omitempty"`
	Gender      string     `json:"gender,omitempty" gorm:"type:varchar(50)"`
	Biography   string     `json:"biography,omitempty" gorm:"type:text"`
	Birthplace  string     `json:"birthplace,omitempty" gorm:"type:varchar(255)"`
	Popularity  float32    `json:"popularity,omitempty"`

	// Professional information
	KnownFor string `json:"knownFor,omitempty" gorm:"type:varchar(255)"` // Primary role (Actor, Director, etc.)

	// Relationships - handled via joins rather than embedding
	Credits []Credit `json:"-" gorm:"foreignKey:PersonID;constraint:OnDelete:CASCADE"` // Use json:"-" to avoid circular references

	// Additional metadata
	Metadata PersonMetadata `json:"metadata,omitempty" gorm:"type:jsonb"`
}

// PersonMetadata holds additional, less structured data about a person
type PersonMetadata struct {
	SocialMedia    SocialMedia    `json:"socialMedia,omitempty"`
	AwardHistory   []Award        `json:"awards,omitempty"`
	ExternalLinks  []ExternalLink `json:"externalLinks,omitempty"`
	AlternateNames []string       `json:"alternateNames,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
}

// Value implements the driver.Valuer interface for database serialization
func (pm PersonMetadata) Value() (driver.Value, error) {
	return json.Marshal(pm)
}

// Scan implements the sql.Scanner interface for database deserialization
func (pm *PersonMetadata) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &pm)
}

// SocialMedia represents a person's social media accounts
type SocialMedia struct {
	Twitter   string `json:"twitter,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	Website   string `json:"website,omitempty"`
}

// Award represents an award a person has won or been nominated for
type Award struct {
	Name       string `json:"name"`
	Year       int    `json:"year"`
	Category   string `json:"category,omitempty"`
	IsWinner   bool   `json:"isWinner"`
	Production string `json:"production,omitempty"` // Movie/show the award was for
}

// ExternalLink represents a link to an external site about the person
type ExternalLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Helper methods for accessing related records
func (p *Person) GetCredits() []Credit {
	return p.Credits
}

func (p *Person) GetCastCredits() []Credit {
	var castCredits []Credit
	for _, credit := range p.Credits {
		if credit.IsCast {
			castCredits = append(castCredits, credit)
		}
	}
	return castCredits
}

func (p *Person) GetCrewCredits() []Credit {
	var crewCredits []Credit
	for _, credit := range p.Credits {
		if credit.IsCrew {
			crewCredits = append(crewCredits, credit)
		}
	}
	return crewCredits
}

func (p *Person) GetDirectingCredits() []Credit {
	var directingCredits []Credit
	for _, credit := range p.Credits {
		if credit.Role == "Director" {
			directingCredits = append(directingCredits, credit)
		}
	}
	return directingCredits
}

func (p *Person) GetWritingCredits() []Credit {
	var writingCredits []Credit
	for _, credit := range p.Credits {
		if credit.Role == "Writer" || credit.Role == "Screenplay" {
			writingCredits = append(writingCredits, credit)
		}
	}
	return writingCredits
}

// GetPersonPublicView returns a version of the person suitable for API responses
func (p *Person) GetPersonPublicView() map[string]any {
	return map[string]any{
		"id":          p.ID,
		"name":        p.Name,
		"photo":       p.Photo,
		"gender":      p.Gender,
		"knownFor":    p.KnownFor,
		"dateOfBirth": p.DateOfBirth,
		"biography":   p.Biography,
		"popularity":  p.Popularity,
	}
}

// NewPerson creates a new person record with the given name
func NewPerson(name string) *Person {
	return &Person{
		Name:        name,
		ClientIDs:   make(ClientIDs, 0),
		ExternalIDs: make(ExternalIDs, 0),
		Metadata: PersonMetadata{
			SocialMedia:    SocialMedia{},
			AwardHistory:   make([]Award, 0),
			ExternalLinks:  make([]ExternalLink, 0),
			AlternateNames: make([]string, 0),
			Tags:           make([]string, 0),
		},
	}
}

// AddExternalID adds an external ID to the person's record
func (p *Person) AddExternalID(source string, id string) {
	for i, extID := range p.ExternalIDs {
		if extID.Source == source {
			// Update existing ID
			p.ExternalIDs[i].ID = id
			return
		}
	}
	// Add new ID
	p.ExternalIDs = append(p.ExternalIDs, ExternalID{
		Source: source,
		ID:     id,
	})
}

// GetExternalID retrieves an external ID by source
func (p *Person) GetExternalID(source string) string {
	return p.ExternalIDs.GetID(source)
}

// AddClientID adds a client ID to the person's record
func (p *Person) AddClientID(clientID uint64, clientType client.ClientType, itemID string) {
	for i, cID := range p.ClientIDs {
		if cID.ID == clientID {
			// Update existing ID
			p.ClientIDs[i].ItemID = itemID
			return
		}
	}
	// Add new ID
	p.ClientIDs = append(p.ClientIDs, ClientID{
		ID:     clientID,
		Type:   clientType,
		ItemID: itemID,
	})
}

// AddAward adds an award to the person's record
func (p *Person) AddAward(name string, year int, category string, isWinner bool) {
	if p.Metadata.AwardHistory == nil {
		p.Metadata.AwardHistory = make([]Award, 0)
	}

	p.Metadata.AwardHistory = append(p.Metadata.AwardHistory, Award{
		Name:     name,
		Year:     year,
		Category: category,
		IsWinner: isWinner,
	})
}

// AddExternalLink adds an external link to the person's record
func (p *Person) AddExternalLink(name, url string) {
	if p.Metadata.ExternalLinks == nil {
		p.Metadata.ExternalLinks = make([]ExternalLink, 0)
	}

	p.Metadata.ExternalLinks = append(p.Metadata.ExternalLinks, ExternalLink{
		Name: name,
		URL:  url,
	})
}

// TableName specifies the database table name for Person
func (Person) TableName() string {
	return "people"
}
