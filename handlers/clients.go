package handlers

import (
	"suasor/client/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"

	"github.com/gin-gonic/gin"
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

// NewClientsHandler creates a new handler for listing all clients
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

// ListAllClients godoc
// @Summary Get all clients
// @Description Retrieves all configured clients across different types for the user
// @Tags clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.ClientsResponse "All user clients with various config types"
// @Failure 401 {object} responses.BasicErrorResponse "Unauthorized"
// @Failure 500 {object} responses.BasicErrorResponse "Server error"
// @Router /clients [get]
// @Example       response - Plex client example
//
//	{
//	  "data": [
//	    {
//	      "id": 1,
//	      "userId": 123,
//	      "name": "My Plex Server",
//	      "clientType": "plex",
//	      "client": {
//	        "type": "plex",
//	        "host": "192.168.1.100",
//	        "port": 32400,
//	        "token": "your-plex-token",
//	        "ssl": false,
//	        "enabled": true
//	      },
//	      "createdAt": "2023-01-01T12:00:00Z",
//	      "updatedAt": "2023-01-01T12:00:00Z"
//	    },
//	    {
//	      "id": 2,
//	      "userId": 123,
//	      "name": "My Jellyfin Server",
//		     "isEnabled": true,
//	      "clientType": "jellyfin",
//	      "client": {
//	        "type": "jellyfin",
//	        "host": "192.168.1.101",
//	        "port": 8096,
//	        "apiKey": "your-jellyfin-api-key",
//	        "username": "admin",
//	        "ssl": false,
//	        "enabled": true
//	      },
//	      "createdAt": "2023-01-01T12:00:00Z",
//	      "updatedAt": "2023-01-01T12:00:00Z"
//	    }
//	  ],
//	  "message": "All clients retrieved successfully"
//	}
func (h *ClientsHandler) ListAllClients(c *gin.Context) {
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

	// Create a slice to hold all client responses
	var allClients []responses.ClientResponse

	// Fetch all client types
	embyClients, _ := h.embyService.GetByUserID(ctx, uid)
	for _, client := range embyClients {
		allClients = append(allClients, toClientResponse(client))
	}

	jellyfinClients, _ := h.jellyfinService.GetByUserID(ctx, uid)
	for _, client := range jellyfinClients {
		allClients = append(allClients, toClientResponse(client))
	}

	plexClients, _ := h.plexService.GetByUserID(ctx, uid)
	for _, client := range plexClients {
		allClients = append(allClients, toClientResponse(client))
	}

	subsonicClients, _ := h.subsonicService.GetByUserID(ctx, uid)
	for _, client := range subsonicClients {
		allClients = append(allClients, toClientResponse(client))
	}

	sonarrClients, _ := h.sonarrService.GetByUserID(ctx, uid)
	for _, client := range sonarrClients {
		allClients = append(allClients, toClientResponse(client))
	}

	radarrClients, _ := h.radarrService.GetByUserID(ctx, uid)
	for _, client := range radarrClients {
		allClients = append(allClients, toClientResponse(client))
	}

	lidarrClients, _ := h.lidarrService.GetByUserID(ctx, uid)
	for _, client := range lidarrClients {
		allClients = append(allClients, toClientResponse(client))
	}

	claudeClients, _ := h.claudeService.GetByUserID(ctx, uid)
	for _, client := range claudeClients {
		allClients = append(allClients, toClientResponse(client))
	}

	openaiClients, _ := h.openaiService.GetByUserID(ctx, uid)
	for _, client := range openaiClients {
		allClients = append(allClients, toClientResponse(client))
	}

	ollamaClients, _ := h.ollamaService.GetByUserID(ctx, uid)
	for _, client := range ollamaClients {
		allClients = append(allClients, toClientResponse(client))
	}

	log.Info().
		Uint64("userID", uid).
		Int("clientCount", len(allClients)).
		Msg("Retrieved all clients")

	responses.RespondOK(c, allClients, "All clients retrieved successfully")
}

// toClientResponse converts a client model to response
func toClientResponse[T types.ClientConfig](client *models.Client[T]) responses.ClientResponse {
	return responses.ClientResponse{
		ID:         client.ID,
		UserID:     client.UserID,
		Name:       client.Name,
		IsEnabled:  client.IsEnabled,
		ClientType: types.MediaClientType(client.Type),
		Client:     client.Config.Data,
		CreatedAt:  client.CreatedAt,
		UpdatedAt:  client.UpdatedAt,
	}
}

// GetClientConfigs godoc
// @Summary Reference for all client config types
// @Description This endpoint doesn't exist but serves as a reference for all client config types
// @Tags swagger-reference
// @Accept json
// @Produce json
// @Success 200 {object} types.EmbyConfig "Emby client config"
// @Success 200 {object} types.JellyfinConfig "Jellyfin client config"
// @Success 200 {object} types.PlexConfig "Plex client config"
// @Success 200 {object} types.SubsonicConfig "Subsonic client config"
// @Success 200 {object} types.SonarrConfig "Sonarr client config"
// @Success 200 {object} types.RadarrConfig "Radarr client config"
// @Success 200 {object} types.LidarrConfig "Lidarr client config"
// @Success 200 {object} types.ClaudeConfig "Claude client config"
// @Success 200 {object} types.OpenAIConfig "OpenAI client config"
// @Success 200 {object} types.OllamaConfig "Ollama client config"
// @Router /docs/client-types [get]
func SwaggerClientTypes() {
	// Define all client config types for Swagger reference
	_ = &types.EmbyConfig{}
	_ = &types.JellyfinConfig{}
	_ = &types.PlexConfig{}
	_ = &types.SubsonicConfig{}
	_ = &types.SonarrConfig{}
	_ = &types.RadarrConfig{}
	_ = &types.LidarrConfig{}
	_ = &types.ClaudeConfig{}
	_ = &types.OpenAIConfig{}
	_ = &types.OllamaConfig{}
}

