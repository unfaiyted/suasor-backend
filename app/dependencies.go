// app/dependencies.go
package app

import (
	"suasor/client"
)

type AppDependencies struct {
	// Repositories
	AppRepository
	UserRepositories
	MediaItemRepositories
	ClientRepositories

	// Services
	AppServices
	UserServices
	SystemServices
	ClientServices
	ClientMediaServices
	MediaItemServices

	// Factories
	ClientFactoryService *client.ClientFactoryService

	// Handlers
	ClientHandlers
	MediaHandlers
	MediaItemHandlers
	UserHandlers
	SystemHandlers
}
