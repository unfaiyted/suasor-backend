package bundles

import (
	"suasor/services"
	"suasor/services/jobs"
	"suasor/services/jobs/recommendation"
)

type JobServices interface {
	JobService() services.JobService
	RecommendationJob() *recommendation.RecommendationJob
	MediaSyncJob() *jobs.MediaSyncJob
	WatchHistorySyncJob() *jobs.WatchHistorySyncJob
	FavoritesSyncJob() *jobs.FavoritesSyncJob
}
