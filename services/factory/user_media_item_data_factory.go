// package factory
//
// import (
// 	"suasor/client/media/types"
// 	"suasor/repository"
// 	"suasor/services"
//
// 	"gorm.io/gorm"
// )
//
// // UserMediaItemDataServiceFactory creates services for user media item data
// type UserMediaItemDataServiceFactory struct {
// 	db *gorm.DB
// }
//
// // NewUserMediaItemDataServiceFactory creates a new user media item data service factory
// func NewUserMediaItemDataServiceFactory(db *gorm.DB) *UserMediaItemDataServiceFactory {
// 	return &UserMediaItemDataServiceFactory{db: db}
// }
//
// // CreateCoreService creates a core user media item data service for the given media type
// // This now uses the correct architecture by passing a CoreMediaItemService instead of repository
// func (f *UserMediaItemDataServiceFactory) CreateCoreService[T types.MediaData]() services.CoreUserMediaItemDataService[T] {
// 	// First create the media item repository and service
// 	mediaItemRepo := repository.NewCoreMediaItemRepository[T](f.db)
// 	mediaItemService := services.NewCoreMediaItemService(mediaItemRepo)
//
// 	// Now create the user media item data service using the media item service
// 	return services.NewCoreUserMediaItemDataService(mediaItemService)
// }
//
// // CreateUserService creates a user media item data service for the given media type
// func (f *UserMediaItemDataServiceFactory) CreateUserService[T types.MediaData]() services.UserMediaItemDataService[T] {
// 	coreService := f.CreateCoreService[T]()
// 	repo := repository.NewUserMediaItemDataRepository[T](f.db)
// 	return services.NewUserMediaItemDataService(coreService, repo)
// }
//
// // CreateClientService creates a client user media item data service for the given media type
// func (f *UserMediaItemDataServiceFactory) CreateClientService[T types.MediaData]() services.ClientUserMediaItemDataService[T] {
// 	userService := f.CreateUserService[T]()
// 	repo := repository.NewClientUserMediaItemDataRepository[T](f.db)
// 	return services.NewClientUserMediaItemDataService(userService, repo)
// }
//
// // CreateMovieDataService creates a specialized client service for movie data
// func (f *UserMediaItemDataServiceFactory) CreateMovieDataService() services.ClientUserMediaItemDataService[*types.Movie] {
// 	return f.CreateClientService[*types.Movie]()
// }
//
// // CreateSeriesDataService creates a specialized client service for series data
// func (f *UserMediaItemDataServiceFactory) CreateSeriesDataService() services.ClientUserMediaItemDataService[*types.Series] {
// 	return f.CreateClientService[*types.Series]()
// }
//
// // CreateEpisodeDataService creates a specialized client service for episode data
// func (f *UserMediaItemDataServiceFactory) CreateEpisodeDataService() services.ClientUserMediaItemDataService[*types.Episode] {
// 	return f.CreateClientService[*types.Episode]()
// }
//
// // CreateMusicDataService creates a specialized client service for music data
// func (f *UserMediaItemDataServiceFactory) CreateMusicDataService() services.ClientUserMediaItemDataService[*types.Track] {
// 	return f.CreateClientService[*types.Track]()
// }
//
// // CreateAlbumDataService creates a specialized client service for album data
// func (f *UserMediaItemDataServiceFactory) CreateAlbumDataService() services.ClientUserMediaItemDataService[*types.Album] {
// 	return f.CreateClientService[*types.Album]()
// }
//
// // CreateArtistDataService creates a specialized client service for artist data
// func (f *UserMediaItemDataServiceFactory) CreateArtistDataService() services.ClientUserMediaItemDataService[*types.Artist] {
// 	return f.CreateClientService[*types.Artist]()
// }
//
// // CreateCollectionDataService creates a specialized client service for collection data
// func (f *UserMediaItemDataServiceFactory) CreateCollectionDataService() services.ClientUserMediaItemDataService[*types.Collection] {
// 	return f.CreateClientService[*types.Collection]()
// }
//
// // CreatePlaylistDataService creates a specialized client service for playlist data
// func (f *UserMediaItemDataServiceFactory) CreatePlaylistDataService() services.ClientUserMediaItemDataService[*types.Playlist] {
// 	return f.CreateClientService[*types.Playlist]()
// }
