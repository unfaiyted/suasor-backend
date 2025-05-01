package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type TrackEntry struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	TrackID uint64 `json:"trackID"`
}

type TrackEntries []*TrackEntry

func (t TrackEntries) AddTrack(id uint64, track Track) {
	t = append(t, &TrackEntry{
		Number:  track.Number,
		Title:   track.Details.Title,
		TrackID: id,
	})
}

func (t *TrackEntries) GetTrackIDs() []uint64 {
	trackIDs := make([]uint64, 0, len(*t))
	for _, track := range *t {
		trackIDs = append(trackIDs, track.TrackID)
	}
	return trackIDs
}

func (t *TrackEntries) Merge(other TrackEntries) {
	for _, otherTrack := range other {
		found := false
		for i, existingTrack := range *t {
			if existingTrack.TrackID == otherTrack.TrackID {
				// Update existing entry
				(*t)[i].Number = otherTrack.Number
				(*t)[i].Title = otherTrack.Title
				found = true
				break
			}
		}
		if !found {
			// Add new entry
			*t = append(*t, otherTrack)
		}
	}
}

// MusicAlbum represents a music album
type Album struct {
	MediaData `json:"-"`

	Details    *MediaDetails `json:"details"`
	ArtistID   uint64        `json:"artistID"`
	ArtistName string        `json:"artistName"`
	TrackCount int           `json:"trackCount"`
	Credits    *Credits      `json:"credits,omitempty"`
	Tracks     *TrackEntries `json:"tracks,omitempty"`
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
	m.Credits.Merge(otherAlbum.Credits)
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

func (m *Album) GetTrackIDs() []uint64 {
	trackIDs := make([]uint64, 0, len(*m.Tracks))
	for _, track := range *m.Tracks {
		trackIDs = append(trackIDs, track.TrackID)
	}
	return trackIDs
}
