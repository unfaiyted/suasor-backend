package types

import "encoding/json"

// @Description Emby media server configuration
type EmbyConfig struct {
	BaseClientMediaConfig
	UserID   string `json:"userID,omitempty" mapstructure:"userID" example:"your-internal-user-id"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
}

func NewEmbyConfig() EmbyConfig {
	return EmbyConfig{
		BaseClientMediaConfig: BaseClientMediaConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeEmby,
			},
			ClientType: ClientMediaTypeEmby,
		},
	}
}

func (c *EmbyConfig) GetUsername() string {
	return c.Username
}

func (c *EmbyConfig) GetUserID() string {
	return c.UserID
}

func (EmbyConfig) GetClientType() ClientMediaType {
	return ClientMediaTypeEmby
}
func (EmbyConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (EmbyConfig) SupportsMovies() bool {
	return true
}
func (EmbyConfig) SupportsSeries() bool {
	return true
}
func (EmbyConfig) SupportsMusic() bool {
	return true
}
func (EmbyConfig) SupportsPlaylists() bool {
	return true
}
func (EmbyConfig) SupportsCollections() bool {
	return true
}
func (EmbyConfig) SupportsHistory() bool {
	return true
}

func (c *EmbyConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid recursion
	type Alias EmbyConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Type is always the correct constant
	c.ClientType = ClientMediaTypeEmby
	c.Type = ClientTypeEmby
	return nil
}
