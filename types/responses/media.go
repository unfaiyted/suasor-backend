package responses

import (
	"suasor/clients/media/types"
)

// MediaItemResponse is used for Swagger documentation to avoid generics
type MediaItemResponse struct {
	ID         uint64          `json:"id,omitempty"`
	Type       types.MediaType `json:"type"`
	ClientID   uint64          `json:"clientId"`
	ClientType string          `json:"clientType"`
	ExternalID string          `json:"externalId"`
	Data       interface{}     `json:"data"`
	CreatedAt  string          `json:"createdAt,omitempty"`
	UpdatedAt  string          `json:"updatedAt,omitempty"`
}
