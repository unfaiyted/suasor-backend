package requests

type ListType string

const (
	ListTypePlaylist   ListType = "playlist"
	ListTypeCollection ListType = "collection"
)

type ListCreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        ListType `json:"type"`
	IsPublic    bool     `json:"isPublic"`
	IsSmart     bool     `json:"isSmart"`
	Genre       string   `json:"genre"`
	Year        int      `json:"year"`
	Rating      float32
	Duration    int
}

type ListUpdateRequest struct {
	Name        string `json:"name"`
	IsPublic    bool   `json:"isPublic"`
	Description string `json:"description"`
}

type ListAddTrackRequest struct {
	TrackID uint64 `json:"trackId"`
}

type ListRemoveTrackRequest struct {
	TrackID uint64 `json:"trackId"`
}

type ListReorderRequest struct {
	ItemIDs []uint64 `json:"itemIds"`
}

type ListSearchRequest struct {
	Query string `json:"query"`
}
