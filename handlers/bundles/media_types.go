package bundles

import (
	// clienttypes "suasor/client/types"
	"suasor/handlers"
)

type MediaTypeHandlers interface {
	CoreMediaTypeHandlers
	UserMediaTypeHandlers
}

type CoreMediaTypeHandlers interface {
	MusicCoreHandler() *handlers.CoreMusicHandler
	MovieCoreHandler() *handlers.CoreMovieHandler
	SeriesCoreHandler() *handlers.CoreSeriesHandler
}

type UserMediaTypeHandlers interface {
	MusicUserHandler() *handlers.UserMusicHandler
	MovieUserHandler() *handlers.UserMovieHandler
	SeriesUserHandler() *handlers.UserSeriesHandler
}
