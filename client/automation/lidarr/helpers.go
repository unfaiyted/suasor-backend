package lidarr

import (
	lidarr "github.com/devopsarr/lidarr-go/lidarr"

	"suasor/client/automation/types"
)

func determineDownloadStatus(stats lidarr.ArtistStatisticsResource) types.DownloadedStatus {
	// Check if all values are properly set
	if *stats.TrackFileCount == 0 || *stats.TrackCount == 0 || *stats.TotalTrackCount == 0 {
		return types.DOWNLOADEDSTATUS_NONE
	}

	trackFileCount := stats.GetTrackFileCount()
	trackCount := stats.GetTrackCount()
	totalTrackCount := stats.GetTotalTrackCount()

	// No tracks downloaded
	if trackFileCount == 0 {
		return types.DOWNLOADEDSTATUS_NONE
	}

	// All tracks the artist has ever released are downloaded
	if trackFileCount == totalTrackCount && totalTrackCount > 0 {
		return types.DOWNLOADEDSTATUS_COMPLETE
	}

	// All monitored/requested tracks are downloaded, but not all tracks the artist has ever released
	if trackFileCount == trackCount && trackCount > 0 && trackCount < totalTrackCount {
		return types.DOWNLOADEDSTATUS_REQUESTED
	}

	// Some tracks are downloaded, but not all monitored tracks
	return types.DOWNLOADEDSTATUS_NONE
}

// Helper function to convert []int64 to []int32
func convertInt64SliceToInt32(in []int64) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}
