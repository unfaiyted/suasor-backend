package models

import (
	"suasor/clients/media/types"
	"time"
)

type MediaItemList[T types.ListData] struct {
	List *MediaItem[T] `json:"list"`

	Movies      map[string]*MediaItem[*types.Movie]      `json:"movies,omitempty"`
	Series      map[string]*MediaItem[*types.Series]     `json:"series,omitempty"`
	Episodes    map[string]*MediaItem[*types.Episode]    `json:"episodes,omitempty"`
	Seasons     map[string]*MediaItem[*types.Season]     `json:"seasons,omitempty"`
	Tracks      map[string]*MediaItem[*types.Track]      `json:"tracks,omitempty"`
	Albums      map[string]*MediaItem[*types.Album]      `json:"albums,omitempty"`
	Artists     map[string]*MediaItem[*types.Artist]     `json:"artists,omitempty"`
	Playlists   map[string]*MediaItem[*types.Playlist]   `json:"playlists,omitempty"`
	Collections map[string]*MediaItem[*types.Collection] `json:"collections,omitempty"`

	ListType     types.ListType `json:"listType"`
	ListOriginID uint64         `json:"listOriginID"` // 0 for internal db, otherwise external client/ProviderID
	OwnerID      uint64         `json:"ownerID"`

	Order ListItems `json:"order"`

	TotalItems int `json:"totalItems"`
}

func NewMediaItemList[T types.ListData](list *MediaItem[T], listOriginID uint64, ownerID uint64) *MediaItemList[T] {

	listType := types.GetListType[T]()

	return &MediaItemList[T]{
		List:         list,
		ListType:     listType,
		ListOriginID: listOriginID,
		OwnerID:      ownerID,
		Order:        ListItems{},
		TotalItems:   0,
	}
}

func (m *MediaItemList[T]) AddListItem(itemUUID string, itemPosition int) {
	m.Order = append(m.Order, ListItem{
		ItemUUID:    itemUUID,
		Position:    itemPosition,
		LastChanged: time.Now(),
	})
}

func (m *MediaItemList[T]) GetTotalItems() int {
	return m.TotalItems
}

func (m *MediaItemList[T]) GetSyncClientItemIDs(clientID uint64) []string {
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
func (m *MediaItemList[T]) AddMovie(item *MediaItem[*types.Movie]) {
	m.Movies[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddMovieList(items []*MediaItem[*types.Movie]) {
	for _, item := range items {
		m.AddMovie(item)
	}
}

// AddSeries adds a series to the media items
func (m *MediaItemList[T]) AddSeries(item *MediaItem[*types.Series]) {
	m.Series[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddSeriesList(items []*MediaItem[*types.Series]) {
	for _, item := range items {
		m.AddSeries(item)
	}
}

// AddSeason adds a season to the media items
func (m *MediaItemList[T]) AddSeason(item *MediaItem[*types.Season]) {
	m.Seasons[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddSeasonList(items []*MediaItem[*types.Season]) {
	for _, item := range items {
		m.AddSeason(item)
	}
}

// AddEpisode adds an episode to the media items
func (m *MediaItemList[T]) AddEpisode(item *MediaItem[*types.Episode]) {
	m.Episodes[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddEpisodeList(items []*MediaItem[*types.Episode]) {
	for _, item := range items {
		m.AddEpisode(item)
	}
}

// AddArtist adds an artist to the media items
func (m *MediaItemList[T]) AddArtist(item *MediaItem[*types.Artist]) {
	m.Artists[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddArtistList(items []*MediaItem[*types.Artist]) {
	for _, item := range items {
		m.AddArtist(item)
	}
}

// AddAlbum adds an album to the media items
func (m *MediaItemList[T]) AddAlbum(item *MediaItem[*types.Album]) {
	m.Albums[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddAlbumList(items []*MediaItem[*types.Album]) {
	for _, item := range items {
		m.AddAlbum(item)
	}
}

// AddTrack adds a track to the media items
func (m *MediaItemList[T]) AddTrack(item *MediaItem[*types.Track]) {
	m.Tracks[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddTrackList(items []*MediaItem[*types.Track]) {
	for _, item := range items {
		m.AddTrack(item)
	}
}

// AddPlaylist adds a playlist to the media items
func (m *MediaItemList[T]) AddPlaylist(item *MediaItem[*types.Playlist]) {
	m.Playlists[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddPlaylistList(items []*MediaItem[*types.Playlist]) {
	for _, item := range items {
		m.AddPlaylist(item)
	}
}

// AddCollection adds a collection to the media items
func (m *MediaItemList[T]) AddCollection(item *MediaItem[*types.Collection]) {
	m.Collections[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemList[T]) AddCollectionList(items []*MediaItem[*types.Collection]) {
	for _, item := range items {
		m.AddCollection(item)
	}
}

// ForEach iterates over all media items in the list in the specified order.
// The callback function receives the UUID, media type, and the item itself.
// If the callback returns false, iteration stops early.
func (m *MediaItemList[T]) ForEach(callback func(uuid string, mediaType types.MediaType, item any) bool) {
	for _, listItem := range m.Order {
		uuid := listItem.ItemUUID

		// Check in each map and call the callback with the appropriate type
		if movie, ok := m.Movies[uuid]; ok {
			if !callback(uuid, types.MediaTypeMovie, movie) {
				return
			}
			continue
		}
		if series, ok := m.Series[uuid]; ok {
			if !callback(uuid, types.MediaTypeSeries, series) {
				return
			}
			continue
		}
		if episode, ok := m.Episodes[uuid]; ok {
			if !callback(uuid, types.MediaTypeEpisode, episode) {
				return
			}
			continue
		}
		if season, ok := m.Seasons[uuid]; ok {
			if !callback(uuid, types.MediaTypeSeason, season) {
				return
			}
			continue
		}
		if track, ok := m.Tracks[uuid]; ok {
			if !callback(uuid, types.MediaTypeTrack, track) {
				return
			}
			continue
		}
		if album, ok := m.Albums[uuid]; ok {
			if !callback(uuid, types.MediaTypeAlbum, album) {
				return
			}
			continue
		}
		if artist, ok := m.Artists[uuid]; ok {
			if !callback(uuid, types.MediaTypeArtist, artist) {
				return
			}
			continue
		}
		if playlist, ok := m.Playlists[uuid]; ok {
			if !callback(uuid, types.MediaTypePlaylist, playlist) {
				return
			}
			continue
		}
		if collection, ok := m.Collections[uuid]; ok {
			if !callback(uuid, types.MediaTypeCollection, collection) {
				return
			}
			continue
		}
	}
}

// IsItemAtPosition checks if a media item is at a specific position
func (m *MediaItemList[T]) IsItemAtPosition(uuid string, position int) bool {
	for _, item := range m.Order {
		if item.ItemUUID == uuid && item.Position == position {
			return true
		}
	}
	return false
}
