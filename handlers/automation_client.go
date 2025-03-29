// handlers/automation_client.go
package handlers

import (
	"strconv"
	automationtypes "suasor/client/automation/types"
	"suasor/services"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// AutomationClientHandler handles automation client API endpoints
type AutomationClientHandler struct {
	service services.AutomationClientService
}

// NewAutomationClientHandler creates a new automation client handler
func NewAutomationClientHandler(service services.AutomationClientService) *AutomationClientHandler {
	return &AutomationClientHandler{
		service: service,
	}
}

// GetSystemStatus godoc
// @Summary Get automation client system status
// @Description Retrieves system status information from the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Success 200 {object} responses.APIResponse "System status retrieved"
// @Failure 400 {object} responses.ErrorResponse "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/status [get]
func (h *AutomationClientHandler) GetSystemStatus(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Retrieving system status from automation client")

	status, err := h.service.GetSystemStatus(ctx, uid, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve system status")
		responses.RespondInternalError(c, err, "Failed to retrieve system status")
		return
	}

	responses.RespondOK(c, status, "System status retrieved successfully")
}

// GetLibraryItems godoc
// @Summary Get library items from automation client
// @Description Retrieves all library items from the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param sortBy query string false "Sort by"
// @Param sortOrder query string false "Sort order"
// @Success 200 {object} responses.APIResponse "Library items retrieved"
// @Failure 400 {object} responses.ErrorResponse "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/library [get]
func (h *AutomationClientHandler) GetLibraryItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	offset, err := strconv.Atoi(c.Query("offset"))
	sortBy := c.Query("sortBy")
	// sortOrder := c.Query("sortOrder")
	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Retrieving library items from automation client")

	opts := automationtypes.LibraryQueryOptions{
		Limit:  limit,
		Offset: offset,
	}

	if sortBy != "" {
		opts.SortBy = sortBy
	}

	items, err := h.service.GetLibraryItems(ctx, uid, clientID, &opts)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve library items")
		responses.RespondInternalError(c, err, "Failed to retrieve library items")
		return
	}

	responses.RespondOK(c, items, "Library items retrieved successfully")
}

// GetMediaByID godoc
// @Summary Get media by ID from automation client
// @Description Retrieves a specific media item from the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param mediaID path string true "Media ID"
// @Success 200 {object} responses.APIResponse "Media retrieved"
// @Failure 400 {object} responses.ErrorResponse "Invalid client or media ID"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/media/{mediaID} [get]
func (h *AutomationClientHandler) GetMediaByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	mediaID := c.Param("mediaID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("mediaID", mediaID).
		Msg("Retrieving media by ID from automation client")

	media, err := h.service.GetMediaByID(ctx, uid, clientID, mediaID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("mediaID", mediaID).
			Msg("Failed to retrieve media")
		responses.RespondInternalError(c, err, "Failed to retrieve media")
		return
	}

	responses.RespondOK(c, media, "Media retrieved successfully")
}

// AddMedia godoc
// @Summary Add media to automation client
// @Description Adds a new media item to the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param request body requests.AddMediaRequest true "Media details"
// @Success 201 {object} responses.APIResponse "Media added"
// @Failure 400 {object} responses.ErrorResponse "Invalid request"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/media [post]
func (h *AutomationClientHandler) AddMedia(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req requests.AutomationMediaAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Interface("request", req).
		Msg("Adding media to automation client")

	result, err := h.service.AddMedia(ctx, uid, clientID, req)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Interface("request", req).
			Msg("Failed to add media")
		responses.RespondInternalError(c, err, "Failed to add media")
		return
	}

	responses.RespondCreated(c, result, "Media added successfully")
}

// UpdateMedia godoc
// @Summary Update media in automation client
// @Description Updates an existing media item in the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param mediaID path string true "Media ID"
// @Param request body requests.UpdateMediaRequest true "Media details"
// @Success 200 {object} responses.APIResponse "Media updated"
// @Failure 400 {object} responses.ErrorResponse "Invalid request"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/media/{mediaID} [put]
func (h *AutomationClientHandler) UpdateMedia(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	mediaID := c.Param("mediaID")

	var req requests.AutomationMediaUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("mediaID", mediaID).
		Interface("request", req).
		Msg("Updating media in automation client")

	result, err := h.service.UpdateMedia(ctx, uid, clientID, mediaID, req)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("mediaID", mediaID).
			Interface("request", req).
			Msg("Failed to update media")
		responses.RespondInternalError(c, err, "Failed to update media")
		return
	}

	responses.RespondOK(c, result, "Media updated successfully")
}

// DeleteMedia godoc
// @Summary Delete media from automation client
// @Description Deletes a media item from the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param mediaID path string true "Media ID"
// @Success 200 {object} responses.APIResponse "Media deleted"
// @Failure 400 {object} responses.ErrorResponse "Invalid client or media ID"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/media/{mediaID} [delete]
func (h *AutomationClientHandler) DeleteMedia(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	mediaID := c.Param("mediaID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("mediaID", mediaID).
		Msg("Deleting media from automation client")

	err = h.service.DeleteMedia(ctx, uid, clientID, mediaID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("mediaID", mediaID).
			Msg("Failed to delete media")
		responses.RespondInternalError(c, err, "Failed to delete media")
		return
	}

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Media deleted successfully")
}

// SearchMedia godoc
// @Summary Search media in automation client
// @Description Searches for media items in the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse "Search results"
// @Failure 400 {object} responses.ErrorResponse "Invalid client ID or query"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/search [get]
func (h *AutomationClientHandler) SearchMedia(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	query := c.Query("q")
	if query == "" {
		log.Warn().Uint64("userID", uid).Msg("Empty search query provided")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("query", query).
		Msg("Searching media in automation client")

	results, err := h.service.SearchMedia(ctx, uid, clientID, query)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("query", query).
			Msg("Failed to search media")
		responses.RespondInternalError(c, err, "Failed to search media")
		return
	}

	responses.RespondOK(c, results, "Search completed successfully")
}

// GetQualityProfiles godoc
// @Summary Get quality profiles from automation client
// @Description Retrieves all quality profiles from the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Success 200 {object} responses.APIResponse "Quality profiles retrieved"
// @Failure 400 {object} responses.ErrorResponse "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/profiles/quality [get]
func (h *AutomationClientHandler) GetQualityProfiles(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Retrieving quality profiles from automation client")

	profiles, err := h.service.GetQualityProfiles(ctx, uid, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve quality profiles")
		responses.RespondInternalError(c, err, "Failed to retrieve quality profiles")
		return
	}

	responses.RespondOK(c, profiles, "Quality profiles retrieved successfully")
}

// GetMetadataProfiles godoc
// @Summary Get metadata profiles from automation client
// @Description Retrieves all metadata profiles from the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Success 200 {object} responses.APIResponse "Metadata profiles retrieved"
// @Failure 400 {object} responses.ErrorResponse "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/profiles/metadata [get]
func (h *AutomationClientHandler) GetMetadataProfiles(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Retrieving metadata profiles from automation client")

	profiles, err := h.service.GetMetadataProfiles(ctx, uid, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve metadata profiles")
		responses.RespondInternalError(c, err, "Failed to retrieve metadata profiles")
		return
	}

	responses.RespondOK(c, profiles, "Metadata profiles retrieved successfully")
}

// GetTags godoc
// @Summary Get tags from automation client
// @Description Retrieves all tags from the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Success 200 {object} responses.APIResponse "Tags retrieved"
// @Failure 400 {object} responses.ErrorResponse "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/tags [get]
func (h *AutomationClientHandler) GetTags(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Retrieving tags from automation client")

	tags, err := h.service.GetTags(ctx, uid, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve tags")
		responses.RespondInternalError(c, err, "Failed to retrieve tags")
		return
	}

	responses.RespondOK(c, tags, "Tags retrieved successfully")
}

// CreateTag godoc
// @Summary Create tag in automation client
// @Description Creates a new tag in the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param request body requests.CreateTagRequest true "Tag details"
// @Success 201 {object} responses.APIResponse "Tag created"
// @Failure 400 {object} responses.ErrorResponse "Invalid request"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/tags [post]
func (h *AutomationClientHandler) CreateTag(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req requests.AutomationCreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Interface("request", req).
		Msg("Creating tag in automation client")

	tag, err := h.service.CreateTag(ctx, uid, clientID, req)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Interface("request", req).
			Msg("Failed to create tag")
		responses.RespondInternalError(c, err, "Failed to create tag")
		return
	}

	responses.RespondCreated(c, tag, "Tag created successfully")
}

// GetCalendar godoc
// @Summary Get calendar from automation client
// @Description Retrieves calendar events from the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param start query string false "Start date (YYYY-MM-DD)"
// @Param end query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} responses.APIResponse "Calendar events retrieved"
// @Failure 400 {object} responses.ErrorResponse "Invalid client ID or dates"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/calendar [get]
func (h *AutomationClientHandler) GetCalendar(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	startDate := c.Query("start")
	endDate := c.Query("end")

	startDateParsed, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		log.Error().Err(err).Str("startDate", startDate).Msg("Invalid start date format")
		responses.RespondBadRequest(c, err, "Invalid start date")
		return
	}
	endDateParsed, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		log.Error().Err(err).Str("endDate", endDate).Msg("Invalid end date format")
		responses.RespondBadRequest(c, err, "Invalid end date")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Time("startDate", startDateParsed).
		Time("endDate", endDateParsed).
		Msg("Retrieving calendar from automation client")

	events, err := h.service.GetCalendar(ctx, uid, clientID, startDateParsed, endDateParsed)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("startDate", startDate).
			Str("endDate", endDate).
			Msg("Failed to retrieve calendar")
		responses.RespondInternalError(c, err, "Failed to retrieve calendar")
		return
	}

	responses.RespondOK(c, events, "Calendar events retrieved successfully")
}

// ExecuteCommand godoc
// @Summary Execute command on automation client
// @Description Executes a command on the automation client
// @Tags automation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param request body requests.ExecuteCommandRequest true "Command details"
// @Success 200 {object} responses.APIResponse "Command executed"
// @Failure 400 {object} responses.ErrorResponse "Invalid request"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Server error"
// @Router /clients/automation/{clientID}/command [post]
func (h *AutomationClientHandler) ExecuteCommand(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req requests.AutomationExecuteCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Interface("request", req).
		Msg("Executing command on automation client")

	result, err := h.service.ExecuteCommand(ctx, uid, clientID, req)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Interface("request", req).
			Msg("Failed to execute command")
		responses.RespondInternalError(c, err, "Failed to execute command")
		return
	}

	responses.RespondOK(c, result, "Command executed successfully")
}

