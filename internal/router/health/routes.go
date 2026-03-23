package health

import (
	healthhandler "sys-admin-serve/internal/handler/health"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine, v1 *gin.RouterGroup, handler *healthhandler.Handler) {
	engine.GET("/health", handler.GetHealth)
	v1.GET("/health", handler.GetHealth)
}
