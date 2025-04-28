package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"suasor/clients/media/types"
	"suasor/types/models"
)

// MediaItemResponse is used for Swagger documentation to avoid generics
type MediaItemResponse struct {
	ID         uint64          `json:"id,omitempty"`
	Type       types.MediaType `json:"type"`
	ClientID   uint64          `json:"clientID"`
	ClientType string          `json:"clientType"`
	ExternalID string          `json:"externalID"`
	Data       any             `json:"data"`
	CreatedAt  string          `json:"createdAt,omitempty"`
	UpdatedAt  string          `json:"updatedAt,omitempty"`
}

// MediaItemListResponse is used for Swagger documentation to avoid generics
type MediaItemList[T types.MediaData] struct {
	Items []*models.MediaItem[T] `json:"items"`
	Total int                    `json:"total"`
}

type MediaDataList[T types.MediaData] struct {
	Items []*models.UserMediaItemData[T] `json:"items"`
	Total int                            `json:"total"`
}

// Convenience functions for success responses
func RespondMediaItemListOK[T types.MediaData](c *gin.Context, data []*models.MediaItem[T], message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	itemList := MediaItemList[T]{
		Items: data,
		Total: len(data),
	}

	RespondSuccess(c, http.StatusOK, itemList, msg)
}

func RespondMediaDataListOK[T types.MediaData](c *gin.Context, data []*models.UserMediaItemData[T], message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	itemList := MediaDataList[T]{
		Items: data,
		Total: len(data),
	}

	RespondSuccess(c, http.StatusOK, itemList, msg)
}
