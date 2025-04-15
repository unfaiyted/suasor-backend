// handlers/client_media_collection.go
package handlers

import (
	"strconv"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/services"
	models "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"

	"github.com/gin-gonic/gin"
)

func createCollectionMediaItem[T mediatypes.Collection](clientID uint64, clientType clienttypes.ClientMediaType, externalID string, data mediatypes.Collection) models.MediaItem[mediatypes.Collection] {
	mediaItem := models.MediaItem[mediatypes.Collection]{
		Type:        mediatypes.MediaTypeCollection,
		Data:        data,
		SyncClients: []models.SyncClient{},
		ExternalIDs: []models.ExternalID{},
	}

	// Set client info and external ID
	mediaItem.SetClientInfo(clientID, clientType, externalID)

	// Only add external ID if provided
	if externalID != "" {
		mediaItem.AddExternalID("client", externalID)
	}

	return mediaItem
}

// ClientMediaCollectionHandler handles collection-related operations for media clients
type ClientMediaCollectionHandler[T clienttypes.ClientMediaConfig] struct {
	collectionService services.ClientMediaCollectionService
}

// NewClientMediaCollectionHandler creates a new media client collection handler
func NewClientMediaCollectionHandler[T clienttypes.ClientMediaConfig](collectionService services.ClientMediaCollectionService) *ClientMediaCollectionHandler[T] {
	return &ClientMediaCollectionHandler[T]{
		collectionService: collectionService,
	}
}

// GetCollectionByID godoc
// @Summary Get collection by ID
// @Description Retrieves a specific collection from the client by ID
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param collectionID path string true "Collection ID"
// @Success 200 {object} responses.APIResponse[responses.MediaItemResponse] "Collection retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /clients/media/{clientID}/collections/{collectionID} [get]
func (h *ClientMediaCollectionHandler[T]) GetCollectionByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting collection by ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access collection without authentication")
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

	collectionID := c.Param("id")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("collectionID", collectionID).
		Msg("Retrieving collection by ID")

	// This is a placeholder. In actual implementations, you would implement a GetCollectionByID method
	// in the services.ClientMediaCollectionService interface.
	responses.RespondNotImplemented(c, nil, "Get collection by ID not implemented")
}

// GetCollections godoc
// @Summary Get all collections
// @Description Retrieves all collections from the client
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param count query int false "Maximum number of collections to return"
// @Success 200 {object} responses.APIResponse[[]responses.MediaItemResponse] "Collections retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /clients/media/{clientID}/collections [get]
func (h *ClientMediaCollectionHandler[T]) GetCollections(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting all collections")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access collections without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Get count parameter
	count := 0
	countParam := c.Query("count")
	if countParam != "" {
		var err error
		count, err = strconv.Atoi(countParam)
		if err != nil {
			log.Error().Err(err).Str("count", countParam).Msg("Invalid count format")
			responses.RespondBadRequest(c, err, "Invalid count")
			return
		}
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving collections")

	// This is a placeholder. In actual implementations, you would implement a GetCollections method
	// in the services.ClientMediaCollectionService interface.
	responses.RespondNotImplemented(c, nil, "Get collections not implemented")
}
