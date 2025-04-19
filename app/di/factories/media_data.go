// app/factory/media_data.go
package factories

import (
	"suasor/app/handlers"
	"suasor/app/repository"
	"suasor/app/services"
)

// MediaDataFactory defines the interface for creating repositories and services
type MediaDataFactory interface {

	// (MediaItem) Repositories
	CreateCoreRepositories() repository.CoreMediaItemRepositories
	CreateUserRepositories() repository.UserMediaItemRepositories
	CreateClientRepositories() repository.ClientUserMediaDataRepositories

	// (MediaData) Repositories
	CreateUserDataRepositories() repository.UserMediaDataRepositories
	CreateCoreDataRepositories() repository.CoreUserMediaItemDataRepositories
	CreateClientDataRepositories() repository.ClientUserMediaDataRepositories

	// (MediaItem) Services
	CreateCoreServices(repos repository.CoreMediaItemRepositories) services.CoreMediaItemServices
	CreateUserServices(coreServices services.CoreMediaItemServices, userRepos repository.UserMediaItemRepositories) services.UserMediaItemServices
	CreateClientServices(coreServices services.CoreMediaItemServices, clientRepos repository.ClientMediaItemRepositories) services.ClientMediaItemServices

	// (MediaData) Services
	CreateCoreDataServices(repos repository.CoreMediaItemRepositories) services.CoreUserMediaItemDataServices
	CreateUserDataServices(coreDataServices services.CoreUserMediaItemDataServices, userRepos repository.UserMediaDataRepositories) services.UserMediaItemDataServices
	CreateClientDataServices(userDataServices services.UserMediaItemDataServices, clientRepos repository.ClientUserMediaDataRepositories) services.ClientUserMediaItemDataServices

	// (ListServices) Services
	CreateCoreListServices(coreServices services.CoreMediaItemServices) services.CoreListServices
	CreateUserListServices(userServices services.UserMediaItemServices, coreListServices services.CoreListServices) services.UserListServices
	CreateClientListServices(clientServices services.ClientMediaItemServices, coreListServices services.CoreListServices) services.ClientListServices

	// --- HANDLER FACTORIES --- //

	// (MediaItem) Handlers
	CreateCoreMediaItemHandlers(
		coreServices services.CoreMediaItemServices,
	) handlers.CoreMediaItemHandlers
	CreateUserMediaItemHandlers(
		userServices services.UserMediaItemServices,
		coreHandlers handlers.CoreMediaItemHandlers,
	) handlers.UserMediaItemHandlers
	CreateClientMediaItemHandlers(clientServices services.ClientMediaItemServices,
		dataServices services.UserMediaItemServices,
		userHandlers handlers.UserMediaItemHandlers,
	) handlers.ClientMediaItemHandlers

	// (MediaData) Handlers
	CreateCoreDataHandlers(coreServices services.CoreUserMediaItemDataServices) handlers.CoreMediaItemDataHandlers
	CreateUserDataHandlers(userServices services.UserMediaItemDataServices, coreDataHandlers handlers.CoreMediaItemDataHandlers) handlers.UserMediaItemDataHandlers
	CreateClientDataHandlers(dataServices services.ClientUserMediaItemDataServices, userDataHandlers handlers.UserMediaItemDataHandlers) handlers.ClientMediaItemDataHandlers
}
