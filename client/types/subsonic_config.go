package types

// @Description Supersonic music server configuration
type SubsonicConfig struct {
	Enabled  bool            `json:"enabled" mapstructure:"enabled" example:"false"`
	Host     string          `json:"host" mapstructure:"host" example:"localhost" binding:"required_if=Enabled true"`
	Port     int             `json:"port" mapstructure:"port" example:"4533" binding:"required_if=Enabled true"`
	Username string          `json:"username" mapstructure:"username" example:"admin" binding:"required_if=Enabled true"`
	Password string          `json:"password" mapstructure:"password" example:"your-password" binding:"required_if=Enabled true"`
	SSL      bool            `json:"ssl" mapstructure:"ssl" example:"false"`
	Type     MediaClientType `json:"type" mapstructure:"type" default:"subsonic"`
}

func NewSubsonicConfig() SubsonicConfig {
	return SubsonicConfig{
		Type: MediaClientTypeSubsonic,
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
