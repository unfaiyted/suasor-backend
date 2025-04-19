package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"suasor/client/media/types"
	client "suasor/client/types"
	"time"
)

type ListItem struct {
	ItemUUID    string    `json:"itemUuid"`
	Position    int       `json:"position"`
	LastChanged time.Time `json:"lastChanged"`
}

type ListItems []ListItem

type ClientMediaItems struct {
	MediaItems
}

type MediaItems struct {
	Details *MediaItem[types.ListData]

	movies      map[string]*MediaItem[*types.Movie]
	series      map[string]*MediaItem[*types.Series]
	seasons     map[string]*MediaItem[*types.Season]
	episodes    map[string]*MediaItem[*types.Episode]
	artists     map[string]*MediaItem[*types.Artist]
	albums      map[string]*MediaItem[*types.Album]
	tracks      map[string]*MediaItem[*types.Track]
	playlists   map[string]*MediaItem[*types.Playlist]
	collections map[string]*MediaItem[*types.Collection]

	ListType     types.ListType
	ListOriginID uint64 // 0 for internal db, otherwise external client/ProviderID
	OwnerID      uint64

	Order ListItems

	totalItems int
}

func (m *MediaItems) AddListItem(itemUUID string, itemPosition int) {
	m.Order = append(m.Order, ListItem{
		ItemUUID:    itemUUID,
		Position:    itemPosition,
		LastChanged: time.Now(),
	})
}

func (m *MediaItems) GetTotalItems() int {
	return m.totalItems
}

// AddMovie adds a movie to the media items
func (m *MediaItems) AddMovie(item *MediaItem[*types.Movie]) {
	m.movies[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// AddSeries adds a series to the media items
func (m *MediaItems) AddSeries(item *MediaItem[*types.Series]) {
	m.series[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// AddSeason adds a season to the media items
func (m *MediaItems) AddSeason(item *MediaItem[*types.Season]) {
	m.seasons[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// AddEpisode adds an episode to the media items
func (m *MediaItems) AddEpisode(item *MediaItem[*types.Episode]) {
	m.episodes[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// AddArtist adds an artist to the media items
func (m *MediaItems) AddArtist(item *MediaItem[*types.Artist]) {
	m.artists[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// AddAlbum adds an album to the media items
func (m *MediaItems) AddAlbum(item *MediaItem[*types.Album]) {
	m.albums[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// AddTrack adds a track to the media items
func (m *MediaItems) AddTrack(item *MediaItem[*types.Track]) {
	m.tracks[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// AddPlaylist adds a playlist to the media items
func (m *MediaItems) AddPlaylist(item *MediaItem[*types.Playlist]) {
	m.playlists[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// AddCollection adds a collection to the media items
func (m *MediaItems) AddCollection(item *MediaItem[*types.Collection]) {
	m.collections[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

// MediaItem is the base type for all media items
type MediaItem[T types.MediaData] struct {
	BaseModel
	ID          uint64      `json:"id" gorm:"primaryKey;autoIncrement"` // Internal ID
	UUID        string      `json:"uuid" gorm:"type:uuid;uniqueIndex"`  // Stable UUID for syncing
	SyncClients SyncClients `json:"syncClients" gorm:"type:jsonb"`      // Client IDs for this item (mapping client to their IDs)
	ExternalIDs ExternalIDs `json:"externalIds" gorm:"type:jsonb"`      // External IDs for this item (TMDB, IMDB, etc.)
	OwnerID     uint64      `json:"ownerId"`                            // ID of the user that owns this item, 0 for system owned items

	Type types.MediaType `json:"type" gorm:"type:varchar(50)"` // Type of media (movie, show, episode, etc.)

	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"releaseDate,omitempty"`
	ReleaseYear int       `json:"releaseYear,omitempty"`

	StreamURL   string `json:"streamUrl,omitempty" gorm:"size:1024"`
	DownloadURL string `json:"downloadUrl,omitempty" gorm:"size:1024"`
	Data        T      `json:"data" gorm:"type:jsonb"` // Type-specific media data
}

func NewMediaItem[T types.MediaData](itemType types.MediaType, data T) *MediaItem[T] {
	// Initialize with empty arrays
	clientIDs := make(SyncClients, 0)
	externalIDs := make(ExternalIDs, 0)
	return &MediaItem[T]{
		UUID:        uuid.New().String(),
		Type:        itemType,
		SyncClients: clientIDs,
		Data:        data,
		ExternalIDs: externalIDs,
	}
}

func (m *MediaItem[T]) SetData(data T) {
	m.Data = data
}

// ExternalID represents an ID that identifies this media item in an external system
type SyncClient struct {
	// ID of the client that this external ID belongs to (optional for service IDs like TMDB)
	ID uint64 `json:"clientId,omitempty"`
	// Type of client this ID belongs to (optional for service IDs)
	Type client.ClientType `json:"clientType,omitempty" gorm:"type:varchar(50)"`
	// The actual ID value in the external system
	ItemID string `json:"itemId"`
}

type SyncClients []SyncClient

func (s SyncClients) AddClient(clientID uint64, clientType client.ClientType, itemID string) {
	s = append(s, SyncClient{
		ID:     clientID,
		Type:   clientType,
		ItemID: itemID,
	})
}

func (s SyncClients) GetClientItemID(clientID uint64) string {
	for _, id := range s {
		if id.ID == clientID {
			return id.ItemID
		}
	}
	return ""
}

func (s SyncClients) GetByClientID(clientID uint64) (*SyncClient, bool) {
	for _, client := range s {
		if client.ID == clientID {
			return &client, true
		}
	}
	return &SyncClient{}, false
}

type ExternalID struct {
	Source string `json:"source"` // e.g., "tmdb", "imdb", "trakt", "tvdb"
	ID     string `json:"id"`     // The actual ID
}

type ExternalIDs []ExternalID

func (ids ExternalIDs) GetID(source string) string {
	for _, id := range ids {
		if id.Source == source {
			return id.ID
		}
	}
	return ""
}

func (MediaItem[T]) TableName() string {
	return "media_items"
}

// Custom serialization for GORM and JSON

// Value implements driver.Valuer for database storage
func (m MediaItem[T]) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(m)
}

// Scan implements sql.Scanner for database retrieval
func (m *MediaItem[T]) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// First unmarshal the common fields
	type Alias MediaItem[T]
	aux := &struct {
		*Alias
		Data json.RawMessage `json:"data"`
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(bytes, &aux); err != nil {
		return fmt.Errorf("error unmarshaling media item: %w", err)
	}

	// Then unmarshal the Data field separately into the appropriate type
	var mediaData T
	if err := json.Unmarshal(aux.Data, &mediaData); err != nil {
		return fmt.Errorf("error unmarshaling data field: %w", err)
	}
	m.Data = mediaData

	return nil
}

// MarshalJSON provides custom JSON serialization
func (m MediaItem[T]) MarshalJSON() ([]byte, error) {
	// Create a temporary structure without the Data field
	type Alias MediaItem[T]

	// Marshal everything together
	return json.Marshal(struct {
		Alias
		Data T `json:"data"`
	}{
		Alias: Alias(m),
		Data:  m.Data,
	})
}

// UnmarshalJSON provides custom JSON deserialization
func (m *MediaItem[T]) UnmarshalJSON(data []byte) error {
	// Create a temporary structure to unmarshall common fields
	type Alias MediaItem[T]
	aux := &struct {
		*Alias
		Data json.RawMessage `json:"data"`
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Unmarshal the data field into the appropriate type
	var mediaData T
	if err := json.Unmarshal(aux.Data, &mediaData); err != nil {
		return fmt.Errorf("error unmarshaling data field: %w", err)
	}
	m.Data = mediaData

	return nil
}

// SetClientInfo adds or updates client ID information for this media item
func (m *MediaItem[T]) SetClientInfo(clientID uint64, clientType client.ClientMediaType, clientItemKey string) {
	// Add to SyncClients
	found := false
	genericType := clientType.AsGenericType()

	for i, id := range m.SyncClients {
		if id.ID == clientID && id.Type == genericType {
			// Update existing entry
			m.SyncClients[i].ItemID = clientItemKey
			found = true
			break
		}
	}

	if !found {
		// Add new entry
		m.SyncClients = append(m.SyncClients, SyncClient{
			ID:     clientID,
			Type:   genericType,
			ItemID: clientItemKey,
		})
	}
}

// AddExternalID adds or updates an external ID for this media item
func (m *MediaItem[T]) AddExternalID(source string, id string) {
	if id == "" {
		return
	}

	// Check if external ID already exists
	found := false
	for i, extID := range m.ExternalIDs {
		if extID.Source == source {
			// Update existing entry
			m.ExternalIDs[i].ID = id
			found = true
			break
		}
	}

	if !found {
		// Add new entry
		m.ExternalIDs = append(m.ExternalIDs, ExternalID{
			Source: source,
			ID:     id,
		})
	}
}

// GetExternalID retrieves an external ID by source
func (m *MediaItem[T]) GetExternalID(source string) (string, bool) {
	for _, extID := range m.ExternalIDs {
		if extID.Source == source {
			return extID.ID, true
		}
	}
	return "", false
}

// GetClientItemID retrieves the item ID for a specific client
func (m *MediaItem[T]) GetClientItemID(clientID uint64) (string, bool) {
	for _, cID := range m.SyncClients {
		if cID.ID == clientID {
			return cID.ItemID, true
		}
	}
	return "", false
}

func (m *MediaItem[T]) AddSyncClient(clientID uint64, clientType client.ClientType, itemID string) {
	if itemID == "" {
		return
	}

	// Check if client ID already exists
	found := false
	for i, cID := range m.SyncClients {
		if cID.ID == clientID && cID.Type == clientType {
			// Update existing entry
			m.SyncClients[i].ItemID = itemID
			found = true
			break
		}
	}

	if !found {
		// Add new entry
		m.SyncClients = append(m.SyncClients, SyncClient{
			ID:     clientID,
			Type:   clientType,
			ItemID: itemID,
		})
	}
}

func (m *MediaItem[T]) GetData() T {
	return m.Data
}

func (m *MediaItem[T]) AsEpisode() (MediaItem[types.Episode], bool) {
	if m.Type != types.MediaTypeEpisode {
		return MediaItem[types.Episode]{}, false
	}
	episode, ok := any(m).(MediaItem[types.Episode])

	return episode, ok
}

func (m *MediaItem[T]) AsMovie() (MediaItem[types.Movie], bool) {
	if m.Type != types.MediaTypeMovie {
		return MediaItem[types.Movie]{}, false
	}
	movie, ok := any(m).(MediaItem[types.Movie])

	return movie, ok
}

func (m *MediaItem[T]) AsSeries() (MediaItem[types.Series], bool) {
	if m.Type != types.MediaTypeSeries {
		return MediaItem[types.Series]{}, false
	}
	show, ok := any(m).(MediaItem[types.Series])

	return show, ok
}

func (m *MediaItem[T]) AsSeason() (MediaItem[types.Season], bool) {
	if m.Type != types.MediaTypeSeason {
		return MediaItem[types.Season]{}, false
	}
	season, ok := any(m).(MediaItem[types.Season])

	return season, ok
}

func (m *MediaItem[T]) AsTrack() (MediaItem[types.Track], bool) {
	if m.Type != types.MediaTypeTrack {
		return MediaItem[types.Track]{}, false
	}
	track, ok := any(m).(MediaItem[types.Track])

	return track, ok
}

func (m *MediaItem[T]) AsAlbum() (MediaItem[types.Album], bool) {
	if m.Type != types.MediaTypeAlbum {
		return MediaItem[types.Album]{}, false
	}
	album, ok := any(m).(MediaItem[types.Album])

	return album, ok
}

func (m *MediaItem[T]) AsArtist() (MediaItem[types.Artist], bool) {
	if m.Type != types.MediaTypeArtist {
		return MediaItem[types.Artist]{}, false
	}
	artist, ok := any(m).(MediaItem[types.Artist])

	return artist, ok
}

func (m *MediaItem[T]) AsCollection() (MediaItem[types.Collection], bool) {
	if m.Type != types.MediaTypeCollection {
		return MediaItem[types.Collection]{}, false
	}
	collection, ok := any(m).(MediaItem[types.Collection])

	return collection, ok
}

func (m *MediaItem[T]) AsPlaylist() (MediaItem[types.Playlist], bool) {
	if m.Type != types.MediaTypePlaylist {
		return MediaItem[types.Playlist]{}, false
	}
	playlist, ok := any(m).(MediaItem[types.Playlist])

	return playlist, ok
}

// IsList returns true if the media item is a list
func (m *MediaItem[T]) IsList() bool {
	return m.Type == types.MediaTypePlaylist || m.Type == types.MediaTypeCollection
}

// IsPlaylist returns true if the media item is a playlist
func (m *MediaItem[T]) IsPlaylist() bool {
	return m.Type == types.MediaTypePlaylist
}
func (m *MediaItem[T]) IsCollection() bool {
	return m.Type == types.MediaTypeCollection
}

// CreateMediaItem creates a new MediaItem of the appropriate type
func CreateMediaItem(mediaType types.MediaType) (any, error) {
	// Initialize with empty arrays for SyncClients and ExternalIDs
	clientIDs := make(SyncClients, 0)
	externalIDs := make(ExternalIDs, 0)

	switch mediaType {
	case types.MediaTypeMovie:
		return &MediaItem[types.Movie]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeSeries:
		return &MediaItem[types.Series]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeEpisode:
		return &MediaItem[types.Episode]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeSeason:
		return &MediaItem[types.Season]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeTrack:
		return &MediaItem[types.Track]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeAlbum:
		return &MediaItem[types.Album]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeArtist:
		return &MediaItem[types.Artist]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeCollection:
		return &MediaItem[types.Collection]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypePlaylist:
		return &MediaItem[types.Playlist]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	default:
		return nil, fmt.Errorf("unknown media type: %s", mediaType)
	}
}
