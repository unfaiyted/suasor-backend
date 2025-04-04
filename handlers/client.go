package handlers

import (
	"strconv"
	"suasor/services"

	"github.com/gin-gonic/gin"
	// "io"
	"suasor/client/types"
	client "suasor/client/types"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils"
)

// ClientHandler handles setting up client API endpoints
type ClientHandler[T client.ClientConfig] struct {
	service services.ClientService[T]
}

// NewClientHandler creates a new client handler
func NewClientHandler[T types.ClientConfig](service services.ClientService[T]) *ClientHandler[T] {
	return &ClientHandler[T]{
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
// @Success 201 {object} responses.APIResponse[models.Client[client.ClientConfig]] "client created"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /client/:clientType [post]
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
	log.Info().
		Msg("Creating new client")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// requestBody, err := io.ReadAll(c.Request.Body)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to read request body")
	// 	responses.RespondInternalError(c, err, "Failed to read request body")
	// 	return
	// }
	// log.Debug().
	// 	Str("requestBody", string(requestBody)).
	// 	Msg("Received request body")

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
// @Router /client/:clientType/{id} [get]
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

// GetAllClients godoc
// @Summary Get all clients
// @Description Retrieves all client configurations for the user
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.APIResponse[[]models.Client[client.ClientConfig]] "Clients retrieved"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /client/:clientType [get]
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
		Msg("Retrieving all clients")

	clients, err := h.service.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}

	responses.RespondOK(c, clients, "clients retrieved successfully")
}

// UpdateClient godoc
// @Summary Update client
// @Description Updates an existing client configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Param request body requests.ClientRequest true "Updated client data"
// @Success 200 {object} responses.APIResponse[models.Client[client.ClientConfig]] "clients updated"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request or client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Client not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /clients/:clientType/{id} [put]
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

	updatedClient, err := h.service.Update(ctx, client)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "client not found" {
			responses.RespondNotFound(c, err, "client not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to update client")
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
// @Param id path int true "Client ID"
// @Success 200 {object} responses.APIResponse[responses.EmptyResponse] "client deleted"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Client not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /clients/:clientType/{id} [delete]
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
		Msg("Deleting client")

	err = h.service.Delete(ctx, clientID, uid)
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
// @Param clientType path string true "Client type"
// @Param id path uint64 true "Client ID"
// @Success 200 {object} responses.APIResponse[responses.TestConnectionResponse] "Connection test result"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /admin/client/:clientType/:clientId/test [get]
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
	clientID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("id")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}
	clientType := c.Param("clientType")
	uid := userID.(uint64)
	cid := uint64(clientID)

	log.Info().
		Uint64("userID", uid).
		Str("clientType", clientType).
		Int("clientID", clientID).
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

// GetClientsByType godoc
// @Summary Get clients by type
// @Description Retrieves all clients of a specific type for the user
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientType path string true "Client type (e.g. 'plex', 'jellyfin', 'emby')"
// @Success 200 {object} responses.APIResponse[[]models.Client[client.ClientConfig]] "Clients retrieved"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /client/{clientType} [get]
func (h *ClientHandler[T]) GetClientsByType(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	clientType := client.ClientType(c.Param("clientType"))
	log.Info().
		Str("clientType", clientType.String()).
		Msg("Retrieving clients")
	clients, err := h.service.GetByType(ctx, clientType, userID.(uint64))
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	responses.RespondOK(c, clients, "clients retrieved successfully")
}
