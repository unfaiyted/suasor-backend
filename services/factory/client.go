package factory

import (
	"fmt"
	"gorm.io/gorm"
	"suasor/client"
	"suasor/client/types"
	"suasor/services"
)

// ClientServiceFactory creates appropriate service instances based on client type
type ClientServiceFactory struct {
	db *gorm.DB
}

// GetServiceForType returns the appropriate service for the client type
func (f *ClientServiceFactory) GetServiceForType(clientType string) (interface{}, error) {

	factory := client.NewClientFactoryService()
	switch clientType {
	case "jellyfin":
		return services.NewClientService[*types.JellyfinConfig](factory, f.db), nil
	case "emby":
		return services.NewClientService[*types.EmbyConfig](factory, f.db), nil
	case "plex":
		return services.NewClientService[*types.PlexConfig](factory, f.db), nil
	case "subsonic":
		return services.NewClientService[*types.SubsonicConfig](factory, f.db), nil
	case "radarr":
		return services.NewClientService[*types.RadarrConfig](factory, f.db), nil
	case "lidarr":
		return services.NewClientService[*types.LidarrConfig](factory, f.db), nil
	case "sonarr":
		return services.NewClientService[*types.SonarrConfig](factory, f.db), nil
	default:
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}
}
