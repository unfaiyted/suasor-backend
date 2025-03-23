package sonarr

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sonarr "github.com/devopsarr/sonarr-go/sonarr"
	"suasor/client/automation/interfaces"
	"suasor/utils"
)

// Configuration holds Sonarr connection settings
type Configuration struct {
	BaseURL string
	APIKey  string
}

// SonarrClient implements the AutomationProvider interface
type SonarrClient struct {
	interfaces.BaseAutomationTool
	client *sonarr.APIClient
	config Configuration
}

func DetermineDownloadStatus(stats sonarr.SeriesStatisticsResource) interfaces.DownloadedStatus {

	allRequestedDownloaded := stats.GetEpisodeFileCount() == stats.GetEpisodeCount()
	allEpisodesDownloaded := stats.GetEpisodeFileCount() == stats.GetTotalEpisodeCount()

	downloadStatus := interfaces.DOWNLOADEDSTATUS_NONE
	if allEpisodesDownloaded {
		downloadStatus = interfaces.DOWNLOADEDSTATUS_COMPLETE
	}
	if allRequestedDownloaded && !allEpisodesDownloaded {
		downloadStatus = interfaces.DOWNLOADEDSTATUS_REQUESTED
	}

	return downloadStatus
}

// NewSonarrClient creates a new Sonarr client instance
func NewSonarrClient(ctx context.Context, clientID uint32, config any) (interfaces.AutomationProvider, error) {
	// Extract config
	cfg, ok := config.(Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Sonarr client")
	}

	// Create API client configuration
	apiConfig := sonarr.NewConfiguration()
	apiConfig.AddDefaultHeader("X-Api-Key", cfg.APIKey)
	apiConfig.Servers = sonarr.ServerConfigurations{
		{
			URL: cfg.BaseURL,
		},
	}

	client := sonarr.NewAPIClient(apiConfig)

	sonarrClient := &SonarrClient{
		BaseAutomationTool: interfaces.BaseAutomationTool{
			ClientID:   clientID,
			ClientType: interfaces.ClientTypeSonarr,
			URL:        cfg.BaseURL,
			APIKey:     cfg.APIKey,
		},
		client: client,
		config: cfg,
	}

	return sonarrClient, nil
}

// Register the provider factory
func init() {
	interfaces.RegisterAutomationProvider(interfaces.ClientTypeSonarr, NewSonarrClient)
}

// Capability methods
func (s *SonarrClient) SupportsMovies() bool  { return false }
func (s *SonarrClient) SupportsTVShows() bool { return true }
func (s *SonarrClient) SupportsMusic() bool   { return false }

// GetSystemStatus retrieves system information from Sonarr
func (s *SonarrClient) GetSystemStatus(ctx context.Context) (interfaces.SystemStatus, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("baseURL", s.URL).
		Msg("Retrieving system status from Sonarr server")

	// Call the Sonarr API
	log.Debug().Msg("Making API request to Sonarr server for system status")

	statusResult, resp, err := s.client.SystemAPI.GetSystemStatus(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", s.URL).
			Str("apiEndpoint", "/system/status").
			Int("statusCode", 0).
			Msg("Failed to fetch system status from Sonarr")
		return interfaces.SystemStatus{}, fmt.Errorf("failed to fetch system status: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("version", statusResult.GetVersion()).
		Msg("Successfully retrieved system status from Sonarr")

	// Convert to our internal type
	status := interfaces.SystemStatus{
		Version:     statusResult.GetVersion(),
		StartupPath: statusResult.GetStartupPath(),
		AppData:     statusResult.GetAppData(),
		OsName:      statusResult.GetOsName(),
		Branch:      statusResult.GetBranch(),
	}

	return status, nil
}

// GetLibraryItems retrieves all series from Sonarr
func (s *SonarrClient) GetLibraryItems(ctx context.Context, options *interfaces.LibraryQueryOptions) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("baseURL", s.URL).
		Msg("Retrieving library items from Sonarr server")

	// Call the Sonarr API
	log.Debug().Msg("Making API request to Sonarr server for series library")

	seriesResult, resp, err := s.client.SeriesAPI.ListSeries(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", s.URL).
			Str("apiEndpoint", "/series").
			Int("statusCode", 0).
			Msg("Failed to fetch series from Sonarr")
		return nil, fmt.Errorf("failed to fetch series: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("seriesCount", len(seriesResult)).
		Msg("Successfully retrieved series from Sonarr")

	// Apply paging if options provided
	var start, end int
	if options != nil {
		if options.Offset > 0 {
			start = options.Offset
		}
		if options.Limit > 0 {
			end = start + options.Limit
			if end > len(seriesResult) {
				end = len(seriesResult)
			}
		} else {
			end = len(seriesResult)
		}
	} else {
		end = len(seriesResult)
	}

	// Ensure valid slice bounds
	if start >= len(seriesResult) {
		start = 0
		end = 0
	}

	// Apply paging
	var pagedSeries []sonarr.SeriesResource
	if start < end {
		pagedSeries = seriesResult[start:end]
	} else {
		pagedSeries = []sonarr.SeriesResource{}
	}

	// Convert to our internal type
	mediaItems := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(pagedSeries))
	for _, series := range pagedSeries {
		mediaItem := s.convertSeriesToMediaItem(&series)
		mediaItems = append(mediaItems, mediaItem)
	}

	log.Info().
		Int("itemsReturned", len(mediaItems)).
		Msg("Completed GetLibraryItems request")

	return mediaItems, nil
}

// Helper function to convert Sonarr series to generic MediaItem
func (s *SonarrClient) convertSeriesToMediaItem(series *sonarr.SeriesResource) interfaces.AutomationMediaItem[interfaces.AutomationData] {
	// Convert images
	images := make([]interfaces.AutomationMediaImage, 0, len(series.GetImages()))
	for _, img := range series.GetImages() {
		images = append(images, interfaces.AutomationMediaImage{
			URL:       img.GetRemoteUrl(),
			CoverType: string(img.GetCoverType()),
		})
	}

	// Get quality profile name
	qualityProfile := interfaces.QualityProfileSummary{
		ID:   int64(series.GetQualityProfileId()),
		Name: "", // We don't have the name in the series object
	}

	// Convert genres
	genres := series.GetGenres()

	// First aired date as release date if available
	var releaseDate time.Time
	if series.FirstAired.IsSet() {
		releaseDate = *series.FirstAired.Get()
	}

	return interfaces.AutomationMediaItem[interfaces.AutomationData]{
		ID:               uint64(series.GetId()),
		Title:            series.GetTitle(),
		Overview:         series.GetOverview(),
		MediaType:        interfaces.AUTOMEDIATYPE_SERIES,
		Year:             series.GetYear(),
		AddedAt:          series.GetAdded(),
		Status:           interfaces.GetStatusFromSeriesStatus(series.GetStatus()),
		Path:             series.GetPath(),
		QualityProfile:   qualityProfile,
		Images:           images,
		Genres:           genres,
		DownloadedStatus: DetermineDownloadStatus(series.GetStatistics()),
		Monitored:        series.GetMonitored(),
		Data: interfaces.AutomationEpisode{
			ReleaseDate: releaseDate,
		},
	}
}

// GetMediaByID retrieves a specific series by ID
func (s *SonarrClient) GetMediaByID(ctx context.Context, id int64) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Int64("seriesID", id).
		Str("baseURL", s.URL).
		Msg("Retrieving specific series from Sonarr server")

	// Call the Sonarr API
	log.Debug().
		Int64("seriesID", id).
		Msg("Making API request to Sonarr server")

	series, resp, err := s.client.SeriesAPI.GetSeriesById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", s.URL).
			Str("apiEndpoint", fmt.Sprintf("/series/%d", id)).
			Int("statusCode", 0).
			Msg("Failed to fetch series from Sonarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to fetch series: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int64("seriesID", id).
		Str("seriesTitle", series.GetTitle()).
		Msg("Successfully retrieved series from Sonarr")

	// Convert to our internal type
	mediaItem := s.convertSeriesToMediaItem(series)

	log.Debug().
		Int64("seriesID", id).
		Str("seriesTitle", mediaItem.Title).
		Int32("year", mediaItem.Year).
		Msg("Successfully returned series data")

	return mediaItem, nil
}

// AddMedia adds a new series to Sonarr
func (s *SonarrClient) AddMedia(ctx context.Context, item interfaces.AutomationMediaAddRequest) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("title", item.Title).
		Msg("Adding series to Sonarr")

	// Create new series resource
	newSeries := sonarr.NewSeriesResource()
	newSeries.SetTitle(item.Title)
	newSeries.SetQualityProfileId(int32(item.QualityProfileID))
	newSeries.SetTvdbId(int32(item.TVDBID))
	newSeries.SetYear(int32(item.Year))
	newSeries.SetMonitored(item.Monitored)
	newSeries.SetRootFolderPath(item.Path)
	newSeries.SetTags(item.Tags)

	// Set series type (standard, anime, daily)
	newSeries.SetSeriesType("standard")

	// Set add options
	options := sonarr.NewAddSeriesOptions()
	options.SetSearchForMissingEpisodes(item.SearchForMedia)
	newSeries.SetAddOptions(*options)

	// Make API request
	result, resp, err := s.client.SeriesAPI.CreateSeries(ctx).SeriesResource(*newSeries).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", s.URL).
			Str("title", item.Title).
			Msg("Failed to add series to Sonarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to add series: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("seriesID", result.GetId()).
		Str("title", result.GetTitle()).
		Msg("Successfully added series to Sonarr")

	return s.convertSeriesToMediaItem(result), nil
}

// UpdateMedia updates an existing series in Sonarr
func (s *SonarrClient) UpdateMedia(ctx context.Context, id int64, item interfaces.AutomationMediaUpdateRequest) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Int64("seriesID", id).
		Msg("Updating series in Sonarr")

	// First get the existing series
	existingSeries, resp, err := s.client.SeriesAPI.GetSeriesById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("seriesID", id).
			Msg("Failed to fetch series for update")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to fetch series for update: %w", err)
	}

	// Update fields as needed
	existingSeries.SetMonitored(item.Monitored)

	if item.QualityProfileID > 0 {
		existingSeries.SetQualityProfileId(int32(item.QualityProfileID))
	}

	if item.Path != "" {
		existingSeries.SetPath(item.Path)
	}

	if item.Tags != nil {
		existingSeries.SetTags(convertInt64SliceToInt32(item.Tags))
	}

	stringId := strconv.FormatInt(id, 10)

	// Send update request
	updatedSeries, resp, err := s.client.SeriesAPI.UpdateSeries(ctx, stringId).SeriesResource(*existingSeries).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("seriesID", id).
			Msg("Failed to update series in Sonarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to update series: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("seriesID", updatedSeries.GetId()).
		Msg("Successfully updated series in Sonarr")

	return s.convertSeriesToMediaItem(updatedSeries), nil
}

// DeleteMedia removes a series from Sonarr
func (s *SonarrClient) DeleteMedia(ctx context.Context, id int64) error {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Int64("seriesID", id).
		Msg("Deleting series from Sonarr")

	// Optional deletion flags
	deleteFiles := false
	addExclusion := false

	resp, err := s.client.SeriesAPI.DeleteSeries(ctx, int32(id)).
		DeleteFiles(deleteFiles).
		AddImportListExclusion(addExclusion).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Int64("seriesID", id).
			Msg("Failed to delete series from Sonarr")
		return fmt.Errorf("failed to delete series: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int64("seriesID", id).
		Msg("Successfully deleted series from Sonarr")

	return nil
}

// SearchMedia searches for series in Sonarr
func (s *SonarrClient) SearchMedia(ctx context.Context, query string, options *interfaces.SearchOptions) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("query", query).
		Msg("Searching for series in Sonarr")

	searchResult, resp, err := s.client.SeriesLookupAPI.ListSeriesLookup(ctx).Term(query).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("query", query).
			Msg("Failed to search for series in Sonarr")
		return nil, fmt.Errorf("failed to search for series: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("resultCount", len(searchResult)).
		Msg("Successfully searched for series in Sonarr")

	// Convert results to MediaItems
	mediaItems := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(searchResult))
	for _, series := range searchResult {
		mediaItem := s.convertSeriesToMediaItem(&series)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}

// GetQualityProfiles retrieves available quality profiles from Sonarr
func (s *SonarrClient) GetQualityProfiles(ctx context.Context) ([]interfaces.QualityProfile, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Msg("Retrieving quality profiles from Sonarr")

	profiles, resp, err := s.client.QualityProfileAPI.ListQualityProfile(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch quality profiles from Sonarr")
		return nil, fmt.Errorf("failed to fetch quality profiles: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("profileCount", len(profiles)).
		Msg("Successfully retrieved quality profiles from Sonarr")

	// Convert to our internal representation
	result := make([]interfaces.QualityProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, interfaces.QualityProfile{
			ID:   int64(profile.GetId()),
			Name: profile.GetName(),
		})
	}

	return result, nil
}

// GetTags retrieves all tags from Sonarr
func (s *SonarrClient) GetTags(ctx context.Context) ([]interfaces.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Msg("Retrieving tags from Sonarr")

	tags, resp, err := s.client.TagAPI.ListTag(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch tags from Sonarr")
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("tagCount", len(tags)).
		Msg("Successfully retrieved tags from Sonarr")

	// Convert to our internal representation
	result := make([]interfaces.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, interfaces.Tag{
			ID:   int64(tag.GetId()),
			Name: tag.GetLabel(),
		})
	}

	return result, nil
}

// CreateTag creates a new tag in Sonarr
func (s *SonarrClient) CreateTag(ctx context.Context, tagName string) (interfaces.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("tagName", tagName).
		Msg("Creating new tag in Sonarr")

	newTag := sonarr.NewTagResource()
	newTag.SetLabel(tagName)

	createdTag, resp, err := s.client.TagAPI.CreateTag(ctx).TagResource(*newTag).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("tagName", tagName).
			Msg("Failed to create tag in Sonarr")
		return interfaces.Tag{}, fmt.Errorf("failed to create tag: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("tagID", createdTag.GetId()).
		Str("tagName", createdTag.GetLabel()).
		Msg("Successfully created tag in Sonarr")

	return interfaces.Tag{
		ID:   int64(createdTag.GetId()),
		Name: createdTag.GetLabel(),
	}, nil
}

// GetCalendar retrieves upcoming releases from Sonarr
func (s *SonarrClient) GetCalendar(ctx context.Context, start, end time.Time) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Time("startDate", start).
		Time("endDate", end).
		Msg("Retrieving calendar from Sonarr")

	calendar, resp, err := s.client.CalendarAPI.ListCalendar(ctx).
		Start(start).
		End(end).
		IncludeSeries(true).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Time("startDate", start).
			Time("endDate", end).
			Msg("Failed to fetch calendar from Sonarr")
		return nil, fmt.Errorf("failed to fetch calendar: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("itemCount", len(calendar)).
		Msg("Successfully retrieved calendar from Sonarr")

	// Convert to our internal representation
	result := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(calendar))
	for _, item := range calendar {
		// Get base series details
		seriesInfo := interfaces.AutomationMediaItem[interfaces.AutomationData]{
			ID:               uint64(item.GetSeriesId()),
			ClientID:         s.ClientID,
			ClientType:       s.ClientType,
			Title:            *item.GetSeries().Title.Get(),
			MediaType:        "episode",
			Status:           interfaces.GetStatusFromSeriesStatus(*item.GetSeries().Status),
			Overview:         item.GetOverview(),
			Year:             *item.GetSeries().Year,
			Monitored:        *item.GetSeries().Monitored,
			DownloadedStatus: DetermineDownloadStatus(item.Series.GetStatistics()),
			Data: interfaces.AutomationEpisode{
				ReleaseDate: item.GetAirDateUtc(),
			},
		}
		result = append(result, seriesInfo)
	}

	return result, nil
}

// ExecuteCommand executes system commands in Sonarr
func (s *SonarrClient) ExecuteCommand(ctx context.Context, command interfaces.Command) (interfaces.CommandResult, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("commandName", command.Name).
		Msg("Executing command in Sonarr")

	// Create command
	newCommand := sonarr.NewCommandResource()
	newCommand.SetName(command.Name)

	// Execute command
	cmdResult, resp, err := s.client.CommandAPI.CreateCommand(ctx).CommandResource(*newCommand).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("commandName", command.Name).
			Msg("Failed to execute command in Sonarr")
		return interfaces.CommandResult{}, fmt.Errorf("failed to execute command: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("commandId", cmdResult.GetId()).
		Str("commandName", cmdResult.GetName()).
		Str("status", string(cmdResult.GetStatus())).
		Msg("Successfully initiated command in Sonarr")

	return interfaces.CommandResult{
		ID:        int64(cmdResult.GetId()),
		Name:      cmdResult.GetName(),
		Status:    string(cmdResult.GetStatus()),
		StartedAt: cmdResult.GetStarted(),
	}, nil
}

// Helper function to convert []int64 to []int32
func convertInt64SliceToInt32(in []int64) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}

func (r *SonarrClient) GetMetadataProfiles(ctx context.Context) ([]interfaces.MetadataProfile, error) {
	return nil, interfaces.ErrAutomationFeatureNotSupported
}
