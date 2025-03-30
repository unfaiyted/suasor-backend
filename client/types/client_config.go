package types

type ClientConfig interface {
	GetName() string
	isClientConfig()
	GetCategory() ClientCategory
	GetType() ClientType
}

type BaseClientConfig struct {
	Type ClientType
	Name string
}

func (c *BaseClientConfig) GetType() ClientType {
	return c.Type
}

func (c *BaseClientConfig) GetCategory() ClientCategory {
	return ClientCategoryUnknown
}

func (BaseClientConfig) isClientConfig() {}

func (c *BaseClientConfig) GetName() string {
	return c.Name
}
