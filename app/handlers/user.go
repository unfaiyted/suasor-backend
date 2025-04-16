package handlers

import (
	"suasor/handlers"
)

type UserHandlers interface {
	AuthHandler() *handlers.AuthHandler
	UserHandler() *handlers.UserHandler
	UserConfigHandler() *handlers.UserConfigHandler
}
