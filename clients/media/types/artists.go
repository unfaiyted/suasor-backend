package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sort"
)

// AlbumEntry represents a single album with its ID and track IDs
type AlbumEntry struct {
	AlbumID   uint64   `json:"albumID"`
	AlbumName string   `json:"albumName"`
	TrackIDs  []uint64 `json:"trackIDs,omitempty"`
}

// AlbumEntries is a collection of AlbumEntry objects
type AlbumEntries []AlbumEntry

// Artist represents a music artist
type Artist struct {
	MediaData `json:"-"`
	Details   *MediaDetails `json:"details"`

	Albums         AlbumEntries `json:"albums,omitempty"`
	AlbumCount     int          `json:"albumCount"`
	TrackCount     int          `json:"trackCount"`
	Biography      string       `json:"biography,omitempty"`
	Genres         []string     `json:"genres,omitempty"`
	SimilarArtists []Person     `json:"similarArtists,omitempty"`
	StartYear      int          `json:"startYear,omitempty"`
	EndYear        int          `json:"endYear,omitempty"`
	Rating         float64      `json:"rating"`
	Credits        Credits      `json:"credits,omitempty"`
}

func (m *Artist) SetDetails(details *MediaDetails) {
	m.Details = details
}

// GetOrderedAlbums returns the album IDs in a sorted manner
func (a *Artist) GetOrderedAlbums() []uint64 {
	if len(a.Albums) == 0 {
		return []uint64{}
	}

	albums := make([]uint64, 0, len(a.Albums))
	for _, album := range a.Albums {
		albums = append(albums, album.AlbumID)
	}

	// Sort by album ID for consistency
	sort.Slice(albums, func(i, j int) bool {
		return albums[i] < albums[j]
	})

	return albums
}

func (*Artist) isMediaData() {}

func (a *Artist) GetTitle() string { return a.Details.Title }

func (a *Artist) GetDetails() *MediaDetails { return a.Details }
func (a *Artist) GetMediaType() MediaType   { return MediaTypeArtist }

// AddAlbumTrackIDs adds track IDs to an album
func (a *Artist) AddAlbumTrackIDs(albumID uint64, albumName string, trackIDs []uint64) {
	if albumID == 0 {
		return
	}

	// Try to find an existing album with the same ID
	for i, existingAlbum := range a.Albums {
		if existingAlbum.AlbumID == albumID {
			// Found an existing album, merge track IDs
			a.Albums[i].TrackIDs = mergeTrackIDs(existingAlbum.TrackIDs, trackIDs)
			return
		}
	}

	// Album doesn't exist, create a new entry
	a.Albums = append(a.Albums, AlbumEntry{
		AlbumID:   albumID,
		AlbumName: albumName,
		TrackIDs:  trackIDs,
	})

	// Update counts
	a.updateCounts()
}

// GetAlbumByID returns a specific album by ID
func (a *Artist) GetAlbumByID(albumID uint64) *AlbumEntry {
	for i, album := range a.Albums {
		if album.AlbumID == albumID {
			return &a.Albums[i]
		}
	}
	return nil
}

// Scan implements the Scanner interface for SQL
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

// Value implements the Valuer interface for SQL
func (m *Artist) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// MergeAlbums merges albums from another artist
func (m *Artist) MergeAlbums(otherArtist *Artist) {
	if otherArtist == nil || len(otherArtist.Albums) == 0 {
		return
	}

	// Process each album from the other artist
	for _, otherAlbum := range otherArtist.Albums {
		existingAlbum := m.GetAlbumByID(otherAlbum.AlbumID)

		if existingAlbum != nil {
			// Merge track IDs
			existingAlbum.TrackIDs = mergeTrackIDs(existingAlbum.TrackIDs, otherAlbum.TrackIDs)

			// Update album name if this one is empty
			if existingAlbum.AlbumName == "" && otherAlbum.AlbumName != "" {
				existingAlbum.AlbumName = otherAlbum.AlbumName
			}
		} else {
			// Add new album
			m.Albums = append(m.Albums, otherAlbum)
		}
	}

	// Update counts
	m.updateCounts()
}

// updateCounts recalculates album and track counts
func (m *Artist) updateCounts() {
	m.AlbumCount = len(m.Albums)
	m.TrackCount = m.CalculateTrackCount()
}

// CalculateTrackCount returns the total number of tracks
func (m *Artist) CalculateTrackCount() int {
	trackCount := 0
	for _, album := range m.Albums {
		trackCount += len(album.TrackIDs)
	}
	return trackCount
}

// MergeTrackIDsByAlbum merges track IDs for a specific album
func (m *Artist) MergeTrackIDsByAlbum(albumID uint64, albumName string, trackIDs []uint64) {
	album := m.GetAlbumByID(albumID)

	if album != nil {
		// Merge into existing album
		album.TrackIDs = mergeTrackIDs(album.TrackIDs, trackIDs)
	} else {
		// Create new album
		m.Albums = append(m.Albums, AlbumEntry{
			AlbumID:   albumID,
			AlbumName: albumName,
			TrackIDs:  trackIDs,
		})
	}

	// Update counts
	m.updateCounts()
}

// Helper function to merge track IDs without duplicates
func mergeTrackIDs(existing, new []uint64) []uint64 {
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

// GetTrackIDsByAlbum returns track IDs for a specific album
func (m *Artist) GetTrackIDsByAlbum(albumID uint64) []uint64 {
	album := m.GetAlbumByID(albumID)
	if album == nil {
		return []uint64{}
	}
	return album.TrackIDs
}

// GetAllTrackIDs returns all track IDs across all albums
func (m *Artist) GetAllTrackIDs() []uint64 {
	var allTrackIDs []uint64
	for _, album := range m.Albums {
		allTrackIDs = append(allTrackIDs, album.TrackIDs...)
	}
	return allTrackIDs
}

// Merge merges this artist with another one
func (m *Artist) Merge(other MediaData) {
	otherArtist, ok := other.(*Artist)
	if !ok {
		return
	}

	// Initialize Details if nil
	if m.Details == nil && otherArtist.Details != nil {
		m.Details = &MediaDetails{}
	}

	if m.Details != nil && otherArtist.Details != nil {
		m.Details.Merge(otherArtist.Details)
	}

	// Merge similar artists if this one doesn't have any
	if len(m.SimilarArtists) == 0 && len(otherArtist.SimilarArtists) > 0 {
		m.SimilarArtists = otherArtist.SimilarArtists
	}

	// Merge biography if this one is empty
	if m.Biography == "" && otherArtist.Biography != "" {
		m.Biography = otherArtist.Biography
	}

	// Merge genres
	if len(m.Genres) == 0 && len(otherArtist.Genres) > 0 {
		m.Genres = otherArtist.Genres
	}

	// Merge years if not set
	if m.StartYear == 0 && otherArtist.StartYear != 0 {
		m.StartYear = otherArtist.StartYear
	}

	if m.EndYear == 0 && otherArtist.EndYear != 0 {
		m.EndYear = otherArtist.EndYear
	}

	// Merge rating if this one is zero
	if m.Rating == 0 && otherArtist.Rating != 0 {
		m.Rating = otherArtist.Rating
	}

	// Merge credits if needed
	if m.Credits == nil && otherArtist.Credits != nil {
		m.Credits = otherArtist.Credits
	}

	m.MergeAlbums(otherArtist)
}
