package app

// import (
// 	"suasor/repository"
// 	"suasor/services"
// 	"suasor/services/jobs"
// 	"suasor/services/jobs/recommendation"
// )
//
// // Job repositories implementation
// type jobRepositoriesImpl struct {
// 	jobRepo repository.JobRepository
// }
//
// func (r *jobRepositoriesImpl) JobRepo() repository.JobRepository {
// 	return r.jobRepo
// }
//
// // Job services implementation
// type jobServicesImpl struct {
// 	jobService          services.JobService
// 	recommendationJob   *recommendation.RecommendationJob
// 	mediaSyncJob        *jobs.MediaSyncJob
// 	watchHistorySyncJob *jobs.WatchHistorySyncJob
// 	favoritesSyncJob    *jobs.FavoritesSyncJob
// }
//
// func (s *jobServicesImpl) JobService() services.JobService {
// 	return s.jobService
// }
//
// func (s *jobServicesImpl) RecommendationJob() *recommendation.RecommendationJob {
// 	return s.recommendationJob
// }
//
// func (s *jobServicesImpl) MediaSyncJob() *jobs.MediaSyncJob {
// 	return s.mediaSyncJob
// }
//
// func (s *jobServicesImpl) WatchHistorySyncJob() *jobs.WatchHistorySyncJob {
// 	return s.watchHistorySyncJob
// }
//
// func (s *jobServicesImpl) FavoritesSyncJob() *jobs.FavoritesSyncJob {
// 	return s.favoritesSyncJob
// }
