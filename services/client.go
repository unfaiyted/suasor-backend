package services

import (
	"context"
	"gorm.io/gorm"
	"suasor/client"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/types/responses"
)

// type ClientService interface {
// 	// GetContentProvider(ctx context.Context, userID, clientID uint64) (types.ContentProvider, error)
// }

type ClientService[T types.ClientConfig] interface {
	Create(ctx context.Context, client models.Client[T]) (*models.Client[T], error)
	Update(ctx context.Context, client models.Client[T]) (*models.Client[T], error)
	GetByID(ctx context.Context, id uint64, userID uint64) (*models.Client[T], error)
	GetByUserID(ctx context.Context, userID uint64) ([]*models.Client[T], error)
	GetByType(ctx context.Context, clientType types.ClientType, userID uint64) ([]*models.Client[T], error)
	Delete(ctx context.Context, id uint64, userID uint64) error
	TestConnection(ctx context.Context, clientID uint64, config *T) (responses.ClientTestResponse, error)
}

// ClientService handles business logic for clients with specific config types
type clientService[T types.ClientConfig] struct {
	repo    repository.ClientRepository[T]
	factory client.ClientFactoryService
	// Other dependencies like validators, API clients, etc.
}

// NewClientService creates a service for a specific client type
func NewClientService[T types.ClientConfig](factory *client.ClientFactoryService, db *gorm.DB) *clientService[T] {
	return &clientService[T]{
		repo:    repository.NewClientRepository[T](db),
		factory: *factory,
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

func (s *clientService[T]) GetByType(ctx context.Context, clientType types.ClientType, userID uint64) ([]*models.Client[T], error) {
	return s.repo.GetByType(ctx, clientType, userID)
}

func (s *clientService[T]) Update(ctx context.Context, client models.Client[T]) (*models.Client[T], error) {
	return s.repo.Update(ctx, client)
}

func (s *clientService[T]) Delete(ctx context.Context, id uint64, userID uint64) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *clientService[T]) TestConnection(ctx context.Context, clientID uint64, config *T) (responses.ClientTestResponse, error) {
	// Get client from factory
	c, err := s.factory.GetClient(ctx, clientID, *config)
	if err != nil {
		return responses.ClientTestResponse{
			Success: false,
			Message: "Failed to create client",
		}, err
	}

	testResult, err := c.TestConnection(ctx)
	if err != nil {
		return responses.ClientTestResponse{
			Success: false,
			Message: "Failed to test connection",
		}, err
	}
	return responses.ClientTestResponse{
		Success: testResult,
		Message: "Successfully connected to " + c.GetType().String(),
	}, nil

}
