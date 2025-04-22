// app/di/services/search.go
package services

import (
	"context"
	"suasor/clients"
	"suasor/di/container"
	"suasor/repository"
	apprepos "suasor/repository/bundles"
	"suasor/services"
)

// RegisterSearchService registers the search service
func registerSearchService(ctx context.Context, c *container.Container) {
	container.RegisterFactory[services.SearchService](c, func(c *container.Container) services.SearchService {
		searchRepo := container.MustGet[repository.SearchRepository](c)
		clientRepos := container.MustGet[apprepos.ClientRepositories](c)
		itemRepos := container.MustGet[apprepos.CoreMediaItemRepositories](c)
		personRepo := container.MustGet[repository.PersonRepository](c)
		clientFactoryService := container.MustGet[*clients.ClientProviderFactoryService](c)
		return services.NewSearchService(
			searchRepo,
			clientRepos,
			itemRepos,
			personRepo,
			clientFactoryService,
		)
	})
}
