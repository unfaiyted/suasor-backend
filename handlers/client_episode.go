// handlers/client_episode.go
package handlers

import (
	"strconv"
	clienttypes "suasor/clients/types"
	"suasor/services"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

// ClientEpisodeHandler defines handler methods for TV episodes
type ClientEpisodeHandler[T clienttypes.ClientMediaConfig] interface {
	GetEpisodeByID(c *gin.Context)
	GetEpisodeByExternalID(c *gin.Context)
	RecordEpisodePlay(c *gin.Context)
	GetEpisodePlaybackState(c *gin.Context)
	UpdateEpisodePlaybackState(c *gin.Context)
}

// clientEpisodeHandler handles episode-related operations for media clients
type clientEpisodeHandler[T clienttypes.ClientMediaConfig] struct {
	seriesService services.ClientSeriesService[T]
}

// NewClientEpisodeHandler creates a new media client episode handler
func NewClientEpisodeHandler[T clienttypes.ClientMediaConfig](
	seriesService services.ClientSeriesService[T],
) ClientEpisodeHandler[T] {
	return &clientEpisodeHandler[T]{
		seriesService: seriesService,
	}
}

// GetEpisodeByID godoc
//
//	@Summary		Get episode by ID
//	@Description	Retrieves a specific TV episode by ID
//	@Tags			episodes, client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int											true	"Client ID"
//	@Param			episodeID	path		string										true	"Episode ID"
//	@Success		200			{object}	responses.APIResponse[types.Episode]	"Episode retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized"
//	@Failure		404			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Episode not found"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/client/{clientID}/media/episode/{episodeID} [get]
func (h *clientEpisodeHandler[T]) GetEpisodeByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting episode by ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access episode without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse client ID from URL
	clientIDStr := c.Param("clientID")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", clientIDStr).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	episodeID := c.Param("episodeID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("episodeID", episodeID).
		Msg("Retrieving episode by ID")

	episode, err := h.seriesService.GetEpisodeByID(ctx, clientID, episodeID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("episodeID", episodeID).
			Msg("Failed to retrieve episode")
		responses.RespondInternalError(c, err, "Failed to retrieve episode")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("episodeID", episodeID).
		Str("episodeTitle", episode.Data.Details.Title).
		Msg("Episode retrieved successfully")
	responses.RespondOK(c, episode, "Episode retrieved successfully")
}

// GetEpisodeByExternalID godoc
//
//	@Summary		Get episode by external ID
//	@Description	Retrieves a specific TV episode by external ID (e.g., TMDB ID)
//	@Tags			episodes, client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int											true	"Client ID"
//	@Param			source		query		string										true	"External ID source (e.g., 'tmdb')"
//	@Param			id			query		string										true	"External ID value"
//	@Success		200			{object}	responses.APIResponse[types.Episode]	"Episode retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid parameters"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized"
//	@Failure		404			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Episode not found"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/client/{clientID}/media/episode/external [get]
func (h *clientEpisodeHandler[T]) GetEpisodeByExternalID(c *gin.Context) {
	// This would be implemented with the search functionality based on external IDs
	// Currently not implemented in the base service
	responses.RespondNotImplemented(c, nil, "Get episode by external ID not implemented")
}

// RecordEpisodePlay godoc
//
//	@Summary		Record episode play
//	@Description	Records a play event for a TV episode
//	@Tags			episodes, client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int											true	"Client ID"
//	@Param			episodeID	path		string										true	"Episode ID"
//	@Success		200			{object}	responses.APIResponse[nil]				"Play recorded"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid parameters"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/client/{clientID}/media/episode/{episodeID}/play [post]
func (h *clientEpisodeHandler[T]) RecordEpisodePlay(c *gin.Context) {
	// This would record a play event for the episode
	// Currently not implemented in the base service
	responses.RespondNotImplemented(c, nil, "Record episode play not implemented")
}

// GetEpisodePlaybackState godoc
//
//	@Summary		Get episode playback state
//	@Description	Retrieves the current playback state for a TV episode
//	@Tags			episodes, client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int											true	"Client ID"
//	@Param			episodeID	path		string										true	"Episode ID"
//	@Success		200			{object}	responses.APIResponse[types.PlaybackState]	"Playback state retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid parameters"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized"
//	@Failure		404			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Episode not found"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/client/{clientID}/media/episode/{episodeID}/state [get]
func (h *clientEpisodeHandler[T]) GetEpisodePlaybackState(c *gin.Context) {
	// This would retrieve the playback state for the episode
	// Currently not implemented in the base service
	responses.RespondNotImplemented(c, nil, "Get episode playback state not implemented")
}

// UpdateEpisodePlaybackState godoc
//
//	@Summary		Update episode playback state
//	@Description	Updates the playback state for a TV episode
//	@Tags			episodes, client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int											true	"Client ID"
//	@Param			episodeID	path		string										true	"Episode ID"
//	@Param			state		body		types.PlaybackState						true	"Playback state"
//	@Success		200			{object}	responses.APIResponse[nil]				"Playback state updated"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid parameters"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized"
//	@Failure		404			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Episode not found"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/client/{clientID}/media/episode/{episodeID}/state [put]
func (h *clientEpisodeHandler[T]) UpdateEpisodePlaybackState(c *gin.Context) {
	// This would update the playback state for the episode
	// Currently not implemented in the base service
	responses.RespondNotImplemented(c, nil, "Update episode playback state not implemented")
}

