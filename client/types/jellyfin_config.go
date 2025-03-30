package types

import "encoding/json"

// @Description Jellyfin media server configuration
type JellyfinConfig struct {
	BaseMediaClientConfig
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
	c.Type = MediaClientTypeJellyfin
	return nil
}

func NewJellyfinConfig() JellyfinConfig {
	return JellyfinConfig{
		BaseMediaClientConfig: BaseMediaClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeJellyfin,
			},
			Type: MediaClientTypeJellyfin,
		},
	}
}

func (JellyfinConfig) GetClientType() MediaClientType {
	return MediaClientTypeJellyfin
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
