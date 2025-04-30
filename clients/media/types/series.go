package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sort"
)

// SeasonEntry represents a single season with its ID and episode IDs
type SeasonEntry struct {
	SeasonNumber int      `json:"seasonNumber"`
	SeasonID     uint64   `json:"seasonID"`
	EpisodeIDs   []uint64 `json:"episodeIDs,omitempty"`
}

// SeasonEntries is a collection of SeasonEntry objects
type SeasonEntries []SeasonEntry

// Series represents a TV series
type Series struct {
	MediaData `json:"-"`
	Details   *MediaDetails `json:"details"`
	// Collection of seasons with their IDs and episode IDs
	Seasons       SeasonEntries `json:"seasons,omitempty"`
	EpisodeCount  int           `json:"episodeCount"`
	SeasonCount   int           `json:"seasonCount"`
	ReleaseYear   int           `json:"releaseYear"`
	ContentRating string        `json:"contentRating"`
	Rating        float64       `json:"rating"`
	Network       string        `json:"network,omitempty"`
	Status        string        `json:"status,omitempty"` // e.g., "Ended", "Continuing"
	Genres        []string      `json:"genres,omitempty"`
	Credits       Credits       `json:"credits,omitempty"`
}

func (m *Series) SetDetails(details *MediaDetails) {
	m.Details = details
}

// GetOrderedSeasons returns the season numbers in ascending order
func (s *Series) GetOrderedSeasons() []int {
	if len(s.Seasons) == 0 {
		return []int{}
	}

	seasons := make([]int, 0, len(s.Seasons))
	for _, season := range s.Seasons {
		seasons = append(seasons, season.SeasonNumber)
	}

	sort.Ints(seasons)
	return seasons
}

func (*Series) isMediaData() {}

func (t *Series) GetTitle() string { return t.Details.Title }

func (t *Series) GetDetails() *MediaDetails { return t.Details }
func (t *Series) GetMediaType() MediaType   { return MediaTypeSeries }

// AddSeasonEpisodeIDs adds episode IDs to a season
func (s *Series) AddSeasonEpisodeIDs(season *Season) {
	if season == nil {
		return
	}

	// Try to find an existing season with the same number
	for i, existingSeason := range s.Seasons {
		if existingSeason.SeasonNumber == season.Number {
			// Found an existing season, merge episode IDs
			s.Seasons[i].EpisodeIDs = mergeEpisodeIDs(existingSeason.EpisodeIDs, season.EpisodeIDs)
			return
		}
	}

	// Season doesn't exist, create a new entry
	s.Seasons = append(s.Seasons, SeasonEntry{
		SeasonNumber: season.Number,
		EpisodeIDs:   season.EpisodeIDs,
	})

	// Update counts
	s.updateCounts()
}

// AddSeasonID is a legacy method that adds a season ID to any season without an ID
// This is kept for backward compatibility
func (s *Series) AddSeasonID(seasonID uint64) {
	if seasonID == 0 {
		return
	}

	// Try to find a season without an ID (to assign this ID to)
	for i := range s.Seasons {
		if s.Seasons[i].SeasonID == 0 {
			s.Seasons[i].SeasonID = seasonID
			return
		}
	}

	// If all seasons have IDs, check if this ID is a duplicate
	for _, season := range s.Seasons {
		if season.SeasonID == seasonID {
			return // ID already exists
		}
	}

	// If we get here, we have an ID without a matching season
	// This shouldn't normally happen, but add a dummy season just in case
	s.Seasons = append(s.Seasons, SeasonEntry{
		SeasonNumber: len(s.Seasons) + 1, // Fallback season number
		SeasonID:     seasonID,
	})
	s.updateCounts()
}

// SetSeasonID sets a specific season ID for a season with the given number
func (s *Series) SetSeasonID(seasonNumber int, seasonID uint64) {
	if seasonID == 0 {
		return
	}

	// Find the season by its number
	for i := range s.Seasons {
		if s.Seasons[i].SeasonNumber == seasonNumber {
			s.Seasons[i].SeasonID = seasonID
			return
		}
	}

	// If we didn't find a matching season, create a new one
	s.Seasons = append(s.Seasons, SeasonEntry{
		SeasonNumber: seasonNumber,
		SeasonID:     seasonID,
	})
	s.updateCounts()
}

// GetSeasonIDs returns all season IDs
func (s *Series) GetSeasonIDs() []uint64 {
	seasonIDs := make([]uint64, 0, len(s.Seasons))
	for _, season := range s.Seasons {
		if season.SeasonID != 0 {
			seasonIDs = append(seasonIDs, season.SeasonID)
		}
	}
	return seasonIDs
}

// GetSeasonByNumber returns a specific season by number
func (s *Series) GetSeasonByNumber(number int) *SeasonEntry {
	for i, season := range s.Seasons {
		if season.SeasonNumber == number {
			return &s.Seasons[i]
		}
	}
	return nil
}

func (s *Season) GetTitle() string { return s.Details.Title }

// Scan implements the Scanner interface for SQL
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

// Value implements the Valuer interface for SQL
func (m *Series) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// MergeSeasons merges seasons from another series
func (m *Series) MergeSeasons(otherSeries *Series) {
	if otherSeries == nil || len(otherSeries.Seasons) == 0 {
		return
	}

	// Process each season from the other series
	for _, otherSeason := range otherSeries.Seasons {
		existingSeason := m.GetSeasonByNumber(otherSeason.SeasonNumber)

		if existingSeason != nil {
			// Merge episode IDs
			existingSeason.EpisodeIDs = mergeEpisodeIDs(existingSeason.EpisodeIDs, otherSeason.EpisodeIDs)

			// Take the other season ID if this one doesn't have one
			if existingSeason.SeasonID == 0 && otherSeason.SeasonID != 0 {
				existingSeason.SeasonID = otherSeason.SeasonID
			}
		} else {
			// Add new season
			m.Seasons = append(m.Seasons, otherSeason)
		}
	}

	// Update counts
	m.updateCounts()
}

// updateCounts recalculates episode and season counts
func (m *Series) updateCounts() {
	m.SeasonCount = len(m.Seasons)
	m.EpisodeCount = m.CalculateEpisodeCount()
}

// CalculateEpisodeCount returns the total number of episodes
func (m *Series) CalculateEpisodeCount() int {
	episodeCount := 0
	for _, season := range m.Seasons {
		episodeCount += len(season.EpisodeIDs)
	}
	return episodeCount
}

// MergeEpisodeIDsBySeason merges episode IDs for a specific season
func (m *Series) MergeEpisodeIDsBySeason(seasonNumber int, episodeIDs []uint64) {
	season := m.GetSeasonByNumber(seasonNumber)

	if season != nil {
		// Merge into existing season
		season.EpisodeIDs = mergeEpisodeIDs(season.EpisodeIDs, episodeIDs)
	} else {
		// Create new season
		m.Seasons = append(m.Seasons, SeasonEntry{
			SeasonNumber: seasonNumber,
			EpisodeIDs:   episodeIDs,
		})
	}

	// Update counts
	m.updateCounts()
}

// Helper function to merge episode IDs without duplicates
func mergeEpisodeIDs(existing, new []uint64) []uint64 {
	// Create a map for fast lookup
	idMap := make(map[uint64]bool)
	for _, id := range existing {
		idMap[id] = true
	}

	// Add new IDs that don't already exist
	for _, id := range new {
		if !idMap[id] {
			existing = append(existing, id)
			idMap[id] = true
		}
	}

	return existing
}

// GetEpisodeIDsBySeason returns episode IDs for a specific season
func (m *Series) GetEpisodeIDsBySeason(seasonNumber int) []uint64 {
	season := m.GetSeasonByNumber(seasonNumber)
	if season == nil {
		return []uint64{}
	}
	return season.EpisodeIDs
}

// GetAllEpisodeIDs returns all episode IDs across all seasons
func (m *Series) GetAllEpisodeIDs() []uint64 {
	var allEpisodeIDs []uint64
	for _, season := range m.Seasons {
		allEpisodeIDs = append(allEpisodeIDs, season.EpisodeIDs...)
	}
	return allEpisodeIDs
}

// Merge merges this series with another one
func (m *Series) Merge(other MediaData) {
	otherSeries, ok := other.(*Series)
	if !ok {
		return
	}

	// Initialize Details if nil
	if m.Details == nil && otherSeries.Details != nil {
		m.Details = &MediaDetails{}
	}

	if m.Details != nil && otherSeries.Details != nil {
		m.Details.Merge(otherSeries.Details)
	}

	m.MergeSeasons(otherSeries)
}
