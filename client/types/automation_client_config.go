package types

type AutomationClientConfig interface {
	ClientConfig
	isAutomationClientConfig()
	GetClientType() AutomationClientType

	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
}
