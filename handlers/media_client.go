package handlers

import (
	"strconv"
	"suasor/services"

	"github.com/gin-gonic/gin"
	"suasor/client/types"
	client "suasor/client/types"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils"
)

// MediaClientHandler handles media client API endpoints
type ClientHandler[T client.ClientConfig] struct {
	service services.ClientService[T]
}

// NewMediaClientHandler creates a new media client handler
func NewClientHandler[T types.ClientConfig](service services.ClientService[T]) *ClientHandler[T] {
	return &ClientHandler[T]{
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
func (h *ClientHandler[T]) CreateClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	var req requests.ClientRequest[T]
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("name", req.Name).
		Str("type", string(req.ClientType)).
		Msg("Creating new media client")

	client := models.Client[T]{
		UserID:   uid,
		Name:     req.Name,
		Category: req.ClientType.AsCategory(),
		Config:   models.ClientConfigWrapper[T]{Data: req.Client},
	}

	client, err := h.service.Create(ctx, client)
	if err != nil {
		responses.RespondInternalError(c, err, err.Error())
		return
	}

	responses.RespondCreated(c, client, "Media client created successfully")
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

func (h *ClientHandler[T]) GetClient(c *gin.Context) {
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
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("id")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Retrieving media client")

	client, err := h.service.GetByID(ctx, uid, clientID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "media client not found" {
			responses.RespondNotFound(c, err, "Media client not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to retrieve media client")
		return
	}

	responses.RespondOK(c, client, "Media client retrieved successfully")
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
func (h *ClientHandler[T]) GetAllClients(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	log.Info().
		Uint64("userID", uid).
		Msg("Retrieving all media clients")

	clients, err := h.service.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve media clients")
		return
	}

	responses.RespondOK(c, clients, "Media clients retrieved successfully")
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

func (h *ClientHandler[T]) UpdateClient(c *gin.Context) {
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
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("id")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req requests.ClientRequest[T]
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("name", req.Name).
		Str("type", string(req.ClientType)).
		Msg("Updating media client")

	client, err := h.service.Update(ctx, req.Client)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "media client not found" {
			responses.RespondNotFound(c, err, "Media client not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to update media client")
		return
	}

	responses.RespondOK(c, client, "Media client updated successfully")
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
func (h *ClientHandler[T]) DeleteClient(c *gin.Context) {
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
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("id")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Deleting media client")

	err = h.service.Delete(ctx, uid, clientID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "media client not found" {
			responses.RespondNotFound(c, err, "Media client not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to delete media client")
		return
	}

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Media client deleted successfully")
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
func (h *ClientHandler[T]) TestConnection(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	var req requests.ClientRequest[T]
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("type", string(req.ClientType)).
		Msg("Testing media client connection")

	result, err := h.service.TestConnection(ctx, req.Client)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to test media client connection")
		return
	}

	responses.RespondOK(c, result, "Connection test completed")
}
