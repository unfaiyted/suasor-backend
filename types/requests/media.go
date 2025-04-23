package requests

type MediaItemCreateRequest struct {
	Type string        `json:"type" binding:"required"`
	Data MediaItemData `json:"data" binding:"required"`
}

type MediaItemUpdateRequest struct {
	Type string        `json:"type" binding:"required"`
	Data MediaItemData `json:"data" binding:"required"`
}

type MediaItemData struct {
	// Base data
	ID uint64 `json:"id,omitempty"`
	// Movie data
	Title       string `json:"title,omitempty"`
	ReleaseYear int    `json:"releaseYear,omitempty"`
	// Series data
	Titles []string `json:"titles,omitempty"`
	// Episode data
	SeasonNumber  int `json:"seasonNumber,omitempty"`
	EpisodeNumber int `json:"episodeNumber,omitempty"`
	// Track data
	TrackNumber int `json:"trackNumber,omitempty"`
	// Album data
	AlbumTitle string `json:"albumTitle,omitempty"`
	// Artist data
	ArtistName string `json:"artistName,omitempty"`
	// Collection data
	CollectionName string `json:"collectionName,omitempty"`
}
