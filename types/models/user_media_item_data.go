package models

import (
	"errors"
	"gorm.io/gorm"
	"suasor/client/media/types"
	"time"
)

type MediaItemDatas struct {
	Movies      []*UserMediaItemData[*types.Movie]
	Series      []*UserMediaItemData[*types.Series]
	Seasons     []*UserMediaItemData[*types.Season]
	Episodes    []*UserMediaItemData[*types.Episode]
	Artists     []*UserMediaItemData[*types.Artist]
	Albums      []*UserMediaItemData[*types.Album]
	Tracks      []*UserMediaItemData[*types.Track]
	Playlists   []*UserMediaItemData[*types.Playlist]
	Collections []*UserMediaItemData[*types.Collection]

	TotalItems int
}

func (m *MediaItemDatas) AddMovie(item *UserMediaItemData[*types.Movie]) {
	m.Movies = append(m.Movies, item)
	m.TotalItems++
}
func (m *MediaItemDatas) AddSeries(item *UserMediaItemData[*types.Series]) {
	m.Series = append(m.Series, item)
	m.TotalItems++
}
func (m *MediaItemDatas) AddSeason(item *UserMediaItemData[*types.Season]) {
	m.Seasons = append(m.Seasons, item)
	m.TotalItems++
}
func (m *MediaItemDatas) AddEpisode(item *UserMediaItemData[*types.Episode]) {
	m.Episodes = append(m.Episodes, item)
	m.TotalItems++
}
func (m *MediaItemDatas) AddArtist(item *UserMediaItemData[*types.Artist]) {
	m.Artists = append(m.Artists, item)
	m.TotalItems++
}
func (m *MediaItemDatas) AddAlbum(item *UserMediaItemData[*types.Album]) {
	m.Albums = append(m.Albums, item)
	m.TotalItems++
}
func (m *MediaItemDatas) AddTrack(item *UserMediaItemData[*types.Track]) {
	m.Tracks = append(m.Tracks, item)
	m.TotalItems++
}
func (m *MediaItemDatas) AddPlaylist(item *UserMediaItemData[*types.Playlist]) {
	m.Playlists = append(m.Playlists, item)
	m.TotalItems++
}
func (m *MediaItemDatas) AddCollection(item *UserMediaItemData[*types.Collection]) {
	m.Collections = append(m.Collections, item)
	m.TotalItems++
}

func (m *MediaItemDatas) GetTotalItems() int {
	return m.TotalItems
}
func (m *MediaItemDatas) GetMovies() []*UserMediaItemData[*types.Movie] {
	return m.Movies
}
func (m *MediaItemDatas) GetSeries() []*UserMediaItemData[*types.Series] {
	return m.Series
}
func (m *MediaItemDatas) GetSeasons() []*UserMediaItemData[*types.Season] {
	return m.Seasons
}
func (m *MediaItemDatas) GetEpisodes() []*UserMediaItemData[*types.Episode] {
	return m.Episodes
}
func (m *MediaItemDatas) GetArtists() []*UserMediaItemData[*types.Artist] {
	return m.Artists
}
func (m *MediaItemDatas) GetAlbums() []*UserMediaItemData[*types.Album] {
	return m.Albums
}
func (m *MediaItemDatas) GetTracks() []*UserMediaItemData[*types.Track] {
	return m.Tracks
}
func (m *MediaItemDatas) GetPlaylists() []*UserMediaItemData[*types.Playlist] {
	return m.Playlists
}
func (m *MediaItemDatas) GetCollections() []*UserMediaItemData[*types.Collection] {
	return m.Collections
}

// Represents the user's personal data for a specific media item
type UserMediaItemData[T types.MediaData] struct {
	ID               uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           uint64          `json:"userId" gorm:"index"`          // Foreign key to User
	MediaItemID      uint64          `json:"mediaItemId" gorm:"index"`     // Foreign key to MediaItem
	Item             *MediaItem[T]   `json:"item" gorm:"-"`                // Not stored in DB, loaded via relationship
	Type             types.MediaType `json:"type" gorm:"type:varchar(50)"` // "movie", "episode", "show", "season"
	PlayedAt         time.Time       `json:"playedAt" gorm:"index"`
	LastPlayedAt     time.Time       `json:"lastPlayedAt" gorm:"index"`
	IsFavorite       bool            `json:"isFavorite,omitempty"`
	IsDisliked       bool            `json:"isDisliked,omitempty"`
	UserRating       float32         `json:"userRating,omitempty"`
	Watchlist        bool            `json:"watchlist,omitempty"`
	PlayedPercentage float64         `json:"playedPercentage,omitempty"`
	PlayCount        int32           `json:"playCount,omitempty"`
	PositionSeconds  int             `json:"positionSeconds"`
	DurationSeconds  int             `json:"durationSeconds"`
	Completed        bool            `json:"completed"`
	CreatedAt        time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}

// Associate links this history record with a media item
func (h *UserMediaItemData[T]) Associate(item *MediaItem[T]) {
	h.MediaItemID = item.ID
	h.Item = item
	h.Type = item.Type
}

// LoadItem loads the associated MediaItem from the database using MediaItemID
// This would be called after retrieving the history record from the database
func (h *UserMediaItemData[T]) LoadItem(db *gorm.DB) error {
	if h.MediaItemID == 0 {
		return errors.New("no media item ID associated with this history record")
	}

	item := &MediaItem[T]{}
	result := db.First(item, h.MediaItemID)
	if result.Error != nil {
		return result.Error
	}

	h.Item = item
	return nil
}

// BeforeSave ensures we have the proper MediaItemID before saving
func (h *UserMediaItemData[T]) BeforeSave(tx *gorm.DB) error {
	if h.MediaItemID == 0 && h.Item != nil {
		h.MediaItemID = h.Item.ID
	}
	return nil
}

// UserMediaItemDataGeneric is a non-generic version of MediaPlayHistory to avoid type issues
// type UserMediaItemDataGeneric struct {
// 	ID               uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
// 	UserID           uint64          `json:"userId" gorm:"index"`
// 	MediaItemID      uint64          `json:"mediaItemId" gorm:"index"`
// 	Type             types.MediaType `json:"type" gorm:"type:varchar(50)"`
// 	PlayedAt         time.Time       `json:"playedAt" gorm:"index"`
// 	LastPlayedAt     time.Time       `json:"lastPlayedAt" gorm:"index"`
// 	IsFavorite       bool            `json:"isFavorite,omitempty"`
// 	IsDisliked       bool            `json:"isDisliked,omitempty"`
// 	UserRating       float32         `json:"userRating,omitempty"`
// 	PlayedPercentage float64         `json:"playedPercentage,omitempty"`
// 	PlayCount        int32           `json:"playCount,omitempty"`
// 	PositionSeconds  int             `json:"positionSeconds"`
// 	DurationSeconds  int             `json:"durationSeconds"`
// 	Completed        bool            `json:"completed"`
// 	CreatedAt        time.Time       `json:"createdAt" gorm:"autoCreateTime"`
// 	UpdatedAt        time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
// }

// UserMediaItemDataRequest is used to record a new play history entry
// type UserMediaItemDataRequest struct {
// 	UserID           uint64          `json:"userId" binding:"required"`
// 	MediaItemID      uint64          `json:"mediaItemId" binding:"required"`
// 	Type             types.MediaType `json:"type" binding:"required"`
// 	IsFavorite       bool            `json:"isFavorite,omitempty"`
// 	UserRating       float32         `json:"userRating,omitempty"`
// 	PlayedPercentage float64         `json:"playedPercentage,omitempty"`
// 	PositionSeconds  int             `json:"positionSeconds"`
// 	DurationSeconds  int             `json:"durationSeconds"`
// 	Completed        bool            `json:"completed"`
// 	Continued        bool            `json:"continued"` // If this is a continuation of a previous play
// }
