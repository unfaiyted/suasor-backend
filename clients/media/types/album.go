package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// MusicAlbum represents a music album
type Album struct {
	MediaData `json:"-"`

	Details    *MediaDetails `json:"details"`
	ArtistID   uint64        `json:"artistID"`
	TrackIDs   []uint64      `json:"trackIDs,omitempty"`
	ArtistName string        `json:"artistName"`
	TrackCount int           `json:"trackCount"`
	Credits    Credits       `json:"credits,omitempty"`
	Tracks     []*Track      `json:"tracks,omitempty"`
}

func (a *Album) isMediaData() {}

func (a *Album) GetDetails() *MediaDetails { return a.Details }
func (a *Album) GetMediaType() MediaType   { return MediaTypeAlbum }

func (a *Album) GetTitle() string { return a.Details.Title }

func (a *Album) SetDetails(details *MediaDetails) {
	a.Details = details
}

// Scan
func (m *Album) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Album) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Album) Merge(other MediaData) {
	otherAlbum, ok := other.(*Album)
	if !ok {
		return
	}
	m.Details.Merge(otherAlbum.Details)
	m.Credits.Merge(&otherAlbum.Credits)
	if m.ArtistID == 0 {
		m.ArtistID = otherAlbum.ArtistID
	}
	if m.ArtistName == "" {
		m.ArtistName = otherAlbum.ArtistName
	}
	if m.TrackCount == 0 {
		m.TrackCount = otherAlbum.TrackCount
	}

}
