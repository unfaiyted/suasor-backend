package bundles

// ApplicationHandlers contains all handlers organized by category
type ApplicationHandlers struct {
	System      SystemHandlers
	Media       MediaHandlers
	MediaData   MediaDataHandlers
	Clients     ClientHandlers
	Specialized SpecializedHandlers
}

func NewApplicationHandlers(
	systemHandlers SystemHandlers,
	mediaHandlers MediaHandlers,
	mediaDataHandlers MediaDataHandlers,
	clientsHandlers ClientHandlers,
	specializedHandlers SpecializedHandlers,
) ApplicationHandlers {
	return ApplicationHandlers{
		System:      systemHandlers,
		Media:       mediaHandlers,
		MediaData:   mediaDataHandlers,
		Clients:     clientsHandlers,
		Specialized: specializedHandlers,
	}
}
