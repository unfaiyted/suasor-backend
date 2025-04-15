package requests

type PlaylistCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
	IsSmart     bool   `json:"isSmart"`
	Genre       string `json:"genre"`
	Year        int    `json:"year"`
	Rating      float32
	Duration    int
}

type PlaylistUpdateRequest struct {
	Name        string `json:"name"`
	IsPublic    bool   `json:"isPublic"`
	Description string `json:"description"`
}

type PlaylistAddTrackRequest struct {
	TrackID uint64 `json:"trackId"`
}

type PlaylistRemoveTrackRequest struct {
	TrackID uint64 `json:"trackId"`
}
