package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// MusicArtist represents a music artist
type Artist struct {
	MediaData `json:"-"`

	Details        MediaDetails      `json:"details"`
	Albums         []*Album          `json:"albums,omitempty"`
	AlbumIDs       []uint64          `json:"albumIDs,omitempty"`
	AlbumCount     int               `json:"albumCount"`
	Biography      string            `json:"biography,omitempty"`
	SimilarArtists []ArtistReference `json:"similarArtists,omitempty"`
}

type ArtistReference struct {
	Name string `json:"name"`
	ID   uint64 `json:"id"`
}

func (*Artist) isMediaData() {}

func (a *Artist) GetDetails() MediaDetails { return a.Details }
func (a *Artist) GetMediaType() MediaType  { return MediaTypeArtist }

func (a *Artist) GetTitle() string { return a.Details.Title }

// Scan
func (m *Artist) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Artist) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
