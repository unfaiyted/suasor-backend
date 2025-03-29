package types

type MediaClientConfig interface {
	ClientConfig
	isMediaClientConfig()
	GetClientType() MediaClientType

	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
	SupportsPlaylists() bool
	SupportsCollections() bool
	SupportsHistory() bool
}
