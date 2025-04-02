// services/automation_client.go
package services

import (
	"context"
	"errors"
	"suasor/client"
	"suasor/client/automation"
	"suasor/client/automation/providers"
	automationtypes "suasor/client/automation/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/types/requests"
	"time"
)

var ErrAutomationUnsupportedFeature = errors.New("feature not supported by this automation client")

// AutomationClientService defines operations for interacting with automation clients
type AutomationClientService interface {
	GetSystemStatus(ctx context.Context, userID uint64, clientID uint64) (*automationtypes.SystemStatus, error)
	GetLibraryItems(ctx context.Context, userID uint64, clientID uint64, options *automationtypes.LibraryQueryOptions) ([]*models.AutomationMediaItem[automationtypes.AutomationData], error)
	GetMediaByID(ctx context.Context, userID uint64, clientID uint64, mediaID string) (*models.AutomationMediaItem[automationtypes.AutomationData], error)
	AddMedia(ctx context.Context, userID uint64, clientID uint64, req requests.AutomationMediaAddRequest) (*models.AutomationMediaItem[automationtypes.AutomationData], error)
	UpdateMedia(ctx context.Context, userID uint64, clientID uint64, mediaID string, req requests.AutomationMediaUpdateRequest) (*models.AutomationMediaItem[automationtypes.AutomationData], error)
	DeleteMedia(ctx context.Context, userID uint64, clientID uint64, mediaID string) error
	SearchMedia(ctx context.Context, userID uint64, clientID uint64, query string) ([]*models.AutomationMediaItem[automationtypes.AutomationData], error)
	GetQualityProfiles(ctx context.Context, userID uint64, clientID uint64) ([]automationtypes.QualityProfile, error)
	GetMetadataProfiles(ctx context.Context, userID uint64, clientID uint64) ([]automationtypes.MetadataProfile, error)
	GetTags(ctx context.Context, userID uint64, clientID uint64) ([]automationtypes.Tag, error)
	CreateTag(ctx context.Context, userID uint64, clientID uint64, req requests.AutomationCreateTagRequest) (*automationtypes.Tag, error)
	GetCalendar(ctx context.Context, userID uint64, clientID uint64, startDate, endDate time.Time) ([]*models.AutomationCalendarItem[automationtypes.AutomationData], error)
	ExecuteCommand(ctx context.Context, userID uint64, clientID uint64, req requests.AutomationExecuteCommandRequest) (*automationtypes.CommandResult, error)
}

type automationClientService struct {
	clientRepo    repository.ClientRepository[types.AutomationClientConfig]
	clientFactory client.ClientFactoryService
}

// NewAutomationClientService creates a new automation client service
func NewAutomationClientService(
	clientRepo repository.ClientRepository[types.AutomationClientConfig],
	clientFactory client.ClientFactoryService,
) AutomationClientService {
	return &automationClientService{
		clientRepo:    clientRepo,
		clientFactory: clientFactory,
	}
}

// getAutomationClient gets a specific automation client for a user
func (s *automationClientService) getAutomationClient(ctx context.Context, userID, clientID uint64) (automation.AutomationClient, error) {
	clientConfig, err := s.clientRepo.GetByID(ctx, clientID, userID)
	if err != nil {
		return nil, err
	}

	autoClient, err := s.clientFactory.GetClient(ctx, clientID, clientConfig.Config.Data.GetType())
	if err != nil {
		return nil, err
	}

	client, ok := autoClient.(automation.AutomationClient)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}
	return client, nil
}

func (s *automationClientService) GetSystemStatus(ctx context.Context, userID uint64, clientID uint64) (*automationtypes.SystemStatus, error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	systemProvider, ok := client.(providers.SystemProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	return systemProvider.GetSystemStatus(ctx)
}

func (s *automationClientService) GetLibraryItems(ctx context.Context, userID uint64, clientID uint64, options *automationtypes.LibraryQueryOptions) ([]*models.AutomationMediaItem[automationtypes.AutomationData], error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	libraryProvider, ok := client.(providers.LibraryProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	library, err := libraryProvider.GetLibraryItems(ctx, options)
	if err != nil {
		return nil, err
	}
	return library, nil
}

func (s *automationClientService) GetMediaByID(ctx context.Context, userID uint64, clientID uint64, mediaID string) (*models.AutomationMediaItem[automationtypes.AutomationData], error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	mediaProvider, ok := client.(providers.MediaProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	media, err := mediaProvider.GetMediaByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	return &media, nil
}

func (s *automationClientService) AddMedia(ctx context.Context, userID uint64, clientID uint64, req requests.AutomationMediaAddRequest) (*models.AutomationMediaItem[automationtypes.AutomationData], error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	mediaProvider, ok := client.(providers.MediaProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	media, err := mediaProvider.AddMedia(ctx, req)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (s *automationClientService) UpdateMedia(ctx context.Context, userID uint64, clientID uint64, mediaID string, req requests.AutomationMediaUpdateRequest) (*models.AutomationMediaItem[automationtypes.AutomationData], error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	mediaProvider, ok := client.(providers.MediaProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	media, err := mediaProvider.UpdateMedia(ctx, mediaID, req)
	return &media, nil
}

func (s *automationClientService) DeleteMedia(ctx context.Context, userID uint64, clientID uint64, mediaID string) error {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return err
	}

	mediaProvider, ok := client.(providers.MediaProvider)
	if !ok {
		return ErrAutomationUnsupportedFeature
	}

	return mediaProvider.DeleteMedia(ctx, mediaID)
}

func (s *automationClientService) SearchMedia(ctx context.Context, userID uint64, clientID uint64, query string) ([]*models.AutomationMediaItem[automationtypes.AutomationData], error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	searchProvider, ok := client.(providers.SearchProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	options := &automationtypes.SearchOptions{
		Limit:  10,
		Offset: 0,
	}
	search, err := searchProvider.SearchMedia(ctx, query, options)
	if err != nil {
		return nil, err
	}
	return search, nil
}

func (s *automationClientService) GetQualityProfiles(ctx context.Context, userID uint64, clientID uint64) ([]automationtypes.QualityProfile, error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	profileProvider, ok := client.(providers.ProfileProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	return profileProvider.GetQualityProfiles(ctx)
}

func (s *automationClientService) GetMetadataProfiles(ctx context.Context, userID uint64, clientID uint64) ([]automationtypes.MetadataProfile, error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	profileProvider, ok := client.(providers.ProfileProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	return profileProvider.GetMetadataProfiles(ctx)
}

func (s *automationClientService) GetTags(ctx context.Context, userID uint64, clientID uint64) ([]automationtypes.Tag, error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	tagProvider, ok := client.(providers.TagProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	return tagProvider.GetTags(ctx)
}

func (s *automationClientService) CreateTag(ctx context.Context, userID uint64, clientID uint64, req requests.AutomationCreateTagRequest) (*automationtypes.Tag, error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	tagProvider, ok := client.(providers.TagProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	tag, err := tagProvider.CreateTag(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *automationClientService) GetCalendar(ctx context.Context, userID uint64, clientID uint64, startDate, endDate time.Time) ([]*models.AutomationCalendarItem[automationtypes.AutomationData], error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	calendarProvider, ok := client.(providers.CalendarProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	calendar, err := calendarProvider.GetCalendar(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return calendar, nil
}

func (s *automationClientService) ExecuteCommand(ctx context.Context, userID uint64, clientID uint64, req requests.AutomationExecuteCommandRequest) (*automationtypes.CommandResult, error) {
	client, err := s.getAutomationClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	commandProvider, ok := client.(providers.CommandProvider)
	if !ok {
		return nil, ErrAutomationUnsupportedFeature
	}

	command := automationtypes.Command{
		Name:       req.Command,
		Parameters: req.Parameters,
	}

	result, err := commandProvider.ExecuteCommand(ctx, command)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
