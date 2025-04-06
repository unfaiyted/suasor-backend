package metadata

import (
	metadataTypes "suasor/client/metadata/types"
)

// Type aliases to make it easier to use the types from the client package
type (
	Movie        = metadataTypes.Movie
	TVShow       = metadataTypes.TVShow
	TVSeason     = metadataTypes.TVSeason
	TVEpisode    = metadataTypes.TVEpisode
	Person       = metadataTypes.Person
	Collection   = metadataTypes.Collection
	MediaImage   = metadataTypes.MediaImage
	Video        = metadataTypes.Video
	Genre        = metadataTypes.Genre
	ExternalIDs  = metadataTypes.ExternalIDs
	Credits      = metadataTypes.Credits
	CastMember   = metadataTypes.CastMember
	CrewMember   = metadataTypes.CrewMember
	MovieCredit  = metadataTypes.MovieCredit
	TVCredit     = metadataTypes.TVCredit
)