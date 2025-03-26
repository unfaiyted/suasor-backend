package providers

// MediaContentProvider combines multiple provider interfaces
// This is useful for clients that implement multiple provider types
type MediaContentProvider interface {
	MovieProvider
	TVShowProvider
	MusicProvider
	PlaylistProvider
	CollectionProvider
	HistoryProvider
}

// MovieMusicProvider combines movie and music provider interfaces
// This is useful for testing functions that need both movie and music capabilities
type MovieMusicProvider interface {
	MovieProvider
	MusicProvider
}
