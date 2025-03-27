package lidarr

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"suasor/client/automation/types"
	"suasor/types/models"
	"suasor/utils"
)

func (l *LidarrClient) GetCalendar(ctx context.Context, start, end time.Time) ([]models.AutomationMediaItem[types.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Time("startDate", start).
		Time("endDate", end).
		Msg("Retrieving calendar from Lidarr")

	calendar, resp, err := l.client.CalendarAPI.ListCalendar(ctx).
		Start(start).
		End(end).
		IncludeArtist(true).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Time("startDate", start).
			Time("endDate", end).
			Msg("Failed to fetch calendar from Lidarr")
		return nil, fmt.Errorf("failed to fetch calendar: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("itemCount", len(calendar)).
		Msg("Successfully retrieved calendar from Lidarr")

	// Convert to our internal representation
	result := make([]models.AutomationMediaItem[types.AutomationData], 0, len(calendar))
	for _, item := range calendar {

		downloadStatus := types.DOWNLOADEDSTATUS_NONE
		if *item.GetStatistics().TrackFileCount >= *item.GetStatistics().TotalTrackCount {
			downloadStatus = types.DOWNLOADEDSTATUS_COMPLETE
		} else if *item.GetStatistics().TrackFileCount > 0 {
			downloadStatus = types.DOWNLOADEDSTATUS_PARTIAL
		}

		artistID := strconv.Itoa(int(item.GetArtistId()))
		artist := item.GetArtist()

		albumData := types.AutomationAlbum{
			ArtistName:  artist.GetArtistName(),
			ArtistID:    artistID,
			ReleaseDate: item.GetReleaseDate(),
		}
		// Get base album details
		albumInfo := models.AutomationMediaItem[types.AutomationData]{
			ID:               uint64(item.GetId()),
			ClientID:         l.ClientID,
			ClientType:       l.ClientType,
			Title:            item.GetTitle(), // Album title
			Overview:         item.GetOverview(),
			Year:             int32(item.GetReleaseDate().Year()),
			Monitored:        item.GetMonitored(),
			ExternalIDs:      l.getAlbumExternalIDs(&item),
			DownloadedStatus: downloadStatus,
		}

		albumInfo.SetData(albumData, types.AUTOMEDIATYPE_ALBUM)

		result = append(result, albumInfo)
	}

	return result, nil
}
