package handlers

import (
	"strconv"
	"suasor/services"

	"github.com/gin-gonic/gin"
	// "io"
	"suasor/clients/types"
	client "suasor/clients/types"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"
)

type ClientHandler[T types.ClientConfig] interface {
	CreateClient(c *gin.Context)

	GetClient(c *gin.Context)
	UpdateClient(c *gin.Context)
	DeleteClient(c *gin.Context)
	TestConnection(c *gin.Context)
}

// clientHandler handles setting up client API endpoints
type clientHandler[T client.ClientConfig] struct {
	service services.ClientService[T]
}

// NewclientHandler creates a new client handler
func NewClientHandler[T types.ClientConfig](service services.ClientService[T]) ClientHandler[T] {
	return &clientHandler[T]{
		service: service,
	}
}

// CreateClient godoc
// @Summary Create a new client
// @Description Creates a new client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.ClientRequest[client.ClientConfig] true "client data"
// @Param clientType path string true "Client type"
// @Success 201 {object} responses.APIResponse[models.Client[client.ClientConfig]] "client created"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientType} [post]
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
func (h *clientHandler[T]) CreateClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Msg("Creating new client")

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)

	var req requests.ClientRequest[T]
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		responses.RespondValidationError(c, err)
		return
	}

	clientType := client.ClientType(req.ClientType)
	category := clientType.AsCategory()

	client := req.Client
	client.SetCategory(category)

	clientOfType := models.Client[T]{
		UserID:   uid,
		Name:     req.Name,
		Type:     clientType,
		Category: category,
		Config:   models.ClientConfigWrapper[T]{Data: client},
	}
	log.Debug().
		Uint64("userID", uid).
		Str("name", req.Name).
		Str("type", string(req.ClientType)).
		Str("category", clientOfType.GetCategory().String()).
		Msg("Creating new client")

	if h.service == nil {
		log.Error().Msg("Client service is not initialized")
		responses.RespondInternalError(c, nil, "Internal server error: service not initialized")
		return
	}

	clnt, err := h.service.Create(ctx, clientOfType)
	if err != nil {
		responses.RespondInternalError(c, err, err.Error())
		return
	}

	responses.RespondCreated(c, clnt, "Client created successfully")
}

// GetClient godoc
// @Summary Get client
// @Description Retrieves a specific client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Success 200 {object} responses.APIResponse[models.Client[client.ClientConfig]] "Client retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Client not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/:clientType/{id} [get]
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
//   "message": "Client retrieved successfully"
// }

func (h *clientHandler[T]) GetClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)

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
		Msg("Retrieving client")

	client, err := h.service.GetByID(ctx, uid, clientID)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "client not found" {
			responses.RespondNotFound(c, err, "client not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to retrieve client")
		return
	}

	responses.RespondOK(c, client, "Client retrieved successfully")
}

// UpdateClient godoc
// @Summary Update client
// @Description Updates an existing client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param clientType path string true "Client type"
// @Param request body requests.ClientRequest true "Updated client data"
// @Success 200 {object} responses.APIResponse[models.Client[client.ClientConfig]] "clients updated"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request or client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Client not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/clients/{clientType}/{clientID} [put]
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

func (h *clientHandler[T]) UpdateClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	// client ID from URL
	clientID, _ := checkClientID(c)

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
		Str("category", req.ClientType.AsCategory().String()).
		Bool("isEnabled", req.IsEnabled).
		Msg("Updating client")

	client := models.Client[T]{
		BaseModel: models.BaseModel{
			ID: clientID,
		},
		UserID:    uid,
		Name:      req.Name,
		Category:  req.ClientType.AsCategory(),
		Type:      req.ClientType,
		Config:    models.ClientConfigWrapper[T]{Data: req.Client},
		IsEnabled: req.IsEnabled,
	}

	updatedClient, e := h.service.Update(ctx, client)
	if e != nil {
		responses.RespondInternalError(c, e, "Failed to update client")
		return
	}

	responses.RespondOK(c, updatedClient, "client updated successfully")
}

// DeleteClient godoc
// @Summary Delete client
// @Description Deletes a client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Success 200 {object} responses.APIResponse[responses.EmptyResponse] "client deleted"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Client not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/admin/client/{clientID} [delete]
func (h *clientHandler[T]) DeleteClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	uid, _ := checkUserAccess(c)
	clientID, _ := checkClientID(c)

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Msg("Deleting client")

	err := h.service.Delete(ctx, clientID, uid)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "client not found" {
			responses.RespondNotFound(c, err, "client not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to delete client")
		return
	}

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "client deleted successfully")
}

// TestConnection godoc
// @Summary Test client connection
// @Description Tests the connection to a client using the provided configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path uint64 true "Client ID"
// @Success 200 {object} responses.APIResponse[responses.TestConnectionResponse] "Connection test result"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/admin/client/{clientID}/test [get]
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
func (h *clientHandler[T]) TestConnection(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	cid, _ := checkClientID(c)

	clientType := c.Param("clientType")

	log.Info().
		Uint64("userID", uid).
		Str("clientType", clientType).
		Uint64("clientID", cid).
		Msg("Testing client connection")

	client, err := h.service.GetByID(ctx, cid, uid)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "client not found" {
			responses.RespondNotFound(c, err, "client not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to retrieve client")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("type", string(client.GetClientType())).
		Msg("Testing client connection")

	result, err := h.service.TestConnection(ctx, cid, &client.Config.Data)
	if err != nil {
		responses.RespondInternalError(c, err, result.Message)
		return
	}

	responses.RespondOK(c, result, "Connection test completed")
}
