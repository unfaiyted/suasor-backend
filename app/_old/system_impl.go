// app/dependencies.go
package app

// import (
// 	"suasor/handlers"
// 	"suasor/repository"
// 	"suasor/services"
// )
//
// type systemRepositoriesImpl struct {
// 	configRepo repository.ConfigRepository
// }
//
// func (r *systemRepositoriesImpl) ConfigRepo() repository.ConfigRepository {
// 	return r.configRepo
// }
//
// type systemHandlersImpl struct {
// 	configHandler  *handlers.ConfigHandler
// 	healthHandler  *handlers.HealthHandler
// 	clientsHandler *handlers.ClientsHandler
// }
//
// func (h *systemHandlersImpl) ConfigHandler() *handlers.ConfigHandler {
// 	return h.configHandler
// }
//
// func (h *systemHandlersImpl) HealthHandler() *handlers.HealthHandler {
// 	return h.healthHandler
// }
//
// func (h *systemHandlersImpl) ClientsHandler() *handlers.ClientsHandler {
// 	return h.clientsHandler
// }
//
// type systemServicesImpl struct {
// 	healthService services.HealthService
// 	configService services.ConfigService
// }
//
// func (s *systemServicesImpl) HealthService() services.HealthService {
// 	return s.healthService
// }
//
// func (s *systemServicesImpl) ConfigService() services.ConfigService {
// 	return s.configService
// }
