package radarr

import (
	"context"
	"fmt"
	"time"

	radarr "github.com/devopsarr/radarr-go/radarr"
	"strconv"
	"suasor/client/automation/interfaces"
	"suasor/utils"
)

// Configuration holds Radarr connection settings
type Configuration struct {
	BaseURL string
	APIKey  string
}

// RadarrClient implements the AutomationProvider interface
type RadarrClient struct {
	interfaces.BaseAutomationTool
	client *radarr.APIClient
	config Configuration
}

// NewRadarrClient creates a new Radarr client instance
func NewRadarrClient(ctx context.Context, clientID uint32, config any) (interfaces.AutomationProvider, error) {
	// Extract config
	cfg, ok := config.(Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Radarr client")
	}

	// Create API client configuration
	apiConfig := radarr.NewConfiguration()
	apiConfig.AddDefaultHeader("X-Api-Key", cfg.APIKey)
	apiConfig.Servers = radarr.ServerConfigurations{
		{
			URL: cfg.BaseURL,
		},
	}

	client := radarr.NewAPIClient(apiConfig)

	radarrClient := &RadarrClient{
		BaseAutomationTool: interfaces.BaseAutomationTool{
			ClientID:   clientID,
			ClientType: interfaces.ClientTypeRadarr,
			URL:        cfg.BaseURL,
			APIKey:     cfg.APIKey,
		},
		client: client,
		config: cfg,
	}

	return radarrClient, nil
}

// Register the provider factory
func init() {
	interfaces.RegisterAutomationProvider(interfaces.ClientTypeRadarr, NewRadarrClient)
}

// Capability methods
func (r *RadarrClient) SupportsMovies() bool  { return true }
func (r *RadarrClient) SupportsTVShows() bool { return false }
func (r *RadarrClient) SupportsMusic() bool   { return false }

// GetSystemStatus retrieves system information from Radarr
func (r *RadarrClient) GetSystemStatus(ctx context.Context) (interfaces.SystemStatus, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("baseURL", r.URL).
		Msg("Retrieving system status from Radarr server")

	// Call the Radarr API
	log.Debug().Msg("Making API request to Radarr server for system status")

	statusResult, resp, err := r.client.SystemAPI.GetSystemStatus(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", r.URL).
			Str("apiEndpoint", "/system/status").
			Int("statusCode", 0).
			Msg("Failed to fetch system status from Radarr")
		return interfaces.SystemStatus{}, fmt.Errorf("failed to fetch system status: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("version", statusResult.GetVersion()).
		Msg("Successfully retrieved system status from Radarr")

	// Convert to our internal type
	status := interfaces.SystemStatus{
		Version:     statusResult.GetVersion(),
		StartupPath: statusResult.GetStartupPath(),
		AppData:     statusResult.GetAppData(),
		OsName:      statusResult.GetOsName(),
		Branch:      statusResult.GetBranch(),
	}

	return status, nil
}

// GetLibraryItems retrieves all movies from Radarr
func (r *RadarrClient) GetLibraryItems(ctx context.Context, options *interfaces.LibraryQueryOptions) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("baseURL", r.URL).
		Msg("Retrieving library items from Radarr server")

	// Call the Radarr API
	log.Debug().Msg("Making API request to Radarr server for movie library")

	moviesResult, resp, err := r.client.MovieAPI.ListMovie(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", r.URL).
			Str("apiEndpoint", "/movie").
			Int("statusCode", 0).
			Msg("Failed to fetch movies from Radarr")
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("movieCount", len(moviesResult)).
		Msg("Successfully retrieved movies from Radarr")

	// Apply paging if options provided
	var start, end int
	if options != nil {
		if options.Offset > 0 {
			start = options.Offset
		}
		if options.Limit > 0 {
			end = start + options.Limit
			if end > len(moviesResult) {
				end = len(moviesResult)
			}
		} else {
			end = len(moviesResult)
		}
	} else {
		end = len(moviesResult)
	}

	// Ensure valid slice bounds
	if start >= len(moviesResult) {
		start = 0
		end = 0
	}

	// Apply paging
	var pagedMovies []radarr.MovieResource
	if start < end {
		pagedMovies = moviesResult[start:end]
	} else {
		pagedMovies = []radarr.MovieResource{}
	}

	// Convert to our internal type
	mediaItems := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(pagedMovies))
	for _, movie := range pagedMovies {
		mediaItem := r.convertMovieToMediaItem(&movie)
		mediaItems = append(mediaItems, mediaItem)
	}

	log.Info().
		Int("itemsReturned", len(mediaItems)).
		Msg("Completed GetLibraryItems request")

	return mediaItems, nil
}

// Helper function to convert Radarr movie to generic MediaItem
func (r *RadarrClient) convertMovieToMediaItem(movie *radarr.MovieResource) interfaces.AutomationMediaItem[interfaces.AutomationData] {
	// Convert images
	images := make([]interfaces.AutomationMediaImage, 0, len(movie.GetImages()))
	for _, img := range movie.GetImages() {
		images = append(images, interfaces.AutomationMediaImage{
			URL:       img.GetRemoteUrl(),
			CoverType: string(img.GetCoverType()),
		})
	}

	// Get quality profile name
	qualityProfile := interfaces.QualityProfileSummary{
		ID:   int64(movie.GetQualityProfileId()),
		Name: "", // We don't have the name in the movie object
	}

	status := interfaces.DOWNLOADEDSTATUS_NONE
	if movie.GetHasFile() {
		status = interfaces.DOWNLOADEDSTATUS_COMPLETE
	}

	return interfaces.AutomationMediaItem[interfaces.AutomationData]{
		ID:               uint64(movie.GetId()),
		Title:            movie.GetTitle(),
		Overview:         movie.GetOverview(),
		MediaType:        "movie",
		AddedAt:          movie.GetAdded(),
		Status:           interfaces.GetStatusFromMovieStatus(movie.GetStatus()),
		Path:             movie.GetPath(),
		QualityProfile:   qualityProfile,
		Images:           images,
		DownloadedStatus: status,
		Monitored:        movie.GetMonitored(),
		Data: interfaces.AutomationMovie{
			ReleaseDate: movie.GetReleaseDate(),
			Year:        movie.GetYear(),
		},
	}
}

// // Helper function to convert []int32 to []int64
// func convertInt32SliceToInt64(in []int32) []int64 {
// 	out := make([]int64, len(in))
// 	for i, v := range in {
// 		out[i] = int64(v)
// 	}
// 	return out
// }

// GetMediaByID retrieves a specific movie by ID
func (r *RadarrClient) GetMediaByID(ctx context.Context, id int64) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {

	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Int64("movieID", id).
		Str("baseURL", r.URL).
		Msg("Retrieving specific movie from Radarr server")

	// Call the Radarr API
	log.Debug().
		Int64("movieID", id).
		Msg("Making API request to Radarr server")

	movie, resp, err := r.client.MovieAPI.GetMovieById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", r.URL).
			Str("apiEndpoint", fmt.Sprintf("/movie/%d", id)).
			Int("statusCode", 0).
			Msg("Failed to fetch movie from Radarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to fetch movie: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int64("movieID", id).
		Str("movieTitle", movie.GetTitle()).
		Msg("Successfully retrieved movie from Radarr")

	// Convert to our internal type
	mediaItem := r.convertMovieToMediaItem(movie)

	log.Debug().
		Int64("movieID", id).
		Str("movieTitle", mediaItem.Title).
		Int32("year", mediaItem.Year).
		Msg("Successfully returned movie data")

	return mediaItem, nil
}

// AddMedia adds a new movie to Radarr
func (r *RadarrClient) AddMedia(ctx context.Context, item interfaces.AutomationMediaAddRequest) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("title", item.Title).
		Msg("Adding movie to Radarr")

	// Create new movie resource
	newMovie := radarr.NewMovieResource()
	newMovie.SetTitle(item.Title)
	newMovie.SetQualityProfileId(int32(item.QualityProfileID))
	newMovie.SetTmdbId(int32(item.TMDBID))
	newMovie.SetYear(int32(item.Year))
	newMovie.SetMonitored(item.Monitored)
	newMovie.SetRootFolderPath(item.Path)
	newMovie.SetTags(item.Tags)

	// Set minimum availability if provided
	// if item.MinimumAvailability != "" {
	// 	newMovie.SetMinimumAvailability(item.MinimumAvailability)
	// }

	// Make API request
	result, resp, err := r.client.MovieAPI.CreateMovie(ctx).MovieResource(*newMovie).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", r.URL).
			Str("title", item.Title).
			Msg("Failed to add movie to Radarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to add movie: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("movieID", result.GetId()).
		Str("title", result.GetTitle()).
		Msg("Successfully added movie to Radarr")

	return r.convertMovieToMediaItem(result), nil
}

// UpdateMedia updates an existing movie in Radarr
func (r *RadarrClient) UpdateMedia(ctx context.Context, id int64, item interfaces.AutomationMediaUpdateRequest) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Int64("movieID", id).
		Msg("Updating movie in Radarr")

	// First get the existing movie
	existingMovie, resp, err := r.client.MovieAPI.GetMovieById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("movieID", id).
			Msg("Failed to fetch movie for update")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to fetch movie for update: %w", err)
	}

	// Update fields as needed
	existingMovie.SetMonitored(item.Monitored)

	if item.QualityProfileID > 0 {
		existingMovie.SetQualityProfileId(int32(item.QualityProfileID))
	}

	if item.Path != "" {
		existingMovie.SetPath(item.Path)
	}

	if item.Tags != nil {
		existingMovie.SetTags(convertInt64SliceToInt32(item.Tags))
	}

	stringId := strconv.FormatInt(id, 10)

	// Send update request
	updatedMovie, resp, err := r.client.MovieAPI.UpdateMovie(ctx, stringId).MovieResource(*existingMovie).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("movieID", id).
			Msg("Failed to update movie in Radarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to update movie: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("movieID", updatedMovie.GetId()).
		Msg("Successfully updated movie in Radarr")

	return r.convertMovieToMediaItem(updatedMovie), nil
}

// DeleteMedia removes a movie from Radarr
func (r *RadarrClient) DeleteMedia(ctx context.Context, id int64) error {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Int64("movieID", id).
		Msg("Deleting movie from Radarr")

	// Optional deletion flags
	deleteFiles := false
	addExclusion := false

	resp, err := r.client.MovieAPI.DeleteMovie(ctx, int32(id)).
		DeleteFiles(deleteFiles).
		AddImportExclusion(addExclusion).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Int64("movieID", id).
			Msg("Failed to delete movie from Radarr")
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int64("movieID", id).
		Msg("Successfully deleted movie from Radarr")

	return nil
}

// SearchMedia searches for movies in Radarr
func (r *RadarrClient) SearchMedia(ctx context.Context, query string, options *interfaces.SearchOptions) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("query", query).
		Msg("Searching for movies in Radarr")

	searchResult, resp, err := r.client.MovieLookupAPI.ListMovieLookup(ctx).Term(query).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("query", query).
			Msg("Failed to search for movies in Radarr")
		return nil, fmt.Errorf("failed to search for movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("resultCount", len(searchResult)).
		Msg("Successfully searched for movies in Radarr")

	// Convert results to MediaItems
	mediaItems := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(searchResult))
	for _, movie := range searchResult {
		mediaItem := r.convertMovieToMediaItem(&movie)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}

// GetQualityProfiles retrieves available quality profiles from Radarr
func (r *RadarrClient) GetQualityProfiles(ctx context.Context) ([]interfaces.QualityProfile, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Msg("Retrieving quality profiles from Radarr")

	profiles, resp, err := r.client.QualityProfileAPI.ListQualityProfile(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch quality profiles from Radarr")
		return nil, fmt.Errorf("failed to fetch quality profiles: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("profileCount", len(profiles)).
		Msg("Successfully retrieved quality profiles from Radarr")

	// Convert to our internal representation
	result := make([]interfaces.QualityProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, interfaces.QualityProfile{
			ID:   int64(profile.GetId()),
			Name: profile.GetName(),
		})
	}

	return result, nil
}

// GetTags retrieves all tags from Radarr
func (r *RadarrClient) GetTags(ctx context.Context) ([]interfaces.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Msg("Retrieving tags from Radarr")

	tags, resp, err := r.client.TagAPI.ListTag(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch tags from Radarr")
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("tagCount", len(tags)).
		Msg("Successfully retrieved tags from Radarr")

	// Convert to our internal representation
	result := make([]interfaces.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, interfaces.Tag{
			ID:   int64(tag.GetId()),
			Name: tag.GetLabel(),
		})
	}

	return result, nil
}

// CreateTag creates a new tag in Radarr
func (r *RadarrClient) CreateTag(ctx context.Context, tagName string) (interfaces.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("tagName", tagName).
		Msg("Creating new tag in Radarr")

	newTag := radarr.NewTagResource()
	newTag.SetLabel(tagName)

	createdTag, resp, err := r.client.TagAPI.CreateTag(ctx).TagResource(*newTag).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("tagName", tagName).
			Msg("Failed to create tag in Radarr")
		return interfaces.Tag{}, fmt.Errorf("failed to create tag: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("tagID", createdTag.GetId()).
		Str("tagName", createdTag.GetLabel()).
		Msg("Successfully created tag in Radarr")

	return interfaces.Tag{
		ID:   int64(createdTag.GetId()),
		Name: createdTag.GetLabel(),
	}, nil
}

// GetCalendar retrieves upcoming releases from Radarr
func (r *RadarrClient) GetCalendar(ctx context.Context, start, end time.Time) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Time("startDate", start).
		Time("endDate", end).
		Msg("Retrieving calendar from Radarr")

	// Format dates as required by Radarr API
	// startStr := start.Format(time.RFC3339)
	// endStr := end.Format(time.RFC3339)

	calendar, resp, err := r.client.CalendarAPI.ListCalendar(ctx).
		Start(start).
		End(end).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Time("startDate", start).
			Time("endDate", end).
			Msg("Failed to fetch calendar from Radarr")
		return nil, fmt.Errorf("failed to fetch calendar: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("itemCount", len(calendar)).
		Msg("Successfully retrieved calendar from Radarr")

	// Convert to our internal representation
	result := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(calendar))
	for _, item := range calendar {

		status := interfaces.DOWNLOADEDSTATUS_NONE
		if item.GetHasFile() {
			status = interfaces.DOWNLOADEDSTATUS_COMPLETE
		}

		result = append(result, interfaces.AutomationMediaItem[interfaces.AutomationData]{
			ID:               uint64(item.GetId()),
			ClientID:         r.ClientID,
			ClientType:       r.ClientType,
			Title:            item.GetTitle(),
			MediaType:        "movie",
			Status:           interfaces.GetStatusFromMovieStatus(item.GetStatus()),
			Overview:         item.GetOverview(),
			Year:             item.GetYear(),
			Monitored:        item.GetMonitored(),
			DownloadedStatus: status,
			Data: interfaces.AutomationMovie{
				ReleaseDate: item.GetPhysicalRelease(),
			},
		})
	}

	return result, nil
}

// ExecuteCommand executes system commands in Radarr
func (r *RadarrClient) ExecuteCommand(ctx context.Context, command interfaces.Command) (interfaces.CommandResult, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("commandName", command.Name).
		Msg("Executing command in Radarr")

	// Create command
	newCommand := radarr.NewCommandResource()
	newCommand.SetName(command.Name)

	// Add command-specific parameters
	// if command.Parameters != nil {
	// 	// Convert map to appropriate format if needed
	// 	switch command.Name {
	// 	case "MoviesSearch":
	// 		if movieIds, ok := command.Parameters["movieIds"].([]int64); ok {
	// 			int32Ids := make([]int32, len(movieIds))
	// 			for i, id := range movieIds {
	// 				int32Ids[i] = int32(id)
	// 			}
	// 			newCommand.SetMovieIds(int32Ids)
	// 		}
	// 		// Add other command-specific parameter handling as needed
	// 	}
	// }

	// Execute command
	cmdResult, resp, err := r.client.CommandAPI.CreateCommand(ctx).CommandResource(*newCommand).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("commandName", command.Name).
			Msg("Failed to execute command in Radarr")
		return interfaces.CommandResult{}, fmt.Errorf("failed to execute command: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("commandId", cmdResult.GetId()).
		Str("commandName", cmdResult.GetName()).
		Str("status", string(cmdResult.GetStatus())).
		Msg("Successfully initiated command in Radarr")

	return interfaces.CommandResult{
		ID:        int64(cmdResult.GetId()),
		Name:      cmdResult.GetName(),
		Status:    string(cmdResult.GetStatus()),
		StartedAt: cmdResult.GetStarted(),
	}, nil
}

// Helper function to convert []int64 to []int32
func convertInt64SliceToInt32(in []int64) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}

func (r *RadarrClient) GetMetadataProfiles(ctx context.Context) ([]interfaces.MetadataProfile, error) {
	return nil, interfaces.ErrAutomationFeatureNotSupported
}
