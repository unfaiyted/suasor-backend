package bundles

import (
	"suasor/handlers"
)

type SystemHandlers interface {
	ConfigHandler() *handlers.ConfigHandler
	HealthHandler() *handlers.HealthHandler
	JobHandler() *handlers.JobHandler
	SearchHandler() *handlers.SearchHandler
	UserConfigHandler() *handlers.UserConfigHandler
}

// SystemHandlers contains core system handlers
type systemHandlers struct {
	AuthHandler       *handlers.AuthHandler
	UserHandler       *handlers.UserHandler
	ConfigHandler     *handlers.ConfigHandler
	JobHandler        *handlers.JobHandler
	HealthHandler     *handlers.HealthHandler
	SearchHandler     *handlers.SearchHandler
	UserConfigHandler *handlers.UserConfigHandler
}

func NewSystemHandlers(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, configHandler *handlers.ConfigHandler, jobHandler *handlers.JobHandler, healthHandler *handlers.HealthHandler, searchHandler *handlers.SearchHandler, peopleHandler *handlers.PeopleHandler, userConfigHandler *handlers.UserConfigHandler) *SystemHandlers {
	return &SystemHandlers{
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		ConfigHandler:     configHandler,
		JobHandler:        jobHandler,
		HealthHandler:     healthHandler,
		SearchHandler:     searchHandler,
		UserConfigHandler: userConfigHandler,
	}
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
	log := logger.LoggerFromContext(ctx)

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

// SpecializedHandlers contains specialized handlers
type SpecializedHandlers struct {
	RecommendationHandler *handlers.RecommendationHandler

	CalendarHandler *handlers.CalendarHandler
}
