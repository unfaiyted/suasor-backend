package plex

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

// SupportsSearch indicates that the Plex client supports search functionality
func (c *PlexClient) SupportsSearch() bool { return true }

// Search searches for media items in Plex
func (c *PlexClient) Search(ctx context.Context, options *types.QueryOptions) (responses.SearchResults, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("query", options.Query).
		Str("mediaType", string(options.MediaType)).
		Msg("Searching media items in Plex server")

	// Initialize the result container
	var results responses.SearchResults

	// We'll take a different approach - use GetMovies, GetSeries, etc. with a title filter
	// Since the search endpoint has format issues

	// Determine which types to search based on media type
	if options.MediaType == "" || options.MediaType == types.MediaTypeAll || options.MediaType == types.MediaTypeMovie {
		// Get and filter movies
		movieResults, err := c.searchMovies(ctx, options)
		if err != nil {
			log.Warn().
				Err(err).
				Str("query", options.Query).
				Msg("Error searching movies, continuing with other types")
		} else {
			results.Movies = movieResults
		}
	}

	if options.MediaType == "" || options.MediaType == types.MediaTypeAll || options.MediaType == types.MediaTypeSeries {
		// Get and filter series
		seriesResults, err := c.searchSeries(ctx, options)
		if err != nil {
			log.Warn().
				Err(err).
				Str("query", options.Query).
				Msg("Error searching series, continuing with other types")
		} else {
			results.Series = seriesResults
		}
	}

	if options.MediaType == "" || options.MediaType == types.MediaTypeAll ||
		options.MediaType == types.MediaTypeArtist || options.MediaType == types.MediaTypeAlbum ||
		options.MediaType == types.MediaTypeTrack {
		// Get and filter music
		artistResults, albumResults, trackResults, err := c.searchMusic(ctx, options)
		if err != nil {
			log.Warn().
				Err(err).
				Str("query", options.Query).
				Msg("Error searching music, continuing with other types")
		} else {
			results.Artists = artistResults
			results.Albums = albumResults
			results.Tracks = trackResults
		}
	}

	log.Info().
		Int("moviesCount", len(results.Movies)).
		Int("seriesCount", len(results.Series)).
		Int("episodesCount", len(results.Episodes)).
		Int("artistsCount", len(results.Artists)).
		Int("albumsCount", len(results.Albums)).
		Int("tracksCount", len(results.Tracks)).
		Int("playlistsCount", len(results.Playlists)).
		Int("collectionsCount", len(results.Collections)).
		Msg("Completed search in Plex server")

	return results, nil
}

// searchMovies searches for movies by title using the Plex search API
func (c *PlexClient) searchMovies(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Movie], error) {
	log := logger.LoggerFromContext(ctx)

	// First, find the movie library section
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to find movie section key")
		return nil, fmt.Errorf("failed to find movie section: %w", err)
	}

	if movieSectionKey == "" {
		log.Info().Msg("No movie library found in Plex")
		return []*models.MediaItem[*types.Movie]{}, nil
	}

	sectionKeyFloat, err := strconv.ParseFloat(movieSectionKey, 64)
	if err != nil {
		log.Error().
			Err(err).
			Str("sectionKey", movieSectionKey).
			Msg("Invalid section key format")
		return nil, fmt.Errorf("invalid section key: %w", err)
	}

	// Set limit if specified
	var limitFloat *float64
	if options.Limit > 0 {
		limit := float64(options.Limit)
		limitFloat = &limit
	}

	// Use the PerformSearch API to search for movies
	log.Debug().
		Str("query", options.Query).
		Float64("sectionID", sectionKeyFloat).
		Msg("Performing Plex search for movies")

	// Try to use the PerformSearch API with detailed logging
	searchResp, err := c.plexAPI.Search.PerformSearch(ctx, options.Query, &sectionKeyFloat, limitFloat)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", options.Query).
			Float64("sectionID", sectionKeyFloat).
			Msg("Failed to search movies in Plex")

		// Log more details about the error for debugging
		log.Debug().
			Str("errorType", fmt.Sprintf("%T", err)).
			Str("errorMessage", err.Error()).
			Msg("Detailed error information")

		// Fall back to loading all and filtering if API search fails
		log.Info().Msg("Falling back to loading all movies and filtering locally")
		return c.fallbackSearchMovies(ctx, options)
	}

	// Process search results
	var movies []*models.MediaItem[*types.Movie]

	// Log response details
	log.Debug().
		Int("statusCode", searchResp.StatusCode).
		Msg("Plex search response received")

	// Use our custom parser to parse the raw response
	plexSearchResp, err := ParsePlexSearchResponse(ctx, searchResp.RawResponse)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to parse Plex search response with custom parser")

		// Log the raw response for further debugging
		if searchResp.RawResponse != nil {
			log.Debug().
				Interface("rawResponse", searchResp.RawResponse).
				Msg("Raw response that failed to parse")
		}

		// Fall back to the old method
		log.Info().Msg("Falling back to loading all movies and filtering locally")
		return c.fallbackSearchMovies(ctx, options)
	}

	// Process the parsed response
	if plexSearchResp != nil && len(plexSearchResp.MediaContainer.Hub) > 0 {
		// Log the hub information
		for i, hub := range plexSearchResp.MediaContainer.Hub {
			log.Debug().
				Int("hubIndex", i).
				Str("hubType", hub.Type).
				Str("hubTitle", hub.Title).
				Int("itemCount", len(hub.Metadata)).
				Msg("Found search hub")

			// Process movie items in this hub
			if hub.Type == "movie" || hub.Type == "1" {
				log.Debug().
					Int("movieCount", len(hub.Metadata)).
					Msg("Processing movie results")

				for _, item := range hub.Metadata {
					if item.RatingKey == "" || item.Title == "" {
						continue
					}

					// Create movie object
					details := &types.MediaDetails{
						Title: item.Title,
					}

					year := item.GetYear()
					if year > 0 {
						details.ReleaseYear = year
					}

					movie := &types.Movie{
						Details: details,
					}

					// Create media item
					mediaItem, err := GetMediaItem[*types.Movie](ctx, c, movie, item.RatingKey)
					if err != nil {
						log.Warn().
							Err(err).
							Str("movieID", item.RatingKey).
							Str("movieTitle", item.Title).
							Msg("Error creating media item for movie")
						continue
					}

					movies = append(movies, mediaItem)
				}
			}
		}
	} else {
		log.Debug().
			Msg("No search hubs or results found in Plex response")

		// Check the structure of the response for debugging
		if plexSearchResp == nil {
			log.Debug().Msg("Parsed Response is nil")
		} else if len(plexSearchResp.MediaContainer.Hub) == 0 {
			log.Debug().Msg("Hub array is empty")
		}

		// Fall back to the old method
		log.Info().Msg("Falling back to loading all movies and filtering locally")
		return c.fallbackSearchMovies(ctx, options)
	}

	log.Info().
		Int("movieResults", len(movies)).
		Str("query", options.Query).
		Msg("Successfully searched for movies in Plex")

	return movies, nil
}

// fallbackSearchMovies is a fallback method that loads all movies and filters them locally
func (c *PlexClient) fallbackSearchMovies(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Movie], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Using fallback search method for movies")

	// Use the GetMovies method, which is known to work
	searchOptions := &types.QueryOptions{
		Limit:  options.Limit,
		Offset: options.Offset,
	}

	// Get all movies and filter them on our side
	movieItems, err := c.GetMovies(ctx, searchOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}

	// If no search query, return all movies
	if options.Query == "" {
		return movieItems, nil
	}

	// Filter movies by title matching the query
	searchLower := strings.ToLower(options.Query)
	var filteredMovies []*models.MediaItem[*types.Movie]

	for _, movie := range movieItems {
		if movie.Data == nil || movie.Data.Details == nil {
			continue
		}

		movieTitle := strings.ToLower(movie.Data.Details.Title)
		if strings.Contains(movieTitle, searchLower) {
			filteredMovies = append(filteredMovies, movie)
		}
	}

	log.Debug().
		Int("searchResults", len(filteredMovies)).
		Int("totalMovies", len(movieItems)).
		Str("query", options.Query).
		Msg("Filtered movie search results (fallback method)")

	return filteredMovies, nil
}

// searchSeries searches for TV series by title using the Plex search API
func (c *PlexClient) searchSeries(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)

	// First, find the TV series library section
	seriesSectionKey, err := c.findLibrarySectionByType(ctx, "show")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to find TV show section key")
		return nil, fmt.Errorf("failed to find TV show section: %w", err)
	}

	if seriesSectionKey == "" {
		log.Info().Msg("No TV show library found in Plex")
		return []*models.MediaItem[*types.Series]{}, nil
	}

	sectionKeyFloat, err := strconv.ParseFloat(seriesSectionKey, 64)
	if err != nil {
		log.Error().
			Err(err).
			Str("sectionKey", seriesSectionKey).
			Msg("Invalid section key format")
		return nil, fmt.Errorf("invalid section key: %w", err)
	}

	// Set limit if specified
	var limitFloat *float64
	if options.Limit > 0 {
		limit := float64(options.Limit)
		limitFloat = &limit
	}

	// Use the PerformSearch API to search for TV shows
	log.Debug().
		Str("query", options.Query).
		Float64("sectionID", sectionKeyFloat).
		Msg("Performing Plex search for TV shows")

	// Try to use the PerformSearch API with detailed logging
	searchResp, err := c.plexAPI.Search.PerformSearch(ctx, options.Query, &sectionKeyFloat, limitFloat)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", options.Query).
			Float64("sectionID", sectionKeyFloat).
			Msg("Failed to search TV shows in Plex")

		// Log more details about the error for debugging
		log.Debug().
			Str("errorType", fmt.Sprintf("%T", err)).
			Str("errorMessage", err.Error()).
			Msg("Detailed error information")

		// Fall back to loading all and filtering if API search fails
		log.Info().Msg("Falling back to loading all TV shows and filtering locally")
		return c.fallbackSearchSeries(ctx, options)
	}

	// Process search results
	var series []*models.MediaItem[*types.Series]

	// Log response details
	log.Debug().
		Int("statusCode", searchResp.StatusCode).
		Msg("Plex search response received")

	// Use our custom parser to parse the raw response
	plexSearchResp, err := ParsePlexSearchResponse(ctx, searchResp.RawResponse)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to parse Plex search response with custom parser")

		// Log the raw response for further debugging
		if searchResp.RawResponse != nil {
			log.Debug().
				Interface("rawResponse", searchResp.RawResponse).
				Msg("Raw response that failed to parse")
		}

		// Fall back to the old method
		log.Info().Msg("Falling back to loading all TV shows and filtering locally")
		return c.fallbackSearchSeries(ctx, options)
	}

	// Process the parsed response
	if plexSearchResp != nil && len(plexSearchResp.MediaContainer.Hub) > 0 {
		// Log the hub information
		for i, hub := range plexSearchResp.MediaContainer.Hub {
			log.Debug().
				Int("hubIndex", i).
				Str("hubType", hub.Type).
				Str("hubTitle", hub.Title).
				Int("itemCount", len(hub.Metadata)).
				Msg("Found search hub")

			// Process TV show items in this hub
			if hub.Type == "show" || hub.Type == "2" {
				log.Debug().
					Int("seriesCount", len(hub.Metadata)).
					Msg("Processing TV show results")

				for _, item := range hub.Metadata {
					if item.RatingKey == "" || item.Title == "" {
						continue
					}

					// Create series object
					details := &types.MediaDetails{
						Title: item.Title,
					}

					year := item.GetYear()
					if year > 0 {
						details.ReleaseYear = year
					}

					seriesItem := &types.Series{
						Details: details,
					}

					// Create media item
					mediaItem, err := GetMediaItem[*types.Series](ctx, c, seriesItem, item.RatingKey)
					if err != nil {
						log.Warn().
							Err(err).
							Str("seriesID", item.RatingKey).
							Str("seriesTitle", item.Title).
							Msg("Error creating media item for TV show")
						continue
					}

					series = append(series, mediaItem)
				}
			}
		}
	} else {
		log.Debug().
			Msg("No search hubs or results found in Plex response")

		// Check the structure of the response for debugging
		if plexSearchResp == nil {
			log.Debug().Msg("Parsed Response is nil")
		} else if len(plexSearchResp.MediaContainer.Hub) == 0 {
			log.Debug().Msg("Hub array is empty")
		}

		// Fall back to the old method
		log.Info().Msg("Falling back to loading all TV shows and filtering locally")
		return c.fallbackSearchSeries(ctx, options)
	}

	log.Info().
		Int("seriesResults", len(series)).
		Str("query", options.Query).
		Msg("Successfully searched for TV shows in Plex")

	return series, nil
}

// fallbackSearchSeries is a fallback method that loads all TV series and filters them locally
func (c *PlexClient) fallbackSearchSeries(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Using fallback search method for TV shows")

	// Use the GetSeries method, which is known to work
	searchOptions := &types.QueryOptions{
		Limit:  options.Limit,
		Offset: options.Offset,
	}

	// Get all series and filter them on our side
	seriesItems, err := c.GetSeries(ctx, searchOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}

	// If no search query, return all series
	if options.Query == "" {
		return seriesItems, nil
	}

	// Filter series by title matching the query
	searchLower := strings.ToLower(options.Query)
	var filteredSeries []*models.MediaItem[*types.Series]

	for _, series := range seriesItems {
		if series.Data == nil || series.Data.Details == nil {
			continue
		}

		seriesTitle := strings.ToLower(series.Data.Details.Title)
		if strings.Contains(seriesTitle, searchLower) {
			filteredSeries = append(filteredSeries, series)
		}
	}

	log.Debug().
		Int("searchResults", len(filteredSeries)).
		Int("totalSeries", len(seriesItems)).
		Str("query", options.Query).
		Msg("Filtered series search results (fallback method)")

	return filteredSeries, nil
}

// searchMusic searches for artists, albums, and tracks by title using the Plex search API
func (c *PlexClient) searchMusic(ctx context.Context, options *types.QueryOptions) (
	[]*models.MediaItem[*types.Artist],
	[]*models.MediaItem[*types.Album],
	[]*models.MediaItem[*types.Track],
	error) {
	log := logger.LoggerFromContext(ctx)

	var artists []*models.MediaItem[*types.Artist]
	var albums []*models.MediaItem[*types.Album]
	var tracks []*models.MediaItem[*types.Track]

	// If the client doesn't support music, return empty results
	if !c.SupportsMusic() {
		return artists, albums, tracks, nil
	}

	// First, find the music library section
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to find music section key")
		return artists, albums, tracks, fmt.Errorf("failed to find music section: %w", err)
	}

	if musicSectionKey == "" {
		log.Info().Msg("No music library found in Plex")
		return artists, albums, tracks, nil
	}

	sectionKeyFloat, err := strconv.ParseFloat(musicSectionKey, 64)
	if err != nil {
		log.Error().
			Err(err).
			Str("sectionKey", musicSectionKey).
			Msg("Invalid section key format")
		return artists, albums, tracks, fmt.Errorf("invalid section key: %w", err)
	}

	// Set limit if specified
	var limitFloat *float64
	if options.Limit > 0 {
		limit := float64(options.Limit)
		limitFloat = &limit
	}

	// Use the PerformSearch API to search for music
	log.Debug().
		Str("query", options.Query).
		Float64("sectionID", sectionKeyFloat).
		Msg("Performing Plex search for music")

	// Try to use the PerformSearch API with detailed logging
	searchResp, err := c.plexAPI.Search.PerformSearch(ctx, options.Query, &sectionKeyFloat, limitFloat)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", options.Query).
			Float64("sectionID", sectionKeyFloat).
			Msg("Failed to search music in Plex")

		// Log more details about the error for debugging
		log.Debug().
			Str("errorType", fmt.Sprintf("%T", err)).
			Str("errorMessage", err.Error()).
			Msg("Detailed error information")

		// Fall back to loading all and filtering if API search fails
		log.Info().Msg("Falling back to loading all music and filtering locally")
		return c.fallbackSearchMusic(ctx, options)
	}

	// Log response details
	log.Debug().
		Int("statusCode", searchResp.StatusCode).
		Msg("Plex search response received")

	// Use our custom parser to parse the raw response
	plexSearchResp, err := ParsePlexSearchResponse(ctx, searchResp.RawResponse)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to parse Plex search response with custom parser")

		// Log the raw response for further debugging
		if searchResp.RawResponse != nil {
			log.Debug().
				Interface("rawResponse", searchResp.RawResponse).
				Msg("Raw response that failed to parse")
		}

		// Fall back to the old method
		log.Info().Msg("Falling back to loading all music and filtering locally")
		return c.fallbackSearchMusic(ctx, options)
	}

	// Process the parsed response
	if plexSearchResp != nil && len(plexSearchResp.MediaContainer.Hub) > 0 {
		// Log the hub information
		for i, hub := range plexSearchResp.MediaContainer.Hub {
			log.Debug().
				Int("hubIndex", i).
				Str("hubType", hub.Type).
				Str("hubTitle", hub.Title).
				Int("itemCount", len(hub.Metadata)).
				Msg("Found search hub")

			// Process the appropriate music hub type
			switch hub.Type {
			case "artist", "8":
				// Process artist results
				log.Debug().
					Int("artistCount", len(hub.Metadata)).
					Msg("Processing artist results")

				for _, item := range hub.Metadata {
					if item.RatingKey == "" || item.Title == "" {
						continue
					}

					// Create artist object
					details := &types.MediaDetails{
						Title: item.Title,
					}

					artist := &types.Artist{
						Details: details,
					}

					// Create media item
					mediaItem, err := GetMediaItem[*types.Artist](ctx, c, artist, item.RatingKey)
					if err != nil {
						log.Warn().
							Err(err).
							Str("artistID", item.RatingKey).
							Str("artistName", item.Title).
							Msg("Error creating media item for artist")
						continue
					}

					artists = append(artists, mediaItem)
				}

			case "album", "9":
				// Process album results
				log.Debug().
					Int("albumCount", len(hub.Metadata)).
					Msg("Processing album results")

				for _, item := range hub.Metadata {
					if item.RatingKey == "" || item.Title == "" {
						continue
					}

					// Create album object
					details := &types.MediaDetails{
						Title: item.Title,
					}

					year := item.GetYear()
					if year > 0 {
						details.ReleaseYear = year
					}

					album := &types.Album{
						Details: details,
					}

					// Set artist name if available
					if item.ParentTitle != "" {
						album.ArtistName = item.ParentTitle
					}

					// Create media item
					mediaItem, err := GetMediaItem[*types.Album](ctx, c, album, item.RatingKey)
					if err != nil {
						log.Warn().
							Err(err).
							Str("albumID", item.RatingKey).
							Str("albumTitle", item.Title).
							Msg("Error creating media item for album")
						continue
					}

					albums = append(albums, mediaItem)
				}

			case "track", "10":
				// Process track results
				log.Debug().
					Int("trackCount", len(hub.Metadata)).
					Msg("Processing track results")

				for _, item := range hub.Metadata {
					if item.RatingKey == "" || item.Title == "" {
						continue
					}

					// Create track object
					details := &types.MediaDetails{
						Title: item.Title,
					}

					track := &types.Track{
						Details: details,
					}

					// Set artist and album names if available
					if item.GrandparentTitle != "" {
						track.ArtistName = item.GrandparentTitle
					}

					if item.ParentTitle != "" {
						track.AlbumName = item.ParentTitle
					}

					// Create media item
					mediaItem, err := GetMediaItem[*types.Track](ctx, c, track, item.RatingKey)
					if err != nil {
						log.Warn().
							Err(err).
							Str("trackID", item.RatingKey).
							Str("trackTitle", item.Title).
							Msg("Error creating media item for track")
						continue
					}

					tracks = append(tracks, mediaItem)
				}
			}
		}
	} else {
		log.Debug().
			Msg("No search hubs or results found in Plex response")

		// Check the structure of the response for debugging
		if plexSearchResp == nil {
			log.Debug().Msg("Parsed Response is nil")
		} else if len(plexSearchResp.MediaContainer.Hub) == 0 {
			log.Debug().Msg("Hub array is empty")
		}

		// Fall back to the old method
		log.Info().Msg("Falling back to loading all music and filtering locally")
		return c.fallbackSearchMusic(ctx, options)
	}

	log.Info().
		Int("artistResults", len(artists)).
		Int("albumResults", len(albums)).
		Int("trackResults", len(tracks)).
		Str("query", options.Query).
		Msg("Successfully searched for music in Plex")

	return artists, albums, tracks, nil
}

// fallbackSearchMusic is a fallback method that loads all music and filters them locally
func (c *PlexClient) fallbackSearchMusic(ctx context.Context, options *types.QueryOptions) (
	[]*models.MediaItem[*types.Artist],
	[]*models.MediaItem[*types.Album],
	[]*models.MediaItem[*types.Track],
	error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Using fallback search method for music")

	var artists []*models.MediaItem[*types.Artist]
	var albums []*models.MediaItem[*types.Album]
	var tracks []*models.MediaItem[*types.Track]

	// Get music tracks using the standard method
	searchOptions := &types.QueryOptions{
		Limit:  options.Limit,
		Offset: options.Offset,
	}

	// Try to get tracks (which should include artists and albums)
	musicItems, err := c.GetMusicTracks(ctx, searchOptions)
	if err != nil {
		return artists, albums, tracks, fmt.Errorf("failed to get music: %w", err)
	}

	// If no search query, return all tracks
	if options.Query == "" {
		return artists, albums, musicItems, nil
	}

	// Filter tracks by title matching the query
	searchLower := strings.ToLower(options.Query)
	var filteredTracks []*models.MediaItem[*types.Track]

	for _, track := range musicItems {
		if track.Data == nil || track.Data.Details == nil {
			continue
		}

		trackTitle := strings.ToLower(track.Data.Details.Title)
		artistName := strings.ToLower(track.Data.ArtistName)
		albumName := strings.ToLower(track.Data.AlbumName)

		if strings.Contains(trackTitle, searchLower) ||
			strings.Contains(artistName, searchLower) ||
			strings.Contains(albumName, searchLower) {
			filteredTracks = append(filteredTracks, track)
		}
	}

	log.Debug().
		Int("searchResults", len(filteredTracks)).
		Int("totalTracks", len(musicItems)).
		Str("query", options.Query).
		Msg("Filtered music search results (fallback method)")

	return artists, albums, filteredTracks, nil
}

// Helper function to parse ints safely
func parseInt(s string) (int, error) {
	// Use a placeholder if the input is empty
	if strings.TrimSpace(s) == "" {
		return 0, fmt.Errorf("empty string cannot be converted to int")
	}

	// Try to convert to integer
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
