// handlers/core_media_item.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"

	"suasor/clients/media/types"
	"suasor/services"
	_ "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

type CoreMediaItemHandler[T types.MediaData] interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	GetMostPlayed(c *gin.Context)
	// GetByClientItemID(c *gin.Context)
	GetByExternalID(c *gin.Context)
	Search(c *gin.Context)
	GetRecentlyAdded(c *gin.Context)
	GetByType(c *gin.Context)
	GetByPerson(c *gin.Context)
	GetByYear(c *gin.Context)
	GetLatestByAdded(c *gin.Context)
	GetByClient(c *gin.Context)
	GetByGenre(c *gin.Context)
	GetByRating(c *gin.Context)
	GetPopular(c *gin.Context)
	GetTopRated(c *gin.Context)

	GetType() string
}

// coreMediaItemHandler is a generic handler for media items in the database
// It provides basic operations that are shared across all media types
// and serves as the base for more specialized media handlers
type coreMediaItemHandler[T types.MediaData] struct {
	mediaService services.CoreMediaItemService[T]
}

// NewCoreMediaItemHandler creates a new core media item handler
func NewCoreMediaItemHandler[T types.MediaData](
	mediaService services.CoreMediaItemService[T],
) CoreMediaItemHandler[T] {
	return &coreMediaItemHandler[T]{
		mediaService: mediaService,
	}
}

// GetAll godoc
// @Summary Get all media items
// @Description Retrieves all media items of a specific type from the database
// @Tags media, core
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType} [get]
func (h *coreMediaItemHandler[T]) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	log.Debug().Msg("Getting all media items")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	// Get all media items
	items, err := h.mediaService.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Int("count", len(items)).
		Msg("Media items retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetByID godoc
// @Summary Get media item by ID
// @Description Retrieves a specific media item by ID
// @Tags media, core
// @Accept json
// @Produce json
// @Param id path int true "Media Item ID"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/{id} [get]
func (h *coreMediaItemHandler[T]) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting media item by ID")

	item, err := h.mediaService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to retrieve media item")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	log.Info().
		Uint64("id", id).
		Msg("Media item retrieved successfully")
	responses.RespondOK(c, item, "Media item retrieved successfully")
}

// GetByExternalID godoc
// @Summary Get media item by external ID
// @Description Retrieves a specific media item by its external ID from a source
// @Tags media, core
// @Accept json
// @Produce json
// @Param source path string true "Source of the external ID (e.g., tmdb, imdb)"
// @Param id path string true "External ID"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/external/{source}/{id} [get]
func (h *coreMediaItemHandler[T]) GetByExternalID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	source := c.Param("source")
	if source == "" {
		log.Warn().Msg("Source is required")
		responses.RespondBadRequest(c, nil, "Source is required")
		return
	}

	externalID := c.Param("id")
	if externalID == "" {
		log.Warn().Msg("External ID is required")
		responses.RespondBadRequest(c, nil, "External ID is required")
		return
	}

	log.Debug().
		Str("source", source).
		Str("externalID", externalID).
		Msg("Getting media item by external ID")

	item, err := h.mediaService.GetByExternalID(ctx, source, externalID)
	if err != nil {
		log.Error().Err(err).
			Str("source", source).
			Str("externalID", externalID).
			Msg("Failed to retrieve media item")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	log.Info().
		Str("source", source).
		Str("externalID", externalID).
		Msg("Media item retrieved successfully")
	responses.RespondOK(c, item, "Media item retrieved successfully")
}

// Search godoc
// @Summary Search media items
// @Description Searches for media items based on query parameters
// @Tags media, core
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/search [get]
func (h *coreMediaItemHandler[T]) Search(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Str("query", query).
		Str("type", string(mediaType)).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Searching media items")

	// Create query options
	options := types.QueryOptions{
		Query:     query,
		MediaType: mediaType,
		Limit:     limit,
		Offset:    offset,
	}

	// Search media items
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Msg("Failed to search media items")
		responses.RespondInternalError(c, err, "Failed to search media items")
		return
	}

	log.Info().
		Str("query", query).
		Int("count", len(items)).
		Msg("Media items search completed successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetRecentlyAdded godoc
// @Summary Get recently added media items
// @Description Retrieves recently added media items of a specific type
// @Tags media, core
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param days query int false "Number of days to look back (default 30)"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/recently-added [get]
func (h *coreMediaItemHandler[T]) GetRecentlyAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil {
		days = 30
	}

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Str("type", string(mediaType)).
		Int("limit", limit).
		Int("days", days).
		Msg("Getting recently added media items")

	// Get recently added media items
	items, err := h.mediaService.GetRecentItems(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve recently added media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Int("count", len(items)).
		Msg("Recently added media items retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetByType godoc
// @Summary Get media items by type
// @Description Retrieves media items of a specific type
// @Tags media, core
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType} [get]
func (h *coreMediaItemHandler[T]) GetByType(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	typeParam := c.Param("type")
	if typeParam == "" {
		log.Warn().Msg("Media type is required")
		responses.RespondBadRequest(c, nil, "Media type is required")
		return
	}

	mediaType := types.MediaType(typeParam)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	log.Debug().
		Str("type", string(mediaType)).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting media items by type")

	// Get media items by type
	items, err := h.mediaService.GetByType(ctx, mediaType)
	if err != nil {
		log.Error().Err(err).
			Str("type", string(mediaType)).
			Msg("Failed to retrieve media items by type")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Str("type", string(mediaType)).
		Int("count", len(items)).
		Msg("Media items by type retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetByPerson godoc
// @Summary Get media items by person
// @Description Retrieves media items associated with a specific person (actor, director, etc.)
// @Tags media, core
// @Accept json
// @Produce json
// @Param personID path int true "Person ID"
// @Param mediaType path string true "Media type"
// @Param role query string false "Role filter (actor, director, etc.)"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/person/{personID} [get]
func (h *coreMediaItemHandler[T]) GetByPerson(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	personID, err := strconv.ParseUint(c.Param("personID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("personID", c.Param("personID")).Msg("Invalid person ID")
		responses.RespondBadRequest(c, err, "Invalid person ID")
		return
	}

	role := c.Query("role") // Optional role filter (actor, director, etc.)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("personID", personID).
		Str("role", role).
		Int("limit", limit).
		Msg("Getting media items by person")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Create query options with person filter
	options := types.QueryOptions{
		MediaType: mediaType,
		PersonID:  personID,
		Role:      role,
		Limit:     limit,
	}

	// Use search with person filter
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("personID", personID).
			Str("role", role).
			Msg("Failed to retrieve media items by person")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Uint64("personID", personID).
		Str("role", role).
		Int("count", len(items)).
		Msg("Media items by person retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetByYear godoc
// @Summary Get media items by release year
// @Description Retrieves media items released in a specific year
// @Tags media, core
// @Accept json
// @Produce json
// @Param year path int true "Release year"
// @Param mediaType path string true "Media type"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/year/{year} [get]
func (h *coreMediaItemHandler[T]) GetByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	yearStr := c.Param("year")
	if yearStr == "" {
		log.Warn().Msg("Year is required")
		responses.RespondBadRequest(c, nil, "Year is required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Warn().Err(err).Str("year", yearStr).Msg("Invalid year format")
		responses.RespondBadRequest(c, err, "Invalid year format")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Int("year", year).
		Int("limit", limit).
		Msg("Getting media items by year")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Create query options with year filter
	options := types.QueryOptions{
		MediaType: mediaType,
		Year:      year,
		Limit:     limit,
	}

	// Search with year filter
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Int("year", year).
			Msg("Failed to retrieve media items by year")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Int("year", year).
		Int("count", len(items)).
		Msg("Media items by year retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetLatestByAdded godoc
// @Summary Get latest added media items
// @Description Retrieves the most recently added media items
// @Tags media, core
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/latest [get]
func (h *coreMediaItemHandler[T]) GetLatestByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Int("limit", limit).
		Msg("Getting latest added media items")

	// Use recent items with a short time window
	days := 90 // Last 90 days by default
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil {
			days = parsedDays
		}
	}
	cutoffDate := time.Now().AddDate(0, 0, -days)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Create query options
	options := types.QueryOptions{
		MediaType:      mediaType,
		Sort:           "created_at",
		DateAddedAfter: cutoffDate,
		SortOrder:      "desc",

		Limit: limit,
	}

	// Search with sorting by creation date
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve latest media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Int("count", len(items)).
		Msg("Latest media items retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetByClient godoc
// @Summary Get media items by client
// @Description Retrieves media items associated with a specific client
// @Tags media, core
// @Accept json
// @Produce json
// @Param clientID path int true "Client ID"
// @Param mediaType path string true "Media type"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/client/{clientID} [get]
func (h *coreMediaItemHandler[T]) GetByClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Getting media items by client")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Create query options with client filter
	options := types.QueryOptions{
		MediaType: mediaType,
		ClientID:  clientID,
		Limit:     limit,
	}

	// Search with client filter
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve media items by client")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Uint64("clientID", clientID).
		Int("count", len(items)).
		Msg("Media items by client retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetByGenre godoc
// @Summary Get media items by genre
// @Description Retrieves media items that match a specific genre
// @Tags media, core
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Param mediaType path string true "Media type"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/genre/{genre} [get]
func (h *coreMediaItemHandler[T]) GetByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	genre := c.Param("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting media items by genre")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Create query options with genre filter
	options := types.QueryOptions{
		MediaType: mediaType,
		Genre:     genre,
		Limit:     limit,
	}

	// Search with genre filter
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to retrieve media items by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", len(items)).
		Msg("Media items by genre retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetByExternalSourceID godoc
// @Summary Get media item by external source ID
// @Description Retrieves a media item using its external source ID (e.g., TMDB ID)
// @Tags media, core
// @Accept json
// @Produce json
// @Param source path string true "External source name (e.g., tmdb, imdb)"
// @Param mediaType path string true "Media type"
// @Param externalId path string true "External ID from the source"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/external/{source}/{externalId} [get]
func (h *coreMediaItemHandler[T]) GetByExternalSourceID(c *gin.Context) {
	// This function is an alias for GetByExternalID to maintain compatibility with both naming schemes
	h.GetByExternalID(c)
}

// GetPopular godoc
// @Summary Get popular media items
// @Description Retrieves popular media items based on views or ratings
// @Tags media, core
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/popular [get]
func (h *coreMediaItemHandler[T]) GetPopular(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Int("limit", limit).
		Msg("Getting popular media items")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Create query options with popularity sorting
	options := types.QueryOptions{
		MediaType: mediaType,
		Sort:      "popularity",
		SortOrder: "desc",
		Limit:     limit,
	}

	// Search with popularity sorting
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve popular media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Int("count", len(items)).
		Msg("Popular media items retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetTopRated godoc
// @Summary Get top rated media items
// @Description Retrieves the highest rated media items
// @Tags media, core
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param mediaType path string true "Media type"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/top-rated [get]
func (h *coreMediaItemHandler[T]) GetTopRated(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Int("limit", limit).
		Msg("Getting top rated media items")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Create query options with rating sorting
	options := types.QueryOptions{
		MediaType: mediaType,
		Sort:      "rating",
		SortOrder: "desc",
		Limit:     limit,
	}

	// Search with rating sorting
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve top rated media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Int("count", len(items)).
		Msg("Top rated media items retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetMostPlayed godoc
// @Summary Get most played media items
// @Description Retrieves the most played media items
// @Tags media, core
// @Accept json
// @Produce json
// @Param userID query int false "User ID"
// @Param mediaType path string true "Media type"
// @Param limit query int false "Maximum quantity items to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/most-played [get]
func (h *coreMediaItemHandler[T]) GetMostPlayed(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting most played media items")

	// Get the most played media items
	items, err := h.mediaService.GetMostPlayed(ctx, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve most played media items")
		responses.RespondInternalError(c, err, "Failed to retrieve most played media items")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(items)).
		Msg("Most played media items retrieved successfully")
	responses.RespondOK(c, items, "Most played media items retrieved successfully")
}

// GetByRating godoc
// @Summary Get media items by rating
// @Description Retrieves media items that match a specific rating
// @Tags media, core
// @Accept json
// @Produce json
// @Param rating path number true "Rating"
// @Param mediaType path string true "Media type"
// @Param limit query int false "Maximum number of media items to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.Movie]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/{mediaType}/rating/{rating} [get]
func (h *coreMediaItemHandler[T]) GetByRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	rating, err := strconv.ParseFloat(c.Param("rating"), 32)
	if err != nil {
		log.Warn().Err(err).Str("rating", c.Param("rating")).Msg("Invalid rating value")
		responses.RespondBadRequest(c, err, "Invalid rating value")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Float64("rating", rating).
		Int("limit", limit).
		Msg("Getting media items by rating")

	// Create query options
	options := types.QueryOptions{
		MinimumRating: float32(rating),
		Limit:         limit,
	}

	// Search media items by rating
	items, err := h.mediaService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Float64("rating", rating).
			Msg("Failed to retrieve media items by rating")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Float64("rating", rating).
		Int("count", len(items)).
		Msg("Media items by rating retrieved successfully")

	responses.RespondOK(c, items, "Media items retrieved successfully")
}

func (h *coreMediaItemHandler[T]) GetType() string {

	var zero T
	types := types.GetMediaTypeFromTypeName(zero)
	return string(types)
}
