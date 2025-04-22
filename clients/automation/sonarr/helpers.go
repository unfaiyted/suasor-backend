package sonarr

import (
	sonarr "github.com/devopsarr/sonarr-go/sonarr"

	"suasor/clients/automation/types"
)

func determineDownloadStatus(stats sonarr.SeriesStatisticsResource) types.DownloadedStatus {

	allRequestedDownloaded := stats.GetEpisodeFileCount() == stats.GetEpisodeCount()
	allEpisodesDownloaded := stats.GetEpisodeFileCount() == stats.GetTotalEpisodeCount()

	downloadStatus := types.DOWNLOADEDSTATUS_NONE
	if allEpisodesDownloaded {
		downloadStatus = types.DOWNLOADEDSTATUS_COMPLETE
	}
	if allRequestedDownloaded && !allEpisodesDownloaded {
		downloadStatus = types.DOWNLOADEDSTATUS_REQUESTED
	}

	return downloadStatus
}

func convertInt64SliceToInt32(in []int64) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}
