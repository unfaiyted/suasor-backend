package types

type ClientConfig interface {
	isClientConfig()
	GetCategory() ClientCategory
	GetType() ClientType
}

type BaseClientConfig struct {
	Type ClientType `json:"type"`
}

func (c *BaseClientConfig) GetType() ClientType {
	return c.Type
}

func (c *BaseClientConfig) GetCategory() ClientCategory {
	return ClientCategoryUnknown
}

func (BaseClientConfig) isClientConfig() {}
