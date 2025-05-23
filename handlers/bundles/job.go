package bundles

import (
	"suasor/handlers"
)

type JobHandlers interface {
	JobHandler() *handlers.JobHandler
	RecommendationHandler() *handlers.RecommendationHandler
}
