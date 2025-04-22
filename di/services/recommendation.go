// app/di/services/recommendation.go
package services

import (
	"context"
	"suasor/di/container"
	"suasor/repository"
	"suasor/services"
)

// RegisterRecommendationService registers the recommendation service
func registerRecommendationService(ctx context.Context, c *container.Container) {
	container.RegisterFactory[services.RecommendationService](c, func(c *container.Container) services.RecommendationService {
		recommendationRepo := container.MustGet[repository.RecommendationRepository](c)
		return services.NewRecommendationService(recommendationRepo)
	})
}
