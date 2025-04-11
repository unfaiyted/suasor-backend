package types

// ClientType represents different types of clients
type ClientCategory string

const (
	ClientCategoryAutomation ClientCategory = "automation"
	ClientCategoryMedia      ClientCategory = "media"
	ClientCategoryAI         ClientCategory = "ai"
	ClientCategoryMetadata   ClientCategory = "metadata"
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

func (c MediaClientType) AsGenericType() ClientType {
	return ClientType(c)
}

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

type MetadataClientType string

const (
	MetadataClientTypeTMDB    MetadataClientType = "tmdb"
	MetadataClientTypeTrakt   MetadataClientType = "trakt"
	MetadataClientTypeUnknown MetadataClientType = "unknown"
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

	ClientTypeTMDB  ClientType = "tmdb"
	ClientTypeTrakt ClientType = "trakt"
)

func (c ClientType) String() string {
	return string(c)
}

func (c ClientType) AsCategory() ClientCategory {
	// Media clients
	mediaClients := map[ClientType]bool{
		ClientTypeEmby: true, ClientTypeJellyfin: true,
		ClientTypePlex: true, ClientTypeSubsonic: true,
	}
	if mediaClients[c] {
		return ClientCategoryMedia
	}

	// Automation clients
	automationClients := map[ClientType]bool{
		ClientTypeRadarr: true, ClientTypeSonarr: true, ClientTypeLidarr: true,
	}
	if automationClients[c] {
		return ClientCategoryAutomation
	}

	// AI clients
	aiClients := map[ClientType]bool{
		ClientTypeClaude: true, ClientTypeOpenAI: true, ClientTypeOllama: true,
	}
	if aiClients[c] {
		return ClientCategoryAI
	}

	// Metadata clients
	metadataClients := map[ClientType]bool{
		ClientTypeTMDB: true, ClientTypeTrakt: true,
	}
	if metadataClients[c] {
		return ClientCategoryMetadata
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

func (c ClientType) AsMediaClientType() MediaClientType {
	switch c {
	case ClientTypeEmby:
		return MediaClientTypeEmby
	case ClientTypeJellyfin:
		return MediaClientTypeJellyfin
	case ClientTypePlex:
		return MediaClientTypePlex
	case ClientTypeSubsonic:
		return MediaClientTypeSubsonic
	default:
		return MediaClientTypeUnknown
	}
}

func (ClientType) isClientConfig() {}

func (c ClientCategory) String() string {
	return string(c)
}

// IsMedia checks if this category is a media client category
func (c ClientCategory) IsMedia() bool {
	return c == ClientCategoryMedia
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

func (c MetadataClientType) AsCategory() ClientCategory {
	return ClientCategoryMetadata
}

func (c MetadataClientType) AsGenericClient() ClientType {
	switch c {
	case MetadataClientTypeTMDB:
		return ClientTypeTMDB
	case MetadataClientTypeTrakt:
		return ClientTypeTrakt
	default:
		return ClientTypeUnknown
	}
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
	case "*types.TMDBConfig":
		return ClientTypeTMDB
	case "*types.TraktConfig":
		return ClientTypeTrakt
	default:
		return ClientTypeUnknown
	}
}
