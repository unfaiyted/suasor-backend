package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"suasor/clients/types"
	"time"
)

type SyncStatus string

const (
	SyncStatusSuccess SyncStatus = "success"
	SyncStatusFailed  SyncStatus = "failed"
	SyncStatusPending SyncStatus = "pending"
	SyncStatusUnknown SyncStatus = "unknown"
)

// ExternalID represents an ID that identifies this media item in an external system
type SyncClient struct {
	// ID of the client that this external ID belongs to (optional for service IDs like TMDB)
	ID uint64 `json:"clientID,omitempty"`
	// Type of client this ID belongs to (optional for service IDs)
	Type types.ClientType `json:"clientType,omitempty" gorm:"type:varchar(50)"`
	// The actual ID value in the external system
	ItemID string `json:"itemID"`
	// The last time this item was synced
	LastSynced time.Time  `json:"lastSynced,omitempty"`
	SyncStatus SyncStatus `json:"syncStatus,omitempty"`
}

type SyncClients []*SyncClient

func (s *SyncClients) AddClient(clientID uint64, clientType types.ClientType, itemID string) {
	*s = append(*s, &SyncClient{
		ID:         clientID,
		Type:       clientType,
		ItemID:     itemID,
		LastSynced: time.Now(),
		SyncStatus: SyncStatusPending,
	})
}

func (s *SyncClients) GetSyncClients() []*SyncClient {
	if s == nil {
		return []*SyncClient{}
	}
	return *s
}

func (s *SyncClients) GetClientItemID(clientID uint64) string {
	if s == nil {
		return ""
	}
	for _, id := range *s {
		if id != nil && id.ID == clientID {
			return id.ItemID
		}
	}
	return ""
}

func (s *SyncClients) GetByClientID(clientID uint64) (*SyncClient, bool) {
	if s == nil {
		return &SyncClient{}, false
	}
	for _, client := range *s {
		if client != nil && client.ID == clientID {
			return client, true
		}
	}
	return &SyncClient{}, false
}

func (s *SyncClients) Merge(other SyncClients) {
	if s == nil {
		return
	}

	for _, otherClient := range other {
		if otherClient == nil {
			continue
		}

		found := false
		for i, client := range *s {
			if client != nil && otherClient != nil &&
			   client.ID == otherClient.ID &&
			   client.Type == otherClient.Type {
				// Update existing entry
				(*s)[i].ItemID = otherClient.ItemID
				found = true
				break
			}
		}
		if !found {
			// Add new entry
			*s = append(*s, otherClient)
		}
	}
}

func (s *SyncClients) Value() (driver.Value, error) {
	if s == nil || len(*s) == 0 {
		return "[]", nil
	}
	// Serialize the entire item to JSON for storage
	jsonData, err := json.Marshal(*s)
	if err != nil {
		return nil, err
	}
	return string(jsonData), nil
}

func (s *SyncClients) Scan(value any) error {
	if value == nil {
		*s = SyncClients{}
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}

	if len(data) > 0 && data[0] != '[' {
		data = append([]byte("["), append(data, ']')...)
	}

	return json.Unmarshal(data, s)
}

func (s *SyncClients) IsClientPresent(clientID uint64) bool {
	if s == nil {
		return false
	}
	for _, client := range *s {
		if client != nil && client.ID == clientID {
			return true
		}
	}
	return false
}

func (s *SyncClients) UpdateSyncStatus(clientID uint64, syncState SyncStatus) {
	if s == nil {
		return
	}
	for i, client := range *s {
		if client != nil && client.ID == clientID {
			(*s)[i].SyncStatus = syncState
			(*s)[i].LastSynced = time.Now()
			return
		}
	}
}

func (s *SyncClients) GetSyncStatus(clientID uint64) (SyncStatus, bool) {
	if s == nil {
		return SyncStatusUnknown, false
	}
	for _, client := range *s {
		if client != nil && client.ID == clientID {
			return client.SyncStatus, true
		}
	}
	// If not found, return pending
	return SyncStatusUnknown, false
}
