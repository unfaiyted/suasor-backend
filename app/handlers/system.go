package handlers

import (
	"suasor/handlers"
)

type SystemHandlers interface {
	ConfigHandler() *handlers.ConfigHandler
	HealthHandler() *handlers.HealthHandler
	ClientsHandler() *handlers.ClientsHandler
}
