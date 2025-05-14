package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"suasor/clients/media/types"
	client "suasor/clients/types"
	"time"
)

// MediaItem is the base type for all media items
type MediaItem[T types.MediaData] struct {
	BaseModel
	ID      uint64 `json:"id" gorm:"primaryKey;autoIncrement"` // Internal ID
	UUID    string `json:"uuid" gorm:"type:uuid;uniqueIndex"`  // Stable UUID for syncing
	OwnerID uint64 `json:"ownerId"`                            // ID of the user that owns this item, 0 for system owned items

	SyncClients SyncClients `json:"syncClients" gorm:"type:jsonb"` // Client IDs for this item (mapping client to their IDs)

	ExternalIDs types.ExternalIDs `json:"externalIds" gorm:"type:jsonb"` // External IDs for this item (TMDB, IMDB, etc.)
	IsPublic    bool              `json:"isPublic"`                      // Whether this item is public or not

	Type types.MediaType `json:"type" gorm:"type:varchar(50)"` // Type of media (movie, show, episode, etc.)

	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"releaseDate,omitempty"`
	ReleaseYear int       `json:"releaseYear,omitempty"`

	StreamURL   string `json:"streamUrl,omitempty" gorm:"size:1024"`
	DownloadURL string `json:"downloadUrl,omitempty" gorm:"size:1024"`
	Data        T      `json:"data" gorm:"type:jsonb"` // Type-specific media data
}

func (MediaItem[T]) TableName() string {
	return "media_items"
}

func NewMediaItem[T types.MediaData](data T) *MediaItem[T] {
	// Initialize with empty arrays
	clientIDs := make(SyncClients, 0)
	externalIDs := make(types.ExternalIDs, 0)
	itemType := types.GetMediaType[T]()

	// Make sure data has a valid Details field
	details := data.GetDetails()
	title := ""
	releaseDate := time.Time{}

	if details != nil {
		// Safe to access Details fields
		if details.ExternalIDs != nil {
			externalIDs = details.ExternalIDs
		}
		title = details.Title
		releaseDate = details.ReleaseDate
	}

	return &MediaItem[T]{
		BaseModel: BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UUID:        uuid.New().String(),
		Type:        itemType,
		SyncClients: clientIDs,
		IsPublic:    true,
		Title:       title,
		ReleaseDate: releaseDate,
		Data:        data,
		ExternalIDs: externalIDs,
	}
}

func NewMediaItemCopy[T types.MediaData, U types.MediaData](item *MediaItem[T]) *MediaItem[U] {
	// Create a proper cast of the data based on media type
	var expectedType U

	expectedType = any(item.Data).(U)

	NewItem := NewMediaItem[U](expectedType)
	NewItem.UUID = item.UUID
	NewItem.SyncClients = item.SyncClients
	NewItem.ExternalIDs = item.ExternalIDs
	NewItem.IsPublic = item.IsPublic
	NewItem.OwnerID = item.OwnerID
	NewItem.Title = item.Title
	NewItem.ReleaseDate = item.ReleaseDate
	NewItem.ReleaseYear = item.ReleaseYear
	NewItem.StreamURL = item.StreamURL
	NewItem.DownloadURL = item.DownloadURL
	NewItem.CreatedAt = item.CreatedAt
	NewItem.UpdatedAt = item.UpdatedAt

	return NewItem
}

func (m *MediaItem[T]) GetData() T {
	return m.Data
}
func (m *MediaItem[T]) GetTitle() string {
	return m.Title
}
func (m *MediaItem[T]) GetDescription() string {
	return m.Data.GetDetails().Description
}

func (m *MediaItem[T]) SetData(data T) {
	m.Data = data
}

func (m *MediaItem[T]) SetIsPublic(isPublic bool) {
	m.IsPublic = isPublic
}

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
func (m *MediaItem[T]) SetClientInfo(clientID uint64, clientType client.ClientType, clientItemKey string) {
	// Add to SyncClients
	found := false

	for i, id := range m.SyncClients {
		if id.ID == clientID && id.Type == clientType {
			// Update existing entry
			m.SyncClients[i].ItemID = clientItemKey
			found = true
			break
		}
	}

	if !found {
		// Add new entry
		m.SyncClients = append(m.SyncClients, &SyncClient{
			ID:     clientID,
			Type:   clientType,
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
		m.ExternalIDs = append(m.ExternalIDs, types.ExternalID{
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
		m.SyncClients = append(m.SyncClients, &SyncClient{
			ID:     clientID,
			Type:   clientType,
			ItemID: itemID,
		})
	}
}

func (m *MediaItem[T]) IsSyncClient(clientID uint64) bool {
	for _, syncClient := range m.SyncClients {
		if syncClient.ID == clientID {
			return true
		}
	}
	return false
}

// IsList returns true if the media item is a list
func (m *MediaItem[T]) IsList() bool {
	return m.Type == types.MediaTypePlaylist || m.Type == types.MediaTypeCollection
}
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
	externalIDs := make(types.ExternalIDs, 0)

	switch mediaType {
	case types.MediaTypeMovie:
		return &MediaItem[*types.Movie]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeSeries:
		return &MediaItem[*types.Series]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeEpisode:
		return &MediaItem[*types.Episode]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeSeason:
		return &MediaItem[*types.Season]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeTrack:
		return &MediaItem[*types.Track]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeAlbum:
		return &MediaItem[*types.Album]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeArtist:
		return &MediaItem[*types.Artist]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypeCollection:
		return &MediaItem[*types.Collection]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	case types.MediaTypePlaylist:
		return &MediaItem[*types.Playlist]{
			Type:        mediaType,
			SyncClients: clientIDs,
			ExternalIDs: externalIDs,
		}, nil
	default:
		return nil, fmt.Errorf("unknown media type: %s", mediaType)
	}
}

func (existingItem *MediaItem[T]) Merge(newItem *MediaItem[T]) *MediaItem[T] {
	// Merge sync clients
	existingItem.SyncClients.Merge(newItem.SyncClients)
	existingItem.ExternalIDs.Merge(newItem.ExternalIDs)

	existingDetails := existingItem.Data.GetDetails()
	newDetails := newItem.Data.GetDetails()
	existingDetails.ExternalIDs = existingItem.ExternalIDs

	// Update data fields
	if existingDetails.Title == "" {
		existingDetails.Title = newDetails.Title
	}
	if existingDetails.Description == "" {
		existingDetails.Description = newDetails.Description
	}
	if existingDetails.ContentRating == "" {
		existingDetails.ContentRating = newDetails.ContentRating
	}
	if existingDetails.ContentRating == "" {
		existingDetails.ContentRating = newDetails.ContentRating
	}
	if existingDetails.Studio == "" {
		existingDetails.Studio = newDetails.Studio
	}

	existingDetails.Genres = mergeStringArray(existingDetails.Genres, newDetails.Genres)
	existingDetails.Ratings = mergeRatings(existingDetails.Ratings, newDetails.Ratings)

	// Artworks
	if existingDetails.Artwork.Poster == "" {
		existingDetails.Artwork.Poster = newDetails.Artwork.Poster
	}
	if existingDetails.Artwork.Banner == "" {
		existingDetails.Artwork.Banner = newDetails.Artwork.Banner
	}
	if existingDetails.Artwork.Thumbnail == "" {
		existingDetails.Artwork.Thumbnail = newDetails.Artwork.Thumbnail
	}
	if existingDetails.Artwork.Logo == "" {
		existingDetails.Artwork.Logo = newDetails.Artwork.Logo
	}

	if existingDetails.ReleaseYear == 0 {
		existingDetails.ReleaseYear = newDetails.ReleaseYear
	}
	if existingDetails.ReleaseDate.IsZero() {
		existingDetails.ReleaseDate = newDetails.ReleaseDate
	}

	if existingItem.Title == "" {
		existingItem.Title = newDetails.Title
	}
	if existingItem.ReleaseYear == 0 {
		existingItem.ReleaseYear = newDetails.ReleaseYear
	}
	if existingItem.ReleaseDate.IsZero() {
		existingItem.ReleaseDate = newDetails.ReleaseDate
	}

	return existingItem
}
