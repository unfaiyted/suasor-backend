package lidarr

import (
	"strconv"

	lidarr "github.com/devopsarr/lidarr-go/lidarr"

	"suasor/client/automation/types"
)

// getExternalIDs extracts all available external IDs from a Lidarr artist resource
func (l *LidarrClient) getArtistExternalIDs(artist *lidarr.ArtistResource) []types.ExternalID {
	var ids []types.ExternalID

	// Add Lidarr's internal ID
	ids = append(ids, types.ExternalID{
		Source: "lidarr",
		Value:  strconv.FormatInt(int64(artist.GetId()), 10),
	})

	// Add MusicBrainz ID (the primary ID in Lidarr)
	if artist.ForeignArtistId.Get() != nil && *artist.ForeignArtistId.Get() != "" {
		ids = append(ids, types.ExternalID{
			Source: "musicbrainz",
			Value:  *artist.ForeignArtistId.Get(),
		})
	}

	// Add MusicBrainz ID (alternative field)
	if artist.MbId.Get() != nil && *artist.MbId.Get() != "" {
		ids = append(ids, types.ExternalID{
			Source: "musicbrainz",
			Value:  *artist.MbId.Get(),
		})
	}

	// Add The Audio DB ID
	if artist.GetTadbId() != 0 {
		ids = append(ids, types.ExternalID{
			Source: "audiodb",
			Value:  strconv.FormatInt(int64(artist.GetTadbId()), 10),
		})
	}

	// Add Discogs ID
	if artist.GetDiscogsId() != 0 {
		ids = append(ids, types.ExternalID{
			Source: "discogs",
			Value:  strconv.FormatInt(int64(artist.GetDiscogsId()), 10),
		})
	}

	// Add AllMusic ID
	if artist.AllMusicId.Get() != nil && *artist.AllMusicId.Get() != "" {
		ids = append(ids, types.ExternalID{
			Source: "allmusic",
			Value:  *artist.AllMusicId.Get(),
		})
	}

	return ids
}

// getAlbumExternalIDs extracts all available external IDs from a Lidarr album resource
func (l *LidarrClient) getAlbumExternalIDs(album *lidarr.AlbumResource) []types.ExternalID {
	var ids []types.ExternalID

	// Add Lidarr's internal album ID
	ids = append(ids, types.ExternalID{
		Source: "lidarr_album",
		Value:  strconv.FormatInt(int64(album.GetId()), 10),
	})

	// Add artist ID reference
	if album.GetArtistId() != 0 {
		ids = append(ids, types.ExternalID{
			Source: "lidarr_artist",
			Value:  strconv.FormatInt(int64(album.GetArtistId()), 10),
		})
	}

	// Add MusicBrainz album ID
	if album.ForeignAlbumId.Get() != nil && *album.ForeignAlbumId.Get() != "" {
		ids = append(ids, types.ExternalID{
			Source: "musicbrainz_album",
			Value:  *album.ForeignAlbumId.Get(),
		})
	}

	// Extract IDs from releases if available
	for _, release := range album.GetReleases() {
		if release.ForeignReleaseId.Get() != nil && *release.ForeignReleaseId.Get() != "" {
			ids = append(ids, types.ExternalID{
				Source: "musicbrainz_release",
				Value:  *release.ForeignReleaseId.Get(),
			})
		}
	}

	// If artist object is included, extract artist IDs too
	if album.Artist != nil {
		artistIDs := l.getArtistExternalIDs(album.Artist)
		ids = append(ids, artistIDs...)
	}

	return ids
}
