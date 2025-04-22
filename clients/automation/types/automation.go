package types

import (
	"time"
)

// AutomationData defines the allowed types for AutomationMediaItem's Data field
type AutomationData interface {
	isAutomationData()
	GetMediaType() AutomationMediaType
}

// Implement the marker method for each allowed type
func (AutomationMovie) isAutomationData()   {}
func (AutomationTVShow) isAutomationData()  {}
func (AutomationEpisode) isAutomationData() {}
func (AutomationArtist) isAutomationData()  {}
func (AutomationAlbum) isAutomationData()   {}
func (AutomationTrack) isAutomationData()   {}

func (AutomationMovie) GetMediaType() AutomationMediaType {
	return AUTOMEDIATYPE_MOVIE
}
func (AutomationTVShow) GetMediaType() AutomationMediaType {
	return AUTOMEDIATYPE_SERIES
}
func (AutomationArtist) GetMediaType() AutomationMediaType {
	return AUTOMEDIATYPE_ARTIST
}
func (AutomationAlbum) GetMediaType() AutomationMediaType {
	return AUTOMEDIATYPE_ALBUM
}

func (AutomationEpisode) GetMediaType() AutomationMediaType {
	return AUTOMEDIATYPE_EPISODE
}

type AutomationMovie struct {
	Year        int32
	ReleaseDate time.Time
}

type AutomationTVShow struct {
	ReleaseDate time.Time
	Year        int32
}

type AutomationEpisode struct {
	ReleaseDate time.Time
}

type AutomationTrack struct {
	AlbumName  string
	ArtistName string
}

type AutomationArtist struct {
	Albums                []AutomationAlbum
	MetadataProfile       MetadataProfile
	MostRecentReleaseDate time.Time
}

type AutomationAlbum struct {
	ArtistName  string
	ArtistID    string
	ReleaseDate time.Time
}
