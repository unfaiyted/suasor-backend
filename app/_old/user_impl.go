// app/dependencies.go
package app

//
// import (
// 	mediatypes "suasor/client/media/types"
// 	"suasor/handlers"
// 	"suasor/repository"
// 	"suasor/services"
// )
//
// // Concrete implementation of UserRepositories
// type userRepositoriesImpl struct {
// 	userRepo       repository.UserRepository
// 	userConfigRepo repository.UserConfigRepository
// 	sessionRepo    repository.SessionRepository
// }
//
// func (r *userRepositoriesImpl) UserRepo() repository.UserRepository {
// 	return r.userRepo
// }
//
// func (r *userRepositoriesImpl) UserConfigRepo() repository.UserConfigRepository {
// 	return r.userConfigRepo
// }
//
// func (r *userRepositoriesImpl) SessionRepo() repository.SessionRepository {
// 	return r.sessionRepo
// }
//
// type userServicesImpl struct {
// 	userService       services.UserService
// 	userConfigService services.UserConfigService
// 	authService       services.AuthService
// }
//
// func (s *userServicesImpl) UserService() services.UserService {
// 	return s.userService
// }
//
// func (s *userServicesImpl) UserConfigService() services.UserConfigService {
// 	return s.userConfigService
// }
//
// func (s *userServicesImpl) AuthService() services.AuthService {
// 	return s.authService
// }
//
// type userHandlersImpl struct {
// 	authHandler       *handlers.AuthHandler
// 	userHandler       *handlers.UserHandler
// 	userConfigHandler *handlers.UserConfigHandler
// }
//
// func (h *userHandlersImpl) AuthHandler() *handlers.AuthHandler {
// 	return h.authHandler
// }
//
// func (h *userHandlersImpl) UserHandler() *handlers.UserHandler {
// 	return h.userHandler
// }
//
// func (h *userHandlersImpl) UserConfigHandler() *handlers.UserConfigHandler {
// 	return h.userConfigHandler
// }
//
// type userMediaItemDataServicesImpl struct {
// 	movieDataService      services.UserMediaItemDataService[*mediatypes.Movie]
// 	seriesDataService     services.UserMediaItemDataService[*mediatypes.Series]
// 	episodeDataService    services.UserMediaItemDataService[*mediatypes.Episode]
// 	trackDataService      services.UserMediaItemDataService[*mediatypes.Track]
// 	albumDataService      services.UserMediaItemDataService[*mediatypes.Album]
// 	artistDataService     services.UserMediaItemDataService[*mediatypes.Artist]
// 	collectionDataService services.UserMediaItemDataService[*mediatypes.Collection]
// 	playlistDataService   services.UserMediaItemDataService[*mediatypes.Playlist]
// }
//
// func (s *userMediaItemDataServicesImpl) MovieDataService() services.UserMediaItemDataService[*mediatypes.Movie] {
// 	return s.movieDataService
// }
//
// func (s *userMediaItemDataServicesImpl) SeriesDataService() services.UserMediaItemDataService[*mediatypes.Series] {
// 	return s.seriesDataService
// }
//
// func (s *userMediaItemDataServicesImpl) EpisodeDataService() services.UserMediaItemDataService[*mediatypes.Episode] {
// 	return s.episodeDataService
// }
//
// func (s *userMediaItemDataServicesImpl) TrackDataService() services.UserMediaItemDataService[*mediatypes.Track] {
// 	return s.trackDataService
// }
//
// func (s *userMediaItemDataServicesImpl) AlbumDataService() services.UserMediaItemDataService[*mediatypes.Album] {
// 	return s.albumDataService
// }
//
// func (s *userMediaItemDataServicesImpl) ArtistDataService() services.UserMediaItemDataService[*mediatypes.Artist] {
// 	return s.artistDataService
// }
//
// func (s *userMediaItemDataServicesImpl) CollectionDataService() services.UserMediaItemDataService[*mediatypes.Collection] {
// 	return s.collectionDataService
// }
//
// func (s *userMediaItemDataServicesImpl) PlaylistDataService() services.UserMediaItemDataService[*mediatypes.Playlist] {
// 	return s.playlistDataService
// }
