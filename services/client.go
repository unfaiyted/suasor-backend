package services

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/types/responses"
)

// type ClientService interface {
// 	// GetContentProvider(ctx context.Context, userID, clientID uint64) (types.ContentProvider, error)
// }

type ClientService[T types.ClientConfig] interface {
	Create(ctx context.Context, client models.Client[T]) (models.Client[T], error)
	Update(ctx context.Context, client types.ClientConfig) (models.Client[T], error)
	GetByID(ctx context.Context, id uint64, userID uint64) (*models.Client[T], error)
	GetByUserID(ctx context.Context, userID uint64) ([]models.Client[T], error)
	GetByType(ctx context.Context, clientType types.ClientType, userID uint64) ([]models.Client[T], error)
	Delete(ctx context.Context, id uint64, userID uint64) error
	TestConnection(ctx context.Context, config T) (responses.ClientTestResponse, error)
}

// ClientService handles business logic for clients with specific config types
type clientService[T types.ClientConfig] struct {
	repo repository.ClientRepository[T]
	// Other dependencies like validators, API clients, etc.
}

// NewClientService creates a service for a specific client type
func NewClientService[T types.ClientConfig](db *gorm.DB) *clientService[T] {
	return &clientService[T]{
		repo: repository.NewClientRepository[T](db),
	}
}

// Create handles client creation with business logic
func (s *clientService[T]) Create(ctx context.Context, client models.Client[T]) (*models.Client[T], error) {
	return s.repo.Create(ctx, client)
}

func (s *clientService[T]) GetByID(ctx context.Context, id uint64, userID uint64) (*models.Client[T], error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *clientService[T]) GetByUserID(ctx context.Context, userID uint64) ([]*models.Client[T], error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *clientService[T]) Update(ctx context.Context, client models.Client[T]) (*models.Client[T], error) {
	return s.repo.Update(ctx, client)
}

func (s *clientService[T]) Delete(ctx context.Context, id uint64, userID uint64) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *clientService[T]) TestConnection(ctx context.Context, config T) (responses.ClientTestResponse, error) {
	switch any(config).(type) {
	case types.JellyfinConfig:
		// return testJellyfinConnection(ctx, config)
	case types.EmbyConfig:
		// return testEmbyConnection(ctx, config)
	case types.SubsonicConfig:
		//		return testSubsonicConnection(ctx, config)
	default:
		return responses.ClientTestResponse{
			Success: false,
			Message: "Unsupported client type",
		}, fmt.Errorf("unsupported client type: %s", config)
	}
	return responses.ClientTestResponse{
		Success: false,
		Message: "Error processing client type",
	}, nil
}

// switch config.Data.GetClientType() {
// case types.ClientTypePlex:
// 	return s.testPlexConnection(ctx, config.Data)
// case types.ClientTypeJellyfin:
// 	return s.testJellyfinConnection(ctx, config.Data)
// case types.ClientTypeEmby:
// 	return s.testEmbyConnection(ctx, config.Data)
// case types.ClientTypeSubsonic:
// 	return s.testSubsonicConnection(ctx, config.Data)
// default:
// 	return types.ClientTestResponse{
// 		Success: false,
// 		Message: "Unsupported client type",
// 	}, fmt.Errorf("unsupported client type: %s", config.Data.GetClientType())
