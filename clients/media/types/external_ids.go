package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

type ExternalID struct {
	Source string `json:"source"` // e.g., "tmdb", "imdb", "trakt", "tvdb"
	ID     string `json:"id"`     // The actual ID
}

type ExternalIDs []ExternalID

func (ids ExternalIDs) GetID(source string) string {
	for _, id := range ids {
		if id.Source == source {
			return id.ID
		}
	}
	return ""
}

func (ids ExternalIDs) Merge(other ExternalIDs) {
	for _, otherID := range other {
		found := false
		for i, id := range ids {
			if id.Source == otherID.Source {
				// Update existing entry
				ids[i].ID = otherID.ID
				found = true
				break
			}
		}
		if !found {
			// Add new entry
			ids = append(ids, otherID)
		}
	}
}

// Value implements driver.Valuer for database storage
func (ids ExternalIDs) Value() (driver.Value, error) {
	if ids == nil {
		return nil, nil
	}
	return json.Marshal(ids)
}

// Scan implements sql.Scanner for database retrieval
func (ids *ExternalIDs) Scan(value any) error {
	if value == nil {
		*ids = make(ExternalIDs, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ids)
}

func (ids ExternalIDs) String() string {
	if len(ids) == 0 {
		return ""
	}
	var strs []string
	for _, id := range ids {
		strs = append(strs, id.Source+":"+id.ID)
	}
	return "[" + strings.Join(strs, ", ") + "]"
}
