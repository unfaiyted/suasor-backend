package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Episode represents a TV episode
type Episode struct {
	Details      MediaDetails
	Number       int64       `json:"number"`
	SeriesID     uint64      `json:"showID"`
	SyncSeries   SyncClients `json:"syncSeries,omitempty"`
	SeasonID     uint64      `json:"seasonID"`
	SyncSeason   SyncClients `json:"syncSeason,omitempty"`
	SeasonNumber int         `json:"seasonNumber"`
	ShowTitle    string      `json:"showTitle,omitempty"`
	Credits      Credits     `json:"credits,omitempty"`
}

func (e *Episode) AddSyncClient(clientID uint64, seriesID string, seasonID string) {
	itemID := e.SyncSeries.GetClientItemID(clientID)
	if itemID == "" {
		e.SyncSeries.AddClient(clientID, seriesID)
	}
	itemID = e.SyncSeason.GetClientItemID(clientID)
	if itemID == "" {
		e.SyncSeason.AddClient(clientID, seasonID)
	}
}

func (e *Episode) isMediaData() {}

func (e *Episode) GetDetails() MediaDetails { return e.Details }
func (e *Episode) GetMediaType() MediaType  { return MediaTypeEpisode }

func (e *Episode) GetTitle() string { return e.Details.Title }

func (e *Episode) SetDetails(details MediaDetails) {
	e.Details = details
}

// the clients id stored in the sync clients
func (e *Episode) GetClientSeriesID(clientID uint64) string {
	return e.SyncSeries.GetClientItemID(clientID)
}

func (e *Episode) GetClientSeasonID(clientID uint64) string {
	return e.SyncSeason.GetClientItemID(clientID)
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
