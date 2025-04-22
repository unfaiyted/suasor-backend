package responses

import (
	"suasor/clients/metadata"
)

// MetadataMovieResponse represents a response containing a movie
type MetadataMovieResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    *metadata.Movie `json:"data,omitempty"`
}

// MetadataMovieSearchResponse represents a response containing a list of movies
type MetadataMovieSearchResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Data    []*metadata.Movie `json:"data,omitempty"`
}

// MetadataTVShowResponse represents a response containing a TV show
type MetadataTVShowResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message,omitempty"`
	Data    *metadata.TVShow `json:"data,omitempty"`
}

// MetadataTVShowSearchResponse represents a response containing a list of TV shows
type MetadataTVShowSearchResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Data    []*metadata.TVShow `json:"data,omitempty"`
}

// MetadataTVSeasonResponse represents a response containing a TV season
type MetadataTVSeasonResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Data    *metadata.TVSeason `json:"data,omitempty"`
}

// MetadataTVEpisodeResponse represents a response containing a TV episode
type MetadataTVEpisodeResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message,omitempty"`
	Data    *metadata.TVEpisode `json:"data,omitempty"`
}

// MetadataPersonResponse represents a response containing a person
type MetadataPersonResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message,omitempty"`
	Data    *metadata.Person `json:"data,omitempty"`
}

// MetadataPersonSearchResponse represents a response containing a list of people
type MetadataPersonSearchResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Data    []*metadata.Person `json:"data,omitempty"`
}

// MetadataCollectionResponse represents a response containing a collection
type MetadataCollectionResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message,omitempty"`
	Data    *metadata.Collection `json:"data,omitempty"`
}

// MetadataCollectionSearchResponse represents a response containing a list of collections
type MetadataCollectionSearchResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    []*metadata.Collection `json:"data,omitempty"`
}

