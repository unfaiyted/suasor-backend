package types

// @Description Plex media server configuration
type PlexConfig struct {
	Enabled bool            `json:"enabled" mapstructure:"enabled" example:"false"`
	Host    string          `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port    int             `json:"port" mapstructure:"port" example:"32400" binding:"required_if=Enabled true"`
	Token   string          `json:"token" mapstructure:"token" example:"your-plex-token" binding:"required_if=Enabled true"`
	SSL     bool            `json:"ssl" mapstructure:"ssl" example:"false"`
	Type    MediaClientType `json:"type" mapstructure:"type" default:"plex"`
}

func NewPlexConfig() PlexConfig {
	return PlexConfig{
		Type: MediaClientTypePlex,
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
