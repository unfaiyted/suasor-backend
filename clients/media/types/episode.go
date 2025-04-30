package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Episode represents a TV episode
type Episode struct {
	MediaData `json:"-"`
	Details   *MediaDetails

	SeriesID uint64 `json:"seriesID"` // itemID of the series
	SeasonID uint64 `json:"seasonID"` // itemID of the season

	Number       int64   `json:"number"`
	SeasonNumber int     `json:"seasonNumber"`
	ShowTitle    string  `json:"showTitle,omitempty"`
	Credits      Credits `json:"credits,omitempty"`
}

func (e *Episode) isMediaData() {}

func (e *Episode) GetDetails() *MediaDetails { return e.Details }
func (e *Episode) GetMediaType() MediaType   { return MediaTypeEpisode }

func (e *Episode) GetSeriesID() uint64         { return e.SeriesID }
func (e *Episode) SetSeriesID(seriesID uint64) { e.SeriesID = seriesID }

func (e *Episode) GetSeasonNumber() int             { return e.SeasonNumber }
func (e *Episode) SetSeasonNumber(seasonNumber int) { e.SeasonNumber = seasonNumber }

func (e *Episode) SetSeasonID(seasonID uint64) { e.SeasonID = seasonID }
func (e *Episode) GetSeasonID() uint64         { return e.SeasonID }

func (e *Episode) GetTitle() string { return e.Details.Title }

func (e *Episode) SetDetails(details *MediaDetails) {
	e.Details = details
}

// Scan
func (m *Episode) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Episode) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Episode) Merge(other MediaData) {
	otherEpisode, ok := other.(*Episode)
	if !ok {
		return
	}
	m.Details.Merge(otherEpisode.Details)
	m.SeasonNumber = otherEpisode.SeasonNumber
}
