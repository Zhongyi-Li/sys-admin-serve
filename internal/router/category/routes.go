package category

import (
	categoryhandler "sys-admin-serve/internal/handler/category"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(v1 *gin.RouterGroup, handler *categoryhandler.Handler, authMiddleware gin.HandlerFunc) {
	categoryGroup := v1.Group("/categories", authMiddleware)
	{
		categoryGroup.GET("/tree", handler.Tree)
		categoryGroup.GET("", handler.List)
		categoryGroup.POST("", handler.Create)
		categoryGroup.PUT("/:id", handler.Update)
	}
}
