package lidarr

import (
	lidarr "github.com/devopsarr/lidarr-go/lidarr"
	"suasor/clients/automation/types"
	"suasor/types/models"
)

func (l *LidarrClient) convertArtistToMediaItem(artist *lidarr.ArtistResource) models.AutomationMediaItem[types.AutomationData] {
	// Convert images
	images := make([]types.AutomationMediaImage, 0, len(artist.GetImages()))
	for _, img := range artist.GetImages() {
		images = append(images, types.AutomationMediaImage{
			URL:       img.GetRemoteUrl(),
			CoverType: string(img.GetCoverType()),
		})
	}

	// Get quality profile
	qualityProfile := types.QualityProfileSummary{
		ID:   int64(artist.GetQualityProfileId()),
		Name: "", // We don't have the name in the artist object
	}

	// Get metadata profile if available
	metadataProfile := types.MetadataProfile{
		ID:   artist.GetMetadataProfileId(),
		Name: "", // We don't have the name in the artist object
	}

	// Convert genres
	genres := artist.GetGenres()

	// Use start year or end year as appropriate
	// var releaseDate time.Time

	releaseDate := *artist.LastAlbum.ReleaseDate.Get()

	return models.AutomationMediaItem[types.AutomationData]{
		ID:       uint64(artist.GetId()),
		Title:    artist.GetArtistName(),
		Overview: artist.GetOverview(),
		Type:     types.AUTOMEDIATYPE_ARTIST,
		// TODO: get the first album the arstist release, set a year
		// Year:             artist.GetYearFormed(),
		AddedAt:          artist.GetAdded(),
		Status:           types.GetStatusFromMusicStatus(artist.GetStatus()),
		Path:             artist.GetPath(),
		QualityProfile:   qualityProfile,
		Images:           images,
		Genres:           genres,
		ExternalIDs:      l.getArtistExternalIDs(artist),
		DownloadedStatus: determineDownloadStatus(artist.GetStatistics()),
		Monitored:        artist.GetMonitored(),
		Data: types.AutomationArtist{
			MetadataProfile:       metadataProfile,
			MostRecentReleaseDate: releaseDate,
		},
	}
}
