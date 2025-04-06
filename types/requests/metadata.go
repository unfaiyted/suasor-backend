package requests

// MetadataMovieRequest represents a request to get a movie
type MetadataMovieRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	MovieID  string `json:"movieId" form:"movieId" validate:"required"`
}

// MetadataMovieSearchRequest represents a request to search for movies
type MetadataMovieSearchRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	Query    string `json:"query" form:"query" validate:"required"`
}

// MetadataMovieRecommendationsRequest represents a request to get movie recommendations
type MetadataMovieRecommendationsRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	MovieID  string `json:"movieId" form:"movieId" validate:"required"`
}

// MetadataPopularMoviesRequest represents a request to get popular movies
type MetadataPopularMoviesRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
}

// MetadataTrendingMoviesRequest represents a request to get trending movies
type MetadataTrendingMoviesRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
}

// MetadataTVShowRequest represents a request to get a TV show
type MetadataTVShowRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	TVShowID string `json:"tvShowId" form:"tvShowId" validate:"required"`
}

// MetadataTVShowSearchRequest represents a request to search for TV shows
type MetadataTVShowSearchRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	Query    string `json:"query" form:"query" validate:"required"`
}

// MetadataTVShowRecommendationsRequest represents a request to get TV show recommendations
type MetadataTVShowRecommendationsRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	TVShowID string `json:"tvShowId" form:"tvShowId" validate:"required"`
}

// MetadataPopularTVShowsRequest represents a request to get popular TV shows
type MetadataPopularTVShowsRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
}

// MetadataTrendingTVShowsRequest represents a request to get trending TV shows
type MetadataTrendingTVShowsRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
}

// MetadataTVSeasonRequest represents a request to get a TV season
type MetadataTVSeasonRequest struct {
	ClientID     uint64 `json:"clientId" form:"clientId" validate:"required"`
	TVShowID     string `json:"tvShowId" form:"tvShowId" validate:"required"`
	SeasonNumber int    `json:"seasonNumber" form:"seasonNumber" validate:"required"`
}

// MetadataTVEpisodeRequest represents a request to get a TV episode
type MetadataTVEpisodeRequest struct {
	ClientID      uint64 `json:"clientId" form:"clientId" validate:"required"`
	TVShowID      string `json:"tvShowId" form:"tvShowId" validate:"required"`
	SeasonNumber  int    `json:"seasonNumber" form:"seasonNumber" validate:"required"`
	EpisodeNumber int    `json:"episodeNumber" form:"episodeNumber" validate:"required"`
}

// MetadataPersonRequest represents a request to get a person
type MetadataPersonRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	PersonID string `json:"personId" form:"personId" validate:"required"`
}

// MetadataPersonSearchRequest represents a request to search for people
type MetadataPersonSearchRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	Query    string `json:"query" form:"query" validate:"required"`
}

// MetadataCollectionRequest represents a request to get a collection
type MetadataCollectionRequest struct {
	ClientID     uint64 `json:"clientId" form:"clientId" validate:"required"`
	CollectionID string `json:"collectionId" form:"collectionId" validate:"required"`
}

// MetadataCollectionSearchRequest represents a request to search for collections
type MetadataCollectionSearchRequest struct {
	ClientID uint64 `json:"clientId" form:"clientId" validate:"required"`
	Query    string `json:"query" form:"query" validate:"required"`
}