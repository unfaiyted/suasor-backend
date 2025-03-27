package responses

import (
	"suasor/types"
)

// ConfigResponse represents the response structure for configuration endpoints
// @Description Configuration response wrapper
type ConfigResponse struct {
	Data  *types.Configuration `json:"data,omitempty"`
	Error string               `json:"error,omitempty"`
}
