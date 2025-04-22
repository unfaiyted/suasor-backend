package bundles

import (
	"suasor/repository"
)

type UserRepositories interface {
	UserRepo() repository.UserRepository
	UserConfigRepo() repository.UserConfigRepository
	SessionRepo() repository.SessionRepository
}
