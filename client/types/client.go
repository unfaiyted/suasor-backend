package types

// ClientType represents different types of clients
type ClientType string

const (
	ClientTypeAutomation ClientType = "automation"
	ClientTypeMedia      ClientType = "media"
	ClientTypeAI         ClientType = "ai"
)

// MediaClientType represents different types of media clients
type MediaClientType string

const (
	MediaClientTypePlex     MediaClientType = "plex"
	MediaClientTypeJellyfin MediaClientType = "jellyfin"
	MediaClientTypeEmby     MediaClientType = "emby"
	MediaClientTypeSubsonic MediaClientType = "subsonic"
	MediaClientTypeUnknown  MediaClientType = "unknown"
)

// ClientType represents different types of download clients
type AutomationClientType string

const (
	AutomationClientTypeRadarr AutomationClientType = "radarr"
	AutomationClientTypeSonarr AutomationClientType = "sonarr"
	AutomationClientTypeLidarr AutomationClientType = "lidarr"
	AutmationClientTypeUnknown AutomationClientType = "unknown"
)

type AIClientType string

const (
	AIClientTypeClaude  AIClientType = "claude"
	AIClientTypeOpenAI  AIClientType = "openai"
	AIClientTypeOllama  AIClientType = "ollama"
	AIClientTypeUnknown AIClientType = "unknown"
)

func (c ClientType) String() string {
	return string(c)
}

func (c MediaClientType) String() string {
	return string(c)
}

func (c AutomationClientType) String() string {
	return string(c)
}

func (c MediaClientType) AsClientType() ClientType {
	return ClientTypeMedia
}

func (c AutomationClientType) AsClientType() ClientType {
	return ClientTypeAutomation
}

func (c AIClientType) AsClientType() ClientType {
	return ClientTypeAI
}
