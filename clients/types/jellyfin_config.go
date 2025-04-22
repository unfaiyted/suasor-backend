package types

import "encoding/json"

// @Description Jellyfin media server configuration
type JellyfinConfig struct {
	BaseClientMediaConfig
	UserID   string `json:"userID,omitempty" mapstructure:"userID" example:"your-internal-user-id"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
}

// func (c *JellyfinConfig) GetName() string {
// 	return c.Name
// }

func (c *JellyfinConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid recursion
	type Alias JellyfinConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Type is always the correct constant
	c.ClientType = ClientMediaTypeJellyfin
	c.Type = ClientTypeJellyfin
	return nil
}

func NewJellyfinConfig() JellyfinConfig {
	return JellyfinConfig{
		BaseClientMediaConfig: BaseClientMediaConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeJellyfin,
			},
			ClientType: ClientMediaTypeJellyfin,
		},
	}
}

func (JellyfinConfig) GetClientType() ClientMediaType {
	return ClientMediaTypeJellyfin
}
func (JellyfinConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (JellyfinConfig) SupportsMovies() bool {
	return true
}
func (JellyfinConfig) SupportsSeries() bool {
	return true
}
func (JellyfinConfig) SupportsMusic() bool {
	return true
}

func (JellyfinConfig) SupportsPlaylists() bool {
	return true
}
func (JellyfinConfig) SupportsCollections() bool {
	return true
}
func (JellyfinConfig) SupportsHistory() bool {
	return true
}

func (j *JellyfinConfig) GetBaseURL() string {
	return j.BaseURL
}
