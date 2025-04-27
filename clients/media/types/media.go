package types

type SyncClient struct {
	// ID of the client that this external ID belongs to (optional for service IDs like TMDB)
	ID uint64 `json:"clientID,omitempty"`
	// The actual ID value in the external system
	ItemID string `json:"itemID"`
}

type SyncClients []SyncClient

func (s SyncClients) AddClient(clientID uint64, itemID string) {
	// check if client ID already exists
	found := false
	for i, cID := range s {
		if cID.ID == clientID {
			// Update existing ID
			s[i].ItemID = itemID
			found = true
			break
		}
	}
	if !found {
		// Add new ID if not found
		s = append(s, SyncClient{
			ID:     clientID,
			ItemID: itemID,
		})
	}
}

func (s SyncClients) GetClientItemID(clientID uint64) string {
	for _, cID := range s {
		if cID.ID == clientID {
			return cID.ItemID
		}
	}
	return ""
}

// Artwork holds different types of artwork
type Artwork struct {
	Poster     string `json:"poster,omitempty"`
	Background string `json:"background,omitempty"`
	Banner     string `json:"banner,omitempty"`
	Thumbnail  string `json:"thumbnail,omitempty"`
	Logo       string `json:"logo,omitempty"`
}

// Person represents someone involved with the media

type MediaData interface {
	isMediaData()
	GetDetails() MediaDetails
	GetMediaType() MediaType
	SetDetails(MediaDetails)
}

func NewItem[T MediaData]() T {
	var zero T
	return zero
}
