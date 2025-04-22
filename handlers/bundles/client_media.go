package bundles

import (
	mediatypes "suasor/clients/media/types"
	"suasor/clients/types"
	"suasor/handlers"
)

type ClientMediaHandlers interface {
	ClientMovieHandlers
	ClientSeriesHandlers
	ClientMediaEpisodeHandlers
	ClientMusicHandlers
	ClientPlaylistHandlers
	ClientCollectionHandlers
}

type ClientMovieHandlers interface {
	EmbyMovieHandler() handlers.ClientMovieHandler[*types.EmbyConfig]
	JellyfinMovieHandler() handlers.ClientMovieHandler[*types.JellyfinConfig]
	PlexMovieHandler() handlers.ClientMovieHandler[*types.PlexConfig]
}

type ClientSeriesHandlers interface {
	EmbySeriesHandler() handlers.ClientSeriesHandler[*types.EmbyConfig]
	JellyfinSeriesHandler() handlers.ClientSeriesHandler[*types.JellyfinConfig]
	PlexSeriesHandler() handlers.ClientSeriesHandler[*types.PlexConfig]
}

type ClientMediaEpisodeHandlers interface {
	// EmbyEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.EmbyConfig]
	// JellyfinEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.JellyfinConfig]
	// PlexEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.PlexConfig]
}

type ClientMusicHandlers interface {
	EmbyMusicHandler() handlers.ClientMusicHandler[*types.EmbyConfig]
	JellyfinMusicHandler() handlers.ClientMusicHandler[*types.JellyfinConfig]
	PlexMusicHandler() handlers.ClientMusicHandler[*types.PlexConfig]
	SubsonicMusicHandler() handlers.ClientMusicHandler[*types.SubsonicConfig]
}

type ClientPlaylistHandlers interface {
	EmbyPlaylistHandler() handlers.ClientListHandler[*types.EmbyConfig, *mediatypes.Playlist]
	JellyfinPlaylistHandler() handlers.ClientListHandler[*types.JellyfinConfig, *mediatypes.Playlist]
	PlexPlaylistHandler() handlers.ClientListHandler[*types.PlexConfig, *mediatypes.Playlist]
	SubsonicPlaylistHandler() handlers.ClientListHandler[*types.SubsonicConfig, *mediatypes.Playlist]
}

type ClientCollectionHandlers interface {
	EmbyCollectionHandler() handlers.ClientListHandler[*types.EmbyConfig, *mediatypes.Collection]
	JellyfinCollectionHandler() handlers.ClientListHandler[*types.JellyfinConfig, *mediatypes.Collection]
	PlexCollectionHandler() handlers.ClientListHandler[*types.PlexConfig, *mediatypes.Collection]
	SubsonicCollectionHandler() handlers.ClientListHandler[*types.SubsonicConfig, *mediatypes.Collection]
}

type ClientMediaTypeHandlers[T types.ClientMediaConfig] interface {
	MusicClientHandler() handlers.ClientMusicHandler[T]
	MovieClientHandler() handlers.ClientMovieHandler[T]
	SeriesClientHandler() handlers.ClientSeriesHandler[T]
}

type ClientListHandlers[T types.ClientMediaConfig] interface {
	PlaylistClientHandler() *handlers.ClientListHandler[T, *mediatypes.Playlist]
	CollectionClientHandler() *handlers.ClientListHandler[T, *mediatypes.Collection]
}

//
// type ClientMediaHandlers interface {
// 	JellyfinMediaTypeHandlers() *ClientMediaTypeHandlers[*clienttypes.JellyfinConfig]
// 	EmbyMediaTypeHandlers() *ClientMediaTypeHandlers[*clienttypes.EmbyConfig]
// 	PlexMediaTypeHandlers() *ClientMediaTypeHandlers[*clienttypes.PlexConfig]
// 	SubsonicMediaTypeHandlers() *ClientMediaTypeHandlers[*clienttypes.SubsonicConfig]
