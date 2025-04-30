package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Track struct {
	Details    *MediaDetails
	AlbumID    uint64  `json:"albumID"`
	ArtistID   uint64  `json:"artistID"`
	AlbumName  string  `json:"albumName"`
	AlbumTitle string  `json:"albumTitle,omitempty"`
	Duration   int     `json:"duration,omitempty"`
	ArtistName string  `json:"artistName,omitempty"`
	Number     int     `json:"trackNumber,omitempty"`
	DiscNumber int     `json:"discNumber,omitempty"`
	Composer   string  `json:"composer,omitempty"`
	Lyrics     string  `json:"lyrics,omitempty"`
	Credits    Credits `json:"credits,omitempty"`
}

func (Track) isMediaData()                {}
func (t Track) GetDetails() *MediaDetails { return t.Details }
func (t Track) GetMediaType() MediaType   { return MediaTypeTrack }
func (t Track) GetTitle() string          { return t.Details.Title }

func (t *Track) SetDetails(details *MediaDetails) {
	t.Details = details
}

// Scan
func (m *Track) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Track) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Track) Merge(other MediaData) {
	otherTrack, ok := other.(*Track)
	if !ok {
		return
	}
	m.Details.Merge(otherTrack.Details)
	m.Credits.Merge(&otherTrack.Credits)
	if m.AlbumID == 0 {
		m.AlbumID = otherTrack.AlbumID
	}
	if m.ArtistID == 0 {
		m.ArtistID = otherTrack.ArtistID
	}
	if m.AlbumName == "" {
		m.AlbumName = otherTrack.AlbumName
	}
	if m.AlbumTitle == "" {
		m.AlbumTitle = otherTrack.AlbumTitle
	}
	if m.Duration == 0 {
		m.Duration = otherTrack.Duration
	}
	if m.ArtistName == "" {
		m.ArtistName = otherTrack.ArtistName
	}
	if m.Number == 0 {
		m.Number = otherTrack.Number
	}
	if m.DiscNumber == 0 {
		m.DiscNumber = otherTrack.DiscNumber
	}
	if m.Composer == "" {
		m.Composer = otherTrack.Composer
	}
	if m.Lyrics == "" {
		m.Lyrics = otherTrack.Lyrics
	}

}
