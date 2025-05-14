package models

import (
	"suasor/clients/media/types"
	"time"
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

	Order ListItems `json:"order"`

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

func (m *MediaItemResults) AddOrderedItem(itemUUID string, itemPosition int) {
	m.Order = append(m.Order, ListItem{
		ItemUUID:    itemUUID,
		Position:    itemPosition,
		LastChanged: time.Now(),
	})
}

// AddMovie adds a movie to the media items
func (m *MediaItemResults) AddMovie(item *MediaItem[*types.Movie]) {
	m.Movies = append(m.Movies, item)
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
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
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
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
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
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
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
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
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemResults) AddArtistList(items []*MediaItem[*types.Artist]) {
	for _, item := range items {
		m.AddOrderedItem(item.UUID, m.TotalItems+1)
		m.AddArtist(item)
	}
}

// AddAlbum adds an album to the media items
func (m *MediaItemResults) AddAlbum(item *MediaItem[*types.Album]) {
	m.Albums = append(m.Albums, item)
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
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
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
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
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemResults) AddPlaylistList(items []*MediaItem[*types.Playlist]) {
	for _, item := range items {
		m.AddOrderedItem(item.UUID, m.TotalItems+1)
		m.AddPlaylist(item)
	}
}

// AddCollection adds a collection to the media items
func (m *MediaItemResults) AddCollection(item *MediaItem[*types.Collection]) {
	m.Collections = append(m.Collections, item)
	m.AddOrderedItem(item.UUID, m.TotalItems+1)
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
func (m *MediaItemResults) ForEachByType(callback func(uuid string, mediaType types.MediaType, item any) bool) {
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

func (m *MediaItemResults) ForEach(callback func(uuid string, mediaType types.MediaType, item any) bool) {

	// convert to types to map for faster lookup
	Movies := convertToMap(m.Movies)
	Series := convertToMap(m.Series)
	Episodes := convertToMap(m.Episodes)
	Seasons := convertToMap(m.Seasons)
	Tracks := convertToMap(m.Tracks)
	Albums := convertToMap(m.Albums)
	Artists := convertToMap(m.Artists)
	Playlists := convertToMap(m.Playlists)
	Collections := convertToMap(m.Collections)

	for _, listItem := range m.Order {
		uuid := listItem.ItemUUID

		// Check in each map and call the callback with the appropriate type
		if movie, ok := Movies[uuid]; ok {
			if !callback(uuid, types.MediaTypeMovie, movie) {
				return
			}
			continue
		}
		if series, ok := Series[uuid]; ok {
			if !callback(uuid, types.MediaTypeSeries, series) {
				return
			}
			continue
		}
		if episode, ok := Episodes[uuid]; ok {
			if !callback(uuid, types.MediaTypeEpisode, episode) {
				return
			}
			continue
		}
		if season, ok := Seasons[uuid]; ok {
			if !callback(uuid, types.MediaTypeSeason, season) {
				return
			}
			continue
		}
		if track, ok := Tracks[uuid]; ok {
			if !callback(uuid, types.MediaTypeTrack, track) {
				return
			}
			continue
		}
		if album, ok := Albums[uuid]; ok {
			if !callback(uuid, types.MediaTypeAlbum, album) {
				return
			}
			continue
		}
		if artist, ok := Artists[uuid]; ok {
			if !callback(uuid, types.MediaTypeArtist, artist) {
				return
			}
			continue
		}
		if playlist, ok := Playlists[uuid]; ok {
			if !callback(uuid, types.MediaTypePlaylist, playlist) {
				return
			}
			continue
		}
		if collection, ok := Collections[uuid]; ok {
			if !callback(uuid, types.MediaTypeCollection, collection) {
				return
			}
			continue
		}
	}
}

// IsItemAtPosition checks if a media item is at a specific position
func (m *MediaItemResults) IsItemAtPosition(uuid string, position int) bool {
	for _, item := range m.Order {
		if item.ItemUUID == uuid && item.Position == position {
			return true
		}
	}
	return false
}

func convertToMap[T types.MediaData](items []*MediaItem[T]) map[string]*MediaItem[T] {
	result := make(map[string]*MediaItem[T])
	for _, item := range items {
		result[item.UUID] = item
	}
	return result
}
