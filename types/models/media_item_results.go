package models

import (
	"suasor/clients/media/types"
)

type MediaItemResults struct {
	Movies      []*MediaItem[*types.Movie]      `json:"movies,omitempty"`
	Series      []*MediaItem[*types.Series]     `json:"series,omitempty"`
	Episodes    []*MediaItem[*types.Episode]    `json:"episodes,omitempty"`
	Seasons     []*MediaItem[*types.Season]     `json:"seasons,omitempty"`
	Tracks      []*MediaItem[*types.Track]      `json:"tracks,omitempty"`
	Albums      []*MediaItem[*types.Album]      `json:"albums,omitempty"`
	Artists     []*MediaItem[*types.Artist]     `json:"artists,omitempty"`
	Playlists   []*MediaItem[*types.Playlist]   `json:"playlists,omitempty"`
	Collections []*MediaItem[*types.Collection] `json:"collections,omitempty"`

	TotalItems int `json:"totalItems"`
}

func NewMediaItemResults() *MediaItemResults {
	return &MediaItemResults{
		TotalItems: 0,
	}
}

func (m *MediaItemResults) Len() int {
	return m.TotalItems
}

func (m *MediaItemResults) GetSyncClientItemIDs(clientID uint64) []string {
	ids := make([]string, 0)
	m.ForEach(func(uuid string, mediaType types.MediaType, item any) bool {
		if mediaItem, ok := item.(*MediaItem[types.MediaData]); ok && mediaItem != nil {
			clientItemID := mediaItem.SyncClients.GetClientItemID(clientID)
			// Only add non-empty client IDs
			if clientItemID != "" {
				ids = append(ids, clientItemID)
			}
		}
		return true
	})

	return ids
}

// AddMovie adds a movie to the media items
func (m *MediaItemResults) AddMovie(item *MediaItem[*types.Movie]) {
	m.Movies = append(m.Movies, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddMovieList(items []*MediaItem[*types.Movie]) {
	for _, item := range items {
		m.AddMovie(item)
	}
}

// AddSeries adds a series to the media items
func (m *MediaItemResults) AddSeries(item *MediaItem[*types.Series]) {
	m.Series = append(m.Series, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddSeriesList(items []*MediaItem[*types.Series]) {
	for _, item := range items {
		m.AddSeries(item)
	}
}

// AddSeason adds a season to the media items
func (m *MediaItemResults) AddSeason(item *MediaItem[*types.Season]) {
	m.Seasons = append(m.Seasons, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddSeasonList(items []*MediaItem[*types.Season]) {
	for _, item := range items {
		m.AddSeason(item)
	}
}

// AddEpisode adds an episode to the media items
func (m *MediaItemResults) AddEpisode(item *MediaItem[*types.Episode]) {
	m.Episodes = append(m.Episodes, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddEpisodeList(items []*MediaItem[*types.Episode]) {
	for _, item := range items {
		m.AddEpisode(item)
	}
}

// AddArtist adds an artist to the media items
func (m *MediaItemResults) AddArtist(item *MediaItem[*types.Artist]) {
	m.Artists = append(m.Artists, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddArtistList(items []*MediaItem[*types.Artist]) {
	for _, item := range items {
		m.AddArtist(item)
	}
}

// AddAlbum adds an album to the media items
func (m *MediaItemResults) AddAlbum(item *MediaItem[*types.Album]) {
	m.Albums = append(m.Albums, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddAlbumList(items []*MediaItem[*types.Album]) {
	for _, item := range items {
		m.AddAlbum(item)
	}
}

// AddTrack adds a track to the media items
func (m *MediaItemResults) AddTrack(item *MediaItem[*types.Track]) {
	m.Tracks = append(m.Tracks, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddTrackList(items []*MediaItem[*types.Track]) {
	for _, item := range items {
		m.AddTrack(item)
	}
}

// AddPlaylist adds a playlist to the media items
func (m *MediaItemResults) AddPlaylist(item *MediaItem[*types.Playlist]) {
	m.Playlists = append(m.Playlists, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddPlaylistList(items []*MediaItem[*types.Playlist]) {
	for _, item := range items {
		m.AddPlaylist(item)
	}
}

// AddCollection adds a collection to the media items
func (m *MediaItemResults) AddCollection(item *MediaItem[*types.Collection]) {
	m.Collections = append(m.Collections, item)
	m.TotalItems++
}

func (m *MediaItemResults) AddCollectionList(items []*MediaItem[*types.Collection]) {
	for _, item := range items {
		m.AddCollection(item)
	}
}

// ForEach iterates over all media items in the list in the specified order.
// The callback function receives the UUID, media type, and the item itself.
// If the callback returns false, iteration stops early.
func (m *MediaItemResults) ForEach(callback func(uuid string, mediaType types.MediaType, item any) bool) {
	// go through all moves, series, episodes, seasons, tracks, albums, artists, playlists, collections
	// and call the callback with the appropriate type
	for _, item := range m.Movies {
		if !callback(item.UUID, types.MediaTypeMovie, item) {
			return
		}
	}
	for _, item := range m.Series {
		if !callback(item.UUID, types.MediaTypeSeries, item) {
			return
		}
	}
	for _, item := range m.Episodes {
		if !callback(item.UUID, types.MediaTypeEpisode, item) {
			return
		}
	}
	for _, item := range m.Seasons {
		if !callback(item.UUID, types.MediaTypeSeason, item) {
			return
		}
	}
	for _, item := range m.Tracks {
		if !callback(item.UUID, types.MediaTypeTrack, item) {
			return
		}
	}
	for _, item := range m.Albums {
		if !callback(item.UUID, types.MediaTypeAlbum, item) {
			return
		}
	}
	for _, item := range m.Artists {
		if !callback(item.UUID, types.MediaTypeArtist, item) {
			return
		}
	}
	for _, item := range m.Playlists {
		if !callback(item.UUID, types.MediaTypePlaylist, item) {
			return
		}
	}
	for _, item := range m.Collections {
		if !callback(item.UUID, types.MediaTypeCollection, item) {
			return
		}
	}
}
