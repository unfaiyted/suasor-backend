// services/download_client.go
package services

import (
	"context"
	"fmt"
	"suasor/models"
	"suasor/repository"
	"time"

	lidarr "github.com/devopsarr/lidarr-go/lidarr"
	radarr "github.com/devopsarr/radarr-go/radarr"
	sonarr "github.com/devopsarr/sonarr-go/sonarr"
)

// DownloadClientService defines the interface for download client operations
type DownloadClientService interface {
	CreateClient(ctx context.Context, userID uint64, req models.DownloadClientRequest) (models.DownloadClient, error)
	GetClientByID(ctx context.Context, userID, clientID uint64) (models.DownloadClient, error)
	GetClientsByUserID(ctx context.Context, userID uint64) ([]models.DownloadClient, error)
	UpdateClient(ctx context.Context, userID, clientID uint64, req models.DownloadClientRequest) (models.DownloadClient, error)
	DeleteClient(ctx context.Context, userID, clientID uint64) error
	TestClientConnection(ctx context.Context, req models.ClientTestRequest) (models.ClientTestResponse, error)
}

type downloadClientService struct {
	repo repository.DownloadClientRepository
}

func NewDownloadClientService(repo repository.DownloadClientRepository) DownloadClientService {
	return &downloadClientService{repo: repo}
}

func (s *downloadClientService) CreateClient(ctx context.Context, userID uint64, req models.DownloadClientRequest) (models.DownloadClient, error) {
	// Test connection first
	testReq := models.ClientTestRequest{
		URL:        req.URL,
		APIKey:     req.APIKey,
		ClientType: req.ClientType,
	}

	testResp, err := s.TestClientConnection(ctx, testReq)
	if err != nil || !testResp.Success {
		return models.DownloadClient{}, fmt.Errorf("failed to connect to client: %v", err)
	}

	client := models.DownloadClient{
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

func (s *downloadClientService) GetClientByID(ctx context.Context, userID, clientID uint64) (models.DownloadClient, error) {
	return s.repo.GetByID(ctx, clientID, userID)
}

func (s *downloadClientService) GetClientsByUserID(ctx context.Context, userID uint64) ([]models.DownloadClient, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *downloadClientService) UpdateClient(ctx context.Context, userID, clientID uint64, req models.DownloadClientRequest) (models.DownloadClient, error) {
	// Test connection with updated information
	testReq := models.ClientTestRequest{
		URL:        req.URL,
		APIKey:     req.APIKey,
		ClientType: req.ClientType,
	}

	testResp, err := s.TestClientConnection(ctx, testReq)
	if err != nil || !testResp.Success {
		return models.DownloadClient{}, fmt.Errorf("failed to connect to updated client: %v", err)
	}

	client := models.DownloadClient{
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
		return models.DownloadClient{}, err
	}

	return client, nil
}

func (s *downloadClientService) DeleteClient(ctx context.Context, userID, clientID uint64) error {
	return s.repo.Delete(ctx, clientID, userID)
}

func (s *downloadClientService) TestClientConnection(ctx context.Context, req models.ClientTestRequest) (models.ClientTestResponse, error) {
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

func (s *downloadClientService) testRadarrConnection(ctx context.Context, url, apiKey string) (models.ClientTestResponse, error) {
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

func (s *downloadClientService) testSonarrConnection(ctx context.Context, url, apiKey string) (models.ClientTestResponse, error) {
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

func (s *downloadClientService) testLidarrConnection(ctx context.Context, url, apiKey string) (models.ClientTestResponse, error) {
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
