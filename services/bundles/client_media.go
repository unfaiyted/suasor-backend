package bundles

import (
	mediatypes "suasor/clients/media/types"
	"suasor/clients/types"
	"suasor/services"
)

type ClientMediaServices interface {
	ClientMovieServices
	ClientSeriesServices
	ClientEpisodeServices
	ClientMusicServices
	ClientPlaylistServices
}

type ClientMovieServices interface {
	EmbyMovieService() services.ClientMovieService[*types.EmbyConfig]
	JellyfinMovieService() services.ClientMovieService[*types.JellyfinConfig]
	PlexMovieService() services.ClientMovieService[*types.PlexConfig]
	SubsonicMovieService() services.ClientMovieService[*types.SubsonicConfig]
}

type ClientSeriesServices interface {
	EmbySeriesService() services.ClientSeriesService[*types.EmbyConfig]
	JellyfinSeriesService() services.ClientSeriesService[*types.JellyfinConfig]
	PlexSeriesService() services.ClientSeriesService[*types.PlexConfig]
	SubsonicSeriesService() services.ClientSeriesService[*types.SubsonicConfig]
}

type ClientEpisodeServices interface {
	// TODO: implement
}

type ClientMusicServices interface {
	EmbyMusicService() services.ClientMusicService[*types.EmbyConfig]
	JellyfinMusicService() services.ClientMusicService[*types.JellyfinConfig]
	PlexMusicService() services.ClientMusicService[*types.PlexConfig]
	SubsonicMusicService() services.ClientMusicService[*types.SubsonicConfig]
}

type ClientPlaylistServices interface {
	EmbyPlaylistService() services.ClientListService[*types.EmbyConfig, *mediatypes.Playlist]
	JellyfinPlaylistService() services.ClientListService[*types.JellyfinConfig, *mediatypes.Playlist]
	PlexPlaylistService() services.ClientListService[*types.PlexConfig, *mediatypes.Playlist]
	SubsonicPlaylistService() services.ClientListService[*types.SubsonicConfig, *mediatypes.Playlist]
}
