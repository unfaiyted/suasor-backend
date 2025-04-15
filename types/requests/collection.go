package requests

// ClientMediaRequest is used for testing a media client connection
type CollectionCreateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Type        string   `json:"type" binding:"required,oneof=movie series music mixed"`
	IsPublic    bool     `json:"isPublic"`
	ItemIDs     []uint64 `json:"itemIds" binding:"required"`
}

type CollectionUpdateRequest struct {
	CollectionID uint64   `json:"collectionId" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	Type         string   `json:"type" binding:"required,oneof=movie series music mixed"`
	IsPublic     bool     `json:"isPublic"`
	ItemIDs      []uint64 `json:"itemIds" binding:"required"`
}

type CollectionAddItemRequest struct {
	CollectionID uint64 `json:"collectionId" binding:"required"`
	ItemID       uint64 `json:"itemId" binding:"required"`
}

type CollectionRemoveItemRequest struct {
	CollectionID uint64 `json:"collectionId" binding:"required"`
	ItemID       uint64 `json:"itemId" binding:"required"`
}

type CollectionUpdateItemRequest struct {
	CollectionID uint64 `json:"collectionId" binding:"required"`
	ItemID       uint64 `json:"itemId" binding:"required"`
	Position     int    `json:"position" binding:"required,min=0"`
}
