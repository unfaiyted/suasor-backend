// app/di/handlers/recommendation.go
package handlers

import (
	"context"
	"suasor/app/container"
	"suasor/handlers"
	"suasor/services"
)

// RegisterRecommendationHandlers registers recommendation-related handlers
func RegisterRecommendationHandlers(ctx context.Context, c *container.Container) {
	container.RegisterFactory[*handlers.RecommendationHandler](c, func(c *container.Container) *handlers.RecommendationHandler {
		recommendationService := container.MustGet[services.RecommendationService](c)
		return handlers.NewRecommendationHandler(recommendationService)
	})
}