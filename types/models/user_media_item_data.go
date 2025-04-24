package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"suasor/clients/media/types"
	"time"
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

// Represents the user's personal data for a specific media item
type UserMediaItemData[T types.MediaData] struct {
	ID               uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID             string          `json:"uuid" gorm:"type:uuid;uniqueIndex"` // Stable UUID for syncing
	UserID           uint64          `json:"userId" gorm:"index"`               // Foreign key to User
	MediaItemID      uint64          `json:"mediaItemId" gorm:"index"`          // Foreign key to MediaItem
	Item             *MediaItem[T]   `json:"item" gorm:"-"`                     // Not stored in DB, loaded via relationship
	Type             types.MediaType `json:"type" gorm:"type:varchar(50)"`      // "movie", "episode", "show", "season"
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

func NewUserMediaItemData[T types.MediaData](item *MediaItem[T], userID uint64) *UserMediaItemData[T] {
	// Create a new user media item data object with the media item
	result := &UserMediaItemData[T]{
		ID:          1, // Placeholder ID
		UUID:        uuid.New().String(),
		UserID:      userID,
		MediaItemID: item.ID,
		Item:        item,
		Type:        item.Type,
		// Default values for other fields
		IsFavorite:       false,
		UserRating:       0,
		PlayedPercentage: 0,
		Watchlist:        false,
		Completed:        false,
	}

	return result
}
