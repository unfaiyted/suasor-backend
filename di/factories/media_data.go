// app/factory/media_data.go
package factories

import (
	clienttypes "suasor/client/types"
	"suasor/handlers"
	handlerbundles "suasor/handlers/bundles"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	svcbundles "suasor/services/bundles"
)

// MediaDataFactory defines the interface for creating repositories and services
type MediaDataFactory interface {

	// (MediaItem) Repositories
	CreateCoreRepositories() repobundles.CoreMediaItemRepositories
	CreateUserRepositories() repobundles.UserMediaItemRepositories
	CreateClientRepositories() repobundles.ClientUserMediaDataRepositories
	CreateClientMediaItemRepositories() repobundles.ClientMediaItemRepositories

	// (MediaData) Repositories
	CreateUserDataRepositories() repobundles.UserMediaDataRepositories
	CreateCoreDataRepositories() repobundles.CoreUserMediaItemDataRepositories
	CreateClientDataRepositories() repobundles.ClientUserMediaDataRepositories

	// (MediaItem) Services
	CreateCoreServices(repos repobundles.CoreMediaItemRepositories) svcbundles.CoreMediaItemServices
	CreateUserServices(coreServices svcbundles.CoreMediaItemServices, userRepos repobundles.UserMediaItemRepositories) svcbundles.UserMediaItemServices
	CreateClientServices(coreServices svcbundles.CoreMediaItemServices, clientRepos repository.ClientRepository[clienttypes.ClientMediaConfig], itemRepos repobundles.ClientMediaItemRepositories) svcbundles.ClientMediaItemServices[clienttypes.ClientMediaConfig]

	// (MediaData) Services
	CreateCoreDataServices(repos repobundles.CoreMediaItemRepositories) svcbundles.CoreUserMediaItemDataServices
	CreateUserDataServices(coreDataServices svcbundles.CoreUserMediaItemDataServices, userRepos repobundles.UserMediaDataRepositories) svcbundles.UserMediaItemDataServices
	CreateClientDataServices(userDataServices svcbundles.UserMediaItemDataServices, clientRepos repobundles.ClientUserMediaDataRepositories) svcbundles.ClientUserMediaItemDataServices

	// (ListServices) Services
	CreateCoreListServices(coreServices svcbundles.CoreMediaItemServices) svcbundles.CoreListServices
	CreateUserListServices(userServices svcbundles.UserMediaItemServices, coreListServices svcbundles.CoreListServices) svcbundles.UserListServices
	CreateClientListServices(clientServices svcbundles.ClientMediaItemServices[clienttypes.ClientMediaConfig], coreListServices svcbundles.CoreListServices) svcbundles.ClientListServices

	// --- HANDLER FACTORIES --- //

	// (MediaItem) Handlers
	CreateCoreMediaItemHandlers(
		coreServices svcbundles.CoreMediaItemServices,
	) handlerbundles.CoreMediaItemHandlers
	CreateUserMediaItemHandlers(
		userServices svcbundles.UserMediaItemServices,
		coreHandlers handlers.CoreMediaItemHandlers,
	) handlers.UserMediaItemHandlers
	CreateClientMediaItemHandlers(clientServices svcbundles.ClientMediaItemServices[clienttypes.ClientMediaConfig],
		dataServices svcbundles.UserMediaItemServices,
		userHandlers handlers.UserMediaItemHandlers,
	) handlers.ClientMediaItemHandlers[clienttypes.ClientMediaConfig]

	// (MediaData) Handlers
	CreateCoreDataHandlers(coreServices svcbundles.CoreUserMediaItemDataServices) handlers.CoreMediaItemDataHandlers
	CreateUserDataHandlers(userServices svcbundles.UserMediaItemDataServices, coreDataHandlers handlers.CoreMediaItemDataHandlers) handlers.UserMediaItemDataHandlers
	CreateClientDataHandlers(dataServices svcbundles.ClientUserMediaItemDataServices, userDataHandlers handlers.UserMediaItemDataHandlers) handlers.ClientMediaItemDataHandlers
}
