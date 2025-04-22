// app/di/init.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/handlers"
	"suasor/services"
	"suasor/utils"
)

// Initialize registers all dependencies in the container
func Initialize(ctx context.Context, db *gorm.DB, configService services.ConfigService) *container.Container {
	// Create a new container
	c := container.NewContainer()
	log := utils.LoggerFromContext(ctx)

	// Register core dependencies
	log.Info().Msg("Registering core dependencies")
	RegisterCore(ctx, c, db, configService)

	// Register repositories
	log.Info().Msg("Registering repositories")
	RegisterRepositories(ctx, c)

	// Register media data factory and repositories
	log.Info().Msg("Registering media data")
	RegisterMediaData(ctx, c)

	// Register services
	log.Info().Msg("Registering services")
	RegisterServices(ctx, c)

	// Register handlers
	log.Info().Msg("Registering handlers")
	RegisterHandlers(ctx, c)

	return c
}

// SystemHandlers contains core system handlers
type SystemHandlers struct {
	AuthHandler       *handlers.AuthHandler
	UserHandler       *handlers.UserHandler
	ConfigHandler     *handlers.ConfigHandler
	JobHandler        *handlers.JobHandler
	HealthHandler     *handlers.HealthHandler
	SearchHandler     *handlers.SearchHandler
	PeopleHandler     *handlers.PeopleHandler
	UserConfigHandler *handlers.UserConfigHandler
}

// GetSystemHandlers returns all system handlers from the container
func GetSystemHandlers(c *container.Container) (*SystemHandlers, error) {
	authHandler, err := container.Get[*handlers.AuthHandler](c)
	if err != nil {
		return nil, err
	}

	userHandler, err := container.Get[*handlers.UserHandler](c)
	if err != nil {
		return nil, err
	}

	configHandler, err := container.Get[*handlers.ConfigHandler](c)
	if err != nil {
		return nil, err
	}

	jobHandler, err := container.Get[*handlers.JobHandler](c)
	if err != nil {
		return nil, err
	}

	healthHandler, err := container.Get[*handlers.HealthHandler](c)
	if err != nil {
		return nil, err
	}

	searchHandler, err := container.Get[*handlers.SearchHandler](c)
	if err != nil {
		return nil, err
	}

	peopleHandler, err := container.Get[*handlers.PeopleHandler](c)
	if err != nil {
		return nil, err
	}

	userConfigHandler, err := container.Get[*handlers.UserConfigHandler](c)
	if err != nil {
		return nil, err
	}

	return &SystemHandlers{
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		ConfigHandler:     configHandler,
		JobHandler:        jobHandler,
		HealthHandler:     healthHandler,
		SearchHandler:     searchHandler,
		PeopleHandler:     peopleHandler,
		UserConfigHandler: userConfigHandler,
	}, nil
}

// GetMediaHandlers returns all media handlers from the container
func GetMediaHandlers(c *container.Container) (*MediaHandlers, error) {
	// Implementation simplified for brevity
	// In a real implementation, we would get these handlers from the container
	return &MediaHandlers{}, nil
}

// GetMediaDataHandlers returns all media data handlers from the container
func GetMediaDataHandlers(c *container.Container) (*MediaDataHandlers, error) {
	// Implementation simplified for brevity
	// In a real implementation, we would get these handlers from the container
	return &MediaDataHandlers{}, nil
}

// GetClientHandlers returns all client handlers from the container
func GetClientHandlers(c *container.Container) (*ClientHandlers, error) {
	// Implementation simplified for brevity
	// In a real implementation, we would get these handlers from the container
	return &ClientHandlers{}, nil
}

// GetSpecializedHandlers returns all specialized handlers from the container
func GetSpecializedHandlers(c *container.Container) (*SpecializedHandlers, error) {
	// Get recommendation handler
	recommendationHandler, err := container.Get[*handlers.RecommendationHandler](c)
	if err != nil {
		return nil, err
	}

	// For simplicity, we're not handling credit and calendar handlers yet
	return &SpecializedHandlers{
		RecommendationHandler: recommendationHandler,
		// CreditHandler and CalendarHandler would be retrieved similarly
	}, nil
}

// GetAllHandlers returns all application handlers organized by category
func GetAllHandlers(ctx context.Context, c *container.Container) (*ApplicationHandlers, error) {
	log := utils.LoggerFromContext(ctx)
	
	// Get system handlers
	log.Info().Msg("Getting system handlers")
	systemHandlers, err := GetSystemHandlers(c)
	if err != nil {
		return nil, err
	}

	// Get media handlers
	log.Info().Msg("Getting media handlers")
	mediaHandlers, err := GetMediaHandlers(c)
	if err != nil {
		return nil, err
	}

	// Get media data handlers
	log.Info().Msg("Getting media data handlers")
	mediaDataHandlers, err := GetMediaDataHandlers(c)
	if err != nil {
		return nil, err
	}

	// Get client handlers
	log.Info().Msg("Getting client handlers")
	clientHandlers, err := GetClientHandlers(c)
	if err != nil {
		return nil, err
	}

	// Get specialized handlers
	log.Info().Msg("Getting specialized handlers")
	specializedHandlers, err := GetSpecializedHandlers(c)
	if err != nil {
		return nil, err
	}

	return &ApplicationHandlers{
		System:      *systemHandlers,
		Media:       *mediaHandlers,
		MediaData:   *mediaDataHandlers,
		Clients:     *clientHandlers,
		Specialized: *specializedHandlers,
	}, nil
}

// MediaHandlers contains media-related handlers
type MediaHandlers struct {
	// Core Handlers for media items
	MovieHandler   handlers.UserMediaItemHandler
	SeriesHandler  handlers.UserMediaItemHandler
	SeasonHandler  handlers.UserMediaItemHandler
	EpisodeHandler handlers.UserMediaItemHandler
	TrackHandler   handlers.UserMediaItemHandler
	AlbumHandler   handlers.UserMediaItemHandler
	ArtistHandler  handlers.UserMediaItemHandler
	
	// List Handlers
	PlaylistHandler   handlers.UserListHandler
	CollectionHandler handlers.UserListHandler
}

// MediaDataHandlers contains handlers for media item data
type MediaDataHandlers struct {
	MovieDataHandler      handlers.UserMediaItemDataHandler
	SeriesDataHandler     handlers.UserMediaItemDataHandler
	SeasonDataHandler     handlers.UserMediaItemDataHandler
	EpisodeDataHandler    handlers.UserMediaItemDataHandler
	TrackDataHandler      handlers.UserMediaItemDataHandler
	AlbumDataHandler      handlers.UserMediaItemDataHandler
	ArtistDataHandler     handlers.UserMediaItemDataHandler
	PlaylistDataHandler   handlers.UserMediaItemDataHandler
	CollectionDataHandler handlers.UserMediaItemDataHandler
}

// ClientHandlers contains client-related handlers
type ClientHandlers struct {
	// Master Handler
	ClientsHandler *handlers.ClientsHandler
	
	// Media Clients
	EmbyHandler     *handlers.ClientHandler
	JellyfinHandler *handlers.ClientHandler
	PlexHandler     *handlers.ClientHandler
	SubsonicHandler *handlers.ClientHandler
	
	// Automation Clients
	RadarrHandler  *handlers.ClientHandler
	SonarrHandler  *handlers.ClientHandler
	LidarrHandler  *handlers.ClientHandler
	
	// AI Clients
	AIHandler     *handlers.AIHandler
	ClaudeHandler *handlers.ClientHandler
	OpenAIHandler *handlers.ClientHandler
	OllamaHandler *handlers.ClientHandler
	
	// Metadata
	MetadataHandler *handlers.MetadataClientHandler
}

// SpecializedHandlers contains specialized handlers
type SpecializedHandlers struct {
	RecommendationHandler *handlers.RecommendationHandler
	CreditHandler         *handlers.CreditHandler
	CalendarHandler       *handlers.CalendarHandler
}

// ApplicationHandlers contains all handlers organized by category
type ApplicationHandlers struct {
	System      SystemHandlers
	Media       MediaHandlers
	MediaData   MediaDataHandlers
	Clients     ClientHandlers
	Specialized SpecializedHandlers
}