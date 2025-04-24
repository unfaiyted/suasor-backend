package jobs

import (
	"time"
	// clienttypes "suasor/clients/types"
)

// MovieSummary contains a summary of a movie for recommendation purposes
type MovieSummary struct {
	Title             string         `json:"title"`
	Year              int            `json:"year"`
	Genres            []string       `json:"genres,omitempty"`
	DetailedRating    *RatingDetails `json:"detailedRating,omitempty"` // Enhanced rating information
	PlayCount         int            `json:"playCount,omitempty"`
	IsFavorite        bool           `json:"isFavorite,omitempty"`
	CompletionPercent float32        `json:"completionPercent,omitempty"`
	WatchDate         int64          `json:"watchDate,omitempty"` // Unix timestamp of last watch
	UserTags          []string       `json:"userTags,omitempty"`  // Custom tags applied by the user
	TMDB_ID           string         `json:"tmdbId,omitempty"`    // External ID for lookup
	Cast              []string       `json:"cast,omitempty"`      // Main actors
	Directors         []string       `json:"directors,omitempty"` // Movie directors
}

// SeriesSummary contains a summary of a TV series for recommendation purposes
type SeriesSummary struct {
	Title           string         `json:"title"`
	Year            int            `json:"year"`
	Genres          []string       `json:"genres,omitempty"`
	Rating          float32        `json:"rating,omitempty"`         // Basic rating for compatibility
	DetailedRating  *RatingDetails `json:"detailedRating,omitempty"` // Enhanced rating information
	Seasons         int            `json:"seasons,omitempty"`
	Status          string         `json:"status,omitempty"` // e.g., "Ended", "Continuing"
	IsFavorite      bool           `json:"isFavorite,omitempty"`
	EpisodesWatched int            `json:"episodesWatched,omitempty"` // Number of episodes watched
	TotalEpisodes   int            `json:"totalEpisodes,omitempty"`   // Total episodes in the series
	LastWatchDate   int64          `json:"lastWatchDate,omitempty"`   // Unix timestamp of last watch
	UserTags        []string       `json:"userTags,omitempty"`        // Custom tags applied by the user
	TMDB_ID         string         `json:"tmdbId,omitempty"`          // External ID for lookup
	Network         string         `json:"network,omitempty"`         // Network or platform that airs the series
	Cast            []string       `json:"cast,omitempty"`            // Main actors
	Showrunners     []string       `json:"showrunners,omitempty"`     // Series creators/showrunners
}

// MusicSummary contains a summary of a music track/artist for recommendation purposes
type MusicSummary struct {
	Title          string         `json:"title"`
	Artist         string         `json:"artist"`
	Album          string         `json:"album,omitempty"`
	Year           int            `json:"year,omitempty"`
	Genres         []string       `json:"genres,omitempty"`
	Rating         float32        `json:"rating,omitempty"`         // Basic rating for compatibility
	DetailedRating *RatingDetails `json:"detailedRating,omitempty"` // Enhanced rating information
	PlayCount      int            `json:"playCount,omitempty"`
	IsFavorite     bool           `json:"isFavorite,omitempty"`
	LastPlayDate   int64          `json:"lastPlayDate,omitempty"` // Unix timestamp of last play
	DurationSec    int            `json:"durationSec,omitempty"`  // Duration in seconds
	UserTags       []string       `json:"userTags,omitempty"`     // Custom tags applied by the user
	ExternalID     string         `json:"externalId,omitempty"`   // ID from music service
	Popularity     float32        `json:"popularity,omitempty"`   // Popularity metric (0-100)
	IsExplicit     bool           `json:"isExplicit,omitempty"`   // Whether the track has explicit content
	Mood           []string       `json:"mood,omitempty"`         // Mood categories (energetic, calm, etc.)
}

// MovieRecommendation represents a movie recommendation from the system
type MovieRecommendation struct {
	Title          string
	Year           int
	PopularityRank int
	Rating         float32
	Genres         []string
}

// UserPreferenceProfile holds user preferences for generating recommendations
type UserPreferenceProfile struct {
	// Common preferences
	PreferredReleaseYears [2]int             // [min, max] years
	ContentRatingRange    [2]string          // [min, max] content ratings
	MinRating             float32            // Minimum rating (0-10)
	PreferredLanguages    map[string]float32 // Language -> weight

	// Recommended content behavior
	ExcludeWatched       bool // Whether to exclude content already watched
	IncludeSimilarItems  bool // Whether to include items similar to favorites
	UseAIRecommendations bool // Whether to use AI for recommendations

	// Content type preferences
	NotifyForMovies bool
	NotifyForSeries bool
	NotifyForMusic  bool

	// Notification settings
	RatingThreshold   float64 // Minimum rating to notify about
	MaxNotifications  int     // Maximum notifications to create
	NotifyForUpcoming bool    // Notify about upcoming releases
	NotifyForRecent   bool    // Notify about recent releases

	// User content in library (to avoid notifying about owned content)
	OwnedMovieIDs  map[string]bool
	OwnedSeriesIDs map[string]bool
	OwnedMusicIDs  map[string]bool

	// Movie preferences
	WatchedMovieIDs         map[uint64]bool    // MediaItemID -> watched
	RecentMovies            []MovieSummary     // Recently watched movies
	TopRatedMovies          []MovieSummary     // Top rated movies
	FavoriteMovieGenres     map[string]float32 // Genre -> weight
	FavoriteActors          map[string]float32 // Actor -> weight
	FavoriteDirectors       map[string]float32 // Director -> weight
	MovieWatchTimes         map[string][]int64 // Time of day preferences (hour -> count)
	MovieWatchDays          map[string]int     // Days of week preferences (day -> count)
	MovieTagPreferences     map[string]float32 // User tag preferences for movies
	ExcludedMovieGenres     []string           // Genres to exclude
	PreferredMovieGenres    []string           // Genres to prefer
	MovieReleaseYearRange   [2]int             // Preferred year range for movies
	MovieDurationPreference [2]int             // Preferred duration range in minutes

	// Series preferences
	WatchedSeriesIDs         map[uint64]bool    // MediaItemID -> watched
	RecentSeries             []SeriesSummary    // Recently watched series
	TopRatedSeries           []SeriesSummary    // Top rated series
	FavoriteSeriesGenres     map[string]float32 // Genre -> weight
	FavoriteShowrunners      map[string]float32 // Showrunner -> weight
	SeriesWatchTimes         map[string][]int64 // Time of day preferences
	SeriesWatchDays          map[string]int     // Days of week preferences
	SeriesTagPreferences     map[string]float32 // User tag preferences for series
	ExcludedSeriesGenres     []string           // Genres to exclude
	PreferredSeriesGenres    []string           // Genres to prefer
	SeriesReleaseYearRange   [2]int             // Preferred year range for series
	SeriesEpisodeLengthRange [2]int             // Preferred episode length in minutes
	PreferredSeriesStatus    []string           // Preferred status (ended, continuing, etc.)

	// Music preferences
	PlayedMusicIDs        map[uint64]bool    // MediaItemID -> played
	RecentMusic           []MusicSummary     // Recently played music
	TopRatedMusic         []MusicSummary     // Top rated music
	FavoriteMusicGenres   map[string]float32 // Genre -> weight
	FavoriteArtists       map[string]float32 // Artist -> weight
	MusicPlayTimes        map[string][]int64 // Time of day preferences
	MusicPlayDays         map[string]int     // Days of week preferences
	MusicTagPreferences   map[string]float32 // User tag preferences for music
	ExcludedMusicGenres   []string           // Genres to exclude
	PreferredMusicGenres  []string           // Genres to prefer
	MusicReleaseDateRange [2]int             // Preferred year range for music
	MusicDurationRange    [2]int             // Preferred track length in seconds
	MusicMoodPreferences  map[string]float32 // Mood preferences for music

	// Watch/play history patterns
	WatchTimeOfDay       map[int]int        // Hour (0-23) -> count
	WatchDayOfWeek       map[string]int     // Day name -> count
	TypicalSessionLength map[string]float32 // Media type -> avg minutes
	BingeWatchingScore   float32            // Score indicating binge watching tendency (0-1)
	ContentRotationFreq  float32            // How often user switches genres/styles (0-1)

	// Advanced metrics
	GenreBreadth         float32            // How diverse their tastes are (0-1)
	ContentCompleter     float32            // How likely to finish what they start (0-1)
	NewContentScore      float32            // Affinity for new vs. classic content (0-1)
	PopularityInfluence  float32            // How much popularity affects choices (0-1)
	RatingInfluence      float32            // How much ratings affect choices (0-1)
	ExplorationScore     float32            // Willingness to try new things (0-1)
	OverallActivityLevel map[string]float32 // Activity level by media type (0-1)

	// Analysis metadata
	AnalysisTimestamp  int64   // When this profile was generated
	ProfileConfidence  float32 // Confidence in this profile (0-1)
	DataPointsAnalyzed int     // Number of history items analyzed
}

// NewReleases holds new release information by media type
type NewReleases struct {
	Movies []NewRelease `json:"movies"`
	Series []NewRelease `json:"series"`
	Music  []NewRelease `json:"music"`
}

type NewRelease struct {
	ID           string      `json:"id"`
	ExternalID   string      `json:"externalId"` // ID in the external system
	Title        string      `json:"title"`
	Description  string      `json:"description,omitempty"`
	ReleaseDate  time.Time   `json:"releaseDate"`
	MediaType    string      `json:"mediaType"` // "movie", "series", "music"
	Genres       []string    `json:"genres,omitempty"`
	Creators     []string    `json:"creators,omitempty"`     // Directors for movies, showrunners for series, artists for music
	ImageURL     string      `json:"imageUrl,omitempty"`     // Poster or cover art
	Rating       float64     `json:"rating,omitempty"`       // Rating if available
	Source       string      `json:"source,omitempty"`       // Source of this new release data
	SourceItemID uint64      `json:"sourceItemID,omitempty"` // ID of item in our system if known
	SourceURL    string      `json:"sourceUrl,omitempty"`    // URL to view details
	Metadata     interface{} `json:"-"`                      // Additional metadata (not stored)
}

// RatingDetails contains detailed rating information
type RatingDetails struct {
	Overall    float32            `json:"overall,omitempty"`    // Overall rating (0-10)
	Categories map[string]float32 `json:"categories,omitempty"` // Ratings for specific categories like "Acting", "Story", etc.
	Source     string             `json:"source,omitempty"`     // Source of the rating (user, aggregated, external service)
	MaxValue   float32            `json:"maxValue,omitempty"`   // Maximum possible rating value (default: 10)
	UserCount  int                `json:"userCount,omitempty"`  // Number of users who rated (for aggregated ratings)
	Timestamp  int64              `json:"timestamp,omitempty"`  // When the rating was last updated
}

// NotificationStats holds statistics about notifications sent
type NotificationStats struct {
	UsersNotified       int `json:"usersNotified"`
	MovieNotifications  int `json:"movieNotifications"`
	SeriesNotifications int `json:"seriesNotifications"`
	MusicNotifications  int `json:"musicNotifications"`
	TotalNotifications  int `json:"totalNotifications"`
}

// PlaylistClientInfo holds basic client information for playlist sync operations
// type PlaylistClientInfo struct {
// 	ClientID   uint64
// 	ClientType clienttypes.ClientMediaType
// 	Name       string
// 	IsPrimary  bool
// }

// PlaylistSyncStats contains statistics about a playlist sync operation
// type PlaylistSyncStats struct {
// 	totalSynced int
// 	created     int
// 	updated     int
// 	conflicts   int
// }

// WatchlistItem represents an item in a user's watchlist
type WatchlistItem struct {
	ID      uint64
	Type    string // movie, series
	Title   string
	Year    int
	TMDB_ID int
}

// AvailabilityChanges tracks changes in content availability
type AvailabilityChanges struct {
	New             int
	Removed         int
	NewServices     []string
	RemovedServices []string
}

// MaintenanceStats holds statistics for various maintenance operations
type MaintenanceStats struct {
	optimized int
	archived  int
	cleaned   int
	fixed     int
}

// CleanupStats holds statistics for cleanup operations
type CleanupStats struct {
	found    int
	resolved int
	removed  int
	fixed    int
}

// MetadataRefreshStats holds statistics for metadata refresh operations
type MetadataRefreshStats struct {
	checked int
	updated int
}

// NotificationType defines the type of notification
type NotificationType string

const (
	NotificationTypeNewRelease     NotificationType = "new_release"
	NotificationTypeUpcoming       NotificationType = "upcoming"
	NotificationTypeRecommendation NotificationType = "recommendation"
)

// UserNotification represents a notification for a user
type UserNotification struct {
	UserID      uint64
	Title       string
	Message     string
	Type        NotificationType
	ContentType string // movie, series, music
	ContentID   string
	ImageURL    string
	ActionURL   string
	Created     time.Time
	Expires     time.Time
	Priority    int // 1-5, with 5 being highest
	Read        bool
	Dismissed   bool
}
