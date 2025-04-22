package types

// Common metadata types

// ExternalIDs represents external IDs for an item
type ExternalIDs struct {
	IMDBID   string `json:"imdbId,omitempty"`
	TMDBID   string `json:"tmdbId,omitempty"`
	TVDBId   string `json:"tvdbId,omitempty"`
	TraktID  string `json:"traktId,omitempty"`
	FanartTV string `json:"fanartTvId,omitempty"`
}

// MediaImage represents an image associated with a media item
type MediaImage struct {
	URL        string  `json:"url"`
	Type       string  `json:"type"` // poster, backdrop, logo, etc.
	Language   string  `json:"language,omitempty"`
	Width      int     `json:"width,omitempty"`
	Height     int     `json:"height,omitempty"`
	AspectRatio float64 `json:"aspectRatio,omitempty"`
}

// Video represents a video associated with a media item (trailer, etc.)
type Video struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Site        string `json:"site"` // YouTube, Vimeo, etc.
	Size        int    `json:"size,omitempty"`
	Type        string `json:"type"` // Trailer, Teaser, etc.
	Official    bool   `json:"official"`
	PublishedAt string `json:"publishedAt,omitempty"`
}

// Genre represents a genre
type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Person represents a person (actor, director, etc.)
type Person struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	ProfilePath     string       `json:"profilePath,omitempty"`
	KnownForDepartment string    `json:"knownForDepartment,omitempty"`
	Biography       string       `json:"biography,omitempty"`
	Birthday        string       `json:"birthday,omitempty"`
	Deathday        string       `json:"deathday,omitempty"`
	PlaceOfBirth    string       `json:"placeOfBirth,omitempty"`
	Gender          int          `json:"gender,omitempty"`
	Popularity      float64      `json:"popularity,omitempty"`
	Images          []MediaImage `json:"images,omitempty"`
	ExternalIDs     ExternalIDs  `json:"externalIds,omitempty"`
}

// CastMember represents a cast member in a movie or TV show
type CastMember struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Character   string `json:"character,omitempty"`
	Order       int    `json:"order,omitempty"`
	ProfilePath string `json:"profilePath,omitempty"`
	Gender      int    `json:"gender,omitempty"`
}

// CrewMember represents a crew member in a movie or TV show
type CrewMember struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Department  string `json:"department,omitempty"`
	Job         string `json:"job,omitempty"`
	ProfilePath string `json:"profilePath,omitempty"`
	Gender      int    `json:"gender,omitempty"`
}

// Credits represents the cast and crew for a movie or TV show
type Credits struct {
	Cast []CastMember `json:"cast,omitempty"`
	Crew []CrewMember `json:"crew,omitempty"`
}

// Movie represents a movie
type Movie struct {
	ID               string       `json:"id"`
	Title            string       `json:"title"`
	OriginalTitle    string       `json:"originalTitle,omitempty"`
	Overview         string       `json:"overview,omitempty"`
	Tagline          string       `json:"tagline,omitempty"`
	ReleaseDate      string       `json:"releaseDate,omitempty"`
	Runtime          int          `json:"runtime,omitempty"`
	Genres           []Genre      `json:"genres,omitempty"`
	PosterPath       string       `json:"posterPath,omitempty"`
	BackdropPath     string       `json:"backdropPath,omitempty"`
	VoteAverage      float64      `json:"voteAverage,omitempty"`
	VoteCount        int          `json:"voteCount,omitempty"`
	Popularity       float64      `json:"popularity,omitempty"`
	ProductionCountries []string  `json:"productionCountries,omitempty"`
	ProductionCompanies []string  `json:"productionCompanies,omitempty"`
	SpokenLanguages    []string   `json:"spokenLanguages,omitempty"`
	Status            string      `json:"status,omitempty"`
	Budget            int64       `json:"budget,omitempty"`
	Revenue           int64       `json:"revenue,omitempty"`
	Adult             bool        `json:"adult,omitempty"`
	Video             bool        `json:"video,omitempty"`
	Images            []MediaImage `json:"images,omitempty"`
	Videos            []Video     `json:"videos,omitempty"`
	Credits           Credits     `json:"credits,omitempty"`
	ExternalIDs       ExternalIDs `json:"externalIds,omitempty"`
	CollectionID      string      `json:"collectionId,omitempty"`
	CollectionName    string      `json:"collectionName,omitempty"`
}

// TVShow represents a TV show
type TVShow struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	OriginalName      string       `json:"originalName,omitempty"`
	Overview          string       `json:"overview,omitempty"`
	Tagline           string       `json:"tagline,omitempty"`
	FirstAirDate      string       `json:"firstAirDate,omitempty"`
	LastAirDate       string       `json:"lastAirDate,omitempty"`
	Genres            []Genre      `json:"genres,omitempty"`
	PosterPath        string       `json:"posterPath,omitempty"`
	BackdropPath      string       `json:"backdropPath,omitempty"`
	VoteAverage       float64      `json:"voteAverage,omitempty"`
	VoteCount         int          `json:"voteCount,omitempty"`
	Popularity        float64      `json:"popularity,omitempty"`
	OriginCountry     []string     `json:"originCountry,omitempty"`
	OriginalLanguage  string       `json:"originalLanguage,omitempty"`
	Status            string       `json:"status,omitempty"`
	Type              string       `json:"type,omitempty"`
	NumberOfSeasons   int          `json:"numberOfSeasons,omitempty"`
	NumberOfEpisodes  int          `json:"numberOfEpisodes,omitempty"`
	InProduction      bool         `json:"inProduction,omitempty"`
	Images            []MediaImage `json:"images,omitempty"`
	Videos            []Video      `json:"videos,omitempty"`
	Credits           Credits      `json:"credits,omitempty"`
	ExternalIDs       ExternalIDs  `json:"externalIds,omitempty"`
	CreatedBy         []Person     `json:"createdBy,omitempty"`
	Networks          []string     `json:"networks,omitempty"`
	Seasons           []TVSeason   `json:"seasons,omitempty"`
}

// TVSeason represents a TV season
type TVSeason struct {
	ID               string       `json:"id"`
	TVShowID         string       `json:"tvShowId,omitempty"`
	Name             string       `json:"name"`
	Overview         string       `json:"overview,omitempty"`
	SeasonNumber     int          `json:"seasonNumber"`
	AirDate          string       `json:"airDate,omitempty"`
	PosterPath       string       `json:"posterPath,omitempty"`
	EpisodeCount     int          `json:"episodeCount,omitempty"`
	VoteAverage      float64      `json:"voteAverage,omitempty"`
	VoteCount        int          `json:"voteCount,omitempty"`
	Images           []MediaImage `json:"images,omitempty"`
	Videos           []Video      `json:"videos,omitempty"`
	Credits          Credits      `json:"credits,omitempty"`
	ExternalIDs      ExternalIDs  `json:"externalIds,omitempty"`
	Episodes         []TVEpisode  `json:"episodes,omitempty"`
}

// TVEpisode represents a TV episode
type TVEpisode struct {
	ID                 string       `json:"id"`
	TVShowID           string       `json:"tvShowId,omitempty"`
	SeasonID           string       `json:"seasonId,omitempty"`
	Name               string       `json:"name"`
	Overview           string       `json:"overview,omitempty"`
	EpisodeNumber      int          `json:"episodeNumber"`
	SeasonNumber       int          `json:"seasonNumber"`
	AirDate            string       `json:"airDate,omitempty"`
	StillPath          string       `json:"stillPath,omitempty"`
	VoteAverage        float64      `json:"voteAverage,omitempty"`
	VoteCount          int          `json:"voteCount,omitempty"`
	Runtime            int          `json:"runtime,omitempty"`
	Images             []MediaImage `json:"images,omitempty"`
	Videos             []Video      `json:"videos,omitempty"`
	Credits            Credits      `json:"credits,omitempty"`
	ExternalIDs        ExternalIDs  `json:"externalIds,omitempty"`
	Crew               []CrewMember `json:"crew,omitempty"`
	GuestStars         []CastMember `json:"guestStars,omitempty"`
}

// Collection represents a collection of movies
type Collection struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Overview     string      `json:"overview,omitempty"`
	PosterPath   string      `json:"posterPath,omitempty"`
	BackdropPath string      `json:"backdropPath,omitempty"`
	Parts        []Movie     `json:"parts,omitempty"`
	Images       []MediaImage `json:"images,omitempty"`
}

// MovieCredit represents a movie credit for a person
type MovieCredit struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Character        string  `json:"character,omitempty"`
	Department       string  `json:"department,omitempty"`
	Job              string  `json:"job,omitempty"`
	PosterPath       string  `json:"posterPath,omitempty"`
	ReleaseDate      string  `json:"releaseDate,omitempty"`
	VoteAverage      float64 `json:"voteAverage,omitempty"`
	VoteCount        int     `json:"voteCount,omitempty"`
	Popularity       float64 `json:"popularity,omitempty"`
}

// TVCredit represents a TV credit for a person
type TVCredit struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Character        string  `json:"character,omitempty"`
	Department       string  `json:"department,omitempty"`
	Job              string  `json:"job,omitempty"`
	PosterPath       string  `json:"posterPath,omitempty"`
	FirstAirDate     string  `json:"firstAirDate,omitempty"`
	VoteAverage      float64 `json:"voteAverage,omitempty"`
	VoteCount        int     `json:"voteCount,omitempty"`
	Popularity       float64 `json:"popularity,omitempty"`
	EpisodeCount     int     `json:"episodeCount,omitempty"`
}

// MetadataClientConfig is the interface for metadata client configurations
type MetadataClientConfig interface {
	SupportsMovieMetadata() bool
	SupportsTVMetadata() bool
	SupportsPersonMetadata() bool
	SupportsCollectionMetadata() bool
}