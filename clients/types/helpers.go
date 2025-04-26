package types

import (
	"suasor/clients/media/types"
)

func DoesClientSupportMediaType(config ClientMediaConfig, mediaType types.MediaType) bool {
	switch mediaType {
	case types.MediaTypeMovie:
		return config.SupportsMovies()
	case types.MediaTypeSeries, types.MediaTypeSeason, types.MediaTypeEpisode:
		return config.SupportsSeries()
	case types.MediaTypeTrack, types.MediaTypeArtist, types.MediaTypeAlbum:
		return config.SupportsMusic()
	case types.MediaTypePlaylist:
		return config.SupportsPlaylists()
	case types.MediaTypeCollection:
		return config.SupportsCollections()
	default:
		return false
	}
}
