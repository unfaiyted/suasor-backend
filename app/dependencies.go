// app/dependencies.go
package app

import (
	"suasor/client"
)

type AppDependencies struct {
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
