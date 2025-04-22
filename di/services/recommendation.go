// app/di/services/recommendation.go
package services

import (
	"context"
	"suasor/container"
	"suasor/repository"
	"suasor/services"
)

// RegisterRecommendationService registers the recommendation service
func RegisterRecommendationService(ctx context.Context, c *container.Container) {
	container.RegisterFactory[services.RecommendationService](c, func(c *container.Container) services.RecommendationService {
		recommendationRepo := container.MustGet[repository.RecommendationRepository](c)
		return services.NewRecommendationService(recommendationRepo)
	})
}

