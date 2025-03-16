package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"suasor/services"
	"suasor/utils"
)

// RequireRole is a middleware that checks if the user has the required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := utils.LoggerFromContext(ctx)

		userRole, exists := c.Get("userRole")
		if !exists {
			log.Warn().Msg("User role not found in context")
			utils.RespondUnauthorized(c, nil, "Authentication required")
			c.Abort()
			return
		}

		role := userRole.(string)
		allowed := false
		for _, r := range roles {
			if r == role {
				allowed = true
				break
			}
		}

		if !allowed {
			log.Warn().Str("role", role).Msg("Insufficient permissions")
			utils.RespondForbidden(c, nil, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// VerifyToken is a middleware that verifies JWT tokens
func VerifyToken(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := utils.LoggerFromContext(ctx)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn().Msg("Authorization header missing")
			utils.RespondUnauthorized(c, nil, "Authorization required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn().Msg("Invalid authorization format")
			utils.RespondUnauthorized(c, nil, "Invalid authorization format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := authService.ValidateToken(ctx, token)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid or expired token")
			utils.RespondUnauthorized(c, err, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user claims in context
		c.Set("userId", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}
