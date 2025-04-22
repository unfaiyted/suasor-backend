package bundles

import (
	"suasor/services"
)

type UserServices interface {
	UserService() services.UserService
	UserConfigService() services.UserConfigService
	AuthService() services.AuthService
}
