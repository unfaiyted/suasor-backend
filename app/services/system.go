package services

import (
	"suasor/services"
)

type SystemServices interface {
	HealthService() services.HealthService
	ConfigService() services.ConfigService
}
