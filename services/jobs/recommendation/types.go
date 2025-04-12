package recommendation

import (
	"time"
)

// UserPreferenceProfile represents a user's media consumption preferences
type UserPreferenceProfile struct {
	// Common preferences
	PreferredLanguages    map[string]float32 // Language -> weight
	PreferredReleaseYears [2]int             // [min, max] year range
	ContentRatingRange    [2]string          // [min, max] content rating range (G, PG, PG-13, R, etc.)
	MinRating             float32            // Minimum rating threshold for recommendations
	WatchTimeOfDay        map[int]int        // Hour of day (0-23) -> count of watches
	TypicalSessionLength  map[string]float32 // Media type -> average session length in seconds
	OverallActivityLevel  map[string]float32 // Media type -> activity level (0-1)

	// Movie preferences
	WatchedMovieIDs         map[uint64]bool    // Watched movie IDs
	OwnedMovieIDs           map[string]bool    // Owned movie IDs
	RecentMovies            []MovieSummary     // Recently watched movies
	TopRatedMovies          []MovieSummary     // Highest rated movies
	FavoriteMovieGenres     map[string]float32 // Genre -> weight
	FavoriteActors          map[string]float32 // Actor name -> weight
	FavoriteDirectors       map[string]float32 // Director name -> weight
	MovieWatchTimes         map[string][]int64 // Hour -> counts
	MovieWatchDays          map[string]int     // Day of week -> count
	MovieTagPreferences     map[string]float32 // Tag -> weight
	ExcludedMovieGenres     []string           // Genres the user dislikes
	PreferredMovieGenres    []string           // Genres the user particularly likes
	MovieReleaseYearRange   [2]int             // [min, max] year range specifically for movies
	MovieDurationPreference [2]int             // [min, max] duration range specifically for movies

	// Series preferences
	WatchedSeriesIDs         map[uint64]bool    // Watched series IDs
	RecentSeries             []SeriesSummary    // Recently watched series
	TopRatedSeries           []SeriesSummary    // Highest rated series
	FavoriteSeriesGenres     map[string]float32 // Genre -> weight
	FavoriteShowrunners      map[string]float32 // Showrunner name -> weight
	SeriesWatchTimes         map[string][]int64 // Hour -> counts
	SeriesWatchDays          map[string]int     // Day of week -> count
	SeriesTagPreferences     map[string]float32 // Tag -> weight
	ExcludedSeriesGenres     []string           // Genres the user dislikes
	PreferredSeriesGenres    []string           // Genres the user particularly likes
	SeriesReleaseYearRange   [2]int             // [min, max] year range specifically for series
	SeriesEpisodeLengthRange [2]int             // [min, max] episode length range specifically for series
	PreferredSeriesStatus    []string           // Preferred series status (e.g., "Continuing", "Ended")
	OwnedSeriesIDs           map[string]bool    // Owned series IDs

	// Music preferences
	PlayedMusicIDs       map[uint64]bool    // Played music IDs
	RecentMusic          []MusicSummary     // Recently played music
	TopRatedMusic        []MusicSummary     // Highest rated music
	FavoriteMusicGenres  map[string]float32 // Genre -> weight
	FavoriteArtists      map[string]float32 // Artist name -> weight
	MusicPlayTimes       map[string][]int64 // Hour -> counts
	MusicPlayDays        map[string]int     // Day of week -> count
	MusicTagPreferences  map[string]float32 // Tag -> weight
	MusicMoodPreferences map[string]float32 // Mood -> weight
	ExcludedMusicGenres  []string           // Genres the user dislikes
	PreferredMusicGenres []string           // Genres the user particularly likes
	MusicDurationRange   [2]int             // [min, max] duration range specifically for music
	OwnedMusicIDs        map[string]bool    // Owned music IDs

	ExcludePlayed      bool    // Exclude previously played content from recommendations
	ProfileConfidence  float32 // User profile confidence
	BingeWatchingScore float32 // Binge-watching score
	ExplorationScore   float32 // Exploration score
	ContentCompleter   float32 // Content completion tendency
	// WTF?
	ContentRotationFreq float32 // Content rotation frequency
}

// RatingDetails contains detailed rating information
type RatingDetails struct {
	Overall    float32 `json:"overall"`    // Overall rating
	Story      float32 `json:"story"`      // Story rating (if applicable)
	Acting     float32 `json:"acting"`     // Acting rating (if applicable)
	Visuals    float32 `json:"visuals"`    // Visuals rating (if applicable)
	Audio      float32 `json:"audio"`      // Audio rating (if applicable)
	Engagement float32 `json:"engagement"` // Engagement rating (if applicable)
	Timestamp  int64   `json:"timestamp"`  // When the rating was given
}

// MovieSummary represents a summary of a movie for preference profiles
type MovieSummary struct {
	ID             uint64         `json:"id"`             // Internal ID
	Title          string         `json:"title"`          // Movie title
	Year           int            `json:"year"`           // Release year
	Genres         []string       `json:"genres"`         // Movie genres
	Director       string         `json:"director"`       // Director name
	Cast           []string       `json:"cast"`           // Main cast names
	Rating         float32        `json:"rating"`         // User rating
	DetailedRating *RatingDetails `json:"detailedRating"` // Detailed rating
	WatchCount     int            `json:"watchCount"`     // Times watched
	LastWatchDate  int64          `json:"lastWatchDate"`  // Last watched timestamp
	DurationMin    int            `json:"durationMin"`    // Duration in minutes
	IsFavorite     bool           `json:"isFavorite"`     // Marked as favorite
	UserTags       []string       `json:"userTags"`       // User-defined tags
	ExternalID     string         `json:"externalId"`     // External ID for metadata
}

// SeriesSummary represents a summary of a TV series for preference profiles
type SeriesSummary struct {
	ID              uint64         `json:"id"`              // Internal ID
	Title           string         `json:"title"`           // Series title
	Year            int            `json:"year"`            // Start year
	EndYear         int            `json:"endYear"`         // End year (if ended)
	Genres          []string       `json:"genres"`          // Series genres
	Creator         string         `json:"creator"`         // Creator/showrunner
	Cast            []string       `json:"cast"`            // Main cast
	Rating          float32        `json:"rating"`          // User rating
	DetailedRating  *RatingDetails `json:"detailedRating"`  // Detailed rating
	EpisodesWatched int            `json:"episodesWatched"` // Episodes watched
	TotalEpisodes   int            `json:"totalEpisodes"`   // Total episodes
	LastWatchDate   int64          `json:"lastWatchDate"`   // Last watched timestamp
	Status          string         `json:"status"`          // Series status (e.g., "Continuing")
	IsFavorite      bool           `json:"isFavorite"`      // Marked as favorite
	UserTags        []string       `json:"userTags"`        // User-defined tags
	ExternalID      string         `json:"externalId"`      // External ID for metadata
}

// MusicSummary represents a summary of a music track for preference profiles
type MusicSummary struct {
	ID             uint64         `json:"id"`             // Internal ID
	Title          string         `json:"title"`          // Track title
	Artist         string         `json:"artist"`         // Artist name
	Album          string         `json:"album"`          // Album name
	Year           int            `json:"year"`           // Release year
	Genres         []string       `json:"genres"`         // Music genres
	Rating         float32        `json:"rating"`         // User rating
	DetailedRating *RatingDetails `json:"detailedRating"` // Detailed rating
	PlayCount      int            `json:"playCount"`      // Times played
	LastPlayDate   int64          `json:"lastPlayDate"`   // Last played timestamp
	DurationSec    int            `json:"durationSec"`    // Duration in seconds
	IsFavorite     bool           `json:"isFavorite"`     // Marked as favorite
	UserTags       []string       `json:"userTags"`       // User-defined tags
	Mood           []string       `json:"mood"`           // Mood tags
	ExternalID     string         `json:"externalId"`     // External ID for metadata
}

// MovieRecommendation represents a recommended movie with explanation
type MovieRecommendation struct {
	MovieID          uint64    `json:"movieId"`          // Movie ID
	Title            string    `json:"title"`            // Movie title
	Year             int       `json:"year"`             // Release year
	Genres           []string  `json:"genres"`           // Movie genres
	Score            float32   `json:"score"`            // Recommendation score (0-1)
	Reasoning        string    `json:"reasoning"`        // Explanation for recommendation
	SimilarToMovies  []string  `json:"similarToMovies"`  // Similar to these movies
	MatchesActors    []string  `json:"matchesActors"`    // Matches these actors
	MatchesDirectors []string  `json:"matchesDirectors"` // Matches these directors
	MatchesGenres    []string  `json:"matchesGenres"`    // Matches these genres
	RecommendedBy    string    `json:"recommendedBy"`    // Recommendation source (AI, similar users, etc.)
	Timestamp        time.Time `json:"timestamp"`        // When the recommendation was generated
}

