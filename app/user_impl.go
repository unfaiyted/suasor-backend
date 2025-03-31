// app/dependencies.go
package app

import (
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"
)

// Concrete implementation of UserRepositories
type userRepositoriesImpl struct {
	userRepo       repository.UserRepository
	userConfigRepo repository.UserConfigRepository
	sessionRepo    repository.SessionRepository
}

func (r *userRepositoriesImpl) UserRepo() repository.UserRepository {
	return r.userRepo
}

func (r *userRepositoriesImpl) UserConfigRepo() repository.UserConfigRepository {
	return r.userConfigRepo
}

func (r *userRepositoriesImpl) SessionRepo() repository.SessionRepository {
	return r.sessionRepo
}

type userServicesImpl struct {
	userService       services.UserService
	userConfigService services.UserConfigService
	authService       services.AuthService
}

func (s *userServicesImpl) UserService() services.UserService {
	return s.userService
}

func (s *userServicesImpl) UserConfigService() services.UserConfigService {
	return s.userConfigService
}

func (s *userServicesImpl) AuthService() services.AuthService {
	return s.authService
}

type userHandlersImpl struct {
	authHandler       *handlers.AuthHandler
	userHandler       *handlers.UserHandler
	userConfigHandler *handlers.UserConfigHandler
}

func (h *userHandlersImpl) AuthHandler() *handlers.AuthHandler {
	return h.authHandler
}

func (h *userHandlersImpl) UserHandler() *handlers.UserHandler {
	return h.userHandler
}

func (h *userHandlersImpl) UserConfigHandler() *handlers.UserConfigHandler {
	return h.userConfigHandler
}
