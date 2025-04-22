// app/di/repositories/client.go
package repositories

import (
	"gorm.io/gorm"
	"suasor/client"
	"suasor/client/types"
	"suasor/container"
	"suasor/repository"
	"suasor/services"
)

// ProvideEmbyClientRepository provides a ClientRepository for EmbyConfig
func ProvideEmbyClientRepository(c *container.Container) repository.ClientRepository[*types.EmbyConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.EmbyConfig](db)
}

// ProvideJellyfinClientRepository provides a ClientRepository for JellyfinConfig
func ProvideJellyfinClientRepository(c *container.Container) repository.ClientRepository[*types.JellyfinConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.JellyfinConfig](db)
}

// ProvidePlexClientRepository provides a ClientRepository for PlexConfig
func ProvidePlexClientRepository(c *container.Container) repository.ClientRepository[*types.PlexConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.PlexConfig](db)
}

// ProvideSubsonicClientRepository provides a ClientRepository for SubsonicConfig
func ProvideSubsonicClientRepository(c *container.Container) repository.ClientRepository[*types.SubsonicConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.SubsonicConfig](db)
}

// ProvideSonarrClientRepository provides a ClientRepository for SonarrConfig
func ProvideSonarrClientRepository(c *container.Container) repository.ClientRepository[*types.SonarrConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.SonarrConfig](db)
}

// ProvideRadarrClientRepository provides a ClientRepository for RadarrConfig
func ProvideRadarrClientRepository(c *container.Container) repository.ClientRepository[*types.RadarrConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.RadarrConfig](db)
}

// ProvideLidarrClientRepository provides a ClientRepository for LidarrConfig
func ProvideLidarrClientRepository(c *container.Container) repository.ClientRepository[*types.LidarrConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.LidarrConfig](db)
}

// ProvideClaudeClientRepository provides a ClientRepository for ClaudeConfig
func ProvideClaudeClientRepository(c *container.Container) repository.ClientRepository[*types.ClaudeConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.ClaudeConfig](db)
}

// ProvideOpenAIClientRepository provides a ClientRepository for OpenAIConfig
func ProvideOpenAIClientRepository(c *container.Container) repository.ClientRepository[*types.OpenAIConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.OpenAIConfig](db)
}

// ProvideOllamaClientRepository provides a ClientRepository for OllamaConfig
func ProvideOllamaClientRepository(c *container.Container) repository.ClientRepository[*types.OllamaConfig] {
	db := container.MustGet[*gorm.DB](c)
	return repository.NewClientRepository[*types.OllamaConfig](db)
}

// ProvideMetadataClientService provides a MetadataClientService for TMDBConfig
func ProvideMetadataClientService(c *container.Container) *services.MetadataClientService[*types.TMDBConfig] {
	factory := container.MustGet[*client.ClientFactoryService](c)
	db := container.MustGet[*gorm.DB](c)
	repo := repository.NewClientRepository[*types.TMDBConfig](db)
	return services.NewMetadataClientService(factory, repo)
}

