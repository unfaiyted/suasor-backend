package jellyfin

import (
	"context"
	"fmt"

	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
)

func (j *JellyfinClient) GetPlaylists(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[t.Playlist], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving playlists from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_PLAYLIST}
	recursive := true

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for playlists")

	itemsRequest := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(recursive).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder)

	result, resp, err := itemsRequest.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch playlists from Jellyfin")
		return nil, fmt.Errorf("failed to fetch playlists: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved playlists from Jellyfin")

	// Convert results to expected format
	playlists := make([]models.MediaItem[t.Playlist], 0)
	for _, item := range result.Items {
		if *item.Type == "Playlist" {
			playlist := models.MediaItem[t.Playlist]{
				Data: t.Playlist{
					Details: t.MediaMetadata{
						Title:       *item.Name.Get(),
						Description: *item.Overview.Get(),
						Artwork:     j.getArtworkURLs(&item),
					},
					ItemCount: int(*item.ChildCount.Get()),
					IsPublic:  true, // Assume public by default in Jellyfin
				},
				Type: "playlist",
			}
			playlist.SetClientInfo(j.ClientID, j.ClientType, *item.Id)
			playlists = append(playlists, playlist)
		}
	}

	log.Info().
		Int("playlistsReturned", len(playlists)).
		Msg("Completed GetPlaylists request")

	return playlists, nil
}
