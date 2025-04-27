package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// MusicAlbum represents a music album
type Album struct {
	Details    MediaDetails
	ArtistID   uint64      `json:"artistID"`
	SyncArtist SyncClients `json:"syncArtist,omitempty"`
	ArtistName string      `json:"artistName"`
	TrackCount int         `json:"trackCount"`
	Credits    Credits     `json:"credits,omitempty"`
	Tracks     []*Track    `json:"tracks,omitempty"`
	TrackIDs   []uint64    `json:"trackIDs,omitempty"`
}

func (a *Album) isMediaData() {}
func (a *Album) AddSyncClient(clientID uint64, artistID string) {
	itemID := a.SyncArtist.GetClientItemID(clientID)
	if itemID == "" {
		a.SyncArtist.AddClient(clientID, artistID)
	}
}

func (a *Album) GetDetails() MediaDetails { return a.Details }
func (a *Album) GetMediaType() MediaType  { return MediaTypeAlbum }

func (a *Album) GetTitle() string { return a.Details.Title }

func (a *Album) SetDetails(details MediaDetails) {
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
