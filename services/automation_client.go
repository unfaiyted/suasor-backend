// services/automation_client.go
package services

import (
	"context"
	"fmt"
	client "suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"time"

	lidarr "github.com/devopsarr/lidarr-go/lidarr"
	radarr "github.com/devopsarr/radarr-go/radarr"
	sonarr "github.com/devopsarr/sonarr-go/sonarr"
)

// AutomationClientService defines the interface for automation client operations
type AutomationClientService interface {
	CreateClient(ctx context.Context, userID uint64, req requests.AutomationClientRequest) (client.AutomationClient, error)
	GetClientByID(ctx context.Context, userID, clientID uint64) (models.AutomationClient, error)
	GetClientsByUserID(ctx context.Context, userID uint64) ([]models.AutomationClient, error)
	UpdateClient(ctx context.Context, userID, clientID uint64, req models.AutomationClientRequest) (models.AutomationClient, error)
	DeleteClient(ctx context.Context, userID, clientID uint64) error
	TestClientConnection(ctx context.Context, req models.ClientTestRequest) (models.ClientTestResponse, error)
}

type automationClientService struct {
	repo repository.AutomationClientRepository
}

func NewAutomationClientService(repo repository.AutomationClientRepository) AutomationClientService {
	return &automationClientService{repo: repo}
}

func (s *automationClientService) CreateClient(ctx context.Context, userID uint64, req models.AutomationClientRequest) (models.AutomationClient, error) {
	// Test connection first
	testReq := models.ClientTestRequest{
		URL:        req.URL,
		APIKey:     req.APIKey,
		ClientType: req.ClientType,
	}

	testResp, err := s.TestClientConnection(ctx, testReq)
	if err != nil || !testResp.Success {
		return models.AutomationClient{}, fmt.Errorf("failed to connect to client: %v", err)
	}

	client := models.AutomationClient{
		UserID:     userID,
		Name:       req.Name,
		ClientType: req.ClientType,
		URL:        req.URL,
		APIKey:     req.APIKey,
		IsEnabled:  req.IsEnabled,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.repo.Create(ctx, client)
}

func (s *automationClientService) GetClientByID(ctx context.Context, userID, clientID uint64) (models.AutomationClient, error) {
	return s.repo.GetByID(ctx, clientID, userID)
}

func (s *automationClientService) GetClientsByUserID(ctx context.Context, userID uint64) ([]models.AutomationClient, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *automationClientService) UpdateClient(ctx context.Context, userID, clientID uint64, req models.AutomationClientRequest) (models.AutomationClient, error) {
	// Test connection with updated information
	testReq := models.ClientTestRequest{
		URL:        req.URL,
		APIKey:     req.APIKey,
		ClientType: req.ClientType,
	}

	testResp, err := s.TestClientConnection(ctx, testReq)
	if err != nil || !testResp.Success {
		return models.AutomationClient{}, fmt.Errorf("failed to connect to updated client: %v", err)
	}

	client := models.AutomationClient{
		ID:         clientID,
		UserID:     userID,
		Name:       req.Name,
		ClientType: req.ClientType,
		URL:        req.URL,
		APIKey:     req.APIKey,
		IsEnabled:  req.IsEnabled,
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return models.AutomationClient{}, err
	}

	return client, nil
}

func (s *automationClientService) DeleteClient(ctx context.Context, userID, clientID uint64) error {
	return s.repo.Delete(ctx, clientID, userID)
}

func (s *automationClientService) TestClientConnection(ctx context.Context, req models.ClientTestRequest) (models.ClientTestResponse, error) {
	switch req.ClientType {
	case models.ClientTypeRadarr:
		return s.testRadarrConnection(ctx, req.URL, req.APIKey)
	case models.ClientTypeSonarr:
		return s.testSonarrConnection(ctx, req.URL, req.APIKey)
	case models.ClientTypeLidarr:
		return s.testLidarrConnection(ctx, req.URL, req.APIKey)
	default:
		return models.ClientTestResponse{
			Success: false,
			Message: "Unsupported client type",
		}, fmt.Errorf("unsupported client type: %s", req.ClientType)
	}
}

// Client-specific test connection methods...

func (s *automationClientService) testRadarrConnection(ctx context.Context, url, apiKey string) (models.ClientTestResponse, error) {
	// Configure the Radarr client
	cfg := radarr.NewConfiguration()
	cfg.AddDefaultHeader("X-Api-Key", apiKey)
	cfg.Servers = radarr.ServerConfigurations{
		{
			URL: url,
		},
	}

	// Create client and test connection by checking system status
	client := radarr.NewAPIClient(cfg)
	_, resp, err := client.SystemAPI.GetSystemStatus(ctx).Execute()

	if err != nil {
		return models.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect to Radarr: %v", err),
		}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return models.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Radarr returned status code %d", resp.StatusCode),
		}, fmt.Errorf("radarr returned status code %d", resp.StatusCode)
	}

	return models.ClientTestResponse{
		Success: true,
		Message: "Successfully connected to Radarr",
	}, nil
}

func (s *automationClientService) testSonarrConnection(ctx context.Context, url, apiKey string) (models.ClientTestResponse, error) {
	// Configure the Sonarr client
	cfg := sonarr.NewConfiguration()
	cfg.AddDefaultHeader("X-Api-Key", apiKey)
	cfg.Servers = sonarr.ServerConfigurations{
		{
			URL: url,
		},
	}

	// Create client and test connection by checking system status
	client := sonarr.NewAPIClient(cfg)
	_, resp, err := client.SystemAPI.GetSystemStatus(ctx).Execute()

	if err != nil {
		return models.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect to Sonarr: %v", err),
		}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return models.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Sonarr returned status code %d", resp.StatusCode),
		}, fmt.Errorf("sonarr returned status code %d", resp.StatusCode)
	}

	return models.ClientTestResponse{
		Success: true,
		Message: "Successfully connected to Sonarr",
	}, nil
}

func (s *automationClientService) testLidarrConnection(ctx context.Context, url, apiKey string) (models.ClientTestResponse, error) {
	// Configure the Lidarr client
	cfg := lidarr.NewConfiguration()
	cfg.AddDefaultHeader("X-Api-Key", apiKey)
	cfg.Servers = lidarr.ServerConfigurations{
		{
			URL: url,
		},
	}

	// Create client and test connection by checking system status
	client := lidarr.NewAPIClient(cfg)
	_, resp, err := client.SystemAPI.GetSystemStatus(ctx).Execute()

	if err != nil {
		return models.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect to Lidarr: %v", err),
		}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return models.ClientTestResponse{
			Success: false,
			Message: fmt.Sprintf("Lidarr returned status code %d", resp.StatusCode),
		}, fmt.Errorf("lidarr returned status code %d", resp.StatusCode)
	}

	return models.ClientTestResponse{
		Success: true,
		Message: "Successfully connected to Lidarr",
	}, nil
}
