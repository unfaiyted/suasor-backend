package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Person struct {
	MediaData `json:"-"`
	Name      string `json:"name"`
	Role      string `json:"role,omitempty"`      // e.g., "Director", "Actor"
	Character string `json:"character,omitempty"` // For actors
	Photo     string `json:"photo,omitempty"`

	IsCast    bool `json:"isCast,omitempty"`
	IsCrew    bool `json:"isCrew,omitempty"`
	IsGuest   bool `json:"isGuest,omitempty"`
	IsCreator bool `json:"isCreator,omitempty"`
	IsArtist  bool `json:"isArtist,omitempty"`
}

// Scan
func (m *Person) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Person) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Person) Merge(other *Person) {
	if m.Name == "" {
		m.Name = other.Name
	}
	if m.Role == "" {
		m.Role = other.Role
	}
	if m.Character == "" {
		m.Character = other.Character
	}
	if m.Photo == "" {
		m.Photo = other.Photo
	}
	if m.IsCast == false {
		m.IsCast = other.IsCast
	}
	if m.IsCrew == false {
		m.IsCrew = other.IsCrew
	}
	if m.IsGuest == false {
		m.IsGuest = other.IsGuest
	}
	if m.IsCreator == false {
		m.IsCreator = other.IsCreator
	}
	if m.IsArtist == false {
		m.IsArtist = other.IsArtist
	}
}
