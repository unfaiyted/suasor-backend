package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	types "suasor/client/automation/types"
	client "suasor/client/types"
	"time"
)

// AutomationMediaItem represents a generic automation media item
type AutomationMediaItem[T types.AutomationData] struct {
	ID         uint64                      `json:"id" gorm:"primaryKey;autoIncrement"`
	ExternalID string                      `json:"externalId" gorm:"index;size:255"` // ID from external media client
	ClientID   uint64                      `json:"clientId" gorm:"index"`
	ClientType client.AutomationClientType `json:"clientType" gorm:"type:varchar(50)"`

	Title            string                 `json:"title" gorm:"type:varchar(255)"`
	Overview         string                 `json:"overview" gorm:"type:text"`
	Year             int32                  `json:"year"`
	AddedAt          time.Time              `json:"addedAt"`
	Ratings          []types.Rating         `json:"ratings" gorm:"type:jsonb"`
	DownloadedStatus types.DownloadedStatus `json:"downloadedStatus" gorm:"type:varchar(50)"`

	ExternalIDs    []types.ExternalID           `json:"externalIds" gorm:"type:jsonb"`
	Status         types.AutomationStatusType   `json:"status" gorm:"type:varchar(50)"`
	Path           string                       `json:"path" gorm:"type:varchar(1024)"`
	Genres         []string                     `json:"genres" gorm:"type:jsonb"`
	QualityProfile types.QualityProfileSummary  `json:"qualityProfile" gorm:"embedded"`
	Images         []types.AutomationMediaImage `json:"images" gorm:"type:jsonb"`
	Monitored      bool                         `json:"monitored"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Type types.AutomationMediaType `json:"mediaType" gorm:"type:varchar(50)"`
	Data T                         `json:"data" gorm:"type:jsonb"`
}

// Value implements driver.Valuer for database storage
func (m AutomationMediaItem[T]) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan implements sql.Scanner for database retrieval
func (m *AutomationMediaItem[T]) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &m)
}

func (m *AutomationMediaItem[T]) GetMediaType() types.AutomationMediaType {
	return m.Type
}

// MarshalJSON provides custom JSON serialization
func (m AutomationMediaItem[T]) MarshalJSON() ([]byte, error) {
	type Alias AutomationMediaItem[T]

	return json.Marshal(struct {
		Alias
		Data T `json:"data"`
	}{
		Alias: Alias(m),
		Data:  m.Data,
	})
}

// UnmarshalJSON provides custom JSON deserialization
func (m *AutomationMediaItem[T]) UnmarshalJSON(data []byte) error {
	type Alias AutomationMediaItem[T]
	aux := &struct {
		*Alias
		Data json.RawMessage `json:"data"`
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var mediaData T
	if err := json.Unmarshal(aux.Data, &mediaData); err != nil {
		return fmt.Errorf("error unmarshaling data field: %w", err)
	}
	m.Data = mediaData

	return nil
}

// SetClientInfo sets the client information
func (m *AutomationMediaItem[T]) SetClientInfo(clientID uint64, clientType client.AutomationClientType) {
	m.ClientID = clientID
	m.ClientType = clientType
}

// GetData returns the data field
func (m *AutomationMediaItem[T]) GetData() T {
	return m.Data
}

// SetData sets the data field
func (m *AutomationMediaItem[T]) SetData(data T, mediaType types.AutomationMediaType) {
	m.Data = data
	m.Type = mediaType
}
