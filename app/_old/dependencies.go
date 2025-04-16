// app/dependencies.go
package app

// import (
// 	"gorm.io/gorm"
// 	"suasor/app/container"
// 	"suasor/client"
// 	mediatypes "suasor/client/media/types"
// 	"suasor/handlers"
// 	"suasor/repository"
// )

// // Only define the type that's not already defined elsewhere
// type clientUserDataRepositoriesImpl struct{}
//
// func (r *clientUserDataRepositoriesImpl) MovieDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Movie] {
// 	return nil
// }
//
// func (r *clientUserDataRepositoriesImpl) SeriesDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Series] {
// 	return nil
// }
//
// func (r *clientUserDataRepositoriesImpl) EpisodeDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Episode] {
// 	return nil
// }
//
// func (r *clientUserDataRepositoriesImpl) TrackDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Track] {
// 	return nil
// }
//
// func (r *clientUserDataRepositoriesImpl) AlbumDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Album] {
// 	return nil
// }
//
// func (r *clientUserDataRepositoriesImpl) ArtistDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Artist] {
// 	return nil
// }
//
// func (r *clientUserDataRepositoriesImpl) CollectionDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Collection] {
// 	return nil
// }
//
// func (r *clientUserDataRepositoriesImpl) PlaylistDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Playlist] {
// 	return nil
// }
//
// // AppDependencies contains all application dependencies
// // It uses a clean three-pronged architecture for media data dependencies
// type AppDependencies struct {
// 	container *container.Container
// 	// Database connection
// 	db *gorm.DB
//
// 	// Core infrastructure repositories and services
// 	SystemRepositories
// 	UserRepositories
// 	JobRepositories
//
// 	// Three-pronged architecture for media data
// 	// Repository layer
// 	interfaces.CoreMediaItemRepositories         // Base storage layer
// 	interfaces.CoreUserMediaItemDataRepositories // Core-User Data storage layer
// 	interfaces.UserRepositoryFactories           // User-specific storage layer
// 	interfaces.UserDataFactories                 // User-specific user data storage layer
// 	interfaces.ClientRepositoryFactories         // Client-specific storage layer
// 	interfaces.ClientUserDataRepositories        // Client-specific user data storage layer
//
// 	// Service layer
// 	interfaces.CoreMediaItemServices // Core business logic
// 	CoreUserMediaItemDataServices    // Core-User Data business logic
// 	interfaces.UserMediaItemServices // User-specific business logic
//
// 	interfaces.UserMediaItemDataServices       // User-specific user data logic
// 	interfaces.ClientMediaItemServices         // Client-specific business logic
// 	interfaces.ClientUserMediaItemDataServices // Client-specific user data logic
// 	MediaCollectionServices                    // Collection/playlist specialized services
//
// 	// Handler layer (presentation)
// 	interfaces.CoreMediaItemHandlers    // Core API endpoints
// 	CoreMediaItemDataHandlers           // Core-Data API endpoints
// 	interfaces.UserMediaItemHandlers    // User-specific API endpoints
// 	interfaces.ClientMediaItemHandlers  // Client-specific API endpoints
// 	interfaces.SpecializedMediaHandlers // Domain-specific API endpoints
//
// 	// Repository collections for convenience
// 	RepositoryCollections
//
// 	// Standard services
// 	UserServices
// 	SystemServices
// 	ClientServices
// 	MediaServices
// 	JobServices
//
// 	// Standard handlers
// 	ClientHandlers
// 	ClientMediaHandlers // Client-specific media type handlers
// 	AIHandlers
// 	UserHandlers
// 	SystemHandlers
// 	JobHandlers
// 	SearchHandler *handlers.SearchHandler
//
// 	// Factories
// 	ClientFactoryService *client.ClientFactoryService
// 	MediaDataFactory     interfaces.MediaDataFactory
// }
//
// // GetDB returns the database connection
// func (a *AppDependencies) GetDB() *gorm.DB {
// 	return a.db
// }
//
// // GetSearchHandler returns the search handler
// func (a *AppDependencies) GetSearchHandler() *handlers.SearchHandler {
// 	return a.SearchHandler
// }
//
// // DI container
// func (a *AppDependencies) Container() *container.Container {
// 	return a.container
// }

