package handlers

import (
	"suasor/clients/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
	"suasor/types/requests"
)

// ClientsHandler is responsible for handling requests related to listing all clients
type ClientsHandler struct {
	embyService     services.ClientService[*types.EmbyConfig]
	jellyfinService services.ClientService[*types.JellyfinConfig]
	plexService     services.ClientService[*types.PlexConfig]
	subsonicService services.ClientService[*types.SubsonicConfig]
	sonarrService   services.ClientService[*types.SonarrConfig]
	radarrService   services.ClientService[*types.RadarrConfig]
	lidarrService   services.ClientService[*types.LidarrConfig]
	claudeService   services.ClientService[*types.ClaudeConfig]
	openaiService   services.ClientService[*types.OpenAIConfig]
	ollamaService   services.ClientService[*types.OllamaConfig]
}

// NewClientsHandler godoc
// @Summary Create a new clients handler
// @Description Creates a new handler for retrieving and managing all client types
// @Tags internal
func NewClientsHandler(
	embyService services.ClientService[*types.EmbyConfig],
	jellyfinService services.ClientService[*types.JellyfinConfig],
	plexService services.ClientService[*types.PlexConfig],
	subsonicService services.ClientService[*types.SubsonicConfig],
	sonarrService services.ClientService[*types.SonarrConfig],
	radarrService services.ClientService[*types.RadarrConfig],
	lidarrService services.ClientService[*types.LidarrConfig],
	claudeService services.ClientService[*types.ClaudeConfig],
	openaiService services.ClientService[*types.OpenAIConfig],
	ollamaService services.ClientService[*types.OllamaConfig],
) *ClientsHandler {
	return &ClientsHandler{
		embyService:     embyService,
		jellyfinService: jellyfinService,
		plexService:     plexService,
		subsonicService: subsonicService,
		sonarrService:   sonarrService,
		radarrService:   radarrService,
		lidarrService:   lidarrService,
		claudeService:   claudeService,
		openaiService:   openaiService,
		ollamaService:   ollamaService,
	}
}

// GetAllClients godoc
// @Summary Get all clients
// @Description Retrieves all client configurations for the user
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.APIResponse[[]models.ClientList] "Clients retrieved"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/admin/clients [get]
func (h *ClientsHandler) GetAllClients(c *gin.Context) {
	ctx := c.Request.Context()
	// log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)

	var clientList models.ClientList

	emby, err := h.embyService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddEmbyArray(emby)
	plex, err := h.plexService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddPlexArray(plex)
	subsonic, err := h.subsonicService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddSubsonicArray(subsonic)
	sonarr, err := h.sonarrService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddSonarrArray(sonarr)
	radarr, err := h.radarrService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddRadarrArray(radarr)
	lidarr, err := h.lidarrService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddLidarrArray(lidarr)
	claude, err := h.claudeService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddClaudeArray(claude)
	openai, err := h.openaiService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddOpenAIArray(openai)
	ollama, err := h.ollamaService.GetByUserID(ctx, uid)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}
	clientList.AddOllamaArray(ollama)

	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return
	}

	responses.RespondOK(c, clientList, "clients retrieved successfully")
}

// TestNewConnection godoc
// @Summary Test client connection
// @Description Tests the connection to a client using the provided configuration
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.ClientTestRequest[client.ClientConfig] true "Updated client data"
// @Param clientType path string true "Client type"
// @Success 200 {object} responses.APIResponse[responses.TestConnectionResponse] "Connection test result"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
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
//
// @Router /api/v1/admin/clients/{clientType}/test [get]
func (h *ClientsHandler) TestNewConnection(c *gin.Context) {
	clientType, _ := checkClientType(c)

	switch clientType {
	case types.ClientTypeEmby:
		testConnection[*types.EmbyConfig](c, clientType, h.embyService)
	case types.ClientTypeJellyfin:
		testConnection[*types.JellyfinConfig](c, clientType, h.jellyfinService)
	case types.ClientTypePlex:
		testConnection[*types.PlexConfig](c, clientType, h.plexService)
	case types.ClientTypeSubsonic:
		testConnection[*types.SubsonicConfig](c, clientType, h.subsonicService)
	case types.ClientTypeSonarr:
		testConnection[*types.SonarrConfig](c, clientType, h.sonarrService)
	case types.ClientTypeRadarr:
		testConnection[*types.RadarrConfig](c, clientType, h.radarrService)
	case types.ClientTypeLidarr:
		testConnection[*types.LidarrConfig](c, clientType, h.lidarrService)
	case types.ClientTypeClaude:
		testConnection[*types.ClaudeConfig](c, clientType, h.claudeService)
	case types.ClientTypeOpenAI:
		testConnection[*types.OpenAIConfig](c, clientType, h.openaiService)
	case types.ClientTypeOllama:
		testConnection[*types.OllamaConfig](c, clientType, h.ollamaService)
	default:
		responses.RespondBadRequest(c, nil, "Unknown client type")
		return
	}
}

// GetClientsByType godoc
// @Summary Get clients by type
// @Description Retrieves all clients of a specific type for the user
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientType path string true "Client type (e.g. 'plex', 'jellyfin', 'emby')"
// @Success 200 {object} responses.APIResponse[[]models.Client[types.EmbyConfig]] "Clients retrieved"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/admin/clients/{clientType} [get]
func (h *ClientsHandler) GetClientsByType(c *gin.Context) {

	clientType, _ := checkClientType(c)

	switch clientType {
	case types.ClientTypeEmby:
		clients := getClientsByType[*types.EmbyConfig](c, clientType, h.embyService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypeJellyfin:
		clients := getClientsByType[*types.JellyfinConfig](c, clientType, h.jellyfinService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypePlex:
		clients := getClientsByType[*types.PlexConfig](c, clientType, h.plexService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypeSubsonic:
		clients := getClientsByType[*types.SubsonicConfig](c, clientType, h.subsonicService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypeSonarr:
		clients := getClientsByType[*types.SonarrConfig](c, clientType, h.sonarrService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypeRadarr:
		clients := getClientsByType[*types.RadarrConfig](c, clientType, h.radarrService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypeLidarr:
		clients := getClientsByType[*types.LidarrConfig](c, clientType, h.lidarrService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypeClaude:
		clients := getClientsByType[*types.ClaudeConfig](c, clientType, h.claudeService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypeOpenAI:
		clients := getClientsByType[*types.OpenAIConfig](c, clientType, h.openaiService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	case types.ClientTypeOllama:
		clients := getClientsByType[*types.OllamaConfig](c, clientType, h.ollamaService)
		responses.RespondOK(c, clients, "clients retrieved successfully")
	default:
		responses.RespondBadRequest(c, nil, "Unknown client type")
		return
	}
}

func testConnection[T types.ClientConfig](c *gin.Context, clientType types.ClientType, service services.ClientService[T]) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Str("clientType", clientType.String()).
		Msg("Testing new client connection")

	var request requests.ClientTestRequest[T]
	if err := c.ShouldBindJSON(&request); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	result, err := service.TestNewConnection(ctx, &request.Client)
	if err != nil {
		responses.RespondInternalError(c, err, result.Message)
		return
	}

	responses.RespondOK(c, result, "Connection test completed")
}

func getClientsByType[T types.ClientConfig](c *gin.Context, clientType types.ClientType, service services.ClientService[T]) []*models.Client[T] {
	ctx := c.Request.Context()

	clients, err := service.GetByType(ctx, clientType, 0)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve clients")
		return nil
	}

	return clients
}
