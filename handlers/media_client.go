package handlers

import (
	"strconv"
	"suasor/models"
	"suasor/services"
	"suasor/utils"

	"github.com/gin-gonic/gin"
)

// MediaClientHandler handles media client API endpoints
type MediaClientHandler struct {
	service services.MediaClientService
}

// NewMediaClientHandler creates a new media client handler
func NewMediaClientHandler(service services.MediaClientService) *MediaClientHandler {
	return &MediaClientHandler{
		service: service,
	}
}

// CreateClient godoc
// @Summary Create a new media client
// @Description Creates a new media client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.MediaClientRequest true "Media client data"
// @Success 201 {object} models.APIResponse[models.MediaClientResponse] "Media client created"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/media [post]
// @Example request - Plex client
//
//	{
//	  "name": "My Plex Server",
//	  "clientType": "plex",
//	  "client": {
//	    "enabled": true,
//	    "host": "192.168.1.100",
//	    "port": 32400,
//	    "token": "your-plex-token",
//	    "ssl": false
//	  }
//	}
func (h *MediaClientHandler) CreateClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	var req models.MediaClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("name", req.Name).
		Str("type", string(req.ClientType)).
		Msg("Creating new media client")

	client, err := h.service.CreateClient(ctx, uid, req)
	if err != nil {
		utils.RespondInternalError(c, err, err.Error())
		return
	}

	utils.RespondCreated(c, client, "Media client created successfully")
}

// GetClient godoc
// @Summary Get media client
// @Description Retrieves a specific media client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Success 200 {object} models.APIResponse[models.MediaClientResponse] "Media client retrieved"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 404 {object} models.ErrorResponse[error] "Client not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/media/{id} [get]
// @Example response
// {
//   "data": {
//     "id": 1,
//     "userId": 123,
//     "name": "My Plex Server",
//     "clientType": "plex",
//     "client": {
//       "enabled": true,
//       "host": "192.168.1.100",
//       "port": 32400,
//       "token": "your-plex-token",
//       "ssl": false
//     },
//     "createdAt": "2023-01-01T12:00:00Z",
//     "updatedAt": "2023-01-01T12:00:00Z"
//   },
//   "message": "Media client retrieved successfully"
// }

func (h *MediaClientHandler) GetClient(c *gin.Context) {
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
		Msg("Retrieving media client")

	client, err := h.service.GetClientByID(ctx, uid, clientID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "media client not found" {
			utils.RespondNotFound(c, err, "Media client not found")
			return
		}
		utils.RespondInternalError(c, err, "Failed to retrieve media client")
		return
	}

	utils.RespondOK(c, client, "Media client retrieved successfully")
}

// GetAllClients godoc
// @Summary Get all media clients
// @Description Retrieves all media client configurations for the user
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse[[]models.MediaClient[models.ClientConfig]] "Media clients retrieved"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/media [get]
func (h *MediaClientHandler) GetAllClients(c *gin.Context) {
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
		Msg("Retrieving all media clients")

	clients, err := h.service.GetClientsByUserID(ctx, uid)
	if err != nil {
		utils.RespondInternalError(c, err, "Failed to retrieve media clients")
		return
	}

	utils.RespondOK(c, clients, "Media clients retrieved successfully")
}

// UpdateClient godoc
// @Summary Update media client
// @Description Updates an existing media client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Param request body models.MediaClientRequest true "Updated client data"
// @Success 200 {object} models.APIResponse[models.MediaClientResponse] "Media client updated"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request or client ID"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 404 {object} models.ErrorResponse[error] "Client not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/media/{id} [put]
// @Example request - Jellyfin client
// {
//   "name": "My Jellyfin Server",
//   "clientType": "jellyfin",
//   "client": {
//     "enabled": true,
//     "host": "192.168.1.101",
//     "port": 8096,
//     "apiKey": "your-jellyfin-apikey",
//     "username": "admin",
//     "ssl": false
//   }
// }

func (h *MediaClientHandler) UpdateClient(c *gin.Context) {
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

	var req models.MediaClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("name", req.Name).
		Str("type", string(req.ClientType)).
		Msg("Updating media client")

	client, err := h.service.UpdateClient(ctx, uid, clientID, req)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "media client not found" {
			utils.RespondNotFound(c, err, "Media client not found")
			return
		}
		utils.RespondInternalError(c, err, "Failed to update media client")
		return
	}

	utils.RespondOK(c, client, "Media client updated successfully")
}

// DeleteClient godoc
// @Summary Delete media client
// @Description Deletes a media client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Success 200 {object} models.APIResponse[models.EmptyResponse] "Media client deleted"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 404 {object} models.ErrorResponse[error] "Client not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/media/{id} [delete]
func (h *MediaClientHandler) DeleteClient(c *gin.Context) {
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
		Msg("Deleting media client")

	err = h.service.DeleteClient(ctx, uid, clientID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "media client not found" {
			utils.RespondNotFound(c, err, "Media client not found")
			return
		}
		utils.RespondInternalError(c, err, "Failed to delete media client")
		return
	}

	utils.RespondOK(c, models.EmptyResponse{Success: true}, "Media client deleted successfully")
}

// TestConnection godoc
// @Summary Test media client connection
// @Description Tests the connection to a media client using the provided configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.MediaClientTestRequest true "Client configuration to test"
// @Success 200 {object} models.APIResponse[models.MediaClientTestResponse] "Connection test result"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /clients/media/test [post]
// @Example request - Emby client test
//
//	{
//	  "url": "http://192.168.1.102:8096",
//	  "clientType": "emby",
//	  "client": {
//	    "apiKey": "your-emby-apikey",
//	    "username": "admin"
//	  }
//	}
//
// @Example response
//
//	{
//	  "data": {
//	    "success": true,
//	    "message": "Successfully connected to Emby server",
//	    "version": "4.7.0"
//	  },
//	  "message": "Connection test completed"
//	}
func (h *MediaClientHandler) TestConnection(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	var req models.MediaClientTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("type", string(req.ClientType)).
		Msg("Testing media client connection")

	result, err := h.service.TestClientConnection(ctx, req)
	if err != nil {
		utils.RespondInternalError(c, err, "Failed to test media client connection")
		return
	}

	utils.RespondOK(c, result, "Connection test completed")
}
