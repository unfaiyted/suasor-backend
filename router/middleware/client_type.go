// clientTypeMiddleware.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	clienttypes "suasor/clients/types"
)

// ClientTypeMiddleware extracts the client type from the database based on clientID
// and adds it to the request context
func ClientTypeMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.Param("clientID")
		if clientID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Client ID is required"})
			return
		}

		// Execute a raw query specifying the table name explicitly
		var clientType clienttypes.ClientType
		if err := db.Table("clients"). // Use your actual table name here
						Where("id = ?", clientID).
						Select("type").
						Scan(&clientType).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}

		// Add client type to context
		c.Set("clientType", clientType)
		c.Set("clientID", clientID)

		c.Next()
	}
}
