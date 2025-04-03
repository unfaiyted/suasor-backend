package types

// ClientType represents different types of clients
type ClientCategory string

const (
	ClientCategoryAutomation ClientCategory = "automation"
	ClientCategoryMedia      ClientCategory = "media"
	ClientCategoryAI         ClientCategory = "ai"
	ClientCategoryUnknown    ClientCategory = "unknown"
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

func (c MediaClientType) AsCategory() ClientCategory {
	return ClientCategoryMedia
}

// ClientType represents different types of download clients
type AutomationClientType string

const (
	AutomationClientTypeRadarr  AutomationClientType = "radarr"
	AutomationClientTypeSonarr  AutomationClientType = "sonarr"
	AutomationClientTypeLidarr  AutomationClientType = "lidarr"
	AutomationClientTypeUnknown AutomationClientType = "unknown"
)

func (c MediaClientType) AsGenericClient() ClientType {
	switch c {
	case MediaClientTypePlex:
		return ClientTypePlex
	case MediaClientTypeJellyfin:
		return ClientTypeJellyfin
	case MediaClientTypeEmby:
		return ClientTypeEmby
	case MediaClientTypeSubsonic:
		return ClientTypeSubsonic
	default:
		return ClientTypeUnknown
	}
}

func (c AutomationClientType) AsGenericClient() ClientType {
	switch c {
	case AutomationClientTypeRadarr:
		return ClientTypeRadarr
	case AutomationClientTypeSonarr:
		return ClientTypeSonarr
	case AutomationClientTypeLidarr:
		return ClientTypeLidarr
	default:
		return ClientTypeUnknown
	}
}

type AIClientType string

const (
	AIClientTypeClaude  AIClientType = "claude"
	AIClientTypeOpenAI  AIClientType = "openai"
	AIClientTypeOllama  AIClientType = "ollama"
	AIClientTypeUnknown AIClientType = "unknown"
)

type ClientType string

const (
	ClientTypeEmby     ClientType = "emby"
	ClientTypeJellyfin ClientType = "jellyfin"
	ClientTypePlex     ClientType = "plex"
	ClientTypeSubsonic ClientType = "subsonic"

	ClientTypeRadarr ClientType = "radarr"
	ClientTypeSonarr ClientType = "sonarr"
	ClientTypeLidarr ClientType = "lidarr"

	ClientTypeUnknown ClientType = "unknown"

	ClientTypeClaude ClientType = "claude"
	ClientTypeOpenAI ClientType = "openai"
	ClientTypeOllama ClientType = "ollama"
)

func (c ClientType) String() string {
	return string(c)
}

func (c ClientType) AsCategory() ClientCategory {
	switch c {
	case ClientTypeEmby:
	case ClientTypeJellyfin:
	case ClientTypePlex:
	case ClientTypeSubsonic:
		return ClientCategoryMedia
	case ClientTypeRadarr:
	case ClientTypeSonarr:
	case ClientTypeLidarr:
		return ClientCategoryAutomation
	case ClientTypeClaude:
	case ClientTypeOpenAI:
	case ClientTypeOllama:
		return ClientCategoryAI
	default:
		return ClientCategoryUnknown
	}
	return ClientCategoryUnknown
}

// Make ClientType implement ClientConfig interface
func (c ClientType) GetType() ClientType {
	return c
}

func (c ClientType) GetCategory() ClientCategory {
	return c.AsCategory()
}

func (ClientType) isClientConfig() {}

func (c ClientCategory) String() string {
	return string(c)
}

func (c MediaClientType) String() string {
	return string(c)
}

func (c AutomationClientType) String() string {
	return string(c)
}

func (c AutomationClientType) AsCategory() ClientCategory {
	return ClientCategoryAutomation
}

func (c AIClientType) AsCategory() ClientCategory {
	return ClientCategoryAI
}

func GetClientTypeFromTypeName(typeName string) ClientType {
	switch typeName {
	case "*types.EmbyConfig":
		return ClientTypeEmby
	case "*types.JellyfinConfig":
		return ClientTypeJellyfin
	case "*types.RadarrConfig":
		return ClientTypeRadarr
	case "*types.SonarrConfig":
		return ClientTypeSonarr
	case "*types.LidarrConfig":
		return ClientTypeLidarr
	case "*types.ClaudeConfig":
		return ClientTypeClaude
	case "*types.OpenAIConfig":
		return ClientTypeOpenAI
	case "*types.OllamaConfig":
		return ClientTypeOllama
	case "*types.SubsonicConfig":
		return ClientTypeSubsonic
	case "*types.PlexConfig":
		return ClientTypePlex
	default:
		return ClientTypeUnknown
	}
}
