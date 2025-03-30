package types

import "encoding/json"

// @Description Supersonic music server configuration
type SubsonicConfig struct {
	BaseMediaClientConfig
	Username string `json:"username" mapstructure:"username" example:"admin" binding:"required_if=Enabled true"`
	Password string `json:"password" mapstructure:"password" example:"your-password" binding:"required_if=Enabled true"`
}

func NewSubsonicConfig() SubsonicConfig {
	return SubsonicConfig{
		BaseMediaClientConfig: BaseMediaClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeSubsonic,
			},
			Type: MediaClientTypeSubsonic,
		},
	}
}

func (SubsonicConfig) isClientConfig()      {}
func (SubsonicConfig) isMediaClientConfig() {}
func (SubsonicConfig) GetClientType() MediaClientType {
	return MediaClientTypeSubsonic
}
func (SubsonicConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (SubsonicConfig) SupportsMovies() bool {
	return false
}
func (SubsonicConfig) SupportsSeries() bool {
	return false
}
func (SubsonicConfig) SupportsMusic() bool {
	return true
}
func (SubsonicConfig) SupportsPlaylists() bool {
	return true
}
func (SubsonicConfig) SupportsCollections() bool {
	return false
}
func (SubsonicConfig) SupportsHistory() bool {
	return true
}

func (c *SubsonicConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid recursion
	type Alias SubsonicConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Type is always the correct constant
	c.Type = MediaClientTypeSubsonic
	return nil
}
