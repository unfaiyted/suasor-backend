package requests

import (
	"time"

	"suasor/types/models"
)

// CreatePersonRequest represents the data for creating a new person
//
//	@Description	Request payload for creating a new person
type CreatePersonRequest struct {
	// Name is the person's name
	//	@Description	Person's name
	//	@Example		"Tom Hanks"
	Name string `json:"name" binding:"required" example:"Tom Hanks"`

	// Photo is the URL or path to the person's photo
	//	@Description	URL or path to the person's photo
	//	@Example		"https://example.com/photos/tom-hanks.jpg"
	Photo string `json:"photo,omitempty" example:"https://example.com/photos/tom-hanks.jpg"`

	// DateOfBirth is the person's date of birth
	//	@Description	Person's date of birth in RFC3339 format
	//	@Example		"1956-07-09T00:00:00Z"
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty" example:"1956-07-09T00:00:00Z"`

	// DateOfDeath is the person's date of death (if applicable)
	//	@Description	Person's date of death in RFC3339 format (if applicable)
	//	@Example		"2056-07-09T00:00:00Z"
	DateOfDeath *time.Time `json:"dateOfDeath,omitempty" example:"2056-07-09T00:00:00Z"`

	// Gender is the person's gender
	//	@Description	Person's gender
	//	@Example		"Male"
	Gender string `json:"gender,omitempty" example:"Male"`

	// Biography is the person's biography
	//	@Description	Person's biography
	//	@Example		"Thomas Jeffrey Hanks is an American actor and filmmaker..."
	Biography string `json:"biography,omitempty" example:"Thomas Jeffrey Hanks is an American actor and filmmaker..."`

	// Birthplace is the person's birthplace
	//	@Description	Person's birthplace
	//	@Example		"Concord, California, USA"
	Birthplace string `json:"birthplace,omitempty" example:"Concord, California, USA"`

	// KnownFor is what the person is primarily known for
	//	@Description	What the person is primarily known for
	//	@Example		"Actor"
	KnownFor string `json:"knownFor,omitempty" example:"Actor"`

	// ExternalIDs contains IDs from external services
	//	@Description	IDs from external services
	ExternalIDs []ExternalIDRequest `json:"externalIds,omitempty"`
}

// UpdatePersonRequest represents the data for updating an existing person
//
//	@Description	Request payload for updating an existing person
type UpdatePersonRequest struct {
	// Name is the person's name
	//	@Description	Person's name
	//	@Example		"Tom Hanks"
	Name string `json:"name,omitempty" example:"Tom Hanks"`

	// Photo is the URL or path to the person's photo
	//	@Description	URL or path to the person's photo
	//	@Example		"https://example.com/photos/tom-hanks.jpg"
	Photo string `json:"photo,omitempty" example:"https://example.com/photos/tom-hanks.jpg"`

	// DateOfBirth is the person's date of birth
	//	@Description	Person's date of birth in RFC3339 format
	//	@Example		"1956-07-09T00:00:00Z"
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty" example:"1956-07-09T00:00:00Z"`

	// DateOfDeath is the person's date of death (if applicable)
	//	@Description	Person's date of death in RFC3339 format (if applicable)
	//	@Example		"2056-07-09T00:00:00Z"
	DateOfDeath *time.Time `json:"dateOfDeath,omitempty" example:"2056-07-09T00:00:00Z"`

	// Gender is the person's gender
	//	@Description	Person's gender
	//	@Example		"Male"
	Gender string `json:"gender,omitempty" example:"Male"`

	// Biography is the person's biography
	//	@Description	Person's biography
	//	@Example		"Thomas Jeffrey Hanks is an American actor and filmmaker..."
	Biography string `json:"biography,omitempty" example:"Thomas Jeffrey Hanks is an American actor and filmmaker..."`

	// Birthplace is the person's birthplace
	//	@Description	Person's birthplace
	//	@Example		"Concord, California, USA"
	Birthplace string `json:"birthplace,omitempty" example:"Concord, California, USA"`

	// KnownFor is what the person is primarily known for
	//	@Description	What the person is primarily known for
	//	@Example		"Actor"
	KnownFor string `json:"knownFor,omitempty" example:"Actor"`

	// ExternalIDs contains IDs from external services
	//	@Description	IDs from external services
	ExternalIDs []ExternalIDRequest `json:"externalIds,omitempty"`
}

// ExternalIDRequest represents an external ID for a person
//
//	@Description	External ID for a person
type ExternalIDRequest struct {
	// Source is the name of the external service
	//	@Description	Name of the external service
	//	@Example		"TMDB"
	Source string `json:"source" binding:"required" example:"TMDB"`

	// ID is the identifier in the external service
	//	@Description	Identifier in the external service
	//	@Example		"31"
	ID string `json:"id" binding:"required" example:"31"`
}

// SearchPersonRequest represents the parameters for searching people
//
//	@Description	Parameters for searching people
type SearchPersonRequest struct {
	// Query is the search term
	//	@Description	Search term
	//	@Example		"Tom"
	Query string `json:"query" binding:"required" example:"Tom"`

	// Limit is the maximum number of results to return
	//	@Description	Maximum number of results to return
	//	@Example		10
	Limit int `json:"limit,omitempty" example:"10"`
}

// ImportPersonRequest represents the data for importing a person from an external source
//
//	@Description	Request payload for importing a person from an external source
type ImportPersonRequest struct {
	// Source is the name of the external service
	//	@Description	Name of the external service
	//	@Example		"TMDB"
	Source string `json:"source" binding:"required" example:"TMDB"`

	// ExternalID is the identifier in the external service
	//	@Description	Identifier in the external service
	//	@Example		"31"
	ExternalID string `json:"externalId" binding:"required" example:"31"`

	// PersonData contains the person data to import
	//	@Description	Person data to import
	PersonData CreatePersonRequest `json:"personData" binding:"required"`
}

// CreateCreditRequest represents the data for creating a new credit
//
//	@Description	Request payload for creating a new credit
type CreateCreditRequest struct {
	// PersonID is the ID of the person
	//	@Description	ID of the person
	//	@Example		1
	PersonID uint64 `json:"personId" binding:"required" example:"1"`

	// MediaItemID is the ID of the media item
	//	@Description	ID of the media item
	//	@Example		2
	MediaItemID uint64 `json:"mediaItemId" binding:"required" example:"2"`

	// Name is the person's name for this credit
	//	@Description	Person's name for this credit
	//	@Example		"Tom Hanks"
	Name string `json:"name" binding:"required" example:"Tom Hanks"`

	// Role is the person's role
	//	@Description	Person's role
	//	@Example		"Actor"
	Role models.MediaRole `json:"role,omitempty" example:"Actor"`

	// Character is the character's name (for acting roles)
	//	@Description	Character's name (for acting roles)
	//	@Example		"Forrest Gump"
	Character string `json:"character,omitempty" example:"Forrest Gump"`

	// Department is the department the person worked in
	//	@Description	Department the person worked in
	//	@Example		"Acting"
	Department models.MediaDepartment `json:"department,omitempty" example:"Acting"`

	// Job is the specific job the person had
	//	@Description	Specific job the person had
	//	@Example		"Lead Actor"
	Job string `json:"job,omitempty" example:"Lead Actor"`

	// Order is the order of importance (lower means more important)
	//	@Description	Order of importance (lower means more important)
	//	@Example		1
	Order int `json:"order,omitempty" example:"1"`

	// IsCast indicates if this is a cast credit
	//	@Description	Indicates if this is a cast credit
	//	@Example		true
	IsCast bool `json:"isCast,omitempty" example:"true"`

	// IsCrew indicates if this is a crew credit
	//	@Description	Indicates if this is a crew credit
	//	@Example		false
	IsCrew bool `json:"isCrew,omitempty" example:"false"`

	// IsGuest indicates if this is a guest credit
	//	@Description	Indicates if this is a guest credit
	//	@Example		false
	IsGuest bool `json:"isGuest,omitempty" example:"false"`

	// IsCreator indicates if this is a creator credit
	//	@Description	Indicates if this is a creator credit
	//	@Example		false
	IsCreator bool `json:"isCreator,omitempty" example:"false"`

	// IsArtist indicates if this is an artist credit
	//	@Description	Indicates if this is an artist credit
	//	@Example		false
	IsArtist bool `json:"isArtist,omitempty" example:"false"`
}

// UpdateCreditRequest represents the data for updating an existing credit
//
//	@Description	Request payload for updating an existing credit
type UpdateCreditRequest struct {
	// PersonID is the ID of the person
	//	@Description	ID of the person
	//	@Example		1
	PersonID uint64 `json:"personId,omitempty" example:"1"`

	// MediaItemID is the ID of the media item
	//	@Description	ID of the media item
	//	@Example		2
	MediaItemID uint64 `json:"mediaItemId,omitempty" example:"2"`

	// Name is the person's name for this credit
	//	@Description	Person's name for this credit
	//	@Example		"Tom Hanks"
	Name string `json:"name,omitempty" example:"Tom Hanks"`

	// Role is the person's role
	//	@Description	Person's role
	//	@Example		"Actor"
	Role models.MediaRole `json:"role,omitempty" example:"Actor"`

	// Character is the character's name (for acting roles)
	//	@Description	Character's name (for acting roles)
	//	@Example		"Forrest Gump"
	Character string `json:"character,omitempty" example:"Forrest Gump"`

	// Department is the department the person worked in
	//	@Description	Department the person worked in
	//	@Example		"Acting"
	Department models.MediaDepartment `json:"department,omitempty" example:"Acting"`

	// Job is the specific job the person had
	//	@Description	Specific job the person had
	//	@Example		"Lead Actor"
	Job string `json:"job,omitempty" example:"Lead Actor"`

	// Order is the order of importance (lower means more important)
	//	@Description	Order of importance (lower means more important)
	//	@Example		1
	Order int `json:"order,omitempty" example:"1"`

	// IsCast indicates if this is a cast credit
	//	@Description	Indicates if this is a cast credit
	//	@Example		true
	IsCast bool `json:"isCast,omitempty" example:"true"`

	// IsCrew indicates if this is a crew credit
	//	@Description	Indicates if this is a crew credit
	//	@Example		false
	IsCrew bool `json:"isCrew,omitempty" example:"false"`

	// IsGuest indicates if this is a guest credit
	//	@Description	Indicates if this is a guest credit
	//	@Example		false
	IsGuest bool `json:"isGuest,omitempty" example:"false"`

	// IsCreator indicates if this is a creator credit
	//	@Description	Indicates if this is a creator credit
	//	@Example		false
	IsCreator bool `json:"isCreator,omitempty" example:"false"`

	// IsArtist indicates if this is an artist credit
	//	@Description	Indicates if this is an artist credit
	//	@Example		false
	IsArtist bool `json:"isArtist,omitempty" example:"false"`
}

// CreateCreditsRequest represents the data for creating multiple credits for a media item
//
//	@Description	Request payload for creating multiple credits for a media item
type CreateCreditsRequest struct {
	// Credits is the list of credits to create
	//	@Description	List of credits to create
	Credits []CreateCreditRequest `json:"credits" binding:"required"`
}

