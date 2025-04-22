package sonarr

import (
	"context"
	"fmt"
	"time"

	"suasor/utils/logger"

	"suasor/clients/automation/types"
	"suasor/types/models"
)

func (s *SonarrClient) GetCalendar(ctx context.Context, start, end time.Time) ([]models.AutomationMediaItem[types.AutomationData], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
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
	result := make([]models.AutomationMediaItem[types.AutomationData], 0, len(calendar))
	for _, item := range calendar {
		// Get base series details
		seriesInfo := models.AutomationMediaItem[types.AutomationData]{
			ID:               uint64(item.GetSeriesId()),
			ClientID:         s.ClientID,
			ClientType:       s.ClientType,
			Title:            *item.GetSeries().Title.Get(),
			Type:             "episode",
			Status:           types.GetStatusFromSeriesStatus(*item.GetSeries().Status),
			Overview:         item.GetOverview(),
			Year:             *item.GetSeries().Year,
			Monitored:        *item.GetSeries().Monitored,
			DownloadedStatus: determineDownloadStatus(item.Series.GetStatistics()),
			Data: types.AutomationEpisode{
				ReleaseDate: item.GetAirDateUtc(),
			},
		}
		result = append(result, seriesInfo)
	}

	return result, nil
}
