// interfaces/media_client.go
package interfaces

import (
	"suasor/client/media/types"
)

// MediaClient defines basic client information that all providers must implement
type MediaClient interface {
	GetClientID() uint64
	GetClientType() types.MediaClientType
}

type MediaData interface {
	isMediaData()
	GetDetails() types.MediaMetadata
	GetMediaType() types.MediaType
}

func (types.Movie) isMediaData()      {}
func (types.TVShow) isMediaData()     {}
func (types.Episode) isMediaData()    {}
func (types.Track) isMediaData()      {}
func (types.Artist) isMediaData()     {}
func (types.Album) isMediaData()      {}
func (types.Season) isMediaData()     {}
func (types.Collection) isMediaData() {}
func (types.Playlist) isMediaData()   {}

// MediaItem is the base type for all media items
type MediaItem[T MediaData] struct {
	ID          uint64                `json:"ID" gorm:"primaryKey"` // internal ID
	ExternalID  string                `json:"externalID" gorm:"index"`
	ClientID    uint64                `json:"clientID"  gorm:"index"` // internal ClientID
	ClientType  types.MediaClientType `json:"clientType"`             // internal Client Type "plex", "jellyfin", etc.
	Type        string                `json:"type"`                   // "movie", "tvshow", "episode", "music","playlist","artist"
	StreamURL   string                `json:"streamUrl,omitempty"`
	DownloadURL string                `json:"downloadUrl,omitempty"`
	Data        T
}

// Implement this interface for MediaItem[T]
func (m *MediaItem[MediaData]) SetClientInfo(clientID uint64, clientType types.MediaClientType, clientItemKey string) {
	m.ClientID = clientID
	m.ClientType = clientType
	m.ExternalID = clientItemKey
}

func (m *MediaItem[MediaData]) GetData() MediaData {
	return m.Data
}
