package types

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

func (a *Album) AddSyncClient(clientID uint64, artistID string) {
	itemID := a.SyncArtist.GetClientItemID(clientID)
	if itemID == "" {
		a.SyncArtist.AddClient(clientID, artistID)
	}
}
