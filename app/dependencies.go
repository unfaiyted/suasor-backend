// app/dependencies.go
package app

import (
	"gorm.io/gorm"
	"suasor/client"
)

type AppDependencies struct {
	// Database
	db *gorm.DB

	// Repositories
	SystemRepositories
	UserRepositories
	MediaItemRepositories
	ClientRepositories
	JobRepositories

	// Services
	UserServices
	SystemServices
	ClientServices
	ClientMediaServices
	MediaItemServices
	MediaServices
	JobServices

	// Factories
	ClientFactoryService *client.ClientFactoryService

	// Handlers
	ClientHandlers
	ClientMediaHandlers
	MediaItemHandlers
	AIHandlers
	UserHandlers
	SystemHandlers
	JobHandlers
}

// GetDB returns the database connection
func (a *AppDependencies) GetDB() *gorm.DB {
	return a.db
}
