// handlers/core_playlists.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	"suasor/types/responses"
	"suasor/utils"
)

// CorePlaylistHandler handles operations for playlists in the database
type CorePlaylistHandler struct {
	playlistService services.CoreMediaItemService[*mediatypes.Playlist]
	coreService     services.PlaylistService
}

// NewCorePlaylistHandler creates a new core playlist handler
func NewCorePlaylistHandler(
	playlistService services.CoreMediaItemService[*mediatypes.Playlist],
	coreService services.PlaylistService,
) *CorePlaylistHandler {
	return &CorePlaylistHandler{
		playlistService: playlistService,
		coreService:     coreService,
	}
}

// GetAll godoc
// @Summary Get all playlists
// @Description Retrieves all playlists in the database
// @Tags playlists
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of playlists to return (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Playlist]] "Playlists retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists [get]
func (h *CorePlaylistHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Debug().Msg("Getting all playlists")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	// Get all playlists
	playlists, err := h.playlistService.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve playlists")
		responses.RespondInternalError(c, err, "Failed to retrieve playlists")
		return
	}

	log.Info().
		Int("count", len(playlists)).
		Msg("Playlists retrieved successfully")
	responses.RespondOK(c, playlists, "Playlists retrieved successfully")
}

// GetByID godoc
// @Summary Get playlist by ID
// @Description Retrieves a specific playlist by ID
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Playlist]] "Playlist retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id} [get]
func (h *CorePlaylistHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting playlist by ID")

	playlist, err := h.playlistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "Playlist not found")
		return
	}

	log.Info().
		Uint64("id", id).
		Msg("Playlist retrieved successfully")
	responses.RespondOK(c, playlist, "Playlist retrieved successfully")
}

// GetPlaylistTracks godoc
// @Summary Get tracks in a playlist
// @Description Retrieves all tracks in a specific playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id}/tracks [get]
func (h *CorePlaylistHandler) GetPlaylistTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting tracks for playlist")

	playlist, err := h.playlistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "Playlist not found")
		return
	}

	items := playlist.Data.Items

	log.Info().
		Uint64("id", id).
		Int("itemCount", len(items)).
		Msg("Playlist tracks retrieved successfully")
	responses.RespondOK(c, items, "Items retrieved successfully")
}

// GetByGenre godoc
// @Summary Get playlists by genre
// @Description Retrieves playlists that match a specific genre
// @Tags playlists
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Playlist]] "Playlists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/genre/{genre} [get]
func (h *CorePlaylistHandler) GetByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	genre := c.Param("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	log.Debug().
		Str("genre", genre).
		Msg("Getting playlists by genre")

	// Create query options
	options := mediatypes.QueryOptions{
		Genre:     genre,
		MediaType: mediatypes.MediaTypePlaylist,
	}

	// Search playlists by genre
	playlists, err := h.playlistService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to retrieve playlists by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve playlists")
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", len(playlists)).
		Msg("Playlists by genre retrieved successfully")
	responses.RespondOK(c, playlists, "Playlists retrieved successfully")
}

// Search godoc
// @Summary Search playlists
// @Description Searches for playlists that match the query
// @Tags playlists
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Playlist]] "Playlists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/search [get]
func (h *CorePlaylistHandler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	log.Debug().
		Str("query", query).
		Msg("Searching playlists")

	// Create query options
	options := mediatypes.QueryOptions{
		Query:     query,
		MediaType: mediatypes.MediaTypePlaylist,
	}

	// Search playlists
	playlists, err := h.playlistService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Msg("Failed to search playlists")
		responses.RespondInternalError(c, err, "Failed to search playlists")
		return
	}

	log.Info().
		Str("query", query).
		Int("count", len(playlists)).
		Msg("Playlists search completed successfully")
	responses.RespondOK(c, playlists, "Playlists retrieved successfully")
}
