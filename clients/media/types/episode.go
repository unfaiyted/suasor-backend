package types

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

// the clients id stored in the sync clients
func (e *Episode) GetClientSeriesID(clientID uint64) string {
	return e.SyncSeries.GetClientItemID(clientID)
}

func (e *Episode) GetClientSeasonID(clientID uint64) string {
	return e.SyncSeason.GetClientItemID(clientID)
}
