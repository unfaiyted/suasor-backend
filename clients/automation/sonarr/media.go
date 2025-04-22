package sonarr

import (
	"context"
	"fmt"
	"strconv"

	sonarr "github.com/devopsarr/sonarr-go/sonarr"
	"suasor/utils/logger"

	"suasor/clients/automation/types"
	"suasor/types/models"
	"suasor/types/requests"
)

func (s *SonarrClient) GetLibraryItems(ctx context.Context, options *types.LibraryQueryOptions) ([]models.AutomationMediaItem[types.AutomationData], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("baseURL", s.config.BaseURL).
		Msg("Retrieving library items from Sonarr server")

	// Call the Sonarr API
	log.Debug().Msg("Making API request to Sonarr server for series library")

	seriesResult, resp, err := s.client.SeriesAPI.ListSeries(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", s.config.BaseURL).
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
	mediaItems := make([]models.AutomationMediaItem[types.AutomationData], 0, len(pagedSeries))
	for _, series := range pagedSeries {
		mediaItem := s.convertSeriesToMediaItem(&series)
		mediaItems = append(mediaItems, mediaItem)
	}

	log.Info().
		Int("itemsReturned", len(mediaItems)).
		Msg("Completed GetLibraryItems request")

	return mediaItems, nil
}

func (s *SonarrClient) GetMediaByID(ctx context.Context, id int64) (models.AutomationMediaItem[types.AutomationData], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Int64("seriesID", id).
		Str("baseURL", s.config.BaseURL).
		Msg("Retrieving specific series from Sonarr server")

	// Call the Sonarr API
	log.Debug().
		Int64("seriesID", id).
		Msg("Making API request to Sonarr server")

	series, resp, err := s.client.SeriesAPI.GetSeriesById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", s.config.BaseURL).
			Str("apiEndpoint", fmt.Sprintf("/series/%d", id)).
			Int("statusCode", 0).
			Msg("Failed to fetch series from Sonarr")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to fetch series: %w", err)
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

func (s *SonarrClient) AddMedia(ctx context.Context, item requests.AutomationMediaAddRequest) (models.AutomationMediaItem[types.AutomationData], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
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
			Str("baseURL", s.config.BaseURL).
			Str("title", item.Title).
			Msg("Failed to add series to Sonarr")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to add series: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("seriesID", result.GetId()).
		Str("title", result.GetTitle()).
		Msg("Successfully added series to Sonarr")

	return s.convertSeriesToMediaItem(result), nil
}

func (s *SonarrClient) UpdateMedia(ctx context.Context, id int64, item requests.AutomationMediaUpdateRequest) (models.AutomationMediaItem[types.AutomationData], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
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
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to fetch series for update: %w", err)
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
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to update series: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("seriesID", updatedSeries.GetId()).
		Msg("Successfully updated series in Sonarr")

	return s.convertSeriesToMediaItem(updatedSeries), nil
}

func (s *SonarrClient) DeleteMedia(ctx context.Context, id int64) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
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

func (s *SonarrClient) SearchMedia(ctx context.Context, query string, options *types.SearchOptions) ([]models.AutomationMediaItem[types.AutomationData], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
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
	mediaItems := make([]models.AutomationMediaItem[types.AutomationData], 0, len(searchResult))
	for _, series := range searchResult {
		mediaItem := s.convertSeriesToMediaItem(&series)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}
