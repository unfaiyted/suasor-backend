package lidarr

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"encoding/json"
	lidarr "github.com/devopsarr/lidarr-go/lidarr"
	"suasor/client/automation/types"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/utils"
)

func (l *LidarrClient) GetLibraryItems(ctx context.Context, options *types.LibraryQueryOptions) ([]models.AutomationMediaItem[types.AutomationData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Msg("Retrieving library items from Lidarr server")

	// Call the Lidarr API
	log.Debug().Msg("Making API request to Lidarr server for artist library")

	artistsResult, resp, err := l.client.ArtistAPI.ListArtist(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
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
	mediaItems := make([]models.AutomationMediaItem[types.AutomationData], 0, len(pagedArtists))
	for _, artist := range pagedArtists {
		mediaItem := l.convertArtistToMediaItem(&artist)
		mediaItems = append(mediaItems, mediaItem)
	}

	log.Info().
		Int("itemsReturned", len(mediaItems)).
		Msg("Completed GetLibraryItems request")

	return mediaItems, nil
}

// GetMediaByID retrieves a specific artist by ID

func (l *LidarrClient) GetMediaByID(ctx context.Context, id int64) (models.AutomationMediaItem[types.AutomationData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Int64("artistID", id).
		Msg("Retrieving specific artist from Lidarr server")

	// Call the Lidarr API
	log.Debug().
		Int64("artistID", id).
		Msg("Making API request to Lidarr server")

	artist, resp, err := l.client.ArtistAPI.GetArtistById(ctx, int32(id)).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("apiEndpoint", fmt.Sprintf("/artist/%d", id)).
			Int("statusCode", 0).
			Msg("Failed to fetch artist from Lidarr")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to fetch artist: %w", err)
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

func (l *LidarrClient) AddMedia(ctx context.Context, req requests.AutomationMediaAddRequest) (models.AutomationMediaItem[types.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("title", req.Title).
		Msg("Adding artist to Lidarr")

	// Create new artist resource
	newArtist := lidarr.NewArtistResource()
	newArtist.SetArtistName(req.Title)
	newArtist.SetQualityProfileId(int32(req.QualityProfileID))
	newArtist.SetForeignArtistId(req.MusicBrainzID)
	newArtist.SetMonitored(req.Monitored)
	newArtist.SetRootFolderPath(req.Path)
	newArtist.SetTags(req.Tags)

	// Set metadata profile if provided
	if req.MetadataProfileID > 0 {
		newArtist.SetMetadataProfileId(int32(req.MetadataProfileID))
	}

	// Set add options
	options := lidarr.NewAddArtistOptions()
	options.SetSearchForMissingAlbums(req.SearchForMedia)
	newArtist.SetAddOptions(*options)

	// Make API request
	result, resp, err := l.client.ArtistAPI.CreateArtist(ctx).ArtistResource(*newArtist).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("title", req.Title).
			Msg("Failed to add artist to Lidarr")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to add artist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("artistID", result.GetId()).
		Str("title", result.GetArtistName()).
		Msg("Successfully added artist to Lidarr")

	return l.convertArtistToMediaItem(result), nil
}

func (l *LidarrClient) UpdateMedia(ctx context.Context, id int64, req requests.AutomationMediaUpdateRequest) (models.AutomationMediaItem[types.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
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
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to fetch artist for update: %w", err)
	}

	// Update fields as needed
	existingArtist.SetMonitored(req.Monitored)

	if req.QualityProfileID > 0 {
		existingArtist.SetQualityProfileId(int32(req.QualityProfileID))
	}

	if req.MetadataProfileID > 0 {
		existingArtist.SetMetadataProfileId(req.MetadataProfileID)
	}

	if req.Path != "" {
		existingArtist.SetPath(req.Path)
	}

	if req.Tags != nil {
		existingArtist.SetTags(convertInt64SliceToInt32(req.Tags))
	}

	stringId := strconv.FormatInt(id, 10)

	// Send update request
	updatedArtist, resp, err := l.client.ArtistAPI.UpdateArtist(ctx, stringId).ArtistResource(*existingArtist).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Int64("artistID", id).
			Msg("Failed to update artist in Lidarr")
		return models.AutomationMediaItem[types.AutomationData]{}, fmt.Errorf("failed to update artist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("artistID", updatedArtist.GetId()).
		Msg("Successfully updated artist in Lidarr")

	return l.convertArtistToMediaItem(updatedArtist), nil
}

func (l *LidarrClient) DeleteMedia(ctx context.Context, id int64) error {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
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

func (l *LidarrClient) SearchMedia(ctx context.Context, query string, options *types.SearchOptions) ([]models.AutomationMediaItem[types.AutomationData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
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
	mediaItems := make([]models.AutomationMediaItem[types.AutomationData], 0, len(artistResults))
	for _, artist := range artistResults {
		mediaItem := l.convertArtistToMediaItem(&artist)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}
