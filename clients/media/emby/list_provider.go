// list_provider.go
package emby

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"suasor/clients/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"
	"suasor/utils/logger"
)

// SupportsLists returns true if the client supports lists
func (e *EmbyClient) SupportsLists() bool {
	return true
}

// determineListType tries to identify if a list is a playlist or collection
func (e *EmbyClient) determineListType(ctx context.Context, listID string) (types.MediaType, error) {
	log := logger.LoggerFromContext(ctx)

	// Get user ID if available
	userID := e.getUserID()

	// Try to get the item details
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids: optional.NewString(listID),
	}

	if userID != "" {
		queryParams.UserId = optional.NewString(userID)
	}

	itemsResult, _, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("listID", listID).
			Msg("Failed to get list details from Emby")
		return "", fmt.Errorf("failed to get list details: %w", err)
	}

	if len(itemsResult.Items) == 0 {
		return "", fmt.Errorf("list not found: %s", listID)
	}

	item := itemsResult.Items[0]

	// Determine the type based on the Emby item type
	switch item.Type_ {
	case "Playlist":
		return types.MediaTypePlaylist, nil
	case "BoxSet":
		return types.MediaTypeCollection, nil
	default:
		return "", fmt.Errorf("unknown list type: %s", item.Type_)
	}
}

// Convert a playlist MediaItem to a generic ListData MediaItem
func convertPlaylistToListData(playlist *models.MediaItem[*types.Playlist]) *models.MediaItem[types.ListData] {
	// We'd need to implement type conversion here to convert to the interface.
	// This is a conceptual example.
	// In reality, this would require more complex type handling.
	return nil
}

// Convert a collection MediaItem to a generic ListData MediaItem
func convertCollectionToListData(collection *models.MediaItem[*types.Collection]) *models.MediaItem[types.ListData] {
	// We'd need to implement type conversion here to convert to the interface.
	// This is a conceptual example.
	// In reality, this would require more complex type handling.
	return nil
}

// Convert a slice of playlist MediaItems to a slice of generic ListData MediaItems
func convertPlaylistsToListData(playlists []*models.MediaItem[*types.Playlist]) []*models.MediaItem[types.ListData] {
	// We'd need to implement type conversion here to convert to the interface.
	// This is a conceptual example.
	// In reality, this would require more complex type handling.
	return nil
}

// Convert a slice of collection MediaItems to a slice of generic ListData MediaItems
func convertCollectionsToListData(collections []*models.MediaItem[*types.Collection]) []*models.MediaItem[types.ListData] {
	// We'd need to implement type conversion here to convert to the interface.
	// This is a conceptual example.
	// In reality, this would require more complex type handling.
	return nil
}

