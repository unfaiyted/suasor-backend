package factory

import (
	"fmt"
	"gorm.io/gorm"
	"suasor/client/types"
	"suasor/services"
)

// ClientServiceFactory creates appropriate service instances based on client type
type ClientServiceFactory struct {
	db *gorm.DB
}

// GetServiceForType returns the appropriate service for the client type
func (f *ClientServiceFactory) GetServiceForType(clientType string) (interface{}, error) {
	switch clientType {
	case "jellyfin":
		return services.NewClientService[types.JellyfinConfig](f.db), nil
	case "emby":
		return services.NewClientService[types.EmbyConfig](f.db), nil
	case "plex":
		return services.NewClientService[types.PlexConfig](f.db), nil
	case "subsonic":
		return services.NewClientService[types.SubsonicConfig](f.db), nil
	case "radarr":
		return services.NewClientService[types.RadarrConfig](f.db), nil
	case "lidarr":
		return services.NewClientService[types.LidarrConfig](f.db), nil
	case "sonarr":
		return services.NewClientService[types.SonarrConfig](f.db), nil
	default:
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}
}
