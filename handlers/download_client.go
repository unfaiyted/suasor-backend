// handlers/download_client.go
package handlers

import (
	"strconv"
	"suasor/models"
	"suasor/services"
	"suasor/utils"

	"github.com/gin-gonic/gin"
)

// DownloadClientHandler handles download client API endpoints
type DownloadClientHandler struct {
	service services.DownloadClientService
}

// NewDownloadClientHandler creates a new download client handler
func NewDownloadClientHandler(service services.DownloadClientService) *DownloadClientHandler {
	return &DownloadClientHandler{
		service: service,
	}
}

// CreateClient godoc
// @Summary Create a new download client
// @Description Creates a new download client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.DownloadClientRequest true "Download client data"
// @Success 201 {object} models.APIResponse[models.DownloadClient] "Download client created"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/download [post]
func (h *DownloadClientHandler) CreateClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	var req models.DownloadClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("name", req.Name).
		Str("type", string(req.ClientType)).
		Msg("Creating new download client")

	client, err := h.service.CreateClient(ctx, uid, req)
	if err != nil {
		utils.RespondInternalError(c, err, err.Error())
		return
	}

	utils.RespondCreated(c, client, "Download client created successfully")
}

// GetClient godoc
// @Summary Get download client
// @Description Retrieves a specific download client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Success 200 {object} models.APIResponse[models.DownloadClient] "Download client retrieved"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 404 {object} models.ErrorResponse[error] "Client not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/download/{id} [get]
func (h *DownloadClientHandler) GetClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("id")).Msg("Invalid client ID format")
		utils.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Retrieving download client")

	client, err := h.service.GetClientByID(ctx, uid, clientID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "download client not found" {
			utils.RespondNotFound(c, err, "Download client not found")
			return
		}
		utils.RespondInternalError(c, err, "Failed to retrieve download client")
		return
	}

	utils.RespondOK(c, client, "Download client retrieved successfully")
}

// GetAllClients godoc
// @Summary Get all download clients
// @Description Retrieves all download client configurations for the user
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse[[]models.DownloadClient] "Download clients retrieved"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/download [get]
func (h *DownloadClientHandler) GetAllClients(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	log.Info().
		Uint64("userID", uid).
		Msg("Retrieving all download clients")

	clients, err := h.service.GetClientsByUserID(ctx, uid)
	if err != nil {
		utils.RespondInternalError(c, err, "Failed to retrieve download clients")
		return
	}

	utils.RespondOK(c, clients, "Download clients retrieved successfully")
}

// UpdateClient godoc
// @Summary Update download client
// @Description Updates an existing download client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Param request body models.DownloadClientRequest true "Updated client data"
// @Success 200 {object} models.APIResponse[models.DownloadClient] "Download client updated"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request or client ID"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 404 {object} models.ErrorResponse[error] "Client not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/download/{id} [put]
func (h *DownloadClientHandler) UpdateClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("id")).Msg("Invalid client ID format")
		utils.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req models.DownloadClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("name", req.Name).
		Str("type", string(req.ClientType)).
		Msg("Updating download client")

	client, err := h.service.UpdateClient(ctx, uid, clientID, req)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "download client not found" {
			utils.RespondNotFound(c, err, "Download client not found")
			return
		}
		utils.RespondInternalError(c, err, "Failed to update download client")
		return
	}

	utils.RespondOK(c, client, "Download client updated successfully")
}

// DeleteClient godoc
// @Summary Delete download client
// @Description Deletes a download client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Success 200 {object} models.APIResponse[any] "Download client deleted"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 404 {object} models.ErrorResponse[error] "Client not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/download/{id} [delete]
func (h *DownloadClientHandler) DeleteClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("id")).Msg("Invalid client ID format")
		utils.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Deleting download client")

	err = h.service.DeleteClient(ctx, uid, clientID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "download client not found" {
			utils.RespondNotFound(c, err, "Download client not found")
			return
		}
		utils.RespondInternalError(c, err, "Failed to delete download client")
		return
	}

	utils.RespondOK(c, models.EmptyResponse{Success: true}, "Download client deleted successfully")
}

// TestConnection godoc
// @Summary Test download client connection
// @Description Tests the connection to a download client using the provided configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ClientTestRequest true "Client configuration to test"
// @Success 200 {object} models.APIResponse[models.ClientTestResponse] "Connection test result"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/download/test [post]
func (h *DownloadClientHandler) TestConnection(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	log.Info().
		Str("userID", userID.(string)).
		Msg("Testing download client connection")

	uid := userID.(uint64)

	var req models.ClientTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("type", string(req.ClientType)).
		Msg("Testing download client connection")

	result, err := h.service.TestClientConnection(ctx, req)
	if err != nil {
		utils.RespondInternalError(c, err, "Failed to test download client connection")
		return
	}

	utils.RespondOK(c, result, "Connection test completed")
}
