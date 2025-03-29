package providers

import (
	"context"
	"suasor/client/automation/types"
	"suasor/types/models"
	"suasor/types/requests"
	"time"
)

// AutomationProvider defines a common interface for media automation tools (Radarr, Sonarr, Lidarr)
type AutomationProvider interface {
	isAutomationProvider()
}

// SystemProvider defines methods to interact with system-level features
type SystemProvider interface {
	GetSystemStatus(ctx context.Context) (*types.SystemStatus, error)
}

// LibraryProvider defines methods to interact with the library
type LibraryProvider interface {
	GetLibraryItems(ctx context.Context, options *types.LibraryQueryOptions) ([]*models.AutomationMediaItem[types.AutomationData], error)
}

// MediaProvider defines methods to interact with media items
type MediaProvider interface {
	// GetsBy ExternalID, not our internal ID
	GetMediaByID(ctx context.Context, id string) (models.AutomationMediaItem[types.AutomationData], error)
	AddMedia(ctx context.Context, req requests.AutomationMediaAddRequest) (models.AutomationMediaItem[types.AutomationData], error)
	UpdateMedia(ctx context.Context, id string, req requests.AutomationMediaUpdateRequest) (models.AutomationMediaItem[types.AutomationData], error)
	DeleteMedia(ctx context.Context, id string) error
}

// SearchProvider defines methods for searching
type SearchProvider interface {
	// Search operations
	SearchMedia(ctx context.Context, query string, options *types.SearchOptions) ([]*models.AutomationMediaItem[types.AutomationData], error)
}

// ProfileProvider defines methods for working with profiles
type ProfileProvider interface {
	// Quality profile methods
	GetQualityProfiles(ctx context.Context) ([]types.QualityProfile, error)
	GetMetadataProfiles(ctx context.Context) ([]types.MetadataProfile, error)
}

// TagProvider defines methods for working with tags
type TagProvider interface {
	GetTags(ctx context.Context) ([]types.Tag, error)
	CreateTag(ctx context.Context, tag string) (types.Tag, error)
}

// CalendarProvider defines methods for working with calendar events
type CalendarProvider interface {
	GetCalendar(ctx context.Context, start, end time.Time) ([]*models.AutomationCalendarItem[types.AutomationData], error)
}

// CommandProvider defines methods for executing commands
type CommandProvider interface {
	ExecuteCommand(ctx context.Context, command types.Command) (types.CommandResult, error)
}
type CombinedProvider interface {
	MediaProvider
	LibraryProvider
	SearchProvider
	ProfileProvider
	TagProvider
	CalendarProvider
	CommandProvider
}

func AsCombinedProvider(client AutomationProvider) (CombinedProvider, bool) {
	provider, ok := client.(CombinedProvider)
	return provider, ok
}
