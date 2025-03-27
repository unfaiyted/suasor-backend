package types

// ClientType represents different types of clients
type ClientType string

const (
	ClientTypeRadarr   ClientType = "radarr"
	ClientTypeSonarr   ClientType = "sonarr"
	ClientTypeLidarr   ClientType = "lidarr"
	ClientTypeSubsonic ClientType = "subsonic"
	ClientTypeEmby     ClientType = "emby"
	ClientTypeJellyfin ClientType = "jellyfin"
	ClientTypePlex     ClientType = "plex"
)

// MediaClientType represents different types of media clients
type MediaClientType string

const (
	MediaClientTypePlex     MediaClientType = "plex"
	MediaClientTypeJellyfin MediaClientType = "jellyfin"
	MediaClientTypeEmby     MediaClientType = "emby"
	MediaClientTypeSubsonic MediaClientType = "subsonic"
)

// ClientType represents different types of download clients
type AutomationClientType string

const (
	AutomationClientTypeRadarr AutomationClientType = "radarr"
	AutomationClientTypeSonarr AutomationClientType = "sonarr"
	AutomationClientTypeLidarr AutomationClientType = "lidarr"
)
