package models

import "time"

type ListItem struct {
	ItemUUID    string    `json:"itemUUID"`
	Position    int       `json:"position"`
	LastChanged time.Time `json:"lastChanged"`
}

type ListItems []ListItem

type ListCollaborator struct {
	BaseModel
	ListID        uint64                 `json:"listID"`
	UserID        uint64                 `json:"userID"`
	Permission    CollaboratorPermission `json:"permission"`
	SharedAt      string                 `json:"sharedAt"`
	SharedBy      uint64                 `json:"sharedBy"`
	LastSynced    string                 `json:"lastSynced"`
	SyncDirection string                 `json:"syncDirection"`
}

type CollaboratorPermission string

const (
	CollaboratorPermissionRead  CollaboratorPermission = "read"
	CollaboratorPermissionWrite CollaboratorPermission = "write"
)

type ListSyncStatus struct {
	ListID       uint64                     `json:"playlistID"`
	LastSynced   string                     `json:"lastSynced"`
	ClientStates map[uint64]ClientSyncState `json:"clientStates"`
}

type ClientSyncState struct {
	ClientID     uint64        `json:"clientID"`
	ClientListID string        `json:"clientListID"`
	Items        SyncListItems `json:"items"`
	LastSynced   string        `json:"lastSynced"`
}

type SyncListItems []SyncListItem

type SyncListItem struct {
	ItemID        string         `json:"itemID"`
	Position      int            `json:"position"`
	LastChanged   time.Time      `json:"lastChanged"`
	ChangeHistory []ChangeRecord `json:"changeHistory"`
}

type ChangeRecord struct {
	ClientID   uint64    `json:"clientID"` // 0 = internal client
	ItemID     string    `json:"itemID,omitempty"`
	ChangeType string    `json:"changeType"` // "add", "remove", "update", "reorder", "sync"
	Timestamp  time.Time `json:"timestamp"`
}
