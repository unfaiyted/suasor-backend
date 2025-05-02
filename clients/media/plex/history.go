package plex

import (
	"context"
	"fmt"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils"

	"github.com/LukeHagar/plexgo/models/operations"
	"suasor/utils/logger"
)

// GetWatchHistory retrieves watch history from Plex
func (c *PlexClient) GetPlayHistory(ctx context.Context, options *types.QueryOptions) (*models.MediaItemDataList, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)
	mediaItemDataList := models.NewMediaItemDataList()

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving watch history from Plex server")

	log.Warn().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Watch history retrieval not yet implemented for Plex")

		// type GetAllLibrariesResponse struct {
		// 	// HTTP response content type for this operation
		// 	ContentType string
		// 	// HTTP response status code for this operation
		// 	StatusCode int
		// 	// Raw HTTP response; suitable for custom response parsing
		// 	RawResponse *http.Response
		// 	// The libraries available on the Server
		// 	Object *GetAllLibrariesResponseBody
		// }

	sort := "viewedAt:desc"

	username := c.config.GetUsername()
	var accountID int64
	// Get all library section IDs
	libraries, err := c.plexAPI.Library.GetAllLibraries(ctx)
	userResponse, err := c.plexAPI.Users.GetUsers(ctx, operations.GetUsersRequest{})
	filter := &operations.QueryParamFilter{}

	// accounts.Body // []byte

	plexUserResponse, err := ParsePlexUsersResponse(userResponse)
	if err != nil {
		return mediaItemDataList, fmt.Errorf("failed to get library sections: %w", err)
	}
	users := plexUserResponse.MediaContainer.Users

	// find user by username
	for _, user := range users {
		if user.Username == username {
			accountID = user.ID
			break
		}
	}

	// Do a history search for each library section based on id.
	for _, library := range libraries.Object.MediaContainer.Directory {
		// Skip if library is not a show
		if library.Type != "show" {
			continue
		}
		librarySectionID := utils.GetInt64(library.Key)

		log.Debug().
			Str("libraryTitle", library.Title).
			Str("libraryKey", library.Key).
			Msg("Searching for watch history for library")
		// Process each item
		// func (s *Sessions) GetSessionHistory(ctx context.Context, sort *string, accountID *int64, filter *operations.QueryParamFilter, librarySectionID *int64, opts ...operations.Option) (*operations.GetSessionHistoryResponse, error) {
		historyResponse, err := c.plexAPI.Sessions.GetSessionHistory(ctx, &sort, &accountID, filter, &librarySectionID)
		if err != nil {
			log.Error().
				Err(err).
				Str("libraryTitle", library.Title).
				Str("libraryKey", library.Key).
				Msg("Failed to get watch history for library")
			continue
		}
		sessionHistory := historyResponse.Object.MediaContainer.Metadata
		for _, session := range sessionHistory {
			log.Debug().
				Str("sessionKey", *session.Key).
				Str("sessionTitle", *session.Title).
				Msg("Found session in watch history")
			// Process each session
		}
		log.Info().
			Str("libraryTitle", library.Title).
			Str("libraryKey", library.Key).
			Msg("Successfully retrieved watch history for library")
	}
	// This would require querying Plex for watch history
	return mediaItemDataList, fmt.Errorf("watch history not yet implemented for Plex")
}
