package types

import "encoding/json"

// @Description Jellyfin media server configuration
type JellyfinConfig struct {
	ClientMediaConfig
	UserID   string `json:"userID,omitempty" mapstructure:"userID" example:"your-internal-user-id"`
	Username string `json:"username" mapstructure:"username" example:"admin"`
}

func NewJellyfinConfig(username string, userID string, baseURL string, APIKey string, enabled bool, validateConn bool) JellyfinConfig {
	clientConfig := NewClientMediaConfig(ClientMediaTypeJellyfin, ClientCategoryMedia, "Jellyfin", baseURL, APIKey, enabled, validateConn)
	return JellyfinConfig{
		ClientMediaConfig: clientConfig,
		Username:          username,
		UserID:            userID,
	}
}

func (c *JellyfinConfig) GetUsername() string {
	return c.Username
}

func (c *JellyfinConfig) GetUserID() string {
	return c.UserID
}

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

	return nil
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
