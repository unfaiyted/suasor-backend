package tmdb

import (
	"context"
	"fmt"
	"strconv"
	"suasor/clients/metadata"
	metadatatypes "suasor/clients/metadata/types"
	"suasor/clients/types"
	"time"

	tmdbClient "github.com/cyruzin/golang-tmdb"
)

type TMDBClient struct {
	metadata.ClientMetadata
	client *tmdbClient.Client
	config *types.TMDBConfig
}

// NewClient creates a new TMDB client
func NewTMDBClient(ctx, clientID uint64, config *types.TMDBConfig) (metadata.ClientMetadata, error) {
	client, err := tmdbClient.Init(config.GetApiKey())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TMDB client: %w", err)
	}

	clientMetadata, err := metadata.NewClientMetadata(ctx, clientID, config)
	if err != nil {
		return nil, err
	}

	tmdbClient := TMDBClient{
		ClientMetadata: clientMetadata,
		client:         client,
		config:         config,
	}

	return &tmdbClient, nil
}

// GetType returns the client type
func (c *TMDBClient) GetType() types.ClientType {
	return types.ClientTypeTMDB
}

// GetConfig returns the client configuration
func (c *TMDBClient) GetConfig() types.ClientConfig {
	return c.config
}

// SupportsMovieMetadata returns true because TMDB supports movie metadata
func (c *TMDBClient) SupportsMovieMetadata() bool {
	return true
}

// SupportsTVMetadata returns true because TMDB supports TV metadata
func (c *TMDBClient) SupportsTVMetadata() bool {
	return true
}

// SupportsPersonMetadata returns true because TMDB supports person metadata
func (c *TMDBClient) SupportsPersonMetadata() bool {
	return true
}

// SupportsCollectionMetadata returns true because TMDB supports collection metadata
func (c *TMDBClient) SupportsCollectionMetadata() bool {
	return true
}

// GetMovie retrieves movie details by ID
func (c *TMDBClient) GetMovie(ctx context.Context, id string) (*metadatatypes.Movie, error) {
	options := map[string]string{
		"append_to_response": "videos,images,credits",
		"language":           "en-US",
	}

	movieID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID format: %w", err)
	}

	movie, err := c.client.GetMovieDetails(movieID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %w", err)
	}

	// Convert to our format
	result := &metadatatypes.Movie{
		ID:            fmt.Sprintf("%d", movie.ID),
		Title:         movie.Title,
		OriginalTitle: movie.OriginalTitle,
		Overview:      movie.Overview,
		Tagline:       movie.Tagline,
		ReleaseDate:   movie.ReleaseDate,
		Runtime:       movie.Runtime,
		PosterPath:    movie.PosterPath,
		BackdropPath:  movie.BackdropPath,
		VoteAverage:   float64(movie.VoteAverage),
		VoteCount:     int(movie.VoteCount),
		Popularity:    float64(movie.Popularity),
		Status:        movie.Status,
		Budget:        int64(movie.Budget),
		Revenue:       int64(movie.Revenue),
		Adult:         movie.Adult,
		Video:         movie.Video,
	}

	// Convert genres
	if movie.Genres != nil {
		genres := make([]metadatatypes.Genre, 0, len(movie.Genres))
		for _, genre := range movie.Genres {
			genres = append(genres, metadatatypes.Genre{
				ID:   fmt.Sprintf("%d", genre.ID),
				Name: genre.Name,
			})
		}
		result.Genres = genres
	}

	// Convert collection if available
	if movie.BelongsToCollection.ID != 0 {
		result.CollectionID = fmt.Sprintf("%d", movie.BelongsToCollection.ID)
		result.CollectionName = movie.BelongsToCollection.Name
	}

	return result, nil
}

// SearchMovies searches for movies by query
func (c *TMDBClient) SearchMovies(ctx context.Context, query string) ([]*metadatatypes.Movie, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	result, err := c.client.GetSearchMovies(query, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	movies := make([]*metadatatypes.Movie, 0, len(result.Results))
	for i := range result.Results {
		// Convert the result directly
		movie := &metadatatypes.Movie{
			ID:            fmt.Sprintf("%d", result.Results[i].ID),
			Title:         result.Results[i].Title,
			OriginalTitle: result.Results[i].OriginalTitle,
			Overview:      result.Results[i].Overview,
			ReleaseDate:   result.Results[i].ReleaseDate,
			PosterPath:    result.Results[i].PosterPath,
			BackdropPath:  result.Results[i].BackdropPath,
			VoteAverage:   float64(result.Results[i].VoteAverage),
			VoteCount:     int(result.Results[i].VoteCount),
			Popularity:    float64(result.Results[i].Popularity),
			Adult:         result.Results[i].Adult,
			Video:         result.Results[i].Video,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetMovieRecommendations gets movie recommendations based on a movie ID
func (c *TMDBClient) GetMovieRecommendations(ctx context.Context, movieID string) ([]*metadatatypes.Movie, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	id, err := strconv.Atoi(movieID)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID format: %w", err)
	}

	result, err := c.client.GetMovieRecommendations(id, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie recommendations: %w", err)
	}

	movies := make([]*metadatatypes.Movie, 0, len(result.Results))
	for i := range result.Results {
		// Convert the result directly
		movie := &metadatatypes.Movie{
			ID:            fmt.Sprintf("%d", result.Results[i].ID),
			Title:         result.Results[i].Title,
			OriginalTitle: result.Results[i].OriginalTitle,
			Overview:      result.Results[i].Overview,
			ReleaseDate:   result.Results[i].ReleaseDate,
			PosterPath:    result.Results[i].PosterPath,
			BackdropPath:  result.Results[i].BackdropPath,
			VoteAverage:   float64(result.Results[i].VoteAverage),
			VoteCount:     int(result.Results[i].VoteCount),
			Popularity:    float64(result.Results[i].Popularity),
			Adult:         result.Results[i].Adult,
			Video:         result.Results[i].Video,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetPopularMovies gets popular movies
func (c *TMDBClient) GetPopularMovies(ctx context.Context) ([]*metadatatypes.Movie, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	result, err := c.client.GetMoviePopular(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular movies: %w", err)
	}

	movies := make([]*metadatatypes.Movie, 0, len(result.Results))
	for i := range result.Results {
		// Convert the result directly
		movie := &metadatatypes.Movie{
			ID:            fmt.Sprintf("%d", result.Results[i].ID),
			Title:         result.Results[i].Title,
			OriginalTitle: result.Results[i].OriginalTitle,
			Overview:      result.Results[i].Overview,
			ReleaseDate:   result.Results[i].ReleaseDate,
			PosterPath:    result.Results[i].PosterPath,
			BackdropPath:  result.Results[i].BackdropPath,
			VoteAverage:   float64(result.Results[i].VoteAverage),
			VoteCount:     int(result.Results[i].VoteCount),
			Popularity:    float64(result.Results[i].Popularity),
			Adult:         result.Results[i].Adult,
			Video:         result.Results[i].Video,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetTrendingMovies gets trending movies
func (c *TMDBClient) GetTrendingMovies(ctx context.Context) ([]*metadatatypes.Movie, error) {
	options := map[string]string{
		"page": "1",
	}

	result, err := c.client.GetTrending("movie", "week", options)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending movies: %w", err)
	}

	// The trending API returns movie results directly
	movies := make([]*metadatatypes.Movie, 0, len(result.Results))
	for _, movieResult := range result.Results {
		// Convert directly from the movie result struct
		movie := &metadatatypes.Movie{
			ID:           fmt.Sprintf("%d", movieResult.ID),
			Title:        movieResult.Title,
			Overview:     movieResult.Overview,
			PosterPath:   movieResult.PosterPath,
			BackdropPath: movieResult.BackdropPath,
			ReleaseDate:  movieResult.ReleaseDate,
			VoteAverage:  float64(movieResult.VoteAverage),
			VoteCount:    int(movieResult.VoteCount),
			Popularity:   float64(movieResult.Popularity),
			Adult:        movieResult.Adult,
			Video:        movieResult.Video,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetTVShow retrieves TV show details by ID
func (c *TMDBClient) GetTVShow(ctx context.Context, id string) (*metadatatypes.TVShow, error) {
	options := map[string]string{
		"append_to_response": "videos,images,credits",
		"language":           "en-US",
	}

	tvID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid TV show ID format: %w", err)
	}

	show, err := c.client.GetTVDetails(tvID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show details: %w", err)
	}

	// Convert to our format
	result := &metadatatypes.TVShow{
		ID:               fmt.Sprintf("%d", show.ID),
		Name:             show.Name,
		OriginalName:     show.OriginalName,
		Overview:         show.Overview,
		Tagline:          show.Tagline,
		FirstAirDate:     show.FirstAirDate,
		LastAirDate:      show.LastAirDate,
		PosterPath:       show.PosterPath,
		BackdropPath:     show.BackdropPath,
		VoteAverage:      float64(show.VoteAverage),
		VoteCount:        int(show.VoteCount),
		Popularity:       float64(show.Popularity),
		OriginCountry:    show.OriginCountry,
		OriginalLanguage: show.OriginalLanguage,
		Status:           show.Status,
		Type:             show.Type,
		NumberOfSeasons:  show.NumberOfSeasons,
		NumberOfEpisodes: show.NumberOfEpisodes,
		InProduction:     show.InProduction,
	}

	// Convert genres
	if show.Genres != nil {
		genres := make([]metadatatypes.Genre, 0, len(show.Genres))
		for _, genre := range show.Genres {
			genres = append(genres, metadatatypes.Genre{
				ID:   fmt.Sprintf("%d", genre.ID),
				Name: genre.Name,
			})
		}
		result.Genres = genres
	}

	// Convert seasons
	if show.Seasons != nil {
		seasons := make([]metadatatypes.TVSeason, 0, len(show.Seasons))
		for _, season := range show.Seasons {
			seasons = append(seasons, metadatatypes.TVSeason{
				ID:           fmt.Sprintf("%d", season.ID),
				TVShowID:     fmt.Sprintf("%d", show.ID),
				Name:         season.Name,
				Overview:     season.Overview,
				SeasonNumber: season.SeasonNumber,
				AirDate:      season.AirDate,
				PosterPath:   season.PosterPath,
				EpisodeCount: season.EpisodeCount,
			})
		}
		result.Seasons = seasons
	}

	// Get external IDs with a separate call
	externalIDOptions := map[string]string{}
	externalIDs, err := c.client.GetTVExternalIDs(tvID, externalIDOptions)
	if err == nil && externalIDs != nil {
		result.ExternalIDs = metadatatypes.ExternalIDs{
			IMDBID: externalIDs.IMDbID,
			TMDBID: fmt.Sprintf("%d", show.ID),
			TVDBId: fmt.Sprintf("%d", externalIDs.TVDBID),
		}
	}

	return result, nil
}

// SearchTVShows searches for TV shows by query
func (c *TMDBClient) SearchTVShows(ctx context.Context, query string) ([]*metadatatypes.TVShow, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	result, err := c.client.GetSearchTVShow(query, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search TV shows: %w", err)
	}

	shows := make([]*metadatatypes.TVShow, 0, len(result.Results))
	for i := range result.Results {
		// Convert to our format
		tvShow := &metadatatypes.TVShow{
			ID:               fmt.Sprintf("%d", result.Results[i].ID),
			Name:             result.Results[i].Name,
			OriginalName:     result.Results[i].OriginalName,
			Overview:         result.Results[i].Overview,
			FirstAirDate:     result.Results[i].FirstAirDate,
			PosterPath:       result.Results[i].PosterPath,
			BackdropPath:     result.Results[i].BackdropPath,
			VoteAverage:      float64(result.Results[i].VoteAverage),
			VoteCount:        int(result.Results[i].VoteCount),
			Popularity:       float64(result.Results[i].Popularity),
			OriginCountry:    result.Results[i].OriginCountry,
			OriginalLanguage: result.Results[i].OriginalLanguage,
		}
		shows = append(shows, tvShow)
	}

	return shows, nil
}

// GetTVShowRecommendations gets TV show recommendations
func (c *TMDBClient) GetTVShowRecommendations(ctx context.Context, tvShowID string) ([]*metadatatypes.TVShow, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	id, err := strconv.Atoi(tvShowID)
	if err != nil {
		return nil, fmt.Errorf("invalid TV show ID format: %w", err)
	}

	result, err := c.client.GetTVRecommendations(id, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show recommendations: %w", err)
	}

	shows := make([]*metadatatypes.TVShow, 0, len(result.Results))
	for i := range result.Results {
		// Convert to our format
		tvShow := &metadatatypes.TVShow{
			ID:               fmt.Sprintf("%d", result.Results[i].ID),
			Name:             result.Results[i].Name,
			OriginalName:     result.Results[i].OriginalName,
			Overview:         result.Results[i].Overview,
			FirstAirDate:     result.Results[i].FirstAirDate,
			PosterPath:       result.Results[i].PosterPath,
			BackdropPath:     result.Results[i].BackdropPath,
			VoteAverage:      float64(result.Results[i].VoteAverage),
			VoteCount:        int(result.Results[i].VoteCount),
			Popularity:       float64(result.Results[i].Popularity),
			OriginCountry:    result.Results[i].OriginCountry,
			OriginalLanguage: result.Results[i].OriginalLanguage,
		}
		shows = append(shows, tvShow)
	}

	return shows, nil
}

// GetPopularTVShows gets popular TV shows
func (c *TMDBClient) GetPopularTVShows(ctx context.Context) ([]*metadatatypes.TVShow, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	result, err := c.client.GetTVPopular(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular TV shows: %w", err)
	}

	shows := make([]*metadatatypes.TVShow, 0, len(result.Results))
	for i := range result.Results {
		// Convert to our format
		tvShow := &metadatatypes.TVShow{
			ID:               fmt.Sprintf("%d", result.Results[i].ID),
			Name:             result.Results[i].Name,
			OriginalName:     result.Results[i].OriginalName,
			Overview:         result.Results[i].Overview,
			FirstAirDate:     result.Results[i].FirstAirDate,
			PosterPath:       result.Results[i].PosterPath,
			BackdropPath:     result.Results[i].BackdropPath,
			VoteAverage:      float64(result.Results[i].VoteAverage),
			VoteCount:        int(result.Results[i].VoteCount),
			Popularity:       float64(result.Results[i].Popularity),
			OriginCountry:    result.Results[i].OriginCountry,
			OriginalLanguage: result.Results[i].OriginalLanguage,
		}
		shows = append(shows, tvShow)
	}

	return shows, nil
}

// GetTrendingTVShows gets trending TV shows
func (c *TMDBClient) GetTrendingTVShows(ctx context.Context) ([]*metadatatypes.TVShow, error) {
	options := map[string]string{
		"page": "1",
	}

	result, err := c.client.GetTrending("tv", "week", options)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending TV shows: %w", err)
	}

	// The trending API returns TV results directly
	shows := make([]*metadatatypes.TVShow, 0, len(result.Results))
	for _, tvResult := range result.Results {
		// Convert directly from the TV result struct
		tvShow := &metadatatypes.TVShow{
			ID:               fmt.Sprintf("%d", tvResult.ID),
			Name:             tvResult.Name,
			OriginalName:     tvResult.OriginalName,
			Overview:         tvResult.Overview,
			PosterPath:       tvResult.PosterPath,
			BackdropPath:     tvResult.BackdropPath,
			FirstAirDate:     tvResult.FirstAirDate,
			VoteAverage:      float64(tvResult.VoteAverage),
			VoteCount:        int(tvResult.VoteCount),
			Popularity:       float64(tvResult.Popularity),
			OriginCountry:    tvResult.OriginCountry,
			OriginalLanguage: tvResult.OriginalLanguage,
		}
		shows = append(shows, tvShow)
	}

	return shows, nil
}

// GetTVSeason retrieves a TV season
func (c *TMDBClient) GetTVSeason(ctx context.Context, tvShowID string, seasonNumber int) (*metadatatypes.TVSeason, error) {
	options := map[string]string{
		"language": "en-US",
	}

	id, err := strconv.Atoi(tvShowID)
	if err != nil {
		return nil, fmt.Errorf("invalid TV show ID format: %w", err)
	}

	season, err := c.client.GetTVSeasonDetails(id, seasonNumber, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV season details: %w", err)
	}

	// We need to convert the season to our format
	// For now, let's create a minimal implementation
	result := &metadatatypes.TVSeason{
		ID:           fmt.Sprintf("%d", season.ID),
		TVShowID:     tvShowID,
		Name:         season.Name,
		Overview:     season.Overview,
		SeasonNumber: season.SeasonNumber,
		AirDate:      season.AirDate,
		PosterPath:   season.PosterPath,
	}

	return result, nil
}

// GetTVEpisode retrieves a TV episode
func (c *TMDBClient) GetTVEpisode(ctx context.Context, tvShowID string, seasonNumber int, episodeNumber int) (*metadatatypes.TVEpisode, error) {
	options := map[string]string{
		"language": "en-US",
	}

	id, err := strconv.Atoi(tvShowID)
	if err != nil {
		return nil, fmt.Errorf("invalid TV show ID format: %w", err)
	}

	episode, err := c.client.GetTVEpisodeDetails(id, seasonNumber, episodeNumber, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV episode details: %w", err)
	}

	// We need to convert the episode to our format
	// For now, let's create a minimal implementation
	result := &metadatatypes.TVEpisode{
		ID:            fmt.Sprintf("%d", episode.ID),
		TVShowID:      tvShowID,
		SeasonID:      "", // We don't have this from the API directly
		Name:          episode.Name,
		Overview:      episode.Overview,
		EpisodeNumber: episode.EpisodeNumber,
		SeasonNumber:  episode.SeasonNumber,
		AirDate:       episode.AirDate,
		StillPath:     episode.StillPath,
		VoteAverage:   float64(episode.VoteAverage),
		VoteCount:     int(episode.VoteCount),
	}

	return result, nil
}

// GetPerson retrieves person details by ID
func (c *TMDBClient) GetPerson(ctx context.Context, id string) (*metadatatypes.Person, error) {
	options := map[string]string{
		"append_to_response": "images",
		"language":           "en-US",
	}

	personID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid person ID format: %w", err)
	}

	person, err := c.client.GetPersonDetails(personID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get person details: %w", err)
	}

	// Convert to our format
	result := &metadatatypes.Person{
		ID:                 fmt.Sprintf("%d", person.ID),
		Name:               person.Name,
		ProfilePath:        person.ProfilePath,
		KnownForDepartment: person.KnownForDepartment,
		Biography:          person.Biography,
		Birthday:           person.Birthday,
		Deathday:           person.Deathday,
		PlaceOfBirth:       person.PlaceOfBirth,
		Gender:             person.Gender,
		Popularity:         float64(person.Popularity),
	}

	// Convert external IDs
	if person.ExternalIDs != nil {
		result.ExternalIDs = metadatatypes.ExternalIDs{
			IMDBID: person.ExternalIDs.IMDbID,
			TMDBID: fmt.Sprintf("%d", person.ID),
		}
	}

	// Convert images if available
	if person.Images != nil && person.Images.Profiles != nil {
		images := make([]metadatatypes.MediaImage, 0, len(person.Images.Profiles))
		for _, profile := range person.Images.Profiles {
			images = append(images, metadatatypes.MediaImage{
				URL:         fmt.Sprintf("https://image.tmdb.org/t/p/original%s", profile.FilePath),
				Type:        "profile",
				Width:       profile.Width,
				Height:      profile.Height,
				AspectRatio: float64(profile.AspectRatio),
			})
		}
		result.Images = images
	}

	return result, nil
}

// SearchPeople searches for people by query
func (c *TMDBClient) SearchPeople(ctx context.Context, query string) ([]*metadatatypes.Person, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	result, err := c.client.GetSearchPeople(query, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search people: %w", err)
	}

	people := make([]*metadatatypes.Person, 0, len(result.Results))
	for i := range result.Results {
		// Convert to our format
		person := &metadatatypes.Person{
			ID:                 fmt.Sprintf("%d", result.Results[i].ID),
			Name:               result.Results[i].Name,
			ProfilePath:        result.Results[i].ProfilePath,
			KnownForDepartment: result.Results[i].KnownForDepartment,
			Gender:             result.Results[i].Gender,
			Popularity:         float64(result.Results[i].Popularity),
		}
		people = append(people, person)
	}

	return people, nil
}

// GetPersonMovieCredits retrieves a person's movie credits
func (c *TMDBClient) GetPersonMovieCredits(ctx context.Context, personID string) ([]*metadatatypes.MovieCredit, error) {
	options := map[string]string{
		"language": "en-US",
	}

	id, err := strconv.Atoi(personID)
	if err != nil {
		return nil, fmt.Errorf("invalid person ID format: %w", err)
	}

	credits, err := c.client.GetPersonMovieCredits(id, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get person movie credits: %w", err)
	}

	// Convert to our format
	result := make([]*metadatatypes.MovieCredit, 0, len(credits.Cast)+len(credits.Crew))

	// Add cast credits
	for _, credit := range credits.Cast {
		result = append(result, &metadatatypes.MovieCredit{
			ID:          fmt.Sprintf("%d", credit.ID),
			Title:       credit.Title,
			Character:   credit.Character,
			PosterPath:  credit.PosterPath,
			ReleaseDate: credit.ReleaseDate,
			VoteAverage: float64(credit.VoteAverage),
			VoteCount:   int(credit.VoteCount),
			Popularity:  float64(credit.Popularity),
		})
	}

	// Add crew credits
	for _, credit := range credits.Crew {
		result = append(result, &metadatatypes.MovieCredit{
			ID:          fmt.Sprintf("%d", credit.ID),
			Title:       credit.Title,
			Department:  credit.Department,
			Job:         credit.Job,
			PosterPath:  credit.PosterPath,
			ReleaseDate: credit.ReleaseDate,
			VoteAverage: float64(credit.VoteAverage),
			VoteCount:   int(credit.VoteCount),
			Popularity:  float64(credit.Popularity),
		})
	}

	return result, nil
}

// GetPersonTVCredits retrieves a person's TV credits
func (c *TMDBClient) GetPersonTVCredits(ctx context.Context, personID string) ([]*metadatatypes.TVCredit, error) {
	options := map[string]string{
		"language": "en-US",
	}

	id, err := strconv.Atoi(personID)
	if err != nil {
		return nil, fmt.Errorf("invalid person ID format: %w", err)
	}

	credits, err := c.client.GetPersonTVCredits(id, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get person TV credits: %w", err)
	}

	// Convert to our format
	result := make([]*metadatatypes.TVCredit, 0, len(credits.Cast)+len(credits.Crew))

	// Add cast credits
	for _, credit := range credits.Cast {
		result = append(result, &metadatatypes.TVCredit{
			ID:           fmt.Sprintf("%d", credit.ID),
			Name:         credit.Name,
			Character:    credit.Character,
			PosterPath:   credit.PosterPath,
			FirstAirDate: credit.FirstAirDate,
			VoteAverage:  float64(credit.VoteAverage),
			VoteCount:    int(credit.VoteCount),
			Popularity:   float64(credit.Popularity),
			EpisodeCount: credit.EpisodeCount,
		})
	}

	// Add crew credits
	for _, credit := range credits.Crew {
		result = append(result, &metadatatypes.TVCredit{
			ID:           fmt.Sprintf("%d", credit.ID),
			Name:         credit.Name,
			Department:   credit.Department,
			Job:          credit.Job,
			PosterPath:   credit.PosterPath,
			FirstAirDate: credit.FirstAirDate,
			VoteAverage:  float64(credit.VoteAverage),
			VoteCount:    int(credit.VoteCount),
			Popularity:   float64(credit.Popularity),
			EpisodeCount: credit.EpisodeCount,
		})
	}

	return result, nil
}

// GetCollection retrieves collection details by ID
func (c *TMDBClient) GetCollection(ctx context.Context, id string) (*metadatatypes.Collection, error) {
	options := map[string]string{
		"language": "en-US",
	}

	collectionID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid collection ID format: %w", err)
	}

	collection, err := c.client.GetCollectionDetails(collectionID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection details: %w", err)
	}

	// Convert to our format
	parts := make([]metadatatypes.Movie, 0, len(collection.Parts))
	for _, part := range collection.Parts {
		parts = append(parts, metadatatypes.Movie{
			ID:           fmt.Sprintf("%d", part.ID),
			Title:        part.Title,
			Overview:     part.Overview,
			PosterPath:   part.PosterPath,
			BackdropPath: part.BackdropPath,
			ReleaseDate:  part.ReleaseDate,
			Adult:        part.Adult,
		})
	}

	return &metadatatypes.Collection{
		ID:           fmt.Sprintf("%d", collection.ID),
		Name:         collection.Name,
		Overview:     collection.Overview,
		PosterPath:   collection.PosterPath,
		BackdropPath: collection.BackdropPath,
		Parts:        parts,
	}, nil
}

// SearchCollections searches for collections by query
func (c *TMDBClient) SearchCollections(ctx context.Context, query string) ([]*metadatatypes.Collection, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	result, err := c.client.GetSearchCollections(query, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search collections: %w", err)
	}

	collections := make([]*metadatatypes.Collection, 0, len(result.Results))
	for i := range result.Results {
		// Convert to our format
		collection := &metadatatypes.Collection{
			ID:           fmt.Sprintf("%d", result.Results[i].ID),
			Name:         result.Results[i].Name,
			Overview:     result.Results[i].Overview,
			PosterPath:   result.Results[i].PosterPath,
			BackdropPath: result.Results[i].BackdropPath,
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

// GetUpcomingMovies gets movies that are coming to theaters in the near future
func (c *TMDBClient) GetUpcomingMovies(ctx context.Context, daysAhead int) ([]*metadatatypes.Movie, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	if daysAhead > 0 {
		// Set the minimum release date to today
		today := time.Now().Format("2006-01-02")
		options["primary_release_date.gte"] = today

		// Set the maximum release date to N days in the future
		future := time.Now().AddDate(0, 0, daysAhead).Format("2006-01-02")
		options["primary_release_date.lte"] = future
	}

	// Call the Upcoming Movies API
	result, err := c.client.GetMovieUpcoming(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming movies: %w", err)
	}

	movies := make([]*metadatatypes.Movie, 0, len(result.Results))
	for i := range result.Results {
		// Convert to our format
		movie := &metadatatypes.Movie{
			ID:            fmt.Sprintf("%d", result.Results[i].ID),
			Title:         result.Results[i].Title,
			OriginalTitle: result.Results[i].OriginalTitle,
			Overview:      result.Results[i].Overview,
			ReleaseDate:   result.Results[i].ReleaseDate,
			PosterPath:    result.Results[i].PosterPath,
			BackdropPath:  result.Results[i].BackdropPath,
			VoteAverage:   float64(result.Results[i].VoteAverage),
			VoteCount:     int(result.Results[i].VoteCount),
			Popularity:    float64(result.Results[i].Popularity),
			Adult:         result.Results[i].Adult,
			Video:         result.Results[i].Video,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetNowPlayingMovies gets movies that are currently playing in theaters
func (c *TMDBClient) GetNowPlayingMovies(ctx context.Context, daysPast int) ([]*metadatatypes.Movie, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	if daysPast > 0 {
		// Set the maximum release date to today
		today := time.Now().Format("2006-01-02")
		options["primary_release_date.lte"] = today

		// Set the minimum release date to N days in the past
		past := time.Now().AddDate(0, 0, -daysPast).Format("2006-01-02")
		options["primary_release_date.gte"] = past
	}

	// Call the Now Playing Movies API
	result, err := c.client.GetMovieNowPlaying(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get now playing movies: %w", err)
	}

	movies := make([]*metadatatypes.Movie, 0, len(result.Results))
	for i := range result.Results {
		// Convert to our format
		movie := &metadatatypes.Movie{
			ID:            fmt.Sprintf("%d", result.Results[i].ID),
			Title:         result.Results[i].Title,
			OriginalTitle: result.Results[i].OriginalTitle,
			Overview:      result.Results[i].Overview,
			ReleaseDate:   result.Results[i].ReleaseDate,
			PosterPath:    result.Results[i].PosterPath,
			BackdropPath:  result.Results[i].BackdropPath,
			VoteAverage:   float64(result.Results[i].VoteAverage),
			VoteCount:     int(result.Results[i].VoteCount),
			Popularity:    float64(result.Results[i].Popularity),
			Adult:         result.Results[i].Adult,
			Video:         result.Results[i].Video,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetRecentTVShows gets TV shows that have recently premiered or returned
func (c *TMDBClient) GetRecentTVShows(ctx context.Context, daysWindow int) ([]*metadatatypes.TVShow, error) {
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
		"sort_by":  "first_air_date.desc",
	}

	if daysWindow > 0 {
		// Set the maximum premiere date to today
		today := time.Now().Format("2006-01-02")
		options["first_air_date.lte"] = today

		// Set the minimum premiere date to N days in the past
		past := time.Now().AddDate(0, 0, -daysWindow).Format("2006-01-02")
		options["first_air_date.gte"] = past
	}

	// Use the discover TV API to find recent shows
	result, err := c.client.GetDiscoverTV(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent TV shows: %w", err)
	}

	shows := make([]*metadatatypes.TVShow, 0, len(result.Results))
	for i := range result.Results {
		// Convert to our format
		tvShow := &metadatatypes.TVShow{
			ID:               fmt.Sprintf("%d", result.Results[i].ID),
			Name:             result.Results[i].Name,
			OriginalName:     result.Results[i].OriginalName,
			Overview:         result.Results[i].Overview,
			FirstAirDate:     result.Results[i].FirstAirDate,
			PosterPath:       result.Results[i].PosterPath,
			BackdropPath:     result.Results[i].BackdropPath,
			VoteAverage:      float64(result.Results[i].VoteAverage),
			VoteCount:        int(result.Results[i].VoteCount),
			Popularity:       float64(result.Results[i].Popularity),
			OriginCountry:    result.Results[i].OriginCountry,
			OriginalLanguage: result.Results[i].OriginalLanguage,
		}
		shows = append(shows, tvShow)
	}

	return shows, nil
}
