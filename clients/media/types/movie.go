package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Movie represents a movie item
type Movie struct {
	MediaData    `json:"-"`
	Details      *MediaDetails `json:"details"`
	Credits      *Credits      `json:"credits,omitempty"`
	TrailerURL   string        `json:"trailerUrl,omitempty"`
	Resolution   string        `json:"resolution,omitempty"` // e.g., "4K", "1080p"
	VideoCodec   string        `json:"videoCodec,omitempty"`
	AudioCodec   string        `json:"audioCodec,omitempty"`
	SubtitleURLs []string      `json:"subtitleUrls,omitempty"`
}

func (m *Movie) SetDetails(details *MediaDetails) {
	m.Details = details
}

func (m *Movie) GetDetails() *MediaDetails {
	return m.Details
}
func (m *Movie) GetMediaType() MediaType {
	return MediaTypeMovie
}
func (m *Movie) isMediaData() {}

func (m *Movie) GetTitle() string { return m.Details.Title }

// Scan
func (m *Movie) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Movie) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Movie) Merge(other MediaData) {
	otherMovie, ok := other.(*Movie)
	if !ok {
		return
	}
	m.Details.Merge(otherMovie.Details)
	m.Credits.Merge(otherMovie.Credits)
	if m.TrailerURL == "" {
		m.TrailerURL = otherMovie.TrailerURL
	}
	if m.Resolution == "" {
		m.Resolution = otherMovie.Resolution
	}
	if m.VideoCodec == "" {
		m.VideoCodec = otherMovie.VideoCodec
	}
	if m.AudioCodec == "" {
		m.AudioCodec = otherMovie.AudioCodec
	}
	if m.SubtitleURLs == nil {
		m.SubtitleURLs = otherMovie.SubtitleURLs
	}
}
