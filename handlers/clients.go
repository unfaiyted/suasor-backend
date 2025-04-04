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

// ListAllClients godoc
// @Summary Get all clients
// @Description Retrieves all configured clients across different types for the user
// @Tags clients
// @Accept json
// @Produce json
// @Param type query string false "Filter by client category (e.g. 'media')"
// @Param clientType query string false "Filter by specific client type (e.g. 'jellyfin')"
// @Security BearerAuth
// @Success 200 {object} responses.ClientsResponse "All user clients with various config types"
// @Failure 401 {object} responses.BasicErrorResponse "Unauthorized"
// @Failure 500 {object} responses.BasicErrorResponse "Server error"
// @Router /clients [get]
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

	// Get filter parameters
	typeFilter := c.Query("type")
	clientTypeFilter := c.Query("clientType")

	log.Info().
		Uint64("userID", uid).
		Str("typeFilter", typeFilter).
		Str("clientTypeFilter", clientTypeFilter).
		Msg("Retrieving filtered clients")

	// Create a slice to hold all client responses
	var allClients []responses.ClientResponse

	// Define client type categories
	mediaClientTypes := map[string]bool{
		"emby":     true,
		"jellyfin": true,
		"plex":     true,
		"subsonic": true,
		"sonarr":   true,
		"radarr":   true,
		"lidarr":   true,
	}

	aiClientTypes := map[string]bool{
		"claude": true,
		"openai": true,
		"ollama": true,
	}

	autoClientTypes := map[string]bool{
		"radarr": true,
		"sonarr": true,
		"lidarr": true,
	}

	// Helper function to check if we should fetch a specific client type
	shouldFetchClientType := func(clientType string) bool {
		// If clientType filter is specified, only return that type
		if clientTypeFilter != "" {
			return clientTypeFilter == clientType
		}

		// If type filter is "media", only return media clients
		if typeFilter == "media" {
			return mediaClientTypes[clientType]
		}
		if typeFilter == "ai" {
			return aiClientTypes[clientType]
		}
		if typeFilter == "automation" {
			return autoClientTypes[clientType]
		}

		// No relevant filters, return all clients
		return true
	}

	// Fetch all client types based on filters
	if shouldFetchClientType("emby") {
		embyClients, _ := h.embyService.GetByUserID(ctx, uid)
		for _, client := range embyClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("jellyfin") {
		jellyfinClients, _ := h.jellyfinService.GetByUserID(ctx, uid)
		for _, client := range jellyfinClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("plex") {
		plexClients, _ := h.plexService.GetByUserID(ctx, uid)
		for _, client := range plexClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("subsonic") {
		subsonicClients, _ := h.subsonicService.GetByUserID(ctx, uid)
		for _, client := range subsonicClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("sonarr") {
		sonarrClients, _ := h.sonarrService.GetByUserID(ctx, uid)
		for _, client := range sonarrClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("radarr") {
		radarrClients, _ := h.radarrService.GetByUserID(ctx, uid)
		for _, client := range radarrClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("lidarr") {
		lidarrClients, _ := h.lidarrService.GetByUserID(ctx, uid)
		for _, client := range lidarrClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("claude") {
		claudeClients, _ := h.claudeService.GetByUserID(ctx, uid)
		for _, client := range claudeClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("openai") {
		openaiClients, _ := h.openaiService.GetByUserID(ctx, uid)
		for _, client := range openaiClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("ollama") {
		ollamaClients, _ := h.ollamaService.GetByUserID(ctx, uid)
		for _, client := range ollamaClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("claude") {
		claudeClients, _ := h.claudeService.GetByUserID(ctx, uid)
		for _, client := range claudeClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	if shouldFetchClientType("openai") {
		openaiClients, _ := h.openaiService.GetByUserID(ctx, uid)
		for _, client := range openaiClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}
	if shouldFetchClientType("ollama") {
		ollamaClients, _ := h.ollamaService.GetByUserID(ctx, uid)
		for _, client := range ollamaClients {
			allClients = append(allClients, toClientResponse(client))
		}
	}

	log.Info().
		Uint64("userID", uid).
		Int("clientCount", len(allClients)).
		Str("typeFilter", typeFilter).
		Str("clientTypeFilter", clientTypeFilter).
		Msg("Retrieved filtered clients")

	responses.RespondOK(c, allClients, "Clients retrieved successfully")
}

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
