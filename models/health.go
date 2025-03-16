package models

// HealthResponse contains information about system health
type HealthResponse struct {
	// Overall status of the system
	Status      string `json:"status" example:"up" binding:"required,oneof=up down degraded"`
	Application bool   `json:"application" example:"true" binding:"required"`
	Database    bool   `json:"database" example:"true" binding:"required"`
}
