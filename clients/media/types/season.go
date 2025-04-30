package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Season represents a TV season
type Season struct {
	MediaData `json:"-"`
	Details   *MediaDetails `json:"details"`

	SeriesID     uint64   `json:"seriesID"`
	EpisodeIDs   []uint64 `json:"episodeIDs,omitempty"`
	EpisodeCount int      `json:"episodeCount"`

	Number      int       `json:"seasonNumber"`
	Title       string    `json:"title,omitempty"`
	Overview    string    `json:"overview,omitempty"`
	Artwork     Artwork   `json:"artwork,omitempty"`
	ReleaseDate time.Time `json:"releaseDate,omitempty"`
	SeriesName  string    `json:"seriesName,omitempty"`

	Credits Credits `json:"credits,omitempty"`
}

func (*Season) isMediaData() {}

// Scan
func (m *Season) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *Season) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *Season) Merge(other MediaData) {
	otherSeason, ok := other.(*Season)
	if !ok {
		return
	}
	m.Details.Merge(otherSeason.Details)
	m.MergeEpisodeIDs(otherSeason.EpisodeIDs)
	m.EpisodeCount = len(m.EpisodeIDs)
	
	// Preserve the season number when merging if the other season has a non-zero number
	if otherSeason.Number > 0 {
		m.Number = otherSeason.Number
	}
}

func (m *Season) MergeEpisodeIDs(otherSeasons []uint64) {
	for _, episodeID := range otherSeasons {
		found := false
		for i, existingEpisodeID := range m.EpisodeIDs {
			if existingEpisodeID == episodeID {
				// Update existing entry
				m.EpisodeIDs[i] = episodeID
				found = true
				break
			}
		}
		if !found {
			// Add new entry
			m.EpisodeIDs = append(m.EpisodeIDs, episodeID)
		}
	}
	// update episode count
	m.EpisodeCount = len(m.EpisodeIDs)
}

func (m *Season) GetSeriesID() uint64         { return m.SeriesID }
func (m *Season) SetSeriesID(seriesID uint64) { m.SeriesID = seriesID }

func (m *Season) GetEpisodeIDs() []uint64 { return m.EpisodeIDs }

func (s *Season) GetDetails() *MediaDetails { return s.Details }
func (s *Season) GetMediaType() MediaType   { return MediaTypeSeason }
