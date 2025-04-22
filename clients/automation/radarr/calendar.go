package radarr

import (
	"context"
	"fmt"
	"time"

	"suasor/clients/automation/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (r *RadarrClient) GetCalendar(ctx context.Context, start, end time.Time) ([]models.AutomationMediaItem[types.AutomationData], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Time("startDate", start).
		Time("endDate", end).
		Msg("Retrieving calendar from Radarr")

	// Format dates as required by Radarr API
	// startStr := start.Format(time.RFC3339)
	// endStr := end.Format(time.RFC3339)

	calendar, resp, err := r.client.CalendarAPI.ListCalendar(ctx).
		Start(start).
		End(end).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Time("startDate", start).
			Time("endDate", end).
			Msg("Failed to fetch calendar from Radarr")
		return nil, fmt.Errorf("failed to fetch calendar: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("itemCount", len(calendar)).
		Msg("Successfully retrieved calendar from Radarr")

	// Convert to our internal representation
	result := make([]models.AutomationMediaItem[types.AutomationData], 0, len(calendar))
	for _, item := range calendar {

		status := types.DOWNLOADEDSTATUS_NONE
		if item.GetHasFile() {
			status = types.DOWNLOADEDSTATUS_COMPLETE
		}

		result = append(result, models.AutomationMediaItem[types.AutomationData]{
			ID:               uint64(item.GetId()),
			ClientID:         r.ClientID,
			ClientType:       r.ClientType,
			Title:            item.GetTitle(),
			Type:             types.AUTOMEDIATYPE_MOVIE,
			Status:           types.GetStatusFromMovieStatus(item.GetStatus()),
			Overview:         item.GetOverview(),
			Year:             item.GetYear(),
			Monitored:        item.GetMonitored(),
			DownloadedStatus: status,
			Data: types.AutomationMovie{
				ReleaseDate: item.GetPhysicalRelease(),
			},
		})
	}

	return result, nil
}
