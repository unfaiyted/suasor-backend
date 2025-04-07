package app

import (
	"suasor/handlers"
)

// -----------------------------
// Job Handlers Implementation
// -----------------------------
type jobHandlersImpl struct {
	jobHandler *handlers.JobHandler
}

func (h *jobHandlersImpl) JobHandler() *handlers.JobHandler {
	return h.jobHandler
}