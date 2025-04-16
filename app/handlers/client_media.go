package handlers

import (
	"suasor/client/types"
	"suasor/handlers"
)

type ClientMediaHandlers interface {
	ClientMovieHandlers
	ClientSeriesHandlers
	ClientMediaEpisodeHandlers
	ClientMusicHandlers
	ClientPlaylistHandlers
}

type ClientMovieHandlers interface {
	EmbyMovieHandler() *handlers.ClientMovieHandler[*types.EmbyConfig]
	JellyfinMovieHandler() *handlers.ClientMovieHandler[*types.JellyfinConfig]
	PlexMovieHandler() *handlers.ClientMovieHandler[*types.PlexConfig]
}

type ClientSeriesHandlers interface {
	EmbySeriesHandler() *handlers.ClientSeriesHandler[*types.EmbyConfig]
	JellyfinSeriesHandler() *handlers.ClientSeriesHandler[*types.JellyfinConfig]
	PlexSeriesHandler() *handlers.ClientSeriesHandler[*types.PlexConfig]
}

type ClientMediaEpisodeHandlers interface {
	// EmbyEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.EmbyConfig]
	// JellyfinEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.JellyfinConfig]
	// PlexEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.PlexConfig]
}

type ClientMusicHandlers interface {
	EmbyMusicHandler() *handlers.ClientMusicHandler[*types.EmbyConfig]
	JellyfinMusicHandler() *handlers.ClientMusicHandler[*types.JellyfinConfig]
	PlexMusicHandler() *handlers.ClientMusicHandler[*types.PlexConfig]
	SubsonicMusicHandler() *handlers.ClientMusicHandler[*types.SubsonicConfig]
}

type ClientPlaylistHandlers interface {
	EmbyPlaylistHandler() *handlers.ClientPlaylistHandler[*types.EmbyConfig]
	JellyfinPlaylistHandler() *handlers.ClientPlaylistHandler[*types.JellyfinConfig]
	PlexPlaylistHandler() *handlers.ClientPlaylistHandler[*types.PlexConfig]
	SubsonicPlaylistHandler() *handlers.ClientPlaylistHandler[*types.SubsonicConfig]
}

type ClientMediaTypeHandlers[T types.ClientMediaConfig] interface {
	MusicClientHandler() *handlers.ClientMusicHandler[T]
	MovieClientHandler() *handlers.ClientMovieHandler[T]
	SeriesClientHandler() *handlers.ClientSeriesHandler[T]
}

type ClientListHandlers[T types.ClientMediaConfig] interface {
	PlaylistClientHandler() *handlers.ClientPlaylistHandler[T]
	CollectionClientHandler() *handlers.ClientCollectionHandler[T]
}

//
// type ClientMediaHandlers interface {
// 	JellyfinMediaTypeHandlers() *ClientMediaTypeHandlers[*clienttypes.JellyfinConfig]
// 	EmbyMediaTypeHandlers() *ClientMediaTypeHandlers[*clienttypes.EmbyConfig]
// 	PlexMediaTypeHandlers() *ClientMediaTypeHandlers[*clienttypes.PlexConfig]
// 	SubsonicMediaTypeHandlers() *ClientMediaTypeHandlers[*clienttypes.SubsonicConfig]
// }
