package types

type ClientConfig interface {
	isClientConfig()
	GetCategory() ClientCategory
	SetCategory(ClientCategory)
	GetType() ClientType
}

type BaseClientConfig struct {
	Type     ClientType     `json:"type"`
	Category ClientCategory `json:"category"`
}

func (c *BaseClientConfig) GetType() ClientType {
	return c.Type
}

func (c *BaseClientConfig) GetCategory() ClientCategory {
	return c.Category
}
func (c *BaseClientConfig) SetCategory(category ClientCategory) {
	c.Category = category
}

func (BaseClientConfig) isClientConfig() {}
