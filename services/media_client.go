package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"suasor/client/media/providers"
	"suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"time"
)

// ClientService defines the interface for media client operations
type ClientService interface {
	CreateClient(ctx context.Context, userID uint64, req requests.ClientRequest) (responses.ClientResponse, error)
	GetClientByID(ctx context.Context, userID, clientID uint64) (responses.ClientResponse, error)
	GetClientsByUserID(ctx context.Context, userID uint64) ([]responses.ClientResponse, error)
	UpdateClient(ctx context.Context, userID, clientID uint64, req requests.ClientRequest) (responses.ClientResponse, error)
	DeleteClient(ctx context.Context, userID, clientID uint64) error
	TestClientConnection(ctx context.Context, req requests.ClientTestRequest) (responses.ClientTestResponse, error)
	GetContentProvider(ctx context.Context, userID, clientID uint64) (types.ContentProvider, error)
}

type mediaClientService struct {
	repo repository.ClientRepository
}

// NewClientService creates a new media client service
func NewClientService(repo repository.MediaClientRepository) MediaClientService {
	return &mediaClientService{repo: repo}
}

// CreateClient creates a new media client configuration
func (s *mediaClientService) CreateClient(ctx context.Context, userID uint64, req requests.ClientRequest) (types.MediaClientResponse, error) {
	// Test connection first
	testReq := types.ClientTestRequest{
		ClientType: req.ClientType,
		Client:     req.Client,
	}

	testResp, err := s.TestClientConnection(ctx, testReq)
	if err != nil || !testResp.Success {
		return responses.ClientResponse{}, fmt.Errorf("failed to connect to client: %v", err)
	}

	var client interface{}

	switch req.ClientType {
	case types.ClientTypePlex:
		config, ok := req.Client.(models.PlexConfig)
		if !ok {
			return responses.ClientResponse{}, fmt.Errorf("invalid Plex configuration")
		}
		mediaClient := types.Client[types.PlexConfig]{
			UserID:     userID,
			Name:       req.Name,
			ClientType: req.ClientType,
			Client:     config,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		client, err = s.repo.CreatePlex(ctx, mediaClient)

	case types.ClientTypeJellyfin:
		config, ok := req.Client.(models.JellyfinConfig)
		if !ok {
			return responses.ClientResponse{}, fmt.Errorf("invalid Jellyfin configuration")
		}
		mediaClient := types.Client[types.JellyfinConfig]{
			UserID:     userID,
			Name:       req.Name,
			ClientType: req.ClientType,
			Client:     config,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		client, err = s.repo.CreateJellyfin(ctx, mediaClient)

	case types.ClientTypeEmby:
		config, ok := req.Client.(models.EmbyConfig)
		if !ok {
			return responses.ClientResponse{}, fmt.Errorf("invalid Emby configuration")
		}
		mediaClient := types.Client[types.EmbyConfig]{
			UserID:     userID,
			Name:       req.Name,
			ClientType: req.ClientType,
			Client:     config,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		client, err = s.repo.CreateEmby(ctx, mediaClient)

	case types.ClientTypeSubsonic:
		config, ok := req.Client.(types.SubsonicConfig)
		if !ok {
			return responses.ClientResponse{}, fmt.Errorf("invalid Subsonic configuration")
		}
		mediaClient := types.Client[types.SubsonicConfig]{
			UserID:     userID,
			Name:       req.Name,
			ClientType: req.ClientType,
			Client:     config,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		client, err = s.repo.CreateNavidrome(ctx, mediaClient)

	default:
		return responses.ClientResponse{}, fmt.Errorf("unsupported client type: %s", req.ClientType)
	}

	if err != nil {
		return responses.ClientResponse{}, err
	}

	// Convert to response type based on the client type
	switch req.ClientType {
	case types.ClientTypePlex:
		return models.ToResponse(client.(*types.Client[models.PlexConfig])), nil
	case types.ClientTypeJellyfin:
		return models.ToResponse(client.(*types.Client[models.JellyfinConfig])), nil
	case types.ClientTypeEmby:
		return models.ToResponse(client.(*types.Client[models.EmbyConfig])), nil
	case types.ClientTypeSubsonic:
		return models.ToResponse(client.(*types.Client[models.NavidromeConfig])), nil
	default:
		return responses.ClientResponse{}, fmt.Errorf("unsupported client type: %s", req.ClientType)
	}
}

// GetClientByID retrieves a media client by ID
func (s *mediaClientService) GetClientByID(ctx context.Context, userID, clientID uint64) (responses.ClientResponse, error) {
	client, err := s.repo.GetByID(ctx, clientID, userID)
	if err != nil {
		return responses.ClientResponse{}, err
	}

	return client, nil
}

// GetClientsByUserID retrieves all media clients for a user
func (s *mediaClientService) GetClientsByUserID(ctx context.Context, userID uint64) ([]responses.ClientResponse, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// UpdateClient updates an existing media client
func (s *mediaClientService) UpdateClient(ctx context.Context, userID, clientID uint64, req requests.ClientRequest) (types.MediaClientResponse, error) {
	// Test connection with updated information
	testReq := types.ClientTestRequest{
		ClientType: req.ClientType,
		Client:     req.Client,
	}

	testResp, err := s.TestClientConnection(ctx, testReq)
	if err != nil || !testResp.Success {
		return responses.ClientResponse{}, fmt.Errorf("failed to connect to updated client: %v", err)
	}

	// Get existing client to verify it exists and belongs to user
	_, err = s.GetClientByID(ctx, userID, clientID)
	if err != nil {
		return responses.ClientResponse{}, err
	}

	var updatedClient interface{}

	switch req.ClientType {
	case types.ClientTypePlex:
		config, ok := req.Client.(models.PlexConfig)
		if !ok {
			return responses.ClientResponse{}, fmt.Errorf("invalid Plex configuration")
		}
		mediaClient := types.Client[types.PlexConfig]{
			ID:         clientID,
			UserID:     userID,
			Name:       req.Name,
			ClientType: req.ClientType,
			Client:     config,
			UpdatedAt:  time.Now(),
		}
		updatedClient, err = s.repo.UpdatePlex(ctx, mediaClient)

	case types.ClientTypeJellyfin:
		config, ok := req.Client.(models.JellyfinConfig)
		if !ok {
			return responses.ClientResponse{}, fmt.Errorf("invalid Jellyfin configuration")
		}
		mediaClient := types.Client[types.JellyfinConfig]{
			ID:         clientID,
			UserID:     userID,
			Name:       req.Name,
			ClientType: req.ClientType,
			Client:     config,
			UpdatedAt:  time.Now(),
		}
		updatedClient, err = s.repo.UpdateJellyfin(ctx, mediaClient)

	case types.ClientTypeEmby:
		config, ok := req.Client.(models.EmbyConfig)
		if !ok {
			return responses.ClientResponse{}, fmt.Errorf("invalid Emby configuration")
		}
		mediaClient := types.Client[types.EmbyConfig]{
			ID:         clientID,
			UserID:     userID,
			Name:       req.Name,
			ClientType: req.ClientType,
			Client:     config,
			UpdatedAt:  time.Now(),
		}
		updatedClient, err = s.repo.UpdateEmby(ctx, mediaClient)

	case types.ClientTypeSubsonic:
		config, ok := req.Client.(types.SubsonicConfig)
		if !ok {
			return responses.ClientResponse{}, fmt.Errorf("invalid Subsonic configuration")
		}
		mediaClient := types.Client[types.SubsonicConfig]{
			ID:         clientID,
			UserID:     userID,
			Name:       req.Name,
			ClientType: req.ClientType,
			Client:     config,
			UpdatedAt:  time.Now(),
		}
		updatedClient, err = s.repo.UpdateNavidrome(ctx, mediaClient)

	default:
		return responses.ClientResponse{}, fmt.Errorf("unsupported client type: %s", req.ClientType)
	}

	if err != nil {
		return responses.ClientResponse{}, err
	}

	// Convert to response type based on the client type
	switch req.ClientType {
	case types.ClientTypePlex:
		return models.ToResponse(updatedClient.(*types.Client[models.PlexConfig])), nil
	case types.ClientTypeJellyfin:
		return models.ToResponse(updatedClient.(*types.Client[models.JellyfinConfig])), nil
	case types.ClientTypeEmby:
		return models.ToResponse(updatedClient.(*types.Client[models.EmbyConfig])), nil
	case types.ClientTypeSubsonic:
		return models.ToResponse(updatedClient.(*types.Client[models.NavidromeConfig])), nil
	default:
		return responses.ClientResponse{}, fmt.Errorf("unsupported client type: %s", req.ClientType)
	}
}

// DeleteClient deletes a media client
func (s *mediaClientService) DeleteClient(ctx context.Context, userID, clientID uint64) error {
	return s.repo.Delete(ctx, clientID, userID)
}

// TestClientConnection tests the connection to a media client
func (s *mediaClientService) TestClientConnection(ctx context.Context, req types.ClientTestRequest) (types.MediaClientTestResponse, error) {
	switch req.ClientType {
	case types.ClientTypePlex:
		return s.testPlexConnection(ctx, req.Client)
	case types.ClientTypeJellyfin:
		return s.testJellyfinConnection(ctx, req.Client)
	case types.ClientTypeEmby:
		return s.testEmbyConnection(ctx, req.Client)
	case types.ClientTypeSubsonic:
		return s.testSubsonicConnection(ctx, req.Client)
	default:
		return types.ClientTestResponse{
			Success: false,
			Message: "Unsupported client type",
		}, fmt.Errorf("unsupported client type: %s", req.ClientType)
	}
}

// Client-specific test connection methods
func (s *mediaClientService) testPlexConnection(ctx context.Context, clientConfig interface{}) (types.ClientTestResponse, error) {
	config, ok := clientConfig.(models.PlexConfig)
	if !ok {
		return types.ClientTestResponse{
			Success: false,
			Message: "Invalid Plex configuration",
		}, fmt.Errorf("invalid plex configuration")
	}

	protocol := "http"
	if config.SSL {
		protocol = "https"
	}

	serverURL := fmt.Sprintf("%s://%s:%d", protocol, config.Host, config.Port)

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL+"/identity", nil)
	if err != nil {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}, err
	}

	req.Header.Add("X-Plex-Token", config.Token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect to Plex server: %v", err),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Plex server returned status code %d", resp.StatusCode),
		}, fmt.Errorf("plex server returned status code %d", resp.StatusCode)
	}

	// Could parse the response for version info but keeping it simple
	return types.ClientTestResponse{
		Success: true,
		Message: "Successfully connected to Plex server",
	}, nil
}

func (s *mediaClientService) testJellyfinConnection(ctx context.Context, clientConfig interface{}) (types.ClientTestResponse, error) {
	config, ok := clientConfig.(models.JellyfinConfig)
	if !ok {
		return types.ClientTestResponse{
			Success: false,
			Message: "Invalid Jellyfin configuration",
		}, fmt.Errorf("invalid jellyfin configuration")
	}

	protocol := "http"
	if config.SSL {
		protocol = "https"
	}

	serverURL := fmt.Sprintf("%s://%s:%d", protocol, config.Host, config.Port)

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL+"/System/Info", nil)
	if err != nil {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}, err
	}

	req.Header.Add("X-API-Key", config.APIKey)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect to Jellyfin server: %v", err),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Jellyfin server returned status code %d", resp.StatusCode),
		}, fmt.Errorf("jellyfin server returned status code %d", resp.StatusCode)
	}

	return types.ClientTestResponse{
		Success: true,
		Message: "Successfully connected to Jellyfin server",
	}, nil
}

func (s *mediaClientService) testEmbyConnection(ctx context.Context, clientConfig interface{}) (types.ClientTestResponse, error) {
	config, ok := clientConfig.(models.EmbyConfig)
	if !ok {
		return types.ClientTestResponse{
			Success: false,
			Message: "Invalid Emby configuration",
		}, fmt.Errorf("invalid emby configuration")
	}

	protocol := "http"
	if config.SSL {
		protocol = "https"
	}

	serverURL := fmt.Sprintf("%s://%s:%d", protocol, config.Host, config.Port)

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL+"/System/Info", nil)
	if err != nil {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}, err
	}

	req.Header.Add("X-API-Key", config.APIKey)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect to Emby server: %v", err),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Emby server returned status code %d", resp.StatusCode),
		}, fmt.Errorf("emby server returned status code %d", resp.StatusCode)
	}

	return types.ClientTestResponse{
		Success: true,
		Message: "Successfully connected to Emby server",
	}, nil
}

func (s *mediaClientService) testSubsonicConnection(ctx context.Context, clientConfig interface{}) (types.ClientTestResponse, error) {
	config, ok := clientConfig.(models.NavidromeConfig)
	if !ok {
		return types.ClientTestResponse{
			Success: false,
			Message: "Invalid Subsonic configuration",
		}, fmt.Errorf("invalid subsonic configuration")
	}

	protocol := "http"
	if config.SSL {
		protocol = "https"
	}

	serverURL := fmt.Sprintf("%s://%s:%d", protocol, config.Host, config.Port)

	// Subsonic API parameters
	params := url.Values{}
	params.Add("u", config.Username)
	params.Add("p", config.Password)
	params.Add("v", "1.16.1")
	params.Add("c", "suasor")
	params.Add("f", "json")

	pingURL := fmt.Sprintf("%s/rest/ping.view?%s", serverURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", pingURL, nil)
	if err != nil {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}, err
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect to Subsonic server: %v", err),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return types.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Subsonic server returned status code %d", resp.StatusCode),
		}, fmt.Errorf("subsonic server returned status code %d", resp.StatusCode)
	}

	return types.ClientTestResponse{
		Success: true,
		Message: "Successfully connected to Subsonic server",
	}, nil
}

// GetContentProvider returns a ContentProvider for the specified client
func (s *mediaClientService) GetContentProvider(ctx context.Context, userID, clientID uint64) (providers.ContentProvider, error) {
	clientResp, err := s.GetClientByID(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	return providers.NewContentProvider(clientID, clientResp.ClientType, clientResp.Client)
}
