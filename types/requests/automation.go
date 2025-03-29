package requests

// MediaAddRequest represents a request to add new media
type AutomationMediaAddRequest struct {
	Title             string
	Year              int
	QualityProfileID  int64
	MetadataProfileID int32 // For Lidarr
	Path              string
	TMDBID            int64  // For Radarr (movies)
	TVDBID            int64  // For Sonarr (TV shows)
	MusicBrainzID     string // For Lidarr (music)
	Tags              []int32
	Monitored         bool
	SearchForMedia    bool // Whether to search for the media after adding
}

// MediaUpdateRequest represents a request to update existing media
type AutomationMediaUpdateRequest struct {
	QualityProfileID  int64
	MetadataProfileID int32 // For Lidarr
	Path              string
	Tags              []int64
	Monitored         bool
}

// CreateTagRequest represents a request to create a tag in an automation client
type AutomationCreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// ExecuteCommandRequest represents a request to execute a command in an automation client
type AutomationExecuteCommandRequest struct {
	Command    string                 `json:"command" binding:"required"`
	Parameters map[string]interface{} `json:"parameters"`
}
