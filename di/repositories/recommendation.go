// app/di/repositories/recommendation.go
package repositories

import (
	"context"
	"suasor/di/container"
	"suasor/repository"

	"gorm.io/gorm"
)

// RegisterRecommendationRepository registers the recommendation repository
func RegisterRecommendationRepository(ctx context.Context, c *container.Container) {

	container.RegisterFactory[repository.RecommendationRepository](c, func(c *container.Container) repository.RecommendationRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewRecommendationRepository(db)
	})

}
