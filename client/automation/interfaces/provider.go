package interfaces

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrAutomationFeatureNotSupported = errors.New("feature not supported by this automation tool")

// AutomationProvider defines a common interface for media automation tools (Radarr, Sonarr, Lidarr)
type AutomationProvider interface {
	// Capability methods
	SupportsMovies() bool
	SupportsTVShows() bool
	SupportsMusic() bool

	// System methods
	GetSystemStatus(ctx context.Context) (SystemStatus, error)

	// Library management
	GetLibraryItems(ctx context.Context, options *LibraryQueryOptions) ([]AutomationMediaItem[AutomationData], error)
	GetMediaByID(ctx context.Context, id int64) (AutomationMediaItem[AutomationData], error)
	AddMedia(ctx context.Context, item AutomationMediaAddRequest) (AutomationMediaItem[AutomationData], error)
	UpdateMedia(ctx context.Context, id int64, item AutomationMediaUpdateRequest) (AutomationMediaItem[AutomationData], error)
	DeleteMedia(ctx context.Context, id int64) error

	// Search operations
	SearchMedia(ctx context.Context, query string, options *SearchOptions) ([]AutomationMediaItem[AutomationData], error)

	// Quality profile methods
	GetQualityProfiles(ctx context.Context) ([]QualityProfile, error)
	// Metadata profiles (Lidarr, only I think)
	GetMetadataProfiles(ctx context.Context) ([]MetadataProfile, error)

	// Tag methods
	GetTags(ctx context.Context) ([]Tag, error)
	CreateTag(ctx context.Context, tag string) (Tag, error)

	// Calendar/upcoming releases
	GetCalendar(ctx context.Context, start, end time.Time) ([]AutomationMediaItem[AutomationData], error)

	// Command execution
	ExecuteCommand(ctx context.Context, command Command) (CommandResult, error)
}

// Default "not supported" implementations
func (b *BaseAutomationTool) SupportsMovies() bool  { return false }
func (b *BaseAutomationTool) SupportsTVShows() bool { return false }
func (b *BaseAutomationTool) SupportsMusic() bool   { return false }

// Default implementation for unsupported features
func (b *BaseAutomationTool) GetSystemStatus(ctx context.Context) (SystemStatus, error) {
	return SystemStatus{}, ErrAutomationFeatureNotSupported
}

// Add default implementations for all other methods...

// Provider factory type definition
type AutomationProviderFactory func(ctx context.Context, clientID uint32, config any) (AutomationProvider, error)

// Registry to store provider factories
var automationProviderFactories = make(map[AutomationClientType]AutomationProviderFactory)

// RegisterAutomationProvider adds a new provider factory to the registry
func RegisterAutomationProvider(clientType AutomationClientType, factory AutomationProviderFactory) {
	automationProviderFactories[clientType] = factory
}

// NewAutomationProvider creates providers using the registry
func NewAutomationProvider(ctx context.Context, clientID uint32, clientType AutomationClientType, config any) (AutomationProvider, error) {
	factory, exists := automationProviderFactories[clientType]
	if !exists {
		return nil, fmt.Errorf("unsupported automation tool type: %s", clientType)
	}
	return factory(ctx, clientID, config)
}
