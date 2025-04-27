package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Track struct {
	Details    MediaDetails
	AlbumID    uint64      `json:"albumID"`
	SyncAlbum  SyncClients `json:"syncAlbum,omitempty"`
	ArtistID   uint64      `json:"artistID"`
	SyncArtist SyncClients `json:"syncArtist,omitempty"`
	AlbumName  string      `json:"albumName"`
	AlbumTitle string      `json:"albumTitle,omitempty"`
	Duration   int         `json:"duration,omitempty"`
	ArtistName string      `json:"artistName,omitempty"`
	Number     int         `json:"trackNumber,omitempty"`
	DiscNumber int         `json:"discNumber,omitempty"`
	Composer   string      `json:"composer,omitempty"`
	Lyrics     string      `json:"lyrics,omitempty"`
	Credits    Credits     `json:"credits,omitempty"`
}

func (t *Track) AddSyncClient(clientID uint64, albumID string, artistID string) {
	itemID := t.SyncAlbum.GetClientItemID(clientID)
	if itemID == "" {
		t.SyncAlbum.AddClient(clientID, albumID)
	}
	itemID = t.SyncArtist.GetClientItemID(clientID)
	if itemID == "" {
		t.SyncArtist.AddClient(clientID, artistID)
	}
}

func (Track) isMediaData()               {}
func (t Track) GetDetails() MediaDetails { return t.Details }
func (t Track) GetMediaType() MediaType  { return MediaTypeTrack }
func (t Track) GetTitle() string         { return t.Details.Title }

func (t *Track) SetDetails(details MediaDetails) {
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
