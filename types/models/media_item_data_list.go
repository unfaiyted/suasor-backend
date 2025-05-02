package models

import (
	"suasor/clients/media/types"
)

type DataListItem struct {
	ItemUUID string `json:"itemUUID"`
	Position int    `json:"position"`
}

type DataListItems []DataListItem

type MediaItemDataList struct {
	movies      map[string]*UserMediaItemData[*types.Movie]
	series      map[string]*UserMediaItemData[*types.Series]
	seasons     map[string]*UserMediaItemData[*types.Season]
	episodes    map[string]*UserMediaItemData[*types.Episode]
	artists     map[string]*UserMediaItemData[*types.Artist]
	albums      map[string]*UserMediaItemData[*types.Album]
	tracks      map[string]*UserMediaItemData[*types.Track]
	playlists   map[string]*UserMediaItemData[*types.Playlist]
	collections map[string]*UserMediaItemData[*types.Collection]

	OwnerID uint64

	Order DataListItems

	TotalItems int
}

func (m *MediaItemDataList) AddListItem(itemUUID string, itemPosition int) {
	m.Order = append(m.Order, DataListItem{
		ItemUUID: itemUUID,
		Position: itemPosition,
	})
}

func (m *MediaItemDataList) AddMovieList(items []*UserMediaItemData[*types.Movie]) {
	for _, item := range items {
		m.AddMovie(item)
	}
}
func (m *MediaItemDataList) AddMovie(item *UserMediaItemData[*types.Movie]) {
	m.movies[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}
func (m *MediaItemDataList) AddSeriesList(items []*UserMediaItemData[*types.Series]) {
	for _, item := range items {
		m.AddSeries(item)
	}
}

func (m *MediaItemDataList) AddSeries(item *UserMediaItemData[*types.Series]) {
	m.series[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemDataList) AddSeason(item *UserMediaItemData[*types.Season]) {
	m.seasons[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemDataList) AddSeasonList(items []*UserMediaItemData[*types.Season]) {
	for _, item := range items {
		m.AddSeason(item)
	}
}

func (m *MediaItemDataList) AddEpisode(item *UserMediaItemData[*types.Episode]) {
	m.episodes[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemDataList) AddEpisodeList(items []*UserMediaItemData[*types.Episode]) {
	for _, item := range items {
		m.AddEpisode(item)
	}
}

func (m *MediaItemDataList) AddArtist(item *UserMediaItemData[*types.Artist]) {
	m.artists[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemDataList) AddArtistList(items []*UserMediaItemData[*types.Artist]) {
	for _, item := range items {
		m.AddArtist(item)
	}
}

func (m *MediaItemDataList) AddAlbum(item *UserMediaItemData[*types.Album]) {
	m.albums[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemDataList) AddAlbumList(items []*UserMediaItemData[*types.Album]) {
	for _, item := range items {
		m.AddAlbum(item)
	}
}

func (m *MediaItemDataList) AddTrack(item *UserMediaItemData[*types.Track]) {
	m.tracks[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemDataList) AddTrackList(items []*UserMediaItemData[*types.Track]) {
	for _, item := range items {
		m.AddTrack(item)
	}
}

func (m *MediaItemDataList) AddPlaylist(item *UserMediaItemData[*types.Playlist]) {
	m.playlists[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
	m.TotalItems++
}

func (m *MediaItemDataList) AddPlaylistList(items []*UserMediaItemData[*types.Playlist]) {
	for _, item := range items {
		m.AddPlaylist(item)
	}
}

func (m *MediaItemDataList) AddCollection(item *UserMediaItemData[*types.Collection]) {
	m.collections[item.UUID] = item
	m.AddListItem(item.UUID, m.TotalItems+1)
}

func (m *MediaItemDataList) AddCollectionList(items []*UserMediaItemData[*types.Collection]) {
	for _, item := range items {
		m.AddCollection(item)
	}
}

func (m *MediaItemDataList) GetTotalItems() int {
	return m.TotalItems
}
