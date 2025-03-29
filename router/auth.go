package router

import (
	"suasor/services"

	"github.com/gin-gonic/gin"

	"suasor/handlers"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, service services.AuthService) {
	authHandlers := handlers.NewAuthHandler(service)
	auths := rg.Group("/auth")
	{
		auths.POST("/register", authHandlers.Register)
		auths.POST("/login", authHandlers.Login)
		auths.POST("/refresh", authHandlers.RefreshToken)
		auths.POST("/logout", authHandlers.Logout)
		auths.GET("/validate", authHandlers.ValidateSession)
	}
}
