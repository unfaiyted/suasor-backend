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

	// Services
	UserServices
	SystemServices
	ClientServices
	ClientMediaServices
	MediaItemServices

	// Factories
	ClientFactoryService *client.ClientFactoryService

	// Handlers
	ClientHandlers
	ClientMediaHandlers
	MediaItemHandlers
	UserHandlers
	SystemHandlers
}
