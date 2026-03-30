package router

import (
	authhandler "sys-admin-serve/internal/handler/auth"
	categoryhandler "sys-admin-serve/internal/handler/category"
	healthhandler "sys-admin-serve/internal/handler/health"
	"sys-admin-serve/internal/middleware"
	"sys-admin-serve/internal/pkg/config"
	authrouter "sys-admin-serve/internal/router/auth"
	categoryrouter "sys-admin-serve/internal/router/category"
	healthrouter "sys-admin-serve/internal/router/health"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Dependencies struct {
	HealthHandler   *healthhandler.Handler
	AuthHandler     *authhandler.Handler
	CategoryHandler *categoryhandler.Handler
	AuthMiddleware  gin.HandlerFunc
}

func New(cfg *config.Config, log *zap.Logger, deps Dependencies) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger(log))

	v1 := engine.Group("/api/v1")
	healthrouter.RegisterRoutes(engine, v1, deps.HealthHandler)
	authrouter.RegisterRoutes(v1, deps.AuthHandler, deps.AuthMiddleware)
	categoryrouter.RegisterRoutes(v1, deps.CategoryHandler, deps.AuthMiddleware)

	return engine
}
