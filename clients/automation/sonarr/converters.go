package sonarr

import (
	"time"

	sonarr "github.com/devopsarr/sonarr-go/sonarr"
	"suasor/clients/automation/types"
	"suasor/types/models"
)

func (s *SonarrClient) convertSeriesToMediaItem(series *sonarr.SeriesResource) models.AutomationMediaItem[types.AutomationData] {
	// Convert images
	images := make([]types.AutomationMediaImage, 0, len(series.GetImages()))
	for _, img := range series.GetImages() {
		images = append(images, types.AutomationMediaImage{
			URL:       img.GetRemoteUrl(),
			CoverType: string(img.GetCoverType()),
		})
	}

	// Get quality profile name
	qualityProfile := types.QualityProfileSummary{
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

	return models.AutomationMediaItem[types.AutomationData]{
		ID:               uint64(series.GetId()),
		Title:            series.GetTitle(),
		Overview:         series.GetOverview(),
		Year:             series.GetYear(),
		AddedAt:          series.GetAdded(),
		Status:           types.GetStatusFromSeriesStatus(series.GetStatus()),
		Path:             series.GetPath(),
		QualityProfile:   qualityProfile,
		Images:           images,
		Genres:           genres,
		DownloadedStatus: determineDownloadStatus(series.GetStatistics()),
		Monitored:        series.GetMonitored(),

		Type: types.AUTOMEDIATYPE_SERIES,
		Data: types.AutomationEpisode{
			ReleaseDate: releaseDate,
		},
	}
}
