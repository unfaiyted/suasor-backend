// playlists.go
package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/utils"
)

// GetPlaylists retrieves playlists from the Emby server
func (e *EmbyClient) GetPlaylists(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Playlist], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving playlists from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Playlist"),
		Recursive:        optional.NewBool(true),
	}

	applyQueryOptions(&queryParams, options)

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Msg("Failed to fetch playlists from Emby")
		return nil, fmt.Errorf("failed to fetch playlists: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved playlists from Emby")

	playlists := make([]types.MediaItem[types.Playlist], 0)
	for _, item := range items.Items {
		if item.Type_ == "Playlist" {
			playlist, err := e.convertToPlaylist(&item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("playlistID", item.Id).
					Str("playlistName", item.Name).
					Msg("Error converting Emby item to playlist format")
				continue
			}
			playlists = append(playlists, playlist)
		}
	}

	log.Info().
		Int("playlistsReturned", len(playlists)).
		Msg("Completed GetPlaylists request")

	return playlists, nil
}
