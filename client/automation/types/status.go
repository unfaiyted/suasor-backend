package types

import (
	lidarr "github.com/devopsarr/lidarr-go/lidarr"
	radarr "github.com/devopsarr/radarr-go/radarr"
	sonarr "github.com/devopsarr/sonarr-go/sonarr"
)

type DownloadedStatus string

const (
	DOWNLOADEDSTATUS_COMPLETE  DownloadedStatus = "complete"
	DOWNLOADEDSTATUS_REQUESTED DownloadedStatus = "requested"
	DOWNLOADEDSTATUS_PARTIAL   DownloadedStatus = "partial"
	DOWNLOADEDSTATUS_NONE      DownloadedStatus = "none"
)

type AutomationMediaType string

const (
	AUTOMEDIATYPE_MOVIE   AutomationMediaType = "movie"
	AUTOMEDIATYPE_SERIES  AutomationMediaType = "series"
	AUTOMEDIATYPE_ARTIST  AutomationMediaType = "artist"
	AUTOMEDIATYPE_ALBUM   AutomationMediaType = "album"
	AUTOMEDIATYPE_TRACK   AutomationMediaType = "track"
	AUTOMEDIATYPE_EPISODE AutomationMediaType = "episode"
)

type AutomationStatusType string

const (
	AUTOSTATUSTYPE_CONTINUING AutomationStatusType = "continuing"
	AUTOSTATUSTYPE_ENDED      AutomationStatusType = "ended"
	AUTOSTATUSTYPE_UPCOMING   AutomationStatusType = "upcoming"
	AUTOSTATUSTYPE_DELETED    AutomationStatusType = "deleted"
	AUTOSTATUSTYPE_IN_CINEMAS AutomationStatusType = "inCinemas"
	AUTOSTATUSTYPE_RELEASED   AutomationStatusType = "released"
)

func GetStatusFromSeriesStatus(status sonarr.SeriesStatusType) AutomationStatusType {
	switch status {
	case sonarr.SERIESSTATUSTYPE_CONTINUING:
		return AUTOSTATUSTYPE_CONTINUING
	case sonarr.SERIESSTATUSTYPE_ENDED:
		return AUTOSTATUSTYPE_ENDED
	case sonarr.SERIESSTATUSTYPE_UPCOMING:
		return AUTOSTATUSTYPE_UPCOMING
	case sonarr.SERIESSTATUSTYPE_DELETED:
		return AUTOSTATUSTYPE_DELETED
	default:
		return "" // or some default value
	}
}

func GetStatusFromMovieStatus(status radarr.MovieStatusType) AutomationStatusType {
	switch status {
	case radarr.MOVIESTATUSTYPE_ANNOUNCED, radarr.MOVIESTATUSTYPE_TBA:
		return AUTOSTATUSTYPE_UPCOMING
	case radarr.MOVIESTATUSTYPE_IN_CINEMAS:
		return AUTOSTATUSTYPE_IN_CINEMAS
	case radarr.MOVIESTATUSTYPE_RELEASED:
		return AUTOSTATUSTYPE_RELEASED // or maybe ENDED is more appropriate here?
	case radarr.MOVIESTATUSTYPE_DELETED:
		return AUTOSTATUSTYPE_DELETED
	default:
		return "" // or some default value
	}
}

func GetStatusFromMusicStatus(status lidarr.ArtistStatusType) AutomationStatusType {
	switch status {
	case lidarr.ARTISTSTATUSTYPE_CONTINUING:
		return AUTOSTATUSTYPE_CONTINUING
	case lidarr.ARTISTSTATUSTYPE_ENDED:
		return AUTOSTATUSTYPE_ENDED
	default:
		return "" // or some default value
	}
}
