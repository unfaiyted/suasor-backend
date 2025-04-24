package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	types "suasor/clients/automation/types"
	client "suasor/clients/types"
	"time"
)

// AutomationCalendarItem represents an upcoming media item
type AutomationCalendarItem[T types.AutomationData] struct {
	ID         uint64                      `json:"ID" gorm:"primaryKey;autoIncrement"`
	ExternalID string                      `json:"externalID" gorm:"index;size:255"` // ID from external media client
	ClientID   uint32                      `json:"clientID" gorm:"index"`
	ClientType client.AutomationClientType `json:"clientType" gorm:"type:varchar(50)"`

	// Foreign key to the media item
	ItemID uint64                    `json:"itemID" gorm:"index"`
	Item   *T                        `json:"item" gorm:"foreignKey:ItemID"`
	Type   types.AutomationMediaType `json:"type" gorm:"type:varchar(50)"`

	AirDate   time.Time `json:"airDate" gorm:"index"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// Value implements driver.Valuer for database storage
func (c AutomationCalendarItem[T]) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (c *AutomationCalendarItem[T]) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &c)
}

// MarshalJSON provides custom JSON serialization
func (c AutomationCalendarItem[T]) MarshalJSON() ([]byte, error) {
	type Alias AutomationCalendarItem[T]
	return json.Marshal(struct {
		Alias
	}{
		Alias: Alias(c),
	})
}

// UnmarshalJSON provides custom JSON deserialization
func (c *AutomationCalendarItem[T]) UnmarshalJSON(data []byte) error {
	type Alias AutomationCalendarItem[T]
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

// SetClientInfo sets the client information
func (c *AutomationCalendarItem[T]) SetClientInfo(clientID uint32, clientType client.AutomationClientType, externalID string) {
	c.ClientID = clientID
	c.ClientType = clientType
	c.ExternalID = externalID
}

func (c *AutomationCalendarItem[T]) GetMediaItem() *T {
	return c.Item
}

// GetMediaItem returns the associated media item
func (c *AutomationCalendarItem[T]) GetTypedMediaItem() interface{} {
	switch item := any(c.Item).(type) {
	case *types.AutomationMovie:
		return item
	case *types.AutomationTVShow:
		return item
	case *types.AutomationEpisode:
		return item
	case *types.AutomationArtist:
		return item
	case *types.AutomationAlbum:
		return item
	case *types.AutomationTrack:
		return item
	default:
		return nil
	}
}

// SetMediaItem sets the associated media item
func (c *AutomationCalendarItem[T]) SetMediaItem(item *T) {
	c.Item = item
	if item != nil {
		if mediaItem, ok := any(item).(interface{ GetID() uint64 }); ok {
			c.ItemID = mediaItem.GetID()
		}
	}
}

func CreateCalendarItem(mediaType types.AutomationMediaType) (any, error) {
	switch mediaType {
	case types.AUTOMEDIATYPE_MOVIE:
		return &AutomationCalendarItem[types.AutomationMovie]{Type: types.AUTOMEDIATYPE_MOVIE}, nil
	case types.AUTOMEDIATYPE_SERIES:
		return &AutomationCalendarItem[types.AutomationTVShow]{Type: types.AUTOMEDIATYPE_SERIES}, nil
	case types.AUTOMEDIATYPE_ARTIST:
		return &AutomationCalendarItem[types.AutomationArtist]{Type: types.AUTOMEDIATYPE_ARTIST}, nil
	default:
		return nil, fmt.Errorf("unknown media type: %s", mediaType)
	}
}
