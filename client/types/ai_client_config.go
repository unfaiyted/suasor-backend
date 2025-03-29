package types

type AIClientConfig interface {
	ClientConfig
	isAutomationClientConfig()
	GetClientType() AutomationClientType
}
