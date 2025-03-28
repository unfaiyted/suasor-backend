package handlers

//
// import (
// 	gin "github.com/gin-gonic/gin"
//
// 	"gorm.io/gorm"
// 	cTypes "suasor/client/types"
// 	"suasor/services"
// 	"suasor/types"
// 	"suasor/types/models"
// )
//
// // ClientHandler handles all client operations
// type ClientHandler struct {
// 	db *gorm.DB
// }
//
// func NewClientHandler(db *gorm.DB) *ClientHandler {
// 	return &ClientHandler{db: db}
// }
//
// // CreateClient handles client creation for any client type
// func (h *ClientHandler) CreateClient(c *gin.Context) {
// 	var req struct {
// 		Type   string          `json:"type" binding:"required"`
// 		Name   string          `json:"name" binding:"required"`
// 		Config json.RawMessage `json:"config" binding:"required"`
// 	}
//
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	userID := getUserIDFromContext(c)
//
// 	// Handle different client types
// 	switch req.Type {
// 	case "jellyfin":
// 		h.createTypedClient[types.JellyfinConfig](c, req.Name, userID, req.Config)
// 	case "emby":
// 		h.createTypedClient[types.EmbyConfig](c, req.Name, userID, req.Config)
// 	case "plex":
// 		h.createTypedClient[types.PlexConfig](c, req.Name, userID, req.Config)
// 	case "radarr":
// 		h.createTypedClient[types.RadarrConfig](c, req.Name, userID, req.Config)
// 	case "sonarr":
// 		h.createTypedClient[types.SonarrConfig](c, req.Name, userID, req.Config)
// 	case "lidarr":
// 		h.createTypedClient[types.LidarrConfig](c, req.Name, userID, req.Config)
// 	case "subsonic":
// 		h.createTypedClient[types.SubsonicConfig](c, req.Name, userID, req.Config)
// 	default:
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported client type"})
// 	}
// }
//
// // Type-specific client creation helper
// func (h *ClientHandler) createTypedClient(c *gin.Context, name string, userID uint64, configData json.RawMessage) {
// 	var config T
// 	if err := json.Unmarshal(configData, &config); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration format"})
// 		return
// 	}
//
// 	clientType := getClientTypeFromConfig(config)
//
// 	client := models.Client[T]{
// 		UserID: userID,
// 		Name:   name,
// 		Type:   clientType,
// 		Config: models.ClientConfigWrapper[T]{Data: config},
// 	}
//
// 	service := services.NewClientService[T](h.db)
// 	result, err := service.Create(c.Request.Context(), client)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	c.JSON(http.StatusCreated, result)
// }
//
// // UpdateClient handles client updates for any client type
// func (h *ClientHandler) UpdateClient(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
// 		return
// 	}
//
// 	var req struct {
// 		Type   string          `json:"type" binding:"required"`
// 		Name   string          `json:"name" binding:"required"`
// 		Config json.RawMessage `json:"config" binding:"required"`
// 	}
//
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	userID := getUserIDFromContext(c)
//
// 	switch req.Type {
// 	case "jellyfin":
// 		h.updateTypedClient[types.JellyfinConfig](c, id, req.Name, userID, req.Config)
// 	case "emby":
// 		h.updateTypedClient[types.EmbyConfig](c, id, req.Name, userID, req.Config)
// 	// Handle other types
// 	default:
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported client type"})
// 	}
// }
//
// // GetClient fetches a client by ID
// func (h *ClientHandler) GetClient(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
// 		return
// 	}
//
// 	clientType := c.Param("type")
// 	userID := getUserIDFromContext(c)
//
// 	switch clientType {
// 	case "jellyfin":
// 		h.getTypedClient[types.JellyfinConfig](c, id, userID)
// 	case "emby":
// 		h.getTypedClient[types.EmbyConfig](c, id, userID)
// 	// Handle other types
// 	default:
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported client type"})
// 	}
// }
//
// // GetClientsByUser gets all clients for the current user
// func (h *ClientHandler) GetClientsByUser(c *gin.Context) {
// 	userID := getUserIDFromContext(c)
// 	clientType := c.Query("type") // Optional type filter
//
// 	if clientType != "" {
// 		// If type is specified, get clients of that type
// 		switch clientType {
// 		case "jellyfin":
// 			h.getClientsByType[types.JellyfinConfig](c, userID, clientType)
// 		case "emby":
// 			h.getClientsByType[types.EmbyConfig](c, userID, clientType)
// 		// Handle other types
// 		default:
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported client type"})
// 		}
// 		return
// 	}
//
// 	// If no type specified, gather all client types
// 	// This is more complex and requires merging results from different repos
// 	// ...implementation depends on your specific requirements
// }
//
// // DeleteClient deletes a client by ID
// func (h *ClientHandler) DeleteClient(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
// 		return
// 	}
//
// 	clientType := c.Param("type")
// 	userID := getUserIDFromContext(c)
//
// 	switch clientType {
// 	case "jellyfin":
// 		h.deleteTypedClient[types.JellyfinConfig](c, id, userID)
// 	case "emby":
// 		h.deleteTypedClient[types.EmbyConfig](c, id, userID)
// 	// Handle other types
// 	default:
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported client type"})
// 	}
// }
//
// // TestClient tests connection to client API
// func (h *ClientHandler) TestClient(c *gin.Context) {
// 	var req struct {
// 		Type   string          `json:"type" binding:"required"`
// 		Config json.RawMessage `json:"config" binding:"required"`
// 	}
//
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	switch req.Type {
// 	case "jellyfin":
// 		h.testTypedClient[types.JellyfinConfig](c, req.Config)
// 	case "emby":
// 		h.testTypedClient[types.EmbyConfig](c, req.Config)
// 	// Handle other types
// 	default:
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported client type"})
// 	}
// }
//
// // Helper methods for typed operations
// func (h *ClientHandler) testTypedClient(c *gin.Context, configData json.RawMessage) {
// 	var config T
// 	if err := json.Unmarshal(configData, &config); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration format"})
// 		return
// 	}
//
// 	service := NewClientService[T](h.db)
// 	if err := service.TestConnection(c.Request.Context(), config); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{"success": true})
// }
//
// // Implement other helper methods (getTypedClient, updateTypedClient, etc.)
