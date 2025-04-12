// app/dependencies.go
package app

import (
	"gorm.io/gorm"
	"suasor/client"
	"suasor/handlers"
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

	// Collections
	RepositoryCollections

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
	searchHandler *handlers.SearchHandler
}

// GetDB returns the database connection
func (a *AppDependencies) GetDB() *gorm.DB {
	return a.db
}

// SearchHandler returns the search handler
func (a *AppDependencies) SearchHandler() *handlers.SearchHandler {
	return a.searchHandler
}
