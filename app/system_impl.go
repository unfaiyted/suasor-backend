// app/dependencies.go
package app

import (
	mediatypes "suasor/client/media/types"
	"suasor/handlers"
	"suasor/services"
)

type systemHandlersImpl struct {
	configHandler *handlers.ConfigHandler
	healthHandler *handlers.HealthHandler
}

func (h *systemHandlersImpl) ConfigHandler() *handlers.ConfigHandler {
	return h.configHandler
}

func (h *systemHandlersImpl) HealthHandler() *handlers.HealthHandler {
	return h.healthHandler
}

type systemServicesImpl struct {
	healthService services.HealthService
	configService services.ConfigService
}

func (s *systemServicesImpl) HealthService() services.HealthService {
	return s.healthService
}

func (s *systemServicesImpl) ConfigService() services.ConfigService {
	return s.configService
}
