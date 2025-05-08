// repository/client.go
package repository

import (
	"context"
	"suasor/clients/types"
	"suasor/utils/logger"

	"gorm.io/gorm"
)

// ClientRepository defines the interface for media client database operations
type ClientHelper interface {
	GetClientTypeByClientID(ctx context.Context, clientID uint64) (types.ClientType, error)
}

type clientHelper struct {
	db *gorm.DB
}

// NewClientRepository creates a new media client repository
func NewClientHelper(db *gorm.DB) ClientHelper {
	return &clientHelper{db: db}
}

func (r *clientHelper) GetClientTypeByClientID(ctx context.Context, clientID uint64) (types.ClientType, error) {
	log := logger.LoggerFromContext(ctx)

	log.Debug().
		Uint64("clientID", clientID).
		Msg("Retrieving client type")

	var clientType types.ClientType
	if err := r.db.Table("clients"). // Use your actual table name here
						Where("id = ?", clientID).
						Select("type").
						Scan(&clientType).Error; err != nil {
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("clientType", clientType.String()).
		Msg("Retrieved client type")

	return clientType, nil
}
