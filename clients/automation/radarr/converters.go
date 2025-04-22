package radarr

import (
	radarr "github.com/devopsarr/radarr-go/radarr"
	"suasor/clients/automation/types"
	"suasor/types/models"
)

func (r *RadarrClient) convertMovieToMediaItem(movie *radarr.MovieResource) models.AutomationMediaItem[types.AutomationData] {
	// Convert images
	images := make([]types.AutomationMediaImage, 0, len(movie.GetImages()))
	for _, img := range movie.GetImages() {
		images = append(images, types.AutomationMediaImage{
			URL:       img.GetRemoteUrl(),
			CoverType: string(img.GetCoverType()),
		})
	}

	// Get quality profile name
	qualityProfile := types.QualityProfileSummary{
		ID:   int64(movie.GetQualityProfileId()),
		Name: "", // We don't have the name in the movie object
	}

	status := types.DOWNLOADEDSTATUS_NONE
	if movie.GetHasFile() {
		status = types.DOWNLOADEDSTATUS_COMPLETE
	}

	return models.AutomationMediaItem[types.AutomationData]{
		ID:       uint64(movie.GetId()),
		Title:    movie.GetTitle(),
		Overview: movie.GetOverview(),

		AddedAt:          movie.GetAdded(),
		Status:           types.GetStatusFromMovieStatus(movie.GetStatus()),
		Path:             movie.GetPath(),
		QualityProfile:   qualityProfile,
		Images:           images,
		DownloadedStatus: status,
		Monitored:        movie.GetMonitored(),

		Type: types.AUTOMEDIATYPE_MOVIE,
		Data: types.AutomationMovie{
			ReleaseDate: movie.GetReleaseDate(),
			Year:        movie.GetYear(),
		},
	}
}
