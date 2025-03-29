package types

type ClientConfig interface {
	isClientConfig()
	GetCategory() ClientCategory
}
