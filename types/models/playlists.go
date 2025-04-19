package models

type ListCollaborator struct {
	BaseModel
	ListID        uint64 `json:"listId"`
	UserID        uint64 `json:"userId"`
	Permission    string `json:"permission"`
	SharedAt      string `json:"sharedAt"`
	SharedBy      uint64 `json:"sharedBy"`
	LastSynced    string `json:"lastSynced"`
	SyncDirection string `json:"syncDirection"`
}

type ListSyncStatus struct {
	ListID       uint64                     `json:"playlistId"`
	LastSynced   string                     `json:"lastSynced"`
	ClientStates map[uint64]ClientSyncState `json:"clientStates"`
}

type ClientSyncState struct {
	ClientID     uint64        `json:"clientId"`
	ClientListID string        `json:"clientListId"`
	Items        SyncListItems `json:"items"`
	LastSynced   string        `json:"lastSynced"`
}

type SyncListItems []SyncListItem

type SyncListItem struct {
	ItemID        string         `json:"itemId"`
	Position      int            `json:"position"`
	LastChanged   string         `json:"lastChanged"`
	ChangeHistory []ChangeRecord `json:"changeHistory"`
}

type ChangeRecord struct {
	ClientID   uint64 `json:"clientId"` // 0 = internal client
	ItemID     string `json:"itemId,omitempty"`
	ChangeType string `json:"changeType"` // "add", "remove", "update", "reorder", "sync"
	Timestamp  string `json:"timestamp"`
}
