package lidarr

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"encoding/json"
	lidarr "github.com/devopsarr/lidarr-go/lidarr"
	"io"
	"suasor/client/automation/interfaces"
	"suasor/utils"
)

// Configuration holds Lidarr connection settings
type Configuration struct {
	BaseURL string
	APIKey  string
}

// LidarrClient implements the AutomationProvider interface
type LidarrClient struct {
	interfaces.BaseAutomationTool
	client *lidarr.APIClient
	config Configuration
}

// getAlbumExternalIDs extracts all available external IDs from a Lidarr album resource
func (l *LidarrClient) getAlbumExternalIDs(album *lidarr.AlbumResource) []interfaces.ExternalID {
	var ids []interfaces.ExternalID

	// Add Lidarr's internal album ID
	ids = append(ids, interfaces.ExternalID{
		Source: "lidarr_album",
		Value:  strconv.FormatInt(int64(album.GetId()), 10),
	})

	// Add artist ID reference
	if album.GetArtistId() != 0 {
		ids = append(ids, interfaces.ExternalID{
			Source: "lidarr_artist",
			Value:  strconv.FormatInt(int64(album.GetArtistId()), 10),
		})
	}

	// Add MusicBrainz album ID
	if album.ForeignAlbumId.Get() != nil && *album.ForeignAlbumId.Get() != "" {
		ids = append(ids, interfaces.ExternalID{
			Source: "musicbrainz_album",
			Value:  *album.ForeignAlbumId.Get(),
		})
	}

	// Extract IDs from releases if available
	for _, release := range album.GetReleases() {
		if release.ForeignReleaseId.Get() != nil && *release.ForeignReleaseId.Get() != "" {
			ids = append(ids, interfaces.ExternalID{
				Source: "musicbrainz_release",
				Value:  *release.ForeignReleaseId.Get(),
			})
		}
	}

	// If artist object is included, extract artist IDs too
	if album.Artist != nil {
		artistIDs := l.getArtistExternalIDs(album.Artist)
		ids = append(ids, artistIDs...)
	}

	return ids
}

// getExternalIDs extracts all available external IDs from a Lidarr artist resource
func (l *LidarrClient) getArtistExternalIDs(artist *lidarr.ArtistResource) []interfaces.ExternalID {
	var ids []interfaces.ExternalID

	// Add Lidarr's internal ID
	ids = append(ids, interfaces.ExternalID{
		Source: "lidarr",
		Value:  strconv.FormatInt(int64(artist.GetId()), 10),
	})

	// Add MusicBrainz ID (the primary ID in Lidarr)
	if artist.ForeignArtistId.Get() != nil && *artist.ForeignArtistId.Get() != "" {
		ids = append(ids, interfaces.ExternalID{
			Source: "musicbrainz",
			Value:  *artist.ForeignArtistId.Get(),
		})
	}

	// Add MusicBrainz ID (alternative field)
	if artist.MbId.Get() != nil && *artist.MbId.Get() != "" {
		ids = append(ids, interfaces.ExternalID{
			Source: "musicbrainz",
			Value:  *artist.MbId.Get(),
		})
	}

	// Add The Audio DB ID
	if artist.GetTadbId() != 0 {
		ids = append(ids, interfaces.ExternalID{
			Source: "audiodb",
			Value:  strconv.FormatInt(int64(artist.GetTadbId()), 10),
		})
	}

	// Add Discogs ID
	if artist.GetDiscogsId() != 0 {
		ids = append(ids, interfaces.ExternalID{
			Source: "discogs",
			Value:  strconv.FormatInt(int64(artist.GetDiscogsId()), 10),
		})
	}

	// Add AllMusic ID
	if artist.AllMusicId.Get() != nil && *artist.AllMusicId.Get() != "" {
		ids = append(ids, interfaces.ExternalID{
			Source: "allmusic",
			Value:  *artist.AllMusicId.Get(),
		})
	}

	return ids
}

func DetermineDownloadStatus(stats lidarr.ArtistStatisticsResource) interfaces.DownloadedStatus {
	// Check if all values are properly set
	if *stats.TrackFileCount == 0 || *stats.TrackCount == 0 || *stats.TotalTrackCount == 0 {
		return interfaces.DOWNLOADEDSTATUS_NONE
	}

	trackFileCount := stats.GetTrackFileCount()
	trackCount := stats.GetTrackCount()
	totalTrackCount := stats.GetTotalTrackCount()

	// No tracks downloaded
	if trackFileCount == 0 {
		return interfaces.DOWNLOADEDSTATUS_NONE
	}

	// All tracks the artist has ever released are downloaded
	if trackFileCount == totalTrackCount && totalTrackCount > 0 {
		return interfaces.DOWNLOADEDSTATUS_COMPLETE
	}

	// All monitored/requested tracks are downloaded, but not all tracks the artist has ever released
	if trackFileCount == trackCount && trackCount > 0 && trackCount < totalTrackCount {
		return interfaces.DOWNLOADEDSTATUS_REQUESTED
	}

	// Some tracks are downloaded, but not all monitored tracks
	return interfaces.DOWNLOADEDSTATUS_NONE
}

// NewLidarrClient creates a new Lidarr client instance
func NewLidarrClient(ctx context.Context, clientID uint32, config any) (interfaces.AutomationProvider, error) {
	// Extract config
	cfg, ok := config.(Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Lidarr client")
	}

	// Create API client configuration
	apiConfig := lidarr.NewConfiguration()
	apiConfig.AddDefaultHeader("X-Api-Key", cfg.APIKey)
	apiConfig.Servers = lidarr.ServerConfigurations{
		{
			URL: cfg.BaseURL,
		},
	}

	client := lidarr.NewAPIClient(apiConfig)

	lidarrClient := &LidarrClient{
		BaseAutomationTool: interfaces.BaseAutomationTool{
			ClientID:   clientID,
			ClientType: interfaces.ClientTypeLidarr,
			URL:        cfg.BaseURL,
			APIKey:     cfg.APIKey,
		},
		client: client,
		config: cfg,
	}

	return lidarrClient, nil
}

// Register the provider factory
func init() {
	interfaces.RegisterAutomationProvider(interfaces.ClientTypeLidarr, NewLidarrClient)
}

// Capability methods
func (l *LidarrClient) SupportsMovies() bool  { return false }
func (l *LidarrClient) SupportsTVShows() bool { return false }
func (l *LidarrClient) SupportsMusic() bool   { return true }

// GetSystemStatus retrieves system information from Lidarr
func (l *LidarrClient) GetSystemStatus(ctx context.Context) (interfaces.SystemStatus, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("baseURL", l.URL).
		Msg("Retrieving system status from Lidarr server")

	// Call the Lidarr API
	log.Debug().Msg("Making API request to Lidarr server for system status")

	statusResult, resp, err := l.client.SystemAPI.GetSystemStatus(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", l.URL).
			Str("apiEndpoint", "/system/status").
			Int("statusCode", 0).
			Msg("Failed to fetch system status from Lidarr")
		return interfaces.SystemStatus{}, fmt.Errorf("failed to fetch system status: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("version", statusResult.GetVersion()).
		Msg("Successfully retrieved system status from Lidarr")

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

// GetLibraryItems retrieves all artists from Lidarr
func (l *LidarrClient) GetLibraryItems(ctx context.Context, options *interfaces.LibraryQueryOptions) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("baseURL", l.URL).
		Msg("Retrieving library items from Lidarr server")

	// Call the Lidarr API
	log.Debug().Msg("Making API request to Lidarr server for artist library")

	artistsResult, resp, err := l.client.ArtistAPI.ListArtist(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", l.URL).
			Str("apiEndpoint", "/artist").
			Int("statusCode", 0).
			Msg("Failed to fetch artists from Lidarr")
		return nil, fmt.Errorf("failed to fetch artists: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("artistCount", len(artistsResult)).
		Msg("Successfully retrieved artists from Lidarr")

	// Apply paging if options provided
	var start, end int
	if options != nil {
		if options.Offset > 0 {
			start = options.Offset
		}
		if options.Limit > 0 {
			end = start + options.Limit
			if end > len(artistsResult) {
				end = len(artistsResult)
			}
		} else {
			end = len(artistsResult)
		}
	} else {
		end = len(artistsResult)
	}

	// Ensure valid slice bounds
	if start >= len(artistsResult) {
		start = 0
		end = 0
	}

	// Apply paging
	var pagedArtists []lidarr.ArtistResource
	if start < end {
		pagedArtists = artistsResult[start:end]
	} else {
		pagedArtists = []lidarr.ArtistResource{}
	}

	// Convert to our internal type
	mediaItems := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(pagedArtists))
	for _, artist := range pagedArtists {
		mediaItem := l.convertArtistToMediaItem(&artist)
		mediaItems = append(mediaItems, mediaItem)
	}

	log.Info().
		Int("itemsReturned", len(mediaItems)).
		Msg("Completed GetLibraryItems request")

	return mediaItems, nil
}

// Helper function to convert Lidarr artist to generic MediaItem
func (l *LidarrClient) convertArtistToMediaItem(artist *lidarr.ArtistResource) interfaces.AutomationMediaItem[interfaces.AutomationData] {
	// Convert images
	images := make([]interfaces.AutomationMediaImage, 0, len(artist.GetImages()))
	for _, img := range artist.GetImages() {
		images = append(images, interfaces.AutomationMediaImage{
			URL:       img.GetRemoteUrl(),
			CoverType: string(img.GetCoverType()),
		})
	}

	// Get quality profile
	qualityProfile := interfaces.QualityProfileSummary{
		ID:   int64(artist.GetQualityProfileId()),
		Name: "", // We don't have the name in the artist object
	}

	// Get metadata profile if available
	metadataProfile := interfaces.MetadataProfile{
		ID:   artist.GetMetadataProfileId(),
		Name: "", // We don't have the name in the artist object
	}

	// Convert genres
	genres := artist.GetGenres()

	// Use start year or end year as appropriate
	// var releaseDate time.Time

	releaseDate := *artist.LastAlbum.ReleaseDate.Get()

	return interfaces.AutomationMediaItem[interfaces.AutomationData]{
		ID:        uint64(artist.GetId()),
		Title:     artist.GetArtistName(),
		Overview:  artist.GetOverview(),
		MediaType: "artist",
		// TODO: get the first album the arstist release, set a year
		// Year:             artist.GetYearFormed(),
		AddedAt:          artist.GetAdded(),
		Status:           interfaces.GetStatusFromMusicStatus(artist.GetStatus()),
		Path:             artist.GetPath(),
		QualityProfile:   qualityProfile,
		Images:           images,
		Genres:           genres,
		ExternalIDs:      l.getArtistExternalIDs(artist),
		DownloadedStatus: DetermineDownloadStatus(artist.GetStatistics()),
		Monitored:        artist.GetMonitored(),
		Data: interfaces.AutomationArtist{
			MetadataProfile:       metadataProfile,
			MostRecentReleaseDate: releaseDate,
		},
	}
}

// GetMediaByID retrieves a specific artist by ID
func (l *LidarrClient) GetMediaByID(ctx context.Context, id int64) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Int64("artistID", id).
		Str("baseURL", l.URL).
		Msg("Retrieving specific artist from Lidarr server")

	// Call the Lidarr API
	log.Debug().
		Int64("artistID", id).
		Msg("Making API request to Lidarr server")

	artist, resp, err := l.client.ArtistAPI.GetArtistById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", l.URL).
			Str("apiEndpoint", fmt.Sprintf("/artist/%d", id)).
			Int("statusCode", 0).
			Msg("Failed to fetch artist from Lidarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to fetch artist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int64("artistID", id).
		Str("artistName", artist.GetArtistName()).
		Msg("Successfully retrieved artist from Lidarr")

	// Convert to our internal type
	mediaItem := l.convertArtistToMediaItem(artist)

	log.Debug().
		Int64("artistID", id).
		Str("artistName", mediaItem.Title).
		Msg("Successfully returned artist data")

	return mediaItem, nil
}

// AddMedia adds a new artist to Lidarr
func (l *LidarrClient) AddMedia(ctx context.Context, item interfaces.AutomationMediaAddRequest) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("title", item.Title).
		Msg("Adding artist to Lidarr")

	// Create new artist resource
	newArtist := lidarr.NewArtistResource()
	newArtist.SetArtistName(item.Title)
	newArtist.SetQualityProfileId(int32(item.QualityProfileID))
	newArtist.SetForeignArtistId(item.MusicBrainzID)
	newArtist.SetMonitored(item.Monitored)
	newArtist.SetRootFolderPath(item.Path)
	newArtist.SetTags(item.Tags)

	// Set metadata profile if provided
	if item.MetadataProfileID > 0 {
		newArtist.SetMetadataProfileId(int32(item.MetadataProfileID))
	}

	// Set add options
	options := lidarr.NewAddArtistOptions()
	options.SetSearchForMissingAlbums(item.SearchForMedia)
	newArtist.SetAddOptions(*options)

	// Make API request
	result, resp, err := l.client.ArtistAPI.CreateArtist(ctx).ArtistResource(*newArtist).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", l.URL).
			Str("title", item.Title).
			Msg("Failed to add artist to Lidarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to add artist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("artistID", result.GetId()).
		Str("title", result.GetArtistName()).
		Msg("Successfully added artist to Lidarr")

	return l.convertArtistToMediaItem(result), nil
}

// UpdateMedia updates an existing artist in Lidarr
func (l *LidarrClient) UpdateMedia(ctx context.Context, id int64, item interfaces.AutomationMediaUpdateRequest) (interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Int64("artistID", id).
		Msg("Updating artist in Lidarr")

	// First get the existing artist
	existingArtist, resp, err := l.client.ArtistAPI.GetArtistById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("artistID", id).
			Msg("Failed to fetch artist for update")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to fetch artist for update: %w", err)
	}

	// Update fields as needed
	existingArtist.SetMonitored(item.Monitored)

	if item.QualityProfileID > 0 {
		existingArtist.SetQualityProfileId(int32(item.QualityProfileID))
	}

	if item.MetadataProfileID > 0 {
		existingArtist.SetMetadataProfileId(item.MetadataProfileID)
	}

	if item.Path != "" {
		existingArtist.SetPath(item.Path)
	}

	if item.Tags != nil {
		existingArtist.SetTags(convertInt64SliceToInt32(item.Tags))
	}

	stringId := strconv.FormatInt(id, 10)

	// Send update request
	updatedArtist, resp, err := l.client.ArtistAPI.UpdateArtist(ctx, stringId).ArtistResource(*existingArtist).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("artistID", id).
			Msg("Failed to update artist in Lidarr")
		return interfaces.AutomationMediaItem[interfaces.AutomationData]{}, fmt.Errorf("failed to update artist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("artistID", updatedArtist.GetId()).
		Msg("Successfully updated artist in Lidarr")

	return l.convertArtistToMediaItem(updatedArtist), nil
}

// DeleteMedia removes an artist from Lidarr
func (l *LidarrClient) DeleteMedia(ctx context.Context, id int64) error {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Int64("artistID", id).
		Msg("Deleting artist from Lidarr")

	// Optional deletion flags
	deleteFiles := false
	addExclusion := false

	resp, err := l.client.ArtistAPI.DeleteArtist(ctx, int32(id)).
		DeleteFiles(deleteFiles).
		AddImportListExclusion(addExclusion).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Int64("artistID", id).
			Msg("Failed to delete artist from Lidarr")
		return fmt.Errorf("failed to delete artist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int64("artistID", id).
		Msg("Successfully deleted artist from Lidarr")

	return nil
}

// SearchMedia searches for artists in Lidarr
func (l *LidarrClient) SearchMedia(ctx context.Context, query string, options *interfaces.SearchOptions) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("query", query).
		Msg("Searching for artists in Lidarr")

	// Call the Lidarr API - this returns http.Response directly
	httpResp, err := l.client.ArtistLookupAPI.GetArtistLookup(ctx).Term(query).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("query", query).
			Msg("Failed to search for artists in Lidarr")
		return nil, fmt.Errorf("failed to search for artists: %w", err)
	}
	defer httpResp.Body.Close()

	// Log the status code and response for debugging
	log.Debug().
		Int("statusCode", httpResp.StatusCode).
		Msg("Received response from Lidarr artist lookup")

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the raw response for debugging
	log.Debug().
		Str("responseBody", string(bodyBytes)).
		Msg("Raw response body from artist lookup")

	// Parse the response into ArtistResource objects
	var artistResults []lidarr.ArtistResource
	err = json.Unmarshal(bodyBytes, &artistResults)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to unmarshal artist lookup response")
		return nil, fmt.Errorf("failed to unmarshal artist lookup response: %w", err)
	}

	log.Info().
		Int("statusCode", httpResp.StatusCode).
		Int("resultCount", len(artistResults)).
		Msg("Successfully searched for artists in Lidarr")

	// Convert results to MediaItems
	mediaItems := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(artistResults))
	for _, artist := range artistResults {
		mediaItem := l.convertArtistToMediaItem(&artist)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}

// GetQualityProfiles retrieves available quality profiles from Lidarr
func (l *LidarrClient) GetQualityProfiles(ctx context.Context) ([]interfaces.QualityProfile, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Msg("Retrieving quality profiles from Lidarr")

	profiles, resp, err := l.client.QualityProfileAPI.ListQualityProfile(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch quality profiles from Lidarr")
		return nil, fmt.Errorf("failed to fetch quality profiles: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("profileCount", len(profiles)).
		Msg("Successfully retrieved quality profiles from Lidarr")

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

// GetMetadataProfiles retrieves available metadata profiles from Lidarr
func (l *LidarrClient) GetMetadataProfiles(ctx context.Context) ([]interfaces.MetadataProfile, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Msg("Retrieving metadata profiles from Lidarr")

	profiles, resp, err := l.client.MetadataProfileAPI.ListMetadataProfile(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch metadata profiles from Lidarr")
		return nil, fmt.Errorf("failed to fetch metadata profiles: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("profileCount", len(profiles)).
		Msg("Successfully retrieved metadata profiles from Lidarr")

	// Convert to our internal representation
	result := make([]interfaces.MetadataProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, interfaces.MetadataProfile{
			ID:   profile.GetId(),
			Name: profile.GetName(),
		})
	}

	return result, nil
}

// GetTags retrieves all tags from Lidarr
func (l *LidarrClient) GetTags(ctx context.Context) ([]interfaces.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Msg("Retrieving tags from Lidarr")

	tags, resp, err := l.client.TagAPI.ListTag(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch tags from Lidarr")
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("tagCount", len(tags)).
		Msg("Successfully retrieved tags from Lidarr")

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

// CreateTag creates a new tag in Lidarr
func (l *LidarrClient) CreateTag(ctx context.Context, tagName string) (interfaces.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("tagName", tagName).
		Msg("Creating new tag in Lidarr")

	newTag := lidarr.NewTagResource()
	newTag.SetLabel(tagName)

	createdTag, resp, err := l.client.TagAPI.CreateTag(ctx).TagResource(*newTag).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("tagName", tagName).
			Msg("Failed to create tag in Lidarr")
		return interfaces.Tag{}, fmt.Errorf("failed to create tag: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("tagID", createdTag.GetId()).
		Str("tagName", createdTag.GetLabel()).
		Msg("Successfully created tag in Lidarr")

	return interfaces.Tag{
		ID:   int64(createdTag.GetId()),
		Name: createdTag.GetLabel(),
	}, nil
}

// ExecuteCommand executes system commands in Lidarr
func (l *LidarrClient) ExecuteCommand(ctx context.Context, command interfaces.Command) (interfaces.CommandResult, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("commandName", command.Name).
		Msg("Executing command in Lidarr")

	// Create command
	newCommand := lidarr.NewCommandResource()
	newCommand.SetName(command.Name)

	// Execute command
	cmdResult, resp, err := l.client.CommandAPI.CreateCommand(ctx).CommandResource(*newCommand).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("commandName", command.Name).
			Msg("Failed to execute command in Lidarr")
		return interfaces.CommandResult{}, fmt.Errorf("failed to execute command: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("commandId", cmdResult.GetId()).
		Str("commandName", cmdResult.GetName()).
		Str("status", string(cmdResult.GetStatus())).
		Msg("Successfully initiated command in Lidarr")

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

// GetCalendar retrieves upcoming album releases from Lidarr
func (l *LidarrClient) GetCalendar(ctx context.Context, start, end time.Time) ([]interfaces.AutomationMediaItem[interfaces.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint32("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Time("startDate", start).
		Time("endDate", end).
		Msg("Retrieving calendar from Lidarr")

	calendar, resp, err := l.client.CalendarAPI.ListCalendar(ctx).
		Start(start).
		End(end).
		IncludeArtist(true).
		Execute()

	if err != nil {
		log.Error().
			Err(err).
			Time("startDate", start).
			Time("endDate", end).
			Msg("Failed to fetch calendar from Lidarr")
		return nil, fmt.Errorf("failed to fetch calendar: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("itemCount", len(calendar)).
		Msg("Successfully retrieved calendar from Lidarr")

	// Convert to our internal representation
	result := make([]interfaces.AutomationMediaItem[interfaces.AutomationData], 0, len(calendar))
	for _, item := range calendar {

		downloadStatus := interfaces.DOWNLOADEDSTATUS_NONE
		if *item.GetStatistics().TrackFileCount >= *item.GetStatistics().TotalTrackCount {
			downloadStatus = interfaces.DOWNLOADEDSTATUS_COMPLETE
		} else if *item.GetStatistics().TrackFileCount > 0 {
			downloadStatus = interfaces.DOWNLOADEDSTATUS_PARTIAL
		}

		artistID := strconv.Itoa(int(item.GetArtistId()))

		// Get base album details
		albumInfo := interfaces.AutomationMediaItem[interfaces.AutomationData]{
			ID:               uint64(item.GetId()),
			ClientID:         l.ClientID,
			ClientType:       l.ClientType,
			Title:            item.GetTitle(), // Album title
			MediaType:        "album",
			Overview:         item.GetOverview(),
			Year:             int32(item.GetReleaseDate().Year()),
			Monitored:        item.GetMonitored(),
			ExternalIDs:      l.getAlbumExternalIDs(&item),
			DownloadedStatus: downloadStatus,
			Data: interfaces.AutomationAlbum{
				ArtistName:  *item.GetArtist().ArtistName.Get(),
				ArtistID:    artistID,
				ReleaseDate: item.GetReleaseDate(),
			},
		}

		result = append(result, albumInfo)
	}

	return result, nil
}
