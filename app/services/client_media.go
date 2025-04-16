package services

import (
	"suasor/client/types"
	"suasor/services"
)

type ClientMediaServices interface {
	ClientMovieServies
	ClientSeriesServices
	ClientEpisodeServices
	ClientMusicServices
	ClientPlaylistServices
}

type ClientMovieServies interface {
	EmbyMovieService() services.ClientMediaMovieService[*types.EmbyConfig]
	JellyfinMovieService() services.ClientMediaMovieService[*types.JellyfinConfig]
	PlexMovieService() services.ClientMediaMovieService[*types.PlexConfig]
	SubsonicMovieService() services.ClientMediaMovieService[*types.SubsonicConfig]
}

type ClientSeriesServices interface {
	EmbySeriesService() services.ClientMediaSeriesService[*types.EmbyConfig]
	JellyfinSeriesService() services.ClientMediaSeriesService[*types.JellyfinConfig]
	PlexSeriesService() services.ClientMediaSeriesService[*types.PlexConfig]
	SubsonicSeriesService() services.ClientMediaSeriesService[*types.SubsonicConfig]
}

type ClientEpisodeServices interface {
	// TODO: implement
}

type ClientMusicServices interface {
	EmbyMusicService() services.ClientMediaMusicService[*types.EmbyConfig]
	JellyfinMusicService() services.ClientMediaMusicService[*types.JellyfinConfig]
	PlexMusicService() services.ClientMediaMusicService[*types.PlexConfig]
	SubsonicMusicService() services.ClientMediaMusicService[*types.SubsonicConfig]
}

type ClientPlaylistServices interface {
	EmbyPlaylistService() services.ClientMediaPlaylistService[*types.EmbyConfig]
	JellyfinPlaylistService() services.ClientMediaPlaylistService[*types.JellyfinConfig]
	PlexPlaylistService() services.ClientMediaPlaylistService[*types.PlexConfig]
	SubsonicPlaylistService() services.ClientMediaPlaylistService[*types.SubsonicConfig]
}
