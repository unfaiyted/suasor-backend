package app

// import (
// 	"suasor/client"
// 	mediatypes "suasor/client/media/types"
// 	clienttypes "suasor/client/types"
// 	"suasor/handlers"
// 	"suasor/repository"
// 	"suasor/services"
// )

// // NewCoreMusicHandler creates a new CoreMusicHandler with the proper services
// // Since the fields are unexported, we'll need to use the constructor
// func NewCoreMusicHandler(
// 	trackService services.ClientMediaItemService[*mediatypes.Track],
// 	albumService services.ClientMediaItemService[*mediatypes.Album],
// 	artistService services.ClientMediaItemService[*mediatypes.Artist],
// ) *handlers.CoreMusicHandler {
// 	// Create a mock struct for now
// 	// In a real implementation, we'd use the proper constructor with the correct parameters
// 	return &handlers.CoreMusicHandler{}
// }
//
// // CreateClientMediaMovieHandler creates a client-specific movie handler from the base three-pronged service
// func CreateClientMediaMovieHandler[T clienttypes.ClientMediaConfig](
// 	clientService services.ClientService[T],
// 	movieService services.ClientMediaItemService[*mediatypes.Movie],
// ) *handlers.ClientMediaMovieHandler[T] {
// 	// Create a specialized service for this client type
// 	specificMovieService := services.NewClientMediaMovieService[T](
// 		repository.NewClientRepository[T](nil), // We'll replace this with a proper repo later
// 		client.GetClientFactoryService(),
// 	)
//
// 	// Return a new handler using the constructor to respect private fields
// 	return handlers.NewClientMediaMovieHandler[T](specificMovieService)
// }
//
// // CreateClientMediaSeriesHandler creates a client-specific series handler from the base three-pronged service
// func CreateClientMediaSeriesHandler[T clienttypes.ClientMediaConfig](
// 	clientService services.ClientService[T],
// 	seriesService services.ClientMediaItemService[*mediatypes.Series],
// ) *handlers.ClientMediaSeriesHandler[T] {
// 	// Create a specialized service for this client type
// 	specificSeriesService := services.NewClientMediaSeriesService[T](
// 		repository.NewClientRepository[T](nil), // We'll replace this with a proper repo later
// 		client.GetClientFactoryService(),
// 	)
//
// 	// Return a new handler using the constructor to respect private fields
// 	return handlers.NewClientMediaSeriesHandler[T](specificSeriesService)
// }
//
// // CreateClientMediaMusicHandler creates a client-specific music handler from the base three-pronged service
// func CreateClientMediaMusicHandler[T clienttypes.ClientMediaConfig](
// 	clientService services.ClientService[T],
// 	trackService services.ClientMediaItemService[*mediatypes.Track],
// 	albumService services.ClientMediaItemService[*mediatypes.Album],
// 	artistService services.ClientMediaItemService[*mediatypes.Artist],
// ) *handlers.ClientMediaMusicHandler[T] {
// 	// Create a specialized service for this client type
// 	specificMusicService := services.NewClientMediaMusicService[T](
// 		repository.NewClientRepository[T](nil), // We'll replace this with a proper repo later
// 		client.GetClientFactoryService(),
// 	)
//
// 	// Return a new handler using the constructor to respect private fields
// 	return handlers.NewClientMediaMusicHandler[T](specificMusicService)
// }
//
// // For simplicity, we'll use a stub implementation rather than a full wrapper
// // Since we're removing legacy code, we can replace this with a proper implementation later
// func NewClientCollectionService(
// 	coreService services.CoreCollectionService,
// 	clientRepo repository.ClientRepository[clienttypes.ClientMediaConfig],
// 	clientFactory *client.ClientFactoryService,
// ) services.ClientMediaCollectionService {
// 	// In a real implementation, we would implement a proper service
// 	// For now, we'll return nil and handle the nil case in dependent code
// 	return nil
// }
