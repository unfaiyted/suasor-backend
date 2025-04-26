package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// @Description Emby media server configuration
type EmbyConfig struct {
	ClientMediaConfig `json:"details"`
	UserID            string `json:"userID,omitempty" mapstructure:"userID" example:"your-internal-user-id"`
	Username          string `json:"username" mapstructure:"username" example:"admin"`
}

func NewEmbyConfig(username string, userID string, baseURL string, APIKey string, enabled bool, validateConn bool) EmbyConfig {
	clientConfig := NewClientMediaConfig(ClientMediaTypeEmby, ClientCategoryMedia, "Emby", baseURL, APIKey, enabled, validateConn)
	return EmbyConfig{
		ClientMediaConfig: clientConfig,
		Username:          username,
		UserID:            userID,
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

func (*EmbyConfig) SupportsMovies() bool {
	return true
}
func (*EmbyConfig) SupportsSeries() bool {
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
	return UnmarshalConfigJSON(data, c)
}

// Value implements driver.Valuer for database storage
func (c *EmbyConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *EmbyConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use the same custom unmarshaling logic we defined in UnmarshalJSON
	return m.UnmarshalJSON(bytes)
}
