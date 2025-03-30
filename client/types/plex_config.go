package types

import "encoding/json"

// @Description Plex media server configuration
type PlexConfig struct {
	BaseMediaClientConfig
	Token string `json:"token" mapstructure:"token" example:"your-plex-token" binding:"required_if=Enabled true"`
}

func NewPlexConfig() PlexConfig {
	return PlexConfig{
		BaseMediaClientConfig: BaseMediaClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypePlex,
			},
			Type: MediaClientTypePlex,
		},
	}
}

func (PlexConfig) isClientConfig()      {}
func (PlexConfig) isMediaClientConfig() {}

func (PlexConfig) GetClientType() MediaClientType {
	return MediaClientTypePlex
}
func (PlexConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (PlexConfig) SupportsMovies() bool {
	return true
}
func (PlexConfig) SupportsSeries() bool {
	return true
}
func (PlexConfig) SupportsMusic() bool {
	return true
}
func (PlexConfig) SupportsPlaylists() bool {
	return true
}
func (PlexConfig) SupportsCollections() bool {
	return true
}
func (PlexConfig) SupportsHistory() bool {
	return true
}

func (c *PlexConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid recursion
	type Alias PlexConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Type is always the correct constant
	c.Type = MediaClientTypePlex
	return nil
}
