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
	Movies      map[string]*UserMediaItemData[*types.Movie]
	Series      map[string]*UserMediaItemData[*types.Series]
	Seasons     map[string]*UserMediaItemData[*types.Season]
	Episodes    map[string]*UserMediaItemData[*types.Episode]
	Artists     map[string]*UserMediaItemData[*types.Artist]
	Albums      map[string]*UserMediaItemData[*types.Album]
	Tracks      map[string]*UserMediaItemData[*types.Track]
	Playlists   map[string]*UserMediaItemData[*types.Playlist]
	Collections map[string]*UserMediaItemData[*types.Collection]

	OwnerID uint64

	Order DataListItems

	totalItems int
}

func NewMediaItemDataList() *MediaItemDataList {
	return &MediaItemDataList{
		Movies:      make(map[string]*UserMediaItemData[*types.Movie]),
		Series:      make(map[string]*UserMediaItemData[*types.Series]),
		Seasons:     make(map[string]*UserMediaItemData[*types.Season]),
		Episodes:    make(map[string]*UserMediaItemData[*types.Episode]),
		Artists:     make(map[string]*UserMediaItemData[*types.Artist]),
		Albums:      make(map[string]*UserMediaItemData[*types.Album]),
		Tracks:      make(map[string]*UserMediaItemData[*types.Track]),
		Playlists:   make(map[string]*UserMediaItemData[*types.Playlist]),
		Collections: make(map[string]*UserMediaItemData[*types.Collection]),

		Order: DataListItems{},

		totalItems: 0,
	}
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
	m.Movies[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}
func (m *MediaItemDataList) AddSeriesList(items []*UserMediaItemData[*types.Series]) {
	for _, item := range items {
		m.AddSeries(item)
	}
}

func (m *MediaItemDataList) AddSeries(item *UserMediaItemData[*types.Series]) {
	m.Series[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

func (m *MediaItemDataList) AddSeason(item *UserMediaItemData[*types.Season]) {
	m.Seasons[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

func (m *MediaItemDataList) AddSeasonList(items []*UserMediaItemData[*types.Season]) {
	for _, item := range items {
		m.AddSeason(item)
	}
}

func (m *MediaItemDataList) AddEpisode(item *UserMediaItemData[*types.Episode]) {
	m.Episodes[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

func (m *MediaItemDataList) AddEpisodeList(items []*UserMediaItemData[*types.Episode]) {
	for _, item := range items {
		m.AddEpisode(item)
	}
}

func (m *MediaItemDataList) AddArtist(item *UserMediaItemData[*types.Artist]) {
	m.Artists[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

func (m *MediaItemDataList) AddArtistList(items []*UserMediaItemData[*types.Artist]) {
	for _, item := range items {
		m.AddArtist(item)
	}
}

func (m *MediaItemDataList) AddAlbum(item *UserMediaItemData[*types.Album]) {
	m.Albums[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

func (m *MediaItemDataList) AddAlbumList(items []*UserMediaItemData[*types.Album]) {
	for _, item := range items {
		m.AddAlbum(item)
	}
}

func (m *MediaItemDataList) AddTrack(item *UserMediaItemData[*types.Track]) {
	m.Tracks[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

func (m *MediaItemDataList) AddTrackList(items []*UserMediaItemData[*types.Track]) {
	for _, item := range items {
		m.AddTrack(item)
	}
}

func (m *MediaItemDataList) AddPlaylist(item *UserMediaItemData[*types.Playlist]) {
	m.Playlists[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
	m.totalItems++
}

func (m *MediaItemDataList) AddPlaylistList(items []*UserMediaItemData[*types.Playlist]) {
	for _, item := range items {
		m.AddPlaylist(item)
	}
}

func (m *MediaItemDataList) AddCollection(item *UserMediaItemData[*types.Collection]) {
	m.Collections[item.UUID] = item
	m.AddListItem(item.UUID, m.totalItems+1)
}

func (m *MediaItemDataList) AddCollectionList(items []*UserMediaItemData[*types.Collection]) {
	for _, item := range items {
		m.AddCollection(item)
	}
}

func (m *MediaItemDataList) GetTotalItems() int {
	return m.totalItems
}

func (m *MediaItemDataList) GetMoviesArray() []*UserMediaItemData[*types.Movie] {
	var movies []*UserMediaItemData[*types.Movie]
	for _, item := range m.Movies {
		movies = append(movies, item)
	}
	return movies
}

func (m *MediaItemDataList) GetSeriesArray() []*UserMediaItemData[*types.Series] {
	var series []*UserMediaItemData[*types.Series]
	for _, item := range m.Series {
		series = append(series, item)
	}
	return series
}

func (m *MediaItemDataList) GetEpisodesArray() []*UserMediaItemData[*types.Episode] {
	var episodes []*UserMediaItemData[*types.Episode]
	for _, item := range m.Episodes {
		episodes = append(episodes, item)
	}
	return episodes
}

func (m *MediaItemDataList) GetTracksArray() []*UserMediaItemData[*types.Track] {
	var tracks []*UserMediaItemData[*types.Track]
	for _, item := range m.Tracks {
		tracks = append(tracks, item)
	}
	return tracks
}

func (m *MediaItemDataList) GetAlbumsArray() []*UserMediaItemData[*types.Album] {
	var albums []*UserMediaItemData[*types.Album]
	for _, item := range m.Albums {
		albums = append(albums, item)
	}
	return albums
}
