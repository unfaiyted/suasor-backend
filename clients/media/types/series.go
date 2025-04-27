package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Season represents a TV season
type Season struct {
	MediaData    `json:"-"`
	Details      MediaDetails `json:"details"`
	Number       int          `json:"seasonNumber"`
	Title        string       `json:"title,omitempty"`
	Overview     string       `json:"overview,omitempty"`
	EpisodeCount int          `json:"episodeCount"`
	Episodes     []*Episode   `json:"episodes,omitempty"`
	EpisodeIDs   []uint64     `json:"episodeIDs,omitempty"`
	Artwork      Artwork      `json:"artwork,omitempty"`
	ReleaseDate  time.Time    `json:"releaseDate,omitempty"`
	SeriesName   string       `json:"seriesName,omitempty"`
	SeriesID     uint64       `json:"seriesID"`
	SyncSeries   SyncClients  `json:"syncSeries,omitempty"`
	Credits      Credits      `json:"credits,omitempty"`
}

// Series represents a TV series
type Series struct {
	MediaData     `json:"-"`
	Details       MediaDetails `json:"details"`
	Seasons       []*Season    `json:"seasons,omitempty"`
	EpisodeCount  int          `json:"episodeCount"`
	SeasonCount   int          `json:"seasonCount"`
	ReleaseYear   int          `json:"releaseYear"`
	ContentRating string       `json:"contentRating"`
	Rating        float64      `json:"rating"`
	Network       string       `json:"network,omitempty"`
	Status        string       `json:"status,omitempty"` // e.g., "Ended", "Continuing"
	Genres        []string     `json:"genres,omitempty"`
	Credits       Credits      `json:"credits,omitempty"`
}

func (m *Series) SetDetails(details MediaDetails) {
	m.Details = details
}

func (*Series) isMediaData() {}

func (t *Series) GetTitle() string { return t.Details.Title }

func (t *Series) GetDetails() MediaDetails { return t.Details }
func (t *Series) GetMediaType() MediaType  { return MediaTypeSeries }

func (s *Season) GetDetails() MediaDetails { return s.Details }
func (s *Season) GetMediaType() MediaType  { return MediaTypeSeason }

func (s *Season) GetTitle() string { return s.Details.Title }

// Scan
func (m *Series) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Series) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
