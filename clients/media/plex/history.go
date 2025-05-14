package plex

import (
	"context"
	"fmt"
	"suasor/clients/media/types"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils"
	"time"

	"github.com/unfaiyted/plexgo/models/operations"
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

	// Check if config is nil before accessing methods
	if c.config == nil {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Plex config is nil, cannot retrieve username for history sync")
		return mediaItemDataList, fmt.Errorf("plex config is nil, username not available")
	}

	username := c.config.GetUsername()
	if username == "" {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Plex username is not set in configuration")
		return mediaItemDataList, fmt.Errorf("plex username not configured")
	}

	var accountID int64
	// Get all library section IDs
	libraries, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		return mediaItemDataList, fmt.Errorf("failed to get libraries: %w", err)
	}

	clientName := "Suasor"
	opts := operations.GetUsersRequest{
		XPlexToken: c.config.GetToken(),
		ClientName: &clientName,
	}

	userResponse, err := c.plexAPI.Users.GetUsers(ctx, opts)
	if err != nil {
		return mediaItemDataList, fmt.Errorf("failed to get users: %w", err)
	}

	filter := &operations.QueryParamFilter{}

	plexUserResponse, err := ParsePlexUsersResponse(ctx, userResponse)
	if err != nil {
		return mediaItemDataList, fmt.Errorf("failed to get library sections: %w", err)
	}
	users := plexUserResponse.MediaContainer.Users

	log.Debug().
		Int("userCount", len(users)).
		Msg("Found users")

	// find user by username
	for _, user := range users {
		log.Debug().
			Str("username", user.Username).
			Int64("userID", user.ID).
			Msg("Checking username")
		if user.Username == username {
			accountID = user.ID
			break
		}
	}

	// Use sort for history items
	sort := "viewedAt:desc"

	// Do a history search for each library section based on id.
	for _, library := range libraries.Object.MediaContainer.Directory {
		librarySectionID := utils.GetInt64(library.Key)

		log.Debug().
			Str("libraryTitle", library.Title).
			Int64("accountID", accountID).
			Str("libraryKey", library.Key).
			Msg("Searching for watch history for library")

		// Get session history for this library section
		historyResponse, err := c.plexAPI.Sessions.GetSessionHistory(ctx, &sort, &accountID, filter, &librarySectionID)
		if err != nil {
			log.Error().
				Err(err).
				Str("libraryTitle", library.Title).
				Str("libraryKey", library.Key).
				Msg("Failed to get watch history for library")
			continue
		}

		// Process each history item based on library type
		sessionHistory := historyResponse.Object.MediaContainer.Metadata
		for _, session := range sessionHistory {
			if session.Key == nil || session.Title == nil {
				continue
			}

			// Determine the media type
			var itemType string
			if session.Type != nil {
				itemType = *session.Type
				log.Debug().
					Str("sessionKey", *session.Key).
					Str("sessionTitle", *session.Title).
					Str("type", itemType).
					Msg("Found session in watch history")
			} else {
				log.Debug().
					Str("sessionKey", *session.Key).
					Str("sessionTitle", *session.Title).
					Msg("Found session in watch history (unknown type)")
				continue
			}

			// Get last viewed timestamp
			var viewedAt time.Time
			if session.ViewedAt != nil {
				viewedAt = time.Unix(int64(*session.ViewedAt), 0)
			}

			// Create appropriate data item based on media type
			switch itemType {
			case "movie":
				// Create a movie with basic details
				movie := &mediatypes.Movie{
					Details: &mediatypes.MediaDetails{
						Title: *session.Title,
					},
				}

				// Add year if available
				// if session.Year != nil {
				// 	movie.Details.ReleaseYear = int(utils.GetInt64(*session.Year))
				// }
				//
				// // Add summary if available
				// if session.Summary != nil {
				// 	movie.Details.Description = *session.Summary
				// }

				// Create a MediaItem for this movie
				mediaItem := models.NewMediaItem(movie)
				mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *session.Key)

				// Create a history data item for this movie
				historyItem := models.NewUserMediaItemData(mediaItem, c.GetClientID())
				historyItem.PlayedAt = viewedAt
				historyItem.PlayedPercentage = 100 // Assume completed

				// Add to the history list
				mediaItemDataList.Movies[*session.Key] = historyItem

			case "episode":
				// Create an episode with basic details
				episode := &mediatypes.Episode{
					Details: &mediatypes.MediaDetails{
						Title: *session.Title,
					},
				}

				// Create a MediaItem for this episode
				mediaItem := models.NewMediaItem(episode)
				mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *session.Key)

				// Create a history data item for this episode
				historyItem := models.NewUserMediaItemData(mediaItem, c.GetClientID())
				historyItem.PlayedAt = viewedAt
				historyItem.PlayedPercentage = 100 // Assume completed

				// Add to the history list
				mediaItemDataList.Episodes[*session.Key] = historyItem

			case "track":
				// Create a track with basic details
				track := &mediatypes.Track{
					Details: &mediatypes.MediaDetails{
						Title: *session.Title,
					},
				}

				// Create a MediaItem for this track
				mediaItem := models.NewMediaItem(track)
				mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *session.Key)

				// Create a history data item for this track
				historyItem := models.NewUserMediaItemData(mediaItem, c.GetClientID())
				historyItem.PlayedAt = viewedAt
				historyItem.PlayedPercentage = 100 // Assume completed

				// Add to the history list
				mediaItemDataList.Tracks[*session.Key] = historyItem
			}
		}

		log.Info().
			Str("libraryTitle", library.Title).
			Str("libraryKey", library.Key).
			Msg("Successfully retrieved watch history for library")
	}

	// Return the populated history list (or an error message if it's empty)
	if mediaItemDataList.GetTotalItems() == 0 {
		log.Warn().Msg("No history items found in any Plex library")
		return mediaItemDataList, fmt.Errorf("no watch history found for Plex user %s", username)
	}

	log.Info().
		Int("totalItems", mediaItemDataList.GetTotalItems()).
		Msg("Successfully retrieved Plex watch history")

	return mediaItemDataList, nil
}
