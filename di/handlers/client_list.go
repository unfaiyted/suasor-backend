// app/di/handlers/client_list.go
package handlers

import (
	"context"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/services"
)

// RegisterClientListHandlers registers handlers for client lists
func RegisterClientListHandlers(ctx context.Context, c *container.Container) {
	// Register handlers for all client types and list types
	registerClientListHandler[*clienttypes.EmbyConfig, *mediatypes.Playlist](c)
	registerClientListHandler[*clienttypes.JellyfinConfig, *mediatypes.Playlist](c)
	registerClientListHandler[*clienttypes.PlexConfig, *mediatypes.Playlist](c)
	registerClientListHandler[*clienttypes.SubsonicConfig, *mediatypes.Playlist](c)

	registerClientListHandler[*clienttypes.EmbyConfig, *mediatypes.Collection](c)
	registerClientListHandler[*clienttypes.JellyfinConfig, *mediatypes.Collection](c)
	registerClientListHandler[*clienttypes.PlexConfig, *mediatypes.Collection](c)
	// Subsonic doesn't support collections
}

func registerClientListHandler[T clienttypes.ClientMediaConfig, U mediatypes.ListData](c *container.Container) {
	container.RegisterFactory[handlers.ClientListHandler[T, U]](c, func(c *container.Container) handlers.ClientListHandler[T, U] {
		coreHandler := container.MustGet[handlers.CoreListHandler[U]](c)
		clientService := container.MustGet[services.ClientListService[T, U]](c)
		return handlers.NewClientListHandler[T, U](coreHandler, clientService)
	})
}