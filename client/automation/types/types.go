package interfaces

import (
	"time"
)

// SystemStatus represents system information from the automation tool
type SystemStatus struct {
	Version     string
	StartupPath string
	AppData     string
	OsName      string
	IsUpdated   bool
	HasUpdate   bool
	Branch      string
}

type ExternalID struct {
	Source string
	Value  string
}

// RatingType the model 'RatingType'
type RatingType string

// List of RatingType
const (
	RATINGTYPE_USER   RatingType = "user"
	RATINGTYPE_CRITIC RatingType = "critic"
)

type DownloadedStatus string

const (
	DOWNLOADEDSTATUS_COMPLETE  DownloadedStatus = "complete"
	DOWNLOADEDSTATUS_REQUESTED DownloadedStatus = "requested"
	DOWNLOADEDSTATUS_PARTIAL   DownloadedStatus = "partial"
	DOWNLOADEDSTATUS_NONE      DownloadedStatus = "none"
)

type Rating struct {
	Votes *int32      `json:"votes,omitempty"`
	Value *float64    `json:"value,omitempty"`
	Type  *RatingType `json:"type,omitempty"`
}

type AutomationMediaType string

const (
	AUTOMEDIATYPE_MOVIE  AutomationMediaType = "movie"
	AUTOMEDIATYPE_SERIES AutomationMediaType = "series"
	AUTOMEDIATYPE_ARTIST AutomationMediaType = "artist"
)

type AutomationStatusType string

const (
	AUTOSTATUSTYPE_CONTINUING AutomationStatusType = "continuing"
	AUTOSTATUSTYPE_ENDED      AutomationStatusType = "ended"
	AUTOSTATUSTYPE_UPCOMING   AutomationStatusType = "upcoming"
	AUTOSTATUSTYPE_DELETED    AutomationStatusType = "deleted"
	AUTOSTATUSTYPE_IN_CINEMAS AutomationStatusType = "inCinemas"
	AUTOSTATUSTYPE_RELEASED   AutomationStatusType = "released"
)

// type AutomationData interface {
// 	AutomationMovie | AutomationTVShow | AutomationArtist | AutomationAlbum | AutomationTrack
// }

// AutomationData defines the allowed types for AutomationMediaItem's Data field
type AutomationData interface {
	isAutomationData()
	GetMediaType() AutomationMediaType
}

// Implement the marker method for each allowed type
func (AutomationMovie) isAutomationData()   {}
func (AutomationTVShow) isAutomationData()  {}
func (AutomationEpisode) isAutomationData() {}
func (AutomationArtist) isAutomationData()  {}
func (AutomationAlbum) isAutomationData()   {}
func (AutomationTrack) isAutomationData()   {}

// MediaItem represents a generic media item (movie, show, or music)
type AutomationMediaItem[T AutomationData] struct {
	ID               uint64
	ClientID         uint32
	ClientType       AutomationClientType
	Title            string
	Overview         string
	MediaType        AutomationMediaType
	Year             int32
	AddedAt          time.Time
	Ratings          []Rating
	DownloadedStatus DownloadedStatus
	ExternalIDs      []ExternalID
	Status           AutomationStatusType
	Path             string
	Genres           []string
	QualityProfile   QualityProfileSummary
	Images           []AutomationMediaImage
	Monitored        bool
	Data             T
}

type AutomationMovie struct {
	Year        int32
	ReleaseDate time.Time
}

type AutomationTVShow struct {
	ReleaseDate time.Time
	Year        int32
}

type AutomationEpisode struct {
	ReleaseDate time.Time
}

type AutomationTrack struct {
	AlbumName  string
	ArtistName string
}

type AutomationArtist struct {
	Albums                []AutomationAlbum
	MetadataProfile       MetadataProfile
	MostRecentReleaseDate time.Time
}

type AutomationAlbum struct {
	ArtistName  string
	ArtistID    string
	ReleaseDate time.Time
}

// MediaImage represents an image associated wit media
type AutomationMediaImage struct {
	URL       string
	CoverType string // "poster", "fanart", "banner", etc.
}

// SearchOptions defines options for search requests
type SearchOptions struct {
	Limit  int
	Offset int
}

// LibraryQueryOptions defines options for library queries
type LibraryQueryOptions struct {
	Limit  int
	Offset int
	SortBy string
	Filter string
}

// QualityProfile represents a quality profile configuration
type QualityProfile struct {
	ID    int64
	Name  string
	Items []QualityItem
}

// QualityProfileSummary is a simplified quality profile reference
type QualityProfileSummary struct {
	ID   int64
	Name string
}

// QualityProfileSummary is a simplified quality profile reference
type MetadataProfile struct {
	ID   int32
	Name string
}

// QualityItem represents an item in a quality profile
type QualityItem struct {
	ID         int64
	Name       string
	Resolution string
	Source     string
}

// Tag represents a media tag
type Tag struct {
	ID   int64
	Name string
}

// CalendarItem represents an upcoming media item
type CalendarItem struct {
	Details AutomationMediaItem[AutomationData]
	AirDate time.Time
}

// Command represents a command to execute
type Command struct {
	Name       string
	MediaIDs   []int64
	Parameters map[string]any
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	ID        int64
	Name      string
	Status    string
	StartedAt time.Time
	Completed time.Time
	Message   string
}

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

// BaseAutomationTool provides common behavior for all automation tools
type BaseAutomationTool struct {
	ClientID   uint32
	ClientType AutomationClientType
	URL        string
	APIKey     string
}
