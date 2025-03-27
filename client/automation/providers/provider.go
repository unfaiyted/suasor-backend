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

	// System methods
	GetSystemStatus(ctx context.Context) (types.SystemStatus, error)

	// Library management
	GetLibraryItems(ctx context.Context, options *types.LibraryQueryOptions) ([]models.AutomationMediaItem[types.AutomationData], error)
	GetMediaByID(ctx context.Context, id int64) (models.AutomationMediaItem[types.AutomationData], error)
	AddMedia(ctx context.Context, req requests.AutomationMediaAddRequest) (models.AutomationMediaItem[types.AutomationData], error)
	UpdateMedia(ctx context.Context, id int64, req requests.AutomationMediaUpdateRequest) (models.AutomationMediaItem[types.AutomationData], error)
	DeleteMedia(ctx context.Context, id int64) error

	// Search operations
	SearchMedia(ctx context.Context, query string, options *types.SearchOptions) ([]models.AutomationMediaItem[types.AutomationData], error)

	// Quality profile methods
	GetQualityProfiles(ctx context.Context) ([]types.QualityProfile, error)
	// Metadata profiles (Lidarr, only I think)
	GetMetadataProfiles(ctx context.Context) ([]types.MetadataProfile, error)

	// Tag methods
	GetTags(ctx context.Context) ([]types.Tag, error)
	CreateTag(ctx context.Context, tag string) (types.Tag, error)

	// Calendar/upcoming releases
	GetCalendar(ctx context.Context, start, end time.Time) ([]models.AutomationMediaItem[types.AutomationData], error)

	// Command execution
	ExecuteCommand(ctx context.Context, command types.Command) (types.CommandResult, error)
}
