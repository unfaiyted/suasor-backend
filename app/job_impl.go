package app

import (
	"suasor/repository"
	"suasor/services"
)

// Job repositories implementation
type jobRepositoriesImpl struct {
	jobRepo repository.JobRepository
}

func (r *jobRepositoriesImpl) JobRepo() repository.JobRepository {
	return r.jobRepo
}

// Job services implementation
type jobServicesImpl struct {
	jobService         services.JobService
	recommendationJob  *services.RecommendationJob
	mediaSyncJob       *services.MediaSyncJob
	watchHistorySyncJob *services.WatchHistorySyncJob
	favoritesSyncJob    *services.FavoritesSyncJob
}

func (s *jobServicesImpl) JobService() services.JobService {
	return s.jobService
}

func (s *jobServicesImpl) RecommendationJob() *services.RecommendationJob {
	return s.recommendationJob
}

func (s *jobServicesImpl) MediaSyncJob() *services.MediaSyncJob {
	return s.mediaSyncJob
}

func (s *jobServicesImpl) WatchHistorySyncJob() *services.WatchHistorySyncJob {
	return s.watchHistorySyncJob
}

func (s *jobServicesImpl) FavoritesSyncJob() *services.FavoritesSyncJob {
	return s.favoritesSyncJob
}