package handlers

import (
	"suasor/client/types"
	"suasor/handlers"
)

type ClientMediaHandlers interface {
	ClientMediaMovieHandlers
	ClientMediaSeriesHandlers
	ClientMediaEpisodeHandlers
	ClientMediaMusicHandlers
	ClientMediaPlaylistHandlers
}

type ClientMediaMovieHandlers interface {
	EmbyMovieHandler() *handlers.ClientMediaMovieHandler[*types.EmbyConfig]
	JellyfinMovieHandler() *handlers.ClientMediaMovieHandler[*types.JellyfinConfig]
	PlexMovieHandler() *handlers.ClientMediaMovieHandler[*types.PlexConfig]
}

type ClientMediaSeriesHandlers interface {
	EmbySeriesHandler() *handlers.ClientMediaSeriesHandler[*types.EmbyConfig]
	JellyfinSeriesHandler() *handlers.ClientMediaSeriesHandler[*types.JellyfinConfig]
	PlexSeriesHandler() *handlers.ClientMediaSeriesHandler[*types.PlexConfig]
}

type ClientMediaEpisodeHandlers interface {
	// EmbyEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.EmbyConfig]
	// JellyfinEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.JellyfinConfig]
	// PlexEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.PlexConfig]
}

type ClientMediaMusicHandlers interface {
	EmbyMusicHandler() *handlers.ClientMediaMusicHandler[*types.EmbyConfig]
	JellyfinMusicHandler() *handlers.ClientMediaMusicHandler[*types.JellyfinConfig]
	PlexMusicHandler() *handlers.ClientMediaMusicHandler[*types.PlexConfig]
	SubsonicMusicHandler() *handlers.ClientMediaMusicHandler[*types.SubsonicConfig]
}

type ClientMediaPlaylistHandlers interface {
	EmbyPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.EmbyConfig]
	JellyfinPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.JellyfinConfig]
	PlexPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.PlexConfig]
	SubsonicPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.SubsonicConfig]
}
