// interfaces/media_client.go
package interfaces

import ()

// MediaClientType represents different types of media clients
type MediaClientType string

const (
	MediaClientTypePlex     MediaClientType = "plex"
	MediaClientTypeJellyfin MediaClientType = "jellyfin"
	MediaClientTypeEmby     MediaClientType = "emby"
	MediaClientTypeSubsonic MediaClientType = "subsonic"
)

// MediaClient defines basic client information that all providers must implement
type MediaClient interface {
	GetClientID() uint64
	GetClientType() MediaClientType
}
