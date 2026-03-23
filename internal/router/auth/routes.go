package auth

import (
	authhandler "sys-admin-serve/internal/handler/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(v1 *gin.RouterGroup, handler *authhandler.Handler, authMiddleware gin.HandlerFunc) {
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", handler.Login)
		authGroup.GET("/me", authMiddleware, handler.Me)
		authGroup.GET("/menus", authMiddleware, handler.Menus)
	}
}
