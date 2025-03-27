package types

import (
	"errors"
	"time"
)

var ErrAutomationFeatureNotSupported = errors.New("feature not supported by this automation tool")

// SystemStatus represents system information from the automation tool
type SystemStatus struct {
	Version     string
	StartupPath string
	AppData     string
	OsName      string
	IsUpdated   bool
	HasUpdate   bool
	Branch      string
}

type ExternalID struct {
	Source string
	Value  string
}

// RatingType the model 'RatingType'
type RatingType string

// List of RatingType
const (
	RATINGTYPE_USER   RatingType = "user"
	RATINGTYPE_CRITIC RatingType = "critic"
)

type Rating struct {
	Votes *int32      `json:"votes,omitempty"`
	Value *float64    `json:"value,omitempty"`
	Type  *RatingType `json:"type,omitempty"`
}

// MediaImage represents an image associated wit media
type AutomationMediaImage struct {
	URL       string
	CoverType string // "poster", "fanart", "banner", etc.
}

// SearchOptions defines options for search requests
type SearchOptions struct {
	Limit  int
	Offset int
}

// LibraryQueryOptions defines options for library queries
type LibraryQueryOptions struct {
	Limit  int
	Offset int
	SortBy string
	Filter string
}

// QualityProfile represents a quality profile configuration
type QualityProfile struct {
	ID    int64
	Name  string
	Items []QualityItem
}

// QualityProfileSummary is a simplified quality profile reference
type QualityProfileSummary struct {
	ID   int64
	Name string
}

// MetadataProfile is a simplified quality profile reference
type MetadataProfile struct {
	ID   int32
	Name string
}

// QualityItem represents an item in a quality profile
type QualityItem struct {
	ID         int64
	Name       string
	Resolution string
	Source     string
}

// Tag represents a media tag
type Tag struct {
	ID   int64
	Name string
}

// Command represents a command to execute
type Command struct {
	Name       string
	MediaIDs   []int64
	Parameters map[string]any
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	ID        int64
	Name      string
	Status    string
	StartedAt time.Time
	Completed time.Time
	Message   string
}
